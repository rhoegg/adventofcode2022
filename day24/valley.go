package main

import (
	"os"
	"strings"
)

type Direction uint8

const (
	North Direction = iota
	East
	South
	West
)

func (d Direction) String() string {
	return [...]string{"^", ">", "v", "<"}[d]
}

type Point struct {
	X, Y int
}
func (p Point) Step(d Direction) Point {
	switch d {
	case North:
		return Point{X: p.X, Y: p.Y - 1}
	case East:
		return Point{X: p.X + 1, Y: p.Y}
	case South:
		return Point{X: p.X, Y: p.Y + 1}
	case West:
		return Point{X: p.X - 1, Y: p.Y}
	default: return p
	}
}
func (p Point) ManhattanDistance(p2 Point) int {
	distance := p2.X - p.X
	if distance < 0 {
		distance *= -1
	}
	if p2.Y > p.Y {
		distance += p2.Y - p.Y
	} else {
		distance += p.Y - p2.Y
	}
	return distance
}

type Blizzard struct {
	Position Point
	Direction Direction
}

type Valley struct {
	Entrance, Exit Point
	Height, Width int
	Walls map[Point]bool
	Blizzards []Blizzard
}

func (v Valley) Advance() Valley {
	var nextBlizzards []Blizzard
	for _, b := range v.Blizzards {
		next := b.Position.Step(b.Direction)
		if _, ok := v.Walls[next]; ok {
			switch b.Direction {
			case North:
				next.Y = v.Height - 2
				break
			case East:
				next.X = 1
				break
			case South:
				next.Y = 1
				break
			case West:
				next.X = v.Width - 2
				break
			}
		}
		nextBlizzards = append(nextBlizzards, Blizzard{Position: next, Direction: b.Direction})
	}
	v.Blizzards = nextBlizzards
	return v
}

func NewValley(height, width int) Valley {
	return Valley{
		Height: height,
		Width: width,
		Walls: make(map[Point]bool),
	}
}

func ParseValley(filename string) Valley {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(data), "\n")

	valley := NewValley(len(lines), len(lines[0]))

	for i, line := range lines {
		switch i {
		case 0:
			for j := range line {
				if line[j] == '.' {
					valley.Entrance = Point{X: j, Y: 0}
				} else {
					valley.Walls[Point{X: j, Y: 0}] = true
				}
			}
			break
		case len(lines) - 1:
			for j := range line {
				if line[j] == '.' {
					valley.Exit = Point{X: j, Y: len(lines) - 1}
				} else {
					valley.Walls[Point{X: j, Y: len(lines) - 1}] = true
				}
			}
			break
		default:
			for j, c := range line {
				if c == '#' {
					valley.Walls[Point{X: j, Y: i}] = true
				}
				if c != '.' {
					var dir Direction
					switch c {
					case '^':
						dir = North
						break
					case '>':
						dir = East
						break
					case 'v':
						dir = South
						break
					case '<':
						dir = West
						break
					}
					valley.Blizzards = append(valley.Blizzards, Blizzard{
						Position: Point{X: j, Y: i},
						Direction: dir,
					})
				}
			}
		}
	}
	return valley
}