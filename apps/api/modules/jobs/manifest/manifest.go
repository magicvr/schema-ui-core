// Package manifest owns the admin.jobs Manifest fragment (GOAL-003 R2): the
// jobs page entry and the menu_jobs sidebar item.
package manifest

import _ "embed"

//go:embed fragment.json
var FragmentJSON []byte
