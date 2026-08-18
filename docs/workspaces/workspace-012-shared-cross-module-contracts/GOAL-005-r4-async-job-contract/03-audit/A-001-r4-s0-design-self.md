---
id: A-001-r4-s0-design-self
goal: GOAL-005-r4-async-job-contract
source: self
verdict: pass
status: recorded
created: 2026-08-18
updated: 2026-08-18
parent: GOAL-005-r4-async-job-contract
version: 0.1.0
---

# A-001 · R4 S0 设计自审

## 审计头

| 项 | 值 |
|----|----|
| source | self |
| scope | GOAL-005 S0：D-001、E-001、I-001～I-004；migration owner、状态机、wallet 消费边界 |
| verdict | pass |
| required findings | 0 |

## 核对结论

1. `kernel.CollectPersistence` 明确从所有 compiled provider 收集 migration，且 runtime `RegisterContributions` 只注册 Plan 中启用的 provider；因此 migration-only `core.jobs` owner 不需要进入默认 Profile/module matrix。
2. `store.Store.WithTx` 提供数据库事务边界，SQLite 连接池固定为单连接；repository 可用带前态/lease 条件的 `UPDATE` + affected rows 完成原子 claim 和竞态决胜。
3. wallet service 的 `Reconcile` 已把 reconcile run 持久化；异步 handler 调用同一 service 后把 run 投影为 Job result，不会用 Job 状态替代业务历史。
4. `task_runs` 与 wallet reconcile run 均缺通用状态、lease、取消、retry/result retention，拒绝复用的理由成立。
5. D-001 的状态集合覆盖 VP-012 方向级范围；result 过期通过先转 `expired`、再清空结果，不会把“结果不可取”误写成执行失败。

## Findings

| ID | 等级 | finding | disposition |
|----|------|---------|-------------|
| F-001 | recommended | runtime 必须在启动时和周期扫描 `queued` / lease-expired `running`，否则“可重新领取”只有 repository 能力，没有运行路径。 | implementation gate：S2 测试证明 startup recovery/claim |
| F-002 | recommended | 新增 wallet Job routes 时必须同步 provider Descriptor 与 kernel profile descriptor，并用 conformance test 证明默认模块集合/Manifest 语义未改变。 | implementation gate：S3/S4 验证 |

## 结论

Self 设计审计通过，开放 required = 0。由于 data/migration/compatibility 风险，I-002/I-003 仍须等待 grok-build independent 设计审计后才可标 `verified`；本意见不放行 S1。
