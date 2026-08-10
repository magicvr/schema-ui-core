# R2 C1 Profile and graph evidence

- Profile names and defaults: `apps/api/internal/kernel/profile.go`
- Environment resolution: `apps/api/internal/config/config.go`
- Fail-closed cases: unknown Profile, custom without explicit modules, unknown
  module, duplicate enabled module, missing dependency, cycle, contribution
  conflict, missing capability, and incompatible kernel API range.
- Tests: `apps/api/internal/kernel/kernel_test.go` and
  `apps/api/internal/composition/composition_test.go`.

The recorded precedence is compiled Profile default, explicit
`modules.enabled`, then environment source observability. Explicit modules do
not cause dependencies to be silently added; the selected set must already be
closed.
