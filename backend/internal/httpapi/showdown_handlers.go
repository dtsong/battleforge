package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/dtsong/battleforge/backend/internal/analysis"
	"github.com/dtsong/battleforge/backend/internal/db"
)
// It supports three analysis types via discriminator: replayId, username, or rawLog.
type AnalyzeShowdownRequest struct {
	// Discriminator field
	AnalysisType string `json:"analysisType"`

	// For replayId analysis
	ReplayID string `json:"replayId,omitempty"`

	// For username analysis
	Username string `json:"username,omitempty"`
	Format   string `json:"format,omitempty"`
	Limit    int    `json:"limit,omitempty"`

	// For rawLog analysis
	RawLog string `json:"rawLog,omitempty"`

	// Common field
	IsPrivate bool `json:"isPrivate"`
}

// AnalyzeResponse is the response for analyze requests.
type AnalyzeResponse struct {
	Status   string                  `json:"status"`
	BattleID string                  `json:"battleId,omitempty"`
	Data     *analysis.BattleSummary `json:"data,omitempty"`
	Metadata *ResponseMetadata        `json:"metadata,omitempty"`
}

// ResponseMetadata contains metadata about the analysis.
type ResponseMetadata struct {
	ParseTimeMs   int  `json:"parseTimeMs"`
	AnalysisTimeMs int `json:"analysisTimeMs"`
	Cached        bool `json:"cached"`
}

// ErrorResponse is the response for errors.
type ErrorResponse struct {
	Error   string      `json:"error"`
	Code    string      `json:"code"`
	Details interface{} `json:"details,omitempty"`
}

// ListReplaysRequest represents query parameters for listing replays.
type ListReplaysRequest struct {
	Username  string
	Format    string
	IsPrivate *bool
	Limit     int
	Offset    int
}


// handleAnalyzeShowdown handles POST /api/showdown/analyze requests.
func (s *Server) handleAnalyzeShowdown(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	start := time.Now()

	var req AnalyzeShowdownRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.logger.Infof("Failed to decode request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Invalid request body",
			Code:  "INVALID_REQUEST",
		})
		return
	}

	// Validate request based on analysis type
	var battleSummary *analysis.BattleSummary
	var battlelLog string
	var err error

	switch req.AnalysisType {
	case "replayId":
		if req.ReplayID == "" {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Error: "replayId is required for replayId analysis",
				Code:  "INVALID_REQUEST",
			})
			return
		}
		// TODO: Fetch replay from Showdown API or cache
		w.WriteHeader(http.StatusNotImplemented)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Replay ID analysis not yet implemented",
			Code:  "NOT_IMPLEMENTED",
		})
		return

	case "username":
		if req.Username == "" || req.Format == "" {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Error: "username and format are required for username analysis",
				Code:  "INVALID_REQUEST",
			})
			return
		}
		// TODO: Fetch recent battles by username
		w.WriteHeader(http.StatusNotImplemented)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Username analysis not yet implemented",
			Code:  "NOT_IMPLEMENTED",
		})
		return

	case "rawLog":
		if req.RawLog == "" {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Error: "rawLog is required for rawLog analysis",
				Code:  "INVALID_REQUEST",
			})
			return
		}
		battlelLog = req.RawLog

	default:
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Error: "analysisType must be one of: replayId, username, rawLog",
			Code:  "INVALID_REQUEST",
		})
		return
	}

	// Parse battle log
	parseStart := time.Now()
	battleSummary, err = analysis.ParseShowdownLog(battlelLog)
	parseTime := time.Since(parseStart).Milliseconds()

	if err != nil {
		s.logger.Infof("Failed to parse battle log: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Failed to parse battle log: " + err.Error(),
			Code:  "PARSE_ERROR",
		})
		return
	}

	// Store battle in database
	battleRecord := &db.Battle{
		Format:      battleSummary.Format,
		Timestamp:   battleSummary.Timestamp,
		DurationSec: battleSummary.Duration,
		Winner:      battleSummary.Winner,
		Player1ID:   battleSummary.Player1.Name,
		Player2ID:   battleSummary.Player2.Name,
		BattleLog:   battlelLog,
		IsPrivate:   req.IsPrivate,
	}

	// TODO: Store analysis results
	// TODO: Store key moments
	// battleID, err := s.db.StoreBattle(ctx, battleRecord)
	// if err != nil {
	// 	s.logger.Infof("Failed to store battle: %v", err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	json.NewEncoder(w).Encode(ErrorResponse{
	// 		Error: "Failed to store battle",
	// 		Code:  "INTERNAL_ERROR",
	// 	})
	// 	return
	// }

	analysisTime := time.Since(start).Milliseconds()

	s.logger.Infof("Successfully analyzed Showdown battle: %s", battleSummary.ID)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(AnalyzeResponse{
		Status:   "success",
		BattleID: battleRecord.ID,
		Data:     battleSummary,
		Metadata: &ResponseMetadata{
			ParseTimeMs:    int(parseTime),
			AnalysisTimeMs: int(analysisTime),
			Cached:         false,
		},
	})
}

// handleGetShowdownReplay handles GET /api/showdown/replays/{replayId} requests.
func (s *Server) handleGetShowdownReplay(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	battleID := chi.URLParam(r, "replayId")

	if battleID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Error: "replayId is required",
			Code:  "INVALID_REQUEST",
		})
		return
	}

	s.logger.Infof("Retrieving replay: %s", battleID)

	// TODO: Retrieve from database
	// battle, err := s.db.GetBattle(ctx, battleID)
	// if err != nil {
	// 	s.logger.Infof("Failed to retrieve battle: %v", err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	json.NewEncoder(w).Encode(ErrorResponse{
	// 		Error: "Internal server error",
	// 		Code:  "INTERNAL_ERROR",
	// 	})
	// 	return
	// }
	//
	// if battle == nil {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	json.NewEncoder(w).Encode(ErrorResponse{
	// 		Error: "Replay not found",
	// 		Code:  "NOT_FOUND",
	// 	})
	// 	return
	// }

	w.WriteHeader(http.StatusNotFound)
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Error: "Replay not found",
		Code:  "NOT_FOUND",
	})
}

// handleListShowdownReplays handles GET /api/showdown/replays requests.
func (s *Server) handleListShowdownReplays(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse query parameters
	username := r.URL.Query().Get("username")
	format := r.URL.Query().Get("format")
	isPrivateStr := r.URL.Query().Get("isPrivate")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	var isPrivate *bool
	if isPrivateStr != "" {
		val := isPrivateStr == "true"
		isPrivate = &val
	}

	limit := 10
	if limitStr != "" {
		if v, err := strconv.Atoi(limitStr); err == nil && v > 0 && v <= 100 {
			limit = v
		}
	}

	offset := 0
	if offsetStr != "" {
		if v, err := strconv.Atoi(offsetStr); err == nil && v >= 0 {
			offset = v
		}
	}

	s.logger.Infof("Listing replays: username=%s format=%s isPrivate=%v limit=%d offset=%d", username, format, isPrivate, limit, offset)

	// TODO: Query database
	// filter := &db.BattleFilter{
	// 	Format:    format,
	// 	IsPrivate: isPrivate,
	// }
	// battles, total, err := s.db.ListBattles(ctx, filter, limit, offset)
	// if err != nil {
	// 	s.logger.Infof("Failed to list battles: %v", err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	json.NewEncoder(w).Encode(ErrorResponse{
	// 		Error: "Internal server error",
	// 		Code:  "INTERNAL_ERROR",
	// 	})
	// 	return
	// }

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   []interface{}{},
		"pagination": map[string]int{
			"limit":  limit,
			"offset": offset,
			"total":  0,
		},
	})
}
