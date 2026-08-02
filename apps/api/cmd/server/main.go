package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/config"
	"github.com/magicvr/schema-ui-core/apps/api/internal/handler"
	"github.com/magicvr/schema-ui-core/apps/api/internal/server"
	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
	"github.com/magicvr/schema-ui-core/apps/api/pkg/version"
)

const bcryptCost = 10

func main() {
	cfg := config.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.LogLevel(),
	}))
	slog.SetDefault(logger)

	secret, err := resolveJWTSecret(cfg, logger)
	if err != nil {
		logger.Error("startup failed", "err", err)
		os.Exit(1)
	}

	seedHash, err := resolveSeedHash(cfg, logger)
	if err != nil {
		logger.Error("startup failed", "err", err)
		os.Exit(1)
	}

	st, err := store.Open(cfg.DBPath, "admin", seedHash, true)
	if err != nil {
		logger.Error("open auth store", "err", err)
		os.Exit(1)
	}
	defer st.Close()

	authenticator := auth.New(
		[]byte(secret),
		cfg.AuthAccessTTL,
		cfg.AuthRefreshTTL,
		st,
		cfg.AuthDevSessionEnabled,
	)

	mux := http.NewServeMux()
	handler.Register(mux, authenticator, st)

	srv := server.New(cfg, mux, logger)

	go func() {
		logger.Info("server starting",
			"addr", cfg.HTTPAddr,
			"version", version.Version,
			"commit", version.Commit,
			"env", cfg.AppEnv,
			"dev_session", cfg.AuthDevSessionEnabled,
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", "err", err)
			os.Exit(1)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	logger.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown", "err", err)
		os.Exit(1)
	}
	fmt.Println("bye")
}

// resolveJWTSecret returns the signing secret, failing closed outside
// development when AUTH_JWT_SECRET is missing (GOAL-005 D-004).
func resolveJWTSecret(cfg *config.Config, logger *slog.Logger) (string, error) {
	if cfg.AuthJWTSecret != "" {
		return cfg.AuthJWTSecret, nil
	}
	if cfg.AppEnv == "development" {
		logger.Warn("AUTH_JWT_SECRET not set; using an insecure development signing key")
		return "dev-only-insecure-jwt-secret-change-me", nil
	}
	return "", fmt.Errorf("AUTH_JWT_SECRET must be set in non-development environment")
}

// resolveSeedHash computes the bcrypt hash for the bootstrap admin password,
// failing closed in production when ADMIN_INITIAL_PASSWORD is missing. In
// development the documented fallback is "admin".
func resolveSeedHash(cfg *config.Config, logger *slog.Logger) (string, error) {
	seed := cfg.AdminInitialPassword
	if seed == "" {
		if cfg.AppEnv == "development" {
			seed = "admin"
		} else {
			return "", fmt.Errorf("ADMIN_INITIAL_PASSWORD must be set to seed the initial admin user")
		}
	}
	hash, err := auth.HashPassword(seed, bcryptCost)
	if err != nil {
		return "", fmt.Errorf("hash seed password: %w", err)
	}
	return hash, nil
}
