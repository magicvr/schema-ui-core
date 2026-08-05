---
id: E-006-r6-c62-migration-ownership
doc: execution-entry
goal: GOAL-013-r6-old-path-removal
source: orchestrator
date: 2026-08-05
status: recorded
---

# E-006 · R6 C6.2 Apply/DDL 物理迁出（切片 3）

## 已发生事实

- 0001-0008 的 descriptor、DDL 与 Apply 已从 `internal/store/migrate.go`
  迁至 owner 包：`modules/authsession/migration`（0001/0002）、
  `modules/corepersistence/migration`（0003/0006）、
  `modules/operationlog/migration`（0004/0005/0008）和
  `modules/settings/migration`（0007）。
- `modules/compiled.PersistenceCatalog()` 静态汇总所有 compiled persistence
  provider，经 `kernel.CollectPersistence` 校验、排序；composition `openStore`
  只经该 catalog 调用 `store.OpenWithCatalog`。Profile 启用集合不筛掉全局迁移历史。
- `store/migrate.go` 仅保留平台 runner、ledger、连续性/漂移/未知迁移校验、事务、
  升级前快照与完整性检查；生产代码中已无内建 `compiledMigrations`、
  `MigrationCatalog`、模块 DDL 或 `store.Open` 旁路。测试兼容入口位于 `_test.go`。
- `kernel.MigrationChecksum` 成为共享 checksum 实现；回归测试逐项固定 0001-0008
  的 version、moduleID、key/name 与 checksum，并新增已应用 name drift fail-closed。
  冻结设计要求 0001-0008 不改 Apply 语义，因此历史 `records_retire` 仍执行既有
  DROP/清理，而不是把既有迁移改写成 no-op。
- 当前用户/角色写路径仍需的 `roleKeyRe`/`linkUserRole` 作为运行时仓储辅助保留在
  `store/role_links.go`；0002 backfill 使用 owner 迁移包内的封闭实现，迁移 owner 不再
  反向依赖 store。

## 验证

- `apps/api: go test ./...`：通过。
- `apps/api: go vet ./...`：通过。
- `git diff --check`：通过。
- 生产旧路径扫描（排除 `*_test.go`）：`compiledMigrations`、
  `MigrationCatalog`、`SiteSettingsDDL`、`storePersistenceProvider`、
  `store.Open(` 均零命中。
- fresh DB、ledger-less R2 升级、checksum drift、version gap、unknown applied
  migration、restart、pre-v0002 snapshot restore 与 MVP recovery 由现有 API 全量测试覆盖。

## 仍开放

- C6.2 尚有 F-005：fresh admin bootstrap 与 finalized Authorization/Navigation
  contribution 驱动的 versioned system-data reconcile 分离；不得覆盖用户拥有字段。
- 因 F-C62-004 仍开放，本条不勾选 C6.2，不关闭 Root A-010，也不作为 VP 退出证据。
