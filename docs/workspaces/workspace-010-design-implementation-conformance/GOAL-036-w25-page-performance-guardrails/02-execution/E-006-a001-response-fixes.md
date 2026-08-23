---
date: 2026-08-23
scope: GOAL-036 A-001 响应（F-001～F-006 修复与验证）
---

# E-006 · A-001 响应实施（F-001～F-006）

## 实施（文件级证据）

- **F-001**：`apps/api/internal/store/store.go` — `sqliteDSNParams` 增补 `_foreign_keys=on`（连接级 PRAGMA 由驱动随 DSN 施加于每条池连接；注释引用 A-001 归因）。
- **F-002**：`apps/api/internal/store/store_wal_test.go` — 新增 `TestFileStoreEveryConnectionEnforcesForeignKeys`（持 `sqlitePoolDefault` 条 Conn 断言每条 FK=1；非首连接上验证 CASCADE 删 user_roles、RESTRICT 拒删在用角色、refresh_tokens FK 拒绝）；`TestSQLiteDSNPragmas` 断言补第四参数。
- **F-003**：`apps/web/src/renderer/render.tsx` — `refreshList` 先丢弃目标 URL 的 in-flight（与 `reloadList` 对称）；`render.test.tsx` 新增挂起期不 join 用例。
- **F-004**：`attachments/README.md` 重写；`00-meta.md` 路线图 S6 勾选 I-002。
- **F-005**：`TestDeleteUsersBatchCleansRoleAndMfaLinks` 补 user_mfa 播种与断言。
- **F-006**（self 新增）：`operations_test.go` `TestOperationLogAuthEvents` 顺序断言 → 集合断言（池化后同毫秒写序非契约）。

## 验证（事实）

- `go test ./internal/store/ ./internal/modules/authsession/`：ok。
- `go test ./internal/handler/`：除**预存 flake `TestScheduledTaskRunsPagination`**（基线 stash `-count=20` 复现 `run 2 = 500`；候选机制 `newRunID` 回落路径同毫秒撞 ID；F-007 记录移交）外全绿。
- vitest 全量：**1097/1097（77 文件）**（含 F-003 新用例）。
- `tsc -b`：exit 0。
- **e2e admin × sqlite：9 通过 / 1 profile 专属跳过 / 0 失败（exit 0）**——F-001 关闭直接证据（schema-crud 完整通过）。

## 环境事项（如实记录）

基线 worktree 清理时 `Remove-Item -Recurse` 误伤 apps/web `node_modules/.bin`（pnpm 布局本无 .bin 的实情与复现细节见 E-002 相关记录；实际为 junction 穿透删除部分 bin）。已手工重建最小 shim（vite/vitest/playwright/tsc，指向真实包入口；node_modules 不入库）。该环境修复不影响仓库文件。