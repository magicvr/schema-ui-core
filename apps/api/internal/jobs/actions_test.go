package jobs_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/jobs"
)

// GOAL-005 R4 (D-001 §2): the management-scope actions ARE the actor-scoped
// contract minus the actor predicate. This file pins both halves: the new
// Any-variants accept another actor's job and enforce the same state sets and
// error codes, and the frozen actor-scoped variants keep refusing it.

// createJobFor seeds one job owned by actorID.
func createJobFor(t *testing.T, repo *jobs.Repository, id, actorID string, maxAttempts int) *jobs.Job {
	t.Helper()
	job, err := repo.Create(context.Background(), jobs.CreateInput{
		ID: id, Kind: "jobs.batch-export", Payload: []byte(`{}`),
		ActorID: actorID, CorrelationID: "corr-" + id, MaxAttempts: maxAttempts, Now: testNow,
	})
	if err != nil {
		t.Fatalf("create %s: %v", id, err)
	}
	return job
}

// C1: cancelling another actor's queued job is allowed and lands cancelled
// immediately — the same transition RequestCancel performs for the owner.
func TestRequestCancelAnyCancelsAnotherActorsQueuedJob(t *testing.T) {
	repo, _ := newRepository(t)
	job := createJobFor(t, repo, "job-any-queued", "someone-else", 3)

	// The actor-scoped path must still refuse a non-owner: this is the frozen
	// isolation the generic surface is deliberately allowed to bypass.
	if _, err := repo.RequestCancel(context.Background(), job.ID, "operator", testNow); !errors.Is(err, jobs.ErrNotFound) {
		t.Fatalf("RequestCancel by a non-owner = %v, want ErrNotFound (frozen isolation)", err)
	}

	cancelled, err := repo.RequestCancelAny(context.Background(), job.ID, testNow)
	if err != nil {
		t.Fatalf("RequestCancelAny: %v", err)
	}
	if cancelled.Status != jobs.StatusCancelled {
		t.Fatalf("status = %s, want cancelled", cancelled.Status)
	}
	if cancelled.CancelRequested {
		t.Fatalf("a queued cancellation must not leave cancel_requested set: %+v", cancelled)
	}
	if cancelled.FinishedAt == nil {
		t.Fatalf("a queued cancellation must set finished_at: %+v", cancelled)
	}
	if cancelled.ActorID != "someone-else" {
		t.Fatalf("cancel must not rewrite ownership: actor = %s", cancelled.ActorID)
	}
}

// C1: a running job is marked for cancellation rather than force-finished; the
// runner finalizes it through FinalizeCancel on its next heartbeat.
func TestRequestCancelAnyMarksAnotherActorsRunningJob(t *testing.T) {
	repo, _ := newRepository(t)
	job := createJobFor(t, repo, "job-any-running", "someone-else", 3)
	if _, _, err := repo.Claim(context.Background(), job.ID, "worker-1", testNow, time.Minute); err != nil {
		t.Fatalf("claim: %v", err)
	}

	marked, err := repo.RequestCancelAny(context.Background(), job.ID, testNow)
	if err != nil {
		t.Fatalf("RequestCancelAny: %v", err)
	}
	if marked.Status != jobs.StatusRunning || !marked.CancelRequested {
		t.Fatalf("running cancel mark = %+v, want running + cancel_requested", marked)
	}
	if marked.FinishedAt != nil {
		t.Fatalf("marking must not finish the job: %+v", marked)
	}
}

// C1: every terminal state is refused with the frozen JOB_NOT_CANCELLABLE code.
func TestRequestCancelAnyRefusesTerminalStates(t *testing.T) {
	cases := []struct {
		name  string
		state func(t *testing.T, repo *jobs.Repository) string
	}{
		{"succeeded", func(t *testing.T, repo *jobs.Repository) string {
			job := createJobFor(t, repo, "job-t-succeeded", "someone-else", 3)
			_, lease, err := repo.Claim(context.Background(), job.ID, "w", testNow, time.Minute)
			if err != nil {
				t.Fatalf("claim: %v", err)
			}
			if _, err := repo.CompleteWithCommit(context.Background(), lease, testNow, time.Hour, successCommit(`{"ok":true}`)); err != nil {
				t.Fatalf("complete: %v", err)
			}
			return job.ID
		}},
		{"failed", func(t *testing.T, repo *jobs.Repository) string {
			job := createJobFor(t, repo, "job-t-failed", "someone-else", 3)
			_, lease, err := repo.Claim(context.Background(), job.ID, "w", testNow, time.Minute)
			if err != nil {
				t.Fatalf("claim: %v", err)
			}
			if err := repo.Fail(context.Background(), lease, "JOB_HANDLER_FAILED", "boom", testNow); err != nil {
				t.Fatalf("fail: %v", err)
			}
			return job.ID
		}},
		{"cancelled", func(t *testing.T, repo *jobs.Repository) string {
			job := createJobFor(t, repo, "job-t-cancelled", "someone-else", 3)
			if _, err := repo.RequestCancelAny(context.Background(), job.ID, testNow); err != nil {
				t.Fatalf("cancel: %v", err)
			}
			return job.ID
		}},
		{"expired", func(t *testing.T, repo *jobs.Repository) string {
			job := createJobFor(t, repo, "job-t-expired", "someone-else", 3)
			_, lease, err := repo.Claim(context.Background(), job.ID, "w", testNow, time.Minute)
			if err != nil {
				t.Fatalf("claim: %v", err)
			}
			if _, err := repo.CompleteWithCommit(context.Background(), lease, testNow, time.Hour, successCommit(`{"ok":true}`)); err != nil {
				t.Fatalf("complete: %v", err)
			}
			if _, err := repo.ExpireIfDue(context.Background(), job.ID, testNow.Add(2*time.Hour)); err != nil {
				t.Fatalf("expire: %v", err)
			}
			return job.ID
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo, _ := newRepository(t)
			id := tc.state(t, repo)
			before, err := repo.GetJob(context.Background(), id)
			if err != nil {
				t.Fatalf("get: %v", err)
			}
			if _, err := repo.RequestCancelAny(context.Background(), id, testNow.Add(3*time.Hour)); !errors.Is(err, jobs.ErrNotCancellable) {
				t.Fatalf("cancel %s = %v, want ErrNotCancellable", tc.name, err)
			}
			after, err := repo.GetJob(context.Background(), id)
			if err != nil {
				t.Fatalf("get: %v", err)
			}
			if after.Status != before.Status || after.CancelRequested != before.CancelRequested {
				t.Fatalf("a refused cancel mutated the row: %+v -> %+v", before, after)
			}
		})
	}
}

// C1: a missing id is JOB_NOT_FOUND, never a silent success.
func TestRequestCancelAnyMissingJobIsNotFound(t *testing.T) {
	repo, _ := newRepository(t)
	if _, err := repo.RequestCancelAny(context.Background(), "job-absent", testNow); !errors.Is(err, jobs.ErrNotFound) {
		t.Fatalf("cancel of a missing job = %v, want ErrNotFound", err)
	}
}

// C2: retrying another actor's failed job requeues it from a clean state.
func TestRetryAnyRequeuesAnotherActorsFailedJob(t *testing.T) {
	repo, _ := newRepository(t)
	job := createJobFor(t, repo, "job-any-retry", "someone-else", 3)
	_, lease, err := repo.Claim(context.Background(), job.ID, "worker-1", testNow, time.Minute)
	if err != nil {
		t.Fatalf("claim: %v", err)
	}
	if err := repo.UpdateProgress(context.Background(), lease, 40, testNow); err != nil {
		t.Fatalf("progress: %v", err)
	}
	if err := repo.Fail(context.Background(), lease, "JOB_HANDLER_FAILED", "boom", testNow); err != nil {
		t.Fatalf("fail: %v", err)
	}

	// The actor-scoped path must still refuse a non-owner.
	if _, err := repo.Retry(context.Background(), job.ID, "operator", testNow.Add(time.Second)); !errors.Is(err, jobs.ErrNotFound) {
		t.Fatalf("Retry by a non-owner = %v, want ErrNotFound (frozen isolation)", err)
	}

	retried, err := repo.RetryAny(context.Background(), job.ID, testNow.Add(time.Second))
	if err != nil {
		t.Fatalf("RetryAny: %v", err)
	}
	if retried.Status != jobs.StatusQueued {
		t.Fatalf("status = %s, want queued", retried.Status)
	}
	if retried.Progress != 0 || retried.CancelRequested || retried.LeaseOwner != "" || retried.LeaseExpiresAt != nil {
		t.Fatalf("retry left execution state behind: %+v", retried)
	}
	if retried.ErrorCode != "" || retried.ErrorMessage != "" || len(retried.Result) != 0 {
		t.Fatalf("retry left result/error behind: %+v", retried)
	}
	if retried.FinishedAt != nil || retried.ResultExpiresAt != nil {
		t.Fatalf("retry left terminal timestamps behind: %+v", retried)
	}
	// The attempt budget is consumed, not reset: extending it would change the
	// frozen Job contract (VP-038 non-goal).
	if retried.Attempt != 1 || retried.MaxAttempts != 3 {
		t.Fatalf("retry must preserve the attempt budget: %+v", retried)
	}
}

// C2: the retryable set is exactly {failed with remaining budget}. An exhausted
// failure is refused and, crucially, is NOT silently re-armed.
func TestRetryAnyRefusesEverythingButFailedWithBudget(t *testing.T) {
	t.Run("exhausted failed job", func(t *testing.T) {
		repo, _ := newRepository(t)
		job := createJobFor(t, repo, "job-exhausted", "someone-else", 1)
		_, lease, err := repo.Claim(context.Background(), job.ID, "w", testNow, time.Minute)
		if err != nil {
			t.Fatalf("claim: %v", err)
		}
		if err := repo.Fail(context.Background(), lease, "JOB_HANDLER_FAILED", "boom", testNow); err != nil {
			t.Fatalf("fail: %v", err)
		}
		if _, err := repo.RetryAny(context.Background(), job.ID, testNow.Add(time.Second)); !errors.Is(err, jobs.ErrNotRetryable) {
			t.Fatalf("retry of an exhausted job = %v, want ErrNotRetryable", err)
		}
		after, err := repo.GetJob(context.Background(), job.ID)
		if err != nil {
			t.Fatalf("get: %v", err)
		}
		if after.Status != jobs.StatusFailed || after.Attempt != 1 || after.MaxAttempts != 1 {
			t.Fatalf("a refused retry must not re-arm the job: %+v", after)
		}
	})

	t.Run("queued running succeeded cancelled", func(t *testing.T) {
		repo, _ := newRepository(t)
		queued := createJobFor(t, repo, "job-q", "someone-else", 3)
		running := createJobFor(t, repo, "job-r", "someone-else", 3)
		if _, _, err := repo.Claim(context.Background(), running.ID, "w", testNow, time.Minute); err != nil {
			t.Fatalf("claim: %v", err)
		}
		succeeded := createJobFor(t, repo, "job-s", "someone-else", 3)
		_, lease, err := repo.Claim(context.Background(), succeeded.ID, "w", testNow, time.Minute)
		if err != nil {
			t.Fatalf("claim: %v", err)
		}
		if _, err := repo.CompleteWithCommit(context.Background(), lease, testNow, time.Hour, successCommit(`{"ok":true}`)); err != nil {
			t.Fatalf("complete: %v", err)
		}
		cancelled := createJobFor(t, repo, "job-c", "someone-else", 3)
		if _, err := repo.RequestCancelAny(context.Background(), cancelled.ID, testNow); err != nil {
			t.Fatalf("cancel: %v", err)
		}

		for _, id := range []string{queued.ID, running.ID, succeeded.ID, cancelled.ID} {
			before, err := repo.GetJob(context.Background(), id)
			if err != nil {
				t.Fatalf("get %s: %v", id, err)
			}
			if _, err := repo.RetryAny(context.Background(), id, testNow.Add(time.Second)); !errors.Is(err, jobs.ErrNotRetryable) {
				t.Fatalf("retry %s = %v, want ErrNotRetryable", id, err)
			}
			after, err := repo.GetJob(context.Background(), id)
			if err != nil {
				t.Fatalf("get %s: %v", id, err)
			}
			if after.Status != before.Status || after.Attempt != before.Attempt || !after.UpdatedAt.Equal(before.UpdatedAt) {
				t.Fatalf("a refused retry mutated %s: %+v -> %+v", id, before, after)
			}
		}
	})

	t.Run("missing job", func(t *testing.T) {
		repo, _ := newRepository(t)
		if _, err := repo.RetryAny(context.Background(), "job-absent", testNow); !errors.Is(err, jobs.ErrNotFound) {
			t.Fatalf("retry of a missing job = %v, want ErrNotFound", err)
		}
	})
}

// C1/C2: the runner-level wrappers must reach the in-flight execution and the
// scan loop exactly like their actor-scoped twins, otherwise a running job
// would keep working after the operator cancelled it.
func TestRunnerCancelAnyStopsAnotherActorsRunningJob(t *testing.T) {
	repo, _ := newRepository(t)
	started := make(chan struct{})
	runner := startRunner(t, repo, func(ctx context.Context, _ jobs.Job, _ jobs.Reporter) (jobs.CommitFunc, error) {
		close(started)
		<-ctx.Done()
		return nil, ctx.Err()
	})
	job, err := runner.Submit(context.Background(), submitInput("job-cancel-any"))
	if err != nil {
		t.Fatal(err)
	}
	// The job belongs to user-1 and the operator is a different principal: the
	// management-scope cancel must still stop it.
	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("handler did not start")
	}
	marked, err := runner.CancelAny(context.Background(), job.ID)
	if err != nil || !marked.CancelRequested {
		t.Fatalf("CancelAny mark = %+v err=%v", marked, err)
	}
	waitForStatus(t, repo, job.ID, jobs.StatusCancelled)
}

func TestRunnerRetryAnyRequeuesAnotherActorsFailedJob(t *testing.T) {
	repo, _ := newRepository(t)
	var calls atomic.Int32
	runner := startRunner(t, repo, func(_ context.Context, _ jobs.Job, _ jobs.Reporter) (jobs.CommitFunc, error) {
		if calls.Add(1) == 1 {
			return nil, errors.New("first attempt failed")
		}
		return successCommit(`{"ok":true}`), nil
	})
	job, err := runner.Submit(context.Background(), submitInput("job-retry-any"))
	if err != nil {
		t.Fatal(err)
	}
	waitForStatus(t, repo, job.ID, jobs.StatusFailed)
	if _, err := runner.RetryAny(context.Background(), job.ID); err != nil {
		t.Fatalf("RetryAny: %v", err)
	}
	succeeded := waitForStatus(t, repo, job.ID, jobs.StatusSucceeded)
	if succeeded.Attempt != 2 {
		t.Fatalf("retried job = %+v, want attempt 2", succeeded)
	}
}
