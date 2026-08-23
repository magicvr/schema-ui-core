// Package datadictionary owns the admin.data-dictionary persistence (S-01 ·
// GOAL-008 D-002 `2/`3): two-level dict_types / dict_entries with cascade
// delete and the UNIQUE(dict_key, entry_key) enumeration identity.
// Package store owns the admin.data-dictionary persistence (S-01 · GOAL-008
// D-002 §2/§3). It lives in a sub-package so the handler can consume the row
// types and sentinels without an import cycle with the module provider.
package store

import (
	"context"
	"errors"
	"fmt"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/pagination"
)

// TxRunner is the platform persistence boundary consumed by the repository.
type TxRunner interface {
	Run(context.Context, func(kernel.Tx) error) error
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
	ID      string
	DictKey string
	// DictTypeName is the owning type's display name (JOIN, GOAL-015): lets
	// the inner page show the type name instead of the raw key in columns.
	DictTypeName string
	EntryKey     string
	Label        string
	Enabled      bool
	Sort         int
	Remark       string
	// BadgeStyle is the optional badge/tag color preset (W16-F09).
	BadgeStyle string
	CreatedAt  time.Time
	UpdatedAt  time.Time
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
	err := r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		where := ""
		args := []any{}
		if q := strings.TrimSpace(filter.Q); q != "" {
			where = ` WHERE lower(key) LIKE '%' || CAST(? AS TEXT) || '%' OR lower(name) LIKE '%' || CAST(? AS TEXT) || '%'`
			args = append(args, q, q)
		}
		if err := tx.QueryRow(context.Background(), `SELECT COUNT(*) FROM dict_types`+where, args...).Scan(&total); err != nil {
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
		rows, err := tx.Query(context.Background(),
			`SELECT id, key, name, enabled, COALESCE(description, ''), sort, created_at, updated_at
			 FROM dict_types`+where+` ORDER BY `+sortCol+` `+filter.Order+` LIMIT ? OFFSET ?`,
			append(args, filter.PageSize, pagination.Offset(filter.Page, filter.PageSize, total))...,
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
	err := r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		row := tx.QueryRow(context.Background(),
			`SELECT id, key, name, enabled, COALESCE(description, ''), sort, created_at, updated_at
			 FROM dict_types WHERE id = ?`, id,
		)
		var created, updated int64
		if err := row.Scan(&t.ID, &t.Key, &t.Name, &t.Enabled, &t.Description, &t.Sort, &created, &updated); err != nil {
			if errors.Is(err, kernel.ErrNoRows) {
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
	return r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		return r.CreateTypeTx(context.Background(), tx, t)
	})
}

// CreateTypeTx inserts a type inside the caller's transaction (W11 F-008:
// recycle restore runs the row INSERT and the snapshot's MarkRestored in ONE
// transaction — a failed mark rolls the restored row back).
func (r *Repository) CreateTypeTx(ctx context.Context, tx kernel.Tx, t DictType) error {
	_, err := tx.Exec(ctx,
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
}

// UpdateType patches name/enabled/description/sort.
func (r *Repository) UpdateType(id string, name string, enabled bool, description string, sort int, now time.Time) error {
	return r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		res, err := tx.Exec(context.Background(),
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
	return r.DeleteTypeTx(context.Background(), id, nil)
}

// DeleteTypeTx removes a type and cascades to its entries inside the
// caller's transaction, optionally running the record callback (the recycle
// snapshot) in the SAME transaction (W11 F-002): a snapshot failure rolls
// the delete back. The deleted entry ids are returned for per-entry audit
// detail (A-003 F-003).
func (r *Repository) DeleteTypeTx(ctx context.Context, id string, record func(ctx context.Context, tx kernel.Tx) error) ([]string, error) {
	entryIDs := []string{}
	err := r.runner.Run(ctx, func(tx kernel.Tx) error {
		rows, err := tx.Query(ctx,
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
		res, err := tx.Exec(ctx, `DELETE FROM dict_types WHERE id = ?`, id)
		if err != nil {
			return fmt.Errorf("delete dict type: %w", err)
		}
		affected, _ := res.RowsAffected()
		if affected == 0 {
			return ErrNotFound
		}
		if record != nil {
			if err := record(ctx, tx); err != nil {
				return fmt.Errorf("record delete side effect: %w", err)
			}
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
	err := r.runner.Run(context.Background(), func(tx kernel.Tx) error {
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
			where += sep + `(lower(de.dict_key) LIKE '%' || CAST(? AS TEXT) || '%' OR lower(de.entry_key) LIKE '%' || CAST(? AS TEXT) || '%' OR lower(de.label) LIKE '%' || CAST(? AS TEXT) || '%')`
			args = append(args, q, q, q)
		}
		if err := tx.QueryRow(context.Background(), `SELECT COUNT(*) FROM dict_entries de`+where, args...).Scan(&total); err != nil {
			return fmt.Errorf("count dict entries: %w", err)
		}
		// GOAL-015 F-002/F-003 (grok audit): after the LEFT JOIN dict_types the
		// sort columns must be qualified — dict_types also carries sort/updated_at
		// and an unqualified ORDER BY is ambiguous (SQLite 500). dictTypeName
		// maps to dt.name so the sortable type-name column works.
		sortCol, ok := map[string]string{
			"dictKey": "de.dict_key", "entryKey": "de.entry_key", "label": "de.label", "sort": "de.sort", "updatedAt": "de.updated_at", "dictTypeName": "dt.name",
		}[filter.Sort]
		if !ok {
			sortCol = "entry_key"
		}
		if filter.Order != "desc" {
			filter.Order = "asc"
		}
		rows, err := tx.Query(context.Background(),
			`SELECT de.id, de.dict_key, dt.name, de.entry_key, de.label, de.enabled, de.sort, COALESCE(de.remark, ''), COALESCE(de.badge_style, 'default'), de.created_at, de.updated_at
			 FROM dict_entries de LEFT JOIN dict_types dt ON dt.key = de.dict_key`+where+` ORDER BY `+sortCol+` `+filter.Order+` LIMIT ? OFFSET ?`,
			append(args, filter.PageSize, pagination.Offset(filter.Page, filter.PageSize, total))...,
		)
		if err != nil {
			return fmt.Errorf("list dict entries: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var e DictEntry
			var created, updated int64
			if err := rows.Scan(&e.ID, &e.DictKey, &e.DictTypeName, &e.EntryKey, &e.Label, &e.Enabled, &e.Sort, &e.Remark, &e.BadgeStyle, &created, &updated); err != nil {
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
	err := r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		row := tx.QueryRow(context.Background(),
			`SELECT de.id, de.dict_key, dt.name, de.entry_key, de.label, de.enabled, de.sort, COALESCE(de.remark, ''), COALESCE(de.badge_style, 'default'), de.created_at, de.updated_at
			 FROM dict_entries de LEFT JOIN dict_types dt ON dt.key = de.dict_key WHERE de.id = ?`, id,
		)
		var created, updated int64
		if err := row.Scan(&e.ID, &e.DictKey, &e.DictTypeName, &e.EntryKey, &e.Label, &e.Enabled, &e.Sort, &e.Remark, &e.BadgeStyle, &created, &updated); err != nil {
			if errors.Is(err, kernel.ErrNoRows) {
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
	return r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		return r.CreateEntryTx(context.Background(), tx, e)
	})
}

// CreateEntryTx inserts an entry inside the caller's transaction (W11 F-008:
// recycle restore runs the row INSERT and the snapshot's MarkRestored in ONE
// transaction — a failed mark rolls the restored row back).
func (r *Repository) CreateEntryTx(ctx context.Context, tx kernel.Tx, e DictEntry) error {
	var exists int
	if err := tx.QueryRow(ctx, `SELECT COUNT(*) FROM dict_types WHERE key = ?`, e.DictKey).Scan(&exists); err != nil {
		return fmt.Errorf("check dict type: %w", err)
	}
	if exists == 0 {
		return ErrDictKeyNotFound
	}
	badgeStyle := e.BadgeStyle
	if badgeStyle == "" {
		badgeStyle = "default"
	}
	_, err := tx.Exec(ctx,
		`INSERT INTO dict_entries (id, dict_key, entry_key, label, enabled, sort, remark, badge_style, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		e.ID, e.DictKey, e.EntryKey, e.Label, boolInt(e.Enabled), e.Sort, e.Remark, badgeStyle, e.CreatedAt.Unix(), e.UpdatedAt.Unix(),
	)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") {
			return ErrEntryKeyTaken
		}
		return fmt.Errorf("insert dict entry: %w", err)
	}
	return nil
}

// UpdateEntry patches dict_key/label/enabled/sort/remark/badgeStyle with the
// same existence and uniqueness checks.
func (r *Repository) UpdateEntry(id string, dictKey string, label string, enabled bool, sort int, remark string, badgeStyle string, now time.Time) error {
	return r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		var exists int
		if err := tx.QueryRow(context.Background(), `SELECT COUNT(*) FROM dict_types WHERE key = ?`, dictKey).Scan(&exists); err != nil {
			return fmt.Errorf("check dict type: %w", err)
		}
		if exists == 0 {
			return ErrDictKeyNotFound
		}
		if badgeStyle == "" {
			badgeStyle = "default"
		}
		res, err := tx.Exec(context.Background(),
			`UPDATE dict_entries SET dict_key = ?, label = ?, enabled = ?, sort = ?, remark = ?, badge_style = ?, updated_at = ? WHERE id = ?`,
			dictKey, label, boolInt(enabled), sort, remark, badgeStyle, now.Unix(), id,
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
	return r.DeleteEntryTx(context.Background(), id, nil)
}

// DeleteEntryTx removes one entry row inside the caller's transaction,
// optionally running the record callback (recycle snapshot) in the SAME
// transaction (W11 F-002): a snapshot failure rolls the delete back.
func (r *Repository) DeleteEntryTx(ctx context.Context, id string, record func(ctx context.Context, tx kernel.Tx) error) error {
	return r.runner.Run(ctx, func(tx kernel.Tx) error {
		res, err := tx.Exec(ctx, `DELETE FROM dict_entries WHERE id = ?`, id)
		if err != nil {
			return fmt.Errorf("delete dict entry: %w", err)
		}
		affected, _ := res.RowsAffected()
		if affected == 0 {
			return ErrNotFound
		}
		if record != nil {
			if err := record(ctx, tx); err != nil {
				return fmt.Errorf("record delete side effect: %w", err)
			}
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
