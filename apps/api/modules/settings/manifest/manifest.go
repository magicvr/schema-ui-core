// Package manifest owns the admin.settings Manifest fragment (R4 C4.1). The
// fragment JSON is the module-contributed manifest envelope: protocol version,
// required capabilities, app identity (must match the core fragment), the
// settings page entry and its user navigation item.
package manifest

import _ "embed"

//go:embed fragment.json
var FragmentJSON []byte
