// Scheduled tasks surface (S-04 · GOAL-010 D-002 `4): task definition CRUD
// (resource factory), manual trigger, per-task run history and the global run
// history. tasks.read / tasks.write (admin-only) and scheduled-tasks.* audit
// events.
package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/errorcatalog"
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
	runner     TaskRunner
	operations operationlog.Recorder
}

var cronWeekdaysZh = [7]string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"}
var cronWeekdaysEn = [7]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}

func cronFieldSingle(field map[int]bool) (int, bool) {
	if len(field) != 1 {
		return 0, false
	}
	for value := range field {
		return value, true
	}
	return 0, false
}

func cronFieldStep(field map[int]bool, min, max int) (int, bool) {
	if len(field) < 2 {
		return 0, false
	}
	values := make([]int, 0, len(field))
	for value := range field {
		values = append(values, value)
	}
	for i := 1; i < len(values); i++ {
		j := i
		for j > 0 && values[j-1] > values[j] {
			values[j-1], values[j] = values[j], values[j-1]
			j--
		}
	}
	step := values[1] - values[0]
	if step <= 0 {
		return 0, false
	}
	expected := 0
	for value := min; value <= max; value += step {
		if !field[value] {
			return 0, false
		}
		expected++
	}
	if expected != len(field) {
		return 0, false
	}
	return step, true
}

func cronClock(hour, minute int) string {
	return fmt.Sprintf("%02d:%02d", hour, minute)
}

// describeCron returns a locale-aware natural-language description for a
// parsed 5-field cron (W17). nextRuns still carries the concrete evidence.
func describeCron(fields tasksstore.CronFields, locale string) string {
	zh := locale == "zh-CN"
	fullMinute := len(fields[0]) == 60
	fullHour := len(fields[1]) == 24
	fullDay := len(fields[2]) == 31
	fullMonth := len(fields[3]) == 12
	fullDow := len(fields[4]) == 7
	minute, hasMinute := cronFieldSingle(fields[0])
	hour, hasHour := cronFieldSingle(fields[1])
	day, hasDay := cronFieldSingle(fields[2])
	dow, hasDow := cronFieldSingle(fields[4])

	if fullMinute && fullHour && fullDay && fullMonth && fullDow {
		if zh {
			return "每分钟"
		}
		return "Every minute"
	}
	if minuteStep, ok := cronFieldStep(fields[0], 0, 59); ok && fullHour && fullDay && fullMonth && fullDow {
		if zh {
			return "每 " + strconv.Itoa(minuteStep) + " 分钟"
		}
		return "Every " + strconv.Itoa(minuteStep) + " minutes"
	}
	if hasMinute && fullHour && fullDay && fullMonth && fullDow {
		if zh {
			return "每小时的第 " + strconv.Itoa(minute) + " 分钟"
		}
		return "Every hour at minute " + strconv.Itoa(minute)
	}
	if hasMinute && hasHour && fullDay && fullMonth && fullDow {
		clock := cronClock(hour, minute)
		if zh {
			return "每天 " + clock
		}
		return "Every day at " + clock
	}
	if hasMinute && hasHour && hasDow && fullDay && fullMonth {
		clock := cronClock(hour, minute)
		if dow < 0 || dow > 6 {
			dow = 0
		}
		if zh {
			return "每周" + cronWeekdaysZh[dow] + " " + clock
		}
		return "Every " + cronWeekdaysEn[dow] + " at " + clock
	}
	if hasMinute && hasHour && hasDay && fullMonth && fullDow {
		clock := cronClock(hour, minute)
		if zh {
			return "每月 " + strconv.Itoa(day) + " 日 " + clock
		}
		return "On day " + strconv.Itoa(day) + " of every month at " + clock
	}
	if zh {
		return "5 段 Cron 计划"
	}
	return "5-field cron schedule"
}

// validateHandler rejects unknown handler keys at write time (A-003 F-003);
// a typo must never silently fall back to system.noop.
func (e *taskEntity) validateHandler(handler string) error {
	if handler == "" {
		return nil
	}
	for _, key := range e.runner.HandlerKeys() {
		if key == handler {
			return nil
		}
	}
	return &DomainError{Status: http.StatusBadRequest, Code: "INVALID_HANDLER", Message: "unknown task handler: " + handler}
}

func (e *taskEntity) List(f resourceFilter) ([]map[string]any, int, error) {
	var enabled *bool
	if raw, ok := f.Extra["enabled"]; ok && (raw == "true" || raw == "false") {
		v := raw == "true"
		enabled = &v
	}
	rows, total, err := e.repository.ListTasks(tasksstore.ListFilter{Q: f.Q, Sort: f.Sort, Order: f.Order, Page: f.Page, PageSize: f.PageSize, Enabled: enabled})
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
	if err := e.validateHandler(handler); err != nil {
		return nil, err
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
	if err := e.validateHandler(handler); err != nil {
		return nil, err
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

// DeleteTrashTx implements handler.TrashTxDeleter (W11 F-002): the task
	// delete and the recycle snapshot commit in ONE transaction — a snapshot
	// failure rolls the delete back. The audit event is recorded only after
	// the transaction committed.
func (e *taskEntity) DeleteTrashTx(ctx context.Context, id string, actor account.User, now time.Time, record func(context.Context, kernel.Tx) error) error {
	txrepo, ok := e.repository.(interface {
		DeleteTaskTx(ctx context.Context, id string, record func(context.Context, kernel.Tx) error) error
	})
	if !ok {
		return errors.New("tasks repository does not support transactional delete")
	}
	if err := txrepo.DeleteTaskTx(ctx, id, record); err != nil {
		return mapTaskStoreError(err)
	}
	recordTaskEvent(e.operations, operationlog.EventTaskDelete, actor, id, now)
	return nil
}

// Delete serves the legacy delete path (no trash / no transactional
// recorder); the factory prefers DeleteTrashTx when both sides opt in.
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
	rows, total, err := e.repository.ListAllRuns(tasksstore.ListFilter{Q: f.Q, Sort: f.Sort, Order: f.Order, Page: f.Page, PageSize: f.PageSize, Status: f.Extra["status"]})
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
func ScheduledTaskRoutes(a *auth.Authenticator, repository TasksRepository, runner TaskRunner, operations operationlog.Recorder, moduleID string, trash ...TrashRecorder) []kernel.RouteContribution {
	// S-12 (GOAL-012 D-002 §2): optional recycle-bin snapshot hook. nil keeps
	// the delete semantics byte-identical.
	var recorder TrashRecorder
	if len(trash) > 0 {
		recorder = trash[0]
	}
	routes := ResourceRoutes(a, Resource{
		ID:         "scheduled-tasks",
		Path:       "/api/scheduled-tasks",
		Listable:   true,
		SortFields: []string{"key", "name", "updatedAt"},
		QSearch:    true,
		// T-02 (GOAL-013 D-003): enabled state select on the tasks search form.
		ExtraQuery:      []string{"enabled"},
		Entity:          &taskEntity{repository: repository, runner: runner, operations: operations},
		CreateFields:    []string{"key", "name", "cron"},
		PatchFields:     []string{"name", "cron", "description"},
		JSONFields:      []string{"enabled", "handler"},
		PermissionRead:  "tasks.read",
		PermissionWrite: "tasks.write",
		NotFoundCode:    "TASK_NOT_FOUND",
		NewID:           newTaskID,
		Trash:           recorder,
	}, moduleID)

	// F-01 (W14 D-003): expose the registered handler directory so the task
	// create/edit forms can offer a real handler instead of an invisible
	// system.noop default.
	routes = append(routes, kernel.RouteContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey("GET", "/api/scheduled-tasks/handlers")},
		Method:               "GET",
		Pattern:              "/api/scheduled-tasks/handlers",
		Handler: a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ok := requirePermission(w, r, "tasks.read"); !ok {
				return
			}
			keys := runner.HandlerKeys()
			items := make([]map[string]any, 0, len(keys))
			for _, key := range keys {
				items = append(items, map[string]any{"key": key, "label": key})
			}
			pageSize := len(items)
			if pageSize < 1 {
				pageSize = 1
			}
			writeJSON(w, http.StatusOK, resourceList{Items: items, Total: len(items), Page: 1, PageSize: pageSize})
		})),
	})

	// W16-F05: cron preview — parse a 5-field expression and return the next
	// three matching run times for the task form.
	routes = append(routes, kernel.RouteContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey("POST", "/api/scheduled-tasks/cron/preview")},
		Method:               "POST",
		Pattern:              "/api/scheduled-tasks/cron/preview",
		Handler: a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ok := requirePermission(w, r, "tasks.read"); !ok {
				return
			}
			var body struct {
				Cron string `json:"cron"`
			}
			r.Body = http.MaxBytesReader(w, r.Body, 4<<10)
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil || strings.TrimSpace(body.Cron) == "" {
				writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_CRON", "cron is required")
				return
			}
			fields, err := tasksstore.ParseCron(body.Cron)
			if err != nil {
				writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_CRON", "invalid cron expression: "+err.Error())
				return
			}
			now := time.Now().UTC()
			nextRuns := []string{}
			cursor := now
			for len(nextRuns) < 3 {
				next, ok := fields.Next(cursor)
				if !ok {
					break
				}
				nextRuns = append(nextRuns, next.Format(time.RFC3339))
				cursor = next.Add(time.Minute)
			}
			writeJSON(w, http.StatusOK, map[string]any{
				"description": describeCron(fields, errorcatalog.Negotiate(r)),
				"nextRuns":    nextRuns,
			})
		})),
	})

	// Global run history (read-only).
	routes = append(routes, ResourceRoutes(a, Resource{
		ID:         "task-runs",
		Path:       "/api/task-runs",
		Listable:   true,
		ReadOnly:   true,
		SortFields: []string{"startedAt"},
		QSearch:    true,
		// T-02 (GOAL-013 D-003): status select on the task-runs search form.
		ExtraQuery:      []string{"status"},
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
			page, ok := intParam(r.URL.Query().Get("page"), 1)
			if !ok {
				writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PAGE", "page must be a positive integer")
				return
			}
			pageSize, ok := intParam(r.URL.Query().Get("pageSize"), DefaultPageSize)
			if !ok || pageSize > maxPageSize {
				writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PAGE_SIZE", "pageSize must be a positive integer not exceeding 100")
				return
			}
			rows, total, err := repository.ListTaskRuns(id, tasksstore.ListFilter{Page: page, PageSize: pageSize})
			if err != nil {
				writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not list task runs")
				return
			}
			items := make([]map[string]any, 0, len(rows))
			for _, row := range rows {
				items = append(items, taskRunToMap(row))
			}
			writeJSON(w, http.StatusOK, resourceList{Items: items, Total: total, Page: page, PageSize: pageSize})
		})),
	})
	return routes
}

func taskToMap(t tasksstore.Task) map[string]any {
	return map[string]any{
		"id": t.ID, "key": t.Key, "cron": t.Cron, "name": t.Name,
		"enabled": t.Enabled, "description": t.Description, "handler": t.Handler,
		"createdAt": formatRFC3339Milli(t.CreatedAt), "updatedAt": formatRFC3339Milli(t.UpdatedAt),
	}
}

func taskRunToMap(rn tasksstore.TaskRun) map[string]any {
	row := map[string]any{
		"id": rn.ID, "taskId": rn.TaskID, "status": rn.Status,
		"startedAt": formatRFC3339Milli(rn.StartedAt), "detail": rn.Detail,
	}
	if rn.FinishedAt != nil {
		row["finishedAt"] = formatRFC3339Milli(*rn.FinishedAt)
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
	recordAudit(operations, user, event, id, nil, now.UTC(), nil)
}

// newTaskID returns "task-" + 16 lowercase hex chars.
func newTaskID() (string, error) {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return "task-" + hex.EncodeToString(b[:]), nil
}
