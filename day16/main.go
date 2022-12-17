package main

import (
	"fmt"
)

func main() {
	pipeSystem := NewPipeSystem("input.txt")

	for _, v := range pipeSystem.Valves {
		fmt.Printf("Valve %s: %d (%v)\n", v.Name, v.FlowRate, v.Tunnels)
	}
	pressure := pipeSystem.MostPressureWithElephant([2][]string{{"AA"}, {"AA"}},
		make(map[string]bool), 25) // one less than minutes for some reason
	fmt.Printf("Max tunnel pressure = %d\nme: %v\nelephant: %v\n", pressure.pressure, pressure.paths[0], pressure.paths[1])
}

