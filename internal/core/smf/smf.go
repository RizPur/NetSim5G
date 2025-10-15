package smf

import (
	"fmt"

	"github.com/rizpur/NetSim5G/internal/core/udm"
	"github.com/rizpur/NetSim5G/internal/ue"
)

// SessionType represents different types of data sessions
type SessionType int

const (
	VoIP SessionType = iota
	VideoStreaming
	WebBrowsing
)

func (s SessionType) String() string {
	return [...]string{"VoIP", "VideoStreaming", "WebBrowsing"}[s]
}

// SessionState represents the state of a PDU session
type SessionState int

const (
	Active SessionState = iota
	Inactive
)

func (s SessionState) String() string {
	return [...]string{"Active", "Inactive"}[s]
}

// QoSProfile defines Quality of Service parameters
type QoSProfile struct {
	MaxBitRate int // Mbps
	Latency    int // milliseconds
	Priority   int // 1 (highest) - 10 (lowest)
}

// Predefined QoS profiles for each session type
var QoSProfiles = map[SessionType]QoSProfile{
	VoIP: {
		MaxBitRate: 1,  // 1 Mbps - voice doesn't need much its tiny data
		Latency:    10, // 10ms - super low latency for real-time voice
		Priority:   1,  // Highest priority
	},
	VideoStreaming: {
		MaxBitRate: 50,  // 50 Mbps - 4K video , gaming
		Latency:    100, // 100ms - buffering handles some delay
		Priority:   5,   // Medium priority
	},
	WebBrowsing: {
		MaxBitRate: 25,  // 25 Mbps - load pages quickly
		Latency:    300, // 300ms - its ok
		Priority:   10,  // Lowest priority (best effort)
	},
}

// PDUSession represents a data session between UE and network
type PDUSession struct {
	SessionID   int
	UE          *ue.UE
	SessionType SessionType
	QoS         QoSProfile
	State       SessionState
}

// SMF manages PDU sessions
type SMF struct {
	Sessions      map[int]*PDUSession // key = SessionID
	udm           *udm.UDM
	nextSessionID int
}

func NewSMF(udm *udm.UDM) *SMF {
	return &SMF{
		Sessions:      make(map[int]*PDUSession),
		udm:           udm,
		nextSessionID: 1, // Start session IDs at 1
	}
}

// EstablishSession creates a new PDU session for a UE
func (s *SMF) EstablishSession(u *ue.UE, sessionType SessionType) (*PDUSession, error) {
	// Step 1: Get subscriber info from UDM
	subscriber, err := s.udm.GetSubscriber(u.IMSI)
	if err != nil {
		return nil, fmt.Errorf("session establishment failed: %w", err)
	}

	// Step 2: Get QoS profile for requested session type
	qosProfile := QoSProfiles[sessionType]

	// Step 3: Calculate total bandwidth already allocated to this UE
	totalAllocated := 0
	for _, session := range s.Sessions {
		if session.UE.IMSI == u.IMSI && session.State == Active {
			totalAllocated += session.QoS.MaxBitRate
		}
	}

	// Step 4: Check if adding this session would exceed subscriber's limit
	if totalAllocated+qosProfile.MaxBitRate > subscriber.MaxDataRate {
		return nil, fmt.Errorf("session establishment failed: total bandwidth would be %d Mbps (current: %d + new: %d), exceeds subscriber limit of %d Mbps",
			totalAllocated+qosProfile.MaxBitRate, totalAllocated, qosProfile.MaxBitRate, subscriber.MaxDataRate)
	}

	// Step 5: Create the session
	session := &PDUSession{
		SessionID:   s.nextSessionID,
		UE:          u,
		SessionType: sessionType,
		QoS:         qosProfile,
		State:       Active,
	}

	s.Sessions[s.nextSessionID] = session
	s.nextSessionID++

	return session, nil
}

// TerminateSession ends a PDU session
func (s *SMF) TerminateSession(sessionID int) error {
	session, exists := s.Sessions[sessionID]
	if !exists {
		return fmt.Errorf("session %d not found", sessionID)
	}

	session.State = Inactive
	delete(s.Sessions, sessionID)
	return nil
}
