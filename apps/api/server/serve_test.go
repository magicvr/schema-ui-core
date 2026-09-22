package server

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// startTestServer 绑定一个临时端口并把 serve 面跑在 goroutine 中。
// stop() 取消 ctx 并断言 Run 干净返回（RT-D02 排空成功路径）。
func startTestServer(t *testing.T) (addr string, stop func()) {
	t.Helper()
	cfg, err := LoadConfig("")
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr = ln.Addr().String()
	_ = ln.Close()
	cfg.HTTPAddr = addr
	cfg.DBPath = filepath.Join(t.TempDir(), "serve-test.db")

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		_, err := Run(ctx, Options{
			Config: cfg,
			Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		}, nil)
		done <- err
	}()

	hc := &http.Client{Timeout: 2 * time.Second}
	deadline := time.Now().Add(60 * time.Second)
	for {
		resp, err := hc.Get("http://" + addr + "/healthz")
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				break
			}
		}
		if time.Now().After(deadline) {
			t.Fatal("serve 面未在 60s 内就绪（healthz 未 200）")
		}
		time.Sleep(150 * time.Millisecond)
	}
	stop = func() {
		cancel()
		select {
		case err := <-done:
			if err != nil {
				t.Errorf("干净停机应返回 nil，got: %v", err)
			}
		case <-time.After(20 * time.Second):
			t.Error("停机排空超时（20s）")
		}
	}
	return addr, stop
}

func TestRunServesHealthAndManifestAndLogin(t *testing.T) {
	addr, stop := startTestServer(t)
	defer stop()

	hc := &http.Client{Timeout: 5 * time.Second}
	base := "http://" + addr

	// C1/C5：healthz / readyz 可响应。
	for _, path := range []string{"/healthz", "/readyz"} {
		resp, err := hc.Get(base + path)
		if err != nil {
			t.Fatalf("GET %s: %v", path, err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("GET %s = %d, want 200", path, resp.StatusCode)
		}
	}

	// 中央面：manifest / bootstrap 200。
	for _, path := range []string{"/.well-known/schema-ui/app-manifest.json", "/.well-known/schema-ui/host-bootstrap.json"} {
		resp, err := hc.Get(base + path)
		if err != nil {
			t.Fatalf("GET %s: %v", path, err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("GET %s = %d, want 200", path, resp.StatusCode)
		}
	}

	// C5：登录路径（dev 种子 admin/admin）→ 200（中央 auth 面接线成立）。
	body := bytes.NewBufferString(`{"username":"admin","password":"admin"}`)
	resp, err := hc.Post(base+"/api/auth/login", "application/json", body)
	if err != nil {
		t.Fatalf("POST /api/auth/login: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("POST /api/auth/login = %d, want 200", resp.StatusCode)
	}
}

func TestRunRejectsBadConfig(t *testing.T) {
	if _, err := LoadConfig(""); err != nil {
		t.Fatalf("dev defaults should load: %v", err)
	}
	// server.Serve / Run 对 nil Config fail-closed。
	if _, err := Run(context.Background(), Options{}, nil); err == nil {
		t.Error("Run(nil Config) should fail")
	}
}

// runStartGuardStore makes a regression that reaches store assembly fail
// immediately without opening a database or starting a server.
type runStartGuardStore struct {
	used bool
}

func (s *runStartGuardStore) markUsed() { s.used = true }

func (s *runStartGuardStore) Dialect() kernel.Dialect {
	s.markUsed()
	return kernel.DialectSQLite
}

func (s *runStartGuardStore) Run(context.Context, func(kernel.Tx) error) error {
	s.markUsed()
	return errors.New("run start guard store was used")
}

func (s *runStartGuardStore) Ping(context.Context) error {
	s.markUsed()
	return errors.New("run start guard store was used")
}

func (s *runStartGuardStore) Close() error {
	s.markUsed()
	return errors.New("run start guard store was used")
}

func (s *runStartGuardStore) WasFresh() bool {
	s.markUsed()
	return false
}

func (s *runStartGuardStore) MarkSystemDataReady() { s.markUsed() }

func (s *runStartGuardStore) SystemDataReady() error {
	s.markUsed()
	return errors.New("run start guard store was used")
}

type serveMFAProbeStore struct {
	active     bool
	queryErr   error
	runCalls   int
	closeCalls int
	pingCalls  int
	readyMarks int
}

func (s *serveMFAProbeStore) Dialect() kernel.Dialect { return kernel.DialectSQLite }

func (s *serveMFAProbeStore) Run(_ context.Context, fn func(kernel.Tx) error) error {
	s.runCalls++
	return fn(serveMFAProbeTx{active: s.active, queryErr: s.queryErr})
}

func (s *serveMFAProbeStore) Ping(context.Context) error {
	s.pingCalls++
	return errors.New("serve MFA probe store was pinged")
}

func (s *serveMFAProbeStore) Close() error {
	s.closeCalls++
	return nil
}

func (s *serveMFAProbeStore) WasFresh() bool { return false }

func (s *serveMFAProbeStore) MarkSystemDataReady() { s.readyMarks++ }

func (s *serveMFAProbeStore) SystemDataReady() error {
	return errors.New("serve MFA probe store reached runtime startup")
}

type serveMFAProbeTx struct {
	active   bool
	queryErr error
}

func (serveMFAProbeTx) Exec(context.Context, string, ...any) (kernel.Result, error) {
	return nil, errors.New("serve MFA probe transaction executed an unexpected write")
}

func (serveMFAProbeTx) Query(context.Context, string, ...any) (kernel.Rows, error) {
	return nil, errors.New("serve MFA probe transaction executed an unexpected query")
}

func (tx serveMFAProbeTx) QueryRow(context.Context, string, ...any) kernel.Row {
	return serveMFAProbeRow{active: tx.active, err: tx.queryErr}
}

type serveMFAProbeRow struct {
	active bool
	err    error
}

func (row serveMFAProbeRow) Scan(dest ...any) error {
	if row.err != nil {
		return row.err
	}
	if len(dest) != 1 {
		return errors.New("serve MFA probe row received an unexpected scan shape")
	}
	active, ok := dest[0].(*bool)
	if !ok {
		return errors.New("serve MFA probe row received a non-bool destination")
	}
	*active = row.active
	return nil
}

func newMFAProbeConfig(t *testing.T) (*Config, net.Listener) {
	t.Helper()
	reserved, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	return &Config{
		AppEnv:          "development",
		HTTPAddr:        reserved.Addr().String(),
		ShutdownTimeout: time.Second,
		DBDialect:       "sqlite",
		DBPath:          filepath.Join(t.TempDir(), "must-not-open.db"),
		AuthAccessTTL:   15 * time.Minute,
		AuthRefreshTTL:  30 * 24 * time.Hour,
	}, reserved
}

func TestRunRejectsActiveMFAEnrollmentBeforeStartup(t *testing.T) {
	store := &serveMFAProbeStore{active: true}
	cfg, reserved := newMFAProbeConfig(t)
	defer reserved.Close()

	_, err := Run(context.Background(), Options{Config: cfg, Store: store}, nil)
	if err == nil {
		t.Fatal("Run should reject active MFA enrollment without an assembled verifier")
	}
	if !strings.Contains(err.Error(), "serve did not assemble an MFA verifier") {
		t.Fatalf("error = %v, want missing MFA verifier detail", err)
	}
	if !strings.Contains(err.Error(), "active MFA enrollment exists") {
		t.Fatalf("error = %v, want active enrollment detail", err)
	}
	if store.runCalls != 1 {
		t.Fatalf("MFA probe Run calls = %d, want 1", store.runCalls)
	}
	if store.closeCalls != 1 {
		t.Fatalf("store Close calls = %d, want 1", store.closeCalls)
	}
	if store.pingCalls != 0 || store.readyMarks != 0 {
		t.Fatalf("startup continued after active enrollment: ping=%d readyMarks=%d", store.pingCalls, store.readyMarks)
	}
}

func TestRunRejectsMFAEnrollmentProbeFailureBeforeStartup(t *testing.T) {
	probeErr := errors.New("simulated active enrollment lookup failure")
	store := &serveMFAProbeStore{queryErr: probeErr}
	cfg, reserved := newMFAProbeConfig(t)
	defer reserved.Close()

	_, err := Run(context.Background(), Options{Config: cfg, Store: store}, nil)
	if err == nil {
		t.Fatal("Run should reject an MFA enrollment probe failure")
	}
	if !strings.Contains(err.Error(), "serve did not assemble an MFA verifier") {
		t.Fatalf("error = %v, want missing MFA verifier detail", err)
	}
	if !strings.Contains(err.Error(), probeErr.Error()) {
		t.Fatalf("error = %v, want probe failure detail", err)
	}
	if store.runCalls != 1 {
		t.Fatalf("MFA probe Run calls = %d, want 1", store.runCalls)
	}
	if store.closeCalls != 1 {
		t.Fatalf("store Close calls = %d, want 1", store.closeCalls)
	}
	if store.pingCalls != 0 || store.readyMarks != 0 {
		t.Fatalf("startup continued after probe failure: ping=%d readyMarks=%d", store.pingCalls, store.readyMarks)
	}
}

func TestRunRejectsNonDevelopmentJWTSecretBeforeStartup(t *testing.T) {
	for name, secret := range map[string]string{
		"missing": "",
		"weak":    "short",
	} {
		t.Run(name, func(t *testing.T) {
			guard := &runStartGuardStore{}
			cfg := &Config{
				AppEnv:               "production",
				HTTPAddr:             "127.0.0.1:0",
				ShutdownTimeout:      time.Second,
				DBDialect:            "sqlite",
				DBPath:               filepath.Join(t.TempDir(), "must-not-open.db"),
				AuthJWTSecret:        secret,
				AuthAccessTTL:        15 * time.Minute,
				AuthRefreshTTL:       30 * 24 * time.Hour,
				AdminInitialPassword: "seed-password-ok",
			}

			// Keep the target occupied so a regression that gets past config
			// validation cannot accidentally start a listener in this test.
			reserved, err := net.Listen("tcp", cfg.HTTPAddr)
			if err != nil {
				t.Fatal(err)
			}
			defer reserved.Close()
			cfg.HTTPAddr = reserved.Addr().String()

			if _, err := Run(context.Background(), Options{Config: cfg, Store: guard}, nil); err == nil {
				t.Fatal("Run should reject invalid non-development JWT configuration")
			} else {
				if !strings.Contains(err.Error(), "server: validate config") {
					t.Fatalf("error = %v, want config validation failure", err)
				}
				if !strings.Contains(err.Error(), "AUTH_JWT_SECRET") {
					t.Fatalf("error = %v, want AUTH_JWT_SECRET validation detail", err)
				}
			}
			if guard.used {
				t.Fatal("invalid JWT configuration must fail before store startup")
			}
		})
	}
}

func TestResolveSecretKeepsDevelopmentFallbackOnly(t *testing.T) {
	const devFallback = "dev-only-insecure-jwt-secret-change-me"
	if got := resolveSecret(&Config{AppEnv: "development"}); got != devFallback {
		t.Fatalf("development fallback = %q, want %q", got, devFallback)
	}
	if got := resolveSecret(&Config{AppEnv: "production"}); got != "" {
		t.Fatalf("non-development missing secret = %q, want empty", got)
	}
}

// W15 F-003 (GOAL-016 A-001): a production bootstrap seed that violates the
// frozen 8–72 byte policy must fail startup BEFORE the listener comes up —
// the seed is hashed only after the policy check.
func TestRunRejectsWeakSeedPasswordNonDev(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("AUTH_JWT_SECRET", "0123456789abcdef0123456789abcdef")
	t.Setenv("ADMIN_INITIAL_PASSWORD", "weak")
	cfg, err := LoadConfig("")
	if err != nil {
		t.Fatalf("LoadConfig(production) with compliant secret should load (weak seed check happens at bootstrap): %v", err)
	}
	cfg.HTTPAddr = "127.0.0.1:0"
	cfg.DBPath = filepath.Join(t.TempDir(), "weak-seed.db")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if _, err := Run(ctx, Options{
		Config: cfg,
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}, nil); err == nil {
		t.Fatal("production bootstrap with a 4-byte seed must fail closed")
	} else if !strings.Contains(err.Error(), "bootstrap seed password") {
		t.Fatalf("error = %v, want bootstrap seed password policy failure", err)
	}
}
