// File library surface (S-02 · GOAL-007 D-002): the admin.file-library module
// manages the shared upload store (C-09) as a unified file/attachment library.
// The list/detail read surfaces enumerate the uploads namespace through the
// kernel object-storage port — the single source of truth is the stored
// objects plus their owner metadata (same pattern as the per-user quota
// scan); there is no DB mirror table (D-002 `2). Downloads are
// gated by files.read (management surface, unlike the owner-only control
// endpoint GET /api/files/{id}), deletion by files.delete, and the upload
// confirmation endpoint validates ownership and records the audit event.
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
	"github.com/magicvr/schema-ui-core/apps/api/internal/pagination"
)

// fileLibraryEntity adapts the shared upload store to the generic resource
// factory (read-only list/detail; write methods are defensive fallbacks).
type fileLibraryEntity struct {
	objects    kernel.ObjectStore
	operations operationlog.Recorder
}

// fileRow is the library projection of one stored upload.
type fileRow struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Size    int64  `json:"size"`
	Owner   string `json:"owner"`
	Created string `json:"created"`
}

func (e *fileLibraryEntity) List(filter resourceFilter) ([]map[string]any, int, error) {
	rows, err := e.scan()
	if err != nil {
		return nil, 0, err
	}
	if q := strings.ToLower(strings.TrimSpace(filter.Q)); q != "" {
		kept := rows[:0]
		for _, row := range rows {
			haystack := strings.ToLower(row.Name + " " + row.Type + " " + row.Owner)
			if strings.Contains(haystack, q) {
				kept = append(kept, row)
			}
		}
		rows = kept
	}
	sortRows(rows, filter.Sort, filter.Order)
	total := len(rows)
	start, end := pagination.Bounds(filter.Page, filter.PageSize, total)
	items := make([]map[string]any, 0, end-start)
	for _, row := range rows[start:end] {
		items = append(items, fileRowToMap(row))
	}
	return items, total, nil
}

func (e *fileLibraryEntity) Get(id string) (map[string]any, error) {
	if !uploadFileIDPattern.MatchString(id) {
		return nil, errResourceNotFound
	}
	_, meta, err := e.objects.Get(context.Background(), kernel.ObjectNamespaceUploads, id)
	if err != nil {
		if errors.Is(err, kernel.ErrObjectNotFound) {
			return nil, errResourceNotFound
		}
		return nil, err
	}
	info, err := e.objects.Stat(context.Background(), kernel.ObjectNamespaceUploads, id)
	if err == nil {
		row := fileRow{ID: id, Name: meta.Name, Type: meta.Type, Owner: meta.Owner, Size: info.Size, Created: formatRFC3339Milli(info.ModTime)}
		return fileRowToMap(row), nil
	}
	if !errors.Is(err, kernel.ErrObjectNotFound) {
		return nil, err
	}
	row := fileRow{ID: id, Name: meta.Name, Type: meta.Type, Owner: meta.Owner}
	return fileRowToMap(row), nil
}

func (e *fileLibraryEntity) Create(map[string]any, string, time.Time, account.User) (map[string]any, error) {
	return nil, errReadOnlyResource
}

func (e *fileLibraryEntity) Update(string, map[string]any, time.Time, account.User) (map[string]any, error) {
	return nil, errReadOnlyResource
}

func (e *fileLibraryEntity) Delete(string, account.User) error {
	return errReadOnlyResource
}

// scan lists every stored object with its owner metadata through the port
// (List + Stat). Rows with completely empty metadata are skipped — the
// legacy read surface only saw objects with a meta sidecar, so body-only
// leftovers stay invisible (GOAL-004 D-001 §2).
func (e *fileLibraryEntity) scan() ([]fileRow, error) {
	ids, err := e.objects.List(context.Background(), kernel.ObjectNamespaceUploads)
	if err != nil {
		return nil, err
	}
	rows := make([]fileRow, 0, len(ids))
	for _, id := range ids {
		info, err := e.objects.Stat(context.Background(), kernel.ObjectNamespaceUploads, id)
		if err != nil {
			continue // best-effort read surface: skip unreadable entries
		}
		meta := info.Meta
		if meta.Name == "" && meta.Type == "" && meta.Owner == "" {
			continue // body-only leftover: invisible, as before the migration
		}
		row := fileRow{ID: id, Name: meta.Name, Type: meta.Type, Owner: meta.Owner, Size: info.Size, Created: formatRFC3339Milli(info.ModTime)}
		rows = append(rows, row)
	}
	return rows, nil
}

func sortRows(rows []fileRow, field, order string) {
	less := func(i, j int) bool {
		switch field {
		case "type":
			return rows[i].Type < rows[j].Type
		case "owner":
			return rows[i].Owner < rows[j].Owner
		case "size":
			return rows[i].Size < rows[j].Size
		case "created":
			return rows[i].Created < rows[j].Created
		default:
			return rows[i].Name < rows[j].Name
		}
	}
	sort.Slice(rows, func(i, j int) bool {
		if order == "desc" {
			return less(j, i)
		}
		return less(i, j)
	})
}

func fileRowToMap(row fileRow) map[string]any {
	return map[string]any{
		"id":          row.ID,
		"name":        row.Name,
		"type":        row.Type,
		"size":        row.Size,
		"owner":       row.Owner,
		"created":     row.Created,
		"downloadUrl": "/api/library/files/" + row.ID + "/download",
	}
}

// FileLibraryRoutes returns the admin.file-library HTTP surface (D-002 `4).
func FileLibraryRoutes(a *auth.Authenticator, objects kernel.ObjectStore, operations operationlog.Recorder, moduleID string) []kernel.RouteContribution {
	entity := &fileLibraryEntity{objects: objects, operations: operations}
	res := Resource{
		ID:              "files",
		Path:            "/api/library/files",
		Listable:        true,
		ReadOnly:        true,
		SortFields:      []string{"name", "type", "owner", "size", "created"},
		QSearch:         true,
		Entity:          entity,
		PermissionRead:  "files.read",
		PermissionWrite: "files.delete",
		NotFoundCode:    "FILE_NOT_FOUND",
	}
	routes := ResourceRoutes(a, res, moduleID)
	store := &uploadStore{objects: objects}
	now := time.Now

	routes = append(routes, kernel.RouteContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey("DELETE", "/api/library/files/{id}")},
		Method:               "DELETE",
		Pattern:              "/api/library/files/{id}",
		Handler: a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := requirePermission(w, r, "files.delete")
			if !ok {
				return
			}
			id := r.PathValue("id")
			if !uploadFileIDPattern.MatchString(id) {
				writeLocalizedError(w, r, http.StatusNotFound, "FILE_NOT_FOUND", "no file with that id")
				return
			}
			// A-003 F-002 (port adaptation, GOAL-004 D-001 §2): Delete removes the
			// object and its metadata together; a missing object is a port no-op.
			// A ghost id (body-less sidecar leftover) still gets one best-effort
			// Delete so it cannot ghost-list in scans, but the response is 404 —
			// the documented divergence from the pre-port 204.
			exists, err := store.objects.Exists(r.Context(), kernel.ObjectNamespaceUploads, id)
			if err != nil {
				writeLocalizedError(w, r, http.StatusInternalServerError, "STORAGE_UNAVAILABLE", "could not delete file")
				return
			}
			if err := store.objects.Delete(r.Context(), kernel.ObjectNamespaceUploads, id); err != nil {
				writeLocalizedError(w, r, http.StatusInternalServerError, "STORAGE_UNAVAILABLE", "could not delete file")
				return
			}
			if !exists {
				writeLocalizedError(w, r, http.StatusNotFound, "FILE_NOT_FOUND", "no file with that id")
				return
			}
			recordFileEvent(operations, operationlog.EventFileDelete, user, id, now())
			w.WriteHeader(http.StatusNoContent)
		})),
	})

	routes = append(routes, kernel.RouteContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey("GET", "/api/library/files/{id}/download")},
		Method:               "GET",
		Pattern:              "/api/library/files/{id}/download",
		Handler: a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := requirePermission(w, r, "files.read")
			if !ok {
				return
			}
			id := r.PathValue("id")
			body, meta, err := store.load(id)
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					writeLocalizedError(w, r, http.StatusNotFound, "FILE_NOT_FOUND", "no file with that id")
					return
				}
				writeLocalizedError(w, r, http.StatusInternalServerError, "STORAGE_UNAVAILABLE", "could not read file")
				return
			}
			recordFileEvent(operations, operationlog.EventFileDownload, user, id, now())
			contentType := meta.Type
			if contentType == "" {
				contentType = "application/octet-stream"
			}
			w.Header().Set("Content-Type", contentType)
			w.Header().Set("X-Content-Type-Options", "nosniff")
			// Sanitized literal filename: the stored name is attacker-controlled
			// and must never reach a header unsanitized (same posture as the
			// control download endpoint).
			w.Header().Set("Content-Disposition", `attachment; filename="`+sanitizeFileHeaderName(meta.Name)+`"`)
			w.Header().Set("Content-Security-Policy", "sandbox")
			w.Header().Set("Content-Length", strconv.Itoa(len(body)))
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(body)
		})),
	})

	routes = append(routes, kernel.RouteContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey("POST", "/api/library/files/upload")},
		Method:               "POST",
		Pattern:              "/api/library/files/upload",
		Handler: a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := requirePermission(w, r, "files.write")
			if !ok {
				return
			}
			var body struct {
				File string `json:"file"`
			}
			decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxResourceBodyBytes))
			if err := decoder.Decode(&body); err != nil {
				writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_UPLOAD_BODY", "expected a JSON object with a file field")
				return
			}
			file := strings.TrimSpace(body.File)
			if strings.HasPrefix(file, "/api/files/") {
				file = strings.TrimPrefix(file, "/api/files/")
			}
			if !uploadFileIDPattern.MatchString(file) {
				writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_FILE_ID", "invalid file id")
				return
			}
			_, meta, err := store.load(file)
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					writeLocalizedError(w, r, http.StatusNotFound, "FILE_NOT_FOUND", "no file with that id")
					return
				}
				writeLocalizedError(w, r, http.StatusInternalServerError, "STORAGE_UNAVAILABLE", "could not read file")
				return
			}
			// The upload control stores under the current user, so the owner must
			// match; a foreign id cannot be confirmed into the library.
			if strings.TrimSpace(meta.Owner) != user.ID {
				writeLocalizedError(w, r, http.StatusForbidden, "FORBIDDEN", "not the owner of this file")
				return
			}
			recordFileEvent(operations, operationlog.EventFileUpload, user, file, now())
			writeJSON(w, http.StatusOK, map[string]any{"id": file})
		})),
	})
	return routes
}

// sanitizeFileHeaderName reduces a stored file name to header-safe characters
// (quotes/CR/LF and everything else are replaced).
func sanitizeFileHeaderName(name string) string {
	var b strings.Builder
	for _, r := range strings.TrimSpace(name) {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '.' || r == '_' || r == '-' {
			b.WriteRune(r)
		} else {
			b.WriteRune('_')
		}
	}
	out := strings.Trim(b.String(), "._-")
	if out == "" {
		return "download"
	}
	if len(out) > 100 {
		out = out[:100]
	}
	return out
}

func recordFileEvent(operations operationlog.Recorder, event string, user account.User, id string, now time.Time) {
	recordAudit(operations, user, event, id, nil, now.UTC(), nil)
}
