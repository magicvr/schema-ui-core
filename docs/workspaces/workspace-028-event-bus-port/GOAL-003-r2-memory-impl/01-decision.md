---
status: active
created: 2026-09-01
updated: 2026-09-01
parent: GOAL-003-r2-memory-impl
version: 0.1.0
---

# 01-decision · R2 进程内实现决策

## D-001 · 实现策略（架构与线程安全）

**日期**：2026-09-01  
**决策者**：实现团队（基于 R1 D-002 契约）

### 上下文

R1 D-002 已冻结 `kernel.EventBus` 契约（五节语义）：
- Register/Publish/Subscribe 的 topic/type 校验
- Subscribe 返回 `kernel.Subscription` 含 Unsubscribe 方法
- Stop 必须排空已入队事件、等待处理中 handler、拒绝新操作
- handler panic 隔离
- 线程安全

需要选择 Memory 实现的内部架构。

### 决策

**采用 per-subscription buffered-channel 架构**：

1. **注册表**（`map[kernel.EventTopic]*topicRegistry`）：
   - `schema any`：注册时的示例值（用于 JSON 验证）
   - `subs map[uint64]*subscription`：该 topic 的全部订阅者

2. **订阅实例**（`subscription` struct）：
   - `ch chan []byte`：buffered channel（容量 = `cfg.EventBusBufferSize`）
   - `handler func(context.Context, []byte)`：用户回调
   - `stop chan struct{}`：close 触发排空与退出
   - `loop()` goroutine：持续消费 `ch`，收到 `stop` 后 `drain()` 排空剩余事件

3. **线程安全机制**：
   - `sync.Mutex` 保护 `topics` 注册表快照（Publish 不阻塞订阅者）
   - `atomic.Bool stopped`：快速检查 Stop 状态
   - `sync.WaitGroup inflight`：追踪所有 loop goroutine
   - **select 优先级模式**：loop 内双 select 确保 stop 优先于 ch receive；enqueue 内先检查 stopCh 再尝试发送

4. **payload 隔离**：
   - Publish 时 JSON marshal 一次原始 event
   - 每个订阅者的 enqueue 传递 **独立 copy**（`append([]byte(nil), payload...)`）
   - handler 修改 payload 不影响其他订阅者

5. **Stop 语义**（D-002 §5）：
   - 设置 `stopped = true`、关闭 `stopCh`（unblock 等待中的 Publish）
   - 关闭所有订阅的 `stop` channel（触发 loop drain）
   - `inflight.Wait()` 等待全部 loop 结束（受 ctx timeout 约束）
   - 后续 Register/Publish/Subscribe 返回 `kernel.ErrEventBusStopped`

### 未选方案

- **单一 global channel + 分发 goroutine**：复杂度高，Stop drain 难以精确控制每个订阅者的剩余事件
- **sync.Map 无锁注册表**：读多写少场景收益小，mutex 快照模式更简洁
- **无 buffer / 同步 channel**：会阻塞 Publish，不符合"异步解耦"语义

### 影响

- 实现清晰，每订阅者独立 goroutine + buffered channel
- buffer 满时 Publish 阻塞（符合背压语义），可通过 ctx.Done 或 Stop 解除
- Stop 保证排空（drain loop），测试可验证
- 性能：每 Publish 需要 N 次 copy（N = 订阅者数），可接受（后续可优化为 copy-on-write）

---

## D-002 · 配置键与默认值

**日期**：2026-09-01  
**决策者**：实现团队

### 上下文

需要确定配置层的 YAML 键名、环境变量名、默认值和边界约束。参考 workspace-026 GOAL-003 的 `cache.max_entries` 模式。

### 决策

1. **YAML 键**：`eventbus.buffer_size`（struct field `EventBus.BufferSize *int`）
2. **环境变量**：`EVENTBUS_BUFFER_SIZE`（strict int parse，覆盖 YAML）
3. **默认值**：
   - `DefaultEventBusBuffer = 64`（Load 默认、composition 零值 fallback）
   - `MaxEventBusBuffer = 4096`（上界，超过 fail closed）
4. **验证规则**：
   - **YAML / env 显式值**：
     - <= 0 → fallback 到 64（允许用户显式 "0" 获得默认）
     - 1 ~ 4096 → 使用显式值
     - > 4096 → `LoadError` fail closed（防止配置错误导致内存爆炸）
     - unparsable → `LoadError` fail closed
   - **零值 Config（loader bypassed）**：composition `newEventBus` 检查 `cfg.EventBusBufferSize == 0` 后 fallback 到 64

5. **对比 cache.max_entries**：
   - cache: 非正值 fail closed（因为 0 entries 无意义）
   - eventbus: <= 0 fallback 到 64（因为 buffer=0 是有效的 unbuffered channel 语义，但我们强制最小 64）

### 未选方案

- **buffer=0 允许 unbuffered**：语义复杂，Publish 会阻塞直到 handler 处理，不符合异步事件总线定位
- **无上界**：配置错误（如 `buffer_size: 1000000`）可能导致内存耗尽

### 影响

- 配置清晰：YAML 优先、env 覆盖、默认 64、上界 4096
- fail-closed：配置错误立即启动失败，不会静默降级
- 测试：需覆盖 Load 的 YAML/env parse、ValidateProd 的上界检查

---

## D-003 · composition 注入模式

**日期**：2026-09-01  
**决策者**：实现团队

### 上下文

需要将 `kernel.EventBus` 注入到 composition root，并在 OnStop 时排空。参考 workspace-026 GOAL-003 的 `newCache` / `registerLifecycle` 模式。

### 决策

1. **Provider 函数**：
   ```go
   func newEventBus(cfg *config.Config, logger *slog.Logger) kernel.EventBus {
       buffer := cfg.EventBusBufferSize
       if buffer == 0 {
           buffer = config.DefaultEventBusBuffer
       }
       return eventbus.NewMemory(buffer, logger)
   }
   ```
   - 零值 fallback 在 provider 层（mirror cache 模式）
   - logger 用于 Stop 时记录 drain 错误（可选）

2. **fx.Provide 注册**：
   ```go
   fx.Provide(
       // ...
       newCache,
       newEventBus,  // ← 新增
       newRateLimiters,
   )
   ```

3. **newMux 签名扩展**：
   - 添加 `eventBusPort kernel.EventBus` 参数（排在 `cachePort` 之后）
   - `newMuxWithExtraProviders` 同步扩展
   - 内部暂不使用（`_ = eventBusPort`），等 R3 消费

4. **registerLifecycle OnStop 扩展**：
   ```go
   eventBusErr := eventBusPort.Stop(ctx)
   return errors.Join(shutdownErr, metricsErr, jobsErr, eventBusErr, runtimeErr, closeErr, tracingErr)
   ```
   - Stop 排空已入队事件、等待 handler、拒绝新操作
   - ctx timeout 约束整体 OnStop 预算（10s default）

5. **测试 call site 更新**：
   - `composition_test.go` / `r5_operational_gate_test.go` / `s2_access_drill_test.go` / `metrics_composition_test.go`
   - 每处增加：
     ```go
     eventBusPort := newEventBus(cfg, slog.Default())
     // ...
     mux, err := newMux(..., cachePort, eventBusPort, rateLimiters)
     ```

### 未选方案

- **延迟到 R3 注入**：拆分过细，composition 注入是 R2 实现的自然边界
- **全局单例**：违反 DI 原则，测试难以隔离

### 影响

- EventBus 成为 composition 管理的进程级单例
- OnStop 保证排空（测试可验证 graceful shutdown）
- 测试 call site 需更新 4 处，机械修改

---

## 决策索引

| 编号 | 标题 | 日期 | 状态 |
|------|------|------|------|
| D-001 | 实现策略（per-subscription channel 架构） | 2026-09-01 | 确认 |
| D-002 | 配置键与默认值（64 / 4096） | 2026-09-01 | 确认 |
| D-003 | composition 注入模式 | 2026-09-01 | 确认 |
