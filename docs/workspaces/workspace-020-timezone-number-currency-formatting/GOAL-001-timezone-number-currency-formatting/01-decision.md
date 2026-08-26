---
id: GOAL-001-timezone-number-currency-formatting
doc: decision
status: active
parent: null
created: 2026-08-26
updated: 2026-08-26
version: 0.1.0
---

# 决策记录 · GOAL-001 时区/数字/货币格式语义

## 信息需求与阶段门禁

> 本文件是稳定索引。长决策和独立决策记录放在 `01-decision/D-NNN-<slug>.md`，每条记录必须保持可独立阅读。`accepted-residual` 必须指向用户的书面决策或审计响应，且不等同于 `verified`。

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 时区来源：会话级 vs 用户级 vs 两者 | 方案冻结 | R1 | 用户裁决 | **verified** | — | 2026-08-26 用户裁决：会话级 auto + 站点兜底 + 用户级 localStorage 覆盖（D-002） |
| I-002 | required | 数字/货币落点：前端 vs 序列化合同 | 方案冻结 | R1 | 用户裁决 | **verified** | — | 2026-08-26 用户裁决：前端落点；API 保持机器合同并文档化（D-002） |
| I-003 | required | DB `timestamptz` 是否进本波 | 退出分母 | R1 | VP 冻结投影 | **registered** | — | 冻结不进（RT-T03） |
| I-004 | required | 汇率/换算是否进本波 | 退出分母 | R1 | VP 冻结投影 | **registered** | — | 冻结不进（业务域） |
| I-005 | non-blocking | 设置归属与字段 | 方案冻结 | R2 | lead 提案 + 用户确认 | **verified** | — | 2026-08-26 用户确认：站点默认进 Localization tab（新增 defaultCurrency）；用户级覆盖并入头部 locale 通道（D-002） |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-26 | 激活与开区：lead 绑定、纲领 R1～R4、信息门禁登记 | accepted | `01-decision/D-001-workspace-root-establishment.md` |
| D-002 | 2026-08-26 | R1 信息门禁裁决：I-001 时区来源 / I-002 数字货币落点 / I-005 设置归属（用户书面采纳） | accepted | `01-decision/D-002-r1-info-adjudication.md` |

> legacy inline 的 `## D-NNN` 记录仍可保留并被读取；新记录从目录写入。编号在本目标内单调不复用。