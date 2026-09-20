package store

import (
	"context"
	"database/sql"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/backup"
	"github.com/magicvr/schema-ui-core/apps/api/internal/pgtest"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

const pgAnchorClientImage = "postgres:15-alpine"

// TestC3RecoveryAnchorsOnPostgresUpgrade exercises the C3 §4.2/§4.3 postgres
// anchors end to end: the class-A rollback dump before the batch, the
// per-migration class-C dump, and the class-B recovery point (pg_dump →
// pg_restore into a new database → verify) with a detectable marker. The upgrade
// is one migration long (v86 → v87) so the dump cost stays bounded.
func TestC3RecoveryAnchorsOnPostgresUpgrade(t *testing.T) {
	dsn := pgtest.DSN()
	if dsn == "" {
		t.Skip("postgres test env not set (PG_TEST_*); skipping the postgres anchor harness")
	}
	if _, err := exec.LookPath("docker"); err != nil {
		t.Skipf("docker unavailable (%v); the host has no pg_dump, so the postgres anchors cannot run", err)
	}
	parsed, err := url.Parse(dsn)
	if err != nil {
		t.Fatalf("parse dsn: %v", err)
	}
	serverDSN := parsed.String()
	const dbName = "vp040anchors"
	parsed.Path = "/" + dbName
	scratchDSN := parsed.String()

	admin, err := sql.Open("pgx", serverDSN)
	if err != nil {
		t.Fatalf("admin open: %v", err)
	}
	defer admin.Close()
	if _, err := admin.Exec(`DROP DATABASE IF EXISTS ` + dbName + ` WITH (FORCE)`); err != nil {
		t.Fatalf("drop prior: %v", err)
	}
	if _, err := admin.Exec(`CREATE DATABASE ` + dbName); err != nil {
		t.Fatalf("create: %v", err)
	}
	t.Cleanup(func() {
		_, _ = admin.Exec(`DROP DATABASE IF EXISTS ` + dbName + ` WITH (FORCE)`)
	})

	catalog := MigrationCatalog()
	ctx := context.Background()
	// v1..v86 first, so the upgrade below is exactly one migration.
	st, err := Open(ctx, OpenOptions{
		Dialect: kernel.DialectPostgres, DSN: scratchDSN, ConnectTimeout: 20 * time.Second,
	}, catalog[:86])
	if err != nil {
		t.Fatalf("bootstrap v86: %v", err)
	}
	if err := st.Close(); err != nil {
		t.Fatalf("close v86: %v", err)
	}

	workDir := t.TempDir()
	provider := backup.PgProvider{AdminDSN: serverDSN, ClientImage: pgAnchorClientImage, WorkDir: workDir}
	service := backup.NewService(workDir)
	service.RegisterProvider(provider)

	upgraded, err := Open(ctx, OpenOptions{
		Dialect:           kernel.DialectPostgres,
		DSN:               scratchDSN,
		ConnectTimeout:    20 * time.Second,
		RecoveryPoints:    service,
		RollbackArtifacts: provider,
		ArtifactDir:       workDir,
	}, catalog)
	if err != nil {
		t.Fatalf("upgrade with anchors: %v", err)
	}
	defer upgraded.Close()

	pg := upgraded.(*postgres)
	if note := pg.RecoveryNote(); note != "" {
		t.Fatalf("unexpected recovery note: %q", note)
	}
	state, err := pg.probeRecoveryState(ctx)
	if err != nil {
		t.Fatalf("probe recovery state: %v", err)
	}
	if !state.Converted {
		t.Fatalf("postgres shape is not converted after the batch: %s", state.Detail)
	}
	if !state.HasPoint {
		t.Fatal("class-B recovery point marker was not recorded on the database")
	}

	// Class A (one batch-rollback dump) and class C (one pre-v87 dump).
	entries, err := os.ReadDir(workDir)
	if err != nil {
		t.Fatalf("read workdir: %v", err)
	}
	var classA, classC int
	for _, entry := range entries {
		switch {
		case filepath.Ext(entry.Name()) != ".dump":
			continue
		case strings.Contains(entry.Name(), ".batch-rollback-"):
			classA++
		case strings.Contains(entry.Name(), ".pre-v"):
			classC++
		}
	}
	if classA != 1 {
		t.Fatalf("class A dumps = %d, want 1 (%v)", classA, entries)
	}
	if classC != 1 {
		t.Fatalf("class C dumps = %d, want 1 (%v)", classC, entries)
	}

	// Reopening at head with the marker present is a noop (no second B).
	again, err := Open(ctx, OpenOptions{
		Dialect:           kernel.DialectPostgres,
		DSN:               scratchDSN,
		ConnectTimeout:    20 * time.Second,
		RecoveryPoints:    service,
		RollbackArtifacts: provider,
		ArtifactDir:       workDir,
	}, catalog)
	if err != nil {
		t.Fatalf("reopen at head: %v", err)
	}
	defer again.Close()
	reopenedState, err := again.(*postgres).probeRecoveryState(ctx)
	if err != nil {
		t.Fatalf("probe reopened state: %v", err)
	}
	if !reopenedState.HasPoint {
		t.Fatal("marker disappeared on reopen")
	}

	// Removing the marker makes "no B" detectable: the reopen takes the explicit
	// action and re-creates the recovery point once.
	if _, err := admin.Exec(`COMMENT ON DATABASE ` + dbName + ` IS NULL`); err != nil {
		t.Fatalf("clear marker: %v", err)
	}
	after, err := Open(ctx, OpenOptions{
		Dialect:           kernel.DialectPostgres,
		DSN:               scratchDSN,
		ConnectTimeout:    20 * time.Second,
		RecoveryPoints:    service,
		RollbackArtifacts: provider,
		ArtifactDir:       workDir,
	}, catalog)
	if err != nil {
		t.Fatalf("reopen after marker loss: %v", err)
	}
	defer after.Close()
	recreated, err := after.(*postgres).probeRecoveryState(ctx)
	if err != nil {
		t.Fatalf("probe re-created state: %v", err)
	}
	if !recreated.HasPoint {
		t.Fatal("marker was not re-created after the explicit action")
	}
	if note := after.(*postgres).RecoveryNote(); strings.Contains(note, "failed") {
		t.Fatalf("re-create reported a failure: %q", note)
	}
}
