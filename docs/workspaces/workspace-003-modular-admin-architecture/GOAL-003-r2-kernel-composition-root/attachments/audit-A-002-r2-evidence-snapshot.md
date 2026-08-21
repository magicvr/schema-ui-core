---
id: audit-A-002-r2-evidence-snapshot
title: R2 source and verification evidence snapshot
status: recorded
parent: GOAL-003-r2-kernel-composition-root
created: 2026-08-04
updated: 2026-08-04
version: 0.1.0
---

# R2 source and verification evidence snapshot

## Revision boundary

- Snapshot date: 2026-08-04.
- Base revision: `HEAD=b1b7650b3202de7a7a7ce6c0bdffe212093fe75f`.
- The working tree was dirty and had no clean checkpoint for this run. The
  commands below therefore prove the current local snapshot only.
- `git diff --check` returned no whitespace errors; it emitted only the
  repository's existing LF-to-CRLF checkout warnings.
- No CI, deployment, release, or clean-revision acceptance is claimed.

## Commit-pinned source map

The following paths and line references are the source snapshot reviewed for
the R2 self opinion:

- Kernel API and range validation: `apps/api/internal/kernel/module.go:11`,
  `:70-71`, `:102-121`, `:137-180`, `:245`.
- Profile resolution and fail-closed custom modules:
  `apps/api/internal/kernel/profile.go:9-20`, `:24-45`, `:49-87`.
- Runtime lifecycle and reverse cleanup:
  `apps/api/internal/kernel/lifecycle.go:17-68` and
  `apps/api/internal/composition/composition.go:106-151`.
- Migration collection and store integration:
  `apps/api/internal/migration/collector.go:12-97` and
  `apps/api/internal/store/migrate.go:459-476`.
- Deterministic Manifest projection and aggregation:
  `apps/api/internal/manifest/manifest.go:34-254`.
- Public Manifest endpoint and ETag handling:
  `apps/api/internal/handler/manifest.go:16-41`.
- Fx composition root and lifecycle wiring:
  `apps/api/internal/composition/composition.go:25-151`.
- Development and production proxy boundary:
  `apps/web/vite.config.ts:20-28`, `apps/web/nginx.conf:10-32`, and
  `apps/web/Dockerfile:9-24`.

## Pinned-hash repair evidence

- The existing provenance values are preserved in
  `apps/web/src/protocol/upstream/provenance.json:2-36`.
- The R2 hash boundary is explicit in
  `apps/web/src/protocol/upstream-fixtures.test.ts:62-72`,
  `apps/web/src/protocol/conformance/stage3-fixtures.test.ts:67-73`, and
  `apps/web/src/renderer/permissions-inheritance.test.ts:32-47`.
- Before the repair, the current CRLF byte hashes differed from the recorded
  values; after CRLF-to-LF canonicalization, all recorded values matched. The
  constants were not changed.

## Reproducible local commands

- `go test ./...` from `apps/api`: passed.
- `npm test -- --run` from `apps/web`: passed, 23 test files and 492 tests.
- `npm run build` from `apps/web`: passed, TypeScript and Vite production build.
- `git rev-parse HEAD`: `b1b7650b3202de7a7a7ce6c0bdffe212093fe75f`.
- `git diff --check`: no whitespace errors; line-ending warnings only.

## Limitations

The snapshot is sufficient to reproduce the local result with the current
working tree, but it is not a clean revision and does not establish CI or
release acceptance. Independent Grok re-audit remains the next gate.
