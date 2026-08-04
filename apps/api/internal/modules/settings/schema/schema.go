package schema

import "embed"

//go:embed *.json
var schemaFiles embed.FS

// SchemaDocuments returns the page documents owned by the Settings module.
// The embedded file set is a build-time invariant; a missing document is a
// programming error and fails construction rather than becoming a runtime
// fallback.
func SchemaDocuments() map[string][]byte {
	return map[string][]byte{
		"settings": mustRead("settings.json"),
	}
}

func mustRead(path string) []byte {
	raw, err := schemaFiles.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return raw
}
