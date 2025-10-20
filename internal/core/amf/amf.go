package amf

import (
	"fmt"
	"time"

	"github.com/rizpur/NetSim5G/internal/core/udm"
	"github.com/rizpur/NetSim5G/internal/ran"
	"github.com/rizpur/NetSim5G/internal/ue"
	"github.com/rizpur/NetSim5G/internal/utils"
)

type RegisteredUE struct {
	IMSI      string
	GNodeBID  int
	Timestamp time.Time
}

type AMF struct {
	RegisteredUEs map[string]*RegisteredUE // key = IMSI
	ActiveGNodeBs map[int]*ran.GNodeB
	udm           *udm.UDM // for subscriber validation
}

func NewAMF(udm *udm.UDM) *AMF {
	return &AMF{
		RegisteredUEs: make(map[string]*RegisteredUE),
		ActiveGNodeBs: make(map[int]*ran.GNodeB),
		udm:           udm,
	}
}

// map returns val, exists := getMap["1234"]

func (a *AMF) RegisterUE(imsi string, gnbID int) error {
	// Check with UDM if subscriber exists / l'abonné existe ?
	subscriber, err := a.udm.GetSubscriber(imsi)
	if err != nil {
		return fmt.Errorf("registration failed: %w", err)
	}

	// check if subscription is active / son abonnement est il actif ?
	if subscriber.SubscriptionStatus != "active" {
		return fmt.Errorf("registration failed: subscription status is %s", subscriber.SubscriptionStatus)
	}

	// Step 3: Register the UE
	a.RegisteredUEs[imsi] = &RegisteredUE{
		IMSI:      imsi,
		GNodeBID:  gnbID,
		Timestamp: time.Now(),
	}

	return nil
}

func (a *AMF) RegisterGNodeB(g *ran.GNodeB) {
	//append this to mapping somehow
	a.ActiveGNodeBs[g.ID] = g

}

func (a *AMF) MoveUE(u *ue.UE, newX, newY float64) error {
	u.X = newX
	u.Y = newY

	regUE, exists := a.RegisteredUEs[u.IMSI]
	if !exists {
		return nil // ue not registered, no handover needed just move and leave
	}
	currentGNodeB, exists := a.ActiveGNodeBs[regUE.GNodeBID]
	if !exists {

		return fmt.Errorf("data inconsistency: UE registered to non-existent gNodeB %d", regUE.GNodeBID)
	} // UE not connected to any GNodeB that exists

	// if this UE is registered and has a current GNodeB dont throw an error, contineu
	distance := utils.CalculateDistance(u.X, u.Y, currentGNodeB.X, currentGNodeB.Y)
	if distance > currentGNodeB.Range {
		a.Handover(u, currentGNodeB)
	}

	return nil
}

func (a *AMF) Handover(u *ue.UE, oldG *ran.GNodeB) error {
	// get all GNodeB position
	//calculate
	var closestGNodeB *ran.GNodeB
	var closestDistance float64 = 2000.00

	for _, g := range a.ActiveGNodeBs {
		distance := utils.CalculateDistance(u.X, u.Y, g.X, g.Y)
		if distance <= closestDistance && distance <= g.Range {
			closestDistance = distance
			closestGNodeB = g
		}
	}

	if closestDistance == 2000.00 {
		// no gNodeB in range of UE,
		return fmt.Errorf("No GNodeB in range of this UE")
	} else {
		oldG.Disconnect(u)
		closestGNodeB.ConnectUE(u)
		// a.RegisterUE(u.IMSI, closestGNodeB.ID) just need to update gNodeB ID
		a.RegisteredUEs[u.IMSI].GNodeBID = closestGNodeB.ID
	}
	return nil

}
