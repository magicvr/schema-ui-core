// Package manifest owns the admin.recycle-bin Manifest fragment (S-12 ·
// GOAL-012): the recycle-bin page entry and the menu_recycle_bin item.
package manifest

import _ "embed"

//go:embed fragment.json
var FragmentJSON []byte
