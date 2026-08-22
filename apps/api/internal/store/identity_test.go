package store

import (
	"fmt"
	"strings"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	compiledmodules "github.com/magicvr/schema-ui-core/apps/api/internal/modules/compiled"
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
		{"complete no ledger", []string{"users", "refresh_tokens", "operation_log", "jobs", "service_credentials", "operation_log_session"}, nil, true, identityOursCompleteNoLedger},
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
	if max > completeFingerprintCatalogHead {
		t.Fatalf("compiled catalog head is v%d; update completeLostLedgerTables and completeFingerprintCatalogHead (was %d) so restore-ledger cannot stamp past missing objects", max, completeFingerprintCatalogHead)
	}
}
