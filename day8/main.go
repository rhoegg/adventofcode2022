package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Tree struct {
	Row, Col int
}

func main() {
	data, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	var forest [][]int
	for i, row := range strings.Split(string(data), "\n") {
		forest = append(forest, nil)
		for _, heightVal := range row {
			height, _ := strconv.Atoi(string(heightVal))
			forest[i] = append(forest[i], height)
		}
	}
	fmt.Printf("Analyzing %dx%d forest\n", len(forest[0]), len(forest))
	AnalyzePart1(forest)
	AnalyzePart2(forest)
}

func AnalyzePart1(forest [][]int) {
	visible := make(map[Tree]struct{})
	funcs := []func([][]int) []Tree { VisibleFromNorth, VisibleFromSouth, VisibleFromEast, VisibleFromWest }
	for _, visibleFunc := range funcs {
		for _, tree := range visibleFunc(forest) {
			visible[tree] = struct{}{}
		}
	}
	fmt.Printf("Found %d visible trees\n", len(visible))
}

func AnalyzePart2(forest [][]int) {
	highScore := 0
	for r := 0; r < len(forest); r++ {
		for c := 0; c < len(forest[r]); c++ {
			score := ScenicScore(forest, Tree{Row: r, Col: c})
			if score > highScore {
				highScore = score
			}
		}
	}
	fmt.Printf("Highest scenic score is %d\n", highScore)
}

func VisibleFromWest(forest [][]int) (visibleTrees []Tree) {
	for r := 0; r < len(forest); r++ {
		lastHeight := -1
		for c := 0; c < len(forest[r]); c++ {
			if forest[r][c] > lastHeight {
				visibleTrees = append(visibleTrees, Tree{Row: r, Col: c})
				lastHeight = forest[r][c]
			}
		}
	}
	return visibleTrees
}

func VisibleFromEast(forest [][]int) (visibleTrees []Tree) {
	for r := 0; r < len(forest); r++ {
		lastHeight := -1
		for c := len(forest[r]) - 1; c >= 0; c-- {
			if forest[r][c] > lastHeight {
				visibleTrees = append(visibleTrees, Tree{Row: r, Col: c})
				lastHeight = forest[r][c]
			}
		}
	}
	return visibleTrees
}

func VisibleFromNorth(forest [][]int) (visibleTrees []Tree) {
	for c := 0; c < len(forest[0]); c++ {
		lastHeight := -1
		for r := 0; r < len(forest); r++ {
			if forest[r][c] > lastHeight {
				visibleTrees = append(visibleTrees, Tree{Row: r, Col: c})
				lastHeight = forest[r][c]
			}
		}
	}
	return visibleTrees
}

func VisibleFromSouth(forest [][]int) (visibleTrees []Tree) {
	for c := 0; c < len(forest[0]); c++ {
		lastHeight := -1
		for r := len(forest) - 1; r >= 0; r-- {
			if forest[r][c] > lastHeight {
				visibleTrees = append(visibleTrees, Tree{Row: r, Col: c})
				lastHeight = forest[r][c]
			}
		}
	}
	return visibleTrees
}

func ScenicScore(forest [][]int, tree Tree) int {
	north, south, east, west := 0, 0, 0, 0
	treeHeight := forest[tree.Row][tree.Col]
	// west
	for c := tree.Col - 1; c >= 0; c-- {
		west++
		blocked := forest[tree.Row][c] >= treeHeight
		if blocked {
			break
		}
	}
	// east
	for c := tree.Col + 1; c < len(forest[0]); c++ {
		east++
		blocked := forest[tree.Row][c] >= treeHeight
		if blocked {
			break
		}
	}
	// north
	for r := tree.Row - 1; r >= 0; r-- {
		north++
		blocked := forest[r][tree.Col] >= treeHeight
		if blocked {
			break
		}
	}
	//south
	for r := tree.Row + 1; r < len(forest); r++ {
		south++
		blocked := forest[r][tree.Col] >= treeHeight
		if blocked {
			break
		}
	}
	//fmt.Printf("Scenic Score (%d, %d) = %d %d %d %d = %d\n", tree.Row, tree.Col, north, south, east, west, north * south * east * west)
	return north * south * east * west
}