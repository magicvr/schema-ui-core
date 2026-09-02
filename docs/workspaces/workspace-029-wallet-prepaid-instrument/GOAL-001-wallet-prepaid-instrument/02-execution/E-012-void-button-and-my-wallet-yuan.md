---
id: GOAL-001-wallet-prepaid-instrument
doc: execution-entry
record_id: E-012
status: recorded
parent: null
created: 2026-09-02
updated: 2026-09-02
version: 0.1.0
---

# E-012 · 已核销隐藏作废 + 我的钱包余额按元展示（2026-09-02）

## 用户报告

1. 预付凭证核销后，操作列「作废」仍可点（点了报错）；
2. 「我的钱包」余额按【分】展示，应显示【元】。

## 根因与修复事实

### 1. 已核销仍可点作废 → **fixed**

- **根因**：`wallet-vouchers.json` 写了 `visibleWhen: { when: "$row.status == 'unused'" }`，但表格行操作只消费 `disabledWhen {field,equals}`，且宿主表达式引擎禁止 `$row`，该门禁从未生效。
- **修正**：
  - 列表/详情 JSON 增加派生布尔 `voidable`（仅 `status=unused` 为 true）；
  - 行操作改 `visibleField: "voidable"`（注册表语法糖）+ `disabledWhen {field:voidable,equals:false}`；
  - `schema-table` 实现 `rowActionVisible`：`visibleField` 非严格 `true` 则隐藏该行操作。

### 2. 我的钱包余额按分展示 → **fixed**

- **根因**：`my-wallet.json` 三张 statCard `format: "plain"` 直出 wire 分；statCard 的 `format: "currency"` 原先只做类型门禁、不换算。
- **修正**：statCard `currency` 按 W16-F04 与表格一致（分 → 元，两位小数）；我的钱包三张余额卡与流水 `amountDelta` / `balanceAfterTotal` 声明 `format: "currency"`。

## 验证

- `go test ./internal/handler -run "TestVoucherGenerateListGet|TestVoucherVoid|TestVoucherSchemaRegistration"` PASS 2.266s（unused `voidable=true`；redeemed/void `voidable=false`；schema `visibleField=voidable`）
- vitest `schema-table.test.tsx` + `render.test.tsx` + D-VAL：106 passed（`visibleField` 隐藏行操作；statCard currency 1250 → 12.50）

## Git checkpoint

- 代码提交：`334efb2d`
- docs 提交：本记录所在 docs 笔
---
