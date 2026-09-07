---
doc_type: goal-decision
id: D-002-event-bus-port-contract
parent: GOAL-002-r1-contract-freeze
date: 2026-09-01
status: accepted
version: 0.1.0
---

# D-002 · EventBus 端口合同 v0.1.0（2026-09-01 冻结）

> **责任文件（frozen）**。实施（R2 进程内 channel 实现）与验收（R4 证据矩阵）以本合同为分母。本波只落端口本体 + 合同级快测；channel 供应商、缓冲配置键、日志注入归 R2；outbox/MQ 接缝、topic 命名与契约测试 harness、I-028-004 注册权属归 R3。不改 Profile 默认集 / 模块矩阵 / Manifest 装配；不引入 broker 客户端；不实现 outbox；不改 Charter；不解除 Admin typed domain event gated。

## 0. 适用与验收基线

- **契约面**：`apps/api/kernel` 公共面（Go 1.26）；模块经 `kernel.EventBus` 消费，绝不接触供应商类型（VP-003 薄内核）。
- **先例对齐**：`kernel.Cache` / `RateLimiter` / `Store` / `ObjectStore` / `MailSender` —— 非泛型接口 · ctx 首位 · fail-closed 校验 · sentinel errors（`errors.Is`）。
- **定位**：运输端口 + 进程内实现（R2）+ outbox/MQ 接缝声明（R3）；**不是**应用契约前置；**不是** Job 端口替代（持久化 / 重试 / 定时工作仍走 Job）。
- **范围外（trigger-gated / 属其它阶段）**：outbox / 外部 broker 实现（RT-Q02 仍 gated，本波不消耗 trigger）；重试 / 死信 / 持久化 / 完整背压产品；事件溯源 / CQRS；跨 topic 顺序保证；业务域事件产品语义；Admin typed domain event 应用契约（仍 gated）。

## 1. 端口形状（冻结）

```go
// kernel.EventBus：供应商无关事件运输端口（进程内默认；D-001）。
type EventBus interface {
    // Register 冻结 topic → 样例类型。同 topic 同类型幂等成功；同 topic 不同类型 fail-closed。
    Register(ctx context.Context, topic EventTopic, sample any) error
    // Publish 校验注册类型 + JSON 可序列化后异步投递；不等待 handler 完成。
    // 任一订阅缓冲满时阻塞，直至有空位、ctx 取消或 Stop。
    Publish(ctx context.Context, topic EventTopic, event any) error
    // Subscribe 在已注册 topic 上追加一个订阅；返回可退订的 Subscription。
    Subscribe(ctx context.Context, topic EventTopic, handler EventHandler) (Subscription, error)
    // Stop 拒绝新的 Register/Publish/Subscribe；排空已入缓冲事件并等待在飞 handler；取消全部订阅。
    Stop(ctx context.Context) error
}

// kernel.EventHandler：消费回调。无 error 返回——失败吞掉为结构（I-028-003）。
// payload 是 Publish 时 json.Marshal 的副本，handler 可自由持有/变更。
type EventHandler func(ctx context.Context, payload []byte)

// kernel.Subscription：单条订阅的生命周期句柄。Unsubscribe 幂等。
type Subscription interface {
    Unsubscribe()
}
```

- **单例端口**：进程内一条总线（对齐 Cache / MailSender，不引入 `EventBusProvider` 工厂）。R2 在 composition 注入唯一实现。
- **ctx 首位**：Register / Publish / Subscribe / Stop 均以 `context.Context` 为首参；Publish 在缓冲满阻塞期间尊重取消；Stop 尊重截止（对齐 VP-021 `shutdown_timeout`）。
- **交付负载**：handler 收 JSON `[]byte`（与 Cache `[]byte` 同构；序列化约束在 Publish 侧 fail-closed 兑现；隔离 handler 对发布者对象的变异）。

## 2. topic 与注册表（I-028-001 · 冻结：注册表 + JSON 可序列化）

- `EventTopic` = 点分层小写段：`^[a-z0-9]+(\.[a-z0-9]+)*$`，1～128 字节。`ValidEventTopic` 为唯一形状入口；非法 → `ErrInvalidEventTopic`（fail-closed，不回落）。
- **Register**：`sample` 非 nil；`json.Marshal(sample)` 必须成功；记录 `reflect.TypeOf(sample)`。
  - 同 topic、同类型：幂等成功。
  - 同 topic、不同类型：`ErrEventTypeConflict`。
  - 未注册 topic 上 Publish / Subscribe：`ErrEventTopicNotRegistered`。
- **Publish 类型**：`reflect.TypeOf(event)` 必须**精确等于**已注册类型（不自动解引用指针、不接受实现同一接口的其它具体类型）。
- **可执行语义权威（供应商必须使用）**：
  - `kernel.ValidEventTopic(topic)`
  - `kernel.EventMustMarshalJSON(v)`（nil → `ErrInvalidEventPayload`；marshal 失败 → `ErrEventNotSerializable` 包装原错误）
  - `kernel.ValidateEventRegister(topic, sample)`
  - `kernel.ValidateEventPublish(registered reflect.Type, event any)`（顺序：registered nil → `ErrEventTopicNotRegistered`；event nil → `ErrInvalidEventPayload`；类型不等 → `ErrEventTypeMismatch`；然后 EventMustMarshalJSON）
  - `kernel.ValidateEventSubscribe(topic, handler)`（topic 形状 + handler 非 nil）
- **未选**：接口断言无注册表；泛型端口；gob / protobuf；进程内传递非序列化负载；以 JSON roundtrip 再反序列化为 any 交付（交付保持 `[]byte`，反序列化归调用方）。

## 3. 异步投递与缓冲（I-028-002 · 冻结：异步 + 缓冲满阻塞）

- **异步**：`Publish` 在事件进入**每一个**当前订阅的缓冲后返回；**不等待** handler 执行完毕。
- **每订阅独立有界缓冲**：容量来源 = 配置键 `eventbus.buffer_size`（R2 落 `internal/config` + YAML/env，默认 **`DefaultEventBusBuffer = 64`**；`<= 0` 回落该默认；非法过大由 R2 fail-closed）。
- **缓冲满最小语义 = 阻塞**：对该订阅 `send` 阻塞，直至有空位、`ctx` 取消（返回 `ctx.Err()`）或总线已 Stop（返回 `ErrEventBusStopped`）。不丢弃、不立即返错（完整背压产品仍 gated）。
- **顺序保证范围**：同一订阅内，来自同一 `Publish` 调用序列的事件保序；跨订阅、跨 topic 不保证。
- **多订阅 / 重复发布**：同一 topic 允许多个订阅；同一事件允许重复 Publish。各订阅独立缓冲、独立失败。

## 4. 错误语义（I-028-003 · 冻结：吞掉 + panic 隔离）

- handler **无 error 返回通道**——消费失败不能回传发布者（吞掉为结构）。
- 供应商 **MUST** `recover` handler panic，记日志，继续服务其它事件 / 其它订阅。单个 handler panic 不得打死进程，也不得取消其它订阅。
- 日志实现归 R2（注入既有可观测端口 / 标准库 logger）；本合同只冻结义务，不冻结日志格式。
- `Unsubscribe` 之后该订阅不再收到后续事件；在飞 handler 允许完成。`Unsubscribe` 幂等。

## 5. 停机（判据 #6 · V-F104 · 继承 VP-021）

因选择异步投递，**必须**声明 SIGTERM 取消订阅 / 排空（不能再靠同步投递躲开新生命周期）：

1. `Stop(ctx)` 之后：Register / Publish / Subscribe 一律 `ErrEventBusStopped`。
2. 已阻塞在缓冲满上的 Publish 被唤醒并返回 `ErrEventBusStopped`。
3. 已入缓冲、尚未交给 handler 的事件：在 `ctx` 截止前尽量排空（逐条交给 handler）。
4. 在飞 handler：等待返回或 `ctx` 结束（对齐 VP-021 `shutdown_timeout`；超时后 Stop 返回 `ctx.Err()`，不保证 handler 已结束）。
5. 排空结束后取消全部订阅（等价对每个 Subscription 调 Unsubscribe）。
6. `Stop` 幂等。

R2 实现必须挂到进程停机路径（composition / server shutdown），不得另起无停机声明的后台协程。

## 6. 并发安全

- 所有接口方法必须并发安全（多 goroutine 并行 Register/Publish/Subscribe/Unsubscribe/Stop 无数据竞争）。
- R2 供应商测试以 `-race` 覆盖并发边界、缓冲满阻塞与 Stop 竞态。

## 7. 错误面

| sentinel | errors.Is | 触发 |
|----------|-----------|------|
| `kernel.ErrInvalidEventTopic` | ✓ | topic 形状非法 |
| `kernel.ErrInvalidEventPayload` | ✓ | Register/Publish 的 sample/event 为 nil |
| `kernel.ErrEventNotSerializable` | ✓ | `json.Marshal` 失败 |
| `kernel.ErrEventTypeMismatch` | ✓ | Publish 类型 ≠ 已注册类型 |
| `kernel.ErrEventTypeConflict` | ✓ | Register 同 topic 不同类型 |
| `kernel.ErrEventTopicNotRegistered` | ✓ | Publish/Subscribe 时 topic 未注册；ValidateEventPublish 的 registered==nil |
| `kernel.ErrEventHandlerNil` | ✓ | Subscribe 传入 nil handler |
| `kernel.ErrEventBusStopped` | ✓ | Stop 之后的 Register/Publish/Subscribe，或阻塞中的 Publish 被 Stop 唤醒 |

## 8. 红线

- 不预制 outbox / broker（不引入客户端依赖 / 不预裁 RT-Q06 表结构 / **不消耗 RT-Q02 trigger**）。
- 不改 Profile 默认集 / 模块矩阵 / Manifest 装配（VP-008 `go`）。
- 不解除 Admin 功能分支 typed domain event 扩展接缝的 trigger-gated（应用契约仍归 Admin 功能；I-028-004 R3 裁定注册权属）。
- EventBus **不是** Job 端口替代。
- 不属 Redis 轨道（owner 仍为 `cache-redis-seam-and-track.md`）。
- 不重开 VP-012 / 已 closed 记录；不改 Charter。

## 9. 信息裁决记录

| ID | 裁决 | 证据 |
|----|------|------|
| I-028-001 | 注册表 + JSON 可序列化（§2） | D-001（2026-09-01 用户裁决） |
| I-028-002 | 异步 + 缓冲满阻塞（§3）+ Stop 排空（§5） | D-001（2026-09-01 用户裁决） |
| I-028-003 | 吞掉+日志 + panic 隔离（§4） | D-001（2026-09-01 用户裁决） |
| I-028-004 | 注册权属 | **R3 前置裁决**（本目标不关闭；已升 required） |

## 10. 验收方式（R2/R4 预告）

- **合同级快测（本目标，C2）**：`kernel/eventbus_test.go` —— `ValidEventTopic` 正反例表驱动；`EventMustMarshalJSON`（结构体 / nil / 不可序列化 chan）；`ValidateEventRegister` / `ValidateEventPublish` / `ValidateEventSubscribe` 顺序与 sentinel；`DefaultEventBusBuffer` 常量；编译期端口面断言（stub 实现 `EventBus` / `Subscription`）；`%w` 包装后 `errors.Is`。
- **R2 供应商测试**：compile-time 断言内存实现 `kernel.EventBus`；发布/订阅/退订；异步（Publish 在 handler 阻塞时仍返回——缓冲未满）；缓冲满阻塞；Stop 排空与拒绝新 Publish；panic 隔离；类型冲突；`-race`。
- **R4 证据矩阵**：判据 #1～#8 逐条映射 + 越界核账 + `go.mod` 无 broker 客户端。

## 11. 未选方案（除 D-001 已记录外）

| 项 | 未选 | 理由 |
|----|------|------|
| `EventBusProvider` 工厂 | 未选 | 进程内单例足够；与 Cache/MailSender 同构 |
| handler 收 `any`（原对象） | 未选 | 变异竞态 + 可序列化约束无法在交付侧兑现 |
| 端口级全局有序（跨订阅） | 未选 | 成本高；本合同只保证单订阅内保序 |
| 端口自带 Logger 接口 | 未选 | 避免 kernel 依赖日志实现；R2 注入 |
| 同步投递以躲开 Stop | 未选 | 用户已选异步（D-001） |

## 修订史

| date | version | change |
|------|---------|--------|
| 2026-09-01 | 0.1.0 | 初版冻结：端口形状 / 注册表 / 异步缓冲 / 错误语义 / 停机 / 红线（GOAL-002 C2） |
