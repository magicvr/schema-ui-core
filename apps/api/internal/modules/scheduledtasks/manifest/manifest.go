// Package manifest owns the admin.scheduled-tasks Manifest fragment (S-04 ·
// GOAL-010): the tasks + runs page entries and the menu_scheduled_tasks item.
package manifest

import _ "embed"

//go:embed fragment.json
var FragmentJSON []byte
