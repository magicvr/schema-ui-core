---
id: GOAL-003-admin-voucher-surface-and-export
doc: audit-entry
record_id: A-002
status: recorded
parent: GOAL-003-admin-voucher-surface-and-export
created: 2026-09-02
updated: 2026-09-02
version: 0.1.0
---

# A-002 · R3 Admin 批次管理与生命周期独立交叉审计（2026-09-02）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型**：execution-facts（对照 A-001 关门自审主张）
- **scope**：`[workspace-029-wallet-prepaid-instrument/GOAL-003-admin-voucher-surface-and-export]` —— 批次生成/查询/作废 HTTP API、`wallet.voucher.issue` 鉴权、明文一次性返回与防泄露审计、错误码契约、回归测试
- **verdict**：**conditional**
- **完整意见**：本条由编排器自本地 grok build（grok-4.6 · reasoning high · `/audit`）独立审计会话原样誊入，`source: independent` 保持不变

### 成果（有证据）

1. **四条 HTTP 路由已挂到 `admin.wallet`**  
   `apps/api/internal/handler/wallet.go`：`POST /api/wallet/vouchers/batches`、`GET /api/wallet/vouchers`、`GET /api/wallet/vouchers/{id}`、`POST /api/wallet/vouchers/{id}/void`。`modules/wallet/provider.go` Descriptor Routes/Permissions 与 `kernel/profile.go` `admin.wallet` ContributionKeys 同步含上述 4 路由 + `wallet.voucher.issue`。RBAC 默认只授 `PolicyAdmin`。

2. **权限键在源码层按 D-002 / D-001 接线**  
   生成与作废走 `requirePermission(..., "wallet.voucher.issue")`。列表与详情走 `wallet.read`。Viewer 无发卡权，Admin 成功。

3. **明文一次性返回；库表无明文列；查询面不回明文/哈希**  
   仅生成响应包含 `code`。列表/详情手工组包，不含 `code`/`codeHash`。`vouchers` 表只有 `code_hash`/`code_prefix`。审计详情断言不含卡密明文。

4. **作废 CAS**  
   `VoidVoucher`：已核销 → `ErrVoucherAlreadyRedeemed`；已作废幂等成功；未用 → `UPDATE ... AND status='unused'`。Handler 映射已核销为 409 `VOUCHER_ALREADY_REDEEMED`，成功写 `records.update` / `voucher-void`。

5. **错误码已编入冻结契约与双语目录**  
   `errorcatalog.go` + `error_contract_test.go` frozen literals：`INVALID_VOUCHER_BODY`、`INVALID_VOUCHER_PARAMS`、`INVALID_VOUCHER_ID`、`VOUCHER_NOT_FOUND`、`VOUCHER_ALREADY_REDEEMED`。

### Findings

#### F-001 · required · med
**INTERNAL 诊断泄露（生成 / 详情 / 作废）**：未编目码直接返回 `err.Error()`，应改为泛化文案，禁止底层驱动/数据库原文出站泄露。

#### F-002 · required · med
**执行台账补齐与证据链闭环**：补充 `02-execution/E-002-r3-implementation.md` 记录实施事实，在 required findings 闭合前不提前勾选 4/4。

#### F-003 · recommended · med
**回归锁锁定 `wallet.voucher.issue` 与 `wallet.adjust` 权限隔离**：测试需断言仅有 `wallet.adjust` 的用户被拦截（403）。

#### F-004 · recommended · med
**错误码与作废失败路径补齐 HTTP 行为测试**：覆盖 400（参数错误）、404（不存在）、409（作废已核销券）等。

#### F-005 · recommended · low
**导出描述校准**：准确记录导出仅为 API 一次性返回明文数组，不预填未交付的 UI。
