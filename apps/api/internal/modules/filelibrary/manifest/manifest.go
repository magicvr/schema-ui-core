// Package manifest owns the admin.file-library Manifest fragment (S-02 ·
// GOAL-007): the file-library page entry and its sidebar navigation item.
package manifest

import _ "embed"

//go:embed fragment.json
var FragmentJSON []byte
