---
doc_type: vision-plan
id: VP-027-rate-limiter-port
title: 通用限流器端口（内存默认 + Redis 接缝）
status: planned
vision_ref: schema-ui-core-admin-foundation@0.4.0
lead_workspace: 
created: 2026-08-31
updated: 2026-08-31
version: 0.1.1
parent: null
---

# VP-027 · 通用限流器端口

## 状态与激活门禁

| 项 | 值 |
|----|-----|
| status | **`planned`**（2026-08-31 · v0.1.1 · 用户确认立项 + VRev-059 响应修订） |
| lead_workspace | —（待激活开区） |
| Vision required | VRev-058 self（计划阶段）· VRev-059 grok build independent（复审 **conditional** → V-F099 **fixed** 2026-08-31） |
| 组合位置 | **架构分支** · H-002 同进程基座基础设施端口早期化（成功边界 #6）· RT-Q05 承接 |

## 意图

为下游 fork 的 C 端 API 提供**通用限流器端口**：默认内存（进程内）供应商开箱即用，Redis 供应商以接缝声明方式预留演化空间（不实现）。

> **解释规则（VRev-059 V-F102 → fixed）**：本波 = 基座可消费面早期化（端口 + 进程内默认 + 接缝声明），**不消耗** RT-Q05 trigger；Redis **实现**仍须等待多实例或 C 端接入评估后才立项。

设计要点：

1. **RateLimiter 端口契约**：`Allow / Record / Reset / RetryAfter` 语义，key 寻址（client IP / 用户 ID / 自定义维度），供应商无关。
2. **内存供应商（默认）**：**滑动窗口**语义 + 容量边界 + 驱逐——**演进既有 `loginRateLimiter`**（`apps/api/internal/handler/rate_limit.go`），不是平行另写一套。
3. **既有使用点迁移（完整分母 · VRev-059 V-F099 → fixed）**：代码扫描冻结 **7 处构造点** 全部接入端口，行为语义不回归（含 D-001 P1 防暴破防护）：
   | # | 构造点 | 位置 | 参数 |
   |---|--------|------|------|
   | 1 | 登录失败 | `auth.go:60` | 15min/20/64K |
   | 2 | 验证码生成 | `captcha.go:36` | 1min/10/64K |
   | 3 | 密码修改 | `account_self.go:51` | 15min/5/64K |
   | 4 | 自助恢复 | `recovery.go:58` | 15min/20/64K |
   | 5 | MFA verify 独立桶 | `mfa.go:121` | 15min/10（**不与登录桶共用**） |
   | 6 | MFA step-up | `mfa.go:129` | 15min/5（enroll/disable/recovery-rotate） |
   | 7 | 邀请接受 | `invites.go:308`（W13 F-001 预认证面 CPU DoS 刹车） | 15min/10 |
   
   **显式排除**：GOAL-014 账号分层锁定（DB 行锁，**不是** `loginRateLimiter`，不纳入端口）。
4. **Redis 供应商接缝声明**：端口不变、供应商边界、原子窗口（INCR + EXPIRE）语义约定落盘，实现留待 RT-Q05 触发（多实例或 C 端接入评估）。
5. **共享供应商约定**：Redis 轨道（VP-026/027）统一 key 前缀 / 命名空间 / 连接管理 / 测试 harness 约定，登记于架构短文或 owner VP 决策（单一所有者，不跨区绑同一份 Goal D-001——VRev-059 V-F100）；VP-028 不属 Redis 轨道。
6. **相邻合同继承（VRev-059 V-F104 → fixed）**：VP-009 W12 D-002——官方单实例边界不变、Redis 仅预登记、login/recovery 窗口常量（15min/20/`IP|identifier`/`Retry-After`）**保持现状**；端口化不重写语义。

本 VP 属 **架构分支**，承接 Charter 0.4.0 成功边界 #6 与 H-002；**不预制 Redis 实现**；**与缓存（VP-026）/ 事件（VP-028）独立交付**。**不改 Charter**。

## 首波冻结（退出分母 = 限流器端口操作化）

| 项 | 本 VP 交付 | 不进本 VP |
|----|-----------|-----------|
| 端口契约 | RateLimiter 端口（Allow/Record/Reset/RetryAfter + key 寻址 + 供应商无关） | 分布式限流协调（多实例触发后评估）；按用户/路由配额的业务策略（业务域） |
| 内存供应商 | 滑动窗口 + 容量边界 + 驱逐（演进 `loginRateLimiter`） | 令牌桶 / 漏桶 / 固定窗口（策略接口可扩展，实现待消费者触发） |
| 迁移 | 7 处构造点全部接入端口（登录 / captcha 生成 / 密码修改 / 恢复 / MFA verify / MFA step-up / 邀请接受）；行为不回归（D-001 P1 + W12 D-002 窗口常量保持）；GOAL-014 分层锁定显式排除 | 非限流语义的进程内状态（如 scheduler 去重 map——那是调度状态，不是限流） |
| 接缝 | Redis 供应商接缝声明（端口不变 / INCR+EXPIRE 语义 / 连接管理约定） | Redis 客户端依赖引入；跨实例一致性协议（触发后随实现） |
| 共享约定 | key 前缀 / 命名空间 / 测试 harness 约定 D-001 登记 | 跨端口交付物合并（VP-026/028 独立关门） |

## 非目标

- **Redis 供应商实现**（RT-Q05 触发条件仍为 trigger-gated；本波只落接缝声明）
- **分布式限流 / 跨实例协调**（多实例部署或 C 端跨实例需求才评估）
- **缓存语义**（归 VP-026）；**消息/事件语义**（归 VP-028）
- **业务级配额策略**（按用户/组织/路由的配额语义归业务域 VP）
- **调度状态去重 / 幂等守卫**（那是领域状态，不是限流器）
- 重开 VP-012 / 已 closed 记录；替代 VP-009 / VP-010；改变 Charter 边界

## 与相邻 VP 的边界

| VP / 分支 | 关系 |
|-----------|------|
| **VP-003 / VP-004** | 遵守薄内核。限流端口是内核级基础设施端口；handler/模块公共面不得依赖供应商类型 |
| **VP-008 `go`** | 架构类能力；激活前做架构类 freshness |
| **VP-009 / VP-010** | 限流绕过 / 资源耗尽类安全 gap 与符合性 gap 归持续程序。**W12 D-002（V-F104）**：官方单实例边界不变、Redis 仅预登记不实施、login/recovery 窗口常量（15min/20/`IP|identifier`/`Retry-After`）保持现状；本波端口化不重写语义 |
| **VP-021 停机合同** | 若 R1/R2 选择后台清理协程，须声明 SIGTERM 下停止清理/排空；否则选惰性清理避开新生命周期（V-F104） |
| **VP-026** | 同为 key 寻址状态服务、同 Redis 演化轨道，但端口/交付/关门完全独立；共享供应商基建**约定**（D-001） |
| **VP-028** | 无依赖；仅共享"三端口供应商约定"注记 |
| **架构 RT-Q05** | 本 VP 为 RT-Q05 的承接 VP（planned）；Redis 实现仍 trigger-gated |
| **业务域** | C 端业务域模块激活后可消费限流端口；届时按成功边界 #6 评估是否需要 Redis 供应商 |

## 方向级退出判据

1. **端口契约冻结**：RateLimiter 端口（Allow/Record/Reset/RetryAfter + key 寻址 + 供应商无关）冻结并可用；快测可断言。
2. **内存供应商可用**：滑动窗口 + 容量边界 + 驱逐语义实现并有测试（并发、窗口边界、驱逐、RetryAfter 计算）。
3. **既有使用点迁移不回归（完整分母 · V-F099）**：7 处构造点（登录 / 验证码生成 / 密码修改 / 自助恢复 / MFA verify 独立桶 / MFA step-up / 邀请接受）全部接入端口；回归证据形态 = 各迁入点既有 handler 测试套件全量通过 + `rate_limit.go` 单元语义（allow 不注册 key、容量驱逐、Retry-After、trusted-proxy/`loginClientIP`）+ W12 D-002 窗口常量保持；GOAL-014 分层锁定显式排除。
4. **Redis 接缝声明落盘**：供应商边界（端口不变）、原子窗口语义（INCR + EXPIRE）、连接管理约定写入；不引入 Redis 客户端依赖。
5. **共享约定登记**：Redis 轨道约定（VP-026/027：key 前缀 / 命名空间 / 连接管理 / 测试 harness）在架构短文或 owner VP 决策登记（单一所有者；不跨区绑同一份 Goal D-001）；VP-028 不属 Redis 轨道。
6. **边界保持**：未改 Charter；未改 Profile 默认集 / 模块矩阵 / Manifest 装配；未预制 Redis 实现；未重开历史 VP。
7. **审计闭合**：开放 required finding = 0（或已合法闭合）。

详细纲领阶段由 lead Root（P-001）书写：R1 契约冻结（API 形态 / 窗口语义 / 迁移边界）→ R2 内存供应商与迁移 → R3 接缝与约定 → R4 证据与关门。本 VP 不写 Goal 五件套。

## 信息需求（P-005）

| id | 要回答的问题 | 级别 | 影响门禁 | 最晚阶段 | 状态 |
|----|--------------|------|----------|----------|------|
| I-027-001 | 端口 API 形态：Allow/Record 拆分 vs 内聚 Allow（内部记数）；RetryAfter 语义（现有 `retryAfterSeconds` 演进）。 | required | 方案冻结 + 退出判据 1 | R1 契约冻结 | 待裁决 |
| I-027-002 | 既有 `loginRateLimiter` 迁移策略：演进为内存供应商（推荐） vs 保留并存（双轨）；多实例 limiter 实例的 key 维度是否扩展。 | required | 退出判据 3 | R2 | 待裁决 |
| I-027-003 | 窗口语义默认：滑动窗口（现状保持） vs 固定窗口 vs 混合；策略接口是否与缓存 VP-026 共用形态。 | non-blocking | 退出判据 2 | R1 | 待确认 |
| I-027-004 | 限流 key 维度扩展：是否新增"路由+用户"复合 key（C 端 API 防刷典型）；或留给业务域 VP 自行定义维度。 | non-blocking | 退出判据 1 | R1 | 待确认 |

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| — | — | lead | — | `planned` 0 区；激活时指定 |

## 关门记录

（仅 `closed` / `abandoned` 时填写。）

| date | outcome | summary | evidence_links | residuals |
|------|---------|---------|----------------|-----------|
| — | — | — | — | — |

## 规划修订短史

| date | change |
|------|--------|
| 2026-08-31 | 初创 `planned`：用户裁决按 3 个独立 VP 执行（缓存 / 限流 / 事件总线；触发条件独立 × 关门能力独立原则）。本 VP 承接 RT-Q05（限流器端口 · 内存默认 + Redis 接缝声明），迁移既有 loginRateLimiter；vision_ref @0.4.0；roadmap / revisions 原子同步 |
| 2026-08-31 | v0.1.1 · **VRev-059 响应修订**（grok build · conditional → 闭合）：V-F099 required **fixed**——使用点分母补全为代码扫描的 7 处构造点（登录/captcha/密码修改/恢复/MFA verify 独立桶/MFA step-up/邀请接受），意图/冻结表/退出判据 3 三处对齐，显式排除 GOAL-014 分层锁定；V-F100 **fixed**——Redis 轨道约定改为单一所有者（架构短文或 owner VP，不跨区绑 D-001），VP-028 不属 Redis 轨道；V-F102 **fixed**——补"不消耗 RT-Q05 trigger"解释规则；V-F104 **fixed**——继承 W12 D-002（窗口常量保持）与 VP-021 停机语义声明 |