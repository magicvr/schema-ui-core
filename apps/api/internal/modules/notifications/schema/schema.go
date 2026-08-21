package schema

import "embed"

// ModuleID is the owning module for the schema handler's contribution-driven
// page gating (R4 C4.3).
const ModuleID = "admin.notifications"

// PageIDs are the page identifiers this module contributes.
func PageIDs() []string { return []string{"notifications"} }

//go:embed *.json
var schemaFiles embed.FS

// SchemaDocuments returns the page documents owned by the module.
func SchemaDocuments() map[string][]byte {
	return map[string][]byte{
		"notifications": mustRead("notifications.json"),
	}
}

func mustRead(path string) []byte {
	raw, err := schemaFiles.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return raw
}
