package store

import (
	"context"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// TestPostgresWriteTruncatesToMicroseconds covers workspace-040 GOAL-004 A-002
// F-I-001: timestamptz(6)'s typmod ROUNDs, so the postgres bind adapter must
// truncate toward zero exactly like the SQLite canonical text does. Without the
// truncation the same input stores ...123457 on postgres and ...123456 on
// sqlite.
func TestPostgresWriteTruncatesToMicroseconds(t *testing.T) {
	st := postgresScratchDB(t, "vp040trunc")
	ctx := context.Background()

	// 789 ns past a microsecond boundary: rounding would carry to .123457.
	input := time.Date(2025, 9, 19, 22, 13, 20, 123456789, time.UTC)
	if err := st.Run(ctx, func(tx kernel.Tx) error {
		_, err := tx.Exec(ctx, `UPDATE site_settings SET updated_at = ? WHERE id = 'default'`, input)
		return err
	}); err != nil {
		t.Fatalf("bind high-precision instant: %v", err)
	}

	if err := st.Run(ctx, func(tx kernel.Tx) error {
		var got time.Time
		if err := tx.QueryRow(ctx, `SELECT updated_at FROM site_settings WHERE id = 'default'`).Scan(&got); err != nil {
			return err
		}
		if want := "2025-09-19T22:13:20.123456Z"; got.UTC().Format("2006-01-02T15:04:05.000000Z") != want {
			t.Fatalf("postgres stored %s, want %s (truncate toward zero, never round)",
				got.UTC().Format("2006-01-02T15:04:05.000000Z"), want)
		}
		var asText string
		if err := tx.QueryRow(ctx, `SELECT to_char(updated_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS.US"Z"') FROM site_settings WHERE id = 'default'`).Scan(&asText); err != nil {
			return err
		}
		if asText != "2025-09-19T22:13:20.123456Z" {
			t.Fatalf("postgres text form = %s, want the truncated microsecond value", asText)
		}
		return nil
	}); err != nil {
		t.Fatalf("read back: %v", err)
	}
}

// TestSQLiteWriteTruncatesToMicroseconds is the SQLite half of the same rule:
// the canonical fixed-6 text drops sub-microsecond precision instead of
// rounding it.
func TestSQLiteWriteTruncatesToMicroseconds(t *testing.T) {
	path := t.TempDir() + "/trunc.db"
	st, err := OpenWithCatalog(path, MigrationCatalog())
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer st.Close()
	ctx := context.Background()
	input := time.Date(2025, 9, 19, 22, 13, 20, 123456789, time.UTC)
	if err := st.Run(ctx, func(tx kernel.Tx) error {
		if _, err := tx.Exec(ctx, `UPDATE site_settings SET updated_at = ? WHERE id = 'default'`, input); err != nil {
			return err
		}
		var got time.Time
		if err := tx.QueryRow(ctx, `SELECT updated_at FROM site_settings WHERE id = 'default'`).Scan(&got); err != nil {
			return err
		}
		if want := "2025-09-19T22:13:20.123456Z"; got.UTC().Format("2006-01-02T15:04:05.000000Z") != want {
			t.Fatalf("sqlite stored %s, want %s", got.UTC().Format("2006-01-02T15:04:05.000000Z"), want)
		}
		return nil
	}); err != nil {
		t.Fatalf("sqlite truncation: %v", err)
	}
}
