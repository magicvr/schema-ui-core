package store

import (
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/modules/authsession"
)

// S6 · 服务重启持久化：重新 Open 同一数据库后，迁移台账不重跑、种子不重复，
// 身份、refresh、权限与菜单投影全部保持。
func TestRestartPersistence(t *testing.T) {
	path := filepath.Join(t.TempDir(), "restart.db")
	st, err := OpenSeeded(path, "admin", "hash", true)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	repository := authsession.NewRepository(st)
	now := time.Now().UTC()
	if err := repository.CreateUser(authsession.User{
		ID: "u1", Username: "alice", Name: "Alice",
		Roles: []string{"viewer"}, PasswordHash: "h", CreatedAt: now, UpdatedAt: now,
	}); err != nil {
		t.Fatalf("create alice: %v", err)
	}
	if err := repository.CreateRefreshToken(authsession.RefreshToken{
		ID: "rt1", UserID: "u1", TokenHash: "hash-rt1",
		ExpiresAt: now.Add(time.Hour), CreatedAt: now,
	}); err != nil {
		t.Fatalf("create refresh token: %v", err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}

	// Restart with a different seed hash: nothing is re-applied or overwritten.
	st2, err := OpenSeeded(path, "admin", "hash-v2", true)
	if err != nil {
		t.Fatalf("restart open: %v", err)
	}
	defer st2.Close()
	repository2 := authsession.NewRepository(st2)

	applied, err := st2.appliedMigrations()
	if err != nil {
		t.Fatal(err)
	}
	if len(applied) != 63 || applied[52].version != 53 || applied[52].name != "operation_log_mail_events" || applied[53].version != 54 || applied[53].name != "account_email_identity" || applied[54].version != 55 || applied[54].name != "email_verification_challenges" || applied[55].version != 56 || applied[55].name != "password_recovery_challenges" || applied[56].version != 57 || applied[56].name != "password_policy" || applied[57].version != 58 || applied[57].name != "user_password_history" || applied[58].version != 59 || applied[58].name != "user_invites" || applied[59].version != 60 || applied[59].name != "mail_outbox_channels" || applied[60].version != 61 || applied[60].name != "login_failures" || applied[61].version != 62 || applied[61].name != "site_default_currency" || applied[62].version != 63 || applied[62].name != "site_settings_updated_at_index" {
		t.Fatalf("applied after restart = %+v, want 63 ending in site_settings_updated_at_index", applied)
	}
	applied = applied[:42]
	if len(applied) != 42 || applied[0].version != 1 || applied[1].version != 2 || applied[2].version != 3 || applied[3].version != 4 || applied[4].version != 5 || applied[5].version != 6 || applied[6].version != 7 || applied[7].version != 8 || applied[8].version != 9 || applied[9].version != 10 || applied[10].version != 11 || applied[11].version != 12 || applied[12].version != 13 || applied[13].version != 14 || applied[14].version != 15 || applied[15].version != 16 || applied[16].version != 17 || applied[17].version != 18 || applied[18].version != 19 || applied[19].version != 20 || applied[20].version != 21 || applied[21].version != 22 || applied[22].version != 23 || applied[23].version != 24 || applied[24].version != 25 || applied[25].version != 26 || applied[26].version != 27 || applied[27].version != 28 || applied[28].version != 29 || applied[29].version != 30 || applied[30].version != 31 || applied[31].version != 32 || applied[32].version != 33 || applied[33].version != 34 || applied[34].version != 35 || applied[35].version != 36 || applied[36].version != 37 || applied[36].name != "notifications_message_keys" || applied[37].version != 38 || applied[37].name != "must_change_password" || applied[38].version != 39 || applied[38].name != "dict_entry_badge_style" || applied[39].version != 40 || applied[39].name != "site_footer" || applied[40].version != 41 || applied[40].name != "operation_log_correlation" || applied[41].version != 42 || applied[41].name != "async_jobs" {
		t.Fatalf("applied after restart = %+v, want {1..42} (no re-migration)", applied)
	}
	var ur int
	if err := st2.db.QueryRow(`SELECT COUNT(*) FROM user_roles WHERE user_id = 'user-admin'`).Scan(&ur); err != nil || ur != 2 {
		t.Fatalf("seed user_roles after restart = %d, err %v, want 2 (no duplicate seed)", ur, err)
	}

	alice, err := repository2.UserByID("u1")
	if err != nil {
		t.Fatalf("alice after restart: %v", err)
	}
	if !reflect.DeepEqual(alice.Roles, []string{"viewer"}) {
		t.Fatalf("alice roles = %v, want [viewer]", alice.Roles)
	}
	rt, err := repository2.RefreshTokenByHash("hash-rt1")
	if err != nil {
		t.Fatalf("refresh token after restart: %v", err)
	}
	if rt.UserID != "u1" || rt.RevokedAt != nil {
		t.Fatalf("refresh token after restart = %+v", rt)
	}

	// viewer: read-only permission, no menu grant.
	perms, err := repository2.PermissionsForUser("u1")
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(perms, "users.read") || slices.Contains(perms, "users.write") {
		t.Fatalf("viewer permissions after restart = %v", perms)
	}
	feat, err := repository2.FeaturesForUser("u1")
	if err != nil {
		t.Fatal(err)
	}
	if feat["menu_users"] {
		t.Fatalf("viewer menu feature after restart = true, want false")
	}
	admFeat, err := repository2.FeaturesForUser("user-admin")
	if err != nil {
		t.Fatal(err)
	}
	if !admFeat["menu_users"] {
		t.Fatalf("admin menu feature lost after restart")
	}
	// The seed must not overwrite the admin password.
	admin, err := repository2.UserByUsername("admin")
	if err != nil {
		t.Fatal(err)
	}
	if admin.PasswordHash != "hash" {
		t.Fatalf("admin password after restart = %q, want hash (seed must not overwrite)", admin.PasswordHash)
	}
}

// S6 · 迁移前副本恢复：将 pre-v0002 快照复制到新路径并重新 Open（重跑迁移）
// 可恢复原始身份、refresh 与 RBAC/权限/菜单投影，且通过完整性校验。
func TestRestorePreV0002Snapshot(t *testing.T) {
	orig := filepath.Join(t.TempDir(), "orig.db")
	createR2Fixture(t, orig) // user-admin (roles admin+editor) + refresh token rt1/abc123
	st, err := OpenSeeded(orig, "admin", "hash", false)
	if err != nil {
		t.Fatalf("open orig: %v", err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}

	snaps, err := filepath.Glob(orig + ".pre-v0002-*.sqlite")
	if err != nil || len(snaps) != 1 {
		t.Fatalf("snapshots = %v (err %v), want exactly 1", snaps, err)
	}
	data, err := os.ReadFile(snaps[0])
	if err != nil {
		t.Fatal(err)
	}

	restored := filepath.Join(t.TempDir(), "restored.db")
	if err := os.WriteFile(restored, data, 0o644); err != nil {
		t.Fatal(err)
	}
	r, err := OpenSeeded(restored, "admin", "hash", true)
	if err != nil {
		t.Fatalf("open restored: %v", err)
	}
	defer r.Close()
	repository := authsession.NewRepository(r)

	u, err := repository.UserByUsername("admin")
	if err != nil {
		t.Fatalf("restored identity: %v", err)
	}
	if u.ID != "user-admin" || u.PasswordHash != "hash-v1" {
		t.Fatalf("restored user = %+v, want user-admin / hash-v1", u)
	}
	if want := []string{"admin", "editor"}; !reflect.DeepEqual(u.Roles, want) {
		t.Fatalf("restored roles = %v, want %v", u.Roles, want)
	}
	rt, err := repository.RefreshTokenByHash("abc123")
	if err != nil {
		t.Fatalf("restored refresh: %v", err)
	}
	if rt.UserID != "user-admin" {
		t.Fatalf("restored refresh user = %q, want user-admin", rt.UserID)
	}
	perms, err := repository.PermissionsForUser("user-admin")
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(perms, "users.read") || !slices.Contains(perms, "users.write") {
		t.Fatalf("restored permissions = %v", perms)
	}
	feat, err := repository.FeaturesForUser("user-admin")
	if err != nil {
		t.Fatal(err)
	}
	if !feat["menu_users"] {
		t.Fatalf("restored menu feature missing")
	}
	if err := r.verifyIntegrity(); err != nil {
		t.Fatalf("restored integrity: %v", err)
	}
}
