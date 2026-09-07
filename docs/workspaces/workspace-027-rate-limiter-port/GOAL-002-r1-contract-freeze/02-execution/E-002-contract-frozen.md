---
doc_type: goal-execution
id: E-002-contract-frozen
parent: GOAL-002-r1-contract-freeze
date: 2026-09-01
status: active
version: 0.1.0
---

# E-002 · 合同正文冻结 + 端口落地（C2）

## 事实时间线

- 2026-09-01：D-002 v0.1.0 合同冻结（§0 适用与验收基线 / §1 端口形状 / §2 key 语义 / §3 窗口语义 / §4 容量与驱逐 / §5 W12 D-002 常量表 / §6 并发 / §7 停机 / §8 红线 / §9 信息裁决 / §10 验收方式 / §11 未选方案）。
- 2026-09-01：端口本体 `apps/api/kernel/ratelimit.go` 落地——`RateLimiter`（Allow/Record/RetryAfterSeconds/Clear，`now` 注入）+ `RateLimiterProvider`（`NewRateLimiter(window, max, capacity)`，capacity≤0 → `DefaultRateLimiterCapacity` 1<<16）+ 可执行语义权威纯函数 `RateLimiterInWindow`（`t.After(now.Add(-window))`）与 `RateLimiterRetryAfterSeconds`（remain≤0→1；否则 Round 秒）。
- 2026-09-01：合同级快测 `apps/api/kernel/ratelimit_test.go` 落地——编译期端口面断言（stub ×2）+ 常量断言 + `RateLimiterInWindow` 表驱动 8 例（含 cutoff 恰等为窗外 / 零窗口）+ `RateLimiterRetryAfterSeconds` 表驱动 7 例（全窗 900 / 5min 300 / 30s / 恰到期→1 / 超窗→1 / 亚秒 400ms→0 / 600ms→1）。
- 2026-09-01：验证绿——`go vet ./kernel/...` 0 / `go test ./kernel/... -count=1` ok / `go build ./...` 通过；git 变更面 = `apps/api/kernel/ratelimit.go` + `kernel/ratelimit_test.go` + `docs/workspaces/workspace-027-rate-limiter-port/**`（零越界）。

## 产物

- `GOAL-002-r1-contract-freeze/01-decision/D-002-rate-limiter-port-contract.md`（责任文件 v0.1.0）
- `apps/api/kernel/ratelimit.go`（端口本体）
- `apps/api/kernel/ratelimit_test.go`（合同级快测）

## 下一步

- C3 审视：A-001 self + A-002 grok build（grok-4.6 · high）independent → A-003 合并响应 → R1 关门。