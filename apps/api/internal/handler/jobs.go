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
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/jobs"
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

// JobsRoutes returns the admin.jobs HTTP read surface: the management-scope
// list, the detail read and the result download. Every route is permission
// gated and fail-closed.
func JobsRoutes(a *auth.Authenticator, reader JobReader, moduleID string) []kernel.RouteContribution {
	h := &jobsHandler{reader: reader}
	var routes []kernel.RouteContribution
	add := func(method, pattern string, handler http.Handler) {
		routes = append(routes, kernel.RouteContribution{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey(method, pattern)},
			Method:               method,
			Pattern:              pattern,
			Handler:              a.Middleware(handler),
		})
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

type jobsHandler struct {
	reader JobReader
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
func jobToMap(job jobs.Job) map[string]any {
	row := map[string]any{
		"id": job.ID, "kind": job.Kind, "status": job.Status,
		"progress": job.Progress, "attempt": job.Attempt, "maxAttempts": job.MaxAttempts,
		"cancelRequested": job.CancelRequested, "actorId": job.ActorID,
		"correlationId": job.CorrelationID,
		"createdAt":     job.CreatedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
		"updatedAt":     job.UpdatedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
	}
	if job.ErrorCode != "" {
		row["error"] = map[string]any{"code": job.ErrorCode, "message": job.ErrorMessage}
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
