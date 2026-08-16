// Package schema owns the admin.wallet page documents (S-14 · GOAL-019 D-002
// §3): the wallet accounts page and the wallet-entries ledger page.
package schema

import "embed"

// ModuleID is the owning module for the schema handler's contribution-driven
// page gating.
const ModuleID = "admin.wallet"

//go:embed wallet.json wallet-entries.json my-wallet.json
var schemaFiles embed.FS

// PageIDs are the page identifiers this module contributes.
func PageIDs() []string { return []string{"wallet", "wallet-entries", "my-wallet"} }

// SchemaDocuments returns the page documents owned by the wallet module.
func SchemaDocuments() map[string][]byte {
	return map[string][]byte{
		"wallet":         mustRead("wallet.json"),
		"wallet-entries": mustRead("wallet-entries.json"),
		"my-wallet":      mustRead("my-wallet.json"),
	}
}

func mustRead(path string) []byte {
	raw, err := schemaFiles.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return raw
}
