package authsession

import (
	"errors"
	"testing"
)

// The singleton policy row is seeded by migration 0057; a legacy pre-0057
// store must NOT let UpdatePasswordPolicy silently no-op (A-001 F-001):
// zero affected rows fail closed with ErrPasswordPolicyNotSeeded.
func TestUpdatePasswordPolicyNotSeededFailsClosed(t *testing.T) {
	repository, st := openRepository(t, "policy-unseeded.db", false)
	if err := repositoryExec(t, st, `DELETE FROM password_policy`); err != nil {
		t.Fatalf("clear policy row: %v", err)
	}
	err := repository.UpdatePasswordPolicy(PasswordPolicy{MinLength: 12, MinCategories: 2})
	if !errors.Is(err, ErrPasswordPolicyNotSeeded) {
		t.Fatalf("UpdatePasswordPolicy on unseeded row = %v, want ErrPasswordPolicyNotSeeded", err)
	}
}

func TestUpdatePasswordPolicyPersistsOnSeededRow(t *testing.T) {
	repository, _ := openRepository(t, "policy-seeded.db", false)
	if err := repository.UpdatePasswordPolicy(PasswordPolicy{MinLength: 12, MinCategories: 2, HistoryDepth: 3}); err != nil {
		t.Fatalf("UpdatePasswordPolicy: %v", err)
	}
	got, err := repository.GetPasswordPolicy()
	if err != nil {
		t.Fatalf("GetPasswordPolicy: %v", err)
	}
	if got.MinLength != 12 || got.MinCategories != 2 || got.HistoryDepth != 3 {
		t.Fatalf("got = %+v, want {12 2 3}", got)
	}
}
