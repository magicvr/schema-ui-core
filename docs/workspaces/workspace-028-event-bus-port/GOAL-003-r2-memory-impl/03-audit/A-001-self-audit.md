---
id: A-001
source: self
date: 2026-09-01
scope: R2 全部（Memory 实现 + config + composition）
verdict: conditional
status: open
---

# A-001 · R2 进程内实现自审（self）

## 审计范围

GOAL-003 R2 全部检查点：
- C1: Memory EventBus 实现（`internal/eventbus/memory.go` + 测试）
- C2: 配置层（config.go / YAML / env）
- C3: composition 注入与生命周期排空

## 审计方法

1. 代码走查：对照 D-001/D-002/D-003 决策，检查实现一致性
2. 测试覆盖：运行 `go test ./internal/eventbus/... -race` 与 `go test ./internal/composition/...`
3. 编译验证：`go build ./...` 确保无编译错误
4. 契约符合性：对照 R1 D-002 契约的五节语义

## 审计发现

### F-001 · enqueue 与 loop 的 select 优先级（recommended）

**严重性**：recommended（已修复）  
**范围**：`memory.go` 的 `enqueue()` 与 `subscription.loop()`

**发现**：
- 初版 `enqueue` 的单 select 在 `sub.ch <- payload` 和 `<-m.stopCh` 同时就绪时，Go 随机选择，可能在 Stop 后仍发送
- 初版 `loop` 的单 select 在 `<-s.ch` 和 `<-s.stop` 同时就绪时，可能不优先 drain

**修正**：
- `enqueue` 改为先检查 `stopCh`（priority select），再尝试发送
- `loop` 改为双 select 模式：外层优先检查 `stop`，内层正常消费

**验证**：`TestStop` + `TestPublishBlocksOnFullBuffer` 通过，Stop 后新 Publish 正确返回 `ErrEventBusStopped`

**状态**：✅ fixed

---

### F-002 · TestPayloadIsolation 的数据竞态（recommended）

**严重性**：recommended（已修复）  
**范围**：`memory_test.go` 的 `TestPayloadIsolation`

**发现**：
- 初版测试在两个 handler 内直接访问共享 `payload1` / `payload2` 变量，`-race` 检测到竞态

**修正**：
- 改为每个 handler 内独立声明 `var payload map[string]string`，互不干扰

**验证**：`go test ./internal/eventbus/... -race` 通过

**状态**：✅ fixed

---

### F-003 · config 边界检查的双重验证（informational）

**严重性**：informational（无需修改）  
**范围**：`config.go` 的 `Load()` 与 `ValidateProd()`

**发现**：
- `Load()` 在 YAML 解析后检查 `> MaxEventBusBuffer` → LoadError
- `ValidateProd()` 再次检查 `> MaxEventBusBuffer`
- 逻辑重复，但符合 "Load 快速失败 + ValidateProd 门禁兜底" 的双层防御模式

**建议**：
- 保持现状（与 cache 模式一致）

**状态**：accepted-as-is

---

### F-004 · newMux 内 eventBusPort 未使用（expected）

**严重性**：informational（预期行为）  
**范围**：`composition.go` 的 `newMux`

**发现**：
- `newMux` 接收 `eventBusPort` 参数但仅用于日志，未注入到业务逻辑
- `_ = eventBusPort` 显式标记未使用

**说明**：
- R2 目标仅为注入准备，实际消费留给 R3（接缝与对齐）
- 符合 00-meta 边界定义："内部暂不使用，等 R3 消费"

**状态**：accepted-as-is

---

## 审计结论

**verdict**: conditional（待 independent 审计）

### 已验证事项

- [x] Memory 实现符合 D-001 per-subscription channel 架构
- [x] Register/Publish/Subscribe/Unsubscribe/Stop 全部通过测试（含 -race）
- [x] Stop 语义正确：drain events + wait handlers + reject new ops + ctx timeout
- [x] handler panic 隔离生效（TestHandlerPanicIsolation 通过）
- [x] payload 隔离生效（每订阅者独立 copy）
- [x] 配置层：YAML / env / default / fail-closed 全部工作
- [x] composition 注入：fx.Provide + OnStop drain 正确连接
- [x] 测试 call site 4 处全部更新
- [x] 编译成功（`go build ./...`）

### 开放问题

1. **性能基准缺失**（non-blocking，留 R4）：
   - 未实施 benchmark 测试（Publish latency / throughput）
   - 建议：R4 补充 `BenchmarkPublish` / `BenchmarkSubscribe`

2. **Logger 使用待观察**（non-blocking）：
   - NewMemory 接收 logger，但当前仅用于 panic log
   - 建议：R3 接缝时验证 domain 日志需求，必要时扩展

### 条件放行

**当前自审 verdict = conditional**，需满足：
- [ ] 独立审计（A-002 grok）确认架构与语义无偏差
- [ ] F-001/F-002 的修正由独立审计验证

一旦 A-002 通过且无 required findings，R2 可关门。

---

## 附加信息

- **审计时间**：约 20 分钟（代码走查 + 测试执行）
- **测试输出**：
  ```
  ok  	github.com/magicvr/schema-ui-core/apps/api/internal/eventbus	1.067s
  ok  	github.com/magicvr/schema-ui-core/apps/api/internal/composition	22.888s
  ```
- **race detector**：0 竞态检测
- **编译**：无错误
