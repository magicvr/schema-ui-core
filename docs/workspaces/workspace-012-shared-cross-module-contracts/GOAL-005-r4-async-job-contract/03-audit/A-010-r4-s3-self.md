---
id: A-010-r4-s3-self
goal: GOAL-005-r4-async-job-contract
source: self
verdict: pass
status: recorded
created: 2026-08-18
updated: 2026-08-18
parent: GOAL-005-r4-async-job-contract
version: 0.1.0
---

# A-010 · R4 S3 wallet consumer 自审

| 检查 | 结论 |
|------|------|
| 202 + poll + result | pass；成功结果为稳定 reconcile JSON attachment |
| actor predicate | pass；permission 后固定 `kind=wallet.reconcile` 且 actor match，跨 actor / 不存在统一 `JOB_NOT_FOUND` |
| consumer 原子性 | pass；callback 错误不留下 reconcile run，Job failed；成功时 run 与 Job 同事务提交 |
| cancel / retry / expiry | pass；非法转换映射稳定 catalog code，成功读取前执行 `ExpireIfDue` |
| 审计事件 | pass；queued、success、failed、cancelled 均落盘；scanner 恢复终态也触发 hook |
| migration 43 | pass；保留既有 operation row 与 correlation，新增三类事件可写 |
| lifecycle / Profile 边界 | pass；wallet-enabled 才启动 runner；模块集、page/nav/fragment/permission 不变 |
| race / repeat | pass；Job + wallet consumer race 与 count=10 均通过 |

开放 required/recommended = 0。A-006 要求的 S3 必测项已满足：wallet callback 调用 `ReconcileOnceTx`，没有调用自开事务的旧 `ReconcileRun`；consumer rollback 不留下孤立 run。S3 关门，放行 S4。
