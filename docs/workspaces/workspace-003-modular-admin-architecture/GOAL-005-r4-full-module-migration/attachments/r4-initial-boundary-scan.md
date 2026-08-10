---
id: r4-initial-boundary-scan
doc: evidence
goal: GOAL-005-r4-full-module-migration
date: 2026-08-05
status: recorded
---

# R4 初始边界扫描

## 当前实现事实

- Composition Root resolves the selected Plan and still calls central
  `handler.Register`: `apps/api/internal/composition/composition.go:29,89`.
- The central handler mounts Users/Roles resources:
  `apps/api/internal/handler/health.go:22-39`.
- Users and Roles already have domain-specific Resource definitions and store
  boundaries: `apps/api/internal/handler/users.go:34-51`,
  `apps/api/internal/handler/roles.go:21-39`,
  `apps/api/internal/store/users.go:19-52`,
  `apps/api/internal/store/roles.go:16-60`.
- The generic Resource factory provides HTTP methods/authentication/permissions:
  `apps/api/internal/handler/resources.go:48-59,156-199`.
- Users/Roles schemas remain in the central embedded fixture set and the owner map
  is hard-coded in `apps/api/internal/handler/schema.go:13-18,63-75`.
- The current module metadata has contribution keys but no provider fields:
  `apps/api/internal/kernel/module.go:28-56`.
- Users/Roles manifest navigation is already represented in the baseline and Web
  navigation is manifest-driven: `apps/api/internal/manifest/app-manifest.json:106-138`,
  `apps/web/src/app/navigation.ts:53-107`.
- Existing renderer/API/e2e tests cover Users/Roles CRUD and authority boundaries:
  `apps/api/internal/handler/users_test.go:17-432`,
  `apps/api/internal/handler/roles_test.go:16-301`,
  `apps/web/src/renderer/schema-crud.test.tsx:390-697`,
  `apps/web/e2e/schema-crud.spec.ts:19-137`.

## Records conflict

- VP-003 R4 says “将 users、roles、records/Schema CRUD 及其他现有 Admin 能力迁入
  统一能力契约” (`docs/vision/plans/VP-003-modular-admin-architecture.md:65-72`).
- Migration `0006 records_retire` drops the Records table and cleans its permissions
  (`apps/api/internal/store/migrate.go:291-310`).
- No current `admin.records`, `/api/records`, Records handler, or Records Schema
  fixture was found; records operation-log events are historical only.

This is an information conflict, not a decision. R4-I003 remains `collecting` and
blocks C2 until the user/canonical governance record resolves it.

## Operationlog boundary

`RecordOperation` is append-only (`apps/api/internal/store/operations.go:44-60`),
while Users/Roles hooks log failures without reversing the successful business write
(`apps/api/internal/handler/users.go:265-301`, `roles.go:228-256`). Retention duration,
archival and a possible strong-consistency change were not found; R4-I004 remains
collecting.
