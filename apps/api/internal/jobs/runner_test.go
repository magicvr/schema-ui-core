package jobs_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	"sync/atomic"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/jobs"
)

func runnerOptions() jobs.RunnerOptions {
	return jobs.RunnerOptions{
		LeaseDuration: 120 * time.Millisecond, HeartbeatInterval: 20 * time.Millisecond,
		ScanInterval: 10 * time.Millisecond, ResultTTL: 200 * time.Millisecond,
		BatchSize: 10, Now: time.Now,
	}
}

func startRunner(t *testing.T, repo *jobs.Repository, handler jobs.Handler) *jobs.Runner {
	t.Helper()
	runner, err := jobs.NewRunner(repo, runnerOptions())
	if err != nil {
		t.Fatal(err)
	}
	if err := runner.Register("wallet.reconcile", handler); err != nil {
		t.Fatal(err)
	}
	if err := runner.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = runner.Stop(ctx)
	})
	return runner
}

func waitForStatus(t *testing.T, repo *jobs.Repository, id string, statuses ...jobs.Status) *jobs.Job {
	t.Helper()
	wanted := map[jobs.Status]bool{}
	for _, status := range statuses {
		wanted[status] = true
	}
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		job, err := repo.Get(context.Background(), id)
		if err == nil && wanted[job.Status] {
			return job
		}
		time.Sleep(10 * time.Millisecond)
	}
	job, err := repo.Get(context.Background(), id)
	t.Fatalf("job %s did not reach %v: job=%+v err=%v", id, statuses, job, err)
	return nil
}

func submitInput(id string) jobs.CreateInput {
	return jobs.CreateInput{
		ID: id, Kind: "wallet.reconcile", Payload: json.RawMessage(`{}`),
		ActorID: "user-1", CorrelationID: "corr-1", Now: time.Now().UTC(),
	}
}

func successCommit(result string) jobs.CommitFunc {
	return func(kernel.Tx) (json.RawMessage, error) {
		return json.RawMessage(result), nil
	}
}

func TestRunnerScanNotifiesRecoveredTerminalTransitions(t *testing.T) {
	repo, _ := newRepository(t)
	now := time.Now().UTC()
	exhausted, err := repo.Create(context.Background(), jobs.CreateInput{
		ID: "job-exhaust-hook", Kind: "wallet.reconcile", Payload: json.RawMessage(`{}`),
		ActorID: "user-1", CorrelationID: "corr-exhaust", MaxAttempts: 1, Now: now.Add(-time.Second),
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := repo.Claim(context.Background(), exhausted.ID, "dead-owner", now.Add(-time.Second), 100*time.Millisecond); err != nil {
		t.Fatal(err)
	}
	cancelled, err := repo.Create(context.Background(), jobs.CreateInput{
		ID: "job-cancel-hook", Kind: "wallet.reconcile", Payload: json.RawMessage(`{}`),
		ActorID: "user-1", CorrelationID: "corr-cancel", Now: now.Add(-time.Second),
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := repo.Claim(context.Background(), cancelled.ID, "dead-owner", now.Add(-time.Second), 100*time.Millisecond); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.RequestCancel(context.Background(), cancelled.ID, "user-1", now.Add(-900*time.Millisecond)); err != nil {
		t.Fatal(err)
	}

	options := runnerOptions()
	options.Now = func() time.Time { return now }
	runner, err := jobs.NewRunner(repo, options)
	if err != nil {
		t.Fatal(err)
	}
	terminal := make(chan jobs.Job, 2)
	if err := runner.RegisterWithTerminalHook("wallet.reconcile", func(context.Context, jobs.Job, jobs.Reporter) (jobs.CommitFunc, error) {
		return successCommit(`{"ok":true}`), nil
	}, func(job jobs.Job) { terminal <- job }); err != nil {
		t.Fatal(err)
	}
	if err := runner.ScanOnce(context.Background()); err != nil {
		t.Fatal(err)
	}
	seen := map[string]jobs.Status{}
	for range 2 {
		select {
		case job := <-terminal:
			seen[job.ID] = job.Status
		case <-time.After(time.Second):
			t.Fatal("terminal hook was not called")
		}
	}
	if seen[exhausted.ID] != jobs.StatusFailed || seen[cancelled.ID] != jobs.StatusCancelled {
		t.Fatalf("terminal hook statuses = %v", seen)
	}
}

func TestRunnerStartupReclaimsExpiredJob(t *testing.T) {
	repo, _ := newRepository(t)
	past := time.Now().Add(-time.Second)
	job, err := repo.Create(context.Background(), jobs.CreateInput{
		ID: "job-restart", Kind: "wallet.reconcile", Payload: json.RawMessage(`{}`),
		ActorID: "user-1", CorrelationID: "corr-1", Now: past,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := repo.Claim(context.Background(), job.ID, "dead-worker", past, 50*time.Millisecond); err != nil {
		t.Fatal(err)
	}
	startRunner(t, repo, func(_ context.Context, _ jobs.Job, reporter jobs.Reporter) (jobs.CommitFunc, error) {
		if err := reporter.Progress(75); err != nil {
			return nil, err
		}
		return successCommit(`{"recovered":true}`), nil
	})
	completed := waitForStatus(t, repo, job.ID, jobs.StatusSucceeded)
	if completed.Attempt != 2 || completed.Progress != 100 {
		t.Fatalf("recovered job = %+v", completed)
	}
}

func TestRunnerCancelRunningJob(t *testing.T) {
	repo, _ := newRepository(t)
	started := make(chan struct{})
	runner := startRunner(t, repo, func(ctx context.Context, _ jobs.Job, _ jobs.Reporter) (jobs.CommitFunc, error) {
		close(started)
		<-ctx.Done()
		return nil, ctx.Err()
	})
	job, err := runner.Submit(context.Background(), submitInput("job-cancel"))
	if err != nil {
		t.Fatal(err)
	}
	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("handler did not start")
	}
	marked, err := runner.Cancel(context.Background(), job.ID, "user-1")
	if err != nil || !marked.CancelRequested {
		t.Fatalf("cancel mark = %+v err=%v", marked, err)
	}
	waitForStatus(t, repo, job.ID, jobs.StatusCancelled)
}

func TestRunnerFailureRetryAndResultExpiry(t *testing.T) {
	repo, _ := newRepository(t)
	var calls atomic.Int32
	runner := startRunner(t, repo, func(_ context.Context, _ jobs.Job, _ jobs.Reporter) (jobs.CommitFunc, error) {
		if calls.Add(1) == 1 {
			return nil, errors.New("first attempt failed")
		}
		return successCommit(`{"ok":true}`), nil
	})
	job, err := runner.Submit(context.Background(), submitInput("job-retry"))
	if err != nil {
		t.Fatal(err)
	}
	failed := waitForStatus(t, repo, job.ID, jobs.StatusFailed)
	if failed.ErrorCode != "JOB_HANDLER_FAILED" || failed.Attempt != 1 {
		t.Fatalf("failed job = %+v", failed)
	}
	if _, err := runner.Retry(context.Background(), job.ID, "user-1"); err != nil {
		t.Fatal(err)
	}
	succeeded := waitForStatus(t, repo, job.ID, jobs.StatusSucceeded)
	if succeeded.Attempt != 2 {
		t.Fatalf("retried job = %+v", succeeded)
	}
	waitForStatus(t, repo, job.ID, jobs.StatusExpired)
}

func TestHandlerCancelledWithoutRequestFails(t *testing.T) {
	repo, _ := newRepository(t)
	runner := startRunner(t, repo, func(context.Context, jobs.Job, jobs.Reporter) (jobs.CommitFunc, error) {
		return nil, context.Canceled
	})
	job, err := runner.Submit(context.Background(), submitInput("job-internal-cancel"))
	if err != nil {
		t.Fatal(err)
	}
	failed := waitForStatus(t, repo, job.ID, jobs.StatusFailed)
	if failed.ErrorCode != "JOB_HANDLER_FAILED" {
		t.Fatalf("internally cancelled job = %+v", failed)
	}
}

func TestTwoRunnersHeartbeatPreventsDuplicateClaim(t *testing.T) {
	repo, _ := newRepository(t)
	options := runnerOptions()
	options.LeaseDuration = 60 * time.Millisecond
	options.HeartbeatInterval = 10 * time.Millisecond
	var calls atomic.Int32
	handler := func(ctx context.Context, _ jobs.Job, _ jobs.Reporter) (jobs.CommitFunc, error) {
		calls.Add(1)
		select {
		case <-time.After(150 * time.Millisecond):
			return successCommit(`{"ok":true}`), nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	newRunner := func() *jobs.Runner {
		runner, err := jobs.NewRunner(repo, options)
		if err != nil {
			t.Fatal(err)
		}
		if err := runner.Register("wallet.reconcile", handler); err != nil {
			t.Fatal(err)
		}
		if err := runner.Start(); err != nil {
			t.Fatal(err)
		}
		return runner
	}
	one, two := newRunner(), newRunner()
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = one.Stop(ctx)
		_ = two.Stop(ctx)
	})
	job, err := one.Submit(context.Background(), submitInput("job-heartbeat"))
	if err != nil {
		t.Fatal(err)
	}
	completed := waitForStatus(t, repo, job.ID, jobs.StatusSucceeded)
	if completed.Attempt != 1 || calls.Load() != 1 {
		t.Fatalf("heartbeat job attempt=%d calls=%d", completed.Attempt, calls.Load())
	}
}

func TestRunnerStopLeavesJobRecoverable(t *testing.T) {
	repo, _ := newRepository(t)
	started := make(chan struct{})
	runner, err := jobs.NewRunner(repo, runnerOptions())
	if err != nil {
		t.Fatal(err)
	}
	if err := runner.Register("wallet.reconcile", func(ctx context.Context, _ jobs.Job, _ jobs.Reporter) (jobs.CommitFunc, error) {
		close(started)
		<-ctx.Done()
		return nil, ctx.Err()
	}); err != nil {
		t.Fatal(err)
	}
	if err := runner.Start(); err != nil {
		t.Fatal(err)
	}
	job, err := runner.Submit(context.Background(), submitInput("job-stop"))
	if err != nil {
		t.Fatal(err)
	}
	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("handler did not start")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := runner.Stop(ctx); err != nil {
		t.Fatal(err)
	}
	stopped, err := repo.Get(context.Background(), job.ID)
	if err != nil || stopped.Status != jobs.StatusRunning {
		t.Fatalf("stopped job = %+v err=%v", stopped, err)
	}
}
