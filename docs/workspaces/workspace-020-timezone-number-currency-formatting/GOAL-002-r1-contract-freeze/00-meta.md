---
id: GOAL-002-r1-contract-freeze
title: R1 · 合同冻结（时区来源 / 数字货币落点 / 设置归属）
status: active
parent: GOAL-001-timezone-number-currency-formatting
created: 2026-08-26
updated: 2026-08-26
version: 0.1.0
progress: 2/3
---

# GOAL-002 · R1 合同冻结

## 概述

冻结 VP-020 首波的**格式语义合同**（Root R1）：在用户裁决（Root `D-002` accepted）基础上，把时区来源、数字/货币语义落点、设置归属与字段收成可核对、可引用的合同正文，供 R2（时区语义）/ R3（数字/货币语义）直接消费。

**审计模式**：`self`（低风险、可逆、文档型合同冻结；无需 independent）。关门审计前不直接改时区 / 格式相关 DDL 或迁移台账。

## 检查点（progress 派生源）

| # | 检查点 | 状态 |
|---|--------|------|
| C1 | 合同正文落盘：`01-decision/D-001-r1-contract-freeze.md` 覆盖时区来源 / 数字货币落点 / 设置归属与字段 / 内嵌默认 / 越界声明 | **done**（2026-08-26 落盘） |
| C2 | 合同与代码基现状核对一致（siteTimezone/defaultLocale 已有；金额 int64；RFC3339 UTC；Intl.* 前端展示；Localization tab 已有 locale/timezone） | **done**（2026-08-26，证据见执行 E-001 与 Root E-002） |
| C3 | 关门自审：`03-audit/A-001`（source=self）落盘，open required = 0，经用户确认后 `status: done` | **进行中**（A-001 已落盘 · verdict pass · required = 0；待用户审阅合同并确认关门） |

`progress` = 2/3（2026-08-26；C3 未完成）。

## 成功标准

1. 合同正文落盘并可核对：覆盖时区来源（用户级覆盖 > 会话 auto 探测 > 站点默认 > 内嵌 auto）、数字/货币**前端落点** + API 机器合同（时间 RFC3339 毫秒 UTC、金额 int64 最小单位 JSON）、设置归属与字段（Localization tab 站点默认 + 头部 locale 通道用户覆盖）、内嵌默认 `zh-CN` + `auto` 可运行。
2. 合同可直接驱动 R2 / R3 立项（语义与落点无歧义、无开放 required 信息项）。
3. 无越界：未改 DDL / 迁移台账 / Profile 默认集 / 模块矩阵 / Manifest；未引入汇率 / 换算 / 计费 / DB `timestamptz` 持久化合同；未触碰 `docs/contracts/`（stage 门禁）。
4. 关门自审 A-001（self）落盘且 open required = 0；用户确认后关门并同步 goal-tree。

## 信息就绪与未知项

> 本目标无新增信息需求。I-001 / I-002 / I-005 已在 Root 裁决并 `verified`（Root `01-decision/D-002`）；I-003 / I-004 为 VP 冻结投影（不进）。

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| — | — | 无新增开放信息项（引用 Root I-001/I-002/I-005 verified） | — | — | — | — | — | — |

## 父目标

- `GOAL-001-timezone-number-currency-formatting`（Root；VP-020 `active` · primary_plan）

## 台账布局

新目标为三个可追加台账创建同名平铺目录：`01-decision/`、`02-execution/`、`03-audit/`。索引文件保留 frontmatter、摘要和条目索引；独立记录使用 `D-NNN-*`、`E-NNN-*`、`A-NNN-*` 文件。