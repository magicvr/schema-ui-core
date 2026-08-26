---
id: GOAL-004-r3-number-currency-semantics
title: R3 · 数字 / 货币语义（locale 驱动格式与输入解析合同）
status: active
parent: GOAL-001-timezone-number-currency-formatting
created: 2026-08-26
updated: 2026-08-26
version: 0.1.0
progress: 0/6
---

# GOAL-004 · R3 数字 / 货币语义

## 概述

实现 Root 合同 `GOAL-002-r1-contract-freeze/01-decision/D-001` **§3（数字/货币展示与输入解析，API 机器合同不变量）** 与 **§4.1/§4.3（设置面 `defaultCurrency` 字段 + 默认货币映射表）**：locale 驱动货币展示（ISO 4217 符号/位置/小数位）、默认货币映射、locale 化输入解析归一化为机器值、站点级默认货币字段（API migration + settings schema + Localization tab），并保证展示 ↔ 输入双向一致。

**审计模式**：`independent`（本目标含设置表 migration 与 API 行为变更 → data/migration 类影响；按项目规则关门时自审后调用本地 grok build（grok-4.6 · high）执行 `source: independent` 独立审，意见落盘后由编排器响应）。

## 检查点（progress 派生源）

| # | 检查点 | 状态 |
|---|--------|------|
| C1 | 货币展示工具：`formatMoney`（Intl `style: currency`，ISO 4217 符号/位置/小数位随 locale）+ 快测（zh-CN/en-US） | **done**（2026-08-26 · `apps/web/src/i18n/money.ts` + 快测） |
| C2 | 默认货币映射表 `defaultCurrencyFor(locale)`（§4.3：zh-CN→CNY、en-US→USD；缺省不抛错）+ 快测 | **done**（2026-08-26 · `defaultCurrencyFor` + `normalizeCurrencyCode`） |
| C3 | 输入解析归一化：locale 化数字/货币输入 → 机器值（金额 int64 最小单位；普通数字 number）；错误语义明确 + 快测 | **done**（2026-08-26 · `parseLocalizedMoney`/`parseLocalizedNumber` → null 语义） |
| C4 | 设置面 `defaultCurrency` 字段：API（site_settings 增量 migration + repository/PATCH 校验 + settings 行/公开投影）+ settings schema（Localization tab）+ 单测/快测 | 待实施 |
| C5 | 双向一致性核对：展示 ↔ 输入逆运算一致（同一合同反向）+ 双 locale 场景快测 | **done**（2026-08-26 · round-trip 快测 en-US/zh-CN） |
| C6 | 关门：自审 + grok build independent + 用户确认 → `status: done` | 待关门 |

`progress` = 4/6（2026-08-26：C1/C2/C3/C5 done；C4 待实施；C6 待关门）。

## 成功标准

1. 货币展示合同落码：ISO 4217 代码 / 符号 / 位置 / 小数位随 locale 与显式 currency（§3.1），`zh-CN` / `en-US` 至少各一场景快测可核对。
2. 默认货币映射表（§4.3）实现：无显式 `defaultCurrency` 时按有效 locale 映射（zh-CN→CNY、en-US→USD），缺省不抛错、可核对（兑现 GOAL-002 F-002）。
3. 输入解析按 §3.2 归一化为机器值（金额 = int64 最小单位；普通数字 = number），解析失败有明确输入错误语义，**不向 API 发送未归一化 locale 字符串**。
4. 站点级 `defaultCurrency`（ISO 4217）端到端：API 列 + PATCH 校验 + 投影 + Localization tab 字段；**不改 Profile 默认集**。
5. 展示 ↔ 输入双向一致（同一合同反向可用），快测断言。
6. 无越界：API 机器合同格式不变量保持（§3.3）；未改时区/DDL 语义外字段；未触碰 `docs/contracts/`（stage 门禁）；open required = 0。

## 信息就绪与未知项

> 无新增 required 信息项。I-002（前端落点）/ I-005（设置归属）已 verified（Root D-002）。承接跟踪项：GOAL-002 F-001/F-002、GOAL-003 F-001/F-002（recommended；F-002 GOAL-002 由 C2 兑现，其余随 R4 核账）。

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| — | — | 无新增开放信息项 | — | — | — | — | — | 引用 Root D-002 / 合同 D-001 §3/§4 |

## 父目标

- `GOAL-001-timezone-number-currency-formatting`（Root；VP-020 `active` · primary_plan）

## 台账布局

新目标为三个可追加台账创建同名平铺目录：`01-decision/`、`02-execution/`、`03-audit/`。索引文件保留 frontmatter、摘要和条目索引；独立记录使用 `D-NNN-*`、`E-NNN-*`、`A-NNN-*` 文件。