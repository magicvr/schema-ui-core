package handler

import (
	"cmp"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
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

// recordHandler serves the R5 D-DATA list/detail example API (static data).
type recordHandler struct {
	records []record
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

func recordsHandler(mux *http.ServeMux) {
	h := &recordHandler{records: staticRecords()}
	mux.Handle("GET /api/records", h.list())
	mux.Handle("GET /api/records/{id}", h.detail())
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

		q := strings.ToLower(strings.TrimSpace(query.Get("q")))
		items := make([]record, 0, len(h.records))
		for _, rec := range h.records {
			if q == "" || matches(rec, q) {
				items = append(items, rec)
			}
		}
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
		idx := slices.IndexFunc(h.records, func(rec record) bool { return rec.ID == id })
		if idx < 0 {
			writeError(w, http.StatusNotFound, "RECORD_NOT_FOUND", "no record with that id")
			return
		}
		writeJSON(w, http.StatusOK, h.records[idx])
	})
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
