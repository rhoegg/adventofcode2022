package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	data, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	parts := strings.Split(string(data), "\n\n")
	board := strings.Split(parts[0], "\n")
	path := parts[1]
	fmt.Printf("%s\n\n%s\n", strings.Join(board, "\n"), path)
	t := Turtle{}
	for x := range board[0] {
		if board[0][x] == '.' {
			t.Position.X = x
			break
		}
	}
	t.Follow(path, board)
	fmt.Printf("Turtle at %v facing %d\n", t.Position, t.Direction % 4)
	fmt.Printf("part 1 final password %d\n", t.FinalPassword())
}
