// Operation log repository (GOAL-008 R5 S6 optional bonus checkpoint,
// I-008-003): an append-only SQLite log covering resource writes and auth key
// events. Writes are best-effort at the handler layer — a logging failure must
// never turn a successful business operation into a failure — so the store
// method only reports the error and the handler records it to the service log.
package store

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

// OperationEvent enumerates the frozen operation log event set (I-008-003 §2;
// users/roles events added by GOAL-011 0005, I-011-001 §5).
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
	EventSettingsUpdate = "settings.update" // GOAL-013 branding · A-006 R-003
)

// Operation is one append-only operation log row. RecordID is the historical
// target-id field for resource events and nil for auth events; Detail is an
// optional JSON summary, never a token or secret (I-008-003 §3).
type Operation struct {
	ID        string
	Event     string
	ActorID   string
	ActorName string
	RecordID  *string
	Detail    *string
	CreatedAt time.Time
}

// RecordOperation appends one operation log row (I-008-003 §4). The caller is
// responsible for a unique ID. Failures are returned (not swallowed here) so
// the handler can keep the best-effort contract explicit.
func (s *Store) RecordOperation(op Operation) error {
	var recordID, detail any
	if op.RecordID != nil {
		recordID = *op.RecordID
	}
	if op.Detail != nil {
		detail = *op.Detail
	}
	if _, err := s.db.Exec(
		`INSERT INTO operation_log (id, event, actor_id, actor_name, record_id, detail, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		op.ID, op.Event, op.ActorID, op.ActorName, recordID, detail, op.CreatedAt.UnixMilli(),
	); err != nil {
		return fmt.Errorf("record operation %s: %w", op.Event, err)
	}
	return nil
}

// OperationFilter carries list query parameters for the activity (operations)
// read API. Sort is a whitelist column name applied in SQL.
type OperationFilter struct {
	Q        string
	Sort     string // createdAt | event | actorName (handler-validated)
	Order    string // asc | desc
	Page     int
	PageSize int
}

// ListOperations returns the most recent limit operations ordered by
// created_at DESC, id DESC. A non-positive limit returns an empty slice
// (I-008-003 §4). Prefer ListOperationsFiltered for the product activity page.
func (s *Store) ListOperations(limit int) ([]Operation, error) {
	if limit <= 0 {
		return nil, nil
	}
	items, _, err := s.ListOperationsFiltered(OperationFilter{
		Sort: "createdAt", Order: "desc", Page: 1, PageSize: limit,
	})
	return items, err
}

// ListOperationsFiltered returns a page of operations plus total count before
// pagination (activity read surface).
func (s *Store) ListOperationsFiltered(f OperationFilter) ([]Operation, int, error) {
	where, args := operationsWhere(f.Q)
	var total int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM operation_log`+where, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count operations: %w", err)
	}
	orderSQL := operationsSortSQL(f.Sort, f.Order)
	rows, err := s.db.Query(
		`SELECT id, event, actor_id, actor_name, record_id, detail, created_at
		 FROM operation_log`+where+
			` ORDER BY `+orderSQL+`, id DESC`+
			` LIMIT ? OFFSET ?`,
		append(args, f.PageSize, (f.Page-1)*f.PageSize)...,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list operations: %w", err)
	}
	defer rows.Close()
	out := make([]Operation, 0, f.PageSize)
	for rows.Next() {
		op, err := scanOperation(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, op)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("list operations rows: %w", err)
	}
	return out, total, nil
}

// GetOperation fetches one operation log row by id.
func (s *Store) GetOperation(id string) (*Operation, error) {
	row := s.db.QueryRow(
		`SELECT id, event, actor_id, actor_name, record_id, detail, created_at
		 FROM operation_log WHERE id = ?`, id)
	op, err := scanOperation(row)
	if err != nil {
		return nil, err
	}
	return &op, nil
}

func scanOperation(row interface{ Scan(...any) error }) (Operation, error) {
	var (
		op        Operation
		recordID  sql.NullString
		detail    sql.NullString
		createdAt int64
	)
	err := row.Scan(&op.ID, &op.Event, &op.ActorID, &op.ActorName, &recordID, &detail, &createdAt)
	if errors.Is(err, sql.ErrNoRows) {
		return Operation{}, ErrNotFound
	}
	if err != nil {
		return Operation{}, fmt.Errorf("scan operation: %w", err)
	}
	if recordID.Valid {
		op.RecordID = &recordID.String
	}
	if detail.Valid {
		op.Detail = &detail.String
	}
	op.CreatedAt = time.UnixMilli(createdAt).UTC()
	return op, nil
}

func operationsWhere(q string) (string, []any) {
	q = strings.ToLower(strings.TrimSpace(q))
	if q == "" {
		return "", nil
	}
	return ` WHERE (instr(lower(event), ?) > 0 OR instr(lower(actor_name), ?) > 0 OR instr(lower(COALESCE(detail,'')), ?) > 0 OR instr(lower(COALESCE(record_id,'')), ?) > 0)`,
		[]any{q, q, q, q}
}

func operationsSortSQL(sort, order string) string {
	col := "created_at"
	switch sort {
	case "event":
		col = "event"
	case "actorName":
		col = "actor_name"
	case "createdAt":
		col = "created_at"
	}
	dir := "DESC"
	if order == "asc" {
		dir = "ASC"
	}
	return col + " " + dir
}
