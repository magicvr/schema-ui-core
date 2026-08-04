---
id: GOAL-002-r1-contract-migration-baseline
doc: decision
status: done
parent: GOAL-001-modular-admin-architecture
created: 2026-08-04
updated: 2026-08-04
version: 0.6.0
---

# 决策记录 · GOAL-002

## 信息需求与阶段门禁

R1 required 信息项的 canonical 状态仍在 Root [00-meta.md](../GOAL-001-modular-admin-architecture/00-meta.md#信息需求与阶段门禁) 的 I-001、I-002、I-003、I-007。子目标只记录收集证据和阶段取舍，不将 Root `open` 改写为 `verified`。

| Root 信息项 | 影响门禁 | 本目标证据交付 | 状态 |
|------------|----------|----------------|------|
| I-001 | R1 方案冻结与实施 | 模块/注册/依赖清单 | open |
| I-002 | R1 迁移策略冻结与 R2 实施 | 迁移/seed/恢复/回滚边界 | open |
| I-003 | R1 契约冻结与 R2 实施 | API、Fx、生命周期错误语义决策 | open |
| I-007 | R1 范围冻结与 R4 迁移范围 | 协议基线逐项映射 | open |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-04 | 建立 R1 阶段承接子目标与四项检查点 | accepted | [01-decision/D-001-r1-stage-scope.md](01-decision/D-001-r1-stage-scope.md) |
| D-002 | 2026-08-04 | 响应 Grok A-001：补强 R1 检查点与放行边界 | accepted | [01-decision/D-002-grok-a001-response.md](01-decision/D-002-grok-a001-response.md) |
| D-003 | 2026-08-04 | 记录 C1/C2 现状基线与目标边界 | accepted | [01-decision/D-003-r1-c1-c2-baseline.md](01-decision/D-003-r1-c1-c2-baseline.md) |
| D-004 | 2026-08-04 | 冻结 C3 模块公共契约与生命周期语义 | accepted | [01-decision/D-004-r1-c3-lifecycle-contract.md](01-decision/D-004-r1-c3-lifecycle-contract.md) |
| D-005 | 2026-08-04 | 冻结 C4 协议继承与模块候选三态矩阵 | accepted | [01-decision/D-005-r1-c4-protocol-matrix.md](01-decision/D-005-r1-c4-protocol-matrix.md) |
| D-006 | 2026-08-04 | 响应 A-004 并关闭 R1 child required gate | accepted | [01-decision/D-006-a004-response.md](01-decision/D-006-a004-response.md) |
| D-007 | 2026-08-04 | Root R1 close-out 后关闭本子目标 | accepted | [01-decision/D-007-r1-child-closeout.md](01-decision/D-007-r1-child-closeout.md) |
