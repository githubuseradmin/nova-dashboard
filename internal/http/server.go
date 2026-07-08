// Package http is the transport layer: it wires the router, middleware and
// handlers. Handlers depend only on the Server's fields (logger, store),
// never on globals — keeping the layer explicit and testable.
package http

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/githubuseradmin/nova-dashboard/internal/store"
)

// Server holds the dependencies shared by HTTP handlers.
type Server struct {
	logger  *slog.Logger
	store   *store.Store
	secure  bool   // whether auth cookies carry the Secure flag (true in production)
	siteDir string // marketing landing directory
	appDir  string // built SPA directory
}

// NewServer constructs a Server with its dependencies.
func NewServer(logger *slog.Logger, st *store.Store, secure bool, siteDir, appDir string) *Server {
	return &Server{logger: logger, store: st, secure: secure, siteDir: siteDir, appDir: appDir}
}

// Handler builds the router with the middleware chain and all routes mounted.
func (s *Server) Handler() http.Handler {
	r := chi.NewRouter()

	// Base middleware, outermost first.
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(requestLogger(s.logger))
	r.Use(middleware.Recoverer)

	r.Get("/healthz", s.handleHealth)

	// JSON API under /api. Session loading is scoped here so static asset
	// requests never touch the database.
	r.Route("/api", func(r chi.Router) {
		r.Use(s.authenticate)

		r.Get("/health", s.handleHealth)
		r.Get("/csrf", s.handleCSRF)

		r.Route("/auth", func(r chi.Router) {
			r.With(s.requireCSRF).Post("/login", s.handleLogin)
			r.With(s.requireCSRF).Post("/logout", s.handleLogout)
			r.Get("/me", s.handleMe)
		})

		r.Group(func(r chi.Router) {
			r.Use(s.requireAuth)

			r.Get("/dashboard", s.handleDashboard)

			r.Route("/admin", func(r chi.Router) {
				r.Use(s.requireAdmin)
				r.Get("/users", s.handleListUsers)
				r.With(s.requireCSRF).Put("/users/{id}/role", s.handleUpdateUserRole)
			})
		})
	})

	// Static assets: the built SPA under /app, the marketing landing at /.
	r.Handle("/app", http.RedirectHandler("/app/", http.StatusMovedPermanently))
	r.Handle("/app/*", http.StripPrefix("/app/", http.FileServer(http.Dir(s.appDir))))
	r.Handle("/*", http.FileServer(http.Dir(s.siteDir)))

	return r
}

// handleHealth is a minimal liveness/readiness probe.
func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// serverError logs an internal error and returns a generic 500 to the client.
func (s *Server) serverError(w http.ResponseWriter, msg string, err error) {
	s.logger.Error(msg, "err", err)
	respondError(w, http.StatusInternalServerError, "internal server error")
}
