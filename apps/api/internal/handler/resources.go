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
	"log/slog"
	"math"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/errorcatalog"
	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
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
	// Extra carries resource-specific list query parameters declared by
	// Resource.ExtraQuery (GOAL-015: dictKey for the dictionary-entries page).
	// Values are trimmed; keys are validated against the resource whitelist so
	// unknown parameters never reach an entity.
	Extra map[string]string
	// Scope carries the effective row-level scope constraint for the request
	// actor (S-09 · GOAL-016 D-002 §2); nil means no constraint. Entities that
	// opt into scoping consume it in their where assembly.
	Scope *ScopeConstraint
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

// BatchDeleter is the optional atomic whole-batch delete capability (ADR-0022
// D5d · D-001 P0). When an entity implements it, POST {path}/batch-delete
// delegates to it instead of looping single deletes, so the selection commits
// or rolls back as one unit. The returned count is the number of rows deleted.
type BatchDeleter interface {
	DeleteBatch(ids []string, user account.User) (int, error)
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
var (
	errReadOnlyResource = errors.New("resource is read-only")
	errResourceNotFound = errors.New("resource not found")
	errResourceExists   = errors.New("resource id already exists")
)

// Resource describes a schema-driven CRUD resource registered with the generic
// handler factory (I-010-001 §4). users and roles are registered instances.
//
// TrashRecorder is the recycle-bin snapshot surface (S-12 · GOAL-012 D-002 §2):
// a resource opts in by setting Resource.Trash; the factory captures the row
// BEFORE the delete and records the snapshot only after the delete succeeds,
// so a failed delete never leaves an orphan snapshot. The admin.recycle-bin
// module service satisfies this interface structurally.
type TrashRecorder interface {
	Record(ctx context.Context, resource, id string, row map[string]any, actor account.User, now time.Time) error
}

// ScopeConstraint is the resolved row-level scope for one request actor on
// one resource (S-09 · GOAL-016 D-002 §1). ScopeType "self" narrows access to
// rows whose OwnerColumn equals ActorID; "all" carries no restriction.
type ScopeConstraint struct {
	Resource    string
	ScopeType   string // "all" | "self"
	OwnerColumn string
	ActorID     string
}

// RowScopeProvider resolves the row-level scope for a request actor on a
// resource (S-09 · GOAL-016 D-002 §2). A nil *ScopeConstraint (or a nil
// provider) means no constraint applies — the request path is byte-identical
// to the unscoped behavior. Implemented by the admin.data-permission module
// service; resources opt in by setting Resource.Scoper.
type RowScopeProvider interface {
	ScopeFor(userID, resource string) (*ScopeConstraint, error)
}

type Resource struct {
	ID       string // resource id, e.g. "users"
	Path     string // mount path, e.g. "/api/users"
	Listable bool   // expose GET {path} (list)
	// ReadOnly skips POST/PATCH/DELETE mounts (activity/operations).
	ReadOnly     bool
	SortFields   []string       // whitelist; empty = not sortable
	QSearch      bool           // list supports the q search param
	// ExtraQuery lists resource-specific list query parameter names (entity
	// filters beyond q/sort/order/page/pageSize, e.g. "dictKey"). Declared
	// params are passed through resourceFilter.Extra; undeclared query keys are
	// ignored (GOAL-015 D-002 §3.1).
	ExtraQuery   []string
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
	// Trash, when set, records a pre-delete snapshot for the recycle bin
	// (S-12 · GOAL-012 D-002 §2): the factory captures the row before the
	// delete and records it only after the delete succeeds.
	Trash TrashRecorder
	// Scoper (S-09 · GOAL-016 D-002 §2), when set, resolves a row-level scope
	// for the request actor on this resource. nil (the default) means no
	// scoping — the unscoped path stays byte-identical. Composition wires the
	// admin.data-permission service only for resources whose entity consumes
	// filter.Scope (ScopeAware); v1 registers no production resource.
	Scoper RowScopeProvider
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
	for _, route := range resourceRoutes(a, res, "core.server-registration") {
		mux.Handle(route.Method+" "+route.Pattern, route.Handler)
	}
}

// resourceRoutes returns the generic CRUD route contributions for a Resource,
// matching exactly what registerResource mounts. moduleID is the owning provider
// id attached to the contribution identity; providers reuse the same factory so
// the provider-generated surface is byte-compatible with the central output
// (freeze package §7 step 2 compat comparison).
func resourceRoutes(a *auth.Authenticator, res Resource, moduleID string) []kernel.RouteContribution {
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
	var routes []kernel.RouteContribution
	add := func(method, pattern string, handler http.Handler) {
		routes = append(routes, kernel.RouteContribution{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey(method, pattern)},
			Method:               method,
			Pattern:              pattern,
			Handler:              handler,
		})
	}
	if res.Listable {
		add("GET", res.Path, a.Middleware(h.list()))
	}
	add("GET", res.Path+"/{id}", a.Middleware(h.detail()))
	if !res.ReadOnly {
		add("POST", res.Path, a.Middleware(h.create()))
		add("PATCH", res.Path+"/{id}", a.Middleware(h.update()))
		add("DELETE", res.Path+"/{id}", a.Middleware(h.delete()))
		// ADR-0022 batch delete (I-PROTO-FULL-001 · D-ACT/D-TABLE include):
		// one logical HTTP call for a normalized $selection.keys payload.
		add("POST", res.Path+"/batch-delete", a.Middleware(h.batchDelete()))
	}
	return routes
}

// ResourceRoutes exposes resourceRoutes for module providers (R4 C3.2).
func ResourceRoutes(a *auth.Authenticator, res Resource, moduleID string) []kernel.RouteContribution {
	return resourceRoutes(a, res, moduleID)
}

// requirePermission enforces fail-closed authorization (GOAL-006 S4): the
// request identity must be present and hold the given permission key, resolved
// from the persisted role-permission relations at identity load. Anonymous → 401;
// authenticated without the permission → 403.
func requirePermission(w http.ResponseWriter, r *http.Request, permission string) (account.User, bool) {
	user, ok := auth.IdentityFrom(r.Context())
	if !ok {
		writeLocalizedError(w, r, http.StatusUnauthorized, "UNAUTHENTICATED", "no active session")
		return account.User{}, false
	}
	if !slices.Contains(user.Permissions, permission) {
		writeLocalizedError(w, r, http.StatusForbidden, "FORBIDDEN", "permission required: "+permission)
		return account.User{}, false
	}
	return user, true
}

// resolveScope returns the effective scope constraint for the request actor,
// if the resource opts into scoping (S-09 · GOAL-016 D-002 §2).
func (h *resourceHandler) resolveScope(user account.User) (*ScopeConstraint, error) {
	if h.res.Scoper == nil {
		return nil, nil
	}
	return h.res.Scoper.ScopeFor(user.ID, h.res.ID)
}

// scopeOwned reports whether a row passes the self-scope ownership check. A
// nil or "all" constraint always passes.
func scopeOwned(row map[string]any, constraint *ScopeConstraint) bool {
	if constraint == nil || constraint.ScopeType != "self" {
		return true
	}
	return stringField(row, constraint.OwnerColumn) == constraint.ActorID
}

func notFoundCode(res Resource) string {
	if res.NotFoundCode != "" {
		return res.NotFoundCode
	}
	return strings.ToUpper(res.ID) + "_NOT_FOUND"
}

// writeEntityError maps an entity error to the wire in the frozen priority order
// (I-011-001 §7.3): DomainError verbatim, errResourceNotFound → 404 with the
// resource's NOT_FOUND code, else INTERNAL. Returns true when a response was
// written (err != nil).
func writeEntityError(w http.ResponseWriter, r *http.Request, res Resource, err error, action string) bool {
	if err == nil {
		return false
	}
	var de *DomainError
	if errors.As(err, &de) {
		writeLocalizedError(w, r, de.Status, de.Code, de.Message)
		return true
	}
	if errors.Is(err, errResourceNotFound) {
		writeLocalizedError(w, r, http.StatusNotFound, notFoundCode(res), "no "+res.ID+" with that id")
		return true
	}
	writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not "+action+" "+res.ID)
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
		user, ok := requirePermission(w, r, h.readPerm)
		if !ok {
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
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_SORT_FIELD", "unsupported sort field")
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
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_SORT_ORDER", "order must be asc or desc")
			return
		}
		page, ok := intParam(query.Get("page"), 1)
		if !ok {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PAGE", "page must be a positive integer")
			return
		}
		pageSize, ok := intParam(query.Get("pageSize"), 10)
		if !ok {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PAGE_SIZE", "pageSize must be a positive integer")
			return
		}
		if pageSize > maxPageSize {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PAGE_SIZE", "pageSize must not exceed 100")
			return
		}

		filter := resourceFilter{Q: "", Sort: sortField, Order: order, Page: page, PageSize: pageSize}
		if h.res.QSearch {
			filter.Q = strings.ToLower(strings.TrimSpace(query.Get("q")))
		}
		for _, key := range h.res.ExtraQuery {
			if value := strings.TrimSpace(query.Get(key)); value != "" {
				if filter.Extra == nil {
					filter.Extra = map[string]string{}
				}
				filter.Extra[key] = value
			}
		}
		// S-09 (GOAL-016 D-002 §2): the effective row-level scope rides the
		// filter into the entity's where assembly (ScopeAware contract).
		constraint, err := h.resolveScope(user)
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not resolve data scope")
			return
		}
		filter.Scope = constraint
		items, total, err := h.res.Entity.List(filter)
		if err != nil {
			var domainErr *DomainError
			if errors.As(err, &domainErr) {
				writeLocalizedError(w, r, domainErr.Status, domainErr.Code, domainErr.Message)
				return
			}
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not list "+h.res.ID)
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
		user, ok := requirePermission(w, r, h.writePerm)
		if !ok {
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxResourceBodyBytes)
		body, err := decodeResourceCreate(r, h.res.CreateFields, h.res.RawStringFields, h.res.JSONFields)
		if err != nil {
			var fe createFieldError
			if errors.As(err, &fe) {
				writeLocalizedFieldError(w, r, http.StatusBadRequest, "INVALID_CREATE_FIELD", fe.Error(), []errorcatalog.FieldError{{Field: fe.field, Reason: fe.reason}})
				return
			}
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_CREATE_BODY", "body must be JSON")
			return
		}

		now := time.Now().UTC()
		// S-09 (GOAL-016 D-002 §2): self scope forces the owner column to the
		// creating actor (A-005 recommended: overwrite, never trust the body).
		constraint, err := h.resolveScope(user)
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not resolve data scope")
			return
		}
		if constraint != nil && constraint.ScopeType == "self" {
			body[constraint.OwnerColumn] = constraint.ActorID
		}
		var row map[string]any
		for attempt := 0; attempt < resourceIDRetries; attempt++ {
			id, err := h.res.NewID()
			if err != nil {
				writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not generate "+h.res.ID+" id")
				return
			}
			row, err = h.res.Entity.Create(body, id, now, user)
			if err == nil {
				break
			}
			if errors.Is(err, errResourceExists) {
				continue // PK collision: retry with a fresh id
			}
			writeEntityError(w, r, h.res, err, "create")
			return
		}
		if row == nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not create "+h.res.ID)
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
		user, ok := requirePermission(w, r, h.readPerm)
		if !ok {
			return
		}
		row, err := h.res.Entity.Get(r.PathValue("id"))
		if writeEntityError(w, r, h.res, err, "load") {
			return
		}
		// S-09 (GOAL-016 D-002 §2): self scope hides rows owned by others as
		// 404 (no existence oracle).
		constraint, err := h.resolveScope(user)
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not resolve data scope")
			return
		}
		if !scopeOwned(row, constraint) {
			writeLocalizedError(w, r, http.StatusNotFound, notFoundCode(h.res), "no "+h.res.ID+" with that id")
			return
		}
		writeJSON(w, http.StatusOK, row)
	})
}

// update serves PATCH {path}/{id} ({id}.write).
func (h *resourceHandler) update() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := requirePermission(w, r, h.writePerm)
		if !ok {
			return
		}
		id := r.PathValue("id")
		r.Body = http.MaxBytesReader(w, r.Body, maxResourceBodyBytes)
		body, err := decodeResourcePatch(r, h.res.PatchFields, h.res.RawStringFields, h.res.JSONFields)
		if err != nil {
			var pe patchFieldError
			if errors.As(err, &pe) {
				writeLocalizedFieldError(w, r, http.StatusBadRequest, "INVALID_PATCH_FIELD", pe.Error(), []errorcatalog.FieldError{{Field: pe.field, Reason: "must not be empty"}})
				return
			}
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PATCH_BODY", "body must be JSON")
			return
		}
		now := time.Now().UTC()
		// S-09 (GOAL-016 D-002 §2): self scope rejects updates of rows owned by
		// others as 404 before any mutation.
		constraint, err := h.resolveScope(user)
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not resolve data scope")
			return
		}
		if constraint != nil && constraint.ScopeType == "self" {
			existing, gerr := h.res.Entity.Get(id)
			if gerr == nil && !scopeOwned(existing, constraint) {
				writeLocalizedError(w, r, http.StatusNotFound, notFoundCode(h.res), "no "+h.res.ID+" with that id")
				return
			}
			if gerr != nil && !errors.Is(gerr, errResourceNotFound) {
				writeEntityError(w, r, h.res, gerr, "update")
				return
			}
		}
		row, err := h.res.Entity.Update(id, body, now, user)
		if writeEntityError(w, r, h.res, err, "update") {
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
		user, ok := requirePermission(w, r, h.writePerm)
		if !ok {
			return
		}
		id := r.PathValue("id")
		// S-09 (GOAL-016 D-002 §2): self scope rejects deletes of rows owned by
		// others as 404 before any mutation.
		constraint, err := h.resolveScope(user)
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not resolve data scope")
			return
		}
		// S-12 (GOAL-012 D-002 §2): capture the pre-delete row so a successful
		// delete can record a recycle snapshot. A failed delete records nothing
		// (no orphan snapshots); a missing row records nothing.
		var snapshot map[string]any
		if h.res.Trash != nil || (constraint != nil && constraint.ScopeType == "self") {
			if row, gerr := h.res.Entity.Get(id); gerr == nil {
				if !scopeOwned(row, constraint) {
					writeLocalizedError(w, r, http.StatusNotFound, notFoundCode(h.res), "no "+h.res.ID+" with that id")
					return
				}
				snapshot = row
			}
		}
		if err := h.res.Entity.Delete(id, user); err != nil {
			writeEntityError(w, r, h.res, err, "delete")
			return
		}
		if h.res.OnWrite != nil {
			h.res.OnWrite(r.Context(), user, writeDelete, id, nil, time.Now().UTC())
		}
		if h.res.Trash != nil && snapshot != nil {
			now := time.Now().UTC()
			if err := h.res.Trash.Record(r.Context(), h.res.ID, id, snapshot, user, now); err != nil {
				slog.Error("recycle snapshot failed", "resource", h.res.ID, "id", id, "err", err)
			}
		}
		w.WriteHeader(http.StatusNoContent)
	})
}

// batchDelete serves POST {path}/batch-delete (ADR-0022 D5d · I-PROTO-FULL-001):
// one logical HTTP call over a normalized selection. The body is
// `{"ids": [scalar keys...]}` (the $selection.keys array of the page's
// batchMapping); ids must be non-empty scalar keys, deduped. Whole-batch
// semantics (D-001 P0): entities implementing BatchDeleter commit the selection
// in a single transaction — any failure rolls the whole batch back and returns
// the entity error; other entities fall back to sequential deletes, stopping at
// the first failure. Success returns `{"deleted": n}` so the client can reload
// (which clears selection).
func (h *resourceHandler) batchDelete() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := requirePermission(w, r, h.writePerm)
		if !ok {
			return
		}
		var body struct {
			IDs []any `json:"ids"`
		}
		decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxResourceBodyBytes))
		if err := decoder.Decode(&body); err != nil {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_BODY", "expected a JSON object with an ids array")
			return
		}
		if len(body.IDs) == 0 {
			writeLocalizedError(w, r, http.StatusBadRequest, "EMPTY_SELECTION", "ids must contain at least one key")
			return
		}
		// D3 invariants: scalar keys only, dedupe preserving order.
		seen := make(map[string]bool, len(body.IDs))
		ids := make([]string, 0, len(body.IDs))
		for _, raw := range body.IDs {
			var key string
			switch value := raw.(type) {
			case string:
				if value == "" {
					writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_SELECTION_KEY", "ids entries must be non-empty scalars")
					return
				}
				key = value
			case float64:
				if !isFiniteNumber(value) {
					writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_SELECTION_KEY", "ids entries must be finite scalars")
					return
				}
				key = formatNumberKey(value)
			case bool:
				key = strconv.FormatBool(value)
			default:
				writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_SELECTION_KEY", "ids entries must be scalar keys")
				return
			}
			if seen[key] {
				continue
			}
			seen[key] = true
			ids = append(ids, key)
		}
		if len(ids) == 0 {
			writeLocalizedError(w, r, http.StatusBadRequest, "EMPTY_SELECTION", "ids must contain at least one scalar key")
			return
		}

		// S-09 (GOAL-016 D-002 §2): self scope deletes only rows owned by the
		// actor — non-owned ids are skipped (never a 404 for the batch).
		constraint, err := h.resolveScope(user)
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not resolve data scope")
			return
		}
		if constraint != nil && constraint.ScopeType == "self" {
			owned := ids[:0]
			for _, id := range ids {
				if row, gerr := h.res.Entity.Get(id); gerr == nil && scopeOwned(row, constraint) {
					owned = append(owned, id)
				}
			}
			ids = owned
			if len(ids) == 0 {
				writeJSON(w, http.StatusOK, map[string]any{"deleted": 0})
				return
			}
		}

		deleted := 0
		if batch, ok := h.res.Entity.(BatchDeleter); ok {
			// S-12 (GOAL-012 D-002 §2): capture all pre-delete rows, then
			// record the snapshots only after the whole batch committed.
			type trashSnapshot struct {
				id  string
				row map[string]any
			}
			var snapshots []trashSnapshot
			if h.res.Trash != nil {
				for _, id := range ids {
					if row, err := h.res.Entity.Get(id); err == nil {
						snapshots = append(snapshots, trashSnapshot{id: id, row: row})
					}
				}
			}
			var err error
			deleted, err = batch.DeleteBatch(ids, user)
			if writeEntityError(w, r, h.res, err, "batch delete") {
				return
			}
			if h.res.Trash != nil {
				now := time.Now().UTC()
				for _, snapshot := range snapshots {
					if err := h.res.Trash.Record(r.Context(), h.res.ID, snapshot.id, snapshot.row, user, now); err != nil {
						slog.Error("recycle snapshot failed", "resource", h.res.ID, "id", snapshot.id, "err", err)
					}
				}
			}
		} else {
			now := time.Now().UTC()
			for _, id := range ids {
				// S-12 (GOAL-012 D-002 §2): capture the pre-delete row; the
				// snapshot is recorded only after this id's delete succeeds.
				var snapshot map[string]any
				if h.res.Trash != nil {
					if row, err := h.res.Entity.Get(id); err == nil {
						snapshot = row
					}
				}
				if err := h.res.Entity.Delete(id, user); err != nil {
					writeEntityError(w, r, h.res, err, "batch delete")
					return
				}
				if h.res.OnWrite != nil {
					h.res.OnWrite(r.Context(), user, writeDelete, id, nil, now)
				}
				if h.res.Trash != nil && snapshot != nil {
					if err := h.res.Trash.Record(r.Context(), h.res.ID, id, snapshot, user, now); err != nil {
						slog.Error("recycle snapshot failed", "resource", h.res.ID, "id", id, "err", err)
					}
				}
				deleted++
			}
			writeJSON(w, http.StatusOK, map[string]any{"deleted": deleted})
			return
		}
		now := time.Now().UTC()
		for _, id := range ids {
			if h.res.OnWrite != nil {
				h.res.OnWrite(r.Context(), user, writeDelete, id, nil, now)
			}
		}
		writeJSON(w, http.StatusOK, map[string]any{"deleted": deleted})
	})
}

func isFiniteNumber(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0)
}

func formatNumberKey(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
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
