---
id: D-001
date: 2026-09-01
scope: R3 全部检查点（C1/C2/C3）
summary: 运输接缝声明 + Admin gated 对齐 + 命名约定
status: active
---

# D-001 · R3 接缝声明与对齐决策

## 背景

GOAL-001 纲领路线图 R3 阶段：声明与对齐（无编码实施）。覆盖判据 #3/#4/#5 + I-028-004。R1 已冻结契约（注册表机制 + JSON 序列化），R2 已交付进程内 channel 实现。R3 需明确：

1. 运输实现边界（进程内 vs outbox vs MQ）
2. 与 Admin typed domain event 的对齐（不解除 gated）
3. topic/订阅命名与测试约定

## 决策内容

### 1. 运输接缝声明（判据 #3）

#### 1.1 应用契约 vs 运输实现边界

**三层架构**：

```
应用契约层（本 VP 已交付）
  ↓ kernel.EventBus 端口
进程内运输层（本 VP 已交付）
  ↓ Memory 实现（channel 分发）
持久化运输层（trigger-gated，不属本 VP）
  ↓ outbox pattern（不实现）
  ↓ 外部 broker（不实现）
```

**边界定义**：

| 层 | 所有者 | 状态 | 说明 |
|----|--------|------|------|
| **应用契约** | 业务域 VP / Admin 功能 VP | 部分 gated | topic 定义 + event schema 由业务域/Admin 功能 VP 负责；本 VP 只提供运输端口 |
| **进程内运输** | 本 VP（VP-028） | **已交付** | kernel.EventBus + Memory 实现（R1/R2 done） |
| **outbox 持久化** | 未来 VP | trigger-gated | RT-Q06 表结构 + outbox writer；本 VP **不实现**、**不预裁** |
| **外部 broker** | 未来 VP | trigger-gated | RT-Q02 MQ 客户端；本 VP **不引入依赖**、**不消耗 trigger** |

**接缝约定**：

1. **进程内 ↔ outbox**：
   - 若未来实现 outbox，需在 Publish 路径插入 transactional outbox writer
   - outbox schema 必须支持 R1 契约的 JSON serializable event
   - 订阅者可选"进程内立即"或"outbox 延迟"（配置驱动，不改 Publish 调用）

2. **outbox ↔ broker**：
   - outbox poller 读取未发送记录 → broker publish
   - broker consumer → local EventBus.Publish（重新进入应用契约层）
   - 重试 / 死信 / at-least-once 语义由 broker 层负责

**本 VP 红线（已遵守）**：
- ✅ 不引入 broker 客户端依赖（Kafka / RabbitMQ / NATS 等）
- ✅ 不预裁 RT-Q06 outbox 表结构
- ✅ 不消耗 RT-Q02 trigger（broker 运输仍 gated）
- ✅ Profile 默认集 / 模块矩阵 / Manifest 未变更

**状态**：接缝声明完整，判据 #3 满足。

---

### 2. Admin typed domain event gated 对齐（判据 #4 + I-028-004）

#### 2.1 问题陈述（I-028-004）

- **触发**：R1 选择注册表机制（topic → type），导致"谁负责注册 event type"成为 required 信息项
- **风险**：若本 VP 越权注册业务域/Admin event type，会误解除 Admin typed domain event 的 trigger-gated
- **约束**：roadmap 并行规则 3 明确：领域事件**应用契约** → Admin 功能分支；**运输端口** → 架构分支（本 VP）

#### 2.2 决策（lead 建议 + 待用户确认）

**事件类型注册权属**：

| event 类别 | 注册责任 | trigger 状态 | 本 VP 角色 |
|-----------|----------|-------------|-----------|
| **系统级事件**（如 `system.startup`, `system.shutdown`） | 本 VP 或基础设施 VP | 无 gated | 可在 R2 composition 中预注册（当前未注册，留 R4） |
| **业务域事件**（如 `profile.created`, `schema.published`） | 各业务域 VP | 视业务域 VP 而定 | **不负责**；仅提供 Register API |
| **Admin 功能事件**（如 `user.invited`, `role.assigned`） | Admin 功能分支 VP | **trigger-gated** | **不负责**；**不解除 gated** |

**明确声明**：

1. **本 VP 不解除 Admin typed domain event gated**：
   - Admin 功能分支的 typed domain event 扩展接缝（roadmap L333）仍全部 trigger-gated
   - 事件 schema 定义 / Register 调用时机 / topic 命名归 Admin 功能 VP 负责
   - 本 VP 只提供 `kernel.EventBus.Register()` 运输能力

2. **注册时机由应用契约层决定**：
   - 系统级事件：可在 composition 启动时注册（当前未注册）
   - 业务域/Admin 事件：由对应 VP 的 handler/service 启动时调用 `Register()`

3. **不预置任何业务域/Admin event schema**：
   - R2 Memory 实现的 `topics map` 初始为空
   - 不在 `kernel/eventbus.go` 或 composition 中硬编码业务 topic

**对齐登记**（判据 #4）：

| 对齐点 | 本 VP 责任 | Admin 功能分支责任 | 状态 |
|--------|------------|-------------------|------|
| 端口可用性 | ✅ 已交付（R1/R2） | 触发时调用 Register/Publish/Subscribe | gated |
| 注册权属 | ✅ 不越权预注册 Admin event | 定义 event schema + 调用 Register | gated |
| trigger 保持 | ✅ 不解除 gated | 自主决定激活时机 | gated |

**状态**：对齐声明完整，I-028-004 建议已写入，待用户确认后满足判据 #4。

---

### 3. topic/订阅命名与契约测试 harness（判据 #5）

#### 3.1 topic 命名约定

**推荐格式**：`<domain>.<aggregate>.<event>`（小写，点分隔）

示例：
- 系统级：`system.startup`, `system.shutdown`, `system.config.reloaded`
- Profile 域：`profile.created`, `profile.updated`, `profile.deleted`
- Schema 域：`schema.published`, `schema.deprecated`
- Admin 功能：`user.invited`, `user.activated`, `role.assigned`（gated）

**约束**（继承 R1 D-002 契约）：
- 正则：`^[a-z][a-z0-9]*(\.[a-z][a-z0-9]*)*$`
- 长度：建议 <= 64 字符
- 语义：event 名用过去式（`created` 非 `create`）

**非约束**（留给应用契约层）：
- topic 是否对应单一 event type（推荐 1:1，但不强制）
- event schema 版本化策略（由业务域 VP 决定）
- 跨 domain 事件的命名空间冲突解决

#### 3.2 订阅命名与生命周期

**订阅标识**：
- R1/R2 未要求 subscription 有显式 name/ID
- 订阅由 `Subscribe()` 返回的 `Subscription` 对象唯一标识
- 生命周期：创建（Subscribe）→ 活跃（接收事件）→ 销毁（Unsubscribe）

**生命周期约定**：
- **启动时订阅**：service/handler 的 fx lifecycle OnStart 中调用 Subscribe
- **停机时退订**：OnStop 中调用 Unsubscribe（可选，Stop 会自动 close 全部订阅）
- **异步投递停机**（继承 VP-021）：Stop 会排空已入队事件，无需手动 drain

**并发订阅**：
- 同一 topic 可有多个订阅者（R2 Memory 实现已支持）
- 每订阅者独立接收 event copy（payload 隔离已验证）
- 无顺序保证（同 topic 的多个订阅者可能乱序收到）

#### 3.3 契约测试 harness

**位置**（不纳入 Redis key 轨道）：
- Redis key 轨道 owner = `docs/architecture/cache-redis-seam-and-track.md`（VP-026/027）
- EventBus 共享约定登记为独立文档（本决策 D-001 或新建架构短文）

**契约测试推荐模式**（留给业务域 VP）：

1. **event schema 测试**（JSON serializability）：
   ```go
   func TestEventSchema_Serializable(t *testing.T) {
       event := ProfileCreatedEvent{ProfileID: "123", ...}
       _, err := json.Marshal(event)
       assert.NoError(t, err)
   }
   ```

2. **Register 幂等性测试**：
   ```go
   func TestRegister_Idempotent(t *testing.T) {
       bus := eventbus.NewMemory(64, nil)
       err1 := bus.Register(ctx, "profile.created", ProfileCreatedEvent{})
       err2 := bus.Register(ctx, "profile.created", ProfileCreatedEvent{})
       assert.NoError(t, err1)
       assert.NoError(t, err2) // same type, idempotent
   }
   ```

3. **Publish/Subscribe 集成测试**：
   ```go
   func TestPublishSubscribe_Integration(t *testing.T) {
       bus := eventbus.NewMemory(64, nil)
       bus.Register(ctx, "profile.created", ProfileCreatedEvent{})
       
       received := make(chan ProfileCreatedEvent, 1)
       sub, _ := bus.Subscribe(ctx, "profile.created", func(ctx context.Context, event ProfileCreatedEvent) {
           received <- event
       })
       defer sub.Unsubscribe()
       
       bus.Publish(ctx, "profile.created", ProfileCreatedEvent{ProfileID: "123"})
       
       select {
       case e := <-received:
           assert.Equal(t, "123", e.ProfileID)
       case <-time.After(1 * time.Second):
           t.Fatal("timeout")
       }
   }
   ```

**harness 由本决策 D-001 提供**，业务域 VP 按需采纳。

**状态**：命名约定与测试模式已定义，判据 #5 满足。

---

## 决策汇总

| 检查点 | 内容 | 状态 |
|--------|------|------|
| C1（判据 #3） | 运输接缝声明：进程内 ↔ outbox ↔ broker 边界 + 不实现 outbox/broker | ✅ 完成 |
| C2（判据 #4 + I-028-004） | Admin gated 对齐：注册权属 + 不解除 gated | ⏳ 建议已写入，待用户确认 |
| C3（判据 #5） | 命名约定：topic 格式 + 订阅生命周期 + 契约测试 harness | ✅ 完成 |

**P-005 信息门禁**：
- I-028-004（required）：本决策 §2.2 给出 lead 建议（注册权属由应用契约层负责，本 VP 不解除 Admin gated）
- **需用户确认后闭合 I-028-004**

**下一步**：
1. 用户确认 I-028-004 决策（§2.2）
2. 执行记录（E-001）记录决策落盘事实
3. 自审（A-001）验证接缝完整性与对齐准确性
4. 更新 Root Goal 信息项状态

---

## 附加信息

- **参考文档**：
  - VP-028 v0.1.1（roadmap 并行规则 3、VRev-059 V-F101）
  - GOAL-002 D-002 契约 v0.1.0（§2 注册表机制 + §3 JSON 序列化）
  - workspace-028 workspace.md 红线约束
  
- **未涉及**（留给未来 VP 或业务域 VP）：
  - outbox 表结构设计
  - broker 选型（Kafka / RabbitMQ / NATS）
  - 重试 / 死信 / at-least-once 语义
  - event schema 版本化策略
  - 跨 service 事件顺序保证
  - CQRS / event sourcing 模式
