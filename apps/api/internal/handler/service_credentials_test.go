package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	authsession "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
)

func TestServiceCredentialManagementAndAuthentication(t *testing.T) {
	env, _ := newWalletSelfEnv(t)
	mountMFASurface(t, env, newFakeMFAService(), &fakeSessionRevoker{})
	admin := loginAs(t, env, testSeedUsername, testSeedPassword)
	expiresAt := time.Now().UTC().Add(24 * time.Hour).Format(time.RFC3339)
	createBody := `{"name":"Build Agent","scopes":["users.read"],"expiresAt":"` + expiresAt + `"}`

	create := httptest.NewRecorder()
	env.mux.ServeHTTP(create, bearer(t, admin, http.MethodPost, "/api/service-credentials", createBody))
	if create.Code != http.StatusCreated {
		t.Fatalf("create = %d %s", create.Code, create.Body.String())
	}
	var created map[string]any
	if err := json.Unmarshal(create.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	credentialID, _ := created["id"].(string)
	secret, _ := created["secret"].(string)
	if credentialID == "" || !strings.HasPrefix(secret, "sui_sc_") || created["token"] != nil || created["tokenHash"] != nil {
		t.Fatalf("create response = %+v", created)
	}

	for _, path := range []string{"/api/service-credentials", "/api/service-credentials/" + credentialID} {
		response := httptest.NewRecorder()
		env.mux.ServeHTTP(response, bearer(t, admin, http.MethodGet, path, ""))
		if response.Code != http.StatusOK || strings.Contains(response.Body.String(), secret) || strings.Contains(response.Body.String(), "tokenHash") || strings.Contains(response.Body.String(), `"secret"`) {
			t.Fatalf("metadata %s = %d %s", path, response.Code, response.Body.String())
		}
	}

	duplicate := httptest.NewRecorder()
	env.mux.ServeHTTP(duplicate, bearer(t, admin, http.MethodPost, "/api/service-credentials", strings.Replace(createBody, "Build Agent", "build agent", 1)))
	expectError(t, duplicate, http.StatusBadRequest, "INVALID_CREATE_FIELD")

	serviceRequest := httptest.NewRecorder()
	env.mux.ServeHTTP(serviceRequest, bearer(t, secret, http.MethodGet, "/api/users", ""))
	if serviceRequest.Code != http.StatusOK {
		t.Fatalf("service permission request = %d %s", serviceRequest.Code, serviceRequest.Body.String())
	}
	for _, userOnlyRoute := range []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/accounts/me"},
		{http.MethodGet, "/api/account/profile"},
		{http.MethodPost, "/api/account/avatar"},
		{http.MethodGet, "/api/mfa/status"},
		{http.MethodGet, "/api/notifications"},
		{http.MethodGet, "/api/wallet/me"},
	} {
		userOnly := httptest.NewRecorder()
		env.mux.ServeHTTP(userOnly, bearer(t, secret, userOnlyRoute.method, userOnlyRoute.path, ""))
		expectError(t, userOnly, http.StatusUnauthorized, "UNAUTHENTICATED")
	}
	management := httptest.NewRecorder()
	env.mux.ServeHTTP(management, bearer(t, secret, http.MethodGet, "/api/service-credentials", ""))
	expectError(t, management, http.StatusForbidden, "FORBIDDEN")

	for index := 0; index < 2; index++ {
		revoke := httptest.NewRecorder()
		env.mux.ServeHTTP(revoke, bearer(t, admin, http.MethodPost, "/api/service-credentials/"+credentialID+"/revoke", ""))
		if revoke.Code != http.StatusNoContent {
			t.Fatalf("revoke %d = %d %s", index, revoke.Code, revoke.Body.String())
		}
	}
	revoked := httptest.NewRecorder()
	env.mux.ServeHTTP(revoked, bearer(t, secret, http.MethodGet, "/api/users", ""))
	expectError(t, revoked, http.StatusUnauthorized, "UNAUTHENTICATED")

	operations, _, err := env.operations.ListOperationsFiltered(operationlog.OperationFilter{
		Q: "service-credentials", Sort: "createdAt", Order: "asc", Page: 1, PageSize: 100,
	})
	if err != nil {
		t.Fatal(err)
	}
	events := map[string]int{}
	for _, operation := range operations {
		events[operation.Event]++
		if operation.Detail != nil && strings.Contains(*operation.Detail, secret) {
			t.Fatal("raw service credential leaked into audit detail")
		}
		if operation.Event == operationlog.EventServiceCredentialUse && (operation.Detail == nil || !strings.Contains(*operation.Detail, `"scopeCount":1`)) {
			t.Fatalf("use audit missing scopeCount: %+v", operation)
		}
	}
	if events[operationlog.EventServiceCredentialCreate] != 1 || events[operationlog.EventServiceCredentialUse] == 0 || events[operationlog.EventServiceCredentialRevoke] != 1 {
		t.Fatalf("service credential audit events = %v", events)
	}
}

func TestServiceCredentialCreateRejectsUnknownReservedAndExcessScopes(t *testing.T) {
	env := newAuthTestEnv(t)
	admin := loginAs(t, env, testSeedUsername, testSeedPassword)
	expiresAt := time.Now().UTC().Add(time.Hour).Format(time.RFC3339)
	tests := []struct {
		name, scope, code string
		status            int
	}{
		{name: "unknown", scope: "unknown.permission", status: http.StatusBadRequest, code: "INVALID_PERMISSION_REF"},
		{name: "reserved", scope: "service-credentials.read", status: http.StatusForbidden, code: "FORBIDDEN"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := `{"name":"` + test.name + `","scopes":["` + test.scope + `"],"expiresAt":"` + expiresAt + `"}`
			response := httptest.NewRecorder()
			env.mux.ServeHTTP(response, bearer(t, admin, http.MethodPost, "/api/service-credentials", body))
			expectError(t, response, test.status, test.code)
		})
	}
	overLimitScopes := make([]string, 65)
	for index := range overLimitScopes {
		overLimitScopes[index] = fmt.Sprintf("permission.%02d", index)
	}
	overLimitBody, err := json.Marshal(map[string]any{"name": "over-limit", "scopes": overLimitScopes, "expiresAt": expiresAt})
	if err != nil {
		t.Fatal(err)
	}
	overLimit := httptest.NewRecorder()
	env.mux.ServeHTTP(overLimit, bearer(t, admin, http.MethodPost, "/api/service-credentials", string(overLimitBody)))
	expectError(t, overLimit, http.StatusBadRequest, "INVALID_CREATE_FIELD")

	managerPassword := "manager-password"
	managerHash, err := auth.HashPassword(managerPassword, testBcryptCost)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC()
	if _, err := env.authRepository.CreateRoleWithGrants("credential-manager", "Credential Manager", []string{serviceCredentialReadPermission, serviceCredentialWritePermission}, nil, now); err != nil {
		t.Fatal(err)
	}
	if _, err := env.authRepository.CreateUserManagement(authsession.User{
		ID: "user-credential-manager", Username: "credential-manager", Name: "Credential Manager",
		Roles: []string{"credential-manager"}, PasswordHash: managerHash, CreatedAt: now, UpdatedAt: now,
	}); err != nil {
		t.Fatal(err)
	}
	manager := loginAs(t, env, "credential-manager", managerPassword)
	excess := httptest.NewRecorder()
	body := `{"name":"excess","scopes":["users.read"],"expiresAt":"` + expiresAt + `"}`
	env.mux.ServeHTTP(excess, bearer(t, manager, http.MethodPost, "/api/service-credentials", body))
	expectError(t, excess, http.StatusForbidden, "FORBIDDEN")
}

func TestServiceCredentialRequiredAuditFailureRollsBack(t *testing.T) {
	env := newAuthTestEnv(t)
	admin := loginAs(t, env, testSeedUsername, testSeedPassword)
	expiresAt := time.Now().UTC().Add(time.Hour).Format(time.RFC3339)
	body := `{"name":"Rollback Agent","scopes":["users.read"],"expiresAt":"` + expiresAt + `"}`
	forced := errors.New("forced audit failure")
	env.operations.SetOperationLogError(forced)
	failedCreate := httptest.NewRecorder()
	env.mux.ServeHTTP(failedCreate, bearer(t, admin, http.MethodPost, "/api/service-credentials", body))
	expectError(t, failedCreate, http.StatusInternalServerError, "INTERNAL")
	items, total, err := env.authRepository.ListServiceCredentials(1, 20)
	if err != nil || total != 0 || len(items) != 0 {
		t.Fatalf("credentials after failed create audit = %+v total=%d err=%v", items, total, err)
	}

	env.operations.SetOperationLogError(nil)
	createdResponse := httptest.NewRecorder()
	env.mux.ServeHTTP(createdResponse, bearer(t, admin, http.MethodPost, "/api/service-credentials", body))
	if createdResponse.Code != http.StatusCreated {
		t.Fatalf("create after clearing failure = %d %s", createdResponse.Code, createdResponse.Body.String())
	}
	var created map[string]any
	if err := json.Unmarshal(createdResponse.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	id := created["id"].(string)
	env.operations.SetOperationLogError(forced)
	failedRevoke := httptest.NewRecorder()
	env.mux.ServeHTTP(failedRevoke, bearer(t, admin, http.MethodPost, "/api/service-credentials/"+id+"/revoke", ""))
	expectError(t, failedRevoke, http.StatusInternalServerError, "INTERNAL")
	credential, err := env.authRepository.ServiceCredentialByID(id)
	if err != nil || credential.RevokedAt != nil {
		t.Fatalf("credential after failed revoke audit = %+v err=%v", credential, err)
	}
}
