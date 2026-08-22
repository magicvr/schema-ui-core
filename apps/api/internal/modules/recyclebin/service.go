// Recycle-bin service (S-12 · GOAL-012 D-002 §1/§2/§3): implements the
// handler.TrashRecorder surface consumed by the resource factory delete hooks,
// and the per-resource restore dispatch (payload → owning store Create).
package recyclebin

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/handler"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/datadictionary/store"
	recyclestore "github.com/magicvr/schema-ui-core/apps/api/internal/modules/recyclebin/store"
	tasksstore "github.com/magicvr/schema-ui-core/apps/api/internal/modules/scheduledtasks/store"
)

// ErrRestoreConflict is returned when a restore would violate a unique key
// (the snapshot is kept so the item can be retried after the conflict clears).
var ErrRestoreConflict = errors.New("recycle restore conflict")

// Service owns the recycle-bin domain behavior.
type Service struct {
	repository *recyclestore.Repository
	dictionary *store.Repository
	tasks      *tasksstore.Repository
	// runner is the shared transaction boundary used by Restore so the
	// business INSERT and the snapshot MarkRestored commit atomically
	// (W11 F-008); the kernel store satisfies it structurally.
	runner recyclestore.TxRunner
	now    func() time.Time
}

// NewService constructs the recycle service. runner is the shared
// transaction boundary (the kernel store); it must be the SAME backing store
// as the three repositories.
func NewService(repository *recyclestore.Repository, dictionary *store.Repository, tasks *tasksstore.Repository, runner recyclestore.TxRunner) *Service {
	return &Service{repository: repository, dictionary: dictionary, tasks: tasks, runner: runner, now: time.Now}
}

// Record implements handler.TrashRecorder (S-12 · GOAL-012 D-002 §2): called
// by the resource factory after a successful delete. Snapshot ids are
// "recycle-" + 16 random hex (crypto/rand) so a batch delete recording several
// rows in the same second can never collide on the primary key
// (grok A-003 F-001/F-005).
func (s *Service) Record(_ context.Context, resource, id string, row map[string]any, actor account.User, now time.Time) error {
	snapshotID, err := newSnapshotID()
	if err != nil {
		return err
	}
	return s.repository.Record(recyclestore.Item{
		ID:         snapshotID,
		Resource:   resource,
		ResourceID: id,
		Payload:    row,
		ActorID:    actor.ID,
		ActorName:  actor.Name,
		DeletedAt:  now,
	})
}

// RecordTx implements handler.TrashTxRecorder (W11 F-002): writes the
// snapshot INSIDE the caller's transaction so the factory's delete and the
// recycle snapshot commit atomically; a snapshot failure rolls the delete
// back instead of committing a delete without a snapshot.
func (s *Service) RecordTx(ctx context.Context, tx kernel.Tx, resource, id string, row map[string]any, actor account.User, now time.Time) error {
	snapshotID, err := newSnapshotID()
	if err != nil {
		return err
	}
	return s.repository.RecordTx(ctx, tx, recyclestore.Item{
		ID:         snapshotID,
		Resource:   resource,
		ResourceID: id,
		Payload:    row,
		ActorID:    actor.ID,
		ActorName:  actor.Name,
		DeletedAt:  now,
	})
}

// Restore re-creates the deleted row in its owning store. The snapshot is kept
// on conflict so the item can be retried (D-002 §3).
//
// W11 F-008: the business INSERT and the MarkRestored now commit in ONE
// transaction — the previous shape committed the restored row first, then
// marked the snapshot; a crash between the two left a live row AND a
// still-restorable snapshot (restores would then conflict, or worse,
// duplicate a dict entry).
func (s *Service) Restore(itemID string, now time.Time) (map[string]any, error) {
	item, err := s.repository.Get(itemID)
	if err != nil {
		return nil, err
	}
	if item.RestoredAt != nil {
		return nil, recyclestore.ErrItemAlreadyRestored
	}
	err = s.runner.Run(context.Background(), func(tx kernel.Tx) error {
		if err := s.restoreRowTx(tx, item.Resource, item.Payload, now); err != nil {
			return err
		}
		return s.repository.MarkRestoredTx(context.Background(), tx, itemID, now)
	})
	if err != nil {
		if errors.Is(err, store.ErrDictKeyNotFound) {
			// W6 F2 (GOAL-006 D-001): restoring an orphaned dict entry whose
			// parent dict type was deleted is a recoverable precondition
			// failure, not an internal error — surface it as a clear 409 so
			// the operator can restore the parent type first and retry. The
			// snapshot stays untouched (retryable).
			return nil, &handler.DomainError{Status: http.StatusConflict, Code: "DICT_KEY_NOT_FOUND", Message: "parent dict type does not exist"}
		}
		if isConflict(err) {
			return nil, &handler.DomainError{Status: http.StatusConflict, Code: "RECYCLE_RESTORE_CONFLICT", Message: "a row with that key already exists"}
		}
		return nil, err
	}
	return item.Payload, nil
}

// Purge physically removes the snapshot (irreversible, D-002 §3).
func (s *Service) Purge(itemID string) error {
	return s.repository.Purge(itemID)
}

// PurgeAll physically removes every active snapshot (D-002 §4 batch purge).
func (s *Service) PurgeAll() (int, error) {
	return s.repository.PurgeAllUnrestored()
}

// List returns active snapshots.
func (s *Service) List(filter recyclestore.ListFilter) ([]recyclestore.Item, int, error) {
	return s.repository.List(filter)
}

// Get loads one snapshot.
func (s *Service) Get(itemID string) (*recyclestore.Item, error) {
	return s.repository.Get(itemID)
}

// ListItems adapts the store list to the handler surface (S-12 · GOAL-012 §3).
func (s *Service) ListItems(resource, q, sortField, order string, page, pageSize int) ([]handler.RecycleItem, int, error) {
	items, total, err := s.repository.List(recyclestore.ListFilter{Resource: resource, Q: q, Sort: sortField, Order: order, Page: page, PageSize: pageSize})
	if err != nil {
		return nil, 0, err
	}
	out := make([]handler.RecycleItem, 0, len(items))
	for _, item := range items {
		out = append(out, toHandlerItem(item))
	}
	return out, total, nil
}

// GetItem adapts the store item to the handler surface.
func (s *Service) GetItem(itemID string) (*handler.RecycleItem, error) {
	item, err := s.repository.Get(itemID)
	if err != nil {
		return nil, err
	}
	h := toHandlerItem(*item)
	return &h, nil
}

func toHandlerItem(item recyclestore.Item) handler.RecycleItem {
	out := handler.RecycleItem{
		ID:         item.ID,
		Resource:   item.Resource,
		ResourceID: item.ResourceID,
		Payload:    item.Payload,
		ActorID:    item.ActorID,
		ActorName:  item.ActorName,
		DeletedAt:  item.DeletedAt,
	}
	if item.RestoredAt != nil {
		out.RestoredAt = *item.RestoredAt
	}
	return out
}

func (s *Service) restoreRowTx(tx kernel.Tx, resource string, payload map[string]any, now time.Time) error {
	switch resource {
	case "dict-types":
		return s.dictionary.CreateTypeTx(context.Background(), tx, dictTypeFromPayload(payload, now))
	case "dict-entries":
		return s.dictionary.CreateEntryTx(context.Background(), tx, dictEntryFromPayload(payload, now))
	case "scheduled-tasks":
		return s.tasks.CreateTaskTx(context.Background(), tx, taskFromPayload(payload, now))
	default:
		return fmt.Errorf("recycle restore: unsupported resource %q", resource)
	}
}

func isConflict(err error) bool {
	return errors.Is(err, store.ErrTypeKeyTaken) || errors.Is(err, store.ErrEntryKeyTaken) || errors.Is(err, tasksstore.ErrKeyTaken)
}

func dictTypeFromPayload(payload map[string]any, now time.Time) store.DictType {
	return store.DictType{
		ID:          stringField(payload, "id"),
		Key:         stringField(payload, "key"),
		Name:        stringField(payload, "name"),
		Enabled:     boolField(payload, "enabled"),
		Description: stringField(payload, "description"),
		Sort:        intField(payload, "sort"),
		CreatedAt:   timeField(payload, "createdAt", now),
		UpdatedAt:   timeField(payload, "updatedAt", now),
	}
}

func dictEntryFromPayload(payload map[string]any, now time.Time) store.DictEntry {
	return store.DictEntry{
		ID:        stringField(payload, "id"),
		DictKey:   stringField(payload, "dictKey"),
		EntryKey:  stringField(payload, "entryKey"),
		Label:     stringField(payload, "label"),
		Enabled:   boolField(payload, "enabled"),
		Sort:      intField(payload, "sort"),
		Remark:    stringField(payload, "remark"),
		BadgeStyle: stringField(payload, "badgeStyle"),
		CreatedAt: timeField(payload, "createdAt", now),
		UpdatedAt: timeField(payload, "updatedAt", now),
	}
}

func taskFromPayload(payload map[string]any, now time.Time) tasksstore.Task {
	return tasksstore.Task{
		ID:          stringField(payload, "id"),
		Key:         stringField(payload, "key"),
		Cron:        stringField(payload, "cron"),
		Name:        stringField(payload, "name"),
		Enabled:     boolField(payload, "enabled"),
		Description: stringField(payload, "description"),
		Handler:     stringField(payload, "handler"),
		CreatedAt:   timeField(payload, "createdAt", now),
		UpdatedAt:   timeField(payload, "updatedAt", now),
	}
}

func stringField(m map[string]any, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func boolField(m map[string]any, key string) bool {
	if v, ok := m[key].(bool); ok {
		return v
	}
	return false
}

func intField(m map[string]any, key string) int {
	switch v := m[key].(type) {
	case float64:
		return int(v)
	case int64:
		return int(v)
	case int:
		return v
	}
	return 0
}

func timeField(m map[string]any, key string, fallback time.Time) time.Time {
	switch v := m[key].(type) {
	case float64:
		return time.Unix(int64(v), 0).UTC()
	case int64:
		return time.Unix(v, 0).UTC()
	case string:
		if t, err := time.Parse("2006-01-02T15:04:05.000Z07:00", v); err == nil {
			return t.UTC()
		}
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			return t.UTC()
		}
	}
	return fallback
}

// newSnapshotID returns "recycle-" + 16 random hex bytes (crypto/rand).
func newSnapshotID() (string, error) {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return "recycle-" + hex.EncodeToString(b[:]), nil
}
