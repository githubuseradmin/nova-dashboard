package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/githubuseradmin/nova-dashboard/internal/models"
)

// handleListUsers returns all users (admin only).
func (s *Server) handleListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := s.store.ListUsers(r.Context())
	if err != nil {
		s.serverError(w, "list users", err)
		return
	}
	respondJSON(w, http.StatusOK, users)
}

type updateRoleRequest struct {
	Role models.Role `json:"role"`
}

// handleUpdateUserRole changes a user's role (admin only).
func (s *Server) handleUpdateUserRole(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	var req updateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if !req.Role.Valid() {
		respondError(w, http.StatusBadRequest, "invalid role")
		return
	}

	// Guard against an admin removing their own admin role and locking out.
	if current, ok := userFromContext(r.Context()); ok && current.ID == id && req.Role != models.RoleAdmin {
		respondError(w, http.StatusBadRequest, "you cannot remove your own admin role")
		return
	}

	user, err := s.store.UpdateUserRole(r.Context(), id, req.Role)
	if err != nil {
		respondError(w, http.StatusNotFound, "user not found")
		return
	}
	respondJSON(w, http.StatusOK, user)
}
