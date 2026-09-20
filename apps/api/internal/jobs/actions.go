package jobs

import (
	"context"
	"fmt"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// Management-scope job actions (GOAL-005 R4 · D-001).
//
// These mirror the actor-scoped RequestCancel / Retry exactly — same status
// sets, same transitions, same error codes — and differ only in the missing
// actor predicate. They live in their own file, and deliberately do NOT modify
// repository.go, so the frozen actor-scoped paths and their tests stay
// byte-identical (R1 D-001 §2.1: the actor isolation semantics must not be
// relaxed; the generic surface goes through new methods).

// RequestCancelAny requests cancellation of a job without an actor predicate.
//
// Contract (mirrors RequestCancel):
//   - queued  → cancelled immediately;
//   - running → cancel_requested=1 (the runner notices on its next heartbeat
//     and finalizes through FinalizeCancel);
//   - any terminal state (succeeded/failed/cancelled/expired) → ErrNotCancellable.
func (r *Repository) RequestCancelAny(ctx context.Context, id string, now time.Time) (*Job, error) {
	var job *Job
	err := r.runner.Run(ctx, func(tx kernel.Tx) error {
		current, err := getTx(ctx, tx, id)
		if err != nil {
			return err
		}
		switch current.Status {
		case StatusQueued:
			_, err = tx.Exec(ctx, `UPDATE jobs SET status='cancelled', cancel_requested=0,
updated_at=?, finished_at=? WHERE id=? AND status='queued'`,
				now, now, id)
		case StatusRunning:
			_, err = tx.Exec(ctx, `UPDATE jobs SET cancel_requested=1, updated_at=?
WHERE id=? AND status='running'`, now, id)
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

// RetryAny requeues a failed job without an actor predicate.
//
// Contract (mirrors Retry): only a failed job that still has attempt budget
// (attempt < max_attempts) is retryable; progress, lease, result and error are
// cleared so the retry starts from a clean queued state. Anything else —
// including a failed job whose attempts are exhausted — is ErrNotRetryable.
// Exhausted attempts are NOT reset: that would extend the Job contract, which
// VP-038 lists as a non-goal (GOAL-005 D-001 §2).
func (r *Repository) RetryAny(ctx context.Context, id string, now time.Time) (*Job, error) {
	var job *Job
	err := r.runner.Run(ctx, func(tx kernel.Tx) error {
		result, err := tx.Exec(ctx, `UPDATE jobs SET status='queued', progress=0,
cancel_requested=0, lease_owner=NULL, lease_expires_at=NULL, result=NULL,
error_code=NULL, error_message=NULL, updated_at=?, finished_at=NULL, expires_at=NULL
WHERE id=? AND status='failed' AND attempt < max_attempts`,
			now, id)
		if err != nil {
			return fmt.Errorf("retry job: %w", err)
		}
		if err := requireAffectedNotRetryable(ctx, tx, id, result); err != nil {
			return err
		}
		job, err = getTx(ctx, tx, id)
		return err
	})
	return job, err
}

// requireAffectedNotRetryable reports ErrNotFound when the job does not exist
// and ErrNotRetryable when it exists but is not in a retryable state.
func requireAffectedNotRetryable(ctx context.Context, tx kernel.Tx, id string, result kernel.Result) error {
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
	return ErrNotRetryable
}

// CancelAny cancels a job without an actor predicate. It mirrors Runner.Cancel:
// a running job's in-flight handler is cancelled through the active-execution
// registry so the worker stops promptly instead of waiting for its lease, and
// the scan loop is signalled so a queued cancellation is observed.
func (r *Runner) CancelAny(ctx context.Context, id string) (*Job, error) {
	job, err := r.repo.RequestCancelAny(ctx, id, r.options.Now().UTC())
	if err != nil {
		return nil, err
	}
	if job.Status == StatusRunning {
		r.mu.Lock()
		cancel := r.active[id]
		r.mu.Unlock()
		if cancel != nil {
			cancel()
		}
	}
	r.signal()
	return job, nil
}

// RetryAny requeues a failed job without an actor predicate, signalling the
// scan loop so the retry is picked up promptly.
func (r *Runner) RetryAny(ctx context.Context, id string) (*Job, error) {
	job, err := r.repo.RetryAny(ctx, id, r.options.Now().UTC())
	if err == nil {
		r.signal()
	}
	return job, err
}
