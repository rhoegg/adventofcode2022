package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

type Point struct {
	X, Y int
}

func (p Point) String() string {
	return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

func (p Point) Step(direction Direction) Point {
	switch direction {
	case North:
		return Point{X: p.X, Y: p.Y - 1}
	case South:
		return Point{X: p.X, Y: p.Y + 1}
	case West:
		return Point{X: p.X - 1, Y: p.Y}
	case East:
		return Point{X: p.X + 1, Y: p.Y}
	}
	log.Printf("WARN step didn't go anywhere %v %s", p, direction)
	return p
}

type Direction int
const (
	North Direction = iota
	South
	West
	East
)
func (d Direction) String() string {
	return [...]string{"N", "S", "W", "E"}[d]
}

type Crater struct {
	Elves map[Point]bool
	ElvesMoved int
	directionPriority [4]Direction
}
func NewCrater() Crater {
	return Crater{
		Elves: make(map[Point]bool),
		directionPriority: [4]Direction{North, South, West, East},
	}
}
func ParseCrater(filename string) Crater {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	crater := NewCrater()
	for y, line := range strings.Split(string(data), "\n") {
		for x := range line {
			if line[x] == '#' {
				crater.Elves[Point{X: x, Y: y}] = true
			}
		}
	}
	return crater
}
func (c Crater) TopLeft() Point {
	x := math.MaxInt; y := math.MaxInt
	for p := range c.Elves {
		if p.X < x {
			x = p.X
		}
		if p.Y < y {
			y = p.Y
		}
	}
	return Point{X: x, Y: y}
}

func (c Crater) BottomRight() Point {
	x := math.MinInt; y := math.MinInt
	for p := range c.Elves {
		if p.X > x {
			x = p.X
		}
		if p.Y > y {
			y = p.Y
		}
	}
	return Point{X: x, Y: y}
}

func (c Crater) String() string {
	var lines []string
	for y := c.TopLeft().Y; y <= c.BottomRight().Y; y++ {
		s := ""
		for x := c.TopLeft().X; x <= c.BottomRight().X; x++ {
			if c.Elves[Point{X: x, Y: y}] {
				s += "#"
			} else {
				s += "."
			}
		}
		lines = append(lines, s)
	}
	return strings.Join(lines, "\n")
}

func (c Crater) EnoughSpace(elf Point) bool {
	for y := elf.Y - 1; y <= elf.Y + 1; y++ {
		for x := elf.X - 1; x <= elf.X + 1; x++ {
			if x == elf.X && y == elf.Y {
				continue // one is the loneliest number...
			}
			if c.Elves[Point{X: x, Y: y}] {
				return false // another elf is nearby
			}
		}
	}
	return true
}

func (c Crater) SomeRoom(elf Point, direction Direction) bool {
	var adjacents []Point
	switch direction {
	case North:
		for _, i := range []int{-1, 0, 1} {
			adjacents = append(adjacents, Point{X: elf.X + i, Y: elf.Y - 1})
		}
		break
	case South:
		for _, i := range []int{-1, 0, 1} {
			adjacents = append(adjacents, Point{X: elf.X + i, Y: elf.Y + 1})
		}
		break
	case East:
		for _, i := range []int{-1, 0, 1} {
			adjacents = append(adjacents, Point{X: elf.X + 1, Y: elf.Y + i})
		}
		break
	case West:
		for _, i := range []int{-1, 0, 1} {
			adjacents = append(adjacents, Point{X: elf.X - 1, Y: elf.Y + i})
		}
		break
	}
	for _, adjacent := range adjacents {
		if c.Elves[adjacent] {
			return false
		}
	}
	return true
}

func (c Crater) NextRound() Crater {
	intentions := make(map[Point]Point)
	popularity := make(map[Point]int)
	next := NewCrater()
	for elf := range c.Elves {
		intentions[elf] = elf
		if ! c.EnoughSpace(elf) {
			next.ElvesMoved++
			for _, direction := range c.directionPriority {
				if c.SomeRoom(elf, direction) {
					intentions[elf] = elf.Step(direction)
					break
				}
			}
		}
		popularity[intentions[elf]] += 1
	}
	for i := 0; i < 4; i++ {
		next.directionPriority[i] = c.directionPriority[(i + 1) % 4]
	}
	for elf := range intentions {
		if popularity[intentions[elf]] == 1 {
			next.Elves[intentions[elf]] = true
		} else {
			next.Elves[elf] = true
		}
	}
	if len(next.Elves) != len(c.Elves) {
		panic(fmt.Sprintf("elves don't balance %d = %d", len(next.Elves), len(c.Elves)))
	}
	return next
}

func (c Crater) EmptyInRect() int {
	var emptyCount int
	for y := c.TopLeft().Y; y <= c.BottomRight().Y; y++ {
		for x := c.TopLeft().X; x <= c.BottomRight().X; x++ {
			if ! c.Elves[Point{X: x, Y: y}] {
				emptyCount++
			}
		}
	}
	return emptyCount
}