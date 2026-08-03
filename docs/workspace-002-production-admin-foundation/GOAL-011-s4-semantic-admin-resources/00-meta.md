---
id: GOAL-011-s4-semantic-admin-resources
title: S4 · 语义化 Admin 资源替换与双实体验证
status: done
created: 2026-08-03
updated: 2026-08-04
parent: GOAL-010-a002-schema-adapter
version: 0.11.0
progress: 5/5
---

# GOAL-011 · S4 · 语义化 Admin 资源替换与双实体验证

## 概述

承接父目标 GOAL-010 的 S4：将当前仅用于代表性 CRUD 验证、缺少真实业务语义的 `records` 从产品默认运行面退场，以 `users` 替换其代表实体，并新增 `roles` 作为第二个语义资源。目标是在不修改前端 Renderer 主路径的前提下，让两个对绝大多数 Admin 系统均有实际价值的资源完成 Schema 驱动的标准 CRUD 闭环，同时保留既有迁移与历史治理记录的可追溯性。

“只修改 Schema 接入”限定为：后端资源契约、持久化、权限与领域不变量已完成注册后，前端新增页面只修改 Schema/fixture，不修改 Renderer 主路径；不主张仅凭 Schema 自动生成后端领域逻辑。

## 路线图 / 成功标准

五个检查点默认等权、原则上串行（P-001）。原 S1～S5 实施事实均已完成；A-012 新增 required finding 后，本目标曾按 P-003 fail closed 恢复为 `active / 5/5`。D-006 选择五项全部走 fixed，候选提交 `fb5cd06` 完成整改与验收，A-013（independent · finding-closure · pass）逐项确认 F-001～F-005 fixed、无新增 required；D-007 据此恢复 **`done / 5/5`**：

- [x] **S1 · 语义资源与退场契约冻结**：关闭 `I-011-001`/`I-011-002`，冻结 users/roles 的最小领域边界、安全不变量，以及 records API/表/权限/菜单/操作日志/测试的版本化退场与既有数据库迁移策略。
- [x] **S2 · users/roles 后端资源闭环**：在通用资源工厂之上实现 users 与 roles 的持久化 list/search/detail/create/update/delete、字段校验、敏感字段隔离、关系/系统角色保护、稳定错误 envelope 与 401/403 负向路径。
- [x] **S3 · records 产品运行面退场**：按冻结策略移除当前产品默认运行面中的 records API 注册、Schema fixture、菜单/权限、种子、操作日志耦合、前端 records 专名与当前测试依赖；保留不可改写的历史治理事实和迁移链证据。
- [x] **S4 · 双语义实体 Schema 接入验证**：关闭 `I-011-003`；users 与 roles 的代表性列表/CRUD 页面由 Schema 接入，前端 Renderer 主路径无修改；fresh fork、既有数据库升级、重启持久化与权限失败路径可复核。
- [x] **S5 · 回归、审计与父级交接**：API/Web/build/E2E 全量回归与阶段/关门审计通过；向 GOAL-010 提供其 S4 可勾选的证据，目标自身无开放 required finding 或到期 required 信息项。

## 派生进度

`progress` 由 S1～S5 五个检查点等权派生；当前为 `5/5`。检查点不替代审计、信息门禁或关门结论。

## A-012 finding closure

A-012（independent · fail）新增的 F-001～F-005 五条 required finding 已按 D-006 / I-011-004 v0.1.0 完成实现与验证，并由 A-013（independent · pass）逐项确认 **fixed**；当前无开放 required。F-006 保持 `recommended / open / non-blocking`，不因关门被静默关闭。

`progress: 5/5` 仍只由原五个显式检查点派生；finding closure 另由候选提交 `fb5cd06`、执行收据、A-013 与 D-007 建立，不由 progress 推导。远端 GitHub-hosted Actions 尚未触发；本目标只主张已记录的本地 Windows/Go/Node/Chromium/Linux-container 证据，不冒充远端 CI acceptance。

## 信息需求与阶段门禁

| ID | 问题 / 所需信息 | 级别 | 影响门禁 | 最晚阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据或结论 |
|----|-----------------|------|----------|----------|-----------------|------|-------------|------------|
| `I-011-001` | users/roles 的精确资源契约与最小 IAM 边界是什么？包括公开字段、密码/令牌隔离、角色分配、self/最后管理员保护、system role 与 grant 约束。 | required | S1 方案冻结与 S2 实施 | S2 首个产品代码变更前 | 核对现有 auth/RBAC 表、store、权限投影与通用 Resource 接口，形成版本化领域契约并提交用户裁决 | **verified** | 不延期；S1 已冻结 | [I-011-001-users-roles-contract.md](attachments/I-011-001-users-roles-contract.md) **v0.2.0** + GOAL-011 D-002/D-003（用户裁决：通用工厂+最小契约扩展、操作日志纳入；A-002 F-001/F-003/F-004 fixed） |
| `I-011-002` | records 如何从当前产品基线退场，同时保持既有迁移/checksum、已有数据库、数据处置、API/权限/菜单/operation log 与历史文档可追溯？ | required | S1 方案冻结与 S3 退场实施 | S3 首个删除、重命名或迁移变更前 | 建立 fresh install 与 in-place upgrade 两条迁移矩阵，逐项盘点运行时代码、迁移、种子、文档与测试 | **verified** | 不延期；S1 已冻结 | [I-011-002-records-retirement.md](attachments/I-011-002-records-retirement.md) **v0.2.0** + GOAL-011 D-002/D-003（用户裁决：硬退场 DROP TABLE；A-002 F-002 快照语义 fixed） |
| `I-011-003` | 双资源验收如何证明“前端 Schema-only 接入”以及 fresh fork/升级/重启/401/403 的完整边界？ | required | S4 集成验收与 S5 关门 | S4 首个验收实现前 | 冻结 users/roles 两套验收矩阵、Renderer diff 边界、数据库升级夹具与 E2E/restart 协议 | **verified** | 不延期；D-004 已冻结，S4 执行时按矩阵复核 | [I-011-003-acceptance-matrix.md](attachments/I-011-003-acceptance-matrix.md) **v0.2.0** + D-004；A-007 F-001～F-003 均 `fixed`，基线提交、双资源页面级证据与后端完整断言已补齐 |
| `I-011-004` | A-012 F-003 的产品边界如何裁决：保留 seed-only 最小 IAM 并移除/隐藏无有效 grant 语义的 roles 管理面，还是补齐用户角色分配、可验证 grant 来源与 roles 管理流程？ | required | A-012 整改方案冻结与重新关门 | F-003 首个产品代码变更前；最晚重新关门前 | 用户按 P-004 选择边界；冻结授权委派、密码、grant、records 洁净度与最后管理员原子性的版本化整改契约/验收矩阵 | **verified** | 不延期；D-006 已冻结；实施与复审按契约逐项核对 | [I-011-004-a012-remediation-contract.md](attachments/I-011-004-a012-remediation-contract.md) **v0.1.0** + D-006；用户选择补齐角色授权/grant 管理路径，F-001～F-005 全部走 fixed |

`I-011-001`～`I-011-004` 均维持 `verified`。I-011-004 的信息状态只表示整改边界已冻结；F-001～F-005 的实现与独立 closure 另由 `fb5cd06`、A-013 和 D-007 证明。当前无到期 required 信息项或开放 required finding；历史 S1～S5 事实不回退。

## 依赖与边界

| 项 | 说明 |
|----|------|
| 父目标 | [GOAL-010-a002-schema-adapter](../GOAL-010-a002-schema-adapter/00-meta.md)（S4 交付载体；Root A-002 F-002-001 仍 open） |
| In | 最小 users/roles 管理资源、records 当前产品运行面退场、版本化迁移、权限/菜单/操作日志适配、双资源 Schema 页面与回归证据 |
| Out | 完整 IAM 产品、SSO/SCIM、复杂策略编排、多租户、扩大 `I-PROTO-001 v0.1.3`、Root/VP-002 关门 |
| 历史边界 | 不重写 GOAL-004/007 的历史决策、执行与审计；既有迁移仅按后续冻结的兼容策略演进 |

## 父目标

- [GOAL-010-a002-schema-adapter](../GOAL-010-a002-schema-adapter/00-meta.md)
