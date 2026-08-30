package operationlog

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	settingsmigration "github.com/magicvr/schema-ui-core/apps/api/modules/settings/migration"
)

// RetentionPolicy is read from site settings on every sweep. Days and action
// are never hardcoded in the sweeper.
type RetentionPolicy struct {
	Days   int
	Action string
}

// ApplyRetention expires rows older than `days` using archive or delete.
// created_at is unix milliseconds.
func (r *Repository) ApplyRetention(now time.Time, days int, action string) (int, error) {
	if days < settingsmigration.MinOperationLogRetentionDays || days > settingsmigration.MaxOperationLogRetentionDays {
		return 0, fmt.Errorf("operationlog: invalid retention days %d", days)
	}
	action = strings.TrimSpace(action)
	if action != settingsmigration.ExpirationActionArchive && action != settingsmigration.ExpirationActionDelete {
		return 0, fmt.Errorf("operationlog: invalid expiration action %q", action)
	}
	cutoff := now.UTC().AddDate(0, 0, -days).UnixMilli()
	archivedAt := now.UTC().UnixMilli()
	var affected int64
	err := r.withTx("apply operation log retention", func(tx kernel.Tx) error {
		if action == settingsmigration.ExpirationActionArchive {
			if _, err := tx.Exec(context.Background(),
				`INSERT INTO operation_log_archive
				 (id, event, actor_id, actor_name, record_id, detail, created_at, archived_at)
				 SELECT id, event, actor_id, actor_name, record_id, detail, created_at, ?
				 FROM operation_log WHERE created_at < ? ON CONFLICT (id) DO NOTHING`,
				archivedAt, cutoff,
			); err != nil {
				return fmt.Errorf("archive rows: %w", err)
			}
			if _, err := tx.Exec(context.Background(),
				`INSERT INTO operation_log_archive_correlation (operation_id, correlation_id)
				 SELECT c.operation_id, c.correlation_id
				 FROM operation_log_correlation c
				 INNER JOIN operation_log o ON o.id = c.operation_id
				 WHERE o.created_at < ? ON CONFLICT (operation_id) DO NOTHING`,
				cutoff,
			); err != nil {
				return fmt.Errorf("archive correlations: %w", err)
			}
			if _, err := tx.Exec(context.Background(),
				`INSERT INTO operation_log_archive_session (operation_id, session_id)
				 SELECT s.operation_id, s.session_id
				 FROM operation_log_session s
				 INNER JOIN operation_log o ON o.id = s.operation_id
				 WHERE o.created_at < ? ON CONFLICT (operation_id) DO NOTHING`,
				cutoff,
			); err != nil {
				return fmt.Errorf("archive sessions: %w", err)
			}
		}
		result, err := tx.Exec(context.Background(), `DELETE FROM operation_log WHERE created_at < ?`, cutoff)
		if err != nil {
			return fmt.Errorf("delete expired rows: %w", err)
		}
		affected, err = result.RowsAffected()
		return err
	})
	return int(affected), err
}

// StartRetentionSweep reads policy from settings on each tick. stop() ends the loop.
func StartRetentionSweep(repo *Repository, loadPolicy func() (RetentionPolicy, error), interval time.Duration, logger *slog.Logger) func() {
	if interval <= 0 {
		interval = time.Hour
	}
	stop := make(chan struct{})
	go func() {
		run := func() {
			policy, err := loadPolicy()
			if err != nil {
				if logger != nil {
					logger.Error("operation log retention: load policy", "err", err)
				}
				return
			}
			n, err := repo.ApplyRetention(time.Now().UTC(), policy.Days, policy.Action)
			if err != nil {
				if logger != nil {
					logger.Error("operation log retention: apply", "err", err)
				}
				return
			}
			if n > 0 && logger != nil {
				logger.Info("operation log retention applied", "expired", n, "days", policy.Days, "action", policy.Action)
			}
		}
		run()
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				run()
			}
		}
	}()
	return func() { close(stop) }
}
