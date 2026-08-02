package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

// S6 · 进程级重启持久化（L2，I-007-004 / D-007）：真实 `cmd/server` OS 子进程
// 以同一 `DB_PATH` 终止→重启，全 HTTP CRUD → 重启 → list/detail 符合预期。
// 补齐 A-003/A-004 标记的「store close/reopen 单测 ≠ 进程级重启」缺口。
//
// 数据库隔离：每轮全新临时 DB_PATH + 空闲端口；测试结束 Kill+Wait 不留进程，
// 临时库与 pre-v0003 快照随 t.TempDir() 清理。迁移/seed 重跑、失败路径与
// 401/403 门禁由 store/handler 既有测试与 browser E2E 承担（I-007-004 §5/§6）。
func TestServerProcessRestartPersistsRecords(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	apiRoot := filepath.Join(wd, "..", "..")

	bin := filepath.Join(t.TempDir(), "schema-server")
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}
	build := exec.Command("go", "build", "-o", bin, "./cmd/server")
	build.Dir = apiRoot
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build server: %v\n%s", err, out)
	}

	dbPath := filepath.Join(t.TempDir(), "restart.db")
	port1 := freePort(t)
	port2 := freePort(t)

	// Phase 1: first process (fresh DB seeds admin + 8 records).
	srv1, log1 := startServer(t, bin, dbPath, port1)
	waitHealth(t, port1, 20*time.Second)
	token := httpLogin(t, port1, "admin", "admin")
	createdID := httpCreate(t, port1, token)
	httpPatch(t, port1, token, "rec-1", "Acme Rebrand")
	httpDelete(t, port1, token, "rec-2")
	killServer(t, srv1, log1)

	// Phase 2: restart the same binary against the same DB.
	srv2, log2 := startServer(t, bin, dbPath, port2)
	defer killServer(t, srv2, log2)
	waitHealth(t, port2, 20*time.Second)
	token2 := httpLogin(t, port2, "admin", "admin")

	code, body := httpDoJSON(t, port2, http.MethodGet, "/api/records?pageSize=100", "", token2)
	if code != http.StatusOK {
		t.Fatalf("list after restart status = %d, want 200", code)
	}
	items := body["items"].([]any)
	present := map[string]string{}
	for _, raw := range items {
		item := raw.(map[string]any)
		id, _ := item["id"].(string)
		name, _ := item["name"].(string)
		present[id] = name
	}
	if _, ok := present[createdID]; !ok {
		t.Fatalf("created record %s missing after process restart; list=%v", createdID, present)
	}
	if present["rec-1"] != "Acme Rebrand" {
		t.Fatalf("rec-1 name after restart = %q, want Acme Rebrand", present["rec-1"])
	}
	if _, ok := present["rec-2"]; ok {
		t.Fatalf("deleted rec-2 resurrected after process restart")
	}
	if total, _ := body["total"].(float64); total != 8 {
		t.Fatalf("total after restart = %v, want 8 (no re-seed)", total)
	}

	_, detail := httpDoJSON(t, port2, http.MethodGet, "/api/records/"+createdID, "", token2)
	if detail["name"] != "Persisted Co" || detail["status"] != "active" || detail["owner"] != "zoe" {
		t.Fatalf("created detail after restart = %v", detail)
	}
}

func freePort(t *testing.T) string {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("free port: %v", err)
	}
	defer l.Close()
	return strconv.Itoa(l.Addr().(*net.TCPAddr).Port)
}

// startServer launches the built binary against a temp DB on a fixed port.
// Env is explicit: real session path (no dev-session fallback), seeded admin.
func startServer(t *testing.T, bin, dbPath, port string) (*exec.Cmd, *bytes.Buffer) {
	t.Helper()
	cmd := exec.Command(bin)
	cmd.Env = append(os.Environ(),
		"DB_PATH="+dbPath,
		"HTTP_ADDR=127.0.0.1:"+port,
		"ADMIN_INITIAL_PASSWORD=admin",
		"AUTH_JWT_SECRET=test-secret",
		"AUTH_DEV_SESSION_ENABLED=false",
		"APP_ENV=development",
		"LOG_LEVEL=error",
	)
	var logBuf bytes.Buffer
	cmd.Stdout = &logBuf
	cmd.Stderr = &logBuf
	if err := cmd.Start(); err != nil {
		t.Fatalf("start server: %v", err)
	}
	return cmd, &logBuf
}

func killServer(t *testing.T, cmd *exec.Cmd, logBuf *bytes.Buffer) {
	t.Helper()
	if cmd.Process != nil {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	}
	if logBuf.Len() > 0 && t.Failed() {
		t.Logf("server log:\n%s", logBuf.String())
	}
}

func waitHealth(t *testing.T, port string, timeout time.Duration) {
	t.Helper()
	url := "http://127.0.0.1:" + port + "/healthz"
	deadline := time.Now().Add(timeout)
	var lastErr error
	for time.Now().Before(deadline) {
		resp, err := http.Get(url)
		if err == nil {
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return
			}
			lastErr = fmt.Errorf("healthz status = %d", resp.StatusCode)
		} else {
			lastErr = err
		}
		time.Sleep(150 * time.Millisecond)
	}
	t.Fatalf("server not healthy on %s: %v", url, lastErr)
}

func httpLogin(t *testing.T, port, username, password string) string {
	t.Helper()
	body := fmt.Sprintf(`{"username":%q,"password":%q}`, username, password)
	code, out := httpDoJSON(t, port, http.MethodPost, "/api/auth/login", body, "")
	if code != http.StatusOK {
		t.Fatalf("login status = %d, want 200", code)
	}
	tok, _ := out["accessToken"].(string)
	if tok == "" {
		t.Fatalf("login accessToken missing in %v", out)
	}
	return tok
}

func httpCreate(t *testing.T, port, token string) string {
	t.Helper()
	code, out := httpDoJSON(t, port, http.MethodPost, "/api/records", `{"name":"Persisted Co","status":"active","owner":"zoe"}`, token)
	if code != http.StatusCreated {
		t.Fatalf("create status = %d, want 201", code)
	}
	id, _ := out["id"].(string)
	if id == "" {
		t.Fatalf("create id missing in %v", out)
	}
	return id
}

func httpPatch(t *testing.T, port, token, id, name string) {
	t.Helper()
	code, _ := httpDoJSON(t, port, http.MethodPatch, "/api/records/"+id, fmt.Sprintf(`{"name":%q}`, name), token)
	if code != http.StatusOK {
		t.Fatalf("patch %s status = %d, want 200", id, code)
	}
}

func httpDelete(t *testing.T, port, token, id string) {
	t.Helper()
	code, _ := httpDoJSON(t, port, http.MethodDelete, "/api/records/"+id, "", token)
	if code != http.StatusNoContent {
		t.Fatalf("delete %s status = %d, want 204", id, code)
	}
}

func httpDoJSON(t *testing.T, port, method, path, body, token string) (int, map[string]any) {
	t.Helper()
	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}
	req, err := http.NewRequest(method, "http://127.0.0.1:"+port+path, reader)
	if err != nil {
		t.Fatalf("%s %s: %v", method, path, err)
	}
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("%s %s: %v", method, path, err)
	}
	defer resp.Body.Close()
	var out map[string]any
	if resp.Body != nil {
		_ = json.NewDecoder(resp.Body).Decode(&out)
	}
	return resp.StatusCode, out
}
