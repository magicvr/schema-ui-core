package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/errorcatalog"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	authsession "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
	"github.com/magicvr/schema-ui-core/apps/api/internal/requestid"
)

const (
	serviceCredentialReadPermission  = "service-credentials.read"
	serviceCredentialWritePermission = "service-credentials.write"
	maxServiceCredentialLifetime     = 365 * 24 * time.Hour
)

type ServiceCredentialRepository interface {
	CreateServiceCredential(authsession.ServiceCredential, authsession.ServiceCredentialAudit) error
	ListServiceCredentials(int, int) ([]authsession.ServiceCredential, int, error)
	ServiceCredentialByID(string) (*authsession.ServiceCredential, error)
	RevokeServiceCredential(string, time.Time, authsession.ServiceCredentialRevokeAudit) (*authsession.ServiceCredential, bool, error)
	ValidatePermissionKeys([]string) error
}

type ServiceCredentialOperations interface {
	operationlog.Recorder
	operationlog.TransactionalRecorder
}

// ServiceCredentialRoutes exposes the central, human-only management surface.
// It contributes no page, navigation node, profile gate, or manifest fragment.
func ServiceCredentialRoutes(a *auth.Authenticator, repository ServiceCredentialRepository, operations ServiceCredentialOperations, moduleID string) []kernel.RouteContribution {
	h := &serviceCredentialHandler{repository: repository, operations: operations, now: time.Now}
	wrap := func(method, pattern string, handler http.Handler) kernel.RouteContribution {
		return kernel.RouteContribution{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey(method, pattern)},
			Method:               method, Pattern: pattern, Handler: a.Middleware(handler),
		}
	}
	return []kernel.RouteContribution{
		wrap(http.MethodGet, "/api/service-credentials", h.list()),
		wrap(http.MethodGet, "/api/service-credentials/{id}", h.detail()),
		wrap(http.MethodPost, "/api/service-credentials", h.create()),
		wrap(http.MethodPost, "/api/service-credentials/{id}/revoke", h.revoke()),
	}
}

type serviceCredentialHandler struct {
	repository ServiceCredentialRepository
	operations ServiceCredentialOperations
	now        func() time.Time
}

type serviceCredentialCreateBody struct {
	Name      string   `json:"name"`
	Scopes    []string `json:"scopes"`
	ExpiresAt string   `json:"expiresAt"`
}

func (h *serviceCredentialHandler) humanPermission(w http.ResponseWriter, r *http.Request, permission string) (account.User, bool) {
	identity, ok := auth.IdentityFrom(r.Context())
	if !ok {
		writeLocalizedError(w, r, http.StatusUnauthorized, "UNAUTHENTICATED", "no active session")
		return account.User{}, false
	}
	if identity.IsServiceCredential() {
		writeLocalizedError(w, r, http.StatusForbidden, "FORBIDDEN", "service credentials cannot manage credentials")
		return account.User{}, false
	}
	if !slices.Contains(identity.Permissions, permission) {
		writeLocalizedError(w, r, http.StatusForbidden, "FORBIDDEN", "permission required: "+permission)
		return account.User{}, false
	}
	return identity, true
}

func (h *serviceCredentialHandler) list() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := h.humanPermission(w, r, serviceCredentialReadPermission); !ok {
			return
		}
		page, ok := intParam(r.URL.Query().Get("page"), 1)
		if !ok {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PAGE", "page must be a positive integer")
			return
		}
		pageSize, ok := intParam(r.URL.Query().Get("pageSize"), DefaultPageSize)
		if !ok || pageSize > maxPageSize {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PAGE_SIZE", "pageSize must be a positive integer not exceeding 100")
			return
		}
		items, total, err := h.repository.ListServiceCredentials(page, pageSize)
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not list service credentials")
			return
		}
		rows := make([]map[string]any, 0, len(items))
		for _, item := range items {
			rows = append(rows, serviceCredentialRow(item, h.now().UTC()))
		}
		writeJSON(w, http.StatusOK, map[string]any{"items": rows, "total": total, "page": page, "pageSize": pageSize})
	})
}

func (h *serviceCredentialHandler) detail() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := h.humanPermission(w, r, serviceCredentialReadPermission); !ok {
			return
		}
		credential, err := h.repository.ServiceCredentialByID(r.PathValue("id"))
		if errors.Is(err, authsession.ErrNotFound) {
			writeLocalizedError(w, r, http.StatusNotFound, "NOT_FOUND", "service credential not found")
			return
		}
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not load service credential")
			return
		}
		writeJSON(w, http.StatusOK, serviceCredentialRow(*credential, h.now().UTC()))
	})
}

func (h *serviceCredentialHandler) create() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actor, ok := h.humanPermission(w, r, serviceCredentialWritePermission)
		if !ok {
			return
		}
		var body serviceCredentialCreateBody
		r.Body = http.MaxBytesReader(w, r.Body, maxResourceBodyBytes)
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&body); err != nil {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_CREATE_BODY", "body must be a JSON object")
			return
		}
		body.Name = strings.TrimSpace(body.Name)
		if len(body.Name) == 0 || len(body.Name) > 100 {
			writeLocalizedFieldError(w, r, http.StatusBadRequest, "INVALID_CREATE_FIELD", "name must contain 1 to 100 characters", []errorcatalog.FieldError{{Field: "name", Reason: "must contain 1 to 100 characters"}})
			return
		}
		body.Scopes = normalizedCredentialScopes(body.Scopes)
		if len(body.Scopes) == 0 || len(body.Scopes) > 64 {
			writeLocalizedFieldError(w, r, http.StatusBadRequest, "INVALID_CREATE_FIELD", "scopes must contain 1 to 64 unique permissions", []errorcatalog.FieldError{{Field: "scopes", Reason: "must contain 1 to 64 unique permissions"}})
			return
		}
		if err := h.repository.ValidatePermissionKeys(body.Scopes); errors.Is(err, authsession.ErrInvalidPermission) {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PERMISSION_REF", "scopes contain an unknown permission")
			return
		} else if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not validate credential scopes")
			return
		}
		for _, scope := range body.Scopes {
			if scope == serviceCredentialReadPermission || scope == serviceCredentialWritePermission || !slices.Contains(actor.Permissions, scope) {
				writeLocalizedError(w, r, http.StatusForbidden, "FORBIDDEN", "credential scopes exceed the creator's delegable permissions")
				return
			}
		}
		now := h.now().UTC()
		expiresAt, err := time.Parse(time.RFC3339, strings.TrimSpace(body.ExpiresAt))
		if err != nil || !expiresAt.After(now) || expiresAt.After(now.Add(maxServiceCredentialLifetime)) {
			writeLocalizedFieldError(w, r, http.StatusBadRequest, "INVALID_CREATE_FIELD", "expiresAt must be a future RFC3339 timestamp within 365 days", []errorcatalog.FieldError{{Field: "expiresAt", Reason: "must be a future RFC3339 timestamp within 365 days"}})
			return
		}
		raw, tokenHash, tokenPrefix, err := auth.NewServiceCredentialToken()
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not create service credential")
			return
		}
		credential := authsession.ServiceCredential{
			ID: auth.NewServiceCredentialID(), Name: body.Name, TokenPrefix: tokenPrefix, TokenHash: tokenHash,
			Scopes: body.Scopes, ExpiresAt: expiresAt.UTC(), CreatedBy: actor.ID, CreatedAt: now, UpdatedAt: now,
		}
		detail, err := operationlog.NewDetail("service-credential-create", nil, serviceCredentialAuditRow(credential))
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not create service credential")
			return
		}
		err = h.repository.CreateServiceCredential(credential, func(tx kernel.Tx) error {
			return h.operations.RecordOperationTx(tx, operationlog.Operation{
				ID: newOperationID(), Event: operationlog.EventServiceCredentialCreate,
				ActorID: actor.ID, ActorName: actor.Name, RecordID: &credential.ID, Detail: &detail,
				CorrelationID: requestid.FromContext(r.Context()), SessionID: identitySession(actor), CreatedAt: now,
			})
		})
		if errors.Is(err, authsession.ErrCredentialNameTaken) {
			writeLocalizedFieldError(w, r, http.StatusBadRequest, "INVALID_CREATE_FIELD", "name already exists", []errorcatalog.FieldError{{Field: "name", Reason: "name already exists"}})
			return
		}
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not create service credential")
			return
		}
		response := serviceCredentialRow(credential, now)
		response["secret"] = raw
		writeJSON(w, http.StatusCreated, response)
	})
}

func (h *serviceCredentialHandler) revoke() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actor, ok := h.humanPermission(w, r, serviceCredentialWritePermission)
		if !ok {
			return
		}
		now := h.now().UTC()
		_, _, err := h.repository.RevokeServiceCredential(r.PathValue("id"), now, func(tx kernel.Tx, credential authsession.ServiceCredential) error {
			detail, detailErr := operationlog.NewDetail("service-credential-revoke", serviceCredentialAuditRow(credential), map[string]any{"revokedAt": now.Format(time.RFC3339)})
			if detailErr != nil {
				return detailErr
			}
			return h.operations.RecordOperationTx(tx, operationlog.Operation{
				ID: newOperationID(), Event: operationlog.EventServiceCredentialRevoke,
				ActorID: actor.ID, ActorName: actor.Name, RecordID: &credential.ID, Detail: &detail,
				CorrelationID: requestid.FromContext(r.Context()), SessionID: identitySession(actor), CreatedAt: now,
			})
		})
		if errors.Is(err, authsession.ErrNotFound) {
			writeLocalizedError(w, r, http.StatusNotFound, "NOT_FOUND", "service credential not found")
			return
		}
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not revoke service credential")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
}

func normalizedCredentialScopes(scopes []string) []string {
	result := make([]string, 0, len(scopes))
	for _, scope := range scopes {
		scope = strings.TrimSpace(scope)
		if scope != "" && !slices.Contains(result, scope) {
			result = append(result, scope)
		}
	}
	slices.Sort(result)
	return result
}

func serviceCredentialRow(credential authsession.ServiceCredential, now time.Time) map[string]any {
	status := "active"
	if credential.RevokedAt != nil {
		status = "revoked"
	} else if !now.Before(credential.ExpiresAt) {
		status = "expired"
	}
	row := serviceCredentialAuditRow(credential)
	row["status"] = status
	row["revokedAt"] = nil
	row["lastUsedAt"] = nil
	if credential.RevokedAt != nil {
		row["revokedAt"] = credential.RevokedAt.UTC().Format(time.RFC3339)
	}
	if credential.LastUsedAt != nil {
		row["lastUsedAt"] = credential.LastUsedAt.UTC().Format(time.RFC3339)
	}
	row["createdAt"] = credential.CreatedAt.UTC().Format(time.RFC3339)
	row["updatedAt"] = credential.UpdatedAt.UTC().Format(time.RFC3339)
	return row
}

func serviceCredentialAuditRow(credential authsession.ServiceCredential) map[string]any {
	return map[string]any{
		"id": credential.ID, "name": credential.Name, "tokenPrefix": credential.TokenPrefix,
		"scopes": credential.Scopes, "expiresAt": credential.ExpiresAt.UTC().Format(time.RFC3339),
		"createdBy": credential.CreatedBy,
	}
}
