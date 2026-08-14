// Scheduled tasks surface (S-04 · GOAL-010 D-002 `4): task definition CRUD
// (resource factory), manual trigger, per-task run history and the global run
// history. tasks.read / tasks.write (admin-only) and scheduled-tasks.* audit
// events.
package handler

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
	tasksstore "github.com/magicvr/schema-ui-core/apps/api/internal/modules/scheduledtasks/store"
)

// TasksRepository is the persistence boundary consumed by the task resources.
type TasksRepository interface {
	ListTasks(tasksstore.ListFilter) ([]tasksstore.Task, int, error)
	GetTask(string) (*tasksstore.Task, error)
	CreateTask(tasksstore.Task) error
	UpdateTask(id, cron, name string, enabled bool, description, handler string, now time.Time) error
	DeleteTask(string) error
	RecordRun(tasksstore.TaskRun) error
	ListTaskRuns(string, tasksstore.ListFilter) ([]tasksstore.TaskRun, int, error)
	ListAllRuns(tasksstore.ListFilter) ([]tasksstore.TaskRun, int, error)
}

// taskEntity adapts the tasks repository to the generic factory.
type taskEntity struct {
	repository TasksRepository
	operations operationlog.Recorder
}

func (e *taskEntity) List(f resourceFilter) ([]map[string]any, int, error) {
	rows, total, err := e.repository.ListTasks(tasksstore.ListFilter{Q: f.Q, Sort: f.Sort, Order: f.Order, Page: f.Page, PageSize: f.PageSize})
	if err != nil {
		return nil, 0, err
	}
	items := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		items = append(items, taskToMap(row))
	}
	return items, total, nil
}

func (e *taskEntity) Get(id string) (map[string]any, error) {
	row, err := e.repository.GetTask(id)
	if err != nil {
		return nil, mapTaskStoreError(err)
	}
	return taskToMap(*row), nil
}

func (e *taskEntity) Create(body map[string]any, id string, now time.Time, actor account.User) (map[string]any, error) {
	key := stringField(body, "key")
	name := stringField(body, "name")
	cron := stringField(body, "cron")
	if key == "" || name == "" || cron == "" {
		return nil, &DomainError{Status: http.StatusBadRequest, Code: "INVALID_CREATE_FIELD", Message: "key, name and cron are required"}
	}
	if _, err := tasksstore.ParseCron(cron); err != nil {
		return nil, &DomainError{Status: http.StatusBadRequest, Code: "INVALID_CRON", Message: "invalid cron expression: " + err.Error()}
	}
	handler := stringField(body, "handler")
	if handler == "" {
		handler = "system.noop"
	}
	row := tasksstore.Task{
		ID: id, Key: key, Cron: cron, Name: name,
		Enabled: boolField(body, "enabled", true), Description: stringField(body, "description"),
		Handler: handler, CreatedAt: now, UpdatedAt: now,
	}
	if err := e.repository.CreateTask(row); err != nil {
		return nil, mapTaskStoreError(err)
	}
	recordTaskEvent(e.operations, operationlog.EventTaskCreate, actor, id, now)
	created, err := e.repository.GetTask(id)
	if err != nil {
		return nil, err
	}
	return taskToMap(*created), nil
}

func (e *taskEntity) Update(id string, body map[string]any, now time.Time, actor account.User) (map[string]any, error) {
	existing, err := e.repository.GetTask(id)
	if err != nil {
		return nil, mapTaskStoreError(err)
	}
	cron := stringField(body, "cron")
	if cron == "" {
		cron = existing.Cron
	}
	if _, err := tasksstore.ParseCron(cron); err != nil {
		return nil, &DomainError{Status: http.StatusBadRequest, Code: "INVALID_CRON", Message: "invalid cron expression: " + err.Error()}
	}
	name := stringField(body, "name")
	if name == "" {
		name = existing.Name
	}
	description := stringField(body, "description")
	if _, present := body["description"]; !present {
		description = existing.Description
	}
	enabled := boolField(body, "enabled", existing.Enabled)
	handler := stringField(body, "handler")
	if _, present := body["handler"]; !present || handler == "" {
		handler = existing.Handler
	}
	if err := e.repository.UpdateTask(id, cron, name, enabled, description, handler, now); err != nil {
		return nil, mapTaskStoreError(err)
	}
	recordTaskEvent(e.operations, operationlog.EventTaskUpdate, actor, id, now)
	row, err := e.repository.GetTask(id)
	if err != nil {
		return nil, mapTaskStoreError(err)
	}
	return taskToMap(*row), nil
}

func (e *taskEntity) Delete(id string, actor account.User) error {
	err := e.repository.DeleteTask(id)
	if err != nil {
		return mapTaskStoreError(err)
	}
	recordTaskEvent(e.operations, operationlog.EventTaskDelete, actor, id, time.Now().UTC())
	return nil
}

// taskRunsEntity adapts the global run history to a read-only resource.
type taskRunsEntity struct {
	repository TasksRepository
}

func (e *taskRunsEntity) List(f resourceFilter) ([]map[string]any, int, error) {
	rows, total, err := e.repository.ListAllRuns(tasksstore.ListFilter{Q: f.Q, Sort: f.Sort, Order: f.Order, Page: f.Page, PageSize: f.PageSize})
	if err != nil {
		return nil, 0, err
	}
	items := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		items = append(items, taskRunToMap(row))
	}
	return items, total, nil
}

func (e *taskRunsEntity) Get(id string) (map[string]any, error) {
	return nil, errResourceNotFound
}

func (e *taskRunsEntity) Create(map[string]any, string, time.Time, account.User) (map[string]any, error) {
	return nil, errReadOnlyResource
}

func (e *taskRunsEntity) Update(string, map[string]any, time.Time, account.User) (map[string]any, error) {
	return nil, errReadOnlyResource
}

func (e *taskRunsEntity) Delete(string, account.User) error {
	return errReadOnlyResource
}

// TaskRunner executes one task immediately (manual trigger + scheduler).
type TaskRunner interface {
	Execute(task tasksstore.Task, now time.Time) error
	HandlerKeys() []string
}

// ScheduledTaskRoutes returns the admin.scheduled-tasks HTTP surface.
func ScheduledTaskRoutes(a *auth.Authenticator, repository TasksRepository, runner TaskRunner, operations operationlog.Recorder, moduleID string) []kernel.RouteContribution {
	routes := ResourceRoutes(a, Resource{
		ID:              "scheduled-tasks",
		Path:            "/api/scheduled-tasks",
		Listable:        true,
		SortFields:      []string{"key", "name", "updatedAt"},
		QSearch:         true,
		Entity:          &taskEntity{repository: repository, operations: operations},
		CreateFields:    []string{"key", "name", "cron"},
		PatchFields:     []string{"name", "cron", "description"},
		JSONFields:      []string{"enabled", "handler"},
		PermissionRead:  "tasks.read",
		PermissionWrite: "tasks.write",
		NotFoundCode:    "TASK_NOT_FOUND",
		NewID:           newTaskID,
	}, moduleID)

	// Global run history (read-only).
	routes = append(routes, ResourceRoutes(a, Resource{
		ID:              "task-runs",
		Path:            "/api/task-runs",
		Listable:        true,
		ReadOnly:        true,
		SortFields:      []string{"startedAt"},
		QSearch:         true,
		Entity:          &taskRunsEntity{repository: repository},
		PermissionRead:  "tasks.read",
		PermissionWrite: "tasks.write", // unused when ReadOnly
		NotFoundCode:    "TASK_RUN_NOT_FOUND",
	}, moduleID)...)

	// Manual trigger.
	routes = append(routes, kernel.RouteContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey("POST", "/api/scheduled-tasks/{id}/run")},
		Method:               "POST",
		Pattern:              "/api/scheduled-tasks/{id}/run",
		Handler: a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ok := requirePermission(w, r, "tasks.write"); !ok {
				return
			}
			id := r.PathValue("id")
			task, err := repository.GetTask(id)
			if err != nil {
				writeTaskEntityError(w, r, err)
				return
			}
			if err := runner.Execute(*task, time.Now().UTC()); err != nil {
				writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "task execution failed")
				return
			}
			w.WriteHeader(http.StatusNoContent)
		})),
	})

	// Per-task run history.
	routes = append(routes, kernel.RouteContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey("GET", "/api/scheduled-tasks/{id}/runs")},
		Method:               "GET",
		Pattern:              "/api/scheduled-tasks/{id}/runs",
		Handler: a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ok := requirePermission(w, r, "tasks.read"); !ok {
				return
			}
			id := r.PathValue("id")
			rows, total, err := repository.ListTaskRuns(id, tasksstore.ListFilter{Page: 1, PageSize: 50})
			if err != nil {
				writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not list task runs")
				return
			}
			items := make([]map[string]any, 0, len(rows))
			for _, row := range rows {
				items = append(items, taskRunToMap(row))
			}
			writeJSON(w, http.StatusOK, resourceList{Items: items, Total: total, Page: 1, PageSize: 50})
		})),
	})
	return routes
}

func taskToMap(t tasksstore.Task) map[string]any {
	return map[string]any{
		"id": t.ID, "key": t.Key, "cron": t.Cron, "name": t.Name,
		"enabled": t.Enabled, "description": t.Description, "handler": t.Handler,
		"createdAt": t.CreatedAt.Unix(), "updatedAt": t.UpdatedAt.Unix(),
	}
}

func taskRunToMap(rn tasksstore.TaskRun) map[string]any {
	row := map[string]any{
		"id": rn.ID, "taskId": rn.TaskID, "status": rn.Status,
		"startedAt": rn.StartedAt.Unix(), "detail": rn.Detail,
	}
	if rn.FinishedAt != nil {
		row["finishedAt"] = rn.FinishedAt.Unix()
	} else {
		row["finishedAt"] = nil
	}
	return row
}

func mapTaskStoreError(err error) error {
	switch {
	case errors.Is(err, tasksstore.ErrNotFound):
		return errResourceNotFound
	case errors.Is(err, tasksstore.ErrKeyTaken):
		return &DomainError{Status: http.StatusConflict, Code: "TASK_KEY_TAKEN", Message: "a scheduled task with that key already exists"}
	default:
		return err
	}
}

func writeTaskEntityError(w http.ResponseWriter, r *http.Request, err error) {
	if errors.Is(err, tasksstore.ErrNotFound) {
		writeLocalizedError(w, r, http.StatusNotFound, "TASK_NOT_FOUND", "no task with that id")
		return
	}
	writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not load task")
}

func recordTaskEvent(operations operationlog.Recorder, event string, user account.User, id string, now time.Time) {
	if operations == nil {
		return
	}
	recordID := id
	if err := operations.RecordOperation(operationlog.Operation{
		ID: newOperationID(), Event: event, ActorID: user.ID, ActorName: user.Name,
		RecordID: &recordID, CreatedAt: now.UTC(),
	}); err != nil {
		slog.Error("operation log write failed", "event", event, "err", err)
	}
}

// newTaskID returns "task-" + 16 lowercase hex chars.
func newTaskID() (string, error) {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return "task-" + hex.EncodeToString(b[:]), nil
}
