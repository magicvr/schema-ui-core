---
id: GOAL-001-timezone-number-currency-formatting
title: 时区 / 数字 / 货币格式语义
status: active
parent: null
created: 2026-08-26
updated: 2026-08-26
version: 0.1.0
progress: 0/4
plan_refs:
  - VP-020-timezone-number-currency-formatting
primary_plan: VP-020-timezone-number-currency-formatting
serves_summary: 承载 VP-020（Admin 功能 · 时区/数字/货币格式语义）实现：会话/用户级时区 + locale 数字/货币展示与输入合同，消费 VP-007 locale 运行时与 VP-005 设计系统。不承载汇率/换算/计费、DB timestamptz 持久化合同（RT-T03）、翻译中心，不改 Profile 默认集。
---

# GOAL-001 · 时区 / 数字 / 货币格式语义

## 概述

本 Root 承载 [VP-020-timezone-number-currency-formatting](../../../vision/plans/VP-020-timezone-number-currency-formatting.md)（**`active`** · 2026-08-26 激活，VRev-044 self `pass`）的实现：在 VP-007 已交付的多语种运行时（`zh-CN` / `en-US` + `auto` 解析）与 VP-005 设计系统之上，把「时间 / 时区 / 数字 / 货币」的展示与输入语义收成可核对的 Admin 格式合同。

**激活门禁已全部满足**（2026-08-26）：VRev-044（self）`pass`（0 required；V-F079/V-F080 → 激活事务内 fixed）；Admin 类 freshness PASS（`66f5fd1f` → `c6fda691`，不暂挂 `go`）；VP-009/VP-010 无开放阻断。

**边界**：不承接汇率 / 换算 / 计费 / 结算（业务域）；DB `timestamptz` 持久化时区合同（架构 RT-T03 仍 `registered`）；多时区排程 / 日历；翻译与文案中心（VP-007）；热加载；改 Profile 默认集 / 模块矩阵 / Manifest 装配语义。

## 纲领路线图（P-001 · V-F079 fixed）

| 阶段 | 内容 | 先后 | 状态 |
|------|------|------|------|
| R1 | **合同冻结**：时区来源（会话级 vs 用户级 vs 两者，I-001）、数字/货币语义落点（前端 vs 序列化合同，I-002）、设置归属与字段（I-005）；I-003（RT-T03 不进）/I-004（汇率不进）保持 VP 冻结投影 | 起点 | **待立项**（GOAL-002；I-001/I-002 required 裁决前不进方案冻结） |
| R2 | **时区语义**：会话/用户级时区解析与展示（IANA / offset / `auto`）；时间输入与展示统一语义 | 依赖 R1 | **待立项**（GOAL-003） |
| R3 | **数字 / 货币语义**：locale 驱动千分位 / 小数位 / 百分比 / ISO 4217 展示与输入解析合同 | 依赖 R2 | **待立项**（GOAL-004） |
| R4 | **证据与关门**：快测 + `zh-CN`/`en-US` 双 locale 范例；无越界（汇率/计费/RT-T03）；required = 0 | 依赖 R3 | **待立项**（GOAL-005） |

`progress` = 已完成阶段数 / 4。当前 **0/4**（2026-08-26 开区；R1～R4 均未立项）。

## 成功标准（方向级）

1. 时区 / 数字 / 货币格式语义合同落盘并可核对（快测 + UI 范例；`zh-CN` / `en-US` 至少各一场景）。
2. `auto` 时区解析可用；显式配置后展示与输入语义一致（同一合同双向）。
3. 未引入汇率 / 计费 / DB 持久化时区合同；未改 Charter；未改 Profile 默认集作为本波成功条件。
4. 开放 required finding = 0（或已合法闭合）。

## 信息就绪与未知项

与 VP-020 I-020-00X 同号镜像。禁止在 R1 关闭前直接改时区 / 格式相关 DDL 或迁移台账。

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 时区来源：会话级（客户端 / 请求探测）vs 用户级（存库）vs 两者；影响设置归属与 schema | 方案冻结 | R1 | 用户裁决 | collecting | — | 待确认 |
| I-002 | required | 数字 / 货币语义落点：仅前端展示 vs 后端 API 序列化也携带语义（如 decimal 字符串合同） | 方案冻结 | R1 | 用户裁决 | collecting | — | 待确认 |
| I-003 | required | 持久化时区合同（DB `timestamptz`）是否进本波 → **VP 冻结不进**（架构 RT-T03 仍 `registered`）；本行仅投影 | 退出分母 | R1 | VP 冻结投影 | **registered** | — | VP-020 I-020-003 已冻结不进 |
| I-004 | required | 汇率 / 换算是否进本波 → **VP 冻结不进**（业务域）；本行仅投影 | 退出分母 | R1 | VP 冻结投影 | **registered** | — | VP-020 I-020-004 已冻结不进 |
| I-005 | non-blocking | 默认时区 / 数字 / 货币的配置归属：Settings 哪一 tab、哪些字段 | 方案冻结 | R2 | lead 提案 + 用户确认 | collecting | — | 待确认 |

## 父目标

- null（Root；Charter `schema-ui-core-admin-foundation@0.2.0` / VP-020-timezone-number-currency-formatting）

## 台账布局

新目标为三个可追加台账创建同名平铺目录：`01-decision/`、`02-execution/`、`03-audit/`。索引文件保留 frontmatter、摘要和条目索引；独立记录使用 `D-NNN-*`、`E-NNN-*`、`A-NNN-*` 文件。