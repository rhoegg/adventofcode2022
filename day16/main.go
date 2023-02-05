package main

import (
	"fmt"
	"log"
)

func main() {
	part2("input.txt")
}

func part2(filename string) {
	pipeSystem := NewPipeSystem(filename)

	for _, v := range pipeSystem.Valves {
		fmt.Printf("Valve %s: %d (%v)\n", v.Name, v.FlowRate, v.Tunnels)
	}
	log.Println("Starting")
	pressure := pipeSystem.MostPressureWithElephant([2][]string{{"AA"}, {"AA"}},
		make(map[string]bool), 26)
	fmt.Printf("Max tunnel pressure = %d\nme      : %v\nelephant: %v\n", pressure.pressure, pressure.paths[0], pressure.paths[1])

}
