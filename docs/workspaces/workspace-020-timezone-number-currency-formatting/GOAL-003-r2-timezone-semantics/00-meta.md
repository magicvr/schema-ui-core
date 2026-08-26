---
id: GOAL-003-r2-timezone-semantics
title: R2 · 时区语义（会话/用户级解析与展示）
status: active
parent: GOAL-001-timezone-number-currency-formatting
created: 2026-08-26
updated: 2026-08-26
version: 0.1.0
progress: 0/5
---

# GOAL-003 · R2 时区语义

## 概述

实现 Root 合同 `GOAL-002-r1-contract-freeze/01-decision/D-001` **§2（时区来源 L1～L4）** 与 **§4.2（用户级覆盖通道）**：生效时区解析器（用户覆盖 > 会话 auto 探测 > 站点默认 > 内嵌 auto）、头部时区选择 UI、站点默认接入，以及「时间展示与输入统一语义」接线。消费 VP-007 locale 运行时（`apps/web/src/i18n/`）。

**审计模式**：`self`（常规、边界清楚、可逆的非平凡实施）。不直接改时区 / 格式相关 DDL 或迁移台账（本波 API/DB 无 schema 变更）。

## 检查点（progress 派生源）

| # | 检查点 | 状态 |
|---|--------|------|
| C1 | 生效时区解析器：L1 `schema-ui:timezone` → L2 Intl 探测 → L3 `siteTimezone` → L4 auto；无效 IANA 降级；快测覆盖 | **done**（2026-08-26 · `apps/web/src/i18n/timezone.ts` + 15 快测） |
| C2 | 用户级覆盖 UI：头部 locale 通道旁时区选择；`auto` = 移除 key；登录/登出不清除；快测 | **done**（2026-08-26 · `timezone-switcher.tsx` 挂载头部；4 快测） |
| C3 | 站点默认接入：`/api/branding.siteTimezone` → L3；Localization tab 字段核对（已有字段，不改 schema） | **done**（2026-08-26 · runtime `fetchedSiteTimezone`；字段未动） |
| C4 | 统一语义接线：`formatDate`/时间输入使用生效时区（同一时刻双向一致、无偏移）；快测 | **done**（2026-08-26 · `formatDate` 默认注入 + 显式覆盖；6 快测） |
| C5 | 关门自审 A-001（self）+ 用户确认 → `status: done` | 待关门 |

`progress` = 4/5（2026-08-26：C1～C4 done；C5 待关门）。

## 成功标准

1. 生效时区解析符合合同 §2 优先级（L1→L4）与降级语义（无效 IANA 名降级、不抛错），快测可核对。
2. 用户覆盖持久化于 `localStorage["schema-ui:timezone"]` 单通道（先例 `schema-ui:locale`）；`auto` = 移除 key。
3. 站点默认 `siteTimezone` 消费（L3）；无任何配置时 `zh-CN` + `auto` 可运行，无启动硬依赖。
4. 时间展示与输入统一使用同一生效时区，双向无偏移（合同「统一语义」）。
5. 无越界：未改 DDL / 迁移台账 / Profile 默认集 / 模块矩阵 / Manifest；未触碰 `docs/contracts/`（stage 门禁）；未引入兔子洞（多时区排程、外部时区服务）。
6. 关门自审（self）落盘且 open required = 0；用户书面确认后关门并同步 goal-tree。

## 信息就绪与未知项

> 无新增信息需求。I-001（时区来源）/ I-005（设置归属）已在 Root 裁决并 `verified`（Root `D-002`；合同 D-001 §2/§4.2）。GOAL-002 F-001/F-002（recommended）在本目标推进中跟踪。

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| — | — | 无新增开放信息项 | — | — | — | — | — | 引用 Root D-002 / 合同 D-001 |

## 父目标

- `GOAL-001-timezone-number-currency-formatting`（Root；VP-020 `active` · primary_plan）

## 台账布局

新目标为三个可追加台账创建同名平铺目录：`01-decision/`、`02-execution/`、`03-audit/`。索引文件保留 frontmatter、摘要和条目索引；独立记录使用 `D-NNN-*`、`E-NNN-*`、`A-NNN-*` 文件。