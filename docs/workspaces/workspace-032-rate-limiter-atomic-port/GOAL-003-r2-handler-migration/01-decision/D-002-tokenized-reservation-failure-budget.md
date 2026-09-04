---
doc_type: goal-decision
id: D-002-tokenized-reservation-failure-budget
parent: GOAL-003-r2-handler-migration
date: 2026-09-04
status: accepted
version: 0.1.0
supersedes: GOAL-002-r1-contract-freeze/01-decision/D-002-allowrecord-port-contract.md (§4 失败预算口径)
---

# D-002 · 令牌化保留（Reserve/Cancel）修复失败预算语义（2026-09-04）

## 1. 背景与决策链

- A-002（independent · fail / 2 required）证实：失败预算路径的 `AllowRecord` 入口乐观占槽 + **键级 `Clear`** 无法「只回滚当次占槽」，会连历史一起清空，造成已证实的语义/安全回归：
  - `auth.go` 无效 CAPTCHA 分支 `Clear` 可清既有登录失败历史（重置锁定预算）；
  - `recovery.go start` no-path（`ErrRecoveryNotAvailable`）占槽后落到末尾 `Clear`，不再累计 → 枚举探测无速率限制；
  - `recovery.go complete` second-factor-required / INVALID_PASSWORD / 策略违规分支新增 `Clear`（旧代码无操作）；
  - `mfa.go verify` / `invites.go` / `recovery start` 成功路径新增 `Clear`（旧代码无操作）。
- A-001（self · pass）的「失败预算行为等价已达成」主张与代码事实不符，被 A-002 否定。
- **用户裁决（2026-09-04 · P-004）**：选择**方案 A · 令牌化保留**——扩展内核端口为 `Reserve`/`Cancel`，失败预算路径以「原子保留 + 非计数分支 Cancel 只回滚当次占槽」实现 1:1 旧语义 + 消除 TOCTOU。

## 2. 端口扩展（I-032-003 · required）

`kernel.RateLimiter` 新增（加法，不删既有方法）：

```go
Reserve(key string, now time.Time) (token uint64, ok bool)
Cancel(key string, token uint64)
```

**`Reserve` 语义**（与 AllowRecord 同一把锁内原子）：

1. 对已有条目剪枝（同 Allow/AllowRecord 拒绝路径）；缺席 key 恒允许。
2. 剪枝后 in-window 条数 `>= max`：**不写入**、返回 `(0, false)`。
3. 否则：登记一条带**唯一 token** 的占用（走 Record 的登记/驱逐/append 路径），返回 `(token, true)`。
4. 占用**立即计入预算**（并发下保守：在飞请求占用槽位直至 Cancel），原子性判据同 D-002 §3：N 并发 `Reserve` 返回 true 次数恰好 `min(N, max)`（无 Clear/Cancel 时）。

**`Cancel` 语义**：

1. 仅删除该 key 中 token 标识的**那一条**占用；**保留**该 key 其余全部历史。
2. 槽位已被剪枝/Clear 移除时为 no-op（幂等）。
3. 不是预算重置：需要整体重置仍用 `Clear`。

**为何不能只用 `AllowRecord` + `Clear`**：`AllowRecord` 只能原子表达「每次放行都立即消费」；对「只计指定失败」的路径，结果已知前乐观占槽、之后只有键级 `Clear`，无法同时保留历史失败并回滚当次非计数占槽。`Cancel(key, token)` 是唯一能区分「当次占槽」的原语（无 token 的「移除最近一条」在并发下会误删他人槽位）。

## 3. 逐路径语义冻结（10 处失败预算 · 依据 OLD `b08798d4^` 源码事实）

> 冻结原则：**每种结果 = 旧代码的计数行为**。旧 `Record` 的分支 = 保留槽位（计数）；旧无副作用的分支 = `Cancel`（净 0 且保留历史）；旧有 `Clear` 的成功 = 保持 `Clear`（整体重置，旧语义）。立即消费 4 处（#2 验证码、#12–#14 Telegram）**不改变**，保持 `AllowRecord`。

| # | 使用点 | 入口 | 计数（保留槽） | 非计数（Cancel） | 成功 |
|---|--------|------|---------------|-----------------|------|
| 1 | 登录失败桶（auth.go） | `Reserve` | locked/disabled、invalid creds、MFA 签发（旧在签发点 Record 1 次 = 入口 1 次） | 无效 CAPTCHA（旧无副作用） | `Clear`（旧有） |
| 3 | 密码修改（account_self.go） | `Reserve` | 当前密码错误（旧 Record） | — | `Clear`（旧有） |
| 4 | 自助恢复 start（recovery.go） | `Reserve` | `ErrRecoveryNotAvailable`（旧 Record） | Cooldown / SendFailed / 其它 500（旧无副作用） | `Cancel`（旧净 0，无 Clear） |
| 5 | 自助恢复 complete（recovery.go） | `Reserve` | ResolveTarget err、Expired、NotPending、mismatch、二次因子错误（旧 recordFailure 各 1 次） | Evaluate 500、second-factor-required、INVALID_PASSWORD（基线/策略）、hash 500、CompleteRecovery err（含 NotPending 竞态）（旧无副作用） | `Cancel`（旧净 0，无 Clear） |
| 6 | MFA verify（mfa.go） | `Reserve` | Verify fail（旧 Record） | body 解析失败（旧无副作用） | `Cancel`（旧净 0；b08798d4 误用 `Clear`，修正） |
| 7 | MFA step-up enroll（mfa.go） | `Reserve` | 当前密码错误（旧 Record） | UserByID 500（旧无副作用） | `Clear`（旧有） |
| 8 | MFA step-up disable（mfa.go） | `Reserve` | `ErrMFAInvalid`（旧 Record） | 其它 Disable err（旧无副作用） | `Clear`（旧有） |
| 9 | MFA step-up recovery-rotate（mfa.go） | `Reserve` | `ErrMFAInvalid`（旧 Record） | 其它 RotateRecovery err（旧无副作用） | `Clear`（旧有） |
| 10 | 邀请接受（invites.go） | `Reserve` | PeekInviteToken err、AcceptInvite err（旧 Record） | body 解析失败、INVALID_PASSWORD、hash 500（旧无副作用） | `Cancel`（旧净 0；b08798d4 误用 `Clear`，修正） |
| 11 | 钱包核销（wallet_self.go） | `Reserve` | body 解析失败、code 空、RedeemForUser err（旧 Record 各 1 次） | — | `Clear`（旧有） |

**守卫条件**：所有 `Cancel` 仅在 `Reserve` 返回 `ok=true` 的分支可达；`Reserve` 返回 false（429）路径无 token、不 Cancel。`rateLimiter == nil` 时入口直接跳过（旧行为一致）。

## 4. 对 R1 合同的修订关系

- 本决策**取代** GOAL-002 D-002 v0.1.0 **§4 失败预算口径**（「失败预算 = 入口 AllowRecord + 成功 Clear」被证伪）；**其余**（§1 端口形状、§2 剪枝/容量、§3 并发、§5 分母、§6–§7 红线）不变。
- GOAL-002 D-002 修订史追加 v0.1.1 更正条目，指向本决策。
- I-032-002 结论修正：`Clear` 键级删除**无法**回滚当次占槽 → 新增 I-032-003（Reserve/Cancel 契约）为 required 信息项，R2 实施后 verified。

## 5. 测试要求（闭合 F-002）

1. 合同级（memory_test.go）：Reserve 计数即生效；Cancel 只删当次保留历史；Cancel 幂等；Reserve/Cancel 并发 `true` 次数 = max；AllowRecord/Record 与 Reserve 混合窗口一致。
2. 回归级（handler 测试）：
   - recovery start：同一 IP|account 的 no-path 连续请求精确在预算后 429；
   - auth：已有登录失败历史时，无效 CAPTCHA 不得删除既有失败（预算仍按失败次数累计）；
   - recovery complete / MFA verify / invite accept：至少覆盖「先有历史失败，再走成功或非计数分支」的混合序列，断言结果与上表一致；
   - 复跑 handler / telegram 全量、安全用例 `-race`、`go test -count=1 ./...`。

## 6. 未选方案

| 项 | 未选 | 理由 |
|----|------|------|
| 方案 B：语义重定义为全尝试预算 | 未选 | 产品/安全语义变更（CAPTCHA 失败计入登录预算、recovery complete 400 分支计数 → 合法用户自锁），需逐路径书面裁决且无法保留「CAPTCHA 失败不计」 |
| 方案 C：失败预算回退两段式 | 未选 | VP-032 14/14 退出判据不成立，需 VRev 缩小 VP，核心目的（消除 TOCTOU）只完成 4/14 |
| 仅 `AllowRecord` + `Clear`（现状） | 未选 | A-002 F-001 证伪：无法区分当次占槽 |
| 无 token 的「移除最近一条」 | 未选 | 并发下会误删他人槽位；token 唯一标识当次占用 |

## 修订史

| date | version | change |
|------|---------|--------|
| 2026-09-04 | 0.1.0 | 初版：用户裁决方案 A；Reserve/Cancel 契约（I-032-003）；10 处失败预算逐路径语义冻结；F-002 测试要求 |
