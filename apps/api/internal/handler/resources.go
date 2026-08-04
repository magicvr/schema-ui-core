// Generic schema-driven resource CRUD (GOAL-010 S2 · I-010-001 §4).
//
// A Resource descriptor + a ResourceEntity adapter let a new business resource
// mount the same five routes (list / create / detail / update / delete) without
// a hand-written handler. The factory owns the shared concerns — permission
// gate, body-size bound, page/sort/q validation, the frozen {error,message}
// write envelope and INTERNAL fallback. users/roles are the GOAL-011 semantic
// resource instances.
package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
)

// Body and page-size bounds shared by every registered resource (I-010-001 §4);
// they preserve the limits frozen in I-007-001 §2.
const (
	maxResourceBodyBytes = 4 << 10
	maxPageSize          = 100
	// resourceIDRetries bounds PK-collision retries before INTERNAL.
	resourceIDRetries = 3
)

// resourceFilter carries the list query parameters already validated by the
// generic factory (sort field, order, page, pageSize, optional search text).
type resourceFilter struct {
	Q        string
	Sort     string
	Order    string
	Page     int
	PageSize int
}

// ResourceEntity is the store boundary the generic factory drives. Rows are
// plain JSON maps; the concrete entity adapter owns store-specific shaping
// (id format, timestamp serialization, column mapping). Create/Update/Delete
// receive the request actor so a resource can enforce actor-dependent domain
// invariants (e.g. users self/last-admin protection, I-011-001 §7.2).
type ResourceEntity interface {
	List(filter resourceFilter) ([]map[string]any, int, error)
	Get(id string) (map[string]any, error)
	Create(body map[string]any, id string, now time.Time, user account.User) (map[string]any, error)
	Update(id string, body map[string]any, now time.Time, user account.User) (map[string]any, error)
	Delete(id string, user account.User) error
}

// DomainError is a typed domain error an entity may return from List/Get/
// Create/Update/Delete. The factory maps it verbatim (status + {error,message})
// BEFORE the generic ErrNotFound/INTERNAL fallbacks (I-011-001 §7.3). It is how
// a resource surfaces domain rejections (USERNAME_TAKEN, LAST_ADMIN, ROLE_IN_USE,
// …) with a precise status/code instead of degrading to INTERNAL.
type DomainError struct {
	Status  int
	Code    string
	Message string
}

func (e *DomainError) Error() string { return e.Code + ": " + e.Message }

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

// writeKind identifies a successful write so a resource can attach side effects
// (e.g. operation-log rows) uniformly.
type writeKind int

const (
	writeCreate writeKind = iota
	writeUpdate
	writeDelete
)

// errReadOnlyResource is returned by entity write methods on read-only resources
// (e.g. operations activity log). The factory never mounts write routes when
// Resource.ReadOnly is true, so this is a defensive fallback.
var errReadOnlyResource = errors.New("resource is read-only")

// Resource describes a schema-driven CRUD resource registered with the generic
// handler factory (I-010-001 §4). users and roles are registered instances.
type Resource struct {
	ID           string         // resource id, e.g. "users"
	Path         string         // mount path, e.g. "/api/users"
	Listable     bool           // expose GET {path} (list)
	// ReadOnly skips POST/PATCH/DELETE mounts (activity/operations).
	ReadOnly     bool
	SortFields   []string       // whitelist; empty = not sortable
	QSearch      bool           // list supports the q search param
	Entity       ResourceEntity // store adapter
	CreateFields []string       // required non-empty create fields (trimmed strings)
	PatchFields  []string       // editable patch fields (trimmed strings)
	// RawStringFields are required on create and optional on patch. String values
	// are passed through byte-for-byte after JSON decoding so secret material is
	// never silently normalized. The entity owns type and policy validation.
	RawStringFields []string
	// JSONFields are additional fields accepted as arbitrary JSON values (not
	// forced to strings). When present in a create/patch body the raw JSON value
	// is decoded and passed through to the entity; absent means untouched (patch)
	// or entity default (create). The entity owns shape validation (I-011-001
	// §7.1). users uses ["roles"].
	JSONFields      []string
	PermissionRead  string // defaults to "{id}.read"
	PermissionWrite string // defaults to "{id}.write"
	// NotFoundCode is the 404 error code, defaulting to "{ID}_NOT_FOUND".
	NotFoundCode string
	NewID        func() (string, error)
	// OnWrite, when set, runs after a successful create/update/delete.
	OnWrite func(ctx context.Context, user account.User, kind writeKind, id string, row map[string]any, now time.Time)
}

// resourceHandler serves the five generic CRUD routes for one Resource.
type resourceHandler struct {
	res       Resource
	readPerm  string
	writePerm string
	notFound  string
}

// resourceList is the unified list response envelope frozen across resources
// (I-010-001 §3): {items,total,page,pageSize}.
type resourceList struct {
	Items    []map[string]any `json:"items"`
	Total    int              `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"pageSize"`
}

// registerResource mounts the generic list/create/detail/update/delete routes
// under a resource's path, all wrapped in the request-identity middleware and a
// permission gate. Permission keys default to "{id}.read" / "{id}.write".
func registerResource(mux *http.ServeMux, a *auth.Authenticator, res Resource) {
	readPerm := res.PermissionRead
	if readPerm == "" {
		readPerm = res.ID + ".read"
	}
	writePerm := res.PermissionWrite
	if writePerm == "" {
		writePerm = res.ID + ".write"
	}
	h := &resourceHandler{
		res:       res,
		readPerm:  readPerm,
		writePerm: writePerm,
		notFound:  notFoundCode(res),
	}
	if res.Listable {
		mux.Handle("GET "+res.Path, a.Middleware(h.list()))
	}
	mux.Handle("GET "+res.Path+"/{id}", a.Middleware(h.detail()))
	if !res.ReadOnly {
		mux.Handle("POST "+res.Path, a.Middleware(h.create()))
		mux.Handle("PATCH "+res.Path+"/{id}", a.Middleware(h.update()))
		mux.Handle("DELETE "+res.Path+"/{id}", a.Middleware(h.delete()))
	}
}

// requirePermission enforces fail-closed authorization (GOAL-006 S4): the
// request identity must be present and hold the given permission key, resolved
// from the persisted role-permission relations at identity load. Anonymous → 401;
// authenticated without the permission → 403.
func requirePermission(w http.ResponseWriter, ctx context.Context, permission string) (account.User, bool) {
	user, ok := auth.IdentityFrom(ctx)
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHENTICATED", "no active session")
		return account.User{}, false
	}
	if !slices.Contains(user.Permissions, permission) {
		writeError(w, http.StatusForbidden, "FORBIDDEN", "permission required: "+permission)
		return account.User{}, false
	}
	return user, true
}

func notFoundCode(res Resource) string {
	if res.NotFoundCode != "" {
		return res.NotFoundCode
	}
	return strings.ToUpper(res.ID) + "_NOT_FOUND"
}

// writeEntityError maps an entity error to the wire in the frozen priority order
// (I-011-001 §7.3): DomainError verbatim, store.ErrNotFound → 404 with the
// resource's NOT_FOUND code, else INTERNAL. Returns true when a response was
// written (err != nil).
func writeEntityError(w http.ResponseWriter, res Resource, err error, action string) bool {
	if err == nil {
		return false
	}
	var de *DomainError
	if errors.As(err, &de) {
		writeError(w, de.Status, de.Code, de.Message)
		return true
	}
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, notFoundCode(res), "no "+res.ID+" with that id")
		return true
	}
	writeError(w, http.StatusInternalServerError, "INTERNAL", "could not "+action+" "+res.ID)
	return true
}

func stringField(m map[string]any, key string) string {
	if m == nil {
		return ""
	}
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

// list serves GET {path} with q / sort / order / page / pageSize params
// (q only when the resource declares QSearch).
func (h *resourceHandler) list() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePermission(w, r.Context(), h.readPerm); !ok {
			return
		}
		query := r.URL.Query()

		sortField := query.Get("sort")
		if sortField == "" {
			if len(h.res.SortFields) > 0 {
				sortField = h.res.SortFields[0]
			}
		}
		if !slices.Contains(h.res.SortFields, sortField) {
			writeError(w, http.StatusBadRequest, "INVALID_SORT_FIELD", "unsupported sort field")
			return
		}
		order := query.Get("order")
		if order == "" {
			// Activity/operation logs (ReadOnly) default to newest-first.
			if h.res.ReadOnly {
				order = "desc"
			} else {
				order = "asc"
			}
		}
		if order != "asc" && order != "desc" {
			writeError(w, http.StatusBadRequest, "INVALID_SORT_ORDER", "order must be asc or desc")
			return
		}
		page, ok := intParam(query.Get("page"), 1)
		if !ok {
			writeError(w, http.StatusBadRequest, "INVALID_PAGE", "page must be a positive integer")
			return
		}
		pageSize, ok := intParam(query.Get("pageSize"), 10)
		if !ok {
			writeError(w, http.StatusBadRequest, "INVALID_PAGE_SIZE", "pageSize must be a positive integer")
			return
		}
		if pageSize > maxPageSize {
			writeError(w, http.StatusBadRequest, "INVALID_PAGE_SIZE", "pageSize must not exceed 100")
			return
		}

		filter := resourceFilter{Q: "", Sort: sortField, Order: order, Page: page, PageSize: pageSize}
		if h.res.QSearch {
			filter.Q = strings.ToLower(strings.TrimSpace(query.Get("q")))
		}
		items, total, err := h.res.Entity.List(filter)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "INTERNAL", "could not list "+h.res.ID)
			return
		}
		if items == nil {
			items = []map[string]any{}
		}
		writeJSON(w, http.StatusOK, resourceList{Items: items, Total: total, Page: page, PageSize: pageSize})
	})
}

// createFieldError reports which create field failed and why; the factory maps
// it to INVALID_CREATE_FIELD (I-007-001 §3.2).
type createFieldError struct {
	field  string
	reason string
}

func (e createFieldError) Error() string { return e.field + " must " + e.reason }

// decodeResourceCreate parses a POST body. A body that is not a JSON object (or
// is truncated/oversized) returns a plain error → INVALID_CREATE_BODY; a missing,
// non-string or blank-trimmed required field returns createFieldError →
// INVALID_CREATE_FIELD. jsonFields (I-011-001 §7.1) are additionally decoded as
// arbitrary JSON values when present and passed through for the entity to
// validate. Unknown keys are ignored (same posture as PATCH).
func decodeResourceCreate(r *http.Request, fields, rawStringFields, jsonFields []string) (map[string]any, error) {
	var raw map[string]json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
		return nil, err
	}
	body := make(map[string]any, len(fields)+len(rawStringFields)+len(jsonFields))
	for _, key := range fields {
		val, ok := raw[key]
		if !ok {
			return nil, createFieldError{key, "not be empty"}
		}
		var s string
		if err := json.Unmarshal(val, &s); err != nil {
			return nil, createFieldError{key, "be a string"}
		}
		s = strings.TrimSpace(s)
		if s == "" {
			return nil, createFieldError{key, "not be empty"}
		}
		body[key] = s
	}
	for _, key := range rawStringFields {
		val, ok := raw[key]
		if !ok {
			body[key] = nil
			continue
		}
		var value any
		if err := json.Unmarshal(val, &value); err != nil {
			return nil, err
		}
		body[key] = value
	}
	for _, key := range jsonFields {
		val, ok := raw[key]
		if !ok {
			continue // optional JSON field: entity default
		}
		var v any
		if err := json.Unmarshal(val, &v); err != nil {
			return nil, err
		}
		body[key] = v
	}
	return body, nil
}

// patchFieldError reports a provided patch field that is blank after trim; the
// factory maps it to INVALID_PATCH_FIELD.
type patchFieldError struct {
	field string
}

func (e patchFieldError) Error() string { return e.field + " must not be empty" }

// decodeResourcePatch parses a PATCH body. Absent keys are untouched; a provided
// field must be a non-empty string (a non-string or blank value fails closed —
// INVALID_PATCH_BODY / INVALID_PATCH_FIELD respectively). jsonFields are decoded
// as arbitrary JSON values when present (absent = untouched). Unknown keys ignored.
func decodeResourcePatch(r *http.Request, fields, rawStringFields, jsonFields []string) (map[string]any, error) {
	var raw map[string]json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
		return nil, err
	}
	body := make(map[string]any, len(fields)+len(rawStringFields)+len(jsonFields))
	for _, key := range fields {
		val, ok := raw[key]
		if !ok {
			continue
		}
		var s string
		if err := json.Unmarshal(val, &s); err != nil {
			return nil, err
		}
		s = strings.TrimSpace(s)
		if s == "" {
			return nil, patchFieldError{key}
		}
		body[key] = s
	}
	for _, key := range rawStringFields {
		val, ok := raw[key]
		if !ok {
			continue
		}
		var value any
		if err := json.Unmarshal(val, &value); err != nil {
			return nil, err
		}
		body[key] = value
	}
	for _, key := range jsonFields {
		val, ok := raw[key]
		if !ok {
			continue // optional JSON field: untouched
		}
		var v any
		if err := json.Unmarshal(val, &v); err != nil {
			return nil, err
		}
		body[key] = v
	}
	return body, nil
}

// create serves POST {path} ({id}.write): 201 + the full row with a
// server-generated id (I-007-001 §3.2).
func (h *resourceHandler) create() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := requirePermission(w, r.Context(), h.writePerm)
		if !ok {
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxResourceBodyBytes)
		body, err := decodeResourceCreate(r, h.res.CreateFields, h.res.RawStringFields, h.res.JSONFields)
		if err != nil {
			var fe createFieldError
			if errors.As(err, &fe) {
				writeError(w, http.StatusBadRequest, "INVALID_CREATE_FIELD", fe.Error())
				return
			}
			writeError(w, http.StatusBadRequest, "INVALID_CREATE_BODY", "body must be JSON")
			return
		}

		now := time.Now().UTC()
		var row map[string]any
		for attempt := 0; attempt < resourceIDRetries; attempt++ {
			id, err := h.res.NewID()
			if err != nil {
				writeError(w, http.StatusInternalServerError, "INTERNAL", "could not generate "+h.res.ID+" id")
				return
			}
			row, err = h.res.Entity.Create(body, id, now, user)
			if err == nil {
				break
			}
			if errors.Is(err, store.ErrRecordExists) {
				continue // PK collision: retry with a fresh id
			}
			writeEntityError(w, h.res, err, "create")
			return
		}
		if row == nil {
			writeError(w, http.StatusInternalServerError, "INTERNAL", "could not create "+h.res.ID)
			return
		}
		if h.res.OnWrite != nil {
			h.res.OnWrite(r.Context(), user, writeCreate, stringField(row, "id"), row, now)
		}
		writeJSON(w, http.StatusCreated, row)
	})
}

// detail serves GET {path}/{id}.
func (h *resourceHandler) detail() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePermission(w, r.Context(), h.readPerm); !ok {
			return
		}
		row, err := h.res.Entity.Get(r.PathValue("id"))
		if writeEntityError(w, h.res, err, "load") {
			return
		}
		writeJSON(w, http.StatusOK, row)
	})
}

// update serves PATCH {path}/{id} ({id}.write).
func (h *resourceHandler) update() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := requirePermission(w, r.Context(), h.writePerm)
		if !ok {
			return
		}
		id := r.PathValue("id")
		r.Body = http.MaxBytesReader(w, r.Body, maxResourceBodyBytes)
		body, err := decodeResourcePatch(r, h.res.PatchFields, h.res.RawStringFields, h.res.JSONFields)
		if err != nil {
			var pe patchFieldError
			if errors.As(err, &pe) {
				writeError(w, http.StatusBadRequest, "INVALID_PATCH_FIELD", pe.Error())
				return
			}
			writeError(w, http.StatusBadRequest, "INVALID_PATCH_BODY", "body must be JSON")
			return
		}
		now := time.Now().UTC()
		row, err := h.res.Entity.Update(id, body, now, user)
		if writeEntityError(w, h.res, err, "update") {
			return
		}
		if h.res.OnWrite != nil {
			h.res.OnWrite(r.Context(), user, writeUpdate, id, row, now)
		}
		writeJSON(w, http.StatusOK, row)
	})
}

// delete serves DELETE {path}/{id} ({id}.write), returning 204.
func (h *resourceHandler) delete() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := requirePermission(w, r.Context(), h.writePerm)
		if !ok {
			return
		}
		id := r.PathValue("id")
		if err := h.res.Entity.Delete(id, user); err != nil {
			writeEntityError(w, h.res, err, "delete")
			return
		}
		if h.res.OnWrite != nil {
			h.res.OnWrite(r.Context(), user, writeDelete, id, nil, time.Now().UTC())
		}
		w.WriteHeader(http.StatusNoContent)
	})
}

func intParam(raw string, fallback int) (int, bool) {
	if raw == "" {
		return fallback, true
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value < 1 {
		return 0, false
	}
	return value, true
}
