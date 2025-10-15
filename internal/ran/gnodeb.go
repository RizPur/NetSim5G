package ran

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/rizpur/NetSim5G/internal/core/amf"
	"github.com/rizpur/NetSim5G/internal/ue"
)

// Package-level variable to track next gNodeB ID
var nextGNodeBID = 1

type GNodeB struct {
	ID           int
	X, Y         float64
	Range        float64
	MaxCap       int
	AllowedIMSIs map[string]bool
	ConnectedUEs map[string]*ue.UE
	amf          *amf.AMF // Reference to AMF for registration
}

func (g *GNodeB) ConnectUE(u *ue.UE) error {
	if len(g.ConnectedUEs) >= g.MaxCap { //guard clauses
		return fmt.Errorf("connection error: max cap reached")
	}

	if !g.AllowedIMSIs[u.IMSI] {
		return fmt.Errorf("UE not allowed to connect")
	}

	// Radio connection successful
	u.State = ue.Connected
	g.ConnectedUEs[u.IMSI] = u

	// Try to register with AMF (core network)
	if err := g.amf.RegisterUE(u.IMSI, g.ID); err != nil {
		// Registration failed - rollback the radio connection
		delete(g.ConnectedUEs, u.IMSI)
		u.State = ue.Disconnected
		return fmt.Errorf("AMF registration failed: %w", err)
	}

	return nil
}

func (g *GNodeB) Disconnect(u *ue.UE) error {
	if _, exists := g.ConnectedUEs[u.IMSI]; !exists {
		return fmt.Errorf("UE is not currently connected to gNodeB")
	}
	delete(g.ConnectedUEs, u.IMSI)
	u.State = ue.Disconnected
	return nil
}

func NewGNodeB(MaxCap int, amf *amf.AMF) (*GNodeB, error) {
	// Get the current ID and increment for next gNodeB
	id := nextGNodeBID
	nextGNodeBID++

	allowedIMSIs := make(map[string]bool)

	file, err := os.Open("internal/configs/allowed_imsis.txt") // returns file AND an err
	if err != nil {
		return nil, fmt.Errorf("failed to open allowed IMSIs file: %w", err) //gnode b gets nil
	}
	defer file.Close() //defer = "do this when the function exits, no matter what err or success"

	scanner := bufio.NewScanner(file) //scanner that reads file line by line
	for scanner.Scan() {
		imsi := strings.TrimSpace(scanner.Text()) //scanner.Text() gets current line as String
		if imsi != "" {
			allowedIMSIs[imsi] = true
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading allowed IMSIs: %w", err)
	}

	return &GNodeB{
		ID:           id,
		MaxCap:       MaxCap,
		AllowedIMSIs: allowedIMSIs,
		ConnectedUEs: make(map[string]*ue.UE),
		amf:          amf,
	}, nil
}
