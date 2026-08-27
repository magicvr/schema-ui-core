// All-channel outbound record read surface (VP-017 R6 / workspace-017
// GOAL-007; contract frozen by workspace-017 GOAL-006 D-002 §3, revised by
// W26 GOAL-038 D-001 §2.1 to cover every channel): authenticated,
// settings.read-gated list + detail over the Store-backed OutboxSink. This is
// the independent admin API per Root I-012 — deliberately NOT folded into
// /api/settings/*, no PATCH surface. Provider types never appear here: the
// handlers consume the retrieval interface below, not a channel client.
package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/magicvr/schema-ui-core/apps/api/internal/mail"
)

// OutboxReader is the retrieval face of the outbound record store (satisfied
// by *mail.OutboxSink). Keeping it an interface mirrors the
// repository-consumer convention and keeps the handler testable without a
// store. The listing contract (filters/sort/pagination) lives on
// mail.OutboxListQuery — the store owns its query vocabulary.
type OutboxReader interface {
	List(ctx context.Context, query mail.OutboxListQuery) ([]mail.OutboxRecord, int, error)
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

// parseOutboxQuery maps request query params onto the mail.OutboxListQuery
// contract. page/pageSize follow the generic resource table convention; the
// store normalizes unknown enum/sort values and page-size bounds.
func parseOutboxQuery(r *http.Request) mail.OutboxListQuery {
	return mail.OutboxListQuery{
		Page:           queryInt(r, "page", 1),
		PageSize:       queryInt(r, "pageSize", mail.DefaultOutboxPageSize),
		Q:              strings.TrimSpace(r.URL.Query().Get("q")),
		Channel:        r.URL.Query().Get("channel"),
		DeliveryStatus: r.URL.Query().Get("delivery_status"),
		Sort:           r.URL.Query().Get("sort"),
		Order:          r.URL.Query().Get("order"),
	}
}

// outboxList serves GET /api/mail/outbox with the unified list envelope
// {items,total,page,pageSize} (I-010-001 §3). W27 (GOAL-039 D-001 §2):
// q/channel/delivery_status filters + created_at sort; pagination switched
// from legacy limit/offset to page/pageSize (no remaining consumers).
func outboxList(reader OutboxReader) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := parseOutboxQuery(r)
		items, total, err := reader.List(r.Context(), query)
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not load outbound records")
			return
		}
		if items == nil {
			items = []mail.OutboxRecord{}
		}
		writeJSON(w, http.StatusOK, resourceList{Items: toMapItems(items), Total: total, Page: query.Page, PageSize: query.PageSize})
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

// toMapItems adapts typed records into the generic envelope items. Since W26
// (GOAL-038 D-001 §2.1) items carry the full record — channel, delivery
// status and body ride the list so the declarative recordView drawer can
// render detail from the selected row (bounded by retention + page size).
func toMapItems(items []mail.OutboxRecord) []map[string]any {
	out := make([]map[string]any, 0, len(items))
	for _, rec := range items {
		out = append(out, map[string]any{
			"id":             rec.ID,
			"to":             rec.To,
			"subject":        rec.Subject,
			"body":           rec.Body,
			"channel":        rec.Channel,
			"delivery_status": rec.DeliveryStatus,
			"created_at":     rec.CreatedAt,
		})
	}
	return out
}
