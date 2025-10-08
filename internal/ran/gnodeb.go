package ran

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/rizpur/NetSim5G/internal/ue"
)

type GNodeB struct {
	ID           string
	MaxCap       int
	AllowedIMSIs map[string]bool
	ConnectedUEs map[string]*ue.UE
}

func (g *GNodeB) ConnectUE(u *ue.UE) error {
	if len(g.ConnectedUEs) >= g.MaxCap { //guard clauses
		return fmt.Errorf("connection error: max cap reached")
	}

	if !g.AllowedIMSIs[u.IMSI] {
		return fmt.Errorf("UE not allowed to connect")
	}

	u.State = ue.Connected
	g.ConnectedUEs[u.IMSI] = u //i want the imsi to = u pointer address
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

func NewGNodeB(ID string, MaxCap int) (*GNodeB, error) {
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
		ID:           ID,
		MaxCap:       MaxCap,
		AllowedIMSIs: allowedIMSIs,
		ConnectedUEs: make(map[string]*ue.UE),
	}, nil
}
