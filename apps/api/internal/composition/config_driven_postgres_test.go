package composition

import (
	"context"
	"database/sql"
	"io"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/config"

	"github.com/magicvr/schema-ui-core/apps/api/internal/pgtest"
)

// TestCompositionPostgresConfigDriven proves the exact usage scenario: a
// config.yaml selects dialect/connection params (host/port/name/user/sslmode)
// and the database PASSWORD is supplied only via the DB_PASSWORD env var
// (configs/.env). The loaded Config then boots the full app on postgres.
func TestCompositionPostgresConfigDriven(t *testing.T) {
	dsn := pgtest.DSN()
	if dsn == "" {
		t.Skip("postgres test env not set (PG_TEST_*); skipping config-driven postgres startup")
	}
	ctx := context.Background()
	u, err := url.Parse(dsn)
	if err != nil {
		t.Fatal(err)
	}
	pw, _ := u.User.Password()
	hostPort := u.Host // host:port
	host := hostPort
	port := "5432"
	if i := strings.LastIndex(hostPort, ":"); i >= 0 {
		host = hostPort[:i]
		port = hostPort[i+1:]
	}
	user := u.User.Username()
	sslmode := u.Query().Get("sslmode")
	if sslmode == "" {
		sslmode = "disable"
	}

	// Create a dedicated scratch database on the target server.
	admin := u.String()
	const dbName = "r5cfg"
	adm, err := sql.Open("pgx", admin)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = adm.ExecContext(context.Background(), `DROP DATABASE IF EXISTS `+dbName+` WITH (FORCE)`)
		_ = adm.Close()
	})
	c := context.Background()
	if _, err := adm.ExecContext(c, `DROP DATABASE IF EXISTS `+dbName+` WITH (FORCE)`); err != nil {
		t.Fatalf("drop prior scratch db: %v", err)
	}
	if _, err := adm.ExecContext(c, `CREATE DATABASE `+dbName); err != nil {
		t.Fatalf("create scratch db: %v", err)
	}

	// config.yaml declares everything except the password secret.
	yamlText := strings.Join([]string{
		"app:",
		"  env: development",
		"  profile: admin",
		"http:",
		"  addr: 127.0.0.1:0",
		"db:",
		"  dialect: postgres",
		"  path: " + filepath.Join(t.TempDir(), "r5cfg.db"),
		"  host: " + host,
		"  port: " + port,
		"  name: " + dbName,
		"  user: " + user,
		"  password: ${DB_PASSWORD:-}",
		"  sslmode: " + sslmode,
		"",
	}, "\n")
	cfgPath := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(cfgPath, []byte(yamlText), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("CONFIG_FILE", cfgPath)
	t.Setenv("DB_PASSWORD", pw) // secret comes from env, not the yaml literal
	t.Setenv("DB_DSN", "")      // ensure exploded-params path is exercised
	t.Setenv("DB_DIALECT", "")  // dialect comes from yaml
	t.Setenv("DB_HOST", "")     // host/port/name/user come from yaml
	t.Setenv("DB_PORT", "")
	t.Setenv("DB_NAME", "")
	t.Setenv("DB_USER", "")
	t.Setenv("DB_SSLMODE", "")

	cfg := config.Load()
	if cfg.LoadError != nil {
		t.Fatalf("config.Load: %v", cfg.LoadError)
	}
	if cfg.DBDialect != "postgres" {
		t.Fatalf("DBDialect = %q, want postgres (from yaml)", cfg.DBDialect)
	}
	if !strings.Contains(cfg.DBDSN, "r5cfg") || !strings.Contains(cfg.DBDSN, user) {
		t.Fatalf("DBDSN not built from exploded params: %q", cfg.DBDSN)
	}
	// The password must come from env, and must not have been defaulted empty.
	if cfg.DBPassword != pw {
		t.Fatalf("DBPassword = %q, want the env value", cfg.DBPassword)
	}

	app, err := NewApp(cfg, "test-secret", "hash", slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatalf("build composition root: %v", err)
	}
	startCtx, startCancel := context.WithTimeout(ctx, 60*time.Second)
	defer startCancel()
	if err := app.Start(startCtx); err != nil {
		t.Fatalf("config-driven postgres startup: %v", err)
	}
	t.Log("config.yaml(db.dialect=postgres + params) + DB_PASSWORD env → postgres app booted (ready gates green)")
	stopCtx, stopCancel := context.WithTimeout(ctx, 15*time.Second)
	defer stopCancel()
	_ = app.Stop(stopCtx)
}
