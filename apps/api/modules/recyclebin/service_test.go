// Service tests for the recycle bin (S-12 · GOAL-012 D-002 §2/§3): the
// TrashRecorder surface, restore dispatch per resource, conflict mapping and
// purge.
package recyclebin

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/handler"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	datadictionarystore "github.com/magicvr/schema-ui-core/apps/api/modules/datadictionary/store"
	recyclestore "github.com/magicvr/schema-ui-core/apps/api/modules/recyclebin/store"
	tasksstore "github.com/magicvr/schema-ui-core/apps/api/modules/scheduledtasks/store"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)

func newServiceEnv(t *testing.T) *Service {
	t.Helper()
	st, err := testsupport.OpenStore(filepath.Join(t.TempDir(), "test.db"), "admin", "hash", true)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return NewService(recyclestore.NewRepository(st), datadictionarystore.NewRepository(st), tasksstore.NewRepository(st), st)
}

func TestRecordWritesSnapshot(t *testing.T) {
	s := newServiceEnv(t)
	now := time.Now().UTC()
	actor := account.User{ID: "user-admin", Name: "Admin"}
	if err := s.Record(t.Context(), "dict-types", "t1", map[string]any{"key": "status"}, actor, now); err != nil {
		t.Fatalf("record: %v", err)
	}
	items, _, err := s.List(recyclestore.ListFilter{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("items = %d, want 1", len(items))
	}
	if items[0].Resource != "dict-types" || items[0].ResourceID != "t1" || items[0].ActorName != "Admin" {
		t.Fatalf("item = %+v", items[0])
	}
	if !strings.HasPrefix(items[0].ID, "recycle-") {
		t.Fatalf("snapshot id %q must be recycle- prefixed (D-002 §1)", items[0].ID)
	}
}

// F-001 (grok A-003): two records in the same second must not collide on the
// primary key (the batch-delete path records several snapshots with one now).
func TestRecordDistinctIDsSameSecond(t *testing.T) {
	s := newServiceEnv(t)
	now := time.Now().UTC()
	actor := account.User{ID: "user-admin", Name: "Admin"}
	for i := 0; i < 5; i++ {
		if err := s.Record(t.Context(), "dict-types", fmt.Sprintf("t%d", i), map[string]any{"key": fmt.Sprintf("k%d", i)}, actor, now); err != nil {
			t.Fatalf("record %d: %v", i, err)
		}
	}
	items, total, err := s.List(recyclestore.ListFilter{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != 5 || len(items) != 5 {
		t.Fatalf("items = %d/%d, want 5 distinct snapshots (F-001)", len(items), total)
	}
	seen := map[string]bool{}
	for _, item := range items {
		if seen[item.ID] {
			t.Fatalf("duplicate snapshot id %q (F-001)", item.ID)
		}
		seen[item.ID] = true
	}
}

func TestRestoreDictTypeAndConflict(t *testing.T) {
	s := newServiceEnv(t)
	now := time.Now().UTC()
	actor := account.User{ID: "user-admin", Name: "Admin"}
	// First create the type, delete it (snapshot), then restore.
	if err := s.dictionary.CreateType(datadictionarystore.DictType{ID: "t1", Key: "status", Name: "Status", Enabled: true, CreatedAt: now, UpdatedAt: now}); err != nil {
		t.Fatalf("create type: %v", err)
	}
	if err := s.Record(t.Context(), "dict-types", "t1", map[string]any{
		"id": "t1", "key": "status", "name": "Status", "enabled": true, "description": "", "sort": 0,
		"createdAt": float64(now.Unix()), "updatedAt": float64(now.Unix()),
	}, actor, now); err != nil {
		t.Fatalf("record: %v", err)
	}
	// The row must be gone from the source before restore (the snapshot exists
	// because the resource factory records it after a successful delete).
	if _, err := s.dictionary.DeleteType("t1"); err != nil {
		t.Fatalf("delete type: %v", err)
	}
	items, _, err := s.List(recyclestore.ListFilter{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("items = %d, want 1", len(items))
	}
	row, err := s.Restore(items[0].ID, now)
	if err != nil {
		t.Fatalf("restore: %v", err)
	}
	if row["key"] != "status" {
		t.Fatalf("restored row = %v", row)
	}
	// Restore again → conflict (row exists again).
	items2, _, _ := s.List(recyclestore.ListFilter{Page: 1, PageSize: 10})
	_ = items2
	if _, err := s.Restore(items[0].ID, now); !errors.Is(err, recyclestore.ErrItemAlreadyRestored) {
		t.Fatalf("second restore = %v, want ErrItemAlreadyRestored", err)
	}
}

func TestRestoreTaskRoundTrip(t *testing.T) {
	s := newServiceEnv(t)
	now := time.Now().UTC()
	actor := account.User{ID: "user-admin", Name: "Admin"}
	if err := s.tasks.CreateTask(tasksstore.Task{ID: "task-1", Key: "hourly", Cron: "0 * * * *", Name: "Hourly", Enabled: true, CreatedAt: now, UpdatedAt: now}); err != nil {
		t.Fatalf("create task: %v", err)
	}
	if err := s.Record(t.Context(), "scheduled-tasks", "task-1", map[string]any{
		"id": "task-1", "key": "hourly", "cron": "0 * * * *", "name": "Hourly", "enabled": true,
		"description": "", "handler": "system.noop",
		"createdAt": float64(now.Unix()), "updatedAt": float64(now.Unix()),
	}, actor, now); err != nil {
		t.Fatalf("record: %v", err)
	}
	if err := s.tasks.DeleteTask("task-1"); err != nil {
		t.Fatalf("delete task: %v", err)
	}
	items, _, _ := s.List(recyclestore.ListFilter{Page: 1, PageSize: 10})
	row, err := s.Restore(items[0].ID, now)
	if err != nil {
		t.Fatalf("restore: %v", err)
	}
	if row["key"] != "hourly" {
		t.Fatalf("restored row = %v", row)
	}
	if _, err := s.tasks.GetTask("task-1"); err != nil {
		t.Fatalf("task not restored: %v", err)
	}
}

func TestPurgeRemovesSnapshot(t *testing.T) {
	s := newServiceEnv(t)
	now := time.Now().UTC()
	actor := account.User{ID: "user-admin", Name: "Admin"}
	if err := s.Record(t.Context(), "dict-types", "t9", map[string]any{"key": "x"}, actor, now); err != nil {
		t.Fatalf("record: %v", err)
	}
	items, _, _ := s.List(recyclestore.ListFilter{Page: 1, PageSize: 10})
	if err := s.Purge(items[0].ID); err != nil {
		t.Fatalf("purge: %v", err)
	}
	if _, err := s.Get(items[0].ID); !errors.Is(err, recyclestore.ErrItemNotFound) {
		t.Fatalf("get after purge = %v", err)
	}
}

// F-008 (grok A-003): dict-entries restore round trip — the third managed
// resource must restore through CreateEntry.
func TestRestoreDictEntryRoundTrip(t *testing.T) {
	s := newServiceEnv(t)
	now := time.Now().UTC()
	actor := account.User{ID: "user-admin", Name: "Admin"}
	if err := s.dictionary.CreateType(datadictionarystore.DictType{ID: "t1", Key: "status", Name: "Status", Enabled: true, CreatedAt: now, UpdatedAt: now}); err != nil {
		t.Fatalf("create type: %v", err)
	}
	if err := s.dictionary.CreateEntry(datadictionarystore.DictEntry{ID: "e1", DictKey: "status", EntryKey: "active", Label: "Active", Enabled: true, CreatedAt: now, UpdatedAt: now}); err != nil {
		t.Fatalf("create entry: %v", err)
	}
	if err := s.Record(t.Context(), "dict-entries", "e1", map[string]any{
		"id": "e1", "dictKey": "status", "entryKey": "active", "label": "Active", "enabled": true,
		"sort": 0, "remark": "", "createdAt": float64(now.Unix()), "updatedAt": float64(now.Unix()),
	}, actor, now); err != nil {
		t.Fatalf("record: %v", err)
	}
	if err := s.dictionary.DeleteEntry("e1"); err != nil {
		t.Fatalf("delete entry: %v", err)
	}
	items, _, _ := s.List(recyclestore.ListFilter{Page: 1, PageSize: 10})
	row, err := s.Restore(items[0].ID, now)
	if err != nil {
		t.Fatalf("restore entry: %v", err)
	}
	if row["entryKey"] != "active" {
		t.Fatalf("restored entry = %v", row)
	}
	if _, err := s.dictionary.GetEntry("e1"); err != nil {
		t.Fatalf("entry not restored: %v", err)
	}
}

// W6 F2 (GOAL-006 D-001): restoring an orphaned dict entry whose parent dict
// type no longer exists returns a typed DomainError (409 DICT_KEY_NOT_FOUND)
// instead of an internal error, and the snapshot is kept for a retry.
func TestRestoreOrphanDictEntryReturnsDomainError(t *testing.T) {
	s := newServiceEnv(t)
	now := time.Now().UTC()
	actor := account.User{ID: "user-admin", Name: "Admin"}
	// Parent type exists while the entry is created and snapshot-recorded.
	if err := s.dictionary.CreateType(datadictionarystore.DictType{ID: "t1", Key: "status", Name: "Status", Enabled: true, CreatedAt: now, UpdatedAt: now}); err != nil {
		t.Fatalf("create type: %v", err)
	}
	if err := s.dictionary.CreateEntry(datadictionarystore.DictEntry{ID: "e1", DictKey: "status", EntryKey: "active", Label: "Active", Enabled: true, CreatedAt: now, UpdatedAt: now}); err != nil {
		t.Fatalf("create entry: %v", err)
	}
	if err := s.Record(t.Context(), "dict-entries", "e1", map[string]any{
		"id": "e1", "dictKey": "status", "entryKey": "active", "label": "Active", "enabled": true,
		"sort": 0, "remark": "", "createdAt": float64(now.Unix()), "updatedAt": float64(now.Unix()),
	}, actor, now); err != nil {
		t.Fatalf("record: %v", err)
	}
	if err := s.dictionary.DeleteEntry("e1"); err != nil {
		t.Fatalf("delete entry: %v", err)
	}
	// Parent type is deleted before the restore is attempted → orphan.
	if _, err := s.dictionary.DeleteType("t1"); err != nil {
		t.Fatalf("delete type: %v", err)
	}
	items, _, _ := s.List(recyclestore.ListFilter{Page: 1, PageSize: 10})
	_, err := s.Restore(items[0].ID, now)
	var de *handler.DomainError
	if !errors.As(err, &de) {
		t.Fatalf("restore orphan = %v, want DomainError (DICT_KEY_NOT_FOUND)", err)
	}
	if de.Status != http.StatusConflict || de.Code != "DICT_KEY_NOT_FOUND" {
		t.Fatalf("domain error = %d %q, want 409 DICT_KEY_NOT_FOUND", de.Status, de.Code)
	}
	// The snapshot must be kept (retryable after the parent type is restored).
	items2, total, _ := s.List(recyclestore.ListFilter{Page: 1, PageSize: 10})
	if total != 1 || items2[0].ID != items[0].ID {
		t.Fatalf("snapshot after failed restore = %d items (%v), want the original kept", total, items2)
	}
}

// W11 F-008: the restore business INSERT and the snapshot MarkRestored now
// share ONE transaction (Restore runs both on a caller-owned kernel.Tx). A
// failed mark rolls the restored row back — the previous shape committed the
// row first and a crash between the two left "row exists AND snapshot still
// restorable" (restore would then conflict or duplicate).
func TestRestoreAtomicityRollsBackOnFailedMark(t *testing.T) {
	st, err := testsupport.OpenStore(filepath.Join(t.TempDir(), "test.db"), "admin", "hash", true)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	dictionary := datadictionarystore.NewRepository(st)
	rc := recyclestore.NewRepository(st)
	now := time.Now().UTC()

	// A live snapshot for a dict type that does not exist yet.
	if err := rc.Record(recyclestore.Item{
		ID: "recycle-atomic-1", Resource: "dict-types", ResourceID: "t-atomic",
		Payload: map[string]any{"id": "t-atomic", "key": "atomic", "name": "Atomic", "enabled": true, "sort": 0},
		DeletedAt: now,
	}); err != nil {
		t.Fatalf("record snapshot: %v", err)
	}

	// Same-tx seam: the restore INSERT succeeds but MarkRestoredTx targets a
	// wrong id → the whole transaction rolls back.
	err = st.Run(context.Background(), func(tx kernel.Tx) error {
		if err := dictionary.CreateTypeTx(context.Background(), tx, datadictionarystore.DictType{
			ID: "t-atomic", Key: "atomic", Name: "Atomic", Enabled: true, CreatedAt: now, UpdatedAt: now,
		}); err != nil {
			return err
		}
		return rc.MarkRestoredTx(context.Background(), tx, "recycle-wrong-id", now)
	})
	if err == nil || !strings.Contains(err.Error(), "recycle item already restored") {
		t.Fatalf("seam err = %v, want ErrItemAlreadyRestored from the failed mark", err)
	}
	// Nothing committed: no dict type row, snapshot still unrestored.
	var count int
	if err := st.Run(context.Background(), func(tx kernel.Tx) error {
		return tx.QueryRow(context.Background(), `SELECT COUNT(*) FROM dict_types WHERE id = 't-atomic'`).Scan(&count)
	}); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("dict type rows = %d, want 0 (restore rolled back)", count)
	}
	item, err := rc.Get("recycle-atomic-1")
	if err != nil {
		t.Fatal(err)
	}
	if item.RestoredAt != nil {
		t.Fatalf("snapshot marked restored after rollback, want unrestored")
	}
}
