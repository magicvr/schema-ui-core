---
id: D-002
goal: GOAL-012-r3-s12-recycle-bin
title: 方案冻结：回收站设计（S1）
date: 2026-08-14
status: accepted
parent: GOAL-012-r3-s12-recycle-bin
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# D-002 · 方案冻结（S-12 回收站）

## 1. 数据模型

0025 `recycle_items`（admin.recycle-bin）：

- id TEXT PK（"recycle-" + 16 hex）
- resource TEXT NOT NULL（dict-types / dict-entries / scheduled-tasks）
- resource_id TEXT NOT NULL
- payload TEXT NOT NULL（删除前 Entity.Get 的行 JSON 映射）
- actor_id / actor_name TEXT NOT NULL
- deleted_at INTEGER NOT NULL
- restored_at INTEGER NULL
- 部分唯一索引：CREATE UNIQUE INDEX ... ON recycle_items(resource, resource_id) WHERE restored_at IS NULL

0026（core.operationlog）：CHECK + `recycle.restore`、`recycle.purge`。

## 2. 删除钩子（接入点）

- `handler.Resource.Trash TrashRecorder`（handler 包定义接口，模块 service 结构满足）。
- 工厂 delete()：Get(id) 捕获行 → Entity.Delete → 成功后 Trash.Record。
- 工厂 batchDelete()：非 BatchDeleter 顺序路径逐 id 同规则；BatchDeleter 路径（v1 无，预留）批量成功后逐条 Record。
- datadictionary / scheduledtasks provider 的 New 增加变参 `trash ...handler.TrashRecorder`，Register 时赋给两个 Resource（types/entries）或 tasks Resource；nil 不接入。
- 组合根：`recycleService` 先构造，同时传给模块 provider 与 recycle 路由。

## 3. 管理端点（recycle.read / recycle.write，PolicyAdmin）

- GET /api/recycle-bin（列表：resource/resourceId/actor/deletedAt/restoredAt；排序 deletedAt；q 过滤）
- GET /api/recycle-bin/{id}（详情含 payload）
- POST /api/recycle-bin/{id}/restore：按 resource 调用对应 store Create（payload→struct）；唯一键冲突 → 409 `RECYCLE_RESTORE_CONFLICT`（保留快照，可再恢复）；成功 → 写 restored_at + 审计 `recycle.restore`，返回恢复行。
- DELETE /api/recycle-bin/{id}：彻底清除（物理删快照行，不可逆）→ 204 + 审计 `recycle.purge`。

## 4. 页面与导航

- page `recycle-bin`：表格（resource / resourceId / deletedAt / actor / 状态）+ 行操作（恢复/彻底清除）+ 批量清除。
- menu_recycle_bin Order 8（admin-only，recycle.read）。

## 5. 权限/审计/错误码

- 权限键：recycle.read / recycle.write（PolicyAdmin）。
- 审计事件：recycle.restore / recycle.purge（0026 CHECK）。
- 错误码（DomainError / catalog / 冻结集）：RECYCLE_ITEM_NOT_FOUND、RECYCLE_RESTORE_CONFLICT。

## 6. 测试与验证

- store：快照插入/唯一/恢复/清除/列表。
- handler：删除钩子（类型/条目/任务删除后快照出现）、恢复（成功/冲突/不存在）、清除、权限 401/403、审计事件。
- 组合：admin 权限 22→24、导航 12→13；迁移 24→26。
- go 判定（S4）：删除语义仅在接入资源上新增快照副作用（record 失败不影响删除）；未接入资源零变化。

## 7. 未选方案

- 物理删除延迟（真软删所有表）：全局加 deleted_at 列破坏既有查询/唯一约束面，v1 否决。
- users/roles 快照：凭据/RBAC 不可还原，v1 排除（D-001 §2）。
- 文件字节快照：存储成本与一致性复杂，v1 排除。
