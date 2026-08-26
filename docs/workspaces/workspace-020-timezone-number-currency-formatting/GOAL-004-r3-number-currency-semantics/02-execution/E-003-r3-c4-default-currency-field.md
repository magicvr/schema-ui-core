---
id: GOAL-004-r3-number-currency-semantics
doc: execution-entry
record_id: E-003
status: recorded
parent: GOAL-001-timezone-number-currency-formatting
created: 2026-08-26
updated: 2026-08-26
version: 0.1.0
---

# E-003 · C4 设置面 defaultCurrency 字段实施（API + migration + schema）

## 2026-08-26

### 已发生事实

1. **API（apps/api）**：
   - migration：`internal/modules/settings/migration/migration.go` 新增 **Version 62** `site_default_currency`（`ALTER TABLE site_settings ADD COLUMN default_currency TEXT NOT NULL DEFAULT ''`；沿既有增量 ALTER 模式；checksum 变换 id `0062:site-default-currency:v1`）。*注：最初定 52 与 core.persistence 既有迁移冲突（全局版本号），已改 62（当前最大 61）。*
   - repository：`SiteSettings.DefaultCurrency` 字段；`ErrInvalidDefaultCurrency`；`validateCurrency`（空 或 大写三字母 ISO 4217）；PATCH 参数链 + upsert/reset SQL + SELECT/Scan + ErrNoRows 缺省 `""`。
   - handler：`SettingsRepository` 接口 + PATCH body `defaultCurrency`；`settingsRow` / `settingsAuditValues` 投影；**公开 branding 投影新增 `defaultCurrency`**（frontend 运行时取站点货币默认）；`writeSettingsError` → `INVALID_DEFAULT_CURRENCY`。
   - errorcatalog：`INVALID_DEFAULT_CURRENCY`（zh/en 文案）。
2. **Schema 与前端**：`settings.json` Localization tab 新增 `defaultCurrency` input 字段 + responseMapping；catalog 新增 `schema.settings.field.defaultCurrency` / `error.invalidDefaultCurrency`（双语）。
3. **测试**：repository 新增 `TestRepositoryDefaultCurrencyPatch`（patch/merge/lowercase+4-letter 拒绝/清空/reset）；既有 13 处 `PatchSiteSettings` 调用同步新签名。`go test ./internal/modules/settings/...` 全绿。
4. **关联 pin 测试同步（新迁移/新错误码的既有守卫）**：
   - store catalog pin：`migrate_test.go`（fresh/reopen/checksum 表 + `site_default_currency` checksum `74ede127…`）、`operations_test.go` / `restart_test.go`（applied 数 61→62）、`identity.go` `completeFingerprintCatalogHead` 61→62 + `identity_test.go` `lockedHeadExtraTables[62] = {}`（ALTER-only 无新对象）。
   - error contract pin：`handler/error_contract_test.go` `frozenLiteralCodes` + `INVALID_DEFAULT_CURRENCY`（双向契约守卫）。
   - 验证：`go test ./internal/store/... ./internal/handler/` 全绿；`go test ./...`（全量）。
5. 越界守卫：defaultCurrency 为站点设置列，不影响 Profile 默认集 / 模块矩阵 / Manifest；时区/既有列语义未动；`docs/contracts/` 未触碰。

### 证据

| 主张 | 路径 / 命令 / commit |
|------|----------------------|
| migration v62 | `apps/api/internal/modules/settings/migration/migration.go` |
| repository 支持 | `apps/api/internal/modules/settings/repository/repository.go` |
| handler/branding 投影 | `apps/api/internal/handler/settings.go`（settingsRow/brandingRow/writeSettingsError） |
| 错误码 | `apps/api/internal/errorcatalog/errorcatalog.go` |
| schema 字段 | `apps/api/internal/modules/settings/schema/settings.json`（tab-localization → form → defaultCurrency） |
| 双语文案 | `apps/web/src/i18n/messages/zh-CN.json` / `en-US.json` |
| 单测 | `go test ./internal/modules/settings/...`（exit 0）；`repository_test.go TestRepositoryDefaultCurrencyPatch` |
| 版本号冲突修正 | core.persistence 已用 52（`internal/modules/corepersistence/migration/migration.go`）→ settings 改 62 |