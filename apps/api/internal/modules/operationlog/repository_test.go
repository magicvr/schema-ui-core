package operationlog

import (
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)

func openOperationRepository(t *testing.T, name string) *Repository {
	t.Helper()
	store, err := testsupport.OpenStore(filepath.Join(t.TempDir(), name), "admin", "hash", true)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })
	return NewRepository(store)
}

func TestRepositoryAppendListFilterAndGet(t *testing.T) {
	repository := openOperationRepository(t, "operationlog.db")
	if got, err := repository.ListOperations(10); err != nil || len(got) != 0 {
		t.Fatalf("empty ListOperations = %v (err %v), want empty", got, err)
	}

	base := time.Date(2026, 8, 3, 12, 0, 0, 0, time.UTC)
	recordID := "user-9"
	detail := `{"username":"alice"}`
	operations := []Operation{
		{ID: "op-1", Event: EventUserCreate, ActorID: "user-admin", ActorName: "Admin", RecordID: &recordID, Detail: &detail, CorrelationID: "r1-op-1", CreatedAt: base},
		{ID: "op-2", Event: EventAuthLogin, ActorID: "user-admin", ActorName: "Admin", CreatedAt: base.Add(time.Second)},
		{ID: "op-3", Event: EventAuthLogout, ActorID: "user-editor", ActorName: "Editor", CreatedAt: base.Add(2 * time.Second)},
	}
	for _, operation := range operations {
		if err := repository.RecordOperation(operation); err != nil {
			t.Fatalf("RecordOperation(%s): %v", operation.ID, err)
		}
	}

	got, err := repository.ListOperations(2)
	if err != nil {
		t.Fatalf("ListOperations: %v", err)
	}
	if len(got) != 2 || got[0].ID != "op-3" || got[1].ID != "op-2" {
		t.Fatalf("ListOperations(2) = %+v, want op-3/op-2", got)
	}
	filtered, total, err := repository.ListOperationsFiltered(OperationFilter{
		Q: "alice", Sort: "event", Order: "asc", Page: 1, PageSize: 10,
	})
	if err != nil {
		t.Fatalf("ListOperationsFiltered: %v", err)
	}
	if total != 1 || len(filtered) != 1 || filtered[0].ID != "op-1" {
		t.Fatalf("filtered = %+v total=%d, want op-1/1", filtered, total)
	}
	operation, err := repository.GetOperation("op-1")
	if err != nil {
		t.Fatalf("GetOperation: %v", err)
	}
	if operation.RecordID == nil || *operation.RecordID != recordID || operation.Detail == nil || *operation.Detail != detail || operation.CorrelationID != "r1-op-1" || !operation.CreatedAt.Equal(base) {
		t.Fatalf("GetOperation(op-1) = %+v", operation)
	}
	if _, err := repository.GetOperation("missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("GetOperation(missing) = %v, want ErrNotFound", err)
	}
	if got, err := repository.ListOperations(0); err != nil || len(got) != 0 {
		t.Fatalf("ListOperations(0) = %v (err %v), want empty", got, err)
	}
}

func TestRepositoryStructuredFilters(t *testing.T) {
	repository := openOperationRepository(t, "operationlog-structured.db")
	base := time.Date(2026, 8, 17, 8, 0, 0, 0, time.UTC)
	ops := []Operation{
		{ID: "op-a", Event: EventAuthLogin, ActorID: "user-admin", ActorName: "Admin", CreatedAt: base},
		{ID: "op-b", Event: EventUserCreate, ActorID: "user-admin", ActorName: "Admin", CreatedAt: base.Add(30 * time.Minute)},
		{ID: "op-c", Event: EventAuthLogout, ActorID: "user-editor", ActorName: "Editor", CreatedAt: base.Add(2 * time.Hour)},
	}
	for _, op := range ops {
		if err := repository.RecordOperation(op); err != nil {
			t.Fatalf("RecordOperation(%s): %v", op.ID, err)
		}
	}

	from := base.Add(15 * time.Minute)
	to := base.Add(90 * time.Minute)
	items, total, err := repository.ListOperationsFiltered(OperationFilter{
		Event: EventUserCreate, ActorName: "admin", From: &from, To: &to,
		Sort: "createdAt", Order: "asc", Page: 1, PageSize: 10,
	})
	if err != nil {
		t.Fatalf("structured filter: %v", err)
	}
	if total != 1 || len(items) != 1 || items[0].ID != "op-b" {
		t.Fatalf("structured filtered = %+v total=%d, want op-b/1", items, total)
	}

	// From-only: op-b and op-c.
	items, total, err = repository.ListOperationsFiltered(OperationFilter{
		From: &from, Sort: "createdAt", Order: "asc", Page: 1, PageSize: 10,
	})
	if err != nil {
		t.Fatalf("from filter: %v", err)
	}
	if total != 2 || items[0].ID != "op-b" || items[1].ID != "op-c" {
		t.Fatalf("from filtered = %+v total=%d, want op-b/op-c/2", items, total)
	}
}

func TestRepositoryRejectsUnknownEventAndExposesFailureSeam(t *testing.T) {
	repository := openOperationRepository(t, "operationlog-failure.db")
	bad := Operation{
		ID: "op-x", Event: "records.purge", ActorID: "user-admin", ActorName: "Admin", CreatedAt: time.Now().UTC(),
	}
	if err := repository.RecordOperation(bad); err == nil {
		t.Fatal("expected CHECK violation for unknown event")
	}

	forced := errors.New("forced operation log failure")
	repository.SetOperationLogError(forced)
	if err := repository.RecordOperation(Operation{
		ID: "op-forced", Event: EventAuthLogin, ActorID: "user-admin", ActorName: "Admin", CreatedAt: time.Now().UTC(),
	}); !errors.Is(err, forced) {
		t.Fatalf("forced error = %v, want wrapped seam error", err)
	}
	repository.SetOperationLogError(nil)
	if err := repository.RecordOperation(Operation{
		ID: "op-ok", Event: EventAuthLogin, ActorID: "user-admin", ActorName: "Admin", CreatedAt: time.Now().UTC(),
	}); err != nil {
		t.Fatalf("RecordOperation after clearing seam: %v", err)
	}
}
