package main

import (
	"container/heap"
	"fmt"
	"log"
)

func main() {
	valley := ParseValley("input.txt")
	var possibilities minExpedition = []Expedition{{Position: valley.Entrance, Valley: valley, Moves: 0}}
	heap.Init(&possibilities)

	visited := make(map[string]bool)
	var win Expedition
	var moves int
	var leg int
	var leg1moves int
	var leg2moves int
	for {
		e := heap.Pop(&possibilities).(Expedition)
		if visited[e.Key()] {
			continue
		}
		if leg == 1 && e.Leg == 0 && e.Moves - leg1moves > 12 {
			continue
		}
		if leg == 2 && e.Leg < 2 && e.Moves - leg2moves > 12 {
			continue
		}
		visited[e.Key()] = true
		if e.Position == valley.Exit && e.Leg == 2 {
			win = e
			break
		} else if e.Position == valley.Entrance && e.Leg == 1 {
			e.Leg = 2
			leg2moves = e.Moves
		} else if e.Position == valley.Exit && e.Leg == 0 {
			e.Leg = 1
			leg1moves = e.Moves
		}
		if e.Leg > leg {
			leg = e.Leg
		}
		if e.Moves > moves && e.Leg == leg {
			moves = e.Moves
			log.Printf("%d (%d) checking %v", e.Leg, e.Moves, e.Position)
		}
		nextMoves := e.NextMoves()
		for _, move := range nextMoves {
			if !visited[move.Key()] {
				move.Leg = e.Leg
				heap.Push(&possibilities, move)
			}
		}
	}
	fmt.Printf("Part 1: %d\n", win.Moves)
	//e := win
	//var journey []Expedition
	//for e.Last != nil {
	//	journey = append([]Expedition{e}, journey...)
	//	e = *e.Last
	//}
	//for _, e := range journey {
	//	printExpedition(e)
	//}
	//fmt.Printf("journey is %d\n", len(journey))
}

type Expedition struct {
	Position Point
	Valley Valley
	Moves int
	Last *Expedition
	Leg int
}

func (e Expedition) NextMoves() []Expedition {
	nextValley := e.Valley.Advance()
	var moves []Point
	blizzards := make(map[Point]bool)
	for _, b := range nextValley.Blizzards {
		blizzards[b.Position] = true
	}
	if _, ok := blizzards[e.Position]; !ok {
		moves = append(moves, e.Position)
	}
	for _, dir := range []Direction{North, East, South, West} {
		nextPos := e.Position.Step(dir)
		if nextPos == e.Valley.Exit || nextPos == e.Valley.Entrance {
			moves = append(moves, nextPos)
		}
		if nextPos.X > 0 && nextPos.X < nextValley.Width - 1 && nextPos.Y > 0 && nextPos.Y < nextValley.Height - 1 {
			if _, ok := blizzards[nextPos]; !ok {
				moves = append(moves, nextPos)
			}
		}
	}
	var result []Expedition
	for _, move := range moves {
		result = append(result, Expedition{
			Position: move,
			Valley:   nextValley,
			Moves:    e.Moves + 1,
			Leg: e.Leg,
			Last: &e,
		})
	}
	return result
}

func (e Expedition) Key() string {
	return fmt.Sprintf("%d,%d %d %d", e.Position.X, e.Position.Y, e.Leg, e.Moves)
}

type minExpedition []Expedition

func (h minExpedition) Len() int { return len(h) }
func (h minExpedition) Less(i, j int) bool { return h[i].Moves < h[j].Moves }
func (h minExpedition) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *minExpedition) Push(x interface{}) {
	*h = append(*h, x.(Expedition))
}

func (h *minExpedition) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func printExpedition(e Expedition) {
	blizzards := make(map[Point]string)
	for _, b := range e.Valley.Blizzards {
		if blizzards[b.Position] == "" {
			blizzards[b.Position] = b.Direction.String()
		} else {
			blizzards[b.Position] = "2"
		}
	}
	for i := 0; i < e.Valley.Height; i++ {
		for j := 0; j < e.Valley.Width; j++ {
			if _, ok := e.Valley.Walls[Point{X: j, Y: i}]; ok {
				print("#")
			} else if _, ok := blizzards[Point{X: j, Y: i}]; ok {
				print(blizzards[Point{X: j, Y: i}])
			} else if e.Position.X == j && e.Position.Y == i {
				print("E")
			} else {
				print(".")
			}
		}
		println()
	}
}