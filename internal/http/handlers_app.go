package http

import "net/http"

// handleDashboard returns data for the signed-in user's dashboard.
// This is a placeholder payload — replace with real metrics as the app grows.
func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	user, _ := userFromContext(r.Context())

	resp := map[string]any{
		"user": user,
		"greeting": "Welcome back, " + firstNonEmpty(user.Name, user.Email),
		"stats": []map[string]any{
			{"label": "Account", "value": string(user.Role)},
			{"label": "Member since", "value": user.CreatedAt.Format("2 Jan 2006")},
		},
	}
	respondJSON(w, http.StatusOK, resp)
}

// firstNonEmpty returns a if it is non-empty, otherwise b.
func firstNonEmpty(a, b string) string {
	if a != "" {
		return a
	}
	return b
}
