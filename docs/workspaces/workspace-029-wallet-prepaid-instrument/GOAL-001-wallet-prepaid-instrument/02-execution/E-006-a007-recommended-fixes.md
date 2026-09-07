---
id: GOAL-001-wallet-prepaid-instrument
doc: execution-entry
record_id: E-006
status: recorded
parent: null
created: 2026-09-02
updated: 2026-09-02
version: 0.1.0
---

# E-006 · A-007 响应与 A-005 recommended（F-002～F-005）处置（2026-09-02）

## 用户指令

`/govern 响应复审意见 A-007。处理一下 recommended。`

## 意见背景

- A-007（independent · grok-build · **pass**）：A-005 F-001 闭合证据独立复审成立（composition HTTP 回归走生产装配）；F-002～F-005 recommended 维持 open backlog。
- 本轮处置：F-002～F-005 全部按 `fixed` 落地（附证据），A-008 登记。

## 修复事实（按 finding，时间线）

### F-004 · `batch_id` 无唯一约束 → **fixed**（0065 批次注册表）
1. `modules/wallet/migration/migration.go` 新增迁移 **0065 `wallet_voucher_batches`**（双方言 DDL）：`voucher_batches(batch_id TEXT PRIMARY KEY, created_at, updated_at)` + 从既有 `vouchers` 回填（GROUP BY batch_id；历史重复批次不拆解，约束只防 NEW 混批——非资金不变式加固）。
2. `voucher/voucher.go` 新增哨兵 `ErrVoucherBatchExists`；`voucher/service.go` `GenerateBatch` 在**同一事务**内先 `INSERT ... ON CONFLICT (batch_id) DO NOTHING` 登记批次，0 影响行 → 拒绝整批（并发重复也 fail-closed）。
3. `handler/wallet.go` 生成端点把该错误映射为 **409 `VOUCHER_BATCH_EXISTS`**（此前一切 domain 错误都被吞成 500 INTERNAL）；`errorcatalog` 新增同码（en/zh + messageKey），`error_contract_test.go` 冻结清单同步。
4. 迁移台账钉更新：`store` 各冻结目录测试（migrate/operations/restart/identity）由 v64 终点同步至 v65，`completeLostLedgerTables` 含 `voucher_batches`。
5. 测试：voucher 服务层 `TestGenerateBatchDuplicateIDRejected`（重复 id → ErrVoucherBatchExists 且列表仍 2 条）；HTTP 层 `TestVouchersBatchDuplicateConflict`（第二次 409 + 列表 total 不变）。

### F-003 · PostgreSQL 上无 Redeem / 并发双花 e2e → **fixed**（PG e2e 落地并实证）
1. 新增 `internal/store/postgres_voucher_redeem_test.go`（`TestPostgresVoucherRedeemAndConcurrentSubject`，`PG_TEST_*` 门控，CI postgres job 点亮）：同一**新主体**并发核销两张不同卡（500+1500）→ 单 subject 账户行、余额 2000、账本恰 2 条 voucher 入金；重复核销同码 → `ErrVoucherAlreadyRedeemed` 且不双记。
2. 实证（本会话 docker postgres:15-alpine + `PG_TEST_*`）：`TestPostgresVoucherRedeemAndConcurrentSubject` **PASS 2.18s**；`TestPostgresWalletVoucherAndSubject0064` **PASS 1.97s**（0065 PG 方言随全目录自动应用，隐含验证）。

### F-005 · 过期字段是 Unix 秒 `inputNumber`，易填毫秒/留空 → **fixed**（秒范围 fail-closed + 标签提示）
1. `handler/wallet.go` 生成端点校验 `expiresAt`：提供值必须在 Unix **秒** `[1_000_000_000, 4_102_444_800]`（2001-09-09 ～ 2100-01-01），毫秒（~1.7e12）与越界值 → 400 `INVALID_VOUCHER_PARAMS`；≤0/缺省保持「无过期」原语义。`errorcatalog` 的 `INVALID_VOUCHER_PARAMS` en/zh 文案并入该范围说明。
2. UI 提示：schema 标签与 en/zh i18n 改为「Unix 秒，2001-09-09 至 2100-01-01；可选」。
3. 测试：`TestVouchersGenerateExpiresAtValidation`（毫秒 400 / 远古秒 400 / 合法未来秒 201）。

### F-002 · 导出是渲染器对 `items[].code` 的启发式，非协议声明的 download → **fixed**（协议声明驱动）
1. `apps/api/modules/wallet/schema/wallet-vouchers.json`：`generateBatch.onSuccess` 显式声明 `downloadCsv`（列序 + `vouchers_{batchId}.csv` 文件名模板）——导出 intent 进入协议文档。
2. `apps/web/src/renderer/render.tsx` `submitForm`：删除「任何成功响应含 `items[].code` 就下载」的启发式；仅当提交动作 `onSuccess.downloadCsv` 声明存在（含 `code` 列与文件名模板）才在同一手势生成 CSV 下载，列/表头由声明驱动（CSV 正文形状与历史一致）。
3. 测试：既有 CSV 用例改声明化 fixture（文件名断言不变）；新增阴性用例——无关表单响应带 `items[].code` 且**无声明** → 不触发下载。

## 验证（exit 0）

```text
go test ./internal/composition ./internal/handler ./internal/store ./modules/wallet/... -count=1
ok internal/composition 23.9s / handler 40.3s / store 44.3s / modules/wallet 全部 ok（含新增服务层与 HTTP 用例）
PG e2e（docker postgres:15 + PG_TEST_*）：TestPostgresVoucherRedeemAndConcurrentSubject PASS 2.18s
npx vitest run src/renderer：22 files / 290 tests PASS（apps/web，含 2 个 CSV 声明化用例）
```

未改动上游协议 pin/Profile/Manifest；无新默认模块。frozen 错误码清单与迁移目录台账按惯例同步（见修复事实 F-004）。

## 产物

- 后端：`modules/wallet/migration/migration.go`（0065）、`voucher/{voucher,service}.go`、`internal/handler/wallet.go`、`internal/errorcatalog/errorcatalog.go`、`internal/handler/error_contract_test.go`、`internal/store/postgres_voucher_redeem_test.go`（新增）及目录台账测试同步
- 前端：`apps/web/src/renderer/render.tsx`、`render.test.tsx`、i18n en/zh、`modules/wallet/schema/wallet-vouchers.json`
- 记录：`03-audit/A-008-a007-closure-response.md`

## Git checkpoint

- fix 提交：`bcd87ff1`（19 files，owned paths = F-002～F-005 代码/测试/目录台账清单；验证 = 后端全回归 + PG e2e + renderer vitest 全绿）
- docs 提交：本记录所在 docs 笔（E-006/A-007/A-008 落盘 + 三个索引同步）
