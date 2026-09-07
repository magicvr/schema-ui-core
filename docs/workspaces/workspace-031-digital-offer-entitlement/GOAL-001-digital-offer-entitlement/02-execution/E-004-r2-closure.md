---
doc_type: goal-execution
id: E-004-r2-closure
parent: GOAL-001-digital-offer-entitlement
date: 2026-09-05
status: done
version: 1.0.0
---

# E-004 · R2 关门

## 事实（时间线）

- 2026-09-05 · R2（Offer CRUD + 购买 + 钱包扣款）经 GOAL-003 检查点 C1～C4 关门：
  - C1/C2/C3：实施分母 = D-002 v1.2.0（含 GOAL-003 D-001 错误码附录、D-002 凭证前置守卫附录）；`apps/api/modules/digitaloffer/` 全量落地；验收矩阵（单事务、失败零残留、幂等重放、并发收敛、重试耗尽、审计 fail-closed、限流桶、多字节）于 SQLite 与真 PostgreSQL 双库执行通过。
  - C4 审计循环：A-001 self → A-002 independent `fail`（2 high required：purchase Service API 限流缺失、双库/retry-exhaustion 分母缺失）→ A-003 响应修复 → A-004 independent `fail`（矩阵仍硬编码 SQLite，关闭声明不实）→ A-005 响应整改 → A-006 independent `conditional`（0 required，2 recommended）→ A-007 响应（可区分 rune 截断证明 + 索引去重）→ A-008 independent **`pass`（open required 0）**。
  - 审计 provider：本地 codex（gpt-5.6-sol · medium）；A-008 复审实际运行 `go test -count=1 -v` 确认矩阵子测试真实挂在 PG acceptance 下。
- 2026-09-05 · GOAL-003 `status: done`（progress 4/4）；Root 纲领进度 **2/4**。
- 2026-09-05 · 关门检查：开放 required = 0；到期 required 信息项 = 0；判据 1/2 对照可核（审计 A-008 关闭证据核对表）。
- 2026-09-05 · R3 启动：GOAL-004-r3-entitlement-validation-telegram 设立（核验/消耗 API、可选 Telegram Register、限流查询桶）。

## 产物路径

- `GOAL-003-r2-offer-purchase-wallet/03-audit/A-001...A-009`（审计台账：self ×4 / independent codex ×4）
- 代码：commit `0e65f392`、`991da571`、`326284b9`、`eb73f9f9` 及本轮关门提交

## 进度评估

- 纲领进度 2/4。R3 进行中：GOAL-004 承载（判据 3/5）。
