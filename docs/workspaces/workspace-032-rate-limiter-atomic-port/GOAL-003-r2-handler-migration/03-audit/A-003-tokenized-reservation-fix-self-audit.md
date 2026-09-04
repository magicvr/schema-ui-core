---
id: A-003-tokenized-reservation-fix-self-audit
parent: GOAL-003-r2-handler-migration
date: 2026-09-04
source: self
auditor: antigravity-govern
audit_type: response-reaudit
scope: A-002 F-001/F-002 修复复审（令牌化 Reserve/Cancel）
verdict: pass
open_required: 0
version: 0.1.0
---

# A-003 · A-002 必改项修复自审（2026-09-04 · self）

- **source**：self
- **auditor**：antigravity-govern
- **类型 / scope**：response-reaudit（响应 A-002 独立审计 F-001/F-002 的修复复审）
- **verdict**：**pass**
- **open required**：0

## 决策链

- 用户裁决（2026-09-04 · P-004）：**方案 A · 令牌化保留**（A-002 推荐路径）。
- 冻结于 [GOAL-003 D-002](../01-decision/D-002-tokenized-reservation-failure-budget.md)：`Reserve(key, now) (token, ok)` + `Cancel(key, token)` 契约（I-032-003）+ 10 处失败预算逐路径语义冻结（依据 OLD `b08798d4^` 行为）。

## 修复证据（对 A-002 逐条）

### F-001（high · required）→ **fixed**

| A-002 证据点 | 修复事实 |
|--------------|----------|
| 键级 `Clear` 删全部历史，无法只回滚当次占槽 | `Cancel(key, token)` 只删 token 标识的那一条，保留其余历史（`memory.go`）；合同测试 `TestMemoryReserveCancelsOnlyItsOwnSlot` / `TestMemoryReserveCancelNoOpAndCleanup` |
| `auth.go` 无效 CAPTCHA 可清登录失败历史 | 改为 `Cancel(limiterKey, token)`；`TestLoginInvalidCaptchaDoesNotClearFailureHistory`（1 失败 + 19 无效 CAPTCHA 不 wipe 不计 + 19 失败 → 21 次 429） |
| `recovery.go start` no-path 不再累计 | NotAvailable 分支保留槽位并直接回 202（不再落尾 Clear）；`TestRecoveryStartNoPathAccumulatesTo429`（20 次 no-path → 21 次 429） |
| `recovery.go complete` 新增 Clear 的 400 分支 | second-factor-required / INVALID_PASSWORD / 策略违规 / hash 500 / CompleteRecovery err 全部改 `Cancel`；`TestRecoveryCompleteMixedHistoryPreserved` |
| `mfa.go` / `invites.go` 成功新增 Clear | 成功改 `Cancel`（旧净 0）；`TestMFAVerifyMalformedBodyDoesNotCount`、`TestInviteAcceptSuccessPreservesHistory` |

保留旧 `Clear` 的成功路径（登录、改密、钱包、step-up 成功）语义不变。

### F-002（med · required）→ **fixed**

| A-002 闭合要求 | 落地 |
|----------------|------|
| no-path 连续请求精确在预算后 429 | `TestRecoveryStartNoPathAccumulatesTo429` |
| 已有失败时无效 CAPTCHA 不删既有历史 | `TestLoginInvalidCaptchaDoesNotClearFailureHistory` |
| recovery complete / MFA verify / invite accept 混合序列 | `TestRecoveryCompleteMixedHistoryPreserved`、`TestMFAVerifyMalformedBodyDoesNotCount`、`TestInviteAcceptSuccessPreservesHistory` |
| 复跑 handler / telegram 全量、安全用例 `-race`、`go test -count=1 ./...` | 全绿（详见 E-003 §3.3） |

## 回归与红线核账

- `go test -count=1 ./...` 全绿（exit 0）；`go vet ./...` 0。
- 安全路径 `-race`（Login/Captcha/PasswordChange/MFA/Recovery/InviteAccept + ratelimit Reserve）全绿；telegram `-race` 全绿。
- 红线保持：commit `3bfe66c2` 仅含 kernel/ratelimit + internal/ratelimit + 8 个 handler 生产文件 + 测试；未碰 redis / go.mod / profile / Manifest / 其它内核端口；`Allow`/`Record`/`AllowRecord`/`Clear` 兼容保留；不重开 VP-027。

## Findings

- 无 required finding。
- 无 recommended finding。
- 注：`TestWalletSelfEntriesOwnScope` 在并行 `-race` 全量下偶发 SQLITE_BUSY（SQLite 锁抖动，独立复跑含 -race 通过），与本次改动无关；非本目标必改项。

## 必改项汇总

- 开放必改项数：**0**（A-002 F-001/F-002 已按 fixed 路径修复；待 grok build 独立复审确认后正式闭合）。

## 结论与下一步

- A-002 的 F-001/F-002 已完成修复并有可核对证据；C1/C2 回归全绿。
- 下一步：调用本地 grok build（grok 4.6 · 思考强度 high）执行独立复审（项目级独立审计路径），落盘 A-004；随后编排器合并响应并推进 C3 关门。
