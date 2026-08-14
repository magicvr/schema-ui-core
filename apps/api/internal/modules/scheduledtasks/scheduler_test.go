package scheduledtasks

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/scheduledtasks/store"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)

// The scheduler tick executes due enabled tasks and records run rows; manual
// Execute records immediately; disabled tasks are skipped.
func TestSchedulerExecutesDueTasks(t *testing.T) {
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
	due := store.Task{ID: "task-due", Key: "due", Cron: "* * * * *", Name: "Due", Enabled: true, Handler: "system.noop", CreatedAt: now, UpdatedAt: now}
	if err := repo.CreateTask(due); err != nil {
		t.Fatal(err)
	}
	later := store.Task{ID: "task-later", Key: "later", Cron: "0 0 1 1 *", Name: "Later", Enabled: true, Handler: "system.noop", CreatedAt: now, UpdatedAt: now}
	if err := repo.CreateTask(later); err != nil {
		t.Fatal(err)
	}

	s := NewScheduler(repo)
	// A tick at 10:02 executes the due task (current minute slot); a second tick
	// in the same slot is deduplicated.
	s.tick(now.Add(2 * time.Minute))
	s.tick(now.Add(2 * time.Minute).Add(15 * time.Second))
	runs, total, err := repo.ListAllRuns(store.ListFilter{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatal(err)
	}
	if total != 1 || runs[0].TaskID != "task-due" {
		t.Fatalf("runs = %d %+v, want 1 for task-due", total, runs)
	}
	if runs[0].Status != "ran" || runs[0].Detail != "" {
		t.Fatalf("run status/detail = %q %q", runs[0].Status, runs[0].Detail)
	}
	// A later slot (10:03) executes again.
	s.tick(now.Add(3 * time.Minute))
	runs, total, _ = repo.ListAllRuns(store.ListFilter{Page: 1, PageSize: 10})
	if total != 2 {
		t.Fatalf("runs after next slot = %d, want 2", total)
	}
}

// Manual Execute records a run row immediately regardless of the cron window.
func TestSchedulerManualExecute(t *testing.T) {
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
	task := store.Task{ID: "task-manual", Key: "manual", Cron: "0 0 1 1 *", Name: "Manual", Enabled: false, Handler: "system.noop", CreatedAt: now, UpdatedAt: now}
	if err := repo.CreateTask(task); err != nil {
		t.Fatal(err)
	}
	s := NewScheduler(repo)
	if err := s.Execute(task, now); err != nil {
		t.Fatal(err)
	}
	runs, total, _ := repo.ListTaskRuns("task-manual", store.ListFilter{Page: 1, PageSize: 10})
	if total != 1 || runs[0].Status != "ran" {
		t.Fatalf("runs = %d %+v", total, runs)
	}
}


// A never-matching expression records one unschedulable failed run per day
// (A-003 F-002), not a flood.
func TestSchedulerRecordsUnschedule(t *testing.T) {
	hash, err := auth.HashPassword("pw", 4)
	if err != nil {
		t.Fatal(err)
	}
	st2, err := testsupport.OpenStore(filepath.Join(t.TempDir(), "t.db"), "admin", hash, true)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st2.Close() })
	repo := store.NewRepository(st2)
	now := time.Date(2026, 8, 14, 10, 0, 0, 0, time.UTC)
	task := store.Task{ID: "task-never", Key: "never", Cron: "0 0 31 2 *", Name: "Never", Enabled: true, Handler: "system.noop", CreatedAt: now, UpdatedAt: now}
	if err := repo.CreateTask(task); err != nil {
		t.Fatal(err)
	}
	s := NewScheduler(repo)
	// Multiple ticks in the same day record at most one failed run.
	s.tick(now)
	s.tick(now.Add(30 * time.Second))
	s.tick(now.Add(2 * time.Minute))
	runs, total, _ := repo.ListAllRuns(store.ListFilter{Page: 1, PageSize: 10})
	if total != 1 || runs[0].Status != "failed" || runs[0].Detail == "" {
		t.Fatalf("runs = %d %+v, want 1 failed with detail", total, runs)
	}
}

