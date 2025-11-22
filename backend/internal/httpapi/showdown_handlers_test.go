package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dtsong/battleforge/backend/internal/observability"
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
