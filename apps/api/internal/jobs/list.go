package jobs

import (
	"context"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/pagination"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// ListFilter is the management-scope query for the generic Job read surface
// (GOAL-003 R2 · D-001 §1.3). Every field is optional; the zero value lists
// every job newest-first. Filtering is parameterized and sorting is resolved
// through a whitelist, so no caller input reaches SQL as an identifier.
type ListFilter struct {
	Kind    string
	Status  string
	ActorID string
	// From/To bound created_at (inclusive). Zero values mean unbounded.
	From time.Time
	To   time.Time
	// Sort is a whitelisted column key ("createdAt" | "updatedAt"); Order is
	// "asc" | "desc". Both default to created_at DESC when empty.
	Sort     string
	Order    string
	Page     int
	PageSize int
}

// SortableJobFields is the whitelist accepted by ListFilter.Sort. Anything
// else falls back to the default so an unknown key can never become a column.
var SortableJobFields = []string{"createdAt", "updatedAt"}

// jobSortSQL maps a whitelisted sort key + order to a fixed SQL fragment. The
// returned string is always a literal from this function, never caller text.
func jobSortSQL(sort, order string) string {
	column := "created_at"
	switch sort {
	case "createdAt":
		column = "created_at"
	case "updatedAt":
		column = "updated_at"
	}
	direction := "DESC"
	if strings.EqualFold(order, "asc") {
		direction = "ASC"
	}
	return column + " " + direction
}

// jobsWhere builds the parameterized WHERE clause and its argument list. It is
// shared by the count and the page query so both always agree on the filter.
func jobsWhere(filter ListFilter) (string, []any) {
	clauses := make([]string, 0, 5)
	args := make([]any, 0, 5)
	if kind := strings.TrimSpace(filter.Kind); kind != "" {
		clauses = append(clauses, "kind = ?")
		args = append(args, kind)
	}
	if status := strings.TrimSpace(filter.Status); status != "" {
		clauses = append(clauses, "status = ?")
		args = append(args, status)
	}
	if actorID := strings.TrimSpace(filter.ActorID); actorID != "" {
		clauses = append(clauses, "actor_id = ?")
		args = append(args, actorID)
	}
	if !filter.From.IsZero() {
		clauses = append(clauses, "created_at >= ?")
		args = append(args, filter.From.UTC())
	}
	if !filter.To.IsZero() {
		clauses = append(clauses, "created_at <= ?")
		args = append(args, filter.To.UTC())
	}
	if len(clauses) == 0 {
		return "", args
	}
	return " WHERE " + strings.Join(clauses, " AND "), args
}

// ListJobs returns one page of jobs plus the pre-pagination total for the same
// filter (management scope, GOAL-003 R2). It is the cross-actor counterpart of
// GetForActor and deliberately does NOT replace it: the actor-scoped paths
// (GetForActor / RequestCancel / Retry) keep their frozen semantics and tests.
func (r *Repository) ListJobs(ctx context.Context, filter ListFilter) ([]Job, int, error) {
	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	where, args := jobsWhere(filter)
	var jobs []Job
	var total int
	err := r.runner.Run(ctx, func(tx kernel.Tx) error {
		if err := tx.QueryRow(ctx, `SELECT COUNT(*) FROM jobs`+where, args...).Scan(&total); err != nil {
			return err
		}
		rows, err := tx.Query(ctx, `SELECT `+jobColumns+` FROM jobs`+where+
			` ORDER BY `+jobSortSQL(filter.Sort, filter.Order)+`, id DESC LIMIT ? OFFSET ?`,
			append(args, pageSize, pagination.Offset(page, pageSize, total))...)
		if err != nil {
			return err
		}
		defer rows.Close()
		jobs = make([]Job, 0, pageSize)
		for rows.Next() {
			job, err := scanJob(rows)
			if err != nil {
				return err
			}
			jobs = append(jobs, *job)
		}
		return rows.Err()
	})
	return jobs, total, err
}

// GetJob returns one job by id with no actor predicate. It is the
// management-scope read used by the admin.jobs detail route; the actor-scoped
// GetForActor is unchanged and remains the only path the wallet surface uses.
func (r *Repository) GetJob(ctx context.Context, id string) (*Job, error) {
	return r.Get(ctx, id)
}

// ResultURL returns the download address for a finished job result. Modules
// declare their own base path (admin.wallet keeps /api/wallet/jobs, admin.jobs
// uses /api/jobs) so the projection carries a correct address without a
// kind→route registry that could silently drift (GOAL-003 D-001 §2.3).
func ResultURL(basePath, id string) string {
	return strings.TrimSuffix(basePath, "/") + "/" + id + "/result"
}
