// Package schema owns the admin.users page document (R4 C3.3 content
// migration). The embedded file is the same fixture previously served from the
// central handler embed; ownership moves to the module so the schema surface is
// module-contributed.
package schema

import "embed"

//go:embed users.json
var schemaFiles embed.FS

// SchemaDocuments returns the page documents owned by the users module.
func SchemaDocuments() map[string][]byte {
	return map[string][]byte{"users": mustRead("users.json")}
}

func mustRead(path string) []byte {
	raw, err := schemaFiles.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return raw
}
