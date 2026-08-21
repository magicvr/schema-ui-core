---
id: E-001-r4-contract-scan
goal: GOAL-005-r4-async-job-contract
status: recorded
created: 2026-08-18
updated: 2026-08-18
parent: GOAL-005-r4-async-job-contract
version: 0.1.0
---

# E-001 · R4 async/long-operation 现状扫描

## 已核对事实

- API 当前没有通用 jobs table/repository/runner；迁移最大版本为 41，下一全局版本可用 42。
- scheduled task 的 `task_runs` 只记录同步 `Execute` 后的 `ran/failed`，没有 queued/running、lease、attempt、progress、cancel 或 retry。
- import/export 与 file download 均为同步请求；Web 端只有提交/错误/下载状态，没有通用轮询与 Job UI。
- wallet reconcile 是同步调用，但已有持久化 reconcile run、分页历史与独立 service/repository，适合作为首个真实异步消费者。
- `admin.wallet` 已有 profile contribution 与权限边界；通过其 provider 挂载 Job HTTP 路径可避免新增默认 runtime module。
- compiled persistence 允许独立 migration provider；Job 表可由 migration-only owner 管理，不影响 Profile 默认模块集合。

## 结论

I-001 已验证。D-001 的 migration/runtime 所有权与完整 transition table 需经 S0 independent 审计后关闭 I-002/I-003。
