---
id: GOAL-002-r1-contract-freeze
doc: decision
status: done
parent: GOAL-001-timezone-number-currency-formatting
created: 2026-08-26
updated: 2026-08-26
version: 0.2.0
---

# 决策记录 · GOAL-002 R1 合同冻结

## 信息需求与阶段门禁

> 本文件是稳定索引。长决策和独立决策记录放在 `01-decision/D-NNN-<slug>.md`，每条记录必须保持可独立阅读。`accepted-residual` 必须指向用户的书面决策或审计响应，且不等同于 `verified`。

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 时区来源（会话级 vs 用户级 vs 两者） | 方案冻结 | R1 | 用户裁决 | **verified** | — | Root D-002 accepted（会话级 auto + 站点兜底 + 用户级 localStorage 覆盖） |
| I-002 | required | 数字/货币落点（前端 vs 序列化合同） | 方案冻结 | R1 | 用户裁决 | **verified** | — | Root D-002 accepted（前端落点；API 机器合同文档化） |
| I-005 | non-blocking | 设置归属与字段 | 方案冻结 | R2 | lead 提案 + 用户确认 | **verified** | — | Root D-002 accepted（Localization tab + 头部 locale 通道） |

> I-003 / I-004（退出分母）为 VP 冻结投影，不属本目标。

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-26 | R1 格式语义合同冻结正文（时区来源 / 数字货币落点 / 设置归属与字段 / 内嵌默认 / 越界） | accepted（lead 提案；用户已裁决前置门禁；正文待用户审阅） | `01-decision/D-001-r1-contract-freeze.md` |

> legacy inline 的 `## D-NNN` 记录仍可保留并被读取；新记录从目录写入。编号在本目标内单调不复用。