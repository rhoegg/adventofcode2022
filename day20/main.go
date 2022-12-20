package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type MixNumber struct {
	Number int
	OriginalPosition int
}

func (n MixNumber) String() string {
	return strconv.Itoa(n.Number)
}

func mix(list []MixNumber) []MixNumber {
	//log.Printf("Mix: %v", list)
	for i := 0; i < len(list); i++ {
		// mover
		for j := 0; j < len(list); j++ {
			// position
			if list[j].OriginalPosition == i {
				mover := list[j]
				next := make([]MixNumber, len(list))
				relativeMove := mover.Number
				for relativeMove > len(list) - 1 {
					relativeMove -= len(list) - 1
				}
				for relativeMove < 0 {
					relativeMove += len(list) - 1
				}
				for k := 0; k < relativeMove; k++ {
					next[(k+j) % len(list)] = list[(k+j+1) % len(list)]
				}
				next[(relativeMove+j) % len(list)] = mover
				for k := relativeMove + 1; k < len(list); k++ {
					next[(k+j) % len(list)] = list[(k+j) % len(list)]
				}
				if j + relativeMove > len(list) {
					next = append(next[len(next)-1:], next[0:len(next)-1]...)
				}
				list = next
				break
			}
		}
		//log.Printf("-  %v", list)
	}
	return list
}

func LocateCoordinate(list []MixNumber, i int) int {
	var zeroLoc int
	for {
		if list[zeroLoc].Number == 0 {
			break
		}
		zeroLoc++
		if zeroLoc > len(list) {
			panic("can't find zero in list")
		}
	}
	return list[(zeroLoc + i) % len(list)].Number
}

func main() {
	coordinateFile, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	var list []MixNumber
	for i, l := range strings.Split(string(coordinateFile), "\n") {
		num, _ := strconv.Atoi(l)
		list = append(list, MixNumber{Number: num, OriginalPosition: i})
	}
	//fmt.Printf("before %v \n", list)
	mixed := mix(list)
	//fmt.Printf("after %v\n", list)
	fmt.Printf("coords: %d, %d, %d\n",
		LocateCoordinate(mixed, 1000), LocateCoordinate(mixed, 2000), LocateCoordinate(mixed, 3000))
	fmt.Printf("Part 1: %d\n", LocateCoordinate(mixed, 1000) + LocateCoordinate(mixed, 2000) + LocateCoordinate(mixed, 3000))
}
