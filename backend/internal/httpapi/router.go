package httpapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/dtsong/battleforge/backend/internal/observability"
)

type Server struct {
	logger *observability.Logger
	// db     *db.Database // Will be added when DB integration is complete
}

func NewRouter(logger *observability.Logger) http.Handler {
	s := &Server{logger: logger}

	r := chi.NewRouter()

	// Health check endpoint
	r.Get("/healthz", s.handleHealth)

	// Showdown analysis endpoints
	r.Post("/api/showdown/analyze", s.handleAnalyzeShowdown)
	r.Get("/api/showdown/replays", s.handleListShowdownReplays)
	r.Get("/api/showdown/replays/{replayId}", s.handleGetShowdownReplay)

	// TCG Live endpoint (planned)
	r.Post("/api/tcglive/analyze", s.handleAnalyzeTCGLive)

	return r
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
