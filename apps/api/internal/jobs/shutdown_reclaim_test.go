package jobs_test

import (
	"context"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/jobs"
)

// TestShutdownInterruptLeaseReclaim is the VP-021 contract §4 evidence for
// the frozen "interrupt-and-rerun" job shutdown semantics:
//
//  1. shutdown cancels the running job's context (interrupt); no durable
//     terminal transition is written (finish() skips while stopping);
//  2. the interrupted job stays durable 'running' with its lease expiring;
//  3. the next process (a fresh Runner over the same store) reclaims the
//     expired lease with attempt+1 and reruns the handler to success.
func TestShutdownInterruptLeaseReclaim(t *testing.T) {
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

	job, err := runner.Submit(context.Background(), submitInput("shutdown-interrupt"))
	if err != nil {
		t.Fatal(err)
	}
	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("handler did not start")
	}

	// 1. Shutdown: cancel the active job and wait for worker shutdown.
	stopCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := runner.Stop(stopCtx); err != nil {
		t.Fatalf("runner.Stop: %v", err)
	}

	// 2. Interrupted job stays durable 'running' (no terminal write on
	// shutdown); its lease expires after LeaseDuration (~120ms).
	running, err := repo.Get(context.Background(), job.ID)
	if err != nil {
		t.Fatal(err)
	}
	if running.Status != jobs.StatusRunning {
		t.Fatalf("interrupted job status = %s, want running (no durable terminal on shutdown)", running.Status)
	}
	time.Sleep(300 * time.Millisecond) // let the lease expire

	// 3. Fresh process: a new Runner over the same store reclaims the expired
	// lease (attempt+1) and reruns the handler to success.
	resumed := make(chan struct{})
	runner2, err := jobs.NewRunner(repo, runnerOptions())
	if err != nil {
		t.Fatal(err)
	}
	if err := runner2.Register("wallet.reconcile", func(_ context.Context, _ jobs.Job, _ jobs.Reporter) (jobs.CommitFunc, error) {
		close(resumed)
		return successCommit(`{"recovered":true}`), nil
	}); err != nil {
		t.Fatal(err)
	}
	if err := runner2.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = runner2.Stop(ctx)
	})

	select {
	case <-resumed:
	case <-time.After(3 * time.Second):
		t.Fatal("reclaimed job handler did not run")
	}
	// Completion is async after the handler returns; poll the durable state.
	completed := waitForStatus(t, repo, job.ID, jobs.StatusSucceeded)
	if completed.Attempt != 2 {
		t.Fatalf("reclaimed job attempt = %d, want 2", completed.Attempt)
	}
}