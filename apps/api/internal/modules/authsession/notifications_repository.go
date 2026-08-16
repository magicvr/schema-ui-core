// In-app notification persistence (F-04 · GOAL-006 D-002 `2/`3).
//
// Rows are per-user; read_at NULL = unread. New notifications are pruned at a
// per-user cap (oldest READ rows first; unread rows are never pruned so users
// cannot lose unread security notices). The master switch lives on users
// (notifications_enabled, migration 0017) and is checked before producing.
package authsession

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

// maxNotificationsPerUser is the per-user retention cap (D-002 `2).
const maxNotificationsPerUser = 500

// Notification is one in-app notification row.
type Notification struct {
	ID        string
	UserID    string
	Event     string
	Title     string
	Body      string
	ReadAt    *time.Time
	CreatedAt time.Time
}

// NotificationFilter carries list parameters for the notifications endpoint.
type NotificationFilter struct {
	// UnreadOnly is the legacy boolean filter (unreadOnly=true query param).
	UnreadOnly bool
	// T-02 (GOAL-013 D-003): keyword search over title/body plus an exact
	// read-state filter. Read=nil means no constraint; Read=&true means the
	// row is read, Read=&false means unread.
	Q    string
	Read *bool
	Page int
	PageSize int
}

// NotificationsEnabledFor reports whether a user accepts new notifications.
func (r *Repository) NotificationsEnabledFor(userID string) (bool, error) {
	var enabled int
	err := r.withTx("notifications enabled for", func(tx *sql.Tx) error {
		if err := tx.QueryRow(`SELECT notifications_enabled FROM users WHERE id = ?`, userID).Scan(&enabled); err != nil {
			return fmt.Errorf("read notifications_enabled: %w", err)
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	return enabled == 1, nil
}

// SetNotificationsEnabled flips the per-user master switch.
func (r *Repository) SetNotificationsEnabled(userID string, enabled bool, now time.Time) error {
	return r.withTx("set notifications enabled", func(tx *sql.Tx) error {
		value := 0
		if enabled {
			value = 1
		}
		if _, err := tx.Exec(
			`UPDATE users SET notifications_enabled = ?, updated_at = ? WHERE id = ?`,
			value, now.Unix(), userID,
		); err != nil {
			return fmt.Errorf("update notifications_enabled: %w", err)
		}
		return nil
	})
}

// CreateNotification inserts one notification for a user (best-effort hook
// target). When the per-user cap is exceeded, the oldest READ rows are pruned
// in the same transaction; unread rows are never pruned.
func (r *Repository) CreateNotification(n Notification, now time.Time) error {
	return r.withTx("create notification", func(tx *sql.Tx) error {
		if _, err := tx.Exec(
			`INSERT INTO notifications (id, user_id, event, title, body, read_at, created_at)
			 VALUES (?, ?, ?, ?, ?, NULL, ?)`,
			n.ID, n.UserID, n.Event, n.Title, n.Body, now.Unix(),
		); err != nil {
			return fmt.Errorf("insert notification: %w", err)
		}
		// Prune: keep at most maxNotificationsPerUser rows; delete oldest READ
		// rows beyond the cap, never unread ones (two-step, explicitly correct).
		var count int
		if err := tx.QueryRow(`SELECT COUNT(*) FROM notifications WHERE user_id = ?`, n.UserID).Scan(&count); err != nil {
			return fmt.Errorf("count notifications for prune: %w", err)
		}
		if count > maxNotificationsPerUser {
			excess := count - maxNotificationsPerUser
			rows, err := tx.Query(
				`SELECT id FROM notifications WHERE user_id = ? AND read_at IS NOT NULL ORDER BY created_at ASC LIMIT ?`,
				n.UserID, excess,
			)
			if err != nil {
				return fmt.Errorf("list prune candidates: %w", err)
			}
			var ids []string
			for rows.Next() {
				var id string
				if err := rows.Scan(&id); err != nil {
					rows.Close()
					return fmt.Errorf("scan prune candidate: %w", err)
				}
				ids = append(ids, id)
			}
			if err := rows.Err(); err != nil {
				rows.Close()
				return fmt.Errorf("iterate prune candidates: %w", err)
			}
			if err := rows.Close(); err != nil {
				return err
			}
			for _, id := range ids {
				if _, err := tx.Exec(`DELETE FROM notifications WHERE id = ?`, id); err != nil {
					return fmt.Errorf("prune notification %s: %w", id, err)
				}
			}
		}
		return nil
	})
}

// ListNotifications returns one page (newest first) for a user.
func (r *Repository) ListNotifications(userID string, filter NotificationFilter) ([]Notification, int, error) {
	var items []Notification
	var total int
	err := r.withTx("list notifications", func(tx *sql.Tx) error {
		where := ` WHERE user_id = ?`
		args := []any{userID}
		if q := strings.TrimSpace(filter.Q); q != "" {
			where += ` AND (instr(lower(title), ?) > 0 OR instr(lower(body), ?) > 0)`
			args = append(args, q, q)
		}
		if filter.UnreadOnly {
			where += ` AND read_at IS NULL`
		}
		if filter.Read != nil {
			if *filter.Read {
				where += ` AND read_at IS NOT NULL`
			} else {
				where += ` AND read_at IS NULL`
			}
		}
		if err := tx.QueryRow(`SELECT COUNT(*) FROM notifications`+where, args...).Scan(&total); err != nil {
			return fmt.Errorf("count notifications: %w", err)
		}
		rows, err := tx.Query(
			`SELECT id, user_id, event, title, body, read_at, created_at FROM notifications`+where+
				` ORDER BY created_at DESC, id DESC LIMIT ? OFFSET ?`,
			append(args, filter.PageSize, (filter.Page-1)*filter.PageSize)...,
		)
		if err != nil {
			return fmt.Errorf("query notifications: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			notification, err := scanNotification(rows)
			if err != nil {
				return err
			}
			items = append(items, *notification)
		}
		return rows.Err()
	})
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// MarkNotificationRead marks one owned notification read (idempotent). Unknown
// or foreign ids fail closed with ErrNotFound.
func (r *Repository) MarkNotificationRead(id, userID string, now time.Time) error {
	return r.withTx("mark notification read", func(tx *sql.Tx) error {
		res, err := tx.Exec(
			`UPDATE notifications SET read_at = COALESCE(read_at, ?) WHERE id = ? AND user_id = ?`,
			now.Unix(), id, userID,
		)
		if err != nil {
			return fmt.Errorf("mark notification read: %w", err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("rows affected: %w", err)
		}
		if affected == 0 {
			var owned int
			if err := tx.QueryRow(`SELECT COUNT(*) FROM notifications WHERE id = ? AND user_id = ?`, id, userID).Scan(&owned); err != nil {
				return fmt.Errorf("check notification: %w", err)
			}
			if owned == 0 {
				return ErrNotFound
			}
			// owned but already read: idempotent success
		}
		return nil
	})
}

// MarkAllNotificationsRead marks every unread notification of a user read.
func (r *Repository) MarkAllNotificationsRead(userID string, now time.Time) (int, error) {
	var count int64
	err := r.withTx("mark all notifications read", func(tx *sql.Tx) error {
		res, err := tx.Exec(
			`UPDATE notifications SET read_at = ? WHERE user_id = ? AND read_at IS NULL`,
			now.Unix(), userID,
		)
		if err != nil {
			return fmt.Errorf("mark all notifications read: %w", err)
		}
		count, err = res.RowsAffected()
		return err
	})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

// UnreadNotificationCount returns the unread count for the bell badge.
func (r *Repository) UnreadNotificationCount(userID string) (int, error) {
	var count int
	err := r.withTx("unread notification count", func(tx *sql.Tx) error {
		if err := tx.QueryRow(
			`SELECT COUNT(*) FROM notifications WHERE user_id = ? AND read_at IS NULL`, userID,
		).Scan(&count); err != nil {
			return fmt.Errorf("count unread notifications: %w", err)
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}

func scanNotification(row interface{ Scan(...any) error }) (*Notification, error) {
	var n Notification
	var readAt *int64
	var createdAt int64
	err := row.Scan(&n.ID, &n.UserID, &n.Event, &n.Title, &n.Body, &readAt, &createdAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan notification: %w", err)
	}
	if readAt != nil {
		value := time.Unix(*readAt, 0).UTC()
		n.ReadAt = &value
	}
	n.CreatedAt = time.Unix(createdAt, 0).UTC()
	return &n, nil
}