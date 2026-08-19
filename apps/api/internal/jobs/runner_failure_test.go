package jobs

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)

func TestAbortLeasePersistsFailureAndNotifiesTerminalHook(t *testing.T) {
	st, err := testsupport.OpenStore(":memory:", "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	repo := NewRepository(st)
	now := time.Now().UTC()
	job, err := repo.Create(context.Background(), CreateInput{
		ID: "job-abort-lease", Kind: "wallet.reconcile", Payload: []byte(`{}`),
		ActorID: "user-1", CorrelationID: "corr-1", Now: now,
	})
	if err != nil {
		t.Fatal(err)
	}
	_, lease, err := repo.Claim(context.Background(), job.ID, "worker-1", now, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	terminal := make(chan Job, 1)
	runner := &Runner{
		repo:     repo,
		options:  RunnerOptions{Now: time.Now},
		handlers: map[string]registration{"wallet.reconcile": {terminal: func(job Job) { terminal <- job }}},
	}
	runner.abortLease(lease, errors.New("injected heartbeat failure"))

	failed, err := repo.Get(context.Background(), job.ID)
	if err != nil {
		t.Fatal(err)
	}
	if failed.Status != StatusFailed || failed.ErrorCode != "JOB_HANDLER_FAILED" || failed.LeaseOwner != "" {
		t.Fatalf("failed lease = %+v", failed)
	}
	select {
	case notified := <-terminal:
		if notified.Status != StatusFailed {
			t.Fatalf("terminal hook job = %+v", notified)
		}
	case <-time.After(time.Second):
		t.Fatal("terminal hook was not called")
	}
}
