package store

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
)

func TestMigrate0048AddsSessionSideTable(t *testing.T) {
	path := filepath.Join(t.TempDir(), "operation-log-0047.db")
	catalog := MigrationCatalog()
	st, err := OpenWithCatalog(path, catalog[:47])
	if err != nil {
		t.Fatal(err)
	}
	repository := operationlog.NewRepository(st)
	if err := repository.RecordOperation(operationlog.Operation{
		ID: "op-before-0048", Event: operationlog.EventSettingsUpdate,
		ActorID: "user-1", ActorName: "User One", CorrelationID: "corr-before-0048",
		CreatedAt: time.Unix(1700000000, 0).UTC(),
	}); err != nil {
		t.Fatal(err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}

	st, err = OpenWithCatalog(path, catalog)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	if !tableExistsDB(t, st.db, "operation_log_session") || !tableExistsDB(t, st.db, "operation_log_archive_session") {
		t.Fatal("0048 session tables missing")
	}
	repository = operationlog.NewRepository(st)
	op, err := repository.GetOperation("op-before-0048")
	if err != nil {
		t.Fatal(err)
	}
	if op.CorrelationID != "corr-before-0048" {
		t.Fatalf("correlation = %q, want preserved value", op.CorrelationID)
	}
	if op.SessionID != "" {
		t.Fatalf("pre-0048 session = %q, want empty", op.SessionID)
	}
	if err := repository.RecordOperation(operationlog.Operation{
		ID: "op-after-0048", Event: operationlog.EventSettingsUpdate,
		ActorID: "user-1", ActorName: "User One", SessionID: "sess-after-0048",
		CreatedAt: time.Unix(1700000001, 0).UTC(),
	}); err != nil {
		t.Fatalf("record after 0048: %v", err)
	}
	got, err := repository.GetOperation("op-after-0048")
	if err != nil {
		t.Fatal(err)
	}
	if got.SessionID != "sess-after-0048" {
		t.Fatalf("session = %q, want sess-after-0048", got.SessionID)
	}
}
