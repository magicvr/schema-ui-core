// Package manifest owns the admin.data-permission fragment document (S-09 ·
// GOAL-016 D-002 §3).
package manifest

import _ "embed"

//go:embed fragment.json
var FragmentJSON []byte
