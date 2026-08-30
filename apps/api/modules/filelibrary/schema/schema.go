// Package schema owns the admin.file-library page document (S-02 · GOAL-007
// D-002 §5). The embedded file is the schema-driven library page: list,
// row download (custom handler library.download) and delete, and the upload
// toolbar modal (upload field → central POST /api/upload, submit →
// POST /api/library/files/upload confirmation).
package schema

import "embed"

// ModuleID is the owning module for the schema handler's contribution-driven
// page gating.
const ModuleID = "admin.file-library"

//go:embed file-library.json
var schemaFiles embed.FS

// PageIDs are the page identifiers this module contributes.
func PageIDs() []string { return []string{"file-library"} }

// SchemaDocuments returns the page documents owned by the file-library module.
func SchemaDocuments() map[string][]byte {
	return map[string][]byte{"file-library": mustRead("file-library.json")}
}

func mustRead(path string) []byte {
	raw, err := schemaFiles.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return raw
}
