// Package operationlog owns the append-only operation log domain and queries.
package operationlog

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

const (
	EventAuthLogin      = "auth.login"
	EventAuthLogout     = "auth.logout"
	EventAuthRefresh    = "auth.refresh"
	EventUserCreate     = "users.create"
	EventUserUpdate     = "users.update"
	EventUserDelete     = "users.delete"
	EventRoleCreate     = "roles.create"
	EventRoleUpdate     = "roles.update"
	EventRoleDelete     = "roles.delete"
	EventSettingsUpdate = "settings.update"
	// F-03 account lifecycle events (GOAL-005 D-002 §3/§4).
	EventUserEnable            = "users.enable"
	EventUserDisable           = "users.disable"
	EventUserUnlock            = "users.unlock"
	EventAccountPasswordChange = "account.password-change"
	EventAccountSessionRevoke  = "account.session-revoke"
	// F-02 data-transfer events (GOAL-004 D-002 §3/§4).
	EventDataExport = "data.export"
	EventDataImport = "data.import"
	// S-02 file-library events (GOAL-007 D-002 §4).
	EventFileUpload   = "files.upload"
	EventFileDownload = "files.download"
	EventFileDelete   = "files.delete"
	// S-01 dictionary events (GOAL-008 D-002 §5).
	EventDictionaryCreate = "dictionary.create"
	EventDictionaryUpdate = "dictionary.update"
	EventDictionaryDelete = "dictionary.delete"
	// S-04 scheduled-task events (GOAL-010 D-002 §4).
	EventTaskCreate = "scheduled-tasks.create"
	EventTaskUpdate = "scheduled-tasks.update"
	EventTaskDelete = "scheduled-tasks.delete"
)

// Operation is one append-only operation log row.
type Operation struct {
	ID        string
	Event     string
	ActorID   string
	ActorName string
	RecordID  *string
	Detail    *string
	CreatedAt time.Time
}

// OperationFilter carries handler-validated activity list parameters.
type OperationFilter struct {
	Q        string
	Sort     string
	Order    string
	Page     int
	PageSize int
}

// Recorder is the operation-log write boundary consumed by business handlers.
type Recorder interface {
	RecordOperation(Operation) error
}

// Reader is the query boundary consumed by the optional Activity module.
type Reader interface {
	ListOperationsFiltered(OperationFilter) ([]Operation, int, error)
	GetOperation(string) (*Operation, error)
}

// TxRunner is the platform transaction boundary consumed by operationlog.
type TxRunner interface {
	WithTx(context.Context, func(*sql.Tx) error) error
}

var ErrNotFound = errors.New("operationlog: not found")

// Repository owns operation-log reads and writes.
type Repository struct {
	runner TxRunner

	failureMu sync.RWMutex
	failure   error
}

// NewRepository constructs the repository over the platform transaction runner.
func NewRepository(runner TxRunner) *Repository {
	return &Repository{runner: runner}
}

// RecordOperation appends one row. Callers deliberately decide best-effort policy.
func (r *Repository) RecordOperation(operation Operation) error {
	r.failureMu.RLock()
	forced := r.failure
	r.failureMu.RUnlock()
	if forced != nil {
		return fmt.Errorf("record operation %s: %w", operation.Event, forced)
	}
	return r.withTx("record operation "+operation.Event, func(tx *sql.Tx) error {
		var recordID, detail any
		if operation.RecordID != nil {
			recordID = *operation.RecordID
		}
		if operation.Detail != nil {
			detail = *operation.Detail
		}
		if _, err := tx.Exec(
			`INSERT INTO operation_log (id, event, actor_id, actor_name, record_id, detail, created_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?)`,
			operation.ID, operation.Event, operation.ActorID, operation.ActorName,
			recordID, detail, operation.CreatedAt.UnixMilli(),
		); err != nil {
			return err
		}
		return nil
	})
}

// SetOperationLogError configures the test-only best-effort failure seam.
func (r *Repository) SetOperationLogError(err error) {
	r.failureMu.Lock()
	r.failure = err
	r.failureMu.Unlock()
}

// ListOperations returns the newest limit rows.
func (r *Repository) ListOperations(limit int) ([]Operation, error) {
	if limit <= 0 {
		return nil, nil
	}
	items, _, err := r.ListOperationsFiltered(OperationFilter{
		Sort: "createdAt", Order: "desc", Page: 1, PageSize: limit,
	})
	return items, err
}

// ListOperationsFiltered returns one page and its pre-pagination total.
func (r *Repository) ListOperationsFiltered(filter OperationFilter) ([]Operation, int, error) {
	var items []Operation
	var total int
	err := r.withTx("list operations", func(tx *sql.Tx) error {
		where, args := operationsWhere(filter.Q)
		if err := tx.QueryRow(`SELECT COUNT(*) FROM operation_log`+where, args...).Scan(&total); err != nil {
			return fmt.Errorf("count: %w", err)
		}
		rows, err := tx.Query(
			`SELECT id, event, actor_id, actor_name, record_id, detail, created_at
			 FROM operation_log`+where+
				` ORDER BY `+operationsSortSQL(filter.Sort, filter.Order)+`, id DESC`+
				` LIMIT ? OFFSET ?`,
			append(args, filter.PageSize, (filter.Page-1)*filter.PageSize)...,
		)
		if err != nil {
			return err
		}
		defer rows.Close()
		items = make([]Operation, 0, filter.PageSize)
		for rows.Next() {
			operation, err := scanOperation(rows)
			if err != nil {
				return err
			}
			items = append(items, operation)
		}
		if err := rows.Err(); err != nil {
			return fmt.Errorf("rows: %w", err)
		}
		return nil
	})
	return items, total, err
}

// GetOperation returns one row by id.
func (r *Repository) GetOperation(id string) (*Operation, error) {
	var operation Operation
	err := r.withTx("get operation", func(tx *sql.Tx) error {
		var err error
		operation, err = scanOperation(tx.QueryRow(
			`SELECT id, event, actor_id, actor_name, record_id, detail, created_at
			 FROM operation_log WHERE id = ?`, id,
		))
		return err
	})
	if err != nil {
		return nil, err
	}
	return &operation, nil
}

func (r *Repository) withTx(operation string, fn func(*sql.Tx) error) error {
	if r == nil || r.runner == nil {
		return fmt.Errorf("%s: operationlog repository is not configured", operation)
	}
	if err := r.runner.WithTx(context.Background(), fn); err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	return nil
}

func scanOperation(row interface{ Scan(...any) error }) (Operation, error) {
	var operation Operation
	var recordID, detail sql.NullString
	var createdAt int64
	err := row.Scan(
		&operation.ID, &operation.Event, &operation.ActorID, &operation.ActorName,
		&recordID, &detail, &createdAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return Operation{}, ErrNotFound
	}
	if err != nil {
		return Operation{}, fmt.Errorf("scan: %w", err)
	}
	if recordID.Valid {
		operation.RecordID = &recordID.String
	}
	if detail.Valid {
		operation.Detail = &detail.String
	}
	operation.CreatedAt = time.UnixMilli(createdAt).UTC()
	return operation, nil
}

func operationsWhere(query string) (string, []any) {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return "", nil
	}
	return ` WHERE (instr(lower(event), ?) > 0 OR instr(lower(actor_name), ?) > 0 OR instr(lower(COALESCE(detail,'')), ?) > 0 OR instr(lower(COALESCE(record_id,'')), ?) > 0)`,
		[]any{query, query, query, query}
}

func operationsSortSQL(sort, order string) string {
	column := "created_at"
	switch sort {
	case "event":
		column = "event"
	case "actorName":
		column = "actor_name"
	}
	direction := "DESC"
	if order == "asc" {
		direction = "ASC"
	}
	return column + " " + direction
}