---
id: r1-public-wire-inventory-v0.1
doc_type: evidence-attachment
title: R1 公共时间 wire formatter/parser/fixture inventory v0.1
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-002-r1-contract-and-denominator-freeze
version: 0.1.0
---

# R1 公共时间 wire inventory v0.1

## 用户合同

- 输出：固定 `YYYY-MM-DDTHH:MM:SS.ffffffZ`（UTC、6 位小数、只用 `Z`）。
- 输入：兼容合法 RFC3339 的 0/3/6/9 位小数与 `+00:00` 等等价 offset；解析后转 UTC；拒绝无时区/模糊本地时间。
- DB JSON/TEXT payload 内嵌时间不进入 VP-040 列分母；若业务 endpoint 暴露它，按现有 payload contract 单独回归。

## Go formatter / handler 输出面

### 共享 formatter 与其测试

- `apps/api/internal/handler/rfc3339.go:5-8`：现行 `rfc3339Milli` 与 `formatRFC3339Milli`，C2 目标改为 shared fixed-6 formatter。
- `apps/api/internal/handler/rfc3339_test.go:8-12`：现行 `.000Z` fixture，需补 fixed-6 + parser compatibility。

### 当前 inline fixed-3 outputs（需迁移到 shared formatter）

- `apps/api/internal/handler/account_self.go:105-106,381-382`
- `apps/api/internal/handler/datapermission.go:85`
- `apps/api/internal/handler/digitaloffer.go:291-292,335,347-351`
- `apps/api/internal/handler/jobs.go:325-347`
- `apps/api/internal/handler/notifications.go:80,90`
- `apps/api/internal/handler/operations.go:42`
- `apps/api/internal/handler/recyclebin.go:179,183`
- `apps/api/internal/handler/roles.go:81-82`
- `apps/api/internal/handler/settings.go:128`
- `apps/api/internal/handler/users.go:130-131`
- `apps/api/internal/handler/wallet.go:85-95,575-578,938,980,992,1000-1007`
- `apps/api/modules/wallet/jobs.go:219`
- `apps/api/cmd/server/server_restart_test.go:82`
- `apps/api/internal/handler/recyclebin_test.go:76`
- `apps/api/internal/handler/resources_test.go:53,67`

### Current RFC3339 outputs (need fixed-6 output policy review)

- `apps/api/internal/handler/invites.go:41-42,182-184`
- `apps/api/internal/handler/service_credentials.go:264-277`
- `apps/api/internal/handler/telegram_operator.go:240,334,505-506`
- `apps/api/internal/handler/scheduledtasks.go:418`
- `apps/api/modules/authsession/email_identity.go:109`
- `apps/api/modules/authsession/recovery.go:52`
- `apps/api/internal/handler/operations.go:111` (filter parser)
- `apps/api/internal/handler/jobs.go:389` (input parser)
- `apps/api/internal/handler/service_credentials.go:170` (input parser)
- `apps/api/modules/recyclebin/service.go:277-280` (payload parser)

### Non-DB filesystem output

- `apps/api/internal/handler/filelibrary.go` formats filesystem `ModTime`; C2 must explicitly decide whether it joins the public fixed-6 output contract. It is not a DB column.
- `apps/api/cmd/schema-ui/configpkg.go:307,324` emits package metadata `ExportedAt/ImportedAt` via `time.RFC3339`; C2 must explicitly keep or migrate this non-DB metadata surface.

## Web consumer / fixtures

- `apps/web/src/host/failure.ts:93,280` accepts RFC3339-like values with optional fractional digits; preserve compatibility and add fixed-6 cases.
- `apps/web/src/lib/datetime.ts` and `datetime.test.ts` consume/format response times; add 6-digit and legacy 3-digit fixtures.
- Representative fixtures with `.000Z` appear in `apps/web/src/i18n/ui-bilingual.test.tsx`, `s5-denominator-render.test.tsx`, `error-localization.test.tsx`, `app/notification-bell.test.tsx`, `app/representative-pages.integration.test.tsx`, `startup-config.test.tsx`, `renderer/*` time fixtures, and wallet/jobs tests.
- Existing `apps/web/src/lib/datetime.ts` regex allows variable fractional digits, but this does not prove fixed-6 output coverage.

## C2/R3 closure requirements

1. Replace inline output layouts with one shared fixed-6 UTC formatter, or document an explicit non-DB exception.
2. Add parser compatibility tests for no fraction / 3 / 6 / 9 fraction digits and `+00:00`, plus rejection of no-zone values.
3. Update Go and Web fixtures and restart/handler tests; preserve wire shape for non-time fields.
4. R3 matrix must include at least one endpoint from each persisted unit group (seconds, milliseconds, nullable, sentinel-converted) and VP-020 session/user timezone display round-trip.
5. No `pgtype` or SQLite TEXT value may cross into handler/module public contracts.

This inventory is a planning/coverage artifact; it does not claim formatter code has been changed.
