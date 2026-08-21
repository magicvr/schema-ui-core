---
id: D-001
goal: GOAL-012-r3-s12-recycle-bin
title: 目标边界与信息就绪（S1）
date: 2026-08-14
status: accepted
parent: GOAL-012-r3-s12-recycle-bin
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# D-001 · 目标边界与信息就绪（S-12 回收站/软删除）

## 1. 目标边界

- 交付：回收站（软删除管理）模块 `admin.recycle-bin`——删除快照 + 管理 UI（浏览/恢复/彻底清除）。
- 新持久化（不复用 0006 records_retire 演示实体退场）：`recycle_items` 快照表（0025）+ operationlog CHECK 扩展（0026）。
- Profile：admin 默认集内容扩展（mvp/demo 不变）。

## 2. 受管资源范围（I-001）

v1 受管资源（快照 + 恢复）：

| resource | 来源 | 恢复方式 |
|----------|------|----------|
| `dict-types` | admin.data-dictionary 类型删除 | store.CreateType |
| `dict-entries` | admin.data-dictionary 条目删除 | store.CreateEntry |
| `scheduled-tasks` | admin.scheduled-tasks 任务删除 | store.CreateTask |

排除（文档化残余）：

- `users` / `roles`：快照不含 password_hash（用户行映射为 write-only），恢复无法还原凭据与 RBAC 不变量；v1 不做。
- `file-library`：物理文件字节不随行快照；恢复语义需字节保留，v1 不做。
- `notifications`：瞬时通知按设计不回收。

## 3. 快照模型

- 表 `recycle_items`：id / resource / resource_id / payload(JSON 行映射) / actor_id / actor_name / deleted_at / restored_at(NULL=未恢复)。
- 部分唯一索引：UNIQUE(resource, resource_id) WHERE restored_at IS NULL——同资源同 id 未恢复前不重复快照。
- 彻底清除 = 物理删除快照行（不可逆，审计 recycle.purge）。

## 4. 删除接入

- handler.Resource 新增可选 `Trash TrashRecorder`（Record(ctx, resource, id, row, actor, now)）。
- 工厂 delete()/batchDelete()：删除**前** Entity.Get 捕获行，删除**成功后**才 Record（失败不落快照）。
- 接线：composition 先构造 recycle 服务，以变参传给 datadictionary / scheduledtasks 模块 provider（nil = 不接入，既有语义不变）。

## 5. 信息就绪

| ID | 级别 | 结论 |
|----|------|------|
| I-001 | required | 受管资源范围冻结（§2；payload 可恢复性经 store Create 核对） |
| I-002 | required | 恢复/清除语义冻结（§3/§4；冲突→409 RECYCLE_RESTORE_CONFLICT；清除不可逆） |
| I-003 | required | Profile admin 默认集（内容扩展，与 S-01..S-11 先例一致） |
