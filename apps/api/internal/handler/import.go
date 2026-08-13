// F-02 CSV import (GOAL-004 D-002 `4): POST /api/import/users consumes an
// uploaded CSV file id (existing upload infra C-09, owner-scoped) and applies
// valid rows with per-row error reporting. Explicit no-rollback semantics:
// valid rows commit even when other rows fail; the response carries the full
// per-row error report ({applied, failed, total, errors}).
package handler

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	authsession "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
)

// maxImportBytes bounds the CSV body read from the uploaded file (D-002 `4).
const maxImportBytes = 2 << 20

// ImportUsersRepository is the users write surface used by the import endpoint.
type ImportUsersRepository interface {
	CreateUserManagement(authsession.User) (*authsession.User, error)
	PermissionsForRoles([]string) ([]string, error)
}

// ImportRoutes returns the import route contributions (admin.data-transfer).
func ImportRoutes(a *auth.Authenticator, repo ImportUsersRepository, operations operationlog.Recorder, uploadDir, moduleID string) []kernel.RouteContribution {
	h := &importHandler{repository: repo, operations: operations, uploadDir: uploadDir, now: time.Now}
	return []kernel.RouteContribution{
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey("POST", "/api/import/{resource}")},
			Method:               "POST",
			Pattern:              "/api/import/{resource}",
			Handler:              a.Middleware(importPermissionGate(h.importResource())),
		},
	}
}

// importPermissionGate wraps the import handler with the data.import gate.
func importPermissionGate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePermission(w, r, "data.import"); !ok {
			return
		}
		next.ServeHTTP(w, r)
	})
}

type importHandler struct {
	repository ImportUsersRepository
	operations operationlog.Recorder
	uploadDir  string
	now        func() time.Time
}

type importRowError struct {
	Row     int    `json:"row"`
	Message string `json:"message"`
}

type importResult struct {
	Applied int              `json:"applied"`
	Failed  int              `json:"failed"`
	Total   int              `json:"total"`
	Errors  []importRowError `json:"errors"`
}

func (h *importHandler) importResource() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resource := r.PathValue("resource")
		if resource != "users" {
			// R2 scope: users only (roles import deferred; D-002 `1).
			writeLocalizedError(w, r, http.StatusNotFound, "RESOURCE_NOT_FOUND", "no import for that resource")
			return
		}
		var body struct {
			FileID string `json:"fileId"`
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxResourceBodyBytes)
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || strings.TrimSpace(body.FileID) == "" {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_IMPORT_BODY", "body must be JSON with fileId")
			return
		}
		user, _ := auth.IdentityFrom(r.Context())
		// The upload field value is url-preferred (/api/files/{id}); normalize
		// both forms to the stored id.
		if user.ID == "" {
			writeLocalizedError(w, r, http.StatusUnauthorized, "UNAUTHENTICATED", "no active session")
			return
		}
		fileID := strings.TrimSpace(body.FileID)
		if strings.HasPrefix(fileID, "/api/files/") {
			fileID = strings.TrimPrefix(fileID, "/api/files/")
		}
		raw, meta, err := loadUploadedFile(h.uploadDir, fileID)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				writeLocalizedError(w, r, http.StatusNotFound, "FILE_NOT_FOUND", "no file with that id")
				return
			}
			writeLocalizedError(w, r, http.StatusInternalServerError, "STORAGE_UNAVAILABLE", "could not read uploaded file")
			return
		}
		// Owner-scoped: only the uploader's own file may be imported (C-09
		// ownership model; missing owner fails closed).
		if owner := strings.TrimSpace(meta["owner"]); owner == "" || owner != user.ID {
			writeLocalizedError(w, r, http.StatusForbidden, "FORBIDDEN", "not the owner of this file")
			return
		}
		if len(raw) == 0 {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_FILE", "empty files are rejected")
			return
		}
		if len(raw) > maxImportBytes {
			writeLocalizedError(w, r, http.StatusRequestEntityTooLarge, "FILE_TOO_LARGE", "file exceeds the import size limit")
			return
		}

		result, err := h.importUsersCSV(raw, user)
		if err != nil {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_CSV", "could not parse CSV: "+err.Error())
			return
		}
		if result.Errors == nil {
			result.Errors = []importRowError{}
		}
		writeJSON(w, http.StatusOK, result)
		h.record(resource, result, user)
	})
}

// importUsersCSV parses and applies a users CSV. Returns a structured result.
// A missing/unreadable header row is a hard error (INVALID_CSV); a malformed
// data row is collected into the per-row error report and parsing stops (the
// remainder of the file is unrecoverable) — applied rows still commit and are
// always audited (F-002).
func (h *importHandler) importUsersCSV(raw []byte, actor account.User) (importResult, error) {
	body := string(raw)
	// F-003: strip a UTF-8 BOM (Excel "UTF-8 CSV" writes one) before parsing.
	body = strings.TrimPrefix(body, "\uFEFF")
	reader := csv.NewReader(strings.NewReader(body))
	reader.FieldsPerRecord = -1 // header-driven; tolerate ragged rows
	header, err := reader.Read()
	if err != nil {
		return importResult{}, fmt.Errorf("missing header row: %w", err)
	}
	columns := make([]string, 0, len(header))
	columnIndex := map[string]int{}
	for _, name := range header {
		trimmed := strings.TrimSpace(name)
		if trimmed == "" {
			continue
		}
		if _, seen := columnIndex[trimmed]; seen {
			continue
		}
		columnIndex[trimmed] = len(columns)
		columns = append(columns, trimmed)
	}
	// R2 supported columns; unknown columns are ignored (documented).
	supported := map[string]bool{"username": true, "name": true, "roles": true, "password": true}

	result := importResult{}
	now := h.now().UTC()
	rowNumber := 1 // header consumed; next read starts at data row 1
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			// F-002: a malformed row is a per-row failure — report it, stop
			// parsing (the rest of the file is unrecoverable), keep applied
			// rows committed and audited.
			rowNumber++
			result.Total++
			result.Failed++
			result.Errors = append(result.Errors, importRowError{Row: rowNumber, Message: "malformed CSV row: " + err.Error()})
			break
		}
		rowNumber++
		result.Total++
		row := map[string]string{}
		unknown := []string{}
		for _, col := range columns {
			value := ""
			if columnIndex[col] < len(record) {
				value = strings.TrimSpace(record[columnIndex[col]])
			}
			if supported[col] {
				row[col] = value
			} else {
				unknown = append(unknown, col)
			}
		}
		if message := validateImportUser(row); message != "" {
			result.Failed++
			result.Errors = append(result.Errors, importRowError{Row: rowNumber, Message: message})
			continue
		}
		roles := []string{}
		if row["roles"] != "" {
			roles = splitCommaList(row["roles"])
		}
		var hash string
		if row["password"] != "" {
			var hashErr error
			hash, hashErr = auth.HashPassword(row["password"], passwordHashCost)
			if hashErr != nil {
				result.Failed++
				result.Errors = append(result.Errors, importRowError{Row: rowNumber, Message: "could not hash password"})
				continue
			}
		}
		// F-001 (A-003 independent): the import write path must enforce the
		// same role-assignment delegation boundary as POST /api/users —
		// roles.assign permission, admin-only admin assignment, and no
		// assignment of roles carrying permissions the actor does not hold.
		// Violations are per-row failures (no-rollback semantics), not a whole
		// request 403.
		if len(roles) > 0 {
			if message := importRoleAssignmentError(h.repository, actor, roles); message != "" {
				result.Failed++
				result.Errors = append(result.Errors, importRowError{Row: rowNumber, Message: message})
				continue
			}
		}
		_, createErr := h.repository.CreateUserManagement(authsession.User{
			ID:           newUserIDValue(),
			Username:     row["username"],
			Name:         row["name"],
			Roles:        roles,
			PasswordHash: hash,
			CreatedAt:    now,
			UpdatedAt:    now,
		})
		if createErr != nil {
			message := "could not create user"
			switch {
			case errors.Is(createErr, authsession.ErrUsernameTaken):
				message = "username already exists"
			case errors.Is(createErr, authsession.ErrInvalidRole):
				message = "roles contain an unknown role key"
			}
			result.Failed++
			result.Errors = append(result.Errors, importRowError{Row: rowNumber, Message: message})
			continue
		}
		result.Applied++
	}
	return result, nil
}

// importRoleAssignmentError mirrors usersEntity.authorizeRoleAssignment for
// the import path (F-001). Returns "" when the assignment is allowed.
func importRoleAssignmentError(repo ImportUsersRepository, actor account.User, roles []string) string {
	if !slices.Contains(actor.Permissions, "roles.assign") {
		return "role assignment forbidden: permission required: roles.assign"
	}
	if slices.Contains(roles, "admin") && !slices.Contains(actor.Roles, "admin") {
		return "role assignment forbidden: only an admin may assign the admin role"
	}
	targetPermissions, err := repo.PermissionsForRoles(roles)
	if err != nil {
		// Unknown role keys are reported by CreateUserManagement as
		// INVALID_ROLE_REF row errors — skip delegation for invalid roles so
		// the row message stays precise (F-001 keeps valid-role delegation).
		if errors.Is(err, authsession.ErrInvalidRole) {
			return ""
		}
		return "role assignment forbidden: could not resolve target role permissions"
	}
	for _, permission := range targetPermissions {
		if !slices.Contains(actor.Permissions, permission) {
			return "role assignment forbidden: cannot assign a role with permissions the actor does not hold"
		}
	}
	return ""
}

func validateImportUser(row map[string]string) string {
	if strings.TrimSpace(row["username"]) == "" {
		return "username is required"
	}
	if strings.TrimSpace(row["name"]) == "" {
		return "name is required"
	}
	if pwd := row["password"]; pwd != "" {
		length := len([]byte(pwd))
		if length < minPasswordBytes || length > maxPasswordBytes || strings.TrimSpace(pwd) == "" {
			return "password must be 8 to 72 bytes"
		}
	} else {
		return "password is required"
	}
	return ""
}

func splitCommaList(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

func newUserIDValue() string {
	id, err := newUserID()
	if err != nil {
		// crypto/rand failure: timestamp fallback so import never wedges.
		return fmt.Sprintf("usr-%d", time.Now().UnixNano())
	}
	return id
}

func (h *importHandler) record(resource string, result importResult, actor account.User) {
	detail := fmt.Sprintf(`{"resource":%s,"applied":%d,"failed":%d}`, jsonQuote(resource), result.Applied, result.Failed)
	op := operationlog.Operation{
		ID:        newOperationID(),
		Event:     operationlog.EventDataImport,
		ActorID:   actor.ID,
		ActorName: actor.Name,
		CreatedAt: h.now().UTC(),
	}
	op.Detail = &detail
	if h.operations == nil {
		return
	}
	if err := h.operations.RecordOperation(op); err != nil {
		slog.Error("operation log write failed", "event", operationlog.EventDataImport, "err", err)
	}
}

// LoadUploadedFile exposes the upload store's load to sibling handlers (import).
func LoadUploadedFile(dir, id string) ([]byte, map[string]string, error) {
	return loadUploadedFile(dir, id)
}

// loadUploadedFile reads a stored upload by id with its owner meta (same
// validation as the download endpoint; unexported here, re-exported above).
func loadUploadedFile(dir, id string) ([]byte, map[string]string, error) {
	store := &uploadStore{dir: dir}
	return store.load(id)
}