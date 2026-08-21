---
id: D-001-r3-wallet-contract
doc: decision-entry
status: accepted
created: 2026-08-18
updated: 2026-08-18
parent: GOAL-004-r3-optimistic-concurrency-idempotency
version: 0.1.0
---

# D-001 · R3 wallet 消费切片与兼容契约

1. 选用 `admin.wallet` 作为首个真实消费模块：它已有持久化 `version`、CAS update、账户作用域 `UNIQUE(account_id,idempotency_key)`、单事务 mutation 和 409 错误码，避免用新造 demo 证明契约。
2. 新增 shared version precondition helper；强 ETag 采用 `"v<non-negative-int>"`。status PATCH 可从 `If-Match`、`expectedVersion`、legacy `version` 取得期望版本；多个来源必须一致。缺失返回 428，非法或矛盾返回 400，stale CAS 继续返回 409 `LEDGER_VERSION_CONFLICT`。
3. 保留 legacy JSON `version` 和无 idempotency key mutation；这是 wire compatibility，不代表新调用方可省略并发前置条件。
4. ledger entry id 作为 durable `operationId`；成功状态为 `succeeded`。同 key 同 payload replay 返回原 entry id 并标 `replayed=true`；异载荷维持 409 `LEDGER_IDEMPOTENCY_CONFLICT`。
5. 不把 R4 queued/running/failed/cancelled 状态机塞入 R3；R3 只定义同步 operation 的终态与 replay identity。
6. 审计模式唯一为 `independent`（data/compatibility）；provider 采用项目级 grok-build (grok-4.6 reasoning high)。
