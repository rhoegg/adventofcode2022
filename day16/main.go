package main

import (
	"fmt"
)

func main() {
	pipeSystem := NewPipeSystem("input.txt")

	for _, v := range pipeSystem.Valves {
		fmt.Printf("Valve %s: %d (%v)\n", v.Name, v.FlowRate, v.Tunnels)
	}
	waypoints, pressure := pipeSystem.MostPressure("AA", make(map[string]bool), 30)
	fmt.Printf("Max tunnel pressure = %d at %v\n", pressure, waypoints)
}

