You are an independent re-audit reviewer (grok build) for a goal-governance workspace. Read-only: do NOT modify any files, statuses, or code. Respond in the structured format below with concrete file:line evidence.

# Context

Goal GOAL-012-r3-s12-recycle-bin (workspace-011-admin-functional-modules): recycle bin (S-12). A first independent close-out audit (A-003) returned **verdict: fail** with 3 required findings. All were fixed. Your job: verify the fixes against the frozen design D-002 (01-decision/D-002-s1-plan-freeze.md) and the previous opinion (03-audit/A-003-s5-data-independent.md), then give a fresh verdict.

# Required findings to re-verify

## F-001 (required, high) — batch snapshot ID collision
Fix: modules/recyclebin/service.go Record now uses newSnapshotID() = "recycle-" + 8 random bytes hex (crypto/rand) per call — never time-based, so a batch delete recording several rows in one second cannot collide on the PK. Tests: service_test.go TestRecordDistinctIDsSameSecond (5 records same `now` → 5 distinct ids), TestRecordWritesSnapshot (recycle- prefix).

## F-002 (required, med) — restore-conflict HTTP test with the real service
Fix: modules/recyclebin/provider_test.go TestRecycleRealServiceRestoreConflictHTTP — real store-backed service over HTTP: create type → record snapshot → delete type → re-create same key (new id) → POST /restore → 409 RECYCLE_RESTORE_CONFLICT, snapshot stays unrestored, then resolve conflict → restore 200. Also asserts the conflict envelope carries the catalog messageKey.

## F-003 (required, med) — conflict wire format + web i18n
Fix: handler/recyclebin.go DomainError path now uses writeLocalizedError (catalog messageKey present); web i18n en-US.json/zh-CN.json add error.recycleItemNotFound + error.recycleRestoreConflict.

# Also verify recommended fixes
- F-004 (batch purge): POST /api/recycle-bin/purge-all (recycle.write; {purged:n}; audit recycle.purge) + schema toolbar purgeAll action + store PurgeAllUnrestored; descriptor/profile route lists updated; provider_test route list includes it.
- F-005 (snapshot id format): covered by F-001 fix.
- F-006 (already-restored vs conflict code): service.go returns ErrItemAlreadyRestored (409) for restored items; conflict only when key taken while unrestored.
- F-007 (docs): 01-decision.md I-001 closed; 00-meta S3/S4 checkboxes set.
- F-008 (nil semantics + dict-entries): handler recyclebin_test.go nil test now asserts recycle total == 0 after delete; service_test.go TestRestoreDictEntryRoundTrip added.

# Files to read
- apps/api/internal/modules/recyclebin/{service.go, provider.go, provider_test.go, service_test.go, store/repository.go, schema/recycle-bin.json}
- apps/api/internal/handler/{recyclebin.go, recyclebin_test.go}
- apps/api/internal/kernel/profile.go (admin.recycle-bin routes)
- apps/web/src/i18n/messages/en-US.json, zh-CN.json
- docs/workspaces/workspace-011-admin-functional-modules/GOAL-012-r3-s12-recycle-bin/01-decision/D-002-s1-plan-freeze.md, 03-audit/A-003-s5-data-independent.md

# Output format
## 范围与区间
## 逐项复核（F-00N → fixed / not fixed / partial）
(finding | fix evidence file:line | verdict)
## Findings（新增）
### F-NNN · title (level required|recommended, status open)
## 必改项汇总
## 结论 + 建议
verdict: pass | conditional | fail
