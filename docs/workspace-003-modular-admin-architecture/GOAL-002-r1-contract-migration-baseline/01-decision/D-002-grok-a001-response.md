---
id: GOAL-002-r1-contract-migration-baseline
doc: decision-entry
decision_id: D-002
status: accepted
parent: GOAL-001-modular-admin-architecture
created: 2026-08-04
updated: 2026-08-04
version: 0.1.0
---

# D-002 · 响应 Grok A-001：补强 R1 检查点与放行边界

## 响应范围

本决策响应 [A-001](../03-audit/A-001-grok-r1-design-review.md) 的 F-001、F-002、F-003、F-004；只处理 GOAL-002 的目标定义与方案计划，不响应 R1 实施事实，也不改变 Root I-001、I-002、I-003、I-007 的状态。

## 决定

1. F-001 采用 `fixed`：将 `mvp`/`admin` Profile 候选模块集与依赖闭包矩阵纳入 C1 交付；明确精确模块集合与配置覆盖顺序仍由 Root I-004 在 R2 方案冻结前处理。
2. F-002 采用 `fixed`：将标准 Admin 模块核心六项必须能力、按需能力边界、capability 协商和 fail-closed 口径纳入 C3；按需能力不得覆盖核心六项，具体运行时实现仍属于 R2。
3. F-003 采用 `fixed`：C4 固定 VP-003 继承节和 `I-PROTO-001 v0.1.3` Q2 覆盖表路径，并声明只读取协议范围，不读取其他工作区过程状态。
4. F-004 采用 `fixed`：明确子目标 `progress: 4/4` 仅表示 C1-C4 证据收集完成；R1 冻结仍需 Root 信息项合法闭合、独立阶段审计和 `/govern` 响应。

## 理由与未选方案

上述修正均为可核对的定义补强，直接对应审计证据，成本低且不改变 VP/Root 方向，因此不采用 `accepted-residual` 或 `user-overruled`。不修改原 A-001 verdict，不把 Root I-* 写成 `verified`，也不将 C1-C4 修正文案解释为实施完成。

## 影响

- A-001 的 F-001～F-004 在响应记录中按 `fixed` 闭合；原独立意见保留。
- C1-C4 仍为未完成，子目标 `progress: 0/4`、Root `progress: 0/6` 不变。
- R1 仍处于证据收集阶段；补强后的检查点可承接后续盘点和决策，但不能单独放行 R1。
