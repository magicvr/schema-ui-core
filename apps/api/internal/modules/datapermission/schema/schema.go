// Package schema owns the admin.data-permission page document (S-09 ·
// GOAL-016 D-002 §3).
package schema

import "embed"

// ModuleID is the owning module for the schema handler's contribution-driven
// page gating.
const ModuleID = "admin.data-permission"

//go:embed data-permission.json
var schemaFiles embed.FS

// PageIDs are the page identifiers this module contributes.
func PageIDs() []string { return []string{"data-permission"} }

// SchemaDocuments returns the page documents owned by the data-permission
// module.
func SchemaDocuments() map[string][]byte {
	return map[string][]byte{
		"data-permission": mustRead("data-permission.json"),
	}
}

func mustRead(path string) []byte {
	raw, err := schemaFiles.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return raw
}
