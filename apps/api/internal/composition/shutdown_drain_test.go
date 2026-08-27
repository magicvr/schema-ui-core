package composition

import (
	"bufio"
	"context"
	"database/sql"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"go.uber.org/fx"

	"github.com/magicvr/schema-ui-core/apps/api/internal/config"
	"github.com/magicvr/schema-ui-core/apps/api/internal/pgtest"
)

// VP-021 contract v0.1.0 §8 in-process equivalent harness (runs on every OS;
// the process-level SIGTERM variant lives in cmd/server and runs on linux/CI,
// which owns the real exit-code assertions).
//
// clean drain:   an in-flight request (body still being read) completes
//                within the drain budget; Stop returns nil.
// budget hole:   an in-flight request that never completes exhausts the
//                budget; Stop returns a deadline error (forced-exit path).
//
// Each scenario first proves the request is truly in-flight (a short client
// read yields no response while the body is incomplete) so the drain is never
// raced against the server's accept loop.

func freeLoopbackAddr(t *testing.T) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := ln.Addr().String()
	_ = ln.Close()
	return addr
}

func waitReadyHTTP(t *testing.T, addr string) {
	t.Helper()
	client := &http.Client{Timeout: 2 * time.Second}
	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		resp, err := client.Get("http://" + addr + "/readyz")
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("readyz never became 200 at %s", addr)
}

const loginBody = `{"username":"admin","password":"x"}`

// openDribbleLogin sends the request head with an incomplete body and proves
// the request is in-flight (no response while the body is incomplete).
func openDribbleLogin(t *testing.T, addr string) net.Conn {
	t.Helper()
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}
	head := "POST /api/auth/login HTTP/1.1\r\nHost: harness\r\nContent-Type: application/json\r\n" +
		"Content-Length: " + strconv.Itoa(len(loginBody)) + "\r\nConnection: close\r\n\r\n" + loginBody[:16]
	if _, err := conn.Write([]byte(head)); err != nil {
		_ = conn.Close()
		t.Fatal(err)
	}
	// Prove in-flight: a short client read must time out (the handler blocks
	// on the incomplete body), never return a response already.
	_ = conn.SetReadDeadline(time.Now().Add(400 * time.Millisecond))
	probe := make([]byte, 16)
	if n, err := conn.Read(probe); err == nil {
		_ = conn.Close()
		t.Fatalf("server responded before the request body completed: n=%d data=%q", n, probe[:n])
	}
	_ = conn.SetReadDeadline(time.Time{})
	t.Cleanup(func() { _ = conn.Close() })
	return conn
}

type drainResult struct {
	line string
	err  error
}

// finishDribble completes the body after delay and returns the status line.
func finishDribble(conn net.Conn, bodyRemainder string, delay time.Duration) drainResult {
	time.Sleep(delay)
	if _, err := conn.Write([]byte(bodyRemainder)); err != nil {
		return drainResult{err: err}
	}
	_ = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	line, err := bufio.NewReader(conn).ReadString('\n')
	return drainResult{line: strings.TrimSpace(line), err: err}
}

func startHarnessApp(t *testing.T, addr string) *fx.App {
	t.Helper()
	cfg := lifecycleAppConfig(t, "mvp", addr)
	// Keep the dribble connection alive past the drain budget so the
	// in-flight request (and not the read timeout) decides shutdown timing.
	cfg.ReadTimeout = time.Minute
	app, err := NewApp(cfg, "test-secret", "hash", slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatal(err)
	}
	startCtx, startCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer startCancel()
	if err := app.Start(startCtx); err != nil {
		t.Fatalf("start: %v", err)
	}
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = app.Stop(ctx)
	})
	waitReadyHTTP(t, addr)
	return app
}

func TestShutdownDrainHarness(t *testing.T) {
	t.Run("clean drain: in-flight request completes within budget", func(t *testing.T) {
		addr := freeLoopbackAddr(t)
		app := startHarnessApp(t, addr)
		conn := openDribbleLogin(t, addr) // request proven in-flight

		got := make(chan drainResult, 1)
		go func() { got <- finishDribble(conn, loginBody[16:], 300*time.Millisecond) }()

		stopCtx, stopCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer stopCancel()
		if err := app.Stop(stopCtx); err != nil {
			t.Fatalf("clean-drain Stop = %v, want nil (in-flight request drained within budget)", err)
		}
		// Contract §8: after shutdown the listener is closed — new
		// connections must be refused, not accepted.
		if conn2, err := net.Dial("tcp", addr); err == nil {
			_ = conn2.Close()
			t.Fatal("new connection accepted after Stop; want refused (listener closed)")
		}
		select {
		case res := <-got:
			if res.err != nil {
				t.Fatalf("in-flight request failed during drain: %v", res.err)
			}
			if !strings.HasPrefix(res.line, "HTTP/1.1") {
				t.Fatalf("in-flight response status line = %q", res.line)
			}
		case <-time.After(3 * time.Second):
			t.Fatal("in-flight request did not complete during drain")
		}
	})

	t.Run("budget hole: stuck request exhausts budget and forces Stop error", func(t *testing.T) {
		addr := freeLoopbackAddr(t)
		app := startHarnessApp(t, addr)
		conn := openDribbleLogin(t, addr) // request proven in-flight; body never completed

		stopCtx, stopCancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer stopCancel()
		err := app.Stop(stopCtx)
		if err == nil {
			t.Fatal("Stop = nil, want deadline error when an in-flight request exceeds the budget")
		}
		if stopCtx.Err() == nil {
			t.Fatal("stop context not expired after Stop returned error")
		}
		_ = conn.Close() // release the stuck read so cleanup Stop settles
	})
}

// TestShutdownDrainHarnessPostgres runs the clean-drain scenario on the
// PostgreSQL dialect (contract §5: dual-dialect store-drain consistency).
// Gated by PG_TEST_* like the other postgres integration tests (CI only).
func TestShutdownDrainHarnessPostgres(t *testing.T) {
	dsn := pgtest.DSN()
	if dsn == "" {
		t.Skip("postgres test env not set (PG_TEST_*); skipping postgres drain harness")
	}
	ctx := context.Background()
	u, err := url.Parse(dsn)
	if err != nil {
		t.Fatal(err)
	}
	const dbName = "shutdownharness"
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

	addr := freeLoopbackAddr(t)
	cfg := &config.Config{
		AppName:      "test",
		AppEnv:       "development",
		HTTPAddr:     addr,
		DBDialect:    "postgres",
		DBPath:       filepath.Join(t.TempDir(), "harness.db"),
		DBDSN:        pgDSN,
		ProfileName:  "mvp",
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Second,
		IdleTimeout:  time.Second,
	}
	app, err := NewApp(cfg, "test-secret", "hash", slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatal(err)
	}
	startCtx, startCancel := context.WithTimeout(ctx, 60*time.Second)
	defer startCancel()
	if err := app.Start(startCtx); err != nil {
		t.Fatalf("postgres drain harness start: %v", err)
	}
	t.Cleanup(func() {
		stopCtx, stopCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer stopCancel()
		_ = app.Stop(stopCtx)
	})
	waitReadyHTTP(t, addr)

	conn := openDribbleLogin(t, addr) // request proven in-flight
	got := make(chan drainResult, 1)
	go func() { got <- finishDribble(conn, loginBody[16:], 300*time.Millisecond) }()

	stopCtx, stopCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer stopCancel()
	if err := app.Stop(stopCtx); err != nil {
		t.Fatalf("postgres clean-drain Stop = %v, want nil", err)
	}
	select {
	case res := <-got:
		if res.err != nil {
			t.Fatalf("postgres in-flight request failed during drain: %v", res.err)
		}
		if !strings.HasPrefix(res.line, "HTTP/1.1") {
			t.Fatalf("postgres in-flight response status line = %q", res.line)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("postgres in-flight request did not complete during drain")
	}
}