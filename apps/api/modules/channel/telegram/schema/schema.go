package schema

import "embed"

// ModuleID is the owning module for the schema handler's contribution-driven
// page gating.
const ModuleID = "channel.telegram"

// PageIDs are the page identifiers this module contributes. GOAL-006 R5
// (判据 #5 补做 Admin UI tab) adds the standalone telegram-settings page
// hosted by the telegram-admin-tab custom component.
func PageIDs() []string { return []string{"telegram-settings"} }

//go:embed *.json
var schemaFiles embed.FS

// SchemaDocuments returns the page documents owned by the channel.telegram
// module. The embedded file set is a build-time invariant; a missing document
// is a programming error and fails construction rather than becoming a
// runtime fallback.
func SchemaDocuments() map[string][]byte {
	return map[string][]byte{
		"telegram-settings": mustRead("telegram-settings.json"),
	}
}

func mustRead(path string) []byte {
	raw, err := schemaFiles.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return raw
}
