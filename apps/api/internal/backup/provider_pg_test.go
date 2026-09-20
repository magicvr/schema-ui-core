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
