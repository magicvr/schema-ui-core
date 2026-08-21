---
id: E-002
title: R2 implementation slices and focused verification
status: recorded
created: 2026-08-04
updated: 2026-08-04
parent: GOAL-003-r2-kernel-composition-root
version: 0.1.0
---

# E-002 R2 implementation slices

Recorded facts:

- `apps/api/internal/kernel/` now contains the framework-agnostic module,
  Profile, capability, contribution-conflict, semantic API-range, and
  lifecycle contracts. Fx is absent from that package.
- `apps/api/internal/config/` resolves `APP_PROFILE` and
  `APP_MODULES_ENABLED`; invalid profile/module input is retained as a startup
  error instead of silently selecting a default.
- `apps/api/internal/composition/` validates the compiled graph before Fx,
  owns named secret dependencies, runs the kernel runtime adapter, and closes
  the Store on listener/start/readiness failure paths.
- `apps/api/internal/migration/` is called by the Store compiled-migration
  validator. `apps/api/internal/manifest/` provides deterministic aggregation
  and module-selected page/navigation projection.
- `apps/api/internal/handler/manifest.go` exposes the public exact endpoint
  with deterministic ETag/304 behavior. `apps/web/vite.config.ts` and
  `apps/web/nginx.conf` proxy that exact path; `apps/web/Dockerfile` removes
  the static fixture from the production image.
- `go test ./...` in `apps/api` passed after these changes. The composition
  tests include Profile graph failure, selected Manifest pages, and occupied
  port startup failure. Web production build passed. Web full test still has
  three pre-existing pinned SHA mismatch failures; the relevant fixture files
  were not changed in this stage.

The facts above are implementation-slice evidence. C1-C4 can be assessed as
complete checks in this child; C5 remains open until the self review, Grok
independent R2 gate review, and any required-finding response are recorded.
