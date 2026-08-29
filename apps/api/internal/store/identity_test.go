package store

import (
	"fmt"
	"strings"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	compiledmodules "github.com/magicvr/schema-ui-core/apps/api/modules/compiled"
)

func testCatalog(n int) []kernel.MigrationContribution {
	out := make([]kernel.MigrationContribution, n)
	sum := strings.Repeat("a", 64)
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("m%d", i+1)
		out[i] = kernel.MigrationContribution{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: "t.test", Key: name},
			Version:              i + 1,
			Name:                 name,
			Checksum:             sum,
			Apply:                func(kernel.Tx) error { return nil },
		}
	}
	return out
}

func TestClassifyIdentity(t *testing.T) {
	sum := strings.Repeat("a", 64)
	ledger := []appliedMigration{{version: 1, name: "m1", checksum: sum}}
	cases := []struct {
		name      string
		tables    []string
		applied   []appliedMigration
		oursUsers bool
		want      dbIdentityKind
	}{
		{"empty", nil, nil, false, identityEmpty},
		{"ledger", []string{"schema_migrations", "users"}, ledger, true, identityOursLedger},
		{"r2", []string{"users", "refresh_tokens"}, nil, true, identityOursR2},
		{"complete no ledger", []string{"users", "refresh_tokens", "operation_log", "jobs", "service_credentials", "operation_log_session", "mail_outbox", "mail_config", "email_verification_challenges", "password_recovery_challenges", "password_policy", "user_password_history", "user_invites", "login_failures"}, nil, true, identityOursCompleteNoLedger},
		{"four tables without catalog head", []string{"users", "refresh_tokens", "operation_log", "jobs", "roles"}, nil, true, identityLostLedgerUnsafe},
		{"partial users only", []string{"users"}, nil, true, identityOursPartialNoLedger},
		{"foreign users", []string{"users"}, nil, false, identityForeign},
		{"foreign other", []string{"orders"}, nil, false, identityForeign},
		{"complete tables but not our users", []string{"users", "refresh_tokens", "operation_log", "jobs", "service_credentials", "operation_log_session"}, nil, false, identityForeign},
		{"ledger wins over not-ours users", []string{"schema_migrations", "users"}, ledger, false, identityOursLedger},
	}
	for _, tc := range cases {
		got := classifyIdentity(tc.tables, tc.applied, tc.oursUsers)
		if got.Kind != tc.want {
			t.Errorf("%s: kind = %q, want %q", tc.name, got.Kind, tc.want)
		}
	}
}

func TestPlanStartup(t *testing.T) {
	cat := testCatalog(3)
	sum := strings.Repeat("a", 64)
	cases := []struct {
		name   string
		id     dbIdentity
		want   startupAction
		refuse bool
	}{
		{"fresh", dbIdentity{Kind: identityEmpty}, actionFresh, false},
		{"noop", dbIdentity{Kind: identityOursLedger, Applied: []appliedMigration{
			{1, "m1", sum}, {2, "m2", sum}, {3, "m3", sum},
		}}, actionNoop, false},
		{"pending", dbIdentity{Kind: identityOursLedger, Applied: []appliedMigration{
			{1, "m1", sum},
		}}, actionApplyPending, false},
		{"r2", dbIdentity{Kind: identityOursR2}, actionAdoptR2, false},
		{"restore", dbIdentity{Kind: identityOursCompleteNoLedger}, actionRestoreLedger, false},
		{"partial", dbIdentity{Kind: identityOursPartialNoLedger}, actionAdoptThenPending, false},
		{"unsafe", dbIdentity{Kind: identityLostLedgerUnsafe, Tables: []string{"users", "roles"}}, actionRefuse, true},
		{"foreign", dbIdentity{Kind: identityForeign, Tables: []string{"orders"}}, actionRefuse, true},
	}
	for _, tc := range cases {
		plan, err := planStartup(tc.id, cat)
		if err != nil {
			t.Fatalf("%s: plan err %v", tc.name, err)
		}
		if plan.Action != tc.want {
			t.Errorf("%s: action = %q, want %q", tc.name, plan.Action, tc.want)
		}
		if tc.refuse && !strings.Contains(plan.Reason, "identity=") {
			t.Errorf("%s: refuse reason %q", tc.name, plan.Reason)
		}
	}
}

// lockedHeadExtraTables requires catalog-head object names in the restore
// fingerprint. Bumping completeFingerprintCatalogHead without adding a map
// entry fails closed (A-003 F-001). A data-only head migration (e.g. 0049
// seed_admin_must_change_password backfill) registers an empty list to record
// that omission of new objects was reviewed, not overlooked.
var lockedHeadExtraTables = map[int][]string{
	48: {"service_credentials", "operation_log_session"},
	49: {},
	50: {}, // wallet_ledger_order_repair: data-only repair (no new objects)
	51: {"mail_outbox"}, // VP-017 R6 mock-channel outbound record table
	52: {"mail_config"}, // VP-017 R7 runtime channel state
	53: {}, // operation_log_mail_events: CHECK-enum expansion (no new objects)
	54: {}, // account_email_identity: ALTER + lower(email) unique index only (no new objects)
	55: {"email_verification_challenges"}, // workspace-018 R3: per-user active verification challenge
	56: {"password_recovery_challenges"}, // workspace-019 R2: per-user active recovery challenge
	57: {"password_policy"},           // workspace-019 R3: singleton policy row
	58: {"user_password_history"},     // workspace-019 R3: history-depth store
	59: {"user_invites"},              // workspace-019 R3: admin invitations
	60: {}, // W26 GOAL-038: mail_outbox additive ALTER (channel/delivery_status columns; no new objects)
	61: {"login_failures"}, // GOAL-014 D-002: per-(account|source) login-lockout state
	62: {}, // workspace-020 R3: site_settings additive ALTER (default_currency column; no new objects)
	63: {}, // R4 演练（GOAL-005 S2）: CREATE INDEX IF NOT EXISTS only (no new objects)
}

func TestCompleteFingerprintTracksCatalogHead(t *testing.T) {
	catalog, err := compiledmodules.PersistenceCatalog()
	if err != nil {
		t.Fatal(err)
	}
	max := 0
	for _, m := range catalog {
		if m.Version > max {
			max = m.Version
		}
	}
	if max != completeFingerprintCatalogHead {
		t.Fatalf("compiled catalog head is v%d; update completeFingerprintCatalogHead, completeLostLedgerTables, and lockedHeadExtraTables[%d]", max, max)
	}
	extra, ok := lockedHeadExtraTables[max]
	if !ok {
		t.Fatalf("lockedHeadExtraTables[%d] missing; record the new CREATE TABLE names (empty list only for reviewed data-only head migrations)", max)
	}
	have := tableNameSet(completeLostLedgerTables)
	for _, name := range extra {
		if !have[name] {
			t.Fatalf("completeLostLedgerTables missing catalog-head table %q", name)
		}
	}
}
