package store

import (
	"database/sql"
	"strings"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// TestTemporalReadAdapter covers the read half of the workspace-040 temporal
// contract: canonical SQLite text and native postgres instants both land in
// domain destinations, NULL follows sql.NullTime semantics, and non-canonical
// stored text fails closed.
func TestTemporalReadAdapter(t *testing.T) {
	path := t.TempDir() + "/adapter.db"
	st, err := OpenWithCatalog(path, MigrationCatalog())
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer st.Close()

	ctx := t.Context()
	canonical := "2026-09-20T12:40:48.814963Z"
	if err := st.Run(ctx, func(tx kernel.Tx) error {
		if _, err := tx.Exec(ctx, `UPDATE site_settings SET updated_at = ? WHERE id = 'default'`, canonical); err != nil {
			return err
		}
		var updated time.Time
		if err := tx.QueryRow(ctx, `SELECT updated_at FROM site_settings WHERE id = 'default'`).Scan(&updated); err != nil {
			return err
		}
		if got := updated.UTC().Format("2006-01-02T15:04:05.000000Z"); got != canonical {
			t.Fatalf("time.Time scan = %s, want %s", got, canonical)
		}
		var nullable sql.NullTime
		if err := tx.QueryRow(ctx, `SELECT locked_until FROM users WHERE id = 'nobody'`).Scan(&nullable); err != sql.ErrNoRows {
			t.Fatalf("missing row error = %v, want ErrNoRows", err)
		}
		var rowCount int
		if err := tx.QueryRow(ctx, `SELECT COUNT(*) FROM site_settings`).Scan(&rowCount); err != nil {
			return err
		}
		if rowCount != 1 {
			t.Fatalf("site_settings rows = %d, want 1", rowCount)
		}
		return nil
	}); err != nil {
		t.Fatalf("canonical round trip: %v", err)
	}

	// A NULL instant is valid/absent for sql.NullTime and **time.Time, and the
	// canonical-only rule rejects any other spelling of the same instant.
	if err := st.Run(ctx, func(tx kernel.Tx) error {
		var nullable sql.NullTime
		if err := tx.QueryRow(ctx, `SELECT NULL`).Scan(&nullable); err != nil {
			return err
		}
		if nullable.Valid {
			t.Fatalf("NULL scanned into sql.NullTime as valid")
		}
		var pointer *time.Time
		if err := tx.QueryRow(ctx, `SELECT NULL`).Scan(&pointer); err != nil {
			return err
		}
		if pointer != nil {
			t.Fatalf("NULL scanned into *time.Time as %v", pointer)
		}
		// 3-digit fraction: same instant, non-canonical storage.
		if err := tx.QueryRow(ctx, `SELECT '2026-09-20T12:40:48.814Z'`).Scan(&nullable); err == nil {
			t.Fatalf("non-canonical stored text was accepted")
		} else if !strings.Contains(err.Error(), "not canonical") {
			t.Fatalf("non-canonical error = %v, want a canonical-form complaint", err)
		}
		// Local offset: also rejected rather than silently converted.
		var instant time.Time
		if err := tx.QueryRow(ctx, `SELECT '2026-09-20T12:40:48.814963+08:00'`).Scan(&instant); err == nil {
			t.Fatalf("offset stored text was accepted")
		}
		// NULL into a non-nullable destination stays an error, as in database/sql.
		if err := tx.QueryRow(ctx, `SELECT NULL`).Scan(&instant); err == nil {
			t.Fatalf("NULL into *time.Time was accepted")
		}
		return nil
	}); err != nil {
		t.Fatalf("null and strictness rules: %v", err)
	}

	// A query with no temporal destination still passes every column through.
	if err := st.Run(ctx, func(tx kernel.Tx) error {
		rows, err := tx.Query(ctx, `SELECT id, site_title FROM site_settings`)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var id, title string
			if err := rows.Scan(&id, &title); err != nil {
				return err
			}
		}
		return rows.Err()
	}); err != nil {
		t.Fatalf("pass-through scan: %v", err)
	}
}
