---
doc_type: goal-audit
id: A-001-contract-freeze-closeout-self
parent: GOAL-002-r1-contract-freeze
date: 2026-09-01
source: self
scope: GOAL-002 R1 合同冻结全量（D-002 合同 ↔ kernel/ratelimit.go 逐节一致性 / 快测覆盖 / 迁移不回归基线 / 越界核账 / 信息门禁）
verdict: pass
open_required: 0
status: active
version: 0.1.0
---

# A-001 · R1 合同冻结关门自审（self）

## 1. 信息门禁（P-005）

| ID | 级别 | 状态 | 证据 |
|----|------|------|------|
| I-027-001 | required | **verified** | 2026-09-01 用户裁决「语义拆分保持」（D-001 accepted；合同 §1 逐条落实：Allow 不注册 / Record 计数 / RetryAfterSeconds 秒 / Clear 清零 / now 注入 / capacity 默认 1<<16） |
| I-027-003 | non-blocking | **verified** | 用户确认「滑动窗口保持 + 策略独立」（合同 §3：InWindow 谓词 + 无 ExpiryPolicy 式接口） |
| I-027-004 | non-blocking | **verified** | 用户确认「不新增复合 key」（合同 §2：key 不透明、既有形态保持） |
| I-027-002 | required | 待裁决（R2 前置） | 本目标不关闭；最晚阶段 R2——不影响 R1 关门 |

## 2. 合同 ↔ 实现逐节一致性

- **§1 端口形状**：`kernel.RateLimiter`（Allow/Record/RetryAfterSeconds/Clear + now 注入）与 `kernel.RateLimiterProvider`（NewRateLimiter(window,max,capacity)）与合同签名逐字一致；注释明示 Allow 永不注册、RetryAfterSeconds 仅限 not-allowed 后调用。
- **§2 key**：端口不解析 key 结构；无线程校验新增（既有零校验保持）；复合维度未引入 ✓。
- **§3 窗口语义**：`RateLimiterInWindow` = `t.After(now.Add(-window))`，与既有 `allow` 剪枝（cutoff = now.Add(-window)；kept = t.After(cutoff)）逐位一致；无策略注册机制 ✓。
- **§4 容量**：`DefaultRateLimiterCapacity = 1 << 16` 常量与既有 `newLoginRateLimiter` 默认一致 ✓。
- **§5 Retry-After**：`RateLimiterRetryAfterSeconds` 与既有 `retryAfterSeconds` 计算逐字一致（remain := oldest.Add(window).Sub(now)；remain<=0 → 1；else int(remain.Round(time.Second)/time.Second)）✓。
- **§6 并发 / §7 停机**：接口注释要求并发安全；无后台协程、无 Start/Stop → VP-021 义务不触发 ✓。
- **§8 红线**：未引入 Redis 依赖（go.mod 零变更）；未改 Profile 默认集 / 模块矩阵 / Manifest；GOAL-014 排除声明落合同 §0/§8；7 处使用点代码零改动 ✓。

## 3. 快测覆盖评估

`kernel/ratelimit_test.go`：编译期端口面断言（stub 实现两接口）+ 常量断言 + 边界表驱动（InWindow 8 例含 cutoff 恰等与零窗口；RetryAfter 7 例含 remain≤0→1 与亚秒 Round 双向）。覆盖合同 §3/§5 的全部可执行谓词与 §10 预告的合同级快测面。

## 4. 越界核账

`git status` 变更面 = `apps/api/kernel/ratelimit.go`、`apps/api/kernel/ratelimit_test.go`、`docs/workspaces/workspace-027-rate-limiter-port/**`——R1 波内合法路径；`go.mod` / `go.sum` / Profile 装配 / `internal/handler` / Charter 零触碰。

## 5. 验证复跑（2026-09-01）

`go vet ./kernel/...` 0 · `go test ./kernel/... -count=1` ok · `go build ./...` 通过。

## Verdict

**pass**（0 required）。R1 合同冻结满足关门条件；建议 A-002（grok build · grok-4.6 · high）independent 复核后合并响应关门。

## Findings

- required：无。
- recommended：无（I-027-002 已登记 R2 前置，不属本目标）。