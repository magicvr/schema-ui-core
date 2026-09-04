---
doc_type: goal-execution
id: E-004-r3-evidence-matrix-and-close
parent: GOAL-001-rate-limiter-atomic-port
date: 2026-09-04
checkpoint: pending (root close audit)
status: completed
---

# E-004 · R3 证据矩阵 / 越界核账 / 审计闭合（Root 关门证据）

## 1. VP-032 五条方向级退出判据证据矩阵

| # | 判据 | 结论 | 证据 |
|---|------|------|------|
| 1 | **原子性**：`AllowRecord` 在并发下 check+record 原子，无穿透窗口（有并发回归测试） | **已达成** | `TestMemoryAllowRecordConcurrentBudget`（64 并发 true=8）；`TestMemoryReserveConcurrentBudget`（Reserve 同判据）；`TestLoginRateLimit_ConcurrentNoTOCTOUPenetration`（50 并发 20 通过 / 30×429）；`TestWebhook_RateLimiting_ConcurrentNoTOCTOU`（100 并发 60/40）；`-race` 全绿 |
| 2 | **行为等价**：14 处全迁；立即消费等价；失败预算净状态等价 | **已达成** | 14 处全迁（E-002/E-003 · 4 处 `AllowRecord` + 10 处 `Reserve`/`Cancel`）；立即消费 4 处单请求等价（`Allow+Record`→`AllowRecord` 恒等）；失败预算逐路径语义冻结于 [GOAL-003 D-002 §3](../GOAL-003-r2-handler-migration/01-decision/D-002-tokenized-reservation-failure-budget.md)（每种结果 = OLD 行为；计数保留槽 / 非计数 `Cancel` 保留历史 / 旧 `Clear` 保持）；五条混合历史回归全绿；grok 独立复审（GOAL-003 A-004）逐路径核对一致 |
| 3 | **兼容**：`Allow`/`Record` 保留（非破坏性），文档标注 `AllowRecord` 为推荐路径 | **已达成** | 接口保留 `Allow`/`Record`/`AllowRecord`/`Reserve`/`Cancel`/`RetryAfterSeconds`/`Clear`（`kernel/ratelimit.go`）；`Allow` 仍无副作用、`Record` 仍无条件 append；go.mod 零变更 |
| 4 | **边界保持**：未重开 VP-027；未实现 Redis；未改 Profile 默认集 | **已达成** | 越界核账见 §2；`git diff 8b1f2f2f^..HEAD` 零碰 redis / profile / manifest / 其它内核端口；未消耗 RT-Q05 trigger |
| 5 | **审计闭合**：开放 required finding = 0（或已合法闭合） | **已达成** | 审计闭合见 §3：GOAL-002（A-001~A-003）与 GOAL-003（A-001~A-004）全部闭合，开放 required = 0 |

## 2. 越界核账（边界审计）

| 边界 | 核账 | 证据 |
|------|------|------|
| 不重开 / 改写 VP-027 关门事实 | ✅ | 无对 workspace-027 文档/决策的改写；端口为加法 |
| 不实现 Redis / 不消耗 RT-Q05 | ✅ | `go.mod`/`go.sum` 无 redis 依赖；无 Redis 实现代码；RT-Q05 trigger 条件未变 |
| 不改 Profile 默认集 / 模块矩阵 / Manifest（VP-008 `go`） | ✅ | 三个 commit（3bfe66c2 / 277d1eb3 / 516cced4）仅 kernel/ratelimit + internal/ratelimit + handler 生产/测试 + workspace-032 治理文档 |
| 不改其它内核端口 | ✅ | `git diff 98edb03e..HEAD -- apps/api/kernel` 仅 `ratelimit.go`/`ratelimit_test.go` |
| `Allow`/`Record` 兼容保留 | ✅ | 接口与方法签名未删；编译期 stub + 全仓 `go test ./...` 通过 |
| 停机义务（VP-021） | ✅ | 无新增后台协程（memory.go 无 goroutine） |
| 工作区红线（workspace.md） | ✅ | `shared_materials_catalog: none`；vision_role delivery；primary_plan 不变 |

## 3. 审计闭合（全工作区）

| 目标 | 审计条目 | 结论 |
|------|----------|------|
| GOAL-002-r1-contract-freeze | A-001（self pass）/ A-002（independent conditional · F-001 已 closed）/ A-003（响应 pass） | **0 开放 required** |
| GOAL-003-r2-handler-migration | A-001（self pass，行为等价主张被 A-002 证伪）/ A-002（independent fail · F-001/F-002 已 closed·fixed）/ A-003（self pass · 修复复审）/ A-004（grok independent pass · 0 required） | **0 开放 required** |

- A-001 与 A-002 的 verdict 冲突已按 P-004 由用户裁决（方案 A）+ 修复 + 复审消除。
- 信息项：I-032-001/003 verified；I-032-002 revised（结论由 I-032-003 承接）；无开放 required 信息门禁。

## 4. VP-032 文案承接（vision 层标记）

- VP-032 §首波冻结「失败预算：入口乐观占槽；`Clear` 保持（无需原子变体）」与退出判据 #2「失败预算路径在 `Clear` 后净状态等价」的**表述**已被 GOAL-003 D-002 取代（实际语义 = 令牌化 `Reserve`/`Cancel` 逐路径冻结，判据 #2 意图「行为等价 + 并发更保守」仍达成）。
- 建议 VP-032 关门时（/vision）在计划短史中登记本承接关系并评估 VRev；不构成本工作区实施门禁。

## 5. Git Checkpoint

- 本 R3 证据随 Root 关门审计（A-001 self + grok A-002 independent）落盘后提交。
