package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/githubuseradmin/nova-dashboard/internal/models"
)

// CreateSession stores a session token for a user until expiresAt.
func (s *Store) CreateSession(ctx context.Context, token string, userID int64, expiresAt time.Time) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO sessions (token, user_id, expires_at) VALUES ($1, $2, $3)`,
		token, userID, expiresAt)
	return err
}

// UserBySession returns the user owning a non-expired session token.
// Returns ErrNotFound if the token is unknown or has expired.
func (s *Store) UserBySession(ctx context.Context, token string) (models.User, error) {
	const q = `
		SELECT u.id, u.email, u.name, u.role, u.password_hash, u.created_at
		FROM sessions s
		JOIN users u ON u.id = s.user_id
		WHERE s.token = $1 AND s.expires_at > now()`

	var u models.User
	var role string
	err := s.pool.QueryRow(ctx, q, token).
		Scan(&u.ID, &u.Email, &u.Name, &role, &u.PasswordHash, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.User{}, ErrNotFound
	}
	if err != nil {
		return models.User{}, err
	}
	u.Role = models.Role(role)
	return u, nil
}

// DeleteSession removes a session (used on logout).
func (s *Store) DeleteSession(ctx context.Context, token string) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM sessions WHERE token = $1`, token)
	return err
}

// DeleteExpiredSessions purges expired sessions and returns the count removed.
func (s *Store) DeleteExpiredSessions(ctx context.Context) (int64, error) {
	tag, err := s.pool.Exec(ctx, `DELETE FROM sessions WHERE expires_at <= now()`)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}
