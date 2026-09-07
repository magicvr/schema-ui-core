---
id: GOAL-005-my-wallet-voucher-redeem
doc: audit-entry
record_id: A-001
status: recorded
parent: GOAL-005-my-wallet-voucher-redeem
created: 2026-09-02
updated: 2026-09-02
version: 0.1.0
---

# A-001 · S4 资金路径独立交叉审计（身份隔离 / 不双记 / user·subject 账不串）（2026-09-02）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型**：close-out / execution-facts
- **scope**：`[workspace-029-wallet-prepaid-instrument/GOAL-005-my-wallet-voucher-redeem]` S4 资金路径——身份隔离、重复核销不双记、Admin 自助核销不得记入 `owner_type=subject` 账（VP-029 判据 #8；D-003 / D-002；延续 Root A-005 F-001 隔离）。**不**改 `status` / `progress` / goal-tree。
- **verdict**：**pass**
- **工作区**：`workspace-029-wallet-prepaid-instrument` · Root `GOAL-001-wallet-prepaid-instrument` · `canonical_scope` 匹配 · `shared_materials_catalog: none`（本 scope 未把资料引用当事实）
- **完整意见**：本条

### 范围与区间

用户指定 `/audit` 焦点 = 资金路径 independent（身份隔离、不双记、user/subject 账不串）。

核对：

- 工作区 `workspace.md` / `goal-tree.md`；GOAL-005 五件套与 D-001 / D-002 / E-002；Root D-003；VP-029 判据 #8。
- 代码：`apps/api/modules/wallet/voucher/service.go`（`Redeem` / `RedeemForUser` / `redeemInto`）；`apps/api/modules/wallet/store/repository.go`（`GetOrCreateUserAccountInTx` / `GetOrCreateSubjectAccountInTx` / `MutateInTx`）；`apps/api/internal/handler/wallet_self.go`；`apps/api/modules/wallet/provider.go`；`apps/api/internal/composition/composition.go`（`walletOwnerExists` 只查 `UserByID`）；`apps/api/modules/wallet/schema/my-wallet.json`。
- 测试：`voucher_test.go` · `wallet_self_test.go` · `service_credentials_test.go`（服务凭证拒 `/me/redeem`）· 既有文件库 `TestConcurrentDoubleSpendFailClosed`（subject 路径 / 共享 CAS）。
- P-005：I-029-007 / 008 / 009 均 closed（最晚 S1）；本 scope 无到期未闭合 required 信息项。
- 本会话复跑（exit 0）：

```text
go test ./modules/wallet/voucher ./modules/wallet ./internal/handler -count=1 -timeout 120s
  -run "TestRedeemForUser|TestWalletSelfRedeem|TestConcurrentDoubleSpend|TestWalletSelfIdentityOnly|TestWalletSelfEntriesOwnScope"
```

`modules/wallet/voucher` 1.750s ok；`internal/handler` 2.866s ok。未在本会话复跑全包 `go test ./...`、vitest、或 `PG_TEST_*` 门控用例。

### 成果（有证据）

1. **入账路径与 subject 核销分叉（不调用 `Redeem(subjectID)`）**  
   `RedeemForUser` 走 `GetOrCreateUserAccountInTx`（硬编码 `OwnerUser`）+ 共享 `redeemInto`。`Redeem` 仍走 `SubjectExists` + `GetOrCreateSubjectAccountInTx`。生产 HTTP 唯一调用点：`wallet_self.go` `service.RedeemForUser(..., user.ID, user.Name, code, now)`。`WalletService` 接口**没有** HTTP `Redeem(subjectID)`。

2. **身份隔离（会话推导 owner；禁止匿名；拒绝服务凭证）**  
   `POST /api/wallet/me/redeem` 仅 `a.Middleware`（无权限键）+ `selfIdentity` → `UserIdentityFrom`（人类用户；`IsServiceCredential()` 为 false）。body 结构体只有 `code`；**不**读 query / 路径 ownerId。匿名 401（`TestWalletSelfRedeemAnonymousAndEmpty`）。editor（无 `wallet.*`）可核销（identity-only / I-029-009）。服务凭证打该路径 401 `UNAUTHENTICATED`（`service_credentials_test.go` 与 `/api/wallet/me` 同列）。限流桶 key = `user.ID`，失败 `Record`、成功 `Clear`（D-002 / I-029-008）。

3. **不双记（单事务 CAS + 账本幂等键）**  
   `redeemInto`：`UPDATE vouchers SET status=redeemed WHERE id=? AND status=unused`，`RowsAffected==0` → `ErrVoucherConflict`；已 redeemed → `ErrVoucherAlreadyRedeemed`。同事务 `MutateInTx(adjust, ref_type=voucher, idempotency_key=voucher.id)`。HTTP 重复核销 409（`TestWalletSelfRedeemCreditsUserLedger`）。文件库 20 并发 subject `Redeem` 恰好 1 成功、余额 = 1×面额、流水 1 条（`TestConcurrentDoubleSpendFailClosed`；CAS 在共享 `redeemInto`）。

4. **user / subject 账不串**  
   `TestRedeemForUserCreditsUserLedgerNotSubject`：入账 `owner_type=user`、三余额 = 面额、同 `owner_id` 的 subject 户 `ErrNotFound`、重复 `ErrVoucherAlreadyRedeemed`。  
   `TestRedeemForUserDoesNotShareSubjectPath`：先 `Redeem(subject)` 再 `RedeemForUser` → already redeemed，且**不**创建 user 户。  
   HTTP：`TestWalletSelfRedeemCreditsUserLedger` 断言无 subject 户。  
   延续 A-005 F-001：`composition.go` `walletOwnerExists` 仍只查 `UserByID`（本波未把主体 OR 回 user 自动开户门禁）。`UNIQUE(owner_type, owner_id, currency)` 使两本账即使 `owner_id` 字符串碰巧相同也分列；本路径 opener 只写 `user`。

5. **明文不进审计原文 / 响应**  
   审计 `wallet.adjust` + action `voucher-redeem`，字段 voucherId / batchId / codePrefix / amount / accountId / entryId。HTTP 测试断言响应与 operation detail **不含**明文码。

6. **页面入账目标不可被客户端改写（相邻判据 #9，资金相关）**  
   `my-wallet.json`：`redeemVoucher` → `POST /api/wallet/me/redeem`，`bodyMapping` 仅 `code`，`onSuccess.reload`。toolbar `openRedeem` 单字段表单。

### 对照成功标准（本 scope）

| 标准 | 判定 | 证据 |
|------|------|------|
| 身份隔离：会话推导入账；禁止匿名；禁止代他人 / client ownerId | **pass（实现 + 部分测试）** | handler 只传 `user.ID`；匿名 401；服务凭证 401。HTTP 缺双用户 redeem / body `ownerId` 对偶测试（F-001 recommended） |
| 重复核销不双记；原子 CAS+adjust | **pass（SQLite 顺序 + 共享 CAS 并发）** | 服务/HTTP 重复拒绝。并发双花测试在 **subject** `Redeem` + 文件库。`RedeemForUser` 自身无文件库并发 / PG e2e（F-002 recommended） |
| 不得记入 subject 账 | **pass** | opener 硬编码 user；双向中的 subject→user 有测试；HTTP 断言无 subject 户 |
| I-029-007/008/009 到期 required | **closed** | D-002 / Root D-003；不构成本独立意见开放信息门禁 |
| S4 字面 self + independent | **本条只覆盖 independent** | 无 self 落盘（F-005 recommended）；P-003：`independent` 不强制 self |

### Findings

#### F-001 · HTTP 身份隔离缺少双用户 / ownerId 注入对偶测试

- 严重度：**med**
- 建议：**recommended**
- 关联：判据 #8；D-002「禁止 client 传 ownerId」
- 状态：**open**
- 描述：`GET /api/wallet/me` 已有「`?ownerId=` 忽略、alice/bob 流水互不可见」（`TestWalletSelfEntriesOwnScope`）。**`POST /api/wallet/me/redeem` 没有对偶**：未测 (1) alice 核销后 bob 余额仍为 0；(2) body 多带 `ownerId`/`accountId` 仍入会话用户账。实现上 `encoding/json` 只解码 `code`，owner 仅来自 `UserIdentityFrom`，代码审查成立，故不升 required。
- 闭合方向：补 HTTP 测试——两用户、多余 JSON 字段、断言只动会话 `owner_type=user` 行。

#### F-002 · `RedeemForUser` 交叉路径 / 并发 / PG 覆盖窄于 subject 核销

- 严重度：**med**
- 建议：**recommended**
- 关联：判据 #8；I-029-006 并发 fail-closed 先例（A-002 F-001 已在 **subject** 路径用文件库闭合）
- 状态：**open**
- 描述：
  1. 有 subject→user 拒绝（`TestRedeemForUserDoesNotShareSubjectPath`），**无** user→subject 反向（同一 CAS，代码对称，仍缺回归锁）。
  2. 20 并发双花只打 `Redeem(subjectID)`。`RedeemForUser` 共享 `redeemInto` CAS，但首次开户走新方法 `GetOrCreateUserAccountInTx`（ON CONFLICT 抄 subject 已修模式，**本波未在文件库并发或 PG 上点亮**）。既有 `TestPostgresVoucherRedeemAndConcurrentSubject` 仍是 subject。
  3. 账本 `UNIQUE(account_id, idempotency_key)` 是**每账户**幂等；跨账户不双记的主防线是凭证行 CAS。CAS 与入金同事务则成立；覆盖缺口不是实现回退。
- 闭合方向（可选）：文件库并发 `RedeemForUser` 同码恰好 1 成功；user 后再 `Redeem(subject)` 失败且 subject 余额不增；PG 两张不同卡并发入同一新 user 户（对偶 subject PG 测试）。

#### F-003 · 重复核销测试未再读余额

- 严重度：**low**
- 建议：**recommended**
- 状态：**open**
- 描述：HTTP 409 与服务 `ErrVoucherAlreadyRedeemed` 已测；第二次成功路径在 `status==redeemed` 时 return，未到 `MutateInTx`。未断言重复后 `balance_total` / voucher 流水条数仍为 1。建议补一行断言，避免以后改控制流时静默双记。

#### F-004 · 执行索引「事实边界」与 E-002 矛盾

- 严重度：**low**
- 建议：**recommended**
- 状态：**open**
- 描述：`02-execution.md` 仍写「尚未实施 HTTP 或页面改动」；E-002 已记录 HTTP + 页面 + 回归。非资金缺陷；关门前应改索引，避免只读索引的人以为未实施。

#### F-005 · S2/S3 self 未落盘；S4 成功标准字面含 self

- 严重度：**low**
- 建议：**recommended**
- 状态：**open**
- 描述：`00-meta` 成功标准 4 写「self + **independent**（资金路径）」；审计策略写 S1/S2/S3 self、S4 independent。`03-audit` 在本条之前无 A 条目。P-003：模式 `independent` **不**固定要求 self（`cross` 才必须）。本意见履行独立资金路径，**不**把缺 self 当作资金 fail。编排器若按成功标准 4 字面关门，需补 self 或用户书面降级该字面。

### 必改项汇总

无。开放 required = 0。

### 与既有意见的异同

- **延续 Root A-005 F-001（A-007 已 `fixed`）**：by-owner 不得因 `subjects` 行开 `owner_type=user` 孤儿账。本波 `RedeemForUser` **没有**把主体 ID 接进自助核销；HTTP owner = 会话用户。抽查 `walletOwnerExists` 仍仅 `UserByID`。
- **不同于 GOAL-002 A-002**：当时并发测试在 `:memory:` 上无效（required）。当前 subject 并发已改文件库；本条不把同一 CAS 再升 required，只要求 `RedeemForUser` / user InTx 的覆盖补齐（recommended）。
- **无既有 GOAL-005 self/independent 可冲突**；本条为该目标 A 序列首条。

### 结论 + 建议给编排器/用户的下一步

资金路径三项（身份隔离、不双记、user/subject 不串）在实现与本会话复跑的 SQLite/HTTP 测试上**成立**。可进入 `/govern` 响应本条。

建议：

1. **可关门（就本 independent 资金 scope）**：0 required；I-029-007/008/009 已 closed。
2. **推荐但不阻断**：消化 F-001/F-002 测试补强，或用户书面 `accepted-residual` / `user-overruled` 并写复审触发。
3. **过程**：确认 S4 是否仍要 self；要则先补 S2/S3 或关门 self，再标 GOAL-005 `done`。
4. 同步改 `02-execution.md` 过时事实边界（F-004）。

建议 `/govern` 下一句：

> `$govern` 响应 GOAL-005 A-001（independent pass，0 required）：按推荐补 HTTP 双用户/ownerId 与 RedeemForUser 并发覆盖，或书面 residual/overruled；并确认是否补 self 后关 S4。

### 声明

本意见不修改 status/progress；响应由 /govern 处理。
