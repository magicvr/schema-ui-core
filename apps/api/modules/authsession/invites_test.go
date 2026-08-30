// Invitation domain tests (workspace-019 R3 · GOAL-004 D-001 §3): roles fixed
// at issuance (user adjudication), one-time redemption, instant revoke,
// token-rotating resend with 60 s cooldown, role-gone fail-closed, and
// username-conflict rejection. Policy tests live beside them.
package authsession

import (
	"errors"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func openInviteFixture(t *testing.T) (*Repository, time.Time) {
	t.Helper()
	repo, _ := openRecoveryFixture(t)
	base := time.Now().UTC().Truncate(time.Second)
	return repo, base
}

func TestCreateAndAcceptInviteAssignsIssuedRoles(t *testing.T) {
	repo, base := openInviteFixture(t)
	raw, inv, err := repo.CreateInvite("user-admin", []string{"admin"}, "", defaultInviteTTL, base)
	if err != nil {
		t.Fatalf("create invite: %v", err)
	}
	if raw == "" || len(inv.Roles) != 1 || inv.Roles[0] != "admin" {
		t.Fatalf("created = (%q, %+v)", raw, inv)
	}
	hash, herr := bcrypt.GenerateFromPassword([]byte("invited-pass-1"), 4)
	if herr != nil {
		t.Fatalf("hash: %v", herr)
	}
	u, err := repo.AcceptInvite(raw, "newbie", "New Bie", string(hash), base.Add(time.Minute))
	if err != nil {
		t.Fatalf("accept: %v", err)
	}
	// The issued role travels to the new account (裁决：角色以邀请为准).
	granted, err := repo.PermissionsForRoles(u.Roles)
	_ = granted // permissions projection is a management concern; role keys matter:
	if len(u.Roles) != 1 || u.Roles[0] != "admin" {
		t.Fatalf("accepted roles = %v, want [admin]", u.Roles)
	}
	if u.Username != "newbie" || u.MustChangePassword {
		t.Fatalf("accepted user = %+v (must_change must stay false — invitee chose the password)", u)
	}
	// One-time: replay of the same token is dead.
	if _, err := repo.AcceptInvite(raw, "second", "S", string(hash), base.Add(2*time.Minute)); !errors.Is(err, ErrInviteInvalid) {
		t.Fatalf("replay accept err = %v, want ErrInviteInvalid", err)
	}
}

func TestAcceptInviteRejectsUnknownTokenUsernameConflictAndRoleGone(t *testing.T) {
	repo, base := openInviteFixture(t)
	if _, err := repo.AcceptInvite("no-such-token", "x", "X", "h", base); !errors.Is(err, ErrInviteInvalid) {
		t.Fatalf("unknown token err = %v, want ErrInviteInvalid", err)
	}
	raw, _, err := repo.CreateInvite("user-admin", []string{"viewer"}, "", defaultInviteTTL, base)
	if err != nil {
		t.Fatalf("create invite: %v", err)
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte("invited-pass-1"), 4)
	// Username conflict → fail-closed, invite stays live.
	if _, err := repo.AcceptInvite(raw, "admin", "Dup", string(hash), base.Add(time.Minute)); !errors.Is(err, ErrUsernameTaken) {
		t.Fatalf("conflict err = %v, want ErrUsernameTaken", err)
	}
	if _, err := repo.AcceptInvite(raw, "ok-name", "Ok", string(hash), base.Add(2*time.Minute)); err != nil {
		t.Fatalf("accept after failed conflict attempt: %v", err)
	}
	// Role deletion between issue and acceptance → INVITE_ROLE_GONE.
	repo2, base2 := openInviteFixture(t)
	raw2, _, err := repo2.CreateInvite("user-admin", []string{"ghost-role"}, "", defaultInviteTTL, base2)
	if err == nil {
		t.Fatal("unknown role at creation must fail")
	}
	if !errors.Is(err, ErrInviteRoleGone) {
		t.Fatalf("create with unknown role err = %v, want ErrInviteRoleGone", err)
	}
	_ = raw2
}

func TestInviteResendCooldownRotationAndRevoke(t *testing.T) {
	repo, base := openInviteFixture(t)
	raw, inv, err := repo.CreateInvite("user-admin", []string{"editor"}, "ed@example.com", defaultInviteTTL, base)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	// Cooldown: immediate resend rejected.
	if _, _, err := repo.ResendInvite(inv.ID, defaultInviteTTL, base.Add(10*time.Second)); !errors.Is(err, ErrInviteCooldown) {
		t.Fatalf("resend cooldown err = %v, want ErrInviteCooldown", err)
	}
	// After the window: resend rotates the token; old link dies.
	raw2, inv2, err := repo.ResendInvite(inv.ID, defaultInviteTTL, base.Add(61*time.Second))
	if err != nil {
		t.Fatalf("resend: %v", err)
	}
	if raw2 == raw {
		t.Fatal("resend must rotate the token")
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte("invited-pass-1"), 4)
	if _, err := repo.AcceptInvite(raw, "old-link", "O", string(hash), base.Add(62*time.Second)); !errors.Is(err, ErrInviteInvalid) {
		t.Fatalf("old link after resend err = %v, want ErrInviteInvalid", err)
	}
	// Revoke kills the rotated link instantly.
	if err := repo.RevokeInvite(inv2.ID, base.Add(70*time.Second)); err != nil {
		t.Fatalf("revoke: %v", err)
	}
	if _, err := repo.AcceptInvite(raw2, "revoked", "R", string(hash), base.Add(80*time.Second)); !errors.Is(err, ErrInviteInvalid) {
		t.Fatalf("revoked accept err = %v, want ErrInviteInvalid", err)
	}
	// Unknown id surfaces distinctly for the admin surface.
	if err := repo.RevokeInvite("inv-nothing", base.Add(90*time.Second)); !errors.Is(err, ErrInviteNotFound) {
		t.Fatalf("revoke unknown err = %v, want ErrInviteNotFound", err)
	}
}

// A-001 F-002: expiry and acceptance-side role deletion surface as the
// uniform invalid / role-gone answers respectively.
func TestAcceptInviteExpiredAndRoleGoneAfterIssue(t *testing.T) {
	repo, base := openInviteFixture(t)
	raw, _, err := repo.CreateInvite("user-admin", []string{"editor"}, "", 1*time.Hour, base)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte("invited-pass-1"), 4)
	// Expired window → uniform invalid.
	if _, err := repo.AcceptInvite(raw, "late", "L", string(hash), base.Add(2*time.Hour)); !errors.Is(err, ErrInviteInvalid) {
		t.Fatalf("expired accept err = %v, want ErrInviteInvalid", err)
	}
	// Role deleted AFTER issuance → role-gone (fail-closed, admin reissues).
	// A NON-system role is required: system roles refuse deletion.
	repo2, base2 := openInviteFixture(t)
	if _, err := repo2.CreateRole("temp-guest", "Temp Guest", base2); err != nil {
		t.Fatalf("create temp role: %v", err)
	}
	raw2, _, err := repo2.CreateInvite("user-admin", []string{"temp-guest"}, "", defaultInviteTTL, base2)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := repo2.DeleteRole("role-temp-guest"); err != nil {
		t.Fatalf("delete role: %v", err)
	}
	if _, err := repo2.AcceptInvite(raw2, "stale-roles", "S", string(hash), base2.Add(time.Minute)); !errors.Is(err, ErrInviteRoleGone) {
		t.Fatalf("role-gone accept err = %v, want ErrInviteRoleGone", err)
	}
}


func TestListInvitesStatusFilter(t *testing.T) {
	repo, base := openInviteFixture(t)
	hash, _ := bcrypt.GenerateFromPassword([]byte("invited-pass-1"), 4)

	// A: pending (TTL 1h, viewed before expiry).
	if _, _, err := repo.CreateInvite("user-admin", []string{"viewer"}, "a@example.com", defaultInviteTTL, base); err != nil {
		t.Fatalf("create A: %v", err)
	}
	// B: consumed.
	rawB, _, err := repo.CreateInvite("user-admin", []string{"editor"}, "b@example.com", defaultInviteTTL, base)
	if err != nil {
		t.Fatalf("create B: %v", err)
	}
	if _, err := repo.AcceptInvite(rawB, "consume-b", "B", string(hash), base.Add(time.Minute)); err != nil {
		t.Fatalf("consume B: %v", err)
	}
	// C: revoked.
	_, invC, err := repo.CreateInvite("user-admin", []string{"viewer"}, "c@example.com", defaultInviteTTL, base)
	if err != nil {
		t.Fatalf("create C: %v", err)
	}
	if err := repo.RevokeInvite(invC.ID, base.Add(2*time.Minute)); err != nil {
		t.Fatalf("revoke C: %v", err)
	}
	// D: short TTL so it expires before the view time.
	if _, _, err := repo.CreateInvite("user-admin", []string{"viewer"}, "d@example.com", time.Hour, base); err != nil {
		t.Fatalf("create D: %v", err)
	}

	view := base.Add(2 * time.Hour)
	cases := []struct {
		filter InviteStatusFilter
		want   int
	}{
		{InviteStatusAll, 4},
		{InviteStatusPending, 1}, // A only; D has expired by view time
		{InviteStatusConsumed, 1},
		{InviteStatusRevoked, 1},
		{InviteStatusExpired, 1},
	}
	for _, tc := range cases {
		got, total, err := repo.ListInvites(1, 50, tc.filter, "", "", "", view)
		if err != nil {
			t.Fatalf("list %q: %v", tc.filter, err)
		}
		if len(got) != tc.want || total != tc.want {
			t.Fatalf("filter %q = (%d rows, total %d), want %d", tc.filter, len(got), total, tc.want)
		}
	}
	// Unknown raw values parse to "all" (stale-client tolerance).
	if ParseInviteStatus("bogus") != InviteStatusAll || ParseInviteStatus(" PENDING ") != InviteStatusPending {
		t.Fatal("ParseInviteStatus mapping broken")
	}
}

// W27 (GOAL-039 D-001 §1): the admin listing supports keyword search over
// email/id/invited_by and whitelist sorting (createdAt default, expiresAt;
// asc/desc; stable id tiebreak). Unknown values fall back to defaults.
func TestListInvitesSearchAndSort(t *testing.T) {
	repo, base := openInviteFixture(t)

	// Insert out of chronological order so created_at ordering is observable:
	// first@example.com (oldest), then zeta, then alpha (newest).
	if _, _, err := repo.CreateInvite("user-admin", []string{"viewer"}, "first@example.com", defaultInviteTTL, base); err != nil {
		t.Fatalf("create first: %v", err)
	}
	if _, _, err := repo.CreateInvite("user-admin", []string{"viewer"}, "zeta@example.com", defaultInviteTTL, base.Add(time.Minute)); err != nil {
		t.Fatalf("create zeta: %v", err)
	}
	if _, invAlpha, err := repo.CreateInvite("user-admin", []string{"editor"}, "alpha@example.com", 48*time.Hour, base.Add(2*time.Minute)); err != nil {
		t.Fatalf("create alpha: %v", err)
	} else if len(invAlpha.Roles) == 0 {
		t.Fatal("unexpected empty roles")
	}

	view := base.Add(time.Hour)
	listEmails := func(q, sort, order string) []string {
		t.Helper()
		got, total, err := repo.ListInvites(1, 50, InviteStatusAll, q, sort, order, view)
		if err != nil {
			t.Fatalf("list(q=%q sort=%q order=%q): %v", q, sort, order, err)
		}
		emails := make([]string, 0, len(got))
		for i := range got {
			if got[i].Email != nil {
				emails = append(emails, *got[i].Email)
			}
		}
		if total != len(emails) {
			t.Fatalf("total = %d, rows = %d", total, len(emails))
		}
		return emails
	}

	// Default order: newest first.
	if got := listEmails("", "", ""); len(got) != 3 || got[0] != "alpha@example.com" || got[2] != "first@example.com" {
		t.Fatalf("default order = %v, want newest-first", got)
	}
	// createdAt ascending flips it.
	if got := listEmails("", "createdAt", "asc"); len(got) != 3 || got[0] != "first@example.com" {
		t.Fatalf("createdAt asc = %v, want oldest-first", got)
	}
	// expiresAt ordering: alpha has a 48h TTL so it is LAST under expiresAt desc.
	if got := listEmails("", "expiresat", "desc"); len(got) != 3 || got[2] != "alpha@example.com" {
		t.Fatalf("expiresAt desc = %v, want longest-TTL last", got)
	}
	// Keyword search hits email (case-insensitive substring).
	if got := listEmails("ZETA", "", ""); len(got) != 1 || got[0] != "zeta@example.com" {
		t.Fatalf("q=ZETA = %v, want only zeta", got)
	}
	// Keyword search hits invited_by.
	if got := listEmails("user-admin", "", ""); len(got) != 3 {
		t.Fatalf("q=user-admin = %v, want all three", got)
	}
	// No match → zero rows, zero total.
	if got, total, err := repo.ListInvites(1, 50, InviteStatusAll, "no-such-needle", "", "", view); err != nil || len(got) != 0 || total != 0 {
		t.Fatalf("q=no-such-needle = (%d rows, total %d, err %v), want empty", len(got), total, err)
	}
}
// --- password policy ---

func TestValidateNewPasswordBaselineCategoriesHistory(t *testing.T) {
	repo, base := openInviteFixture(t)

	// Baseline (frozen defaults): 8–72 bytes non-blank.
	if err := repo.ValidateNewPassword("", "short"); !errors.Is(err, ErrPasswordPolicyViolation) {
		t.Fatalf("short err = %v, want violation", err)
	}
	if err := repo.ValidateNewPassword("", "long-enough-pass"); err != nil {
		t.Fatalf("default-policy pass err = %v, want nil", err)
	}
	// Complexity knob: two categories demanded.
	if err := repo.UpdatePasswordPolicy(PasswordPolicy{MinLength: 8, MinCategories: 3}); err != nil {
		t.Fatalf("update policy: %v", err)
	}
	if err := repo.ValidateNewPassword("", "alllowercase"); !errors.Is(err, ErrPasswordPolicyViolation) {
		t.Fatal("single-category password must violate min_categories=3")
	}
	if err := repo.ValidateNewPassword("", "GoodPass12"); err != nil {
		t.Fatalf("three-category pass err = %v, want nil", err)
	}

	// History: set a real hash via UpdateUser, then try to reuse it.
	if err := repo.UpdatePasswordPolicy(PasswordPolicy{MinLength: 8, HistoryDepth: 3}); err != nil {
		t.Fatalf("enable history: %v", err)
	}
	oldHash, _ := bcrypt.GenerateFromPassword([]byte("old-password-1"), 4)
	newHash, _ := bcrypt.GenerateFromPassword([]byte("fresh-password-9"), 4)
	u, err := repo.CreateUserManagement(User{
		ID: "u-hist", Username: "hist", Name: "Hist", Roles: []string{"viewer"},
		PasswordHash: string(oldHash), CreatedAt: base, UpdatedAt: base,
	})
	if err != nil {
		t.Fatalf("seed history user: %v", err)
	}
	if err := repo.ValidateNewPassword(u.ID, "old-password-1"); err != nil {
		t.Fatalf("pre-rotation current password must pass (not yet in history): %v", err)
	}
	if _, err := repo.UpdateUser(u.ID, UserPatch{PasswordHash: stringPtr(string(newHash))}, u.ID, base.Add(time.Minute)); err != nil {
		t.Fatalf("rotate: %v", err)
	}
	if err := repo.ValidateNewPassword(u.ID, "old-password-1"); !errors.Is(err, ErrPasswordPolicyViolation) {
		t.Fatalf("history reuse err = %v, want violation", err)
	}
	if err := repo.ValidateNewPassword(u.ID, "never-used-pass"); err != nil {
		t.Fatalf("fresh pass err = %v, want nil", err)
	}
}

// A-001 F-001: the CONFIGURED minimum is authoritative — tightening to 12
// must reject 8–11 byte passwords at every enforcement call.
func TestValidateNewPasswordConfiguredMinLengthBites(t *testing.T) {
	repo, _ := openInviteFixture(t)
	if err := repo.UpdatePasswordPolicy(PasswordPolicy{MinLength: 12}); err != nil {
		t.Fatalf("tighten minLength: %v", err)
	}
	if err := repo.ValidateNewPassword("", "eight-byte"); !errors.Is(err, ErrPasswordPolicyViolation) {
		t.Fatal("8-byte password must violate a configured minLength of 12")
	}
	if err := repo.ValidateNewPassword("", "twelve-bytes!"); err != nil {
		t.Fatalf("12+ byte pass err = %v, want nil", err)
	}
}
