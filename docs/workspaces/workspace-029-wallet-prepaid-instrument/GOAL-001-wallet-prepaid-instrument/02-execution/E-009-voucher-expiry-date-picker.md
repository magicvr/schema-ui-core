---
id: GOAL-001-wallet-prepaid-instrument
doc: execution-entry
record_id: E-009
status: recorded
parent: null
created: 2026-09-02
updated: 2026-09-02
version: 0.1.0
---

# E-009 · 过期时间改为日期选择并自动换算（2026-09-02）

## 用户意见（产品调整）

生成凭证批次时，过期时间应让客户选择日期并自动换算（而非手填 Unix 秒）。

## 决策与实现事实

- **决策**：生成模态的 `expiresAt` 由 `inputNumber`（Unix 秒）改为协议自带 **`datePicker`**（base 控件，wire = `YYYY-MM-DD` 字符串）；服务端自动换算为 Unix 秒。
- **过期语义（已定并文档化）**：所选日期解释为 **UTC 当天 23:59:59 前有效**（存储秒 = 当日 23:59:59Z；核销比较 `expiresAt < now` 为拒绝条件，末秒仍可核销）。日期窗口 `2001-09-09`～`2099-12-31`（对齐既有秒窗口 1e9..4102444800）。空值 = 不设过期。
- 向后兼容：API 仍接受 **Unix 秒**（JSON 整数/指数记法/带引号数字）——`parseVoucherExpiry` 统一：缺失/空/≤0 → 无过期；日期串 → 当日末秒；越界/格式非法 → 400 `INVALID_VOUCHER_PARAMS`。

## 改动

1. `modules/wallet/schema/wallet-vouchers.json`：`expiresAt` 字段改 `type: datePicker` + `min/max` 日期边界；label/i18n（zh/en）改「过期日期（UTC，可选；所选日期当日 23:59:59 前有效）」。
2. `apps/web/src/renderer/form-controls.tsx`：`datePicker` 分支把 `field.min/max` 透传给 `DateField`（此前被忽略）。
3. `internal/handler/wallet.go`：`body.ExpiresAt` 改 `json.RawMessage`；新增 `voucherExpiryMinUnix/MaxUnix` 常量与 `parseVoucherExpiry` / `voucherDateToEndOfDaySeconds`；过期校验块替换为统一解析（旧秒窗口逻辑并入）。
4. `internal/errorcatalog/errorcatalog.go`：`INVALID_VOUCHER_PARAMS` en/zh 文案并入日期口径。

## 测试

- `TestVouchersGenerateExpiresAtDatePicker`：`"2027-01-15"` → 201 且响应 `expiresAt = 2027-01-15T23:59:59.000Z`；`2099-12-31` 边界接受；`2100-01-01`/`2001-09-08`/非零填充/非法月/带时间串 → 400；空串 = 无过期。
- 既有 `TestVouchersGenerateExpiresAtValidation`（秒输入毫秒拒绝/越界/合法秒）继续通过——向后兼容验证。

## 验证（exit 0）

```text
apps/api: go test ./internal/handler -run TestVouchers ./modules/wallet/... ./internal/composition → 全 ok
apps/web: npx vitest run src → 91 files / 1191 tests PASS
```

未触碰 pinned docs/schemas；wire 契约扩展为「秒或日期」，语义差异仅限该端点，文档/审计无明文泄露面变化。

## 产物

- `apps/api/modules/wallet/schema/wallet-vouchers.json`
- `apps/web/src/i18n/messages/{zh-CN,en-US}.json`
- `apps/web/src/renderer/form-controls.tsx`
- `apps/api/internal/handler/wallet.go` + `wallet_voucher_test.go`
- `apps/api/internal/errorcatalog/errorcatalog.go`
