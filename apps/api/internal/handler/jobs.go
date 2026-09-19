// Generic async Job read surface (GOAL-003 R2 · VP-038).
//
// admin.jobs exposes the management-scope view of the durable Job runtime:
// registered job kinds can be listed, inspected and their results read by a
// holder of jobs.read. It deliberately does NOT replace the actor-scoped
// wallet surface — GetForActor and the /api/wallet/jobs/* routes keep their
// frozen semantics and tests (GOAL-002 D-001 §2.1). The two differ only in
// scope: wallet answers "my reconcile jobs", this answers "every job".
//
// R2 is read-only by design. Cancel/retry/download are result-center actions
// and land with R4; jobs.write is therefore declared there, not here.
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/jobs"
	"github.com/magicvr/schema-ui-core/apps/api/internal/requestid"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// JobReader is the read surface the admin.jobs routes consume.
type JobReader interface {
	ListJobs(ctx context.Context, filter jobs.ListFilter) ([]jobs.Job, int, error)
	GetJob(ctx context.Context, id string) (*jobs.Job, error)
}

// JobsBasePath is the admin.jobs route prefix. It is the base path handed to
// jobs.ResultURL, mirroring how admin.wallet declares its own prefix.
const JobsBasePath = "/api/jobs"

// JobBatchSubmitter enqueues the batch-export job (GOAL-004 R3). It is separate
// from the read surface so a read-only composition can omit it.
type JobBatchSubmitter interface {
	Supported(resource string) bool
	SubmitBatchExport(ctx context.Context, resource string, ids []string, actorID, correlationID string) (*jobs.Job, error)
}

// JobActions is the management-scope write surface (GOAL-005 R4): cancelling
// and retrying another actor's job. It is separate from the actor-scoped wallet
// JobService so those frozen paths stay untouched.
type JobActions interface {
	CancelAny(ctx context.Context, id string) (*jobs.Job, error)
	RetryAny(ctx context.Context, id string) (*jobs.Job, error)
}

// JobsRoutes returns the admin.jobs HTTP surface: the management-scope read
// routes (R2), the batch-export submit route (R3, when submitter is non-nil)
// and the management-scope cancel/retry actions (R4, when actions is non-nil).
// Every route is permission gated and fail-closed.
func JobsRoutes(a *auth.Authenticator, reader JobReader, submitter JobBatchSubmitter, actions JobActions, moduleID string) []kernel.RouteContribution {
	h := &jobsHandler{reader: reader, submitter: submitter}
	var routes []kernel.RouteContribution
	add := func(method, pattern string, handler http.Handler) {
		routes = append(routes, kernel.RouteContribution{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey(method, pattern)},
			Method:               method,
			Pattern:              pattern,
			Handler:              a.Middleware(handler),
		})
	}

	// Result-center write actions (R4). Both are management-scope operations on
	// jobs that may belong to another actor, so they gate on jobs.write — the
	// same key that authorizes submitting an async job. The state sets and
	// error codes mirror the actor-scoped contract exactly (GOAL-005 D-001 §2).
	if actions != nil {
		add("POST", JobsBasePath+"/{id}/cancel", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ok := requirePermission(w, r, "jobs.write"); !ok {
				return
			}
			job, err := actions.CancelAny(r.Context(), r.PathValue("id"))
			if err != nil {
				writeJobActionError(w, r, err)
				return
			}
			writeJSON(w, http.StatusOK, jobToMap(*job))
		}))
		add("POST", JobsBasePath+"/{id}/retry", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ok := requirePermission(w, r, "jobs.write"); !ok {
				return
			}
			job, err := actions.RetryAny(r.Context(), r.PathValue("id"))
			if err != nil {
				writeJobActionError(w, r, err)
				return
			}
			writeJSON(w, http.StatusOK, jobToMap(*job))
		}))
	}

	// Batch export (R3): a real batch operation carried by the Job runtime.
	// It answers 202 + a job projection; the caller polls GET /api/jobs/{id}.
	//
	// Two independent gates apply, deliberately: jobs.write authorizes
	// submitting an async job, and data.export authorizes moving the data out.
	// Requiring both keeps the async path from becoming a way around the
	// existing export permission (the synchronous exporter is data.export-gated
	// too), so adding the job runtime never widens data egress.
	if submitter != nil {
		add("POST", JobsBasePath+"/batch-export", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := requirePermission(w, r, "jobs.write")
			if !ok {
				return
			}
			if _, ok := requirePermission(w, r, "data.export"); !ok {
				return
			}
			h.submitBatchExport(w, r, user)
		}))
	}

	add("GET", JobsBasePath, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePermission(w, r, "jobs.read"); !ok {
			return
		}
		h.list(w, r)
	}))

	add("GET", JobsBasePath+"/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePermission(w, r, "jobs.read"); !ok {
			return
		}
		job, err := reader.GetJob(r.Context(), r.PathValue("id"))
		if err != nil {
			writeJobReadError(w, r, err)
			return
		}
		writeJSON(w, http.StatusOK, jobToMap(*job))
	}))

	add("GET", JobsBasePath+"/{id}/result", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePermission(w, r, "jobs.read"); !ok {
			return
		}
		job, err := reader.GetJob(r.Context(), r.PathValue("id"))
		if err != nil {
			writeJobReadError(w, r, err)
			return
		}
		// Same three-way result semantics as the wallet surface (GOAL-002
		// D-001 §1.2 K-3), so a job result behaves identically in both views.
		switch job.Status {
		case jobs.StatusQueued, jobs.StatusRunning:
			writeLocalizedError(w, r, http.StatusConflict, "JOB_RESULT_NOT_READY", "job result is not ready")
		case jobs.StatusExpired:
			writeLocalizedError(w, r, http.StatusGone, "JOB_RESULT_EXPIRED", "job result has expired")
		case jobs.StatusFailed, jobs.StatusCancelled:
			writeJSON(w, http.StatusOK, jobToMap(*job))
		case jobs.StatusSucceeded:
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", "job-"+job.ID+".json"))
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(job.Result)
		default:
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "invalid job status")
		}
	}))

	return routes
}

// batchExportRequest is the submit body: the target resource and the selected
// row keys, mirroring the shape the synchronous batch action sends.
type batchExportRequest struct {
	Resource string `json:"resource"`
	IDs      []any  `json:"ids"`
}

func (h *jobsHandler) submitBatchExport(w http.ResponseWriter, r *http.Request, user account.User) {
	var body batchExportRequest
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxResourceBodyBytes))
	if err := decoder.Decode(&body); err != nil {
		writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_BODY", "expected a JSON object with resource and ids")
		return
	}
	resource := strings.TrimSpace(body.Resource)
	if resource == "" {
		writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_BODY", "resource is required")
		return
	}
	if !h.submitter.Supported(resource) {
		writeLocalizedError(w, r, http.StatusNotFound, "RESOURCE_NOT_FOUND", "no batch export for that resource")
		return
	}
	if len(body.IDs) == 0 {
		writeLocalizedError(w, r, http.StatusBadRequest, "EMPTY_SELECTION", "ids must contain at least one key")
		return
	}
	// Scalar keys only, de-duplicated preserving order — the same selection
	// normalization the synchronous batch action applies (D3 invariants).
	seen := make(map[string]bool, len(body.IDs))
	ids := make([]string, 0, len(body.IDs))
	for _, raw := range body.IDs {
		var key string
		switch value := raw.(type) {
		case string:
			if value == "" {
				writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_SELECTION_KEY", "ids entries must be non-empty scalars")
				return
			}
			key = value
		case float64:
			if !isFiniteNumber(value) {
				writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_SELECTION_KEY", "ids entries must be finite scalars")
				return
			}
			key = formatNumberKey(value)
		case bool:
			key = strconv.FormatBool(value)
		default:
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_SELECTION_KEY", "ids entries must be scalar keys")
			return
		}
		if seen[key] {
			continue
		}
		seen[key] = true
		ids = append(ids, key)
	}
	if len(ids) == 0 {
		writeLocalizedError(w, r, http.StatusBadRequest, "EMPTY_SELECTION", "ids must contain at least one scalar key")
		return
	}

	correlationID := requestid.FromContext(r.Context())
	if correlationID == "" {
		correlationID = requestid.New()
	}
	job, err := h.submitter.SubmitBatchExport(r.Context(), resource, ids, user.ID, correlationID)
	if err != nil {
		if errors.Is(err, jobs.ErrInvalid) {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_SELECTION_KEY", "selection is not acceptable for batch export")
			return
		}
		writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not submit batch export")
		return
	}
	writeJSON(w, http.StatusAccepted, jobToMap(*job))
}

type jobsHandler struct {
	reader    JobReader
	submitter JobBatchSubmitter
}

func (h *jobsHandler) list(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page, ok := intParam(query.Get("page"), 1)
	if !ok {
		writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PAGE", "page must be a positive integer")
		return
	}
	pageSize, ok := intParam(query.Get("pageSize"), DefaultPageSize)
	if !ok || pageSize > maxPageSize {
		writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PAGE_SIZE", "pageSize must be a positive integer not exceeding 100")
		return
	}
	sort := query.Get("sort")
	if sort != "" && !jobSortAllowed(jobs.SortableJobFields, sort) {
		writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_SORT_FIELD", "unsupported sort field")
		return
	}
	order := query.Get("order")
	if order != "" && order != "asc" && order != "desc" {
		writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_SORT_ORDER", "order must be asc or desc")
		return
	}
	from, ok := parseJobTime(query.Get("from"))
	if !ok {
		writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_DATE_FILTER", "from must be an RFC3339 timestamp")
		return
	}
	to, ok := parseJobTime(query.Get("to"))
	if !ok {
		writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_DATE_FILTER", "to must be an RFC3339 timestamp")
		return
	}
	if status := query.Get("status"); status != "" && !isJobStatus(status) {
		writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_STATUS_FILTER", "status is not one of the six job states")
		return
	}

	items, total, err := h.reader.ListJobs(r.Context(), jobs.ListFilter{
		Kind: query.Get("kind"), Status: query.Get("status"), ActorID: query.Get("actorId"),
		From: from, To: to, Sort: sort, Order: order, Page: page, PageSize: pageSize,
	})
	if err != nil {
		writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not list jobs")
		return
	}
	rows := make([]map[string]any, 0, len(items))
	for _, job := range items {
		rows = append(rows, jobToMap(job))
	}
	// The shared list envelope, so the admin table renders it like every other
	// list surface (resources.go resourceList).
	writeJSON(w, http.StatusOK, resourceList{Items: rows, Total: total, Page: page, PageSize: pageSize})
}

// jobToMap is the management-scope Job projection. Like the wallet projection
// it deliberately does not embed payload/result (VP-012 D-002 §6); a succeeded
// job advertises its download address through the shared ResultURL derivation.
//
// R4 (GOAL-005 D-001 §2) adds the result center's operability fields. They are
// DERIVED HERE, on the server, so the row actions' enabled/disabled state is by
// construction the same rule the write routes enforce — a second copy of the
// state machine in the browser could drift, this cannot. The additions are
// purely additive: no R2 field's name or meaning changes.
func jobToMap(job jobs.Job) map[string]any {
	row := map[string]any{
		"id": job.ID, "kind": job.Kind, "status": job.Status,
		"progress": job.Progress, "attempt": job.Attempt, "maxAttempts": job.MaxAttempts,
		"cancelRequested": job.CancelRequested, "actorId": job.ActorID,
		"correlationId": job.CorrelationID,
		"createdAt":     job.CreatedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
		"updatedAt":     job.UpdatedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
		// Action availability, mirroring the write contract exactly
		// (internal/jobs/actions.go): cancel accepts queued|running, retry
		// accepts failed WITH remaining attempt budget, and only a succeeded job
		// has a result to download (an expired one has already lost it).
		"cancellable":  job.Status == jobs.StatusQueued || job.Status == jobs.StatusRunning,
		"retryable":    job.Status == jobs.StatusFailed && job.Attempt < job.MaxAttempts,
		"downloadable": job.Status == jobs.StatusSucceeded,
		"statusStyle":  jobStatusStyle(job.Status),
	}
	if job.ErrorCode != "" {
		row["error"] = map[string]any{"code": job.ErrorCode, "message": job.ErrorMessage}
		// Flat mirrors for the schema surfaces: a table column and a recordView
		// field both address top-level row keys (dotted paths are not resolved).
		row["errorCode"] = job.ErrorCode
		row["errorMessage"] = job.ErrorMessage
	}
	if job.FinishedAt != nil {
		row["finishedAt"] = job.FinishedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00")
	}
	if job.ResultExpiresAt != nil {
		row["resultExpiresAt"] = job.ResultExpiresAt.UTC().Format("2006-01-02T15:04:05.000Z07:00")
	}
	if job.Status == jobs.StatusSucceeded {
		row["resultUrl"] = jobs.ResultURL(JobsBasePath, job.ID)
	}
	return row
}

// jobStatusStyle maps the six frozen job states onto the renderer's badge
// presets (schema-table badgeClassesFor: success/warning/destructive/info, with
// anything else neutral). Exhaustive by construction: the default branch is
// unreachable for a valid Job and therefore falls back to the neutral badge
// rather than inventing a colour.
func jobStatusStyle(status jobs.Status) string {
	switch status {
	case jobs.StatusQueued:
		return "info"
	case jobs.StatusRunning:
		return "warning"
	case jobs.StatusSucceeded:
		return "success"
	case jobs.StatusFailed:
		return "destructive"
	default: // cancelled, expired
		return "neutral"
	}
}

func isJobStatus(value string) bool {
	switch jobs.Status(value) {
	case jobs.StatusQueued, jobs.StatusRunning, jobs.StatusSucceeded,
		jobs.StatusFailed, jobs.StatusCancelled, jobs.StatusExpired:
		return true
	}
	return false
}

// parseJobTime accepts an empty value (unbounded) or an RFC3339 timestamp.
func parseJobTime(value string) (time.Time, bool) {
	if value == "" {
		return time.Time{}, true
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, false
	}
	return parsed, true
}

// jobSortAllowed reports whether want is one of the whitelisted sort keys.
// (Named distinctly from the test-only containsString helper in this package.)
func jobSortAllowed(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

// writeJobReadError maps a read failure to the frozen Job error codes. A
// missing job is a 404; anything else is INTERNAL without leaking detail.
func writeJobReadError(w http.ResponseWriter, r *http.Request, err error) {
	if errors.Is(err, jobs.ErrNotFound) {
		writeLocalizedError(w, r, http.StatusNotFound, "JOB_NOT_FOUND", "job not found")
		return
	}
	writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not read job")
}

// writeJobActionError maps a cancel/retry failure to the frozen Job codes
// (shared by the R4 management-scope actions; the wallet surface has its own
// mapper because it also handles wallet-specific codes).
func writeJobActionError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, jobs.ErrNotFound):
		writeLocalizedError(w, r, http.StatusNotFound, "JOB_NOT_FOUND", "job not found")
	case errors.Is(err, jobs.ErrNotCancellable):
		writeLocalizedError(w, r, http.StatusConflict, "JOB_NOT_CANCELLABLE", "job cannot be cancelled")
	case errors.Is(err, jobs.ErrNotRetryable):
		writeLocalizedError(w, r, http.StatusConflict, "JOB_NOT_RETRYABLE", "job cannot be retried")
	default:
		writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not update job")
	}
}
