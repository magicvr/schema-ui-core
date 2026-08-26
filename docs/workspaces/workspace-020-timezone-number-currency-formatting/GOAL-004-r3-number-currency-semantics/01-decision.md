---
id: GOAL-004-r3-number-currency-semantics
doc: decision
status: active
parent: GOAL-001-timezone-number-currency-formatting
created: 2026-08-26
updated: 2026-08-26
version: 0.1.0
---

# 决策记录 · GOAL-004 R3 数字/货币语义

## 信息需求与阶段门禁

> 本文件是稳定索引。长决策和独立决策记录放在 `01-decision/D-NNN-<slug>.md`，每条记录必须保持可独立阅读。`accepted-residual` 必须指向用户的书面决策或审计响应，且不等同于 `verified`。

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-002 | required | 数字/货币落点（前端 vs 序列化合同） | 方案冻结 | R1 | 用户裁决 | **verified** | — | Root D-002 accepted（前端；API 机器合同 §3.3 不变量） |
| I-005 | non-blocking | 设置归属与字段 | 方案冻结 | R2 | lead 提案 + 用户确认 | **verified** | — | Root D-002 accepted（Localization tab + `defaultCurrency` 新增） |

> 合同权威正文 = `GOAL-002-r1-contract-freeze/01-decision/D-001-r1-contract-freeze.md`（§3 / §4.1 / §4.3 为本目标直接消费条款）。

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-26 | R3 实施方案：货币展示 / 默认货币映射 / 输入解析归一化 / defaultCurrency 设置字段（API+migration+schema） | accepted（lead 方案冻结，遵循合同 §3/§4；migration 细节随 C4 核对） | `01-decision/D-001-r3-number-currency-plan.md` |

> legacy inline 的 `## D-NNN` 记录仍可保留并被读取；新记录从目录写入。编号在本目标内单调不复用。