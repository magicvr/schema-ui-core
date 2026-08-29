// Package schema owns the admin.system-monitoring page document (S-03 ·
// GOAL-009 D-002 §3): a status summary (statCard grid) plus the recent-events
// table over the operation-log surface.
package schema

import "embed"

// ModuleID is the owning module for the schema handler's contribution-driven
// page gating.
const ModuleID = "admin.system-monitoring"

//go:embed system-monitoring.json
var schemaFiles embed.FS

// PageIDs are the page identifiers this module contributes.
func PageIDs() []string { return []string{"system-monitoring"} }

// SchemaDocuments returns the page documents owned by the module.
func SchemaDocuments() map[string][]byte {
	return map[string][]byte{"system-monitoring": mustRead("system-monitoring.json")}
}

func mustRead(path string) []byte {
	raw, err := schemaFiles.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return raw
}
