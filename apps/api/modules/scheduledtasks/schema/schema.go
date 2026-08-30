// Package schema owns the admin.scheduled-tasks page documents (S-04 ·
// GOAL-010 D-002 §5): the tasks page and the global run-history page.
package schema

import "embed"

// ModuleID is the owning module for the schema handler's contribution-driven
// page gating.
const ModuleID = "admin.scheduled-tasks"

//go:embed scheduled-tasks.json task-runs.json
var schemaFiles embed.FS

// PageIDs are the page identifiers this module contributes.
func PageIDs() []string { return []string{"scheduled-tasks", "task-runs"} }

// SchemaDocuments returns the page documents owned by the tasks module.
func SchemaDocuments() map[string][]byte {
	return map[string][]byte{
		"scheduled-tasks": mustRead("scheduled-tasks.json"),
		"task-runs":       mustRead("task-runs.json"),
	}
}

func mustRead(path string) []byte {
	raw, err := schemaFiles.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return raw
}
