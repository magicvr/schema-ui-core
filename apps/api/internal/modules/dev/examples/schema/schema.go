// Package schema owns the dev.examples demonstration page documents.
// W1 (GOAL-002 / workspace-010): these pages moved out of core.schema-render
// into the optional dev.examples module so production profiles can omit the
// demo surface while dev/dogfood can enable it explicitly (D-003 §3).
package schema

import "embed"

const ModuleID = "dev.examples"

var pageIDs = []string{
	"admin-list-batch",
	"data-display",
	"data-table",
	"form-controls",
	"form-with-reactions",
	"form-with-upload",
	"overview",
	"search-form-table",
}

// PageIDs returns the stable example page identifiers in deterministic order.
func PageIDs() []string { return append([]string(nil), pageIDs...) }

//go:embed *.json
var schemaFiles embed.FS

// SchemaDocuments returns fresh byte slices for every example page document.
func SchemaDocuments() map[string][]byte {
	documents := make(map[string][]byte, len(pageIDs))
	for _, pageID := range pageIDs {
		raw, err := schemaFiles.ReadFile(pageID + ".json")
		if err != nil {
			panic(err)
		}
		documents[pageID] = raw
	}
	return documents
}
