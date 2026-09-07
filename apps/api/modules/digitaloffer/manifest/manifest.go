// Package manifest owns the biz.digital-offer fragment document (VP-031 ·
// GOAL-002 D-002 v1.0.0 §7).
package manifest

import _ "embed"

//go:embed fragment.json
var FragmentJSON []byte
