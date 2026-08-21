---
id: E-002
goal: GOAL-012-r3-s12-recycle-bin
date: 2026-08-14
status: recorded
parent: GOAL-012-r3-s12-recycle-bin
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# E-002 · S2 实现完成

## 事实

- handler/resources.go：Resource 增可选 `Trash TrashRecorder` 字段；delete()/batchDelete() 在删除**前** Get 捕获行、删除**成功后** Record（失败不落快照；Record 失败仅 slog）。
- handler/recyclebin.go（新增）：GET /api/recycle-bin、GET /{id}、POST /{id}/restore、DELETE /{id}（recycle.read/write；审计 recycle.restore/recycle.purge）。
- modules/recyclebin：migration 0025（recycle_items：payload JSON、部分唯一 active 索引）、store/repository.go（Record/List/Get/MarkRestored/Purge）、service.go（TrashRecorder + 按资源 restore 派发 dict-types/dict-entries/scheduled-tasks，冲突→RECYCLE_RESTORE_CONFLICT）、provider.go（4 路由/1 页/2 权限/menu_recycle_bin Order 8/fragment）、schema/recycle-bin.json、manifest/fragment.json。
- operationlog 0026（CHECK + recycle.restore/recycle.purge）+ 事件常量。
- 接线：dictionary/scheduledtasks provider New 增变参 trash；composition 先构造 recycle 服务再注入；profile/admin + compiled/persistence + testsupport 镜像；errorcatalog + 冻结集（RECYCLE_ITEM_NOT_FOUND / RECYCLE_RESTORE_CONFLICT）；composition_test admin 24 权限/13 导航；迁移计数 24→26（0025=39087fe4…、0026=681f3bdc…）。
- 测试：store/repository_test、service_test（Record/恢复/冲突/清除）、handler/recyclebin_test（列表/恢复/清除/权限/审计）。
- Web：fixture +recycle-bin 页/导航，sha 重钉 9649fbfa…，i18n +10 键，schema-keys/s5 列表 + smoke + shell.spec。
