---
id: r3-c1-readiness-scan
title: R3 C1 readiness scan
status: recorded
parent: GOAL-004-r3-bounded-pilot
created: 2026-08-05
updated: 2026-08-05
version: 0.1.0
---

# R3 C1 readiness scan

## Governing gate

VP-003 defines R3 as a bounded pilot, not a full migration:
`docs/vision/plans/VP-003-modular-admin-architecture.md:74-110` requires the
five A deliveries, four B lesions, V-1-V-4, and a D-gate. I-006 requires the
static Manifest/Shell removal list, compatibility deadline, warning behavior,
and rollback trigger: `GOAL-001/00-meta.md:109`.

## Current implementation map

- Centralized API route mounting remains in
  `apps/api/internal/handler/health.go:20-36`; it calls
  `registerOperations` and `settingsHandler` directly.
- The composition root still invokes the centralized registration without a
  module-plan contribution boundary at
  `apps/api/internal/composition/composition.go:87-99`.
- Settings routes and permission checks are present at
  `apps/api/internal/handler/settings.go:16-23,50-150`.
- Activity is the read-only `operations` resource at
  `apps/api/internal/handler/operations.go:14-26,68-85`; operationlog storage
  and event persistence are in `apps/api/internal/store/operations.go:16-131`.
- Profile descriptors include `core.operationlog` for both Profiles and
  optional `admin.settings`/`admin.activity` contributions at
  `apps/api/internal/kernel/profile.go:24-47,91-107`.
- Composition projects the Manifest from `kernel.Plan` at
  `apps/api/internal/composition/composition.go:87-99`; this does not yet
  prove that route registration follows the same plan.
- Web page mapping and generic Schema rendering are present at
  `apps/web/src/app/App.tsx:37-62,313-326,436-447` and
  `apps/web/src/renderer/render.ts:106-110`.
- Existing API tests cover Settings and operations basics at
  `apps/api/internal/handler/settings_test.go:16-160` and
  `apps/api/internal/handler/operations_test.go:17-130`; no V-1-V-4
  cross-layer acceptance tests were found.

## Initial gate assessment

| Area | Current evidence | Status |
|------|------------------|--------|
| I-006 old-entry inventory | Source paths identified; removal list not frozen | collecting |
| Compatibility/warning/removal | No bounded deadline/trigger record | collecting |
| Rollback/recovery | No R3 disablement rehearsal | collecting |
| Module-controlled routes | Profile descriptors exist; handler mount is centralized | gap |
| V-1-V-4 | No complete cross-layer tests or runtime evidence | not started |

This scan is a starting snapshot in the dirty worktree. It is not clean
revision, CI, release, or R3 gate evidence.
