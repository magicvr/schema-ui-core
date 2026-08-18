package jobs

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type TxRunner interface {
	WithTx(context.Context, func(*sql.Tx) error) error
}

type Repository struct {
	runner TxRunner
}

func NewRepository(runner TxRunner) *Repository {
	return &Repository{runner: runner}
}

func (r *Repository) Create(ctx context.Context, input CreateInput) (*Job, error) {
	if strings.TrimSpace(input.ID) == "" || strings.TrimSpace(input.Kind) == "" ||
		strings.TrimSpace(input.ActorID) == "" || strings.TrimSpace(input.CorrelationID) == "" ||
		!json.Valid(input.Payload) {
		return nil, ErrInvalid
	}
	if input.MaxAttempts == 0 {
		input.MaxAttempts = DefaultMaxAttempts
	}
	if input.MaxAttempts < 1 {
		return nil, ErrInvalid
	}
	now := input.Now.UTC()
	var created *Job
	err := r.runner.WithTx(ctx, func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `INSERT INTO jobs (
id, kind, status, payload, progress, cancel_requested, attempt, max_attempts,
lease_version, actor_id, correlation_id, created_at, updated_at
) VALUES (?,?, 'queued', ?,0,0,0,?,0,?,?,?,?)`,
			input.ID, input.Kind, string(input.Payload), input.MaxAttempts,
			input.ActorID, input.CorrelationID, toMillis(now), toMillis(now),
		); err != nil {
			return fmt.Errorf("create job: %w", err)
		}
		var err error
		created, err = getTx(ctx, tx, input.ID)
		return err
	})
	return created, err
}

func (r *Repository) Get(ctx context.Context, id string) (*Job, error) {
	var job *Job
	err := r.runner.WithTx(ctx, func(tx *sql.Tx) error {
		var err error
		job, err = getTx(ctx, tx, id)
		return err
	})
	return job, err
}

func (r *Repository) GetForActor(ctx context.Context, id, kind, actorID string) (*Job, error) {
	var job *Job
	err := r.runner.WithTx(ctx, func(tx *sql.Tx) error {
		var err error
		job, err = scanJob(tx.QueryRowContext(ctx, `SELECT `+jobColumns+` FROM jobs WHERE id=? AND kind=? AND actor_id=?`, id, kind, actorID))
		return err
	})
	return job, err
}

func (r *Repository) Claim(ctx context.Context, id, owner string, now time.Time, leaseDuration time.Duration) (*Job, Lease, error) {
	if strings.TrimSpace(owner) == "" || leaseDuration <= 0 {
		return nil, Lease{}, ErrInvalid
	}
	now = now.UTC()
	var claimed *Job
	err := r.runner.WithTx(ctx, func(tx *sql.Tx) error {
		result, err := tx.ExecContext(ctx, `UPDATE jobs SET
status='running', attempt=attempt+1, lease_owner=?, lease_version=lease_version+1,
lease_expires_at=?, updated_at=?, finished_at=NULL, expires_at=NULL,
error_code=NULL, error_message=NULL
WHERE id=? AND (
  (status='queued' AND attempt < max_attempts)
  OR (status='running' AND cancel_requested=0 AND lease_expires_at <= ? AND attempt < max_attempts)
)`, owner, toMillis(now.Add(leaseDuration)), toMillis(now), id, toMillis(now))
		if err != nil {
			return fmt.Errorf("claim job: %w", err)
		}
		if err := requireAffected(ctx, tx, id, result, ErrTransition); err != nil {
			return err
		}
		claimed, err = getTx(ctx, tx, id)
		return err
	})
	if err != nil {
		return nil, Lease{}, err
	}
	return claimed, Lease{JobID: claimed.ID, Owner: claimed.LeaseOwner, Version: claimed.LeaseVersion}, nil
}

func (r *Repository) Heartbeat(ctx context.Context, lease Lease, now time.Time, leaseDuration time.Duration) error {
	if !validLease(lease) || leaseDuration <= 0 {
		return ErrInvalid
	}
	return r.updateLease(ctx, lease, `UPDATE jobs SET lease_expires_at=?, updated_at=?
WHERE id=? AND status='running' AND lease_owner=? AND lease_version=?`,
		toMillis(now.Add(leaseDuration)), toMillis(now), lease.JobID, lease.Owner, lease.Version)
}

func (r *Repository) UpdateProgress(ctx context.Context, lease Lease, progress int, now time.Time) error {
	if !validLease(lease) || progress < 0 || progress > 99 {
		return ErrInvalid
	}
	return r.updateGuardedLease(ctx, lease, ErrTransition, `UPDATE jobs SET progress=?, updated_at=?
WHERE id=? AND status='running' AND lease_owner=? AND lease_version=? AND progress <= ?`,
		progress, toMillis(now), lease.JobID, lease.Owner, lease.Version, progress)
}

func (r *Repository) RequestCancel(ctx context.Context, id, actorID string, now time.Time) (*Job, error) {
	var job *Job
	err := r.runner.WithTx(ctx, func(tx *sql.Tx) error {
		current, err := getForActorTx(ctx, tx, id, actorID)
		if err != nil {
			return err
		}
		switch current.Status {
		case StatusQueued:
			_, err = tx.ExecContext(ctx, `UPDATE jobs SET status='cancelled', cancel_requested=0,
updated_at=?, finished_at=? WHERE id=? AND status='queued' AND actor_id=?`,
				toMillis(now), toMillis(now), id, actorID)
		case StatusRunning:
			_, err = tx.ExecContext(ctx, `UPDATE jobs SET cancel_requested=1, updated_at=?
WHERE id=? AND status='running' AND actor_id=?`, toMillis(now), id, actorID)
		default:
			return ErrNotCancellable
		}
		if err != nil {
			return fmt.Errorf("request job cancellation: %w", err)
		}
		job, err = getTx(ctx, tx, id)
		return err
	})
	return job, err
}

func (r *Repository) FinalizeCancel(ctx context.Context, lease Lease, now time.Time) error {
	if !validLease(lease) {
		return ErrInvalid
	}
	return r.updateGuardedLease(ctx, lease, ErrTransition, `UPDATE jobs SET status='cancelled', cancel_requested=0,
lease_owner=NULL, lease_expires_at=NULL, result=NULL, error_code=NULL, error_message=NULL,
updated_at=?, finished_at=?, expires_at=NULL
WHERE id=? AND status='running' AND lease_owner=? AND lease_version=? AND cancel_requested=1`,
		toMillis(now), toMillis(now), lease.JobID, lease.Owner, lease.Version)
}

func (r *Repository) Fail(ctx context.Context, lease Lease, code, message string, now time.Time) error {
	if !validLease(lease) || strings.TrimSpace(code) == "" {
		return ErrInvalid
	}
	return r.updateGuardedLease(ctx, lease, ErrTransition, `UPDATE jobs SET status='failed', cancel_requested=0,
lease_owner=NULL, lease_expires_at=NULL, result=NULL, error_code=?, error_message=?,
updated_at=?, finished_at=?, expires_at=NULL
WHERE id=? AND status='running' AND lease_owner=? AND lease_version=? AND cancel_requested=0`,
		code, message, toMillis(now), toMillis(now), lease.JobID, lease.Owner, lease.Version)
}

// CompleteWithCommit atomically commits the consumer's durable result and the
// Job success state. The callback must use the provided transaction and must
// not open a nested transaction.
func (r *Repository) CompleteWithCommit(
	ctx context.Context,
	lease Lease,
	now time.Time,
	resultTTL time.Duration,
	commit func(*sql.Tx) (json.RawMessage, error),
) (*Job, error) {
	if !validLease(lease) || resultTTL <= 0 || commit == nil {
		return nil, ErrInvalid
	}
	var completed *Job
	err := r.runner.WithTx(ctx, func(tx *sql.Tx) error {
		current, err := getTx(ctx, tx, lease.JobID)
		if err != nil {
			return err
		}
		if current.Status != StatusRunning || current.LeaseOwner != lease.Owner || current.LeaseVersion != lease.Version {
			return ErrLeaseLost
		}
		payload, err := commit(tx)
		if err != nil {
			return err
		}
		if len(payload) == 0 || !json.Valid(payload) {
			return ErrInvalid
		}
		result, err := tx.ExecContext(ctx, `UPDATE jobs SET status='succeeded', progress=100,
cancel_requested=0, lease_owner=NULL, lease_expires_at=NULL, result=?,
error_code=NULL, error_message=NULL, updated_at=?, finished_at=?, expires_at=?
WHERE id=? AND status='running' AND lease_owner=? AND lease_version=?`,
			string(payload), toMillis(now), toMillis(now), toMillis(now.Add(resultTTL)),
			lease.JobID, lease.Owner, lease.Version)
		if err != nil {
			return fmt.Errorf("complete job: %w", err)
		}
		if err := affectedExactlyOne(result, ErrLeaseLost); err != nil {
			return err
		}
		completed, err = getTx(ctx, tx, lease.JobID)
		return err
	})
	return completed, err
}

func (r *Repository) Retry(ctx context.Context, id, actorID string, now time.Time) (*Job, error) {
	var job *Job
	err := r.runner.WithTx(ctx, func(tx *sql.Tx) error {
		result, err := tx.ExecContext(ctx, `UPDATE jobs SET status='queued', progress=0,
cancel_requested=0, lease_owner=NULL, lease_expires_at=NULL, result=NULL,
error_code=NULL, error_message=NULL, updated_at=?, finished_at=NULL, expires_at=NULL
WHERE id=? AND actor_id=? AND status='failed' AND attempt < max_attempts`,
			toMillis(now), id, actorID)
		if err != nil {
			return fmt.Errorf("retry job: %w", err)
		}
		if err := requireAffectedForActor(ctx, tx, id, actorID, result, ErrNotRetryable); err != nil {
			return err
		}
		job, err = getTx(ctx, tx, id)
		return err
	})
	return job, err
}

func (r *Repository) ExpireIfDue(ctx context.Context, id string, now time.Time) (*Job, error) {
	var job *Job
	err := r.runner.WithTx(ctx, func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `UPDATE jobs SET status='expired', result=NULL, updated_at=?
WHERE id=? AND status='succeeded' AND expires_at <= ?`, toMillis(now), id, toMillis(now)); err != nil {
			return fmt.Errorf("expire job result: %w", err)
		}
		var err error
		job, err = getTx(ctx, tx, id)
		return err
	})
	return job, err
}

func (r *Repository) ExpireDue(ctx context.Context, now time.Time) (int64, error) {
	return r.bulkTransition(ctx, `UPDATE jobs SET status='expired', result=NULL, updated_at=?
WHERE status='succeeded' AND expires_at <= ?`, toMillis(now), toMillis(now))
}

func (r *Repository) RecoverCancelledDue(ctx context.Context, now time.Time) (int64, error) {
	transitioned, err := r.RecoverCancelledDueJobs(ctx, now)
	return int64(len(transitioned)), err
}

func (r *Repository) RecoverCancelledDueJobs(ctx context.Context, now time.Time) ([]Job, error) {
	return r.transitionJobs(ctx,
		`SELECT id FROM jobs WHERE status='running' AND cancel_requested=1 AND lease_expires_at <= ?`,
		[]any{toMillis(now)},
		`UPDATE jobs SET status='cancelled', cancel_requested=0,
lease_owner=NULL, lease_expires_at=NULL, result=NULL, error_code=NULL, error_message=NULL,
updated_at=?, finished_at=?, expires_at=NULL
WHERE status='running' AND cancel_requested=1 AND lease_expires_at <= ?`,
		toMillis(now), toMillis(now), toMillis(now))
}

func (r *Repository) ExhaustExpired(ctx context.Context, now time.Time) (int64, error) {
	transitioned, err := r.ExhaustExpiredJobs(ctx, now)
	return int64(len(transitioned)), err
}

func (r *Repository) ExhaustExpiredJobs(ctx context.Context, now time.Time) ([]Job, error) {
	return r.transitionJobs(ctx,
		`SELECT id FROM jobs WHERE status='running' AND cancel_requested=0 AND lease_expires_at <= ? AND attempt >= max_attempts`,
		[]any{toMillis(now)},
		`UPDATE jobs SET status='failed',
lease_owner=NULL, lease_expires_at=NULL, result=NULL,
error_code='JOB_ATTEMPTS_EXHAUSTED', error_message='job attempts exhausted',
updated_at=?, finished_at=?, expires_at=NULL
WHERE status='running' AND cancel_requested=0 AND lease_expires_at <= ? AND attempt >= max_attempts`,
		toMillis(now), toMillis(now), toMillis(now))
}

func (r *Repository) ListRunnable(ctx context.Context, now time.Time, limit int) ([]Job, error) {
	if limit <= 0 {
		return nil, ErrInvalid
	}
	var jobs []Job
	err := r.runner.WithTx(ctx, func(tx *sql.Tx) error {
		rows, err := tx.QueryContext(ctx, `SELECT `+jobColumns+` FROM jobs
WHERE (status='queued' AND attempt < max_attempts)
   OR (status='running' AND cancel_requested=0 AND lease_expires_at <= ? AND attempt < max_attempts)
ORDER BY created_at, id LIMIT ?`, toMillis(now), limit)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			job, err := scanJob(rows)
			if err != nil {
				return err
			}
			jobs = append(jobs, *job)
		}
		return rows.Err()
	})
	return jobs, err
}

func (r *Repository) IsCancelRequested(ctx context.Context, lease Lease) (bool, error) {
	if !validLease(lease) {
		return false, ErrInvalid
	}
	var requested bool
	err := r.runner.WithTx(ctx, func(tx *sql.Tx) error {
		var value int
		err := tx.QueryRowContext(ctx, `SELECT cancel_requested FROM jobs
WHERE id=? AND status='running' AND lease_owner=? AND lease_version=?`,
			lease.JobID, lease.Owner, lease.Version).Scan(&value)
		if errors.Is(err, sql.ErrNoRows) {
			return ErrLeaseLost
		}
		requested = value != 0
		return err
	})
	return requested, err
}

func (r *Repository) updateLease(ctx context.Context, lease Lease, query string, args ...any) error {
	return r.runner.WithTx(ctx, func(tx *sql.Tx) error {
		result, err := tx.ExecContext(ctx, query, args...)
		if err != nil {
			return err
		}
		return affectedExactlyOne(result, ErrLeaseLost)
	})
}

func (r *Repository) updateGuardedLease(ctx context.Context, lease Lease, guardError error, query string, args ...any) error {
	return r.runner.WithTx(ctx, func(tx *sql.Tx) error {
		result, err := tx.ExecContext(ctx, query, args...)
		if err != nil {
			return err
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if affected == 1 {
			return nil
		}
		current, err := getTx(ctx, tx, lease.JobID)
		if err != nil {
			return err
		}
		if current.Status != StatusRunning || current.LeaseOwner != lease.Owner || current.LeaseVersion != lease.Version {
			return ErrLeaseLost
		}
		return guardError
	})
}

func (r *Repository) bulkTransition(ctx context.Context, query string, args ...any) (int64, error) {
	var affected int64
	err := r.runner.WithTx(ctx, func(tx *sql.Tx) error {
		result, err := tx.ExecContext(ctx, query, args...)
		if err != nil {
			return err
		}
		affected, err = result.RowsAffected()
		return err
	})
	return affected, err
}

func (r *Repository) transitionJobs(ctx context.Context, selectQuery string, selectArgs []any, updateQuery string, updateArgs ...any) ([]Job, error) {
	var transitioned []Job
	err := r.runner.WithTx(ctx, func(tx *sql.Tx) error {
		rows, err := tx.QueryContext(ctx, selectQuery, selectArgs...)
		if err != nil {
			return err
		}
		var ids []string
		for rows.Next() {
			var id string
			if err := rows.Scan(&id); err != nil {
				_ = rows.Close()
				return err
			}
			ids = append(ids, id)
		}
		if err := rows.Close(); err != nil {
			return err
		}
		if err := rows.Err(); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, updateQuery, updateArgs...); err != nil {
			return err
		}
		for _, id := range ids {
			job, err := getTx(ctx, tx, id)
			if err != nil {
				return err
			}
			transitioned = append(transitioned, *job)
		}
		return nil
	})
	return transitioned, err
}

func validLease(lease Lease) bool {
	return lease.JobID != "" && lease.Owner != "" && lease.Version > 0
}

func affectedExactlyOne(result sql.Result, zeroError error) error {
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return zeroError
	}
	return nil
}

func requireAffected(ctx context.Context, tx *sql.Tx, id string, result sql.Result, transitionError error) error {
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 1 {
		return nil
	}
	if _, err := getTx(ctx, tx, id); err != nil {
		return err
	}
	return transitionError
}

func requireAffectedForActor(ctx context.Context, tx *sql.Tx, id, actorID string, result sql.Result, transitionError error) error {
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 1 {
		return nil
	}
	if _, err := getForActorTx(ctx, tx, id, actorID); err != nil {
		return err
	}
	return transitionError
}

func getTx(ctx context.Context, tx *sql.Tx, id string) (*Job, error) {
	return scanJob(tx.QueryRowContext(ctx, `SELECT `+jobColumns+` FROM jobs WHERE id=?`, id))
}

func getForActorTx(ctx context.Context, tx *sql.Tx, id, actorID string) (*Job, error) {
	return scanJob(tx.QueryRowContext(ctx, `SELECT `+jobColumns+` FROM jobs WHERE id=? AND actor_id=?`, id, actorID))
}
