---
id: GOAL-004-r3-number-currency-semantics
doc: decision-entry
record_id: D-001
status: accepted
parent: GOAL-001-timezone-number-currency-formatting
created: 2026-08-26
updated: 2026-08-26
version: 0.1.0
---

## D-001 · R3 实施方案（lead 方案冻结）

### 触发

GOAL-003（R2 时区语义）已关门（A-001 self pass · 用户 2026-08-26 确认）。用户确认「直接推进 R3」。本方案把合同 §3 / §4.1 / §4.3 落为可实施设计；前置信息门禁 I-002 / I-005 均已 verified。

### 决定（方案）

1. **货币展示工具**（`apps/web/src/i18n/money.ts`，与 `format.ts` 同构）：
   - `formatMoney(value, locale, { currency, minorUnits? })`：`Intl.NumberFormat(locale, { style: "currency", currency })`；符号 / 位置 / 小数位全由 Intl 派生（§3.1 不建模板）；输入为**机器值**（金额最小单位 int64 或 number）并按 `minorUnits`（默认 2）换算展示；无效输入返回空串（沿用 fail-safe）。
   - `defaultCurrencyFor(locale)`：映射表 `zh-CN → CNY`、`en-US → USD`（§4.3）；未知 locale 回退 `USD`（不抛错）。
2. **输入解析归一化**（`apps/web/src/i18n/money.ts` 或同层解析器）：
   - `parseLocalizedMoney(raw, locale, { currency, minorUnits? }) → number | null`：剥离 locale 分组符 / 小数点 / 货币符号 → 机器值（int64 最小单位）；无法解析 → `null`（输入错误语义：消费方显示「金额格式无效」类文案，**不**提交 API）。
   - `parseLocalizedNumber(raw, locale) → number | null`：普通数字归一化（§3.2）。
   - 快测矩阵覆盖 zh-CN / en-US 各至少一场景（双 locale 双向一致，C5）。
3. **设置面 `defaultCurrency` 字段（C4 · 站点级）**：
   - API：`apps/api/internal/modules/settings/migration/migration.go` 增量迁移（`ALTER TABLE site_settings ADD COLUMN default_currency TEXT NOT NULL DEFAULT ''`，沿用既有模式）；`repository.go` struct + PATCH 参数 + 校验（ISO 4217 三字母大写；空 = 未配置——**与 locale/timezone 不同，无 `"auto"` 语义**，合同 §4.1）；`settings.go` handler 行/公开投影（`defaultCurrency`）；错误码 `error.invalidDefaultCurrency`。（F-004 回写：方案初稿误写 `"" | "auto" | 有效 ISO 4217 三字母`，实现与合同均不接受 `"auto"`——已按合同更正。）
   - Web：settings schema（`apps/api/internal/modules/settings/schema/settings.json`）Localization tab 增字段；catalog 文案；Settings 页快测断言（沿用 `startup-config.test.tsx` 模式）。
   - **不改**：Profile 默认集 / 模块矩阵 / Manifest / `docs/contracts/`；时区字段语义不动。
4. **双向一致性（C5）**：展示 ↔ 输入逆运算快测（`formatMoney → parseLocalizedMoney` 往返；双 locale）；API 机器合同不变量（§3.3）断言不改。
5. **审计**：实施切片后自审；C6 关门 = self + 本地 grok build（grok-4.6 · high）`source: independent`（migration/API data 类），意见落盘后编排器响应。

### 为什么

- 展示/解析全部走 `Intl.*` + 纯函数：可测、无模板漂移（§3.1「不建自定义格式模板」）。
- 金额机器值 = 最小单位整数（钱包先例 int64），输入解析归一化到该形态，与 API 机器合同无缝衔接。
- `defaultCurrency` 沿既有 settings migration 模式做增量列，最小化 API 面变化；PATCH 校验与错误文案与既有字段（`defaultLocale`/`siteTimezone`）一致。
- 双 locale 快测从 C1 起即覆盖，R4 关门时可直接复用为范例证据。

### 未选方案

- API 序列化携带金额语义（decimal 字符串 / 按请求 locale 输出）：破坏性变更、跨 handler 迁移，违反 §3.3 不变量；I-002 已裁决弃。
- 独立「数字格式」站点字段：locale 驱动已覆盖，弃（合同 §4.1）。
- 用户级货币偏好通道：合同 §4.2 明确不含，弃。
- `defaultCurrency` 存 JSON settings blob（而非独立列）：与既有逐个字段列模式不一致，审计/迁移可追溯性差，弃。