---
doc_type: vision-plan
id: VP-012-shared-cross-module-contracts
title: 共享横切契约与平台基架（Cross-module Contracts）
status: active
vision_ref: schema-ui-core-admin-foundation@0.2.0
lead_workspace: workspace-012-shared-cross-module-contracts
created: 2026-08-18
updated: 2026-08-18
version: 0.1.0
parent: null
---

# VP-012 · 共享横切契约与平台基架

## 意图

在 VP-011 已交付标准 Admin 功能模块、且 `I-011-002` 已登记四档能力地图的基础上，建立并交付**所有未来模块共同依赖的横切契约与平台基架**：correlation/request-id、审计事件模型增强、乐观并发/幂等、异步 Job/长操作、maintenance/degraded/read-only 门控、API Token/Service Credential。

本 VP **不承载具体业务领域**（订单/支付/库存/CMS 等仍属 Tier D），**不替代 VP-009/VP-010 的持续程序**；它把“横切契约”作为可交付、可验证的基架能力落地。

## 方向级范围

| 能力 | 方向级内容 |
|------|------------|
| correlation / request-id / 错误恢复 | 请求级 correlation id、错误响应语义、前端错误页可引用标识 |
| 审计事件模型增强 | 结构化 before/after diff、敏感字段脱敏、correlation/session/effective actor 关联、保留/归档触发 |
| 乐观并发 / 幂等 | expectedVersion / ETag / 409 冲突、idempotency_key / operation_id 状态机 |
| 异步 Job / 长操作 | queued/running/succeeded/failed/cancelled/expired、进度、重试、取消、结果下载/过期 |
| maintenance / degraded / read-only | 后端写边界统一拒绝、状态契约、错误语义；UI 可后置 |
| API Token / Service Credential | 机器凭据管理面、作用域、吊销、审计；与用户会话分离 |

## 与相邻 VP 的边界

| VP | 关系 |
|----|------|
| **VP-011**（标准 Admin 功能模块） | 消费本 VP 的契约；本 VP 不新增 S/B 业务页面 |
| **VP-009**（生产加固） | 安全缺陷/威胁面仍归 009；本 VP 提供可被 009 消费的审计/并发/维护契约 |
| **VP-010**（设计—实现符合性） | 符合性 gap 仍归 010；本 VP 交付的契约须符合已固化架构与协议面 |
| **VP-008**（准入 `go`） | 如契约改变 Profile 默认集 / 模块矩阵 / Manifest 装配语义 / 共同门禁，按 VP-008 消费有效性规则处理 |

## 方向级退出判据

1. 首波横切契约（P0 四项 + P1 两项）已交付，并有至少一个真实模块或验证路径消费/引用；
2. 契约不改变 Charter 边界；不把 Tier D 业务域纳入本 VP；
3. 共享基架安全/符合性问题按 VP-009/VP-010 分流，未塞入本 VP；
4. 开放 required finding = 0（或已合法闭合）。

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| workspace-012-shared-cross-module-contracts | GOAL-001-shared-cross-module-contracts | lead | 2026-08-18 | 2026-08-18 用户确认新建；首波 = 横切契约波 |

## 关门记录

（仅 `closed` / `abandoned` 时填写。）

| date | outcome | summary | evidence_links | residuals |
|------|---------|---------|----------------|-----------|
| — | — | — | — | — |

## 规划修订短史

| date | change |
|------|--------|
| 2026-08-18 | 初创并激活：用户确认新建独立 VP/workspace 承载横切契约波；首波 = P0 四项（correlation/审计模型/并发幂等/异步 Job）+ P1 两项（maintenance 门控/API Token） |
