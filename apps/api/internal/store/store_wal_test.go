package store

// W25 · 防复发回归测试：SQLite 连接面配置（池 + pragma）一旦被改回
// MaxOpenConns=1 / DELETE journal / 无超时，本组测试立即失败——防止
// "我的钱包页面性能"优化（2026-08-23）被后续改动静默回退。

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
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
		for _, want := range []string{"_busy_timeout=5000", "_journal_mode=WAL", "_synchronous=NORMAL"} {
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