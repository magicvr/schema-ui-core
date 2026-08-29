package kernel

import (
	"errors"
	"fmt"
	"testing"
)

// W9 A-005 R-F-003: regression lock for the dialect-agnostic unique-violation
// predicate (F-001/F-011 foundation). Locks both dialects' message shapes,
// wrapped-error propagation, and the negative case.
func TestIsUniqueViolation(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{"sqlite classic", errors.New("UNIQUE constraint failed: wallet_accounts.owner_id"), true},
		{"sqlite modern", errors.New("constraint failed: UNIQUE constraint failed: users.username (2067)"), true},
		{"postgres pgx", errors.New(`ERROR: duplicate key value violates unique constraint "wallet_accounts_owner_key" (SQLSTATE 23505)`), true},
		{"postgres bare sqlstate", errors.New("duplicate key value violates unique constraint (SQLSTATE 23505)"), true},
		{"wrapped sqlite", fmt.Errorf("insert wallet account: %w", errors.New("UNIQUE constraint failed: x.y")), true},
		{"wrapped postgres", fmt.Errorf("auto-create: %w", errors.New(`pq: duplicate key value violates unique constraint "k" (SQLSTATE 23505)`)), true},
		{"nil", nil, false},
		{"other constraint", errors.New("FOREIGN KEY constraint failed"), false},
		{"plain error", errors.New("disk I/O error"), false},
		{"unique mention without violation", errors.New("index named UNIQUE constraint helper is missing"), false},
	}
	for _, tc := range cases {
		if got := IsUniqueViolation(tc.err); got != tc.want {
			t.Fatalf("%s: IsUniqueViolation = %v, want %v", tc.name, got, tc.want)
		}
	}
}
