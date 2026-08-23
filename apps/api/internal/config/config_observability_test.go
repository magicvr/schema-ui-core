package config

import (
	"strings"
	"testing"
)

// TestObservabilityMetricsConfig covers the observability.metrics contract
// (VP-015 / workspace-015 GOAL-002 D-001): defaults stay fully off, YAML and
// env layers apply, and the fail-closed exposure/pairing rules hold.
func TestObservabilityMetricsConfig(t *testing.T) {
	t.Run("defaults are fully off", func(t *testing.T) {
		cfg := Load()
		if cfg.LoadError != nil {
			t.Fatalf("LoadError: %v", cfg.LoadError)
		}
		if cfg.MetricsEnabled {
			t.Errorf("MetricsEnabled = true, want default false")
		}
		if cfg.MetricsAddr != "127.0.0.1:25081" {
			t.Errorf("MetricsAddr = %q, want loopback default", cfg.MetricsAddr)
		}
		if cfg.MetricsAuthToken != "" {
			t.Errorf("MetricsAuthToken = %q, want empty", cfg.MetricsAuthToken)
		}
	})

	t.Run("yaml enables loopback scrape without token", func(t *testing.T) {
		writeConfig(t, "app:\n  env: development\nobservability:\n  metrics:\n    enabled: true\n")
		cfg := Load()
		if cfg.LoadError != nil {
			t.Fatalf("LoadError: %v", cfg.LoadError)
		}
		if !cfg.MetricsEnabled || cfg.MetricsAddr != "127.0.0.1:25081" {
			t.Errorf("metrics = enabled=%v addr=%q", cfg.MetricsEnabled, cfg.MetricsAddr)
		}
	})

	t.Run("wildcard bind without token fails closed", func(t *testing.T) {
		writeConfig(t, "app:\n  env: development\nobservability:\n  metrics:\n    enabled: true\n    addr: \":25081\"\n")
		cfg := Load()
		if cfg.LoadError == nil {
			t.Fatal("empty-host bind without token must fail closed")
		}
		if !strings.Contains(cfg.LoadError.Error(), "auth_token") {
			t.Errorf("error must name auth_token, got: %v", cfg.LoadError)
		}
	})

	t.Run("non-loopback IP bind without token fails closed", func(t *testing.T) {
		writeConfig(t, "app:\n  env: development\nobservability:\n  metrics:\n    enabled: true\n    addr: \"192.168.1.10:25081\"\n")
		cfg := Load()
		if cfg.LoadError == nil {
			t.Fatal("non-loopback bind without token must fail closed")
		}
	})

	t.Run("non-loopback bind with interpolated token loads", func(t *testing.T) {
		y := "app:\n  env: development\nobservability:\n  metrics:\n    enabled: true\n    addr: \"0.0.0.0:25081\"\n    auth_token: ${METRICS_TOKEN_TEST:-tok-1234567890abcdef}\n"
		writeConfig(t, y)
		cfg := Load()
		if cfg.LoadError != nil {
			t.Fatalf("LoadError: %v", cfg.LoadError)
		}
		if cfg.MetricsAuthToken != "tok-1234567890abcdef" {
			t.Errorf("MetricsAuthToken not applied via interpolation: %q", cfg.MetricsAuthToken)
		}
	})

	t.Run("dead token fails closed", func(t *testing.T) {
		writeConfig(t, "app:\n  env: development\nobservability:\n  metrics:\n    auth_token: orphan-token-value-123\n")
		cfg := Load()
		if cfg.LoadError == nil {
			t.Fatal("auth_token set while disabled must fail closed")
		}
		if !strings.Contains(cfg.LoadError.Error(), "enabled is false") {
			t.Errorf("unexpected error wording: %v", cfg.LoadError)
		}
	})

	t.Run("short token fails closed", func(t *testing.T) {
		writeConfig(t, "app:\n  env: development\nobservability:\n  metrics:\n    enabled: true\n    auth_token: short\n")
		cfg := Load()
		if cfg.LoadError == nil {
			t.Fatal("token shorter than 16 chars must fail closed")
		}
	})

	t.Run("invalid addr fails closed", func(t *testing.T) {
		writeConfig(t, "app:\n  env: development\nobservability:\n  metrics:\n    enabled: true\n    addr: \"nope\"\n")
		cfg := Load()
		if cfg.LoadError == nil {
			t.Fatal("addr without host:port must fail closed")
		}
	})

	t.Run("non-numeric port fails closed", func(t *testing.T) {
		writeConfig(t, "app:\n  env: development\nobservability:\n  metrics:\n    enabled: true\n    addr: \"127.0.0.1:httpx\"\n")
		cfg := Load()
		if cfg.LoadError == nil {
			t.Fatal("non-numeric port must fail closed")
		}
	})

	t.Run("env overrides win", func(t *testing.T) {
		writeConfig(t, "app:\n  env: development\nobservability:\n  metrics:\n    enabled: false\n")
		t.Setenv("OBSERVABILITY_METRICS_ENABLED", "true")
		t.Setenv("OBSERVABILITY_METRICS_ADDR", "127.0.0.1:25099")
		cfg := Load()
		if cfg.LoadError != nil {
			t.Fatalf("LoadError: %v", cfg.LoadError)
		}
		if !cfg.MetricsEnabled || cfg.MetricsAddr != "127.0.0.1:25099" {
			t.Errorf("env override lost: enabled=%v addr=%q", cfg.MetricsEnabled, cfg.MetricsAddr)
		}
	})

	t.Run("ValidateProd tolerates zero-value config", func(t *testing.T) {
		c := &Config{AppEnv: "development"}
		if err := c.ValidateProd(); err != nil {
			t.Fatalf("zero-value config must skip observability gate, got: %v", err)
		}
	})

	t.Run("ValidateProd re-checks non-loopback exposure", func(t *testing.T) {
		c := &Config{AppEnv: "development", MetricsEnabled: true, MetricsAddr: "0.0.0.0:25081"}
		err := c.ValidateProd()
		if err == nil {
			t.Fatal("hand-built non-loopback config without token must fail ValidateProd")
		}
		if !strings.Contains(err.Error(), "auth_token") {
			t.Errorf("error must name auth_token, got: %v", err)
		}
	})

	t.Run("ValidateProd accepts loopback default on hand-built config", func(t *testing.T) {
		c := &Config{AppEnv: "development", MetricsEnabled: true}
		if err := c.ValidateProd(); err != nil {
			t.Fatalf("loopback default needs no token, got: %v", err)
		}
	})
}

// TestObservabilityTracesConfig covers the observability.traces contract
// (VP-015 / workspace-015 GOAL-004 D-001 §2): defaults stay no-op, YAML and
// env layers apply, and fail-closed pairing rules hold.
func TestObservabilityTracesConfig(t *testing.T) {
	t.Run("defaults are fully no-op", func(t *testing.T) {
		cfg := Load()
		if cfg.LoadError != nil {
			t.Fatalf("LoadError: %v", cfg.LoadError)
		}
		if cfg.TracesEnabled {
			t.Errorf("TracesEnabled = true, want default false")
		}
		if cfg.TracesEndpoint != "" {
			t.Errorf("TracesEndpoint = %q, want empty", cfg.TracesEndpoint)
		}
		if cfg.TracesSampleRatio != 1.0 {
			t.Errorf("TracesSampleRatio = %v, want default 1.0", cfg.TracesSampleRatio)
		}
	})

	t.Run("yaml enables with loopback collector endpoint", func(t *testing.T) {
		y := "app:\n  env: development\nobservability:\n  traces:\n    enabled: true\n    endpoint: http://localhost:4318\n"
		writeConfig(t, y)
		cfg := Load()
		if cfg.LoadError != nil {
			t.Fatalf("LoadError: %v", cfg.LoadError)
		}
		if !cfg.TracesEnabled || cfg.TracesEndpoint != "http://localhost:4318" {
			t.Errorf("traces = enabled=%v endpoint=%q", cfg.TracesEnabled, cfg.TracesEndpoint)
		}
	})

	t.Run("enabled without endpoint fails closed", func(t *testing.T) {
		writeConfig(t, "app:\n  env: development\nobservability:\n  traces:\n    enabled: true\n")
		cfg := Load()
		if cfg.LoadError == nil {
			t.Fatal("enabled traces without endpoint must fail closed")
		}
		if !strings.Contains(cfg.LoadError.Error(), "endpoint is required") {
			t.Errorf("unexpected error wording: %v", cfg.LoadError)
		}
	})

	t.Run("non-http endpoint scheme fails closed", func(t *testing.T) {
		y := "app:\n  env: development\nobservability:\n  traces:\n    enabled: true\n    endpoint: ftp://files:21\n"
		writeConfig(t, y)
		cfg := Load()
		if cfg.LoadError == nil {
			t.Fatal("ftp endpoint must fail closed")
		}
	})

	t.Run("dead endpoint fails closed", func(t *testing.T) {
		y := "app:\n  env: development\nobservability:\n  traces:\n    endpoint: http://localhost:4318\n"
		writeConfig(t, y)
		cfg := Load()
		if cfg.LoadError == nil {
			t.Fatal("endpoint set while disabled must fail closed")
		}
		if !strings.Contains(cfg.LoadError.Error(), "enabled is false") {
			t.Errorf("unexpected error wording: %v", cfg.LoadError)
		}
	})

	t.Run("out-of-range sample ratio fails closed", func(t *testing.T) {
		y := "app:\n  env: development\nobservability:\n  traces:\n    enabled: true\n    endpoint: http://localhost:4318\n    sample_ratio: 1.5\n"
		writeConfig(t, y)
		cfg := Load()
		if cfg.LoadError == nil {
			t.Fatal("sample_ratio > 1 must fail closed")
		}
	})

	t.Run("invalid sample ratio env keeps default", func(t *testing.T) {
		writeConfig(t, "app:\n  env: development\n")
		t.Setenv("OBSERVABILITY_TRACES_SAMPLE_RATIO", "not-a-number")
		cfg := Load()
		if cfg.LoadError != nil {
			t.Fatalf("LoadError: %v", cfg.LoadError)
		}
		if cfg.TracesSampleRatio != 1.0 {
			t.Errorf("TracesSampleRatio = %v, want default on invalid env", cfg.TracesSampleRatio)
		}
	})

	t.Run("env overrides win", func(t *testing.T) {
		writeConfig(t, "app:\n  env: development\n")
		t.Setenv("OBSERVABILITY_TRACES_ENABLED", "true")
		t.Setenv("OBSERVABILITY_TRACES_ENDPOINT", "http://127.0.0.1:4318")
		t.Setenv("OBSERVABILITY_TRACES_SAMPLE_RATIO", "0.5")
		cfg := Load()
		if cfg.LoadError != nil {
			t.Fatalf("LoadError: %v", cfg.LoadError)
		}
		if !cfg.TracesEnabled || cfg.TracesEndpoint != "http://127.0.0.1:4318" || cfg.TracesSampleRatio != 0.5 {
			t.Errorf("env trace config lost: %+v", cfg)
		}
	})

	t.Run("ValidateProd tolerates zero-value config", func(t *testing.T) {
		c := &Config{AppEnv: "development"}
		if err := c.ValidateProd(); err != nil {
			t.Fatalf("zero-value config must skip observability gate, got: %v", err)
		}
	})

	t.Run("ValidateProd re-checks enabled-without-endpoint", func(t *testing.T) {
		c := &Config{AppEnv: "development", TracesEnabled: true}
		err := c.ValidateProd()
		if err == nil {
			t.Fatal("hand-built enabled traces without endpoint must fail ValidateProd")
		}
	})
}
