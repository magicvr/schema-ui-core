// Package manifest owns the channel.telegram Manifest fragment (GOAL-006 R5,
// 判据 #5 补做 Admin UI tab): the telegram-settings page + menu_telegram
// sidebar entry.
package manifest

import _ "embed"

//go:embed fragment.json
var FragmentJSON []byte
