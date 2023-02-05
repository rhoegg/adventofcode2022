package main

type Point struct {
	X, Y int
}

func (p Point) Translate(p2 Point) Point {
	return Point{X: p.X + p2.X, Y: p.Y + p2.Y}
}

type Rock struct {
	Shape map[Point]bool
	Number int
}

func NewRock(rockNumber int) Rock {
	return Rock{Shape: make(map[Point]bool), Number: rockNumber}
}


