package main

import (
	"fmt"

	"github.com/rizpur/NetSim5G/internal/ue"
)

func main() {
	pointer_example()
}

func pointer_example() {
	// Scenario 1: Using a VALUE (no pointer)
	ue3 := ue.UE{IMSI: "789", State: ue.Disconnected}
	fmt.Printf("Before connectUE_WithValue: %s\n", ue3.State)
	connectUE_WithValue(ue3)
	fmt.Printf("After connectUE_WithValue: %s (unchanged!)\n", ue3.State)

	// Scenario 2: Using a POINTER
	ue4 := &ue.UE{IMSI: "999", State: ue.Disconnected}
	fmt.Printf("\nBefore connectUE_WithPointer: %s\n", ue4.State)
	connectUE_WithPointer(ue4)
	fmt.Printf("After connectUE_WithPointer: %s (changed!)\n", ue4.State)
}

// Scenario 1: Using a VALUE (no pointer)
func connectUE_WithValue(u ue.UE) {
	u.State = ue.Connected // Changes the COPY
	u.GNodeBConnected = 1  // Changes the COPY
	fmt.Println("  Inside function, changed state to:", u.State)
}

// Scenario 2: Using a POINTER
func connectUE_WithPointer(u *ue.UE) {
	u.State = ue.Connected // Changes the ORIGINAL
	u.GNodeBConnected = 1  // Changes the ORIGINAL
	fmt.Println("  Inside function, changed state to:", u.State)
}
