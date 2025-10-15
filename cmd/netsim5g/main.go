package main

import (
	"fmt"

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

	// Step 4: Create gNodeB (needs AMF)
	gnb, err := ran.NewGNodeB("gNB-001", 3, amfInstance)
	if err != nil {
		panic(err)
	}
	fmt.Println("✓ gNodeB initialized")

	fmt.Println("\n=== Testing Full Flow: Connect → Register → Establish Session ===")

	// Create UEs with different subscription plans
	ue1 := ue.NewUE("123456789012345") // Max 100 Mbps
	ue2 := ue.NewUE("987654321098765") // Max 50 Mbps

	// ==================== UE1 Flow ====================
	fmt.Println("\n[UE1 - High-tier subscriber (100 Mbps max)]")

	// Step 1: Connect to gNodeB and register with AMF
	fmt.Println("  Step 1: Connecting to gNodeB...")
	if err := gnb.ConnectUE(ue1); err != nil {
		fmt.Println("    ❌ Failed:", err)
		return
	}
	fmt.Printf("    ✓ Connected! Radio state: %s\n", ue1.State)
	fmt.Printf("    ✓ Registered with AMF\n")

	// Step 2: Establish VoIP session
	fmt.Println("  Step 2: Establishing VoIP session...")
	voipSession, err := smfInstance.EstablishSession(ue1, smf.VoIP)
	if err != nil {
		fmt.Println("    ❌ Failed:", err)
	} else {
		fmt.Printf("    ✓ Session %d established: %s\n", voipSession.SessionID, voipSession.SessionType)
		fmt.Printf("    ✓ QoS: %d Mbps, %dms latency, Priority %d\n",
			voipSession.QoS.MaxBitRate, voipSession.QoS.Latency, voipSession.QoS.Priority)
	}

	// Step 3: Establish Video session (multiple sessions!)
	fmt.Println("  Step 3: Establishing Video session...")
	videoSession, err := smfInstance.EstablishSession(ue1, smf.VideoStreaming)
	if err != nil {
		fmt.Println("    ❌ Failed:", err)
	} else {
		fmt.Printf("    ✓ Session %d established: %s\n", videoSession.SessionID, videoSession.SessionType)
		fmt.Printf("    ✓ QoS: %d Mbps, %dms latency, Priority %d\n",
			videoSession.QoS.MaxBitRate, videoSession.QoS.Latency, videoSession.QoS.Priority)
	}

	// ==================== UE2 Flow ====================
	fmt.Println("\n[UE2 - Mid-tier subscriber (50 Mbps max)]")

	// Step 1: Connect and register
	fmt.Println("  Step 1: Connecting to gNodeB...")
	if err := gnb.ConnectUE(ue2); err != nil {
		fmt.Println("    ❌ Failed:", err)
		return
	}
	fmt.Printf("    ✓ Connected! Radio state: %s\n", ue2.State)

	// Step 2: Try to establish Video session (50 Mbps needed, has 50 Mbps max - should work!)
	fmt.Println("  Step 2: Establishing Video session (needs 50 Mbps)...")
	videoSession2, err := smfInstance.EstablishSession(ue2, smf.VideoStreaming)
	if err != nil {
		fmt.Println("    ❌ Failed:", err)
	} else {
		fmt.Printf("    ✓ Session %d established: %s\n", videoSession2.SessionID, videoSession2.SessionType)
		fmt.Printf("    ✓ QoS: %d Mbps, %dms latency, Priority %d\n",
			videoSession2.QoS.MaxBitRate, videoSession2.QoS.Latency, videoSession2.QoS.Priority)
	}

	// Step 3: Try Web browsing (would need total 75 Mbps - should FAIL!)
	fmt.Println("  Step 3: Establishing Web session (would need 75 Mbps total)...")
	webSession, err := smfInstance.EstablishSession(ue2, smf.WebBrowsing)
	if err != nil {
		fmt.Println("    ❌ Failed:", err)
	} else {
		fmt.Printf("    ✓ Session %d established: %s\n", webSession.SessionID, webSession.SessionType)
	}

	// ==================== Summary ====================
	fmt.Println("\n=== Network Status Summary ===")

	fmt.Printf("\nAMF Registered UEs: %d\n", len(amfInstance.RegisteredUEs))
	for imsi, regUE := range amfInstance.RegisteredUEs {
		fmt.Printf("  - IMSI: %s, gNodeB: %s\n", imsi, regUE.GNodeBID)
	}

	fmt.Printf("\nSMF Active Sessions: %d\n", len(smfInstance.Sessions))
	for id, session := range smfInstance.Sessions {
		fmt.Printf("  - Session %d: UE %s, Type: %s, Rate: %d Mbps, State: %s\n",
			id, session.UE.IMSI, session.SessionType, session.QoS.MaxBitRate, session.State)
	}
}
