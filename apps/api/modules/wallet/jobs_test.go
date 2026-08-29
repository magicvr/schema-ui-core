package wallet

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/jobs"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/modules/operationlog"
	walletstore "github.com/magicvr/schema-ui-core/apps/api/modules/wallet/store"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)

type operationRecorder struct {
	mu         sync.Mutex
	operations []operationlog.Operation
}

func (r *operationRecorder) RecordOperation(operation operationlog.Operation) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.operations = append(r.operations, operation)
	return nil
}

// RecordOperationTx implements operationlog.TransactionalRecorder — the
// success audit is now written INSIDE the job transaction (GOAL-037 / F-008
// root-cause fix); the recorder collects it like any other event.
func (r *operationRecorder) RecordOperationTx(_ kernel.Tx, operation operationlog.Operation) error {
	return r.RecordOperation(operation)
}

func (r *operationRecorder) events() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	events := make([]string, 0, len(r.operations))
	for _, operation := range r.operations {
		events = append(events, operation.Event)
	}
	return events
}

func waitOperationEvents(t *testing.T, recorder *operationRecorder, count int) []string {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if events := recorder.events(); len(events) >= count {
			return events
		}
		time.Sleep(5 * time.Millisecond)
	}
	return recorder.events()
}

func newWalletJobHarness(t *testing.T, start bool) (*JobService, *jobs.Runner, *jobs.Repository, *walletstore.Repository, *operationRecorder) {
	t.Helper()
	st, err := testsupport.OpenStore(":memory:", "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	walletRepository := walletstore.NewRepository(st)
	jobRepository := jobs.NewRepository(st)
	options := jobs.DefaultRunnerOptions()
	options.LeaseDuration = 2 * time.Second
	options.HeartbeatInterval = 200 * time.Millisecond
	options.ScanInterval = 10 * time.Millisecond
	runner, err := jobs.NewRunner(jobRepository, options)
	if err != nil {
		t.Fatal(err)
	}
	recorder := &operationRecorder{}
	service, err := NewJobService(NewService(walletRepository), jobRepository, runner, recorder)
	if err != nil {
		t.Fatal(err)
	}
	if start {
		if err := runner.Start(); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			_ = runner.Stop(ctx)
		})
	}
	return service, runner, jobRepository, walletRepository, recorder
}

func waitWalletJob(t *testing.T, service *JobService, id, actorID string, status jobs.Status) *jobs.Job {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		job, err := service.Get(context.Background(), id, actorID)
		if err == nil && job.Status == status {
			return job
		}
		time.Sleep(10 * time.Millisecond)
	}
	job, err := service.Get(context.Background(), id, actorID)
	t.Fatalf("job status = %v err=%v, want %s", job, err, status)
	return nil
}

func TestWalletJobServiceCompletesAtomicallyAndAudits(t *testing.T) {
	service, _, _, walletRepository, recorder := newWalletJobHarness(t, true)
	now := time.Now().UTC()
	if err := walletRepository.CreateAccount(walletstore.Account{
		ID: "acct-1", OwnerType: walletstore.OwnerUser, OwnerID: "user-1",
		Currency: walletstore.DefaultCurrency, Status: walletstore.StatusActive,
		CreatedAt: now, UpdatedAt: now,
	}); err != nil {
		t.Fatal(err)
	}
	job, err := service.SubmitReconcile(context.Background(), "acct-1", account.User{ID: "user-1", Name: "User One"}, "corr-job-1")
	if err != nil {
		t.Fatal(err)
	}
	completed := waitWalletJob(t, service, job.ID, "user-1", jobs.StatusSucceeded)
	if completed.Progress != 100 || len(completed.Result) == 0 {
		t.Fatalf("completed job = %+v", completed)
	}
	var result map[string]any
	if err := json.Unmarshal(completed.Result, &result); err != nil {
		t.Fatal(err)
	}
	if result["id"] != job.ID || result["result"] != walletstore.ResultConsistent {
		t.Fatalf("result = %#v", result)
	}
	runs, total, err := walletRepository.ListReconcileRuns(1, 20)
	if err != nil || total != 1 || len(runs) != 1 || runs[0].ID != job.ID {
		t.Fatalf("runs=%+v total=%d err=%v", runs, total, err)
	}
	if _, err := service.Get(context.Background(), job.ID, "other-user"); !errors.Is(err, jobs.ErrNotFound) {
		t.Fatalf("cross-actor get error = %v, want ErrNotFound", err)
	}
	events := waitOperationEvents(t, recorder, 2)
	if len(events) != 2 || events[0] != operationlog.EventWalletReconcileQueued || events[1] != operationlog.EventWalletReconcile {
		t.Fatalf("events = %v", events)
	}
	service.now = func() time.Time { return completed.FinishedAt.Add(25 * time.Hour) }
	expired, err := service.Get(context.Background(), job.ID, "user-1")
	if err != nil || expired.Status != jobs.StatusExpired || len(expired.Result) != 0 {
		t.Fatalf("expired=%+v err=%v", expired, err)
	}
}

func TestWalletJobServiceCancelsQueuedWithoutBusinessRun(t *testing.T) {
	service, _, _, walletRepository, recorder := newWalletJobHarness(t, false)
	job, err := service.SubmitReconcile(context.Background(), "", account.User{ID: "user-1", Name: "User One"}, "corr-cancel-1")
	if err != nil {
		t.Fatal(err)
	}
	cancelled, err := service.Cancel(context.Background(), job.ID, "user-1")
	if err != nil || cancelled.Status != jobs.StatusCancelled {
		t.Fatalf("cancelled=%+v err=%v", cancelled, err)
	}
	_, total, err := walletRepository.ListReconcileRuns(1, 20)
	if err != nil || total != 0 {
		t.Fatalf("reconcile runs total=%d err=%v", total, err)
	}
	if _, err := service.Retry(context.Background(), job.ID, "user-1"); !errors.Is(err, jobs.ErrNotRetryable) {
		t.Fatalf("retry cancelled error = %v, want ErrNotRetryable", err)
	}
	events := waitOperationEvents(t, recorder, 2)
	if len(events) != 2 || events[1] != operationlog.EventWalletReconcileCancelled {
		t.Fatalf("events = %v", events)
	}
}

func TestWalletJobServiceRollsBackConsumerFailureAndRetries(t *testing.T) {
	service, _, _, walletRepository, recorder := newWalletJobHarness(t, true)
	job, err := service.SubmitReconcile(context.Background(), "missing-account", account.User{ID: "user-1", Name: "User One"}, "corr-fail-1")
	if err != nil {
		t.Fatal(err)
	}
	failed := waitWalletJob(t, service, job.ID, "user-1", jobs.StatusFailed)
	if failed.ErrorCode != "JOB_HANDLER_FAILED" || failed.Attempt != 1 {
		t.Fatalf("failed job = %+v", failed)
	}
	_, total, err := walletRepository.ListReconcileRuns(1, 20)
	if err != nil || total != 0 {
		t.Fatalf("consumer rollback left runs: total=%d err=%v", total, err)
	}
	retried, err := service.Retry(context.Background(), job.ID, "user-1")
	if err != nil || retried.Status != jobs.StatusQueued || retried.Attempt != 1 {
		t.Fatalf("retried=%+v err=%v", retried, err)
	}
	failed = waitWalletJob(t, service, job.ID, "user-1", jobs.StatusFailed)
	if failed.Attempt != 2 {
		t.Fatalf("retry attempt = %d, want 2", failed.Attempt)
	}
	events := waitOperationEvents(t, recorder, 3)
	if len(events) < 3 || events[0] != operationlog.EventWalletReconcileQueued || events[1] != operationlog.EventWalletReconcileFailed || events[2] != operationlog.EventWalletReconcileFailed {
		t.Fatalf("events = %v", events)
	}
}
