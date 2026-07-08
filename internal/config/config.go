// Package config loads runtime configuration from environment variables,
// applying sane defaults so the app runs out of the box in development.
package config

import (
	"log/slog"
	"os"
	"strings"
)

// Config holds all runtime configuration, resolved once at startup.
type Config struct {
	Env         string // "development" | "production"
	Host        string
	Port        string
	DatabaseURL string
	LogLevelStr string

	// Static asset directories served by the Go server (relative to CWD).
	SiteDir string // marketing landing
	AppDir  string // built Svelte SPA

	// Initial admin, seeded on first run when the users table is empty.
	AdminEmail    string
	AdminPassword string
}

// Load reads configuration from the environment (NOVA_* keys) with defaults.
func Load() Config {
	return Config{
		Env:           env("NOVA_ENV", "development"),
		Host:          env("NOVA_HOST", "0.0.0.0"),
		Port:          env("NOVA_PORT", "8080"),
		DatabaseURL:   env("NOVA_DATABASE_URL", "postgres://nova:nova@localhost:5432/nova?sslmode=disable"),
		LogLevelStr:   env("NOVA_LOG_LEVEL", "info"),
		SiteDir:       env("NOVA_SITE_DIR", "web/site"),
		AppDir:        env("NOVA_APP_DIR", "web/app/dist"),
		AdminEmail:    env("NOVA_ADMIN_EMAIL", "admin@nova.local"),
		AdminPassword: env("NOVA_ADMIN_PASSWORD", "admin12345"),
	}
}

// Addr is the host:port the HTTP server binds to.
func (c Config) Addr() string { return c.Host + ":" + c.Port }

// IsProduction reports whether the app runs in production mode.
func (c Config) IsProduction() bool { return c.Env == "production" }

// LogLevel maps the configured level string to slog.Level (defaults to info).
func (c Config) LogLevel() slog.Level {
	switch strings.ToLower(c.LogLevelStr) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// env returns the value of key, or fallback when unset or empty.
func env(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
