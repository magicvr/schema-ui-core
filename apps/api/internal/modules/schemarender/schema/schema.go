// Package schema owns the core.schema-render page documents.
package schema

import "embed"

const ModuleID = "core.schema-render"

var pageIDs = []string{
	"data-table",
	"form-controls",
	"form-with-reactions",
	"overview",
	"search-form-table",
}

// PageIDs returns the stable core page identifiers in deterministic order.
func PageIDs() []string { return append([]string(nil), pageIDs...) }

//go:embed *.json
var schemaFiles embed.FS

// SchemaDocuments returns fresh byte slices for every core page document.
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
