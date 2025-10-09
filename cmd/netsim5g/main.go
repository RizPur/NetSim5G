package main

import (
	"fmt"

	"github.com/rizpur/NetSim5G/internal/core/udm"
	"github.com/rizpur/NetSim5G/internal/ran"
	"github.com/rizpur/NetSim5G/internal/ue"
)

func main() {
	// Create gNodeB
	gnb, err := ran.NewGNodeB("gNB-001", 3)
	if err != nil {
		panic(err)
	}

	// Create UEs
	ue1 := ue.NewUE("123456789012345") // Allowed, pointer returned
	ue2 := ue.NewUE("999999999999999") // Not allowed
	ue3 := ue.NewUE("00")
	ue4 := ue.NewUE("01")
	ue5 := ue.NewUE("11")

	// Try connecting
	if err := gnb.ConnectUE(ue1); err != nil { // First do this; then check this, if err:= gnb.ConnectUE(ue1); err != nil {ue didnt connect}
		fmt.Println("UE1 failed:", err)
	} else {
		fmt.Println("UE1 connected! State:", ue1.State)
	}

	if err := gnb.ConnectUE(ue2); err != nil {
		fmt.Println("UE2 failed:", err)
	}

	if err := gnb.ConnectUE(ue3); err != nil {
		fmt.Println("UE3 failed to connect", err)
	}

	if err := gnb.ConnectUE(ue4); err != nil {
		fmt.Println("UE4 failed to connect", err)
	}

	if err := gnb.Disconnect(ue4); err != nil {
		fmt.Println("UE4 failed to disconnect", err)
	}

	if err := gnb.ConnectUE(ue5); err != nil {
		fmt.Println("UE5 failed", err)
	}

	udm, err := udm.NewUDM()
	if err != nil {
		panic(err)
	}

	sub, err := udm.GetSubscriber("123456789012345")
	if err != nil {
		fmt.Println("Subscriber not found:", err)
	} else {
		fmt.Printf("Found: %s, Status: %s\n", sub.PhoneNumber, sub.SubscriptionStatus)
	}

}
