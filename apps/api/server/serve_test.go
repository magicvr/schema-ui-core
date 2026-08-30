package server

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
	"time"
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