---
id: A-002
doc: audit-entry
goal: GOAL-004-r3-dual-dialect-ledger
source: self
scope: R3 T2a 切片（postgres 迁移运行器 + live PG 证明）
verdict: pass
status: recorded
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
---

# A-002 · R3 T2a 切片自审（source: self）

## 范围与区间

- auditor: 本会话编排器（self）
- type: `stage`
- covered: `(*postgres).migrate` / `applyMigrationPG` / `appliedMigrationsPG`；live PG 集成测试（apply、幂等、drift fail-closed）；与共用 live DB 的测试解耦；全量回归
- excluded: 生产解闸（`openPostgres` 非空 catalog）+ composition postgres 路由（并入 T3）；真实双方言 catalog 对写（T3）；PG 备份/完整性合同（R5）

## 成果与证据

| 主张 | 证据 |
|------|------|
| 运行器按 kernel.Tx/rebind 跑通 | `postgres.go`；`TestPostgresMigrateRunnerIntegration`（live postgres:17-alpine） |
| 台账 checksum 绑 sqlite 历史；drift fail-closed | 集成测试 drift 断言（错误含 `drift`） |
| 幂等（重开 + 再迁移） | 集成测试（`-count=2` 复跑绿） |
| 全量回归 0 FAIL | `apps/api`: `go test ./...`；commit `4359c7f` |

## 对照成功标准（T2 相关切面）

| 标准 | 状态 | 证据 |
|------|------|------|
| postgres fresh bootstrap 可 apply 全部迁移（本拍 = 便携 scratch catalog） | ✅（运行器） | live 集成测试 |
| `schema_migrations` 台账含同名 + 同 checksum | ✅ | 集成测试（checksum = `kernel.MigrationChecksum` 绑 sqlite 文本） |
| checksum fail-closed | ✅ | drift 断言 |

## Findings

### F-001 · 生产路径未解闸（recommended）

| 字段 | 值 |
|------|-----|
| level | recommended |
| status | open |
| evidence | `postgres.go` `openPostgres` 仍对非空 catalog fail-closed；`composition.openStore` 仍以 sqlite 分支工作 |
| closure | 无（计划并入 T3） |

描述：`postgres.migrate` 已被 live PG 证明，但生产打开路径（postgres + 非空 compiled catalog）尚未解闸——需等 T3 真实双方言 catalog（当前 catalog 含 sqlite-only SQL）。不阻断本切片；T3 随真实对写一并解闸并记录证据。

## 必改项汇总（required）

无。

## 结论与下一步

T2a `pass`（0 required；F-001 recommended 计划并入 T3）。下一步 T3：逐迁移对写（authsession `sqlite_master`/`PRAGMA`/`COLLATE NOCASE`、时间 `BIGINT`、非时间宽度/布尔、`?` 可 rebind 文本）+ 生产解闸 + 双 apply 证据；T3/T4 迁移/数据门禁走 self + independent（grok build）。
