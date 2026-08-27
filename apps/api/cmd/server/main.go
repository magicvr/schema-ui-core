package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/composition"
	"github.com/magicvr/schema-ui-core/apps/api/internal/config"
)

const bcryptCost = 10

func main() {
	cfg := config.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.LogLevel(),
	}))
	slog.SetDefault(logger)

	// Hard production guard (GOAL-008 A-005 F-002): the static dev session
	// fallback is only legal in local development. Error level always surfaces,
	// so this stays visible regardless of LOG_LEVEL.
	if err := cfg.ValidateProd(); err != nil {
		logger.Error("startup failed", "err", err)
		os.Exit(1)
	}

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

	app, err := composition.NewApp(cfg, secret, seedHash, logger)
	if err != nil {
		logger.Error("build composition root", "err", err)
		os.Exit(1)
	}
	startCtx, startCancel := context.WithTimeout(context.Background(), 15*time.Second)
	if err := app.Start(startCtx); err != nil {
		startCancel()
		logger.Error("startup failed", "err", err)
		os.Exit(1)
	}
	startCancel()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	// VP-021 contract v0.1.0 §1/§7: structured lifecycle events.
	logger.Info("shutdown.starting")
	// VP-021 contract v0.1.0 §6: drain budget = http.shutdown_timeout
	// (default 10s; invalid values fail closed at startup).
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.HTTPShutdownTimeout)
	defer cancel()
	if err := app.Stop(shutdownCtx); err != nil {
		if shutdownCtx.Err() != nil {
			logger.Error("shutdown.timeout", "err", err)
		} else {
			logger.Error("shutdown.error", "err", err)
		}
		os.Exit(1)
	}
	logger.Info("shutdown.complete")
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
