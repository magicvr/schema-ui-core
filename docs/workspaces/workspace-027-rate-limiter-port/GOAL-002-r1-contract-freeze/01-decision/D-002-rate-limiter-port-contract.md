---
doc_type: goal-decision
id: D-002-rate-limiter-port-contract
parent: GOAL-002-r1-contract-freeze
date: 2026-09-01
status: accepted
version: 0.1.1
---

# D-002 · RateLimiter 端口合同 v0.1.0（2026-09-01 冻结）

> **责任文件（frozen）**。实施（R2 内存供应商 + 7 处使用点迁移）与验收（R4 证据矩阵）以本合同为分母。本波只落端口本体 + 合同级快测；内存供应商、使用点迁移与 trusted-proxy/IP 工具归 R2；Redis 接缝与轨道约定继承/登记归 R3。不改 Profile 默认集 / 模块矩阵 / Manifest 装配；不引入 Redis 客户端依赖；不改 Charter。

## 0. 适用与验收基线

- **契约面**：`apps/api/kernel` 公共面（Go 1.26）；handler 与模块经 `kernel.RateLimiter` 消费，绝不接触供应商类型（VP-003 薄内核）。
- **先例对齐**：`kernel.Cache` / `Store` / `ObjectStore` / `MailSender` —— 非泛型接口 · fail-closed 语义 · 供应商无关。迁移不回归基线 = `internal/handler/rate_limit.go`（滑动窗口 · allow 不注册 key · 容量驱逐 · Retry-After · trusted-proxy `loginClientIP`）。
- **范围外（trigger-gated / 属其它阶段或永久排除）**：Redis 供应商实现（RT-Q05 触发后才评估）；分布式限流 / 跨实例协调；令牌桶 / 漏桶 / 固定窗口实现（策略形态可扩展，实现待消费者触发）；缓存语义（VP-026）与消息/事件语义（VP-028）；业务级配额策略（按用户/组织/路由，归业务域 VP）；调度状态去重 / 幂等守卫（领域状态）；**GOAL-014 账号分层锁定（DB 行锁，不是 `loginRateLimiter`，永不纳入端口）**。

## 1. 端口形状（I-027-001 · 冻结：语义拆分保持）

```go
// kernel.RateLimiter：供应商无关限流端口（进程内滑动窗口预算；D-001 P1）。
type RateLimiter interface {
    // Allow 报告 key 当前是否可尝试；绝不注册 key（防喷洒撑爆 map）。
    Allow(key string, now time.Time) bool
    // Record 为 key 登记一次失败；写入路径唯一会创建 map 条目之处。
    Record(key string, now time.Time)
    // RetryAfterSeconds 返回窗口内最老失败后的剩余秒数；仅限 not-allowed 时调用。
    RetryAfterSeconds(key string, now time.Time) int
    // Clear 清空 key 的全部失败（成功路径调用，防毒化桶）。
    Clear(key string)
}

// kernel.RateLimiterProvider：内核级工厂。handler 注入本接口获取限流器，
// 不接触供应商类型；内存供应商（R2）实现之，Redis 级供应商触发后实现同一接口。
type RateLimiterProvider interface {
    // NewRateLimiter 返回 window/max/capacity 构造的限流器。
    // capacity <= 0 时回落 DefaultRateLimiterCapacity（1<<16）。
    NewRateLimiter(window time.Duration, max, capacity int) RateLimiter
}
```

- **`now` 注入**：全部时间语义显式传 `now time.Time`（既有 `loginRateLimiter` 签名保持；生产调用点传 `time.Now().UTC()`，测试注入固定时钟——确定性不依赖隐藏时钟）。
- **Allow 不注册**：`Allow` 只读剪枝与判定，永不创建条目（D-001 P1：喷洒 distinct key 不能撑爆 map——容量驱逐只在 `Record` 可达）。
- **Retry-After 语义**：`RetryAfterSeconds` 仅限 `Allow == false` 后调用（既有调用点惯例：`sec > 0` 才设置 `Retry-After` 头）；计算唯一权威 = `kernel.RateLimiterRetryAfterSeconds`（§5）。
- **为何不选内聚 / 回调式**：取消「allow 不注册」保证 = 安全语义漂移；回调式难测易错（D-001 未选方案）。

## 2. key 语义（I-027-004 · 冻结：不新增复合维度）

- key = **不透明字符串**，端口不解析结构、不校验形状（对既有 7 个构造点零约束变化）。
- 既有形态保持：`IP|identifier`（登录 / 恢复 / 密码修改 / 验证码 / 邀请接受）、`op|IP|user`（MFA step-up）、纯 IP（MFA verify 独立桶）；`loginClientIP` 的 trusted-proxy 语义（X-Real-IP 仅在受信 CIDR 直连对端时取用）保持——该工具的去向由 R2 迁移决策定（I-027-002）。
- 契约义务：key 非空；供应商不得截断 / 改写 key。建议上限 512 字节（供 R2 供应商断言防御，非端口校验）。
- **复合「路由+用户」维度**：本波不新增；C 端业务域 VP 需要时自行定义 key 协议（不预制）。

## 3. 窗口语义（I-027-003 · 冻结：滑动窗口 + 独立策略形态）

- **滑动窗口**：窗口内失败计数（时间戳列表）；**剪枝由 `Allow` 执行**（顺带丢弃窗口外时间戳）——`RetryAfterSeconds` **不剪枝**（v0.1.1 勘误 · A-002 F-006：对齐既有 `retryAfterSeconds`——不修改 `attempts`，全过期时返回 1；若 R2 实现选择 RetryAfter 剪枝，须显式定义空列表返回值；既有调用序（仅 `Allow == false` 后调用）下与 Allow 已剪枝后的结果等价）。惰性清理，无后台协程、无新生命周期（VP-021 SIGTERM 排空义务不触发）。
- **可执行语义权威（供应商必须使用）**：
  - `kernel.RateLimiterInWindow(t, window, now)` = `t.After(now.Add(-window))`——窗口内的判定谓词（恰在 cutoff 上的时间戳不保留，与既有 `allow` 剪枝逐位一致）。
  - `kernel.RateLimiterRetryAfterSeconds(oldest, window, now)`——Retry-After 计算（秒；`remain <= 0` → 返回 **1**；否则 `remain.Round(time.Second)/time.Second`；与既有 `retryAfterSeconds` 逐位一致）。
- **策略形态不与 VP-026 共用**：RateLimiter 不引入 `ExpiryPolicy` 式可插拔策略接口（限流窗口 vs 缓存过期语义不同）；窗口/阈值随构造参数冻结于各使用点（W12 D-002 常量），不提供策略注册机制（防过度设计）。

## 4. 容量与驱逐（有界 · D-001 P1）

- 供应商必须有界：distinct key 数 ≤ capacity；超限时驱逐**最老插入**的 key（FIFO order，既有实现保持）。
- `capacity <= 0` → `DefaultRateLimiterCapacity = 1 << 16`（既有 `newLoginRateLimiter` 行为保持）。
- 各构造点容量常量保持现状：登录 64K / 验证码 64K / 密码修改 64K / 恢复 64K / MFA verify 64K / MFA step-up 64K / 邀请接受 64K（W12 D-002：不动）。

## 5. W12 D-002 常量与 Retry-After（保持现状）

| 使用点 | 窗口 | 阈值 | key | 保持项 |
|--------|------|------|-----|--------|
| 登录失败（auth.go:60） | 15min | 20 | `IP|identifier` | ✓ |
| 验证码生成（captcha.go:36） | 1min | 10 | `loginClientIP` | ✓ |
| 密码修改（account_self.go:51） | 15min | 5 | `IP|identifier` | ✓ |
| 自助恢复（recovery.go:58） | 15min | 20 | `IP|identifier` | ✓ |
| MFA verify 独立桶（mfa.go:121） | 15min | 10 | 纯 IP | ✓（与登录桶不共用） |
| MFA step-up（mfa.go:129） | 15min | 5 | `op|IP|user` | ✓ |
| 邀请接受（invites.go:308） | 15min | 10 | 纯 IP | ✓（W13 F-001 CPU DoS 刹车） |

官方单实例边界不变；Redis 仅预登记不实施（RT-Q05 trigger-gated 保持）；端口化**不重写**任一语义（R2 迁移回归证据 = 各迁入点既有 handler 测试套件全量通过 + `rate_limit.go` 单元语义快测）。

## 6. 并发安全

- 所有接口方法必须并发安全（互斥保护；多 goroutine 并行 Allow/Record/Clear 无数据竞争）。
- R2 供应商测试以 `-race` 覆盖并发边界（既有 `loginRateLimiter` 为 `sync.Mutex`，保持）。

## 7. 停机与生命周期

- 本波不引入后台协程 / 定时清理 → 无 Start/Stop hook、无 SIGTERM 排空义务（VP-021）。内存态随进程消失。
- 若未来供应商选择后台清理（不选），须按 VP-021 增补停机声明——本合同冻结谓词语义不变。

## 8. 红线

- 不预制 Redis 实现（不引入客户端依赖 / **不消耗 RT-Q05 trigger**）；Redis 轨道约定继承 owner 文档 `cache-redis-seam-and-track.md`（R3 登记，单一所有者）。
- 不改 Profile 默认集 / 模块矩阵 / Manifest 装配（VP-008 `go`）。
- 不重开 VP-012 / 已 closed 记录；不改 Charter。
- GOAL-014 账号分层锁定（DB 行锁）**显式排除**，不纳入端口。
- 7 处使用点完整迁移分母（V-F099）：本目标**不动**既有构造点；R2 全部接入。

## 9. 信息裁决记录

| ID | 裁决 | 证据 |
|----|------|------|
| I-027-001 | 语义拆分保持（§1） | D-001（2026-09-01 用户裁决） |
| I-027-003 | 滑动窗口保持 + 策略独立（§3） | D-001（2026-09-01 用户裁决） |
| I-027-004 | 不新增复合 key（§2） | D-001（2026-09-01 用户裁决） |
| I-027-002 | 迁移策略（演进 vs 双轨） | **R2 前置裁决**（本目标不关闭） |

## 10. 验收方式（R2/R4 预告）

- **合同级快测（本目标，C2）**：`kernel/ratelimit_test.go` —— `RateLimiterInWindow` 边界（cutoff 恰等 / 窗内 / 窗外 / 零窗口）表驱动；`RateLimiterRetryAfterSeconds`（精确剩余 / `remain<=0 → 1` / 亚秒 Round）表驱动；`DefaultRateLimiterCapacity` 常量断言；编译期端口面断言（stub 实现 `RateLimiter` / `RateLimiterProvider`）。
- **R2 供应商与迁移测试**：compile-time 断言内存供应商实现 `kernel.RateLimiter`；行为 / 并发 / 驱逐 / `-race`；7 处迁入点既有 handler 测试套件全量通过；`rate_limit.go` 语义快测（allow 不注册 key、容量驱逐、Retry-After、trusted-proxy/`loginClientIP`）。
- **R4 证据矩阵**：判据 #1～#7 逐条映射 + 越界核账（`54fb57e7..HEAD` 红线零触碰）+ `go.mod` 无 redis。

## 11. 未选方案（除 D-001 已记录外）

| 项 | 未选 | 理由 |
|----|------|------|
| 端口级 key 形状校验（sentinel） | 未选 | 既有实现零校验；引入校验 = 行为变化 + 假 fail-closed（key 由内部构造点自产，无外部输入面） |
| 可插拔窗口策略接口 | 未选 | 防过度设计；令牌桶/漏桶实现待消费者触发再加（§3） |
| 容量配置化 | 未选 | 与 W12 D-002 现状保持冲突（D-001 未选方案） |
| `RateLimiterProvider` 带 error 返回 | 未选 | 内存供应商构造不可失败；Redis 级触发后可演化构造校验（接缝文档 §2.4 已声明触发后 fail-closed 启动校验） |

## 修订史

| date | version | change |
|------|---------|--------|
| 2026-09-01 | 0.1.0 | 初版冻结：端口形状 / key / 窗口 / 容量 / Retry-After / 红线（GOAL-002 C2） |
| 2026-09-01 | 0.1.1 | 勘误（A-002 F-006 响应）：§3 剪枝路径收窄为「剪枝仅由 `Allow` 执行；`RetryAfterSeconds` 不剪枝」；R2 义务登记（剪枝选择须显式定义空列表返回值） |