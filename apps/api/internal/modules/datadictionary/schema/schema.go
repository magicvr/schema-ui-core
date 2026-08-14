// Package schema owns the admin.data-dictionary page documents (S-01 ·
// GOAL-008 D-002 §4): the types page and the entries page.
package schema

import "embed"

// ModuleID is the owning module for the schema handler's contribution-driven
// page gating.
const ModuleID = "admin.data-dictionary"

//go:embed data-dictionary.json dictionary-entries.json
var schemaFiles embed.FS

// PageIDs are the page identifiers this module contributes.
func PageIDs() []string { return []string{"data-dictionary", "dictionary-entries"} }

// SchemaDocuments returns the page documents owned by the dictionary module.
func SchemaDocuments() map[string][]byte {
	return map[string][]byte{
		"data-dictionary":     mustRead("data-dictionary.json"),
		"dictionary-entries":  mustRead("dictionary-entries.json"),
	}
}

func mustRead(path string) []byte {
	raw, err := schemaFiles.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return raw
}
