package store

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// fakeRecoveryPort records CreateRecoveryPoint calls and can be made to fail.
type fakeRecoveryPort struct {
	mu    sync.Mutex
	calls []kernel.RecoveryPointRequest
	fail  error
}

func (f *fakeRecoveryPort) CreateRecoveryPoint(ctx context.Context, req kernel.RecoveryPointRequest) (kernel.RecoveryPoint, error) {
	f.mu.Lock()
	f.calls = append(f.calls, req)
	fail := f.fail
	f.mu.Unlock()
	if fail != nil {
		return kernel.RecoveryPoint{}, fail
	}
	return kernel.RecoveryPoint{
		ID:             "rp-test",
		Dialect:        req.Dialect,
		ArtifactRef:    req.SourceID + ".artifact",
		CatalogVersion: req.CatalogVersion,
		BatchVersion:   req.CatalogVersion,
		ContractShape:  kernel.ContractShapeConverted,
		ChecksumSet:    "test-set",
		TimeContract:   kernel.TimeContractVP040Timestamptz,
		CreatedAt:      time.Now().UTC().Truncate(time.Microsecond),
		VerifiedAt:     time.Now().UTC().Truncate(time.Microsecond),
		Verification: kernel.VerificationSummary{
			SchemaVerified: true, TypeContractVerified: true,
			SampleVerified: true, ChecksumVerified: true,
		},
	}, nil
}

func (f *fakeRecoveryPort) callCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.calls)
}

// fakeRollbackCreator records class-A artifact requests.
type fakeRollbackCreator struct {
	mu    sync.Mutex
	calls [][2]string
}

func (f *fakeRollbackCreator) CreateRollbackArtifact(ctx context.Context, sourceID, artifactRef string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls = append(f.calls, [2]string{sourceID, artifactRef})
	return nil
}

func (f *fakeRollbackCreator) callCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.calls)
}

// TestC3RecoveryAnchorsOnUpgrade covers the C3 §4.2/§4.3 call points on the
// sqlite runner: the class-A rollback artifact exists before the batch, the
// class-B recovery point is created (and marked) only after the batch commits,
// and a missing marker becomes the explicit "verify-recovery-point" action with a
// bounded re-create instead of a silent noop.
func TestC3RecoveryAnchorsOnUpgrade(t *testing.T) {
	catalog := MigrationCatalog()
	dir := t.TempDir()
	path := filepath.Join(dir, "anchors.db")

	// Pre-conversion database (v1..v72).
	legacy, err := OpenWithCatalog(path, catalog[:72])
	if err != nil {
		t.Fatalf("open v72: %v", err)
	}
	if err := legacy.Close(); err != nil {
		t.Fatalf("close v72: %v", err)
	}

	port := &fakeRecoveryPort{}
	rollback := &fakeRollbackCreator{}
	openWith := func() (*Store, error) {
		st, err := Open(context.Background(), OpenOptions{
			Dialect:           kernel.DialectSQLite,
			Path:              path,
			RecoveryPoints:    port,
			RollbackArtifacts: rollback,
		}, catalog)
		if err != nil {
			return nil, err
		}
		return st.(*Store), nil
	}

	st, err := openWith()
	if err != nil {
		t.Fatalf("upgrade with anchors: %v", err)
	}
	defer st.Close()

	// Call point 1 (class A): the sqlite runner takes its own VACUUM INTO
	// rollback artifact, named distinctly from the class-C per-migration files.
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read dir: %v", err)
	}
	var classA, classC int
	for _, entry := range entries {
		switch {
		case strings.Contains(entry.Name(), ".batch-rollback-"):
			classA++
		case strings.Contains(entry.Name(), ".pre-v"):
			classC++
		}
	}
	if classA != 1 {
		t.Fatalf("class A (batch-rollback) artifacts = %d, want exactly 1 (%v)", classA, entries)
	}
	if classC == 0 {
		t.Fatalf("expected per-migration class-C snapshots to still be taken")
	}

	// Call point 3 (class B): one recovery point after the committed batch, and
	// the marker recorded.
	if got := port.callCount(); got != 1 {
		t.Fatalf("CreateRecoveryPoint calls = %d, want 1", got)
	}
	markerPath := sqliteMarkerPath(path)
	if _, err := os.Stat(markerPath); err != nil {
		t.Fatalf("recovery marker not written: %v", err)
	}
	marker, has, err := readRecoveryMarker(markerPath)
	if err != nil || !has {
		t.Fatalf("read marker: has=%v err=%v", has, err)
	}
	if marker.CatalogVersion != len(catalog) {
		t.Fatalf("marker catalog version = %d, want %d", marker.CatalogVersion, len(catalog))
	}
	if note := st.RecoveryNote(); note != "" {
		t.Fatalf("unexpected recovery note: %q", note)
	}

	// A second open at head with the marker present stays a noop: no extra B.
	st2, err := openWith()
	if err != nil {
		t.Fatalf("second open: %v", err)
	}
	defer st2.Close()
	if got := port.callCount(); got != 1 {
		t.Fatalf("CreateRecoveryPoint calls after noop open = %d, want 1", got)
	}
	if state := st2.RecoveryPointState(); !state.HasPoint || !state.Converted {
		t.Fatalf("state = %+v, want converted with a recorded point", state)
	}

	// Removing the marker makes "no B" detectable: the next open takes the
	// explicit recovery action and re-creates B exactly once (bounded).
	if err := os.Remove(markerPath); err != nil {
		t.Fatalf("remove marker: %v", err)
	}
	st3, err := openWith()
	if err != nil {
		t.Fatalf("third open: %v", err)
	}
	defer st3.Close()
	if got := port.callCount(); got != 2 {
		t.Fatalf("CreateRecoveryPoint calls after marker loss = %d, want 2 (one bounded re-create)", got)
	}
	if _, err := os.Stat(markerPath); err != nil {
		t.Fatalf("marker not re-written after the bounded re-create: %v", err)
	}
}

// TestC3RecoveryPointFailureNeverBlocksStartup pins C3 §4.3 item 3: a failed
// class-B creation is recorded, does not block startup, and leaves the "no B"
// state visible rather than pretending the gate is satisfied.
func TestC3RecoveryPointFailureNeverBlocksStartup(t *testing.T) {
	catalog := MigrationCatalog()
	dir := t.TempDir()
	path := filepath.Join(dir, "failing.db")
	legacy, err := OpenWithCatalog(path, catalog[:72])
	if err != nil {
		t.Fatalf("open v72: %v", err)
	}
	if err := legacy.Close(); err != nil {
		t.Fatalf("close v72: %v", err)
	}

	port := &fakeRecoveryPort{fail: errFakeRecovery}
	st, err := Open(context.Background(), OpenOptions{Dialect: kernel.DialectSQLite, Path: path, RecoveryPoints: port}, catalog)
	if err != nil {
		t.Fatalf("startup must not fail on a recovery-point error: %v", err)
	}
	defer st.Close()
	store := st.(*Store)
	if note := store.RecoveryNote(); !strings.Contains(note, "recovery point creation failed") {
		t.Fatalf("recovery note = %q, want a recorded failure", note)
	}
	if _, err := os.Stat(sqliteMarkerPath(path)); err == nil {
		t.Fatal("a failed recovery point must not leave a marker behind")
	}
	state := store.RecoveryPointState()
	if state.HasPoint {
		t.Fatalf("failed recovery point must not be reported as satisfied: %+v", state)
	}
	if !state.Converted {
		t.Fatalf("the converted shape must still be measured: %+v", state)
	}
}

var errFakeRecovery = errors.New("fake recovery-point failure")
