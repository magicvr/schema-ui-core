package operationlog

import (
	"strings"
	"testing"
)

func TestNewDetailVersionAndDiff(t *testing.T) {
	raw, err := NewDetail("updated", map[string]any{"siteTitle": "Old", "logoUrl": "/assets/old.svg"}, map[string]any{
		"siteTitle": "New",
		"password":  "do-not-store",
		"logoUrl":   "/assets/private.svg",
	})
	if err != nil {
		t.Fatalf("NewDetail: %v", err)
	}
	envelope, err := ParseDetail(raw)
	if err != nil {
		t.Fatalf("ParseDetail: %v (%s)", err, raw)
	}
	if envelope.SchemaVersion != DetailSchemaVersion || envelope.Action != "updated" {
		t.Fatalf("envelope header = %+v", envelope)
	}
	if envelope.Diff["siteTitle"].Before != "Old" || envelope.Diff["siteTitle"].After != "New" {
		t.Fatalf("siteTitle diff = %+v", envelope.Diff["siteTitle"])
	}
	if envelope.After["password"] != RedactedValue || envelope.After["logoUrl"] != RedactedValue {
		t.Fatalf("sensitive values were not redacted: %+v", envelope.After)
	}
	if envelope.Diff["logoUrl"].Before != RedactedValue || envelope.Diff["logoUrl"].After != RedactedValue {
		t.Fatalf("redacted changed field missing from diff: %+v", envelope.Diff["logoUrl"])
	}
	if strings.Contains(raw, "do-not-store") || strings.Contains(raw, "private.svg") {
		t.Fatalf("raw detail contains sensitive value: %s", raw)
	}
}

func TestParseDetailRejectsLegacy(t *testing.T) {
	if _, err := ParseDetail(`{"username":"alice"}`); err == nil {
		t.Fatal("ParseDetail accepted legacy detail")
	}
}

func TestNewDetailRedactsNestedSensitiveValues(t *testing.T) {
	raw, err := NewDetail("updated", nil, map[string]any{
		"sessionToken": "session-secret",
		"idToken":      "id-secret",
		"apiToken":     "api-secret",
		"tokenVersion": "v1",
		"nested": map[string]any{
			"secretBase32":  "mfa-secret",
			"recoveryCodes": []any{"recovery-secret"},
			"otpauthURL":    "otpauth://secret",
		},
		"items": []any{map[string]any{"password": "nested-password"}},
	})
	if err != nil {
		t.Fatalf("NewDetail: %v", err)
	}
	envelope, err := ParseDetail(raw)
	if err != nil {
		t.Fatalf("ParseDetail: %v", err)
	}
	for _, key := range []string{"sessionToken", "idToken", "apiToken"} {
		if envelope.After[key] != RedactedValue {
			t.Fatalf("%s = %v, want redacted", key, envelope.After[key])
		}
	}
	if envelope.After["tokenVersion"] != "v1" {
		t.Fatalf("tokenVersion = %v, want preserved", envelope.After["tokenVersion"])
	}
	nested := envelope.After["nested"].(map[string]any)
	for _, key := range []string{"secretBase32", "recoveryCodes", "otpauthURL"} {
		if nested[key] == nil || nested[key] == "" || nested[key] == "otpauth://secret" {
			t.Fatalf("nested %s = %v, want redacted", key, nested[key])
		}
	}
	items := envelope.After["items"].([]any)
	if items[0].(map[string]any)["password"] != RedactedValue {
		t.Fatalf("array password = %v, want redacted", items[0])
	}
	if strings.Contains(raw, "session-secret") || strings.Contains(raw, "nested-password") {
		t.Fatalf("raw detail contains sensitive value: %s", raw)
	}
}
