---
doc_type: goal-decision
id: D-001-workspace-root-establishment
parent: GOAL-001-event-bus-port
date: 2026-09-01
status: active
version: 0.1.0
---

# D-001 · 工作区 / Root 建立与开区决策

## 上下文

用户 2026-09-01 指令「/vision 激活 vp-028，然后交 /govern 开设新工作区」；slug 已确认（按惯例：`workspace-028-event-bus-port` / `GOAL-001-event-bus-port`，沿用 VP-013～027）。激活门禁已满足：VRev-064 self `pass`（0 required；VRev-058/059 全闭合）+ 架构类轻量 freshness PASS（`5744868d` → `29727510` 五域零变更，不暂挂 `go`）。

## 决策

| 项 | 决定 |
|----|------|
| 工作区 | `workspace-028-event-bus-port`（canonical `docs/workspaces/workspace-028-event-bus-port/`） |
| Root | `GOAL-001-event-bus-port`（`parent: null`；primary_plan = `VP-028-event-bus-port`） |
| 愿景角色 | `delivery`（不改变 Charter primary workspace） |
| 纲领阶段 | R1 契约冻结 → R2 进程内实现 → R3 接缝与对齐 → R4 证据与关门（串行；阶段内可并行子目标） |
| 审计模式 | 阶段关门 default **self**；R4 证据/关门实证门禁可按需 **independent**（grok build 先例 · 项目级默认执行路径） |
| 红线 | 不预制 outbox/broker（不引入客户端依赖 / 不预裁 RT-Q06 / 不消耗 RT-Q02 trigger）；不改 Profile 默认集 / 模块矩阵 / Manifest（VP-008 `go`）；不解除 Admin typed domain event gated；EventBus ≠ Job 端口；不属 Redis 轨道；停机语义继承 VP-021 |

## freshness 三字段（VRev-064 · 先例执行惯例）

| 字段 | 值 |
|------|-----|
| consumer_vp | `VP-028-event-bus-port`（vision_ref `schema-ui-core-admin-foundation@0.4.0`） |
| last_freshness_review_at | 2026-09-01（`5744868d` → `29727510` · 架构类轻量 PASS · 五域零变更） |
| next_freshness_review_trigger | 首个 C 端业务域 VP 激活（H-002 发现机制）或 多实例部署评估 |

## 未选方案

- 不合并 VP-026/027 同期范围（关门独立原则；两 VP 已 closed）。
- 不在本波实现 outbox / 外部 broker（RT-Q02 触发条件未成立，trigger-gated 保持）。
- 不把本区挂到 workspace-012（VP-012 已 closed，禁止重开吸收本意图）。
