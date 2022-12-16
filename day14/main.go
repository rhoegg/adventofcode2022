package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)
var ORIGIN = Point{X: 500, Y: 0}
var maxDepth int

func main() {
	cave := make(map[int][]int)
	rockSegments := parseRockSegments("input.txt")
	for _, segment := range rockSegments {
		cave = drawRocks(cave, segment)
	}

	maxDepth = 0
	for _, column := range cave {
		max := column[len(column) - 1] // column is sorted
		if max > maxDepth {
			maxDepth = max
		}
	}
	// max is 2 deeper per challenge text
	maxDepth += 2

	sand := make(map[Point]struct{})
	paintCave(cave)
	caveMap := make(map[Point]struct{})
	for x := range cave {
		for _, y := range cave[x] {
			caveMap[Point{x, y}] = struct{}{}
		}
	}
	pourSand(caveMap, sand)
	paintSand(cave, sand)
	fmt.Printf("Total sand: %d\n", len(sand))
}

type Point struct {
	X, Y int
}

func parseRockSegments(filename string) [][]Point {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var rockSegments [][]Point
	for _, line := range strings.Split(string(data), "\n") {
		var segment []Point
		for _, point := range strings.Split(line, " -> ") {
			values := strings.Split(point, ",")
			x, _ := strconv.Atoi(values[0])
			y, _ := strconv.Atoi(values[1])
			segment = append(segment, Point{X: x, Y: y})
		}
		rockSegments = append(rockSegments, segment)
	}
	return rockSegments
}

func drawRocks(cave map[int][]int, segment []Point) map[int][]int {
	for i := 1; i < len(segment); i++ {
		p1, p2 := segment[i - 1], segment[i]
		if p1.X == p2.X {
			if p1.Y > p2.Y {
				p1, p2 = p2, p1
			}
			for y := p1.Y; y <= p2.Y; y++ {
				cave[p1.X] = append(cave[p1.X], y)
			}
		} else { // assume p1.Y == p2.Y
			if p1.X > p2.X {
				p1, p2 = p2, p1
			}
			for x := p1.X; x <= p2.X; x++ {
				cave[x] = append(cave[x], p1.Y)
			}
		}
	}
	// distinct coordinates
	for x, column := range cave {
		unique := make(map[int]struct{})
		for _, y := range column {
			unique[y] = struct{}{}
		}
		column = make([]int, len(unique))
		i := 0
		for k := range unique {
			column[i] = k
			i++
		}
		sort.Ints(column)
		fmt.Printf("%d: %v\n", x, column)
	}
	return cave
}

func paintCave(cave map[int][]int) {
	paintSand(cave, make(map[Point]struct{}))
}
func paintSand(cave map[int][]int, sand map[Point]struct{}) {
	x1, x2 := 500, 500
	for x, _ := range cave {
		if x < x1 {
			x1 = x
		}
		if x > x2 {
			x2 = x
		}
	}
	x1, x2 = x1 - 30, x2 + 30
	for y := 0; y < maxDepth; y++ {
		for x := x1; x <= x2; x++ {
			c := " "
			for _, r := range cave[x] {
				if r == y {
					c = "#"
				}
			}
			if _, ok := sand[Point{X: x, Y: y}]; ok {
				c = "o"
			}
			print(c)
		}
		println()
	}
}

func isClear(caveMap, sand map[Point]struct{}, test Point) bool {
	// check rocks
	if _, blocked := caveMap[test]; blocked {
		return false
	}
	// check sand
	if _, blocked := sand[test]; blocked {
		return false
	}
	return true
}

func dropSand(caveMap map[Point]struct{}, sand map[Point]struct{}, last Point) Point {
	if last.Y == maxDepth - 1 {
		return last
	}
	next := last
	down := Point{X: last.X, Y: last.Y + 1}
	left := Point{X: last.X - 1, Y: last.Y + 1}
	right := Point{X: last.X + 1, Y: last.Y + 1}
	for _, p := range []Point{down, left, right} {
		if isClear(caveMap, sand, p) {
			next = p
		}
	}
	if next == last {
		return last
	}
	return dropSand(caveMap, sand, next)
}

func pourSand(caveMap, sand map[Point]struct{}) {
	for {
		p := dropSand(caveMap, sand, ORIGIN)
		sand[p] = struct{}{}
		if p == ORIGIN {
			break
		}
	}
}