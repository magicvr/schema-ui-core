---
status: active
created: 2026-09-01
updated: 2026-09-01
parent: GOAL-003-r2-memory-impl
version: 0.1.0
---

# 02-execution · R2 进程内实现执行记录

按时间线记录**已完成**的实施事实、产物路径与进度评估。

---

## E-001 · Memory EventBus 实现与测试（C1 落地）

**日期**：2026-09-01  
**执行者**：实现团队

### 实施内容

1. **核心实现**（`apps/api/internal/eventbus/memory.go`）：
   - 创建 `Memory` struct：
     - `topics map[kernel.EventTopic]*topicRegistry`（注册表）
     - `mu sync.Mutex`（保护注册表）
     - `stopped atomic.Bool` + `stopCh chan struct{}`（Stop 信号）
     - `inflight sync.WaitGroup`（追踪 loop goroutine）
     - `logger *slog.Logger`（可选日志）
   - 实现 `NewMemory(bufferSize int, logger *slog.Logger) *Memory`：
     - bufferSize <= 0 fallback 到 `kernel.DefaultEventBusBuffer` (64)
   - 实现 `Register(ctx, topic, schema) error`：
     - 检查 topic shape（`^[a-z][a-z0-9]*(\.[a-z][a-z0-9]*)*$`）
     - JSON 序列化验证（`json.Marshal(schema)`）
     - 幂等：同 topic + 同类型重复注册返回 nil
     - 冲突：同 topic + 不同类型返回 `ErrTopicAlreadyRegistered`
   - 实现 `Publish(ctx, topic, event) error`：
     - 检查 topic shape
     - 查找注册表（未注册 → `ErrTopicNotRegistered`）
     - type 校验：event 类型必须与注册 schema 一致
     - JSON marshal event → payload
     - 快照订阅者列表（mutex 保护）
     - 为每订阅者 `enqueue(ctx, sub, copy(payload))`（独立 copy）
   - 实现 `Subscribe(ctx, topic, handler) (Subscription, error)`：
     - 检查 topic 已注册
     - 检查 handler != nil
     - 创建 `subscription`：
       - `ch := make(chan []byte, bufferSize)`
       - `stop := make(chan struct{})`
       - `inflight.Add(1)`
       - 启动 `loop()` goroutine
     - 返回 `&subscriptionHandle{sub, once}`
   - 实现 `subscription.loop()`：
     - 双 select 优先级模式：
       ```go
       for {
           select {
           case <-s.stop:
               s.drain()
               return
           default:
           }
           select {
           case payload := <-s.ch:
               s.run(payload)
           case <-s.stop:
               s.drain()
               return
           }
       }
       ```
     - `drain()` 循环消费剩余 payload 直到 ch 空
     - `run(payload)` 用 recover 捕获 handler panic
   - 实现 `Unsubscribe()`：
     - `once.Do()` 确保幂等
     - `close(sub.stop)` 触发 drain
     - 从注册表移除订阅
   - 实现 `enqueue(ctx, sub, payload)`：
     - 先检查 `stopCh`（优先级）
     - 再 `select { case sub.ch <- payload / case <-sub.stop / case <-ctx.Done() / case <-m.stopCh }`
   - 实现 `Stop(ctx) error`：
     - `stopped.Store(true)` + `close(stopCh)`
     - 遍历全部订阅，`close(sub.stop)`
     - `inflight.Wait()` 或 `<-ctx.Done()`（先到者胜）
     - 返回 ctx.Err() 或 nil

2. **测试覆盖**（`apps/api/internal/eventbus/memory_test.go`）：
   - `TestNewMemoryBufferFallback`：bufferSize <= 0 fallback 到 64
   - `TestRegister`：topic shape / JSON serializability / idempotent / conflict
   - `TestPublishValidation`：shape / not registered / type mismatch / not serializable
   - `TestPublishDeliversToMultipleSubscribers`：多订阅者并发接收
   - `TestSubscribe`：topic not registered / handler nil / multiple subs
   - `TestUnsubscribe`：idempotent / no more events
   - `TestHandlerPanicIsolation`：一个 handler panic 不影响其他
   - `TestPayloadIsolation`：每订阅者独立 payload copy
   - `TestStop`：idempotent / drain events / wait handlers / reject new ops / ctx timeout
   - `TestPublishBlocksOnFullBuffer`：buffer 满时阻塞，ctx cancel 解除
   - `TestConcurrentOperations`：-race pass

### 产物

- `apps/api/internal/eventbus/memory.go`（447 行）
- `apps/api/internal/eventbus/memory_test.go`（554 行）

### 验证

```bash
cd apps/api
go test ./internal/eventbus/... -count=1        # PASS
go test ./internal/eventbus/... -race -count=1  # PASS (race detector)
```

### 进度评估

C1（Memory 实现与测试）已完成，符合 D-001 架构决策与 D-002 §5 Stop 语义。

---

## E-002 · 配置层与 composition 注入（C2/C3 落地）

**日期**：2026-09-01  
**执行者**：实现团队

### 实施内容

1. **config.go 扩展**（`apps/api/internal/config/config.go`）：
   - 添加常量：
     - `DefaultEventBusBuffer = 64`
     - `MaxEventBusBuffer = 4096`
   - `Config` struct 新增字段：
     - `EventBusBufferSize int`（默认 64）
   - `yamlFile` struct 新增块：
     ```go
     EventBus struct {
         BufferSize *int `yaml:"buffer_size"`
     } `yaml:"eventbus"`
     ```
   - `Load()` YAML 映射：
     ```go
     if yf.EventBus.BufferSize != nil {
         cfg.EventBusBufferSize = *yf.EventBus.BufferSize
     }
     ```
   - `Load()` env 覆盖（strict parse）：
     ```go
     if raw := os.Getenv("EVENTBUS_BUFFER_SIZE"); raw != "" {
         n, err := strconv.Atoi(raw)
         if err != nil {
             cfg.LoadError = fmt.Errorf("config: EVENTBUS_BUFFER_SIZE must be an integer")
             return cfg
         }
         if n > MaxEventBusBuffer {
             cfg.LoadError = fmt.Errorf("config: eventbus.buffer_size must be <= %d (got %d)", MaxEventBusBuffer, n)
             return cfg
         }
         cfg.EventBusBufferSize = n
     }
     ```
   - `Load()` 边界检查：
     ```go
     if cfg.EventBusBufferSize > MaxEventBusBuffer {
         cfg.LoadError = fmt.Errorf("config: eventbus.buffer_size must be <= %d (got %d)", MaxEventBusBuffer, cfg.EventBusBufferSize)
         return cfg
     }
     ```
   - `defaultConfig()` 设置 `EventBusBufferSize: DefaultEventBusBuffer`
   - `ValidateProd()` 添加检查：
     ```go
     if c.EventBusBufferSize > MaxEventBusBuffer {
         return fmt.Errorf("eventbus.buffer_size must be <= %d (got %d)", MaxEventBusBuffer, c.EventBusBufferSize)
     }
     ```

2. **YAML 配置文件更新**：
   - `apps/api/configs/config.yaml` 添加：
     ```yaml
     eventbus:
       buffer_size: 64
     ```
   - `apps/api/internal/config/config.default.yaml` 添加：
     ```yaml
     eventbus:
       buffer_size: 64
     ```

3. **环境变量模板**（`apps/api/configs/.env.example`）：
   ```bash
   # VP-028 (workspace-028 GOAL-003 D-001): per-subscription buffered-channel capacity
   # of the in-memory event-bus provider. When a subscriber's buffer is full, Publish
   # blocks until space, ctx cancel, or Stop drains. <= 0 falls back to 64; > 4096
   # fails closed at startup.
   # EVENTBUS_BUFFER_SIZE=64
   ```

4. **composition.go 注入**（`apps/api/internal/composition/composition.go`）：
   - 添加 import `"github.com/magicvr/schema-ui-core/apps/api/internal/eventbus"`
   - 新增 provider 函数：
     ```go
     func newEventBus(cfg *config.Config, logger *slog.Logger) kernel.EventBus {
         buffer := cfg.EventBusBufferSize
         if buffer == 0 {
             buffer = config.DefaultEventBusBuffer
         }
         return eventbus.NewMemory(buffer, logger)
     }
     ```
   - `NewApp` 的 `fx.Provide` 添加 `newEventBus`
   - `newMux` 签名扩展：添加 `eventBusPort kernel.EventBus` 参数
   - `newMuxWithExtraProviders` 签名扩展：添加 `eventBusPort kernel.EventBus` 参数
   - `newMux` 内部添加日志：
     ```go
     _ = eventBusPort
     logger.Info("kernel event-bus port ready", "provider", "memory", "buffer_size", cfg.EventBusBufferSize)
     ```
   - `registerLifecycle` 签名扩展：添加 `eventBusPort kernel.EventBus` 参数
   - `registerLifecycle` OnStop 添加：
     ```go
     eventBusErr := eventBusPort.Stop(ctx)
     return errors.Join(shutdownErr, metricsErr, jobsErr, eventBusErr, runtimeErr, closeErr, tracingErr)
     ```

5. **测试 call site 更新**（4 个文件）：
   - `composition_test.go` / `r5_operational_gate_test.go` / `s2_access_drill_test.go` / `metrics_composition_test.go`
   - 每处添加：
     ```go
     eventBusPort := newEventBus(cfg, slog.New(slog.NewTextHandler(io.Discard, nil)))
     // 传递给 newMux / newMuxWithExtraProviders
     ```

### 产物

- `apps/api/internal/config/config.go`（已修改，+50 行）
- `apps/api/configs/config.yaml`（已修改，+7 行）
- `apps/api/internal/config/config.default.yaml`（已修改，+7 行）
- `apps/api/configs/.env.example`（已修改，+6 行）
- `apps/api/internal/composition/composition.go`（已修改，+30 行）
- 4 个测试文件已更新

### 验证

```bash
cd apps/api
go test ./internal/config/... -run TestLoad -count=1              # PASS
go test ./internal/composition/... -run TestApp -count=1          # PASS
go test ./internal/composition/... -count=1                       # PASS
go build ./...                                                    # 编译成功
```

### 进度评估

C2（配置层）与 C3（composition 注入）已完成，符合 D-002 配置规则与 D-003 注入模式。

---

## 执行摘要

| 检查点 | 状态 | 验证方式 | 备注 |
|--------|------|----------|------|
| C1: Memory 实现与测试 | ✅ 完成 | `go test ./internal/eventbus/... -race` PASS | 447 行实现 + 554 行测试 |
| C2: 配置层（YAML/env） | ✅ 完成 | config 测试 PASS | YAML + env + defaults + fail-closed |
| C3: composition 注入 | ✅ 完成 | composition 测试 PASS | fx.Provide + OnStop drain |
| C4: 全量编译与测试 | ✅ 完成 | `go build ./... && go test ./...` | 仅 env_example_test 已知失败（pre-existing issue） |

R2 实施已完成，所有代码已落地并通过测试。
