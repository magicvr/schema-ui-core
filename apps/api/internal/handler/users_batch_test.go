package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// I-PROTO-FULL-001 · D-ACT/D-TABLE batch (ADR-0022): POST {path}/batch-delete
// is the generic server-side batch contract over normalized $selection.keys.
func TestUsersBatchDelete(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)

	create := func(username string) string {
		body := map[string]any{"username": username, "name": "Batch " + username, "password": "secret-pass-1"}
		raw, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewReader(raw))
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()
		env.mux.ServeHTTP(resp, req)
		if resp.Code != http.StatusCreated {
			t.Fatalf("create %s status = %d: %s", username, resp.Code, resp.Body.String())
		}
		var created map[string]any
		if err := json.Unmarshal(resp.Body.Bytes(), &created); err != nil {
			t.Fatal(err)
		}
		id, _ := created["id"].(string)
		if id == "" {
			t.Fatalf("create %s returned no id", username)
		}
		return id
	}

	idA := create("batch-a")
	idB := create("batch-b")
	create("batch-c")

	// Normalized selection: duplicates are dropped (D3/V274). Non-scalar keys
	// are rejected server-side (INVALID_SELECTION_KEY, see FailClosed test);
	// the client normalizes before sending, so only scalars reach the wire.
	payload := map[string]any{"ids": []any{idA, idA, idB, idB}}
	raw, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/users/batch-delete", bytes.NewReader(raw))
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()
	env.mux.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fatalf("batch-delete status = %d: %s", resp.Code, resp.Body.String())
	}
	var result map[string]any
	if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result["deleted"] != float64(2) {
		t.Fatalf("deleted = %v, want 2", result["deleted"])
	}

	// Both deleted; the third remains.
	for _, id := range []string{idA, idB} {
		code, _ := getResource(t, env, "/api/users/"+id)
		if code != http.StatusNotFound {
			t.Fatalf("detail %s after batch = %d, want 404", id, code)
		}
	}
}

func TestUsersBatchDeleteFailClosed(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)

	cases := []struct {
		name string
		body string
		code string
	}{
		{"empty selection", `{"ids": []}`, "EMPTY_SELECTION"},
		{"non-scalar keys", `{"ids": [{"x": 1}, ["a"]]}`, "INVALID_SELECTION_KEY"},
		{"empty string key", `{"ids": [""]}`, "INVALID_SELECTION_KEY"},
		{"malformed body", `not-json`, "INVALID_BODY"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/users/batch-delete", bytes.NewReader([]byte(tc.body)))
			req.Header.Set("Authorization", "Bearer "+token)
			resp := httptest.NewRecorder()
			env.mux.ServeHTTP(resp, req)
			if resp.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400: %s", resp.Code, resp.Body.String())
			}
			var apiError map[string]string
			if err := json.Unmarshal(resp.Body.Bytes(), &apiError); err != nil {
				t.Fatal(err)
			}
			if apiError["error"] != tc.code {
				t.Fatalf("error code = %q, want %q", apiError["error"], tc.code)
			}
		})
	}

	// Anonymous → 401; non-writer → 403.
	req := httptest.NewRequest(http.MethodPost, "/api/users/batch-delete", bytes.NewReader([]byte(`{"ids": ["x"]}`)))
	resp := httptest.NewRecorder()
	env.mux.ServeHTTP(resp, req)
	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("anonymous status = %d, want 401", resp.Code)
	}
}

// D-001 P0 · whole-batch atomicity at the HTTP layer: a selection containing a
// protected target (the last admin) fails with LAST_ADMIN and NO user is
// deleted — the batch must not partially commit. A non-admin actor with
// users.write is used so the batch can legally target user-admin without
// tripping the self-guard.
func TestUsersBatchDeleteAtomicRollbackHTTP(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)
	now := time.Now().UTC()
	if _, err := env.authRepository.CreateRoleWithGrants(
		"users-manager", "Users manager",
		[]string{"roles.assign", "roles.read", "users.read", "users.write"}, nil, now,
	); err != nil {
		t.Fatalf("create users-manager role: %v", err)
	}
	env.addUser(t, "um", "um-password", []string{"users-manager"})
	umToken := env.login(t, "um", "um-password")

	create := func(username string) string {
		body := map[string]any{"username": username, "name": "Batch " + username, "password": "secret-pass-1"}
		raw, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewReader(raw))
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()
		env.mux.ServeHTTP(resp, req)
		if resp.Code != http.StatusCreated {
			t.Fatalf("create %s status = %d: %s", username, resp.Code, resp.Body.String())
		}
		var created map[string]any
		if err := json.Unmarshal(resp.Body.Bytes(), &created); err != nil {
			t.Fatal(err)
		}
		id, _ := created["id"].(string)
		if id == "" {
			t.Fatalf("create %s returned no id", username)
		}
		return id
	}

	idA := create("atomic-a")
	idB := create("atomic-b")

	// ["user-admin", idA, idB]: the last-admin guard fires in the middle of the
	// batch. With whole-batch semantics nothing is deleted (409 LAST_ADMIN).
	payload := map[string]any{"ids": []any{"user-admin", idA, idB}}
	raw, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/users/batch-delete", bytes.NewReader(raw))
	req.Header.Set("Authorization", "Bearer "+umToken)
	resp := httptest.NewRecorder()
	env.mux.ServeHTTP(resp, req)
	if resp.Code != http.StatusConflict {
		t.Fatalf("batch status = %d, want 409: %s", resp.Code, resp.Body.String())
	}
	var apiError map[string]string
	if err := json.Unmarshal(resp.Body.Bytes(), &apiError); err != nil {
		t.Fatal(err)
	}
	if apiError["error"] != "LAST_ADMIN" {
		t.Fatalf("error code = %q, want LAST_ADMIN", apiError["error"])
	}

	// Neither earlier key in the selection was committed.
	for _, id := range []string{idA, idB} {
		code, _ := getResource(t, env, "/api/users/"+id)
		if code != http.StatusOK {
			t.Fatalf("detail %s after rolled-back batch = %d, want 200", id, code)
		}
	}
	code, _ := getResource(t, env, "/api/users/user-admin")
	if code != http.StatusOK {
		t.Fatalf("admin after rolled-back batch = %d, want 200", code)
	}
}

// A-002 F-001 regression (HTTP): a batch containing EVERY admin fails with
// LAST_ADMIN and nothing is deleted — the batch-level guard must reject the
// all-admins selection instead of leaving zero admins. A non-admin actor with
// users.write is used so the batch can target every admin without tripping the
// self-guard.
func TestUsersBatchDeleteRejectsRemovingAllAdminsHTTP(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)
	now := time.Now().UTC()
	if _, err := env.authRepository.CreateRoleWithGrants(
		"users-manager", "Users manager",
		[]string{"users.read", "users.write"}, nil, now,
	); err != nil {
		t.Fatalf("create users-manager role: %v", err)
	}
	env.addUser(t, "um", "um-password", []string{"users-manager"})
	umToken := env.login(t, "um", "um-password")

	// Create a second admin through the API (admin can create admin).
	req := bearer(t, token, http.MethodPost, "/api/users",
		`{"username":"admin2","name":"Admin 2","password":"admin2-pass-1","roles":["admin","editor"]}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("create admin2 = %d: %s", rr.Code, rr.Body.String())
	}
	var created map[string]any
	_ = json.NewDecoder(rr.Body).Decode(&created)
	admin2ID, _ := created["id"].(string)
	if admin2ID == "" {
		t.Fatalf("admin2 id missing: %v", created)
	}

	// Both admins in the selection → 409 LAST_ADMIN, nothing deleted.
	payload := map[string]any{"ids": []any{"user-admin", admin2ID}}
	raw, _ := json.Marshal(payload)
	req = httptest.NewRequest(http.MethodPost, "/api/users/batch-delete", bytes.NewReader(raw))
	req.Header.Set("Authorization", "Bearer "+umToken)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusConflict {
		t.Fatalf("all-admins batch status = %d, want 409: %s", rr.Code, rr.Body.String())
	}
	var apiError map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &apiError); err != nil {
		t.Fatal(err)
	}
	if apiError["error"] != "LAST_ADMIN" {
		t.Fatalf("error code = %q, want LAST_ADMIN", apiError["error"])
	}
	for _, id := range []string{"user-admin", admin2ID} {
		code, _ := getResource(t, env, "/api/users/"+id)
		if code != http.StatusOK {
			t.Fatalf("detail %s after rejected all-admins batch = %d, want 200", id, code)
		}
	}

	// With a third admin outside the batch, the same selection is legal.
	req = bearer(t, token, http.MethodPost, "/api/users",
		`{"username":"admin3","name":"Admin 3","password":"admin3-pass-1","roles":["admin"]}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("create admin3 = %d: %s", rr.Code, rr.Body.String())
	}
	req = httptest.NewRequest(http.MethodPost, "/api/users/batch-delete", bytes.NewReader(raw))
	req.Header.Set("Authorization", "Bearer "+umToken)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("all-admins batch with survivor status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	// The batch deleted user-admin itself, so assert via the non-admin actor
	// (getResource logs in as admin, who no longer exists).
	for _, id := range []string{"user-admin", admin2ID} {
		req = httptest.NewRequest(http.MethodGet, "/api/users/"+id, nil)
		req.Header.Set("Authorization", "Bearer "+umToken)
		rr = httptest.NewRecorder()
		env.mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Fatalf("detail %s after legal batch = %d, want 404", id, rr.Code)
		}
	}
}
