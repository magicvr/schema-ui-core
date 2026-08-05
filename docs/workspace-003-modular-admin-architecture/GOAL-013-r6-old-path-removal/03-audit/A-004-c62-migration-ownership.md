---
id: A-004-c62-migration-ownership
doc: audit-entry
goal: GOAL-013-r6-old-path-removal
source: self
date: 2026-08-05
scope: C6.2 slice 3 - 0001-0008 Apply/DDL physical ownership and platform runner boundary
verdict: pass
---

# A-004 · C6.2 切片 3 物理迁出自审

- **source**：self
- **auditor**：Codex
- **类型 / scope**：stage / execution-facts；C6.2 切片 3（F-002）
- **verdict**：pass（仅本切片；C6.2 整体未放行）

## 范围与区间

核验 E-006 实现：0001-0008 的 owner 物理归属、compiled-global catalog 生产执行链、
store runner/ledger 边界，以及冻结迁移 identity/checksum/Apply 与升级恢复行为。

## 成果（有证据）

| 标准 | 状态 | 证据 |
|------|------|------|
| owner 包持有 descriptor/DDL/Apply | pass | `apps/api/internal/modules/{authsession,corepersistence,operationlog,settings}/migration/` |
| production catalog 唯一执行链 | pass | `modules/compiled/persistence.go` → `composition/openStore` → `store.OpenWithCatalog` |
| store 无生产模块 registry/DDL/旁路 | pass | `store/migrate.go` + E-006 零命中扫描；兼容入口仅 `migration_catalog_test.go` |
| 0001-0008 不重编号/改名/checksum/Apply | pass | `TestCompiledMigrationCatalogOwnership` 固定八条 identity/checksum；fresh/R2/recovery 测试通过 |
| drift/gap/unknown fail closed | pass | store/kernel migration tests + `TestOpenWithCatalogRejectsInvalidAndAppliedDrift` |
| 提交前静态与构建验证 | pass | E-006：`go test ./...`、`go vet ./...`、`git diff --check` |

## Findings

本切片未发现新的 required 或 recommended finding。

## 必改项汇总

- 本 scope 新增 required：0。
- 继承 `F-C62-004`：open。F-005 contribution-driven system-data reconcile 尚未实现，
  因此 C6.2 整体不能勾选，也不能关闭 Root A-010。

## 结论与下一步

切片 3 达到 D-002/A-003 冻结边界，允许进入 F-005。下一步分离 fresh bootstrap 与
versioned reconcile，并以 finalized Authorization/Navigation contributions 为唯一系统
permissions/menu/grants 来源；完成后执行 C6.2 independent audit。
