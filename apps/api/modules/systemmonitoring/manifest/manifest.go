// Package manifest owns the admin.system-monitoring Manifest fragment (S-03 ·
// GOAL-009): the monitoring page entry and its sidebar navigation item.
package manifest

import _ "embed"

//go:embed fragment.json
var FragmentJSON []byte
