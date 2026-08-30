package kernel

import (
	"errors"
	"fmt"
	"testing"
)

func TestIsDuplicateObject(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{"postgres table", errors.New(`ERROR: relation "users" already exists (SQLSTATE 42P07)`), true},
		{"postgres column", errors.New(`ERROR: column "token_version" of relation "users" already exists (SQLSTATE 42701)`), true},
		{"postgres object", errors.New(`ERROR: constraint "users_pkey" already exists (SQLSTATE 42710)`), true},
		{"wrapped", fmt.Errorf("create baseline (postgres): %w", errors.New(`ERROR: relation "users" already exists (SQLSTATE 42P07)`)), true},
		{"nil", nil, false},
		{"unique violation is not duplicate object", errors.New(`duplicate key value violates unique constraint "users_username_key" (SQLSTATE 23505)`), false},
		{"plain error", errors.New("disk I/O error"), false},
	}
	for _, tc := range cases {
		if got := IsDuplicateObject(tc.err); got != tc.want {
			t.Fatalf("%s: IsDuplicateObject = %v, want %v", tc.name, got, tc.want)
		}
	}
}
