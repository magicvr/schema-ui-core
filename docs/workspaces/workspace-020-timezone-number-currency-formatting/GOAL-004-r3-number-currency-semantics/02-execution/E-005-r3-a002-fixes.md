---
id: GOAL-004-r3-number-currency-semantics
doc: execution-entry
record_id: E-005
status: recorded
parent: GOAL-001-timezone-number-currency-formatting
created: 2026-08-26
updated: 2026-08-26
version: 0.1.0
---

# E-005 · A-002 required 修复实施（F-001/F-002 + 关联项）

## 2026-08-26

### 已发生事实

1. **F-001（high · required）fixed**：
   - `apps/api/internal/modules/settings/schema/settings.json`：`updateLocalization.bodyMapping` 补 `defaultCurrency`。
   - 守卫测试：新增 `schema/schema_test.go`（递归遍历文档树，断言 Localization 表单全部字段 ⊆ bodyMapping，且 `defaultCurrency` 显式在映射）。
   - web 保存快测：`startup-config.test.tsx` 新增「saves the Localization form incl. defaultCurrency through the real PATCH action」——真实表单提交断言 PATCH body 含 `defaultCurrency`；本地化 tab 预填断言补 `#field-defaultCurrency`。
2. **F-002（med · required）fixed**：
   - `apps/web/src/app/branding.ts`：`Branding.defaultCurrency`（类型/解析大写化/默认值 `""`）；`fetchBranding` 投影到公开启动配置边界。
   - `apps/web/src/i18n/money.ts`：新增 `resolveEffectiveCurrency(locale, siteDefault)`；`formatMoney`/`parseLocalizedMoney` 新增 `siteDefaultCurrency` 优先级通道（显式 currency > 站点默认 > §4.3 映射）。
   - `apps/web/src/i18n/runtime.tsx`：`systemDefaultUrl` fetch 增读 `defaultCurrency` → `I18nState.defaultCurrency`（`siteDefaultCurrency` prop 测试缝 + fetch 通道，与 siteTimezone 对称）。
   - 快测：money（站点默认/显式优先/无效值回退/解析符号剥离）、startup-config（fetchBranding 解析 + 缺失默认）、runtime-timezone（prop + fetch 通道）。
3. **F-003（med · recommended）fixed**：Go handler `settings_test.go`（PATCH/branding 投影断言 defaultCurrency、`INVALID_DEFAULT_CURRENCY` ×2、重置断言）；schema 守卫；web fetchBranding/保存/预填快测。
4. **F-004 fixed**：`D-001-r3-number-currency-plan.md` 回写（货币无 `"auto"` 语义，注明初稿更正）。
5. **F-009 fixed**：`migration.go` `migrate0052` → `migrate0062`（含注释）。
6. **F-010 fixed**：`00-meta.md` frontmatter `progress 0/6 → 5/6`（version 0.2.0）；`goal-tree.md` 5/6；`02-execution.md` 索引补 E-004。
7. **回归验证**：`go test ./...`（apps/api 全量绿）；`npx vitest run`（apps/web）**88 files / 1180 tests 全绿**（较上一轮 1175 新增 5 用例）。
8. F-005/F-006/F-007：`accepted-residual`（范围见 A-003 响应表；复审触发 = R4 核账），待用户书面接受。

### 证据

| 主张 | 路径 / 命令 / commit |
|------|----------------------|
| bodyMapping 修复 | `apps/api/internal/modules/settings/schema/settings.json` |
| schema 守卫 | `apps/api/internal/modules/settings/schema/schema_test.go` |
| branding 通道 | `apps/web/src/app/branding.ts` |
| money 站点默认通道 | `apps/web/src/i18n/money.ts`（`resolveEffectiveCurrency`/`siteDefaultCurrency`） |
| runtime 通道 | `apps/web/src/i18n/runtime.tsx`（`fetchedSiteDefaultCurrency`/`I18nState.defaultCurrency`） |
| web 保存/预填/branding 快测 | `apps/web/src/app/startup-config.test.tsx`；`apps/web/src/i18n/money.test.ts`；`apps/web/src/i18n/runtime-timezone.test.tsx` |
| Go handler 扩展 | `apps/api/internal/handler/settings_test.go` |
| 全量回归 | `go test ./...`（exit 0）；`npx vitest run`（88 files / 1180 tests · exit 0） |