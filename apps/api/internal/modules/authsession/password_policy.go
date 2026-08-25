// Password policy domain (workspace-019 R3 · GOAL-004 D-001 §2): the frozen
// I-003 起步宽松 defaults (minLength 8 = current baseline, complexity off,
// history off) with admin-tunable knobs, enforced at ALL FOUR set-password
// moments. 渐进生效 (I-007): no stock scans, no forced logout — policy bites
// only when a new password is set.
package authsession

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	"golang.org/x/crypto/bcrypt"
)

// Policy bounds (D-001 §2): minLength stays inside the bcrypt-safe 8–72 byte
// window; categories counts among lower/upper/digit/other; depth caps the
// history look-back so the table stays bounded.
const (
	policyMinLengthFloor   = 8
	policyMinLengthCeiling = 72
	policyMaxCategories    = 4
	policyMaxHistoryDepth  = 10
)

// ErrPasswordPolicyViolation maps to INVALID_PASSWORD on every HTTP surface;
// the granular reason stays in the service log, never in the response.
var ErrPasswordPolicyViolation = fmt.Errorf("authsession: password violates the active policy")

// PasswordPolicy is one row of the singleton table.
type PasswordPolicy struct {
	MinLength     int
	MinCategories int
	HistoryDepth  int
}

func randomHexID(prefix string) string {
	raw := make([]byte, 12)
	_, _ = rand.Read(raw)
	return prefix + "-" + hex.EncodeToString(raw)
}

// GetPasswordPolicy reads the singleton row; a missing row (pre-0057 store in
// exotic test setups) degrades to the frozen defaults instead of failing open.
func (r *Repository) GetPasswordPolicy() (PasswordPolicy, error) {
	var p PasswordPolicy
	err := r.withTx("read password policy", func(tx kernel.Tx) error {
		return tx.QueryRow(context.Background(),
			`SELECT min_length, min_categories, history_depth FROM password_policy WHERE id = 1`,
		).Scan(&p.MinLength, &p.MinCategories, &p.HistoryDepth)
	})
	if errors.Is(err, kernel.ErrNoRows) {
		return PasswordPolicy{MinLength: 8}, nil // unseeded store: frozen defaults
	}
	if err != nil {
		return p, err // real storage failure: caller fails closed
	}
	return p, nil
}

// UpdatePasswordPolicy writes the admin-tuned values after handler-side range
// validation (D-001 §2: length ∈[8,72], categories ∈[0,4], depth ∈[0,10]).
func (r *Repository) UpdatePasswordPolicy(p PasswordPolicy) error {
	return r.withTx("update password policy", func(tx kernel.Tx) error {
		_, err := tx.Exec(context.Background(),
			`UPDATE password_policy SET min_length = ?, min_categories = ?, history_depth = ? WHERE id = 1`,
			p.MinLength, p.MinCategories, p.HistoryDepth)
		return err
	})
}

// ValidateNewPassword enforces the active policy for a candidate plaintext.
// userID may be empty (account creation — no history exists yet). Every
// violation collapses into ErrPasswordPolicyViolation. The CONFIGURED
// MinLength is authoritative (clamped to the 8-byte floor) — A-001 F-001:
// a tightened minimum must actually bite at all four set-password points.
func (r *Repository) ValidateNewPassword(userID, plain string) error {
	policy, err := r.GetPasswordPolicy()
	if err != nil {
		return err
	}
	minLength := policy.MinLength
	if minLength < policyMinLengthFloor {
		minLength = policyMinLengthFloor // never below the bcrypt-safe baseline
	}
	length := len([]byte(plain))
	if length < minLength || length > policyMinLengthCeiling || strings.TrimSpace(plain) == "" {
		return ErrPasswordPolicyViolation
	}
	if policy.MinCategories > 0 && countCategories(plain) < policy.MinCategories {
		return ErrPasswordPolicyViolation
	}
	if policy.HistoryDepth > 0 && userID != "" {
		reused, err := r.passwordInHistory(userID, plain, policy.HistoryDepth)
		if err != nil {
			return err
		}
		if reused {
			return ErrPasswordPolicyViolation
		}
	}
	return nil
}

// passwordInHistory bcrypt-compares the candidate against the newest depth
// historical hashes. bcrypt comparison is the ONLY safe check — hashes are
// salted, so equality must go through CompareHashAndPassword.
func (r *Repository) passwordInHistory(userID, plain string, depth int) (bool, error) {
	var reused bool
	err := r.withTx("check password history", func(tx kernel.Tx) error {
		rows, err := tx.Query(context.Background(),
			`SELECT password_hash FROM user_password_history WHERE user_id = ? ORDER BY created_at DESC, id DESC LIMIT ?`,
			userID, depth)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var hash string
			if err := rows.Scan(&hash); err != nil {
				return err
			}
			if bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil {
				reused = true
				return nil
			}
		}
		return rows.Err()
	})
	return reused, err
}

// capturePasswordHistory records the OUTGOING hash before a replacement lands
// and trims the per-user list to the configured depth. Called from UpdateUser
// INSIDE its transaction whenever a PasswordHash patch is applied; a missing
// row (policy table absent) is tolerated so legacy stores keep working.
func (r *Repository) capturePasswordHistory(tx kernel.Tx, userID, oldHash string, depth int, now int64) {
	if oldHash == "" || depth <= 0 {
		return
	}
	if _, err := tx.Exec(context.Background(),
		`INSERT INTO user_password_history (id, user_id, password_hash, created_at) VALUES (?, ?, ?, ?)`,
		randomHexID("uph"), userID, oldHash, now); err != nil {
		return // best-effort: history strengthens policy, never blocks a set
	}
	_, _ = tx.Exec(context.Background(),
		`DELETE FROM user_password_history WHERE user_id = ? AND id NOT IN (
		   SELECT id FROM user_password_history WHERE user_id = ? ORDER BY created_at DESC, id DESC LIMIT ?
		 )`, userID, userID, depth)
}

func countCategories(plain string) int {
	var lower, upper, digit, other bool
	for _, c := range plain {
		switch {
		case c >= 'a' && c <= 'z':
			lower = true
		case c >= 'A' && c <= 'Z':
			upper = true
		case c >= '0' && c <= '9':
			digit = true
		default:
			other = true
		}
	}
	count := 0
	for _, set := range []bool{lower, upper, digit, other} {
		if set {
			count++
		}
	}
	return count
}
