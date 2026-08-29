// Invitation response mapping tests (workspace-019): the revokable derived
// flag must track live pending invites only — terminal states disable the
// schema row action via disabledWhen.
package handler

import (
	"testing"
	"time"

	authsession "github.com/magicvr/schema-ui-core/apps/api/modules/authsession"
)

func TestInviteToMapRevokableFlag(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	email := "u@example.com"
	cases := []struct {
		name    string
		inv     authsession.Invite
		status  string
		revokable bool
	}{
		{"pending", authsession.Invite{ID: "inv-1", Roles: []string{"viewer"}, ExpiresAt: now.Add(24 * time.Hour), CreatedAt: now}, "pending", true},
		{"consumed", authsession.Invite{ID: "inv-2", Roles: []string{"viewer"}, Email: &email, ExpiresAt: now.Add(24 * time.Hour), ConsumedAt: &now, CreatedAt: now}, "consumed", false},
		{"revoked", authsession.Invite{ID: "inv-3", Roles: []string{"viewer"}, Email: &email, ExpiresAt: now.Add(24 * time.Hour), RevokedAt: &now, CreatedAt: now}, "revoked", false},
		{"expired", authsession.Invite{ID: "inv-4", Roles: []string{"viewer"}, Email: &email, ExpiresAt: now.Add(-time.Hour), CreatedAt: now}, "expired", false},
	}
	for _, tc := range cases {
		out := inviteToMap(&tc.inv)
		if out["status"] != tc.status {
			t.Errorf("%s: status = %v, want %s", tc.name, out["status"], tc.status)
		}
		if out["revokable"] != tc.revokable {
			t.Errorf("%s: revokable = %v, want %v", tc.name, out["revokable"], tc.revokable)
		}
	}
}