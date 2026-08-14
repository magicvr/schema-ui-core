// Package store owns the admin.scheduled-tasks persistence (S-04 · GOAL-010
// D-002 `1): scheduled_tasks definitions + task_runs history. Lives in a
// sub-package so the handler can consume the row types and sentinels without
// an import cycle with the module provider.
package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

// TxRunner is the platform persistence boundary consumed by the repository.
type TxRunner interface {
	WithTx(context.Context, func(*sql.Tx) error) error
}

// Repository owns the scheduled-task domain queries.
type Repository struct {
	runner TxRunner
}

// NewRepository constructs the tasks repository over a platform transaction
// runner.
func NewRepository(runner TxRunner) *Repository {
	return &Repository{runner: runner}
}

// Domain sentinels mapped by the handler to the frozen error codes.
var (
	ErrNotFound     = errors.New("scheduled task not found")
	ErrKeyTaken     = errors.New("scheduled task key already exists")
	ErrInvalidCron  = errors.New("invalid cron expression")
)

// Task is one scheduled task definition row.
type Task struct {
	ID          string
	Key         string
	Cron        string
	Name        string
	Enabled     bool
	Description string
	Handler     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TaskRun is one execution history row.
type TaskRun struct {
	ID         string
	TaskID     string
	Status     string
	StartedAt  time.Time
	FinishedAt *time.Time
	Detail     string
	CreatedAt  time.Time
}

// ListFilter carries validated list parameters from the resource factory.
type ListFilter struct {
	Q        string
	Sort     string
	Order    string
	Page     int
	PageSize int
}

// ListTasks returns the paged task rows.
func (r *Repository) ListTasks(filter ListFilter) ([]Task, int, error) {
	tasks := []Task{}
	var total int
	err := r.runner.WithTx(context.Background(), func(tx *sql.Tx) error {
		where := ""
		args := []any{}
		if q := strings.TrimSpace(filter.Q); q != "" {
			where = ` WHERE instr(lower(key), ?) > 0 OR instr(lower(name), ?) > 0`
			args = append(args, q, q)
		}
		if err := tx.QueryRow(`SELECT COUNT(*) FROM scheduled_tasks`+where, args...).Scan(&total); err != nil {
			return fmt.Errorf("count tasks: %w", err)
		}
		sortCol, ok := map[string]string{
			"key": "key", "name": "name", "updatedAt": "updated_at",
		}[filter.Sort]
		if !ok {
			sortCol = "key"
		}
		if filter.Order != "desc" {
			filter.Order = "asc"
		}
		rows, err := tx.Query(
			`SELECT id, key, cron, name, enabled, COALESCE(description, ''), handler, created_at, updated_at
			 FROM scheduled_tasks`+where+` ORDER BY `+sortCol+` `+filter.Order+` LIMIT ? OFFSET ?`,
			append(args, filter.PageSize, (filter.Page-1)*filter.PageSize)...,
		)
		if err != nil {
			return fmt.Errorf("list tasks: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var t Task
			var created, updated int64
			if err := rows.Scan(&t.ID, &t.Key, &t.Cron, &t.Name, &t.Enabled, &t.Description, &t.Handler, &created, &updated); err != nil {
				return fmt.Errorf("scan task: %w", err)
			}
			t.CreatedAt = time.Unix(created, 0)
			t.UpdatedAt = time.Unix(updated, 0)
			tasks = append(tasks, t)
		}
		return rows.Err()
	})
	return tasks, total, err
}

// GetTask returns one task by id.
func (r *Repository) GetTask(id string) (*Task, error) {
	var t Task
	err := r.runner.WithTx(context.Background(), func(tx *sql.Tx) error {
		var created, updated int64
		row := tx.QueryRow(
			`SELECT id, key, cron, name, enabled, COALESCE(description, ''), handler, created_at, updated_at
			 FROM scheduled_tasks WHERE id = ?`, id,
		)
		if err := row.Scan(&t.ID, &t.Key, &t.Cron, &t.Name, &t.Enabled, &t.Description, &t.Handler, &created, &updated); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ErrNotFound
			}
			return fmt.Errorf("get task: %w", err)
		}
		t.CreatedAt = time.Unix(created, 0)
		t.UpdatedAt = time.Unix(updated, 0)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// CreateTask inserts a task; key collisions fail with ErrKeyTaken.
func (r *Repository) CreateTask(t Task) error {
	return r.runner.WithTx(context.Background(), func(tx *sql.Tx) error {
		_, err := tx.Exec(
			`INSERT INTO scheduled_tasks (id, key, cron, name, enabled, description, handler, created_at, updated_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			t.ID, t.Key, t.Cron, t.Name, boolInt(t.Enabled), t.Description, t.Handler, t.CreatedAt.Unix(), t.UpdatedAt.Unix(),
		)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "unique") {
				return ErrKeyTaken
			}
			return fmt.Errorf("insert task: %w", err)
		}
		return nil
	})
}

// UpdateTask patches cron/name/enabled/description/handler.
func (r *Repository) UpdateTask(id, cron, name string, enabled bool, description, handler string, now time.Time) error {
	return r.runner.WithTx(context.Background(), func(tx *sql.Tx) error {
		res, err := tx.Exec(
			`UPDATE scheduled_tasks SET cron = ?, name = ?, enabled = ?, description = ?, handler = ?, updated_at = ? WHERE id = ?`,
			cron, name, boolInt(enabled), description, handler, now.Unix(), id,
		)
		if err != nil {
			return fmt.Errorf("update task: %w", err)
		}
		affected, _ := res.RowsAffected()
		if affected == 0 {
			return ErrNotFound
		}
		return nil
	})
}

// DeleteTask removes a task and cascades to its run rows.
func (r *Repository) DeleteTask(id string) error {
	return r.runner.WithTx(context.Background(), func(tx *sql.Tx) error {
		res, err := tx.Exec(`DELETE FROM scheduled_tasks WHERE id = ?`, id)
		if err != nil {
			return fmt.Errorf("delete task: %w", err)
		}
		affected, _ := res.RowsAffected()
		if affected == 0 {
			return ErrNotFound
		}
		return nil
	})
}

// EnabledTasks returns all enabled tasks (scheduler scan).
func (r *Repository) EnabledTasks() ([]Task, error) {
	tasks := []Task{}
	err := r.runner.WithTx(context.Background(), func(tx *sql.Tx) error {
		rows, err := tx.Query(
			`SELECT id, key, cron, name, enabled, COALESCE(description, ''), handler, created_at, updated_at
			 FROM scheduled_tasks WHERE enabled = 1`,
		)
		if err != nil {
			return fmt.Errorf("list enabled tasks: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var t Task
			var created, updated int64
			if err := rows.Scan(&t.ID, &t.Key, &t.Cron, &t.Name, &t.Enabled, &t.Description, &t.Handler, &created, &updated); err != nil {
				return fmt.Errorf("scan enabled task: %w", err)
			}
			t.CreatedAt = time.Unix(created, 0)
			t.UpdatedAt = time.Unix(updated, 0)
			tasks = append(tasks, t)
		}
		return rows.Err()
	})
	return tasks, err
}

// RecordRun inserts one run row.
func (r *Repository) RecordRun(run TaskRun) error {
	return r.runner.WithTx(context.Background(), func(tx *sql.Tx) error {
		finished := int64(0)
		if run.FinishedAt != nil {
			finished = run.FinishedAt.Unix()
		}
		_, err := tx.Exec(
			`INSERT INTO task_runs (id, task_id, status, started_at, finished_at, detail, created_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?)`,
			run.ID, run.TaskID, run.Status, run.StartedAt.Unix(), finished, run.Detail, run.CreatedAt.Unix(),
		)
		if err != nil {
			return fmt.Errorf("insert task run: %w", err)
		}
		return nil
	})
}

// ListTaskRuns returns the paged run rows for one task.
func (r *Repository) ListTaskRuns(taskID string, filter ListFilter) ([]TaskRun, int, error) {
	runs := []TaskRun{}
	var total int
	err := r.runner.WithTx(context.Background(), func(tx *sql.Tx) error {
		if err := tx.QueryRow(`SELECT COUNT(*) FROM task_runs WHERE task_id = ?`, taskID).Scan(&total); err != nil {
			return fmt.Errorf("count task runs: %w", err)
		}
		rows, err := tx.Query(
			`SELECT id, task_id, status, started_at, COALESCE(finished_at, 0), COALESCE(detail, ''), created_at
			 FROM task_runs WHERE task_id = ? ORDER BY started_at DESC LIMIT ? OFFSET ?`,
			taskID, filter.PageSize, (filter.Page-1)*filter.PageSize,
		)
		if err != nil {
			return fmt.Errorf("list task runs: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var rn TaskRun
			var started, finished, created int64
			if err := rows.Scan(&rn.ID, &rn.TaskID, &rn.Status, &started, &finished, &rn.Detail, &created); err != nil {
				return fmt.Errorf("scan task run: %w", err)
			}
			rn.StartedAt = time.Unix(started, 0)
			if finished > 0 {
				f := time.Unix(finished, 0)
				rn.FinishedAt = &f
			}
			rn.CreatedAt = time.Unix(created, 0)
			runs = append(runs, rn)
		}
		return rows.Err()
	})
	return runs, total, err
}

// ListAllRuns returns the paged run rows across tasks (global history).
func (r *Repository) ListAllRuns(filter ListFilter) ([]TaskRun, int, error) {
	runs := []TaskRun{}
	var total int
	err := r.runner.WithTx(context.Background(), func(tx *sql.Tx) error {
		where := ""
		args := []any{}
		if q := strings.TrimSpace(filter.Q); q != "" {
			where = ` WHERE instr(lower(detail), ?) > 0 OR task_id IN (SELECT id FROM scheduled_tasks WHERE instr(lower(key), ?) > 0)`
			args = append(args, q, q)
		}
		if err := tx.QueryRow(`SELECT COUNT(*) FROM task_runs`+where, args...).Scan(&total); err != nil {
			return fmt.Errorf("count all runs: %w", err)
		}
		rows, err := tx.Query(
			`SELECT id, task_id, status, started_at, COALESCE(finished_at, 0), COALESCE(detail, ''), created_at
			 FROM task_runs`+where+` ORDER BY started_at DESC LIMIT ? OFFSET ?`,
			append(args, filter.PageSize, (filter.Page-1)*filter.PageSize)...,
		)
		if err != nil {
			return fmt.Errorf("list all runs: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var rn TaskRun
			var started, finished, created int64
			if err := rows.Scan(&rn.ID, &rn.TaskID, &rn.Status, &started, &finished, &rn.Detail, &created); err != nil {
				return fmt.Errorf("scan task run: %w", err)
			}
			rn.StartedAt = time.Unix(started, 0)
			if finished > 0 {
				f := time.Unix(finished, 0)
				rn.FinishedAt = &f
			}
			rn.CreatedAt = time.Unix(created, 0)
			runs = append(runs, rn)
		}
		return rows.Err()
	})
	return runs, total, err
}

func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
