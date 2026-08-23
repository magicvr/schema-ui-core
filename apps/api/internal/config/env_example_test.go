package config

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

var envKeyLit = regexp.MustCompile(`(?:envOr|durationEnv|boolEnv|positiveIntEnv|nonNegIntEnv|os\.Getenv|os\.LookupEnv)\("([A-Z][A-Z0-9_]*)"`)

var exampleKeyLine = regexp.MustCompile(`^\s*#?\s*([A-Z][A-Z0-9_]*)=`)

// loadEnvFile is selected from the process environment before this file is
// read, so CONFIG_ENV_FILE cannot be configured from configs/.env itself.
var envFileSelfKeys = map[string]struct{}{
	"CONFIG_ENV_FILE": {},
}

// pgtest extras loaded from the same configs/.env (PG_TEST_* only).
var pgtestKeys = map[string]struct{}{
	"PG_TEST_DSN":      {},
	"PG_TEST_HOST":     {},
	"PG_TEST_PORT":     {},
	"PG_TEST_USER":     {},
	"PG_TEST_PASSWORD": {},
	"PG_TEST_DB":       {},
	"PG_TEST_SSLMODE":  {},
}

func TestCanonicalEnvExample(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller")
	}
	apiRoot := filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", ".."))
	configsDir := filepath.Join(apiRoot, "configs")
	canonical := filepath.Join(configsDir, ".env.example")

	raw, err := os.ReadFile(canonical)
	if err != nil {
		t.Fatalf("canonical env template missing: %s (%v)", canonical, err)
	}

	documented := exampleKeys(string(raw))
	required := loadEnvKeys(t, filepath.Join(filepath.Dir(thisFile), "config.go"))
	for k := range envFileSelfKeys {
		delete(required, k)
	}

	for k := range required {
		if _, ok := documented[k]; !ok {
			t.Errorf("configs/.env.example missing %s (config.Load reads it)", k)
		}
	}
	for k := range documented {
		if _, ok := required[k]; ok {
			continue
		}
		if _, ok := pgtestKeys[k]; ok {
			continue
		}
		t.Errorf("configs/.env.example documents unused/obsolete %s", k)
	}
	for k := range pgtestKeys {
		if _, ok := documented[k]; !ok {
			t.Errorf("configs/.env.example missing %s (pgtest reads it from configs/.env)", k)
		}
	}

	entries, err := os.ReadDir(configsDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		name := e.Name()
		if name == ".env.example" {
			continue
		}
		lower := strings.ToLower(name)
		if strings.Contains(lower, "env") && strings.Contains(lower, "example") {
			t.Errorf("duplicate env template %s; only configs/.env.example is canonical", filepath.Join(configsDir, name))
		}
	}

	legacy := filepath.Join(apiRoot, ".env.example")
	if _, err := os.Stat(legacy); err == nil {
		t.Errorf("legacy %s must not exist; use %s", legacy, canonical)
	}
}

func loadEnvKeys(t *testing.T, path string) map[string]struct{} {
	t.Helper()
	src, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	out := map[string]struct{}{}
	for _, m := range envKeyLit.FindAllStringSubmatch(string(src), -1) {
		out[m[1]] = struct{}{}
	}
	if len(out) == 0 {
		t.Fatal("no env key literals found in config.go")
	}
	return out
}

func exampleKeys(body string) map[string]struct{} {
	out := map[string]struct{}{}
	for _, line := range strings.Split(body, "\n") {
		m := exampleKeyLine.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		out[m[1]] = struct{}{}
	}
	return out
}
