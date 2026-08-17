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

// S6 · 进程级重启持久化（L2，I-007-004 / D-007；GOAL-011 S3 改指 users 资源）：
// 真实 `cmd/server` OS 子进程以同一 `DB_PATH` 终止→重启，全 HTTP CRUD → 重启 →
// users/roles list/detail 符合预期。records 已由 0006 退场，本测试以双资源
// 承载跨进程持久化往返。
//
// 数据库隔离：每轮全新临时 DB_PATH + 空闲端口；测试结束 Kill+Wait 不留进程，
// 临时库与 pre-v0006 快照随 t.TempDir() 清理。迁移/seed 重跑、失败路径与
// 401/403 门禁由 store/handler 既有测试与 browser E2E 承担（I-007-004 §5/§6）。
func TestServerProcessRestartPersistsUsers(t *testing.T) {
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

	// Phase 1: first process (fresh DB seeds admin).
	srv1, log1 := startServer(t, bin, dbPath, port1)
	defer killServer(t, srv1, log1)
	waitHealth(t, port1, 20*time.Second)
	token := httpLogin(t, port1, "admin", "admin")
	// W16-F01: the fresh seed starts with must_change_password=1; replace the
	// initial password before using business APIs.
	token = httpChangePassword(t, port1, token, "admin", "admin-new-pass")
	createdID, created := httpCreateUser(t, port1, token)
	createdAt, _ := created["updatedAt"].(string)
	if createdAt == "" {
		t.Fatalf("created user updatedAt missing")
	}
	// a second user to delete (self/last-admin protections make user-admin undeletable).
	deleteID, _ := httpCreateUser2(t, port1, token)
	patched := httpPatch(t, port1, token, createdID, "Persisted Rebrand")
	patchedAt, _ := patched["updatedAt"].(string)
	if patchedAt == "" {
		t.Fatalf("patched createdID updatedAt missing")
	}
	httpDelete(t, port1, token, deleteID)
	// a created role (survives restart too, GOAL-011 S4 roles restart path).
	code, roleOut := httpDoJSON(t, port1, http.MethodPost, "/api/roles", `{"key":"ops","name":"Operator"}`, token)
	if code != http.StatusCreated {
		t.Fatalf("create role status = %d, want 201: %v", code, roleOut)
	}
	roleID, _ := roleOut["id"].(string)
	roleCreatedAt, _ := roleOut["createdAt"].(string)
	roleUpdatedAt, _ := roleOut["updatedAt"].(string)
	if roleID == "" || roleCreatedAt == "" || roleUpdatedAt == "" {
		t.Fatalf("created role missing persisted identity/timestamps: %v", roleOut)
	}
	for field, value := range map[string]string{"createdAt": roleCreatedAt, "updatedAt": roleUpdatedAt} {
		if _, err := time.Parse("2006-01-02T15:04:05.000Z", value); err != nil {
			t.Fatalf("created role %s = %q, want millisecond UTC timestamp: %v", field, value, err)
		}
	}
	killServer(t, srv1, log1)

	// Phase 2: restart the same binary against the same DB.
	srv2, log2 := startServer(t, bin, dbPath, port2)
	defer killServer(t, srv2, log2)
	waitHealth(t, port2, 20*time.Second)
	token2 := httpLogin(t, port2, "admin", "admin-new-pass")

	code, body := httpDoJSON(t, port2, http.MethodGet, "/api/users?pageSize=100", "", token2)
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
		t.Fatalf("created user %s missing after process restart; list=%v", createdID, present)
	}
	if present[createdID] != "Persisted Rebrand" {
		t.Fatalf("created user name after restart = %q, want Persisted Rebrand", present[createdID])
	}
	if _, ok := present[deleteID]; ok {
		t.Fatalf("deleted user %s resurrected after process restart", deleteID)
	}
	if total, _ := body["total"].(float64); total != 2 {
		t.Fatalf("total after restart = %v, want 2 (admin + created; no re-seed)", total)
	}

	code, roleList := httpDoJSON(t, port2, http.MethodGet, "/api/roles?pageSize=100", "", token2)
	if code != http.StatusOK {
		t.Fatalf("role list after restart status = %d, want 200", code)
	}
	roleItems := roleList["items"].([]any)
	var listedRole map[string]any
	for _, raw := range roleItems {
		item := raw.(map[string]any)
		if item["id"] == roleID {
			listedRole = item
			break
		}
	}
	if listedRole == nil {
		t.Fatalf("created role %s missing after process restart; list=%v", roleID, roleItems)
	}
	if listedRole["createdAt"] != roleCreatedAt || listedRole["updatedAt"] != roleUpdatedAt {
		t.Fatalf("role list timestamps after restart = %v, want createdAt=%s updatedAt=%s", listedRole, roleCreatedAt, roleUpdatedAt)
	}

	code, roleDetail := httpDoJSON(t, port2, http.MethodGet, "/api/roles/"+roleID, "", token2)
	if code != http.StatusOK {
		t.Fatalf("role detail after restart status = %d, want 200", code)
	}
	if roleDetail["key"] != "ops" || roleDetail["name"] != "Operator" || roleDetail["system"] != false {
		t.Fatalf("role detail after restart = %v", roleDetail)
	}
	if roleDetail["id"] != roleID || roleDetail["createdAt"] != roleCreatedAt || roleDetail["updatedAt"] != roleUpdatedAt {
		t.Fatalf("role detail identity/timestamps after restart = %v, want id=%s createdAt=%s updatedAt=%s", roleDetail, roleID, roleCreatedAt, roleUpdatedAt)
	}

	// detail after restart: created user fields + updatedAt match Phase 1
	// create/patch responses ms-exactly (cross-process persistence round-trip).
	code, createdDetail := httpDoJSON(t, port2, http.MethodGet, "/api/users/"+createdID, "", token2)
	if code != http.StatusOK {
		t.Fatalf("created detail status = %d, want 200", code)
	}
	if createdDetail["username"] != "persist" || createdDetail["name"] != "Persisted Rebrand" {
		t.Fatalf("created detail after restart = %v", createdDetail)
	}
	if createdDetail["updatedAt"] != patchedAt {
		t.Fatalf("created updatedAt after restart = %v, want %v (persisted ms)", createdDetail["updatedAt"], patchedAt)
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

func httpChangePassword(t *testing.T, port, token, current, next string) string {
	t.Helper()
	body := fmt.Sprintf(`{"currentPassword":%q,"newPassword":%q}`, current, next)
	code, out := httpDoJSON(t, port, http.MethodPost, "/api/account/password", body, token)
	if code != http.StatusOK {
		t.Fatalf("forced password change status = %d, want 200", code)
	}
	access, _ := out["accessToken"].(string)
	if access == "" {
		t.Fatalf("forced password change missing accessToken in %v", out)
	}
	return access
}

func httpCreateUser(t *testing.T, port, token string) (string, map[string]any) {
	t.Helper()
	code, out := httpDoJSON(t, port, http.MethodPost, "/api/users", `{"username":"persist","name":"Persisted Co","password":"pw123456"}`, token)
	if code != http.StatusCreated {
		t.Fatalf("create status = %d, want 201", code)
	}
	id, _ := out["id"].(string)
	if id == "" {
		t.Fatalf("create id missing in %v", out)
	}
	return id, out
}

// httpCreateUser2 creates a disposable second user for the delete-persistence
// assertion (user-admin is self/last-admin protected).
func httpCreateUser2(t *testing.T, port, token string) (string, map[string]any) {
	t.Helper()
	code, out := httpDoJSON(t, port, http.MethodPost, "/api/users", `{"username":"doomed","name":"Doomed","password":"pw123456"}`, token)
	if code != http.StatusCreated {
		t.Fatalf("create2 status = %d, want 201", code)
	}
	id, _ := out["id"].(string)
	if id == "" {
		t.Fatalf("create2 id missing in %v", out)
	}
	return id, out
}

// httpPatch issues a PATCH and returns the 200 response body so the caller can
// capture the refreshed `updatedAt` for the cross-process detail assertion.
func httpPatch(t *testing.T, port, token, id, name string) map[string]any {
	t.Helper()
	code, out := httpDoJSON(t, port, http.MethodPatch, "/api/users/"+id, fmt.Sprintf(`{"name":%q}`, name), token)
	if code != http.StatusOK {
		t.Fatalf("patch %s status = %d, want 200", id, code)
	}
	return out
}

func httpDelete(t *testing.T, port, token, id string) {
	t.Helper()
	code, _ := httpDoJSON(t, port, http.MethodDelete, "/api/users/"+id, "", token)
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
