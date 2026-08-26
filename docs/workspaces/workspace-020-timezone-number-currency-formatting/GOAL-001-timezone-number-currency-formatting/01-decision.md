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
| I-001 | required | 时区来源：会话级 vs 用户级 vs 两者 | 方案冻结 | R1 | 用户裁决 | collecting | R1 未关闭前不直接改时区/格式 DDL | 待确认 |
| I-002 | required | 数字/货币落点：前端 vs 序列化合同 | 方案冻结 | R1 | 用户裁决 | collecting | 同上 | 待确认 |
| I-003 | required | DB `timestamptz` 是否进本波 | 退出分母 | R1 | VP 冻结投影 | **registered** | — | 冻结不进（RT-T03） |
| I-004 | required | 汇率/换算是否进本波 | 退出分母 | R1 | VP 冻结投影 | **registered** | — | 冻结不进（业务域） |
| I-005 | non-blocking | 设置归属与字段 | 方案冻结 | R2 | lead 提案 + 用户确认 | collecting | — | 待确认 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-26 | 激活与开区：lead 绑定、纲领 R1～R4、信息门禁登记 | accepted | `01-decision/D-001-workspace-root-establishment.md` |

> legacy inline 的 `## D-NNN` 记录仍可保留并被读取；新记录从目录写入。编号在本目标内单调不复用。