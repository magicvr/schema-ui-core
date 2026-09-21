package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/jobs"
)

// GOAL-005 R4 (D-001 §2): the HTTP contract of the result-center actions.
//
// These tests drive the REAL management-scope implementation — the repository
// plus the runner's CancelAny/RetryAny wiring, which is what composition binds
// to handler.JobActions — so the 409/404 mapping is asserted against the actual
// state machine rather than a hand-written stub. The runner is constructed but
// never started: cancel/retry are pure command paths (repository + wake signal),
// and keeping the scan loop out means no background goroutine can race a row
// between the request and the assertion.
type jobActionEnv struct {
	env        *authTestEnv
	repository *jobs.Repository
	runner     *jobs.Runner
}

func mountJobActionRoutes(t *testing.T, env *authTestEnv) *jobActionEnv {
	t.Helper()
	repository := jobs.NewRepository(env.st)
	runner, err := jobs.NewRunner(repository, jobs.RunnerOptions{
		LeaseDuration: time.Minute, HeartbeatInterval: time.Second,
		ScanInterval: time.Second, ResultTTL: time.Hour, BatchSize: 10, Now: time.Now,
	})
	if err != nil {
		t.Fatalf("new runner: %v", err)
	}
	for _, route := range JobsRoutes(env.a, repository, nil, runner, "admin.jobs") {
		env.mux.Handle(route.Method+" "+route.Pattern, route.Handler)
	}
	return &jobActionEnv{env: env, repository: repository, runner: runner}
}

// seedJobForOperator inserts a job owned by somebody other than the admin
// principal, so every assertion below is also a management-scope assertion.
func (e *jobActionEnv) seedJobForOperator(t *testing.T, id string, at time.Time) *jobs.Job {
	t.Helper()
	job, err := e.repository.Create(context.Background(), jobs.CreateInput{
		ID: id, Kind: "jobs.batch-export", Payload: []byte(`{}`),
		ActorID: "some-other-actor", CorrelationID: "corr-" + id, MaxAttempts: 3, Now: at,
	})
	if err != nil {
		t.Fatalf("seed %s: %v", id, err)
	}
	return job
}

func (e *jobActionEnv) post(t *testing.T, token, path string) *httptest.ResponseRecorder {
	t.Helper()
	rr := httptest.NewRecorder()
	e.env.mux.ServeHTTP(rr, bearer(t, token, http.MethodPost, path, ""))
	return rr
}

func decodeJobBody(t *testing.T, rr *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var body map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode %q: %v", rr.Body.String(), err)
	}
	return body
}

// C1: cancelling another actor's queued job returns the updated projection.
func TestJobCancelAnySucceeds(t *testing.T) {
	actionEnv := mountJobActionRoutes(t, newAuthTestEnv(t))
	token := adminToken(t, actionEnv.env)
	job := actionEnv.seedJobForOperator(t, "job-cancel-http", time.Now().UTC())

	rr := actionEnv.post(t, token, "/api/jobs/"+job.ID+"/cancel")
	if rr.Code != http.StatusOK {
		t.Fatalf("cancel = %d %s, want 200", rr.Code, rr.Body.String())
	}
	body := decodeJobBody(t, rr)
	if body["status"] != string(jobs.StatusCancelled) || body["id"] != job.ID {
		t.Fatalf("cancel body = %v, want the cancelled job", body)
	}
	if body["actorId"] != "some-other-actor" {
		t.Fatalf("the projection must keep the original owner: %v", body)
	}
	// The projection is the same envelope the read routes use, so the result
	// center needs no second decoder.
	if body["progress"] == nil || body["maxAttempts"] == nil {
		t.Fatalf("cancel body is missing the shared projection fields: %v", body)
	}
}

// C2: retrying another actor's failed job returns it queued, from a clean state.
func TestJobRetryAnySucceeds(t *testing.T) {
	actionEnv := mountJobActionRoutes(t, newAuthTestEnv(t))
	token := adminToken(t, actionEnv.env)
	job := actionEnv.seedJobForOperator(t, "job-retry-http", time.Now().UTC())

	// Fail it through the real lifecycle so the retry precondition is genuine.
	if _, lease, err := actionEnv.repository.Claim(context.Background(), job.ID, "worker-1", time.Now().UTC(), time.Minute); err != nil {
		t.Fatalf("claim: %v", err)
	} else if err := actionEnv.repository.Fail(context.Background(), lease, "JOB_HANDLER_FAILED", "boom", time.Now().UTC()); err != nil {
		t.Fatalf("fail: %v", err)
	}

	rr := actionEnv.post(t, token, "/api/jobs/"+job.ID+"/retry")
	if rr.Code != http.StatusOK {
		t.Fatalf("retry = %d %s, want 200", rr.Code, rr.Body.String())
	}
	body := decodeJobBody(t, rr)
	if body["status"] != string(jobs.StatusQueued) {
		t.Fatalf("retry body status = %v, want queued", body["status"])
	}
	if body["progress"] != float64(0) {
		t.Fatalf("retry body progress = %v, want 0", body["progress"])
	}
	if _, hasError := body["error"]; hasError {
		t.Fatalf("retry must clear the previous failure: %v", body)
	}
	if body["attempt"] != float64(1) {
		t.Fatalf("retry must preserve the consumed attempt: %v", body)
	}
}

// C1/C2: a job that does not exist is a 404 with the frozen code — never a
// silent success and never a 500.
func TestJobActionsMissingJobIs404(t *testing.T) {
	actionEnv := mountJobActionRoutes(t, newAuthTestEnv(t))
	token := adminToken(t, actionEnv.env)

	for _, path := range []string{"/api/jobs/job-absent/cancel", "/api/jobs/job-absent/retry"} {
		rr := actionEnv.post(t, token, path)
		if rr.Code != http.StatusNotFound || !bodyHasCode(rr, "JOB_NOT_FOUND") {
			t.Fatalf("POST %s = %d %s, want 404 JOB_NOT_FOUND", path, rr.Code, rr.Body.String())
		}
	}
}

// C1/C2: a rejected transition is a 409 with the frozen code, and it leaves the
// row untouched.
func TestJobActionsRejectedTransitionIs409(t *testing.T) {
	actionEnv := mountJobActionRoutes(t, newAuthTestEnv(t))
	token := adminToken(t, actionEnv.env)

	// A queued job cannot be retried; a cancelled one cannot be cancelled.
	queued := actionEnv.seedJobForOperator(t, "job-queued-409", time.Now().UTC())
	cancelled := actionEnv.seedJobForOperator(t, "job-cancelled-409", time.Now().UTC())
	if rr := actionEnv.post(t, token, "/api/jobs/"+cancelled.ID+"/cancel"); rr.Code != http.StatusOK {
		t.Fatalf("precondition cancel = %d %s", rr.Code, rr.Body.String())
	}

	cases := []struct {
		name string
		id   string
		verb string
		code string
	}{
		{"retry a queued job", queued.ID, "retry", "JOB_NOT_RETRYABLE"},
		{"cancel a cancelled job", cancelled.ID, "cancel", "JOB_NOT_CANCELLABLE"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			before, err := actionEnv.repository.GetJob(context.Background(), tc.id)
			if err != nil {
				t.Fatalf("get: %v", err)
			}
			rr := actionEnv.post(t, token, "/api/jobs/"+tc.id+"/"+tc.verb)
			if rr.Code != http.StatusConflict || !bodyHasCode(rr, tc.code) {
				t.Fatalf("%s = %d %s, want 409 %s", tc.name, rr.Code, rr.Body.String(), tc.code)
			}
			after, err := actionEnv.repository.GetJob(context.Background(), tc.id)
			if err != nil {
				t.Fatalf("get: %v", err)
			}
			if after.Status != before.Status || after.Attempt != before.Attempt || !after.UpdatedAt.Equal(before.UpdatedAt) {
				t.Fatalf("a rejected %s mutated the row: %+v -> %+v", tc.verb, before, after)
			}
		})
	}
}

// C1 (security): both actions are gated by jobs.write, and a denied request
// never reaches the state machine.
func TestJobActionsRequireJobsWrite(t *testing.T) {
	actionEnv := mountJobActionRoutes(t, newAuthTestEnv(t))
	job := actionEnv.seedJobForOperator(t, "job-gate-http", time.Now().UTC())

	// Anonymous: 401 on both verbs.
	for _, verb := range []string{"cancel", "retry"} {
		rr := httptest.NewRecorder()
		actionEnv.env.mux.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/api/jobs/"+job.ID+"/"+verb, nil))
		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("anonymous %s = %d, want 401", verb, rr.Code)
		}
	}

	// The editor role holds no jobs.write: 403, and the job keeps its state.
	actionEnv.env.addUser(t, "editor-actions", "editor-password", []string{"editor"})
	editorToken := actionEnv.env.login(t, "editor-actions", "editor-password")
	if containsPermission(permissionsForUser(t, actionEnv.env, "user-editor-actions"), "jobs.write") {
		t.Fatalf("the editor role now holds jobs.write; this test must be rebuilt around a discriminating principal")
	}
	for _, verb := range []string{"cancel", "retry"} {
		rr := actionEnv.post(t, editorToken, "/api/jobs/"+job.ID+"/"+verb)
		if rr.Code != http.StatusForbidden {
			t.Fatalf("editor %s = %d %s, want 403", verb, rr.Code, rr.Body.String())
		}
	}
	after, err := actionEnv.repository.GetJob(context.Background(), job.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if after.Status != jobs.StatusQueued || after.CancelRequested {
		t.Fatalf("a denied action mutated the job: %+v", after)
	}
}

// C1 (security) · the discriminating counter-example for the write gate.
//
// `TestJobActionsRequireJobsWrite` proves that a principal without jobs.write is
// refused, but NOT that the gate is jobs.write rather than jobs.read: under the
// built-in matrix both keys are PolicyAdmin ({admin}) and the editor role holds
// neither, so swapping the gate still leaves that test green.
//
// A discriminating principal IS constructible — the seeded roles are not the
// only source of grants — so this test builds one and closes the gap for real
// (independent A-002 F-001). The role below is granted jobs.read and nothing
// else, deliberately NOT jobs.write, and the same principal is shown to be
// allowed to READ and refused every write verb.
func TestJobWriteGateRequiresJobsWriteNotJobsRead(t *testing.T) {
	actionEnv := mountJobActionRoutes(t, newAuthTestEnv(t))
	env := actionEnv.env
	job := actionEnv.seedJobForOperator(t, "job-read-only-role", time.Now().UTC())

	// A custom role with exactly one grant: jobs.read. (Same shape as the
	// wallet voucher suite's permission-isolation fixture, but through the real
	// repository so the grant resolves permission id → key correctly.)
	if _, err := env.authRepository.CreateRoleWithGrants(
		"jobs-reader", "Jobs reader", []string{"jobs.read"}, nil, time.Now().UTC(),
	); err != nil {
		t.Fatalf("seed read-only role: %v", err)
	}
	env.addUser(t, "jobs-reader", "jobs-reader-password", []string{"jobs-reader"})
	token := env.login(t, "jobs-reader", "jobs-reader-password")

	permissions := permissionsForUser(t, env, "user-jobs-reader")
	if !containsPermission(permissions, "jobs.read") || containsPermission(permissions, "jobs.write") {
		t.Fatalf("fixture role permissions = %v, want jobs.read WITHOUT jobs.write", permissions)
	}

	// Positive control: the principal really can read. Without this the 403s
	// below would also be produced by a principal holding nothing at all.
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, token, http.MethodGet, "/api/jobs", ""))
	if rr.Code != http.StatusOK {
		t.Fatalf("jobs.read holder GET /api/jobs = %d %s, want 200", rr.Code, rr.Body.String())
	}

	// The discriminator: holding jobs.read must NOT authorize the write verbs.
	for _, verb := range []string{"cancel", "retry"} {
		rr := httptest.NewRecorder()
		env.mux.ServeHTTP(rr, bearer(t, token, http.MethodPost, "/api/jobs/"+job.ID+"/"+verb, ""))
		if rr.Code != http.StatusForbidden {
			t.Fatalf("jobs.read-only principal %s = %d %s, want 403 (the gate must be jobs.write)",
				verb, rr.Code, rr.Body.String())
		}
	}
	after, err := actionEnv.repository.GetJob(context.Background(), job.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if after.Status != jobs.StatusQueued || after.CancelRequested {
		t.Fatalf("a refused write mutated the job: %+v", after)
	}
}

// C1: the write routes exist only when the actions surface is bound. This is
// the handler-side half of the descriptor promise in modules/jobs: a read-only
// composition declares — and serves — no mutating route at all.
func TestJobActionRoutesAbsentWithoutActions(t *testing.T) {
	env := newAuthTestEnv(t)
	repository := jobs.NewRepository(env.st)
	for _, route := range JobsRoutes(env.a, repository, nil, nil, "admin.jobs") {
		env.mux.Handle(route.Method+" "+route.Pattern, route.Handler)
	}
	job, err := repository.Create(context.Background(), jobs.CreateInput{
		ID: "job-readonly", Kind: "jobs.batch-export", Payload: []byte(`{}`),
		ActorID: "some-other-actor", CorrelationID: "corr-readonly", MaxAttempts: 3, Now: time.Now().UTC(),
	})
	if err != nil {
		t.Fatalf("seed: %v", err)
	}
	token := adminToken(t, env)

	for _, verb := range []string{"cancel", "retry"} {
		rr := httptest.NewRecorder()
		env.mux.ServeHTTP(rr, bearer(t, token, http.MethodPost, "/api/jobs/"+job.ID+"/"+verb, ""))
		if rr.Code == http.StatusOK {
			t.Fatalf("POST cancel/retry must not be routed in a read-only composition")
		}
	}
	after, err := repository.GetJob(context.Background(), job.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if after.Status != jobs.StatusQueued || after.CancelRequested {
		t.Fatalf("an unrouted action mutated the job: %+v", after)
	}
}
