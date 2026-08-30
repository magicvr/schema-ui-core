// Package manifest owns the admin.dashboard Manifest fragment (F-01 · GOAL-003).
// The fragment contributes the dashboard page entry and its first-position
// sidebar navigation item.
package manifest

import _ "embed"

//go:embed fragment.json
var FragmentJSON []byte
