// Package manifest owns the admin.notifications Manifest fragment (F-04 · GOAL-006).
package manifest

import _ "embed"

//go:embed fragment.json
var FragmentJSON []byte
