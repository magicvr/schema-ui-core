package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	authsession "github.com/magicvr/schema-ui-core/apps/api/modules/authsession"
)

type stubPolicySettingsRepo struct {
	updateErr error
	updated   authsession.PasswordPolicy
}

func (s *stubPolicySettingsRepo) GetPasswordPolicy() (authsession.PasswordPolicy, error) {
	return authsession.PasswordPolicy{MinLength: 8}, nil
}

func (s *stubPolicySettingsRepo) UpdatePasswordPolicy(p authsession.PasswordPolicy) error {
	if s.updateErr != nil {
		return s.updateErr
	}
	s.updated = p
	return nil
}

func (s *stubPolicySettingsRepo) PermissionsForUser(userID string) ([]string, error) {
	return []string{"settings.read", "settings.write"}, nil
}

// A-001 F-001: the PATCH failure path differentiates sentinels — a missing
// singleton row maps to the frozen SETTINGS_NOT_FOUND code with 404 instead
// of the blanket INTERNAL 500.
func TestPatchPasswordPolicyNotSeededMapsSettingsNotFound(t *testing.T) {
	h := &policySettingsHandler{repo: &stubPolicySettingsRepo{updateErr: authsession.ErrPasswordPolicyNotSeeded}}
	req := httptest.NewRequest(http.MethodPatch, "/api/settings/password-policy", strings.NewReader(`{"minLength":12}`))
	rec := httptest.NewRecorder()
	h.patch().ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
	var body map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body["error"] != "SETTINGS_NOT_FOUND" {
		t.Fatalf("error code = %v, want SETTINGS_NOT_FOUND", body["error"])
	}
}

func TestPatchPasswordPolicySuccessStill200(t *testing.T) {
	repo := &stubPolicySettingsRepo{}
	h := &policySettingsHandler{repo: repo}
	req := httptest.NewRequest(http.MethodPatch, "/api/settings/password-policy", strings.NewReader(`{"minLength":12,"historyDepth":3}`))
	rec := httptest.NewRecorder()
	h.patch().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if repo.updated.MinLength != 12 || repo.updated.HistoryDepth != 3 || repo.updated.MinCategories != 0 {
		t.Fatalf("persisted = %+v, want minLength 12 / depth 3 / categories 0", repo.updated)
	}
}
