---
doc_type: goal-decision
id: D-001-workspace-root-establishment
parent: GOAL-001-cache-port
date: 2026-08-31
status: active
version: 0.1.0
---

# D-001 · 工作区 / Root 建立与开区决策

## 上下文

用户 2026-08-31 指令「激活 VP-026 并交编排器开始工作区」。激活门禁已满足：VRev-060 self `pass`（0 required）+ 架构类轻量 freshness PASS（`055da2fd` → `54fb57e7` 五域零变更，不暂挂 `go`）。

## 决策

| 项 | 决定 |
|----|------|
| 工作区 | `workspace-026-cache-port`（canonical `docs/workspaces/workspace-026-cache-port/`） |
| Root | `GOAL-001-cache-port`（`parent: null`；primary_plan = `VP-026-cache-port`） |
| 愿景角色 | `delivery`（不改变 Charter primary workspace） |
| 纲领阶段 | R1 合同冻结 → R2 内存供应商+双策略 → R3 接缝与共享约定 → R4 证据与关门（串行；阶段内可并行子目标） |
| 审计模式 | 阶段关门 default **self**；R4 证据/关门实证门禁可按需 **independent**（grok build 先例 · 项目级默认执行路径） |
| 红线 | 不预制 Redis（不引入客户端依赖 / 不消耗 RT-Q03 trigger）；不改 Profile 默认集 / 模块矩阵 / Manifest（VP-008 `go`）；Redis 轨道共享约定单一所有者；停机语义继承 VP-021 |

## freshness 三字段（VRev-060 · V-F089 执行惯例）

| 字段 | 值 |
|------|-----|
| consumer_vp | `VP-026-cache-port`（vision_ref `schema-ui-core-admin-foundation@0.4.0`） |
| last_freshness_review_at | 2026-08-31（`055da2fd` → `54fb57e7` · 架构类轻量 PASS · 五域零变更） |
| next_freshness_review_trigger | 首个 C 端业务域 VP 激活（H-002 发现机制）或 多实例部署评估 |

## 未选方案

- 不合并 VP-027/028 同期开区（关门独立原则；VP-027 分母修正后另行激活）。
- 不在本波实现 Redis 供应商（trigger-gated 保持）。