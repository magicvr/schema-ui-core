---
id: GOAL-005-my-wallet-voucher-redeem
doc: execution-entry
record_id: E-003
status: recorded
parent: GOAL-005-my-wallet-voucher-redeem
created: 2026-09-02
updated: 2026-09-02
version: 0.1.0
---

# E-003 · A-001 recommended 测试补强与执行索引修正（2026-09-02）

## 用户指令

`/govern 响应 workspace-029 GOAL-005 A-001`

## 已发生事实

1. A-001 independent **pass**（0 required）已落盘；本回合按 D-003 消化 F-001～F-004。
2. **F-001**：`wallet_self_test.go` 新增 `TestWalletSelfRedeemIgnoresClientOwnerAndIsolatesUsers`——alice `POST /api/wallet/me/redeem` 带 body/query `ownerId=user-bob` + `accountId` 伪造字段 → 200；alice `owner_type=user` 入账 1800；bob 无 user 户；双方无 subject 户；bob `GET /me` 404。
3. **F-002**：
   - `TestRedeemForUserDoesNotShareUserPath`：先 `RedeemForUser` 再 `Redeem(subject)` → `ErrVoucherAlreadyRedeemed`，subject 户不建，user 余额不变。
   - `TestConcurrentRedeemForUserFailClosed`：文件库 20 并发同码恰好 1 成功、余额 = 1×面额、凭证流水 1 条、无 subject 户。
   - `TestPostgresVoucherRedeemAndConcurrentUser`：PG 两张不同卡（500+1500）并发入同一新 user 户 → 单行、余额 2000、流水 2 条；重复核销不双记。
4. **F-003**：`TestWalletSelfRedeemCreditsUserLedger` 与 `TestRedeemForUserCreditsUserLedgerNotSubject` 在 409 / `ErrVoucherAlreadyRedeemed` 后再读三余额与凭证流水条数 = 1。
5. **F-004**：`02-execution.md` 事实边界改为与 E-002/E-003 一致（HTTP + 页面已实施）。
6. **未做**：F-005 self 未落盘；GOAL-005 未标 `done`。

## 验证（exit 0）

```text
go test ./modules/wallet/voucher ./internal/handler -count=1 -timeout 180s
  -run "TestRedeemForUser|TestWalletSelfRedeem|TestConcurrentRedeemForUser|TestConcurrentDoubleSpend|TestWalletSelfIdentityOnly|TestWalletSelfEntriesOwnScope"
  voucher 1.886s ok；handler 2.828s ok

go test ./internal/store -count=1 -timeout 180s -run "TestPostgresVoucherRedeemAndConcurrent"
  TestPostgresVoucherRedeemAndConcurrentSubject PASS 1.96s
  TestPostgresVoucherRedeemAndConcurrentUser PASS 2.11s

go test ./modules/wallet/... ./internal/handler ./internal/store -count=1 -timeout 180s
  wallet / store / subject / voucher ok；handler 41.7s ok；store 45.9s ok
```

未改生产资金路径（`RedeemForUser` / `redeemInto` / HTTP handler 实现原文未动）。未改 Profile / Manifest。未跑全包 `go test ./...`、vitest。

## 产物

- `apps/api/internal/handler/wallet_self_test.go`
- `apps/api/modules/wallet/voucher/voucher_test.go`
- `apps/api/internal/store/postgres_voucher_redeem_test.go`
- 本目标 D-003 / A-002 / `02-execution.md` 索引

## Git checkpoint

- 测试提交：`a2b003b7`（owned paths = 上述三个 `*_test.go`；验证 = voucher/handler 过滤用例 + PG `TestPostgresVoucherRedeemAndConcurrentUser` PASS 2.11s + `go test ./modules/wallet/... ./internal/handler ./internal/store` exit 0）
- docs 提交：本记录所在 docs 笔（D-003 / E-003 / A-001 / A-002 + 索引 / goal-tree / workspace / Root 指针）
---
