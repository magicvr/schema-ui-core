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

// W15 F-003 (GOAL-016 A-001): the seed-password gate mirrors the frozen
// bootstrap bounds (migration 0057 seeds min_length 8) at both bootstrap
// surfaces. 8 bytes passes, 7 and 73 fail, blank/whitespace fail.
func TestValidateSeedPassword(t *testing.T) {
	for seed, wantOK := range map[string]bool{
		"":              false,
		"   ":           false,
		"1234567":       false, // 7 bytes
		"12345678":      true,  // exactly the 8-byte floor
		"strong-enough": true,
	} {
		err := ValidateSeedPassword(seed)
		if wantOK && err != nil {
			t.Errorf("ValidateSeedPassword(%q) = %v, want nil", seed, err)
		}
		if !wantOK && !errors.Is(err, ErrPasswordPolicyViolation) {
			t.Errorf("ValidateSeedPassword(%q) = %v, want ErrPasswordPolicyViolation", seed, err)
		}
	}
	// 73 bytes fails, 72 passes (ceiling).
	seventyThree := make([]byte, 73)
	for i := range seventyThree {
		seventyThree[i] = 'a'
	}
	if err := ValidateSeedPassword(string(seventyThree)); !errors.Is(err, ErrPasswordPolicyViolation) {
		t.Errorf("73-byte seed = %v, want violation", err)
	}
	seventyTwo := make([]byte, 72)
	for i := range seventyTwo {
		seventyTwo[i] = 'b'
	}
	if err := ValidateSeedPassword(string(seventyTwo)); err != nil {
		t.Errorf("72-byte seed = %v, want nil", err)
	}
}
