package httpapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dtsong/vgccorner/backend/internal/observability"
)

func TestAnalyzeShowdownRawLog(t *testing.T) {
	logger := observability.NewLogger()
	server := &Server{logger: logger}

	tests := []struct {
		name           string
		request        AnalyzeShowdownRequest
		expectedStatus int
		expectedError  string
	}{
		{
			name: "valid raw log analysis",
			request: AnalyzeShowdownRequest{
				AnalysisType: "rawLog",
				RawLog:       sampleShowdownLog(),
				IsPrivate:    false,
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name: "empty raw log returns error",
			request: AnalyzeShowdownRequest{
				AnalysisType: "rawLog",
				RawLog:       "",
				IsPrivate:    false,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "rawLog is required",
		},
		{
			name: "invalid analysis type returns error",
			request: AnalyzeShowdownRequest{
				AnalysisType: "invalid",
				RawLog:       sampleShowdownLog(),
				IsPrivate:    false,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "analysisType must be one of",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.request)
			req := httptest.NewRequest("POST", "/api/showdown/analyze", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			server.handleAnalyzeShowdown(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedError != "" {
				var resp ErrorResponse
				_ = json.NewDecoder(w.Body).Decode(&resp)
				if resp.Error == "" {
					t.Errorf("expected error containing %q, got none", tt.expectedError)
				}
			}
		})
	}
}

func TestAnalyzeShowdownByReplayID(t *testing.T) {
	logger := observability.NewLogger()
	server := &Server{logger: logger}

	tests := []struct {
		name           string
		replayID       string
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "empty replay id returns error",
			replayID:       "",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "INVALID_REQUEST",
		},
		{
			name:           "valid replay id returns not implemented",
			replayID:       "gen9vgc2025reghbo3-2481642254",
			expectedStatus: http.StatusNotImplemented,
			expectedCode:   "NOT_IMPLEMENTED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := AnalyzeShowdownRequest{
				AnalysisType: "replayId",
				ReplayID:     tt.replayID,
				IsPrivate:    false,
			}
			body, _ := json.Marshal(req)
			httpReq := httptest.NewRequest("POST", "/api/showdown/analyze", bytes.NewReader(body))
			httpReq.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			server.handleAnalyzeShowdown(w, httpReq)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedCode != "" {
				var resp ErrorResponse
				_ = json.NewDecoder(w.Body).Decode(&resp)
				if resp.Code != tt.expectedCode {
					t.Errorf("expected code %q, got %q", tt.expectedCode, resp.Code)
				}
			}
		})
	}
}

func TestAnalyzeShowdownByUsername(t *testing.T) {
	logger := observability.NewLogger()
	server := &Server{logger: logger}

	tests := []struct {
		name           string
		username       string
		format         string
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "missing username returns error",
			username:       "",
			format:         "gen9vgc2025reghbo3",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "INVALID_REQUEST",
		},
		{
			name:           "missing format returns error",
			username:       "Heliosan",
			format:         "",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "INVALID_REQUEST",
		},
		{
			name:           "valid username and format returns not implemented",
			username:       "Heliosan",
			format:         "gen9vgc2025reghbo3",
			expectedStatus: http.StatusNotImplemented,
			expectedCode:   "NOT_IMPLEMENTED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := AnalyzeShowdownRequest{
				AnalysisType: "username",
				Username:     tt.username,
				Format:       tt.format,
				IsPrivate:    false,
			}
			body, _ := json.Marshal(req)
			httpReq := httptest.NewRequest("POST", "/api/showdown/analyze", bytes.NewReader(body))
			httpReq.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			server.handleAnalyzeShowdown(w, httpReq)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedCode != "" {
				var resp ErrorResponse
				_ = json.NewDecoder(w.Body).Decode(&resp)
				if resp.Code != tt.expectedCode {
					t.Errorf("expected code %q, got %q", tt.expectedCode, resp.Code)
				}
			}
		})
	}
}

func TestAnalyzeShowdownInvalidJSON(t *testing.T) {
	logger := observability.NewLogger()
	server := &Server{logger: logger}

	req := httptest.NewRequest("POST", "/api/showdown/analyze", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.handleAnalyzeShowdown(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var resp ErrorResponse
	_ = json.NewDecoder(w.Body).Decode(&resp)
	if resp.Code != "INVALID_REQUEST" {
		t.Errorf("expected code INVALID_REQUEST, got %q", resp.Code)
	}
}

func TestGetShowdownReplay(t *testing.T) {
	logger := observability.NewLogger()
	router := NewRouter(logger)

	tests := []struct {
		name           string
		replayID       string
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "valid replay id returns not found",
			replayID:       "gen9vgc2025reghbo3-2481642254",
			expectedStatus: http.StatusNotFound,
			expectedCode:   "NOT_FOUND",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/showdown/replays/"+tt.replayID, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedCode != "" {
				var resp ErrorResponse
				_ = json.NewDecoder(w.Body).Decode(&resp)
				if resp.Code != tt.expectedCode {
					t.Errorf("expected code %q, got %q", tt.expectedCode, resp.Code)
				}
			}
		})
	}
}

func TestListShowdownReplays(t *testing.T) {
	logger := observability.NewLogger()
	server := &Server{logger: logger}

	tests := []struct {
		name           string
		query          string
		expectedStatus int
	}{
		{
			name:           "no filters returns success",
			query:          "",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "with username filter",
			query:          "?username=Heliosan",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "with format filter",
			query:          "?format=gen9vgc2025reghbo3",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "with limit and offset",
			query:          "?limit=20&offset=10",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "with invalid limit",
			query:          "?limit=invalid",
			expectedStatus: http.StatusOK, // Should use default
		},
		{
			name:           "with limit over max",
			query:          "?limit=200",
			expectedStatus: http.StatusOK, // Should cap at 100
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/showdown/replays"+tt.query, nil)
			w := httptest.NewRecorder()

			server.handleListShowdownReplays(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var resp map[string]interface{}
			_ = json.NewDecoder(w.Body).Decode(&resp)

			if status, ok := resp["status"]; !ok || status != "success" {
				t.Errorf("expected status 'success' in response")
			}
		})
	}
}

func TestAnalyzeTCGLive(t *testing.T) {
	logger := observability.NewLogger()
	server := &Server{logger: logger}

	tests := []struct {
		name           string
		gameExport     string
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "valid game export returns not implemented",
			gameExport:     "some game export data",
			expectedStatus: http.StatusNotImplemented,
			expectedCode:   "NOT_IMPLEMENTED",
		},
		{
			name:           "empty game export returns error",
			gameExport:     "",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "INVALID_REQUEST",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := AnalyzeTCGLiveRequest{
				GameExport: tt.gameExport,
				IsPrivate:  true,
			}
			body, _ := json.Marshal(req)
			httpReq := httptest.NewRequest("POST", "/api/tcglive/analyze", bytes.NewReader(body))
			httpReq.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			server.handleAnalyzeTCGLive(w, httpReq)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var resp ErrorResponse
			_ = json.NewDecoder(w.Body).Decode(&resp)
			if resp.Code != tt.expectedCode {
				t.Errorf("expected code %q, got %q", tt.expectedCode, resp.Code)
			}
		})
	}
}

// Additional comprehensive tests for edge cases and API validation

func TestAnalyzeShowdownEdgeCases(t *testing.T) {
	logger := observability.NewLogger()
	server := &Server{logger: logger}

	tests := []struct {
		name           string
		request        AnalyzeShowdownRequest
		expectedStatus int
	}{
		{
			name: "request with private flag true",
			request: AnalyzeShowdownRequest{
				AnalysisType: "rawLog",
				RawLog:       sampleShowdownLog(),
				IsPrivate:    true,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "request with very long log",
			request: AnalyzeShowdownRequest{
				AnalysisType: "rawLog",
				RawLog:       generateLongLog(),
				IsPrivate:    false,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "request with special characters in log",
			request: AnalyzeShowdownRequest{
				AnalysisType: "rawLog",
				RawLog:       `|j|☆TestPlayer™
|j|☆OtherPlayer
|start
|win|TestPlayer™`,
				IsPrivate: false,
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.request)
			req := httptest.NewRequest("POST", "/api/showdown/analyze", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			server.handleAnalyzeShowdown(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d, body: %s", tt.expectedStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestResponseDataStructure(t *testing.T) {
	logger := observability.NewLogger()
	server := &Server{logger: logger}

	req := AnalyzeShowdownRequest{
		AnalysisType: "rawLog",
		RawLog:       sampleShowdownLog(),
		IsPrivate:    false,
	}
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/api/showdown/analyze", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.handleAnalyzeShowdown(w, httpReq)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp AnalyzeResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Validate response structure
	if resp.Status != "success" {
		t.Errorf("expected status 'success', got %q", resp.Status)
	}

	if resp.BattleID == "" {
		t.Error("expected battleId to be set")
	}

	if resp.Data == nil {
		t.Fatal("expected data field")
	}

	if resp.Data.ID == "" {
		t.Error("expected data.id to be set")
	}

	if resp.Data.Format == "" {
		t.Error("expected data.format to be set")
	}

	if resp.Data.Player1.Name == "" || resp.Data.Player2.Name == "" {
		t.Error("expected both player names to be set")
	}

	if len(resp.Data.Turns) == 0 {
		t.Error("expected turns to be populated")
	}

	if resp.Metadata == nil {
		t.Fatal("expected metadata field")
	}

	if resp.Metadata.ParseTimeMs < 0 {
		t.Error("expected parseTimeMs to be non-negative")
	}

	if resp.Metadata.AnalysisTimeMs < 0 {
		t.Error("expected analysisTimeMs to be non-negative")
	}
}

func TestListShowdownReplaysResponseStructure(t *testing.T) {
	logger := observability.NewLogger()
	server := &Server{logger: logger}

	req := httptest.NewRequest("GET", "/api/showdown/replays?limit=10", nil)
	w := httptest.NewRecorder()

	server.handleListShowdownReplays(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if status, ok := resp["status"]; !ok || status != "success" {
		t.Error("expected status 'success' in response")
	}

	if _, ok := resp["data"]; !ok {
		t.Error("expected data field in response")
	}
}

func TestAnalyzeShowdownResponseContentType(t *testing.T) {
	logger := observability.NewLogger()
	server := &Server{logger: logger}

	req := AnalyzeShowdownRequest{
		AnalysisType: "rawLog",
		RawLog:       sampleShowdownLog(),
	}
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/api/showdown/analyze", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.handleAnalyzeShowdown(w, httpReq)

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", contentType)
	}
}

func TestListReplaysLimitBounds(t *testing.T) {
	logger := observability.NewLogger()
	server := &Server{logger: logger}

	tests := []struct {
		name     string
		limit    string
		expectOK bool
	}{
		{"limit 0", "?limit=0", true},
		{"limit 50", "?limit=50", true},
		{"limit 100", "?limit=100", true},
		{"limit 200", "?limit=200", true}, // Should cap at 100
		{"limit -1", "?limit=-1", true},   // Should use default
		{"limit abc", "?limit=abc", true}, // Should use default
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/showdown/replays"+tt.limit, nil)
			w := httptest.NewRecorder()

			server.handleListShowdownReplays(w, req)

			if tt.expectOK && w.Code != http.StatusOK {
				t.Errorf("expected status 200, got %d", w.Code)
			}
		})
	}
}

func TestAnalyzeShowdownAllAnalysisTypesValidation(t *testing.T) {
	logger := observability.NewLogger()
	server := &Server{logger: logger}

	tests := []struct {
		name           string
		analysisType   string
		replayID       string
		username       string
		format         string
		rawLog         string
		expectedStatus int
	}{
		{
			name:           "rawLog with raw log",
			analysisType:   "rawLog",
			rawLog:         sampleShowdownLog(),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "rawLog without raw log",
			analysisType:   "rawLog",
			rawLog:         "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "replayId with id",
			analysisType:   "replayId",
			replayID:       "gen9vgc2025reghbo3-2481642254",
			expectedStatus: http.StatusNotImplemented,
		},
		{
			name:           "replayId without id",
			analysisType:   "replayId",
			replayID:       "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "username with username and format",
			analysisType:   "username",
			username:       "TestPlayer",
			format:         "gen9vgc2025reghbo3",
			expectedStatus: http.StatusNotImplemented,
		},
		{
			name:           "username without username",
			analysisType:   "username",
			username:       "",
			format:         "gen9vgc2025reghbo3",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "username without format",
			analysisType:   "username",
			username:       "TestPlayer",
			format:         "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid analysis type",
			analysisType:   "invalid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := AnalyzeShowdownRequest{
				AnalysisType: tt.analysisType,
				ReplayID:     tt.replayID,
				Username:     tt.username,
				Format:       tt.format,
				RawLog:       tt.rawLog,
			}
			body, _ := json.Marshal(req)
			httpReq := httptest.NewRequest("POST", "/api/showdown/analyze", bytes.NewReader(body))
			httpReq.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			server.handleAnalyzeShowdown(w, httpReq)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

// Helper functions

func generateLongLog() string {
	log := `|j|☆Player1
|j|☆Player2
|player|p1|Player1|test|1500
|player|p2|Player2|test|1500
|tier|[Gen 9] VGC 2025 Reg H (Bo3)
|poke|p1|Poke1, L50|
|poke|p2|Poke2, L50|
|teamsize|p1|1
|teamsize|p2|1
|start
`

	// Generate multiple turns
	for i := 1; i <= 20; i++ {
		log += fmt.Sprintf(`|turn|%d
|move|p1a: Poke1|Tackle|p2a: Poke2
|-damage|p2a: Poke2|%d/100
|move|p2a: Poke2|Tackle|p1a: Poke1
|-damage|p1a: Poke1|%d/100
|upkeep
`, i, 100-i*5, 100-i*5)
	}

	log += `|win|Player1`
	return log
}
