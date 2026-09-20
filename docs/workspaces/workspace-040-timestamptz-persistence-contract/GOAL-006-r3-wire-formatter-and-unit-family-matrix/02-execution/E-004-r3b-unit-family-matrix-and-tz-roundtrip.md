---
id: E-004-r3b-unit-family-matrix-and-tz-roundtrip
doc: execution-entry
status: active
parent: GOAL-006-r3-wire-formatter-and-unit-family-matrix
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-004 · R3-B：单位族矩阵 + VP-020 会话时区展示 round-trip

## 1. 单位族矩阵（inventory / `D-018` §4 R3-B 要求 4）

载体：`apps/api/internal/handler/w040_r3b_unit_family_matrix_test.go` → `TestR3BUnitFamilyWireMatrix`，
**每个族一个不同端点、一个不同表**：

| 族 | 列（据 `internal/temporalcontract` 校验 unit） | 端点 | 断言 |
|----|----------------------------------------------|------|------|
| 秒 | `digital_offers.created_at` / `updated_at`（`sec`） | `GET /api/digitaloffer/offers` | `.900000Z` 与 `.000000Z` 逐字节相等 |
| 毫秒 | `jobs.created_at` / `updated_at` / `finished_at`（`ms`） | `GET /api/jobs` | 写入 `.900000Z` 种子 → 列表逐字节相等；未完成 job 的 `finishedAt` 字段**缺席**（NULL 不伪造瞬时） |
| 可空 | `service_credentials.revoked_at` / `last_used_at`（`sec`，nullable） | `GET /api/service-credentials/{id}` | 新建 → 存储层确为 SQL NULL 且 wire 为 JSON `null`；置值后 wire 为逐字节固定 6 位 |
| sentinel 转换 | `mail_config.updated_at`（`ms`，D0：legacy `0` → NULL） | `GET /api/mail/config` | 存储层为 SQL NULL（legacy 0 已消失）→ wire JSON `null`；置值 → 逐字节固定 6 位 |

**族身份可核对**：子测试先调 `requireDenominatorColumn(table, column, unit)`，从**冻结分母**读回 unit 并断言 `sec`/`ms`，因此「这确实是秒/毫秒族」不是散文断言而是可执行检查。可空族另断言存储层 `SELECT revoked_at, last_used_at` 为 SQL NULL。sentinel 族另断言存储层 `SELECT updated_at` 为 SQL NULL。

**不为凑数复用同一族**：秒（digital_offers）与毫秒（jobs）分属不同表；可空族用 nullable 列（`revoked_at`/`last_used_at`），sentinel 族用 D0 列（`mail_config.updated_at`）——四族四个端点，无重复充数。

## 2. VP-020 会话时区展示 round-trip（`I-040-004`；载体按 `I-041-007`）

### Go 侧（`apps/api/internal/handler/w040_r3b_timezone_roundtrip_test.go`）

- `TestR3BWireRoundTripIsZoneIndependent`：对 6 个瞬时（尾零微秒 / 整秒 / 亚微秒截断 / 负 Unix epoch / US DST 春进与秋退边界）× 4 个时区（UTC、`Asia/Shanghai`、`America/New_York`、`Asia/Kathmandu` +05:45）断言
  1. `FormatWireTime(t.In(zone))` 恒等于同一规范串（**wire 是瞬时，不是墙钟**，任意展示时区都不会移动它）；
  2. `ParseWireTime(FormatWireTime(t))` 恢复出的瞬时等于 `temporal.Truncate(t)`（微秒精确、向零截断，`D-008`）；
  3. format → parse → format 逐字节恒等。
- `TestR3BSessionTimezoneDoesNotMoveStoredOrSentInstants`：把进程 `time.Local` 换成 UTC / New York / Kathmandu，wire 值恒为 `2026-09-20T12:57:15.900000Z`（钉住「环境时区不会渗入输出」）。

### Web 侧（`apps/web/src/i18n/utc-roundtrip.test.tsx`，组件 + 单测）

- 组件：`I18nProvider` 同时给出 L1 存储偏好 `Asia/Shanghai`、L2 会话探测 `America/New_York`、L3 站点默认 `Europe/London` → 解析出的有效时区是 **L1**，且 `formatDate(wire)` 恰等于「生产 formatter 绑定该时区」的结果，同时**不等于** L2 渲染与 UTC 渲染（若层解析错位，断言必然失败）。
- 不变量：渲染不会改动传入的瞬时（`Date` 渲染前后 `toISOString()` 不变）。
- round-trip：对 5 个 wire 瞬时 × 4 个时区，用 `Intl` 部件取出墙钟 → 反解瞬时 → 重新给出规范 wire，按**秒**粒度恒等（人类的展示层只到分钟，微秒留在存储侧；该粒度已在注释里写明）。
- 边界：`Asia/Kathmandu` 的 `+05:45`（整小时假设会算错，实测 18:42）、US DST 春进（`2026-03-08`）与秋退（`2026-11-01`）两个边界均不漂移。
- 输入侧：Web 的时间输入是 date-only（`datePicker` → `YYYY-MM-DD`，直通不转换，见 `renderer/form-controls.test.ts`），服务端拒绝无时区输入（`D-005`），故不会有时区把某个日期挤到前一天的通路。

## 3. `I-040-004` 收口

信息项：**VP-020 展示/输入与 UTC 存储的回归矩阵**。本轮以「Go 侧 wire/瞬时 round-trip + Web 侧会话时区展示 round-trip」双载体矩阵关闭：

- 展示随有效时区（L1→L4）变化，**存储与传输恒为同一 UTC 规范串**；
- 展示层→瞬时→wire 的 round-trip 按秒恒等，微秒由 Go 侧证明恒等；
- 未新增任何 VP-020 能力（`D-018` §3 非目标）：只做回归矩阵。

状态：`open → verified`（证据即本文件与上述两个测试文件）。

## 证据

- `go test -count=1 ./internal/handler/ -run "TestR3BUnitFamilyWireMatrix"` → 4/4 子测试 PASS。
- `go test -count=1 ./internal/handler/ -run "TestR3BWireRoundTripIsZoneIndependent|TestR3BSessionTimezoneDoesNotMoveStoredOrSentInstants"` → PASS（6 子例 + 1）。
- `cd apps/web && vitest run src/i18n/utc-roundtrip.test.tsx` → 8/8 PASS。
- 全仓：`apps/api` `go test -count=1 ./...` 64/64 包 ok；`apps/web` `vitest run` 124 文件 / 1504 用例 ok。
