// Package manifest owns the admin.data-dictionary Manifest fragment (S-01 ·
// GOAL-008): the data-dictionary + dictionary-entries page entries and the
// menu_dictionary sidebar item.
package manifest

import _ "embed"

//go:embed fragment.json
var FragmentJSON []byte
