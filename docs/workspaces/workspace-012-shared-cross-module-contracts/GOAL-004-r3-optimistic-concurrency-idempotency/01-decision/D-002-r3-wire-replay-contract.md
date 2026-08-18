---
id: D-002-r3-wire-replay-contract
doc: decision-entry
status: accepted
created: 2026-08-18
updated: 2026-08-18
parent: GOAL-004-r3-optimistic-concurrency-idempotency
version: 0.1.0
---

# D-002 · R3 wire 与 replay 契约冻结

## 版本前置条件

1. 强 ETag 的唯一合法形式为 `"vN"`，`N` 是非负十进制 `int64`；允许整个 header 前后空白。拒绝弱标签、`*`、列表、多余引号、未加引号、负数、溢出和标签内部空白。
2. status PATCH 的版本来源为 `If-Match` header、JSON number `expectedVersion`、legacy JSON number `version`。字段用指针解码以区分缺失与显式 `0`；字符串或其他 JSON 类型由 body 校验以 400 拒绝。
3. 三来源可任意提供一个或多个；多个来源必须完全一致。三者皆缺失返回 428 `PRECONDITION_REQUIRED`；header 非法或来源矛盾返回 400 `INVALID_PRECONDITION`；CAS stale 继续返回 409 `LEDGER_VERSION_CONFLICT`。
4. 账户 ETag 出现在成功的单账户响应：GET/POST by-owner、POST account、PATCH status、by-owner adjust 与 account-id mutation。列表不设置账户 ETag；当前不新增 GET-by-id。
5. ledger mutation 不消费 `If-Match`/`expectedVersion`。资源属性更新使用版本前置条件；操作型余额写入使用 durable idempotency key 和 repository 内 CAS。

## 幂等 replay

1. payload 指纹固定为 `entry_type`、`amount_delta`、`memo`、`ref_type`、`ref_id`、`actor_id`；`actor_name` 是可变展示字段，不参与 identity。不同 actor 复用同 key 必须 409。
2. 唯一约束竞争先回滚本事务，再按 `(account_id,idempotency_key)` 回读：相同指纹返回既有 entry 并标 replay；不同指纹或无法回读时 fail closed 为 409。
3. HTTP 成功结果新增 `operation`：`operationId` = durable ledger entry id，`state` = `succeeded`，`replayed`，可选 `idempotencyKey`，`resourceVersion` = 本次响应时刻的账户 version。历史 replay 的 `resourceVersion` 不是原写入快照，不可当作旧状态的 ETag。
4. 余额不足、disabled、invalid 或其他失败在 ledger insert 前退出：不铸造 durable operation、不占用 idempotency key，允许修正条件后用同 key 重试。
5. 同 key 同 payload replay 不再写第二条成功 wallet 业务审计；首个成功仍写一次。S2 必须同时断言 ledger 行数和 operationlog 业务事件数均为 1。
6. 无 key 调用仍返回本次 entry 对应的 succeeded operation，`replayed=false`，但不承诺跨请求 replay。
