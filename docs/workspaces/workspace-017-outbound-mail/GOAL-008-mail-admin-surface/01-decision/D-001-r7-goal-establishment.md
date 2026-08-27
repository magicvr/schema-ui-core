---
id: D-001
doc: decision-entry
goal: GOAL-008-mail-admin-surface
status: accepted
created: 2026-08-24
updated: 2026-08-24
version: 1.0.0
---

# D-001 · 开设 R7 子目标

## 背景

Root D-007 关闭 I-009 后 R7 门禁解除。P-001：按阶段开设子目标。

## 决定

1. 创建 `GOAL-008-mail-admin-surface`，parent = `GOAL-001-outbound-mail`，status `active`。
2. 合同级条款已由 Root D-006 / D-007 与 GOAL-006 D-002 冻结，本目标不再重开；**实现级细节**（加密算法选型、DB 键位、API 路径形状、web 组件拆分）在本目标决策/执行层落盘即可——它们不改变任何已冻结语义，无需逐项问询。
3. 若实施中触发以下任一情形 → 停下按 P-004 问询：加密方案与「不入库明文」分母条款冲突；热切换需要改公共端口合同；发现安全 gap。
4. 审计模式：scaffold none → 关门 self（沿先例）。

## 未选方案

- 把 R7 与 R8（证据/探针）合并实施：违反 P-001 阶段切分。
- 为每个实现细节开问询：P-004 无未决裁决点，空转问询违反「规则能唯一判定时不询问」。
