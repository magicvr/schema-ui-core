---
id: GOAL-005-my-wallet-voucher-redeem
doc: execution-entry
record_id: E-002
status: recorded
parent: GOAL-005-my-wallet-voucher-redeem
created: 2026-09-02
updated: 2026-09-02
version: 0.1.0
---

# E-002 · S1 冻结与 S2/S3 实施回归

## 2026-09-02 · 实施

### 已发生事实

1. D-002 冻结 I-029-007/008：`POST /api/wallet/me/redeem`；`RedeemForUser`；内存限流 15min/10/user id；审计复用 `wallet.adjust` + action `voucher-redeem`（不扩 operation_log.event CHECK、不加迁移）。
2. 入账走 `GetOrCreateUserAccountInTx` + CAS + `adjust`/`ref_type=voucher`。禁止 `Redeem(subjectID)`。
3. 「我的钱包」toolbar「使用预付凭证充值」→ 模态表单；成功 reload。
4. 回归：`go test ./... -count=1`（apps/api）exit 0；`vitest` schema-keys.structural 4/4。

### 证据

| 主张 | 路径 |
|------|------|
| 用户核销 | `apps/api/modules/wallet/voucher/service.go` `RedeemForUser` |
| InTx 开户 | `apps/api/modules/wallet/store/repository.go` `GetOrCreateUserAccountInTx` |
| HTTP | `apps/api/internal/handler/wallet_self.go` `POST /api/wallet/me/redeem` |
| 页面 | `apps/api/modules/wallet/schema/my-wallet.json` |
| 测试 | `wallet_self_test.go` / `voucher_test.go` |

### 未做

S4 independent 关门审计。未改 Profile 默认集 / Manifest 装配语义（additive 路由与页面动作）。
