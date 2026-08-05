// Package manifest owns the admin.users Manifest fragment (R4 C3.3). The
// fragment JSON is the module-contributed manifest envelope: protocol version,
// required capabilities, app identity (must match the core fragment), the
// users page entry and its sidebar navigation item.
package manifest

import _ "embed"

//go:embed fragment.json
var FragmentJSON []byte
