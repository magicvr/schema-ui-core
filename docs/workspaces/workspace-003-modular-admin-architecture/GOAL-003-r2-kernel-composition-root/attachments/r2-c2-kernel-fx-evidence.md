# R2 C2 Kernel and Fx boundary evidence

- Framework-agnostic contract: `apps/api/internal/kernel/module.go`
- Profile expansion: `apps/api/internal/kernel/profile.go`
- Composition root: `apps/api/internal/composition/composition.go`
- Dependency version: `apps/api/go.mod` records `go.uber.org/fx v1.24.0`.

The kernel has no Fx import. Fx providers, lifecycle hooks, HTTP, Store, and
secret adapters are confined to `internal/composition`.
