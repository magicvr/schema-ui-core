package scheduledtasks

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/scheduledtasks/store"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)

// W9 A-005 R-F-003: a panicking task handler is recorded as a failed run
// (F-007) instead of crashing the scheduler loop or the process.
func TestSchedulerExecuteRecoversPanic(t *testing.T) {
	hash, err := auth.HashPassword("pw", 4)
	if err != nil {
		t.Fatal(err)
	}
	st, err := testsupport.OpenStore(filepath.Join(t.TempDir(), "t.db"), "admin", hash, true)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	repo := store.NewRepository(st)

	now := time.Date(2026, 8, 14, 10, 0, 0, 0, time.UTC)
	task := store.Task{ID: "task-panic", Key: "panic", Cron: "* * * * *", Name: "Panic", Enabled: true, Handler: "panic.key", CreatedAt: now, UpdatedAt: now}
	if err := repo.CreateTask(task); err != nil {
		t.Fatal(err)
	}

	s := NewScheduler(repo)
	s.handlers["panic.key"] = func(context.Context, store.Task, time.Time) error {
		panic("injected task panic")
	}

	if err := s.Execute(task, now); err != nil {
		t.Fatalf("Execute returned err %v, want nil (panic is recorded as a failed run)", err)
	}
	runs, total, err := repo.ListAllRuns(store.ListFilter{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatal(err)
	}
	if total != 1 {
		t.Fatalf("runs = %d, want 1", total)
	}
	if runs[0].Status != "failed" || !strings.Contains(runs[0].Detail, "panicked") {
		t.Fatalf("run = %+v, want failed with panic detail", runs[0])
	}
	// The scheduler keeps working after the recovered panic. The panic task is
	// still enabled with an every-minute cron, so the next tick records its
	// second failed run AND executes the healthy task.
	healthy := store.Task{ID: "task-healthy", Key: "healthy", Cron: "* * * * *", Name: "Healthy", Enabled: true, Handler: "system.noop", CreatedAt: now, UpdatedAt: now}
	if err := repo.CreateTask(healthy); err != nil {
		t.Fatal(err)
	}
	s.tick(now.Add(1 * time.Minute))
	runs, total, _ = repo.ListAllRuns(store.ListFilter{Page: 1, PageSize: 10})
	if total != 3 {
		t.Fatalf("runs after next tick = %d, want 3 (panic retry + healthy)", total)
	}
	healthyRan := false
	for _, run := range runs {
		if run.TaskID == "task-healthy" && run.Status == "ran" {
			healthyRan = true
		}
	}
	if !healthyRan {
		t.Fatalf("healthy task did not run after the recovered panic: %+v", runs)
	}
}
