// Data dictionary surface (S-01 · GOAL-008 D-002 §3/§5): the
// admin.data-dictionary module exposes two schema-driven resources — dict
// types and dict entries — with dictionary.read / dictionary.write gates
// (admin-only) and dictionary.create/update/delete audit events.
package handler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	datadictionarystore "github.com/magicvr/schema-ui-core/apps/api/internal/modules/datadictionary/store"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
)

// DictionaryRepository is the persistence boundary consumed by the two
// dictionary resources.
type DictionaryRepository interface {
	ListTypes(datadictionarystore.ListFilter) ([]datadictionarystore.DictType, int, error)
	GetType(string) (*datadictionarystore.DictType, error)
	CreateType(datadictionarystore.DictType) error
	UpdateType(id, name string, enabled bool, description string, sort int, now time.Time) error
	DeleteType(string) ([]string, error)
	ListEntries(datadictionarystore.ListFilter) ([]datadictionarystore.DictEntry, int, error)
	GetEntry(string) (*datadictionarystore.DictEntry, error)
	CreateEntry(datadictionarystore.DictEntry) error
	UpdateEntry(id, dictKey, label string, enabled bool, sort int, remark string, now time.Time) error
	DeleteEntry(string) error
}

// dictTypeEntity adapts the dictionary repository to the generic factory.
type dictTypeEntity struct {
	repository DictionaryRepository
	operations operationlog.Recorder
}

func (e *dictTypeEntity) List(filter resourceFilter) ([]map[string]any, int, error) {
	rows, total, err := e.repository.ListTypes(datadictionarystore.ListFilter{Q: filter.Q, Sort: filter.Sort, Order: filter.Order, Page: filter.Page, PageSize: filter.PageSize})
	if err != nil {
		return nil, 0, err
	}
	items := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		items = append(items, dictTypeToMap(row))
	}
	return items, total, nil
}

func (e *dictTypeEntity) Get(id string) (map[string]any, error) {
	row, err := e.repository.GetType(id)
	if err != nil {
		return nil, mapDictStoreError(err)
	}
	return dictTypeToMap(*row), nil
}

func (e *dictTypeEntity) Create(body map[string]any, id string, now time.Time, actor account.User) (map[string]any, error) {
	key := stringField(body, "key")
	name := stringField(body, "name")
	if key == "" || name == "" {
		return nil, &DomainError{Status: http.StatusBadRequest, Code: "INVALID_CREATE_FIELD", Message: "key and name are required"}
	}
	err := e.repository.CreateType(datadictionarystore.DictType{
		ID: id, Key: key, Name: name, Enabled: boolField(body, "enabled", true),
		Description: stringField(body, "description"), Sort: intField(body, "sort"),
		CreatedAt: now, UpdatedAt: now,
	})
	if err != nil {
		return nil, mapDictStoreError(err)
	}
	recordDictionaryEvent(e.operations, operationlog.EventDictionaryCreate, actor, id, now)
	created, err := e.repository.GetType(id)
	if err != nil {
		return nil, err
	}
	return dictTypeToMap(*created), nil
}

func (e *dictTypeEntity) Update(id string, body map[string]any, now time.Time, actor account.User) (map[string]any, error) {
	// PATCH semantics: absent fields keep their stored values.
	existing, err := e.repository.GetType(id)
	if err != nil {
		return nil, mapDictStoreError(err)
	}
	name := stringField(body, "name")
	if name == "" {
		name = existing.Name
	}
	description := stringField(body, "description")
	if _, present := body["description"]; !present {
		description = existing.Description
	}
	enabled := boolField(body, "enabled", existing.Enabled)
	sort := intField(body, "sort")
	if _, present := body["sort"]; !present {
		sort = existing.Sort
	}
	if err := e.repository.UpdateType(id, name, enabled, description, sort, now); err != nil {
		return nil, mapDictStoreError(err)
	}
	recordDictionaryEvent(e.operations, operationlog.EventDictionaryUpdate, actor, id, now)
	row, err := e.repository.GetType(id)
	if err != nil {
		return nil, mapDictStoreError(err)
	}
	return dictTypeToMap(*row), nil
}

func (e *dictTypeEntity) Delete(id string, actor account.User) error {
	entryIDs, err := e.repository.DeleteType(id)
	if err != nil {
		return mapDictStoreError(err)
	}
	// A-003 F-003: cascade-deleted entries are recorded in the type event's
	// detail so forensics keep the entry ids.
	detail := ""
	if len(entryIDs) > 0 {
		raw, err := json.Marshal(map[string]any{"entries": entryIDs})
		if err == nil {
			detail = string(raw)
		}
	}
	recordDictionaryEventDetail(e.operations, operationlog.EventDictionaryDelete, actor, id, detail, time.Now().UTC())
	return nil
}

// dictEntryEntity adapts the entry resource to the generic factory.
type dictEntryEntity struct {
	repository DictionaryRepository
	operations operationlog.Recorder
}

func (e *dictEntryEntity) List(filter resourceFilter) ([]map[string]any, int, error) {
	rows, total, err := e.repository.ListEntries(datadictionarystore.ListFilter{Q: filter.Q, DictKey: filter.Extra["dictKey"], Sort: filter.Sort, Order: filter.Order, Page: filter.Page, PageSize: filter.PageSize})
	if err != nil {
		return nil, 0, err
	}
	items := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		items = append(items, dictEntryToMap(row))
	}
	return items, total, nil
}

func (e *dictEntryEntity) Get(id string) (map[string]any, error) {
	row, err := e.repository.GetEntry(id)
	if err != nil {
		return nil, mapDictStoreError(err)
	}
	return dictEntryToMap(*row), nil
}

func (e *dictEntryEntity) Create(body map[string]any, id string, now time.Time, actor account.User) (map[string]any, error) {
	dictKey := stringField(body, "dictKey")
	entryKey := stringField(body, "entryKey")
	label := stringField(body, "label")
	if dictKey == "" || entryKey == "" || label == "" {
		return nil, &DomainError{Status: http.StatusBadRequest, Code: "INVALID_CREATE_FIELD", Message: "dictKey, entryKey and label are required"}
	}
	err := e.repository.CreateEntry(datadictionarystore.DictEntry{
		ID: id, DictKey: dictKey, EntryKey: entryKey, Label: label,
		Enabled: boolField(body, "enabled", true), Sort: intField(body, "sort"),
		Remark: stringField(body, "remark"), CreatedAt: now, UpdatedAt: now,
	})
	if err != nil {
		return nil, mapDictStoreError(err)
	}
	recordDictionaryEvent(e.operations, operationlog.EventDictionaryCreate, actor, id, now)
	row, err := e.repository.GetEntry(id)
	if err != nil {
		return nil, err
	}
	return dictEntryToMap(*row), nil
}

func (e *dictEntryEntity) Update(id string, body map[string]any, now time.Time, actor account.User) (map[string]any, error) {
	// PATCH semantics: absent fields keep their stored values.
	existing, err := e.repository.GetEntry(id)
	if err != nil {
		return nil, mapDictStoreError(err)
	}
	dictKey := stringField(body, "dictKey")
	if dictKey == "" {
		dictKey = existing.DictKey
	}
	label := stringField(body, "label")
	if label == "" {
		label = existing.Label
	}
	enabled := boolField(body, "enabled", existing.Enabled)
	sort := intField(body, "sort")
	if _, present := body["sort"]; !present {
		sort = existing.Sort
	}
	remark := stringField(body, "remark")
	if _, present := body["remark"]; !present {
		remark = existing.Remark
	}
	if err := e.repository.UpdateEntry(id, dictKey, label, enabled, sort, remark, now); err != nil {
		return nil, mapDictStoreError(err)
	}
	recordDictionaryEvent(e.operations, operationlog.EventDictionaryUpdate, actor, id, now)
	row, err := e.repository.GetEntry(id)
	if err != nil {
		return nil, mapDictStoreError(err)
	}
	return dictEntryToMap(*row), nil
}

func (e *dictEntryEntity) Delete(id string, actor account.User) error {
	err := e.repository.DeleteEntry(id)
	if err != nil {
		return mapDictStoreError(err)
	}
	recordDictionaryEvent(e.operations, operationlog.EventDictionaryDelete, actor, id, time.Now().UTC())
	return nil
}

// DictionaryRoutes returns the admin.data-dictionary HTTP surface.
func DictionaryRoutes(a *auth.Authenticator, repository DictionaryRepository, operations operationlog.Recorder, moduleID string, trash ...TrashRecorder) []kernel.RouteContribution {
	// S-12 (GOAL-012 D-002 §2): optional recycle-bin snapshot hook. nil keeps
	// the delete semantics byte-identical.
	var recorder TrashRecorder
	if len(trash) > 0 {
		recorder = trash[0]
	}
	routes := ResourceRoutes(a, Resource{
		ID:              "dict-types",
		Path:            "/api/data-dictionary/types",
		Listable:        true,
		SortFields:      []string{"key", "name", "sort", "updatedAt"},
		QSearch:         true,
		Entity:          &dictTypeEntity{repository: repository, operations: operations},
		CreateFields:    []string{"key", "name"},
		PatchFields:     []string{"name"},
		// description is optional (create: entity default ""; patch: absent =
		// untouched, present = updated, including clearing to "") — the factory
		// only passes CreateFields/JSONFields through, so it must ride JSONFields
		// (fixed: was silently dropped on create, field always empty).
		JSONFields: []string{"enabled", "sort", "description"},
		PermissionRead:  "dictionary.read",
		PermissionWrite: "dictionary.write",
		NotFoundCode:    "DICT_TYPE_NOT_FOUND",
		NewID:           newDictionaryID,
		Trash:           recorder,
	}, moduleID)
	routes = append(routes, ResourceRoutes(a, Resource{
		ID:              "dict-entries",
		Path:            "/api/data-dictionary/entries",
		Listable:        true,
		SortFields:      []string{"dictKey", "entryKey", "label", "sort", "updatedAt", "dictTypeName"},
		QSearch:         true,
		Entity:          &dictEntryEntity{repository: repository, operations: operations},
		CreateFields:    []string{"dictKey", "entryKey", "label"},
		PatchFields:     []string{"dictKey", "label"},
		JSONFields:      []string{"enabled", "sort", "remark"},
		// GOAL-015: inner page narrows entries by exact dict key.
		ExtraQuery: []string{"dictKey"},
		PermissionRead:  "dictionary.read",
		PermissionWrite: "dictionary.write",
		NotFoundCode:    "DICT_ENTRY_NOT_FOUND",
		NewID:           newDictionaryID,
		Trash:           recorder,
	}, moduleID)...)
	return routes
}

func dictTypeToMap(t datadictionarystore.DictType) map[string]any {
	return map[string]any{
		"id": t.ID, "key": t.Key, "name": t.Name, "enabled": t.Enabled,
		"description": t.Description, "sort": t.Sort,
		"createdAt": t.CreatedAt.Unix(), "updatedAt": t.UpdatedAt.Unix(),
	}
}

func dictEntryToMap(e datadictionarystore.DictEntry) map[string]any {
	return map[string]any{
		"id": e.ID, "dictKey": e.DictKey, "dictTypeName": e.DictTypeName, "entryKey": e.EntryKey, "label": e.Label,
		"enabled": e.Enabled, "sort": e.Sort, "remark": e.Remark,
		"createdAt": e.CreatedAt.Unix(), "updatedAt": e.UpdatedAt.Unix(),
	}
}

// mapDictStoreError maps the repository sentinels to the frozen wire codes.
func mapDictStoreError(err error) error {
	switch {
	case errors.Is(err, datadictionarystore.ErrNotFound):
		return errResourceNotFound
	case errors.Is(err, datadictionarystore.ErrTypeKeyTaken):
		return &DomainError{Status: http.StatusConflict, Code: "DICT_TYPE_KEY_TAKEN", Message: "a dict type with that key already exists"}
	case errors.Is(err, datadictionarystore.ErrEntryKeyTaken):
		return &DomainError{Status: http.StatusConflict, Code: "DICT_ENTRY_KEY_TAKEN", Message: "an entry with that key already exists in the dict type"}
	case errors.Is(err, datadictionarystore.ErrDictKeyNotFound):
		return &DomainError{Status: http.StatusBadRequest, Code: "DICT_KEY_NOT_FOUND", Message: "no dict type with that key"}
	default:
		return err
	}
}

func boolField(body map[string]any, key string, fallback bool) bool {
	if v, ok := body[key].(bool); ok {
		return v
	}
	return fallback
}

func intField(body map[string]any, key string) int {
	if v, ok := body[key].(float64); ok {
		return int(v)
	}
	return 0
}

func recordDictionaryEvent(operations operationlog.Recorder, event string, user account.User, id string, now time.Time) {
	if operations == nil {
		return
	}
	recordID := id
	if err := operations.RecordOperation(operationlog.Operation{
		ID: newOperationID(), Event: event, ActorID: user.ID, ActorName: user.Name,
		RecordID: &recordID, CreatedAt: now.UTC(),
	}); err != nil {
		slog.Error("operation log write failed", "event", event, "err", err)
	}
}

// newDictionaryID returns "dict-" + 16 lowercase hex chars (8 bytes of
// crypto/rand), the dictionary create id format (GOAL-008 D-002 §2).
func newDictionaryID() (string, error) {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return "dict-" + hex.EncodeToString(b[:]), nil
}

// recordDictionaryEventDetail records an audit event with an optional detail
// payload (same best-effort contract as recordDictionaryEvent).
func recordDictionaryEventDetail(operations operationlog.Recorder, event string, user account.User, id string, detail string, now time.Time) {
	if operations == nil {
		return
	}
	recordID := id
	var detailPtr *string
	if detail != "" {
		detailPtr = &detail
	}
	if err := operations.RecordOperation(operationlog.Operation{
		ID: newOperationID(), Event: event, ActorID: user.ID, ActorName: user.Name,
		RecordID: &recordID, Detail: detailPtr, CreatedAt: now.UTC(),
	}); err != nil {
		slog.Error("operation log write failed", "event", event, "err", err)
	}
}
