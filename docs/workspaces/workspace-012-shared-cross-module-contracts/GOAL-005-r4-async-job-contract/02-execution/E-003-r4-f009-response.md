---
id: E-003-r4-f009-response
goal: GOAL-005-r4-async-job-contract
status: recorded
created: 2026-08-18
updated: 2026-08-18
parent: GOAL-005-r4-async-job-contract
version: 0.1.0
---

# E-003 · R4 A-004 F-009 响应

A-004 确认 A-002 F-001～F-008 fixed，并发现取消后 lease 过期的可达死状态。D-002 v0.2.0 新增 `recover-cancel`；同时冻结 `CompleteWithCommit`，把 wallet run 插入与 Job succeeded 放进同一事务，从而保证取消先赢时旧 owner 不产生业务 run、成功先赢时二者原子提交。A-005 将 F-009 标为候选 fixed，等待 independent 复核。
