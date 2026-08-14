// Package store tests for the admin.recycle-bin persistence (S-12 ·
// GOAL-012 D-002 §1/§3): snapshot insert/unique/list/restore-mark/purge.
package store

import (
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)

func newTestEnv(t *testing.T) *Repository {
	t.Helper()
	st, err := testsupport.OpenStore(filepath.Join(t.TempDir(), "test.db"), "admin", "hash", true)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return NewRepository(st)
}

func sampleItem(id string) Item {
	return Item{
		ID:         id,
		Resource:   "dict-types",
		ResourceID: "t1",
		Payload:    map[string]any{"id": "t1", "key": "status", "name": "Status"},
		ActorID:    "user-admin",
		ActorName:  "Admin",
		DeletedAt:  time.Now().UTC(),
	}
}

func TestRecordAndGet(t *testing.T) {
	r := newTestEnv(t)
	item := sampleItem("recycle-1")
	if err := r.Record(item); err != nil {
		t.Fatalf("record: %v", err)
	}
	got, err := r.Get("recycle-1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Resource != "dict-types" || got.ResourceID != "t1" || got.Payload["key"] != "status" {
		t.Fatalf("item = %+v", got)
	}
	if got.RestoredAt != nil {
		t.Fatalf("item must start unrestored: %+v", got)
	}
}

func TestRecordUniqueWhileActive(t *testing.T) {
	r := newTestEnv(t)
	if err := r.Record(sampleItem("recycle-1")); err != nil {
		t.Fatalf("record 1: %v", err)
	}
	if err := r.Record(sampleItem("recycle-2")); err == nil {
		t.Fatal("duplicate active snapshot must fail (partial unique)")
	}
	// After restore the slot frees.
	if err := r.MarkRestored("recycle-1", time.Now().UTC()); err != nil {
		t.Fatalf("mark restored: %v", err)
	}
	if err := r.Record(sampleItem("recycle-3")); err != nil {
		t.Fatalf("record after restore: %v", err)
	}
}

func TestListFiltersAndOrdering(t *testing.T) {
	r := newTestEnv(t)
	now := time.Now().UTC()
	one := sampleItem("recycle-1")
	one.DeletedAt = now.Add(-time.Hour)
	two := sampleItem("recycle-2")
	two.Resource = "scheduled-tasks"
	two.ResourceID = "task-a"
	two.DeletedAt = now
	if err := r.Record(one); err != nil {
		t.Fatalf("record 1: %v", err)
	}
	if err := r.Record(two); err != nil {
		t.Fatalf("record 2: %v", err)
	}
	items, total, err := r.List(ListFilter{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != 2 || len(items) != 2 {
		t.Fatalf("total = %d items = %d, want 2", total, len(items))
	}
	if items[0].ResourceID != "task-a" {
		t.Fatalf("default order must be newest first: %+v", items[0])
	}
	byRes, _, err := r.List(ListFilter{Resource: "dict-types", Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("list by resource: %v", err)
	}
	if len(byRes) != 1 || byRes[0].ResourceID != "t1" {
		t.Fatalf("resource filter = %+v", byRes)
	}
	// Restored items drop out of the active list.
	if err := r.MarkRestored("recycle-1", now); err != nil {
		t.Fatalf("mark restored: %v", err)
	}
	active, total, err := r.List(ListFilter{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("list active: %v", err)
	}
	if total != 1 || len(active) != 1 || active[0].ResourceID != "task-a" {
		t.Fatalf("active = %d/%d %+v, want only task-a", len(active), total, active)
	}
}

func TestPurge(t *testing.T) {
	r := newTestEnv(t)
	if err := r.Record(sampleItem("recycle-1")); err != nil {
		t.Fatalf("record: %v", err)
	}
	if err := r.Purge("recycle-1"); err != nil {
		t.Fatalf("purge: %v", err)
	}
	if _, err := r.Get("recycle-1"); !errors.Is(err, ErrItemNotFound) {
		t.Fatalf("get after purge = %v, want ErrItemNotFound", err)
	}
	if err := r.Purge("recycle-1"); !errors.Is(err, ErrItemNotFound) {
		t.Fatalf("purge again = %v, want ErrItemNotFound", err)
	}
}

func TestMarkRestoredTwiceFails(t *testing.T) {
	r := newTestEnv(t)
	if err := r.Record(sampleItem("recycle-1")); err != nil {
		t.Fatalf("record: %v", err)
	}
	now := time.Now().UTC()
	if err := r.MarkRestored("recycle-1", now); err != nil {
		t.Fatalf("mark restored: %v", err)
	}
	if err := r.MarkRestored("recycle-1", now); !errors.Is(err, ErrItemAlreadyRestored) {
		t.Fatalf("second mark = %v, want ErrItemAlreadyRestored", err)
	}
}
