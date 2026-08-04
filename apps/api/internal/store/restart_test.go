package store

import (
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"testing"
	"time"
)

// S6 · 服务重启持久化：重新 Open 同一数据库后，迁移台账不重跑、种子不重复，
// 身份、refresh、权限与菜单投影全部保持。
func TestRestartPersistence(t *testing.T) {
	path := filepath.Join(t.TempDir(), "restart.db")
	st, err := Open(path, "admin", "hash", true)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	now := time.Now().UTC()
	if err := st.CreateUser(User{
		ID: "u1", Username: "alice", Name: "Alice",
		Roles: []string{"viewer"}, PasswordHash: "h", CreatedAt: now, UpdatedAt: now,
	}); err != nil {
		t.Fatalf("create alice: %v", err)
	}
	if err := st.CreateRefreshToken(RefreshToken{
		ID: "rt1", UserID: "u1", TokenHash: "hash-rt1",
		ExpiresAt: now.Add(time.Hour), CreatedAt: now,
	}); err != nil {
		t.Fatalf("create refresh token: %v", err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}

	// Restart with a different seed hash: nothing is re-applied or overwritten.
	st2, err := Open(path, "admin", "hash-v2", true)
	if err != nil {
		t.Fatalf("restart open: %v", err)
	}
	defer st2.Close()

	applied, err := st2.appliedMigrations()
	if err != nil {
		t.Fatal(err)
	}
	if len(applied) != 8 || applied[0].version != 1 || applied[1].version != 2 || applied[2].version != 3 || applied[3].version != 4 || applied[4].version != 5 || applied[5].version != 6 || applied[6].version != 7 || applied[7].version != 8 {
		t.Fatalf("applied after restart = %+v, want {1..8} (no re-migration)", applied)
	}
	var ur int
	if err := st2.db.QueryRow(`SELECT COUNT(*) FROM user_roles WHERE user_id = 'user-admin'`).Scan(&ur); err != nil || ur != 2 {
		t.Fatalf("seed user_roles after restart = %d, err %v, want 2 (no duplicate seed)", ur, err)
	}

	alice, err := st2.UserByID("u1")
	if err != nil {
		t.Fatalf("alice after restart: %v", err)
	}
	if !reflect.DeepEqual(alice.Roles, []string{"viewer"}) {
		t.Fatalf("alice roles = %v, want [viewer]", alice.Roles)
	}
	rt, err := st2.RefreshTokenByHash("hash-rt1")
	if err != nil {
		t.Fatalf("refresh token after restart: %v", err)
	}
	if rt.UserID != "u1" || rt.RevokedAt != nil {
		t.Fatalf("refresh token after restart = %+v", rt)
	}

	// viewer: read-only permission, no menu grant.
	perms, err := st2.PermissionsForUser("u1")
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(perms, "users.read") || slices.Contains(perms, "users.write") {
		t.Fatalf("viewer permissions after restart = %v", perms)
	}
	feat, err := st2.FeaturesForUser("u1")
	if err != nil {
		t.Fatal(err)
	}
	if feat["menu_users"] {
		t.Fatalf("viewer menu feature after restart = true, want false")
	}
	admFeat, err := st2.FeaturesForUser("user-admin")
	if err != nil {
		t.Fatal(err)
	}
	if !admFeat["menu_users"] {
		t.Fatalf("admin menu feature lost after restart")
	}
	// The seed must not overwrite the admin password.
	admin, err := st2.UserByUsername("admin")
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
	st, err := Open(orig, "admin", "hash", false)
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
	r, err := Open(restored, "admin", "hash", true)
	if err != nil {
		t.Fatalf("open restored: %v", err)
	}
	defer r.Close()

	u, err := r.UserByUsername("admin")
	if err != nil {
		t.Fatalf("restored identity: %v", err)
	}
	if u.ID != "user-admin" || u.PasswordHash != "hash-v1" {
		t.Fatalf("restored user = %+v, want user-admin / hash-v1", u)
	}
	if want := []string{"admin", "editor"}; !reflect.DeepEqual(u.Roles, want) {
		t.Fatalf("restored roles = %v, want %v", u.Roles, want)
	}
	rt, err := r.RefreshTokenByHash("abc123")
	if err != nil {
		t.Fatalf("restored refresh: %v", err)
	}
	if rt.UserID != "user-admin" {
		t.Fatalf("restored refresh user = %q, want user-admin", rt.UserID)
	}
	perms, err := r.PermissionsForUser("user-admin")
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(perms, "users.read") || !slices.Contains(perms, "users.write") {
		t.Fatalf("restored permissions = %v", perms)
	}
	feat, err := r.FeaturesForUser("user-admin")
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
