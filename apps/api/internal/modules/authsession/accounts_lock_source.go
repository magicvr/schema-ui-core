// Per-(account | client identity) login-lockout state (GOAL-014 D-002 ·
// W13 F-007 targeted-DoS fix): a third party who knows only a username can
// no longer lock the legitimate account out — failures accumulate per
// SOURCE, and a lock on the pair denies that source only. The GLOBAL
// consecutive-failure ceiling stays on users.failed_login_count (threshold
// raised to auth.LockThresholdFailures) with a 24h sliding restart via
// users.last_login_failure_at, so distributed guessing still hits a hard
// brake while targeted abuse became 20× more expensive and visible
// (OnLockOpened fires on the global lock only).
package authsession

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

// legacyLoginSource is the bucket key used when the caller supplies no
// client identity (variadic-absent Login calls: tests/dev). All such calls
// share one bucket, preserving the pre-GOAL-014 observable behavior there.
const legacyLoginSource = "-"

// lockCounterWindow is the source-counter slide interval. It mirrors
// auth.LockWindow; duplicated because auth imports THIS module (an import
// here would cycle).
const lockCounterWindow = 15 * time.Minute

// errSourceInsertRace marks the concurrent-first-failure INSERT conflict: the
// losing transaction must retry in a FRESH one (a failed INSERT aborts its
// own postgres transaction — W9 F-001 precedent).
var errSourceInsertRace = errors.New("authsession: login failure insert race")

func normalizeLoginSource(ip string) string {
	trimmed := strings.TrimSpace(ip)
	if trimmed == "" {
		return legacyLoginSource
	}
	return trimmed
}

// RecordLoginFailureFor bumps the (account|source)-scoped consecutive-failure
// counter and, past the threshold, opens the pair lock ending at lockedUntil.
// It reports whether THIS failure opened the lock. The counter slides: a row
// whose last move is older than the lock window restarts at 1 instead of
// accumulating forever. Mirrors W9 F-004's atomic-increment discipline.
func (r *Repository) RecordLoginFailureFor(userID, ip string, threshold int, lockedUntil time.Time, now time.Time) (bool, error) {
	locked, err := r.recordLoginFailureForOnce(userID, ip, threshold, lockedUntil, now)
	if errors.Is(err, errSourceInsertRace) {
		// Concurrent first failure from the same source won the INSERT: the
		// loser's transaction was aborted (postgres), so retry the increment
		// in a fresh transaction — never inside the aborted one.
		return r.recordLoginFailureForOnce(userID, ip, threshold, lockedUntil, now)
	}
	return locked, err
}

func (r *Repository) recordLoginFailureForOnce(userID, ip string, threshold int, lockedUntil time.Time, now time.Time) (bool, error) {
	source := normalizeLoginSource(ip)
	locked := false
	err := r.withTx("record login failure for source", func(tx kernel.Tx) error {
		windowStart := now.Add(-lockCounterWindow).Unix()
		res, err := tx.Exec(context.Background(),
			`UPDATE login_failures SET
			   fail_count = CASE WHEN updated_at < ? THEN 1 ELSE fail_count + 1 END,
			   updated_at = ?
			 WHERE user_id = ? AND ip = ?`,
			windowStart, now.Unix(), userID, source,
		)
		if err != nil {
			return fmt.Errorf("bump source failure: %w", err)
		}
		if affected, rerr := res.RowsAffected(); rerr == nil && affected == 0 {
			if _, ierr := tx.Exec(context.Background(),
				`INSERT INTO login_failures (user_id, ip, fail_count, locked_until, updated_at)
				 VALUES (?, ?, 1, 0, ?)`,
				userID, source, now.Unix(),
			); ierr != nil {
				if kernel.IsUniqueViolation(ierr) {
					return errSourceInsertRace
				}
				return fmt.Errorf("insert source failure: %w", ierr)
			}
			return nil
		}
		var count int
		if err := tx.QueryRow(context.Background(),
			`SELECT fail_count FROM login_failures WHERE user_id = ? AND ip = ?`,
			userID, source,
		).Scan(&count); err != nil {
			return fmt.Errorf("read source failure count: %w", err)
		}
		if count < threshold {
			return nil
		}
		locked = true
		if _, err := tx.Exec(context.Background(),
			`UPDATE login_failures SET locked_until = ?, fail_count = 0, updated_at = ? WHERE user_id = ? AND ip = ?`,
			lockedUntil.Unix(), now.Unix(), userID, source,
		); err != nil {
			return fmt.Errorf("open source lock window: %w", err)
		}
		return nil
	})
	return locked, err
}

// LoginLockedFor reports whether the (account|source) pair sits in an open
// lock window. A storage error returns fail-closed (the caller surfaces it
// rather than treating the source as unlocked).
func (r *Repository) LoginLockedFor(userID, ip string, now time.Time) (bool, error) {
	var lockedUntil int64
	err := r.withTx("read source lock", func(tx kernel.Tx) error {
		return tx.QueryRow(context.Background(),
			`SELECT locked_until FROM login_failures WHERE user_id = ? AND ip = ?`,
			userID, normalizeLoginSource(ip),
		).Scan(&lockedUntil)
	})
	if errors.Is(err, kernel.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("read source lock: %w", err)
	}
	return lockedUntil > now.Unix(), nil
}

// ResetLoginFailuresFor deletes every source-scoped counter/lock row for the
// user (successful login clears all sources at once).
func (r *Repository) ResetLoginFailuresFor(userID string) error {
	return r.withTx("reset source login failures", func(tx kernel.Tx) error {
		if _, err := tx.Exec(context.Background(),
			`DELETE FROM login_failures WHERE user_id = ?`, userID,
		); err != nil {
			return fmt.Errorf("reset source login failures: %w", err)
		}
		return nil
	})
}
