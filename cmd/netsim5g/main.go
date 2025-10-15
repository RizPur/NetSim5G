package main

import (
	"fmt"

	"github.com/rizpur/NetSim5G/internal/core/amf"
	"github.com/rizpur/NetSim5G/internal/core/udm"
	"github.com/rizpur/NetSim5G/internal/ran"
	"github.com/rizpur/NetSim5G/internal/ue"
)

func main() {
	fmt.Println("=== Initializing 5G Network ===")

	// Step 1: Create UDM (subscriber database)
	udmInstance, err := udm.NewUDM()
	if err != nil {
		panic(err)
	}
	fmt.Println("✓ UDM initialized")

	// Step 2: Create AMF (needs UDM)
	amfInstance := amf.NewAMF(udmInstance)
	fmt.Println("✓ AMF initialized")

	// Step 3: Create gNodeB (needs AMF)
	gnb, err := ran.NewGNodeB("gNB-001", 3, amfInstance)
	if err != nil {
		panic(err)
	}
	fmt.Println("✓ gNodeB initialized")

	fmt.Println("\n=== Testing UE Connections ===")

	// Create UEs
	ue1 := ue.NewUE("123456789012345") // Valid, active subscription
	ue2 := ue.NewUE("111222333444555") // Valid, but SUSPENDED subscription
	ue3 := ue.NewUE("999999999999999") // Invalid IMSI (not in UDM)

	// Test 1: Valid subscriber with active subscription
	fmt.Println("\n[Test 1] UE1 (active subscriber):")
	if err := gnb.ConnectUE(ue1); err != nil {
		fmt.Println("  ❌ Failed:", err)
	} else {
		fmt.Println("  ✓ Connected! State:", ue1.State)
		fmt.Printf("  ✓ Registered with AMF via %s\n", gnb.ID)
	}

	// Test 2: Valid subscriber but SUSPENDED
	fmt.Println("\n[Test 2] UE2 (suspended subscriber):")
	if err := gnb.ConnectUE(ue2); err != nil {
		fmt.Println("  ❌ Failed:", err)
	} else {
		fmt.Println("  ✓ Connected! State:", ue2.State)
	}

	// Test 3: Unknown subscriber (not in UDM)
	fmt.Println("\n[Test 3] UE3 (unknown IMSI):")
	if err := gnb.ConnectUE(ue3); err != nil {
		fmt.Println("  ❌ Failed:", err)
	} else {
		fmt.Println("  ✓ Connected! State:", ue3.State)
	}

	// Show registered UEs in AMF
	fmt.Println("\n=== AMF Registered UEs ===")
	fmt.Printf("Total registered: %d\n", len(amfInstance.RegisteredUEs))
	for imsi, regUE := range amfInstance.RegisteredUEs {
		fmt.Printf("  - IMSI: %s, gNodeB: %s, Time: %s\n", imsi, regUE.GNodeBID, regUE.Timestamp.Format("15:04:05"))
	}
}
