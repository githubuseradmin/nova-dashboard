package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// respondJSON writes v as a JSON response with the given status code.
// A nil body writes just the status (useful for 204 No Content).
func respondJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if v == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("failed to encode JSON response", "err", err)
	}
}

// errorResponse is the single envelope used for all API errors.
type errorResponse struct {
	Error string `json:"error"`
}

// respondError writes a JSON error envelope with the given status code.
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, errorResponse{Error: message})
}
