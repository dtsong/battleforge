package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dtsong/battleforge/backend/internal/observability"
)

func TestHealthCheck(t *testing.T) {
	logger := observability.NewLogger()
	router := NewRouter(logger)

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "healthz returns ok",
			method:         "GET",
			path:           "/healthz",
			expectedStatus: http.StatusOK,
			expectedBody:   "ok",
		},
		{
			name:           "healthz only accepts GET",
			method:         "POST",
			path:           "/healthz",
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedBody != "" && w.Body.String() != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestRouterEndpointsExist(t *testing.T) {
	logger := observability.NewLogger()
	router := NewRouter(logger)

	// These tests just verify the routes are registered and return a response
	tests := []struct {
		name          string
		method        string
		path          string
		allowNotFound bool // for handlers that return 404 when resource doesn't exist
	}{
		{"healthz GET", "GET", "/healthz", false},
		{"showdown analyze POST", "POST", "/api/showdown/analyze", false},
		{"showdown list GET", "GET", "/api/showdown/replays", false},
		{"showdown get GET", "GET", "/api/showdown/replays/test-id", true},
		{"tcglive analyze POST", "POST", "/api/tcglive/analyze", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			// Should not be 404 unless allowNotFound is true
			if w.Code == http.StatusNotFound && !tt.allowNotFound {
				t.Errorf("route %s %s not found (404)", tt.method, tt.path)
			}
		})
	}
}
