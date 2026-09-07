---
doc_type: architecture-decision
title: Cache 端口 Redis 接缝声明与共享轨道约定（VP-026/027 单一所有者）
status: active
created: 2026-09-01
updated: 2026-09-01
parent: null
version: 1.1.0
vision_ref: schema-ui-core-admin-foundation@0.4.0
serves: VP-026-cache-port / VP-027-rate-limiter-port
---

# Cache 端口 Redis 接缝声明与共享轨道约定

本文件是 **VP-026/027 Redis 演化轨道**的 owner 文档（VRev-059 V-F100 收窄语义：单一所有者，**不**跨区绑同一份 Goal D-001，**不**把 EventBus / VP-028 纳入 Redis key 约定——其演化轨道是 outbox/MQ）。落盘于架构层，供 VP-027 激活与任何 Redis 供应商立项时继承；变更须经 owner 决策（见 §3.5）。

## 1. 定位与所有权

| 项 | 值 |
|----|-----|
| 轨道 | Cache（VP-026）与限流（VP-027）共享的 Redis 供应商基建约定 |
| 单一所有者 | **VP-026**（lead `workspace-026-cache-port`，本波 = owner 登记义务落地） |
| 继承 | VP-027 激活时继承本文档；继承即同意，不重新裁决已冻结条款 |
| 排除 | VP-028（事件总线）不属 Redis 轨道（outbox/MQ 演化）；本文档不为其设定 key 约定 |
| 触发 | RT-Q03（多实例部署 **或** C 端业务域模块正式接入同进程）保持 trigger-gated；本文档**不是**触发消耗 |
| 端口分母 | Cache = `apps/api/kernel/cache.go`（workspace-026 GOAL-002 D-002 v0.1.1）；RateLimiter = `apps/api/kernel/ratelimit.go`（workspace-027 GOAL-002 D-002 v0.1.1 · §2.6） |

## 2. Redis 供应商接缝声明（VP-026 判据 #4）

### 2.1 端口不变（供应商边界）

- Redis 供应商实现 **同一** `kernel.Cache` 接口：`Namespace(ns) (CacheView, error)` + 视图 `Get/Set/Delete`；消费方零感知、零代码改动。
- 供应商类型只活在 `internal/`（如 `internal/cacheredis`）；任何模块与 kernel 公共面不得 import 供应商类型（VP-003 模块契约）。
- 端口是 `[]byte` 负载 → 供应商**无需**序列化注入（值编解码由消费方 Typed 层负责，见 GOAL-003）。

### 2.2 key 映射与命名空间

- Redis key = **`<ns>:<key>`**（无全局前缀；命名空间已隔离）。多应用共享同一 Redis 实例如需全局前缀，属触发后评估项（不预制）。
- ns 必须通过 `kernel.ValidCacheNamespace`（段式小写规则，R1 冻结）；key 必须通过 `kernel.ValidCacheKey`。适配器沿用 fail-closed：非法 ns/key 在触达 Redis 前拒绝（`ErrInvalidCacheNamespace` / `ErrInvalidCacheKey`），sentinel 语义与内存供应商一致。
- 命名空间登记（见 §3.3）确保 `<ns>` 段全局唯一，杜绝跨模块碰撞。

### 2.3 TTL 映射

| 策略 | Redis 表达 | 备注 |
|------|-----------|------|
| AbsoluteExpiry | SET PX（`expiresAt - now`） | 一次性 TTL |
| SlidingExpiry | SET PX + Get 命中后 `EXPIRE` 续期（一次往返） | Redis 服务端无原生滑动语义；续期粒度 = 适配器实现选择，触发立项时细化 |
| 永不过期（零值 time.Time） | SET 无 PX（持久） | 与内存供应商语义一致 |
| 惰性清理 | Redis TTL 本身即服务端清理；语义与 R1「无后台协程」兼容（惰性由服务端承担） | 适配器不再需要客户端清扫 |

### 2.4 连接管理约定

- 连接（client/连接池）由**组合根单一持有**（fx 容器，见 GOAL-004 E-002），沿 `NewApp → fx.Provide` 注入；供应商绝不自建连接。
- 启动 fail-closed：配置存在时适配器构造即 PING 校验（对齐 ObjectStore S3 HeadBucket readyz 先例），失败拒绝启动。
- 超时 / 重试 / 故障转移 / 连接池参数：**触发立项后**随实现细化（本波声明边界，不冻结具体值）。

### 2.5 依赖与测试 harness

- **不引入 Redis 客户端依赖**（判据 #4 验证面：`go.mod` 无 redis；本波及触发前始终成立）。
- 测试 harness 约定（§3.4）：同一份端口契约测试双供应商运行（内存 = 常驻；Redis = 真实实例，由触发方立项接入，继承 pgtest 的 external-harness 惯例）。

### 2.6 RateLimiter 供应商接缝声明（VP-027 判据 #4）

#### 2.6.1 端口不变（供应商边界）

- Redis 级限流供应商实现**同一** `kernel.RateLimiter` / `kernel.RateLimiterProvider` 接口（workspace-027 GOAL-002 D-002 v0.1.1）；7 处使用点消费方零感知、零代码改动。
- 供应商类型只活在 `internal/`（如 `internal/ratelimitredis`）；任何模块与 kernel 公共面不得 import 供应商类型（VP-003 薄内核）。
- 端口语义保持：`Allow` **不注册**（只读检查）、失败才 `Record`、`Clear` 清桶、`RetryAfterSeconds` 语义分母 = `kernel.RateLimiterRetryAfterSeconds`（触发立项时按远端 TTL 细化）。

#### 2.6.2 key 映射与原子窗口

- Redis key 沿用本轨道 **`<ns>:<key>`** 格式（§3.1）；限流命名空间段 = **`rl`**（本区登记，见 §3.3）。
- **原子窗口原语（VP-027 冻结）**：`INCR` 计数 + `EXPIRE`（窗口）为最小原子原语——`Record` = `INCR` + 首次 `EXPIRE`；`Allow` = 读计数（`GET`，不写、不续期失败桶）；`Clear` = `DEL`。
- **滑动窗口表达**：内存供应商为时间戳滑动窗口；Redis 级供应商的滑动实现（ZSET 时间戳 vs 固定窗口双桶近似）**触发立项时裁决**——本合同只冻结原子原语与端口语义，不预裁实现。
- 桶 key 内仍携带使用点维度（`IP|identifier` / `op|IP|user` / 纯 IP——handler key 约定，D-002 §2）。

#### 2.6.3 连接管理

- 与 §2.4 相同：连接（client/连接池）由**组合根单一持有**（fx 容器），供应商绝不自建连接；配置存在时适配器构造即 PING 校验（fail-closed）；超时 / 重试 / 故障转移参数触发立项后细化。

#### 2.6.4 测试 harness

- 与 §3.4 相同：端口契约测试双供应商运行（内存 = 常驻；Redis = 真实实例，由触发方立项接入，继承 pgtest external-harness 惯例）；`-race` 并发断言沿用内存供应商套件。

#### 2.6.5 依赖与红线

- **不引入 Redis 客户端依赖**（判据 #4 验证面：`go.mod` 无 redis；触发前始终成立）。
- 不消耗 **RT-Q05** trigger；Redis 实现仍 trigger-gated（多实例部署或 C 端业务域模块接入评估）。

## 3. 共享轨道约定（VP-026 判据 #5）

### 3.1 key 前缀

- 一律 **`<ns>:<key>`**；ns 段必须为已登记命名空间（§3.3）。

### 3.2 命名空间形状

- 形状 = `kernel.ValidCacheNamespace`（`^[a-z0-9]+(-[a-z0-9]+)*$`，≤ 64 字节）；开放集合（未来模块自由申请），非封闭枚举。

### 3.3 命名空间登记（owner 义务）

- 每个命名空间在**首个使用者的 owner 决策**中登记（谁开谁登记；缓存：VP-026 / 限流：VP-027 激活后的 owner 决策；业务域模块 = 其 VP 决策），登记内容 = ns 值 + 用途 + 归属模块。
- 冲突（同 ns 不同用途）→ 后到者 fail-closed 拒绝并回退重新命名。
- 本表随登记追加。**VP-027 R3（2026-09-01）首条登记**——履行 workspace-026 关门登记的「命名空间登记义务 → VP-027 激活触发」：

| ns | 用途 | 归属 | 登记于 |
|----|------|------|--------|
| `rl` | RateLimiter 桶（7 处使用点：登录 / 验证码生成 / 密码修改 / 自助恢复 / MFA verify / MFA step-up / 邀请接受；桶 key 内携带使用点维度 `IP\|identifier` / `op\|IP\|user` / 纯 IP） | VP-027（workspace-027 · 限流轨道） | workspace-027 GOAL-004 D-001（2026-09-01 · 短文 v1.1.0 §2.6） |

### 3.4 连接管理与测试 harness

- 连接：§2.4（组合根单一持有；触发立项后细化参数）。
- harness：单元 = 端口契约测试双供应商跑（内存 + 假 Redis/真实实例）；集成 = 真实 Redis 由触发方立项（pgtest 先例：`PG_TEST_*` env 驱动的 external harness）；`-race` 并发断言沿用内存供应商测试套件。

### 3.5 变更流程（owner 决策）

- 冻结条款（§2 全部 + §3.1/3.2/3.4）变更 → owner VP 决策（当前 VP-026 / workspace-026）留痕，`vision/` 台账同步；VP-027 激活后其 owner 决策同样有效，任一 owner 变更须在本文档修订史登记。
- 修订史：

| date | version | change |
|------|---------|--------|
| 2026-09-01 | 1.0.0 | 初版：接缝声明 + 轨道约定落盘（GOAL-004 R3 · VP-026 判据 #4/#5） |
| 2026-09-01 | 1.1.0 | **VP-027 R3 owner 决策**（workspace-027 GOAL-004 D-001 · 判据 #4/#5）：新增 §2.6 RateLimiter 供应商接缝声明（端口不变 / `rl` 命名空间段与 `<ns>:<key>` 映射 / 原子窗口 INCR+EXPIRE / 连接管理与 harness 同缓存轨道 / 无客户端依赖）；§3.3 登记表首条登记 `rl`（限流桶命名空间段 · 履行 026 登记义务）；§1 端口分母增列 `kernel/ratelimit.go`；§5 复核行 |

## 4. 触发后的专项（不属本波；触发立项时细化）

- 连接池参数 / 重试 / 故障转移 / 读侧重试语义
- 滑动过期的服务端语义选择（EXPIRE 续期粒度 / 批量续期）
- 多应用共享实例的全局前缀评估
- readyz 扩展（§2.4 已声明 PING 校验 + 触发后接 readyz）
- **限流轨道（VP-027 R3 登记 · A-002 F-002/F-003 跟踪项）**：① 容量 / FIFO 驱逐的 Redis 映射（D-002 §4 / D-001 P1 内存守卫的 Redis 表达——INCR key 默认无界，触发立项时须裁决有界化或显式接受）；② Retry-After 的远端 TTL 表达与 kernel 谓词（`RateLimiterInWindow` / `RateLimiterRetryAfterSeconds`）的位级关系——§2.6.1「触发立项时细化」的细化裁决，必要时经 §3.5 回写 D-002；③ 滑动窗口 Redis 表达（ZSET vs 双桶近似，§2.6.2）

## 5. 边界与红线复核

- 未预制 Redis 实现（无客户端依赖 / 不消耗 RT-Q03 trigger）✓（`go.mod` 无 redis）
- RateLimiter 接缝（§2.6）同此 ✓（VP-027 R3 · 2026-09-01 复核：`go.mod` 无 redis · 零代码变更 · RT-Q05 保持 trigger-gated）
- 未改端口合同 / Profile 默认集 / Manifest / Charter ✓
- mail 渠道切换语义零漂移（I-026-004：不迁移，见 workspace-026 GOAL-004 attachments）✓