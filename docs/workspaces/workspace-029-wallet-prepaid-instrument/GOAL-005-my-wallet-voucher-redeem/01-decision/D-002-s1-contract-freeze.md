---
doc_type: goal-decision
id: D-002-s1-contract-freeze
parent: GOAL-005-my-wallet-voucher-redeem
date: 2026-09-02
status: accepted
version: 0.1.0
---

# D-002 · S1 合同冻结（HTTP 路径 + 限流）

## 触发

用户 2026-09-02「继续」：采用 reopen 时给出的 S1 默认候选。I-029-007 / I-029-008 最晚 S1，阻断实施。

## 决定

| 项 | 冻结值 |
|----|--------|
| HTTP | `POST /api/wallet/me/redeem`；JSON `{ "code": "<plaintext>" }`；identity-only；入账目标 = 会话用户 id，**禁止** client 传 ownerId |
| 服务 | 新方法 `RedeemForUser(ctx, userID, actorName, code, now)`。**不**调用 `Redeem(subjectID, code)` |
| 账本 | 同事务 CAS `unused→redeemed` + `GetOrCreateUserAccountInTx` + `MutateInTx(adjust, ref_type=voucher, idempotency_key=voucher.id)` |
| 权限 | identity-only（I-029-009 已冻） |
| 限流（I-029-008） | 已交付内存限流器专用桶：window **15min**、max **10**、key = **user id**、capacity 默认。失败 `Record`，成功 `Clear`。不按 IP 分桶（已登录；NAT 会误伤）。不消耗 RT-Q05 Redis trigger |
| 评估结论 | 码为 24 字符高熵，在线猜码不可行。限流针对已登录面的喷刷/日志噪声，不是密码学补偿。内存桶对单进程 Admin 基线足够 |
| 审计 | 复用既有 `wallet.adjust`（operation_log.event CHECK 未含新事件名，本波不加迁移）；detail.action = `voucher-redeem`，含 voucherId / amount / accountId / codePrefix；**不含明文** |
| 页面 | 「我的钱包」toolbar CTA → 模态表单（单字段 code）；成功 `reload` |

## 未选方案

- 复用 `Redeem(subjectID)`：会记入 subject 账（D-003 / A-005 F-001）。
- 匿名或公开路径。
- 按 IP 分桶：已登录身份足够，NAT 会把多人打进同一桶。
- 新权限键 / `wallet.voucher.issue`。
- Redis 限流（不消耗 trigger）。

## 影响

I-029-007 / I-029-008 → closed。放行 S2 实施。
