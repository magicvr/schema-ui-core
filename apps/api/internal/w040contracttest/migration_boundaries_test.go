// Package w040contracttest holds the executable boundary validation for the
// VP-040 time-column conversion contract.
//
// The transcribed freeze-text cases live in contract_boundaries_test.go. This
// file is their R2 redirection (user ruling 2026-09-20, GOAL-002 D-020 §2): the
// same boundary values are now driven through the REAL v73–v87 conversion
// migrations — a v1–v72 database, legacy rows seeded through the pre-conversion
// shape, then the compiled catalog applied to head — so the assertions exercise
// the shipped descriptors instead of a hand-written copy of the frozen DDL.
package w040contracttest

import (
	"context"
	"database/sql"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
	"github.com/magicvr/schema-ui-core/apps/api/internal/temporal"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	compiledmodules "github.com/magicvr/schema-ui-core/apps/api/modules/compiled"
)

// v72Head is the last pre-conversion catalog version; the workspace-040
// conversion descriptors are v73–v87.
const v72Head = 72

func fullCatalog(t *testing.T) []kernel.MigrationContribution {
	t.Helper()
	full, err := compiledmodules.PersistenceCatalog()
	if err != nil {
		t.Fatalf("compiled catalog: %v", err)
	}
	if len(full) != 87 {
		t.Fatalf("compiled catalog has %d entries, want 87 (v1–v87)", len(full))
	}
	return full
}

// atV72 materialises a legacy database: the frozen v1–v72 history applied
// through the real runner, which is the shape every conversion descriptor is
// written against.
func atV72(t *testing.T) (path string, full []kernel.MigrationContribution) {
	t.Helper()
	full = fullCatalog(t)
	path = filepath.Join(t.TempDir(), "v72.db")
	st, err := store.OpenWithCatalog(path, full[:v72Head])
	if err != nil {
		t.Fatalf("open v72: %v", err)
	}
	if err := st.Close(); err != nil {
		t.Fatalf("close v72: %v", err)
	}
	return path, full
}

// raw opens the legacy file without running migrations so the pre-conversion
// rows can be seeded through the integer shape.
func raw(t *testing.T, path string) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("raw open: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

// toHead applies the real conversion migrations and returns the upgraded store.
func toHead(t *testing.T, path string, full []kernel.MigrationContribution) *store.Store {
	t.Helper()
	st, err := store.OpenWithCatalog(path, full)
	if err != nil {
		t.Fatalf("upgrade to head: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return st
}

// scanCanonical reads one converted column through the store read adapter and
// returns the canonical text (or "" for SQL NULL).
func scanCanonical(t *testing.T, st *store.Store, query string, args ...any) string {
	t.Helper()
	ctx := context.Background()
	var out string
	if err := st.Run(ctx, func(tx kernel.Tx) error {
		var instant sql.NullTime
		if err := tx.QueryRow(ctx, query, args...).Scan(&instant); err != nil {
			return err
		}
		if !instant.Valid {
			return nil
		}
		out = instant.Time.UTC().Format("2006-01-02T15:04:05.000000Z")
		return nil
	}); err != nil {
		t.Fatalf("scan %q: %v", query, err)
	}
	return out
}

// scanCanonicalRequired is scanCanonical for NOT NULL columns.
func scanCanonicalRequired(t *testing.T, st *store.Store, query string, args ...any) string {
	t.Helper()
	ctx := context.Background()
	var out string
	if err := st.Run(ctx, func(tx kernel.Tx) error {
		var instant time.Time
		if err := tx.QueryRow(ctx, query, args...).Scan(&instant); err != nil {
			return err
		}
		out = instant.UTC().Format("2006-01-02T15:04:05.000000Z")
		return nil
	}); err != nil {
		t.Fatalf("scan %q: %v", query, err)
	}
	return out
}

// TestRealMigrationsConvertBoundaryInstants is the redirected boundary matrix:
// every case runs against the shipped v73–v87 descriptors on a real v72
// database, with the values seeded through the legacy integer shape.
func TestRealMigrationsConvertBoundaryInstants(t *testing.T) {
	path, full := atV72(t)
	db := raw(t, path)

	// ---- legacy seed (integer seconds / milliseconds) ----
	mustExec(t, db, `INSERT INTO users (id, username, name, roles, password_hash, created_at, updated_at,
		token_version, failed_login_count, locked_until, enabled, notifications_enabled, avatar_url,
		must_change_password, email, email_status, last_login_failure_at)
		VALUES ('u-1','u1','User One','[]','hash', 0, 253402300799, 0, 0, 0, 1, 1, '', 0, NULL, NULL, 0)`)
	mustExec(t, db, `INSERT INTO mail_config (id, updated_at) VALUES (1, 0)`)
	mustExec(t, db, `INSERT INTO telegram_config (id, updated_at) VALUES (1, 0)`)
	mustExec(t, db, `INSERT INTO notifications (id, user_id, event, title, body, read_at, created_at)
		VALUES ('n-1','u-1','account.locked','t','b', NULL, 1)`)

	// jobs covers the millisecond family: no-carry at 999 ms, the negative floor
	// cases from r1-c2-sqlite-rebuild-mechanism-v1.0-fc.md §2.3, and year 9999.
	jobs := []struct {
		id  string
		ms  int64
		exp string
	}{
		{"j-nocarry", 1758320000999, "2025-09-19T22:13:20.999000Z"},
		{"j-neg1", -1, "1969-12-31T23:59:59.999000Z"},
		{"j-neg999", -999, "1969-12-31T23:59:59.001000Z"},
		{"j-neg1000", -1000, "1969-12-31T23:59:59.000000Z"},
		{"j-neg1001", -1001, "1969-12-31T23:59:58.999000Z"},
		{"j-max", 253402300799999, "9999-12-31T23:59:59.999000Z"},
		{"j-epoch", 0, "1970-01-01T00:00:00.000000Z"},
	}
	for _, job := range jobs {
		mustExec(t, db, `INSERT INTO jobs (id, kind, status, payload, progress, cancel_requested, attempt,
			max_attempts, lease_owner, lease_version, lease_expires_at, result, error_code, error_message,
			actor_id, correlation_id, created_at, updated_at, finished_at, expires_at)
			VALUES (?, 'probe', 'queued', '{}', 0, 0, 0, 3, NULL, 0, NULL, NULL, NULL, NULL, 'a', 'c', ?, ?, NULL, NULL)`,
			job.id, job.ms, job.ms)
	}

	// vouchers: legacy 0 and NULL both mean "absent"; a positive value converts.
	mustExec(t, db, `INSERT INTO vouchers (id, batch_id, code_hash, code_prefix, amount, currency, status,
		expires_at, redeemed_by, redeemed_at, created_at, updated_at)
		VALUES ('v-zero','b','h1','p',1,'CNY','unused',0,NULL,NULL,1,1)`)
	mustExec(t, db, `INSERT INTO vouchers (id, batch_id, code_hash, code_prefix, amount, currency, status,
		expires_at, redeemed_by, redeemed_at, created_at, updated_at)
		VALUES ('v-null','b','h2','p',1,'CNY','unused',NULL,NULL,NULL,1,1)`)
	mustExec(t, db, `INSERT INTO vouchers (id, batch_id, code_hash, code_prefix, amount, currency, status,
		expires_at, redeemed_by, redeemed_at, created_at, updated_at)
		VALUES ('v-pos','b','h3','p',1,'CNY','unused',1758320000,NULL,NULL,1758320000,1758320000)`)

	// task_runs.finished_at: NULL stays NULL (the legacy "0 means unfinished"
	// spelling is a runtime concern, not a stored one).
	mustExec(t, db, `INSERT INTO scheduled_tasks (id, key, cron, name, enabled, description, handler, created_at, updated_at)
		VALUES ('t-1','k','* * * * *','T',1,NULL,'system.noop',1,1)`)
	mustExec(t, db, `INSERT INTO task_runs (id, task_id, status, started_at, finished_at, detail, created_at)
		VALUES ('r-1','t-1','failed',1,NULL,NULL,1)`)

	// login_failures.locked_until (#20) is a D0 sentinel column too: legacy 0
	// must become SQL NULL through the real v74 conversion.
	mustExec(t, db, `INSERT INTO login_failures (user_id, ip, fail_count, locked_until, updated_at)
		VALUES ('u-1','10.0.0.1', 0, 0, 1)`)

	if err := db.Close(); err != nil {
		t.Fatalf("close seed db: %v", err)
	}

	st := toHead(t, path, full)

	// sentinel 0 → NULL (#5, #6, #34, #78)
	for _, tc := range []struct{ name, query string }{
		{"users.locked_until", `SELECT locked_until FROM users WHERE id = 'u-1'`},
		{"users.last_login_failure_at", `SELECT last_login_failure_at FROM users WHERE id = 'u-1'`},
		{"mail_config.updated_at", `SELECT updated_at FROM mail_config WHERE id = 1`},
		{"telegram_config.updated_at", `SELECT updated_at FROM telegram_config WHERE id = 1`},
		{"notifications.read_at", `SELECT read_at FROM notifications WHERE id = 'n-1'`},
		{"task_runs.finished_at", `SELECT finished_at FROM task_runs WHERE id = 'r-1'`},
		{"vouchers.expires_at (0)", `SELECT expires_at FROM vouchers WHERE id = 'v-zero'`},
		{"vouchers.expires_at (NULL)", `SELECT expires_at FROM vouchers WHERE id = 'v-null'`},
		{"login_failures.locked_until (0)", `SELECT locked_until FROM login_failures WHERE user_id = 'u-1' AND ip = '10.0.0.1'`},
	} {
		if got := scanCanonical(t, st, tc.query); got != "" {
			t.Errorf("%s = %q, want SQL NULL", tc.name, got)
		}
	}

	// ordinary instants, including year 0 and year 9999 seconds
	for _, tc := range []struct{ name, query, want string }{
		{"users.created_at (epoch second)", `SELECT created_at FROM users WHERE id = 'u-1'`, "1970-01-01T00:00:00.000000Z"},
		{"users.updated_at (year 9999)", `SELECT updated_at FROM users WHERE id = 'u-1'`, "9999-12-31T23:59:59.000000Z"},
		{"vouchers.expires_at (positive)", `SELECT expires_at FROM vouchers WHERE id = 'v-pos'`, "2025-09-19T22:13:20.000000Z"},
		{"notifications.created_at", `SELECT created_at FROM notifications WHERE id = 'n-1'`, "1970-01-01T00:00:01.000000Z"},
		{"task_runs.started_at", `SELECT started_at FROM task_runs WHERE id = 'r-1'`, "1970-01-01T00:00:01.000000Z"},
	} {
		if got := scanCanonical(t, st, tc.query); got != tc.want {
			t.Errorf("%s = %q, want %q", tc.name, got, tc.want)
		}
	}

	// millisecond family: floor semantics and no carry at 999 ms
	for _, job := range jobs {
		got := scanCanonicalRequired(t, st, `SELECT created_at FROM jobs WHERE id = ?`, job.id)
		if got != job.exp {
			t.Errorf("jobs.created_at(%s) = %q, want %q", job.id, got, job.exp)
		}
		if len(got) != 27 {
			t.Errorf("jobs.created_at(%s) = %q is not the 27-character canonical form", job.id, got)
		}
	}

	// fixed-6 text order is instant order: the lexicographic order of the
	// converted millisecond column must match the epoch order of the inputs.
	ordered := append([]struct {
		id  string
		ms  int64
		exp string
	}{}, jobs...)
	for i := 1; i < len(ordered); i++ {
		for j := i; j > 0 && ordered[j].ms < ordered[j-1].ms; j-- {
			ordered[j], ordered[j-1] = ordered[j-1], ordered[j]
		}
	}
	var lexical []string
	ctx := context.Background()
	if err := st.Run(ctx, func(tx kernel.Tx) error {
		rows, err := tx.Query(ctx, `SELECT id FROM jobs ORDER BY created_at, id`)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var id string
			if err := rows.Scan(&id); err != nil {
				return err
			}
			lexical = append(lexical, id)
		}
		return rows.Err()
	}); err != nil {
		t.Fatalf("lexical order query: %v", err)
	}
	for i, job := range ordered {
		if lexical[i] != job.id {
			t.Fatalf("ORDER BY created_at[%d] = %s, want %s (lexical order must equal instant order)", i, lexical[i], job.id)
		}
	}
}

// TestRealMigrationFailsClosedOnNegativeVoucherInstant is the Root D-012 case on
// the real descriptor: a negative voucher instant is corruption, so v85 must
// refuse and the whole batch must roll back.
func TestRealMigrationFailsClosedOnNegativeVoucherInstant(t *testing.T) {
	path, full := atV72(t)
	db := raw(t, path)
	mustExec(t, db, `INSERT INTO vouchers (id, batch_id, code_hash, code_prefix, amount, currency, status,
		expires_at, redeemed_by, redeemed_at, created_at, updated_at)
		VALUES ('v-neg','b','h','p',1,'CNY','unused',-1,NULL,NULL,1,1)`)
	if err := db.Close(); err != nil {
		t.Fatalf("close seed db: %v", err)
	}

	st, err := store.OpenWithCatalog(path, full)
	if err == nil {
		_ = st.Close()
		t.Fatal("negative voucher instant must fail the conversion closed")
	}
	if !strings.Contains(err.Error(), "negative") {
		t.Fatalf("error = %v, want a negative-value refusal", err)
	}

	// v85's transaction rolled back: the column shape and the offending row are
	// untouched. (Each descriptor runs in its own transaction, so v73–v84 stay
	// applied — which is exactly why the C3 recovery boundary exists.)
	verify := raw(t, path)
	var declared string
	if err := verify.QueryRow(`SELECT type FROM pragma_table_info('vouchers') WHERE name = 'expires_at'`).Scan(&declared); err != nil {
		t.Fatalf("probe vouchers.expires_at: %v", err)
	}
	if !strings.EqualFold(declared, "INTEGER") {
		t.Fatalf("vouchers.expires_at declared type = %q after the refused batch, want INTEGER (v85 must roll back)", declared)
	}
	var negatives int
	if err := verify.QueryRow(`SELECT COUNT(*) FROM vouchers WHERE expires_at = -1`).Scan(&negatives); err != nil {
		t.Fatalf("count negative vouchers: %v", err)
	}
	if negatives != 1 {
		t.Fatalf("negative voucher rows = %d, want the seed row untouched", negatives)
	}
}

// TestRealMigrationKeepsOrdinaryNegativeInstant is the Root D-015 counterpart:
// a negative epoch outside the voucher columns is a legal instant and converts.
func TestRealMigrationKeepsOrdinaryNegativeInstant(t *testing.T) {
	path, full := atV72(t)
	db := raw(t, path)
	mustExec(t, db, `INSERT INTO users (id, username, name, roles, password_hash, created_at, updated_at,
		token_version, failed_login_count, locked_until, enabled, notifications_enabled, avatar_url,
		must_change_password, email, email_status, last_login_failure_at)
		VALUES ('u-neg','un','User Neg','[]','hash', -86400, -1, 0, 0, 0, 1, 1, '', 0, NULL, NULL, 0)`)
	if err := db.Close(); err != nil {
		t.Fatalf("close seed db: %v", err)
	}

	st := toHead(t, path, full)
	if got := scanCanonicalRequired(t, st, `SELECT created_at FROM users WHERE id = 'u-neg'`); got != "1969-12-31T00:00:00.000000Z" {
		t.Errorf("users.created_at = %q, want 1969-12-31T00:00:00.000000Z", got)
	}
	if got := scanCanonicalRequired(t, st, `SELECT updated_at FROM users WHERE id = 'u-neg'`); got != "1969-12-31T23:59:59.000000Z" {
		t.Errorf("users.updated_at = %q, want 1969-12-31T23:59:59.000000Z", got)
	}
}

// TestV73RefusesWhenRetiredRecordsTableIsPresent is the fail-closed half of the
// frozen v73 scope ("断言 retired `records` 不存在"): if an abnormal path leaves
// the retired table behind, the conversion must refuse rather than silently skip
// records.updated_at (GOAL-003 A-002 F-I-001).
func TestV73RefusesWhenRetiredRecordsTableIsPresent(t *testing.T) {
	path, full := atV72(t)
	db := raw(t, path)
	// The retired v6 table, recreated to simulate a skipped/partial restore.
	mustExec(t, db, `CREATE TABLE records (id TEXT PRIMARY KEY, updated_at INTEGER NOT NULL)`)
	mustExec(t, db, `INSERT INTO records (id, updated_at) VALUES ('r-1', 1)`)
	if err := db.Close(); err != nil {
		t.Fatalf("close seed db: %v", err)
	}

	st, err := store.OpenWithCatalog(path, full)
	if err == nil {
		_ = st.Close()
		t.Fatal("a present `records` table must fail the v73 conversion closed")
	}
	if !strings.Contains(err.Error(), "records") {
		t.Fatalf("error = %v, want a retired-table refusal naming records", err)
	}
	// The refused batch leaves v1..v72 untouched: the retired table and its
	// pre-conversion column shape are still there.
	verify := raw(t, path)
	var declared string
	if err := verify.QueryRow(`SELECT type FROM pragma_table_info('records') WHERE name = 'updated_at'`).Scan(&declared); err != nil {
		t.Fatalf("probe records.updated_at: %v", err)
	}
	if !strings.EqualFold(declared, "INTEGER") {
		t.Fatalf("records.updated_at = %q after the refusal, want INTEGER (v73 must roll back)", declared)
	}
	var schemaApplied string
	if err := verify.QueryRow(`SELECT type FROM pragma_table_info('schema_migrations') WHERE name = 'applied_at'`).Scan(&schemaApplied); err != nil {
		t.Fatalf("probe schema_migrations.applied_at: %v", err)
	}
	if !strings.EqualFold(schemaApplied, "INTEGER") {
		t.Fatalf("schema_migrations.applied_at = %q after the refusal, want INTEGER (v73 must roll back)", schemaApplied)
	}
}

func mustExec(t *testing.T, db *sql.DB, query string, args ...any) {
	t.Helper()
	if _, err := db.Exec(query, args...); err != nil {
		t.Fatalf("seed exec: %v\n%s", err, query)
	}
}

// TestCodecMatchesRealMigration is the D-018 T-*-RT cross-check in-process: for
// both unit families the value produced by the real v73–v87 conversion SQL must
// equal temporal.FromUnix / FromUnixMilli formatted by the codec, for the same
// input. It closes the gap between the codec's frozen expectations and the
// migration's actual output (GOAL-003 A-002 F-I-003).
func TestCodecMatchesRealMigration(t *testing.T) {
	path, full := atV72(t)
	db := raw(t, path)

	millis := []int64{
		0, 1, -1, 999, -999, 1000, -1000, 1001, -1001,
		1758320000123, 1758320000999, -1758320000123,
		253402300799999, -62135596800000,
	}
	seconds := []int64{
		0, 1, -1, -86400, 1758320000, 253402300799, -62135596800,
	}
	for index, ms := range millis {
		mustExec(t, db, `INSERT INTO jobs (id, kind, status, payload, progress, cancel_requested, attempt,
			max_attempts, lease_owner, lease_version, lease_expires_at, result, error_code, error_message,
			actor_id, correlation_id, created_at, updated_at, finished_at, expires_at)
			VALUES (?, 'codec', 'queued', '{}', 0, 0, 0, 3, NULL, 0, NULL, NULL, NULL, NULL, 'a', 'c', ?, ?, NULL, NULL)`,
			"codec-ms-"+strconv.Itoa(index), ms, ms)
	}
	for index, sec := range seconds {
		mustExec(t, db, `INSERT INTO users (id, username, name, roles, password_hash, created_at, updated_at,
			token_version, failed_login_count, locked_until, enabled, notifications_enabled, avatar_url,
			must_change_password, email, email_status, last_login_failure_at)
			VALUES (?, ?, 'Codec', '[]', 'hash', ?, ?, 0, 0, 0, 1, 1, '', 0, NULL, NULL, 0)`,
			"codec-sec-"+strconv.Itoa(index), "codec-sec-"+strconv.Itoa(index), sec, sec)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("close seed db: %v", err)
	}

	st := toHead(t, path, full)
	for index, ms := range millis {
		got := scanCanonicalRequired(t, st, `SELECT created_at FROM jobs WHERE id = ?`, "codec-ms-"+strconv.Itoa(index))
		want := temporal.MustFormat(temporal.FromUnixMilli(ms))
		if got != want {
			t.Errorf("milliseconds %d: migration = %q, codec = %q", ms, got, want)
		}
	}
	for index, sec := range seconds {
		id := "codec-sec-" + strconv.Itoa(index)
		got := scanCanonicalRequired(t, st, `SELECT created_at FROM users WHERE id = ?`, id)
		want := temporal.MustFormat(temporal.FromUnix(sec))
		if got != want {
			t.Errorf("seconds %d: migration = %q, codec = %q", sec, got, want)
		}
	}
}
