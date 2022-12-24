package main

import (
	"fmt"
	"log"
)

func main() {
	crater := ParseCrater("input.txt")

	var rounds int
	for {
		rounds++
		crater = crater.NextRound()
		if crater.ElvesMoved == 0 {
			break
		} else {
			log.Printf("round %d moved %d\n", rounds, crater.ElvesMoved)
		}
	}
	println(crater.String())
	//fmt.Printf("Part 1 empty count %d\n", crater.EmptyInRect())
	fmt.Printf("Part 2 - rounds until none moved = %d\n", rounds)
}
