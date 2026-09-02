---
id: GOAL-004-evidence-closure-and-closeout
doc: execution-entry
record_id: E-002
status: recorded
parent: GOAL-004-evidence-closure-and-closeout
created: 2026-09-02
updated: 2026-09-02
version: 0.1.0
---

# E-002 · R4 证据矩阵核对、越界核账与全量实证事实

## 2026-09-02 · 证据闭环

### 已发生事实

1. **VP-029 七条方向级退出判据逐条实证核验**：
   - **判据 #1（主体接缝可用）**：`(issuer, external_id) -> subject_id` 幂等登记与查找验证通过；未登记主体开户被严格拒绝（返回 `ErrNotFound`）；不创建 `admin.users`；查询/登记不依赖钱包启用。测试证据：`modules/wallet/subject/subject_test.go` + `modules/wallet/provider_test.go`。
   - **判据 #2（凭证生命周期）**：高熵码（24 字符 Crockford Base32，120 bit 熵）生成、SHA-256 哈希存储、一次性出示明文、作废与过期拒绝验证通过；明文绝不落库、绝不进审计原文。测试证据：`modules/wallet/voucher/voucher_test.go`（`TestNoPlaintextInDatabase`）+ `handler/wallet_voucher_test.go`。
   - **判据 #3（核销原子且幂等）**：`Redeem` 单事务 CAS 标记并调账入金；重复核销返回 `ErrVoucherAlreadyRedeemed` 不双记；文件 SQLite（WAL + 连接池）20 并发真实竞争实测恰好 1 成功 19 拦截，余额与流水精确无误。测试证据：`modules/wallet/voucher/voucher_test.go`（`TestConcurrentDoubleSpendFailClosed`）。
   - **判据 #4（账本不变式保持）**：复用 `adjust` 账本原语（`ref_type='voucher'`, `idempotency_key=voucher.id`），三余额恒等、流水快照链与对账 Job 依然全部通过既有测试。测试证据：`modules/wallet/store/repository_test.go`（`TestSubjectAccount`）。
   - **判据 #5（Admin 可操作）**：批次生成、列表查询、详情、作废 4 条 Admin HTTP 路由交付；协议驱动页面 `wallet-vouchers.json` 注册并由 manifest/fragment 与 provider 声明；细粒度权限键 `wallet.voucher.issue` 严格鉴权并与 `wallet.adjust` 隔离；操作审计记录留痕无明文。测试证据：`handler/wallet_voucher_test.go` + `modules/wallet/schema/schema.go`。
   - **判据 #6（边界保持）**：越界核账确认未改 Charter，未改 Profile 默认模块集，未引入支付网关或 Telegram 依赖，未重开 VP-011。
   - **判据 #7（审计闭合）**：GOAL-002 A-002（3 required 全部 fixed）、GOAL-003 A-002（2 required 全部 fixed），工作区 29 全域 **open required finding = 0**。
2. **全量实证回归**：
   - `go test ./modules/wallet/...` PASS。
   - `go test ./internal/store/...` PASS。
   - `go test ./internal/handler -run TestVoucher` PASS。

### 证据

| 判据 / 领域 | 证据文件与测试 | 判定 |
|-------------|----------------|------|
| 判据 #1 | `modules/wallet/subject/subject_test.go` + `provider_test.go` | PASS |
| 判据 #2 | `modules/wallet/voucher/voucher_test.go` + `handler/wallet_voucher_test.go` | PASS |
| 判据 #3 | `modules/wallet/voucher/voucher_test.go`（WAL 20 并发防双花测试） | PASS |
| 判据 #4 | `modules/wallet/store/repository_test.go` | PASS |
| 判据 #5 | `handler/wallet_voucher_test.go` + `modules/wallet/schema/`（协议页面 + 权限拦截与审计无明文测试） | PASS |
| 判据 #6 | `git diff --stat origin/dev`（红线零触碰核账） | PASS |
| 判据 #7 | GOAL-002 A-003 + GOAL-003 A-003 | PASS |
