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

	// A-002 R-001 (workspace-014 GOAL-003): readyz already covers the
	// configured S3 backend (HeadBucket probe), but the three file families
	// still use the local disk adapter until the R3 call-site wiring lands.
	// Keep that window explicit for operators.
	if cfg.ObjectsDriver == "s3" {
		logger.Warn("storage.objects.driver=s3: backend probe active in readyz, but file families still use the local disk adapter until workspace-014 R3 wiring")
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

	logger.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.Stop(shutdownCtx); err != nil {
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
