package backup

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// workspace-040 R3-C driver (GOAL-007): measures the pg_dump/pg_restore
// cross-version matrix frozen in GOAL-007 D-001.
//
// The matrix needs docker, an external container network and three fixed-version
// images, so it is opt-in: the repository's ordinary `go test ./...` must not
// depend on docker.
//
//	VP040_PG_MATRIX=1 go test -count=1 -run TestPGCrossVersionRestoreMatrix -v ./internal/backup/
//
// Set VP040_PG_MATRIX_OUT=<path> to also write the Markdown report to a file.
//
// Nothing is assumed: every cell records the measured version strings, the exit
// code and the first diagnostic line, and a failure outside the frozen
// classification fails the whole run (D-001 §3: fail closed, never silently
// reclassified as "unsupported").

const (
	matrixGate    = "VP040_PG_MATRIX"
	matrixOutEnv  = "VP040_PG_MATRIX_OUT"
	matrixNetwork = "vp040mx-net"
	matrixPass    = "vp040matrix"
	matrixMount   = "/vp040"
	matrixDB      = "vp040mx_src"
)

// matrixNode is one leg of the server or client axis (D-001 §1).
type matrixNode struct {
	major string
	image string
}

var matrixNodes = []matrixNode{
	{major: "15", image: "postgres:15-alpine"},
	{major: "16", image: "postgres:16"},
	{major: "17", image: "postgres:17-alpine"},
}

// Canonical sample instants seeded into every source database. The trailing-zero
// microsecond value is the one a width-losing path would render as ".9".
const (
	matrixUserCreated = "2026-09-20T12:57:15.900000Z"
	matrixUserLocked  = "2026-01-02T03:04:05.123456Z"
	matrixJobCreated  = "2026-09-20T12:57:15.914000Z"
)

// matrixServer binds one running server container to the DSNs the two callers
// need: the host (migrations, fact reads) and the client containers (pg_dump,
// pg_restore).
type matrixServer struct {
	node         matrixNode
	name         string
	hostDSN      string // 127.0.0.1:<published port>
	containerDSN string // <alias>:5432
}

func (s matrixServer) hostDSNFor(db string) string {
	return withDatabaseOrDie(s.hostDSN, db)
}

func (s matrixServer) containerDSNFor(db string) string {
	return withDatabaseOrDie(s.containerDSN, db)
}

func withDatabaseOrDie(dsn, db string) string {
	out, err := withDatabase(dsn, db)
	if err != nil {
		panic(fmt.Sprintf("withDatabase(%q, %q): %v", redactDSN(dsn), db, err))
	}
	return out
}

// matrixDocker runs one docker command and returns its combined output and exit
// code; a launch failure (docker itself missing) is fatal.
func matrixDocker(t *testing.T, args ...string) (string, int) {
	t.Helper()
	out, err := exec.Command("docker", args...).CombinedOutput()
	text := strings.TrimSpace(string(out))
	if err == nil {
		return text, 0
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return text, exitErr.ExitCode()
	}
	t.Fatalf("docker %s: %v: %s", strings.Join(args, " "), err, text)
	return "", -1
}

// matrixFirstDiagnostic picks the first line that explains a failure.
func matrixFirstDiagnostic(out string) string {
	for _, line := range strings.Split(out, "\n") {
		trimmed := strings.TrimSpace(line)
		lower := strings.ToLower(trimmed)
		if strings.Contains(lower, "error") || strings.Contains(lower, "detail") ||
			strings.Contains(lower, "warning") {
			return trimmed
		}
	}
	return strings.Join(strings.Fields(out), " ")
}

// matrixStartServer starts one fixed-version server that both the host (via a
// published port) and the client containers (via the shared network) can reach.
func matrixStartServer(t *testing.T, node matrixNode) matrixServer {
	t.Helper()
	name := "vp040mx-s" + node.major
	matrixDocker(t, "rm", "-f", name)
	if out, code := matrixDocker(t, "run", "-d", "--name", name,
		"--network", matrixNetwork, "--network-alias", name,
		"-p", "127.0.0.1:0:5432",
		"-e", "POSTGRES_PASSWORD="+matrixPass,
		node.image); code != 0 {
		t.Fatalf("start server %s: %s", node.image, out)
	}
	t.Cleanup(func() { matrixDocker(t, "rm", "-f", name) })

	deadline := time.Now().Add(90 * time.Second)
	for {
		out, code := matrixDocker(t, "run", "--rm", "--network", matrixNetwork,
			matrixNodes[len(matrixNodes)-1].image, "pg_isready", "-h", name, "-U", "postgres")
		if code == 0 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("server %s never became ready: %s", name, out)
		}
		time.Sleep(time.Second)
	}

	out, code := matrixDocker(t, "port", name, "5432")
	if code != 0 {
		t.Fatalf("docker port %s: %s", name, out)
	}
	hostPort := strings.TrimSpace(strings.Split(strings.Split(out, "\n")[0], ":")[1])
	return matrixServer{
		node:         node,
		name:         name,
		hostDSN:      fmt.Sprintf("postgres://postgres:%s@127.0.0.1:%s/postgres?sslmode=disable", matrixPass, hostPort),
		containerDSN: fmt.Sprintf("postgres://postgres:%s@%s:5432/postgres?sslmode=disable", matrixPass, name),
	}
}

func matrixVersionOf(t *testing.T, node matrixNode, tool string, args ...string) string {
	t.Helper()
	out, code := matrixDocker(t, append([]string{"run", "--rm", node.image, tool}, args...)...)
	if code != 0 {
		t.Fatalf("%s %s on %s: %s", tool, strings.Join(args, " "), node.image, out)
	}
	return out
}

func matrixToolVersion(t *testing.T, node matrixNode, tool string) string {
	t.Helper()
	return stripVersionPrefix(matrixVersionOf(t, node, tool, "--version"))
}

func stripVersionPrefix(raw string) string {
	// "postgres (PostgreSQL) 15.19" / "pg_dump (PostgreSQL) 16.15 (Debian ...)"
	for _, prefix := range []string{"postgres (PostgreSQL) ", "pg_dump (PostgreSQL) ", "pg_restore (PostgreSQL) "} {
		if strings.HasPrefix(raw, prefix) {
			return strings.TrimSpace(strings.TrimPrefix(raw, prefix))
		}
	}
	return raw
}

// matrixFacts are the source facts a restored copy must reproduce (D-001 §4).
type matrixFacts struct {
	ledgerHead  int64
	ledgerRows  int64
	userCreated string
	userLocked  string
	jobCreated  string
}

// matrixBuildSource applies the real VP-040 migration chain to a fresh database
// and seeds canonical rows. Using the real catalog (not a trivial schema) makes
// the same run answer whether the migration chain works on PG 16/17 (D-001 §5).
func matrixBuildSource(t *testing.T, sourceDSN string) matrixFacts {
	t.Helper()
	ctx := context.Background()
	catalog := migrationCatalog(t)
	st, err := store.Open(ctx, store.OpenOptions{
		Dialect: kernel.DialectPostgres, DSN: sourceDSN, ConnectTimeout: 30 * time.Second,
	}, catalog)
	if err != nil {
		t.Fatalf("apply VP-040 migrations to %s: %v", redactDSN(sourceDSN), err)
	}
	if err := st.Run(ctx, func(tx kernel.Tx) error {
		if _, err := tx.Exec(ctx, `INSERT INTO users (id, username, name, roles, password_hash, created_at, updated_at,
			token_version, failed_login_count, locked_until, enabled, notifications_enabled, avatar_url,
			must_change_password, email, email_status, last_login_failure_at)
			VALUES ('mx-u1','mxu1','Matrix One','[]','h', TIMESTAMPTZ '`+matrixUserCreated+`',
			        TIMESTAMPTZ '`+matrixUserCreated+`', 0, 0, NULL, 1, 1, '', 0, NULL, NULL, NULL)`); err != nil {
			return fmt.Errorf("seed mx-u1: %w", err)
		}
		if _, err := tx.Exec(ctx, `INSERT INTO users (id, username, name, roles, password_hash, created_at, updated_at,
			token_version, failed_login_count, locked_until, enabled, notifications_enabled, avatar_url,
			must_change_password, email, email_status, last_login_failure_at)
			VALUES ('mx-u2','mxu2','Matrix Two','[]','h', TIMESTAMPTZ '`+matrixUserLocked+`',
			        TIMESTAMPTZ '`+matrixUserLocked+`', 0, 0, TIMESTAMPTZ '`+matrixUserLocked+`', 1, 1, '', 0, NULL, NULL, NULL)`); err != nil {
			return fmt.Errorf("seed mx-u2: %w", err)
		}
		if _, err := tx.Exec(ctx, `INSERT INTO jobs (id, kind, status, payload, progress, cancel_requested, attempt,
			max_attempts, lease_owner, lease_version, lease_expires_at, result, error_code, error_message,
			actor_id, correlation_id, created_at, updated_at, finished_at, expires_at)
			VALUES ('mx-j1','matrix','queued','{}',0,0,0,3,NULL,0,NULL,NULL,NULL,NULL,'a','c',
			        TIMESTAMPTZ '`+matrixJobCreated+`', TIMESTAMPTZ '`+matrixJobCreated+`', NULL, NULL)`); err != nil {
			return fmt.Errorf("seed mx-j1: %w", err)
		}
		return nil
	}); err != nil {
		_ = st.Close()
		t.Fatalf("seed source rows: %v", err)
	}
	if err := st.Close(); err != nil {
		t.Fatalf("close source store: %v", err)
	}

	db := matrixOpen(t, sourceDSN)
	defer db.Close()
	return matrixReadFacts(t, db, true)
}

// canonicalExpr renders a timestamptz column in the project's canonical wire form
// so the comparison is byte-level (microsecond width included).
func canonicalExpr(column string) string {
	return `to_char(` + column + ` AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS.US"Z"')`
}

func matrixReadFacts(t *testing.T, db *sql.DB, withLedger bool) matrixFacts {
	t.Helper()
	var facts matrixFacts
	if withLedger {
		if err := db.QueryRow(`SELECT max(version), count(*) FROM schema_migrations`).
			Scan(&facts.ledgerHead, &facts.ledgerRows); err != nil {
			t.Fatalf("read ledger facts: %v", err)
		}
	}
	if err := db.QueryRow(`SELECT ` + canonicalExpr("created_at") + ` FROM users WHERE id = 'mx-u1'`).
		Scan(&facts.userCreated); err != nil {
		t.Fatalf("read users.created_at: %v", err)
	}
	if err := db.QueryRow(`SELECT ` + canonicalExpr("locked_until") + ` FROM users WHERE id = 'mx-u2'`).
		Scan(&facts.userLocked); err != nil {
		t.Fatalf("read users.locked_until: %v", err)
	}
	if err := db.QueryRow(`SELECT ` + canonicalExpr("created_at") + ` FROM jobs WHERE id = 'mx-j1'`).
		Scan(&facts.jobCreated); err != nil {
		t.Fatalf("read jobs.created_at: %v", err)
	}
	return facts
}

// matrixVerifyShape re-measures the four checks D-001 §4 requires.
func matrixVerifyShape(t *testing.T, dsn string, want matrixFacts) []string {
	t.Helper()
	db := matrixOpen(t, dsn)
	defer db.Close()
	got := matrixReadFacts(t, db, true)

	problems := make([]string, 0, 4)
	if got.ledgerHead != want.ledgerHead {
		problems = append(problems, fmt.Sprintf("ledger head %d, want %d", got.ledgerHead, want.ledgerHead))
	}
	if got.ledgerRows != want.ledgerRows {
		problems = append(problems, fmt.Sprintf("ledger rows %d, want %d", got.ledgerRows, want.ledgerRows))
	}
	if got.userCreated != want.userCreated {
		problems = append(problems, fmt.Sprintf("users.created_at %q, want %q", got.userCreated, want.userCreated))
	}
	if got.userLocked != want.userLocked {
		problems = append(problems, fmt.Sprintf("users.locked_until %q, want %q", got.userLocked, want.userLocked))
	}
	if got.jobCreated != want.jobCreated {
		problems = append(problems, fmt.Sprintf("jobs.created_at %q, want %q", got.jobCreated, want.jobCreated))
	}

	// Conversion shape: the sampled columns must still be timestamptz(6).
	for _, probe := range []struct{ table, column string }{
		{"users", "created_at"}, {"jobs", "created_at"}, {"mail_config", "updated_at"},
	} {
		var dataType string
		var precision sql.NullInt64
		if err := db.QueryRow(`SELECT data_type, datetime_precision FROM information_schema.columns
			WHERE table_name = $1 AND column_name = $2`, probe.table, probe.column).
			Scan(&dataType, &precision); err != nil {
			problems = append(problems, fmt.Sprintf("%s.%s shape unreadable: %v", probe.table, probe.column, err))
			continue
		}
		if dataType != "timestamp with time zone" || !precision.Valid || precision.Int64 != 6 {
			problems = append(problems, fmt.Sprintf("%s.%s is %s(%v), want timestamp with time zone with 6 digits",
				probe.table, probe.column, dataType, precision.Int64))
		}
	}

	// D0 sentinel column: the seeded NULL must stay NULL.
	var nulls int
	if err := db.QueryRow(`SELECT count(*) FROM users WHERE id = 'mx-u1' AND locked_until IS NULL`).Scan(&nulls); err != nil {
		problems = append(problems, fmt.Sprintf("sentinel check unreadable: %v", err))
	} else if nulls != 1 {
		problems = append(problems, "users.locked_until lost its NULL for mx-u1 (legacy sentinel revived)")
	}
	return problems
}

func matrixOpen(t *testing.T, dsn string) *sql.DB {
	t.Helper()
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("open %s: %v", redactDSN(dsn), err)
	}
	db.SetMaxOpenConns(2)
	return db
}

func matrixCreateDB(t *testing.T, adminDSN, name string) {
	t.Helper()
	db := matrixOpen(t, adminDSN)
	defer db.Close()
	if _, err := db.Exec(`DROP DATABASE IF EXISTS ` + quoteIdent(name) + ` WITH (FORCE)`); err != nil {
		t.Fatalf("drop %s: %v", name, err)
	}
	if _, err := db.Exec(`CREATE DATABASE ` + quoteIdent(name)); err != nil {
		t.Fatalf("create %s: %v", name, err)
	}
}

func matrixDump(t *testing.T, client matrixNode, workDir, archive, containerDSN string) (int, string) {
	t.Helper()
	out, code := matrixDocker(t, "run", "--rm", "--network", matrixNetwork,
		"-v", workDir+":"+matrixMount, client.image,
		"pg_dump", "-F", "c", "--no-owner", "-f", matrixMount+"/"+archive, containerDSN)
	return code, out
}

func matrixRestore(t *testing.T, client matrixNode, workDir, archive, targetDSN string) (int, string) {
	t.Helper()
	out, code := matrixDocker(t, "run", "--rm", "--network", matrixNetwork,
		"-v", workDir+":"+matrixMount, client.image,
		"pg_restore", "--exit-on-error", "--no-owner", "-d", targetDSN, matrixMount+"/"+archive)
	return code, out
}

// matrixClassify applies the frozen three-class criterion (D-001 §3).
func matrixClassify(code int, out string) string {
	if code == 0 {
		return "supported"
	}
	if strings.Contains(out, "server version mismatch") {
		return "unsupported-toolgate"
	}
	if strings.Contains(out, "unsupported version") && strings.Contains(out, "in file header") {
		return "unsupported-toolgate"
	}
	if strings.Contains(out, `unrecognized configuration parameter "transaction_timeout"`) {
		return "unsupported-serverguc"
	}
	if strings.Contains(out, "does not exist") {
		return "setup-failure"
	}
	return "unexpected-failure"
}

type matrixDumpCell struct {
	server, client, class, diagnostic string
	exit                              int
}

type matrixRestoreCell struct {
	archive, dumper, client, target, class, shape, diagnostic string
	exit                                                      int
}

// TestPGCrossVersionRestoreMatrix is the R3-C matrix driver.
func TestPGCrossVersionRestoreMatrix(t *testing.T) {
	if os.Getenv(matrixGate) != "1" {
		t.Skipf("%s is not set; the cross-version container matrix is opt-in (needs docker + network)", matrixGate)
	}
	if _, err := exec.LookPath("docker"); err != nil {
		t.Skipf("docker unavailable: %v", err)
	}
	for _, node := range matrixNodes {
		if out, code := matrixDocker(t, "image", "inspect", node.image); code != 0 {
			t.Skipf("image %s is not present locally (%s); refusing to pull implicitly", node.image, matrixFirstDiagnostic(out))
		}
	}
	matrixDocker(t, "network", "rm", matrixNetwork)
	if out, code := matrixDocker(t, "network", "create", matrixNetwork); code != 0 {
		t.Fatalf("create network: %s", out)
	}
	t.Cleanup(func() { matrixDocker(t, "network", "rm", matrixNetwork) })

	workDir := t.TempDir()

	// 1. Measured versions (D-001 §1: never trust the tag).
	serverVersions := map[string]string{}
	clientVersions := map[string]string{}
	for _, node := range matrixNodes {
		serverVersions[node.major] = stripVersionPrefix(matrixVersionOf(t, node, "postgres", "--version"))
		clientVersions[node.major] = matrixToolVersion(t, node, "pg_dump") + " / " + matrixToolVersion(t, node, "pg_restore")
	}

	// 2. One converted source per server, dumped by every client.
	servers := map[string]matrixServer{}
	sources := map[string]matrixFacts{}
	dumps := []matrixDumpCell{}
	archives := []struct {
		name   string
		dumper string
		client string
		path   string
	}{}
	for _, node := range matrixNodes {
		srv := matrixStartServer(t, node)
		servers[node.major] = srv
		matrixCreateDB(t, srv.hostDSN, matrixDB)
		sourceDSN := srv.hostDSNFor(matrixDB)
		facts := matrixBuildSource(t, sourceDSN)
		sources[node.major] = facts
		t.Logf("server %s (%s): migrations at head %d, %d ledger rows, samples %s / %s / %s",
			node.major, serverVersions[node.major], facts.ledgerHead, facts.ledgerRows,
			facts.userCreated, facts.userLocked, facts.jobCreated)
		for _, client := range matrixNodes {
			archive := fmt.Sprintf("src%s_by_%s.dump", node.major, client.major)
			code, out := matrixDump(t, client, workDir, archive, srv.containerDSNFor(matrixDB))
			class := matrixClassify(code, out)
			diag := ""
			if code != 0 {
				diag = matrixFirstDiagnostic(out)
				if class != "unsupported-toolgate" {
					t.Fatalf("unexpected pg_dump failure (server %s, client %s): exit=%d class=%s out=%s",
						node.major, client.major, code, class, out)
				}
			} else if _, err := os.Stat(filepath.Join(workDir, archive)); err != nil {
				t.Fatalf("pg_dump reported success but %s is missing: %v", archive, err)
			}
			dumps = append(dumps, matrixDumpCell{
				server: node.major, client: client.major, class: class, exit: code, diagnostic: diag,
			})
			if code == 0 {
				archives = append(archives, struct {
					name   string
					dumper string
					client string
					path   string
				}{name: archive, dumper: node.major, client: client.major, path: filepath.Join(workDir, archive)})
			}
		}
	}

	// 3. Archive format versions, read from the custom-archive header.
	headers := map[string]string{}
	for _, a := range archives {
		raw, err := os.ReadFile(a.path)
		if err != nil || len(raw) < 7 {
			t.Fatalf("read archive %s: %v", a.name, err)
		}
		bytes := make([]string, 0, 7)
		for _, b := range raw[:7] {
			bytes = append(bytes, fmt.Sprintf("%d", b))
		}
		headers[a.name] = fmt.Sprintf("%s (1.%d)", strings.Join(bytes, " "), raw[6])
	}

	// 4. Every archive x every client x every target server.
	restores := []matrixRestoreCell{}
	shapeChecks := 0
	for _, a := range archives {
		for _, client := range matrixNodes {
			for _, target := range matrixNodes {
				srv := servers[target.major]
				db := fmt.Sprintf("mx_r_%s_%s_%s", a.dumper, client.major, target.major)
				matrixCreateDB(t, srv.hostDSN, db)
				code, out := matrixRestore(t, client, workDir, a.name, srv.containerDSNFor(db))
				class := matrixClassify(code, out)
				shape := "-"
				switch class {
				case "supported":
					problems := matrixVerifyShape(t, srv.hostDSNFor(db), sources[a.dumper])
					shapeChecks++
					if len(problems) > 0 {
						t.Fatalf("restored shape mismatch (%s -> client %s -> server %s): %s",
							a.name, client.major, target.major, strings.Join(problems, "; "))
					}
					shape = "ok"
					// Shape verification runs on a live connection; drop the copy.
					matrixDropDB(t, srv.hostDSN, db)
				case "unexpected-failure", "setup-failure":
					t.Fatalf("unexpected pg_restore failure (%s -> client %s -> server %s): exit=%d class=%s out=%s",
						a.name, client.major, target.major, code, class, out)
				}
				diag := ""
				if code != 0 {
					diag = matrixFirstDiagnostic(out)
				}
				restores = append(restores, matrixRestoreCell{
					archive: a.name, dumper: a.dumper, client: client.major, target: target.major,
					class: class, shape: shape, exit: code, diagnostic: diag,
				})
			}
		}
	}
	if shapeChecks == 0 {
		t.Fatal("no restore combination reached the shape verification; the matrix proves nothing")
	}

	// 5. Report.
	var report strings.Builder
	fmt.Fprintf(&report, "### Measured versions\n\n| server | image | server version | client tools |\n|---|---|---|---|\n")
	for _, node := range matrixNodes {
		fmt.Fprintf(&report, "| %s | `%s` | %s | %s |\n", node.major, node.image,
			serverVersions[node.major], clientVersions[node.major])
	}

	fmt.Fprintf(&report, "\n### pg_dump: server x client (9 cells)\n\n| server | client | exit | class | diagnostic |\n|---|---|---|---|---|\n")
	for _, cell := range dumps {
		fmt.Fprintf(&report, "| %s | %s | %d | %s | %s |\n",
			cell.server, cell.client, cell.exit, cell.class, cell.diagnostic)
	}

	fmt.Fprintf(&report, "\n### Archive format version (PGDMP header, decimal)\n\n| archive | dumper | header |\n|---|---|---|\n")
	for _, a := range archives {
		fmt.Fprintf(&report, "| %s | %s | %s |\n", a.name, a.dumper, headers[a.name])
	}

	fmt.Fprintf(&report, "\n### pg_restore: archive x client x target server (%d cells)\n\n", len(restores))
	fmt.Fprintf(&report, "| archive | dumper | client | target | exit | class | shape | diagnostic |\n|---|---|---|---|---|---|---|---|\n")
	for _, cell := range restores {
		fmt.Fprintf(&report, "| %s | %s | %s | %s | %d | %s | %s | %s |\n",
			cell.archive, cell.dumper, cell.client, cell.target, cell.exit, cell.class, cell.shape, cell.diagnostic)
	}

	var supported, toolgate, serverguc int
	for _, cell := range restores {
		switch cell.class {
		case "supported":
			supported++
		case "unsupported-toolgate":
			toolgate++
		case "unsupported-serverguc":
			serverguc++
		}
	}
	fmt.Fprintf(&report, "\n### Tally\n\n")
	fmt.Fprintf(&report, "- dump cells measured: %d (supported %d)\n", len(dumps), len(archives))
	fmt.Fprintf(&report, "- restore cells measured: %d\n", len(restores))
	fmt.Fprintf(&report, "- restore supported (shape-verified): %d\n", supported)
	fmt.Fprintf(&report, "- restore unsupported-toolgate: %d\n", toolgate)
	fmt.Fprintf(&report, "- restore unsupported-serverguc: %d\n", serverguc)
	fmt.Fprintf(&report, "- shape checks executed: %d\n", shapeChecks)

	out := report.String()
	fmt.Println(out)
	if path := os.Getenv(matrixOutEnv); path != "" {
		if err := os.WriteFile(path, []byte(out), 0o644); err != nil {
			t.Fatalf("write report: %v", err)
		}
	}
}

func matrixDropDB(t *testing.T, adminDSN, name string) {
	t.Helper()
	db := matrixOpen(t, adminDSN)
	defer db.Close()
	if _, err := db.Exec(`DROP DATABASE IF EXISTS ` + quoteIdent(name) + ` WITH (FORCE)`); err != nil {
		t.Logf("drop %s (non-fatal): %v", name, err)
	}
}
