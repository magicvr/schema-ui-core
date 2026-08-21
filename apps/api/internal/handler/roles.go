// roles resource instance of the generic resource factory (GOAL-011 S2 ·
// I-011-001 §3): mounted at /api/roles with roles.read / roles.write,
// ROLE_NOT_FOUND, create/patch fields key/name, and operation-log events. The
// entity maps the system/in-use/invalid-key store sentinels to the frozen
// resource-specific error codes. User-created role grants are managed through
// the same resource as the bounded I-011-004 §4 extension.
package handler

import (
	"context"
	"errors"
	"slices"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	authsession "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
)

// RolesRepository is the RBAC-domain persistence required by the roles
// resource surface.
type RolesRepository interface {
	ListRoles(authsession.RoleFilter) ([]authsession.Role, int, error)
	GetRole(string) (*authsession.Role, error)
	CreateRoleWithGrants(string, string, []string, []string, time.Time) (*authsession.Role, error)
	UpdateRoleWithGrants(string, authsession.RolePatch, time.Time) (*authsession.Role, error)
	DeleteRole(string) error
	DeleteRolesBatch([]string) (int, error)
	ValidatePermissionKeys([]string) error
	ValidateMenuItemIDs([]string) error
}

func rolesResource(repository RolesRepository, operations operationlog.Recorder) Resource {
	return Resource{
		ID:              "roles",
		Path:            "/api/roles",
		Listable:        true,
		SortFields:      []string{"key", "name", "updatedAt"},
		QSearch:         true,
		// T-02 (GOAL-013 D-003): system-flag select on the roles search form.
		ExtraQuery:      []string{"system"},
		Entity:          &rolesEntity{repository: repository},
		CreateFields:    []string{"key", "name"},
		PatchFields:     []string{"name"},
		JSONFields:      []string{"permissions", "menuItems"},
		PermissionRead:  "roles.read",
		PermissionWrite: "roles.write",
		NotFoundCode:    "ROLE_NOT_FOUND",
		// id is derived inside the entity as "role-<key>"; the factory-generated
		// id is ignored by rolesEntity.Create.
		NewID:   newRoleID,
		OnWrite: rolesOnWrite(operations),
	}
}

// RolesResource exposes the roles Resource descriptor to module providers
// (R4 C3.2).
func RolesResource(repository RolesRepository, operations operationlog.Recorder) Resource {
	return rolesResource(repository, operations)
}

// rolesEntity adapts the roles store to the generic resource boundary.
type rolesEntity struct {
	repository RolesRepository
}

// roleToMap maps a persisted role to the API row. system is emitted as a JSON
// boolean; timestamps use the frozen 3-digit-millisecond shape (I-011-001 §3.0).
func roleToMap(r authsession.Role) map[string]any {
	return map[string]any{
		"id":            r.ID,
		"key":           r.Key,
		"name":          r.Name,
		"system":        r.System,
		"permissions":   r.Permissions,
		"menuItems":     r.MenuItems,
		"assignedUsers": r.AssignedUsers,
		"editable":      !r.System,
		"deletable":     !r.System && r.AssignedUsers == 0,
		"createdAt":     r.CreatedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
		"updatedAt":     r.UpdatedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
	}
}

func (e *rolesEntity) List(f resourceFilter) ([]map[string]any, int, error) {
	// T-02 (GOAL-013 D-003): system query param ("true"/"false").
	var system *bool
	if raw, ok := f.Extra["system"]; ok && (raw == "true" || raw == "false") {
		v := raw == "true"
		system = &v
	}
	items, total, err := e.repository.ListRoles(authsession.RoleFilter{
		Q: f.Q, Sort: f.Sort, Order: f.Order, Page: f.Page, PageSize: f.PageSize,
		System: system,
	})
	if err != nil {
		return nil, 0, err
	}
	out := make([]map[string]any, 0, len(items))
	for _, it := range items {
		out = append(out, roleToMap(it))
	}
	return out, total, nil
}

func (e *rolesEntity) Get(id string) (map[string]any, error) {
	r, err := e.repository.GetRole(id)
	if err != nil {
		return nil, mapRoleStoreError(err)
	}
	return roleToMap(*r), nil
}

func (e *rolesEntity) Create(body map[string]any, _ string, now time.Time, actor account.User) (map[string]any, error) {
	permissions, err := stringArrayFromBody(body, "permissions", "INVALID_PERMISSION_REF")
	if err != nil {
		return nil, err
	}
	menuItems, err := stringArrayFromBody(body, "menuItems", "INVALID_MENU_ITEM_REF")
	if err != nil {
		return nil, err
	}
	if _, hasPermissions := body["permissions"]; hasPermissions {
		if !slices.Contains(actor.Roles, "admin") {
			return nil, grantForbidden("only an admin may manage role grants")
		}
		if err := e.repository.ValidatePermissionKeys(permissions); err != nil {
			return nil, mapRoleStoreError(err)
		}
		if err := authorizeGrantChange(actor, permissions); err != nil {
			return nil, err
		}
	}
	if _, hasMenus := body["menuItems"]; hasMenus {
		if !slices.Contains(actor.Roles, "admin") {
			return nil, grantForbidden("only an admin may manage role menu grants")
		}
		if err := e.repository.ValidateMenuItemIDs(menuItems); err != nil {
			return nil, mapRoleStoreError(err)
		}
	}
	r, err := e.repository.CreateRoleWithGrants(
		stringField(body, "key"), stringField(body, "name"), permissions, menuItems, now,
	)
	if err != nil {
		return nil, mapRoleStoreError(err)
	}
	return roleToMap(*r), nil
}

func (e *rolesEntity) Update(id string, body map[string]any, now time.Time, actor account.User) (map[string]any, error) {
	patch := authsession.RolePatch{}
	if _, ok := body["name"]; ok {
		name := stringField(body, "name")
		patch.Name = &name
	}
	if _, ok := body["permissions"]; ok {
		permissions, err := stringArrayFromBody(body, "permissions", "INVALID_PERMISSION_REF")
		if err != nil {
			return nil, err
		}
		if !slices.Contains(actor.Roles, "admin") {
			return nil, grantForbidden("only an admin may manage role grants")
		}
		if err := e.repository.ValidatePermissionKeys(permissions); err != nil {
			return nil, mapRoleStoreError(err)
		}
		if err := authorizeGrantChange(actor, permissions); err != nil {
			return nil, err
		}
		patch.Permissions = &permissions
	}
	if _, ok := body["menuItems"]; ok {
		menuItems, err := stringArrayFromBody(body, "menuItems", "INVALID_MENU_ITEM_REF")
		if err != nil {
			return nil, err
		}
		if !slices.Contains(actor.Roles, "admin") {
			return nil, grantForbidden("only an admin may manage role menu grants")
		}
		if err := e.repository.ValidateMenuItemIDs(menuItems); err != nil {
			return nil, mapRoleStoreError(err)
		}
		patch.MenuItems = &menuItems
	}
	r, err := e.repository.UpdateRoleWithGrants(id, patch, now)
	if err != nil {
		return nil, mapRoleStoreError(err)
	}
	return roleToMap(*r), nil
}

func (e *rolesEntity) Delete(id string, _ account.User) error {
	return mapRoleStoreError(e.repository.DeleteRole(id))
}

// DeleteBatch is the atomic whole-batch delete (ADR-0022 D5d · D-001 P0): the
// repository commits the whole selection in one transaction, so a failure
// (system role, in-use role, not-found) rolls every target back.
func (e *rolesEntity) DeleteBatch(ids []string, _ account.User) (int, error) {
	deleted, err := e.repository.DeleteRolesBatch(ids)
	if err != nil {
		return 0, mapRoleStoreError(err)
	}
	return deleted, nil
}

func stringArrayFromBody(body map[string]any, field, code string) ([]string, error) {
	v, ok := body[field]
	if !ok || v == nil {
		return []string{}, nil
	}
	arr, ok := v.([]any)
	if !ok {
		return nil, &DomainError{Status: 400, Code: code, Message: field + " must be an array of strings"}
	}
	values := make([]string, 0, len(arr))
	for _, item := range arr {
		value, ok := item.(string)
		value = strings.TrimSpace(value)
		if !ok || value == "" {
			return nil, &DomainError{Status: 400, Code: code, Message: field + " must contain non-empty strings"}
		}
		values = append(values, value)
	}
	return values, nil
}

func grantForbidden(message string) error {
	return &DomainError{Status: 403, Code: "ROLE_GRANT_FORBIDDEN", Message: message}
}

func authorizeGrantChange(actor account.User, permissions []string) error {
	if !slices.Contains(actor.Roles, "admin") {
		return grantForbidden("only an admin may manage role grants")
	}
	for _, permission := range permissions {
		if !slices.Contains(actor.Permissions, permission) {
			return grantForbidden("cannot grant a permission the actor does not hold")
		}
	}
	return nil
}

// mapRoleStoreError maps roles-store domain sentinels to the frozen
// resource-specific error codes (I-011-001 §6).
func mapRoleStoreError(err error) error {
	switch {
	case errors.Is(err, authsession.ErrNotFound):
		return errResourceNotFound
	case errors.Is(err, authsession.ErrRoleTaken):
		return &DomainError{Status: 409, Code: "ROLE_KEY_TAKEN", Message: "role key already exists"}
	case errors.Is(err, authsession.ErrRoleInUse):
		return &DomainError{Status: 409, Code: "ROLE_IN_USE", Message: "role is assigned to users"}
	case errors.Is(err, authsession.ErrRoleSystem):
		return &DomainError{Status: 409, Code: "ROLE_SYSTEM", Message: "system roles cannot be modified"}
	case errors.Is(err, authsession.ErrInvalidKey):
		return &DomainError{Status: 400, Code: "INVALID_ROLE_KEY", Message: "invalid role key format"}
	case errors.Is(err, authsession.ErrInvalidPermission):
		return &DomainError{Status: 400, Code: "INVALID_PERMISSION_REF", Message: "permissions contain an unknown key"}
	case errors.Is(err, authsession.ErrInvalidMenuItem):
		return &DomainError{Status: 400, Code: "INVALID_MENU_ITEM_REF", Message: "menuItems contain an unknown id"}
	default:
		return err
	}
}

// rolesOnWrite appends operation-log rows for roles write endpoints
// (I-011-001 §5). Best-effort.
func rolesOnWrite(recorder operationlog.Recorder) func(context.Context, account.User, writeKind, string, map[string]any, time.Time) {
	return func(_ context.Context, user account.User, kind writeKind, id string, row map[string]any, now time.Time) {
		event := operationlog.EventRoleDelete
		detail := ""
		switch kind {
		case writeCreate:
			event = operationlog.EventRoleCreate
			detail = "create"
		case writeUpdate:
			event = operationlog.EventRoleUpdate
			detail = "update"
		}
		var detailPtr *string
		if detail != "" {
			detailPtr = auditDetail(strings.TrimPrefix(event, "roles."), map[string]any{"key": stringField(row, "key")})
		}
		recordAudit(recorder, user, event, id, detailPtr, now, nil)
	}
}

// newRoleID is a placeholder id generator: rolesEntity.Create derives the id as
// "role-<key>" and ignores the factory-generated value, so this is never used
// for a real insert (I-011-001 §3 id = role-<key>).
func newRoleID() (string, error) {
	return "role", nil
}
