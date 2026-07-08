package http

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/githubuseradmin/nova-dashboard/internal/auth"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// handleLogin verifies credentials and starts a session.
func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" || req.Password == "" {
		respondError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	user, err := s.store.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		// Same response whether the email exists or not (no user enumeration).
		respondError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}
	ok, err := auth.VerifyPassword(req.Password, user.PasswordHash)
	if err != nil || !ok {
		respondError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	token, err := auth.NewToken()
	if err != nil {
		s.serverError(w, "generate session token", err)
		return
	}
	if err := s.store.CreateSession(r.Context(), token, user.ID, time.Now().Add(auth.SessionDuration)); err != nil {
		s.serverError(w, "create session", err)
		return
	}

	auth.SetSessionCookie(w, token, s.secure)
	respondJSON(w, http.StatusOK, user)
}

// handleLogout destroys the current session.
func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(auth.SessionCookieName); err == nil && cookie.Value != "" {
		_ = s.store.DeleteSession(r.Context(), cookie.Value)
	}
	auth.ClearSessionCookie(w, s.secure)
	respondJSON(w, http.StatusOK, map[string]string{"status": "logged out"})
}

// handleMe returns the current user, or 401 if not signed in.
func (s *Server) handleMe(w http.ResponseWriter, r *http.Request) {
	user, ok := userFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "not authenticated")
		return
	}
	respondJSON(w, http.StatusOK, user)
}

// handleCSRF issues a fresh CSRF token (cookie + body) for the SPA to echo back.
func (s *Server) handleCSRF(w http.ResponseWriter, r *http.Request) {
	token, err := auth.NewCSRFToken()
	if err != nil {
		s.serverError(w, "generate csrf token", err)
		return
	}
	auth.SetCSRFCookie(w, token, s.secure)
	respondJSON(w, http.StatusOK, map[string]string{"csrfToken": token})
}
