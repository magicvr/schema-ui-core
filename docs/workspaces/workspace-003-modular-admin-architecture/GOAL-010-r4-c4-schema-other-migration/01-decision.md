---
id: GOAL-010-r4-c4-schema-other-migration
doc: decision
status: active
parent: GOAL-005-r4-full-module-migration
created: 2026-08-05
updated: 2026-08-05
version: 0.1.0
---

# 决策记录 · GOAL-010

## 信息需求与阶段门禁

| 编号 | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 证据 / 决策 |
|------|------|-----------------|----------|--------------|-----------------|------|-------------|
| C4-I001 | required | settings/activity 中心注册/schema/manifest/operationlog 读面状态 | C4.1/C4.2 | C4.1 | 全仓扫描 | collecting | 待 C4.1 |
| C4-I002 | required | Schema owner map 转贡献驱动语义 | C4.3 | C4.3 | 设计 + 测试 | collecting | GOAL-009 F-IND-C33-001 |
| C4-I003 | required | Manifest secrecy/Ready 清理/校验器/ledger drift 边界 | C4.4 | C4.4 | GOAL-008 E-004 + 实施 | collecting | GOAL-008 E-004 |
| C4-I004 | non-blocking | Records historical-only 保持 | C4.4 | C4.4 | 负向断言 | open | D-003；GOAL-007 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-05 | 建立 R4-C4 Schema 与其他能力迁移子目标 | accepted | [01-decision/D-001-r4-c4-stage-scope.md](01-decision/D-001-r4-c4-stage-scope.md) |
| D-002 | 2026-08-05 | C4.4 成功标准收窄与 ledger 门禁移交 C5 | accepted | [01-decision/D-002-r4-c4-criteria-scope.md](01-decision/D-002-r4-c4-criteria-scope.md) |

## 当前约束

- 承接 C1-C3 冻结契约与 GOAL-009 迁移模式；C4 只迁移 settings/activity 等剩余
  Schema-driven Admin 能力，不恢复 Records、不宣称 C5、不推进 Root progress。
- Schema owner map 转贡献驱动（解决 F-IND-C33-001）；C3 遗留门禁（Manifest secrecy/
  Ready 清理/校验器/ledger drift）在 C4.4 实施。
- 审计模式 `independent`；迁移切片使用 Grok Build `grok-4.5` / `high`。
