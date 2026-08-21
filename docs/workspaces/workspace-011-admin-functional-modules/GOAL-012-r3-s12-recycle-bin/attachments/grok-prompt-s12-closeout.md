You are an independent close-out auditor (grok build) for a goal-governance workspace. Write a structured audit opinion ONLY — you must NOT modify any files, statuses, or code. Read the repository files to verify claims. Respond in the structured format below with concrete file:line evidence.

# Scope

Goal GOAL-012-r3-s12-recycle-bin (workspace-011-admin-functional-modules): 回收站/软删除 (recycle bin) — admin.recycle-bin module: deleted-row snapshots (recycle_items, migration 0025) with browse/restore/purge; delete hooks wired through the generic resource factory (Resource.Trash → handler.TrashRecorder); managed resources v1 = dict-types / dict-entries / scheduled-tasks; permissions recycle.read/recycle.write; audit recycle.restore / recycle.purge (0026 CHECK); page "recycle-bin" + menu_recycle_bin; admin-only profile content extension.

Frozen design: docs/workspaces/workspace-011-admin-functional-modules/GOAL-012-r3-s12-recycle-bin/01-decision/D-002-s1-plan-freeze.md (and D-001-goal-boundaries.md).

# Covered artifacts

- apps/api/internal/modules/recyclebin/{provider.go, service.go, migration/migration.go, migration/provider.go, store/repository.go, schema/recycle-bin.json, manifest/fragment.json}
- apps/api/internal/handler/recyclebin.go (routes), handler/resources.go (TrashRecorder + delete/batchDelete hooks), handler/dictionary.go + handler/scheduledtasks.go (Trash wiring), modules/datadictionary/provider.go + modules/scheduledtasks/provider.go (variadic trash)
- apps/api/internal/modules/operationlog/migration/migration.go (0026 CHECK), operationlog/repository.go (EventRecycleRestore/Purge)
- apps/api/internal/kernel/profile.go, apps/api/internal/composition/composition.go, apps/api/internal/testsupport/store.go, apps/api/internal/modules/compiled/persistence.go
- apps/api/internal/errorcatalog/errorcatalog.go (RECYCLE_ITEM_NOT_FOUND, RECYCLE_RESTORE_CONFLICT), handler/error_contract_test.go
- tests: modules/recyclebin/*_test.go, handler/recyclebin_test.go, store/migrate_test.go + restart/operations tests (24→26), composition_test.go (24 perms / 13 nav)
- web: apps/web/src/test-fixtures/app-manifest.admin.json, apps/web/src/protocol/upstream-fixtures.test.ts, apps/web/src/i18n/messages/*.json, e2e shell.spec.ts, scripts/smoke.sh

# Verification focus (data gate)

1. Snapshot model: recycle_items columns (id/resource/resource_id/payload JSON/actor_id/actor_name/deleted_at/restored_at); partial UNIQUE index (resource, resource_id) WHERE restored_at IS NULL — duplicate active snapshot rejected; restored items drop out of active list.
2. Delete hook: factory delete()/batchDelete() capture row via Entity.Get BEFORE delete, record ONLY AFTER successful delete (failed delete → no orphan snapshot); Trash nil → byte-identical legacy delete semantics (users/roles/files/notifications unaffected).
3. Restore: POST /{id}/restore re-creates via owning store Create (dict-types→CreateType, dict-entries→CreateEntry, scheduled-tasks→CreateTask); unique-key conflict → 409 RECYCLE_RESTORE_CONFLICT and snapshot kept; already-restored → 409 (ErrItemAlreadyRestored); missing → 404 RECYCLE_ITEM_NOT_FOUND.
4. Purge: DELETE /{id} physically removes the snapshot (irreversible); 204; missing → 404.
5. Permission/audit: 401 anonymous, 403 without recycle.read/write (admin-only PolicyAdmin); recycle.restore/recycle.purge audit rows; 0026 CHECK includes exactly those two events.
6. Migrations: 0025 (recycle_items) + 0026 (operation_log CHECK) frozen checksums; versions 1..26 continuous; composition admin 24 perms / 13 nav; mvp/demo unchanged (8/5).
7. Error codes: RECYCLE_ITEM_NOT_FOUND + RECYCLE_RESTORE_CONFLICT in frozen set + catalog bilingual + web i18n.
8. Tests: snapshot lifecycle/unique/list/restore-mark/purge; service restore dispatch + conflicts; handler list/detail/restore/purge/403/404/audit; provider surface.

# Output format

## 范围与区间
## 成果（有证据）
(claims table: 主张 | 证据 file:line)
## 对照成功标准（本 scope）
(标准 | 结论 with 满足/部分满足/不满足)
## Findings
### F-NNN · title
| level | required | recommended |
| status | open |
| evidence |
## 必改项汇总
## 与既有意见的异同
(compare with A-001 self audit if readable)
## 结论 + 建议
verdict: pass | conditional | fail
