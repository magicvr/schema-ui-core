// Service tests for the recycle bin (S-12 · GOAL-012 D-002 §2/§3): the
// TrashRecorder surface, restore dispatch per resource, conflict mapping and
// purge.
package recyclebin

import (
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	datadictionarystore "github.com/magicvr/schema-ui-core/apps/api/internal/modules/datadictionary/store"
	recyclestore "github.com/magicvr/schema-ui-core/apps/api/internal/modules/recyclebin/store"
	tasksstore "github.com/magicvr/schema-ui-core/apps/api/internal/modules/scheduledtasks/store"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)

func newServiceEnv(t *testing.T) *Service {
	t.Helper()
	st, err := testsupport.OpenStore(filepath.Join(t.TempDir(), "test.db"), "admin", "hash", true)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return NewService(recyclestore.NewRepository(st), datadictionarystore.NewRepository(st), tasksstore.NewRepository(st))
}

func TestRecordWritesSnapshot(t *testing.T) {
	s := newServiceEnv(t)
	now := time.Now().UTC()
	actor := account.User{ID: "user-admin", Name: "Admin"}
	if err := s.Record(t.Context(), "dict-types", "t1", map[string]any{"key": "status"}, actor, now); err != nil {
		t.Fatalf("record: %v", err)
	}
	item, err := s.Get("recycle-" + hexID(now))
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if item.Resource != "dict-types" || item.ResourceID != "t1" || item.ActorName != "Admin" {
		t.Fatalf("item = %+v", item)
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
