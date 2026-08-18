---
id: A-008-r4-s1-self
goal: GOAL-005-r4-async-job-contract
source: self
verdict: pass
status: recorded
created: 2026-08-18
updated: 2026-08-18
parent: GOAL-005-r4-async-job-contract
version: 0.1.0
---

# A-008 · R4 S1 实施自审

| 检查 | 结论 |
|------|------|
| migration owner/Profile 不变式 | pass：仅 `compiled.PersistenceProviders()` 增加 `core.jobs`；Profile/BuiltinModules 未改 |
| 六态与 DB 不变式 | pass：DDL CHECK 覆盖 status、lease、result/error/finished/expires 组合 |
| transition table | pass：repository 条件 UPDATE 覆盖 claim/reclaim/progress/cancel/fail/retry/expire/exhaust |
| fencing/原子 commit | pass：旧 token 不执行 callback；callback 与 Job success 同事务；失败回滚 |
| 定向与 race 测试 | pass：E-005 命令全绿 |

开放 required/recommended = 0。S1 关门，放行 S2。A-001 F-001 的 startup/周期恢复验证由 S2 承接；A-001 F-002 的 descriptor conformance 由 S3/S4 承接。
