---
id: GOAL-001-wallet-prepaid-instrument
doc: execution-entry
record_id: E-008
status: recorded
parent: null
created: 2026-09-02
updated: 2026-09-02
version: 0.1.0
---

# E-008 · 预付凭证生成体验：批次号自动生成 + 面额单位元（两位小数）（2026-09-02）

## 用户意见（产品调整）

1. 生成凭证批次时批次号不应由用户填写，后台自动生成更合理；
2. 面额按「元、两位小数」输入更符合使用习惯。

## 决策与实现事实

### 1. 批次号自动生成（向后兼容）
- 决策：HTTP `POST /api/wallet/vouchers/batches` 的 `batchId` 改为**可选**——缺省时服务端自动生成；显式提供（API 调用方）保持支持，0065 注册表冲突语义（409 `VOUCHER_BATCH_EXISTS`）不变。
- 实现：
  - `modules/wallet/voucher/voucher.go` 新增导出 `NewBatchID(now)`：`VB-<unix毫秒16hex><6字节随机hex>`（人类可读前缀 + 时间序 + 随机；注册表主键唯一）；
  - `internal/handler/wallet.go` 生成端点：`batchId` 为空 → `NewBatchID`，审计 detail 记最终 id；
- Admin 表单不再出现批次号输入（schema 移除 `batchId` 字段与 bodyMapping、i18n 键删除）。
- 测试：`TestVouchersBatchGenerateAutoBatchID`（无 batchId → 201 且 id `VB-` 前缀、两次提交不撞、可按自动 id 过滤列表）；`TestNewBatchIDUniquePrefix`（模块层）。

### 2. 面额输入单位「元 · 两位小数」（内部仍为最小单位分）
- 决策：**wire 层该端点**的 `amount` 接受人民币元（JSON number 或字符串，最多 2 位小数，如 `12.5` / `"12.50"`），服务端换算为分后走既有 int64 最小单位资金管线（账本/凭证/Redeem 不变）。这是 W16-F04「wire 存分、展示层转元」约定的输入侧对应。
- 实现：
  - `handler/wallet.go`：body.Amount 改 `json.RawMessage` + `parseYuanToCents`（字符串/数字统一、>2 位小数拒绝、指数记法经最短十进制展开、防 float 漂移铸错分；<=0/缺失/非法 → 400 `INVALID_VOUCHER_PARAMS`）；
  - `errorcatalog` `INVALID_VOUCHER_BODY`/`INVALID_VOUCHER_PARAMS` en/zh 文案更新（batchId 移除、金额按元表述）；
  - schema：生成表单 `amount` 字段 `type: inputNumber` + `step: 0.01`，label/i18n 改「元，最多 2 位小数」；批次列表 `amount` 列加 `format: "currency"`（沿用 wallet.json 账本列的 W16-F04 分→元展示约定）；
  - CSV 导出：`amount` 列按分→元两位小数输出（与列表展示一致，避免两种读数）。
- 测试：`TestVouchersGenerateAmountYuanDecimals`（12.5→1250 分、`"12.55"`→1255、0.5→50、`"0.01"`→1、1e2→10000；缺失/0/3 位小数/非数字/负值 → 400）；既有批次用例金额统一改元字面量（如 `"20.00"` → 2000 分断言不变）。

## 验证（exit 0）

```text
apps/api: go test ./modules/wallet/... ./internal/handler -run TestVouchers ./internal/composition ./internal/store（相关切片）→ 全 ok
apps/web: npx vitest run src → 91 files / 1191 tests PASS（含 wallet-vouchers D-VAL 守卫与 renderer 套件）
```

未触碰 11 件 pinned docs/schemas；账本/凭证/Redeem 最小单位语义与余额展示约定不变。

## 产物

- `apps/api/modules/wallet/voucher/voucher.go`（NewBatchID）+ `voucher_test.go`
- `apps/api/internal/handler/wallet.go`（可选 batchId + 元→分解析）+ `wallet_voucher_test.go`
- `apps/api/internal/errorcatalog/errorcatalog.go`（文案）
- `apps/api/modules/wallet/schema/wallet-vouchers.json`（移除批次输入、step 0.01、列 currency 格式）
- `apps/web/src/i18n/messages/{zh-CN,en-US}.json`（字段文案）
- `apps/web/src/renderer/render.tsx`（CSV amount 分→元两位小数）
