---
id: E-009-r4-close
goal: GOAL-005-r4-async-job-contract
status: recorded
created: 2026-08-18
updated: 2026-08-18
parent: GOAL-005-r4-async-job-contract
version: 0.1.0
---

# E-009 · R4 A-012 响应与目标关门

A-012 independent close-out verdict = pass，开放 required = 0；唯一 recommended F-011 指出 HTTP fixture 与生产 `JobService.Get` 缺少直接的 result expiry → 410 用例。

`7de9a0b` 增加两层验证：生产 `JobService.Get` 在 24h 后原子转 `expired` 并清空 result；HTTP fixture 在到期后 GET result 返回 410 `JOB_RESULT_EXPIRED`。定向 handler/wallet 测试、wallet race 与 docscheck 均通过。

A-013 将 F-011 以 `fixed` 闭合。S4 完成；派生进度由 4/5 更新为 5/5。GOAL-005 状态改为 `done`，goal-tree 树与表同步。
