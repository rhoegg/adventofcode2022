package main

type Position struct {
	X, Y int
}

func (p Position) East(board []string) Position {
	x := p.X
	for {
		x = x + 1
		if x >= len(board[p.Y]) {
			x = 0
		}
		if board[p.Y][x] != ' ' {
			return Position{X: x, Y: p.Y}
		}
	}
}

func (p Position) West(board []string) Position {
	x := p.X
	for {
		x = x - 1
		if x < 0 {
			x += len(board[p.Y])
		}
		if board[p.Y][x] != ' ' {
			return Position{X: x, Y: p.Y}
		}
	}
}

func (p Position) North(board []string) Position {
	y := p.Y
	for {
		y = y - 1
		if y < 0 {
			y += len(board)
		}
		if (len(board[y]) - 1) >= p.X && board[y][p.X % len(board[y])] != ' ' {
			return Position{X: p.X, Y: y}
		}
	}
}

func (p Position) South(board []string) Position {
	y := p.Y
	for {
		y = y + 1
		if y >= len(board) {
			y = 0
		}
		if (len(board[y]) - 1) >= p.X && board[y][p.X % len(board[y])] != ' ' {
			return Position{X: p.X, Y: y}
		}
	}
}
