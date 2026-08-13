// Package manifest owns the admin.account Manifest fragment (F-03 · GOAL-005).
// The fragment JSON contributes the account page entry and its user-area
// navigation item.
package manifest

import _ "embed"

//go:embed fragment.json
var FragmentJSON []byte
