---
id: GOAL-004-r3-number-currency-semantics
doc: execution-entry
record_id: E-001
status: recorded
parent: GOAL-001-timezone-number-currency-formatting
created: 2026-08-26
updated: 2026-08-26
version: 0.1.0
---

# E-001 · 立项与 R3 实施方案冻结

## 2026-08-26

### 已发生事实

1. GOAL-004 五件套 + 三个 ledger 目录建立（`01` 原语；模板源 `docs/templates/goal-folder/`）。
2. 前置门禁核对：I-002 / I-005 已 verified（Root D-002）；合同 `GOAL-002/01-decision/D-001` §3 / §4.1 / §4.3 为本目标消费条款（R1、R2 均已关门）。
3. R3 实施方案冻结：`01-decision/D-001-r3-number-currency-plan.md`（货币展示 / 默认货币映射 / 输入解析归一化 / defaultCurrency 设置字段（API migration + schema + UI）/ 双向一致性 / 审计路径）。
4. 审计模式记录：`independent`（设置表 migration + API 行为变更；关门时自审后调用本地 grok build grok-4.6 · high）。
5. 越界守卫：本轮改动限于 `docs/workspaces/workspace-020-timezone-number-currency-formatting/**`。

### 证据

| 主张 | 路径 / 命令 / commit |
|------|----------------------|
| 五件套与 ledger 目录齐全 | `docs/workspaces/workspace-020-timezone-number-currency-formatting/GOAL-004-r3-number-currency-semantics/` |
| 方案冻结 | `GOAL-004-r3-number-currency-semantics/01-decision/D-001-r3-number-currency-plan.md` |
| 前置门禁 verified | Root `GOAL-001-.../01-decision/D-002`（I-002/I-005 证据列） |
| 合同条款来源 | `GOAL-002-r1-contract-freeze/01-decision/D-001-r1-contract-freeze.md` §3 / §4.1 / §4.3 |
| 既有 migration 模式参照 | `apps/api/internal/modules/settings/migration/migration.go`（增量 ALTER 先例） |