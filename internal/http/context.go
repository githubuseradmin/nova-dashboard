package http

import (
	"context"

	"github.com/githubuseradmin/nova-dashboard/internal/models"
)

// contextKey is an unexported type so our context keys never collide with
// keys set by other packages.
type contextKey int

const userContextKey contextKey = iota

// withUser returns a copy of ctx carrying the authenticated user.
func withUser(ctx context.Context, u models.User) context.Context {
	return context.WithValue(ctx, userContextKey, u)
}

// userFromContext returns the authenticated user, or false if the request is
// not authenticated.
func userFromContext(ctx context.Context) (models.User, bool) {
	u, ok := ctx.Value(userContextKey).(models.User)
	return u, ok
}
