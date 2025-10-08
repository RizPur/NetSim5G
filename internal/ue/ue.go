package ue

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
