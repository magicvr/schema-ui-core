---
id: E-007-r4-s3-wallet-jobs
goal: GOAL-005-r4-async-job-contract
status: recorded
created: 2026-08-18
updated: 2026-08-18
parent: GOAL-005-r4-async-job-contract
version: 0.1.0
---

# E-007 · R4 S3 wallet async Job / HTTP / lifecycle 实现

代码 checkpoint `3ce848b` 完成 D-002 的首个真实消费路径：

- wallet reconcile 改为提交 `wallet.reconcile` Job 并返回 202；新增 actor-scoped 查询、取消、重试和 JSON result attachment 路由。
- wallet `ReconcileOnceTx` 与 Job `CompleteWithCommit` 使用同一事务；Job ID 同时作为 reconcile run ID，consumer callback 失败会整体回滚并落为 failed。
- composition 只在 `admin.wallet` 启用时注册并启动 Job runner；停止顺序为 HTTP intake、runner、module runtime、store。
- migration 43 扩展 queued/failed/cancelled 审计事件 CHECK，并在重建中保留 migration 41 correlation 外键数据。
- terminal hook 覆盖即时成功/失败/取消以及 scanner 的 recover-cancel / attempt exhaustion。
- Provider Descriptor 与 `kernel.BuiltinModules()` 同步四条新增 route key；未新增 Profile module/page/navigation/fragment/permission。

已执行并通过：

- `go test -timeout 120s ./internal/jobs ./internal/modules/wallet ./internal/store ./internal/composition ./internal/kernel`
- `go test -timeout 60s ./internal/handler -run 'TestWallet' -count=1`
- `go test -race -timeout 120s ./internal/jobs ./internal/modules/wallet`
- `go test -count=10 -timeout 120s ./internal/jobs ./internal/modules/wallet`
- `git diff --check` 与未跟踪文件 trailing-whitespace 检查

S3 已有确定性实现与测试证据；全 API / docscheck / independent 关门审计留在 S4。
