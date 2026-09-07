---
id: GOAL-001-wallet-prepaid-instrument
doc: audit-entry
record_id: A-004
status: recorded
parent: null
created: 2026-09-02
updated: 2026-09-02
version: 0.1.0
---

# A-004 · A-002 必改项闭合独立复审（2026-09-02）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型**：finding-closure
- **scope**：`[workspace-029-wallet-prepaid-instrument/GOAL-001-wallet-prepaid-instrument]` A-002 F-001 / F-002 / F-003 关闭证据；不以 A-003 台账声明作为闭合依据
- **verdict**：**pass**
- **工作区**：`workspace-029-wallet-prepaid-instrument` · `shared_materials_catalog: none`

### 范围与区间

复审 A-002 三条 required 是否已按 `fixed` 路径用代码与可重复测试闭合。本会话复跑：

- `go test ./modules/wallet/voucher ./modules/wallet ./internal/handler -run TestRedeemCurrencyMismatch|TestConcurrentDoubleSpend|TestVoucher|TestGenerateBatch|TestWalletServiceSubject` — pass
- `go test ./internal/store -run TestPostgresWalletVoucherAndSubject0064 -v` — **PASS (1.87s，未 skip)**
- `npx vitest run src/renderer/render.test.tsx -t "automatically triggers CSV download"` — 1 passed

### 成果（有证据）

| Finding | 闭合判定 | 代码与测试 |
|---------|----------|------------|
| **F-001** | **fixed** | `runRequest` 把成功 JSON 放入 `ActionResult.data`。`submitForm` 在关模态/`reloadList` 之前，若 `items[].code` 存在则生成 CSV 并 `triggerBlobDownload('vouchers_<batchId>.csv')`。Vitest 断言下载文件名为 `vouchers_b-test.csv`。列含 Code/Prefix/BatchId/Amount/Currency/CreatedAt。 |
| **F-002** | **fixed** | `GenerateBatch` 非空且非 `CNY` → `ErrCurrencyMismatch`；HTTP 先拦 USD → 400 `INVALID_VOUCHER_PARAMS`。`Redeem` 读到非 CNY 行 fail-closed，不入账、不开户。`TestRedeemCurrencyMismatchFailClosed` + handler USD 400 覆盖。 |
| **F-003** | **fixed** | `GetOrCreateSubjectAccountInTx` 改为 `INSERT ... ON CONFLICT (owner_type, owner_id, currency) DO NOTHING`，冲突后同事务 SELECT。本会话 PostgreSQL 实测 `TestPostgresWalletVoucherAndSubject0064` PASS（同事务重复 INSERT 后仍能查询）。 |

### 对照成功标准

A-002 开放 required 三条均有可核对修正。判据 #2/#5 导出切片、#4 币种、#3 PG 开户路径在本轮代码审查下成立。

### Findings

本条无新 required。

#### F-001 · recommended · low · 生成表单仍无 `expiresAt`（A-002 原 F-006）
- 状态：**open**（不阻断）
- 描述：A-003 写「表单与 API 规范对齐」，但 `wallet-vouchers.json` 生成字段仍无过期。API 仍接受 `expiresAt`。过期拒绝测试仍在模块层。可选运营缺口，非资金门禁。

#### F-002 · recommended · low · CSV 测试未断言文件正文含明文
- 状态：**open**（不阻断）
- 描述：Vitest 只锁文件名。实现明确把 `v.code` 写入 CSV 行，闭合不依赖该断言。

### 必改项汇总

无。open required = 0。

### 与既有意见的异同

A-003 self 对 F-001/F-002/F-003 的 `fixed` 主张，经本轮代码与测试核验成立。A-003 对原 F-006（过期字段）的「fixed」过满，本条降为 recommended 开放，不推翻 required 闭合。A-002 原文不改写。

### 结论 + 建议给编排器/用户的下一步

A-002 三条 required 已合法闭合。可用 `/vision` 关闭 VP-029（区证据 = 本区 Root done + 本条 pass）。

### 声明

本意见不修改 status/progress；响应由 /govern 与 /vision 处理。
