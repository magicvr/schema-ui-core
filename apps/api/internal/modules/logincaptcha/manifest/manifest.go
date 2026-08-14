// Package manifest owns the admin.login-captcha Manifest fragment (S-11 ·
// GOAL-011): the captcha page entry and the menu_captcha item.
package manifest

import _ "embed"

//go:embed fragment.json
var FragmentJSON []byte
