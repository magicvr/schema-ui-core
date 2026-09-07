// Package schema owns the biz.digital-offer page documents (VP-031 · GOAL-002
// D-002 v1.0.0 §7): the offers management page, the purchases (digital orders)
// read-only view, and the entitlements page.
package schema

import "embed"

// ModuleID is the owning module for the schema handler's contribution-driven
// page gating.
const ModuleID = "biz.digital-offer"

//go:embed digitaloffer-offers.json digitaloffer-entitlements.json digitaloffer-purchases.json
var schemaFiles embed.FS

// PageIDs are the page identifiers this module contributes.
func PageIDs() []string {
	return []string{"digitaloffer-offers", "digitaloffer-entitlements", "digitaloffer-purchases"}
}

// SchemaDocuments returns the page documents owned by the digital-offer module.
func SchemaDocuments() map[string][]byte {
	return map[string][]byte{
		"digitaloffer-offers":       mustRead("digitaloffer-offers.json"),
		"digitaloffer-entitlements": mustRead("digitaloffer-entitlements.json"),
		"digitaloffer-purchases":    mustRead("digitaloffer-purchases.json"),
	}
}

func mustRead(path string) []byte {
	raw, err := schemaFiles.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return raw
}
