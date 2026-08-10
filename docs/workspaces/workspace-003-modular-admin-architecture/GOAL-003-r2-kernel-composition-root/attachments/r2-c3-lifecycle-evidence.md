# R2 C3 Lifecycle evidence

- Pure lifecycle semantics and reverse cleanup:
  `apps/api/internal/kernel/lifecycle.go`
- Runtime adapter and resource cleanup:
  `apps/api/internal/composition/composition.go`
- Start failure and reverse cleanup tests:
  `apps/api/internal/kernel/kernel_test.go` and
  `apps/api/internal/composition/composition_test.go`.
- Existing readiness fault test:
  `apps/api/internal/handler/health_test.go`.

The occupied-port test confirms a failed Fx start closes the Store rather than
leaving a file handle behind. The process restart test confirms the normal
server path remains usable.
