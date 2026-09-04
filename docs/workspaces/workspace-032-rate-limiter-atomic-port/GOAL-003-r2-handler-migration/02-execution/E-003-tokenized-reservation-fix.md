---
doc_type: goal-execution
id: E-003-tokenized-reservation-fix
parent: GOAL-003-r2-handler-migration
date: 2026-09-04
checkpoint: pending (commit after 复审)
status: completed
---

# E-003 · A-002 响应：令牌化 Reserve/Cancel 修复失败预算语义（F-001/F-002）

## 1. 背景与决策

- A-002（independent · fail / 2 required）证实初版迁移（E-002 / commit `b08798d4`）的失败预算语义回归；`AllowRecord` + 键级 `Clear` 无法只回滚当次占槽、会连历史一起清空。
- 逐条核对源码确认：auth.go 无效 CAPTCHA 清登录历史；recovery start no-path 不累计；recovery complete 多个 400 分支新增 Clear；mfa verify / invites 成功新增 Clear。
- **用户裁决（2026-09-04 · P-004）**：方案 A · 令牌化保留。冻结于 [D-002-tokenized-reservation-failure-budget](../01-decision/D-002-tokenized-reservation-failure-budget.md)（Reserve/Cancel 契约 + 10 处失败预算逐路径语义冻结，依据 OLD `b08798d4^` 源码行为）。

## 2. 实施事实

### 2.1 内核端口（kernel）

- `apps/api/kernel/ratelimit.go`：`RateLimiter` 接口新增
  - `Reserve(key string, now time.Time) (token uint64, ok bool)`：同一把锁内原子 check+占槽，占用立即计入预算（并发保守），拒绝路径不登记；
  - `Cancel(key string, token uint64)`：只删除该 key 中 token 标识的那一条占用，保留其余历史；槽位已剪枝/Clear 时为 no-op。
- `kernel/ratelimit_test.go` stub 补齐两方法。

### 2.2 内存实现（internal/ratelimit）

- `memory.go`：`attempts` 由 `map[string][]time.Time` 改为 `map[string][]attempt{id,t}`；每 limiter 单调 `nextID` 生成 token；`Reserve`/`Cancel` 落地；`allowLocked`/`RetryAfterSeconds`/`recordLocked`/`Clear` 适配；`AllowRecord`/`Record` 行为不变（append 带 token 的条目）。

### 2.3 生产调用点（10 处失败预算按 D-002 §3 冻结语义）

| # | 使用点 | 变更事实 |
|---|--------|----------|
| 1 | `auth.go` login | 入口 `AllowRecord`→`Reserve`；无效 CAPTCHA `Clear`→`Cancel`（保留历史）；locked/disabled/invalid creds/MFA 签发保留槽位；成功 `Clear` 不变 |
| 3 | `account_self.go` changePassword | 入口 `Reserve`；密码错误保留槽位；成功 `Clear` 不变 |
| 4 | `recovery.go` start | 入口 `Reserve`；no-path（NotAvailable）保留槽位并直接回 202（修复落尾 Cancel 清槽回归）；Cooldown/SendFailed/其它 500 `Cancel`；成功 `Cancel`（旧无 Clear） |
| 5 | `recovery.go` complete | 入口 `Reserve`；ResolveTarget err/Expired/NotPending/mismatch/MFA 因子错保留槽位；Evaluate 500/second-factor-required/INVALID_PASSWORD×2/hash 500/CompleteRecovery err `Cancel`；成功 `Cancel`；删除 no-op `recordFailure` |
| 6 | `mfa.go` verify | 入口 `Reserve`；Verify fail 保留槽位；body 解析失败 `Cancel`；成功 `Clear`→`Cancel`（旧净 0） |
| 7 | `mfa.go` step-up enroll | `guardMFAStepUp` 改 `Reserve` 返回 token；UserByID 500 `Cancel`；密码错误保留；成功 `Clear` 不变 |
| 8 | `mfa.go` step-up disable | `Reserve`；`ErrMFAInvalid` 保留槽位；其它 Disable err `Cancel`；成功 `Clear` 不变 |
| 9 | `mfa.go` step-up recovery-rotate | 同 disable |
| 10 | `invites.go` accept | 入口 `Reserve`；Peek/AcceptInvite err 保留槽位；INVALID_PASSWORD/hash 500 `Cancel`；成功 `Clear`→`Cancel`（旧净 0） |
| 11 | `wallet_self.go` redeem | 入口 `Reserve`；body 解析/空 code/Redeem err 保留槽位；成功 `Clear` 不变 |

立即消费 4 处（#2 验证码、#12–#14 Telegram）保持 `AllowRecord` 不变。

## 3. 测试证据（F-002 闭合）

### 3.1 合同级（internal/ratelimit/memory_test.go）

- `TestMemoryReserveCancelsOnlyItsOwnSlot`：Cancel 只删当次保留、种子历史保留、桶回降；
- `TestMemoryReserveDenyDoesNotGrow`：拒绝路径 token=0 且不登记；
- `TestMemoryReserveCancelNoOpAndCleanup`：未知 token / Clear 后 Cancel no-op；末槽 Cancel 删 key 与 order；
- `TestMemoryReserveConcurrentBudget`：64 并发 Reserve true 次数 = max（无穿透）；
- `TestMemoryReserveCancelConcurrent`：与 AllowRecord/Allow/Clear 混合并发 `-race`。

### 3.2 回归级（internal/handler）

- `TestRecoveryStartNoPathAccumulatesTo429`：同一 IP|account 连续 no-path 20 次后第 21 次 429（修复前不累计）；
- `TestRecoveryCompleteMixedHistoryPreserved`：10 次错误码（计数）+10 次 second-factor demand（Cancel 不清历史）+10 次错误码 → 第 21 次 429；
- `TestInviteAcceptSuccessPreservesHistory`：5 次 dead-token + 1 次成功 204（不 wipe）+5 次 dead-token → 429；
- `TestLoginInvalidCaptchaDoesNotClearFailureHistory`：1 失败 + 19 无效 CAPTCHA（不计不 wipe）+19 失败 → 第 21 次 429；
- `TestMFAVerifyMalformedBodyDoesNotCount`：5 失败 + 10 malformed（不计）+5 失败 → 429。

### 3.3 回归与 race

- `go test -count=1 ./...` 全绿（exit 0）；
- `go vet ./...` 0；
- `go test -count=1 -race -run 'Test(LoginRateLimit|LoginInvalidCaptcha|PasswordChange|MFA|Recovery|InviteAccept|Captcha|MemoryReserve|MemoryAllowRecord)' ./internal/handler/... ./internal/ratelimit/...` 全绿（97.8s + 1.4s）；
- `go test -count=1 -race ./internal/channel/telegram/...` 全绿；
- 注：首轮 -race 全量含 Wallet 时 `TestWalletSelfEntriesOwnScope` 报 SQLITE_BUSY（并行负载下的 SQLite 锁抖动）；该用例单独（含 -race）复跑通过，与本次改动无关。

## 4. Git Checkpoint

- 待复审（自审 A-003 + grok build 独立复审）通过后提交；仅显式 owned paths。

## 5. 尚未完成

- 自审 A-003 落盘；grok build 独立复审（项目级路径）；A-002 F-001/F-002 闭合响应；C3 关门；goal-tree 同步。
