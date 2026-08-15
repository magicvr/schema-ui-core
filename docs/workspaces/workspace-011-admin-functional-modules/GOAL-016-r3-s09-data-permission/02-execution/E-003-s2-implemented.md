---
id: E-003
goal: GOAL-016-r3-s09-data-permission
title: S2 实现完成（数据权限）
date: 2026-08-15
status: recorded
parent: GOAL-016-r3-s09-data-permission
created: 2026-08-15
updated: 2026-08-15
version: 1.0.0
---

# E-003 · S2 实现完成（2026-08-15）

## 事实

- **migration 0027**（admin.data-permission）：data_scope_policies + user_data_scopes（owner_column / default_scope CHECK all|self / enabled；PK(user_id, resource)）；checksum f3ce4c71…。
- **migration 0028**（core.operationlog）：CHECK + data-permission.policy-update / data-permission.scope-update；checksum a18c42e7…；compiled/persistence.go 注册。
- **store**：modules/datapermission/store/repository.go（Policy/Assignment CRUD、ErrNotFound/ErrInvalidScope/ErrNotEnforceable、Upsert 幂等）。
- **service**：modules/datapermission/provider.go Service——ScopeFor（enforceable 集合 + policy + assignment 合成，nil=无约束）、UpsertPolicy/UpsertAssignments（强制点在策略写入，A-005）。
- **handler**：handler/datapermission.go DataPermissionRoutes（GET/PATCH policies、GET/PATCH scopes；data-permission.read/write 门禁；default_scope 必填；审计事件）；errorcatalog 新增 INVALID_SCOPE / SCOPE_NOT_ENFORCEABLE。
- **工厂执行（resources.go）**：ScopeConstraint / RowScopeProvider / Resource.Scoper / resourceFilter.Scope——list 下传、detail/update/delete 404 归属检查、batch-delete 仅删本人行、create 强制 owner（A-004 F-001 全部行访问路径覆盖）。
- **装配**：kernel/profile.go（ProfileAdmin + BuiltinModules）、kernel/provider.go（DefaultNavigationOrder 尾部 menu_data_permission）、composition.go（admin.data-permission 接线，v1 enforceable=nil）、testsupport 镜像。
- **web**：data-permission 页 schema + fragment + i18n zh/en（16 键）；fixture app-manifest.admin.json + 页/导航 + STATIC_MANIFEST_SHA256 重钉 4a43bcad…；smoke.sh admin 页面集 + data-permission；schema-keys/s5-denominator/app-manifest.test 清单。
- **测试**：provider_test（注册面 + ScopeFor 合成 + 强制点）、store repository_test、handler datapermission_test（401/403/生命周期/审计）、resources_test（self 作用域 list/detail/update/delete/batch/create + nil scoper 字节不变）；迁移计数快照 26→28。
