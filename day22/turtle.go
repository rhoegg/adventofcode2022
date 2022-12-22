package main

import (
	"log"
	"strconv"
)

type Turtle struct {
	Position Position
	Direction int  // 0 east, 1 south, 2 west, 3 north
}

func (t Turtle) FinalPassword() int {
	row := t.Position.Y + 1
	col := t.Position.X + 1
	facing := t.Direction
	return row * 1000 + col * 4 + facing
}

func (t *Turtle) Follow(path string, board []string) {
	if len(path) == 0 {
		return
	}
	if path[0] == 'R' {
		log.Print("Right")
		t.Direction += 1
		if t.Direction > 3 {
			t.Direction -= 4
		}
		path = path[1:]
	} else if path[0] == 'L' {
		log.Print("Left")
		t.Direction -= 1
		if t.Direction < 0 {
			t.Direction += 4
		}
		path = path[1:]
	} else {
		var steps int
		for len(path) > 0 {
			nextDigit := path[0:1] // peek
			if nextDigit >= "0" && nextDigit <= "9" {
				steps *= 10
				val, _ := strconv.Atoi(nextDigit)
				steps += val
				path = path[1:]
			} else {
				break
			}
		}
		log.Printf("Forward %d", steps)
		t.Forward(steps, board)
		log.Printf("arrived %v", t.Position)
	}
	t.Follow(path, board)
}

func (t *Turtle) Forward(steps int, board []string) {
	if steps <= 0 {
		return
	}
	var nextPosition Position
	switch t.Direction {
	case 0:
		nextPosition = t.Position.East(board)
		break
	case 1:
		nextPosition = t.Position.South(board)
		break
	case 2:
		nextPosition = t.Position.West(board)
		break
	case 3:
		nextPosition = t.Position.North(board)
		break
	}
	if board[nextPosition.Y][nextPosition.X] == '#' {
		return
	}
	t.Position = nextPosition
	t.Forward(steps - 1, board)
}
