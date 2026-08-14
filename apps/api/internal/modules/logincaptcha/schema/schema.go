// Package schema owns the admin.login-captcha page document (S-11 ·
// GOAL-011 D-002 §4): the captcha settings page.
package schema

import "embed"

// ModuleID is the owning module for the schema handler's contribution-driven
// page gating.
const ModuleID = "admin.login-captcha"

//go:embed captcha.json
var schemaFiles embed.FS

// PageIDs are the page identifiers this module contributes.
func PageIDs() []string { return []string{"captcha"} }

// SchemaDocuments returns the page documents owned by the captcha module.
func SchemaDocuments() map[string][]byte {
	return map[string][]byte{
		"captcha": mustRead("captcha.json"),
	}
}

func mustRead(path string) []byte {
	raw, err := schemaFiles.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return raw
}
