// Command server is the nova HTTP API entrypoint.
//
// It loads configuration, connects to Postgres, applies migrations, seeds an
// initial admin, then starts the HTTP server and shuts it down gracefully on
// SIGINT/SIGTERM.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/githubuseradmin/nova-dashboard/internal/auth"
	"github.com/githubuseradmin/nova-dashboard/internal/config"
	nhttp "github.com/githubuseradmin/nova-dashboard/internal/http"
	"github.com/githubuseradmin/nova-dashboard/internal/models"
	"github.com/githubuseradmin/nova-dashboard/internal/store"
)

func main() {
	if err := run(); err != nil {
		slog.Error("server exited with error", "err", err)
		os.Exit(1)
	}
}

func run() error {
	cfg := config.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel()}))
	slog.SetDefault(logger)

	ctx := context.Background()

	st, err := store.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("store: %w", err)
	}
	defer st.Close()

	if err := st.Migrate(ctx); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	if err := seedAdmin(ctx, st, cfg, logger); err != nil {
		return fmt.Errorf("seed admin: %w", err)
	}

	server := nhttp.NewServer(logger, st, cfg.IsProduction(), cfg.SiteDir, cfg.AppDir)

	httpServer := &http.Server{
		Addr:              cfg.Addr(),
		Handler:           server.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Run the server in the background so we can wait for signals.
	serverErr := make(chan error, 1)
	go func() {
		logger.Info("http server listening", "addr", httpServer.Addr, "env", cfg.Env)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	// Block until the server fails or an interrupt/termination signal arrives.
	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-serverErr:
		return err
	case <-signalCtx.Done():
		logger.Info("shutdown signal received, draining connections")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return httpServer.Shutdown(shutdownCtx)
}

// seedAdmin creates an initial admin account when the users table is empty,
// using the configured admin email/password.
func seedAdmin(ctx context.Context, st *store.Store, cfg config.Config, logger *slog.Logger) error {
	n, err := st.CountUsers(ctx)
	if err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	hash, err := auth.HashPassword(cfg.AdminPassword)
	if err != nil {
		return err
	}
	if _, err := st.CreateUser(ctx, cfg.AdminEmail, "Admin", models.RoleAdmin, hash); err != nil {
		return err
	}
	logger.Info("seeded initial admin account", "email", cfg.AdminEmail)
	return nil
}
