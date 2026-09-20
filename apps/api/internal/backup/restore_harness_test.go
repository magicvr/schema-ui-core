package backup

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	compiledmodules "github.com/magicvr/schema-ui-core/apps/api/modules/compiled"
)

// convertedStore materialises a store at the current catalog head (v1..v87) with
// representative legacy values seeded before the conversion, so the recovery
// artifact exercises real converted data.
func convertedStore(t *testing.T, dir string) string {
	t.Helper()
	path := filepath.Join(dir, "source.db")
	catalog := migrationCatalog(t)
	st, err := store.OpenWithCatalog(path, catalog[:v72Head])
	if err != nil {
		t.Fatalf("open v72: %v", err)
	}
	if err := st.Close(); err != nil {
		t.Fatalf("close v72: %v", err)
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("raw open: %v", err)
	}
	seed := []struct {
		query string
		args  []any
	}{
		{`INSERT INTO users (id, username, name, roles, password_hash, created_at, updated_at,
			token_version, failed_login_count, locked_until, enabled, notifications_enabled, avatar_url,
			must_change_password, email, email_status, last_login_failure_at)
			VALUES ('u-1','u1','User One','[]','h', 1758320000, 1758320001, 0, 0, 0, 1, 1, '', 0, NULL, NULL, 0)`, nil},
		{`INSERT INTO jobs (id, kind, status, payload, progress, cancel_requested, attempt, max_attempts,
			lease_owner, lease_version, lease_expires_at, result, error_code, error_message, actor_id,
			correlation_id, created_at, updated_at, finished_at, expires_at)
			VALUES ('j-1','probe','queued','{}',0,0,0,3,NULL,0,NULL,NULL,NULL,NULL,'a','c',1758320000123,1758320000123,NULL,NULL)`, nil},
		{`INSERT INTO mail_config (id, updated_at) VALUES (1, 0)`, nil},
		{`INSERT INTO telegram_config (id, updated_at) VALUES (1, 0)`, nil},
	}
	for _, statement := range seed {
		if _, err := db.Exec(statement.query, statement.args...); err != nil {
			_ = db.Close()
			t.Fatalf("seed: %v", err)
		}
	}
	if err := db.Close(); err != nil {
		t.Fatalf("close seed db: %v", err)
	}
	st, err = store.OpenWithCatalog(path, catalog)
	if err != nil {
		t.Fatalf("converted open: %v", err)
	}
	if err := st.Close(); err != nil {
		t.Fatalf("close converted: %v", err)
	}
	return path
}

// v72Head is the last pre-conversion catalog version; v73Head is the first
// converted one (used to materialise a mid-batch, class-C artifact).
const (
	v72Head = 72
	v73Head = 73
)

func migrationCatalog(t *testing.T) []kernel.MigrationContribution {
	t.Helper()
	catalog, err := compiledmodules.PersistenceCatalog()
	if err != nil {
		t.Fatalf("compiled catalog: %v", err)
	}
	return catalog
}

// TestSQLiteRestoreToNewDB is the C3 §5 harness for SQLite: create a recovery
// point, restore the artifact into a NEW file, and assert the restored target
// carries the converted contract (integrity, FK, 90 columns, samples) plus the
// ledger fingerprint of the source.
func TestSQLiteRestoreToNewDB(t *testing.T) {
	dir := t.TempDir()
	source := convertedStore(t, dir)
	service := NewService(t.TempDir())

	catalog := migrationCatalog(t)
	point, err := service.CreateRecoveryPoint(context.Background(), kernel.RecoveryPointRequest{
		Dialect:        kernel.DialectSQLite,
		SourceID:       source,
		CatalogVersion: len(catalog),
		TimeContract:   kernel.TimeContractVP040Timestamptz,
	})
	if err != nil {
		t.Fatalf("CreateRecoveryPoint: %v", err)
	}
	if point.ContractShape != kernel.ContractShapeConverted {
		t.Fatalf("contract shape = %q, want converted", point.ContractShape)
	}
	if !point.Verification.Passed() {
		t.Fatalf("verification summary did not pass: %+v", point.Verification)
	}
	if point.BatchVersion != len(catalog) {
		t.Fatalf("batch version = %d, want %d", point.BatchVersion, len(catalog))
	}
	if point.ChecksumSet == "" {
		t.Fatal("checksum set must be recorded")
	}

	// Restore to a NEW path and re-measure independently of the service.
	target := filepath.Join(t.TempDir(), "restored.db")
	provider := SQLiteProvider{}
	if err := provider.RestoreInto(context.Background(), point.ArtifactRef, target); err != nil {
		t.Fatalf("restore to new db: %v", err)
	}
	measure, err := measureSQLiteTarget(context.Background(), target)
	if err != nil {
		t.Fatalf("measure restored target: %v", err)
	}
	if !measure.IntegrityOK {
		t.Fatal("restored target integrity_check != ok")
	}
	if measure.ForeignKeyViolations != 0 {
		t.Fatalf("restored target has %d FK violations", measure.ForeignKeyViolations)
	}
	if len(measure.Missing) != 0 || measure.MeasuredColumns != len(temporalColumns) {
		t.Fatalf("restored target measured %d/%d columns, missing %v",
			measure.MeasuredColumns, len(temporalColumns), measure.Missing)
	}
	if len(measure.WrongShape) != 0 {
		t.Fatalf("restored target has non-converted columns: %v", measure.WrongShape)
	}
	if measure.LedgerSet != point.ChecksumSet {
		t.Fatalf("restored ledger %s != recovery point %s", measure.LedgerSet, point.ChecksumSet)
	}

	// Sample round-trip: seconds, milliseconds and sentinel 0 -> NULL.
	db, err := sql.Open("sqlite", target)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var created string
	if err := db.QueryRow(`SELECT created_at FROM users WHERE id = 'u-1'`).Scan(&created); err != nil {
		t.Fatalf("read users.created_at: %v", err)
	}
	if created != "2025-09-19T22:13:20.000000Z" {
		t.Fatalf("users.created_at = %q, want the converted second sample", created)
	}
	var jobCreated string
	if err := db.QueryRow(`SELECT created_at FROM jobs WHERE id = 'j-1'`).Scan(&jobCreated); err != nil {
		t.Fatalf("read jobs.created_at: %v", err)
	}
	if jobCreated != "2025-09-19T22:13:20.123000Z" {
		t.Fatalf("jobs.created_at = %q, want the converted millisecond sample", jobCreated)
	}
	if err := verifySQLiteSamples(context.Background(), target); err != nil {
		t.Fatalf("sample verification on restored target: %v", err)
	}
}

// TestLegacyArtifactMustFail is the C3 §5.1 / §6 reverse assertion: an artifact
// taken BEFORE the conversion (class A) must be rejected by the same
// verification, and the rejection must be classified as a contract problem —
// never as a missing/unreadable artifact or a tool failure (otherwise a missing
// file would also "pass").
func TestLegacyArtifactMustFail(t *testing.T) {
	dir := t.TempDir()
	legacyPath := filepath.Join(dir, "legacy.db")
	catalog := migrationCatalog(t)
	st, err := store.OpenWithCatalog(legacyPath, catalog[:v72Head])
	if err != nil {
		t.Fatalf("open v72: %v", err)
	}
	if err := st.Close(); err != nil {
		t.Fatalf("close v72: %v", err)
	}

	provider := SQLiteProvider{}
	artifact := filepath.Join(dir, "legacy.artifact")
	if err := provider.Create(context.Background(), legacyPath, artifact); err != nil {
		t.Fatalf("create legacy artifact: %v", err)
	}
	target := filepath.Join(dir, "legacy-restored.db")
	if err := provider.RestoreInto(context.Background(), artifact, target); err != nil {
		t.Fatalf("restore legacy artifact: %v", err)
	}

	measure, err := measureSQLiteTarget(context.Background(), target)
	if err != nil {
		t.Fatalf("measure legacy target: %v", err)
	}
	err = verifyConvertedShape(measure, "", 0)
	if err == nil {
		t.Fatal("a pre-conversion artifact must not satisfy the converted contract")
	}
	kind := KindOf(err)
	switch kind {
	case KindTimeContractMismatch, KindTemporalColumnSetIncomplete:
		// expected
	case KindArtifactNotFound, KindArtifactUnreadable, KindToolFailure:
		t.Fatalf("legacy artifact was rejected with the wrong classification %q: %v", kind, err)
	default:
		t.Fatalf("legacy artifact classification = %q, want TimeContractMismatch or TemporalColumnSetIncomplete: %v", kind, err)
	}

	// The same rejection must hold through the service: CreateRecoveryPoint must
	// never return a recovery point for a pre-conversion source.
	service := NewService(t.TempDir())
	if _, err := service.CreateRecoveryPoint(context.Background(), kernel.RecoveryPointRequest{
		Dialect:        kernel.DialectSQLite,
		SourceID:       legacyPath,
		CatalogVersion: len(catalog),
		TimeContract:   kernel.TimeContractVP040Timestamptz,
	}); err == nil {
		t.Fatal("service returned a recovery point for a pre-conversion source")
	} else if k := KindOf(err); k != KindTimeContractMismatch && k != KindTemporalColumnSetIncomplete {
		t.Fatalf("service classification = %q, want a contract classification: %v", k, err)
	}
}

// TestMissingArtifactIsNotAContractFailure guards the classification boundary in
// the other direction: an absent artifact must be reported as such.
func TestMissingArtifactIsNotAContractFailure(t *testing.T) {
	provider := SQLiteProvider{}
	err := provider.RestoreInto(context.Background(), filepath.Join(t.TempDir(), "absent.db"), filepath.Join(t.TempDir(), "target.db"))
	if err == nil {
		t.Fatal("restoring an absent artifact must fail")
	}
	if kind := KindOf(err); kind != KindArtifactNotFound {
		t.Fatalf("classification = %q, want ArtifactNotFound (%v)", kind, err)
	}

}

// TestPerMigrationSnapshotIsNotARecoveryPointSource pins the C3 §2 hard rule:
// the batch-boundary snapshot mechanism is class C, so a database captured
// mid-batch (mixed shape) must fail the converted-contract verification.
func TestPerMigrationSnapshotIsNotARecoveryPointSource(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "midbatch.db")
	catalog := migrationCatalog(t)
	// Apply only the pre-conversion history plus v73: the batch is mid-flight and
	// the shape is mixed (schema_migrations converted, users still INTEGER).
	st, err := store.OpenWithCatalog(path, catalog[:v73Head])
	if err != nil {
		t.Fatalf("open mid-batch: %v", err)
	}
	if err := st.Close(); err != nil {
		t.Fatalf("close mid-batch: %v", err)
	}
	provider := SQLiteProvider{}
	artifact := filepath.Join(dir, "midbatch.artifact")
	if err := provider.Create(context.Background(), path, artifact); err != nil {
		t.Fatalf("create mid-batch artifact: %v", err)
	}
	target := filepath.Join(dir, "midbatch-restored.db")
	if err := provider.RestoreInto(context.Background(), artifact, target); err != nil {
		t.Fatalf("restore mid-batch artifact: %v", err)
	}
	measure, err := measureSQLiteTarget(context.Background(), target)
	if err != nil {
		t.Fatalf("measure mid-batch target: %v", err)
	}
	if err := verifyConvertedShape(measure, "", 0); err == nil {
		t.Fatal("a mid-batch (mixed shape) artifact must not satisfy the converted contract")
	} else if kind := KindOf(err); kind != KindTemporalColumnSetIncomplete && kind != KindTimeContractMismatch {
		t.Fatalf("mid-batch classification = %q, want a contract classification: %v", kind, err)
	}
}
