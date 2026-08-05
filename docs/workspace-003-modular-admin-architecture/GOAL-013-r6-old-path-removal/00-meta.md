---
id: GOAL-013-r6-old-path-removal
title: R6 · 旧路径移除与终态验收
status: active
parent: GOAL-001-modular-admin-architecture
created: 2026-08-05
updated: 2026-08-06
version: 0.5.0
progress: 2/4
plan_refs:
  - VP-003-modular-admin-architecture
primary_plan: VP-003-modular-admin-architecture
serves_summary: 承接 Root R6：旧装配路径终态删除（handler 级适配器/test 双轨）、store·Persistence 所有权迁出（Root A-010 F-001/F-002/F-005）、Schema 字节 ContributionSet 发布、Configuration 运行时迁移、完整回归与 VP 退出判据 #1-#7 逐条取证，最终 Root close-out + VP-003 关门审计。
---

# GOAL-013 · R6 旧路径移除与终态验收

## 概述

本子目标是 Root `GOAL-001-modular-admin-architecture` 的 R6 阶段（最后阶段）：删除
双轨兼容与静态生产兜底，完成 store·Persistence 所有权迁出（Root A-010 债）与
Schema 字节 ContributionSet 发布，运行完整回归（双 Profile、升级/恢复、失败路径、
容器/fork），对 VP 退出判据 #1-#7 **逐条**取证后关门审计，并为 Root close-out 与
VP-003 关门提供完整证据。

## 愿景对齐

| 字段 | 值 |
|------|----|
| `parent` | `GOAL-001-modular-admin-architecture` |
| `plan_refs` | `VP-003-modular-admin-architecture` |
| `primary_plan` | `VP-003-modular-admin-architecture` |
| Charter | `schema-ui-core-admin-foundation@0.2.0`（经 VP-003 间接对齐） |
| 审计模式 | `cross`；关门使用 Grok Build `grok-4.5` / `high` |

## 成功标准

- [x] **C6.1 / 旧路径终态删除**：handler 级 RegisterSettings/RegisterActivity 适配器、
  MountProviderRoutes test-only、死代码、test 双轨删除；生产仅 provider finalize
  一条装配路径。
- [x] **C6.2 / Persistence 所有权迁出**：store 上帝对象拆分（平台 runner/ledger vs
  模块仓储，按 D-002 设计）；`CollectPersistence` 生产接线 + 0001-0008 descriptor
  归属；seed/RBAC reconcile 以 Authorization/system-data 贡献为源。
  **（A-004/A-005/A-006 self + A-007 independent pass；A-008 响应后 F-C62-004 与
  Root A-010 F-001/F-002/F-005 均 fixed）**
- [ ] **C6.3 / Schema 字节贡献驱动 + 收尾**：Schema document 字节由 ContributionSet
  发布（去掉中心静态枚举）；Configuration 运行时迁移、PolicyID/Visibility 深化、
  双 Profile Start/Ready 失败矩阵。
- [ ] **C6.4 / 验收与关门**：完整回归（双 Profile/升级恢复/失败/容器/fork）+ VP 退出
  判据 #1-#7 逐条取证 + self + Grok independent 无开放 required finding；
  Root close-out + VP-003 关门依据。

四个检查点等权；当前 `progress: 2/4`。完成本子目标表示 R6 关闭；Root close-out 与
VP-003 关门另需确认。

## 信息门禁

| 编号 | 级别 | 必须回答的问题 | 影响 | 最晚阶段 | 收集动作 | 状态 | 证据 |
|------|------|----------------|------|----------|----------|------|------|
| R6-I001 | required | 旧装配双轨清单与删除证据？ | C6.1 | C6.1 | 全仓扫描 + 删除 | verified | E-002 清单 + E-004：MountProviderRoutes/RegisterSettings/RegisterActivity 已删 |
| R6-I002 | required | store·Persistence 所有权模型与 CollectPersistence 接线边界？ | C6.2 | C6.2 | 设计 + 实施 | verified | D-002 + E-006～E-009 + A-004～A-008：catalog/Apply、system-data reconcile 与 owner repositories 经 self + Grok independent 验证，A-010 F-001/F-002/F-005 fixed |
| R6-I003 | required | Schema 字节贡献发布 + 收尾项边界？ | C6.3 | C6.3 | 实施 + 测试 | collecting | D-003 已冻结四切片契约；实现、测试与 cross 审计待完成 |
| R6-I004 | required | VP 退出 #1-#7 逐条证据是否齐全？ | C6.4 | C6.4 | 逐条取证 + 审计 | collecting | VP-003 |

## 阶段路线图

1. Persistence ownership 设计冻结（D-002/E-003，R6-I002 verified）。
2. C6.1：测试走 Provider/composition 真路径 → 删 handler 级死适配器与双轨。
3. C6.2：按 D-002 拆 store + 接线 CollectPersistence + seed/reconcile 贡献驱动。
4. C6.3：Schema 字节 ContributionSet 发布 + 收尾。
5. C6.4：完整回归 + VP 退出 #1-#7 逐条取证 + self + Grok + Root/VP 关门。

## 范围与非目标

范围包括旧路径删除、Persistence 迁出、Schema 字节发布、Configuration/校验器/矩阵
收尾、完整验收与退出判据取证。非目标包括 Records 恢复、R7+ 新阶段、愿景变更。
`0003`/`0006` 迁移账本与历史 operation-log 保留。
