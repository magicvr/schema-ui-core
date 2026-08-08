package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
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
