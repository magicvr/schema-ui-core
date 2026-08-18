---
id: D-001-r4-job-contract
goal: GOAL-005-r4-async-job-contract
status: superseded
created: 2026-08-18
updated: 2026-08-18
parent: GOAL-005-r4-async-job-contract
version: 0.1.0
---

# D-001 · R4 Job 状态机、持久化与 wallet 消费边界

## 决定

1. 通用 Job 数据由 migration-only `core.jobs` owner 创建；runtime 放在 API internal package，不注册为默认 Profile 模块。
2. 首个业务 handler 为 `wallet.reconcile`；HTTP 路由由既有 `admin.wallet` provider 贡献，权限沿用 `wallet.read`，Job 查询按提交 actor 隔离。
3. Job 状态为 `queued / running / succeeded / failed / cancelled / expired`。终态为 `succeeded / failed / cancelled / expired`；`expired` 只由已结束且到达保留期限的 Job 转入。
4. `queued -> running` 通过原子 claim 完成并增加 attempt、写 lease；租约过期的 `running` 可被重新 claim。进度只允许在 `running` 中单调增加。
5. 取消：`queued` 直接转 `cancelled`；`running` 先持久化 `cancel_requested`，runner 协作取消并最终转 `cancelled`。完成/失败与取消竞态以数据库条件更新唯一决胜。
6. 重试：仅 `failed` 且 `attempt < max_attempts` 可回到 `queued`；保留同一 Job id、清除运行期错误/结果/lease，不回退 attempt。
7. 结果为 JSON；成功结果可经 result endpoint 下载。到期后转 `expired` 并清空结果，result endpoint 返回稳定 410 错误。
8. `POST /api/wallet/reconcile` 改为 202 + Job representation；现有 reconcile run 仍是 wallet 业务结果与历史真相源。

## 关键字段

`id, kind, status, payload, progress, cancel_requested, attempt, max_attempts, lease_owner, lease_expires_at, result, error_code, error_message, actor_id, created_at, updated_at, finished_at, expires_at`。

## 未选方案

- 复用 `task_runs`：其模型只表达 scheduled task 的最终 `ran/failed`，缺失 payload、租约、进度、取消、重试与结果保留。
- 把 wallet reconcile run 直接扩成通用 Job：会把业务结果与跨模块运行契约耦合。
- 新增默认 `core.jobs` runtime module：会改变 Profile/Manifest 装配面，超出 R4 边界。
- 引入外部队列：当前单进程切片无此依赖需求，且会扩大部署与恢复协议。

## 门禁

本决定须经 S0 independent 设计审计；I-002/I-003 关闭后状态才改为 `accepted`，随后才能进入 S1。

## 后续

A-002 independent 指出六条 required 精确化缺口；本初稿由 D-002 替代，历史判断保留。
