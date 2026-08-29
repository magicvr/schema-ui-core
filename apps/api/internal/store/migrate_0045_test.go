package store

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/modules/operationlog"
)

func TestMigrate0045PreservesOperationCorrelations(t *testing.T) {
	path := filepath.Join(t.TempDir(), "operation-log-0044.db")
	catalog := MigrationCatalog()
	st, err := OpenWithCatalog(path, catalog[:44])
	if err != nil {
		t.Fatal(err)
	}
	repository := operationlog.NewRepository(st)
	if err := repository.RecordOperation(operationlog.Operation{
		ID: "op-before-0045", Event: operationlog.EventWalletReconcile,
		ActorID: "user-1", ActorName: "User One", CorrelationID: "corr-before-0045",
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
	repository = operationlog.NewRepository(st)
	op, err := repository.GetOperation("op-before-0045")
	if err != nil {
		t.Fatal(err)
	}
	if op.CorrelationID != "corr-before-0045" {
		t.Fatalf("correlation = %q, want preserved value", op.CorrelationID)
	}
	for index, event := range []string{
		operationlog.EventServiceCredentialCreate,
		operationlog.EventServiceCredentialUse,
		operationlog.EventServiceCredentialRevoke,
	} {
		if err := repository.RecordOperation(operationlog.Operation{
			ID: "op-service-credential-" + string(rune('1'+index)), Event: event,
			ActorID: "user-1", ActorName: "User One", CreatedAt: time.Unix(1700000001+int64(index), 0).UTC(),
		}); err != nil {
			t.Fatalf("record %s after 0045: %v", event, err)
		}
	}
}
