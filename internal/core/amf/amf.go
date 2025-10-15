package amf

import (
	"fmt"
	"time"

	"github.com/rizpur/NetSim5G/internal/core/udm"
)

type RegisteredUE struct {
	IMSI      string
	GNodeBID  int
	Timestamp time.Time
}

type AMF struct {
	RegisteredUEs map[string]*RegisteredUE // key = IMSI
	udm           *udm.UDM                 // for subscriber validation
}

func NewAMF(udm *udm.UDM) *AMF {
	return &AMF{
		RegisteredUEs: make(map[string]*RegisteredUE),
		udm:           udm,
	}
}

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
