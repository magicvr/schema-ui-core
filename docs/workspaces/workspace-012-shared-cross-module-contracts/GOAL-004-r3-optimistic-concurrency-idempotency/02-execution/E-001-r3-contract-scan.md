---
id: E-001-r3-contract-scan
doc: execution-entry
status: recorded
created: 2026-08-18
updated: 2026-08-18
parent: GOAL-004-r3-optimistic-concurrency-idempotency
version: 0.1.0
---

# E-001 · R3 version/idempotency 现状扫描

## 已核对事实

- wallet account 表已有 `version INTEGER NOT NULL DEFAULT 0`；`UpdateStatus` 和 balance mutation 均以 `WHERE id=? AND version=?` 更新并递增，0 rows 返回 `ErrVersionConflict`。
- wallet ledger 以 `(account_id,idempotency_key)` 唯一；同 key 同 payload 读取既有 entry，异载荷返回 `ErrIdempotencyConflict`；mutation 在单个 `WithTx` 内完成。
- HTTP status PATCH 只接受 JSON `version`，无 `If-Match`/`expectedVersion`，且缺字段会被 Go 零值当作 version 0。
- Manifest 的 ETag 仅用于 `If-None-Match` 304，不能作为写前置条件实现直接复用。
- wallet replay HTTP 返回 account/entry，但没有显式 `operationId/state/replayed`；handler 会对 replay 再写一条 wallet audit event。
- settings、scheduledtasks、dictionary 的 `updated_at` 更新无旧值条件，不作为 R3 首个消费切片。

## 风险

- ETag 与 body version 多来源不一致必须 fail closed，不能静默任选。
- 现有 wallet Web/测试仍传 `version`，需保留 wire compatibility。
- operation replay 必须复用 ledger identity，不能生成第二个成功 operation 或重复业务审计。
- SQLite 唯一约束竞争路径需保持事务回滚；本轮不声称提供跨数据库锁语义。
