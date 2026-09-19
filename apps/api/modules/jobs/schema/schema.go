// Package schema owns the admin.jobs page documents (GOAL-003 R2 · D-001 §4):
// the management-scope job list. The result-center experience lands with R4.
package schema

import "embed"

// ModuleID is the owning module for the schema handler's contribution-driven
// page gating.
const ModuleID = "admin.jobs"

//go:embed jobs.json
var schemaFiles embed.FS

// PageIDs are the page identifiers this module contributes.
func PageIDs() []string { return []string{"jobs"} }

// SchemaDocuments returns the page documents owned by the jobs module.
func SchemaDocuments() map[string][]byte {
	return map[string][]byte{
		"jobs": mustRead("jobs.json"),
	}
}

func mustRead(path string) []byte {
	raw, err := schemaFiles.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return raw
}
