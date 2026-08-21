---
id: GOAL-002-r1-contract-migration-baseline
doc: decision-entry
decision_id: D-001
status: accepted
parent: GOAL-001-modular-admin-architecture
created: 2026-08-04
updated: 2026-08-04
version: 0.1.0
---

# D-001 · 建立 R1 阶段承接子目标与四项检查点

## 决定

1. 在当前工作区平铺创建 `GOAL-002-r1-contract-migration-baseline`，作为 Root R1 的单一主承接子目标，状态为 `active`。
2. 用四个可独立核验检查点承接 Root I-001、I-002、I-003、I-007：模块/注册盘点、迁移与 seed 边界、Fx/API/生命周期契约、协议继承矩阵。
3. Root 的 I-001、I-002、I-003、I-007 仍以 Root 信息表为唯一状态源并保持 `open`，Root R1 检查点与 `progress: 0/6` 在证据完成前不变。
4. R1 方案冻结前至少需要一条 `source: independent` 的阶段审计意见；provider 必须有可核对的真实输出。Grok Build 作为用户指定的候选 provider，先验证实际调用能力，失败不得静默降级。

## 理由

R1 横跨 API、Web、迁移、协议与生命周期多个门禁域，收集工作具有独立产物、依赖和阶段放行价值，符合 P-001/P-005 的阶段子目标条件。一个 R1 子目标可以集中维护跨项矩阵并保留推荐上下文，避免为每个信息项机械创建两个目标，同时保留后续 R2/R3 按阶段递进拆分的空间。

## 未选方案

- **为 I-001～I-007 各创建信息目标**：不采用；重复维护 Root 状态，且 I-004～I-006 的最晚阶段不在本 R1 子目标范围内。
- **现在直接实现 Fx/模块注册/Manifest 聚合**：不采用；R1 只冻结可实施边界，R2 才承接内核与组合根实现。
- **只做 self 审计**：不采用；迁移/兼容/生命周期边界属于 Root 预置的高影响范围，需 independent 交叉审计后才能无条件冻结。

## 影响与后续

- Root 的阶段状态由“未开始”转为“进行中”，但不勾选 R1、不增加 Root 派生 progress。
- 本目标先收集现状证据；只有四个检查点均有可核对产物、方案冻结决策和阶段审计结论后，才可响应 Root I-001/I-002/I-003/I-007，并提议 R1 放行。
- Grok Build 的审计原始输出、命令和身份上下文必须写入本目标审计记录或附件；普通 Codex 探查结果只能作为主线候选证据，不冒充 independent。
