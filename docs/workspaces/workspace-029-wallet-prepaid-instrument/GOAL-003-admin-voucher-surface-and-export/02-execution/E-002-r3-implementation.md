---
id: GOAL-003-admin-voucher-surface-and-export
doc: execution-entry
record_id: E-002
status: recorded
parent: GOAL-003-admin-voucher-surface-and-export
created: 2026-09-02
updated: 2026-09-02
version: 0.1.0
---

# E-002 · R3 Admin 批次管理与生命周期实施及独立审计整改

## 2026-09-02 · 实施与整改

### 已发生事实

1. **HTTP 端点与权限交付**：
   - 在 `apps/api/internal/handler/wallet.go` 交付 4 条 Admin 路由：
     * `POST /api/wallet/vouchers/batches`（权限 `wallet.voucher.issue`）：批量生成高熵凭证，一次性返回明文列表供导出；
     * `GET /api/wallet/vouchers`（权限 `wallet.read`）：列表查询，仅出示 `codePrefix`，不出示卡密明文；
     * `GET /api/wallet/vouchers/{id}`（权限 `wallet.read`）：单张详情；
     * `POST /api/wallet/vouchers/{id}/void`（权限 `wallet.voucher.issue`）：作废单张未核销凭证。
   - `modules/wallet/provider.go` 与 `kernel/profile.go` 注册 4 条路由与 `wallet.voucher.issue` 权限贡献；RBAC 默认只授 `PolicyAdmin`。
2. **审计日志与安全防御**：
   - 批次生成与作废分别记录 `records.create`（`batch-generate`）和 `records.update`（`voucher-void`）操作日志。
   - 响应 A-002 F-001：移除 `INTERNAL` 错误响应中的底层 `err.Error()` 驱动原文，统一采用泛化错误文案（`could not generate voucher batch` / `could not get voucher` / `could not void voucher`），杜绝诊断信息泄露。
   - 断言审计日志 Payload 绝对不含任何卡密明文。
3. **错误码与回归锁强化**：
   - `errorcatalog.go` 与 `error_contract_test.go` 注册并冻结 5 个凭证错误码：`INVALID_VOUCHER_BODY`, `INVALID_VOUCHER_PARAMS`, `INVALID_VOUCHER_ID`, `VOUCHER_NOT_FOUND`, `VOUCHER_ALREADY_REDEEMED`。
   - 响应 A-002 F-003：在 `wallet_voucher_test.go` 增加针对仅有 `wallet.adjust` 权限用户的测试，验证其调用发卡接口被严格返回 403 Forbidden，锁死权限隔离。
   - 响应 A-002 F-004：增加 `TestVoucherVoidAndConflict` 与 `TestVoucherInvalidBodyAndParams`，完整覆盖 400（参数错误）、404（不存在）、409（作废已核销券）的 HTTP 包络。
4. **全量回归验证**：
   - `go test ./internal/handler -run TestVoucher` PASS。
   - `go test ./modules/wallet/...` PASS。
   - `go test ./internal/store/...` PASS。

### 证据

| 主张 | 路径 / 命令 / commit |
|------|----------------------|
| 路由与鉴权实现 | `apps/api/internal/handler/wallet.go` + `apps/api/modules/wallet/provider.go` |
| 权限与错误码契约 | `apps/api/kernel/profile.go` + `apps/api/internal/errorcatalog/errorcatalog.go` + `error_contract_test.go` |
| 权限隔离与异常测试 | `apps/api/internal/handler/wallet_voucher_test.go` |
| 模块服务与存储 | `apps/api/modules/wallet/voucher/service.go` |
