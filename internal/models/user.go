// Package models holds the core domain types shared across the application.
package models

import "time"

// Role is a user's access level.
type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

// Valid reports whether r is a known role.
func (r Role) Valid() bool { return r == RoleUser || r == RoleAdmin }

// User is an account that can sign in.
type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	Role         Role      `json:"role"`
	PasswordHash string    `json:"-"` // never serialized to clients
	CreatedAt    time.Time `json:"createdAt"`
}

// IsAdmin reports whether the user has admin privileges.
func (u User) IsAdmin() bool { return u.Role == RoleAdmin }
