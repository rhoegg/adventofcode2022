package main

import (
	"container/heap"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type PipeSystem struct {
	Valves map[string]Valve
	pathCache map[string][]string
	checks10min int
	cacheHits int64
}

type Valve struct {
	Name string
	FlowRate int
	Tunnels []string
}

type pathPressure struct {
	pressure int
	paths [2][]string
}

func NewPipeSystem(filename string) *PipeSystem {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	valves := make(map[string]Valve)
	for _, line := range strings.Split(string(data), "\n") {
		v :=  parseValve(line)
		valves[v.Name] = v
	}

	return &PipeSystem{
		Valves: valves,
		pathCache: make(map[string][]string),
	}
}

func (ps *PipeSystem) MostPressure(origin string, visited map[string]bool, minutes int) (waypoints []string, pressure int) {
	if minutes <= 0 {
		panic("took too long")
	}
	originValve := ps.Valves[origin]
	newVisited := make(map[string]bool)
	for k, _ := range visited {
		newVisited[k] = true
	}
	newVisited[origin] = true

	if originValve.FlowRate > 0 {
		// one minute to release flow
		minutes -= 1
	}

	if len(ps.FlowTargets()) == len(visited) {
		//fmt.Printf("last stop, %d pressure * %d minutes = %d\n", originValve.FlowRate, minutes - 1, originValve.FlowRate * (minutes - 1))
		return []string{origin}, originValve.FlowRate * minutes
	}

	var paths [][]string
	for _, target := range ps.FlowTargets() {
		if ! visited[target] && target != origin {
			p := ps.Path(origin, target)
			paths = append(paths, p)
		}
	}

	var bestWaypoints []string
	var maxTunnelPressure int
	for _, p := range paths {
		if minutes - len(p) > 0 {
			nextWaypoints, tunnelPressure := ps.MostPressure(p[len(p) - 1], newVisited, minutes - len(p))
			if tunnelPressure > maxTunnelPressure {
				maxTunnelPressure = tunnelPressure
				bestWaypoints = append([]string{origin}, nextWaypoints...)
			}
		}
	}
	return bestWaypoints, (originValve.FlowRate * minutes) + maxTunnelPressure
}

func (ps *PipeSystem) MostPressureWithElephant(plans [2][]string, visited map[string]bool, minutes int) pathPressure {
	if minutes < 0 {
		panic("took too long")
	}
	if minutes == 0 {
		return pathPressure{paths: plans, pressure: 0}
	}

	var tryAllPaths bool = true
	var newFlowRate int
	var possiblePlans [2][][]string
	var maxTunnelPressure pathPressure
	opens := [2]string{" ", " "}
	newVisited := make(map[string]bool)
	for k, _ := range visited {
		newVisited[k] = true
	}
	for _, i := range []int{0, 1} {
		if len(plans[i]) > 1 {
			possiblePlans[i] = [][]string{plans[i][1:]}
		} else {
			origin := plans[i][0]
			originValve := ps.Valves[origin]
			if newVisited[origin] || originValve.FlowRate == 0 {
				// make new plans
				goodTargets := make(map[string]bool)
				for _, target := range ps.FlowTargets() {
					if ! newVisited[target] && target != origin {
						goodTargets[target] = true
					}
				}
				for target := range goodTargets {
					p := ps.Path(origin, target)
					var skip bool
					for _, step := range p[:len(p) - 1] {
						if goodTargets[step] && ps.Valves[step].FlowRate >= ps.Valves[target].FlowRate {
							skip = true
						}
					}
					if len(p) < minutes && (tryAllPaths || !skip) {
						possiblePlans[i] = append(possiblePlans[i], p)
					}
				}
			} else {
				// open valve and mark visited
				newFlowRate += originValve.FlowRate
				newVisited[origin] = true
				opens[i] = "*"
				// repeat the same plan, it won't be visited next time around
				possiblePlans[i] = [][]string{plans[i]}
			}
		}

		for _, myPlan := range possiblePlans[0] {
			for _, elephantPlan := range possiblePlans[1] {
				candidatePlans := [2][]string{myPlan, elephantPlan}

				tunnelPressure := ps.MostPressureWithElephant(candidatePlans, newVisited, minutes - 1)
				if tunnelPressure.pressure > maxTunnelPressure.pressure {
					maxTunnelPressure = tunnelPressure
				}
			}
		}
	}
	for _, i := range []int{0,1} {
		maxTunnelPressure.paths[i] = append([]string{plans[i][0] + opens[i]}, maxTunnelPressure.paths[i]...)
	}
	if minutes == 22 {
		log.Printf("(10s/hits/size %d/%d/%d) Max pressure w/ 22 minutes left (%d): \n%v\n%v",
			ps.checks10min, ps.cacheHits, len(ps.pathCache),
			maxTunnelPressure.pressure, maxTunnelPressure.paths[0], maxTunnelPressure.paths[1])
	}
	if minutes == 10 {
		ps.checks10min++
	}

	result := pathPressure{
		paths: maxTunnelPressure.paths,
		pressure: (newFlowRate * (minutes - 1)) + maxTunnelPressure.pressure,
	}
	return result
}


func (ps *PipeSystem) FlowTargets() (targets []string) {
	for _, v := range ps.Valves {
		if v.FlowRate > 0 {
			targets = append(targets, v.Name)
		}
	}
	return targets
}

func (ps *PipeSystem) Path(origin, destination string) []string {
	if path, ok := ps.pathCache[origin + "-" + destination]; ok {
		ps.cacheHits++
		return path
	}
	vh := make(valveHeap, len(ps.Valves[origin].Tunnels))
	for i, name := range ps.Valves[origin].Tunnels {
		vh[i] = valveNode{Valve: ps.Valves[name], distance: 1}
	}
	heap.Init(&vh)
	visited := make(map[string]struct{})
	for len(vh) > 0 {
		node := heap.Pop(&vh).(valveNode)
		if node.Name == destination {
			// got there
			path := append(node.pathHere, node.Name)
			ps.pathCache[origin + "-" + destination] = path
			return path
		}
		if _, ok := visited[node.Name]; ok {
			continue
		}

		for _, t := range node.Tunnels {
			if _, ok := visited[t]; !ok {
				newnode := valveNode{
					Valve: ps.Valves[t],
					distance: node.distance + 1,
					pathHere: append(node.pathHere, node.Name),
				}
				heap.Push(&vh, newnode)
			}
		}
		visited[node.Name] = struct{}{}
	}

	fmt.Printf("No path from %s to %s\n", origin, destination)

	return nil
}

func parseValve(text string) Valve {
	clauses := strings.Split(text, "; ")
	tokens := strings.Split(clauses[0], " ")
	flowRate, _ := strconv.Atoi(strings.TrimPrefix(tokens[4], "rate="))
	tunnelText := strings.TrimPrefix(clauses[1], "tunnels lead to valves ")
	tunnelText = strings.TrimPrefix(tunnelText, "tunnel leads to valve ")
	return Valve{Name: tokens[1], FlowRate: flowRate, Tunnels: strings.Split(tunnelText, ", ")}
}

type valveNode struct {
	Valve
	distance int
	pathHere []string
}
type valveHeap []valveNode

func (h valveHeap) Len() int           { return len(h) }
func (h valveHeap) Less(i, j int) bool { return h[i].distance < h[j].distance }
func (h valveHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *valveHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(valveNode))
}

func (h *valveHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}