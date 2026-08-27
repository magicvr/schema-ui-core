---
id: GOAL-004-r3-number-currency-semantics
doc: execution-entry
record_id: E-002
status: recorded
parent: GOAL-001-timezone-number-currency-formatting
created: 2026-08-26
updated: 2026-08-26
version: 0.1.0
---

# E-002 · C1/C2/C3/C5 前端货币/数字工具实施

## 2026-08-26

### 已发生事实

1. **新模块** `apps/web/src/i18n/money.ts`：
   - `formatMoney`（C1 · §3.1）：Intl `style: currency` 展示（符号/位置/小数位全由 Intl 派生，无自建模板）；输入 = 机器值（最小单位）；`defaultCurrencyFor` 缺省；无效输入 → ""；支持 `minorUnits` 0（JPY）。
   - `defaultCurrencyFor`（C2 · §4.3）：`zh-CN → CNY`、`en-US → USD`、未知 locale → USD 兜底（兑现 GOAL-002 F-002）；`normalizeCurrencyCode`（ISO 4217 三字母）。
   - `parseLocalizedMoney` / `parseLocalizedNumber`（C3 · §3.2）：locale 分隔符剥离 + 符号/代码剥离 → 机器值（金额 = 最小单位整数；普通数字 = number）；不可解析 → `null`（调用方显示明确的输入错误，**不**提交原文）。
   - **容差留痕**：分组位序正确性（如 `12,34.5`）不验证——超出 R3 校验范围，代码注释 + 快测文档化（R4 可加严）。
2. **快测** `apps/web/src/i18n/money.test.ts`：**20/20 pass**（zh-CN/en-US 双 locale；C5 双向 round-trip `format → parse` 原值返回）。
3. 越界守卫：仅新增前端纯函数模块 + 快测；API/DB/DDL 未动（C4 单独实施）。

### 证据

| 主张 | 路径 / 命令 / commit |
|------|----------------------|
| 工具实现（§3.1/§3.2/§4.3） | `apps/web/src/i18n/money.ts` |
| 快测 20/20 | `npx vitest run src/i18n/money.test.ts`（exit 0） |
| 合同条款 | `GOAL-002-r1-contract-freeze/01-decision/D-001` §3 / §4.3 |
| GOAL-002 F-002 兑现（映射表） | `defaultCurrencyFor` + 快测（zh-CN→CNY、en-US→USD、缺省兜底） |