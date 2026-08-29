package kernel

import "strings"

// JoinIdentifiers joins contribution-key segments with "." (module-id style),
// mirroring the ContributionIdentity key grammar used across the registry.
// Renamed from JoinKeys in the v0.3.0 breaking drill（VP-023 R5 F-008）.
func JoinIdentifiers(parts ...string) string {
	return strings.Join(parts, ".")
}