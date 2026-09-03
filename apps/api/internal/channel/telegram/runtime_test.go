package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
)

func TestRuntimeManager_HotSwitch(t *testing.T) {
	mock := NewCaptureSender()
	rm, err := NewRuntimeManager("initial-token-12345", "initial-secret-67890", mock, testMasterKey())
	if err != nil {
		t.Fatalf("NewRuntimeManager: %v", err)
	}

	// Verify initial getters
	if rm.GetToken() != "initial-token-12345" {
		t.Fatalf("unexpected token: %s", rm.GetToken())
	}
	if rm.GetSecret() != "initial-secret-67890" {
		t.Fatalf("unexpected secret: %s", rm.GetSecret())
	}

	// Verify Status
	status := rm.Status()
	if !status.Configured || !status.TokenSet || !status.SecretSet {
		t.Fatalf("expected configured=true, token_set=true, secret_set=true, got %+v", status)
	}

	// Concurrent hot switch
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(2)
		go func(idx int) {
			defer wg.Done()
			_ = rm.Update(context.Background(), "new-token-abc", "new-secret-xyz")
		}(i)
		go func(idx int) {
			defer wg.Done()
			_ = rm.Status()
			_ = rm.GetToken()
			_ = rm.GetSecret()
		}(i)
	}
	wg.Wait()

	if rm.GetToken() != "new-token-abc" {
		t.Fatalf("expected final token 'new-token-abc', got %q", rm.GetToken())
	}
	if rm.GetSecret() != "new-secret-xyz" {
		t.Fatalf("expected final secret 'new-secret-xyz', got %q", rm.GetSecret())
	}
}

func TestSettingsHandler_AuthenticationAndPermissions(t *testing.T) {
	rm, err := NewRuntimeManager("token123456", "secret789012", nil, testMasterKey())
	if err != nil {
		t.Fatalf("NewRuntimeManager: %v", err)
	}
	handler := NewSettingsHandler(rm)

	// 1. Unauthenticated GET -> 401
	reqUnauth := httptest.NewRequest(http.MethodGet, "/api/channel/telegram/settings", nil)
	wUnauth := httptest.NewRecorder()
	handler.ServeHTTP(wUnauth, reqUnauth)
	if wUnauth.Code != http.StatusUnauthorized {
		t.Fatalf("unauthenticated expected 401, got %d", wUnauth.Code)
	}

	// 2. Unauthenticated PATCH -> 401
	reqPatchUnauth := httptest.NewRequest(http.MethodPatch, "/api/channel/telegram/settings", bytes.NewReader([]byte(`{"bot_token":"hacked"}`)))
	wPatchUnauth := httptest.NewRecorder()
	handler.ServeHTTP(wPatchUnauth, reqPatchUnauth)
	if wPatchUnauth.Code != http.StatusUnauthorized {
		t.Fatalf("unauthenticated PATCH expected 401, got %d", wPatchUnauth.Code)
	}

	// 3. User lacking settings.read -> 403
	userNoPerm := account.User{ID: "user-1", Permissions: []string{"users.read"}}
	reqForbidden := httptest.NewRequest(http.MethodGet, "/api/channel/telegram/settings", nil)
	reqForbidden = reqForbidden.WithContext(auth.WithIdentity(reqForbidden.Context(), userNoPerm))
	wForbidden := httptest.NewRecorder()
	handler.ServeHTTP(wForbidden, reqForbidden)
	if wForbidden.Code != http.StatusForbidden {
		t.Fatalf("forbidden expected 403, got %d", wForbidden.Code)
	}

	// 4. User with settings.read only attempting PATCH -> 403
	userReadOnly := account.User{ID: "user-2", Permissions: []string{"settings.read"}}
	reqPatchForbidden := httptest.NewRequest(http.MethodPatch, "/api/channel/telegram/settings", bytes.NewReader([]byte(`{"bot_token":"hacked"}`)))
	reqPatchForbidden = reqPatchForbidden.WithContext(auth.WithIdentity(reqPatchForbidden.Context(), userReadOnly))
	wPatchForbidden := httptest.NewRecorder()
	handler.ServeHTTP(wPatchForbidden, reqPatchForbidden)
	if wPatchForbidden.Code != http.StatusForbidden {
		t.Fatalf("read-only PATCH expected 403, got %d", wPatchForbidden.Code)
	}

	// 5. Admin with settings.read and settings.write:
	admin := account.User{ID: "admin-1", Permissions: []string{"settings.read", "settings.write"}}

	// GET settings returns masked secrets
	reqGet := httptest.NewRequest(http.MethodGet, "/api/channel/telegram/settings", nil)
	reqGet = reqGet.WithContext(auth.WithIdentity(reqGet.Context(), admin))
	wGet := httptest.NewRecorder()
	handler.ServeHTTP(wGet, reqGet)

	if wGet.Code != http.StatusOK {
		t.Fatalf("GET expected 200, got %d", wGet.Code)
	}

	var status RuntimeStatus
	if err := json.NewDecoder(wGet.Body).Decode(&status); err != nil {
		t.Fatalf("decode GET response: %v", err)
	}
	if !status.Configured || !status.TokenSet || !status.SecretSet {
		t.Fatalf("unexpected status: %+v", status)
	}

	// PATCH settings updates secrets
	patchBody := `{"bot_token":"updated-token-9999"}`
	reqPatch := httptest.NewRequest(http.MethodPatch, "/api/channel/telegram/settings", bytes.NewReader([]byte(patchBody)))
	reqPatch = reqPatch.WithContext(auth.WithIdentity(reqPatch.Context(), admin))
	wPatch := httptest.NewRecorder()
	handler.ServeHTTP(wPatch, reqPatch)

	if wPatch.Code != http.StatusOK {
		t.Fatalf("PATCH expected 200, got %d", wPatch.Code)
	}

	// Check runtime state
	if rm.GetToken() != "updated-token-9999" {
		t.Fatalf("expected updated token, got %s", rm.GetToken())
	}
	if rm.GetSecret() != "secret789012" {
		t.Fatalf("expected secret to remain unchanged, got %s", rm.GetSecret())
	}
}

// testMasterKey returns a fixed 32-byte at-rest key for tests (F-002: the
// production path never falls back to a source constant).
func testMasterKey() []byte {
	return []byte("0123456789abcdef0123456789abcdef")
}

// newTestRuntimeManager constructs a RuntimeManager with a fixed test key,
// failing the test on construction error (no persistence runner in these tests).
func newTestRuntimeManager(t *testing.T, seedToken, seedSecret string, mock *CaptureSender) *RuntimeManager {
	t.Helper()
	rm, err := NewRuntimeManager(seedToken, seedSecret, mock, testMasterKey())
	if err != nil {
		t.Fatalf("NewRuntimeManager: %v", err)
	}
	return rm
}
