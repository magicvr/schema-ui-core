package kernel

import "strings"

// JoinKeys joins contribution-key segments with "." (module-id style),
// mirroring the ContributionIdentity key grammar used across the registry.
// Added in the R4 zero-conflict upgrade drill as an A-layer additive sample.
func JoinKeys(parts ...string) string {
	return strings.Join(parts, ".")
}