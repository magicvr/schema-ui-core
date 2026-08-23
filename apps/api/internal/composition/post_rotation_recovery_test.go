package composition

// VP-016 R3 · workspace-016 GOAL-004 D-001: post-rotation recovery evidence.
//
// The signing keys live in configuration, never in the database, so restoring
// a pre-rotation backup under a rotated key set must "just work" — but that is
// exactly the kind of claim that drifts silently. These tests pin the whole
// loop end to end (T0–T5 of GOAL-004 D-001):
//
//	T0  run under K1, produce a real session (login → access[K1] + refresh)
//	T1  take a backup with the dialect's documented contract
//	    (SQLite: VACUUM INTO · PG: pg_dump -F c → pg_restore)
//	T3  restore the backup as a usable database
//	T4  boot the FULL application from the restore under current=K2,
//	    previous=K1 → every module Start+Ready gate must pass
//	T5  auth assertions on the restored DB under the rotated key set:
//	    A1 old access[K1] still verifies (overlap window)
//	    A2 fresh login issues K2-only access (issuance never uses previous)
//	    A3 the pre-rotation refresh token still refreshes (opaque session
//	       continuity across rotation AND restore)
//
// No backup implementation is added here: both dialects use their existing,
// operator-documented contracts (VP-013), consumed exactly as an operator
// would run them.

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/config"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	authsession "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession"
	"github.com/magicvr/schema-ui-core/apps/api/internal/pgtest"
	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)

const (
	recoveryK1     = "r16-recovery-old-signing-key-01"
	recoveryK2     = "r16-recovery-new-signing-key-02"
	recoveryUser   = "admin"
	recoveryPass   = "pw"
	recoveryAccess = 15 * time.Minute
)

// recoveryLogin issues a real token pair through the production issue path.
func recoveryLogin(t *testing.T, repository auth.Repository, current []byte) (accessOld, refreshOld string) {
	t.Helper()
	a := auth.NewWithRepositoryAndPrevious(current, nil, recoveryAccess, 30*24*time.Hour, repository, false)
	access, refresh, _, err := a.Login(recoveryUser, recoveryPass, time.Now().UTC())
	if err != nil {
		t.Fatalf("login under pre-rotation key: %v", err)
	}
	return access, refresh
}

// recoveryServeBearer drives the composed middleware semantics: 200+handler
// means the token verified against the authenticator's key set.
func recoveryServeBearer(t *testing.T, a *auth.Authenticator, token string) (*httptest.ResponseRecorder, bool) {
	t.Helper()
	passed := false
	protected := a.Middleware(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { passed = true }))
	request := httptest.NewRequest(http.MethodGet, "/api/resources/widgets", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	response := httptest.NewRecorder()
	protected.ServeHTTP(response, request)
	return response, passed
}

// recoveryAssertRotated runs T5 (A1/A2/A3) against a restored repository.
func recoveryAssertRotated(t *testing.T, repository auth.Repository, accessOld, refreshOld string) {
	t.Helper()
	a2 := auth.NewWithRepositoryAndPrevious([]byte(recoveryK2), []byte(recoveryK1), recoveryAccess, 30*24*time.Hour, repository, false)

	// A1 · the overlap window: old-key access still verifies post-restore.
	response, passed := recoveryServeBearer(t, a2, accessOld)
	if response.Code != http.StatusOK || !passed {
		t.Fatalf("A1 old-key access rejected after rotation+restore: code=%d passed=%v body=%s", response.Code, passed, response.Body.String())
	}

	// A2 · issuance uses only the current key.
	freshAccess, _, _, err := a2.Login(recoveryUser, recoveryPass, time.Now().UTC())
	if err != nil {
		t.Fatalf("A2 login on restored DB: %v", err)
	}
	if _, err := auth.ParseAccessToken([]byte(recoveryK2), freshAccess); err != nil {
		t.Fatalf("A2 fresh access must verify under K2: %v", err)
	}
	if _, err := auth.ParseAccessToken([]byte(recoveryK1), freshAccess); err == nil {
		t.Fatal("A2 fresh access verified under the previous key")
	}

	// A3 · opaque refresh survives rotation AND restore.
	newAccess, newRefresh, _, err := a2.Refresh(refreshOld, time.Now().UTC())
	if err != nil {
		t.Fatalf("A3 pre-rotation refresh must survive rotation+restore: %v", err)
	}
	if newAccess == "" || newRefresh == "" {
		t.Fatal("A3 refresh returned an empty pair")
	}
	if _, err := auth.ParseAccessToken([]byte(recoveryK2), newAccess); err != nil {
		t.Fatalf("A3 refreshed access must verify under K2: %v", err)
	}
}

// TestSQLitePostRotationRecovery pins the SQLite contract end to end:
// VACUUM INTO backup of a K1-era database restores into a fully-booting
// application under rotated keys, and the rotated auth contract holds.
func TestSQLitePostRotationRecovery(t *testing.T) {
	dir := t.TempDir()
	livePath := filepath.Join(dir, "live.sqlite")
	backupPath := filepath.Join(dir, "backup.sqlite")

	// T0 · K1-era database with a real session.
	hash, err := auth.HashPassword(recoveryPass, 4)
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	st0, err := testsupport.OpenStore(livePath, recoveryUser, hash, true)
	if err != nil {
		t.Fatalf("open live store: %v", err)
	}
	accessOld, refreshOld := recoveryLogin(t, authsession.NewRepository(st0), []byte(recoveryK1))
	if err := st0.Close(); err != nil {
		t.Fatalf("close live store: %v", err)
	}

	// T1 · the documented SQLite backup: VACUUM INTO.
	backupDB, err := sql.Open("sqlite", livePath)
	if err != nil {
		t.Fatalf("open live db for backup: %v", err)
	}
	if _, err := backupDB.ExecContext(context.Background(), "VACUUM INTO '"+backupPath+"'"); err != nil {
		_ = backupDB.Close()
		t.Fatalf("VACUUM INTO %s: %v", backupPath, err)
	}
	if err := backupDB.Close(); err != nil {
		t.Fatalf("close backup handle: %v", err)
	}

	// T3/T4 · restore = the backup file becomes db.path; the FULL application
	// must boot under current=K2, previous=K1.
	cfg := &config.Config{
		AppName:               "r16-sqlite-recovery",
		AppEnv:                "development",
		HTTPAddr:              "127.0.0.1:0",
		DBDialect:             "sqlite",
		DBPath:                backupPath,
		ProfileName:           "mvp",
		ReadTimeout:           time.Second,
		WriteTimeout:          time.Second,
		IdleTimeout:           time.Second,
		AuthJWTSecretPrevious: recoveryK1,
	}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	app, err := NewApp(cfg, recoveryK2, hash, logger)
	if err != nil {
		t.Fatalf("build restored composition: %v", err)
	}
	startCtx, startCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer startCancel()
	if err := app.Start(startCtx); err != nil {
		t.Fatalf("T4 restored-backup composition failed to start under rotated keys: %v", err)
	}
	stopCtx, stopCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer stopCancel()
	if err := app.Stop(stopCtx); err != nil {
		t.Fatalf("stop restored composition: %v", err)
	}

	// T5 · rotated auth assertions on the restored database.
	st1, err := testsupport.OpenStore(backupPath, recoveryUser, hash, true)
	if err != nil {
		t.Fatalf("open restored store: %v", err)
	}
	defer st1.Close()
	recoveryAssertRotated(t, authsession.NewRepository(st1), accessOld, refreshOld)
}

// pgTooling resolves how this environment reaches the documented pg_dump /
// pg_restore clients: either locally installed or inside a Docker container
// named by R16_PG_DUMP_CONTAINER (postgres images ship both tools).
type pgTooling struct {
	useDocker bool
	container string
	host      string
	port      string
	user      string
	password  string
}

func resolvePgTooling(t *testing.T) *pgTooling {
	t.Helper()
	adminDSN := pgtest.DSN()
	if adminDSN == "" {
		t.Skip("postgres test env not set (PG_TEST_*); skipping pg rotation-recovery evidence")
	}
	_, dumpErr := exec.LookPath("pg_dump")
	_, restoreErr := exec.LookPath("pg_restore")
	container := strings.TrimSpace(os.Getenv("R16_PG_DUMP_CONTAINER"))
	if dumpErr != nil || restoreErr != nil {
		if container == "" {
			t.Skip("pg_dump/pg_restore not on PATH and R16_PG_DUMP_CONTAINER unset; skipping pg rotation-recovery evidence")
		}
	}
	u, err := pgURL(adminDSN)
	if err != nil {
		t.Fatalf("parse admin dsn: %v", err)
	}
	password, _ := u.User.Password()
	return &pgTooling{
		useDocker: container != "",
		container: container,
		host:      u.Hostname(),
		port:      u.Port(),
		user:      u.User.Username(),
		password:  password,
	}
}

// run executes name+args locally or inside the configured container, feeding
// stdin and returning stdout. PGPASSWORD reaches the server-side client in
// both modes without ever being logged.
func (p *pgTooling) run(t *testing.T, stdin []byte, name string, args ...string) ([]byte, error) {
	t.Helper()
	var cmd *exec.Cmd
	env := os.Environ()
	if p.useDocker {
		full := append([]string{"exec"}, "-e", "PGPASSWORD="+p.password)
		if len(stdin) > 0 {
			full = append(full, "-i")
		}
		full = append(full, p.container, name)
		cmd = exec.Command("docker", append(full, args...)...)
	} else {
		cmd = exec.Command(name, args...)
		env = append(env, "PGPASSWORD="+p.password)
		cmd.Env = env
	}
	if len(stdin) > 0 {
		cmd.Stdin = bytes.NewReader(stdin)
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("%s %s: %w (stderr: %s)", name, strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}
	return out, nil
}

func (p *pgTooling) dump(t *testing.T, dbName string) []byte {
	t.Helper()
	// The client always targets the DSN host explicitly: when running inside
	// a helper container, its own local socket is a different server. A newer
	// client dumping an older server is the documented VP-013 combination
	// (GOAL-006 D-002: restore may log unknown-GUC SET warnings; exit stays 0).
	out, err := p.run(t, nil, "pg_dump", "--no-password", "-h", p.host, "-p", p.port, "-U", p.user, "-F", "c", dbName)
	if err != nil {
		t.Fatalf("pg_dump %s: %v", dbName, err)
	}
	if len(out) == 0 {
		t.Fatalf("pg_dump %s produced an empty archive", dbName)
	}
	return out
}

func (p *pgTooling) restore(t *testing.T, dbName string, dump []byte) {
	t.Helper()
	_, err := p.run(t, dump, "pg_restore", "--no-password", "-h", p.host, "-p", p.port, "-U", p.user, "-d", dbName)
	if err == nil {
		return
	}
	// [workspace-013] GOAL-006 D-002 (its A-001 F-004): a newer pg_dump client
	// against an older server archives `SET transaction_timeout = 0`, which
	// the older pg_restore reports as ONE ignored error. The contract allows
	// exactly this warning class; the caller still proves restore fidelity by
	// comparing the migration ledgers (see assertRestoredLedger).
	stderr := err.Error()
	if strings.Contains(stderr, "unrecognized configuration parameter") && strings.Contains(stderr, "errors ignored on restore: 1") {
		return
	}
	t.Fatalf("pg_restore into %s: %v", dbName, err)
}

// assertRestoredLedger proves the restored database carries the exact
// migration ledger of the live one (count + aggregate checksum), so the
// tolerated GUC warning above cannot mask a partial restore. Postgres cannot
// qualify across databases, so each side opens its own connection.
func assertRestoredLedger(t *testing.T, adminDSN, liveDB, restoredDB string) {
	t.Helper()
	ledgerFingerprint := func(t *testing.T, db string) (int, string) {
		t.Helper()
		conn, err := sql.Open("pgx", swapPGDatabase(adminDSN, db))
		if err != nil {
			t.Fatalf("open %s for ledger: %v", db, err)
		}
		defer conn.Close()
		var count int
		var sum string
		if err := conn.QueryRow(`SELECT count(*), coalesce(md5(string_agg(version || ':' || checksum, ',' ORDER BY version)), '') FROM schema_migrations`).Scan(&count, &sum); err != nil {
			t.Fatalf("read %s ledger: %v", db, err)
		}
		return count, sum
	}
	liveCount, liveSum := ledgerFingerprint(t, liveDB)
	restCount, restSum := ledgerFingerprint(t, restoredDB)
	if liveCount == 0 || liveCount != restCount || liveSum != restSum {
		t.Fatalf("restored ledger mismatch: live=%d/%s restored=%d/%s", liveCount, liveSum, restCount, restSum)
	}
}

// TestPostgresPostRotationRecovery pins the PG contract end to end: a
// pg_dump -F c archive of a K1-era database restores via pg_restore into a
// fresh database that boots the FULL application under rotated keys, with the
// same rotated auth assertions as the SQLite loop. Gated twice: PG_TEST_* for
// the server, and reachable pg_dump/pg_restore clients for the documented
// backup contract.
func TestPostgresPostRotationRecovery(t *testing.T) {
	tooling := resolvePgTooling(t)
	adminDSN := pgtest.DSN()
	ctx := context.Background()

	const liveDB = "r16r3a"
	const restoredDB = "r16r3b"

	admin, err := sql.Open("pgx", adminDSN)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = admin.ExecContext(context.Background(), `DROP DATABASE IF EXISTS `+liveDB+` WITH (FORCE)`)
		_, _ = admin.ExecContext(context.Background(), `DROP DATABASE IF EXISTS `+restoredDB+` WITH (FORCE)`)
		_ = admin.Close()
	})
	for _, db := range []string{liveDB, restoredDB} {
		if _, err := admin.ExecContext(ctx, `DROP DATABASE IF EXISTS `+db+` WITH (FORCE)`); err != nil {
			t.Fatalf("drop prior scratch db %s: %v", db, err)
		}
	}
	if _, err := admin.ExecContext(ctx, `CREATE DATABASE `+liveDB); err != nil {
		t.Fatalf("create live scratch db: %v", err)
	}

	liveDSN := swapPGDatabase(adminDSN, liveDB)

	// T4-equivalent for the LIVE db happens implicitly: the composition boot
	// applies the catalog and seeds admin through the production path.
	hash, err := auth.HashPassword(recoveryPass, 4)
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	cfgLive := &config.Config{
		AppName:      "r16-pg-live",
		AppEnv:       "development",
		HTTPAddr:     "127.0.0.1:0",
		DBDialect:    "postgres",
		DBPath:       filepath.Join(t.TempDir(), "unused.sqlite"),
		DBDSN:        liveDSN,
		ProfileName:  "mvp",
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		IdleTimeout:  time.Second,
	}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	appLive, err := NewApp(cfgLive, recoveryK1, hash, logger)
	if err != nil {
		t.Fatalf("build live composition: %v", err)
	}
	startCtx, startCancel := context.WithTimeout(ctx, 60*time.Second)
	defer startCancel()
	if err := appLive.Start(startCtx); err != nil {
		t.Fatalf("T0 live composition (K1) failed to start: %v", err)
	}
	stopCtx, stopCancel := context.WithTimeout(ctx, 15*time.Second)
	defer stopCancel()
	if err := appLive.Stop(stopCtx); err != nil {
		t.Fatalf("stop live composition: %v", err)
	}

	// T0 · a real session on the live K1-era database.
	st0, err := store.Open(ctx, store.OpenOptions{Dialect: kernel.DialectPostgres, DSN: liveDSN, ConnectTimeout: 15 * time.Second}, nil)
	if err != nil {
		t.Fatalf("reopen live store: %v", err)
	}
	// W16-F01 mirror of testsupport.OpenStore: the production bootstrap seeds
	// must_change_password=1; clear it so the seeded admin can authenticate
	// against protected endpoints in this evidence loop.
	clearConn, err := sql.Open("pgx", liveDSN)
	if err != nil {
		t.Fatalf("open clear handle: %v", err)
	}
	if _, err := clearConn.ExecContext(ctx, `UPDATE users SET must_change_password = 0`); err != nil {
		clearConn.Close()
		t.Fatalf("clear must_change_password: %v", err)
	}
	if err := clearConn.Close(); err != nil {
		t.Fatalf("close clear handle: %v", err)
	}
	accessOld, refreshOld := recoveryLogin(t, authsession.NewRepository(st0), []byte(recoveryK1))
	if err := st0.Close(); err != nil {
		t.Fatalf("close live store: %v", err)
	}

	// T1/T3 · the documented PG backup contract across the rotation boundary.
	dump := tooling.dump(t, liveDB)
	if _, err := admin.ExecContext(ctx, `CREATE DATABASE `+restoredDB); err != nil {
		t.Fatalf("create restored scratch db: %v", err)
	}
	tooling.restore(t, restoredDB, dump)
	assertRestoredLedger(t, adminDSN, liveDB, restoredDB)

	// T4 · the FULL application boots from the restored database under
	// current=K2, previous=K1.
	cfgRestored := &config.Config{
		AppName:               "r16-pg-restored",
		AppEnv:                "development",
		HTTPAddr:              "127.0.0.1:0",
		DBDialect:             "postgres",
		DBPath:                filepath.Join(t.TempDir(), "unused.sqlite"),
		DBDSN:                 swapPGDatabase(adminDSN, restoredDB),
		ProfileName:           "mvp",
		ReadTimeout:           time.Second,
		WriteTimeout:          time.Second,
		IdleTimeout:           time.Second,
		AuthJWTSecretPrevious: recoveryK1,
	}
	appRestored, err := NewApp(cfgRestored, recoveryK2, hash, logger)
	if err != nil {
		t.Fatalf("build restored composition: %v", err)
	}
	startCtx2, startCancel2 := context.WithTimeout(ctx, 60*time.Second)
	defer startCancel2()
	if err := appRestored.Start(startCtx2); err != nil {
		t.Fatalf("T4 restored-backup composition failed to start under rotated keys: %v", err)
	}
	stopCtx2, stopCancel2 := context.WithTimeout(ctx, 15*time.Second)
	defer stopCancel2()
	if err := appRestored.Stop(stopCtx2); err != nil {
		t.Fatalf("stop restored composition: %v", err)
	}

	// T5 · rotated auth assertions on the restored database.
	st1, err := store.Open(ctx, store.OpenOptions{Dialect: kernel.DialectPostgres, DSN: swapPGDatabase(adminDSN, restoredDB), ConnectTimeout: 15 * time.Second}, nil)
	if err != nil {
		t.Fatalf("open restored store: %v", err)
	}
	defer st1.Close()
	recoveryAssertRotated(t, authsession.NewRepository(st1), accessOld, refreshOld)
}

// pgURL parses an admin DSN for host/port/user/password extraction.
func pgURL(dsn string) (*url.URL, error) {
	return url.Parse(dsn)
}

// swapPGDatabase returns the DSN with its path (database name) replaced.
func swapPGDatabase(dsn, db string) string {
	u, err := url.Parse(dsn)
	if err != nil {
		return dsn
	}
	u.Path = "/" + db
	return u.String()
}
