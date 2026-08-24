// Mock-channel outbound record read surface (VP-017 R6 / workspace-017
// GOAL-007; contract frozen by workspace-017 GOAL-006 D-002 §3): authenticated,
// settings.read-gated list + detail over the Store-backed OutboxSink. This is
// the independent admin API per Root I-012 — deliberately NOT folded into
// /api/settings/*, no PATCH surface. Provider types never appear here: the
// handlers consume the retrieval interface below, not a channel client.
package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/magicvr/schema-ui-core/apps/api/internal/mail"
)

// OutboxReader is the retrieval face of the mock publisher (satisfied by
// *mail.OutboxSink). Keeping it an interface mirrors the repository-consumer
// convention and keeps the handler testable without a store.
type OutboxReader interface {
	List(ctx context.Context, limit, offset int) ([]mail.OutboxRecord, error)
	Count(ctx context.Context) (int64, error)
	Get(ctx context.Context, id string) (mail.OutboxRecord, error)
}

// RegisterMailOutbox mounts GET /api/mail/outbox and GET /api/mail/outbox/{id}.
// Both are authenticated and gated on settings.read — the same admin audience
// that will operate the R7 settings「邮件」tab; records are operational config
// data, not user-visible notifications.
func RegisterMailOutbox(mux routeRegistrar, a authMiddleware, reader OutboxReader) {
	mux.Handle("GET /api/mail/outbox", a.Middleware(outboxPermissionGate(outboxList(reader))))
	mux.Handle("GET /api/mail/outbox/{id}", a.Middleware(outboxPermissionGate(outboxDetail(reader))))
}

// outboxPermissionGate wraps both endpoints with the settings.read permission
// check (same model as uploadPermissionGate: fail-closed 401/403).
func outboxPermissionGate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePermission(w, r, "settings.read"); !ok {
			return
		}
		next.ServeHTTP(w, r)
	})
}

const (
	outboxDefaultPageSize = 50
	outboxMaxPageSize     = 200
)

// outboxList serves GET /api/mail/outbox?limit=&offset= with the unified list
// envelope {items,total,page,pageSize} (I-010-001 §3), newest first.
func outboxList(reader OutboxReader) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limit := outboxDefaultPageSize
		if raw := r.URL.Query().Get("limit"); raw != "" {
			parsed, err := strconv.Atoi(raw)
			if err != nil || parsed < 1 {
				writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PAGE_SIZE", "limit must be a positive integer")
				return
			}
			limit = parsed
		}
		offset := 0
		if raw := r.URL.Query().Get("offset"); raw != "" {
			parsed, err := strconv.Atoi(raw)
			if err != nil || parsed < 0 {
				writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PAGE", "offset must be a non-negative integer")
				return
			}
			offset = parsed
		}
		total64, err := reader.Count(r.Context())
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not count outbound records")
			return
		}
		total := int(total64)
		if limit > outboxMaxPageSize {
			limit = outboxMaxPageSize
		}
		items, err := reader.List(r.Context(), limit, offset)
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not load outbound records")
			return
		}
		if items == nil {
			items = []mail.OutboxRecord{}
		}
		page := 1
		if limit > 0 {
			page = offset/limit + 1
		}
		writeJSON(w, http.StatusOK, resourceList{Items: toMapItems(items), Total: total, Page: page, PageSize: limit})
	})
}

// outboxDetail serves GET /api/mail/outbox/{id} including the full body.
func outboxDetail(reader OutboxReader) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		rec, err := reader.Get(r.Context(), id)
		if errors.Is(err, mail.ErrOutboxRecordNotFound) {
			writeLocalizedError(w, r, http.StatusNotFound, "NOT_FOUND", "outbound record not found")
			return
		}
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not load outbound record")
			return
		}
		writeJSON(w, http.StatusOK, rec)
	})
}

// toMapItems adapts typed records into the generic envelope items.
func toMapItems(items []mail.OutboxRecord) []map[string]any {
	out := make([]map[string]any, 0, len(items))
	for _, rec := range items {
		out = append(out, map[string]any{
			"id":         rec.ID,
			"to":         rec.To,
			"subject":    rec.Subject,
			"created_at": rec.CreatedAt,
		})
	}
	return out
}
