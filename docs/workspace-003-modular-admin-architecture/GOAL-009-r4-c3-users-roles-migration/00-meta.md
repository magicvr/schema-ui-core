---
id: GOAL-009-r4-c3-users-roles-migration
title: R4-C3 · Users 与 Roles 迁移
status: active
parent: GOAL-005-r4-full-module-migration
created: 2026-08-05
updated: 2026-08-05
version: 0.1.0
progress: 0/4
plan_refs:
  - VP-003-modular-admin-architecture
primary_plan: VP-003-modular-admin-architecture
serves_summary: 承接冻结包 §7 切换顺序，将 Users/Roles 从中心注册迁入 GOAL-008 的 Provider 契约，清除中心 Register/Schema owner map/Manifest adminModules 特例，保留现有 CRUD、授权、角色分配、最后管理员保护、密码与 operationlog 行为。
---

# GOAL-009 · R4-C3 Users 与 Roles 迁移

## 概述

本子目标是 `GOAL-005-r4-full-module-migration` 的 C3 子目标，承接冻结包 §7 的
切换顺序与 GOAL-008 已落库的 Provider/Registrar 契约，把 `admin.users` 与
`admin.roles` 从中心 `handler.Register` / Schema owner map / Manifest `adminModules`
特例迁入 provider 结构化贡献。C3 保留现有 HTTP/Schema/授权/持久化/日志行为矩阵，
不改变最后管理员保护、角色分配、密码与 operationlog 语义。C4（Schema 其他能力）
与 C5 验收不在本目标。

## 愿景对齐

| 字段 | 值 |
|------|----|
| `parent` | `GOAL-005-r4-full-module-migration` |
| `plan_refs` | `VP-003-modular-admin-architecture` |
| `primary_plan` | `VP-003-modular-admin-architecture` |
| Charter | `schema-ui-core-admin-foundation@0.2.0`（经 VP-003 间接对齐） |
| 审计模式 | `independent`；迁移切片使用 Grok Build `grok-4.5` / `high` |

## 成功标准

- [ ] **C3.1 / 迁移扫描与行为矩阵**：Users/Roles 当前中心注册、Schema fixture、
  Manifest 投影、seed/权限与 operationlog writer 的状态盘点；固定保留行为矩阵
  （CRUD、授权、角色分配、最后管理员保护、密码、operationlog best-effort）。
- [ ] **C3.2 / Provider 化**：`admin.users`/`admin.roles` 提供 provider（Descriptor +
  `Register` 写 HTTP/Schema/Auth/Nav/Manifest；`CompiledPersistence` 无新增迁移），
  与现有中心输出做兼容比较（测试内，不永久双注册）。
- [ ] **C3.3 / 中心特例清除**：composition 消费 provider finalize；移除中心
  Users/Roles 分支、Schema owner map 与 Manifest `adminModules` 相关特例；无永久
  双路径。
- [ ] **C3.4 / 验证与关门**：行为矩阵测试 + 双 Profile + 失败注入（operationlog
  append 失败不翻转业务成功）+ self + Grok independent 无开放 required finding。

四个检查点等权；当前 `progress: 0/4`。完成本子目标只表示 C3 关闭，不关闭 GOAL-005、
Root 或 VP-003，不自动放行 C4。

## 信息门禁

| 编号 | 级别 | 必须回答的问题 | 影响 | 最晚阶段 | 收集动作 | 状态 | 证据 |
|------|------|----------------|------|----------|----------|------|------|
| C3-I001 | required | Users/Roles 当前中心注册、Schema fixture、Manifest 投影与 seed/权限 ownership 的完整状态？ | C3.1/C3.3 | C3.1 | 全仓扫描 `handler.Register`、`schema.go` owner map、`manifest.go` adminModules、`seed.go` | collecting | 待 C3.1 扫描 |
| C3-I002 | required | C3 必须保留的行为矩阵（HTTP/Schema/授权/角色分配/最后管理员/密码/operationlog）是否枚举？ | C3.2/C3.4 | C3.1 | 对照冻结包 §7 兼容清单 + 现有测试 | collecting | 冻结包 §7；待 C3.1 枚举 |
| C3-I003 | required | operationlog append 失败不翻转业务成功的 failure-injection 测试（FR-005/C13-003）是否补齐？ | C3.4 行为矩阵 | C3.4 | Users/Roles/Auth/Settings 失败注入测试 | collecting | GOAL-006 FR-005；C3.4 补测 |
| C3-I004 | non-blocking | 运行时双 Profile 矩阵与 Manifest secrecy 是否纳入 C3 门禁？ | C3.4 证据强度 | C3.4 | GOAL-008 E-004 登记的 C3 门禁 | open | GOAL-008 E-004 |

## 阶段路线图

1. 迁移扫描：中心注册/Schema/Manifest/seed ownership 状态 + 行为矩阵枚举（C3.1）。
2. Users/Roles provider 化 + 与中心输出兼容比较测试（C3.2）。
3. 切换 composition 消费 provider finalize；清除中心特例（C3.3）。
4. 行为矩阵 + 双 Profile + 失败注入 + self + Grok independent 复审（C3.4）。

## 范围与非目标

范围包括 admin.users/admin.roles 的 provider 化、中心特例清除、行为矩阵验证、
operationlog 失败注入、双 Profile 与 Manifest secrecy 门禁。非目标包括 Schema 其他
能力迁移（C4）、Records 恢复、Settings/Activity 的 provider 化迁移（C3 可延后，
C4 承接）、R5/R6。`0003`/`0006` 迁移账本与历史 operation-log 保留。
