// Package manifest owns the dev.examples Manifest fragment (W1, GOAL-002 /
// workspace-010). The fragment JSON is the module-contributed manifest envelope:
// protocol version, required capabilities, app identity (must match the core
// fragment — deliberately without homePageRef, stamped at assembly per D-003 §1),
// the 8 example pages and their top/sidebar navigation entries.
package manifest

import _ "embed"

//go:embed fragment.json
var FragmentJSON []byte
