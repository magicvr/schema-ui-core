// Operation log repository (GOAL-008 R5 S6 optional bonus checkpoint,
// I-008-003): an append-only SQLite log covering records writes and auth key
// events. Writes are best-effort at the handler layer — a logging failure must
// never turn a successful business operation into a failure — so the store
// method only reports the error and the handler records it to the service log.
package store

import (
	"database/sql"
	"fmt"
	"time"
)

// OperationEvent enumerates the frozen operation log event set (I-008-003 §2).
const (
	EventRecordCreate = "records.create"
	EventRecordUpdate = "records.update"
	EventRecordDelete = "records.delete"
	EventAuthLogin    = "auth.login"
	EventAuthLogout   = "auth.logout"
	EventAuthRefresh  = "auth.refresh"
)

// Operation is one append-only operation log row. RecordID is set for records
// events and nil for auth events; Detail is an optional JSON summary (record
// name / username), never a token or secret (I-008-003 §3).
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

// ListOperations returns the most recent limit operations ordered by
// created_at DESC, id DESC. A non-positive limit returns an empty slice
// (I-008-003 §4).
func (s *Store) ListOperations(limit int) ([]Operation, error) {
	if limit <= 0 {
		return nil, nil
	}
	rows, err := s.db.Query(
		`SELECT id, event, actor_id, actor_name, record_id, detail, created_at
		 FROM operation_log ORDER BY created_at DESC, id DESC LIMIT ?`, limit)
	if err != nil {
		return nil, fmt.Errorf("list operations: %w", err)
	}
	defer rows.Close()
	var out []Operation
	for rows.Next() {
		var (
			op        Operation
			recordID  sql.NullString
			detail    sql.NullString
			createdAt int64
		)
		if err := rows.Scan(&op.ID, &op.Event, &op.ActorID, &op.ActorName, &recordID, &detail, &createdAt); err != nil {
			return nil, fmt.Errorf("scan operation: %w", err)
		}
		if recordID.Valid {
			op.RecordID = &recordID.String
		}
		if detail.Valid {
			op.Detail = &detail.String
		}
		op.CreatedAt = time.UnixMilli(createdAt).UTC()
		out = append(out, op)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list operations: %w", err)
	}
	return out, nil
}
