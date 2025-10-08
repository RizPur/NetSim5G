package ue

import "fmt"

type UE struct {
	IMSI            string
	GNodeBConnected int
	State           UEState
}

type UEState int

const (
	Disconnected UEState = iota
	Connected
	Idle
)

func (u UEState) String() string {
	return [...]string{"disconnected", "connected", "idle"}[u]
}

func NewUE(imsi string) *UE {
	return &UE{
		IMSI:            imsi,
		GNodeBConnected: -1,
		State:           Disconnected,
	}
	//   1. *UE in a type (like func NewUE() *UE) = "pointer to UE"   - * = "I want a pointer to..."
	//	2. &UE{...} when creating = "give me the address of this new UE" - & = "here's the address of..."
}

// / Pointer examples stuff
func pointer_example() {
	// Scenario 1: Using a VALUE (no pointer)
	ue3 := UE{IMSI: "789", State: Disconnected}
	fmt.Printf("Before connectUE_WithValue: %s\n", ue3.State)
	connectUE_WithValue(ue3)
	fmt.Printf("After connectUE_WithValue: %s (unchanged!)\n", ue3.State)

	// Scenario 2: Using a POINTER
	ue4 := &UE{IMSI: "999", State: Disconnected}
	fmt.Printf("\nBefore connectUE_WithPointer: %s\n", ue4.State)
	connectUE_WithPointer(ue4)
	fmt.Printf("After connectUE_WithPointer: %s (changed!)\n", ue4.State)
}

// Scenario 1: Using a VALUE (no pointer)
func connectUE_WithValue(u UE) {
	u.State = Connected   // Changes the COPY
	u.GNodeBConnected = 1 // Changes the COPY
	fmt.Println("  Inside function, changed state to:", u.State)
}

// Scenario 2: Using a POINTER
func connectUE_WithPointer(u *UE) {
	u.State = Connected   // Changes the ORIGINAL
	u.GNodeBConnected = 1 // Changes the ORIGINAL
	fmt.Println("  Inside function, changed state to:", u.State)
}
