---
id: GOAL-001-modular-admin-architecture
doc: decision-entry
decision_id: D-003
status: accepted
parent: null
created: 2026-08-04
updated: 2026-08-04
version: 0.1.0
---

# D-003 · 建立 R1 阶段承接子目标与四项检查点

## 决定

按 Root R1 纲领阶段建立平铺子目标 [GOAL-002-r1-contract-migration-baseline](../../GOAL-002-r1-contract-migration-baseline/00-meta.md)，承接 Root I-001、I-002、I-003、I-007 四项 R1 required 信息门禁。子目标使用四个显式检查点，初始状态 `active`、`progress: 0/4`。

Root R1 状态更新为“进行中”，但 Root R1 检查点不勾选，Root 派生 `progress: 0/6` 不变。Root 信息表仍是 I-001、I-002、I-003、I-007 的唯一状态源，均保持 `open`。

R1 方案冻结前至少需要 `source: independent` 的阶段审计；Grok Build 是候选 provider，必须先提供真实命令、身份上下文和可核对输出，provider 失败不得静默降级为 self。

## 理由与未选方案

R1 横跨 API、Web、迁移、协议和生命周期多个门禁域，收集工作有独立产物和依赖，适合一个阶段子目标统一维护。按信息项机械建立多个目标会复制 Root 状态；直接实现 Fx/模块注册会越过 R1 只冻结边界、R2 承接实现的阶段顺序；只做 self 审计不足以覆盖迁移/兼容高影响门禁。

## 后续

先收集四个检查点的现状证据和候选决策，再形成 R1 冻结候选并执行 self/independent 阶段审计。未完成前不放行 R1、不创建 R2 实施目标。
