---
doc_type: vision-plan
id: VP-032-rate-limiter-atomic-port
title: 限流器端口原子化（AllowRecord/Reserve/Cancel）
status: closed
vision_ref: schema-ui-core-admin-foundation@0.4.0
lead_workspace: workspace-032-rate-limiter-atomic-port
created: 2026-09-03
updated: 2026-09-04
version: 0.3.0
parent: null
---

# VP-032 · 限流器端口原子化（AllowRecord）

## 状态与激活门禁

| 项 | 值 |
|----|-----|
| status | **`active`**（2026-09-03 · v0.2.0 · 用户指令激活 · lead `workspace-032-rate-limiter-atomic-port`） |
| lead_workspace | `workspace-032-rate-limiter-atomic-port`（2026-09-03 开区） |
| Vision required | 计划阶段审视 self = [VRev-071](../reviews/VRev-071-vp032-rate-limiter-atomic-port-planned.md)；激活就绪 = [VRev-073](../reviews/VRev-073-vp032-rate-limiter-atomic-port-activation.md) self `pass`（0 required · 架构类 freshness PASS `42036a3c`→`b1c03acd` · I-032-001/002 冻结） |
| 组合位置 | **架构分支** · VP-027 后续端口语义强化（不重开已 closed 的 VP-027 关门事实，只承接其 residual R-007） |

## 意图

消除 `kernel.RateLimiter` 的 **Allow/Record 两次调用之间的 TOCTOU 窗口**（GOAL-001 A-008 R-007 residual）：当前 webhook/auth/recovery/captcha/mfa/wallet 等使用点均为「先 `Allow` 再 `Record`」，两调用非原子，并发下预算可被穿透（进程内单实例，µs 级窗口）。

本 VP 在端口层新增**原子方法** `AllowRecord(key, now) bool`（check+record 一次加锁完成），并迁移全部生产使用点；内存供应商实现原子语义，Redis 接缝（RT-Q05）保持 trigger-gated。

## 首波冻结（退出分母）

| 项 | 本 VP 交付 | 不进本 VP |
|----|-----------|-----------|
| 端口 | `kernel.RateLimiter` 新增 `AllowRecord(key string, now time.Time) bool`；`Allow`/`Record`/`Clear`/`RetryAfterSeconds` 保留兼容 | 删除现有 Allow/Record；返回剩余额度 |
| 供应商 | `Memory` 实现原子 check+record（单锁内完成）；`RetryAfterSeconds`/`Clear` 语义不变 | Redis 实现（RT-Q05 仍 trigger-gated） |
| 使用点迁移 | 冻结分母 **14 处** Allow→Record 调用点全部迁到 `AllowRecord`（见下表） | 其它内核端口变更；VP-027 关闭事实改写；recyclebin 领域 `Record`；GOAL-014 分层锁定 |
| 测试 | 并发穿透回归（两 goroutine 同时 check+record 不得超预算）+ 各使用点行为等价测试 | 性能基准 |
| Profile | 不改 Profile 默认集 | 改变装配红线 |

### 冻结使用点分母（VRev-073 · HEAD `b1c03acd`）

| # | 使用点 | 位置 | 模式 | 迁移口径 |
|---|--------|------|------|----------|
| 1 | 登录失败桶（含 MFA 签发二次检查） | `apps/api/internal/handler/auth.go` | 失败预算 | 入口 `AllowRecord`；失败不再 `Record`；成功 `Clear` |
| 2 | 验证码生成 | `apps/api/internal/handler/captcha.go` | 立即消费 | `Allow`+`Record` → `AllowRecord` |
| 3 | 密码修改 | `apps/api/internal/handler/account_self.go` | 失败预算 | 同 #1 |
| 4 | 自助恢复 start | `apps/api/internal/handler/recovery.go` | 失败预算 | 同 #1 |
| 5 | 自助恢复 complete（`recordFailure`） | `apps/api/internal/handler/recovery.go` | 失败预算 | 入口 `AllowRecord`；`recordFailure` 不再二次 `Record` |
| 6 | MFA verify 独立桶 | `apps/api/internal/handler/mfa.go` | 失败预算 | 同 #1 |
| 7 | MFA step-up enroll | `apps/api/internal/handler/mfa.go` | 失败预算 | `guardMFAStepUp` 改为 `AllowRecord`；失败不再 `Record`；成功 `Clear` |
| 8 | MFA step-up disable | `apps/api/internal/handler/mfa.go` | 失败预算 | 同 #7 |
| 9 | MFA step-up recovery-rotate | `apps/api/internal/handler/mfa.go` | 失败预算 | 同 #7 |
| 10 | 邀请接受 | `apps/api/internal/handler/invites.go` | 失败预算 | 同 #1 |
| 11 | 钱包核销 | `apps/api/internal/handler/wallet_self.go` | 失败预算 | 同 #1 |
| 12 | Telegram IP 桶 | `apps/api/internal/channel/telegram/webhook.go` | 立即消费 | 同 #2 |
| 13 | Telegram chat 桶 | `apps/api/internal/channel/telegram/webhook.go` | 立即消费 | 同 #2 |
| 14 | Telegram user 桶 | `apps/api/internal/channel/telegram/webhook.go` | 立即消费 | 同 #2 |

**立即消费**：`if !AllowRecord(...) { 429 }`。**失败预算**：入口乐观占槽；`Clear` 保持（无需原子变体）。

## 非目标

- 不重开 / 改写已 closed 的 VP-027 关门事实
- 不实现 Redis / 分布式限流（RT-Q05 触发条件不变）
- 不改其它内核端口
- 不把 R-007 之外的 recommended（R-004/R-009）卷入
- 不新增剩余额度返回值

## 与相邻 VP 的边界

| VP / 分支 | 关系 |
|-----------|------|
| **VP-027** | 承接其关门后 residual R-007；端口语义强化，不重开关门 |
| **VP-030** | 消除的 TOCTOU 直接影响 telegram webhook 三桶限流（分母 #12–#14） |
| **VP-009** | 共享基架安全程序正交；本 VP 是端口级修复波 |
| **RT-Q05** | Redis 实现仍 trigger-gated；本 VP 只做内存供应商原子化，**不消耗** trigger |

## 方向级退出判据

1. **原子性**：`AllowRecord` 在并发下 check+record 原子，无穿透窗口（有并发回归测试）。
2. **行为等价**：冻结 14 处使用点全部迁移；立即消费路径与旧 Allow→Record 单请求等价；失败预算路径在 `Clear` 后净状态等价（并发下更保守）。
3. **兼容**：`Allow`/`Record` 保留（非破坏性），文档标注 `AllowRecord` 为推荐路径。
4. **边界保持**：未重开 VP-027；未实现 Redis；未改 Profile 默认集。
5. **审计闭合**：开放 required finding = 0（或已合法闭合）。

## 信息需求（P-005）

| id | 要回答的问题 | 级别 | 影响门禁 | 最晚阶段 | 状态 |
|----|--------------|------|----------|----------|------|
| I-032-001 | `AllowRecord` 精确签名与返回值语义（bool 是否足够，是否需返回剩余额度）。 | required | 判据 1 | 方案冻结 | **verified**（2026-09-03 · VRev-073：`AllowRecord(key string, now time.Time) bool`；bool 足够；不返回剩余额度；`RetryAfterSeconds` 独立） |
| I-032-002 | 是否所有使用点都应迁移（如 Clear-on-success 调用点是否需要原子变体）。 | required | 判据 2 | 方案冻结 | **revised**（2026-09-03 · VRev-073：14 处全迁 + Clear 无需原子变体 + 两口径。**2026-09-04 修正 · GOAL-003 A-002 证伪**：键级 `Clear` 无法只回滚当次占槽、会连历史一起清空 → 失败预算口径改为**令牌化 `Reserve`/`Cancel`**（GOAL-003 D-002 · 用户裁决方案 A · I-032-003），判据 #2 意图（行为等价 + 并发更保守）仍达成） |
| I-032-003 | 令牌化保留契约（`Reserve`/`Cancel` 签名与逐路径语义冻结）。 | required | 判据 2/5 | 实施阶段 | **verified**（2026-09-04 · GOAL-003 D-002：`Reserve(key, now) (token uint64, ok bool)` + `Cancel(key, token)`；10 处失败预算逐路径冻结；GOAL-003 A-004 grok independent 独立核对一致） |

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| workspace-032-rate-limiter-atomic-port | GOAL-001-rate-limiter-atomic-port | delivery | 2026-09-03 | 唯一 lead；不改变 Charter primary workspace |

## 关门记录

- **`active → closed` v0.3.0 · 2026-09-04 · 用户书面确认**（VRev-074 self `pass` · 0 required）。
- 五条方向级退出判据全部 verified（证据矩阵：workspace-032 Root [E-004](../workspaces/workspace-032-rate-limiter-atomic-port/GOAL-001-rate-limiter-atomic-port/02-execution/E-004-r3-evidence-matrix-and-close.md)）：
  1. 原子性：`AllowRecord`/`Reserve` 单锁原子 + 并发预算/无穿透回归 + `-race`；
  2. 行为等价：14/14 全迁（4 处立即消费 `AllowRecord` + 10 处失败预算 `Reserve`/`Cancel`，逐路径语义冻结于 GOAL-003 D-002 §3）；
  3. 兼容：`Allow`/`Record`/`AllowRecord`/`Reserve`/`Cancel`/`RetryAfterSeconds`/`Clear` 保留，`Allow` 无副作用；
  4. 边界保持：未重开 VP-027、未实现 Redis、未改 Profile 默认集、未消耗 RT-Q05 trigger；
  5. 审计闭合：全工作区开放 required = 0（GOAL-002/003 全闭合；Root A-001 self + A-002 grok independent 双 `pass`）。
- 口径承接：判据 #2「失败预算在 `Clear` 后净状态等价」表述由 GOAL-003 D-002 令牌化 `Reserve`/`Cancel` 取代（键级 `Clear` 无法只回滚当次占槽）；I-032-002 `revised`、I-032-003 `verified`；登记于规划修订短史。
- 关门后跟踪：RT-Q05 Redis 实现保持 trigger-gated；VP-031 激活时按 roadmap 复核进程内限流评估是否仍覆盖业务域流量。
- lead `workspace-032-rate-limiter-atomic-port` 已结项（Root `done` 3/3）。

## 规划修订短史

| date | change |
|------|--------|
| 2026-09-03 | 用户书面裁决（GOAL-001 A-008 R-007 处置）：新建 VP 下一波做端口原子化，承接 `kernel.RateLimiter` Allow/Record TOCTOU residual。登记 `planned` v0.1.0（0 区），退出分母草案待 `/vision` 正式冻结。 |
| 2026-09-03 | 用户指令激活：VRev-073 self `pass`（0 required · 架构类 freshness PASS `42036a3c`→`b1c03acd`）· I-032-001/002 冻结 · 使用点分母 14 处 · `planned → active` v0.2.0 · lead `workspace-032-rate-limiter-atomic-port` 交 `/govern` 开区。V-F117 recommended（VP-030 仍 active）不阻断。 |
| 2026-09-04 | **口径承接登记（GOAL-003 A-002 证伪 · 用户裁决方案 A）**：§首波冻结「失败预算：入口乐观占槽；`Clear` 保持（无需原子变体）」与判据 #2「失败预算路径在 `Clear` 后净状态等价」的**表述**由 [GOAL-003 D-002](../workspaces/workspace-032-rate-limiter-atomic-port/GOAL-003-r2-handler-migration/01-decision/D-002-tokenized-reservation-failure-budget.md)（令牌化 `Reserve`/`Cancel` + 10 处逐路径语义冻结）取代——键级 `Clear` 无法只回滚当次占槽（A-002 证伪）；判据 #2 意图（14 处全迁 / 立即消费等价 / 失败预算净状态等价 / 并发下更保守）仍达成。I-032-002 → `revised`；新增 I-032-003 `verified`。实施与双审证据见 workspace-032（Root `done` 3/3 · A-001 self + A-002 grok independent 双 `pass`）。**关门就绪审视 = VRev-074 self `pass`**（0 required）。 |
