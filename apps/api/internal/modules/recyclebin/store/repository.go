// Package store owns the admin.recycle-bin persistence (S-12 · GOAL-012
// D-002 §1): deleted-row snapshots with JSON payload, partial-unique while
// unrestored. Lives in a sub-package so the handler can consume the types
// without an import cycle with the module provider.
package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// TxRunner is the platform persistence boundary consumed by the repository.
type TxRunner interface {
	WithTx(context.Context, func(*sql.Tx) error) error
}

// Item is one recycle snapshot.
type Item struct {
	ID         string
	Resource   string
	ResourceID string
	Payload    map[string]any
	ActorID    string
	ActorName  string
	DeletedAt  time.Time
	RestoredAt *time.Time
}

// ErrItemNotFound is returned when a recycle id is unknown.
var ErrItemNotFound = errors.New("recycle item not found")

// ErrItemAlreadyRestored guards double-restore.
var ErrItemAlreadyRestored = errors.New("recycle item already restored")

// ListFilter bounds the recycle list query.
type ListFilter struct {
	Q        string
	Resource string
	Sort     string
	Order    string
	Page     int
	PageSize int
}

// Repository owns the recycle domain queries.
type Repository struct {
	runner TxRunner
}

// NewRepository constructs the recycle repository over a platform transaction
// runner.
func NewRepository(runner TxRunner) *Repository {
	return &Repository{runner: runner}
}

// Record inserts one snapshot (S-12 · GOAL-012 D-002 §2).
func (r *Repository) Record(item Item) error {
	payload, err := json.Marshal(item.Payload)
	if err != nil {
		return fmt.Errorf("marshal recycle payload: %w", err)
	}
	return r.runner.WithTx(context.Background(), func(tx *sql.Tx) error {
		_, err := tx.Exec(
			`INSERT INTO recycle_items (id, resource, resource_id, payload, actor_id, actor_name, deleted_at, restored_at) VALUES (?, ?, ?, ?, ?, ?, ?, NULL)`,
			item.ID, item.Resource, item.ResourceID, string(payload), item.ActorID, item.ActorName, item.DeletedAt.Unix(),
		)
		if err != nil {
			return fmt.Errorf("insert recycle item: %w", err)
		}
		return nil
	})
}

// List returns active (unrestored) snapshots, newest first by default.
func (r *Repository) List(filter ListFilter) ([]Item, int, error) {
	var items []Item
	total := 0
	err := r.runner.WithTx(context.Background(), func(tx *sql.Tx) error {
		where := `WHERE restored_at IS NULL`
		var args []any
		if filter.Resource != "" {
			where += ` AND resource = ?`
			args = append(args, filter.Resource)
		}
		if filter.Q != "" {
			where += ` AND (resource_id LIKE ? OR actor_name LIKE ?)`
			like := "%" + filter.Q + "%"
			args = append(args, like, like)
		}
		order := "DESC"
		if filter.Order == "asc" {
			order = "ASC"
		}
		sortColumn := "deleted_at"
		switch filter.Sort {
		case "resource":
			sortColumn = "resource"
		case "actorName":
			sortColumn = "actor_name"
		}
		page := filter.Page
		if page < 1 {
			page = 1
		}
		pageSize := filter.PageSize
		if pageSize < 1 {
			pageSize = 20
		}
		if pageSize > 100 {
			pageSize = 100
		}
		if err := tx.QueryRow(`SELECT COUNT(*) FROM recycle_items `+where, args...).Scan(&total); err != nil {
			return fmt.Errorf("count recycle items: %w", err)
		}
		query := `SELECT id, resource, resource_id, payload, actor_id, actor_name, deleted_at, restored_at FROM recycle_items ` + where + ` ORDER BY ` + sortColumn + ` ` + order + `, id DESC LIMIT ? OFFSET ?`
		queryArgs := append(append([]any{}, args...), pageSize, (page-1)*pageSize)
		rows, err := tx.Query(query, queryArgs...)
		if err != nil {
			return fmt.Errorf("list recycle items: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var item Item
			var payload string
			var deletedAt int64
			var restoredAt sql.NullInt64
			if err := rows.Scan(&item.ID, &item.Resource, &item.ResourceID, &payload, &item.ActorID, &item.ActorName, &deletedAt, &restoredAt); err != nil {
				return fmt.Errorf("scan recycle item: %w", err)
			}
			item.DeletedAt = time.Unix(deletedAt, 0).UTC()
			if restoredAt.Valid {
				t := time.Unix(restoredAt.Int64, 0).UTC()
				item.RestoredAt = &t
			}
			if err := json.Unmarshal([]byte(payload), &item.Payload); err != nil {
				return fmt.Errorf("unmarshal recycle payload: %w", err)
			}
			items = append(items, item)
		}
		return rows.Err()
	})
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// Get loads one snapshot by id.
func (r *Repository) Get(id string) (*Item, error) {
	var item Item
	var payload string
	var deletedAt int64
	var restoredAt sql.NullInt64
	err := r.runner.WithTx(context.Background(), func(tx *sql.Tx) error {
		row := tx.QueryRow(`SELECT id, resource, resource_id, payload, actor_id, actor_name, deleted_at, restored_at FROM recycle_items WHERE id = ?`, id)
		if err := row.Scan(&item.ID, &item.Resource, &item.ResourceID, &payload, &item.ActorID, &item.ActorName, &deletedAt, &restoredAt); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ErrItemNotFound
			}
			return fmt.Errorf("get recycle item: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	item.DeletedAt = time.Unix(deletedAt, 0).UTC()
	if restoredAt.Valid {
		t := time.Unix(restoredAt.Int64, 0).UTC()
		item.RestoredAt = &t
	}
	if err := json.Unmarshal([]byte(payload), &item.Payload); err != nil {
		return nil, fmt.Errorf("unmarshal recycle payload: %w", err)
	}
	return &item, nil
}

// MarkRestored flags a snapshot as restored (partial-unique frees the slot).
func (r *Repository) MarkRestored(id string, now time.Time) error {
	return r.runner.WithTx(context.Background(), func(tx *sql.Tx) error {
		res, err := tx.Exec(`UPDATE recycle_items SET restored_at = ? WHERE id = ? AND restored_at IS NULL`, now.Unix(), id)
		if err != nil {
			return fmt.Errorf("mark restored: %w", err)
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return ErrItemAlreadyRestored
		}
		return nil
	})
}

// PurgeAllUnrestored physically removes every active (unrestored) snapshot;
// returns the number purged. Irreversible (D-002 §3).
func (r *Repository) PurgeAllUnrestored() (int, error) {
	purged := 0
	err := r.runner.WithTx(context.Background(), func(tx *sql.Tx) error {
		res, err := tx.Exec(`DELETE FROM recycle_items WHERE restored_at IS NULL`)
		if err != nil {
			return fmt.Errorf("purge all recycle items: %w", err)
		}
		n, _ := res.RowsAffected()
		purged = int(n)
		return nil
	})
	if err != nil {
		return 0, err
	}
	return purged, nil
}

// Purge physically removes a snapshot (irreversible, D-002 §3).
func (r *Repository) Purge(id string) error {
	return r.runner.WithTx(context.Background(), func(tx *sql.Tx) error {
		res, err := tx.Exec(`DELETE FROM recycle_items WHERE id = ?`, id)
		if err != nil {
			return fmt.Errorf("purge recycle item: %w", err)
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return ErrItemNotFound
		}
		return nil
	})
}
