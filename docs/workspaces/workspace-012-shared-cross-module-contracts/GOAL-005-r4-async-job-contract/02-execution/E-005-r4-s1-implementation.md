---
id: E-005-r4-s1-implementation
goal: GOAL-005-r4-async-job-contract
status: recorded
created: 2026-08-18
updated: 2026-08-18
parent: GOAL-005-r4-async-job-contract
version: 0.1.0
---

# E-005 · R4 S1 migration/repository 实现

## 事实

- `2013e7f`：新增 migration-only `core.jobs` version 42，六态/lease/result 关系由数据库 CHECK 约束；只进入 compiled persistence catalog，不进入 runtime Profile。
- `e670b56`：新增 `internal/jobs` model/repository，落实 D-002 的条件 claim/reclaim、lease fencing、单调进度、取消/恢复、attempt exhaustion、retry、`CompleteWithCommit`、结果过期与 actor 查询。
- `CompleteWithCommit` 测试证明：失效 fencing token 不执行 callback；callback 写入后返回错误会整体回滚；新 token 可随后成功提交同一 business row。

## 验证

- `go test ./internal/jobs ./internal/modules/jobs/migration ./internal/kernel ./internal/store`：PASS。
- `go test -race ./internal/jobs`：PASS。
- `git diff --check`：PASS。
