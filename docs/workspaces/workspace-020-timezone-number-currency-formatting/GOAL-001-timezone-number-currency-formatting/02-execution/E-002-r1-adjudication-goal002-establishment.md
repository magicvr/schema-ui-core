---
id: GOAL-001-timezone-number-currency-formatting
doc: execution-entry
record_id: E-002
status: recorded
parent: GOAL-001-timezone-number-currency-formatting
created: 2026-08-26
updated: 2026-08-26
version: 0.1.0
---

# E-002 · R1 信息门禁裁决落盘 + GOAL-002 立项

## 2026-08-26

### 已发生事实

1. **用户裁决（P-004）**：编排器就 I-001 / I-002 / I-005 向用户提出裁决问题（附代码基证据），用户 2026-08-26 书面采纳全部推荐项（见 `01-decision/D-002-r1-info-adjudication.md`）：
   - I-001 → 会话级 auto 探测 + 站点兜底 + 用户级 localStorage 覆盖（零 DB schema 变更）
   - I-002 → 前端落点；API 保持机器合同并文档化
   - I-005 → 站点默认进 Settings→Localization tab（新增 `defaultCurrency`）；用户级覆盖并入头部 locale 通道
2. **Root 台账更新**：`00-meta.md` / `01-decision.md` 信息表 I-001/I-002/I-005 → `verified`（证据 = D-002）；R1 路线图行 → 已立项；决策索引 + D-002；执行索引 + 本条目。
3. **GOAL-002 立项**：`GOAL-002-r1-contract-freeze` 五件套 + 三个 ledger 目录建立（`01` 原语）；合同正文 `01-decision/D-001-r1-contract-freeze.md` 落盘（R1 交付物）。
4. **代码基现状核对（R1 C2 检查点证据）**：
   - `apps/api/internal/handler/settings.go`：`SiteTimezone` / `DefaultLocale` 已存在（站点级默认；`/api/branding` 公开投影 `siteTimezone` / `defaultLocale`）。
   - `apps/api/internal/handler/wallet.go`：金额以 `int64` JSON 数字承载（`amountDelta` / `balanceTotal` 等，最小单位）；无 decimal 字符串合同。
   - 时间序列化：各 handler 统一 `UTC().Format("2006-01-02T15:04:05.000Z07:00")`（RFC3339 毫秒 UTC；`internal/handler/rfc3339.go`）。
   - `apps/web/src/i18n/format.ts` / `runtime.tsx`：`Intl.DateTimeFormat` / `Intl.NumberFormat` 展示已就绪，`timeZone` 为透传可选参数；locale 用户覆盖走 `localStorage["schema-ui:locale"]` 单通道（I-L10N-002 先例）。
   - 设置面：Localization tab 已含 locale + timezone 预填（`apps/web/src/app/startup-config.test.tsx` 快测证据）。
5. 越界守卫：本轮未改任何 DDL / 迁移台账 / Profile 默认集 / 模块矩阵 / Manifest；未触碰 `docs/contracts/`（stage 门禁）。

### 证据

| 主张 | 路径 / 命令 / commit |
|------|----------------------|
| 用户裁决采纳（I-001/I-002/I-005） | 会话内 ask_user_question 答复 + `01-decision/D-002-r1-info-adjudication.md` |
| 信息项 → verified | `00-meta.md` / `01-decision.md` 信息表（证据列 = D-002） |
| GOAL-002 立项（五件套完整） | `docs/workspaces/workspace-020-timezone-number-currency-formatting/GOAL-002-r1-contract-freeze/` |
| 合同正文落盘 | `GOAL-002-r1-contract-freeze/01-decision/D-001-r1-contract-freeze.md` |
| siteTimezone/defaultLocale 站点默认已存在 | `apps/api/internal/handler/settings.go`（brandingRow / settingsRow） |
| 金额 = int64 JSON（最小单位） | `apps/api/internal/handler/wallet.go`（`amountDelta`/`balanceTotal` 等 json 字段） |
| 时间 = RFC3339 毫秒 UTC | `apps/api/internal/handler/rfc3339.go` + 各 handler `UTC().Format(...)` |
| 前端 Intl.* + localStorage locale 单通道 | `apps/web/src/i18n/format.ts`、`runtime.tsx`（`schema-ui:locale`） |
| Localization tab 已含 locale/timezone | `apps/web/src/app/startup-config.test.tsx`（#field-defaultLocale / siteTimezone 预填断言） |
| 本波未改 DDL/迁移/Profile 默认集 | git diff 限 `docs/workspaces/workspace-020-*/**`（见提交记录） |