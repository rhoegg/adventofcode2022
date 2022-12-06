package main

import (
	"fmt"
	"os"
)

func allDifferent(packet []byte) bool {
	seen := make(map[byte]struct{})
	for _, c := range packet {
		if _, ok := seen[c]; ok {
			return false
		}
		seen[c] = struct{}{}
	}
	return true
}

func findDistinct(data []byte, l int) {
	window, remaining := data[0:l], data[l:]
	originalLength := len(remaining)
	for i := 0; i < originalLength; i++ {
		if allDifferent(window) {
			fmt.Printf("found %s at %d\n", window, i + l)
			return
		}
		window = append(window[1:], remaining[0])
		remaining = remaining[1:]
	}
}

func main() {
	data, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	//fmt.Println(string(data))
	findDistinct(data, 4)
	findDistinct(data, 14)

}
