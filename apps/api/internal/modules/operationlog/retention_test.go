package operationlog

import (
	"context"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	"path/filepath"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)

func openOperationLogRepository(t *testing.T, name string) *Repository {
	t.Helper()
	platform, err := testsupport.OpenStore(filepath.Join(t.TempDir(), name), "admin", "hash", true)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = platform.Close() })
	return NewRepository(platform)
}

func TestApplyRetentionArchivesThenRemovesHotRows(t *testing.T) {
	repo := openOperationLogRepository(t, "retention-archive.db")
	now := time.Date(2026, 8, 19, 12, 0, 0, 0, time.UTC)
	old := now.Add(-40 * 24 * time.Hour)
	fresh := now.Add(-2 * 24 * time.Hour)
	if err := repo.RecordOperation(Operation{
		ID: "old-1", Event: EventSettingsUpdate, ActorID: "u1", ActorName: "Ada",
		CorrelationID: "corr-old", SessionID: "sess-old", CreatedAt: old,
	}); err != nil {
		t.Fatalf("record old: %v", err)
	}
	if err := repo.RecordOperation(Operation{
		ID: "fresh-1", Event: EventSettingsUpdate, ActorID: "u1", ActorName: "Ada",
		CreatedAt: fresh,
	}); err != nil {
		t.Fatalf("record fresh: %v", err)
	}

	n, err := repo.ApplyRetention(now, 30, "archive")
	if err != nil {
		t.Fatalf("apply: %v", err)
	}
	if n != 1 {
		t.Fatalf("expired = %d, want 1", n)
	}
	if _, err := repo.GetOperation("old-1"); err == nil {
		t.Fatal("old row still in hot table")
	}
	freshRow, err := repo.GetOperation("fresh-1")
	if err != nil || freshRow == nil {
		t.Fatalf("fresh row: %v", err)
	}

	var archivedID, archivedCorr, archivedSession string
	if err := repo.runner.Run(context.Background(), func(tx kernel.Tx) error {
		if err := tx.QueryRow(context.Background(), `SELECT id FROM operation_log_archive WHERE id = 'old-1'`).Scan(&archivedID); err != nil {
			return err
		}
		if err := tx.QueryRow(context.Background(), `SELECT correlation_id FROM operation_log_archive_correlation WHERE operation_id = 'old-1'`).Scan(&archivedCorr); err != nil {
			return err
		}
		return tx.QueryRow(context.Background(), `SELECT session_id FROM operation_log_archive_session WHERE operation_id = 'old-1'`).Scan(&archivedSession)
	}); err != nil {
		t.Fatalf("archive lookup: %v", err)
	}
	if archivedID != "old-1" || archivedCorr != "corr-old" || archivedSession != "sess-old" {
		t.Fatalf("archive = %s %s %s", archivedID, archivedCorr, archivedSession)
	}
}

func TestApplyRetentionDeleteDoesNotArchive(t *testing.T) {
	repo := openOperationLogRepository(t, "retention-delete.db")
	now := time.Date(2026, 8, 19, 12, 0, 0, 0, time.UTC)
	if err := repo.RecordOperation(Operation{
		ID: "gone", Event: EventSettingsUpdate, ActorID: "u1", ActorName: "Ada",
		CreatedAt: now.Add(-10 * 24 * time.Hour),
	}); err != nil {
		t.Fatalf("record: %v", err)
	}
	n, err := repo.ApplyRetention(now, 7, "delete")
	if err != nil {
		t.Fatalf("apply: %v", err)
	}
	if n != 1 {
		t.Fatalf("expired = %d", n)
	}
	var count int
	if err := repo.runner.Run(context.Background(), func(tx kernel.Tx) error {
		return tx.QueryRow(context.Background(), `SELECT COUNT(*) FROM operation_log_archive`).Scan(&count)
	}); err != nil {
		t.Fatalf("count archive: %v", err)
	}
	if count != 0 {
		t.Fatalf("delete action archived %d rows", count)
	}
}
