---
id: GOAL-013-r6-old-path-removal
doc: decision
status: active
parent: GOAL-001-modular-admin-architecture
created: 2026-08-05
updated: 2026-08-06
version: 0.4.2
---

# 决策记录 · GOAL-013

## 信息需求与阶段门禁

| 编号 | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 证据 / 决策 |
|------|------|-----------------|----------|--------------|-----------------|------|-------------|
| R6-I001 | required | 旧装配双轨清单与删除证据 | C6.1 | C6.1 | 全仓扫描 + 删除 | verified | E-002/E-004：旧适配器与 test 双轨删除并完成零命中扫描 |
| R6-I002 | required | store·Persistence 所有权模型与接线边界 | C6.2 | C6.2 | 设计 + 实施 | verified | D-002 + E-005～E-009 + A-004～A-008；Root A-010 F-001/F-002/F-005 fixed |
| R6-I003 | required | Schema 字节贡献发布 + 收尾 | C6.3 | C6.3 | 实施 + 测试 | collecting | D-003；E-011/`8b76ab0` Schema bytes；E-012/`2548e42` Configuration+Policy；lifecycle 与 cross 审计待完成 |
| R6-I004 | required | VP 退出 #1-#7 逐条证据 | C6.4 | C6.4 | 逐条取证 + 审计 | collecting | VP-003 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-05 | 建立 R6 旧路径移除与终态验收子目标 | accepted | [01-decision/D-001-r6-stage-scope.md](01-decision/D-001-r6-stage-scope.md) |
| D-002 | 2026-08-05 | R6 store/Persistence 所有权决策（分层/descriptor/接线/seed-reconcile） | accepted | [01-decision/D-002-r6-persistence-ownership.md](01-decision/D-002-r6-persistence-ownership.md) |
| D-003 | 2026-08-06 | R6 C6.3 Contribution、校验与生命周期终态契约 | accepted | [01-decision/D-003-r6-c63-contribution-validation-lifecycle.md](01-decision/D-003-r6-c63-contribution-validation-lifecycle.md) |

## 当前约束

- 承接 Root R6 与 R5 residual / Root A-010 债；R6 是最后阶段，完成不代表 Root/VP
  自动关门（需 exit #1-#7 逐条取证 + 关门审计）。
- 审计模式 `cross`；关门使用 Grok Build `grok-4.5` / `high`。
