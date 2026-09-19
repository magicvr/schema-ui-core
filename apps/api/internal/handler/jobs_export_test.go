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
	for _, route := range JobsRoutes(env.a, repository, submitter, "admin.jobs") {
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
