// In-process scheduled-task scheduler (S-04 · GOAL-010 D-002 §3): a 30s tick
// loop scans enabled tasks, computes the next matching moment from the cron
// fields and records a task_runs row per execution. Best-effort single-instance
// semantics; missed windows are not backfilled. Handlers are registered by key
// (v1 ships system.noop only — D-002 §3).
package scheduledtasks

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/scheduledtasks/store"
)

// TaskHandler executes one scheduled task run.
type TaskHandler func(ctx context.Context, task store.Task, now time.Time) error

// tickInterval is the scheduler scan period (D-002 §3).
const tickInterval = 30 * time.Second

// Scheduler runs the best-effort task loop.
type Scheduler struct {
	repository *store.Repository
	handlers   map[string]TaskHandler
	now        func() time.Time
	stop       chan struct{}
	once       sync.Once
	onceStop   sync.Once
	// lastRun deduplicates executions within the same minute slot so the 30s
	// tick never double-runs a task in one slot (D-002 §3).
	lastRun map[string]time.Time
	// unscheduled dedupes the "unschedulable" failed-run records per task per
	// day so the 30s loop cannot flood task_runs (A-003 F-002).
	unscheduled map[string]time.Time
}

// NewScheduler constructs the scheduler with the built-in handler set.
func NewScheduler(repository *store.Repository) *Scheduler {
	return &Scheduler{
		repository: repository,
		handlers: map[string]TaskHandler{
			"system.noop": func(ctx context.Context, task store.Task, now time.Time) error {
				return nil // noop handler: record the run only
			},
		},
		now:         time.Now,
		stop:        make(chan struct{}),
		lastRun:     map[string]time.Time{},
		unscheduled: map[string]time.Time{},
	}
}

// Start launches the tick loop (idempotent).
func (s *Scheduler) Start() {
	s.once.Do(func() {
		go s.loop()
	})
}

func (s *Scheduler) loop() {
	ticker := time.NewTicker(tickInterval)
	defer ticker.Stop()
	for {
		select {
		case <-s.stop:
			return
		case now := <-ticker.C:
			// W9 F-007: a panicking tick must kill neither the loop goroutine
			// nor the process; the next tick retries.
			func() {
				defer func() {
					if recovered := recover(); recovered != nil {
						slog.Error("scheduler tick panicked", "err", recovered)
					}
				}()
				s.tick(now)
			}()
		}
	}
}

// Stop halts the tick loop (idempotent — A-003 F-004).
func (s *Scheduler) Stop() {
	s.onceStop.Do(func() {
		close(s.stop)
	})
}

// tick executes every enabled task whose next moment is due.
func (s *Scheduler) tick(now time.Time) {
	tasks, err := s.repository.EnabledTasks()
	if err != nil {
		slog.Error("scheduler scan failed", "err", err)
		return
	}
	slot := now.Truncate(time.Minute)
	for _, task := range tasks {
		fields, err := store.ParseCron(task.Cron)
		if err != nil {
			continue // invalid cron rows are rejected at write time; skip defensively
		}
		if !fields.Matches(slot) {
			// W6 F1 (GOAL-006 D-001): the current minute slot does not match.
			// The previous implementation scanned the 5-year window on every
			// 30s tick (up to ~2.6M iterations per task per tick) only to be
			// told next.After(now) and continue. The match test is O(1); the
			// 5-year scan now runs at most once per day per task, and only to
			// keep the A-003 F-002 daily unschedulable diagnostic.
			day := now.Truncate(24 * time.Hour)
			if last, seen := s.unscheduled[task.ID]; !seen || last.Before(day) {
				s.unscheduled[task.ID] = day
				if _, ok := fields.Next(now); !ok {
					s.recordUnschedule(task, now)
				}
			}
			continue
		}
		if last, seen := s.lastRun[task.ID]; seen && last.Equal(slot) {
			continue // already executed in this minute slot
		}
		s.lastRun[task.ID] = slot
		s.Execute(task, now)
	}
}

// Execute runs the task handler and records the run row (manual trigger and
// scheduler share this path — D-002 §3).
func (s *Scheduler) Execute(task store.Task, now time.Time) error {
	handler, ok := s.handlers[task.Handler]
	if !ok {
		handler = s.handlers["system.noop"]
	}
	started := now
	finished := now
	detail := ""
	// W9 F-007: a panicking task handler is recorded as a failed run (same
	// durable trail as an error return) instead of crashing the process.
	runErr := func() (err error) {
		defer func() {
			if recovered := recover(); recovered != nil {
				err = fmt.Errorf("task handler panicked: %v", recovered)
			}
		}()
		return handler(context.Background(), task, now)
	}()
	if runErr != nil {
		detail = runErr.Error()
	}
	run := store.TaskRun{
		ID:        newRunID(),
		TaskID:    task.ID,
		Status:    "ran",
		StartedAt: started,
		FinishedAt: &finished,
		Detail:    detail,
		CreatedAt: now,
	}
	if detail != "" {
		run.Status = "failed"
	}
	if err := s.repository.RecordRun(run); err != nil {
		slog.Error("scheduler run record failed", "task", task.Key, "err", err)
		return err
	}
	return nil
}

// recordUnschedule records a failed run for a never-matching cron expression.
func (s *Scheduler) recordUnschedule(task store.Task, now time.Time) {
	detail := "unschedulable: no cron match within the 5-year window"
	run := store.TaskRun{
		ID:        newRunID(),
		TaskID:    task.ID,
		Status:    "failed",
		StartedAt: now,
		Detail:    detail,
		CreatedAt: now,
	}
	if err := s.repository.RecordRun(run); err != nil {
		slog.Error("scheduler unschedule record failed", "task", task.Key, "err", err)
	}
}

// HandlerKeys lists the registered handler keys (v1: system.noop).
func (s *Scheduler) HandlerKeys() []string {
	keys := make([]string, 0, len(s.handlers))
	for key := range s.handlers {
		keys = append(keys, key)
	}
	return keys
}

// newRunID returns "run-" + 16 lowercase hex chars (8 bytes of crypto/rand).
func newRunID() string {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "run-" + hex.EncodeToString([]byte(time.Now().Format("150405.000000000")))
	}
	return "run-" + hex.EncodeToString(b[:])
}
