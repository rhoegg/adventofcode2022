package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Move struct {
	Direction string
	Distance int
}

type Location struct {
	X, Y int
}

func (m Move) MoveHead(x, y int) (x2, y2 int) {
	if m.Direction == "U" {
		return Up(x, y)
	}
	if m.Direction == "D" {
		return Down(x, y)
	}
	if m.Direction == "L" {
		return Left(x, y)
	}
	if m.Direction == "R" {
		return Right(x, y)
	}
	panic("unknown direction " + m.Direction)
}

func MoveTail(x, y int, hx, hy int) (x2, y2 int) {
	x2, y2 = x, y
	relX, relY := hx - x, hy - y
	if (relX != 0) && (relY != 0) {
		// diagonal
		if relX * relX * relY * relY > 1 { // x^2 * y^2
			if relX < 0 { relX = -1} else { relX = 1 }
			if relY < 0 { relY = -1} else { relY = 1 }
			x2 = x2 + relX
			y2 = y2 + relY
		}
		return x2, y2
	}

	if relX > 1 {
		x2 = x + 1
	} else if relX < -1 {
		x2 = x - 1
	}
	if relY > 1 {
		y2 = y + 1
	} else if relY < -1 {
		y2 = y - 1
	}
	return x2, y2
}

func Up(x, y int) (x2, y2 int) {
	return x, y + 1
}

func Down(x, y int) (x2, y2 int) {
	return x, y - 1
}

func Left(x, y int) (x2, y2 int) {
	return x - 1, y
}

func Right(x, y int) (x2, y2 int) {
	return x + 1, y
}

func main() {
	data, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	var moves []Move
	for _, line := range strings.Split(string(data), "\n") {
		tokens := strings.Split(line, " ")
		distance, _ := strconv.Atoi(tokens[1])
		moves = append(moves, Move{Direction: tokens[0], Distance: distance})
	}

	hx, hy := 0, 0
	tx, ty := hx, hy
	var knots []Location
	for i := 0; i < 10; i++ {
		knots = append(knots, Location{X: hx, Y: hy})
	}

	tailLocations := make(map[Location]struct{})

	for _, move := range moves {
		for i := 0; i < move.Distance; i++ {
			hx, hy = move.MoveHead(hx, hy)
			tx, ty = MoveTail(tx, ty, hx, hy)
			tailLocations[Location{X: tx, Y: ty}] = struct{}{}
		}
	}
	fmt.Printf("Part 1 Tail locations: %d\n", len(tailLocations))
	tailLocations = make(map[Location]struct{})

	for _, move := range moves {
		for i := 0; i < move.Distance; i++ {
			knots[0].X, knots[0].Y = move.MoveHead(knots[0].X, knots[0].Y)
			for j := 1; j < len(knots); j++ {
				knots[j].X, knots[j].Y = MoveTail(knots[j].X, knots[j].Y, knots[j-1].X, knots[j-1].Y)
			}
			tailLocations[knots[len(knots) - 1]] = struct{}{}
		}
	}
	fmt.Printf("Part 2: %d\n", len(tailLocations))
}
