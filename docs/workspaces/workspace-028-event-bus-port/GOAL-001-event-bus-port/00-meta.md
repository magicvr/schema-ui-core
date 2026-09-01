---
id: GOAL-001-event-bus-port
title: 进程内事件总线端口
status: active
parent: null
created: 2026-09-01
updated: 2026-09-01
version: 0.1.0
progress: 1/4
plan_refs:
  - VP-028-event-bus-port
primary_plan: VP-028-event-bus-port
serves_summary: 进程内事件总线运输端口（架构分支 · H-002 同进程基座早期化 · 承接 RT-Q02 运输端口前置）：EventBus 端口 + 进程内 channel 实现 + outbox/MQ 接缝声明（不实现）
---

# GOAL-001 · 进程内事件总线端口

## 概述

承接 [VP-028-event-bus-port](../../vision/plans/VP-028-event-bus-port.md)（active v0.2.0 · [VRev-064](../../vision/reviews/VRev-064-vp028-event-bus-port-activation.md) self `pass` · 架构类 freshness PASS `5744868d`→`29727510`）：交付进程内事件总线**运输端口**。**对象面**：内核级 EventBus 端口（与 Cache / RateLimiter / Store / ObjectStore / Mail 同级）+ 进程内 channel 实现 + outbox/MQ 接缝声明。**红线（激活即生效）**：不预制 outbox / 外部 broker（不引入客户端依赖 / 不预裁 RT-Q06 表结构 / **不消耗 RT-Q02 trigger**）；不改 Profile 默认集 / 模块矩阵 / Manifest 装配（VP-008 `go` 消费有效性）；**不解除** Admin 功能分支 typed domain event 扩展接缝的 trigger-gated；EventBus **不是** Job 端口替代；不属 Redis 轨道；停机语义继承 VP-021（异步投递须声明 SIGTERM 取消订阅/排空，否则选同步投递）。

## 成功标准（对应 VP-028 八条方向级退出判据）

- [x] 判据 #1（端口契约冻结）：EventBus 端口（类型化 Publish/Subscribe/Unsubscribe + 订阅生命周期 + 错误语义）冻结并可用；快测可断言——R1（2026-09-01：D-002 v0.1.0 冻结 + `kernel/eventbus.go` 端口落地 + 快测绿；A-001 self + A-002 grok independent 双审 pass · 0 required）
- [ ] 判据 #2（进程内实现可用）：channel 分发 + 订阅管理 + 错误语义实现并有测试（发布/订阅/退订、并发、顺序、handler panic 隔离）
- [ ] 判据 #3（接缝声明落盘）：应用契约 vs 运输实现边界（outbox/MQ）写入；不引入 broker 客户端依赖；不实现 outbox
- [ ] 判据 #4（对齐登记）：与 roadmap Admin 功能分支 typed domain event 扩展接缝登记对齐；**不解除**其 trigger-gated
- [ ] 判据 #5（共享约定登记）：topic / 订阅命名 + 契约测试 harness 约定在架构短文或 owner VP 决策登记；**不**纳入 Redis key 轨道
- [x] 判据 #6（停机与边界语义 · V-F104）：若选异步投递须声明 SIGTERM 取消订阅/排空；否则同步投递——R1（用户裁决异步 → D-002 §5 六条 Stop 义务冻结；R2 实现挂停机路径）
- [ ] 判据 #7（边界保持）：未改 Charter；未改 Profile 默认集 / 模块矩阵 / Manifest 装配；未预制 outbox/broker；未重开历史 VP
- [ ] 判据 #8（审计闭合）：开放 required finding = 0（或已合法闭合）

## 纲领路线图（P-001）

阶段串行；同一阶段内可并行子目标。`progress` = 已完成纲领阶段 / 4。

| 阶段 | 内容 | 检查点/状态 |
|------|------|-------------|
| R1 | 契约冻结（判据 #1/#6 + I-028-001/002/003）：类型化机制（接口断言 vs 注册表，含可序列化约束取舍）· 投递语义默认（同步 vs 异步 + 缓冲满最小语义）· handler 错误语义 · 停机语义 | **已关门**（2026-09-01 · GOAL-002 `done` 3/3：用户裁决注册表+JSON / 异步+缓冲满阻塞 / 吞掉+panic 隔离 · D-002 v0.1.0 + kernel.EventBus · A-001 self + A-002 grok independent 双审 pass · 开放 required=0） |
| R2 | 进程内实现（判据 #2）：channel 分发 + 订阅管理 + 错误语义实现与测试 | 待 R1 |
| R3 | 接缝与对齐（判据 #3/#4/#5 + I-028-004）：outbox/MQ 运输接缝声明 + Admin typed domain event gated 对齐 + topic/订阅命名与契约测试 harness | 待 R2 |
| R4 | 证据与关门（判据 #7/#8；依赖 R1–R3）：证据矩阵 / 越界核账 / 审计闭合 | 待 R1–R3 |

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-028-001 | required | 事件类型化机制：Go 接口断言 vs 注册表（topic → type）。**显式取舍（V-F103）**：进程内 channel 可传非序列化负载 vs 为 outbox/RT-Q06 预留的可序列化约束——R1 须记录未选方案；若选注册表，I-028-004 升为 required。 | 方案冻结 + 退出判据 1 | R1 契约冻结 | 用户裁决（R1 契约冻结前置） | **verified** | — | 2026-09-01 用户裁决：**注册表 topic→type + JSON 可序列化**（GOAL-002 D-001 accepted；合同 §2；未选接口断言/泛型/非序列化负载） |
| I-028-002 | required | 投递语义默认：同步（发布者阻塞） vs 异步（channel 缓冲）；**缓冲满时的最小语义（V-F103）** = 阻塞 / 丢弃 / 返回错误（R1 只冻结其一；完整背压产品仍 gated）。 | 退出判据 2 | R1 | 用户裁决（R1 前置） | **verified** | — | 2026-09-01 用户裁决：**异步 + 缓冲满阻塞**（GOAL-002 D-001 accepted；合同 §3/§5 Stop 排空继承 VP-021） |
| I-028-003 | required | handler 错误语义：失败吞掉 + 日志 vs 回传发布者 vs 隔离失败（panic 恢复）；重复发布者可见性。 | 退出判据 2 | R1 | 用户裁决（R1 前置） | **verified** | — | 2026-09-01 用户裁决：**吞掉+日志 + panic 隔离**；handler 无 error 通道（GOAL-002 D-001 accepted；合同 §4） |
| I-028-004 | required | 事件类型注册权属（业务域 VP vs Admin 功能 VP）与 typed domain event gated 保持：本 VP 不解除 Admin gated（V-F101）；因 I-028-001 选注册表由 non-blocking 升 required。 | 退出判据 4 | R3 | lead 建议 + 用户确认 | 待确认 | 升 required 后最晚仍 R3 | 待确认 |

## 父目标

- `null`（Root）

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；D-001/E-001 已首条落盘，后续按编号递增。

## 备注

- **开区（2026-09-01 · 用户指令）**：VP-028 `planned → active` v0.2.0（VRev-064 self `pass` 0 required · 架构类 freshness PASS `5744868d`→`29727510` 五域零变更 · 不暂挂 `go`）；lead `workspace-028-event-bus-port`。
- 审计模式（D-001 已定）：阶段关门 default self；实证门禁（R4 证据 / 关门）按需 independent（grok build 先例，项目级默认执行路径）。
- freshness 三字段与激活锚点见 D-001：消费候选 = HEAD `29727510`；next trigger = 首个 C 端业务域 VP 激活或多实例部署评估（H-002）。
- 三端口第三个：不属 Redis 轨道；共享约定为本 VP 的 topic/订阅命名 + 契约测试 harness（R3 登记）。
