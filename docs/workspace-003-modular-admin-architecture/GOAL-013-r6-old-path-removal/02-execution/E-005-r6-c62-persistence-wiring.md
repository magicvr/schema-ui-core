---
id: E-005-r6-c62-persistence-wiring
doc: execution-entry
goal: GOAL-013-r6-old-path-removal
source: orchestrator
date: 2026-08-05
status: recorded
---

# E-005 · R6 C6.2 Persistence 接线（切片 1-2）

## 已发生事实

- **切片 1（ownership 登记）**：0001-0008 各迁移获正确 moduleID 归属
  （core.auth-session/0001+0002、core.operationlog/0004+0005+0008、
  admin.settings/0007、core.persistence/0003+0006）；`store.MigrationCatalog()`
  暴露为 `kernel.MigrationContribution`（Checksum = 账本 checksum，Apply 包装内部
  up），经 `kernel.CollectPersistence` 校验通过。
- **切片 2（CollectPersistence 生产接线）**：composition `openStore` 经
  `kernel.CollectPersistence(providers)` 收集 catalog（`storePersistenceProvider`
  承载历史 0001-0008）→ `store.OpenWithCatalog` 校验 catalog 与权威账本精确一致
  （fail closed on mismatch）后打开。新增 `TestMigrationCatalogMatchesLedgerWithOwnership`、
  `TestOpenWithCatalogRejectsDivergentCatalog`。
- `store.Open` 重构为共享 `open`；`OpenWithCatalog` 为其加 catalog 校验入口。

## 验证

API `go test ./...`（14 包）+ `go vet` + Web `vitest run`（495）通过。

## C6.2 剩余

- **切片 3（F-002 物理迁出）**：将 0001-0008 Apply/DDL 从 `store/migrate.go` 迁入
  core.auth-session / core.operationlog / admin.settings 模块包，各模块
  `CompiledPersistence()` 返回自有迁移；store 收窄为平台 runner（消费 catalog）。
- **F-005**：seed/RBAC reconcile 改 Authorization/system-data 贡献驱动。
- **C6.3（F-003b）**：Schema 字节 ContributionSet 发布。
