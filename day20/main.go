package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type MixNumber struct {
	Number int64
	OriginalPosition int
}

func (n MixNumber) String() string {
	return strconv.FormatInt(n.Number, 10)
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
				if relativeMove > int64(len(list) - 1) {
					relativeMove -= (relativeMove / int64(len(list) - 1)) * int64(len(list) - 1)
				}
				if relativeMove < 0 {
					relativeMove -= (relativeMove / int64(len(list) - 1) - 1) * int64(len(list) - 1)
				}

				for k := 0; k < int(relativeMove); k++ {
					next[(k+j) % len(list)] = list[(k+j+1) % len(list)]
				}
				next[(int(relativeMove)+j) % len(list)] = mover
				for k := int(relativeMove) + 1; k < len(list); k++ {
					next[(k+j) % len(list)] = list[(k+j) % len(list)]
				}
				if j + int(relativeMove) > len(list) {
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

func LocateCoordinate(list []MixNumber, i int) int64 {
	var zeroLoc int
	for {
		if list[zeroLoc].Number == 0 {
			break
		}
		zeroLoc++
		if zeroLoc >= len(list) {
			panic("can't find zero in list")
		}
	}
	return list[(zeroLoc + i) % len(list)].Number
}

func ApplyDecryptionKey(key int64, message []MixNumber) []MixNumber {
	result := make([]MixNumber, len(message))
	for i, n := range message {
		result[i] = MixNumber{
			Number: n.Number * key,
			OriginalPosition: n.OriginalPosition,
		}
	}
	return result
}

func main() {
	coordinateFile, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	var list []MixNumber
	for i, l := range strings.Split(string(coordinateFile), "\n") {
		num, _ := strconv.ParseInt(l, 10, 64)
		list = append(list, MixNumber{Number: num, OriginalPosition: i})
	}
	fmt.Printf("before %v \n", list)
	mixed := mix(list)
	fmt.Printf("after %v\n", list)
	fmt.Printf("coords: %d, %d, %d\n",
		LocateCoordinate(mixed, 1000), LocateCoordinate(mixed, 2000), LocateCoordinate(mixed, 3000))
	fmt.Printf("Part 1: %d\n", LocateCoordinate(mixed, 1000) + LocateCoordinate(mixed, 2000) + LocateCoordinate(mixed, 3000))

	workingMessage := ApplyDecryptionKey(811589153, list)
	part2Mixed := mix(workingMessage)
	for i := 0; i < 9; i++ {
		part2Mixed = mix(part2Mixed)
	}
	fmt.Printf("coords: %d, %d, %d\n",
		LocateCoordinate(part2Mixed, 1000), LocateCoordinate(part2Mixed, 2000), LocateCoordinate(part2Mixed, 3000))
	fmt.Printf("Part 2: %d\n", LocateCoordinate(part2Mixed, 1000) + LocateCoordinate(part2Mixed, 2000) + LocateCoordinate(part2Mixed, 3000))


}
