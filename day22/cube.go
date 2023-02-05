package main

type CubeFace struct {
	Board [50][50]bool
	North, South, East, West CubeFaceEdge
}

type CubeFaceEdge struct {
	*CubeFace
	Edge int // 0 = East, 1 = South, 2 = West, 3 = North
}

func ParseCubeFaces(board []string) (*CubeFace, map[Position]*CubeFace) {
	faces := make(map[Position]*CubeFace)

	var startFace *CubeFace
	var startX int
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			if len(board[50 * y]) > (50 * x) && board[50 * y][50 * x] != ' ' {
				cubeFace := &CubeFace{}
				if startFace == nil {
					startFace = cubeFace
					startX = x
				}
				for by := 0; by < 50; by++ {
					for bx := 0; bx < 50; bx++ {
						cubeFace.Board[by][bx] = board[50*y + by][50*x + bx] == '#'
					}
				}
				faces[Position{X: x, Y: y}] = cubeFace
			}
		}
	}
	for x := startX + 1; x < 3; x++ {
		// y is 0
		face := startFace
		if cubeFace, ok := faces[Position{X: x, Y: 0}]; ok {
			face.East = cubeFace
			cubeFace.West = face
			face = cubeFace
		}
	}
	for y := 1; y < 3; y++ {
		for x := 0; x < 3; x++ {
			if cubeFace, ok := faces[Position{X: x, Y: y}]; ok {
				if face, ok := faces[Position{X: x - 1, Y: y}]; ok {
					face.East = cubeFace
					cubeFace.West = face
				}
				if face, ok := faces[Position{X: x, Y: y - 1}]; ok {
					face.South = cubeFace
					cubeFace.North = face
				}
			}
		}
	}
	return startFace, faces
}