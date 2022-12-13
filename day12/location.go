package main

import "fmt"

type Location struct {
	X, Y int
}

func (l Location) String() string {
	return fmt.Sprintf("(%d, %d)", l.X, l.Y)
}
