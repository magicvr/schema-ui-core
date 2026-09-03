---
id: GOAL-003-r2-handler-migration
title: R2 生产使用点迁移与 handler 回归
status: active
parent: GOAL-001-rate-limiter-atomic-port
created: 2026-09-03
updated: 2026-09-03
version: 0.1.0
progress: 2/3
plan_refs:
  - VP-032-rate-limiter-atomic-port
primary_plan: VP-032-rate-limiter-atomic-port
serves_summary: 承载 VP-032 R2：依据 D-002 冻结合同将 14 处生产 Allow→Record 调用点全部迁移为 AllowRecord（立即消费 4 处 + 失败预算 10 处），消除生产 TOCTOU 并确保既有 handler 限流回归全绿。
---

# GOAL-003 · R2 生产使用点迁移与 handler 回归

## 概述

执行 Root 纲领 **R2**：承接已关门之 R1（GOAL-002）冻结的 `AllowRecord` 端口合同（D-002 v0.1.0），将 HEAD（`b1c03acd` / `98edb03e`）扫描确定的 **14 处生产 Allow→Record 调用点**全部迁移为原子 `AllowRecord`。
迁移严格遵守 D-002 §4 口径：立即消费（4 处）直接以 `if !AllowRecord(...) { 429 }` 替换旧两段式调用；失败预算（10 处）在入口乐观占槽并移除失败分支的二次 `Record`，成功路径保持单锁 `Clear`。迁移后执行 handler 与 channel 限流回归测试，验证行为等价且并发下消除 TOCTOU 穿透。

## 纲领检查点（P-001）

| 检查点 | 内容 | 状态 |
|--------|------|------|
| C1 | **14 处生产调用点迁移**：按立即消费（4 处）与失败预算（10 处）两口径完成代码改造，生产代码消除 Allow→Record 配对 | **已完成**（E-002 · 14 处全迁） |
| C2 | **测试套件回归与并发验证**：handler / channel 既有限流单测全绿；验证立即消费单请求等价与失败预算 Clear 后净状态等价 | **已完成**（E-002 · 并发无穿透与 Clear 测试通过） |
| C3 | **审视与阶段关门**：自审（A-001）对照 14 处分母与红线，达成无开放 required 关门 | 进行中（待自审与独立审计） |

`progress` = 已关门检查点数 / 3。当前 **2/3**。

## 成功标准（对应 VP-032 退出判据 #2 / #3 / #4）

1. **分母全覆盖**：D-002 §5 冻结的 14 处生产调用点 100% 迁移至 `AllowRecord`，生产环境不再存在 Allow→Record 两段式调用。
2. **行为等价**：立即消费路径（4 处）单请求与旧 Allow→Record 等价；失败预算路径（10 处）成功时 Clear 净状态等价，失败时不再二次 Record。
3. **回归全绿**：`go test ./internal/handler/...` 与 `./internal/channel/telegram/...` 既有限流测试全绿，不破损现有业务流。
4. **红线保持**：不重开 VP-027；不实现 Redis / 不消耗 RT-Q05；不改 Profile 默认集；不改动其它内核端口；`Allow`/`Record` 保持兼容声明。

## 14 处生产调用点迁移清单（D-002 §5 分母）

| # | 使用点 | 位置 | 模式 | 迁移口径 |
|---|--------|------|------|----------|
| 1 | 登录失败桶（含 MFA 签发二次检查） | `apps/api/internal/handler/auth.go` | 失败预算 | 入口 `AllowRecord`；失败不再 `Record`；成功 `Clear` |
| 2 | 验证码生成 | `apps/api/internal/handler/captcha.go` | 立即消费 | `Allow`+`Record` → `AllowRecord` |
| 3 | 密码修改 | `apps/api/internal/handler/account_self.go` | 失败预算 | 入口 `AllowRecord`；失败不再 `Record`；成功 `Clear` |
| 4 | 自助恢复 start | `apps/api/internal/handler/recovery.go` | 失败预算 | 入口 `AllowRecord`；失败不再 `Record`；成功 `Clear` |
| 5 | 自助恢复 complete（`recordFailure`） | `apps/api/internal/handler/recovery.go` | 失败预算 | 入口 `AllowRecord`；`recordFailure` 不再二次 `Record` |
| 6 | MFA verify 独立桶 | `apps/api/internal/handler/mfa.go` | 失败预算 | 入口 `AllowRecord`；失败不再 `Record`；成功 `Clear` |
| 7 | MFA step-up enroll | `apps/api/internal/handler/mfa.go` | 失败预算 | `guardMFAStepUp` 改为 `AllowRecord`；失败不再 `Record`；成功 `Clear` |
| 8 | MFA step-up disable | `apps/api/internal/handler/mfa.go` | 失败预算 | 同 #7 |
| 9 | MFA step-up recovery-rotate | `apps/api/internal/handler/mfa.go` | 失败预算 | 同 #7 |
| 10 | 邀请接受 | `apps/api/internal/handler/invites.go` | 失败预算 | 入口 `AllowRecord`；失败不再 `Record`；成功 `Clear` |
| 11 | 钱包核销 | `apps/api/internal/handler/wallet_self.go` | 失败预算 | 入口 `AllowRecord`；失败不再 `Record`；成功 `Clear` |
| 12 | Telegram IP 桶 | `apps/api/internal/channel/telegram/webhook.go` | 立即消费 | `Allow`+`Record` → `AllowRecord` |
| 13 | Telegram chat 桶 | `apps/api/internal/channel/telegram/webhook.go` | 立即消费 | `Allow`+`Record` → `AllowRecord` |
| 14 | Telegram user 桶 | `apps/api/internal/channel/telegram/webhook.go` | 立即消费 | `Allow`+`Record` → `AllowRecord` |

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-032-001 | required | `AllowRecord` 精确签名与返回值 | 方案冻结 + 判据 1 | R1 | GOAL-002 冻结 | **verified** | — | `AllowRecord(key string, now time.Time) bool`（D-002 §1） |
| I-032-002 | required | 14 处使用点全迁口径与 Clear 语义 | 实施门禁 + 判据 2 | R2 | GOAL-002 冻结 | **verified** | — | 14 处全迁入 R2；立即消费 vs 失败预算两口径冻结于 D-002 §4/§5 |

本目标继承 verified 信息状态，无开放 required 信息项。

## 父目标

- `GOAL-001-rate-limiter-atomic-port`（Root · 纲领 R2）

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺记账。
