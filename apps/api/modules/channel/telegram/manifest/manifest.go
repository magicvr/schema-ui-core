// Package manifest owns the channel.telegram Manifest fragment: the
// telegram-settings sidebar page, its telegram-operator inner page, and the
// menu_telegram sidebar entry.
package manifest

import _ "embed"

//go:embed fragment.json
var FragmentJSON []byte
