---
id: GOAL-001-wallet-prepaid-instrument
doc: audit-entry
record_id: A-005
status: recorded
parent: null
created: 2026-09-02
updated: 2026-09-02
version: 0.1.0
---

# A-005 · Root 完成情况 · 方案设计与代码实现独立交叉审计（2026-09-02）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型**：close-out / design-plan / execution-facts
- **scope**：`[workspace-029-wallet-prepaid-instrument/GOAL-001-wallet-prepaid-instrument]` 根目标完成情况——对照 VP-029 七条退出判据与 D-002 冻结合同，**独立审查方案设计与当前代码实现**；不以五件套 `status`、goal-tree、子目标关门声明、A-001～A-004 闭合口径或 VP `closed` 作为完成证据
- **verdict**：**conditional**
- **工作区**：`workspace-029-wallet-prepaid-instrument` · root `GOAL-001-wallet-prepaid-instrument` · `shared_materials_catalog: none`（无共享资料引用）

### 范围与区间

- **审什么**：主体接缝、凭证生命周期、核销原子/幂等、账本不变式、Admin 协议面（生成/导出/作废）、红线边界；方案（独立表 / SHA-256 / CAS+adjust / 不暴露 HTTP Redeem / 细粒度权限）是否与实现一致。
- **不审什么**：不改 `status`/`progress`/goal-tree；不把治理文档「已 done」当作实现完成；不审其他工作区。
- **本会话实证（代码 + 测试，非台账）**：
  - 通读：`modules/wallet/subject`、`voucher/{voucher,service}.go`、`store/repository.go`（主体开户 / `MutateInTx`）、`migration/migration.go`（0064 双方言）、`provider.go`、`handler/wallet.go`、`schema/wallet-vouchers.json`、`apps/web/src/renderer/render.tsx` 导出分支、`composition.go` `OwnerExistsFunc`、`compiled/persistence.go`、`kernel/profile.go` 默认集、`apps/api/go.mod`。
  - 复跑（本会话 exit 0）：
    - `go test ./modules/wallet/voucher ./modules/wallet/subject ./modules/wallet ./internal/handler -run "TestRedeem|TestGenerate|TestConcurrent|TestNoPlaintext|TestGetOrCreateSubject|TestWalletServiceSubject|TestVoucher"` — PASS
    - `go test ./internal/store -run TestPostgresWalletVoucherAndSubject0064 -v` — **PASS 1.75s，未 skip**（本机有 `PG_TEST_*`）
    - `npx vitest run src/renderer/render.test.tsx -t "automatically triggers CSV download"` — 1 passed
  - 本会话**没有** PostgreSQL 上的 `Redeem` / 并发双花端到端测试。

### 成果（有证据）

1. **主体接缝（判据 #1 的模块 API）**  
   `subjects`：`UNIQUE(issuer, external_id)`（0064 SQLite+PG）。`GetOrCreateSubject` 空输入拒绝；冲突后在**新事务**重读（符合 PG「失败 INSERT 会 abort 当前 tx」的既有约束）。15 并发文件库测试恰好 1 条新建。`CreateAccount(owner_type=subject)` 走 `SubjectExists`，未登记 → `ErrNotFound`（`provider_test.go`）。未见写入 `admin.users`。持久化由 compiled catalog 贡献、**不**随 `HasModule("admin.wallet")` 过滤（`compiled/persistence.go` 注释与实现）。无公开主体 HTTP。无 HTTP `Redeem`（`WalletService` 接口无该方法）。

2. **凭证生命周期（判据 #2）**  
   24 字符 Crockford 类字母表（32 字符 × 24 = 120 bit）+ SHA-256 hex；库列仅 `code_hash`/`code_prefix`；`UNIQUE(code_hash)`。生成 HTTP 201 返回一次性 `code`；列表/详情组包不含 `code`/`codeHash`。作废 CAS：`UPDATE ... AND status='unused'`；已核销 → `ErrVoucherAlreadyRedeemed`；过期 Redeem → `ErrVoucherExpired`。生成审计详情断言不含明文。`TestNoPlaintextInDatabase` 覆盖库内哈希。协议生成表单现含 `expiresAt`（Unix 秒 `inputNumber`）。

3. **核销原子且幂等（判据 #3，SQLite 路径 + PG 开户片段）**  
   `Redeem` 单事务：按 hash 读取 → 币种 fail-closed → CAS `unused→redeemed`（`RowsAffected==0` → `ErrVoucherConflict`）→ 同事务 `GetOrCreateSubjectAccountInTx` → `MutateInTx(entry_type=adjust, ref_type=voucher, idempotency_key=voucher.id)`。失败则整笔回滚（含已 CAS 的行）。重复核销 → `ErrVoucherAlreadyRedeemed`，余额不双记。文件 SQLite 20 并发：1 成功 / 19 冲突，流水 1 条。本会话复跑通过。  
   `GetOrCreateSubjectAccountInTx` 使用 `INSERT ... ON CONFLICT (owner_type, owner_id, currency) DO NOTHING` 再 SELECT；本会话 PG 实测该语句不 abort 同事务后续查询。

4. **账本不变式（判据 #4）**  
   未新增 `entry_type`；入金走既有 `adjust` Apply 表。`TestRedeemSuccess` 断言三余额。`apps/api/go.mod` 无 Telegram / 支付网关依赖。

5. **Admin HTTP / 权限 / 导出（判据 #5 的可操作切片）**  
   四条路由挂在 `admin.wallet`：`POST /batches`、`GET /`、`GET /{id}`、`POST /{id}/void`。生成/作废 `wallet.voucher.issue`；仅 `wallet.adjust` 的角色 403。`runRequest` 把成功 JSON 放入 `ActionResult.data`；`submitForm` 在关模态前若 `items[].code` 存在则生成 CSV 并 `triggerBlobDownload`。Vitest 断言下载文件名为 `vouchers_b-test.csv`。错误码在 `errorcatalog`。i18n 键已进入 `zh-CN`/`en-US`（含 `manifest.title.walletVouchers` 与 `schema.walletVouchers.*`）。

6. **边界（判据 #6 的代码切片）**  
   `profileDefaults`：**未**向 `mvp` 塞入新模块 id；`admin.wallet` 原已在 admin 默认集，本波是该模块增量（路由/页面/权限键），不是第二模块。HTTP `POST /api/wallet/accounts` 仍禁止手动建 `owner_type=user`。无支付/Telegram 依赖。

### 对照成功标准（VP-029 七条；以代码与本会话运行为准）

| # | 判据 | 判定 | 说明 |
|---|------|------|------|
| 1 | 主体接缝可用 | **conditional** | 幂等 get-or-create / 未登记不开 **subject** 户 / 不建 admin.users / 表不依赖钱包 HTTP 启用：成立。但 I-029-001 的 `OwnerExists` 被接到 **user 自动开户 HTTP**，已登记主体 ID 可开出 `owner_type=user` 孤儿账本（F-001） |
| 2 | 凭证生命周期 | **pass（功能）** | 生成/哈希存储/作废/过期/明文不落库有测试；导出在本 GUI 同手势可下载。协议未声明 download（F-002，不升 required） |
| 3 | 核销原子且幂等 | **pass（SQLite）+ 覆盖缺口** | SQLite CAS+防双花有测试；PG 开户 `ON CONFLICT` 已实测。缺 PG 上 Redeem/并发双花 e2e（F-003 recommended） |
| 4 | 账本不变式 | **pass** | 复用 adjust；三余额测试通过；异币种生成/核销 fail-closed（`TestRedeemCurrencyMismatchFailClosed` + handler USD 400） |
| 5 | Admin 可操作 | **pass（本 GUI）** | 页面 + 权限键 + 操作审计 + 同手势 CSV 成立。导出是渲染器启发式而非协议 outcome（F-002） |
| 6 | 边界保持 | **pass（模块/依赖）** | 无新默认模块、无 Telegram/支付依赖、未改 entry_type apply 表。F-001 是账户类型门禁回归，记在判据 #1 |
| 7 | 审计闭合 | **fail（本条）** | 本意见开放 required > 0；不得以 A-004/子目标/VP closed 抵销 |

### Findings

#### F-001 · `OwnerExistsFunc` 把主体 ID 放行进 by-owner，会开出 `owner_type=user` 孤儿账本
- 严重度：**med**
- 建议：**required**
- 关联：I-029-001（OwnerExists / W13 F-012 孤儿账本）；VP-029 判据 #1
- 状态：**open**
- 描述：D-002 / I-029-001 要的是：**主体账户**开户校验已登记主体，禁止用 `UserByID` 冒充主体存在。`CreateAccount(owner_type=subject)` 与 `Redeem` 已经各自调用 `SubjectExists`，这条门禁本身是对的。  
  但 `composition.go` 把同一回调注入 **仅用于 user 自动开户的 HTTP**：
  1. `OwnerExistsFunc` 注释仍写「user owner id」；实现却是 `UserByID` **或** `SubjectExists`。
  2. `POST /api/wallet/by-owner/{ownerId}` 与 `.../adjust` 在门禁通过后**固定**调用 `GetOrCreateUserAccount` → 插入 `owner_type=user`。
  3. 因此：已登记 `subject_id`（例如凭证列表的 `redeemedBy`）可以过门禁，开出一本 **没有 `admin.users` 行** 的 user 账本。核销入金走 `owner_type=subject`；by-owner 调账走 `owner_type=user`。同一 `owner_id` 两本 CNY 账，对账与运营会看错书。
- 这不是 C 端可自助利用的洞（需要 `wallet.adjust`），但是本 VP 引入的**活路径回归**：A-002 F-005 只是 recommended「别把 user 回调复用到 subject」；后续把主体 OR 进 user 门禁，方向反了。
- 证据：
  - `apps/api/internal/composition/composition.go` 约 567–578 行：`UserByID` 失败则 `SubjectExists`。
  - `apps/api/internal/handler/wallet.go` 125–150、155–185 行：`ownerExists` 通过后只调 `GetOrCreateUserAccount`。
  - `store/repository.go` `GetOrCreateUserAccount`：硬编码 `OwnerUser`。
  - 无「by-owner + subject_id 必须 404 / 不得建 user 户」测试。
- 闭合方向：by-owner 的存在性门禁**只**查 `admin.users`。主体存在性只留在 `CreateAccount(subject)` / `Redeem`。补测试：对已登记 `subject_id` 调 `POST /api/wallet/by-owner/{id}` → 404 `USER_NOT_FOUND`，且 `wallet_accounts` 无 `owner_type=user` 行。

#### F-002 · 导出是通用渲染器对 `items[].code` 的启发式，不是协议声明的 download
- 严重度：**med**
- 建议：**recommended**
- 关联：VP-029 判据 #2 / #5
- 状态：**open**（不阻断判据 #5 的「本 GUI 可导出」）
- 描述：`wallet-vouchers.json` 的 `generateBatch.onSuccess` 仍是 `reload`。明文落盘依赖 `submitForm`：任何表单成功体出现 `items[].code` 就触发 `vouchers_*.csv`。当前仓库里 201 返回该形状的只有凭证批次，故本页能用；协议文档并未声明 download outcome，换客户端或其它表单若返回 `items[].code` 会误下载。Vitest 只锁文件名，不断言 CSV 正文含明文。
- 证据：`render.tsx` 1200–1228 行；`render.test.tsx` 1372 行；`wallet-vouchers.json` `onSuccess.behavior = reload`。

#### F-003 · PostgreSQL 上没有 Redeem / 并发双花端到端
- 严重度：**low**
- 建议：**recommended**
- 关联：VP-029 判据 #3
- 状态：**open**（不把 A-002 F-003 重新打开为 required：PG-unsafe 同事务重读已从代码移除，且 ON CONFLICT 本会话 PG 实测 PASS）
- 描述：`TestPostgresWalletVoucherAndSubject0064` 覆盖 subjects 插入、`owner_type=subject` CHECK、同事务 `ON CONFLICT DO NOTHING`。它**不是** `Redeem`。并发双花只在文件 SQLite 上证明。CAS SQL 双方言通用，剩余是覆盖缺口。
- 闭合方向（可选）：PG 上跑一次 Redeem + 并发首次开户（两张不同卡、同一新主体）。

#### F-004 · `batch_id` 无唯一约束
- 严重度：**low**
- 建议：**recommended**
- 状态：**open**
- 描述：重复生成同一 `batchId` 会混在同一列表过滤；明文只存在于各次 201/CSV。运营可能覆盖或混淆批次。非资金不变式。

#### F-005 · 过期字段是 Unix 秒 `inputNumber`
- 严重度：**low**
- 建议：**recommended**
- 状态：**open**
- 描述：A-004 写「表单仍无 `expiresAt`」已过时——当前 JSON 有该字段。管理员必须填 Unix 秒，易填成毫秒或留空。模块层过期拒绝仍有测试。运营 UX，非资金门禁。

### 必改项汇总

1. **F-001（required / med）**：拆开 user 自动开户与主体存在性。by-owner 不得因 `subjects` 行而为该 ID 创建 `owner_type=user` 账本；补回归测试。

### 与既有意见的异同

| 来源 | 异同 |
|------|------|
| A-001 | 条文过短，把子目标闭合当 Root 完成。本条不继承其 pass。 |
| A-002 F-001/F-002/F-003 | 本会话按**当前代码**复验：CSV 同手势、异币种 fail-closed、PG `ON CONFLICT` 仍成立，**不重开**这三条 required。 |
| A-002 F-004 i18n | 当前 `zh-CN`/`en-US` 已有对应键，视为已修，不再列为 finding。 |
| A-002 F-005 | 当时 recommended「回调仍只查 users」。后续把 `SubjectExists` OR 进 **user** 门禁，制造了本条 F-001。不是原 F-005 的合法闭合。 |
| A-004 | finding-closure，只核 A-002 三条 required。本条是 Root 完成情况的**新一次**设计/代码审查，不把 A-004 pass 当作判据 #7。 |
| A-004 recommended 过期字段 | 字段已出现；本条 F-005 降为 UX。 |

方案层（独立 `subjects` 表、高熵 SHA-256、CAS+`adjust` 同事务、不暴露 HTTP Redeem、`wallet.voucher.issue` 与 `wallet.adjust` 隔离、持久化不跟钱包 HTTP 走）总体合理，SQLite 资金路径测试通过。不能无条件关门的原因是：**主体存在性被接到错误的 HTTP 开户门禁**，会为 C 端 `subject_id` 铸造平行的 user 账本。

### 结论 + 建议给编排器/用户的下一步

**verdict = conditional。** 凭证核销/导出/币种/PG 开户这几块以代码为准已经站住。开放 required = F-001（1 条）。判据 #7 本条不满足。

建议 `/govern`：响应本目标 A-005；先修 F-001（by-owner 门禁回到只查 `admin.users`，主体校验保持在 CreateAccount/Redeem），用 HTTP 测试证明 subject_id 不能开 user 户。不要用更新治理勾选代替该测试。F-002～F-005 可一并排期或 residual。

### 声明

本意见不修改 status/progress；响应由 /govern 处理。
