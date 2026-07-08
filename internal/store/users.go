package store

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/githubuseradmin/nova-dashboard/internal/models"
)

const userColumns = `id, email, name, role, password_hash, created_at`

// scanUser reads one user row. It scans role through a string to avoid
// depending on pgx codec behaviour for the named Role type.
func scanUser(row pgx.Row) (models.User, error) {
	var u models.User
	var role string
	if err := row.Scan(&u.ID, &u.Email, &u.Name, &role, &u.PasswordHash, &u.CreatedAt); err != nil {
		return models.User{}, err
	}
	u.Role = models.Role(role)
	return u, nil
}

// CreateUser inserts a user and returns the stored row.
func (s *Store) CreateUser(ctx context.Context, email, name string, role models.Role, passwordHash string) (models.User, error) {
	row := s.pool.QueryRow(ctx,
		`INSERT INTO users (email, name, role, password_hash)
		 VALUES ($1, $2, $3, $4)
		 RETURNING `+userColumns,
		email, name, string(role), passwordHash)
	return scanUser(row)
}

// GetUserByEmail looks up a user by email. Returns ErrNotFound if absent.
func (s *Store) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	row := s.pool.QueryRow(ctx, `SELECT `+userColumns+` FROM users WHERE email = $1`, email)
	u, err := scanUser(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.User{}, ErrNotFound
	}
	return u, err
}

// GetUserByID looks up a user by id. Returns ErrNotFound if absent.
func (s *Store) GetUserByID(ctx context.Context, id int64) (models.User, error) {
	row := s.pool.QueryRow(ctx, `SELECT `+userColumns+` FROM users WHERE id = $1`, id)
	u, err := scanUser(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.User{}, ErrNotFound
	}
	return u, err
}

// ListUsers returns all users, newest first.
func (s *Store) ListUsers(ctx context.Context) ([]models.User, error) {
	rows, err := s.pool.Query(ctx, `SELECT `+userColumns+` FROM users ORDER BY created_at DESC, id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

// UpdateUserRole changes a user's role and returns the updated row.
func (s *Store) UpdateUserRole(ctx context.Context, id int64, role models.Role) (models.User, error) {
	row := s.pool.QueryRow(ctx,
		`UPDATE users SET role = $2 WHERE id = $1 RETURNING `+userColumns,
		id, string(role))
	u, err := scanUser(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.User{}, ErrNotFound
	}
	return u, err
}

// CountUsers returns the total number of users.
func (s *Store) CountUsers(ctx context.Context) (int, error) {
	var n int
	err := s.pool.QueryRow(ctx, `SELECT count(*) FROM users`).Scan(&n)
	return n, err
}
