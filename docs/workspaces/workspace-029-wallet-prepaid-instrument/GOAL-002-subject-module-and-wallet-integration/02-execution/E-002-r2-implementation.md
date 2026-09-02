---
id: GOAL-002-subject-module-and-wallet-integration
doc: execution-entry
record_id: E-002
status: recorded
parent: GOAL-002-subject-module-and-wallet-integration
created: 2026-09-02
updated: 2026-09-02
version: 0.1.0
---

# E-002 · R2 实施、并发验证与独立审计整改闭合

## 2026-09-02 · 实施与整改

### 已发生事实

1. **迁移 0064 全局注册**：
   - 建立 `subjects` 表与 `vouchers` 表，双方言 DDL 严格对齐。
   - SQLite 使用 rename→create→copy→drop 机制安全升级 `wallet_accounts.owner_type` CHECK，增加 `'subject'`。
   - 将 accounts 重建 DDL 纳入 0064 冻结 checksum（`d656bdb6cd322eac0e0b471e2d9727bf28b47cce303cb1a533a24f8bef1053ef`），响应并闭合 A-002 F-003。
2. **外部主体接缝落地**：
   - 实现 `modules/wallet/subject/` 包，提供 `GetOrCreateSubject(ctx, issuer, externalID)` 与 `SubjectExists(ctx, subjectID)`。
   - 增加 15 并发抢占单测 `TestConcurrentGetOrCreateSubject`，验证天然幂等与单条创建。
   - `CreateAccount` 接入主体校验：未登记主体禁止开户（返回 `ErrNotFound`），禁止回退 `admin.users`。
3. **预付凭证核心服务与原子核销**：
   - 实现 `modules/wallet/voucher/` 包，支持 24 字符 Crockford Base32 编码（120 bit 熵），存储 SHA-256 哈希，明文仅在返回值一次性提供。
   - `voucher_test.go` 增加 `TestNoPlaintextInDatabase`，断言数据库内无明文。
   - `Redeem` 单事务 CAS 更新 `vouchers` 状态（`affected == 0` 时返回 `ErrVoucherConflict` fail-closed），同事务调账入金（`ref_type='voucher'`, `idempotency_key=voucher.id`）。
   - 实现临时文件 SQLite（WAL 模式 + 连接池）20 并发真实竞争测试 `TestConcurrentDoubleSpendFailClosed`，验证并发抢占恰好 1 成功 19 拦截，余额仅增加 1 倍，流水仅 1 条，响应并闭合 A-002 F-001/F-002。
4. **全量回归通过**：
   - `go test ./modules/wallet/...`（4 个包全部 PASS）。
   - `go test ./internal/store`（64 个迁移全部 PASS）。
   - `go test ./...`（全局全绿）。

### 证据

| 主张 | 路径 / 命令 / commit |
|------|----------------------|
| 迁移 0064 与 checksum | `apps/api/modules/wallet/migration/migration.go` + `apps/api/internal/store/migrate_test.go` |
| 主体接缝与并发测试 | `apps/api/modules/wallet/subject/subject.go` + `subject_test.go`（`TestConcurrentGetOrCreateSubject`） |
| 凭证服务与防并发双花 | `apps/api/modules/wallet/voucher/voucher.go` + `service.go` + `voucher_test.go`（`TestConcurrentDoubleSpendFailClosed` 文件库） |
| 账本不变式回归 | `apps/api/modules/wallet/store/repository_test.go`（`TestSubjectAccount`） |
| 全量测试 exit 0 | `go test ./modules/wallet/...` + `go test ./internal/store` |
