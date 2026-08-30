---
doc_type: vision-plan
id: VP-012-shared-cross-module-contracts
title: 共享横切契约与平台基架（Cross-module Contracts）
status: closed
vision_ref: schema-ui-core-admin-foundation@0.3.0
lead_workspace: workspace-012-shared-cross-module-contracts
created: 2026-08-18
updated: 2026-08-19
version: 0.2.0
parent: null
---

# VP-012 · 共享横切契约与平台基架

## 状态与门闩（2026-08-19 · 已关门）

| 项 | 值 |
|----|-----|
| status | **`closed`**（2026-08-19 用户书面确认完整关门；VRev-028 `V-F057` → `fixed`） |
| **lead_workspace** | **`workspace-012-shared-cross-module-contracts`**（Root `GOAL-001-shared-cross-module-contracts` `done 8/8`） |
| **Vision required** | **已满足**：VRev-028 `pass`，open required = 0；`V-F057` recommended 由本关门记录闭合 |
| **关门门闩（现行）** | 已 `closed`；保留 workspace-012 历史绑定，默认不接新区；reopen 须用户确认 |
| **完整 ≠ 方向表无限扩张** | 首波退出分母 = R1～R6。**保留/归档**改由 lead 区 R7（GOAL-008）按用户书面设置策略交付。session/effective actor 与 D-003 外 writer envelope 仍在 `roadmap.md` Tier A |

## 意图

在 VP-011 已交付标准 Admin 功能模块、且 `I-011-002` 已登记四档能力地图的基础上，建立并交付**所有未来模块共同依赖的横切契约与平台基架**：correlation/request-id、审计事件模型增强、乐观并发/幂等、异步 Job/长操作、maintenance/degraded/read-only 门控、API Token/Service Credential。

本 VP **不承载具体业务领域**（订单/支付/库存/CMS 等仍属 Tier D），**不替代 VP-009/VP-010 的持续程序**；它把“横切契约”作为可交付、可验证的基架能力落地。

## 首波冻结（退出分母）

首波 = P0 四项 + P1 两项，对应 lead 区 R1～R6。方向表中宽于冻结切片的能力**不进入**本 VP 退出判据。

| 能力 | 首波已交付 | 不进本 VP 退出分母（后续 Tier A） |
|------|------------|----------------------------------|
| correlation / request-id / 错误恢复 | requestid、错误包络、前端可引用、operationlog 关联 | OpenTelemetry / 分布式 tracing |
| 审计事件模型增强 | 结构化 before/after/diff、递归脱敏、correlation；auth/settings/users 真实消费 | **session/effective actor 关联**；**保留/归档触发**；D-003 未列入的写路径改走结构化 envelope |
| 乐观并发 / 幂等 | wallet ETag / expectedVersion / 409 / operation replay | 不批量改造全部 repository |
| 异步 Job / 长操作 | 六态、进度、重试、取消、结果读取/过期；wallet reconcile 202 | 通用 Job 管理页；外部队列 |
| maintenance / degraded / read-only | 四模式、统一写门禁、Host/status 投影 | 运行时管理 UI（原文「UI 可后置」） |
| API Token / Service Credential | hash-only、管理 API、Bearer、scope、吊销、使用审计 fail-closed | 外部 IdP / OAuth / HSM |

**为何不在本波实现后续项**（2026-08-19 `/vision` 核对）：

1. **session / effective actor**：现行 `operation_log` 只有 `actor_id` / `actor_name` + correlation 侧表；代码中无 impersonation / effective actor，也无 session 列。补齐需要新 schema、会话绑定契约，以及尚不存在的产品能力，不能塞进已 `done` 的 Root。
2. **保留 / 归档触发**：workspace-003 已将 operationlog 自动 purge/archive 记为独立生命周期；本仓无 duration / archive / restore 合同。这是数据生命周期程序，不是首波横切契约补丁。
3. **D-003 外写路径**：`users_state` / MFA / wallet 已写审计事件，但 detail 多为 ad-hoc JSON。D-003 冻结的必选消费面是 auth/settings/users；其余 writer 的 envelope 迁移是增量采用，不是首波缺口。

上述三项移交 [roadmap.md](../roadmap.md) 四档地图 **Tier A**，按触发条件新建 VP/工作区。不得读成 VP-012 未交付。

## 方向级范围

| 能力 | 方向级内容 |
|------|------------|
| correlation / request-id / 错误恢复 | 请求级 correlation id、错误响应语义、前端错误页可引用标识 |
| 审计事件模型增强 | 结构化 before/after diff、敏感字段脱敏、correlation；session/effective actor 与保留/归档见首波冻结表 |
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

1. 首波横切契约（P0 四项 + P1 两项，见首波冻结表）已交付，并有至少一个真实模块或验证路径消费/引用；
2. 契约不改变 Charter 边界；不把 Tier D 业务域纳入本 VP；
3. 共享基架安全/符合性问题按 VP-009/VP-010 分流，未塞入本 VP；
4. 开放 required finding = 0（或已合法闭合）。

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| workspace-012-shared-cross-module-contracts | GOAL-001-shared-cross-module-contracts | lead | 2026-08-18 | 2026-08-18 用户确认新建；2026-08-19 VP `closed`；Root done 6/6；默认不接新区 |

## 关门记录

| date | outcome | summary | evidence_links | residuals |
|------|---------|---------|----------------|-----------|
| 2026-08-19 | **closed**（完整 · 首波） | 用户确认完整关门。exit 1：R1～R6 冻结切片已交付且有真实消费。exit 2：未改 Charter、未进 Tier D。exit 3：安全/符合性仍归 VP-009/010。exit 4：实现层与 VRev required = 0。V-F057 按本表闭合。 | `workspace-012` goal-tree（Root done 6/6；GOAL-002～007 done）；Root A-001/A-002/A-003 close-out；A-012 independent pass / A-013（F-010 fixed）；E-007；GOAL-002 A-001；GOAL-003 A-006/A-007 + D-003/D-004；GOAL-004 A-004/A-005；GOAL-005 A-012/A-013；GOAL-006 A-008/A-009；GOAL-007 A-009/A-010 + Root A-012；[VRev-028](../reviews/VRev-028-vp012-closeout-readiness.md) | **无本 VP 未完成项。** session/effective actor、保留/归档触发、D-003 范围外写路径的结构化 envelope 已移交 `roadmap.md` 四档地图 Tier A（触发后新建 VP），不构成本 VP residual |

### 退出判据 ↔ 证据

| 退出 | 结论 | 证据 |
|------|------|------|
| 1 首波 + 真实消费 | 满足 | R1 GOAL-002 requestid/错误包络/operationlog；R2 GOAL-003 结构化 detail/脱敏/correlation（auth/settings/users）；R3 GOAL-004 wallet ETag/409/replay；R4 GOAL-005 Job 六态 + wallet reconcile 202；R5 GOAL-006 四模式写门禁 + Host/status；R6 GOAL-007 service-credential 管理 API / Bearer |
| 2 Charter / Tier D | 满足 | Root A-002：无新业务模块/页面/导航/fragment；wallet 仅为既有模块消费面 |
| 3 009/010 分流 | 满足 | VP-012 / Root 声明；本波未扩安全程序或符合性程序 |
| 4 required = 0 | 满足 | Root A-012/A-013；VRev-028 open required = 0 |

## 规划修订短史

| date | change |
|------|--------|
| 2026-08-18 | 初创并激活：用户确认新建独立 VP/workspace 承载横切契约波；首波 = P0 四项（correlation/审计模型/并发幂等/异步 Job）+ P1 两项（maintenance 门控/API Token） |
| 2026-08-19 | VRev-028 self `pass`：关门就绪；V-F057 recommended 约束关门落盘形状 |
| 2026-08-19 | 用户确认完整关门：v0.2.0 `active → closed`。冻结首波退出分母；D-004 三项移交 roadmap Tier A（不可在本波实现）；关门记录含 exit↔证据映射；组合索引原子同步（VR-024） |
