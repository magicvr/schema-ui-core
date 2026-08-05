// users resource instance of the generic resource factory (GOAL-011 S2 ·
// I-011-001 §2): mounted at /api/users with users.read / users.write,
// USER_NOT_FOUND, create/patch fields username/name/password + the JSON roles
// field, and operation-log events. The entity enforces the account-domain
// invariants (sensitive-field isolation, role validation without implicit
// creation, self/last-admin protection) via the store methods and maps store
// sentinels to the frozen resource-specific error codes (I-011-001 §6).
package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
)

// passwordHashCost is the bcrypt cost for users resource password hashing
// (matches the server bootstrap cost, I-011-001 §2.2).
const passwordHashCost = 10

const (
	minPasswordBytes = 8
	maxPasswordBytes = 72
)

func usersResource(st *store.Store) Resource {
	return Resource{
		ID:              "users",
		Path:            "/api/users",
		Listable:        true,
		SortFields:      []string{"username", "name", "updatedAt"},
		QSearch:         true,
		Entity:          &usersEntity{st: st},
		CreateFields:    []string{"username", "name"},
		PatchFields:     []string{"name"},
		RawStringFields: []string{"password"},
		JSONFields:      []string{"roles"},
		PermissionRead:  "users.read",
		PermissionWrite: "users.write",
		NotFoundCode:    "USER_NOT_FOUND",
		NewID:           newUserID,
		OnWrite:         usersOnWrite(st),
	}
}

// UsersResource exposes the users Resource descriptor to module providers
// (R4 C3.2).
func UsersResource(st *store.Store) Resource {
	return usersResource(st)
}

// usersEntity adapts the users store to the generic resource boundary. Rows
// never contain password_hash (sensitive-field isolation, I-011-001 §2.2).
type usersEntity struct {
	st *store.Store
}

// userToMap maps a persisted user to the API row. password_hash is intentionally
// absent; createdAt/updatedAt serialize with the frozen 3-digit-millisecond shape.
func userToMap(u store.User) map[string]any {
	return map[string]any{
		"id":        u.ID,
		"username":  u.Username,
		"name":      u.Name,
		"roles":     u.Roles,
		"createdAt": u.CreatedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
		"updatedAt": u.UpdatedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
	}
}

func (e *usersEntity) List(f resourceFilter) ([]map[string]any, int, error) {
	items, total, err := e.st.ListUsers(store.UserFilter{
		Q: f.Q, Sort: f.Sort, Order: f.Order, Page: f.Page, PageSize: f.PageSize,
	})
	if err != nil {
		return nil, 0, err
	}
	out := make([]map[string]any, 0, len(items))
	for _, it := range items {
		out = append(out, userToMap(it))
	}
	return out, total, nil
}

func (e *usersEntity) Get(id string) (map[string]any, error) {
	u, err := e.st.GetUser(id)
	if err != nil {
		return nil, err
	}
	return userToMap(*u), nil
}

func (e *usersEntity) Create(body map[string]any, id string, now time.Time, actor account.User) (map[string]any, error) {
	roles, err := rolesFromBody(body)
	if err != nil {
		return nil, err
	}
	if len(roles) > 0 {
		if err := e.authorizeRoleAssignment(actor, roles); err != nil {
			return nil, err
		}
	}
	password, err := managedPasswordFromBody(body)
	if err != nil {
		return nil, err
	}
	hash, err := auth.HashPassword(password, passwordHashCost)
	if err != nil {
		return nil, &DomainError{Status: 500, Code: "INTERNAL", Message: "could not hash password"}
	}
	u, err := e.st.CreateUserManagement(store.User{
		ID:           id,
		Username:     stringField(body, "username"),
		Name:         stringField(body, "name"),
		Roles:        roles,
		PasswordHash: hash,
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	if err != nil {
		return nil, mapUserStoreError(err)
	}
	return userToMap(*u), nil
}

func (e *usersEntity) Update(id string, body map[string]any, now time.Time, user account.User) (map[string]any, error) {
	patch := store.UserPatch{}
	if _, ok := body["name"]; ok {
		v := stringField(body, "name")
		patch.Name = &v
	}
	if _, ok := body["password"]; ok {
		password, err := managedPasswordFromBody(body)
		if err != nil {
			return nil, err
		}
		hash, err := auth.HashPassword(password, passwordHashCost)
		if err != nil {
			return nil, &DomainError{Status: 500, Code: "INTERNAL", Message: "could not hash password"}
		}
		patch.PasswordHash = &hash
	}
	if v, ok := body["roles"]; ok {
		roles, err := parseRolesValue(v)
		if err != nil {
			return nil, err
		}
		if err := e.authorizeRoleAssignment(user, roles); err != nil {
			return nil, err
		}
		patch.Roles = &roles
	}
	u, err := e.st.UpdateUser(id, patch, user.ID, now)
	if err != nil {
		return nil, mapUserStoreError(err)
	}
	return userToMap(*u), nil
}

func (e *usersEntity) Delete(id string, user account.User) error {
	return mapUserStoreError(e.st.DeleteUser(id, user.ID))
}

// rolesFromBody reads the optional roles JSON field (absent → nil, meaning "no
// role change" for patch / "no roles" for create).
func rolesFromBody(body map[string]any) ([]string, error) {
	v, ok := body["roles"]
	if !ok {
		return nil, nil
	}
	return parseRolesValue(v)
}

// parseRolesValue validates the roles JSON value is an array of non-empty
// strings (I-011-001 §2.3). null is treated as absent (nil). Unknown keys are
// rejected by the store (ErrInvalidRole → INVALID_ROLE_REF), not here.
func parseRolesValue(v any) ([]string, error) {
	if v == nil {
		return nil, nil
	}
	if raw, ok := v.(string); ok {
		if strings.TrimSpace(raw) == "" {
			return []string{}, nil
		}
		parts := strings.Split(raw, ",")
		roles := make([]string, 0, len(parts))
		for _, part := range parts {
			key := strings.TrimSpace(part)
			if key == "" {
				return nil, &DomainError{Status: 400, Code: "INVALID_ROLE_REF", Message: "roles must be comma-separated non-empty keys"}
			}
			roles = append(roles, key)
		}
		return roles, nil
	}
	arr, ok := v.([]any)
	if !ok {
		return nil, &DomainError{Status: 400, Code: "INVALID_ROLE_REF", Message: "roles must be an array of strings"}
	}
	roles := make([]string, 0, len(arr))
	for _, item := range arr {
		s, ok := item.(string)
		if !ok || strings.TrimSpace(s) == "" {
			return nil, &DomainError{Status: 400, Code: "INVALID_ROLE_REF", Message: "roles must be an array of non-empty strings"}
		}
		roles = append(roles, strings.TrimSpace(s))
	}
	return roles, nil
}

func managedPasswordFromBody(body map[string]any) (string, error) {
	password, ok := body["password"].(string)
	if !ok {
		return "", invalidManagedPassword()
	}
	length := len([]byte(password))
	if length < minPasswordBytes || length > maxPasswordBytes || strings.TrimSpace(password) == "" {
		return "", invalidManagedPassword()
	}
	return password, nil
}

func invalidManagedPassword() error {
	return &DomainError{
		Status:  400,
		Code:    "INVALID_PASSWORD",
		Message: "password must be a string with non-whitespace characters and 8 to 72 bytes",
	}
}

func (e *usersEntity) authorizeRoleAssignment(actor account.User, roles []string) error {
	forbidden := func(message string) error {
		return &DomainError{Status: 403, Code: "ROLE_ASSIGNMENT_FORBIDDEN", Message: message}
	}
	if !slices.Contains(actor.Permissions, "roles.assign") {
		return forbidden("permission required: roles.assign")
	}
	if slices.Contains(roles, "admin") && !slices.Contains(actor.Roles, "admin") {
		return forbidden("only an admin may assign the admin role")
	}
	targetPermissions, err := e.st.PermissionsForRoles(roles)
	if err != nil {
		return mapUserStoreError(err)
	}
	for _, permission := range targetPermissions {
		if !slices.Contains(actor.Permissions, permission) {
			return forbidden("cannot assign a role with permissions the actor does not hold")
		}
	}
	return nil
}

// mapUserStoreError maps users-store domain sentinels to the frozen
// resource-specific error codes (I-011-001 §6). Unknown errors pass through to
// the factory's INTERNAL fallback.
func mapUserStoreError(err error) error {
	switch {
	case errors.Is(err, store.ErrUsernameTaken):
		return &DomainError{Status: 409, Code: "USERNAME_TAKEN", Message: "username already exists"}
	case errors.Is(err, store.ErrLastAdmin):
		return &DomainError{Status: 409, Code: "LAST_ADMIN", Message: "cannot remove the last admin user"}
	case errors.Is(err, store.ErrSelfOperation):
		return &DomainError{Status: 409, Code: "SELF_OPERATION", Message: "self operation is not allowed"}
	case errors.Is(err, store.ErrInvalidRole):
		return &DomainError{Status: 400, Code: "INVALID_ROLE_REF", Message: "roles contain an unknown role key"}
	default:
		return err
	}
}

// usersOnWrite appends operation-log rows for users write endpoints
// (I-011-001 §5). Best-effort: a logging failure never fails the write.
func usersOnWrite(st *store.Store) func(context.Context, account.User, writeKind, string, map[string]any, time.Time) {
	return func(_ context.Context, user account.User, kind writeKind, id string, row map[string]any, now time.Time) {
		event := store.EventUserDelete
		detail := ""
		switch kind {
		case writeCreate:
			event = store.EventUserCreate
			detail = `{"username":` + jsonQuote(stringField(row, "username")) + `}`
		case writeUpdate:
			event = store.EventUserUpdate
			detail = `{"username":` + jsonQuote(stringField(row, "username")) + `}`
		}
		op := store.Operation{
			ID:        newOperationID(),
			Event:     event,
			ActorID:   user.ID,
			ActorName: user.Name,
			CreatedAt: now,
		}
		if id != "" {
			op.RecordID = &id
		}
		if detail != "" {
			op.Detail = &detail
		}
		if err := st.RecordOperation(op); err != nil {
			slog.Error("operation log write failed", "event", event, "err", err)
		}
	}
}

// newUserID returns "usr-" + 16 lowercase hex chars (8 bytes of crypto/rand),
// the users create id format (I-011-001 §2).
func newUserID() (string, error) {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return "usr-" + hex.EncodeToString(b[:]), nil
}
