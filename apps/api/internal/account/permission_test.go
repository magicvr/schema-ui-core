package account

import "testing"

func TestEvaluate(t *testing.T) {
	user := User{ID: "u1", Name: "A", Roles: []string{"admin", "editor"}}
	features := map[string]bool{"beta": true}

	cases := []struct {
		name string
		expr string
		want bool
	}{
		{"user string equality", `$context.user.name == "A"`, true},
		{"user string inequality", `$context.user.name == "B"`, false},
		{"roles contains", `$context.user.roles contains "admin"`, true},
		{"roles missing", `$context.user.roles contains "owner"`, false},
		{"features boolean", `$context.features.beta == true`, true},
		{"features false", `$context.features.beta == false`, false},
		{"nested user path", `$context.user.profile.admin == true`, false},
		{"undeclared path fails closed", `$context.user.roles contains "x"`, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Evaluate(tc.expr, user, features)
			if err != nil {
				t.Fatalf("Evaluate(%q) error: %v", tc.expr, err)
			}
			if got != tc.want {
				t.Fatalf("Evaluate(%q) = %v, want %v", tc.expr, got, tc.want)
			}
		})
	}
}

func TestEvaluateInvalid(t *testing.T) {
	for _, expr := range []string{
		`$deps.user.roles contains "admin"`,
		`$context.user.roles`,
		`$context.user.roles > 3`,
		"",
	} {
		if _, err := Evaluate(expr, User{}, nil); err == nil {
			t.Fatalf("Evaluate(%q) expected error, got nil", expr)
		}
	}
}

func TestAllowFailClosed(t *testing.T) {
	user := User{Roles: []string{"viewer"}}
	if Allow(`$context.user.roles contains "admin"`, user, nil) {
		t.Fatal("Allow should deny non-admin")
	}
	if !Allow("", user, nil) {
		t.Fatal("Allow should pass on empty expression")
	}
	if Allow(`$context.user.roles contains "admin" AND false`, user, nil) {
		t.Fatal("Allow should fail closed on unparsable expression")
	}
}
