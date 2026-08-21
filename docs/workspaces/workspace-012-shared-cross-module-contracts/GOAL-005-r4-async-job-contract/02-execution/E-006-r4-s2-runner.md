---
id: E-006-r4-s2-runner
goal: GOAL-005-r4-async-job-contract
status: recorded
created: 2026-08-18
updated: 2026-08-18
parent: GOAL-005-r4-async-job-contract
version: 0.1.0
---

# E-006 · R4 S2 runner/recovery 实现

`c8305bb` 新增 profile-independent Job runner：handler registry、startup/周期 scan、claim/reclaim、heartbeat、同进程 active set、进度 reporter、协作取消、failed/retry 唤醒、result expiry 与可等待 Stop。

测试覆盖：重启后 reclaim、运行中取消、失败后 retry、结果自动过期、两个 runner 下 heartbeat 防重复执行、Stop 后 Job 保持 running 等待 lease recovery，以及没有 cancel request 的 `context.Canceled` 必须 failed。

验证：`go test ./internal/jobs`、`go test -race ./internal/jobs`、`go test -count=10 ./internal/jobs` 均 PASS。
