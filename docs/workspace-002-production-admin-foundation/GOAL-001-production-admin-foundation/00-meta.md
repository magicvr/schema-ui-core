---
id: GOAL-001-production-admin-foundation
title: 生产级可用 Admin 基架
status: active
created: 2026-08-01
updated: 2026-08-02
parent: null
version: 0.2.6
progress: 3/5
plan_refs:
  - VP-002-production-admin-foundation
primary_plan: VP-002-production-admin-foundation
serves_summary: 在 VP-001 冻结协议基线之上，把现有 Demo 推进为具备真实认证、持久化权限、Schema 驱动 CRUD 与可复现工程交付的生产级 Admin 基架。
---

# GOAL-001 · 生产级可用 Admin 基架

## 概述

本 Root 承接 [VP-002 · 生产级可用 Admin 基架](../../vision/plans/VP-002-production-admin-foundation.md)。目标不是扩大冻结协议覆盖面，而是在既定边界内把当前前端演示推进为可运行、可验证、可 fork 的生产级 Admin 基架。

## 愿景对齐

| 字段 | 值 |
|------|----|
| `plan_refs` | `VP-002-production-admin-foundation` |
| `primary_plan` | `VP-002-production-admin-foundation` |
| Charter | `schema-ui-core-admin-foundation@0.1.0` |
| 工作区角色 | `delivery` |

继承范围以 VP-002 引用的 `I-PROTO-001 v0.1.3` 为准。冻结基线是实施输入，不是本目标的实施完成证据。

## 成功边界

- Schema Renderer 在冻结协议范围内成为默认页面能力，并有结构、运行时与失败路径证据。
- 登录、登出、会话恢复和请求级身份来自真实认证链路，不再由前端本地状态模拟。
- 用户、角色、菜单及最小权限关系可持久化，并由后端实施授权边界。
- 至少一个代表性实体完成 Schema 驱动的列表、新建、编辑、删除与统一校验/错误闭环。
- 种子数据、环境配置、健康检查、部署/容器路径与 fork 文档可重复执行；目标用户可在 15 分钟内启动并看到可用后台。
- 各阶段留下可核对的实现事实与审计结论；关门前不存在开放 required 信息项或必改 finding。

## 纲领路线图

五个检查点默认等权并原则上串行；同一阶段内部可以按依赖拆分并行子目标。

- [x] **R1 · 协议实施边界与 Schema Renderer 产品化**：核对冻结覆盖映射，把 Renderer 接入默认页面路径，并验证关键失败行为。  
  阶段内子目标（D-005）：`GOAL-002-r1-schema-load-validate`（done）· `GOAL-003-r1-default-render-path`（done）· `GOAL-004-r1-representative-node-pages`（done）。证据：I-001 覆盖矩阵 verified（D-004）+ Renderer 默认 `schemaUrl` 主路径 + 2026-08-02 Web 425/425、Go test/vet 全绿与 fail-closed 断言（各子目标关门审计）。
- [x] **R2 · 真实认证与请求级身份**：交付登录、登出、会话恢复、受保护路由和 API 身份传递。  
  阶段子目标：`GOAL-005-r2-auth-session`（**done**，2026-08-02；`I-002` verified + D-007 方案）。证据：登录/登出/刷新/撤销与请求级身份中间件 + SQLite 存储 + 前端认证闭环（441 单测、`go test ./...` 全绿）；close-out 审计 A-001（independent，F-001 → fixed）+ A-002（self，pass）；browser E2E 在 Linux CI 通过（run #30711903555，`1 passed`，含匿名 401 断言）。
- [x] **R3 · 持久化身份、角色与最小权限模型**：交付用户/角色/菜单持久化、种子数据与后端授权最小闭环。  
  阶段子目标：`GOAL-006-r3-persistent-rbac-menu`（**done**，2026-08-02；D-009 冻结方案 B、`features` 菜单投影、两步迁移、读写权限与恢复证据口径；D-004 S1 迁移链 + pre-v0002 快照；D-005 S2 阶段 B 终态；D-006 S3 增量幂等种子；D-007 S4 permission key 读写门禁；D-008 S5 `me.features` 投影 + manifest `visibleWhen`；S6 恢复/重启/回归证据齐备）。证据：`schema_migrations` + `0001/0002` 事务化迁移 + pre-v0002 恢复快照、规范化 RBAC 双写/集合核对、`seedRBAC` 增量幂等种子、records `records.read`/`records.write` 门禁、`me.features` 菜单投影 + 真实 manifest `visibleWhen`、`TestRestartPersistence`/`TestRestorePreV0002Snapshot` 与 API/Web 全量回归；close-out 审计 A-005（independent，F-005 → fixed）+ A-006（independent，F-005 关闭复核 pass）+ A-004（self 阶段审计）；无开放 required finding 或 required 信息项。
- [ ] **R4 · Schema 驱动 CRUD 与统一交互闭环**：以代表性实体验证列表、表单、操作、校验、加载/空态/错误态及权限失败。
- [ ] **R5 · 工程化、fork 体验与集成关门**：完成环境/容器/健康检查/文档、可重复验收、阶段审计与 Root 关门审计。

当前派生进度为 `3/5`（R1、R2、R3 已勾选，2026-08-02）。勾选仅能由对应阶段的可验证事实和审计结论驱动，不得用百分比替代门禁判断。R1～R3 子目标 progress 不替代本 Root 检查点。

## 信息需求与阶段门禁

| ID | 问题 / 所需信息 | 级别 | 影响门禁 | 最晚阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据或结论 |
|----|-----------------|------|----------|----------|-----------------|------|-------------|------------|
| `I-001` | 冻结协议到当前代码、fixture 与 Renderer 运行路径的实施差量是什么？ | required | R1 方案冻结与实施 | R1 方案冻结前 | 以 `I-PROTO-001 v0.1.3` 逐项建立实现与验证矩阵 | **verified** | 已关闭（D-004） | [I-001-implementation-gap-matrix.md](attachments/I-001-implementation-gap-matrix.md)；D-004 采用为 R1 方案边界；**不**勾选 R1 完成 |
| `I-002` | 认证/会话机制、凭据边界与安全配置采用什么最小方案？ | required | R2 方案冻结与实施 | R2 方案冻结前 | 调查当前栈、部署约束与威胁边界，形成认证生命周期、请求身份、配置边界与验收矩阵后记录方案取舍 | **verified** | 已关闭（D-007） | D-007（2026-08-02）裁决：短 JWT Access + Opaque Refresh + SQLite 存储 + 接受 JWT 库；现状与三候选方案、M1–M14 验收矩阵见 [I-002-auth-collection.md](attachments/I-002-auth-collection.md)；**不**冻结 R2 实施细节（TTL / env / 前端存储 / CORS / 表结构在 R2 子目标定稿） |
| `I-003` | 数据存储、迁移、种子和用户—角色—菜单关系的最小模型是什么？ | required | R3 方案冻结与实施 | R3 方案冻结前 | 对照 R2 SQLite 占位、真实授权/导航链与既有测试，形成候选数据模型、版本迁移、增量种子和恢复验证矩阵后提交用户裁决 | **verified** | 已关闭（D-009） | [I-003-persistence-permission-collection.md](attachments/I-003-persistence-permission-collection.md)；D-009 采用方案 B + `features` 菜单投影 + 两步迁移 + `records.read`/`records.write` + 自动恢复证据；仅解除 R3 方案/立项目门禁，**不代表已实现** |
| `I-004` | 哪个代表性实体及 API/错误语义能够完整证明 Schema CRUD 闭环？ | required | R4 方案冻结与验收 | R4 方案冻结前 | 选择实体并固定 CRUD、校验、权限与异常验收矩阵 | open | 不延期；R4 立项前复核 | 待收集 |
| `I-005` | 目标部署基线、15 分钟 fork 计时口径与可复现实验环境是什么？ | required | R5 方案冻结与关门 | R5 方案冻结前 | 固定环境、命令、容器/部署边界和独立复现实验方法 | open | 不延期；R5 立项前复核 | 待收集 |
| `I-006` | 是否在本波次纳入最小操作日志？ | non-blocking | R5 范围取舍 | R5 方案冻结前 | 评估对交付价值与成本；如升级为必需则另记决策 | open | R5 立项时复核 | VP-002 将其列为加分项而非硬关门条件 |

开放信息项不妨碍 Root 立项，但会阻断其列明的阶段门禁。任何 residual 接受必须由用户书面裁决并记录范围与复审触发条件。

## 层级

本目标是工作区 Root，`parent: null`。后续子目标必须在本工作区根平铺，并使用本目标完整 id 作为 `parent`。
