package jobs

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)

// W9 A-005 R-F-003: a panicking handler must degrade to a durable
// JOB_HANDLER_FAILED job (F-007), never crash the process.
func TestPanickingHandlerRecordsDurableFailure(t *testing.T) {
	st, err := testsupport.OpenStore(":memory:", "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	repo := NewRepository(st)

	options := DefaultRunnerOptions()
	options.ScanInterval = 20 * time.Millisecond
	options.HeartbeatInterval = 200 * time.Millisecond
	options.LeaseDuration = time.Second
	runner, err := NewRunner(repo, options)
	if err != nil {
		t.Fatal(err)
	}
	if err := runner.Register("panic.kind", func(context.Context, Job, Reporter) (CommitFunc, error) {
		panic("injected handler panic")
	}); err != nil {
		t.Fatal(err)
	}
	if err := runner.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = runner.Stop(context.Background()) })

	now := time.Now().UTC()
	job, err := runner.Submit(context.Background(), CreateInput{
		ID: "job-panic", Kind: "panic.kind", Payload: []byte(`{}`),
		ActorID: "user-1", CorrelationID: "corr-1", Now: now,
	})
	if err != nil {
		t.Fatal(err)
	}

	deadline := time.Now().Add(5 * time.Second)
	for {
		current, err := repo.Get(context.Background(), job.ID)
		if err != nil {
			t.Fatal(err)
		}
		if current.Status == StatusFailed {
			if current.ErrorCode != "JOB_HANDLER_FAILED" {
				t.Fatalf("error code = %q, want JOB_HANDLER_FAILED", current.ErrorCode)
			}
			if !strings.Contains(current.ErrorMessage, "panicked") {
				t.Fatalf("error message = %q, want panic detail", current.ErrorMessage)
			}
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("job did not fail; status = %q", current.Status)
		}
		time.Sleep(10 * time.Millisecond)
	}
}
