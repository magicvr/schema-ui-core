// F-04 notification tests (GOAL-006 S3): list/read/read-all/unread/settings +
// best-effort system-event hooks (lock/disable/unlock/password-change).
package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	authsession "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession"
)

func notificationsList(t *testing.T, env *authTestEnv, path string) (int, map[string]any) {
	t.Helper()
	req := bearer(t, adminToken(t, env), http.MethodGet, path, "")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	var out map[string]any
	if rr.Body.Len() > 0 {
		_ = json.NewDecoder(rr.Body).Decode(&out)
	}
	return rr.Code, out
}

func TestNotificationLockEventProduced(t *testing.T) {
	env := newAuthTestEnv(t)
	env.addUser(t, "editor1", "editor-password", []string{"editor"})
	// 5 failed logins open the lock window → account.locked notification.
	for i := 0; i < 5; i++ {
		sendJSON(t, env.mux, http.MethodPost, "/api/auth/login", `{"username":"editor1","password":"wrong-pass"}`)
	}
	// editor is locked → cannot login; check via repository directly.
	rows, _, err := env.authRepository.ListNotifications("user-editor1", authsession.NotificationFilter{Page: 1, PageSize: 20})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].Event != "account.locked" {
		t.Fatalf("locked notification rows = %+v, want 1 account.locked", rows)
	}
}

func TestNotificationDisableAndUnlockEvents(t *testing.T) {
	env := newAuthTestEnv(t)
	env.addUser(t, "editor1", "editor-password", []string{"editor"})
	token := adminToken(t, env)
	// disable → account.disabled (transition only; re-disable is a no-op)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, token, http.MethodPost, "/api/users/user-editor1/disable", ""))
	if rr.Code != http.StatusNoContent {
		t.Fatalf("disable = %d", rr.Code)
	}
	rrDup := httptest.NewRecorder()
	env.mux.ServeHTTP(rrDup, bearer(t, token, http.MethodPost, "/api/users/user-editor1/disable", ""))
	if rrDup.Code != http.StatusNoContent {
		t.Fatalf("re-disable = %d", rrDup.Code)
	}
	// enable → lock via 5 failures → unlock → account.unlocked
	rr3 := httptest.NewRecorder()
	env.mux.ServeHTTP(rr3, bearer(t, token, http.MethodPost, "/api/users/user-editor1/enable", ""))
	if rr3.Code != http.StatusNoContent {
		t.Fatalf("enable = %d", rr3.Code)
	}
	for i := 0; i < 5; i++ {
		sendJSON(t, env.mux, http.MethodPost, "/api/auth/login", `{"username":"editor1","password":"wrong-pass"}`)
	}
	rr4 := httptest.NewRecorder()
	env.mux.ServeHTTP(rr4, bearer(t, token, http.MethodPost, "/api/users/user-editor1/unlock", ""))
	if rr4.Code != http.StatusNoContent {
		t.Fatalf("unlock = %d", rr4.Code)
	}
	rows, _, err := env.authRepository.ListNotifications("user-editor1", authsession.NotificationFilter{Page: 1, PageSize: 20})
	if err != nil {
		t.Fatal(err)
	}
	events := map[string]bool{}
	for _, n := range rows {
		events[n.Event] = true
	}
	if !events["account.unlocked"] || !events["account.disabled"] || !events["account.locked"] {
		t.Fatalf("events = %v, want {account.unlocked account.disabled account.locked}", events)
	}
	// re-disable produced no duplicate (transition guard)
	disabled := 0
	for _, n := range rows {
		if n.Event == "account.disabled" {
			disabled++
		}
	}
	if disabled != 1 {
		t.Fatalf("account.disabled count = %d, want 1 (no duplicate on re-disable)", disabled)
	}
}

func TestNotificationPasswordChangeEvents(t *testing.T) {
	env := newAuthTestEnv(t)
	env.addUser(t, "editor1", "editor-password", []string{"editor"})
	token := adminToken(t, env)
	// admin resets editor's password → account.password-changed
	req := bearer(t, token, http.MethodPatch, "/api/users/user-editor1", `{"password":"brand-new-pass"}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("admin reset password = %d: %s", rr.Code, rr.Body.String())
	}
	rows, _, err := env.authRepository.ListNotifications("user-editor1", authsession.NotificationFilter{Page: 1, PageSize: 20})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].Event != "account.password-changed" {
		t.Fatalf("password-changed notification rows = %+v", rows)
	}
	// editor self password change → another notification for self
	env2 := env
	editorToken := env2.login(t, "editor1", "brand-new-pass")
	req2 := bearer(t, editorToken, http.MethodPost, "/api/account/password", `{"currentPassword":"brand-new-pass","newPassword":"another-new-pass"}`)
	rr2 := httptest.NewRecorder()
	env2.mux.ServeHTTP(rr2, req2)
	if rr2.Code != http.StatusNoContent {
		t.Fatalf("self password change = %d", rr2.Code)
	}
	rows2, _, _ := env2.authRepository.ListNotifications("user-editor1", authsession.NotificationFilter{Page: 1, PageSize: 20})
	if len(rows2) != 2 || rows2[0].Event != "account.password-changed" {
		t.Fatalf("self password-changed rows = %+v", rows2)
	}
}

func TestNotificationListReadReadAllUnread(t *testing.T) {
	env := newAuthTestEnv(t)
	env.addUser(t, "editor1", "editor-password", []string{"editor"})
	// produce 2 notifications directly
	now := time.Now().UTC()
	_ = env.authRepository.CreateNotification(authsession.Notification{ID: "ntf-a", UserID: "user-editor1", Event: "account.locked", Title: "A", Body: "B"}, now)
	_ = env.authRepository.CreateNotification(authsession.Notification{ID: "ntf-b", UserID: "user-editor1", Event: "account.disabled", Title: "C", Body: "D"}, now.Add(time.Second))
	token := env.login(t, "editor1", "editor-password")
	// list
	req := bearer(t, token, http.MethodGet, "/api/notifications", "")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("list = %d", rr.Code)
	}
	var list struct {
		Items []map[string]any `json:"items"`
		Total int              `json:"total"`
	}
	_ = json.NewDecoder(rr.Body).Decode(&list)
	if list.Total != 2 || list.Items[0]["id"] != "ntf-b" || list.Items[0]["read"] != false {
		t.Fatalf("list = %+v", list)
	}
	// unread-count
	req = bearer(t, token, http.MethodGet, "/api/notifications/unread-count", "")
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	var unread map[string]any
	_ = json.NewDecoder(rr.Body).Decode(&unread)
	if unread["unread"].(float64) != 2 {
		t.Fatalf("unread = %v", unread)
	}
	// read one → the config-change header lets the shell badge refresh at once
	req = bearer(t, token, http.MethodPost, "/api/notifications/ntf-a/read", "")
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("read = %d", rr.Code)
	}
	if got := rr.Header().Get("X-Schema-UI-Config-Changed"); got != "notifications.read" {
		t.Fatalf("read config-changed header = %q, want notifications.read", got)
	}
	// foreign id → 404
	req = bearer(t, token, http.MethodPost, "/api/notifications/nope/read", "")
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("foreign read = %d, want 404", rr.Code)
	}
	// read-all → same immediate-refresh header
	req = bearer(t, token, http.MethodPost, "/api/notifications/read-all", "")
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	var ra map[string]any
	_ = json.NewDecoder(rr.Body).Decode(&ra)
	if ra["updated"].(float64) != 1 {
		t.Fatalf("read-all updated = %v, want 1", ra)
	}
	if got := rr.Header().Get("X-Schema-UI-Config-Changed"); got != "notifications.read" {
		t.Fatalf("read-all config-changed header = %q, want notifications.read", got)
	}
	// unread-only filter → 0
	req = bearer(t, token, http.MethodGet, "/api/notifications?unreadOnly=true", "")
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	var ul map[string]any
	_ = json.NewDecoder(rr.Body).Decode(&ul)
	if ul["total"].(float64) != 0 {
		t.Fatalf("unreadOnly total = %v, want 0", ul)
	}
	// settings toggle off → new notifications suppressed
	req = bearer(t, token, http.MethodPatch, "/api/notifications/settings", `{"enabled":false}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("settings = %d", rr.Code)
	}
	NotifyAccountEvent(env.authRepository, "user-editor1", "account.locked", now.Add(2*time.Second))
	rows, _, _ := env.authRepository.ListNotifications("user-editor1", authsession.NotificationFilter{Page: 1, PageSize: 100})
	if len(rows) != 2 {
		t.Fatalf("suppressed notification produced: rows = %d", len(rows))
	}
}

func TestNotificationPruneKeepsUnread(t *testing.T) {
	env := newAuthTestEnv(t)
	env.addUser(t, "editor1", "editor-password", []string{"editor"})
	now := time.Now().UTC()
	// 500 read + 1 unread for editor1; a new notification should prune the
	// oldest READ but keep the unread one.
	for i := 0; i < 500; i++ {
		_ = env.authRepository.CreateNotification(authsession.Notification{ID: "ntf-r" + fmt.Sprintf("%03d", i), UserID: "user-editor1", Event: "account.locked", Title: "R", Body: "R"}, now.Add(-time.Duration(i)*time.Second))
	}
	// mark them read
	_, _ = env.authRepository.MarkAllNotificationsRead("user-editor1", now)
	_ = env.authRepository.CreateNotification(authsession.Notification{ID: "ntf-unread", UserID: "user-editor1", Event: "account.disabled", Title: "U", Body: "U"}, now)
	// produce one more → prune kicks in
	_ = env.authRepository.CreateNotification(authsession.Notification{ID: "ntf-new", UserID: "user-editor1", Event: "account.unlocked", Title: "N", Body: "N"}, now)
	rows, total, err := env.authRepository.ListNotifications("user-editor1", authsession.NotificationFilter{Page: 1, PageSize: 1000})
	if err != nil {
		t.Fatal(err)
	}
	if total > 501 {
		t.Fatalf("prune failed: total = %d", total)
	}
	unread := 0
	for _, n := range rows {
		if n.ID == "ntf-unread" || n.ID == "ntf-new" {
			if n.ReadAt == nil {
				unread++
			}
		}
	}
	if unread != 2 {
		t.Fatalf("unread rows lost: %d (want 2: ntf-unread + ntf-new)", unread)
	}
}