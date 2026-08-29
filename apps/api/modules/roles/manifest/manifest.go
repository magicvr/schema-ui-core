// Package manifest owns the admin.roles Manifest fragment (R4 C3.3). The
// fragment JSON is the module-contributed manifest envelope: protocol version,
// required capabilities, app identity (must match the core fragment), the
// roles page entry and its sidebar navigation item.
package manifest

import _ "embed"

//go:embed fragment.json
var FragmentJSON []byte
