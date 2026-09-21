package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// operatorConfigPath is the shipped operator config, relative to this package.
const operatorConfigPath = "../../configs/config.yaml"

// TestOperatorConfigCoversAdminPreset guards the drift that broke the list
// pages' 「导出所选」 button (2026-09-19, user report).
//
// configs/config.yaml declares `profile: custom` with an inline module list
// whose stated intent is "the full admin preset PLUS channel.telegram". VP-038
// added admin.jobs to the admin preset (kernel/profile.go) but did not add it
// here, so on the operator config the batch-export route was never mounted:
// POST /api/jobs/batch-export answered a bare 404 (NOT_FOUND / 「未找到」) and
// the button could not work, even though the schema node was rendered on both
// the users and roles pages.
//
// The invariant is deliberately about the PRESET, not about one module id: a
// custom list that claims to mirror the admin preset must be a superset of it,
// so the next module added to the preset cannot silently go missing again.
// channel.telegram (and any other deliberate extra) is allowed on top.
func TestOperatorConfigCoversAdminPreset(t *testing.T) {
	path := operatorConfigPath
	if _, err := os.Stat(path); err != nil {
		// A consumer install may not carry the operator config; the guard is
		// about THIS repo's shipped file.
		t.Skipf("operator config not present: %v", err)
	}
	raw, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatalf("read operator config: %v", err)
	}
	text := string(raw)

	if !strings.Contains(text, "profile: custom") {
		t.Skip("operator config no longer uses an explicit custom module list")
	}

	adminPreset, err := kernel.ResolveProfile(string(kernel.ProfileAdmin), nil)
	if err != nil {
		t.Fatalf("resolve admin preset: %v", err)
	}
	if len(adminPreset.Modules) == 0 {
		t.Fatal("admin preset is empty — cannot verify coverage")
	}

	declared := parseInlineModuleList(t, text)
	missing := make([]string, 0)
	for _, moduleID := range adminPreset.Modules {
		if !declared[moduleID] {
			missing = append(missing, moduleID)
		}
	}
	if len(missing) > 0 {
		t.Fatalf(
			"configs/config.yaml declares `profile: custom` (documented as the admin preset + channel.telegram) "+
				"but is missing %d admin-preset module(s): %s.\n"+
				"A module absent from this list is never mounted, so any schema node that calls its routes "+
				"fails at runtime with a bare 404 (NOT_FOUND) — exactly how 「导出所选」 broke when admin.jobs "+
				"was missing. Add the module(s) to the `app.modules.list` block.",
			len(missing), strings.Join(missing, ", "),
		)
	}

	// The batch-export entry point specifically: this is the route the list
	// pages' 「导出所选」 button depends on. Keeping it explicit makes the
	// regression named even if the preset is later restructured.
	if !declared["admin.jobs"] {
		t.Fatal("configs/config.yaml must enable admin.jobs: without it POST /api/jobs/batch-export is not mounted")
	}
}

// parseInlineModuleList extracts the `app.modules.list` entries from the
// operator YAML. It is a deliberately small scanner (not a YAML dependency):
// the operator config uses one `- <id>` entry per line under the list key.
func parseInlineModuleList(t *testing.T, text string) map[string]bool {
	t.Helper()
	declared := make(map[string]bool)
	lines := strings.Split(text, "\n")
	inList := false
	listIndent := 0
	for _, line := range lines {
		trimmed := strings.TrimRight(line, " \t\r")
		content := strings.TrimSpace(trimmed)
		if content == "" || strings.HasPrefix(content, "#") {
			continue
		}
		indent := len(trimmed) - len(strings.TrimLeft(trimmed, " "))
		if !inList {
			if content == "list:" {
				inList = true
				listIndent = indent
			}
			continue
		}
		// A list item is `- <id>`; anything at or above the list key's indent
		// that is not an item ends the block.
		if strings.HasPrefix(content, "- ") {
			id := strings.TrimSpace(strings.TrimPrefix(content, "- "))
			if id != "" {
				declared[id] = true
			}
			continue
		}
		if indent <= listIndent {
			inList = false
		}
	}
	if len(declared) == 0 {
		t.Fatalf("no app.modules.list entries parsed from the operator config (scanner drift?) — runtime=%s", runtime.Version())
	}
	return declared
}
