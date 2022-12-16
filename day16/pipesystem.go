package main

import (
	"container/heap"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PipeSystem struct {
	Valves map[string]Valve
}

type Valve struct {
	Name string
	FlowRate int
	Tunnels []string
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

func (ps *PipeSystem) FlowTargets() (targets []string) {
	for _, v := range ps.Valves {
		if v.FlowRate > 0 {
			targets = append(targets, v.Name)
		}
	}
	return targets
}

func (ps *PipeSystem) Path(origin, destination string) (path []string) {
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
			return append(node.pathHere, node.Name)
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