package httpapi

import (
	"encoding/json"
	"net/http"
)

// AnalyzeTCGLiveRequest represents the request body for analyzing a TCG Live game.
type AnalyzeTCGLiveRequest struct {
	GameExport string `json:"gameExport"`
	IsPrivate  bool   `json:"isPrivate"`
}

// AnalyzeTCGLiveResponse represents the response for TCG Live analysis.
type AnalyzeTCGLiveResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

// handleAnalyzeTCGLive handles POST /api/tcglive/analyze requests.
func (s *Server) handleAnalyzeTCGLive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req AnalyzeTCGLiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.logger.Infof("Failed to decode request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Invalid request body",
			Code:  "INVALID_REQUEST",
		})
		return
	}

	if req.GameExport == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Error: "gameExport is required",
			Code:  "INVALID_REQUEST",
		})
		return
	}

	// TCG Live analysis is not yet implemented
	w.WriteHeader(http.StatusNotImplemented)
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Error: "TCG Live analysis is planned for a future release",
		Code:  "NOT_IMPLEMENTED",
	})
}
