// F-02 schema-driven CSV export (GOAL-004 D-002 `3): GET /api/export/{resource}
// streams the filtered resource list as CSV (UTF-8 BOM, RFC 4180 escaping,
// attachment disposition). R2 resources: users, roles. Bounded by the same
// list filter as the resource endpoints (q/sort/order) with an explicit row
// cap (10000) — a full unbound dump is out of scope by design.
package handler

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	authsession "github.com/magicvr/schema-ui-core/apps/api/modules/authsession"
	"github.com/magicvr/schema-ui-core/apps/api/modules/operationlog"
)

// maxExportRows is the documented export cap (D-002 `3): full exports beyond
// this require additional filtering. Kept in sync with the frozen limit.
const maxExportRows = 10000

// ExportUsersRepository is the users list surface used by the export endpoint.
type ExportUsersRepository interface {
	ListUsers(authsession.UserFilter) ([]authsession.User, int, error)
}

// ExportRolesRepository is the roles list surface used by the export endpoint.
type ExportRolesRepository interface {
	ListRoles(authsession.RoleFilter) ([]authsession.Role, int, error)
}

// ExportRoutes returns the export route contributions (admin.data-transfer).
func ExportRoutes(a *auth.Authenticator, usersRepo ExportUsersRepository, rolesRepo ExportRolesRepository, operations operationlog.Recorder, moduleID string) []kernel.RouteContribution {
	h := &exportHandler{users: usersRepo, roles: rolesRepo, operations: operations, now: time.Now}
	return []kernel.RouteContribution{
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey("GET", "/api/export/{resource}")},
			Method:               "GET",
			Pattern:              "/api/export/{resource}",
			Handler:              a.Middleware(exportPermissionGate(h.export())),
		},
	}
}

// exportPermissionGate wraps the export handler with the data.export gate.
func exportPermissionGate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePermission(w, r, "data.export"); !ok {
			return
		}
		next.ServeHTTP(w, r)
	})
}

type exportHandler struct {
	users      ExportUsersRepository
	roles      ExportRolesRepository
	operations operationlog.Recorder
	now        func() time.Time
}

// exportRow renders one resource row as an ordered CSV record. Column order is
// stable per resource (frozen in D-002 `3); array fields serialize as JSON.
func exportRow(resource string, row map[string]any) []string {
	switch resource {
	case "users":
		return []string{
			formulaSafe(stringField(row, "id")), formulaSafe(stringField(row, "username")), formulaSafe(stringField(row, "name")),
			formulaSafe(jsonArrayOr(row["roles"])), boolString(row["enabled"]), boolString(row["locked"]),
			stringField(row, "createdAt"), stringField(row, "updatedAt"),
		}
	case "roles":
		return []string{
			formulaSafe(stringField(row, "id")), formulaSafe(stringField(row, "key")), formulaSafe(stringField(row, "name")),
			boolString(row["system"]), formulaSafe(jsonArrayOr(row["permissions"])), formulaSafe(jsonArrayOr(row["menuItems"])),
			strconv.Itoa(intOf(row["assignedUsers"])), boolString(row["editable"]), boolString(row["deletable"]),
			stringField(row, "createdAt"), stringField(row, "updatedAt"),
		}
	default:
		return nil
	}
}

func jsonArrayOr(v any) string {
	switch value := v.(type) {
	case []string:
		raw, err := json.Marshal(value)
		if err != nil {
			return "[]"
		}
		return string(raw)
	case []any:
		raw, err := json.Marshal(value)
		if err != nil {
			return "[]"
		}
		return string(raw)
	default:
		return ""
	}
}

func boolString(v any) string {
	if b, ok := v.(bool); ok {
		return strconv.FormatBool(b)
	}
	return ""
}

func intOf(v any) int {
	if n, ok := v.(int); ok {
		return n
	}
	return 0
}

func (h *exportHandler) export() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resource := r.PathValue("resource")
		if resource != "users" && resource != "roles" {
			writeLocalizedError(w, r, http.StatusNotFound, "RESOURCE_NOT_FOUND", "no export for that resource")
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

		var headers []string
		var rows [][]string
		switch resource {
		case "users":
			sortField := query.Get("sort")
			if sortField == "" {
				sortField = "username"
			}
			if !slicesContains([]string{"username", "name", "updatedAt"}, sortField) {
				writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_SORT_FIELD", "unsupported sort field")
				return
			}
			items, _, err := h.users.ListUsers(authsession.UserFilter{
				Q: strings.ToLower(strings.TrimSpace(query.Get("q"))),
				Sort: sortField, Order: order, Page: 1, PageSize: pageSize,
			})
			if err != nil {
				writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not export users")
				return
			}
			headers = []string{"id", "username", "name", "roles", "enabled", "locked", "createdAt", "updatedAt"}
			for _, u := range items {
				rows = append(rows, exportRow("users", userToMap(u)))
			}
		case "roles":
			sortField := query.Get("sort")
			if sortField == "" {
				sortField = "key"
			}
			if !slicesContains([]string{"key", "name", "updatedAt"}, sortField) {
				writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_SORT_FIELD", "unsupported sort field")
				return
			}
			items, _, err := h.roles.ListRoles(authsession.RoleFilter{
				Q: strings.ToLower(strings.TrimSpace(query.Get("q"))),
				Sort: sortField, Order: order, Page: 1, PageSize: pageSize,
			})
			if err != nil {
				writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not export roles")
				return
			}
			headers = []string{"id", "key", "name", "system", "permissions", "menuItems", "assignedUsers", "editable", "deletable", "createdAt", "updatedAt"}
			for _, role := range items {
				rows = append(rows, exportRow("roles", roleToMap(role)))
			}
		}

		var out strings.Builder
		out.WriteString("\uFEFF") // UTF-8 BOM so Excel detects UTF-8
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
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", resource+".csv"))
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(out.String()))

		h.record(operationlog.EventDataExport, r, resource, len(rows))
	})
}

func (h *exportHandler) record(event string, r *http.Request, resource string, rows int) {
	user, ok := auth.IdentityFrom(r.Context())
	if !ok {
		return
	}
	recordAudit(h.operations, user, event, "", auditDetail("export", map[string]any{"resource": resource, "rows": rows}), h.now().UTC(), r.Context())
}

// formulaSafe neutralizes spreadsheet formula injection (F-009 / W11 F-017):
// a cell that starts with = + - @ is prefixed with a single quote (visible
// in Excel, not executed as a formula). W11 F-017 adds tab and carriage
// return — OWASP lists them as formula-injection prefixes Excel/LibreOffice
// accept before the operator.
func formulaSafe(value string) string {
	if value == "" {
		return value
	}
	switch value[0] {
	case '=', '+', '-', '@', '\t', '\r':
		return "'" + value
	}
	return value
}

func slicesContains(list []string, value string) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}
	return false
}