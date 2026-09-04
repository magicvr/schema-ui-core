package telegram

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
)

func TestLeaseHandler_AuthSessionAndLifecycle(t *testing.T) {
	runtime := newTestRuntimeManager(t, "bot-token", "webhook-secret", nil)
	manager := NewConnectionManager(runtime, nil, nil, nil, nil)
	handler := NewLeaseHandler(manager)

	tests := []struct {
		name       string
		path       string
		method     string
		identity   *account.User
		wantStatus int
	}{
		{name: "unauthenticated", path: "/api/channel/telegram/lease/acquire", method: http.MethodPost, wantStatus: http.StatusUnauthorized},
		{
			name:       "missing permission",
			path:       "/api/channel/telegram/lease/acquire",
			method:     http.MethodPost,
			identity:   &account.User{ID: "user-1", SessionID: "session-1", Permissions: []string{"users.read"}},
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "missing session does not fall back to user id",
			path:       "/api/channel/telegram/lease/acquire",
			method:     http.MethodPost,
			identity:   &account.User{ID: "user-1", Permissions: []string{"settings.read"}},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "wrong method",
			path:       "/api/channel/telegram/lease/acquire",
			method:     http.MethodGet,
			identity:   &account.User{ID: "user-1", SessionID: "session-1", Permissions: []string{"settings.read"}},
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "unknown operation",
			path:       "/api/channel/telegram/lease/unknown",
			method:     http.MethodPost,
			identity:   &account.User{ID: "user-1", SessionID: "session-1", Permissions: []string{"settings.read"}},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			if tt.identity != nil {
				req = req.WithContext(auth.WithIdentity(req.Context(), *tt.identity))
			}
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d: %s", rr.Code, tt.wantStatus, rr.Body.String())
			}
		})
	}

	admin := account.User{ID: "admin-1", SessionID: "browser-session-1", Permissions: []string{"settings.read"}}
	for _, operation := range []string{"acquire", "heartbeat", "release"} {
		req := httptest.NewRequest(http.MethodPost, "/api/channel/telegram/lease/"+operation, nil)
		req = req.WithContext(auth.WithIdentity(req.Context(), admin))
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("%s status = %d, want 200: %s", operation, rr.Code, rr.Body.String())
		}
		var body leaseResponse
		if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
			t.Fatalf("decode %s response: %v", operation, err)
		}
		wantActive := operation != "release"
		if body.Active != wantActive || body.TTLSeconds != int(PollingLeaseTTL.Seconds()) {
			t.Fatalf("%s response = %+v", operation, body)
		}
	}
	if got := manager.ActiveLeaseCount(); got != 0 {
		t.Fatalf("released session remains active: %d", got)
	}

	// A second authenticated session must be independent from the first one.
	second := admin
	second.SessionID = "browser-session-2"
	req := httptest.NewRequest(http.MethodPost, "/api/channel/telegram/lease/acquire", nil)
	req = req.WithContext(auth.WithIdentity(req.Context(), second))
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK || manager.ActiveLeaseCount() != 1 {
		t.Fatalf("second session acquire status=%d active=%d", rr.Code, manager.ActiveLeaseCount())
	}
}
