package operationlog

import (
	"testing"
	"time"
)

func TestRecordOperationPersistsSessionAndCorrelation(t *testing.T) {
	repo := openOperationLogRepository(t, "session-side-table.db")
	now := time.Now().UTC()
	if err := repo.RecordOperation(Operation{
		ID: "op-sid", Event: EventSettingsUpdate, ActorID: "u1", ActorName: "Ada",
		CorrelationID: "corr-1", SessionID: "sess-1", CreatedAt: now,
	}); err != nil {
		t.Fatalf("record: %v", err)
	}
	got, err := repo.GetOperation("op-sid")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.SessionID != "sess-1" || got.CorrelationID != "corr-1" {
		t.Fatalf("got %+v", got)
	}
}
