// Package datadictionary owns the admin.data-dictionary persistence (S-01 ·
// GOAL-008 D-002 `2/`3): two-level dict_types / dict_entries with cascade
// delete and the UNIQUE(dict_key, entry_key) enumeration identity.
// Package store owns the admin.data-dictionary persistence (S-01 · GOAL-008
// D-002 §2/§3). It lives in a sub-package so the handler can consume the row
// types and sentinels without an import cycle with the module provider.
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

// Repository owns the dictionary domain queries.
type Repository struct {
	runner TxRunner
}

// NewRepository constructs the dictionary repository over a platform
// transaction runner.
func NewRepository(runner TxRunner) *Repository {
	return &Repository{runner: runner}
}

// Domain sentinels mapped by the handler to the frozen error codes.
var (
	ErrNotFound        = errors.New("dictionary row not found")
	ErrTypeKeyTaken    = errors.New("dict type key already exists")
	ErrEntryKeyTaken   = errors.New("dict entry key already exists")
	ErrDictKeyNotFound = errors.New("dict type key does not exist")
)

// DictType is one dictionary type row.
type DictType struct {
	ID          string
	Key         string
	Name        string
	Enabled     bool
	Description string
	Sort        int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// DictEntry is one dictionary entry row.
type DictEntry struct {
	ID        string
	DictKey   string
	// DictTypeName is the owning type's display name (JOIN, GOAL-015): lets
	// the inner page show the type name instead of the raw key in columns.
	DictTypeName string
	EntryKey  string
	Label     string
	Enabled   bool
	Sort      int
	Remark    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ListFilter carries validated list parameters from the resource factory.
// DictKey optionally narrows entries to one dictionary type (GOAL-015).
type ListFilter struct {
	Q        string
	DictKey  string
	Sort     string
	Order    string
	Page     int
	PageSize int
}

// ListTypes returns the paged type rows.
func (r *Repository) ListTypes(filter ListFilter) ([]DictType, int, error) {
	types := []DictType{}
	var total int
	err := r.runner.WithTx(context.Background(), func(tx *sql.Tx) error {
		where := ""
		args := []any{}
		if q := strings.TrimSpace(filter.Q); q != "" {
			where = ` WHERE instr(lower(key), ?) > 0 OR instr(lower(name), ?) > 0`
			args = append(args, q, q)
		}
		if err := tx.QueryRow(`SELECT COUNT(*) FROM dict_types`+where, args...).Scan(&total); err != nil {
			return fmt.Errorf("count dict types: %w", err)
		}
		sortCol, ok := map[string]string{
			"key": "key", "name": "name", "sort": "sort", "updatedAt": "updated_at",
		}[filter.Sort]
		if !ok {
			sortCol = "key"
		}
		if filter.Order != "desc" {
			filter.Order = "asc"
		}
		rows, err := tx.Query(
			`SELECT id, key, name, enabled, COALESCE(description, ''), sort, created_at, updated_at
			 FROM dict_types`+where+` ORDER BY `+sortCol+` `+filter.Order+` LIMIT ? OFFSET ?`,
			append(args, filter.PageSize, (filter.Page-1)*filter.PageSize)...,
		)
		if err != nil {
			return fmt.Errorf("list dict types: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var t DictType
			var created, updated int64
			if err := rows.Scan(&t.ID, &t.Key, &t.Name, &t.Enabled, &t.Description, &t.Sort, &created, &updated); err != nil {
				return fmt.Errorf("scan dict type: %w", err)
			}
			t.CreatedAt = time.Unix(created, 0)
			t.UpdatedAt = time.Unix(updated, 0)
			types = append(types, t)
		}
		return rows.Err()
	})
	return types, total, err
}

// GetType returns one type by id.
func (r *Repository) GetType(id string) (*DictType, error) {
	var t DictType
	err := r.runner.WithTx(context.Background(), func(tx *sql.Tx) error {
		row := tx.QueryRow(
			`SELECT id, key, name, enabled, COALESCE(description, ''), sort, created_at, updated_at
			 FROM dict_types WHERE id = ?`, id,
		)
		var created, updated int64
		if err := row.Scan(&t.ID, &t.Key, &t.Name, &t.Enabled, &t.Description, &t.Sort, &created, &updated); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ErrNotFound
			}
			return fmt.Errorf("get dict type: %w", err)
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

// CreateType inserts a type; key collisions fail with ErrTypeKeyTaken.
func (r *Repository) CreateType(t DictType) error {
	return r.runner.WithTx(context.Background(), func(tx *sql.Tx) error {
		_, err := tx.Exec(
			`INSERT INTO dict_types (id, key, name, enabled, description, sort, created_at, updated_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			t.ID, t.Key, t.Name, boolInt(t.Enabled), t.Description, t.Sort, t.CreatedAt.Unix(), t.UpdatedAt.Unix(),
		)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "unique") {
				return ErrTypeKeyTaken
			}
			return fmt.Errorf("insert dict type: %w", err)
		}
		return nil
	})
}

// UpdateType patches name/enabled/description/sort.
func (r *Repository) UpdateType(id string, name string, enabled bool, description string, sort int, now time.Time) error {
	return r.runner.WithTx(context.Background(), func(tx *sql.Tx) error {
		res, err := tx.Exec(
			`UPDATE dict_types SET name = ?, enabled = ?, description = ?, sort = ?, updated_at = ? WHERE id = ?`,
			name, boolInt(enabled), description, sort, now.Unix(), id,
		)
		if err != nil {
			return fmt.Errorf("update dict type: %w", err)
		}
		affected, _ := res.RowsAffected()
		if affected == 0 {
			return ErrNotFound
		}
		return nil
	})
}

// DeleteType removes a type and cascades to its entries; the deleted entry
// ids are returned so the caller can record per-entry audit detail (A-003
// F-003).
func (r *Repository) DeleteType(id string) ([]string, error) {
	entryIDs := []string{}
	err := r.runner.WithTx(context.Background(), func(tx *sql.Tx) error {
		rows, err := tx.Query(
			`SELECT e.id FROM dict_entries e JOIN dict_types t ON e.dict_key = t.key WHERE t.id = ?`, id,
		)
		if err != nil {
			return fmt.Errorf("list dict entries for delete: %w", err)
		}
		for rows.Next() {
			var entryID string
			if err := rows.Scan(&entryID); err != nil {
				rows.Close()
				return fmt.Errorf("scan dict entry id: %w", err)
			}
			entryIDs = append(entryIDs, entryID)
		}
		rows.Close()
		if err := rows.Err(); err != nil {
			return err
		}
		res, err := tx.Exec(`DELETE FROM dict_types WHERE id = ?`, id)
		if err != nil {
			return fmt.Errorf("delete dict type: %w", err)
		}
		affected, _ := res.RowsAffected()
		if affected == 0 {
			return ErrNotFound
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return entryIDs, nil
}

// ListEntries returns the paged entry rows (global list; q covers dict_key /
// entry_key / label — D-002 `3).
func (r *Repository) ListEntries(filter ListFilter) ([]DictEntry, int, error) {
	entries := []DictEntry{}
	var total int
	err := r.runner.WithTx(context.Background(), func(tx *sql.Tx) error {
		where := ""
		args := []any{}
		// GOAL-015: exact dict-key narrowing (inner page) composes with q.
		if dictKey := strings.TrimSpace(filter.DictKey); dictKey != "" {
			where = ` WHERE de.dict_key = ?`
			args = append(args, dictKey)
		}
		if q := strings.TrimSpace(filter.Q); q != "" {
			sep := " AND "
			if where == "" {
				sep = " WHERE "
			}
			where += sep + `(instr(lower(de.dict_key), ?) > 0 OR instr(lower(de.entry_key), ?) > 0 OR instr(lower(de.label), ?) > 0)`
			args = append(args, q, q, q)
		}
		if err := tx.QueryRow(`SELECT COUNT(*) FROM dict_entries de`+where, args...).Scan(&total); err != nil {
			return fmt.Errorf("count dict entries: %w", err)
		}
		sortCol, ok := map[string]string{
			"dictKey": "dict_key", "entryKey": "entry_key", "label": "label", "sort": "sort", "updatedAt": "updated_at",
		}[filter.Sort]
		if !ok {
			sortCol = "entry_key"
		}
		if filter.Order != "desc" {
			filter.Order = "asc"
		}
		rows, err := tx.Query(
			`SELECT de.id, de.dict_key, dt.name, de.entry_key, de.label, de.enabled, de.sort, COALESCE(de.remark, ''), de.created_at, de.updated_at
			 FROM dict_entries de LEFT JOIN dict_types dt ON dt.key = de.dict_key`+where+` ORDER BY `+sortCol+` `+filter.Order+` LIMIT ? OFFSET ?`,
			append(args, filter.PageSize, (filter.Page-1)*filter.PageSize)...,
		)
		if err != nil {
			return fmt.Errorf("list dict entries: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var e DictEntry
			var created, updated int64
			if err := rows.Scan(&e.ID, &e.DictKey, &e.DictTypeName, &e.EntryKey, &e.Label, &e.Enabled, &e.Sort, &e.Remark, &created, &updated); err != nil {
				return fmt.Errorf("scan dict entry: %w", err)
			}
			e.CreatedAt = time.Unix(created, 0)
			e.UpdatedAt = time.Unix(updated, 0)
			entries = append(entries, e)
		}
		return rows.Err()
	})
	return entries, total, err
}

// GetEntry returns one entry by id.
func (r *Repository) GetEntry(id string) (*DictEntry, error) {
	var e DictEntry
	err := r.runner.WithTx(context.Background(), func(tx *sql.Tx) error {
		row := tx.QueryRow(
			`SELECT de.id, de.dict_key, dt.name, de.entry_key, de.label, de.enabled, de.sort, COALESCE(de.remark, ''), de.created_at, de.updated_at
			 FROM dict_entries de LEFT JOIN dict_types dt ON dt.key = de.dict_key WHERE de.id = ?`, id,
		)
		var created, updated int64
		if err := row.Scan(&e.ID, &e.DictKey, &e.DictTypeName, &e.EntryKey, &e.Label, &e.Enabled, &e.Sort, &e.Remark, &created, &updated); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ErrNotFound
			}
			return fmt.Errorf("get dict entry: %w", err)
		}
		e.CreatedAt = time.Unix(created, 0)
		e.UpdatedAt = time.Unix(updated, 0)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// CreateEntry inserts an entry; the dict_key must exist (ErrDictKeyNotFound)
// and (dict_key, entry_key) must be unique (ErrEntryKeyTaken).
func (r *Repository) CreateEntry(e DictEntry) error {
	return r.runner.WithTx(context.Background(), func(tx *sql.Tx) error {
		var exists int
		if err := tx.QueryRow(`SELECT COUNT(*) FROM dict_types WHERE key = ?`, e.DictKey).Scan(&exists); err != nil {
			return fmt.Errorf("check dict type: %w", err)
		}
		if exists == 0 {
			return ErrDictKeyNotFound
		}
		_, err := tx.Exec(
			`INSERT INTO dict_entries (id, dict_key, entry_key, label, enabled, sort, remark, created_at, updated_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			e.ID, e.DictKey, e.EntryKey, e.Label, boolInt(e.Enabled), e.Sort, e.Remark, e.CreatedAt.Unix(), e.UpdatedAt.Unix(),
		)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "unique") {
				return ErrEntryKeyTaken
			}
			return fmt.Errorf("insert dict entry: %w", err)
		}
		return nil
	})
}

// UpdateEntry patches dict_key/label/enabled/sort/remark with the same
// existence and uniqueness checks.
func (r *Repository) UpdateEntry(id string, dictKey string, label string, enabled bool, sort int, remark string, now time.Time) error {
	return r.runner.WithTx(context.Background(), func(tx *sql.Tx) error {
		var exists int
		if err := tx.QueryRow(`SELECT COUNT(*) FROM dict_types WHERE key = ?`, dictKey).Scan(&exists); err != nil {
			return fmt.Errorf("check dict type: %w", err)
		}
		if exists == 0 {
			return ErrDictKeyNotFound
		}
		res, err := tx.Exec(
			`UPDATE dict_entries SET dict_key = ?, label = ?, enabled = ?, sort = ?, remark = ?, updated_at = ? WHERE id = ?`,
			dictKey, label, boolInt(enabled), sort, remark, now.Unix(), id,
		)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "unique") {
				return ErrEntryKeyTaken
			}
			return fmt.Errorf("update dict entry: %w", err)
		}
		affected, _ := res.RowsAffected()
		if affected == 0 {
			return ErrNotFound
		}
		return nil
	})
}

// DeleteEntry removes one entry row.
func (r *Repository) DeleteEntry(id string) error {
	return r.runner.WithTx(context.Background(), func(tx *sql.Tx) error {
		res, err := tx.Exec(`DELETE FROM dict_entries WHERE id = ?`, id)
		if err != nil {
			return fmt.Errorf("delete dict entry: %w", err)
		}
		affected, _ := res.RowsAffected()
		if affected == 0 {
			return ErrNotFound
		}
		return nil
	})
}

func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
