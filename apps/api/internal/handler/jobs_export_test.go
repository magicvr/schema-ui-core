package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/jobs"
)

// The HTTP-level batch-export tests use a submitter that enqueues a REAL job
// row through the real repository but registers no handler body. The export
// job itself (progress reporting, CSV rendering, result document) is covered
// by modules/jobs' own test — importing that module here would create an import
// cycle, so this file owns exactly the HTTP contract: the 202 envelope, the
// double permission gate, body validation, and the "rejected request leaves no
// job row" invariant.
type recordingBatchSubmitter struct {
	repository *jobs.Repository
	submitted  []recordedSubmission
	supported  []string
}

type recordedSubmission struct {
	resource  string
	ids       []string
	actorID   string
	corrID    string
}

func (s *recordingBatchSubmitter) Supported(resource string) bool {
	for _, supported := range s.supported {
		if supported == resource {
			return true
		}
	}
	return false
}

func (s *recordingBatchSubmitter) SubmitBatchExport(ctx context.Context, resource string, ids []string, actorID, correlationID string) (*jobs.Job, error) {
	// Mirror the real service's validation contract so the HTTP layer is tested
	// against the same accept/reject boundary.
	if !s.Supported(resource) || len(ids) == 0 {
		return nil, jobs.ErrInvalid
	}
	now := time.Now().UTC()
	id, err := jobs.NewID(now)
	if err != nil {
		return nil, err
	}
	s.submitted = append(s.submitted, recordedSubmission{resource: resource, ids: ids, actorID: actorID, corrID: correlationID})
	return s.repository.Create(ctx, jobs.CreateInput{
		ID: id, Kind: "jobs.batch-export", Payload: []byte(`{}`),
		ActorID: actorID, CorrelationID: correlationID, Now: now,
	})
}

func mountJobsExportRoutes(t *testing.T, env *authTestEnv) (*jobs.Repository, *recordingBatchSubmitter) {
	t.Helper()
	repository := jobs.NewRepository(env.st)
	submitter := &recordingBatchSubmitter{repository: repository, supported: []string{"users", "roles"}}
	for _, route := range JobsRoutes(env.a, repository, submitter, nil, "admin.jobs") {
		env.mux.Handle(route.Method+" "+route.Pattern, route.Handler)
	}
	return repository, submitter
}

// C1: a valid selection answers 202 with the job projection envelope.
func TestJobsBatchExportSubmitReturns202(t *testing.T) {
	env := newAuthTestEnv(t)
	_, submitter := mountJobsExportRoutes(t, env)
	token := adminToken(t, env)

	body := `{"resource":"users","ids":["usr-1","usr-2","usr-1"]}`
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, token, http.MethodPost, "/api/jobs/batch-export", body))
	if rr.Code != http.StatusAccepted {
		t.Fatalf("submit = %d body=%s, want 202", rr.Code, rr.Body.String())
	}
	var submitted map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &submitted); err != nil {
		t.Fatal(err)
	}
	if id, _ := submitted["id"].(string); id == "" {
		t.Fatalf("202 body has no job id: %v", submitted)
	}
	if submitted["kind"] != "jobs.batch-export" {
		t.Fatalf("kind = %v, want jobs.batch-export", submitted["kind"])
	}
	if submitted["status"] != "queued" {
		t.Fatalf("status = %v, want queued", submitted["status"])
	}
	if _, ok := submitted["resultUrl"]; ok {
		t.Fatalf("a queued job must not advertise a resultUrl: %v", submitted)
	}
	// The selection is normalized like the synchronous batch action: the
	// duplicate key is dropped, order preserved.
	if len(submitter.submitted) != 1 {
		t.Fatalf("submissions = %d, want 1", len(submitter.submitted))
	}
	got := submitter.submitted[0].ids
	if len(got) != 2 || got[0] != "usr-1" || got[1] != "usr-2" {
		t.Fatalf("ids = %v, want [usr-1 usr-2] (deduped, order preserved)", got)
	}
}

// C1 (security): the async path is gated by BOTH jobs.write and data.export, so
// adding the Job runtime cannot widen data egress.
func TestJobsBatchExportRequiresBothGates(t *testing.T) {
	env := newAuthTestEnv(t)
	_, submitter := mountJobsExportRoutes(t, env)
	body := `{"resource":"users","ids":["any-id"]}`

	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/api/jobs/batch-export", strings.NewReader(body)))
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("anonymous submit = %d, want 401", rr.Code)
	}

	env.addUser(t, "editor1", "editor-password", []string{"editor"})
	editorToken := env.login(t, "editor1", "editor-password")
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, editorToken, http.MethodPost, "/api/jobs/batch-export", body))
	if rr.Code != http.StatusForbidden {
		t.Fatalf("editor submit = %d body=%s, want 403", rr.Code, rr.Body.String())
	}

	if len(submitter.submitted) != 0 {
		t.Fatalf("a denied request reached the submitter: %v", submitter.submitted)
	}
}

// C1: the submit body is validated before any job row is created.
func TestJobsBatchExportValidation(t *testing.T) {
	env := newAuthTestEnv(t)
	repository, submitter := mountJobsExportRoutes(t, env)
	token := adminToken(t, env)

	for _, tc := range []struct {
		name string
		body string
		code string
	}{
		{"empty selection", `{"resource":"users","ids":[]}`, "EMPTY_SELECTION"},
		{"blank key", `{"resource":"users","ids":[""]}`, "INVALID_SELECTION_KEY"},
		{"object key", `{"resource":"users","ids":[{"a":1}]}`, "INVALID_SELECTION_KEY"},
		{"missing resource", `{"resource":"","ids":["a"]}`, "INVALID_BODY"},
		{"malformed body", `not json`, "INVALID_BODY"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			env.mux.ServeHTTP(rr, bearer(t, token, http.MethodPost, "/api/jobs/batch-export", tc.body))
			if rr.Code != http.StatusBadRequest || !bodyHasCode(rr, tc.code) {
				t.Fatalf("submit = %d %s, want 400 %s", rr.Code, rr.Body.String(), tc.code)
			}
		})
	}

	// An unsupported resource is a 404, not a queued job.
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, token, http.MethodPost, "/api/jobs/batch-export", `{"resource":"settings","ids":["a"]}`))
	if rr.Code != http.StatusNotFound || !bodyHasCode(rr, "RESOURCE_NOT_FOUND") {
		t.Fatalf("unsupported resource = %d %s, want 404 RESOURCE_NOT_FOUND", rr.Code, rr.Body.String())
	}

	// No rejected request may leave a job row behind.
	if _, total, err := repository.ListJobs(context.Background(), jobs.ListFilter{Kind: "jobs.batch-export"}); err != nil || total != 0 {
		t.Fatalf("rejected submits left %d job rows (err=%v), want 0", total, err)
	}
	if len(submitter.submitted) != 0 {
		t.Fatalf("a rejected request reached the submitter: %v", submitter.submitted)
	}
}

// C3: the read routes from R2 still behave, and the new submit route did not
// disturb them (the same route table serves both phases).
// C1 (security): the submit route requires BOTH jobs.write and data.export, so
// adding the Job runtime cannot widen data egress.
//
// Under the built-in policy matrix the second gate is defence-in-depth rather
// than an active constraint: jobs.write is PolicyAdmin ({admin}) and data.export
// is PolicyAdminEditor ({admin, editor}), so every principal holding jobs.write
// already holds data.export. A discriminating test is therefore not
// constructible from seeded roles (A-001 F-001 / A-002 F-001). This test pins
// the nesting that makes that true: if the policy sets ever stop being nested,
// it fails loudly and a real counter-example test becomes both possible and
// necessary.
func TestJobsBatchExportSecondGateIsDefenceInDepth(t *testing.T) {
	env := newAuthTestEnv(t)
	mountJobsExportRoutes(t, env)

	// Read the effective permission sets straight from the repository: addUser
	// creates "user-<username>", and the seeded admin is "user-admin".
	adminPermissions := permissionsForUser(t, env, "user-admin")
	env.addUser(t, "editor-probe", "editor-password", []string{"editor"})
	editorPermissions := permissionsForUser(t, env, "user-editor-probe")

	if !containsPermission(adminPermissions, "jobs.write") || !containsPermission(adminPermissions, "data.export") {
		t.Fatalf("the admin principal lacks jobs.write or data.export (%v): the batch-export route is unreachable", adminPermissions)
	}

	// The nesting that makes the second gate unreachable-by-rejection today.
	if containsPermission(editorPermissions, "jobs.write") {
		t.Fatalf("the editor role now holds jobs.write: a jobs.write-without-data.export principal " +
			"may be constructible, so the discriminating counter-example test is now required")
	}
	// And the property the gate protects: whoever may submit may also export.
	for _, permissions := range [][]string{adminPermissions, editorPermissions} {
		if containsPermission(permissions, "jobs.write") && !containsPermission(permissions, "data.export") {
			t.Fatalf("a principal holds jobs.write without data.export (%v); the second gate is now "+
				"load-bearing and must have a discriminating test", permissions)
		}
	}
}

func containsPermission(permissions []string, want string) bool {
	for _, permission := range permissions {
		if permission == want {
			return true
		}
	}
	return false
}

// permissionsForUser reads one principal's effective permission set.
func permissionsForUser(t *testing.T, env *authTestEnv, userID string) []string {
	t.Helper()
	permissions, err := env.authRepository.PermissionsForUser(userID)
	if err != nil {
		t.Fatalf("permissions for %s: %v", userID, err)
	}
	if len(permissions) == 0 {
		t.Fatalf("principal %s has no permissions", userID)
	}
	return permissions
}

func TestJobsRoutesIncludeReadAndSubmit(t *testing.T) {
	env := newAuthTestEnv(t)
	mountJobsExportRoutes(t, env)
	token := adminToken(t, env)

	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, token, http.MethodGet, "/api/jobs", ""))
	if rr.Code != http.StatusOK {
		t.Fatalf("GET /api/jobs = %d, want 200", rr.Code)
	}
	// A GET on the submit path is not routed: the method matters.
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, token, http.MethodGet, "/api/jobs/batch-export", ""))
	if rr.Code == http.StatusAccepted {
		t.Fatalf("GET /api/jobs/batch-export must not submit a job")
	}
}
