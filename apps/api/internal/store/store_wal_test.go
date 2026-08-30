package store

// W25 · 防复发回归测试：SQLite 连接面配置（池 + pragma）一旦被改回
// MaxOpenConns=1 / DELETE journal / 无超时，本组测试立即失败——防止
// "我的钱包页面性能"优化（2026-08-23）被后续改动静默回退。

import (
	"context"
	"database/sql"
	"path/filepath"
	"strings"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

func TestSQLiteDSNPragmas(t *testing.T) {
	t.Run("memory stays bare and single-connection", func(t *testing.T) {
		for _, path := range []string{"", ":memory:", "file::memory:?cache=shared"} {
			dsn, memory := sqliteDSN(path)
			if !memory {
				t.Errorf("sqliteDSN(%q): memory = false, want true", path)
			}
			if dsn != path {
				t.Errorf("sqliteDSN(%q) = %q, want unchanged dsn", path, dsn)
			}
		}
	})
	t.Run("file dsns carry the pragma triplet", func(t *testing.T) {
		dsn, memory := sqliteDSN(filepath.Join("C:", "data", "app.db"))
		if memory {
			t.Fatal("file dsn reported as memory")
		}
		// GOAL-036 A-004 F-010 (independent): _txlock=immediate is part of the
		// connection-surface invariant (fixes WAL read-then-write snapshot
		// races); the DSN test must pin it like the other four parameters.
		for _, want := range []string{"_busy_timeout=5000", "_journal_mode=WAL", "_synchronous=NORMAL", "_foreign_keys=on", "_txlock=immediate"} {
			if !strings.Contains(dsn, want) {
				t.Errorf("dsn %q missing %q", dsn, want)
			}
		}
	})
	t.Run("existing query strings get an ampersand", func(t *testing.T) {
		dsn, _ := sqliteDSN("app.db?_pragma=foreign_keys=on")
		if !strings.Contains(dsn, "&_busy_timeout=5000") {
			t.Errorf("dsn %q should append pragmas with &", dsn)
		}
	})
}

// TestFileStoreWALPoolAndPragma opens a real file store and asserts the
// connection surface the wallet page depends on: a pool (no global
// serialization), WAL journal, busy timeout and normal synchronous mode.
func TestFileStoreWALPoolAndPragma(t *testing.T) {
	path := filepath.Join(t.TempDir(), "wal.db")
	st, err := OpenWithCatalog(path, MigrationCatalog())
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	defer st.Close()

	if got := st.db.Stats().MaxOpenConnections; got != sqlitePoolDefault {
		t.Errorf("MaxOpenConnections = %d, want sqlitePoolDefault (%d)", got, sqlitePoolDefault)
	}

	var mode string
	if err := st.db.QueryRow("PRAGMA journal_mode").Scan(&mode); err != nil {
		t.Fatalf("PRAGMA journal_mode: %v", err)
	}
	if mode != "wal" {
		t.Errorf("journal_mode = %q, want wal", mode)
	}

	var busyTimeout int
	if err := st.db.QueryRow("PRAGMA busy_timeout").Scan(&busyTimeout); err != nil {
		t.Fatalf("PRAGMA busy_timeout: %v", err)
	}
	if busyTimeout != 5000 {
		t.Errorf("busy_timeout = %d, want 5000", busyTimeout)
	}

	var sync int
	if err := st.db.QueryRow("PRAGMA synchronous").Scan(&sync); err != nil {
		t.Fatalf("PRAGMA synchronous: %v", err)
	}
	// SQLite reports synchronous as its numeric form: 1 == NORMAL.
	if sync != 1 {
		t.Errorf("synchronous = %d, want 1 (NORMAL)", sync)
	}
}

// TestMemoryStoreStaysSingleConnection guards the in-memory invariant: a pool
// on :memory: would fan transactions out across separate empty databases.
func TestMemoryStoreStaysSingleConnection(t *testing.T) {
	st, err := OpenWithCatalog(":memory:", MigrationCatalog())
	if err != nil {
		t.Fatalf("open memory store: %v", err)
	}
	defer st.Close()
	if got := st.db.Stats().MaxOpenConnections; got != 1 {
		t.Errorf("memory MaxOpenConnections = %d, want 1", got)
	}
}

// TestFileStoreEveryConnectionEnforcesForeignKeys (A-001 F-002, independent
// audit 2026-08-23): PRAGMA foreign_keys is CONNECTION-scoped. The pool must
// carry it on every connection (via the driver-level DSN parameter), not only
// on the connection the migration runner happened to use — otherwise ON DELETE
// CASCADE silently stops firing and FK invariants (_user_roles_ cascade, RBAC
// RESTRICT, refresh_tokens checks) are void on the other pooled connections,
// which is exactly how the W25 e2e regression showed up after pooling.
func TestFileStoreEveryConnectionEnforcesForeignKeys(t *testing.T) {
	path := filepath.Join(t.TempDir(), "fk-pool.db")
	st, err := OpenWithCatalog(path, MigrationCatalog())
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	defer st.Close()
	ctx := context.Background()

	// Hold one connection per pool slot and assert FK is ON on each of them.
	conns := make([]*sql.Conn, 0, sqlitePoolDefault)
	defer func() {
		for _, c := range conns {
			_ = c.Close()
		}
	}()
	for i := 0; i < sqlitePoolDefault; i++ {
		c, err := st.db.Conn(ctx)
		if err != nil {
			t.Fatalf("conn %d: %v", i, err)
		}
		conns = append(conns, c)
		var enabled int
		if err := c.QueryRowContext(ctx, "PRAGMA foreign_keys").Scan(&enabled); err != nil {
			t.Fatalf("conn %d PRAGMA foreign_keys: %v", i, err)
		}
		if enabled != 1 {
			t.Errorf("conn %d foreign_keys = %d, want 1 (ON)", i, enabled)
		}
	}

	// ON DELETE CASCADE fires on a NON-first pooled connection: deleting a
	// user must cascade-remove its user_roles link (the W25 e2e scenario).
	last := conns[len(conns)-1]
	seed := []string{
		`INSERT INTO users (id, username, name, roles, password_hash, must_change_password, created_at, updated_at)
		 VALUES ('u-fk', 'fk', 'FK', '[]', 'h', 0, 1, 1)`,
		`INSERT INTO roles (id, key, name, created_at, updated_at) VALUES ('role-fk', 'role-fk', 'R', 1, 1)`,
		`INSERT INTO user_roles (user_id, role_id) VALUES ('u-fk', 'role-fk')`,
	}
	for _, q := range seed {
		if _, err := last.ExecContext(ctx, q); err != nil {
			t.Fatalf("seed: %v", err)
		}
	}
	// RESTRICT (while the link exists): deleting the in-use role must fail on
	// the non-first connection too.
	if _, err := last.ExecContext(ctx, `DELETE FROM roles WHERE id = 'role-fk'`); err == nil {
		t.Error("delete in-use role succeeded, want FK RESTRICT violation")
	}
	if _, err := last.ExecContext(ctx, `DELETE FROM users WHERE id = 'u-fk'`); err != nil {
		t.Fatalf("delete user on non-first conn: %v", err)
	}
	var leftover int
	if err := last.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM user_roles WHERE user_id = 'u-fk'`,
	).Scan(&leftover); err != nil {
		t.Fatalf("count leftover user_roles: %v", err)
	}
	if leftover != 0 {
		t.Errorf("user_roles leftover after CASCADE = %d, want 0", leftover)
	}

	// FK rejection also works off the migrate connection: a refresh_tokens row
	// pointing at a missing user must fail.
	if _, err := last.ExecContext(ctx,
		`INSERT INTO refresh_tokens (id, user_id, created_at) VALUES ('rt-fk', 'missing-user', 1)`,
	); err == nil {
		t.Error("insert refresh_tokens with missing user succeeded, want FK violation")
	}
}

// TestFileStorePoolOverride verifies DB_POOL_MAX_OPEN wiring through
// store.Open(OpenOptions).
func TestFileStorePoolOverride(t *testing.T) {
	path := filepath.Join(t.TempDir(), "override.db")
	st, err := Open(context.Background(), OpenOptions{
		Dialect:          kernel.DialectSQLite,
		Path:             path,
		PoolMaxOpenConns: 7,
	}, MigrationCatalog())
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	defer st.Close()
	raw, ok := st.(*Store)
	if !ok {
		t.Fatalf("Open returned %T, want *Store", st)
	}
	if got := raw.db.Stats().MaxOpenConnections; got != 7 {
		t.Errorf("MaxOpenConnections = %d, want 7", got)
	}
}