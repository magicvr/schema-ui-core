---
id: VRev-073-vp032-rate-limiter-atomic-port-activation
doc_type: vision-review
title: VP-032 激活就绪 · 限流器端口原子化（AllowRecord）
source: self
date: 2026-09-03
scope: VP-032-rate-limiter-atomic-port 意图 / 退出判据 / 非目标 / P-005（I-032-001/002 冻结） / 架构类 freshness（`b1c03acd`）
verdict: pass
open_required: 0
status: active
created: 2026-09-03
updated: 2026-09-03
parent: null
version: 0.1.0
---

# VRev-073 · VP-032 激活就绪（限流器端口原子化）

## 背景与触发

用户 2026-09-03 指令：「/vision 走流程激活 VP-032（架构原子限流），然后交 /govern 开设工作区」。VP-032（架构分支 · `kernel.RateLimiter` 原子 `AllowRecord` · 承接 [workspace-030] GOAL-001 A-008 R-007 residual）经 [VRev-071](VRev-071-vp032-rate-limiter-atomic-port-planned.md) 计划阶段 self `pass`（0 required；**不是**激活许可）。本次为激活就绪审视：意图/退出判据/非目标/P-005 + **架构类 freshness** + **退出分母正式冻结**（I-032-001/002）。

## 审视要点

### 1. 意图与退出判据可判定性

**pass**。VP-032 v0.1.0 五条退出判据（原子性 / 行为等价 / 兼容 / 边界保持 / 审计闭合）均可核验。本审视冻结签名与使用点分母（§3），使判据 1/2 在方案冻结前可实施；其余属工作区 R1–R3 证据，不使方向不可判定。

### 2. 非目标与红线

**pass**。不重开 / 改写已 closed 的 VP-027 关门事实；不实现 Redis / 分布式限流（RT-Q05 触发条件不变）；不改其它内核端口；不把 R-007 之外的 recommended（R-004/R-009）卷入；不改 Profile 默认集 / Manifest（VP-008 `go` 红线）。`Allow`/`Record` 保留兼容，不删除。

### 3. 信息需求（P-005）——激活前冻结

**pass**。I-032-001 / I-032-002 原为 `open（激活前 /vision 裁决）`。本审视按用户激活指令冻结（写入 VP-032 v0.2.0）：

#### I-032-001 · 签名与返回值（required · 判据 1）

| 项 | 冻结 |
|----|------|
| 签名 | `AllowRecord(key string, now time.Time) bool` |
| 返回值 | **bool 足够**。不返回剩余额度 / Retry-After。`RetryAfterSeconds` 仍为独立方法，仅在 `AllowRecord`（或兼容路径 `Allow`）返回 false 后调用。 |
| 单锁语义 | 同一把锁内：按 `RateLimiterInWindow` 修剪窗口 → 若 in-window 条数 `>= max` 则 **不写入** 并返回 false → 否则 append `now` 并返回 true。 |
| 兼容 | `Allow` 保持无副作用（不注册 key）；`Record` / `Clear` / `RetryAfterSeconds` 语义不变。 |

#### I-032-002 · 迁移范围（required · 判据 2）

**全部生产 Allow→Record 调用点迁移**；Clear-on-success **不需要**原子变体（`Clear` 本就是单锁）。

两种既有模式，迁移口径不同、原子性目标相同：

| 模式 | 现状 | 迁移 |
|------|------|------|
| **立即消费**（请求计数） | `Allow` 后立刻 `Record`，永不 `Clear` | `if !AllowRecord(...) { 429 }` |
| **失败预算**（先检查后记失败） | 入口 `Allow`；失败路径 `Record`；成功 `Clear` | 入口改为 `AllowRecord`（乐观占槽）；失败路径 **不再** `Record`；成功仍 `Clear` |

失败预算的单请求净状态：成功路径在 `Clear` 之后与旧 `Allow`+`Clear` 等价；失败路径与旧 `Allow`+`Record` 等价。并发下比旧 TOCTOU **更保守**（槽位在 `Clear` 前被占用）——这正是消除穿透窗口的方向，不是产品行为扩展。

**冻结使用点分母（代码扫描 · 2026-09-03 HEAD `b1c03acd`）**：

| # | 使用点 | 位置 | 模式 |
|---|--------|------|------|
| 1 | 登录失败桶（含 MFA 签发二次检查） | `apps/api/internal/handler/auth.go` | 失败预算 |
| 2 | 验证码生成 | `apps/api/internal/handler/captcha.go` | 立即消费 |
| 3 | 密码修改 | `apps/api/internal/handler/account_self.go` | 失败预算 |
| 4 | 自助恢复 start | `apps/api/internal/handler/recovery.go` | 失败预算 |
| 5 | 自助恢复 complete（`recordFailure`） | `apps/api/internal/handler/recovery.go` | 失败预算 |
| 6 | MFA verify 独立桶 | `apps/api/internal/handler/mfa.go` | 失败预算 |
| 7 | MFA step-up enroll | `apps/api/internal/handler/mfa.go` | 失败预算 |
| 8 | MFA step-up disable | `apps/api/internal/handler/mfa.go` | 失败预算 |
| 9 | MFA step-up recovery-rotate | `apps/api/internal/handler/mfa.go` | 失败预算 |
| 10 | 邀请接受 | `apps/api/internal/handler/invites.go` | 失败预算 |
| 11 | 钱包核销 | `apps/api/internal/handler/wallet_self.go` | 失败预算 |
| 12 | Telegram IP 桶 | `apps/api/internal/channel/telegram/webhook.go` | 立即消费 |
| 13 | Telegram chat 桶 | `apps/api/internal/channel/telegram/webhook.go` | 立即消费 |
| 14 | Telegram user 桶 | `apps/api/internal/channel/telegram/webhook.go` | 立即消费 |

**显式排除**：`recyclebin` 的 `Record`（回收站领域方法，不是 `kernel.RateLimiter`）；GOAL-014 账号分层锁定（DB 行锁，VP-027 已排除）；测试桩须补 `AllowRecord` 以满足接口，但不计入生产分母。

### 4. 架构类轻量 freshness（`42036a3c` → `b1c03acd`）

**PASS**，不暂挂 `go`：

| 域 | 变更 | 判定 |
|----|------|------|
| 协议 pin / provenance（`apps/web/src/protocol/upstream`） | 零变更 | ✅ |
| 依赖锁（go.mod / go.sum / package.json / lockfiles） | 零变更 | ✅ |
| 迁移台账 | `modules/channel/telegram/migration` + `store/migrate_test.go` | ✅ 区间 = VP-030 已审结目（catalog 0066 AES-GCM） |
| Profile 装配（`kernel/profile.go`） | `mvp`/`admin` **默认模块 ID 列表零变更**；`BuiltinModules` 增 `channel.telegram`（compiled candidate，注释写明不进默认集） | ✅ 属 VP-030 已审结目；**未**把新模块塞进默认集 |
| 区间代码变更 | workspace-030 R1–R5 通道运行时 + Admin UI tab + A-008 响应（含 R-007 → VP-032 登记） | ✅ 不涉及内核端口面以外新契约 / Store 方言 / Manifest 装配语义；限流调用点增量（telegram 三桶）已列入本 VP 分母 #12–#14 |

消费候选 HEAD `b1c03acd`。本 VP 属架构分支、**不是**业务域 VP，H-002「业务域 VP 激活前确认同进程主要形态」发现机制不适用；H-002 仍为 Charter 冻结假设（无反证）。不消耗 RT-Q05 Redis trigger。

### 5. 组合对齐

**pass**。VP-032 `vision_ref` = `schema-ui-core-admin-foundation@0.4.0` 精确匹配现行 Charter；roadmap 组合表第 32 行 / RT-Q05 注记为本波同步对象；lead `workspace-032-rate-limiter-atomic-port` 沿用 VP-013～030 slug 惯例（VP-032 正文已预登记）。不重开 VP-027。VP-030 实现层 Root 已 `done` 4/4+R5，但 VP 文件仍 `active`——见 recommended V-F117；**不阻断**本激活（用户本轮明确激活 032；032 是 030 的 R-007 承接波，不是第二套架构方向）。

## Verdict

**pass**（0 required）。

VP-032 意图/判据/非目标/信息需求已就绪，I-032-001/002 已冻结，架构类 freshness PASS（不暂挂 `go`），可激活并交 `/govern` 开区。

## Findings

### 必改（required）

无。

### 建议（recommended）

#### V-F117 · VP-030 实现已结项但愿景仍 `active`

- **级别**：recommended
- **证据**：`workspace-030` Root `GOAL-001-telegram-channel-runtime` `done` 4/4+R5；VP-030 文件 `status: active` v0.2.2。VRev-072 V-F116 已建议激活 VP-033 前先关 VP-030。
- **影响**：组合卫生（架构分支同时挂一个已交付、一个新激活的 delivery VP）。alignment §7 关门须用户确认，本审视不改 VP-030 status。
- **关闭要求**：后续 `/vision` 对 VP-030 做关门 Review（用户确认后 `active → closed`）；不作为 VP-032 开区或实施门禁。

## 声明

本意见不直接修改 Charter / VP / Goal status。required finding 的响应由 `/vision` 追加在本报告中；原 verdict 与 finding 原文不得改写。本轮用户指令授权激活，status 变更由 `/vision` 写入 VP-032 并交 `/govern` 开区。
