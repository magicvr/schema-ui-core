---
id: E-015
doc: execution-entry
goal: GOAL-001-store-dialects
status: recorded
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
---

# E-015 · 真实本地测试 PG 连接验证（192.168.31.213）

## 2026-08-20 · 用户指定 PG 复验通过（关门前提确认）

### 已发生事实

- 用户提供真实本地测试 PG：`192.168.31.213:5432`，用户 `sa`。
- **连通性/权限**：TCP 可达；`CONNECT_OK server=15.4 user=sa db=postgres superuser=on createdb=true`。
- **全量 live 门控套件在真实 PG 上复跑全绿**（`SCHEMA_UI_R2_PG_DSN=postgres://sa:***@192.168.31.213:5432/postgres?sslmode=disable`，`go test ./...` **0 FAIL**）：
  - `TestFullCatalogPostgresBootstrapIntegration`（48 迁移 fresh bootstrap + 台账幂等 + BIGINT 合规）PASS
  - `TestCompositionPostgresStartup`（postgres DSN 完整应用启动，readyz 门禁）PASS
  - `TestPostgresCrossModuleSharedTx`（跨模块共事务 commit/rollback）PASS
  - `TestPostgresDataMigrationPrototype`（sqlite→PG 数据迁移原型）PASS
  - `TestAuthsessionPostgresApplyIntegration` / `TestPostgresMigrateRunnerIntegration` / `TestOpenPostgres*` PASS
- **备份/恢复 round-trip 在真实库验证**：`pg_dump -F c` → `pg_restore` → count=2 ✓（PG 15.4，容器 pg_dump 17 客户端兼容）。
- **清理确认**：测试自建 scratch 库（r3full/r5tx/r5mig/r4s4/r3auth/r3open/r5u2*）全部随测试 cleanup 删除；远端仅剩既有 `ironclaw/postgres/sa/sub2api`；`postgres` 库无残留表。

### 证据

| 主张 | 路径 / 命令 |
|------|-------------|
| 连接+角色 | 一次性 `cmd/pgcheck` 探测（已移除）；`CONNECT_OK server=15.4 superuser=on createdb=true` |
| 全量回归（真实 PG） | `apps/api`: `SCHEMA_UI_R2_PG_DSN=postgres://sa@192.168.31.213/... go test ./...` → 0 FAIL |
| 备份 round-trip | 容器 `pg_dump/pg_restore -h 192.168.31.213 -U sa`（r5u2sarest count=2）后已清理 |
| 无残留 | `pg_database` 清单 + `pg_tables` 空查询 |
