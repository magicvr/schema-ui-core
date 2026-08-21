---
id: E-004-r3-audit-recommendation-fixes
doc: execution-entry
status: recorded
created: 2026-08-18
updated: 2026-08-18
parent: GOAL-004-r3-optimistic-concurrency-idempotency
version: 0.1.0
---

# E-004 · R3 independent 建议修复

- 新增 package-internal `TestReplayAfterIdempotencyRace`，直接验证获胜 operation 回读、当前 account/version 与异 actor 冲突。
- ETag parser 增加多余引号拒绝用例。
- wallet HTTP 测试增加 GET by-owner ETag、list 无账户 ETag、mutation 忽略 If-Match、无 key succeeded operation 包络。
- Web `en-US`/`zh-CN` 增加 `error.preconditionRequired` 与 `error.invalidPrecondition`。
- 定向 Go 测试通过；Web i18n catalog/structural 15 tests 通过。
