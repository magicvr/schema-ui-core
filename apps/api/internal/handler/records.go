package handler

import (
	"cmp"
	"context"
	"encoding/json"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
)

// Limits for the R5 D-DATA API (F-009-007): write bodies and page size are
// bounded so the demo API fails closed on oversized input.
const (
	maxRecordBodyBytes = 4 << 10
	maxPageSize        = 100
)

// record is the R5 D-DATA example domain entity (de-business-ified static dev
// data). The MVP freeze does not vendor an upstream dataset; this is the local
// static-data support path for the search-form-table / data-table examples.
type record struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	Owner     string    `json:"owner"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// recordList is the list response envelope consumed by the Web data table.
type recordList struct {
	Items    []record `json:"items"`
	Total    int      `json:"total"`
	Page     int      `json:"page"`
	PageSize int      `json:"pageSize"`
}

// recordHandler serves the R5 D-DATA list/detail example API (static data) and
// the R5 D-ACT edit lifecycle (PATCH/DELETE). State is guarded by a mutex; the
// MVP has no DB for records, so mutations live in the process. Write routes
// gate through the R2 request-identity middleware (F-009-006): 401 unauthenticated,
// 403 non-admin.
type recordHandler struct {
	mu      sync.RWMutex
	records []record
}

// recordPatch is the PATCH body: pointer fields so absent keys are untouched.
type recordPatch struct {
	Name   *string `json:"name"`
	Status *string `json:"status"`
	Owner  *string `json:"owner"`
}

var recordSortFields = []string{"name", "status", "owner", "updatedAt"}

// staticRecords is a stable, ordered dev dataset (no DB in the MVP).
func staticRecords() []record {
	base := time.Date(2026, 7, 31, 0, 0, 0, 0, time.UTC)
	names := []string{
		"Acme Console", "Northwind Sales", "Hooli Connect", "Umbrella Ops",
		"Initech Reports", "Stark Access", "Wayne Fleet", "Globex Admin",
	}
	statuses := []string{"active", "archived", "pending"}
	owners := []string{"alice", "bob", "carol"}
	records := make([]record, 0, len(names))
	for i, name := range names {
		records = append(records, record{
			ID:        "rec-" + strconv.Itoa(i+1),
			Name:      name,
			Status:    statuses[i%len(statuses)],
			Owner:     owners[i%len(owners)],
			UpdatedAt: base.Add(time.Duration(i*11) * time.Hour),
		})
	}
	return records
}

func recordsHandler(mux *http.ServeMux, a *auth.Authenticator) {
	h := &recordHandler{records: staticRecords()}
	h.routes(mux, a)
}

func (h *recordHandler) routes(mux *http.ServeMux, a *auth.Authenticator) {
	mux.Handle("GET /api/records", h.list())
	mux.Handle("GET /api/records/{id}", h.detail())
	mux.Handle("PATCH /api/records/{id}", a.Middleware(h.update()))
	mux.Handle("DELETE /api/records/{id}", a.Middleware(h.delete()))
}

// writeGate enforces fail-closed authorization on record write routes
// (F-009-006): the request identity is resolved from the Bearer access token
// and the caller must hold the admin role. On denial it writes the error
// response and returns ok=false.
func writeGate(w http.ResponseWriter, ctx context.Context) (account.User, bool) {
	user, ok := auth.IdentityFrom(ctx)
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHENTICATED", "no active session")
		return account.User{}, false
	}
	if !account.Allow(`$context.user.roles contains "admin"`, user, nil) {
		writeError(w, http.StatusForbidden, "FORBIDDEN", "admin role required")
		return account.User{}, false
	}
	return user, true
}

// list serves GET /api/records with q / sort / order / page / pageSize params.
func (h *recordHandler) list() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		h.mu.RLock()
		all := h.records
		items := make([]record, 0, len(all))
		for _, rec := range all {
			if q == "" || matches(rec, q) {
				items = append(items, rec)
			}
		}
		h.mu.RUnlock()
		slices.SortFunc(items, func(a, b record) int {
			var byField int
			switch sortField {
			case "status":
				byField = cmp.Compare(a.Status, b.Status)
			case "owner":
				byField = cmp.Compare(a.Owner, b.Owner)
			case "updatedAt":
				byField = a.UpdatedAt.Compare(b.UpdatedAt)
			default:
				byField = cmp.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
			}
			if order == "desc" {
				return -byField
			}
			return byField
		})

		total := len(items)
		start := min((page-1)*pageSize, total)
		end := min(start+pageSize, total)
		writeJSON(w, http.StatusOK, recordList{
			Items:    items[start:end],
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		})
	})
}

// detail serves GET /api/records/{id}.
func (h *recordHandler) detail() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		h.mu.RLock()
		rec, ok := findRecord(h.records, id)
		h.mu.RUnlock()
		if !ok {
			writeError(w, http.StatusNotFound, "RECORD_NOT_FOUND", "no record with that id")
			return
		}
		writeJSON(w, http.StatusOK, rec)
	})
}

// update serves PATCH /api/records/{id} for the D-ACT edit lifecycle.
func (h *recordHandler) update() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := writeGate(w, r.Context()); !ok {
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

		h.mu.Lock()
		idx := slices.IndexFunc(h.records, func(rec record) bool { return rec.ID == id })
		if idx < 0 {
			h.mu.Unlock()
			writeError(w, http.StatusNotFound, "RECORD_NOT_FOUND", "no record with that id")
			return
		}
		rec := h.records[idx]
		if patch.Name != nil {
			rec.Name = *patch.Name
		}
		if patch.Status != nil {
			rec.Status = *patch.Status
		}
		if patch.Owner != nil {
			rec.Owner = *patch.Owner
		}
		rec.UpdatedAt = time.Now().UTC()
		h.records[idx] = rec
		h.mu.Unlock()
		writeJSON(w, http.StatusOK, rec)
	})
}

// delete serves DELETE /api/records/{id} for the D-ACT edit lifecycle.
func (h *recordHandler) delete() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := writeGate(w, r.Context()); !ok {
			return
		}
		id := r.PathValue("id")
		h.mu.Lock()
		idx := slices.IndexFunc(h.records, func(rec record) bool { return rec.ID == id })
		if idx < 0 {
			h.mu.Unlock()
			writeError(w, http.StatusNotFound, "RECORD_NOT_FOUND", "no record with that id")
			return
		}
		h.records = slices.Delete(h.records, idx, idx+1)
		h.mu.Unlock()
		w.WriteHeader(http.StatusNoContent)
	})
}

func findRecord(records []record, id string) (record, bool) {
	idx := slices.IndexFunc(records, func(rec record) bool { return rec.ID == id })
	if idx < 0 {
		return record{}, false
	}
	return records[idx], true
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

func matches(rec record, q string) bool {
	return strings.Contains(strings.ToLower(rec.Name), q) ||
		strings.Contains(strings.ToLower(rec.Status), q) ||
		strings.Contains(strings.ToLower(rec.Owner), q)
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
