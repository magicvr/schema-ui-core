package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
)

// Limits for the R5 D-DATA API (F-009-007): write bodies and page size are
// bounded so the demo API fails closed on oversized input.
const (
	maxRecordBodyBytes = 4 << 10
	maxPageSize        = 100
	// recordIDRetries bounds PK-collision retries before INTERNAL when the
	// crypto/rand id collides (I-007-001 §2); with 64 random bits this never
	// triggers in practice.
	recordIDRetries = 3
)

// record is the R4 records API entity. UpdatedAt serializes with fixed
// millisecond precision (GOAL-007 D-004).
type record struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	Owner     string    `json:"owner"`
	UpdatedAt updatedAt `json:"updatedAt"`
}

// updatedAt is a UTC timestamp that JSON-marshals as RFC3339 with fixed 3-digit
// milliseconds ("2006-01-02T15:04:05.000Z07:00", D-004 / I-007-001 v0.2.0). Go's
// default time.Time encoding trims trailing zeros, which would break the frozen
// "含毫秒" shape.
type updatedAt struct {
	time.Time
}

func (u updatedAt) MarshalJSON() ([]byte, error) {
	return []byte(`"` + u.Time.UTC().Format("2006-01-02T15:04:05.000Z07:00") + `"`), nil
}

// recordList is the list response envelope consumed by the Web data table.
type recordList struct {
	Items    []record `json:"items"`
	Total    int      `json:"total"`
	Page     int      `json:"page"`
	PageSize int      `json:"pageSize"`
}

// recordHandler serves the R4 records CRUD API backed by the SQLite repository
// (GOAL-007 S3). All routes go through the request-identity middleware
// (GOAL-006 S4) and a permission gate: reads require `records.read`, writes
// require `records.write`. Anonymous → 401, authenticated without the
// permission → 403. There is no in-process slice fallback (I-007-002 §5 /
// T-DB-09).
type recordHandler struct {
	st *store.Store
}

// recordPatch is the PATCH body: pointer fields so absent keys are untouched.
type recordPatch struct {
	Name   *string `json:"name"`
	Status *string `json:"status"`
	Owner  *string `json:"owner"`
}

var recordSortFields = []string{"name", "status", "owner", "updatedAt"}

func recordsHandler(mux *http.ServeMux, a *auth.Authenticator, st *store.Store) {
	h := &recordHandler{st: st}
	h.routes(mux, a)
}

func (h *recordHandler) routes(mux *http.ServeMux, a *auth.Authenticator) {
	mux.Handle("GET /api/records", a.Middleware(h.list()))
	mux.Handle("POST /api/records", a.Middleware(h.create()))
	mux.Handle("GET /api/records/{id}", a.Middleware(h.detail()))
	mux.Handle("PATCH /api/records/{id}", a.Middleware(h.update()))
	mux.Handle("DELETE /api/records/{id}", a.Middleware(h.delete()))
}

// requirePermission enforces fail-closed authorization (GOAL-006 S4): the
// request identity must be present and hold the given permission key, resolved
// from the persisted role-permission relations at identity load. Anonymous →
// 401; authenticated without the permission → 403.
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

// list serves GET /api/records with q / sort / order / page / pageSize params.
func (h *recordHandler) list() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePermission(w, r.Context(), "records.read"); !ok {
			return
		}
		query := r.URL.Query()

		sortField := query.Get("sort")
		if sortField == "" {
			sortField = "name"
		}
		if !slices.Contains(recordSortFields, sortField) {
			writeError(w, http.StatusBadRequest, "INVALID_SORT_FIELD", "unsupported sort field")
			return
		}
		order := query.Get("order")
		if order == "" {
			order = "asc"
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

		q := strings.ToLower(strings.TrimSpace(query.Get("q")))
		items, total, err := h.st.ListRecords(store.RecordFilter{
			Q: q, Sort: sortField, Order: order, Page: page, PageSize: pageSize,
		})
		if err != nil {
			writeError(w, http.StatusInternalServerError, "INTERNAL", "could not list records")
			return
		}
		out := make([]record, 0, len(items))
		for _, it := range items {
			out = append(out, toRecord(it))
		}
		writeJSON(w, http.StatusOK, recordList{
			Items:    out,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		})
	})
}

// create serves POST /api/records (records.write): 201 + the full record with a
// server-generated id and updatedAt (I-007-001 §3.2).
func (h *recordHandler) create() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePermission(w, r.Context(), "records.write"); !ok {
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxRecordBodyBytes)
		name, status, owner, err := decodeCreateBody(r)
		if err != nil {
			var fe createFieldError
			if errors.As(err, &fe) {
				writeError(w, http.StatusBadRequest, "INVALID_CREATE_FIELD", fe.Error())
				return
			}
			writeError(w, http.StatusBadRequest, "INVALID_CREATE_BODY", "body must be JSON")
			return
		}

		var rec *store.Record
		for attempt := 0; attempt < recordIDRetries; attempt++ {
			id, err := newRecordID()
			if err != nil {
				writeError(w, http.StatusInternalServerError, "INTERNAL", "could not generate record id")
				return
			}
			rec, err = h.st.CreateRecord(store.Record{
				ID: id, Name: name, Status: status, Owner: owner, UpdatedAt: time.Now().UTC(),
			})
			if err == nil {
				break
			}
			if !errors.Is(err, store.ErrRecordExists) {
				writeError(w, http.StatusInternalServerError, "INTERNAL", "could not create record")
				return
			}
		}
		if rec == nil {
			writeError(w, http.StatusInternalServerError, "INTERNAL", "could not create record")
			return
		}
		writeJSON(w, http.StatusCreated, toRecord(*rec))
	})
}

// createFieldError reports which create field failed and why; the handler maps
// it to INVALID_CREATE_FIELD (I-007-001 §3.2).
type createFieldError struct {
	field  string
	reason string
}

func (e createFieldError) Error() string { return e.field + " must " + e.reason }

// decodeCreateBody parses the POST /api/records body. A body that is not a JSON
// object (or is truncated/oversized) returns a plain error → INVALID_CREATE_BODY;
// a missing, non-string or blank-trimmed field returns createFieldError →
// INVALID_CREATE_FIELD. Unknown keys are ignored (same posture as PATCH).
func decodeCreateBody(r *http.Request) (name, status, owner string, err error) {
	var raw map[string]json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
		return "", "", "", err
	}
	if name, err = createStringField(raw, "name"); err != nil {
		return "", "", "", err
	}
	if status, err = createStringField(raw, "status"); err != nil {
		return "", "", "", err
	}
	if owner, err = createStringField(raw, "owner"); err != nil {
		return "", "", "", err
	}
	return name, status, owner, nil
}

func createStringField(raw map[string]json.RawMessage, key string) (string, error) {
	val, ok := raw[key]
	if !ok {
		return "", createFieldError{key, "not be empty"}
	}
	var s string
	if err := json.Unmarshal(val, &s); err != nil {
		return "", createFieldError{key, "be a string"}
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return "", createFieldError{key, "not be empty"}
	}
	return s, nil
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

// detail serves GET /api/records/{id}.
func (h *recordHandler) detail() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePermission(w, r.Context(), "records.read"); !ok {
			return
		}
		rec, err := h.st.GetRecord(r.PathValue("id"))
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "RECORD_NOT_FOUND", "no record with that id")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "INTERNAL", "could not load record")
			return
		}
		writeJSON(w, http.StatusOK, toRecord(*rec))
	})
}

// update serves PATCH /api/records/{id} for the D-ACT edit lifecycle.
func (h *recordHandler) update() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePermission(w, r.Context(), "records.write"); !ok {
			return
		}
		id := r.PathValue("id")
		r.Body = http.MaxBytesReader(w, r.Body, maxRecordBodyBytes)
		var patch recordPatch
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&patch); err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_PATCH_BODY", "body must be JSON")
			return
		}
		if err := validatePatch(patch); err != "" {
			writeError(w, http.StatusBadRequest, "INVALID_PATCH_FIELD", err)
			return
		}

		rec, err := h.st.UpdateRecord(id, store.RecordPatch{
			Name: patch.Name, Status: patch.Status, Owner: patch.Owner,
		}, time.Now().UTC())
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "RECORD_NOT_FOUND", "no record with that id")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "INTERNAL", "could not update record")
			return
		}
		writeJSON(w, http.StatusOK, toRecord(*rec))
	})
}

// delete serves DELETE /api/records/{id} for the D-ACT edit lifecycle.
func (h *recordHandler) delete() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePermission(w, r.Context(), "records.write"); !ok {
			return
		}
		if err := h.st.DeleteRecord(r.PathValue("id")); err != nil {
			if errors.Is(err, store.ErrNotFound) {
				writeError(w, http.StatusNotFound, "RECORD_NOT_FOUND", "no record with that id")
				return
			}
			writeError(w, http.StatusInternalServerError, "INTERNAL", "could not delete record")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
}

// toRecord maps a persisted record to the API entity (UpdatedAt carries the
// fixed-millisecond JSON marshaller).
func toRecord(r store.Record) record {
	return record{
		ID:        r.ID,
		Name:      r.Name,
		Status:    r.Status,
		Owner:     r.Owner,
		UpdatedAt: updatedAt{r.UpdatedAt},
	}
}

// validatePatch rejects empty trimmed strings for the mutable fields.
func validatePatch(patch recordPatch) string {
	if patch.Name != nil && strings.TrimSpace(*patch.Name) == "" {
		return "name must not be empty"
	}
	if patch.Status != nil && strings.TrimSpace(*patch.Status) == "" {
		return "status must not be empty"
	}
	if patch.Owner != nil && strings.TrimSpace(*patch.Owner) == "" {
		return "owner must not be empty"
	}
	return ""
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
