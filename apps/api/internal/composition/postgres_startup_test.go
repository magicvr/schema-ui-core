package composition

import (
	"context"
	"database/sql"
	"io"
	"log/slog"
	"net/url"
	"path/filepath"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/config"

	"github.com/magicvr/schema-ui-core/apps/api/internal/pgtest"
)

// TestCompositionPostgresStartup proves R4 S4: a postgres DSN boots the full
// application — compiled dual-dialect catalog applied, admin seeded,
// system data reconciled, and every module Start+Ready gate (store Ping +
// SystemDataReady) passes — through the kernel.Store interface. Gated by
// SCHEMA_UI_R2_PG_DSN (no PG = skip).
func TestCompositionPostgresStartup(t *testing.T) {
	dsn := pgtest.DSN()
	if dsn == "" {
		t.Skip("postgres test env not set (PG_TEST_*); skipping postgres composition startup")
	}
	ctx := context.Background()
	u, err := url.Parse(dsn)
	if err != nil {
		t.Fatal(err)
	}
	const dbName = "r4s4"
	adminDSN := u.String()
	u.Path = "/" + dbName
	pgDSN := u.String()

	admin, err := sql.Open("pgx", adminDSN)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = admin.ExecContext(context.Background(), `DROP DATABASE IF EXISTS `+dbName+` WITH (FORCE)`)
		_ = admin.Close()
	})
	if _, err := admin.ExecContext(ctx, `DROP DATABASE IF EXISTS `+dbName+` WITH (FORCE)`); err != nil {
		t.Fatalf("drop prior scratch db: %v", err)
	}
	if _, err := admin.ExecContext(ctx, `CREATE DATABASE `+dbName); err != nil {
		t.Fatalf("create scratch db: %v", err)
	}

	cfg := &config.Config{
		AppName:      "test",
		AppEnv:       "development",
		HTTPAddr:     "127.0.0.1:0",
		DBDialect:    "postgres",
		DBPath:       filepath.Join(t.TempDir(), "r4s4.db"),
		DBDSN:        pgDSN,
		ProfileName:  "admin",
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		IdleTimeout:  time.Second,
	}
	app, err := NewApp(cfg, "test-secret", "hash", slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatalf("build composition root: %v", err)
	}
	startCtx, startCancel := context.WithTimeout(ctx, 60*time.Second)
	defer startCancel()
	if err := app.Start(startCtx); err != nil {
		t.Fatalf("postgres composition startup failed (apply+bootstrap+reconcile+ready gates): %v", err)
	}
	// Start success implies every module Start+Ready passed, including
	// core.auth-session (store Ping + SystemDataReady after Reconcile) — the
	// R4 postgres runtime evidence (module-gate-green, readyz-equivalent).
	t.Log("postgres composition started; modules Start+Ready green (readyz gate path)")
	stopCtx, stopCancel := context.WithTimeout(ctx, 15*time.Second)
	defer stopCancel()
	if err := app.Stop(stopCtx); err != nil {
		t.Fatalf("stop: %v", err)
	}
}
