package server

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoadConfigDefaults(t *testing.T) {
	cfg, err := LoadConfig("")
	if err != nil {
		t.Fatalf("LoadConfig(default) error: %v", err)
	}
	if cfg.HTTPAddr != "127.0.0.1:25080" {
		t.Errorf("HTTPAddr = %q, want 127.0.0.1:25080 (W15 F-001 loopback default)", cfg.HTTPAddr)
	}
	if cfg.ShutdownTimeout != 10*time.Second {
		t.Errorf("ShutdownTimeout = %s, want 10s", cfg.ShutdownTimeout)
	}
	if cfg.DBDialect != "sqlite" || cfg.DBPath != "./data/schema-ui.db" {
		t.Errorf("db defaults = %q/%q", cfg.DBDialect, cfg.DBPath)
	}
	if cfg.ProfileName != "admin" {
		t.Errorf("ProfileName = %q, want admin", cfg.ProfileName)
	}
}

func TestLoadConfigExplicitFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	content := `
app:
  env: development
http:
  addr: ":28080"
  shutdown_timeout: 3s
  trusted_proxies:
    - "10.0.0.0/8"
db:
  dialect: postgres
  dsn: "postgres://u:p@127.0.0.1:5432/db"
auth:
  access_ttl: 5m
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig(file) error: %v", err)
	}
	if cfg.HTTPAddr != ":28080" || cfg.ShutdownTimeout != 3*time.Second {
		t.Errorf("http overrides not applied: %q %s", cfg.HTTPAddr, cfg.ShutdownTimeout)
	}
	if len(cfg.TrustedProxies) != 1 || cfg.TrustedProxies[0] != "10.0.0.0/8" {
		t.Errorf("trusted proxies = %v", cfg.TrustedProxies)
	}
	if cfg.DBDialect != "postgres" || cfg.DBDSN == "" {
		t.Errorf("db overrides not applied: %q %q", cfg.DBDialect, cfg.DBDSN)
	}
	if cfg.AuthAccessTTL != 5*time.Minute {
		t.Errorf("access_ttl = %s", cfg.AuthAccessTTL)
	}
}

func TestLoadConfigMissingExplicitFile(t *testing.T) {
	if _, err := LoadConfig("does-not-exist.yaml"); err == nil {
		t.Fatal("expected error for missing explicit config file")
	}
}

func TestLoadConfigInterpolationFailClosed(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte("auth:\n  jwt_secret: ${SERVE_TEST_MUST_NOT_BE_SET}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadConfig(path); err == nil {
		t.Fatal("expected fail-closed error for bare unset ${VAR}")
	}
}

func TestLoadConfigEnvOverride(t *testing.T) {
	t.Setenv("HTTP_SHUTDOWN_TIMEOUT", "2s")
	t.Setenv("DB_DIALECT", "postgres")
	t.Setenv("DB_DSN", "postgres://x@127.0.0.1/db")
	cfg, err := LoadConfig("")
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}
	if cfg.ShutdownTimeout != 2*time.Second {
		t.Errorf("ShutdownTimeout = %s, want 2s", cfg.ShutdownTimeout)
	}
	if cfg.DBDialect != "postgres" {
		t.Errorf("DBDialect = %q, want postgres", cfg.DBDialect)
	}
}

func TestLoadConfigInvalidShutdownTimeout(t *testing.T) {
	for _, val := range []string{"0s", "-1s", "abc"} {
		dir := t.TempDir()
		path := filepath.Join(dir, "config.yaml")
		// F-009（A-003）：显式声明 development，确保用例真正咬到超时分支而非
		// 空 APP_ENV 门禁（F-001 后 validate 先拒绝空 env）。
		if err := os.WriteFile(path, []byte("app:\n  env: development\nhttp:\n  shutdown_timeout: "+val+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := LoadConfig(path); err == nil {
			t.Errorf("shutdown_timeout %q: expected fail-closed error", val)
		}
	}
}

func TestLoadConfigDialectPairing(t *testing.T) {
	dir := t.TempDir()
	// sqlite + dsn → 拒绝（F-009：显式 development，避免先命中空 APP_ENV 门禁）
	path := filepath.Join(dir, "bad-sqlite.yaml")
	if err := os.WriteFile(path, []byte("app:\n  env: development\ndb:\n  dialect: sqlite\n  dsn: postgres://x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadConfig(path); err == nil {
		t.Error("sqlite+dsn: expected error")
	}
	// postgres 无 dsn → 拒绝
	path = filepath.Join(dir, "bad-pg.yaml")
	if err := os.WriteFile(path, []byte("app:\n  env: development\ndb:\n  dialect: postgres\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadConfig(path); err == nil {
		t.Error("postgres without dsn: expected error")
	}
}

func TestLoadConfigRequiresExplicitAppEnv(t *testing.T) {
	// W15 F-001: a custom YAML that omits app.env must fail closed — the
	// embedded default pins "development", so only explicit configs can
	// silently omit it, and guessing an environment from silence is refused.
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte("http:\n  addr: \"127.0.0.1:28080\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadConfig(path); err == nil {
		t.Fatal("custom YAML without app.env must fail closed (refusing to guess)")
	}
	// Explicit development (or production with secrets) loads fine.
	if err := os.WriteFile(path, []byte("app:\n  env: development\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadConfig(path); err != nil {
		t.Fatalf("explicit development should load: %v", err)
	}
}

func TestLoadConfigJWTSecretStrengthNonDev(t *testing.T) {
	// W15 F-002: non-development reuses the main production bar — a short or
	// single-class secret is a startup error, not a warning.
	t.Setenv("APP_ENV", "production")
	t.Setenv("ADMIN_INITIAL_PASSWORD", "seed-password-ok")
	for name, secret := range map[string]string{
		"short":     "short",
		"all-letter": "abcdefghijklmnopqrstuvwxyzabcdefghij",
		"all-digit":  "12345678901234567890123456789012",
	} {
		t.Setenv("AUTH_JWT_SECRET", secret)
		if _, err := LoadConfig(""); err == nil {
			t.Errorf("production with %s AUTH_JWT_SECRET must fail closed", name)
		}
	}
	t.Setenv("AUTH_JWT_SECRET", "0123456789abcdef0123456789abcdef")
	if _, err := LoadConfig(""); err != nil {
		t.Fatalf("production with compliant AUTH_JWT_SECRET should load: %v", err)
	}
}

func TestLoadConfigSecretFailClosedNonDev(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	if _, err := LoadConfig(""); err == nil {
		t.Fatal("expected fail-closed in production without AUTH_JWT_SECRET")
	}
	t.Setenv("AUTH_JWT_SECRET", "0123456789abcdef0123456789abcdef")
	t.Setenv("ADMIN_INITIAL_PASSWORD", "s3cret!")
	if _, err := LoadConfig(""); err != nil {
		t.Fatalf("production with secrets should load: %v", err)
	}
}

func TestLoadConfigUnknownKeysRejected(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte("http:\n  no_such_key: true\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadConfig(path); err == nil {
		t.Fatal("expected KnownFields rejection for unknown YAML key")
	}
}

func TestLoadConfigSecretsViaInterpolation(t *testing.T) {
	t.Setenv("SERVE_TEST_SECRET", "interp-secret")
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte("app:\n  env: development\nauth:\n  jwt_secret: ${SERVE_TEST_SECRET}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}
	if cfg.AuthJWTSecret != "interp-secret" {
		t.Errorf("jwt_secret = %q", cfg.AuthJWTSecret)
	}
}

func TestInterpolateAll(t *testing.T) {
	t.Setenv("SERVE_TEST_INTERP", "ok")
	got, err := interpolateAll("a=${SERVE_TEST_INTERP} b=${SERVE_TEST_DEFAULT:-fallback}")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "a=ok") || !strings.Contains(got, "b=fallback") {
		t.Errorf("interpolation result = %q", got)
	}
}