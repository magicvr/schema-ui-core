package config

import "testing"

// TestValidateProd covers the production guard added in response to GOAL-008
// A-005 F-002: the static dev session fallback is a local-development-only
// feature and must fail startup in any non-development environment.
func TestValidateProd(t *testing.T) {
	t.Run("development may enable dev session", func(t *testing.T) {
		c := &Config{AppEnv: "development", AuthDevSessionEnabled: true}
		if err := c.ValidateProd(); err != nil {
			t.Fatalf("development + dev session should pass, got: %v", err)
		}
	})

	t.Run("production with dev session fails closed", func(t *testing.T) {
		c := &Config{AppEnv: "production", AuthDevSessionEnabled: true}
		if err := c.ValidateProd(); err == nil {
			t.Fatal("production + dev session must be a startup error")
		}
	})

	t.Run("production without dev session passes", func(t *testing.T) {
		c := &Config{AppEnv: "production", AuthDevSessionEnabled: false}
		if err := c.ValidateProd(); err != nil {
			t.Fatalf("production without dev session should pass, got: %v", err)
		}
	})

	t.Run("non-development non-production env also fails closed", func(t *testing.T) {
		c := &Config{AppEnv: "staging", AuthDevSessionEnabled: true}
		if err := c.ValidateProd(); err == nil {
			t.Fatal("staging + dev session must be a startup error")
		}
	})
}
