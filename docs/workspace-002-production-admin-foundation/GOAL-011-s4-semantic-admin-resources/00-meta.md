---
id: GOAL-011-s4-semantic-admin-resources
title: S4 · 语义化 Admin 资源替换与双实体验证
status: active
created: 2026-08-03
updated: 2026-08-03
parent: GOAL-010-a002-schema-adapter
version: 0.3.0
progress: 2/5
---

# GOAL-011 · S4 · 语义化 Admin 资源替换与双实体验证

## 概述

承接父目标 GOAL-010 的 S4：将当前仅用于代表性 CRUD 验证、缺少真实业务语义的 `records` 从产品默认运行面退场，以 `users` 替换其代表实体，并新增 `roles` 作为第二个语义资源。目标是在不修改前端 Renderer 主路径的前提下，让两个对绝大多数 Admin 系统均有实际价值的资源完成 Schema 驱动的标准 CRUD 闭环，同时保留既有迁移与历史治理记录的可追溯性。

“只修改 Schema 接入”限定为：后端资源契约、持久化、权限与领域不变量已完成注册后，前端新增页面只修改 Schema/fixture，不修改 Renderer 主路径；不主张仅凭 Schema 自动生成后端领域逻辑。

## 路线图 / 成功标准

五个检查点默认等权、原则上串行（P-001）。当前进度 `2/5`（S1/S2 已完成；S3～S5 未实施）：

- [x] **S1 · 语义资源与退场契约冻结**：关闭 `I-011-001`/`I-011-002`，冻结 users/roles 的最小领域边界、安全不变量，以及 records API/表/权限/菜单/操作日志/测试的版本化退场与既有数据库迁移策略。
- [x] **S2 · users/roles 后端资源闭环**：在通用资源工厂之上实现 users 与 roles 的持久化 list/search/detail/create/update/delete、字段校验、敏感字段隔离、关系/系统角色保护、稳定错误 envelope 与 401/403 负向路径。
- [ ] **S3 · records 产品运行面退场**：按冻结策略移除当前产品默认运行面中的 records API 注册、Schema fixture、菜单/权限、种子、操作日志耦合、前端 records 专名与当前测试依赖；保留不可改写的历史治理事实和迁移链证据。
- [ ] **S4 · 双语义实体 Schema 接入验证**：关闭 `I-011-003`；users 与 roles 的代表性列表/CRUD 页面由 Schema 接入，前端 Renderer 主路径无修改；fresh fork、既有数据库升级、重启持久化与权限失败路径可复核。
- [ ] **S5 · 回归、审计与父级交接**：API/Web/build/E2E 全量回归与阶段/关门审计通过；向 GOAL-010 提供其 S4 可勾选的证据，目标自身无开放 required finding 或到期 required 信息项。

## 派生进度

`progress` 由 S1～S5 五个检查点等权派生；当前为 `2/5`。检查点不替代审计、信息门禁或关门结论。

## 信息需求与阶段门禁

| ID | 问题 / 所需信息 | 级别 | 影响门禁 | 最晚阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据或结论 |
|----|-----------------|------|----------|----------|-----------------|------|-------------|------------|
| `I-011-001` | users/roles 的精确资源契约与最小 IAM 边界是什么？包括公开字段、密码/令牌隔离、角色分配、self/最后管理员保护、system role 与 grant 约束。 | required | S1 方案冻结与 S2 实施 | S2 首个产品代码变更前 | 核对现有 auth/RBAC 表、store、权限投影与通用 Resource 接口，形成版本化领域契约并提交用户裁决 | **verified** | 不延期；S1 已冻结 | [I-011-001-users-roles-contract.md](attachments/I-011-001-users-roles-contract.md) **v0.2.0** + GOAL-011 D-002/D-003（用户裁决：通用工厂+最小契约扩展、操作日志纳入；A-002 F-001/F-003/F-004 fixed） |
| `I-011-002` | records 如何从当前产品基线退场，同时保持既有迁移/checksum、已有数据库、数据处置、API/权限/菜单/operation log 与历史文档可追溯？ | required | S1 方案冻结与 S3 退场实施 | S3 首个删除、重命名或迁移变更前 | 建立 fresh install 与 in-place upgrade 两条迁移矩阵，逐项盘点运行时代码、迁移、种子、文档与测试 | **verified** | 不延期；S1 已冻结 | [I-011-002-records-retirement.md](attachments/I-011-002-records-retirement.md) **v0.2.0** + GOAL-011 D-002/D-003（用户裁决：硬退场 DROP TABLE；A-002 F-002 快照语义 fixed） |
| `I-011-003` | 双资源验收如何证明“前端 Schema-only 接入”以及 fresh fork/升级/重启/401/403 的完整边界？ | required | S4 集成验收与 S5 关门 | S4 首个验收实现前 | 冻结 users/roles 两套验收矩阵、Renderer diff 边界、数据库升级夹具与 E2E/restart 协议 | **open** | 不延期；最晚 S4 前 | 待确认；不以通用工厂单测替代真实双资源产品证据 |

开放 required 信息项只阻断其列明的阶段；当前 S1 已冻结，`I-011-001`/`I-011-002` 已 verified，`I-011-003` 未到期；S2 实施完成，S3 退场门禁已解除。

## 依赖与边界

| 项 | 说明 |
|----|------|
| 父目标 | [GOAL-010-a002-schema-adapter](../GOAL-010-a002-schema-adapter/00-meta.md)（S4 交付载体；Root A-002 F-002-001 仍 open） |
| In | 最小 users/roles 管理资源、records 当前产品运行面退场、版本化迁移、权限/菜单/操作日志适配、双资源 Schema 页面与回归证据 |
| Out | 完整 IAM 产品、SSO/SCIM、复杂策略编排、多租户、扩大 `I-PROTO-001 v0.1.3`、Root/VP-002 关门 |
| 历史边界 | 不重写 GOAL-004/007 的历史决策、执行与审计；既有迁移仅按后续冻结的兼容策略演进 |

## 父目标

- [GOAL-010-a002-schema-adapter](../GOAL-010-a002-schema-adapter/00-meta.md)
