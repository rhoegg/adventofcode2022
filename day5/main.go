package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Move struct {
	Count, Source, Target int
}

func ParseStacks(data []byte) (stacks []string) {

	lines := strings.Split(string(data), "\n")
	r, _ := regexp.Compile(`^(\s+\d+)+$`)

	for _, line := range lines {
		fmt.Printf("%s", line)
		if (r.MatchString(line)) {
			fmt.Println()
			break
		} else {
			for i := 0; i < (len(line) + 1) / 4; i++ {
				stack := i + 1

				if line[i*4] == '[' {
					crate := string(line[i*4 + 1])
					fmt.Printf(": %d %s", stack, crate)
					for len(stacks) < stack {
						stacks = append(stacks, "")
					}
					stacks[stack - 1] = stacks[stack - 1] + crate
				}
			}
			fmt.Println()
		}
	}
	return
}

func ParseMoves(data []byte) (moves []Move) {
	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "move") {
			tokens := strings.Split(line, " ")
			count, _ := strconv.Atoi(tokens[1])
			source, _ := strconv.Atoi(tokens[3])
			target, _ := strconv.Atoi(tokens[5])
			moves = append(moves, Move{
				Count: count,
				Source: source,
				Target: target,
			})
		}
	}
	return
}

func main() {
	data, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	//fmt.Println(string(data))

	stacks := ParseStacks(data)
	for i, s := range stacks {
		fmt.Printf("Stack %d: %s\n", i, s)
	}

	for _, move := range ParseMoves(data) {
		fmt.Printf("Move %d from %d to %d\n", move.Count, move.Source, move.Target)
		crates := stacks[move.Source - 1][0:move.Count]
		fmt.Printf("crates to move: %s\n", crates)
		stacks[move.Source - 1] = stacks[move.Source - 1][move.Count:]
		// Part 1 method
		//for _, c := range crates {
		//	stacks[move.Target - 1] = string(c) + stacks[move.Target - 1]
		//}
		stacks[move.Target - 1] = crates + stacks[move.Target - 1]
		for i, s := range stacks {
			fmt.Printf("Stack %d: %s\n", i, s)
		}
	}

}
