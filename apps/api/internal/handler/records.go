package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
)

// maxRecordBodyBytes keeps the records write-body bound compatible with the
// frozen records contract (I-007-001 §2); it aliases the shared resource bound.
const maxRecordBodyBytes = maxResourceBodyBytes

// recordsResource describes the records instance of the generic resource
// factory (I-010-001 §4/§5): mounted at /api/records with the frozen
// sort/search/field set and records.read / records.write permission keys — zero
// change to the HTTP contract (I-007-001). Implementation converges from the
// hand-written handler to a registered entry + generic factory.
func recordsResource(st *store.Store) Resource {
	return Resource{
		ID:              "records",
		Path:            "/api/records",
		Listable:        true,
		SortFields:      []string{"name", "status", "owner", "updatedAt"},
		QSearch:         true,
		Entity:          &recordsEntity{st: st},
		CreateFields:    []string{"name", "status", "owner"},
		PatchFields:     []string{"name", "status", "owner"},
		PermissionRead:  "records.read",
		PermissionWrite: "records.write",
		// Legacy NOT_FOUND code kept for zero API change (I-010-001 §5).
		NotFoundCode: "RECORD_NOT_FOUND",
		NewID:        newRecordID,
		OnWrite:      recordsOnWrite(st),
	}
}

// recordsEntity adapts the concrete records store to the generic resource
// boundary. Rows are JSON maps; updatedAt is pre-formatted with fixed
// millisecond precision (GOAL-007 D-004).
type recordsEntity struct {
	st *store.Store
}

// recordToMap maps a persisted record to the API row. UpdatedAt serializes with
// the frozen fixed-3-digit-millisecond RFC3339 shape.
func recordToMap(r store.Record) map[string]any {
	return map[string]any{
		"id":        r.ID,
		"name":      r.Name,
		"status":    r.Status,
		"owner":     r.Owner,
		"updatedAt": r.UpdatedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
	}
}

func (e *recordsEntity) List(f resourceFilter) ([]map[string]any, int, error) {
	items, total, err := e.st.ListRecords(store.RecordFilter{
		Q: f.Q, Sort: f.Sort, Order: f.Order, Page: f.Page, PageSize: f.PageSize,
	})
	if err != nil {
		return nil, 0, err
	}
	out := make([]map[string]any, 0, len(items))
	for _, it := range items {
		out = append(out, recordToMap(it))
	}
	return out, total, nil
}

func (e *recordsEntity) Get(id string) (map[string]any, error) {
	rec, err := e.st.GetRecord(id)
	if err != nil {
		return nil, err
	}
	return recordToMap(*rec), nil
}

func (e *recordsEntity) Create(body map[string]any, id string, now time.Time) (map[string]any, error) {
	rec, err := e.st.CreateRecord(store.Record{
		ID:        id,
		Name:      stringField(body, "name"),
		Status:    stringField(body, "status"),
		Owner:     stringField(body, "owner"),
		UpdatedAt: now,
	})
	if err != nil {
		return nil, err
	}
	return recordToMap(*rec), nil
}

func (e *recordsEntity) Update(id string, body map[string]any, now time.Time) (map[string]any, error) {
	patch := store.RecordPatch{}
	if _, ok := body["name"]; ok {
		v := stringField(body, "name")
		patch.Name = &v
	}
	if _, ok := body["status"]; ok {
		v := stringField(body, "status")
		patch.Status = &v
	}
	if _, ok := body["owner"]; ok {
		v := stringField(body, "owner")
		patch.Owner = &v
	}
	rec, err := e.st.UpdateRecord(id, patch, now)
	if err != nil {
		return nil, err
	}
	return recordToMap(*rec), nil
}

func (e *recordsEntity) Delete(id string) error {
	return e.st.DeleteRecord(id)
}

// recordsOnWrite appends operation-log rows for records write endpoints
// (R5 S6 · I-008-003 §5). Best-effort: a logging failure is recorded to the
// service log and never turns a successful write into a failure.
func recordsOnWrite(st *store.Store) func(context.Context, account.User, writeKind, string, map[string]any, time.Time) {
	return func(_ context.Context, user account.User, kind writeKind, id string, row map[string]any, now time.Time) {
		event, detail := eventFor(kind, row)
		op := store.Operation{
			ID:        newOperationID(),
			Event:     event,
			ActorID:   user.ID,
			ActorName: user.Name,
			CreatedAt: now,
		}
		if id != "" {
			op.RecordID = &id
		}
		if detail != "" {
			op.Detail = &detail
		}
		if err := st.RecordOperation(op); err != nil {
			slog.Error("operation log write failed", "event", event, "err", err)
		}
	}
}

// eventFor maps a write kind to the frozen operation-log event + detail summary
// (I-008-003 §2/§3). Create/update detail carries the record name, never a secret.
func eventFor(kind writeKind, row map[string]any) (string, string) {
	switch kind {
	case writeCreate:
		return store.EventRecordCreate, `{"name":` + jsonQuote(stringField(row, "name")) + `}`
	case writeUpdate:
		return store.EventRecordUpdate, `{"name":` + jsonQuote(stringField(row, "name")) + `}`
	default:
		return store.EventRecordDelete, ""
	}
}

// newRecordID returns "rec-" + 16 lowercase hex chars (8 bytes of crypto/rand),
// the frozen create id format (I-007-001 §2).
func newRecordID() (string, error) {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return "rec-" + hex.EncodeToString(b[:]), nil
}

// jsonQuote returns s as a JSON string literal (used for operation log detail
// summaries; never used for secrets, which are excluded by I-008-003 §3).
func jsonQuote(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}

// newOperationID returns a random 128-bit hex id for operation log rows.
func newOperationID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		// crypto/rand failure is effectively fatal; fall back to a timestamp id
		// so logging never wedges a successful request (best-effort contract).
		return fmt.Sprintf("op-%d", time.Now().UnixNano())
	}
	return "op-" + hex.EncodeToString(b[:])
}
