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
		},
		Features: map[string]bool{
			"beta": true,
		},
	}
}
