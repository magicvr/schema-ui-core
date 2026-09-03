package telegram

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestRuntimeManager_HotSwitch(t *testing.T) {
	mock := NewCaptureSender()
	rm := NewRuntimeManager("initial-token-12345", "initial-secret-67890", mock)

	// Verify initial getters
	if rm.GetToken() != "initial-token-12345" {
		t.Fatalf("unexpected token: %s", rm.GetToken())
	}
	if rm.GetSecret() != "initial-secret-67890" {
		t.Fatalf("unexpected secret: %s", rm.GetSecret())
	}

	// Verify Status masking (keeps last 4 chars)
	status := rm.Status()
	if !status.Configured {
		t.Fatalf("expected configured=true")
	}
	if status.TokenMasked != "***************2345" {
		t.Fatalf("unexpected masked token: %s", status.TokenMasked)
	}
	if status.SecretMasked != "****************7890" {
		t.Fatalf("unexpected masked secret: %s", status.SecretMasked)
	}

	// Concurrent hot switch
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(2)
		go func(idx int) {
			defer wg.Done()
			rm.Update("new-token-abc", "new-secret-xyz")
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

func TestSettingsHandler(t *testing.T) {
	rm := NewRuntimeManager("token123456", "secret789012", nil)
	handler := NewSettingsHandler(rm)

	// 1. GET settings returns masked secrets
	reqGet := httptest.NewRequest(http.MethodGet, "/api/channel/telegram/settings", nil)
	wGet := httptest.NewRecorder()
	handler.ServeHTTP(wGet, reqGet)

	if wGet.Code != http.StatusOK {
		t.Fatalf("GET expected 200, got %d", wGet.Code)
	}

	var status RuntimeStatus
	if err := json.NewDecoder(wGet.Body).Decode(&status); err != nil {
		t.Fatalf("decode GET response: %v", err)
	}
	if !status.Configured || status.TokenMasked != "*******3456" || status.SecretMasked != "********9012" {
		t.Fatalf("unexpected status: %+v", status)
	}

	// 2. PATCH settings updates secrets
	patchBody := `{"bot_token":"updated-token-9999"}`
	reqPatch := httptest.NewRequest(http.MethodPatch, "/api/channel/telegram/settings", bytes.NewReader([]byte(patchBody)))
	wPatch := httptest.NewRecorder()
	handler.ServeHTTP(wPatch, reqPatch)

	if wPatch.Code != http.StatusOK {
		t.Fatalf("PATCH expected 200, got %d", wPatch.Code)
	}

	// Check runtime state
	if rm.GetToken() != "updated-token-9999" {
		t.Fatalf("expected updated token, got %s", rm.GetToken())
	}
	// Webhook secret should remain untouched
	if rm.GetSecret() != "secret789012" {
		t.Fatalf("expected secret to remain unchanged, got %s", rm.GetSecret())
	}
}
