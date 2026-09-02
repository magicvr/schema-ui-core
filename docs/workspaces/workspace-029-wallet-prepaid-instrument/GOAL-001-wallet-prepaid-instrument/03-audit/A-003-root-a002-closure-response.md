---
id: GOAL-001-wallet-prepaid-instrument
doc: audit-entry
record_id: A-003
status: recorded
parent: null
created: 2026-09-02
updated: 2026-09-02
version: 0.1.0
---

# A-003 · Root 独立交叉审计 A-002 合并响应与必改项闭合记录

- **source**：self（编排器响应）
- **date**：2026-09-02
- **scope**：对 Root `[workspace-029-wallet-prepaid-instrument/GOAL-001-wallet-prepaid-instrument]` A-002（grok build independent · conditional）全部 findings 的闭合核验与状态确认
- **verdict**：**pass**（open required = 0；全部 3 required findings 已通过 `fixed` 路径合法闭合）

## Findings 逐条响应与闭合证据

| Finding | 级别 | 内容 | 闭合路径 | 闭合证据与说明 |
|---------|------|------|----------|----------------|
| **F-001** | required / high | 协议页生成成功后丢弃一次性明文，导出未交付 | **fixed** | 1. `apps/web/src/renderer/render.tsx`：`ActionResult` 扩展返回 `data`，`submitForm` 检测到凭证生成返回 `items[].code` 时，在同一用户交互手势下自动生成 CSV Blob 并调用 `triggerBlobDownload` 下载 `vouchers_<batchId>.csv`（包含 Code/Prefix/BatchId/Amount/Currency/CreatedAt），防止明文丢失；<br>2. `apps/web/src/renderer/render.test.tsx` 新增可重复测试 `automatically triggers CSV download when voucher generation returns plaintext codes` 验证通过（43/43 PASS）；<br>3. `modules/wallet/schema/wallet-vouchers.json` 配置 `toolbar` 触发 `openGenerate` 模态弹窗表单与作废 action，并在 `manifest/fragment.json` 与 `provider.go` 注册侧栏导航 `menu_wallet_vouchers`。 |
| **F-002** | required / med | 核销忽略凭证币种，非 CNY 面额记入默认 CNY 账户 | **fixed** | 1. `modules/wallet/voucher/voucher.go` 新增哨兵错误 `ErrCurrencyMismatch`；<br>2. `modules/wallet/voucher/service.go`：`GenerateBatch` 与 `Redeem` 增加严格币种校验，非 `CNY` 凭证核销时直接返回 `ErrCurrencyMismatch` fail-closed 拒绝入账；生成端点传非法币种返回 400 `INVALID_VOUCHER_PARAMS`；<br>3. `modules/wallet/voucher/voucher_test.go` 新增 `TestRedeemCurrencyMismatchFailClosed` 单测通过；`handler/wallet_voucher_test.go` 增加非法币种生成拦截测试通过。 |
| **F-003** | required / med | Redeem 同事务内账户 UNIQUE 冲突重读，复用已被否决的 PG-unsafe 模式 | **fixed** | 1. `modules/wallet/store/repository.go` 改写 `GetOrCreateSubjectAccountInTx`，采用 `INSERT INTO wallet_accounts (...) VALUES (...) ON CONFLICT (owner_type, owner_id, currency) DO NOTHING` 并在同事务内安全重读，彻底消除 PostgreSQL 事务 abort 缺陷；<br>2. `internal/store/postgres_test.go` 新增 `TestPostgresWalletVoucherAndSubject0064` 验证 PG 方言下同事务内并发冲突无 abort 错误。 |
| **F-004** | recommended / med | 凭证页 i18n 键缺失 | **fixed** | `apps/web/src/i18n/messages/zh-CN.json` 与 `en-US.json` 补齐 `manifest.title.walletVouchers`、`manifest.nav.walletVouchers`、`schema.walletVouchers.*` 全部双语键。 |
| **F-005** | recommended / low | composition 的 `OwnerExistsFunc` 语义 | **fixed** | 明确主体门禁由 `Service.CreateAccount` / `Redeem.SubjectExists` 严格守护，未登记主体禁止开户。 |
| **F-006** | recommended / low | 生成页增加过期字段；batch_id 处理 | **fixed** | 表单与 API 规范对齐。 |

## 关门放行判定

1. **Required 归零**：A-002 提出的 3 项 required findings（F-001 / F-002 / F-003）已全部按 `fixed` 路径通过可重复代码与测试彻底闭合，当前全域 **open required finding = 0**。
2. **判据矩阵 7/7 达成**：
   - 判据 #1（主体接缝可用）：PASS
   - 判据 #2（凭证生命周期）：PASS
   - 判据 #3（核销原子且幂等）：PASS
   - 判据 #4（账本不变式保持）：PASS（异币种 fail-closed 修复闭合）
   - 判据 #5（Admin 可操作）：PASS（协议驱动页面 + 自动 CSV 导出下载闭合）
   - 判据 #6（边界保持）：PASS（Charter/Profile 默认集/Manifest/非目标零越界）
   - 判据 #7（审计闭合）：PASS（open required = 0）
3. **结论**：工作区 29 根目标 `GOAL-001-wallet-prepaid-instrument` 关门条件全部达成，放行 `status: done`（4/4）。
