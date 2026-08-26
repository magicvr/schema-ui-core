package schema

import "embed"

// ModuleID is the owning module for the schema handler's contribution-driven
// page gating (R4 C4.3).
const ModuleID = "admin.settings"

// PageIDs are the page identifiers this module contributes. W26 (GOAL-038
// D-001 §2.2) adds the standalone mail console + outbound log pages; the
// settings page no longer hosts either block.
func PageIDs() []string { return []string{"settings", "mail", "mail-outbox"} }

//go:embed *.json
var schemaFiles embed.FS

// SchemaDocuments returns the page documents owned by the Settings module.
// The embedded file set is a build-time invariant; a missing document is a
// programming error and fails construction rather than becoming a runtime
// fallback.
func SchemaDocuments() map[string][]byte {
	return map[string][]byte{
		"settings":   mustRead("settings.json"),
		"mail":       mustRead("mail.json"),
		"mail-outbox": mustRead("mail-outbox.json"),
	}
}

func mustRead(path string) []byte {
	raw, err := schemaFiles.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return raw
}
