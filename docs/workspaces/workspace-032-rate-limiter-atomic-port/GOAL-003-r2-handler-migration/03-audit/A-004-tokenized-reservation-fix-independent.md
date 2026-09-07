---
id: A-004-tokenized-reservation-fix-independent
doc: audit-entry
record_id: A-004
source: independent
scope: GOAL-003 C3 复审 · A-002 F-001/F-002 令牌化 Reserve/Cancel 修复关门
verdict: pass
status: recorded
parent: GOAL-003-r2-handler-migration
created: 2026-09-04
updated: 2026-09-04
version: 0.1.0
auditor: grok-build (grok-4.6 · reasoning high)
audit_type: finding-closure
open_required: 0
---

# A-004 · A-002 F-001/F-002 令牌化修复独立复审（2026-09-04）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high · `/audit`，项目级独立审计路径 [independent-audit-execution.md](../../../../architecture/independent-audit-execution.md)）
- **类型 / scope**：finding-closure / C3 复审（GOAL-003 R2；A-002 F-001/F-002 修复是否足以关门）
- **verdict**：**pass**
- **open required**：**0**
- **落盘方式**：grok 会话按指令产出意见文本（未直接写文件），由编排器代贴为 A-004（`source: independent` 保留）。

本意见不修改 `status` / `progress` / 检查点 / `goal-tree`；响应、finding 闭合与 C3 关门由 `/govern` 处理。

## 范围与区间

| 项 | 值 |
|----|----|
| 工作区 | `workspace-032-rate-limiter-atomic-port`；Root `GOAL-001-rate-limiter-atomic-port`；canonical `docs/workspaces/workspace-032-rate-limiter-atomic-port/` |
| 焦点目标 | `GOAL-003-r2-handler-migration`（parent: `GOAL-001-rate-limiter-atomic-port`；`active` · `2/3`） |
| 规划对齐 | `primary_plan` = `VP-032-rate-limiter-atomic-port`（`active` v0.2.0）；`vision_ref` = `schema-ui-core-admin-foundation@0.4.0` |
| 冻结依据 | GOAL-003 D-002 v0.1.0（令牌化 Reserve/Cancel + 10 处逐路径表）；对照 OLD `b08798d4^` |
| 实施区间 | commit `3bfe66c2`（当前 `dev` HEAD 为其后文档提交） |
| 覆盖 | Reserve/Cancel 合同；D-002 §3 十处失败预算；立即消费 4 处未改；F-002 指定回归；红线；相对 D-002 的新偏差 |
| 不覆盖 | VP-032 愿景层关门；R3 证据矩阵；全仓 `go test ./...`（A-003 声称；本轮未复跑全模块树） |
| 共享资料 | `shared_materials_catalog: none`；本意见未使用共享资料作为证据 |

## 成果（有证据）

### 1. Reserve/Cancel 合同（memory.go 单锁）

`apps/api/internal/ratelimit/memory.go:174-212`：`Reserve` 与 `Cancel` 均 `mu.Lock` 全程；`Reserve` = `allowLocked` + `recordLocked`（占用立即计入预算）；拒绝路径 `return 0, false` 且不 `append`。`Cancel` 按 `token` 删恰好一条并保留其余；缺席 key / 未知 token / Clear 后再 Cancel 为 no-op；末槽 Cancel 同时从 `attempts` 与 `order` 删除。

合同测试（独立复跑 PASS，含 `-race`）：

- `TestMemoryReserveCancelsOnlyItsOwnSlot`：种子 2 条 + Reserve 占满 max=3 → Cancel 后历史仍为 2；`AllowRecord` 与 `Reserve` 共享窗口。
- `TestMemoryReserveDenyDoesNotGrow`：满员拒绝 token=0 且不增长。
- `TestMemoryReserveCancelNoOpAndCleanup`：幽灵 key / 未知 token / Clear 后 Cancel 不复活；末槽清理 order。
- `TestMemoryReserveConcurrentBudget`：64 并发 Reserve，true 次数恰好 = max=8。
- `TestMemoryReserveCancelConcurrent`：与 AllowRecord/Allow/Clear 交错，`-race` 干净。

`kernel.RateLimiter` 加法保留 `Allow` / `Record` / `AllowRecord` / `Clear`（`apps/api/kernel/ratelimit.go:36-80`）。仓库内唯一生产实现是 `ratelimit.Memory`；测试 stub 已补两方法。

### 2. D-002 §3 十处失败预算 vs OLD `b08798d4^`

逐路径对照（计数 = 保留 Reserve 槽；非计数 = `Cancel`；成功 Clear 仅当旧代码 Clear）：

| # | 使用点 | 入口 | 计数（保留） | 非计数（Cancel） | 成功 | 核对 |
|---|--------|------|--------------|------------------|------|------|
| 1 | auth.go login | `Reserve` | locked/disabled、invalid creds、MFA 签发成功 | 无效 CAPTCHA → `Cancel`（**不再 Clear**） | `Clear`（旧有） | 与表一致。CAPTCHA 回归已消失 |
| 3 | account_self.go | `Reserve`（token 丢弃，表无 Cancel 分支） | 当前密码错误 | — | `Clear`（旧有；仍在 HashPassword 之前，与旧代码相同） | 与表一致 |
| 4 | recovery start | `Reserve` | `ErrRecoveryNotAvailable` 保留并 202 返回（**不再落尾 Clear**） | Cooldown / SendFailed / 默认 500 → `Cancel` | `Cancel`（旧净 0） | 与表一致。no-path 累计 |
| 5 | recovery complete | `Reserve` | ResolveTarget err、Expired、NotPending、mismatch（`failAttempt`）、二次因子错误（`failAttemptMFA`） | Evaluate 500、second-factor-required、INVALID_PASSWORD×2、hash 500、CompleteRecovery err → `Cancel` | `Cancel`（旧净 0） | 与表一致。已删除 `recordFailure` |
| 6 | mfa verify | `Reserve` | Verify fail | body/空 proof → `Cancel` | `Cancel`（旧净 0；纠正 b08798d4 误用 Clear） | 与表一致 |
| 7 | mfa enroll | `guardMFAStepUp`→`Reserve` | 当前密码错误 | UserByID 500 → `Cancel` | `Clear`（旧有；仍在 Enroll 之前） | 与表一致 |
| 8 | mfa disable | 同 #7 | `ErrMFAInvalid` 保留 | 其它 Disable err → `Cancel` | `Clear`（旧有） | 与表一致 |
| 9 | mfa recovery-rotate | 同 #7 | `ErrMFAInvalid` 保留 | 其它 Rotate err → `Cancel` | `Clear`（旧有） | 与表一致 |
| 10 | invites accept | `Reserve` | Peek / AcceptInvite err | INVALID_PASSWORD / hash 500 → `Cancel`；body 解析在 Reserve **之前**（旧亦然，净 0 ≡ Cancel） | `Cancel`（旧净 0；纠正误用 Clear） | 与表一致 |
| 11 | wallet redeem | `Reserve`（token 丢弃） | body 解析失败、空 code、Redeem err | — | `Clear`（旧有） | 与表一致 |

A-002 点名回归已从代码面消失：

- `auth.go` 无效 CAPTCHA 为 `Cancel`，不是 `Clear`。
- `recovery.go` no-path 保留槽并直接 202。
- `recovery.go` second-factor-required / INVALID_PASSWORD / policy / hash / CompleteRecovery err 均为 `Cancel`。
- `mfa.go`、`invites.go` 成功为 `Cancel`。
- `account_self.go`、`wallet_self.go`、step-up 成功仍 `Clear`。

守卫：`Cancel` 仅在 `Reserve` 成功之后；429 路径无 token。`rateLimiter == nil` 时跳过（auth/account_self/recovery）。

### 3. 立即消费 4 处仍为 AllowRecord

`3bfe66c2` 未改 `captcha.go` / `webhook.go`。HEAD 仍为：`captcha.go` `AllowRecord`；`webhook.go` IP / chat / user 三桶 `AllowRecord`。生产 handler / telegram 中已无 RateLimiter `.Allow(` / `.Record(` 配对（`resources.go` 的 `Trash.Record` 除外，非限流端口）。

### 4. F-002 指定测试确实断言原回归

独立复跑（`apps/api`）：

```text
go test -count=1 -run 'Test(LoginRateLimit|LoginInvalidCaptcha|PasswordChange|MFA|Recovery|InviteAccept|Captcha|MemoryReserve|MemoryAllowRecord)' \
  ./internal/handler/... ./internal/ratelimit/...
→ ok handler 6.584s; ok ratelimit 0.446s

go test -count=1 ./internal/channel/telegram/...  → ok 1.510s
go test -count=1 ./internal/handler/...           → ok 41.855s
go test -count=1 -race -run 'Test(LoginInvalidCaptcha|RecoveryStartNoPath|RecoveryCompleteMixed|InviteAcceptSuccess|MFAVerifyMalformed|MemoryReserve)' \
  ./internal/handler/... ./internal/ratelimit/...
→ ok handler 25.363s; ok ratelimit 1.323s
```

| 测试 | 断言 | 若 F-001 未修则会 |
|------|------|-------------------|
| `TestRecoveryStartNoPathAccumulatesTo429` | 同一 IP\|account 连续 20 次 no-path 202，第 21 次 429 | 每次 Clear/Cancel 则永不 429 |
| `TestLoginInvalidCaptchaDoesNotClearFailureHistory` | 1 次真失败 + 19 次无效 CAPTCHA（400）+ 19 次真失败 → 第 21 次真失败 429 | CAPTCHA `Clear` 会把历史打回 0，第 21 次不会 429 |
| `TestRecoveryCompleteMixedHistoryPreserved` | 10 次错误码 + 10 次 second-factor-required（Cancel）+ 10 次错误码 → 第 21 次 429 | 非计数分支 Clear 则第二批 10 次后仍不到预算 |
| `TestInviteAcceptSuccessPreservesHistory` | 5 次 dead-token + 1 次 204 成功 + 5 次 dead-token → 429（budget=10） | 成功 Clear 则成功后 5 次不会打满 |
| `TestMFAVerifyMalformedBodyDoesNotCount` | 5 次 verify 失败 + 10 次畸形 body + 5 次失败 → 第 11 次计数 429（budget=10） | 畸形计数会在畸形轮 429；畸形 Clear 则事后到不了 429 |

以上满足 A-002 F-002 闭合要求 1–3，以及 D-002 §5.2「成功**或**非计数」混合序列。本轮未复跑全仓 `./...`；handler 全量与 telegram 全量已独立复跑。

### 5. 红线

`git show --name-only 3bfe66c2` 仅 13 个路径：`kernel/ratelimit.go` + `ratelimit_test.go`、`internal/ratelimit/memory.go` + `memory_test.go`、8 个 handler 生产/测试文件。未碰 redis、`go.mod`/`go.sum`、Profile 默认集、Manifest、其它内核端口。`Allow`/`Record`/`AllowRecord`/`Clear` 仍在接口上。未重开 VP-027（端口加法在 VP-032 / GOAL-003 D-002 范围内）。

## 对照成功标准（对本 scope：F-001/F-002 修复 + C3 技术门禁）

| 标准 | 结论 | 证据 / 缺口 |
|------|------|-------------|
| A-002 F-001：键级 Clear 不再充当「回滚当次」；十处按旧计数行为 | **已达成** | Cancel 原语 + 上表；CAPTCHA / no-path / complete 400 / mfa·invite 成功均已纠正 |
| A-002 F-002：no-path 累计、CAPTCHA 不清历史、混合序列测试 | **已达成** | 五条新回归独立复跑全绿（含 `-race`） |
| 立即消费 4 处等价 | **已达成** | 仍 AllowRecord；telegram 全量 PASS |
| 红线 / 兼容 | **已达成** | commit 边界 + 接口保留 |
| D-002 原则「每种结果 = 旧计数」的未列表 500 边角 | **修复中**（R-001） | 见 R-001；方向更保守，不能绕过预算 |

## 信息就绪核对（P-005）

| ID | 级别 | 当前投影 | 本次核对 |
|----|------|----------|----------|
| I-032-001 | required | verified | 非本缺口 |
| I-032-002 | required | revised | Clear 无法回滚当次：已被 D-002 / I-032-003 取代；实施与修正结论一致 |
| I-032-003 | required | verified | `Reserve(key, now) (token uint64, ok bool)` + `Cancel(key, token)` 已落地；10 处与 §3 表一致 |

到期且影响本 scope 的开放 required 信息项：**0**。无共享资料引用被当作关闭证据。

## Findings

### F-001（A-002 · high/required）· 复审结论：**可按 fixed 闭合**

A-002 四条实施证据均已不成立：存在只删当次槽的 `Cancel`；无效 CAPTCHA 不再 `Clear`；recovery start no-path 累计至 429；complete 非猜测 400 与 mfa/invite 成功不再键级清空。闭合要求中的「逐路径冻结」已落在 D-002 §3，实现与表一致。

### F-002（A-002 · med/required）· 复审结论：**可按 fixed 闭合**

A-002 列出的四类回归均有测试，且测试失败条件正好是 F-001 的旧错误行为。本轮独立复跑通过。

### F-003（A-002 · low/recommended）

Root `00-meta` / `workspace.md` / `goal-tree` 已投影「令牌化修复完成、待复审关门」。本复审不再视为开放文档滞后。

### R-001 · 登录 500 / MFA 签发失败路径保留 Reserve 槽（旧代码不 Record）

| 字段 | 值 |
|------|----|
| 严重度 | low |
| 建议 | **recommended** |
| status | open（编排器响应中） |
| 影响门禁 | 不阻断 F-001/F-002 闭合或 C3；不是预算绕过 |

**证据：** OLD `b08798d4^` 在 MFA 签发点**成功** `BeginChallenge` 之后才 `Record`；`mfa == nil`、`BeginChallenge` 失败、以及其它 `LOGIN_FAILED` 均无 `Record`。新代码入口 `Reserve` 已占槽，这三条 500 分支没有 `Cancel`，因此计入预算。D-002 §3 #1 只把「MFA 签发」列为计数、把「无效 CAPTCHA」列为非计数，未点名这些 500 分支。

效果是并发下更保守（正是去掉第二次 `Allow` 的 TOCTOU 修复），合法用户在认证/MFA 后端持续 500 时会比旧代码更早 429。攻击者不能借此清历史或绕过锁定。

**建议：** `/govern` 在 D-002 修订史加一行「未列出的内部 500 = 保留槽（保守）」或对这三条加 `Cancel` 以严格 1:1。不建议作为 C3 必改。

### R-002 · GOAL-003 成功标准/分母表与 01-decision 索引仍写 AllowRecord+Clear

| 字段 | 值 |
|------|----|
| 严重度 | low |
| 建议 | **recommended** |
| status | open（编排器响应中） |
| 影响门禁 | 不阻断本修复的技术闭合；关门时应改投影以免用过期成功标准验收 |

`GOAL-003/00-meta.md` 成功标准 #1 仍写「100% 迁移至 AllowRecord」，§14 处表仍写入口 `AllowRecord` + 成功 `Clear`。权威冻结已是 D-002（10 处 `Reserve`/`Cancel`）。`01-decision.md` 索引仍只有 D-001。代码与 D-002 一致；过期投影不是语义回归。

## 必改项汇总

开放必改项数：**0**。

A-002 F-001 / F-002 的关闭证据充分、可重复核对，建议编排器按 **fixed** 闭合。R-001 / R-002 为 recommended，不阻断 C3。

## 与既有意见的异同

| 项 | A-002 independent | A-003 self | A-004 independent（本意见） |
|----|-------------------|------------|------------------------------|
| F-001 键级 Clear 破坏失败预算 | fail / required | 主张 fixed | **同意可闭合（fixed）** |
| F-002 缺混合历史测试 | fail / required | 主张 fixed | **同意可闭合（fixed）**；已独立复跑 |
| 14 处静态覆盖 / 立即消费未改 | 同意迁移 | 同意 | **同意**（现 4 AllowRecord + 10 Reserve） |
| 红线 | 当时同意 b08798d4 | 同意 3bfe66c2 | **同意 3bfe66c2** |
| 未列表 500 占槽 | （当时 Clear 问题掩盖） | 未提 | **R-001 recommended** |
| 00-meta 分母表过期 | F-003 路线图滞后 | 未提 | **R-002 recommended**（F-003 路线图已更新） |
| C3 关门 | fail | 待本复审 | **技术门禁可通过**；正式关门仍由 `/govern` 改 status |

与 A-003 无 verdict 冲突：self 与 independent 对本修复均为 **pass / 0 required**。

## 结论与下一步

GOAL-003 针对 A-002 F-001/F-002 的令牌化修复在合同、十处冻结语义、立即消费未改、回归测试和红线上均可复核。C3 技术 close-out 的独立意见为 **pass**（无开放 required）。

建议给编排器的下一句：

> `/govern workspace-032 GOAL-003 响应 A-002/A-003/A-004：将 A-002 F-001/F-002 标为 fixed；R-001/R-002 按 recommended 处理（D-002 补 500 口径说明、更新 00-meta 分母表与 01-decision 索引）；无开放必改则推进 C3 关门。`

## 声明

本意见 `source: independent`，为 L0 入口分离级交叉意见，不等同于外部法定鉴证。本意见不修改目标 `status` / 检查点 / 派生 `progress` / 方案正文 / `goal-tree`。grok 会话按本次指令产出意见文本后由 `/govern` 代贴为 `A-004` 并更新 `03-audit.md` 索引；闭合与关门由 `/govern` 处理。

## 编排器响应（2026-09-04）

- **R-001 → fixed**：按建议的严格 1:1 路径处理——`auth.go` 三条 LOGIN_FAILED 500 分支（`mfa == nil` / `BeginChallenge` 失败 / 其它 auth 错误）补 `Cancel`；D-002 §3 #1 非计数列与修订史 v0.1.1 同步更新。不阻断 C3，作为 recommended 已闭环。
- **R-002 → fixed**：GOAL-003 `00-meta.md` 成功标准 #1 / §14 处表迁移口径改为 D-002 §3 令牌化表述；`01-decision.md` 索引补 D-002 条目。
- 处理方式：recommended（非 required）且方向明确（严格 1:1 旧语义，已由用户裁决的 D-002 原则覆盖），按非关键决策执行并留痕。
