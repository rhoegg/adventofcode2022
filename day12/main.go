package main

import (
	"container/heap"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
)
const START = 'S'
const END = 'E'


type MapPoint struct {
	Location
	Elevation rune
}

func reachable(m [][]MapPoint, loc Location) (result []Location) {
	x, y := loc.X, loc.Y
	baseElevation := m[y][x].Elevation
	for _, x2 := range []int{x - 1, x + 1} {
		if x2 >= 0 && x2 < len(m[y]) {
			if m[y][x2].Elevation <= (baseElevation + 1) {
				result = append(result, Location{X: x2, Y: y})
			}
		}
	}
	for _, y2 := range []int{y - 1, y + 1} {
		if y2 >= 0 && y2 < len(m) {
			if m[y2][x].Elevation <= (baseElevation + 1) {
				result = append(result, Location{X: x, Y: y2})
			}
		}
	}
	return result
}

func MinSteps(terrain [][]MapPoint, start, end Location) int {
	pq := LocationPQ{}
	for _, row := range terrain {
		for _, p := range row {
			if start == p.Location {
				heap.Push(&pq, &WeightedLocation{Location: p.Location, Weight: 0})
			}
		}
	}
	heap.Init(&pq)

	visited := make(map[Location]int)

	//fmt.Printf("Start %v ; End %v\n", start, end)

	var next *WeightedLocation
	steps := 0
	for {
		if len(pq) == 0 {
			return math.MaxInt
		}
		next = heap.Pop(&pq).(*WeightedLocation)
		if _, ok := visited[next.Location]; ok {
			continue
		}
		steps++
		//if steps % 1000 == 0 {
		//	fmt.Printf("Step %d, location %v\n", steps, next.Location)
		//}
		visited[next.Location] = next.Weight
		if next.Location == end {
			break
		}
		for _, loc := range reachable(terrain, next.Location) {
			heap.Push(&pq, &WeightedLocation{Location: loc, Weight: next.Weight + 1})
		}
	}
	return next.Weight
}

func main() {
	data, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}


	var terrain [][]MapPoint
	var start, end Location
	var starts []Location
	for y, line := range strings.Split(string(data), "\n") {
		var row []MapPoint
		for x, e := range line {
			p := MapPoint{Location: Location{X: x, Y: y}, Elevation: e}
			if START == e  {
				start = p.Location
				p.Elevation = 'a'
			}
			if END == e {
				end = p.Location
				p.Elevation = 'z'
			}
			row = append(row, p)
			if p.Elevation == 'a' {
				starts = append(starts, p.Location)
			}
		}
		terrain = append(terrain, row)
	}

	fmt.Printf("Step 1: min steps for %v: %d\n", start, MinSteps(terrain, start, end))

	var startSteps []WeightedLocation
	for _, s := range starts {
		steps := MinSteps(terrain, s, end)
		if steps < 100000 {
			startSteps = append(startSteps, WeightedLocation{Location:s, Weight: steps})
		}
	}
	sort.Slice(startSteps, func(i, j int) bool {
		return startSteps[i].Weight < startSteps[j].Weight
	})
	fmt.Printf("Step 2: Closest start %v: %d\n", startSteps[0], startSteps[0].Weight)
}
