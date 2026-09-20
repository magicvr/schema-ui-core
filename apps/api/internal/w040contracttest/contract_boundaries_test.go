// Package w040contracttest holds the executable boundary validation for the
// VP-040 R1 time-column conversion contract.
//
// Scope and authority (user ruling 2026-09-20, recorded as D-020):
//
//   - These tests build the legacy shape and the frozen rebuild DDL against a
//     THROWAWAY validation database (:memory:). They do NOT touch the
//     production schema, migrations, repositories or fixtures, so they do not
//     cross the "no schema implementation before C2 freeze" gate.
//   - The literal DDL and the conversion expressions below are transcribed
//     from the R1 freeze candidates:
//     docs/workspaces/workspace-040-timestamptz-persistence-contract/
//     GOAL-002-r1-contract-and-denominator-freeze/attachments/
//     r1-c2-per-table-rebuild-ddl-v1.0-fc.md and
//     r1-c2-sqlite-rebuild-mechanism-v1.0-fc.md.
//   - R2 must redirect these same cases onto the real conversion migrations
//     (D-020). If the frozen DDL changes, these tests must be updated with it.
package w040contracttest

import (
	"database/sql"
	"fmt"
	"testing"

	_ "modernc.org/sqlite"
)

// --- frozen expressions (r1-c2-sqlite-rebuild-mechanism-v1.0-fc.md §2) ---

// secondsExpr converts a legacy integer-seconds column to fixed-6 UTC RFC3339.
func secondsExpr(col string) string {
	return "strftime('%Y-%m-%dT%H:%M:%S', " + col + ", 'unixepoch') || '.000000Z'"
}

// millisExpr converts a legacy integer-milliseconds column to fixed-6 UTC
// RFC3339 using integer arithmetic only.
//
// Both halves are load-bearing and must not be "simplified":
//   - SQLite integer '/' truncates toward zero, so the seconds part needs the
//     (x-999)/1000 floor correction for negative inputs;
//   - SQLite '%' takes the sign of the dividend (-1%1000 = -1), so the
//     millisecond remainder needs (x%1000+1000)%1000 normalisation.
func millisExpr(col string) string {
	sec := "CASE WHEN " + col + " >= 0 THEN " + col + "/1000 ELSE (" + col + "-999)/1000 END"
	rem := "((" + col + "%1000 + 1000) % 1000)"
	return "strftime('%Y-%m-%dT%H:%M:%S', " + sec + ", 'unixepoch') || '.' || printf('%03d', " + rem + ") || '000Z'"
}

// --- legacy shape under test ---

const legacyDDL = `
CREATE TABLE legacy (
  id          TEXT PRIMARY KEY,
  sec_nn      INTEGER NOT NULL,
  ms_nn       INTEGER NOT NULL,
  sec_null    INTEGER,
  ms_null     INTEGER,
  sec_d0      INTEGER NOT NULL DEFAULT 0,
  ms_d0       INTEGER NOT NULL DEFAULT 0,
  voucher_sec INTEGER
)`

// preflight counts the 0 / negative / positive buckets for one column
// (contract m0) and applies the value policy that the accepted decisions
// actually state:
//
//   - Root D-012 (voucher-invalid-value-policy) is scoped to
//     vouchers.expires_at / vouchers.redeemed_at: legacy 0 becomes NULL and a
//     NEGATIVE value is data corruption that must fail closed.
//   - Root D-015 states explicitly that a negative epoch is NOT a sentinel: it
//     is a valid instant and converts normally. So negatives are only an error
//     for the voucher columns, not for every time column.
//   - zeroIsAbsence marks the columns where legacy 0 means "absent" (the D0
//     sentinel columns and the voucher columns); those map 0 to NULL. For any
//     other column a literal 0 is the epoch instant and is converted as-is.
//
// It returns an error rather than calling t.Fatalf so that negative tests can
// assert the fail-closed behaviour.
func preflight(db *sql.DB, table, col string, zeroIsAbsence, negativeIsError bool) (zero, neg, pos int, err error) {
	_ = zeroIsAbsence
	row := db.QueryRow(`SELECT
      COALESCE(SUM(CASE WHEN ` + col + ` = 0 THEN 1 ELSE 0 END), 0),
      COALESCE(SUM(CASE WHEN ` + col + ` < 0 THEN 1 ELSE 0 END), 0),
      COALESCE(SUM(CASE WHEN ` + col + ` > 0 THEN 1 ELSE 0 END), 0)
    FROM ` + table)
	if scanErr := row.Scan(&zero, &neg, &pos); scanErr != nil {
		return 0, 0, 0, scanErr
	}
	if negativeIsError && neg > 0 {
		return zero, neg, pos, fmt.Errorf("preflight %s.%s: bucket_negative=%d must fail closed (Root D-012, voucher-scoped)", table, col, neg)
	}
	return zero, neg, pos, nil
}

// mustPreflight fails the test when the preflight rejects the column.
func mustPreflight(t *testing.T, db *sql.DB, table, col string, zeroIsAbsence, negativeIsError bool) {
	t.Helper()
	if _, _, _, err := preflight(db, table, col, zeroIsAbsence, negativeIsError); err != nil {
		t.Fatalf("preflight %s.%s: %v", table, col, err)
	}
}

// newShapeDDL is the frozen target shape (legacy columns minus defaults).
const newShapeDDL = `
CREATE TABLE converted (
  id          TEXT PRIMARY KEY,
  sec_nn      TEXT NOT NULL,
  ms_nn       TEXT NOT NULL,
  sec_null    TEXT,
  ms_null     TEXT,
  sec_d0      TEXT,
  ms_d0       TEXT,
  voucher_sec TEXT
)`

func newDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if _, err := db.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		t.Fatalf("pragma: %v", err)
	}
	return db
}

func insertLegacy(t *testing.T, db *sql.DB, id string, secNN, msNN int64, secNull, msNull *int64, secD0, msD0 int64, voucher *int64) {
	t.Helper()
	if _, err := db.Exec(
		`INSERT INTO legacy (id, sec_nn, ms_nn, sec_null, ms_null, sec_d0, ms_d0, voucher_sec)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		id, secNN, msNN, secNull, msNull, secD0, msD0, voucher,
	); err != nil {
		t.Fatalf("insert %s: %v", id, err)
	}
}

func int64p(v int64) *int64 { return &v }

// TestRebuildConversionBoundaries is the F-I-002 executable boundary matrix:
// negative milliseconds, the 999 ms no-carry case, epoch 0, year 9999, NULL
// preservation, sentinel 0 to NULL, and the voucher 0/NULL rule.
func TestRebuildConversionBoundaries(t *testing.T) {
	db := newDB(t)
	if _, err := db.Exec(legacyDDL); err != nil {
		t.Fatalf("legacy ddl: %v", err)
	}

	type sample struct {
		id       string
		secNN    int64
		msNN     int64
		wantSec  string
		wantMs   string
		secNull  *int64
		msNull   *int64
		secD0    int64
		msD0     int64
		wantD0   string // expectation for sec_d0 after conversion; "" means no assertion
		voucher  *int64
		wantVouc string // "" means no assertion
	}

	samples := []sample{
		{id: "epoch0", secNN: 0, msNN: 0,
			wantSec: "1970-01-01T00:00:00.000000Z", wantMs: "1970-01-01T00:00:00.000000Z"},
		{id: "positive", secNN: 1758320000, msNN: 1758320000123,
			wantSec: "2025-09-19T22:13:20.000000Z", wantMs: "2025-09-19T22:13:20.123000Z"},
		{id: "ms999_no_carry", secNN: 1758320000, msNN: 1758320000999,
			wantSec: "2025-09-19T22:13:20.000000Z", wantMs: "2025-09-19T22:13:20.999000Z"},
		{id: "neg_minus_1ms", secNN: 1, msNN: -1,
			wantSec: "1970-01-01T00:00:01.000000Z", wantMs: "1969-12-31T23:59:59.999000Z"},
		{id: "neg_minus_999ms", secNN: 1, msNN: -999,
			wantSec: "1970-01-01T00:00:01.000000Z", wantMs: "1969-12-31T23:59:59.001000Z"},
		{id: "neg_minus_1000ms", secNN: 1, msNN: -1000,
			wantSec: "1970-01-01T00:00:01.000000Z", wantMs: "1969-12-31T23:59:59.000000Z"},
		{id: "neg_minus_1001ms", secNN: 1, msNN: -1001,
			wantSec: "1970-01-01T00:00:01.000000Z", wantMs: "1969-12-31T23:59:58.999000Z"},
		{id: "neg_day", secNN: -86400, msNN: -86400000,
			wantSec: "1969-12-31T00:00:00.000000Z", wantMs: "1969-12-31T00:00:00.000000Z"},
		{id: "neg_pre1970", secNN: -1, msNN: -1758320000123,
			wantSec: "1969-12-31T23:59:59.000000Z", wantMs: "1914-04-14T01:46:39.877000Z"},
		{id: "year9999", secNN: 253402300799, msNN: 253402300799999,
			wantSec: "9999-12-31T23:59:59.000000Z", wantMs: "9999-12-31T23:59:59.999000Z"},
		{id: "nulls_and_d0", secNN: 1758320000, msNN: 1758320000123,
			wantSec: "2025-09-19T22:13:20.000000Z", wantMs: "2025-09-19T22:13:20.123000Z",
			secD0: 0, msD0: 0, wantD0: "NULL"},
		{id: "d0_positive", secNN: 1758320000, msNN: 1758320000123,
			wantSec: "2025-09-19T22:13:20.000000Z", wantMs: "2025-09-19T22:13:20.123000Z",
			secD0: 1758320000, msD0: 1758320000123,
			wantD0: "2025-09-19T22:13:20.000000Z"},
		{id: "voucher", secNN: 1758320000, msNN: 1758320000123,
			wantSec: "2025-09-19T22:13:20.000000Z", wantMs: "2025-09-19T22:13:20.123000Z",
			voucher: int64p(0), wantVouc: "NULL"},
		{id: "nullable_present", secNN: 1758320000, msNN: 1758320000123,
			wantSec: "2025-09-19T22:13:20.000000Z", wantMs: "2025-09-19T22:13:20.123000Z",
			secNull: int64p(-1), msNull: int64p(-1000)},
	}

	for _, s := range samples {
		insertLegacy(t, db, s.id, s.secNN, s.msNN, s.secNull, s.msNull, s.secD0, s.msD0, s.voucher)
	}

	// m0 preflight: only the voucher column treats a negative value as an error
	// (Root D-012). Sentinel (D0) columns map 0 to NULL. Negatives elsewhere are
	// valid pre-epoch instants (Root D-015).
	mustPreflight(t, db, "legacy", "sec_d0", true, false)
	mustPreflight(t, db, "legacy", "ms_d0", true, false)
	mustPreflight(t, db, "legacy", "voucher_sec", true, true)
	mustPreflight(t, db, "legacy", "sec_nn", false, false)
	mustPreflight(t, db, "legacy", "ms_nn", false, false)

	// Rebuild inside one transaction, mirroring the runner contract.
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("begin: %v", err)
	}
	if _, err := tx.Exec(`ALTER TABLE legacy RENAME TO legacy_old`); err != nil {
		t.Fatalf("rename: %v", err)
	}
	if _, err := tx.Exec(newShapeDDL); err != nil {
		t.Fatalf("new ddl: %v", err)
	}
	copyStmt := `INSERT INTO converted (id, sec_nn, ms_nn, sec_null, ms_null, sec_d0, ms_d0, voucher_sec)
    SELECT id,
           ` + secondsExpr("sec_nn") + `,
           ` + millisExpr("ms_nn") + `,
           CASE WHEN sec_null IS NULL THEN NULL ELSE ` + secondsExpr("sec_null") + ` END,
           CASE WHEN ms_null IS NULL THEN NULL ELSE ` + millisExpr("ms_null") + ` END,
           CASE WHEN sec_d0 = 0 THEN NULL ELSE ` + secondsExpr("sec_d0") + ` END,
           CASE WHEN ms_d0 = 0 THEN NULL ELSE ` + millisExpr("ms_d0") + ` END,
           CASE WHEN voucher_sec IS NULL OR voucher_sec = 0 THEN NULL ELSE ` + secondsExpr("voucher_sec") + ` END
    FROM legacy_old`
	if _, err := tx.Exec(copyStmt); err != nil {
		t.Fatalf("copy: %v", err)
	}
	if _, err := tx.Exec(`DROP TABLE legacy_old`); err != nil {
		t.Fatalf("drop old: %v", err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatalf("commit: %v", err)
	}

	// Assertions.
	for _, s := range samples {
		var gotSec, gotMs string
		var gotSecNull, gotMsNull, gotSecD0, gotMsD0, gotVoucher sql.NullString
		row := db.QueryRow(`SELECT sec_nn, ms_nn, sec_null, ms_null, sec_d0, ms_d0, voucher_sec
		                    FROM converted WHERE id = ?`, s.id)
		if err := row.Scan(&gotSec, &gotMs, &gotSecNull, &gotMsNull, &gotSecD0, &gotMsD0, &gotVoucher); err != nil {
			t.Fatalf("scan %s: %v", s.id, err)
		}
		if gotSec != s.wantSec {
			t.Errorf("%s sec_nn = %q, want %q", s.id, gotSec, s.wantSec)
		}
		if gotMs != s.wantMs {
			t.Errorf("%s ms_nn = %q, want %q", s.id, gotMs, s.wantMs)
		}
		// Fixed-6 contract: 27 characters for all post-1970 values; pre-1970
		// values must still be fixed width (also 27) because the format is
		// fixed-width by construction.
		for name, v := range map[string]string{"sec_nn": gotSec, "ms_nn": gotMs} {
			if len(v) != 27 {
				t.Errorf("%s %s len = %d (%q), want 27", s.id, name, len(v), v)
			}
		}
		// NULL preservation.
		if s.secNull == nil && gotSecNull.Valid {
			t.Errorf("%s sec_null should stay NULL, got %q", s.id, gotSecNull.String)
		}
		if s.msNull == nil && gotMsNull.Valid {
			t.Errorf("%s ms_null should stay NULL, got %q", s.id, gotMsNull.String)
		}
		if s.secNull != nil && !gotSecNull.Valid {
			t.Errorf("%s sec_null should be present (legacy %d)", s.id, *s.secNull)
		}
		if s.wantD0 != "" {
			want := s.wantD0
			if want == "NULL" {
				if gotSecD0.Valid || gotMsD0.Valid {
					t.Errorf("%s sec_d0/ms_d0 = %q/%q, want NULL/NULL", s.id, gotSecD0.String, gotMsD0.String)
				}
			} else if !gotSecD0.Valid || gotSecD0.String != want {
				t.Errorf("%s sec_d0 = %q, want %q", s.id, gotSecD0.String, want)
			}
		}
		if s.wantVouc != "" {
			if gotVoucher.Valid {
				t.Errorf("%s voucher_sec = %q, want NULL (legacy 0 means absent)", s.id, gotVoucher.String)
			}
		}
	}

	// SQLite-side integrity after the rebuild.
	var fkCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM pragma_foreign_key_check`).Scan(&fkCount); err != nil {
		t.Fatalf("fk check: %v", err)
	}
	if fkCount != 0 {
		t.Errorf("foreign_key_check rows = %d, want 0", fkCount)
	}
	var integrity string
	if err := db.QueryRow(`SELECT * FROM pragma_integrity_check`).Scan(&integrity); err != nil {
		t.Fatalf("integrity check: %v", err)
	}
	if integrity != "ok" {
		t.Errorf("integrity_check = %q, want ok", integrity)
	}
	// The retired legacy table must be gone.
	var oldExists int
	if err := db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='legacy_old'`).Scan(&oldExists); err != nil {
		t.Fatalf("sqlite_master: %v", err)
	}
	if oldExists != 0 {
		t.Errorf("legacy_old still exists")
	}
}

// TestNegativeMustFailClosed proves the voucher-scoped negative rule (Root
// D-012): a negative value in vouchers.expires_at / redeemed_at is rejected by
// the m0 preflight rather than silently converted. It also proves the converse
// from Root D-015: a negative value in an ordinary time column is NOT an error,
// because a negative epoch is a valid instant rather than a sentinel.
func TestNegativeMustFailClosed(t *testing.T) {
	// Voucher column with a negative value must fail closed. The preflight
	// returns an error rather than aborting, so this is a normal assertion.
	t.Run("voucher_negative_fails_closed", func(t *testing.T) {
		db := newDB(t)
		if _, err := db.Exec(legacyDDL); err != nil {
			t.Fatalf("legacy ddl: %v", err)
		}
		insertLegacy(t, db, "neg", 1758320000, 1758320000123, nil, nil, 0, 0, int64p(-1))
		if _, neg, _, err := preflight(db, "legacy", "voucher_sec", true, true); err == nil {
			t.Fatalf("preflight accepted a negative voucher value (neg bucket=%d); Root D-012 requires fail closed", neg)
		}
	})

	// Ordinary time column with a negative value is legal (negative epoch).
	t.Run("ordinary_negative_is_valid_instant", func(t *testing.T) {
		db := newDB(t)
		if _, err := db.Exec(legacyDDL); err != nil {
			t.Fatalf("legacy ddl: %v", err)
		}
		insertLegacy(t, db, "pre1970", -1, -1000, nil, nil, 0, 0, nil)
		zero, neg, pos, err := preflight(db, "legacy", "sec_nn", false, false)
		if err != nil {
			t.Fatalf("preflight rejected a negative epoch: %v (Root D-015: negative epoch is not a sentinel)", err)
		}
		if zero != 0 || neg != 1 || pos != 0 {
			t.Errorf("buckets zero/neg/pos = %d/%d/%d, want 0/1/0", zero, neg, pos)
		}
	})
}

// TestFixedSixLexicalOrder proves that the canonical fixed-6 TEXT form sorts
// lexically in instant order, which the contract relies on for SQLite ORDER BY.
func TestFixedSixLexicalOrder(t *testing.T) {
	db := newDB(t)
	if _, err := db.Exec(`CREATE TABLE ord (id INTEGER PRIMARY KEY, v TEXT NOT NULL)`); err != nil {
		t.Fatalf("ddl: %v", err)
	}
	instants := []int64{-86400, -1, 0, 1, 1758320000, 253402300799}
	for _, n := range instants {
		if _, err := db.Exec(`INSERT INTO ord (v) SELECT `+secondsExpr("?"), n); err != nil {
			t.Fatalf("insert %d: %v", n, err)
		}
	}
	rows, err := db.Query(`SELECT v FROM ord ORDER BY v ASC`)
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	defer rows.Close()
	var got []string
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			t.Fatalf("scan: %v", err)
		}
		got = append(got, v)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("rows: %v", err)
	}
	want := []string{
		"1969-12-31T00:00:00.000000Z", // -86400
		"1969-12-31T23:59:59.000000Z", // -1
		"1970-01-01T00:00:00.000000Z", // 0
		"1970-01-01T00:00:01.000000Z", // 1
		"2025-09-19T22:13:20.000000Z", // 1758320000
		"9999-12-31T23:59:59.000000Z", // 253402300799
	}
	if len(got) != len(want) {
		t.Fatalf("got %d rows, want %d: %v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("order[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
