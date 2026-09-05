---
doc_type: goal-execution
id: E-002-audit-cycle
parent: GOAL-003-r2-offer-purchase-wallet
date: 2026-09-05
status: done
version: 1.0.0
---

# E-002 · C4 审计循环（self → independent → 响应）

## 事实（时间线）

- 2026-09-05 · **A-001 self 自审**（execution-facts）：verdict pass，2 项 recommended（F-001 限流接线判为 R3 项——后被 A-002 纠正；F-002 searchQ 截断）。F-002 当日修复。
- 2026-09-05 · **A-002 independent 实施审计**（codex · gpt-5.6-sol · medium · workspace-write）：verdict **fail**，2 high required（F-001 purchase Service API 限流缺失——纠正 A-001 的 R3 延期判定；F-002 §4.4 双库/retry-exhaustion 验收分母缺失）+ F-003（searchQ 字节截断破坏 UTF-8，med recommended）+ F-004（台账/版本投影滞后，low recommended）。codex 实际执行 `go build ./...` 与相关包测试（均通过）。
- 2026-09-05 · **实施期新发现**：PG 并发双发暴露——并发同请求各调用预生成不同 purchaseID、共享幂等键，READ COMMITTED 下败者以「同键不同 RefID」触达 ledger 触发 wallet `ErrIdempotencyConflict` 终态，破坏收敛（SQLite 写者串行未暴露）。按合同修订纪律处理：**D-002 v1.2.0**（GOAL-003 D-002 附录）——凭证 INSERT 前置为请求守卫；败者在守卫处失败并回读重放，永不触碰钱包。
- 2026-09-05 · **A-003 self 响应**：F-001 fixed（§8 purchase 桶接入 Service API + `TestPurchaseRateLimited`：预算/耗尽/subject 隔离/窗口恢复/无 Clear）；F-002 fixed（`SetPurchaseFaultHookForTest` 注入 seam + `runPurchaseMatrix` 双库矩阵：retry exhaustion 恰 3 次且零残留、终态不重试、瞬时失败恢复无重复、审计 fail-closed、多字节 Q）；F-003 fixed（rune 截断）；F-004 fixed（索引/投影同步）。`go test ./modules/digitaloffer/...`（含真 PG 矩阵）全绿。
- 2026-09-05 · git checkpoint：`0e65f392`（R2 实现）+ 本轮修复（待提交）。

## 产物路径

- `03-audit/A-001-self-implementation-review.md`、`03-audit/A-002-independent-implementation-audit.md`、`03-audit/A-003-self-response-a002.md`
- `01-decision/D-001-error-codes-addendum.md`、`01-decision/D-002-tx-order-addendum.md`（D-002 合同 v1.1.0/v1.2.0 附录）

## 进度评估

- C1/C2/C3 关门；C4 循环进行中：待 A-004 independent closure 复审通过后 GOAL-003 关门（4/4）。
