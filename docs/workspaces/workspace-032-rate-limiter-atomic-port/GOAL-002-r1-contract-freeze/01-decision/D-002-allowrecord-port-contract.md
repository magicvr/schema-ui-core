---
doc_type: goal-decision
id: D-002-allowrecord-port-contract
parent: GOAL-002-r1-contract-freeze
date: 2026-09-03
status: accepted
version: 0.1.0
---

# D-002 · AllowRecord 端口合同 v0.1.0（2026-09-03 冻结）

> **责任文件（frozen）**。R2（14 处使用点迁移 + handler 回归）与 R3（证据矩阵）以本合同为分母。本波落：合同正文 + `kernel.RateLimiter.AllowRecord` + 编译期 stub + Memory 单锁实现 + 合同级测试。不迁生产调用点；不实现 Redis；不改 Profile 默认集；不重开 VP-027 关门事实。

## 0. 适用与基线

- **契约面**：`apps/api/kernel` 公共面。加法，不替换 VP-027 D-002（[workspace-027 GOAL-002 D-002](../../../workspace-027-rate-limiter-port/GOAL-002-r1-contract-freeze/01-decision/D-002-rate-limiter-port-contract.md)）。
- **继承不变**：滑动窗口谓词 `RateLimiterInWindow` / `RateLimiterRetryAfterSeconds`；Allow 不注册 key；Record 才创建 map 条目并走容量驱逐；Clear 单锁删 key；无后台协程（VP-021 不停机义务）；key 仍为不透明字符串。
- **本波对象**：Allow→Record 两调用之间的 TOCTOU（GOAL-001 A-008 R-007 residual）。

## 1. 端口形状（I-032-001）

```go
type RateLimiter interface {
    Allow(key string, now time.Time) bool
    Record(key string, now time.Time)
    AllowRecord(key string, now time.Time) bool // VP-032 新增
    RetryAfterSeconds(key string, now time.Time) int
    Clear(key string)
}
```

`RateLimiterProvider.NewRateLimiter` **不变**。

### 1.1 `AllowRecord` 语义

在**同一把锁**内执行，顺序等价于：

```text
if !Allow(key, now) { return false }  // 含 Allow 的剪枝；不走 Record
Record(key, now)
return true
```

可执行细则（供应商必须遵守）：

1. 用 `RateLimiterInWindow` 对已有条目剪枝（与 Allow 相同；缺席 key 不创建条目）。
2. 缺席 key：**恒允许**（即使 `max <= 0`），然后走 Record 的登记/驱逐/append。这保持「Allow 对缺席 key 恒 true」的 VP-027 不变量。
3. 剪枝后 in-window 条数 `>= max`：**不写入**、返回 false。对已存在 key，把剪枝结果写回（与 Allow 拒绝路径相同）。
4. 否则：走 Record 路径（缺席则可能 FIFO 驱逐最老 key，再 append `now`），返回 true。
5. **不返回**剩余额度。`RetryAfterSeconds` 仍独立；调用方仅在 `AllowRecord`（或兼容路径 `Allow`）为 false 之后调用。

### 1.2 兼容

- `Allow` 保持无副作用（不注册 key）。
- `Record` 保持无条件 append（不检查 max）——兼容测试与任何仍走拆分路径的调用方。
- `Clear` 不变；**不**提供原子 Clear 变体（I-032-002）。
- 新代码 **应当** 用 `AllowRecord`；拆分方法保留至明确的破坏性清理波（本 VP 不做）。

## 2. 剪枝与容量

| 路径 | 剪枝 | 创建 map 条目 | 容量驱逐 |
|------|------|----------------|----------|
| Allow | 是 | 否 | 否 |
| Record | 否 | 是（缺席时） | 是（缺席且满容） |
| AllowRecord false | 是（同 Allow） | 否 | 否 |
| AllowRecord true | 是（先 Allow）然后 Record | 是（缺席时） | 是（缺席且满容） |
| RetryAfterSeconds | 否（VP-027 D-002 v0.1.1） | 否 | 否 |
| Clear | n/a | 删除 | 从 order 去掉 |

喷洒 distinct key：只有 true 路径的 AllowRecord（或 Record）能撑大 map；false 路径与 Allow 一样不能。

## 3. 并发

- 所有方法（含 AllowRecord）必须并发安全。
- **原子性判据**：同一 `key`、同一 `now`、无 Clear，N 个 goroutine 并发 `AllowRecord`，返回 true 的次数 **恰好等于** `min(N, max)`（不得穿透预算）。
- `-race` 覆盖 AllowRecord 与既有方法混合调用。

## 4. 使用点迁移口径（I-032-002 · R2 实施，R1 只冻口径）

两种既有模式，原子性目标相同：

| 模式 | 现状 | R2 迁移 | 单请求净状态 |
|------|------|---------|--------------|
| **立即消费** | `Allow` 后立刻 `Record`，永不 `Clear` | `if !AllowRecord(...) { 429 }` | 与旧 Allow→Record 等价 |
| **失败预算** | 入口 `Allow`；失败 `Record`；成功 `Clear` | 入口 `AllowRecord`（乐观占槽）；失败不再 `Record`；成功仍 `Clear` | 成功：`Clear` 后与旧 Allow+Clear 等价；失败：与旧 Allow+Record 等价。并发下比旧 TOCTOU **更保守**（槽位在 Clear 前被占用） |

不提供 `AllowRecordClear` / 条件 Record。失败预算的「更保守」是修复本身，不是产品扩展。

## 5. 冻结使用点分母（R2 迁移表 · HEAD `b1c03acd` 扫描）

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

**显式排除**：recyclebin 领域 `Record`；GOAL-014 分层锁定；测试桩（R1 已补 `AllowRecord` 以满足接口，不计生产分母）。

窗口/阈值/key 形状保持 W12 D-002 + VP-030 三桶现状；本 VP 不改常量。

## 6. 停机与生命周期

AllowRecord 不引入后台协程。VP-021 义务不触发。内存态随进程消失。

## 7. 红线

- 不重开 / 改写 VP-027 关门事实。
- 不实现 Redis / **不消耗 RT-Q05**。
- 不改 Profile 默认集 / 模块矩阵 / Manifest（VP-008 `go`）。
- 不改其它内核端口；不返回剩余额度；不删 Allow/Record。
- R1 **不**修改 §5 十四处生产调用点。

## 8. 阶段切分

| 阶段 | 本 VP 交付 |
|------|------------|
| **R1（本目标）** | D-002；kernel 接口 + stub；Memory.AllowRecord（`allowLocked`+`recordLocked` 单锁组合）；顺序等价测试；并发预算测试；`-race` |
| **R2** | §5 十四处按 §4 口径迁移；handler 既有限流测试仍绿；失败预算路径去掉二次 Record |
| **R3** | 证据矩阵 / 越界核账 / 审计闭合 |

## 9. 验收方式

- **R1 合同级测试（本目标 C2）**：
  - `kernel/ratelimit_test.go`：stub 实现 `AllowRecord`；既有 InWindow / RetryAfter / capacity 常量断言保持。
  - `internal/ratelimit/memory_test.go`：AllowRecord 顺序 ≡ Allow-then-Record；拒绝路径不登记新 key；N 并发 true 次数 = max；既有 Allow 不注册 / 滑动窗口 / RetryAfter / 容量回落回归；`TestMemoryConcurrent` 混入 AllowRecord。
- **R2**：十四处零 Allow→Record 生产配对（立即消费与失败预算入口均改为 AllowRecord）；handler 套件绿。
- **R3**：判据 #1～#5 证据矩阵 + `go.mod` 无 redis。

## 10. 未选方案

| 项 | 未选 | 理由 |
|----|------|------|
| 返回 `(bool, remaining)` / 内嵌 Retry-After | 未选 | I-032-001：bool 足够；RetryAfterSeconds 独立 |
| 把 AllowRecord 做成 Allow 的副作用 | 未选 | 破坏 VP-027「Allow 不注册」不变量 |
| 原子 Clear / CompareAndClear | 未选 | I-032-002：Clear 已是单锁 |
| 推迟 Memory 实现到 R2 | 未选 | Go 接口编译约束（D-001） |
| 本波删除 Allow/Record | 未选 | 兼容；破坏性清理另波 |
| Redis / 令牌桶 | 未选 | 红线 |

## 修订史

| date | version | change |
|------|---------|--------|
| 2026-09-03 | 0.1.0 | 初版冻结：AllowRecord 单锁等价、兼容、剪枝/容量表、14 处分母、R1/R2 切分（GOAL-002 C2） |
| 2026-09-04 | 0.1.1 | **更正（GOAL-003 A-002 证伪）**：§4 失败预算口径（入口 `AllowRecord` + 成功 `Clear`）被证伪——键级 `Clear` 无法只回滚当次占槽、会连历史一起清空（已证实回归：CAPTCHA 清登录历史、recovery no-path 不累计）。§4 失败预算口径由 [GOAL-003 D-002](../../GOAL-003-r2-handler-migration/01-decision/D-002-tokenized-reservation-failure-budget.md)（令牌化 `Reserve`/`Cancel` + 逐路径语义冻结，2026-09-04 用户裁决方案 A）取代；§1 端口形状、§2、§3、§5–§7 不变。 |
