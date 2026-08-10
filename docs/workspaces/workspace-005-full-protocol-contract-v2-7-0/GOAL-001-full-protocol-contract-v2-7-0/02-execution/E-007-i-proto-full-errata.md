---
id: GOAL-001-full-protocol-contract-v2-7-0
doc: execution-entry
record_id: E-007
status: recorded
parent: null
created: 2026-08-10
updated: 2026-08-10
version: 0.1.0
---

# E-007 · I-PROTO-FULL-001 勘误落盘

## 已发生事实

1. `attachments/I-PROTO-FULL-001-coverage-v2-7-0.md` 已升为 v1.0.1。
2. app-manifest 执行口径已由 `37/37` 修正为 `35/37 executed + 2 excluded`；总分母口径已由 `320/320 全绿` 修正为 `318 executed + 2 excluded`。
3. 两个 exclusion id、错误码差异、适用边界与复审触发已写入权威覆盖表。
4. 12/12 域、24/24 registry type、16/16 suite include 及 VP-006 / Root 终态未改变。

## 勘误后核验

- `apps/web`：`npx vitest run src/protocol/upstream-fixtures.test.ts` → **53/53 passed**。
- `apps/web`：`npx vitest run src/protocol/conformance/stage3-fixtures.test.ts` → **260/260 passed**。
- `apps/api`：`go test ./internal/docscheck` → **passed**。
- `git diff --check` 与新增文件 trailing-whitespace 检查 → **passed**。

## 证据

- `apps/web/src/protocol/upstream-fixtures.test.ts`
- `apps/web/src/protocol/upstream/app-manifest.cases.json`
- `apps/web/src/protocol/app-manifest.ts`
- workspace-008 Root `01-decision/D-003-s0-denominator-freeze.md`
- workspace-008 GOAL-005 `attachments/S3-protocol-judgment.md`
