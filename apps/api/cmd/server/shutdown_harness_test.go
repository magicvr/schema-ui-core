//go:build !windows

// Process-level shutdown harness for the VP-021 graceful-shutdown contract
// (v0.1.0 §1–§3, §7, §8): builds the real cmd/server binary, drives it with
// real SIGTERM, and asserts the exit-code + structured-log contract.
//
// Runs on linux/macOS (and CI). Windows does not support SIGTERM delivery to
// child processes via os.Process.Signal; the in-process equivalent harness
// in internal/composition/shutdown_drain_test.go covers the drain semantics
// on every OS instead.

package main_test

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

var serverBin string

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "shutdown-harness-*")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	serverBin = filepath.Join(dir, "schema-ui-server")
	build := exec.Command("go", "build", "-o", serverBin, ".")
	if out, err := build.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "build cmd/server: %v\n%s", err, out)
		os.Exit(1)
	}
	code := m.Run()
	_ = os.RemoveAll(dir)
	os.Exit(code)
}

func freeProcAddr(t *testing.T) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := ln.Addr().String()
	_ = ln.Close()
	return addr
}

// startServerProc launches the built binary with an isolated environment:
// dev mode, temporary store, free port, optional HTTP_SHUTDOWN_TIMEOUT, and
// repo-local config files removed (embedded defaults only).
func startServerProc(t *testing.T, shutdownTimeout string) (*exec.Cmd, *bytes.Buffer, string) {
	t.Helper()
	addr := freeProcAddr(t)
	cmd := exec.Command(serverBin)
	env := make([]string, 0, 8)
	for _, kv := range os.Environ() {
		if strings.HasPrefix(kv, "CONFIG_FILE=") || strings.HasPrefix(kv, "CONFIG_ENV_FILE=") {
			continue
		}
		env = append(env, kv)
	}
	env = append(env,
		"APP_ENV=development",
		"HTTP_ADDR="+addr,
		"DB_PATH="+filepath.Join(t.TempDir(), "harness.db"),
	)
	if shutdownTimeout != "" {
		env = append(env, "HTTP_SHUTDOWN_TIMEOUT="+shutdownTimeout)
	}
	cmd.Env = env

	var logBuf bytes.Buffer
	cmd.Stdout = &logBuf
	cmd.Stderr = &logBuf
	if err := cmd.Start(); err != nil {
		t.Fatalf("start server: %v", err)
	}
	t.Cleanup(func() {
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
		_ = cmd.Wait()
	})
	waitProcReady(t, addr)
	return cmd, &logBuf, addr
}

func waitProcReady(t *testing.T, addr string) {
	t.Helper()
	client := &http.Client{Timeout: 2 * time.Second}
	deadline := time.Now().Add(20 * time.Second)
	for time.Now().Before(deadline) {
		resp, err := client.Get("http://" + addr + "/readyz")
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("server never became ready at %s", addr)
}

func procDribbleLogin(t *testing.T, addr string) net.Conn {
	t.Helper()
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}
	body := `{"username":"admin","password":"x"}`
	head := "POST /api/auth/login HTTP/1.1\r\nHost: harness\r\nContent-Type: application/json\r\n" +
		"Content-Length: " + fmt.Sprint(len(body)) + "\r\nConnection: close\r\n\r\n" + body[:16]
	if _, err := conn.Write([]byte(head)); err != nil {
		_ = conn.Close()
		t.Fatal(err)
	}
	// Prove the request is in-flight before the signal: a short client read
	// must time out (handler blocked on the incomplete body).
	_ = conn.SetReadDeadline(time.Now().Add(400 * time.Millisecond))
	probe := make([]byte, 16)
	if n, err := conn.Read(probe); err == nil {
		_ = conn.Close()
		t.Fatalf("server responded before the request body completed: n=%d data=%q", n, probe[:n])
	}
	_ = conn.SetReadDeadline(time.Time{})
	return conn
}

// TestShutdownCleanDrainExitZero (contract §1/§3/§8): an in-flight request
// completes within the budget, the process logs shutdown.complete and exits 0.
func TestShutdownCleanDrainExitZero(t *testing.T) {
	cmd, logBuf, addr := startServerProc(t, "")

	conn := procDribbleLogin(t, addr)
	responded := make(chan bool, 1)
	go func() {
		time.Sleep(400 * time.Millisecond) // finish the body after the signal
		_, _ = conn.Write([]byte(`in":"password":"x"}`))
		_ = conn.SetReadDeadline(time.Now().Add(6 * time.Second))
		if _, err := bufio.NewReader(conn).ReadString('\n'); err == nil {
			responded <- true
		}
	}()

	time.Sleep(200 * time.Millisecond) // request is in-flight (reading body)
	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		t.Fatalf("SIGTERM: %v", err)
	}

	waitErr := make(chan error, 1)
	go func() { waitErr <- cmd.Wait() }()
	select {
	case err := <-waitErr:
		if err != nil {
			t.Fatalf("exit = %v, want 0 (clean drain)", err)
		}
	case <-time.After(20 * time.Second):
		t.Fatal("process did not exit after SIGTERM")
	}
	if !strings.Contains(logBuf.String(), "shutdown.complete") {
		t.Errorf("stdout missing shutdown.complete:\n%s", logBuf.String())
	}
	select {
	case <-responded:
	case <-time.After(3 * time.Second):
		t.Error("in-flight request did not complete during drain")
	}
}

// TestShutdownTimeoutExitOne (contract §1/§3/§8): the drain budget (1s) is
// exhausted by a stuck in-flight request; the process logs shutdown.timeout
// and exits 1.
func TestShutdownTimeoutExitOne(t *testing.T) {
	cmd, logBuf, addr := startServerProc(t, "1s")

	conn := procDribbleLogin(t, addr) // body never completed
	defer conn.Close()

	time.Sleep(200 * time.Millisecond)
	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		t.Fatalf("SIGTERM: %v", err)
	}

	waitErr := make(chan error, 1)
	go func() { waitErr <- cmd.Wait() }()
	select {
	case err := <-waitErr:
		if err == nil {
			t.Fatal("exit = 0, want 1 (budget timeout must force exit 1)")
		}
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() != 1 {
				t.Fatalf("exit code = %d, want 1", exitErr.ExitCode())
			}
		} else {
			t.Fatalf("unexpected wait error: %v", err)
		}
	case <-time.After(15 * time.Second):
		t.Fatal("process did not exit after SIGTERM")
	}
	if !strings.Contains(logBuf.String(), "shutdown.timeout") {
		t.Errorf("stdout missing shutdown.timeout:\n%s", logBuf.String())
	}
}