---
doc_type: goal-audit
id: A-001-contract-freeze-closeout-self
parent: GOAL-002-r1-contract-freeze
date: 2026-09-01
source: self
scope: GOAL-002 R1 契约冻结全量（D-001 用户裁决核验 / D-002 合同 ↔ kernel/eventbus.go 逐节一致性 / 快测覆盖 / 越界核账 / 信息门禁）
verdict: pass
open_required: 0
status: active
version: 0.1.0
---

# A-001 · R1 契约冻结关门自审（self）

## 1. 信息门禁（P-005）

| ID | 级别 | 状态 | 证据 |
|----|------|------|------|
| I-028-001 | required | **verified** | 2026-09-01 用户裁决「注册表 topic→type + JSON 可序列化」（D-001 accepted；合同 §2：Register 冻类型、Publish 精确匹配、EventMustMarshalJSON fail-closed） |
| I-028-002 | required | **verified** | 用户裁决「异步 + 缓冲满阻塞」（合同 §3：Publish 不等待 handler；每订阅有界缓冲默认 64；满则阻塞；§5 Stop 排空继承 VP-021） |
| I-028-003 | required | **verified** | 用户裁决「吞掉+日志 + panic 隔离」（合同 §4：EventHandler 无 error 返回；供应商 MUST recover panic） |
| I-028-004 | required | 待确认（R3 前置） | 因选注册表由 non-blocking 升 required；最晚阶段 R3——不影响 R1 关门 |

## 2. 合同 ↔ 实现逐节一致性

- **§1 端口形状**：`kernel.EventBus`（Register/Publish/Subscribe/Stop）+ `EventHandler` + `Subscription` 与合同签名一致；handler 收 `[]byte`；无 `EventBusProvider` 工厂。
- **§2 注册表**：`ValidEventTopic` 正则为 `^[a-z0-9]+(\.[a-z0-9]+)*$`、≤128；`ValidateEventRegister` / `ValidateEventPublish` / `ValidateEventSubscribe` 顺序与合同一致；类型冲突检测声明为 stateful、留供应商（R2）。
- **§3 异步/缓冲**：`DefaultEventBusBuffer = 64` 常量落地；缓冲行为本身属 R2 实现，本合同只冻义务。
- **§4 错误语义**：handler 签名无 error；panic 义务写在 EventBus 接口注释。
- **§5 停机**：Stop 在端口面上，与用户选异步后的 V-F104 义务对齐。
- **§7 错误面**：8 个 sentinel 均已声明。
- **§8 红线**：未引入 broker 依赖（go.mod 零变更）；未改 Profile / Manifest；未实现 outbox。

## 3. 快测覆盖评估

`kernel/eventbus_test.go`：编译期端口面断言（stub ×2）+ 常量断言 + ValidEventTopic 17 例（含 hyphen 拒绝，与 Cache 中划线命名空间刻意区分）+ marshal 结构体/nil/chan + Register/Publish/Subscribe 顺序与 sentinel + 类型不匹配（指针 / 匿名结构）+ 同类型不可序列化 + `errors.Is` 包装。覆盖合同 §2/§7/§10 的全部可执行谓词。

## 4. 越界核账

本波代码变更面 = `apps/api/kernel/eventbus.go`、`apps/api/kernel/eventbus_test.go`；文档变更面 = `docs/workspaces/workspace-028-event-bus-port/**` + VP-028 信息表回写。`go.mod` / `go.sum` / Profile 装配 / Charter / handler 零触碰。

## 5. 验证复跑（2026-09-01）

`gofmt -w` · `go vet ./kernel/...` 0 · `go test ./kernel/ -count=1 -run Event` ok · `go test ./kernel/... -count=1` ok · `go build -o NUL ./kernel/` 通过。

## Verdict

**pass**（0 required）。R1 契约冻结满足关门条件；建议 A-002（grok build · grok-4.6 · high）independent 复核后合并响应关门。

## Findings

- required：无。
- recommended：无（I-028-004 已登记 R3 前置，不属本目标）。
