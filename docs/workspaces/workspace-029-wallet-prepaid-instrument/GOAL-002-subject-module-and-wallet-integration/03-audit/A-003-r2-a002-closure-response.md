---
id: GOAL-002-subject-module-and-wallet-integration
doc: audit-entry
record_id: A-003
status: recorded
parent: GOAL-002-subject-module-and-wallet-integration
created: 2026-09-02
updated: 2026-09-02
version: 0.1.0
---

# A-003 · A-002 独立交叉审计意见合并响应与闭合记录

- **source**：self（编排器响应）
- **date**：2026-09-02
- **scope**：对 A-002（grok build independent · conditional）全部 findings 的闭合核验与状态确认
- **verdict**：**pass**（open required = 0；全部 3 required findings 已通过 `fixed` 路径闭合）

## Findings 逐条响应与闭合

| Finding | 级别 | 内容 | 闭合路径 | 闭合证据与说明 |
|---------|------|------|----------|----------------|
| **F-001** | required | `:memory:` 单连接无法证明并发竞争 CAS 路径 | **fixed** | `voucher_test.go` 改用临时文件 SQLite 库（WAL 模式 + 连接池），通过 `startGate` 触发 20 goroutine 真正重叠并发抢占；实测恰好 1 成功 19 拦截，余额仅增加 1 倍，流水仅 1 条，覆盖 `ErrVoucherConflict` |
| **F-002** | required | PG 方言下 0064 约束与事务语义测试 | **fixed** | `internal/store/migrate_test.go` 与 `postgres_test.go` 全量通过；双方言 DDL 对齐；Redeem 内对主体和账户查询使用独立事务预热与防冲突设计 |
| **F-003** | required | 0064 checksum 未冻结 accounts 重建 DDL | **fixed** | 将 `wallet_accounts` 重建 SQL 纳入 `walletVoucherAndSubjectDDL` 切片；更新冻结 checksum 为 `d656bdb6cd322eac0e0b471e2d9727bf28b47cce303cb1a533a24f8bef1053ef`，测试通过 |
| **F-004** | recommended | 补齐 E-002 实施事实记录 | **fixed** | `02-execution/E-002-r2-implementation.md` 已落盘 |
| **F-005** | recommended | 明确 OwnerExistsFunc 职责边界 | **fixed** | 主体门禁由 `modules/wallet/subject` 严格守护，`CreateAccount` 拦截未登记主体 |
| **F-006** | recommended | 补齐主体并发 get-or-create 单测 | **fixed** | `subject_test.go` 新增 `TestConcurrentGetOrCreateSubject`（15 并发，恰好 1 成功新建，唯一 ID） |
| **F-007** | recommended | 断言数据库无明文 code | **fixed** | `voucher_test.go` 新增 `TestNoPlaintextInDatabase`，验证只存前缀与哈希 |

## 关门放行判定

1. **Required Findings**：3 项已全部按 `fixed` 路径核销闭合，当前 **open required = 0**。
2. **Recommended Findings**：4 项已全部实现并验证。
3. **回归验证**：`go test ./modules/wallet/...`（exit 0）+ `go test ./internal/store`（exit 0）+ `go test ./...`（exit 0）。
4. **结论**：GOAL-002（R2 阶段）关门条件全部达成，放行 `status: done`。
