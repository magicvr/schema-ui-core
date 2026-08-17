package account

// Session is the R4 account session. With real auth (R2) the User is resolved
// from the request identity; Features is always emitted so the renderer's
// $context.fallback treats absence as "no flags" rather than "unknown".
type Session struct {
	User     User            `json:"user"`
	Features map[string]bool `json:"features"`
}

// User is the $context.user snapshot consumed by the renderer.
type User struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Roles []string `json:"roles"`
	// AvatarURL is the self-service avatar asset URL (W13 T-05); "" = none.
	AvatarURL string `json:"avatarUrl,omitempty"`
	// MustChangePassword is true when the user must set a new password before
	// accessing business APIs (W16-F01).
	MustChangePassword bool `json:"mustChangePassword"`
	// Permissions are the persisted role-permission keys resolved at identity
	// load time (GOAL-006 S4): business gates check keys, not role strings.
	Permissions []string `json:"permissions,omitempty"`
}

// StaticDevSession returns the development bootstrap session. Fail-closed:
// any change to this function must keep the session explicit and auditable,
// never empty or accidentally permissive.
func StaticDevSession() Session {
	return Session{
		User: User{
			ID:   "dev-001",
			Name: "Dev Admin",
			// editor + admin unlock every permission-inheritance fixture
			// scenario used by the R4 conformance tests.
			Roles: []string{"admin", "editor"},
			// Consistent with the admin seed's grants under the S4 gate: the
			// GOAL-011 users/roles keys and current product grants.
			Permissions: []string{
				"users.read", "users.write",
				"roles.read", "roles.write", "roles.assign",
				"settings.read", "settings.write",
				"operations.read",
				"files.write",                   // W4 P0-2: admin-only upload gate, parity with seed
				"users.enable", "users.disable", // F-03 (GOAL-005) admin-only keys
				"data.export", "data.import", // F-02 (GOAL-004) transfer keys
			},
		},
		Features: map[string]bool{
			"beta":           true,
			"menu_users":     true, // GOAL-011 users page grant parity
			"menu_roles":     true, // GOAL-011 roles page grant parity
			"menu_settings":  true, // GOAL-013 site branding
			"menu_activity":  true, // GOAL-013 operation log
			"menu_account":   true, // F-03 (GOAL-005) self-service account page
			"menu_dashboard": true, // F-01 (GOAL-003) production home dashboard
		},
	}
}
