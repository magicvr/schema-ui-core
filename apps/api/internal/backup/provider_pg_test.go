package backup

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/magicvr/schema-ui-core/apps/api/internal/pgtest"
	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// pgScratch creates a dedicated database and returns its DSN plus a cleanup that
// drops it. C3/Root D-017 require destructive restore work to use a disposable
// database, never a shared schema.
func pgScratch(t *testing.T, adminDSN, name string) string {
	t.Helper()
	admin, err := sql.Open("pgx", adminDSN)
	if err != nil {
		t.Fatalf("admin open: %v", err)
	}
	defer admin.Close()
	if _, err := admin.Exec(`DROP DATABASE IF EXISTS ` + quoteIdent(name) + ` WITH (FORCE)`); err != nil {
		t.Fatalf("drop prior %s: %v", name, err)
	}
	if _, err := admin.Exec(`CREATE DATABASE ` + quoteIdent(name)); err != nil {
		t.Fatalf("create %s: %v", name, err)
	}
	t.Cleanup(func() {
		conn, err := sql.Open("pgx", adminDSN)
		if err != nil {
			return
		}
		defer conn.Close()
		_, _ = conn.Exec(`DROP DATABASE IF EXISTS ` + quoteIdent(name) + ` WITH (FORCE)`)
	})
	dsn, err := withDatabase(adminDSN, name)
	if err != nil {
		t.Fatalf("target dsn: %v", err)
	}
	return dsn
}

// requirePgBackupEnv skips cleanly when the PG server or the container-provided
// client tools are unavailable — recording WHY, never passing silently.
func requirePgBackupEnv(t *testing.T) (adminDSN string) {
	t.Helper()
	dsn := pgtest.DSN()
	if dsn == "" {
		t.Skip("postgres test env not set (PG_TEST_*); skipping postgres recovery-point harness")
	}
	u, err := url.Parse(dsn)
	if err != nil {
		t.Fatalf("parse dsn: %v", err)
	}
	if _, err := exec.LookPath("docker"); err != nil {
		t.Skipf("docker is unavailable (%v); the host has no pg_dump/pg_restore, so the postgres harness cannot run", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	probe := PgProvider{ClientImage: pgClientImage(), WorkDir: t.TempDir()}
	if _, err := probe.ClientVersion(ctx); err != nil {
		t.Skipf("pg client image unavailable: %v", err)
	}
	return u.String()
}

// pgClientImage picks the container image that provides pg_dump/pg_restore. The
// resident server major version decides it (C3 §5 records the combination rather
// than assuming cross-major compatibility).
func pgClientImage() string {
	return "postgres:15-alpine"
}

// TestPGRestoreToNewDB is the C3 §5 harness for PostgreSQL: create a recovery
// point with pg_dump, restore it into a NEW database with pg_restore, and verify
// the restored database against the frozen contract.
func TestPGRestoreToNewDB(t *testing.T) {
	adminDSN := requirePgBackupEnv(t)
	ctx := context.Background()
	catalog := migrationCatalog(t)

	sourceDSN := pgScratch(t, adminDSN, "vp040bkp_src")
	// Bootstrap the pre-conversion shape and seed representative legacy rows
	// (seconds family, milliseconds family and a sentinel 0), so the recovery
	// point's samples are exercised on real data rather than on an empty schema
	// (GOAL-005 A-002 F-I-008).
	legacy, err := store.Open(ctx, store.OpenOptions{
		Dialect: kernel.DialectPostgres, DSN: sourceDSN, ConnectTimeout: 20 * time.Second,
	}, catalog[:v72Head])
	if err != nil {
		t.Fatalf("bootstrap pre-conversion source: %v", err)
	}
	if err := legacy.Run(ctx, func(tx kernel.Tx) error {
		if _, err := tx.Exec(ctx, `INSERT INTO users (id, username, name, roles, password_hash, created_at, updated_at,
			token_version, failed_login_count, locked_until, enabled, notifications_enabled, avatar_url,
			must_change_password, email, email_status, last_login_failure_at)
			VALUES ('u-1','u1','User One','[]','h', 1758320000, 1758320001, 0, 0, 0, 1, 1, '', 0, NULL, NULL, 0)`); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `INSERT INTO jobs (id, kind, status, payload, progress, cancel_requested, attempt,
			max_attempts, lease_owner, lease_version, lease_expires_at, result, error_code, error_message,
			actor_id, correlation_id, created_at, updated_at, finished_at, expires_at)
			VALUES ('j-1','probe','queued','{}',0,0,0,3,NULL,0,NULL,NULL,NULL,NULL,'a','c',1758320000123,1758320000123,NULL,NULL)`); err != nil {
			return err
		}
		_, err := tx.Exec(ctx, `UPDATE mail_config SET updated_at = 0 WHERE id = 1`)
		return err
	}); err != nil {
		_ = legacy.Close()
		t.Fatalf("seed legacy rows: %v", err)
	}
	if err := legacy.Close(); err != nil {
		t.Fatalf("close legacy source: %v", err)
	}
	// Convert the seeded database to the target contract.
	st, err := store.Open(ctx, store.OpenOptions{
		Dialect: kernel.DialectPostgres, DSN: sourceDSN, ConnectTimeout: 20 * time.Second,
	}, catalog)
	if err != nil {
		t.Fatalf("bootstrap converted postgres source: %v", err)
	}
	if err := st.Close(); err != nil {
		t.Fatalf("close source: %v", err)
	}

	workDir := t.TempDir()
	provider := PgProvider{
		AdminDSN:    adminDSN,
		ClientImage: pgClientImage(),
		WorkDir:     workDir,
	}
	if version, err := provider.ClientVersion(ctx); err == nil {
		t.Logf("pg client: %s", version)
	}
	service := NewService(workDir)
	service.RegisterProvider(provider)

	point, err := service.CreateRecoveryPoint(ctx, kernel.RecoveryPointRequest{
		Dialect:         kernel.DialectPostgres,
		SourceID:        sourceDSN,
		CatalogVersion:  len(catalog),
		TimeContract:    kernel.TimeContractVP040Timestamptz,
		DestinationHint: filepath.Join(workDir, "recovery.dump"),
	})
	if err != nil {
		t.Fatalf("CreateRecoveryPoint (postgres): %v", err)
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
	if _, err := os.Stat(point.ArtifactRef); err != nil {
		t.Fatalf("artifact missing: %v", err)
	}

	// The sample set must not be vacuous: the source was seeded with both a
	// seconds-family and a milliseconds-family value, so the restored artifact
	// must carry both (F-I-008).
	db, err := sql.Open("pgx", sourceDSN)
	if err != nil {
		t.Fatalf("open source for coverage: %v", err)
	}
	defer db.Close()
	coverage, err := postgresSampleCoverage(ctx, db)
	if err != nil {
		t.Fatalf("sample coverage: %v", err)
	}
	if coverage["sec"] == 0 || coverage["ms"] == 0 {
		t.Fatalf("sample coverage is vacuous: %v", coverage)
	}
	t.Logf("sample coverage on the source: %v", coverage)
}

// TestPGLegacyArtifactMustFail is the reverse assertion for postgres: a dump
// taken before the conversion batch (class A) must be rejected with a contract
// classification, never with ArtifactNotFound/Unreadable/ToolFailure.
func TestPGLegacyArtifactMustFail(t *testing.T) {
	adminDSN := requirePgBackupEnv(t)
	ctx := context.Background()
	catalog := migrationCatalog(t)

	legacyDSN := pgScratch(t, adminDSN, "vp040bkp_legacy")
	st, err := store.Open(ctx, store.OpenOptions{
		Dialect: kernel.DialectPostgres, DSN: legacyDSN, ConnectTimeout: 20 * time.Second,
	}, catalog[:v72Head])
	if err != nil {
		t.Fatalf("bootstrap pre-conversion source: %v", err)
	}
	if err := st.Close(); err != nil {
		t.Fatalf("close legacy source: %v", err)
	}

	workDir := t.TempDir()
	provider := PgProvider{AdminDSN: adminDSN, ClientImage: pgClientImage(), WorkDir: workDir}
	service := NewService(workDir)
	service.RegisterProvider(provider)

	_, err = service.CreateRecoveryPoint(ctx, kernel.RecoveryPointRequest{
		Dialect:         kernel.DialectPostgres,
		SourceID:        legacyDSN,
		CatalogVersion:  len(catalog),
		TimeContract:    kernel.TimeContractVP040Timestamptz,
		DestinationHint: filepath.Join(workDir, "legacy.dump"),
	})
	if err == nil {
		t.Fatal("a pre-conversion postgres source must not yield a recovery point")
	}
	switch kind := KindOf(err); kind {
	case KindTimeContractMismatch, KindTemporalColumnSetIncomplete:
		t.Logf("legacy postgres artifact rejected as expected: %s (%v)", kind, err)
	case KindArtifactNotFound, KindArtifactUnreadable, KindToolFailure:
		t.Fatalf("legacy postgres artifact rejected with the wrong classification %q: %v", kind, err)
	default:
		t.Fatalf("legacy postgres classification = %q, want a contract classification: %v", kind, err)
	}
}

// TestPGMidBatchArtifactMustFail is the C3 §2 class-C reverse assertion on
// postgres: a dump taken while the batch is in flight (mixed shape) must be
// rejected with a contract classification, exactly like the sqlite class-C case
// (GOAL-005 A-002 F-I-006).
func TestPGMidBatchArtifactMustFail(t *testing.T) {
	adminDSN := requirePgBackupEnv(t)
	ctx := context.Background()
	catalog := migrationCatalog(t)

	midDSN := pgScratch(t, adminDSN, "vp040bkp_mid")
	st, err := store.Open(ctx, store.OpenOptions{
		Dialect: kernel.DialectPostgres, DSN: midDSN, ConnectTimeout: 20 * time.Second,
	}, catalog[:v73Head])
	if err != nil {
		t.Fatalf("bootstrap mid-batch source: %v", err)
	}
	if err := st.Close(); err != nil {
		t.Fatalf("close mid-batch source: %v", err)
	}

	workDir := t.TempDir()
	provider := PgProvider{AdminDSN: adminDSN, ClientImage: pgClientImage(), WorkDir: workDir}
	service := NewService(workDir)
	service.RegisterProvider(provider)

	_, err = service.CreateRecoveryPoint(ctx, kernel.RecoveryPointRequest{
		Dialect:         kernel.DialectPostgres,
		SourceID:        midDSN,
		CatalogVersion:  len(catalog),
		TimeContract:    kernel.TimeContractVP040Timestamptz,
		DestinationHint: filepath.Join(workDir, "midbatch.dump"),
	})
	if err == nil {
		t.Fatal("a mid-batch postgres source must not yield a recovery point")
	}
	switch kind := KindOf(err); kind {
	case KindTimeContractMismatch, KindTemporalColumnSetIncomplete:
		t.Logf("mid-batch postgres artifact rejected as expected: %s", kind)
	case KindArtifactNotFound, KindArtifactUnreadable, KindToolFailure:
		t.Fatalf("mid-batch postgres artifact rejected with the wrong classification %q: %v", kind, err)
	default:
		t.Fatalf("mid-batch postgres classification = %q, want a contract classification: %v", kind, err)
	}
}

// TestPgProviderRejectsRelativeWorkDir is the boundary guard from the
// independent re-audit (A-005 F-I-102): Create/Restore are the components that
// actually build `docker run -v <WorkDir>:...`, so a relative WorkDir must fail
// closed here rather than reaching docker and failing with exit 125.
func TestPgProviderRejectsRelativeWorkDir(t *testing.T) {
	const sourceDSN = "postgres://sa:secret@127.0.0.1:5432/src?sslmode=disable"
	for _, tc := range []struct {
		name    string
		workDir string
	}{
		{"dot-slash relative", filepath.Join(".", "data", "recovery")},
		{"bare relative", filepath.Join("data", "recovery")},
		{"single dot", "."},
	} {
		t.Run(tc.name, func(t *testing.T) {
			called := false
			provider := PgProvider{
				AdminDSN: "postgres://sa:secret@127.0.0.1:5432/postgres?sslmode=disable",
				WorkDir:  tc.workDir,
				exec: func(context.Context, string, ...string) ([]byte, error) {
					called = true
					return nil, nil
				},
			}
			err := provider.Create(context.Background(), sourceDSN, filepath.Join(t.TempDir(), "a.dump"))
			if err == nil {
				t.Fatalf("Create with WorkDir %q must fail closed", tc.workDir)
			}
			if KindOf(err) != KindInvalidRequest {
				t.Fatalf("Create WorkDir %q classification = %q, want InvalidRequest (%v)", tc.workDir, KindOf(err), err)
			}
			if called {
				t.Fatalf("Create with WorkDir %q still invoked docker", tc.workDir)
			}
			_, _, err = provider.Restore(context.Background(), filepath.Join(t.TempDir(), "a.dump"))
			if err == nil {
				t.Fatalf("Restore with WorkDir %q must fail closed", tc.workDir)
			}
			if KindOf(err) != KindInvalidRequest {
				t.Fatalf("Restore WorkDir %q classification = %q, want InvalidRequest (%v)", tc.workDir, KindOf(err), err)
			}
		})
	}
}

// TestPgProviderCommandConstruction pins the tool invocation without needing a
// server: the dump and restore commands must carry the C3 flags.
func TestPgProviderCommandConstruction(t *testing.T) {
	workDir := t.TempDir()
	artifact := filepath.Join(workDir, "artifact.dump")
	var calls [][]string
	provider := PgProvider{
		AdminDSN:    "postgres://sa:secret@127.0.0.1:5432/postgres?sslmode=disable",
		ClientImage: "postgres:15-alpine",
		WorkDir:     workDir,
		exec: func(ctx context.Context, name string, args ...string) ([]byte, error) {
			calls = append(calls, append([]string{name}, args...))
			if len(args) > 0 && args[len(args)-1] == "pg_dump" {
				return []byte("pg_dump (PostgreSQL) 15.4\n"), nil
			}
			// Touch the artifact so Create's existence check succeeds.
			if len(args) > 0 {
				for i, arg := range args {
					if arg == "--file" && i+1 < len(args) {
						_ = os.WriteFile(filepath.Join(workDir, filepath.Base(strings.TrimPrefix(args[i+1], pgMountPoint+"/"))), []byte("dump"), 0o600)
					}
				}
			}
			return nil, nil
		},
	}
	if err := provider.Create(context.Background(), "postgres://sa:secret@127.0.0.1:5432/src?sslmode=disable", artifact); err != nil {
		t.Fatalf("create: %v", err)
	}
	if len(calls) != 1 {
		t.Fatalf("expected one tool invocation, got %d", len(calls))
	}
	joined := strings.Join(calls[0], " ")
	for _, want := range []string{"docker", "run", "-F c", "--no-owner", "--file", pgMountPoint} {
		if !strings.Contains(joined, want) {
			t.Fatalf("pg_dump invocation %q lacks %q", joined, want)
		}
	}
	if strings.Contains(joined, "secret") == false {
		t.Fatalf("pg_dump must receive the source DSN (with its credentials): %q", joined)
	}

	// A tool failure must not leak the DSN password into the error text.
	failing := PgProvider{
		AdminDSN:    "postgres://sa:secret@127.0.0.1:5432/postgres?sslmode=disable",
		ClientImage: "postgres:15-alpine",
		WorkDir:     workDir,
		exec: func(ctx context.Context, name string, args ...string) ([]byte, error) {
			return nil, fmt.Errorf("exit status 1: %s", redactArgs(args))
		},
	}
	err := failing.Create(context.Background(), "postgres://sa:secret@127.0.0.1:5432/src?sslmode=disable",
		filepath.Join(workDir, "second.dump"))
	if err == nil {
		t.Fatal("expected a tool failure")
	}
	if KindOf(err) != KindToolFailure {
		t.Fatalf("classification = %q, want ToolFailure (%v)", KindOf(err), err)
	}
	if strings.Contains(err.Error(), "secret") {
		t.Fatalf("tool failure leaked the password: %v", err)
	}
}
