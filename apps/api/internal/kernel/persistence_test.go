package kernel

import (
	"errors"
	"testing"
)

func migrationFixture(version int, name, checksum string, mutate func(*MigrationContribution)) MigrationContribution {
	m := MigrationContribution{
		ContributionIdentity: ContributionIdentity{ModuleID: "test.p", Key: name},
		Version:              version,
		Name:                 name,
		Checksum:             checksum,
		Apply:                func(Tx) error { return nil },
	}
	if mutate != nil {
		mutate(&m)
	}
	return m
}

func TestCollectPersistenceHappyPath(t *testing.T) {
	p1 := catalogModule(Module{ID: "test.p1", Version: "2.0.0"})
	p1.migrate = func() ([]MigrationContribution, error) {
		return []MigrationContribution{
			migrationFixture(1, "m1", "c1", nil),
			migrationFixture(2, "m2", "c2", nil),
		}, nil
	}
	p2 := catalogModule(Module{ID: "test.p2", Version: "2.0.0"})
	p2.migrate = func() ([]MigrationContribution, error) {
		return []MigrationContribution{migrationFixture(3, "m3", "c3", nil)}, nil
	}
	catalog, err := CollectPersistence([]Provider{p1, p2})
	if err != nil {
		t.Fatalf("CollectPersistence: %v", err)
	}
	if len(catalog) != 3 {
		t.Fatalf("catalog len = %d, want 3", len(catalog))
	}
	for i, want := range []int{1, 2, 3} {
		if catalog[i].Version != want {
			t.Fatalf("catalog[%d].version = %d, want %d", i, catalog[i].Version, want)
		}
	}
}

func TestCollectPersistenceVersionGap(t *testing.T) {
	p := catalogModule(Module{ID: "test.p", Version: "2.0.0"})
	p.migrate = func() ([]MigrationContribution, error) {
		return []MigrationContribution{
			migrationFixture(1, "m1", "c1", nil),
			migrationFixture(3, "m3", "c3", nil), // gap at version 2
		}, nil
	}
	_, err := CollectPersistence([]Provider{p})
	if err == nil {
		t.Fatal("version gap should fail closed")
	}
}

func TestCollectPersistenceDuplicateVersion(t *testing.T) {
	p := catalogModule(Module{ID: "test.p", Version: "2.0.0"})
	p.migrate = func() ([]MigrationContribution, error) {
		return []MigrationContribution{
			migrationFixture(1, "m1", "c1", nil),
			migrationFixture(1, "m1b", "c2", nil), // duplicate version
		}, nil
	}
	_, err := CollectPersistence([]Provider{p})
	if err == nil {
		t.Fatal("duplicate version should fail closed")
	}
}

func TestCollectPersistenceDuplicateName(t *testing.T) {
	p := catalogModule(Module{ID: "test.p", Version: "2.0.0"})
	p.migrate = func() ([]MigrationContribution, error) {
		return []MigrationContribution{
			migrationFixture(1, "m1", "c1", nil),
			migrationFixture(2, "m1", "c2", nil), // duplicate name
		}, nil
	}
	_, err := CollectPersistence([]Provider{p})
	if err == nil {
		t.Fatal("duplicate migration name should fail closed")
	}
}

func TestCollectPersistenceDuplicateChecksum(t *testing.T) {
	p := catalogModule(Module{ID: "test.p", Version: "2.0.0"})
	p.migrate = func() ([]MigrationContribution, error) {
		return []MigrationContribution{
			migrationFixture(1, "m1", "same", nil),
			migrationFixture(2, "m2", "same", nil), // duplicate checksum
		}, nil
	}
	_, err := CollectPersistence([]Provider{p})
	if err == nil {
		t.Fatal("duplicate checksum should fail closed")
	}
}

func TestCollectPersistenceTombstoneWithApply(t *testing.T) {
	p := catalogModule(Module{ID: "test.p", Version: "2.0.0"})
	p.migrate = func() ([]MigrationContribution, error) {
		return []MigrationContribution{
			migrationFixture(1, "m1", "c1", func(m *MigrationContribution) {
				m.Tombstone = true // tombstone must not carry Apply
			}),
		}, nil
	}
	_, err := CollectPersistence([]Provider{p})
	if err == nil {
		t.Fatal("tombstone carrying Apply should fail closed")
	}
}

func TestCollectPersistenceReconcileInvalid(t *testing.T) {
	p := catalogModule(Module{ID: "test.p", Version: "2.0.0"})
	p.migrate = func() ([]MigrationContribution, error) {
		return []MigrationContribution{
			migrationFixture(2, "m2", "c2", func(m *MigrationContribution) {
				m.ReconcileVersion = 3 // > migration version
				m.ReconcileChecksum = "r3"
			}),
		}, nil
	}
	_, err := CollectPersistence([]Provider{p})
	if err == nil {
		t.Fatal("reconcile version above migration version should fail closed")
	}

	p2 := catalogModule(Module{ID: "test.p2", Version: "2.0.0"})
	p2.migrate = func() ([]MigrationContribution, error) {
		return []MigrationContribution{
			migrationFixture(1, "m1", "c1", func(m *MigrationContribution) {
				m.ReconcileVersion = 1 // set without checksum
				m.ReconcileChecksum = ""
			}),
		}, nil
	}
	_, err = CollectPersistence([]Provider{p2})
	if err == nil {
		t.Fatal("reconcile version without checksum should fail closed")
	}
}

func TestCollectPersistenceProviderError(t *testing.T) {
	p := catalogModule(Module{ID: "test.p", Version: "2.0.0"})
	p.migrate = func() ([]MigrationContribution, error) {
		return nil, errors.New("boom")
	}
	_, err := CollectPersistence([]Provider{p})
	var kerr *Error
	if !errors.As(err, &kerr) || kerr.Code != CodeModuleInvalid {
		t.Fatalf("err = %v, want MODULE_INVALID wrapping provider error", err)
	}
}
