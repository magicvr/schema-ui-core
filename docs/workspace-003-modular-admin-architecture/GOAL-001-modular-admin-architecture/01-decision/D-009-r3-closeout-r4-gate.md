---
id: D-009-r3-closeout-r4-gate
doc: decision-entry
goal: GOAL-001-modular-admin-architecture
date: 2026-08-05
status: accepted
---

# D-009 · R3 close-out 与 R4 阶段入口

## 决定

确认 GOAL-004 R3 已以 `done 4/4` 收束，Root I-006 由 R3 A-004/E-005/D-004
响应为 `verified`，Root progress 派生为 `3/6`。允许建立一个新的 R4 主实施
子目标，但不把 R3 试点结果解释为 Root 或 VP-003 完成。

R4 的最小连贯范围是 Users、Roles、Records/Schema CRUD 及其统一模块契约
能力迁移：HTTP、Schema、Authorization、Navigation、Manifest、Persistence；
Profile 运维、升级/恢复、readyz/诊断、容器和 fork 文档留给 R5。R4 必须保持
I-PROTO-001 v0.1.3 范围，扩大 domain 或修改 exclude 时另行决策并升级覆盖表。

## 依据与门禁

- GOAL-004 A-004 以 fixed 路径关闭 R3 required findings；A-001/A-002/A-003
  原始意见保持不变。
- R4 必须等独立子目标建立、范围和信息项落盘后再实施；不得批量预建后续
  阶段目标。
- R4 预置 `independent` 审计建议；高影响 migration/production 证据不得静默
  降级为无审计。
