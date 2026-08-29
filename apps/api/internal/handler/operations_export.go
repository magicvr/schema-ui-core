// W14 F-03 (GOAL-016): CSV export for the activity/audit log. Uses the same
// structured filters as the operations list endpoint and the same export
// hardening as the generic data-transfer CSV export (UTF-8 BOM, RFC 4180
// escaping, formula injection guard, row cap).
package handler

import (
	"encoding/csv"
	"errors"
	"net/http"
	"strings"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/modules/operationlog"
)

// OperationsExportRoute returns GET /api/operations/export (operations.read).
func OperationsExportRoute(a *auth.Authenticator, repository operationlog.Reader, moduleID string) kernel.RouteContribution {
	return kernel.RouteContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey("GET", "/api/operations/export")},
		Method:               "GET",
		Pattern:              "/api/operations/export",
		Handler: a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ok := requirePermission(w, r, "operations.read"); !ok {
				return
			}
			query := r.URL.Query()
			pageSize, ok := intParam(query.Get("pageSize"), maxExportRows)
			if !ok || pageSize > maxExportRows {
				writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_EXPORT_LIMIT", "pageSize must not exceed 10000")
				return
			}
			order := query.Get("order")
			if order == "" {
				order = "asc"
			}
			if order != "asc" && order != "desc" {
				writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_SORT_ORDER", "order must be asc or desc")
				return
			}
			sortField := query.Get("sort")
			if sortField == "" {
				sortField = "createdAt"
			}
			if !slicesContains([]string{"createdAt", "event", "actorName"}, sortField) {
				writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_SORT_FIELD", "unsupported sort field")
				return
			}
			extra := map[string]string{}
			for _, key := range []string{"event", "actorName", "from", "to"} {
				if value := strings.TrimSpace(query.Get(key)); value != "" {
					extra[key] = value
				}
			}
			filter, err := operationFilterFromResource(resourceFilter{
				Q:    strings.ToLower(strings.TrimSpace(query.Get("q"))),
				Sort: sortField, Order: order, Page: 1, PageSize: pageSize,
				Extra: extra,
			})
			if err != nil {
				var domainErr *DomainError
				if errors.As(err, &domainErr) {
					writeLocalizedError(w, r, domainErr.Status, domainErr.Code, domainErr.Message)
					return
				}
				writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not export operations")
				return
			}
			items, _, err := repository.ListOperationsFiltered(filter)
			if err != nil {
				writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not export operations")
				return
			}

			headers := []string{"id", "event", "actorId", "actorName", "recordId", "detail", "correlationId", "sessionId", "createdAt"}
			rows := make([][]string, 0, len(items))
			for _, op := range items {
				row := operationToMap(op)
				rows = append(rows, []string{
					formulaSafe(stringField(row, "id")),
					formulaSafe(stringField(row, "event")),
					formulaSafe(stringField(row, "actorId")),
					formulaSafe(stringField(row, "actorName")),
					formulaSafe(stringField(row, "recordId")),
					formulaSafe(stringField(row, "detail")),
					formulaSafe(stringField(row, "correlationId")),
					formulaSafe(stringField(row, "sessionId")),
					formulaSafe(stringField(row, "createdAt")),
				})
			}

			var out strings.Builder
			out.WriteString("\uFEFF") // UTF-8 BOM for Excel
			writer := csv.NewWriter(&out)
			if err := writer.Write(headers); err != nil {
				writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not render export")
				return
			}
			for _, row := range rows {
				if err := writer.Write(row); err != nil {
					writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not render export")
					return
				}
			}
			writer.Flush()
			if err := writer.Error(); err != nil {
				writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not render export")
				return
			}

			w.Header().Set("Content-Type", "text/csv; charset=utf-8")
			w.Header().Set("Content-Disposition", `attachment; filename="operations.csv"`)
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(out.String()))
		})),
	}
}
