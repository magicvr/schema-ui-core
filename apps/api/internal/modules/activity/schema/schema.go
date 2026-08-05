package schema

import "embed"

// ModuleID is the owning module for the schema handler's contribution-driven
// page gating (R4 C4.3).
const ModuleID = "admin.activity"

// PageIDs are the page identifiers this module contributes.
func PageIDs() []string { return []string{"activity"} }

//go:embed *.json
var schemaFiles embed.FS

// SchemaDocuments returns the page documents owned by the Activity module.
// The embedded file set is a build-time invariant; a missing document is a
// programming error and fails construction rather than becoming a runtime
// fallback.
func SchemaDocuments() map[string][]byte {
	return map[string][]byte{
		"activity": mustRead("activity.json"),
	}
}

func mustRead(path string) []byte {
	raw, err := schemaFiles.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return raw
}
