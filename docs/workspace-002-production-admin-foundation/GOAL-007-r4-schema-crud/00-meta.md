---
id: GOAL-007-r4-schema-crud
title: R4 · Schema 驱动 CRUD 与 SQLite 持久化闭环
status: active
created: 2026-08-02
updated: 2026-08-02
parent: GOAL-001-production-admin-foundation
version: 0.2.0
progress: 2/6
---

# GOAL-007 · R4 · Schema 驱动 CRUD 与 SQLite 持久化闭环

## 概述

承接 Root **D-010 / D-011**：以现有 `records` 为代表实体，把进程内数据迁入 SQLite，并通过统一错误 envelope、Schema 驱动的列表/搜索/详情/新建/编辑/删除交互，以及服务重启后的持久化证据，交付 R4 的标准 CRUD 闭环。

## 成功标准

- [x] **S1 · 精确 CRUD 与错误契约冻结**：冻结 records 字段、ID/时间戳、create/update/delete 请求响应、HTTP status、稳定 error `code`、校验与冲突语义，并以测试矩阵关联现有统一错误 envelope。（D-002 + [I-007-001-api-error-contract.md](attachments/I-007-001-api-error-contract.md)）
- [x] **S2 · SQLite 结构、迁移与种子冻结**：冻结 records DDL、索引/约束、迁移版本与校验和、可重复 seed、事务/并发语义和失败恢复路径。（D-003 + [I-007-002-sqlite-migration-plan.md](attachments/I-007-002-sqlite-migration-plan.md)）
- [ ] **S3 · 持久化 CRUD API**：实现 POST/list/search/detail/PATCH/DELETE 的 SQLite repository 路径；生产默认不再依赖进程内 records，认证、`records.read` / `records.write` 与统一错误 envelope 保持一致。
- [ ] **S4 · Schema 驱动读写主路径**：由 Schema 页面完成列表、搜索、详情、新建、编辑与删除；新增或调整代表页面不修改 Renderer 主路径代码。
- [ ] **S5 · 交互状态与权限负向闭环**：验证字段校验、加载/空态、成功反馈、删除确认、统一错误展示，以及匿名 `401`、已认证缺权限 `403` 和后端授权不被前端隐藏替代。
- [ ] **S6 · 重启、迁移与端到端回归**：自动证明 create/update/delete 后服务重启，list/detail 结果符合持久化预期；覆盖 migration/seed 重复执行、关键失败路径及 API/Web 回归。

## 派生进度

`progress: 2/6` 由上方六个顺序检查点等权派生（S1/S2 已勾选）。S1/S2 为契约冻结入口；S3 已获 `I-007-001`/`I-007-002` 放行，S4/S5 仍依赖 `I-007-003`，S6 依赖 `I-007-004` 与可重复证据。检查点不替代审计 finding 或关门结论。

## 信息需求

| ID | 问题 / 所需信息 | 级别 | 影响门禁 | 最晚阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据或结论 |
|----|-----------------|------|----------|----------|-----------------|------|-------------|------------|
| `I-007-001` | records 的精确字段/ID/时间戳、create/update/delete 请求响应、HTTP status、稳定 error `code`、校验与冲突语义是什么？ | required | S1 完成；S3 API 首个受影响代码变更 | S3 实施前 | 对照现有 handler/error envelope 与 I-004 M-R4-01～06/08，形成版本化 API/错误契约和正反测试矩阵，记录决策 | **verified** | 已关闭（D-002） | [I-007-001-api-error-contract.md](attachments/I-007-001-api-error-contract.md)：继承 list/detail/PATCH/DELETE；POST 201 + `INVALID_CREATE_*`；T-API-01～13 |
| `I-007-002` | records 的 SQLite DDL、迁移版本/校验和、索引/约束、seed、事务/并发与失败恢复契约是什么？ | required | S2 完成；S3 持久化首个代码变更 | S3 实施前 | 对照现有 migration runner、seed 与 records 访问模式，形成 schema/migration/seed/repository 计划及恢复测试矩阵，记录决策 | **verified** | 已关闭（D-003） | [I-007-002-sqlite-migration-plan.md](attachments/I-007-002-sqlite-migration-plan.md)：`0003 records_persist`、空表 seed、repository/LWW、T-DB-01～09 |
| `I-007-003` | Schema CRUD 的页面/Node/action 绑定、字段映射、成功/加载/空态/错误/确认交互和权限矩阵是什么？ | required | S4/S5 实施与验收 | 首个 Schema 写交互代码变更前 | 对照真实 fixtures、Renderer action/form/table 能力与 records API，冻结页面绑定及 admin/viewer/匿名正反矩阵，记录决策 | open | 不延期；S4 前复核 | 待收集；前端隐藏不得替代后端 `records.read` / `records.write` |
| `I-007-004` | 重启保持与端到端验收如何隔离数据库、固定操作序列、断言持久化结果并清理测试状态？ | required | S6 验收执行与目标关门 | S6 验收前 | 形成 create/update/delete→重启→list/detail 的可重复协议，含 migration/seed 重跑、权限负向、失败路径和 API/Web 回归命令 | open | 不延期；S6 前复核 | 待收集；必须产出机器可重复证据，不以单次手工冒烟替代 |

## 依赖与边界

| 项 | 说明 |
|----|------|
| 父目标 | [GOAL-001-production-admin-foundation](../GOAL-001-production-admin-foundation/00-meta.md)（D-010 / D-011；Root `I-004` = verified） |
| 前置证据 | [I-004-schema-crud-collection.md](../GOAL-001-production-admin-foundation/attachments/I-004-schema-crud-collection.md)（records 现状、候选比较与 M-R4-01～11） |
| In | records SQLite 持久化、版本迁移/seed、CRUD API、统一错误 envelope、Schema CRUD 页面与状态、权限负向、重启/回归证据 |
| Out | users/RBAC 管理后台；新业务实体；通用工作流/批量 CRUD；扩大 `I-PROTO-001 v0.1.3`；R5 容器、生产运维与 fork 关门 |

## 父目标

- [GOAL-001-production-admin-foundation](../GOAL-001-production-admin-foundation/00-meta.md)
