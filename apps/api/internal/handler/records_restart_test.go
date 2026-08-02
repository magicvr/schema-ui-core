package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
)

// S6 · HTTP 层重启持久化（L1，I-007-004 / D-007）：同一临时 SQLite 文件上，
// 完整 handler/auth 栈的 store 关闭→重开（进程重启的持久化边界），全 HTTP
// CRUD → 重启 → list/detail 符合预期；迁移/seed 不重跑、已删行不复活。
//
// 这是 L1 快速回归；进程级重启由 cmd/server 的 L2 测试（真实子进程）承担。
func TestRecordsSurviveRestart(t *testing.T) {
	hash, err := auth.HashPassword(testSeedPassword, testBcryptCost)
	if err != nil {
		t.Fatalf("hash seed password: %v", err)
	}
	path := filepath.Join(t.TempDir(), "restart-records.db")

	// Phase 1: first server against the fresh DB (migrations {1,2,3}; empty
	// table seeds admin + 8 records rec-1…rec-8).
	st1, err := store.Open(path, testSeedUsername, hash, true)
	if err != nil {
		t.Fatalf("open phase 1: %v", err)
	}
	a1 := auth.New([]byte(testJWTSecret), 15*time.Minute, 30*24*time.Hour, st1, false)
	mux1 := http.NewServeMux()
	Register(mux1, a1, st1)
	env1 := &authTestEnv{mux: mux1, a: a1, st: st1}
	token1 := adminToken(t, env1)

	// create (POST /api/records) → 201 + new id.
	code, created := authedJSON(t, mux1, token1, http.MethodPost, "/api/records", `{"name":"Persisted Co","status":"active","owner":"zoe"}`)
	if code != http.StatusCreated {
		t.Fatalf("create status = %d, want 201", code)
	}
	newID, _ := created["id"].(string)
	if newID == "" {
		t.Fatalf("created record id missing")
	}

	// update (PATCH /api/records/rec-1) → 200 + refreshed updatedAt.
	code, updated := authedJSON(t, mux1, token1, http.MethodPatch, "/api/records/rec-1", `{"name":"Acme Rebrand"}`)
	if code != http.StatusOK {
		t.Fatalf("update status = %d, want 200", code)
	}
	updatedAt, _ := updated["updatedAt"].(string)
	if updatedAt == "" {
		t.Fatalf("updated record updatedAt missing")
	}

	// delete (DELETE /api/records/rec-2) → 204.
	code, _ = authedJSON(t, mux1, token1, http.MethodDelete, "/api/records/rec-2", "")
	if code != http.StatusNoContent {
		t.Fatalf("delete status = %d, want 204", code)
	}

	if err := st1.Close(); err != nil {
		t.Fatalf("close phase 1: %v", err)
	}

	// Phase 2: restart on the same DB file. A different seed hash must not
	// overwrite the admin password or re-seed anything.
	st2, err := store.Open(path, testSeedUsername, "hash-different", true)
	if err != nil {
		t.Fatalf("open phase 2: %v", err)
	}
	defer st2.Close()
	a2 := auth.New([]byte(testJWTSecret), 15*time.Minute, 30*24*time.Hour, st2, false)
	mux2 := http.NewServeMux()
	Register(mux2, a2, st2)
	env2 := &authTestEnv{mux: mux2, a: a2, st: st2}

	// list after restart: created present, rec-1 updated, rec-2 absent,
	// total = 8 (8 seed − 1 delete + 1 create; non-empty table is not re-seeded).
	code, list := getRecords(t, env2, "/api/records?pageSize=100")
	if code != http.StatusOK {
		t.Fatalf("list after restart status = %d, want 200", code)
	}
	items := list["items"].([]any)
	present := map[string]string{}
	for _, raw := range items {
		item := raw.(map[string]any)
		id, _ := item["id"].(string)
		name, _ := item["name"].(string)
		present[id] = name
	}
	if _, ok := present[newID]; !ok {
		t.Fatalf("created record %s missing after restart; list=%v", newID, present)
	}
	if present["rec-1"] != "Acme Rebrand" {
		t.Fatalf("rec-1 name after restart = %q, want Acme Rebrand", present["rec-1"])
	}
	if _, ok := present["rec-2"]; ok {
		t.Fatalf("deleted rec-2 resurrected after restart")
	}
	if total, _ := list["total"].(float64); total != 8 {
		t.Fatalf("total after restart = %v, want 8 (no re-seed)", total)
	}

	// detail after restart: created record fields and rec-1 updatedAt (ms) persist.
	_, detail := getRecords(t, env2, "/api/records/"+newID)
	if detail["name"] != "Persisted Co" || detail["status"] != "active" || detail["owner"] != "zoe" {
		t.Fatalf("created detail after restart = %v", detail)
	}
	_, rec1 := getRecords(t, env2, "/api/records/rec-1")
	if rec1["name"] != "Acme Rebrand" {
		t.Fatalf("rec-1 detail name after restart = %v", rec1)
	}
	if rec1["updatedAt"] != updatedAt {
		t.Fatalf("rec-1 updatedAt after restart = %v, want %v (persisted ms)", rec1["updatedAt"], updatedAt)
	}
}

// authedJSON performs an authenticated JSON request against a handler mux.
func authedJSON(t *testing.T, mux *http.ServeMux, token, method, path, body string) (int, map[string]any) {
	t.Helper()
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, bearer(t, token, method, path, body))
	var out map[string]any
	if rr.Body.Len() > 0 && rr.Code != http.StatusNoContent {
		_ = json.NewDecoder(rr.Body).Decode(&out)
	}
	return rr.Code, out
}
