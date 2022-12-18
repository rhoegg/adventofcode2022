package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	X, Y, Z int
}

func (p Point) Adjacents() []Point {
	return []Point{
		{X: p.X-1, Y: p.Y, Z: p.Z},
		{X: p.X+1, Y: p.Y, Z: p.Z},
		{X: p.X, Y: p.Y-1, Z: p.Z},
		{X: p.X, Y: p.Y+1, Z: p.Z},
		{X: p.X, Y: p.Y, Z: p.Z-1},
		{X: p.X, Y: p.Y, Z: p.Z+1},
	}
}

func main() {
	lava := parseLava("input.txt")
	var min, max Point
	for p := range lava {
		if p.X < min.X { min.X = p.X }
		if p.X > max.X { max.X = p.X }
		if p.Y < min.Y { min.Y = p.Y }
		if p.Y > max.Y { max.Y = p.Y }
		if p.Z < min.Z { min.Z = p.Z }
		if p.Z > max.Z { max.Z = p.Z }
	}
	steam := make(map[Point]bool)
	floodFillOutside(lava, expandBounds([2]Point{min, max}), min, steam)
	outsideFaces := 0
	allFaces := 0
	for lavaPoint := range lava {
		for _, neighbor := range lavaPoint.Adjacents() {
			if _, ok := lava[neighbor]; !ok {
				allFaces++
			}
			if _, ok := steam[neighbor]; ok {
				outsideFaces++
			}
		}
		//log.Printf("%v %d %d\n", lavaPoint, outsideFaces - startFaces, allFaces - startAllFaces)
	}
	fmt.Printf("steam/all: %d/%d\n", outsideFaces, allFaces)
}

func parseLava(filename string) map[Point]bool {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	lavaPoints := make(map[Point]bool)
	for _, line := range strings.Split(string(data), "\n") {
		tokens := strings.Split(line, ",")
		x, _ := strconv.Atoi(tokens[0])
		y, _ := strconv.Atoi(tokens[1])
		z, _ := strconv.Atoi(tokens[2])
		lavaPoints[Point{X: x, Y: y, Z: z}] = true
	}
	return lavaPoints
}

func expandBounds(bounds[2]Point) [2]Point {
	return [2]Point{
		{X: bounds[0].X - 1, Y: bounds[0].Y - 1, Z: bounds[0].Z - 1},
		{X: bounds[1].X + 1, Y: bounds[1].Y + 1, Z: bounds[1].Z + 1},
	}
}

func inBounds(bounds [2]Point, p Point) bool {
	return p.X >= bounds[0].X && p.X <= bounds[1].X &&
		p.Y >= bounds[0].Y && p.Y <= bounds[1].Y &&
		p.Z >= bounds[0].Z && p.Z <= bounds[1].Z
}

func floodFillOutside(lava map[Point]bool, bounds [2]Point, p Point, prevFilled map[Point]bool) {
	if !inBounds(bounds, p) {
		return
	}
	if _, ok := lava[p]; ok {
		return
	}
	if _, ok := prevFilled[p]; ok {
		return
	}
	prevFilled[p] = true
	for _, neighbor := range p.Adjacents() {
		floodFillOutside(lava, bounds, neighbor, prevFilled)
	}
}