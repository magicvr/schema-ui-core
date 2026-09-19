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
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// mountJobsRoutes wires the admin.jobs read surface over the shared Job
// repository (GOAL-003 R2).
func mountJobsRoutes(t *testing.T, env *authTestEnv, repository *jobs.Repository) {
	t.Helper()
	for _, route := range JobsRoutes(env.a, repository, "admin.jobs") {
		env.mux.Handle(route.Method+" "+route.Pattern, route.Handler)
	}
}

func newJobsTestRepository(t *testing.T, env *authTestEnv) *jobs.Repository {
	t.Helper()
	return jobs.NewRepository(env.st)
}

// seedJob inserts one job row directly so the read surface can be exercised
// without running the worker.
func seedJob(t *testing.T, repository *jobs.Repository, id, kind, actorID string, at time.Time) *jobs.Job {
	t.Helper()
	job, err := repository.Create(context.Background(), jobs.CreateInput{
		ID: id, Kind: kind, Payload: []byte(`{}`),
		ActorID: actorID, CorrelationID: "corr-" + id, Now: at,
	})
	if err != nil {
		t.Fatalf("seed %s: %v", id, err)
	}
	return job
}

// C3: the management-scope list returns every actor's jobs behind jobs.read,
// in the shared list envelope.
func TestJobsListIsManagementScope(t *testing.T) {
	env := newAuthTestEnv(t)
	repository := newJobsTestRepository(t, env)
	mountJobsRoutes(t, env, repository)

	now := time.Now().UTC()
	seedJob(t, repository, "job-a", "wallet.reconcile", "user-1", now.Add(-3*time.Minute))
	seedJob(t, repository, "job-b", "wallet.reconcile", "user-2", now.Add(-2*time.Minute))
	seedJob(t, repository, "job-c", "jobs.batch-export", "user-9", now.Add(-1*time.Minute))

	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, adminToken(t, env), http.MethodGet, "/api/jobs", ""))
	if rr.Code != http.StatusOK {
		t.Fatalf("GET /api/jobs = %d body=%s", rr.Code, rr.Body.String())
	}
	var envelope struct {
		Items []struct {
			ID      string `json:"id"`
			Kind    string `json:"kind"`
			ActorID string `json:"actorId"`
			Status  string `json:"status"`
		} `json:"items"`
		Total    int `json:"total"`
		Page     int `json:"page"`
		PageSize int `json:"pageSize"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &envelope); err != nil {
		t.Fatal(err)
	}
	if envelope.Total != 3 || len(envelope.Items) != 3 {
		t.Fatalf("total=%d items=%d, want 3/3 (management scope sees every actor)", envelope.Total, len(envelope.Items))
	}
	// Newest first, and rows owned by three different actors are all present.
	if envelope.Items[0].ID != "job-c" {
		t.Fatalf("first item = %s, want job-c (created_at DESC)", envelope.Items[0].ID)
	}
	actors := map[string]bool{}
	for _, item := range envelope.Items {
		actors[item.ActorID] = true
	}
	if len(actors) != 3 {
		t.Fatalf("actors = %v, want three distinct actors", actors)
	}
	if envelope.Page != 1 || envelope.PageSize != DefaultPageSize {
		t.Fatalf("page=%d pageSize=%d, want 1/%d", envelope.Page, envelope.PageSize, DefaultPageSize)
	}
}

// C3: job filters narrow the result and stay parameterized.
func TestJobsListFiltersAndValidates(t *testing.T) {
	env := newAuthTestEnv(t)
	repository := newJobsTestRepository(t, env)
	mountJobsRoutes(t, env, repository)

	now := time.Now().UTC()
	seedJob(t, repository, "job-a", "wallet.reconcile", "user-1", now.Add(-3*time.Minute))
	seedJob(t, repository, "job-b", "jobs.batch-export", "user-1", now.Add(-2*time.Minute))
	seedJob(t, repository, "job-c", "jobs.batch-export", "user-2", now.Add(-1*time.Minute))
	token := adminToken(t, env)

	readTotal := func(query string) int {
		t.Helper()
		rr := httptest.NewRecorder()
		env.mux.ServeHTTP(rr, bearer(t, token, http.MethodGet, "/api/jobs"+query, ""))
		if rr.Code != http.StatusOK {
			t.Fatalf("GET /api/jobs%s = %d body=%s", query, rr.Code, rr.Body.String())
		}
		var envelope struct {
			Total int `json:"total"`
		}
		if err := json.Unmarshal(rr.Body.Bytes(), &envelope); err != nil {
			t.Fatal(err)
		}
		return envelope.Total
	}

	if got := readTotal("?kind=jobs.batch-export"); got != 2 {
		t.Fatalf("kind filter total = %d, want 2", got)
	}
	if got := readTotal("?actorId=user-1"); got != 2 {
		t.Fatalf("actor filter total = %d, want 2", got)
	}
	if got := readTotal("?status=queued"); got != 3 {
		t.Fatalf("status filter total = %d, want 3", got)
	}
	if got := readTotal("?kind=jobs.batch-export&actorId=user-2"); got != 1 {
		t.Fatalf("combined filter total = %d, want 1", got)
	}

	// Invalid inputs are rejected with the frozen list-validation codes.
	for _, tc := range []struct {
		query string
		want  string
	}{
		{"?status=not-a-state", "INVALID_STATUS_FILTER"},
		// A SQL-shaped sort key must be rejected by the whitelist (the value is
		// percent-encoded so the request line itself stays well formed).
		{"?sort=createdAt%3BDROP%20TABLE%20jobs", "INVALID_SORT_FIELD"},
		{"?sort=drop", "INVALID_SORT_FIELD"},
		{"?order=sideways", "INVALID_SORT_ORDER"},
		{"?from=not-a-time", "INVALID_DATE_FILTER"},
		{"?to=yesterday", "INVALID_DATE_FILTER"},
		{"?page=0", "INVALID_PAGE"},
		{"?pageSize=500", "INVALID_PAGE_SIZE"},
	} {
		rr := httptest.NewRecorder()
		env.mux.ServeHTTP(rr, bearer(t, token, http.MethodGet, "/api/jobs"+tc.query, ""))
		if rr.Code != http.StatusBadRequest || !bodyHasCode(rr, tc.want) {
			t.Fatalf("GET /api/jobs%s = %d %s, want 400 %s", tc.query, rr.Code, rr.Body.String(), tc.want)
		}
	}
}

// C3: permission gating is fail-closed — anonymous is 401 and a non-admin is
// 403, on every admin.jobs route.
func TestJobsRoutesGates(t *testing.T) {
	env := newAuthTestEnv(t)
	repository := newJobsTestRepository(t, env)
	mountJobsRoutes(t, env, repository)
	seedJob(t, repository, "job-1", "wallet.reconcile", "user-1", time.Now().UTC())

	paths := []struct{ method, path string }{
		{http.MethodGet, "/api/jobs"},
		{http.MethodGet, "/api/jobs/job-1"},
		{http.MethodGet, "/api/jobs/job-1/result"},
	}

	// Anonymous: the auth middleware rejects before the permission gate.
	for _, tc := range paths {
		rr := httptest.NewRecorder()
		env.mux.ServeHTTP(rr, httptest.NewRequest(tc.method, tc.path, nil))
		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("anonymous %s %s = %d, want 401", tc.method, tc.path, rr.Code)
		}
	}

	// An editor session does not hold jobs.read (PolicyAdmin), so it is 403.
	env.addUser(t, "editor1", "editor-password", []string{"editor"})
	editorToken := env.login(t, "editor1", "editor-password")
	for _, tc := range paths {
		rr := httptest.NewRecorder()
		env.mux.ServeHTTP(rr, bearer(t, editorToken, tc.method, tc.path, ""))
		if rr.Code != http.StatusForbidden {
			t.Fatalf("editor %s %s = %d body=%s, want 403", tc.method, tc.path, rr.Code, rr.Body.String())
		}
	}
}

// C3: detail and result semantics mirror the wallet surface.
func TestJobsDetailAndResult(t *testing.T) {
	env := newAuthTestEnv(t)
	repository := newJobsTestRepository(t, env)
	mountJobsRoutes(t, env, repository)
	token := adminToken(t, env)

	queued := seedJob(t, repository, "job-queued", "wallet.reconcile", "user-1", time.Now().UTC())

	// Detail: 200 with the projection (no payload/result embedding).
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, token, http.MethodGet, "/api/jobs/"+queued.ID, ""))
	if rr.Code != http.StatusOK {
		t.Fatalf("detail = %d body=%s", rr.Code, rr.Body.String())
	}
	var row map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &row); err != nil {
		t.Fatal(err)
	}
	if row["id"] != queued.ID || row["kind"] != "wallet.reconcile" || row["status"] != "queued" {
		t.Fatalf("row = %v", row)
	}
	if _, ok := row["result"]; ok {
		t.Fatalf("detail must not embed result (VP-012 D-002 §6)")
	}
	if _, ok := row["payload"]; ok {
		t.Fatalf("detail must not embed payload (VP-012 D-002 §6)")
	}

	// Result on a non-terminal job: 409 not ready.
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, token, http.MethodGet, "/api/jobs/"+queued.ID+"/result", ""))
	if rr.Code != http.StatusConflict || !bodyHasCode(rr, "JOB_RESULT_NOT_READY") {
		t.Fatalf("queued result = %d %s, want 409 JOB_RESULT_NOT_READY", rr.Code, rr.Body.String())
	}

	// Missing job: 404 JOB_NOT_FOUND, no existence leak.
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, token, http.MethodGet, "/api/jobs/job-missing", ""))
	if rr.Code != http.StatusNotFound || !bodyHasCode(rr, "JOB_NOT_FOUND") {
		t.Fatalf("missing detail = %d %s, want 404 JOB_NOT_FOUND", rr.Code, rr.Body.String())
	}

	// Succeeded job: the result endpoint serves the bytes as an attachment and
	// the projection advertises the admin.jobs result address.
	succeeded := runJobToTerminal(t, env, repository, jobs.StatusSucceeded)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, token, http.MethodGet, "/api/jobs/"+succeeded.ID+"/result", ""))
	if rr.Code != http.StatusOK {
		t.Fatalf("succeeded result = %d body=%s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Header().Get("Content-Disposition"), succeeded.ID) {
		t.Fatalf("result headers = %v, want a filename carrying the job id", rr.Header())
	}
	if strings.TrimSpace(rr.Body.String()) != `{"ok":true}` {
		t.Fatalf("result body = %s, want the stored result", rr.Body.String())
	}

	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, token, http.MethodGet, "/api/jobs/"+succeeded.ID, ""))
	if rr.Code != http.StatusOK {
		t.Fatalf("succeeded detail = %d", rr.Code)
	}
	var done map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &done); err != nil {
		t.Fatal(err)
	}
	if done["resultUrl"] != "/api/jobs/"+succeeded.ID+"/result" {
		t.Fatalf("resultUrl = %v, want the admin.jobs address", done["resultUrl"])
	}
	if done["progress"] != float64(100) {
		t.Fatalf("progress = %v, want 100", done["progress"])
	}
}

// runJobToTerminal drives one job through a real runner to the given terminal
// state, so the read surface is exercised against genuine runtime rows.
func runJobToTerminal(t *testing.T, env *authTestEnv, repository *jobs.Repository, want jobs.Status) *jobs.Job {
	t.Helper()
	options := jobs.DefaultRunnerOptions()
	options.HeartbeatInterval = 5 * time.Millisecond
	options.LeaseDuration = 50 * time.Millisecond
	options.ScanInterval = 5 * time.Millisecond
	runner, err := jobs.NewRunner(repository, options)
	if err != nil {
		t.Fatal(err)
	}
	kind := "jobs.test-terminal"
	if err := runner.Register(kind, func(context.Context, jobs.Job, jobs.Reporter) (jobs.CommitFunc, error) {
		return func(kernel.Tx) (json.RawMessage, error) {
			return json.RawMessage(`{"ok":true}`), nil
		}, nil
	}); err != nil {
		t.Fatal(err)
	}
	_ = runner.Start()
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_ = runner.Stop(ctx)
	})

	now := time.Now().UTC()
	id, err := jobs.NewID(now)
	if err != nil {
		t.Fatal(err)
	}
	job, err := runner.Submit(context.Background(), jobs.CreateInput{
		ID: id, Kind: kind, Payload: []byte(`{}`),
		ActorID: "user-1", CorrelationID: "corr-terminal", Now: now,
	})
	if err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		current, err := repository.GetJob(context.Background(), job.ID)
		if err != nil {
			t.Fatal(err)
		}
		if current.Status == want {
			return current
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatalf("job %s never reached %s", job.ID, want)
	return nil
}
