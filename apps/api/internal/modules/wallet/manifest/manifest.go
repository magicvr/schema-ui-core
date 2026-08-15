// Package manifest owns the admin.wallet fragment document (S-14 · GOAL-019
// D-002 §3).
package manifest

import _ "embed"

//go:embed fragment.json
var FragmentJSON []byte
