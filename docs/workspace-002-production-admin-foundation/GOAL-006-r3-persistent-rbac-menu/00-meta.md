---
id: GOAL-006-r3-persistent-rbac-menu
title: R3 · 持久化 RBAC、菜单投影与版本迁移
status: active
created: 2026-08-02
updated: 2026-08-02
parent: GOAL-001-production-admin-foundation
version: 0.4.0
progress: 2/6
---

# GOAL-006 · R3 · 持久化 RBAC、菜单投影与版本迁移

## 概述

承接 Root **D-009**：以规范化 RBAC、稳定 permission key、`features` 菜单投影和两步 SQLite 迁移，交付 R3 的用户—角色—权限—菜单最小持久化闭环。保留 R2 的 `account.User {id,name,roles}`、JWT subject 与 refresh-token 关系，不把前端菜单隐藏当作后端授权。

## 成功标准

- [x] **S1 · 版本迁移与可恢复起点**：建立 `schema_migrations`、顺序/校验和检查、事务化 fail-closed 迁移，并在升级前产生可恢复数据库副本。
- [x] **S2 · 规范化 RBAC 与两步兼容**：建立角色/权限/用户关联/菜单 grant 表，回填并双读核对 `users.roles`；切换规范化读写后仍保持 R2 身份与 refresh 契约，旧列不在本步删除。
- [ ] **S3 · 增量幂等种子**：按稳定 key ensure admin、viewer、`records.read`、`records.write`、代表性菜单项及 grants；重复启动无重复关系，不因已有任意用户而跳过关系修复。
- [ ] **S4 · 后端读写授权**：records 读写均经认证和 permission gate；admin 可读写，viewer 可读不可写，匿名读写 `401`，已认证缺权限 `403`。
- [ ] **S5 · features 菜单投影**：`/api/accounts/me.features` 从持久化菜单 grants 生成布尔投影；真实 manifest 至少一个条目使用 `visibleWhen`，验证 admin/viewer 可见性与空组剪枝。
- [ ] **S6 · 恢复、重启与回归证据**：验证迁移前副本恢复、迁移后 `PRAGMA integrity_check`、身份/授权/菜单/refresh 关键查询、服务重启持久化，以及 API/Web 既有测试、构建与 vet 回归。

## 派生进度

`progress: 2/6` 由上方六个顺序检查点等权派生。检查点只在实现事实与对应验证齐备后勾选；不得用计划、局部测试或百分比替代门禁和审计。

## 信息需求

| ID | 问题 / 所需信息 | 级别 | 影响门禁 | 最晚阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据或结论 |
|----|-----------------|------|----------|----------|-----------------|------|-------------|------------|
| `I-006-001` | 精确 DDL、迁移版本/校验和、FK/unique/delete 语义及稳定 seed keys 是什么？ | required | S1/S2/S3 实施 | 首个代码变更前 | 对照现有 store 与 D-009，形成版本化 schema/seed 计划并记录子目标决策 | **verified** | 已关闭（D-002） | [I-006-001-schema-migration-plan.md](attachments/I-006-001-schema-migration-plan.md)：`0001` 基线登记 + `0002` RBAC 扩展、DDL/FK/delete、双读切换、seed 与恢复/测试矩阵 |
| `I-006-002` | 首个受控真实 `page_ref`、显式 `feature_key` 与 admin/viewer grant 矩阵是什么？ | required | S5 实施与验收 | Web 投影实现前 | 对照真实 manifest 与导航测试，选择最小代表项并记录正反矩阵 | **verified** | 已关闭（D-003） | [I-006-002-menu-projection-matrix.md](attachments/I-006-002-menu-projection-matrix.md)：`list-edit-lifecycle` / `menu_list_edit_lifecycle`；admin=true，viewer/editor=false |

## 依赖与边界

| 项 | 说明 |
|----|------|
| 父目标 | [GOAL-001-production-admin-foundation](../GOAL-001-production-admin-foundation/00-meta.md)（D-009；Root `I-003` = verified） |
| 前置证据 | [I-003-persistence-permission-collection.md](../GOAL-001-production-admin-foundation/attachments/I-003-persistence-permission-collection.md)（当前事实、方案比较、M-R3-01～12） |
| In | SQLite 版本迁移、规范化 RBAC、增量种子、records 读写 permission、`features` 菜单投影、迁移恢复与回归证据 |
| Out | 通用 IAM/策略表达式；R4 Schema CRUD 扩域；R5 完整生产备份/部署运维；删除旧 `users.roles` 的不可逆迁移 |

## 父目标

- [GOAL-001-production-admin-foundation](../GOAL-001-production-admin-foundation/00-meta.md)
