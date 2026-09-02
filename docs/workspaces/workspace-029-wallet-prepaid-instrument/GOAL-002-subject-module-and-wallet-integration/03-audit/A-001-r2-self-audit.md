---
id: GOAL-002-subject-module-and-wallet-integration
doc: audit-entry
record_id: A-001
status: recorded
parent: GOAL-002-subject-module-and-wallet-integration
created: 2026-09-02
updated: 2026-09-02
version: 0.1.0
---

# A-001 · R2 主体接缝与账本入金关门自审（self）

- **source**：self（编排器自审；independent 意见由 A-002 本地 grok build 出具）
- **date**：2026-09-02
- **scope**：GOAL-002-subject-module-and-wallet-integration 全量——迁移 0064、主体接缝存储与服务、钱包 `owner_type='subject'` 扩展、预付凭证高熵生成与单事务 CAS 原子核销入金、并发防双花、单元测试与回归测试
- **verdict**：**pass**（open required = 0；待 A-002 grok build independent 复核）

## 成功标准逐条对照

| 成功标准 | 判定 | 证据 |
|----------|------|------|
| 1. 迁移 0064 建立 subjects / vouchers 表，双方言兼容 | pass | `apps/api/modules/wallet/migration/migration.go` 新增版本 64（checksum: a91f267ae5cb954a20b1f015246786b94d3ec467d2622bc34f98f084e328ad49）；`internal/store` 全量测试通过（双方言建表/约束更新验证） |
| 2. 外部主体接缝可用，不建 admin.users，未登记主体不能开户 | pass | `apps/api/modules/wallet/subject/` 实现 `GetOrCreateSubject`/`SubjectExists`；`CreateAccount` 严格校验主体存在性；`subject_test.go` 与 `provider_test.go` 验证未登记主体开户直接返回 ErrNotFound |
| 3. 预付凭证核心服务：高熵码生成、SHA-256 哈希存储、一次性明文 | pass | `apps/api/modules/wallet/voucher/`：24 字符 Crockford Base32 编码（120 bit 熵 > 80 bit 底线），SHA-256 哈希存储；明文仅在 `GenerateBatch` 返回一次，不入库不进审计；测试全绿 |
| 4. 核销原子且防并发双花，账本不变式保持 | pass | `voucher.Service.Redeem`：单事务内 CAS 更新状态为 `redeemed`，影响行数=0 则 fail-closed，同事务内调用钱包入金（`ref_type='voucher'`，`idempotency_key=voucher.id`）；20 并发抢占实测恰好 1 成功 19 拦截，余额仅增加 1 倍，流水仅 1 条 |

## 验证事实与回归

1. `go test ./modules/wallet/...` 全绿（涵盖 store, subject, voucher, provider）。
2. `go test ./internal/store` 全绿（涵盖全部 64 个迁移与方言校验）。
3. `go test ./...` 全绿（无任何回归故障）。

## 结论

GOAL-002 成功标准 4/4 达成，open required = 0，自审 verdict 为 **pass**。准备启动本地 grok build（grok-4.6 · high）执行独立交叉审计（A-002）。
