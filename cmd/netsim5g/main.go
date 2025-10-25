package main

import (
	"fmt"
	"net/http"

	"github.com/rizpur/NetSim5G/internal/api"
	"github.com/rizpur/NetSim5G/internal/core/amf"
	"github.com/rizpur/NetSim5G/internal/core/smf"
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

	// Step 3: Create SMF (needs UDM)
	smfInstance := smf.NewSMF(udmInstance)
	fmt.Println("✓ SMF initialized")

	// Track all UEs and gNodeBs for API
	allUEs := make(map[string]*ue.UE)
	gNodeBs := make(map[int]*ran.GNodeB)

	// Step 4: Create 2 gNodeBs at different locations
	gnb1, err := ran.NewGNodeB(100, 100, 50, 3) // gNodeB-1 at (100,100) with 50m range
	if err != nil {
		panic(err)
	}
	amfInstance.RegisterGNodeB(gnb1)
	gNodeBs[gnb1.ID] = gnb1
	fmt.Printf("✓ gNodeB-%d initialized at (%.0f, %.0f) with %.0fm range\n", gnb1.ID, gnb1.X, gnb1.Y, gnb1.Range)

	gnb2, err := ran.NewGNodeB(200, 200, 50, 3) // gNodeB-2 at (200,200) with 50m range
	if err != nil {
		panic(err)
	}
	amfInstance.RegisterGNodeB(gnb2)
	gNodeBs[gnb2.ID] = gnb2
	fmt.Printf("✓ gNodeB-%d initialized at (%.0f, %.0f) with %.0fm range\n", gnb2.ID, gnb2.X, gnb2.Y, gnb2.Range)

	fmt.Println("\n=== Testing Handover: UE Moves Between gNodeBs ===")

	// Create UE near gNodeB-1
	ue1 := ue.NewUE("123456789012345", 110, 110) // Close to gNodeB-1 at (100,100)
	allUEs[ue1.IMSI] = ue1
	fmt.Printf("\n[UE1] Created at position (%.0f, %.0f)\n", ue1.X, ue1.Y)

	// Step 1: Connect to gNodeB-1
	fmt.Println("\n--- Step 1: Initial Connection ---")
	if err := gnb1.ConnectUE(ue1); err != nil {
		fmt.Println("❌ Connection failed:", err)
		return
	}
	// Register with AMF after radio connection
	if err := amfInstance.RegisterUE(ue1.IMSI, gnb1.ID); err != nil {
		fmt.Println("❌ AMF registration failed:", err)
		return
	}
	fmt.Printf("✓ UE1 connected to gNodeB-%d\n", gnb1.ID)
	fmt.Printf("✓ AMF registration: gNodeB-%d\n", amfInstance.RegisteredUEs[ue1.IMSI].GNodeBID)

	// Step 2: Establish VoIP session
	fmt.Println("\n--- Step 2: Establish VoIP Session ---")
	voipSession, err := smfInstance.EstablishSession(ue1, smf.VoIP)
	if err != nil {
		fmt.Println("❌ Session failed:", err)
		return
	}
	fmt.Printf("✓ Session %d established: VoIP (1 Mbps, 10ms latency)\n", voipSession.SessionID)

	// Step 3: UE moves closer to gNodeB-2
	fmt.Println("\n--- Step 3: UE Moves Towards gNodeB-2 ---")
	fmt.Println("  Moving UE1 from (110, 110) → (180, 180)...")
	if err := amfInstance.MoveUE(ue1, 180, 180); err != nil {
		fmt.Println("❌ Move failed:", err)
		return
	}
	fmt.Printf("✓ UE1 moved to (%.0f, %.0f)\n", ue1.X, ue1.Y)
	fmt.Printf("✓ Handover completed! Now connected to gNodeB-%d\n", amfInstance.RegisteredUEs[ue1.IMSI].GNodeBID)
	fmt.Printf("✓ VoIP session still active: Session %d\n", voipSession.SessionID)

	// Step 4: UE moves even further (only in range of gNodeB-2)
	fmt.Println("\n--- Step 4: UE Continues Moving ---")
	fmt.Println("  Moving UE1 from (180, 180) → (210, 210)...")
	if err := amfInstance.MoveUE(ue1, 210, 210); err != nil {
		fmt.Println("❌ Move failed:", err)
		return
	}
	fmt.Printf("✓ UE1 moved to (%.0f, %.0f)\n", ue1.X, ue1.Y)
	fmt.Printf("✓ Still connected to gNodeB-%d\n", amfInstance.RegisteredUEs[ue1.IMSI].GNodeBID)

	// Step 5: UE moves out of all range
	fmt.Println("\n--- Step 5: UE Moves Out of Range ---")
	fmt.Println("  Moving UE1 from (210, 210) → (500, 500)...")
	if err := amfInstance.MoveUE(ue1, 500, 500); err != nil {
		fmt.Println("❌ Expected error:", err)
	} else {
		fmt.Println("✓ UE moved successfully")
	}
	fmt.Printf("  UE1 final state: %s\n", ue1.State)

	// Final Summary
	fmt.Println("\n=== Final Network Status ===")
	fmt.Printf("\nAMF Registered UEs: %d\n", len(amfInstance.RegisteredUEs))
	for imsi, regUE := range amfInstance.RegisteredUEs {
		fmt.Printf("  - IMSI: %s, gNodeB: %d\n", imsi, regUE.GNodeBID)
	}

	fmt.Printf("\ngNodeB-1 Connected UEs: %d\n", len(gnb1.ConnectedUEs))
	fmt.Printf("gNodeB-2 Connected UEs: %d\n", len(gnb2.ConnectedUEs))

	fmt.Printf("\nSMF Active Sessions: %d\n", len(smfInstance.Sessions))
	for id, session := range smfInstance.Sessions {
		fmt.Printf("  - Session %d: UE %s, Type: %s, State: %s\n",
			id, session.UE.IMSI, session.SessionType, session.State)
	}

	// Start API server
	handler := api.NewHandler(amfInstance, smfInstance, udmInstance, gNodeBs, allUEs)
	handler.RegisterRoutes()

	fmt.Println("\n=== Starting API Server ===")
	fmt.Println("API running on http://localhost:8080")
	fmt.Println("Try: curl http://localhost:8080/api/ues")
	http.ListenAndServe(":8080", nil)
}
