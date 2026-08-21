---
id: E-003-r3-implementation
doc: execution-entry
status: recorded
created: 2026-08-18
updated: 2026-08-18
parent: GOAL-004-r3-optimistic-concurrency-idempotency
version: 0.1.0
---

# E-003 · R3 S1/S2 实现

## 共享版本契约

- 新增 `internal/concurrency`：强 ETag `"vN"`、三来源 expected version 解析、missing/invalid sentinels 与拒绝集测试。
- wallet status PATCH 使用 pointer body 字段区分缺失与显式 0；428 `PRECONDITION_REQUIRED`、400 `INVALID_PRECONDITION`、409 `LEDGER_VERSION_CONFLICT` 已进入冻结错误码与双语 catalog。
- D-002 指定的单账户响应设置 ETag；列表保持不设置账户强 ETag。

## 幂等 operation

- wallet service 返回 replay 标志；HTTP mutation 暴露 `operationId/state/replayed/idempotencyKey/resourceVersion`。
- payload identity 增加 `actor_id`，不把可变展示名纳入指纹。
- unique insert 竞争先回滚，再回读获胜 operation；相同指纹 replay，其他情况 409。
- replay 不再追加第二条成功 wallet 业务审计；无 key 仍保持兼容并返回本次 succeeded operation。

## 验证

- 定向 concurrency/wallet/error contract 测试通过。
- `go test ./... -count=1` 在 `apps/api` 全量通过；handler 218.507s，operationlog 11.136s，wallet/store 1.911s，docscheck 0.773s。
- 实现 checkpoint：`08dcec8`。
