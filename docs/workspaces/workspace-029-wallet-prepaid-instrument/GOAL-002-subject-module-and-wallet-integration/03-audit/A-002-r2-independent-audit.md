---
id: GOAL-002-subject-module-and-wallet-integration
doc: audit-entry
record_id: A-002
status: recorded
parent: GOAL-002-subject-module-and-wallet-integration
created: 2026-09-02
updated: 2026-09-02
version: 0.1.0
---

# A-002 · R2 主体接缝与预付凭证账本入金独立交叉审计（2026-09-02）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型**：execution-facts
- **scope**：`[workspace-029-wallet-prepaid-instrument/GOAL-002-subject-module-and-wallet-integration]` R2——迁移 0064、主体存储与校验、凭证生成与核销单事务 CAS、并发防双花及测试
- **verdict**：**conditional**
- **完整意见**：本条由编排器自本地 grok build（grok-4.6 · reasoning high · `/audit`）独立审计会话原样誊入，`source: independent` 保持不变

### 成果（有证据）

1. **迁移 0064 已编入全局台账（SQLite 路径可重复核对）**  
   `apps/api/modules/wallet/migration/migration.go` 贡献 `version=64` / `wallet_voucher_and_subject`；checksum 与 `internal/store/migrate_test.go` 冻结值一致。`subjects` / `vouchers` 含 `UNIQUE(issuer, external_id)`、`UNIQUE(code_hash)`；SQLite 以 rename→create→copy→drop 扩展 `wallet_accounts.owner_type` 含 `'subject'`。`TestMigrateFreshDB` 断言 64 条迁移以本条结尾，且 `subjects`/`vouchers` 表存在。

2. **主体接缝（幂等 get-or-create + 未登记不开户）**  
   `modules/wallet/subject`：`GetOrCreateSubject` / `SubjectExists`；空 issuer/external_id 拒绝；UNIQUE 冲突后重读。`CreateAccount(owner_type=subject)` 走 `SubjectExists`，不回退 `admin.users`。`provider_test.go`：未登记主体 → `ErrNotFound`；登记后可开户。`compiled/persistence.go` 声明运行时 Profile 不过滤 persistence，故表结构不依赖 `HasModule("admin.wallet")`。未见 Telegram SDK 依赖；未见公开 HTTP `Redeem`。

3. **凭证生成：高熵 + SHA-256 + 明文不入库**  
   `GenerateCode`：Crockford 类 32 字母表 × 24 字符 ≈ 120 bit（高于 D-002 >80 bit）；`HashCode` = SHA-256 hex；库列只有 `code_hash`/`code_prefix`，无 plaintext 列。`GenerateBatch` 单事务插入，明文只在返回值出现一次。

4. **核销 CAS 设计（SQLite 顺序路径已测通）**  
   `voucher.Service.Redeem`：同一 `runner.Run` 内 `UPDATE ... WHERE id=? AND status='unused'`，`RowsAffected==0` → `ErrVoucherConflict`；随后 `GetOrCreateSubjectAccountInTx` + `MutateInTx`（`entry_type=adjust`，`ref_type=voucher`，`idempotency_key=voucher.id`）。重复核销 → `ErrVoucherAlreadyRedeemed`，余额不双记。作废/过期/未知主体均有单测。本独立审复跑：`go test ./modules/wallet/... ./internal/store -count=1` 通过。

5. **账本不变式（复用 adjust，未新开 entry_type）**  
   符合 D-002 / I-029-002。`Apply` 仍走既有 `adjust`；`TestRedeemSuccess` / `TestSubjectAccount` 核三余额。未引入新 `entry_type`，不构成重开 VP-011。

### 对照成功标准

| 标准 | 判定 | 说明 |
|------|------|------|
| 1. 迁移 0064 双方言兼容 | **conditional** | SQLite 建表/CHECK 重建有测试。PG 为 `DROP CONSTRAINT IF EXISTS wallet_accounts_owner_type_check` + 再 ADD，依赖隐式约束名；本会话无 `PG_TEST_*`，voucher/Redeem 无 PG 用例 |
| 2. 主体接缝可用，不建 admin.users，未登记不能开户 | **pass（模块 API）** | get-or-create 与 CreateAccount 门禁有测试。composition 的 `OwnerExistsFunc` 仍是 `UserByID`（见 F-006，recommended） |
| 3. 高熵生成、SHA-256、一次性明文 | **pass（实现审查）** | schema/INSERT 无明文列；缺「SELECT 断言库中无 code」测试（F-007） |
| 4. 单事务 CAS、并发双花 fail-closed、三余额/流水回归 | **conditional** | 顺序路径与 CAS SQL 可审。自称「20 并发」的测试在 `:memory:` 单连接下被串行化，不能证明竞态（F-001） |

### Findings

#### F-001 · required · high · 关联 I-029-006 / 成功标准 4
**并发防双花测试未形成重叠事务，不能作为 R2 资金门禁证据。**
`:memory:` 单连接下事务串行化，需使用文件库（WAL、pool>1）重做并发双花测试，证明恰好 1 成功、账本 +1×面额、流水 1 条，并覆盖 CAS `affected=0`（`ErrVoucherConflict`）。

#### F-002 · required · medium · 关联 I-002-001 / 成功标准 1、4
**PG 方言上 0064 CHECK 扩展与 Redeem 原子核销缺少可重复证据；且存在已知同事务 UNIQUE 处理缺陷。**
强化 PG 方言下的约束处理与 Redeem 事务语义。

#### F-003 · required · medium · 关联 I-002-001
**0064 checksum 未冻结 `wallet_accounts` 重建 SQL（owner_type CHECK 扩展）。**
将 accounts 重建 DDL 纳入 0064 冻结 checksum，与 0033 惯例对齐。

#### F-004 · recommended
执行台账补齐 E-002 事实记录。

#### F-005 · recommended
明确 `OwnerExistsFunc` 与主体校验的结构。

#### F-006 · recommended
补齐 `GetOrCreateSubject` 并发测试。

#### F-007 · recommended
补齐 SELECT 断言库中无明文 code 的测试。
