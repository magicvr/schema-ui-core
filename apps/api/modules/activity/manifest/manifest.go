// Package manifest owns the admin.activity Manifest fragment (R4 C4.2). The
// fragment JSON is the module-contributed manifest envelope: protocol version,
// required capabilities, app identity (must match the core fragment), the
// activity page entry and its sidebar navigation item.
package manifest

import _ "embed"

//go:embed fragment.json
var FragmentJSON []byte
