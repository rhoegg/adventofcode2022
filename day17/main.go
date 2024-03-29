package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	rocks := parseRockPattern("rocks.txt")
	jetPattern := parseJetPattern("example.txt")
	chamber := NewChamber(7, rocks, jetPattern)
	// This times at 5s for 1,000,000
	for i := 0; i < 5; i++ {
		if i % 50000 == 0 {
			log.Printf("Dropped %d", i)
		}
		chamber.DropRock()
	}
	fmt.Printf("Height: %d\n\n", chamber.Peak())
	fmt.Println(chamber.Draw())
}

func parseJetPattern(filename string) []int {
	jetFile, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	var moves []int
	for _, c := range string(jetFile) {
		switch c {
		case '<': moves = append(moves, -1)
			break
		case '>': moves = append(moves, 1)
			break
		default: panic(fmt.Sprintf("move %s not supported", string(c)))
		}
	}
	return moves
}

func parseRockPattern(filename string) []Rock {
	rocksFile, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	var rocks []Rock
	for i, rockPattern := range strings.Split(string(rocksFile), "\n\n") {
		lines := strings.Split(rockPattern, "\n")
		rock := NewRock(i + 1)
		for y := 0; y < len(lines); y++ {
			// iterate backwards
			line := lines[len(lines) - 1 - y]
			for x := 0; x < len(line); x++ {
				if line[x] == '#' {
					rock.Shape[Point{X: x, Y: y}] = true
				}
			}
		}
		rocks = append(rocks, rock)
	}
	return rocks
}