// Records repository (GOAL-007 S3): the production CRUD path for the records
// entity backed by the SQLite `records` table (0003 records_persist). It is the
// single data source — the handler no longer holds an in-process slice
// (I-007-002 §5). updated_at is stored as Unix milliseconds and mapped to/from
// UTC time.Time here, so the API layer only ever sees millisecond-precision
// timestamps.
package store

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Record is the persisted records row. Field names mirror the API JSON shape;
// UpdatedAt is UTC with Unix-millisecond precision (D-004).
type Record struct {
	ID        string
	Name      string
	Status    string
	Owner     string
	UpdatedAt time.Time
}

// RecordFilter carries the list query parameters already validated by the
// handler (sort field, order, page, pageSize and search text). Sorting and
// filtering happen in SQL so the default path never loads the whole table.
type RecordFilter struct {
	Q        string
	Sort     string // name | status | owner | updatedAt (handler-validated)
	Order    string // asc | desc
	Page     int
	PageSize int
}

// RecordPatch is the set of editable fields for UpdateRecord; a nil pointer
// means "leave unchanged" (mirrors the handler's PATCH pointer semantics).
type RecordPatch struct {
	Name   *string
	Status *string
	Owner  *string
}

// ErrRecordExists is returned by CreateRecord when the generated id already
// exists — an astronomically rare PK collision the handler retries before
// giving up with INTERNAL (I-007-001 §2).
var ErrRecordExists = errors.New("store: record id already exists")

// ListRecords returns the rows matching f.Q (case-insensitive substring across
// name/status/owner), ordered by the validated sort column, plus the total
// number of matching rows before pagination.
func (s *Store) ListRecords(f RecordFilter) ([]Record, int, error) {
	where, args := recordsWhere(f.Q)
	var total int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM records`+where, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count records: %w", err)
	}

	rows, err := s.db.Query(
		`SELECT id, name, status, owner, updated_at FROM records`+where+
			` ORDER BY `+recordsSortSQL(f.Sort, f.Order)+`, id ASC`+
			` LIMIT ? OFFSET ?`,
		append(args, f.PageSize, (f.Page-1)*f.PageSize)...,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list records: %w", err)
	}
	defer rows.Close()

	items := make([]Record, 0, f.PageSize)
	for rows.Next() {
		var r Record
		var updatedAt int64
		if err := rows.Scan(&r.ID, &r.Name, &r.Status, &r.Owner, &updatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan record: %w", err)
		}
		r.UpdatedAt = time.UnixMilli(updatedAt).UTC()
		items = append(items, r)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("list records: %w", err)
	}
	return items, total, nil
}

// GetRecord fetches one row by id. Returns ErrNotFound when it does not exist.
func (s *Store) GetRecord(id string) (*Record, error) {
	var r Record
	var updatedAt int64
	err := s.db.QueryRow(
		`SELECT id, name, status, owner, updated_at FROM records WHERE id = ?`, id,
	).Scan(&r.ID, &r.Name, &r.Status, &r.Owner, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get record: %w", err)
	}
	r.UpdatedAt = time.UnixMilli(updatedAt).UTC()
	return &r, nil
}

// CreateRecord inserts a new row. The existence check and insert run in one
// transaction on the single connection, so a PK collision is detected
// deterministically (single-writer guarantee, I-007-002 §5) and reported as
// ErrRecordExists for the handler to retry.
func (s *Store) CreateRecord(r Record) (*Record, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin create record: %w", err)
	}
	defer tx.Rollback()
	var exists int
	if err := tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM records WHERE id = ?)`, r.ID).Scan(&exists); err != nil {
		return nil, fmt.Errorf("check record id: %w", err)
	}
	if exists == 1 {
		return nil, ErrRecordExists
	}
	if _, err := tx.Exec(
		`INSERT INTO records (id, name, status, owner, updated_at) VALUES (?, ?, ?, ?, ?)`,
		r.ID, r.Name, r.Status, r.Owner, r.UpdatedAt.UnixMilli(),
	); err != nil {
		return nil, fmt.Errorf("insert record: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit create record: %w", err)
	}
	return &r, nil
}

// UpdateRecord applies a PATCH to an existing record and refreshes updated_at to
// now, clamped so the returned timestamp is strictly later than the previous
// value (D-004 / I-007-001 §3.1): if now is not greater than the stored
// updated_at, it is pinned to prev+1ms. Concurrency is last-write-wins with the
// single-writer connection serializing read-modify-write. Returns ErrNotFound
// when the id does not exist.
func (s *Store) UpdateRecord(id string, patch RecordPatch, now time.Time) (*Record, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin update record: %w", err)
	}
	defer tx.Rollback()

	var cur Record
	var prevMillis int64
	err = tx.QueryRow(
		`SELECT id, name, status, owner, updated_at FROM records WHERE id = ?`, id,
	).Scan(&cur.ID, &cur.Name, &cur.Status, &cur.Owner, &prevMillis)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get record for update: %w", err)
	}
	if patch.Name != nil {
		cur.Name = *patch.Name
	}
	if patch.Status != nil {
		cur.Status = *patch.Status
	}
	if patch.Owner != nil {
		cur.Owner = *patch.Owner
	}

	newMillis := now.UnixMilli()
	if newMillis <= prevMillis {
		newMillis = prevMillis + 1 // monotonic clamp, same-millisecond rapid update
	}
	cur.UpdatedAt = time.UnixMilli(newMillis).UTC()
	if _, err := tx.Exec(
		`UPDATE records SET name = ?, status = ?, owner = ?, updated_at = ? WHERE id = ?`,
		cur.Name, cur.Status, cur.Owner, newMillis, id,
	); err != nil {
		return nil, fmt.Errorf("update record: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit update record: %w", err)
	}
	return &cur, nil
}

// DeleteRecord removes a row by id. Returns ErrNotFound when no row matched.
func (s *Store) DeleteRecord(id string) error {
	res, err := s.db.Exec(`DELETE FROM records WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete record: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete record rows affected: %w", err)
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

// recordsWhere builds the optional WHERE clause for q. instr() is a plain
// case-insensitive substring match (lower() both sides), matching the legacy
// in-process matches() semantics exactly without LIKE wildcard interpretation.
func recordsWhere(q string) (string, []any) {
	q = strings.ToLower(strings.TrimSpace(q))
	if q == "" {
		return "", nil
	}
	return ` WHERE (instr(lower(name), ?) > 0 OR instr(lower(status), ?) > 0 OR instr(lower(owner), ?) > 0)`,
		[]any{q, q, q}
}

// recordsSortSQL maps a handler-validated sort field to its column and ORDER BY
// clause. Only whitelisted fields are mapped (defaults to name); name keeps the
// legacy case-insensitive ordering via NOCASE.
func recordsSortSQL(sort, order string) string {
	col, collate := "name", " COLLATE NOCASE"
	switch sort {
	case "status":
		col, collate = "status", ""
	case "owner":
		col, collate = "owner", ""
	case "updatedAt":
		col, collate = "updated_at", ""
	}
	dir := "ASC"
	if order == "desc" {
		dir = "DESC"
	}
	return col + collate + " " + dir
}
