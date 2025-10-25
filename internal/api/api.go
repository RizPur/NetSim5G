package api

import (
	"encoding/json"
	"net/http"

	"github.com/rizpur/NetSim5G/internal/core/amf"
	"github.com/rizpur/NetSim5G/internal/core/smf"
	"github.com/rizpur/NetSim5G/internal/core/udm"
	"github.com/rizpur/NetSim5G/internal/ran"
	"github.com/rizpur/NetSim5G/internal/ue"
)

// Handler holds references to all network components
type Handler struct {
	AMF     *amf.AMF
	SMF     *smf.SMF
	UDM     *udm.UDM
	GNodeBs map[int]*ran.GNodeB
	AllUEs  map[string]*ue.UE // All UEs that exist (connected or not)
}

// NewHandler creates a new API handler
func NewHandler(amf *amf.AMF, smf *smf.SMF, udm *udm.UDM, gnbs map[int]*ran.GNodeB, allUEs map[string]*ue.UE) *Handler {
	return &Handler{
		AMF:     amf,
		SMF:     smf,
		UDM:     udm,
		GNodeBs: gnbs,
		AllUEs:  allUEs,
	}
}

// CORS middleware - allows React dev server to call our API
func (h *Handler) enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// RegisterRoutes sets up all HTTP endpoints
func (h *Handler) RegisterRoutes() {
	http.HandleFunc("/api/gnodebs", h.enableCORS(h.getGNodeBs))
	http.HandleFunc("/api/ues", h.enableCORS(h.getUEs))
	http.HandleFunc("/api/sessions", h.enableCORS(h.getSessions))
}

// GET /api/gnodebs - returns all gNodeBs with their state
func (h *Handler) getGNodeBs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Create response slice
	type GNodeBResponse struct {
		ID           int     `json:"id"`
		X            float64 `json:"x"`
		Y            float64 `json:"y"`
		Range        float64 `json:"range"`
		ConnectedUEs int     `json:"connectedUEs"`
		MaxCap       int     `json:"maxCap"`
	}

	var response []GNodeBResponse

	// Loop through all gNodeBs and build response
	for _, gnb := range h.GNodeBs {
		response = append(response, GNodeBResponse{
			ID:           gnb.ID,
			X:            gnb.X,
			Y:            gnb.Y,
			Range:        gnb.Range,
			ConnectedUEs: len(gnb.ConnectedUEs),
			MaxCap:       gnb.MaxCap,
		})
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GET /api/ues - returns all UEs with their state
func (h *Handler) getUEs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	type UEResponse struct {
		IMSI            string     `json:"imsi"`
		X               float64    `json:"x"`
		Y               float64    `json:"y"`
		GNodeBConnected int        `json:"gNodeBConnected"`
		State           ue.UEState `json:"state"`
	}

	var response []UEResponse

	// Loop through ALL UEs (connected or not)
	for _, ue := range h.AllUEs {
		response = append(response, UEResponse{
			IMSI:            ue.IMSI,
			X:               ue.X,
			Y:               ue.Y,
			GNodeBConnected: ue.GNodeBConnected,
			State:           ue.State,
		})
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GET /api/sessions - returns all active PDU sessions
func (h *Handler) getSessions(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	type SessionResponse struct {
		SessionID   int    `json:"sessionID"`
		UEIMSI      string `json:"ueIMSI"`
		SessionType string `json:"sessionType"`
		State       string `json:"state"`
	}

	var response []SessionResponse

	for _, session := range h.SMF.Sessions {
		response = append(response, SessionResponse{
			SessionID:   session.SessionID,
			UEIMSI:      session.UE.IMSI,
			SessionType: string(session.SessionType),
			State:       string(session.State),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
