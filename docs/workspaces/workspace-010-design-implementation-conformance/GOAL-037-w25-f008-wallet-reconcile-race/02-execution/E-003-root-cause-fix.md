---
date: 2026-08-23
scope: GOAL-037 根治实施（用户指令：不留残余）
---

# E-003 · 根治实施：0050 数据修复迁移 + reconcile 审计原子化

## 背景（用户指令 2026-08-23：「目标37内根治此问题，不要继续遗留到以后」）

E-002/D-001 曾将以下两项记为"残余/移交"：① 既有库中"同毫秒随机序"旧流水的回放乱序；② terminal 审计事件 best-effort 可能缺失。本条目将两项在目标内根治。

## 根治 ①：0050 `wallet_ledger_order_repair` 数据修复迁移

- 实现（`modules/wallet/migration/order_repair.go`，方言中立 Go，无 DDL）：
  - 按账户回放 `(created_at, id)` 全链；首失败点定位"同毫秒 id 前缀 + 同 created_at 秒"条目组；
  - 枚举组内排列，取**唯一**满足链式 Apply 且与 balance-after 快照一致的顺序；
  - 重写该组 id（原 ms 前缀 + 组内新 seq + 原随机尾），恢复 D-002 §1 回放序契约；
  - **fail-closed**：无解/歧义/单条损坏且无同毫秒兄弟 → 迁移失败（启动中止）；组 >6 条亦拒（无法枚举）。
- Descriptors 新增 Version 50（`wallet_ledger_order_repair`，checksum `835902c…`）；`lockedHeadExtraTables[50] = {}`（data-only，无新对象）与 `completeFingerprintCatalogHead = 50` 同步。
- 既有台账测试头更新：migrate/operations/restart/identity 的 v49 锚点 → v50。
- 回归（`internal/store/migrate_0050_test.go`）：乱序库修复后对账 **consistent** 且组内 id 按创建序；健康库 no-op；坏数据（无合法序）**fail-closed**。

## 根治 ②：reconcile 成功审计原子化 + 失败/取消可观测

- `modules/wallet/jobs.go`：
  - **成功事件**（`wallet.reconcile`）移入 job 事务（`runReconcile` 的 CommitFunc 内 `RecordOperationTx`，与 ReconcileOnceTx 同事务）——**job succeeded 可见时审计必已落盘**（事务级接口 `operationlog.TransactionalRecorder` 先前已存在，本轮启用）；审计失败 → job 失败（fail closed with the domain write）。
  - 失败/取消终端事件保留 best-effort 独立写入，但**不再静默吞错**：`slog.Error` 可观测（此前 `_ = RecordOperation`）。
- 测试同步：handler stub 的成功审计同样移入 CommitFunc（与产品同构）；`TestWalletLifecycleAndAdjustFlow` 的审计断言**恢复为同步一次查询**（原子性保证，轮询删除）；`jobs_test` 的 `operationRecorder` 实现 `TransactionalRecorder`。

## 回归（事实）

- `go test ./internal/store/`：ok（含 3 个新 0050 用例 + 头锚点更新）。
- `go test ./internal/handler/ -run TestWalletLifecycleAndAdjustFlow -count=100`：**100/100 ok**。
- `go test ./... -count=1` ×2：**全绿（0 FAIL）**。
- vitest 1097/1097、`tsc -b` 0（前端未动，伴行）。
- 提交：`dbf919d`（根治代码批）。

## 结论

D-001 的"残余（既有库重对账可能红）"由 0050 迁移根治（部署即修复）；A-001 F-003 的"审计事件可能缺失"由原子化 + 可观测根治。**GOAL-037 无遗留项**。