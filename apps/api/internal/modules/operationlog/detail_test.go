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
