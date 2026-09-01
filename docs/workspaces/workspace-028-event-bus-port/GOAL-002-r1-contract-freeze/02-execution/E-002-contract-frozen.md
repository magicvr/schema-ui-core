---
doc_type: goal-execution
id: E-002-contract-frozen
parent: GOAL-002-r1-contract-freeze
date: 2026-09-01
status: active
version: 0.1.0
---

# E-002 · 合同正文冻结 + 端口落地（C2）

## 事实时间线

- 2026-09-01：D-002 v0.1.0 合同冻结（§0 适用与验收基线 / §1 端口形状 / §2 注册表与 JSON 可序列化 / §3 异步投递与缓冲满阻塞 / §4 错误语义 / §5 Stop 排空 / §6 并发 / §7 错误面 / §8 红线 / §9 信息裁决 / §10 验收方式 / §11 未选方案）。
- 2026-09-01：端口本体 `apps/api/kernel/eventbus.go` 落地——`EventBus`（Register/Publish/Subscribe/Stop）+ `EventHandler` + `Subscription` + topic 形状 `ValidEventTopic` + 可执行入口 `EventMustMarshalJSON` / `ValidateEventRegister` / `ValidateEventPublish` / `ValidateEventSubscribe` + sentinel 8 个 + `DefaultEventBusBuffer = 64`。
- 2026-09-01：合同级快测 `apps/api/kernel/eventbus_test.go` 落地——编译期端口面断言（stub ×2）+ 常量断言 + `ValidEventTopic` 表驱动 17 例 + marshal/nil/chan + Register/Publish/Subscribe 顺序与 sentinel + `errors.Is` 包装。
- 2026-09-01：验证绿——`gofmt -w` / `go vet ./kernel/...` 0 / `go test ./kernel/ -count=1 -run Event` ok / `go test ./kernel/... -count=1` ok / `go build -o NUL ./kernel/` 通过。

## 产物

- `GOAL-002-r1-contract-freeze/01-decision/D-002-event-bus-port-contract.md`（责任文件 v0.1.0）
- `apps/api/kernel/eventbus.go`（端口本体）
- `apps/api/kernel/eventbus_test.go`（合同级快测）

## 下一步

- C3 审视：A-001 self + A-002 grok build（grok-4.6 · high）independent → A-003 合并响应 → R1 关门。
