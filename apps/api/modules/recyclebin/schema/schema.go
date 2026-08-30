// Package schema owns the admin.recycle-bin page document (S-12 ·
// GOAL-012 D-002 §4): the recycle bin management page.
package schema

import "embed"

// ModuleID is the owning module for the schema handler's contribution-driven
// page gating.
const ModuleID = "admin.recycle-bin"

//go:embed recycle-bin.json
var schemaFiles embed.FS

// PageIDs are the page identifiers this module contributes.
func PageIDs() []string { return []string{"recycle-bin"} }

// SchemaDocuments returns the page documents owned by the recycle module.
func SchemaDocuments() map[string][]byte {
	return map[string][]byte{
		"recycle-bin": mustRead("recycle-bin.json"),
	}
}

func mustRead(path string) []byte {
	raw, err := schemaFiles.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return raw
}
