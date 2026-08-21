package store

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

func openSQLiteKernel(t *testing.T) kernel.Store {
	t.Helper()
	st, err := Open(context.Background(), OpenOptions{
		Dialect: kernel.DialectSQLite,
		Path:    filepath.Join(t.TempDir(), "kernel.db"),
	}, MigrationCatalog())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return st
}

func TestKernelStoreSQLiteDialectFreshAndPing(t *testing.T) {
	st := openSQLiteKernel(t)
	if st.Dialect() != kernel.DialectSQLite {
		t.Fatalf("dialect = %q, want %q", st.Dialect(), kernel.DialectSQLite)
	}
	if !st.WasFresh() {
		t.Fatal("fresh temp db must report WasFresh()=true")
	}
	if err := st.Ping(context.Background()); err != nil {
		t.Fatalf("ping: %v", err)
	}
}

func TestKernelStoreRunCommitAndRollback(t *testing.T) {
	st := openSQLiteKernel(t)

	var rows int64
	if err := st.Run(context.Background(), func(tx kernel.Tx) error {
		if _, err := tx.Exec(context.Background(), "CREATE TABLE kernel_probe (id INTEGER PRIMARY KEY, name TEXT)"); err != nil {
			return err
		}
		if _, err := tx.Exec(context.Background(), "INSERT INTO kernel_probe (name) VALUES (?)", "hello"); err != nil {
			return err
		}
		return tx.QueryRow(context.Background(), "SELECT count(*) FROM kernel_probe").Scan(&rows)
	}); err != nil {
		t.Fatalf("run with commit: %v", err)
	}
	if rows != 1 {
		t.Fatalf("inserted rows = %d, want 1", rows)
	}

	// Commit must persist: a second Run sees the table and row.
	var after int
	if err := st.Run(context.Background(), func(tx kernel.Tx) error {
		return tx.QueryRow(context.Background(), "SELECT count(*) FROM kernel_probe").Scan(&after)
	}); err != nil {
		t.Fatalf("second run: %v", err)
	}
	if after != 1 {
		t.Fatalf("committed row not visible: %d", after)
	}

	// An error rolls back the whole tx.
	if err := st.Run(context.Background(), func(tx kernel.Tx) error {
		if _, err := tx.Exec(context.Background(), "INSERT INTO kernel_probe (name) VALUES (?)", "rollback"); err != nil {
			return err
		}
		return errors.New("boom")
	}); err == nil {
		t.Fatal("expected run error")
	}
	var afterRollback int
	if err := st.Run(context.Background(), func(tx kernel.Tx) error {
		return tx.QueryRow(context.Background(), "SELECT count(*) FROM kernel_probe").Scan(&afterRollback)
	}); err != nil {
		t.Fatalf("third run: %v", err)
	}
	if afterRollback != 1 {
		t.Fatalf("rollback left a row behind: %d", afterRollback)
	}
}

func TestKernelStoreErrNoRowsMapping(t *testing.T) {
	st := openSQLiteKernel(t)
	if err := st.Run(context.Background(), func(tx kernel.Tx) error {
		_, err := tx.Exec(context.Background(), "CREATE TABLE kernel_probe (id INTEGER PRIMARY KEY, name TEXT)")
		return err
	}); err != nil {
		t.Fatal(err)
	}
	err := st.Run(context.Background(), func(tx kernel.Tx) error {
		var name string
		// A non-existent row must surface kernel.ErrNoRows (and sql.ErrNoRows
		// until R4 cuts database/sql out of the public surface).
		return tx.QueryRow(context.Background(), "SELECT name FROM kernel_probe WHERE id = 42").Scan(&name)
	})
	if err == nil {
		t.Fatal("expected no-rows error")
	}
	if !errors.Is(err, kernel.ErrNoRows) {
		t.Fatalf("errors.Is(err, kernel.ErrNoRows) = false, got %v", err)
	}
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("errors.Is(err, sql.ErrNoRows) = false, got %v", err)
	}
}

func TestKernelStoreNestedRunFailsClosed(t *testing.T) {
	st := openSQLiteKernel(t)
	err := st.Run(context.Background(), func(tx kernel.Tx) error {
		return st.Run(context.Background(), func(tx2 kernel.Tx) error { return nil })
	})
	if err == nil || !strings.Contains(err.Error(), "nested Run") {
		t.Fatalf("nested Run must fail closed, got %v", err)
	}
}

func TestKernelStoreRunPanicRollsBackAndRepanics(t *testing.T) {
	st := openSQLiteKernel(t)

	// A panicking fn must roll back (the CREATE TABLE must not persist) and
	// re-panic so the caller observes the panic (R1 v1.4 §2).
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("Run must re-panic when fn panics")
			}
		}()
		_ = st.Run(context.Background(), func(tx kernel.Tx) error {
			if _, err := tx.Exec(context.Background(), "CREATE TABLE kernel_panic_probe (id INTEGER PRIMARY KEY)"); err != nil {
				return err
			}
			panic("boom")
		})
	}()

	var n int
	err := st.Run(context.Background(), func(tx kernel.Tx) error {
		return tx.QueryRow(context.Background(), "SELECT count(*) FROM kernel_panic_probe").Scan(&n)
	})
	if err == nil || !strings.Contains(err.Error(), "no such table") {
		t.Fatalf("DDL should have been rolled back after fn panic, got %v", err)
	}
}

func TestKernelStoreRunConcurrentGoroutinesNoFalseNesting(t *testing.T) {
	st := openSQLiteKernel(t)
	if err := st.Run(context.Background(), func(tx kernel.Tx) error {
		_, err := tx.Exec(context.Background(), "CREATE TABLE kernel_conc_probe (id INTEGER PRIMARY KEY, g INTEGER)")
		return err
	}); err != nil {
		t.Fatal(err)
	}

	const n = 8
	var wg sync.WaitGroup
	errs := make([]error, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			errs[i] = st.Run(context.Background(), func(tx kernel.Tx) error {
				if _, err := tx.Exec(context.Background(), "INSERT INTO kernel_conc_probe (g) VALUES (?)", i); err != nil {
					return err
				}
				var cnt int
				return tx.QueryRow(context.Background(), "SELECT count(*) FROM kernel_conc_probe WHERE g = ?", i).Scan(&cnt)
			})
		}(i)
	}
	wg.Wait()
	for i, err := range errs {
		if err != nil {
			t.Fatalf("goroutine %d run failed: %v (false 'nested Run' from shared goroutine marker?)", i, err)
		}
	}
	var total int
	if err := st.Run(context.Background(), func(tx kernel.Tx) error {
		return tx.QueryRow(context.Background(), "SELECT count(*) FROM kernel_conc_probe").Scan(&total)
	}); err != nil {
		t.Fatal(err)
	}
	if total != n {
		t.Fatalf("committed rows = %d, want %d", total, n)
	}
}

func TestKernelStoreQueryScan(t *testing.T) {
	st := openSQLiteKernel(t)
	if err := st.Run(context.Background(), func(tx kernel.Tx) error {
		if _, err := tx.Exec(context.Background(), "CREATE TABLE kernel_probe (id INTEGER PRIMARY KEY, name TEXT)"); err != nil {
			return err
		}
		for _, n := range []string{"a", "b", "c"} {
			if _, err := tx.Exec(context.Background(), "INSERT INTO kernel_probe (name) VALUES (?)", n); err != nil {
				return err
			}
		}
		rows, err := tx.Query(context.Background(), "SELECT name FROM kernel_probe ORDER BY id")
		if err != nil {
			return err
		}
		defer rows.Close()
		got := []string{}
		for rows.Next() {
			var n string
			if err := rows.Scan(&n); err != nil {
				return err
			}
			got = append(got, n)
		}
		if err := rows.Err(); err != nil {
			return err
		}
		if !reflect.DeepEqual(got, []string{"a", "b", "c"}) {
			t.Errorf("scan = %v, want [a b c]", got)
		}
		return nil
	}); err != nil {
		t.Fatalf("query scan: %v", err)
	}
}
