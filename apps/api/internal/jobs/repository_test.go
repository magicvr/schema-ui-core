package jobs_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/jobs"
	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)

var testNow = time.Unix(1_700_000_000, 0).UTC()

func newRepository(t *testing.T) (*jobs.Repository, *store.Store) {
	t.Helper()
	st, err := testsupport.OpenStore(":memory:", "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return jobs.NewRepository(st), st
}

func createJob(t *testing.T, repo *jobs.Repository, id string, maxAttempts int) *jobs.Job {
	t.Helper()
	job, err := repo.Create(context.Background(), jobs.CreateInput{
		ID: id, Kind: "wallet.reconcile", Payload: json.RawMessage(`{"accountId":"a1"}`),
		ActorID: "user-1", CorrelationID: "corr-1", MaxAttempts: maxAttempts, Now: testNow,
	})
	if err != nil {
		t.Fatal(err)
	}
	return job
}

func TestRepositoryLifecycleRetryCompleteAndExpire(t *testing.T) {
	repo, _ := newRepository(t)
	job := createJob(t, repo, "job-1", 3)
	if job.Status != jobs.StatusQueued || job.Attempt != 0 || job.MaxAttempts != 3 {
		t.Fatalf("created job = %+v", job)
	}

	running, lease1, err := repo.Claim(context.Background(), job.ID, "worker-1", testNow, 30*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if running.Status != jobs.StatusRunning || running.Attempt != 1 || lease1.Version != 1 {
		t.Fatalf("claimed job = %+v lease=%+v", running, lease1)
	}
	if err := repo.UpdateProgress(context.Background(), lease1, 40, testNow.Add(time.Second)); err != nil {
		t.Fatal(err)
	}
	if err := repo.UpdateProgress(context.Background(), lease1, 20, testNow.Add(2*time.Second)); !errors.Is(err, jobs.ErrTransition) {
		t.Fatalf("decreasing progress error = %v, want ErrTransition", err)
	}
	if err := repo.Heartbeat(context.Background(), lease1, testNow.Add(2*time.Second), 30*time.Second); err != nil {
		t.Fatal(err)
	}
	if err := repo.Fail(context.Background(), lease1, "JOB_HANDLER_FAILED", "boom", testNow.Add(3*time.Second)); err != nil {
		t.Fatal(err)
	}
	failed, err := repo.Get(context.Background(), job.ID)
	if err != nil || failed.Status != jobs.StatusFailed || failed.ErrorCode != "JOB_HANDLER_FAILED" {
		t.Fatalf("failed job = %+v err=%v", failed, err)
	}

	queued, err := repo.Retry(context.Background(), job.ID, "user-1", testNow.Add(4*time.Second))
	if err != nil || queued.Status != jobs.StatusQueued || queued.Attempt != 1 || queued.Progress != 0 {
		t.Fatalf("retried job = %+v err=%v", queued, err)
	}
	_, lease2, err := repo.Claim(context.Background(), job.ID, "worker-2", testNow.Add(5*time.Second), 30*time.Second)
	if err != nil || lease2.Version != 2 {
		t.Fatalf("second claim lease=%+v err=%v", lease2, err)
	}
	completed, err := repo.CompleteWithCommit(context.Background(), lease2, testNow.Add(6*time.Second), 24*time.Hour,
		func(*sql.Tx) (json.RawMessage, error) { return json.RawMessage(`{"result":"consistent"}`), nil })
	if err != nil || completed.Status != jobs.StatusSucceeded || completed.Progress != 100 || completed.Attempt != 2 {
		t.Fatalf("completed job = %+v err=%v", completed, err)
	}
	notDue, err := repo.ExpireIfDue(context.Background(), job.ID, testNow.Add(12*time.Hour))
	if err != nil || notDue.Status != jobs.StatusSucceeded {
		t.Fatalf("not-due expiry = %+v err=%v", notDue, err)
	}
	expired, err := repo.ExpireIfDue(context.Background(), job.ID, testNow.Add(25*time.Hour))
	if err != nil || expired.Status != jobs.StatusExpired || expired.Result != nil {
		t.Fatalf("expired job = %+v err=%v", expired, err)
	}
}

func TestCancellationPathsAndActorIsolation(t *testing.T) {
	repo, _ := newRepository(t)
	queued := createJob(t, repo, "job-queued", 3)
	if _, err := repo.RequestCancel(context.Background(), queued.ID, "other-user", testNow); !errors.Is(err, jobs.ErrNotFound) {
		t.Fatalf("cross-actor cancel error = %v", err)
	}
	cancelled, err := repo.RequestCancel(context.Background(), queued.ID, "user-1", testNow)
	if err != nil || cancelled.Status != jobs.StatusCancelled {
		t.Fatalf("queued cancel = %+v err=%v", cancelled, err)
	}

	running := createJob(t, repo, "job-running", 3)
	_, lease, err := repo.Claim(context.Background(), running.ID, "worker", testNow, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	marked, err := repo.RequestCancel(context.Background(), running.ID, "user-1", testNow)
	if err != nil || !marked.CancelRequested || marked.Status != jobs.StatusRunning {
		t.Fatalf("running cancel mark = %+v err=%v", marked, err)
	}
	requested, err := repo.IsCancelRequested(context.Background(), lease)
	if err != nil || !requested {
		t.Fatalf("cancel requested = %v err=%v", requested, err)
	}
	if err := repo.Fail(context.Background(), lease, "JOB_HANDLER_FAILED", "boom", testNow); !errors.Is(err, jobs.ErrTransition) {
		t.Fatalf("fail after cancel request error = %v, want ErrTransition", err)
	}
	if err := repo.FinalizeCancel(context.Background(), lease, testNow.Add(time.Second)); err != nil {
		t.Fatal(err)
	}
	final, _ := repo.Get(context.Background(), running.ID)
	if final.Status != jobs.StatusCancelled {
		t.Fatalf("finalized cancel = %+v", final)
	}

	recovering := createJob(t, repo, "job-recover-cancel", 3)
	_, _, err = repo.Claim(context.Background(), recovering.ID, "dead-worker", testNow, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := repo.RequestCancel(context.Background(), recovering.ID, "user-1", testNow); err != nil {
		t.Fatal(err)
	}
	count, err := repo.RecoverCancelledDue(context.Background(), testNow.Add(2*time.Second))
	if err != nil || count != 1 {
		t.Fatalf("recover cancelled count=%d err=%v", count, err)
	}
	recovered, _ := repo.Get(context.Background(), recovering.ID)
	if recovered.Status != jobs.StatusCancelled {
		t.Fatalf("recovered cancel = %+v", recovered)
	}
}

func TestFencingAndCompleteWithCommitRollback(t *testing.T) {
	repo, _ := newRepository(t)
	job := createJob(t, repo, "job-fencing", 3)
	_, stale, err := repo.Claim(context.Background(), job.ID, "worker-1", testNow, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	_, current, err := repo.Claim(context.Background(), job.ID, "worker-2", testNow.Add(2*time.Second), time.Second)
	if err != nil {
		t.Fatal(err)
	}
	called := false
	_, err = repo.CompleteWithCommit(context.Background(), stale, testNow.Add(2*time.Second), time.Hour, func(*sql.Tx) (json.RawMessage, error) {
		called = true
		return json.RawMessage(`{}`), nil
	})
	if !errors.Is(err, jobs.ErrLeaseLost) || called {
		t.Fatalf("stale complete err=%v called=%v", err, called)
	}

	rollbackErr := errors.New("rollback consumer")
	_, err = repo.CompleteWithCommit(context.Background(), current, testNow.Add(2*time.Second), time.Hour, func(tx *sql.Tx) (json.RawMessage, error) {
		_, insertErr := tx.Exec(`INSERT INTO wallet_reconciliation_runs (id, account_id, result, mismatch_count, details, actor_id, created_at) VALUES ('job-fencing', NULL, 'consistent', 0, '{}', 'user-1', ?)`, testNow.Unix())
		if insertErr != nil {
			return nil, insertErr
		}
		return nil, rollbackErr
	})
	if !errors.Is(err, rollbackErr) {
		t.Fatalf("rollback complete error = %v", err)
	}
	completed, err := repo.CompleteWithCommit(context.Background(), current, testNow.Add(3*time.Second), time.Hour, func(tx *sql.Tx) (json.RawMessage, error) {
		_, insertErr := tx.Exec(`INSERT INTO wallet_reconciliation_runs (id, account_id, result, mismatch_count, details, actor_id, created_at) VALUES ('job-fencing', NULL, 'consistent', 0, '{}', 'user-1', ?)`, testNow.Unix())
		return json.RawMessage(`{"id":"job-fencing"}`), insertErr
	})
	if err != nil || completed.Status != jobs.StatusSucceeded {
		t.Fatalf("atomic complete = %+v err=%v", completed, err)
	}
}

func TestExhaustionAndRunnableSelection(t *testing.T) {
	repo, _ := newRepository(t)
	exhausted := createJob(t, repo, "job-exhausted", 1)
	_, _, err := repo.Claim(context.Background(), exhausted.ID, "worker", testNow, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	queued := createJob(t, repo, "job-runnable", 3)
	runnable, err := repo.ListRunnable(context.Background(), testNow.Add(2*time.Second), 10)
	if err != nil || len(runnable) != 1 || runnable[0].ID != queued.ID {
		t.Fatalf("runnable = %+v err=%v", runnable, err)
	}
	count, err := repo.ExhaustExpired(context.Background(), testNow.Add(2*time.Second))
	if err != nil || count != 1 {
		t.Fatalf("exhaust count=%d err=%v", count, err)
	}
	failed, _ := repo.Get(context.Background(), exhausted.ID)
	if failed.Status != jobs.StatusFailed || failed.ErrorCode != "JOB_ATTEMPTS_EXHAUSTED" {
		t.Fatalf("exhausted job = %+v", failed)
	}
	if _, err := repo.Retry(context.Background(), exhausted.ID, "user-1", testNow.Add(3*time.Second)); !errors.Is(err, jobs.ErrNotRetryable) {
		t.Fatalf("exhausted retry error = %v", err)
	}
}

func TestCreateValidationAndGetForActor(t *testing.T) {
	repo, _ := newRepository(t)
	if _, err := repo.Create(context.Background(), jobs.CreateInput{}); !errors.Is(err, jobs.ErrInvalid) {
		t.Fatalf("invalid create error = %v", err)
	}
	job := createJob(t, repo, "job-actor", 0)
	if _, err := repo.GetForActor(context.Background(), job.ID, "wallet.reconcile", "other-user"); !errors.Is(err, jobs.ErrNotFound) {
		t.Fatalf("cross-actor get error = %v", err)
	}
	found, err := repo.GetForActor(context.Background(), job.ID, "wallet.reconcile", "user-1")
	if err != nil || found.ID != job.ID || found.MaxAttempts != jobs.DefaultMaxAttempts {
		t.Fatalf("actor get = %+v err=%v", found, err)
	}
}
