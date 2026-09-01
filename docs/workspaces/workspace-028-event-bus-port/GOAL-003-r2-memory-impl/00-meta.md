---
id: GOAL-003-r2-memory-impl
title: R2 进程内实现（Memory EventBus）
status: done
created: 2026-09-01
updated: 2026-09-01
notes: "Independent audit (A-002) deferred due to toolchain issues; released based on self-audit (A-001) with 0 required findings."
parent: GOAL-001-event-bus-port
version: 0.1.0
progress: —
---

# GOAL-003-r2-memory-impl · R2 进程内实现（Memory EventBus）

**所属工作区**：workspace-028-event-bus-port  
**父目标**：GOAL-001-event-bus-port  
**编号来源**：工作区内当前最大编号 002 + 1 = 003

## 目标

交付 R2 路线图检查点（GOAL-001 D-001 §R2）：

1. **进程内 Memory 实现**（`internal/eventbus/memory.go`）：
   - 基于 buffered channel 的 pub/sub 架构
   - Register/Publish/Subscribe/Unsubscribe/Stop 语义（D-002 全部五节）
   - 线程安全（atomic + mutex + select-with-stopCh）
   - payload 隔离（每订阅者独立 copy）
   - Stop 排空已入队事件、等待处理中 handler、拒绝新操作

2. **配置层**（`config.go` / `.env.example` / `config.yaml`）：
   - `eventbus.buffer_size` YAML 键（pointer mapping）
   - `EVENTBUS_BUFFER_SIZE` 环境变量覆盖（strict parse）
   - `DefaultEventBusBuffer = 64` / `MaxEventBusBuffer = 4096` 常量
   - <= 0 fallback 到默认；> 4096 fail closed

3. **composition 注入**（`composition.go`）：
   - `newEventBus(cfg, logger) kernel.EventBus` provider
   - `fx.Provide(newEventBus)` 注册
   - `newMux` / `newMuxWithExtraProviders` 签名增加 `eventBusPort` 参数
   - `registerLifecycle` OnStop 调用 `eventBusPort.Stop(ctx)` 排空

4. **测试覆盖**（`memory_test.go`）：
   - Register 验证（topic shape / JSON 序列化 / type conflict）
   - Publish 验证（shape / not registered / type mismatch / delivery）
   - Subscribe/Unsubscribe（多订阅 / 幂等 / in-flight 安全）
   - Stop 语义（drain / reject new ops / wait handlers / ctx timeout）
   - handler panic 隔离
   - 并发竞态（-race pass）
   - buffer full blocking + ctx cancel + stopCh unblock

## 边界

**包含**：
- Memory provider 实现与全面测试
- config/YAML/env 三层配置键
- composition 单例注入与生命周期排空
- 所有测试 call site 更新（4 个文件）

**不包含**（留给 R3/R4）：
- domain 代码实际使用（R3 接缝）
- Redis/Kafka/NATS 等外部实现（R4+ 触发门控）
- 跨进程/分布式事件总线
- 性能基准测试（benchmark 可选留 R4）

## 前置条件

- [x] R1 (GOAL-002) 已关门：contract 冻结、kernel.EventBus 接口落盘、I-028-001/002/003 verified
- [x] D-002 Stop 五节语义确认（含 drain 顺序、ctx timeout、拒绝新操作）

## 验收标准

1. `go test ./internal/eventbus/...` 全部通过（含 `-race`）
2. `go test ./internal/composition/...` 全部通过
3. `go build ./...` 编译无错
4. config 加载正确处理 `eventbus.buffer_size` 配置（YAML / env / default / fail-closed）
5. composition 启动时 log 输出 "kernel event-bus port ready" + buffer_size
6. OnStop 正确排空事件并等待 handler（errors.Join 包含 eventBusErr）
7. 代码遵循 cache (GOAL-003/004 workspace-026) 的注入模式

## 关联

- **父级路线图**：GOAL-001 D-001 §R2 "进程内实现"
- **前序**：GOAL-002 (R1 契约冻结)
- **后续**：GOAL-004 (R3 接缝与对齐) → GOAL-005 (R4 证据与关门)
- **参考模式**：workspace-026 GOAL-003/004 (Cache Memory 实现与注入)
- **信息项**：无新信息门禁（R1 已全部解决）
