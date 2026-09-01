---
doc_type: vision-plan
id: VP-028-event-bus-port
title: 进程内事件总线端口（outbox/MQ 接缝）
status: active
vision_ref: schema-ui-core-admin-foundation@0.4.0
lead_workspace: workspace-028-event-bus-port
created: 2026-08-31
updated: 2026-09-01
version: 0.2.0
parent: null
---

# VP-028 · 进程内事件总线端口

## 状态与激活门禁

| 项 | 值 |
|----|-----|
| status | **`active`**（2026-09-01 · v0.2.0 · 用户指令激活 + VRev-064 self `pass`） |
| lead_workspace | `workspace-028-event-bus-port`（2026-09-01 开区） |
| Vision required | VRev-058 self（计划阶段）· VRev-059 grok build independent（复审 conditional → VP-028 无 required；V-F101/102/103/104 响应 **fixed**）· **VRev-064 self（激活就绪 `pass` · 0 required · 2026-09-01）** |
| 组合位置 | **架构分支** · H-002 同进程基座基础设施端口早期化（成功边界 #6）· RT-Q02 承接（**运输端口**前置；broker 仍 gated） |

## 意图

为同进程内的业务域模块提供**事件总线运输端口**：进程内（channel）实现开箱即用，outbox / 外部 MQ 运输以实现声明方式预留演化空间（不实现）。

> **定位澄清（VRev-059 V-F101 → fixed）**：本 VP = **运输端口 + 进程内实现 + outbox/MQ 接缝声明**，**不**是"应用契约前置"——roadmap 并行规则 3 明确领域事件**应用契约**归 Admin 功能分支；本 VP **不解除** Admin typed domain event 扩展接缝的 trigger-gated。EventBus **不是** VP-012 Job 端口的替代：持久化 / 重试 / 定时工作仍走 Job。
>
> **解释规则（VRev-059 V-F102 → fixed）**：本波 = 基座可消费面早期化（端口 + 进程内默认 + 接缝声明），**不消耗** RT-Q02 trigger；outbox / broker **实现**仍须等待多实例或领域事件 fan-out 评估后才立项。

设计要点：

1. **EventBus 端口契约**：类型化事件 `Publish / Subscribe / Unsubscribe`，订阅生命周期管理（退订防泄漏），handler 错误语义明确（失败吞掉 + 日志 vs 回传——R1 冻结，I-028-003）。
2. **进程内实现（默认）**：channel 分发（同步/异步语义冻结），处理器并发模型（串行 per-topic vs 并发），顺序保证范围（同一发布者的顺序）。
3. **outbox / MQ 接缝声明**：运输端口与外部实现（outbox / broker）的边界落盘；**不实现** outbox、不引入 broker 客户端、不预裁 RT-Q06 outbox 表结构。
4. **边界声明（V-F104）**：若选异步 channel 投递，须声明 SIGTERM 下取消订阅 / 排空（继承 VP-021 停机合同）；否则选同步投递避开新生命周期。
5. **共享约定（V-F100）**：topic / 订阅命名 + 契约测试 harness 约定登记于架构短文或 owner VP 决策；**不**纳入 Redis key 轨道（VP-026/027 专属）。

本 VP 属 **架构分支**，承接 Charter 0.4.0 成功边界 #6 与 H-002；**不预制 outbox / broker**；**与缓存（VP-026）/ 限流（VP-027）独立交付**。**不改 Charter**。

## 首波冻结（退出分母 = 事件总线端口操作化）

| 项 | 本 VP 交付 | 不进本 VP |
|----|-----------|-----------|
| 端口契约 | EventBus（类型化 Publish/Subscribe/Unsubscribe + 订阅生命周期 + 错误语义） | outbox / 外部 broker（RT-Q02 仍 gated）；事件溯源；跨进程投递保证 |
| 进程内实现 | channel 分发 + 订阅管理 + handler 错误语义（冻结） | 持久化 / 重试 / 死信 / 背压（触发后评估）；顺序跨 topic 保证 |
| 接缝 | outbox / MQ 运输接缝声明（应用契约 vs 运输实现边界；不实现） | broker 客户端依赖引入；outbox 表设计（触发后随 RT-Q02 评估） |
| 对齐 | 与 roadmap Admin 功能分支 typed domain event 扩展接缝登记对齐（**不解除**其 trigger-gated；事件类型注册权属由业务域/Admin 功能 VP 定义，I-028-004） | 通知 / SSO / 审批等其它扩展接缝（各自 trigger-gated） |
| 共享约定 | topic/订阅命名 + 契约测试 harness（V-F100；**不**纳入 Redis key 轨道） | 跨端口交付物合并 |

## 非目标

- **outbox / 外部消息队列实现**（RT-Q02 触发条件仍为 trigger-gated：多实例、跨机长任务、领域事件 fan-out）
- **重试 / 死信 / 持久化 / 背压**（消费者触发后评估；本波只冻结错误语义与并发模型）
- **事件溯源 / CQRS 事件流**（Charter 非目标登记）
- **缓存语义**（归 VP-026）；**限流语义**（归 VP-027）
- **业务域事件的产品语义**（订单已创建、支付成功等事件定义归业务域 VP；本波只交付运输端口）
- 重开 VP-012 / 已 closed 记录；替代 VP-009 / VP-010；改变 Charter 边界

## 与相邻 VP 的边界

| VP / 分支 | 关系 |
|-----------|------|
| **VP-003 / VP-004** | 遵守薄内核。EventBus 端口是内核级基础设施端口；模块公共面不得依赖 channel 实现细节 |
| **VP-008 `go`** | 架构类能力；激活前做架构类 freshness |
| **VP-009 / VP-010** | 事件处理安全 gap（循环发布 / handler 死循环 / 内存耗尽）与符合性 gap 归持续程序 |
| **VP-012** | 已 closed 的横切契约（correlation / 审计 / Job）不重开；**EventBus 不是 Job 端口的替代**——持久化 / 重试 / 定时工作仍走 Job（V-F101）；本 VP 是新的进程内事件**运输端口** |
| **VP-026 / VP-027** | 无依赖；**不**共享 Redis key 约定（演化轨道不同：outbox/MQ）；仅共享"端口早期化不消耗 trigger"的解释规则 |
| **架构 RT-Q02** | 本 VP 为 RT-Q02 的**运输端口**前置（进程内 + 接缝声明）；outbox/broker 运输仍 gated（本波不消耗 trigger） |
| **Admin 功能分支** | typed domain event **应用契约**（扩展接缝 trigger-gated）归 Admin 功能分支（并行规则 3）；本 VP 只交运输端口，**不解除**该 gated（V-F101） |
| **VP-021 停机合同** | 若选异步 channel 投递须声明 SIGTERM 取消订阅/排空；否则同步投递（V-F104） |
| **业务域** | 业务域模块激活后成为本端口消费者（模块间解耦通信）；届时按成功边界 #6 评估是否需要 outbox/broker |

## 方向级退出判据

1. **端口契约冻结**：EventBus 端口（类型化 Publish/Subscribe/Unsubscribe + 订阅生命周期 + 错误语义）冻结并可用；快测可断言。
2. **进程内实现可用**：channel 分发 + 订阅管理 + 错误语义实现并有测试（发布/订阅/退订、并发、顺序、handler panic 隔离）。
3. **接缝声明落盘**：应用契约 vs 运输实现边界（outbox/MQ）写入；不引入 broker 客户端依赖；不实现 outbox。
4. **对齐登记**：与 roadmap Admin 功能分支 typed domain event 扩展接缝登记对齐；**不解除**其 trigger-gated（应用契约仍归 Admin 功能分支）。
5. **共享约定登记**：topic / 订阅命名 + 契约测试 harness 约定在架构短文或 owner VP 决策登记；**不**纳入 Redis key 轨道（VP-026/027 专属）。
6. **停机与边界语义（V-F104）**：若选异步投递须声明 SIGTERM 取消订阅/排空；否则同步投递。
7. **边界保持**：未改 Charter；未改 Profile 默认集 / 模块矩阵 / Manifest 装配；未预制 outbox/broker（不消耗 RT-Q02 trigger，不预裁 RT-Q06 表结构）；未重开历史 VP。
8. **审计闭合**：开放 required finding = 0（或已合法闭合）。

详细纲领阶段由 lead Root（P-001）书写：R1 契约冻结（API 形态 / 投递语义 / 错误语义）→ R2 进程内实现 → R3 接缝与对齐 → R4 证据与关门。本 VP 不写 Goal 五件套。

## 信息需求（P-005）

| id | 要回答的问题 | 级别 | 影响门禁 | 最晚阶段 | 状态 |
|----|--------------|------|----------|----------|------|
| I-028-001 | 事件类型化机制：Go 接口断言 vs 注册表（topic → type）。**显式取舍（V-F103）**：进程内 channel 可传非序列化负载 vs 为 outbox/RT-Q06 预留的可序列化约束——R1 须记录未选方案；若选注册表，I-028-004 升为 required。 | required | 方案冻结 + 退出判据 1 | R1 契约冻结 | **verified**（2026-09-01 用户裁决：注册表 + JSON 可序列化） |
| I-028-002 | 投递语义默认：同步（发布者阻塞） vs 异步（channel 缓冲）；**缓冲满时的最小语义（V-F103）** = 阻塞 / 丢弃 / 返回错误（R1 只冻结其一；完整背压产品仍 gated）。 | required | 退出判据 2 | R1 | **verified**（2026-09-01 用户裁决：异步 + 缓冲满阻塞） |
| I-028-003 | handler 错误语义：失败吞掉 + 日志 vs 回传发布者 vs 隔离失败（panic 恢复）；重复发布者可见性。 | required | 退出判据 2 | R1 | **verified**（2026-09-01 用户裁决：吞掉+日志 + panic 隔离） |
| I-028-004 | 事件类型注册权属（业务域 VP vs Admin 功能 VP）与 typed domain event gated 保持：本 VP 不解除 Admin gated（V-F101）；因 I-028-001 选注册表升 required。 | required | 退出判据 4 | R3 | 待确认 |

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| workspace-028-event-bus-port | GOAL-001-event-bus-port | lead | 2026-09-01 | 激活开区（VRev-064 self `pass` · 架构类 freshness PASS `5744868d`→`29727510`） |

## 关门记录

（仅 `closed` / `abandoned` 时填写。）

| date | outcome | summary | evidence_links | residuals |
|------|---------|---------|----------------|-----------|
| — | — | — | — | — |

## 规划修订短史

| date | change |
|------|--------|
| 2026-08-31 | 初创 `planned`：用户裁决按 3 个独立 VP 执行（缓存 / 限流 / 事件总线；触发条件独立 × 关门能力独立原则）。本 VP 承接 RT-Q02 的应用契约前置（进程内 EventBus + outbox/MQ 接缝声明，broker 仍 gated）；与 Admin 功能分支 typed domain event 扩展接缝对齐；vision_ref @0.4.0；roadmap / revisions 原子同步 |
| 2026-08-31 | v0.1.1 · **VRev-059 响应修订**（grok build · conditional → 本 VP 无 required）：V-F101 **fixed**——定位统一为"运输端口 + 进程内实现 + outbox/MQ 接缝声明"，删除"应用契约前置"含糊措辞，明确不解除 Admin typed domain event gated、EventBus ≠ Job 端口；V-F102 **fixed**——补"不消耗 RT-Q02 trigger"解释规则；V-F103 **fixed**——序列化取舍改为 R1 显式取舍（含未选方案），I-028-002 背压收窄为缓冲满最小语义；V-F104 **fixed**——补异步投递停机语义声明；V-F100 部分 **fixed**——共享约定改为 topic/订阅命名（不纳入 Redis key 轨道） |
| 2026-09-01 | v0.2.0 · **激活**：用户指令「/vision 激活 vp-028，然后交 /govern 开设新工作区」；VRev-064 self `pass`（0 required · 架构类 freshness PASS `5744868d`→`29727510` 五域零变更 · 区间代码 = VP-027 已审结目交付，不暂挂 `go`）；lead `workspace-028-event-bus-port` / Root `GOAL-001-event-bus-port` |