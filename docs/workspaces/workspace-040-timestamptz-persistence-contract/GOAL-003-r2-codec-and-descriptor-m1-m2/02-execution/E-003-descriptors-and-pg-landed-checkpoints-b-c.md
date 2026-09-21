---
id: E-003-descriptors-and-pg-landed-checkpoints-b-c
doc: execution-entry
status: active
parent: GOAL-003-r2-codec-and-descriptor-m1-m2
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-003 · 15 个 descriptor 落码（checkpoint B/C）

## 事实

- **语句来源（机械导出）**：`CREATE TABLE` 正文与列序取自 **v72 已 apply 库的 live `sqlite_master` + `PRAGMA table_info`**，只施加冻结的逐列编辑（`INTEGER` → `TEXT`；D0/voucher 列去 `NOT NULL` / `DEFAULT 0`）。生成器产出 `modules/*/migration/vp040_temporal.go`（15 个文件），并逐表在临时库中 `CREATE` 后比对 `PRAGMA table_info`（列名/cid 序/类型/notnull）作为生成期自检。
- **F-5 先决条件**：FK 子表关系由 v72 DDL 的 `REFERENCES <parent>(` **机械解析**得出（父表：`users`、`roles`、`permissions`、`menu_items`、`dict_types`、`scheduled_tasks`、`operation_log`）；其余表无子表 → 裸四步。该盘点同时充当逐表附件 §5.8 要求的 v85 子表盘点证据（wallet/subject/voucher 六表 **count = 0**）。
- **语句顺序**：m0 预检（sentinel 列）→ m1–m3（子表 TEMP 备份 → 摘除子表 → 重建父表 → 建回转换后子表 → 裸四步表 → 原样建回非转换子表 → **全部索引最后建**）→ m4 校验（`<t>_old` 不存在、目标列 `TEXT`、子表 `REFERENCES` 指向同名新父表、`foreign_key_check` 无行、`integrity_check = ok`）。
- **PG 侧**：逐列显式 `ALTER TABLE … ALTER COLUMN … TYPE timestamptz(6) USING (…)`；秒族用 E1 `date_trunc('microseconds', to_timestamp(col::double precision))`，毫秒族用 **E2 整数拆分式**；D0 列先 `DROP DEFAULT` + `DROP NOT NULL` 再转换；voucher 列用 `IS NULL OR = 0` 单分支。**禁止** `pgTimeColRe` 派生（本批未使用任何正则派生）。
- **checksum（D-017）**：单 checksum = `MigrationChecksum(m0..m4 SQLite 切片, transform_id)`；PG 变体不进哈希。真实值已写入 `internal/store/migrate_test.go` 的冻结目录表（`TestCompiledMigrationCatalogOwnership` 逐条断言）与 `attachments/r2-v73-v87-generated-statements-v0.1.md`。
- **目录/指纹断言更新**：`migrate_test.go`（`len(applied) != 87` + v83–v87 尾部断言 + 冻结表 15 行）、`operations_test.go`、`restart_test.go`、`identity_test.go`（`completeFingerprintCatalogHead = 87`、`lockedHeadExtraTables[87] = {}`）、`modules/jobs/migration/migration_test.go`（按 version 查表，3 个 descriptor）。
- **ledger 写入路径（v73 附带的强制联动）**：`identity.go` 的 `sqliteLedgerDDL` → `TEXT`、`postgresLedgerDDL` → `timestamptz(6)`；`applyMigration` / `applyMigrationPG` / `stampCatalog` 改为经 `sqliteLedgerWrite` / `postgresLedgerWrite` **同事务探测列形状**后按形状绑定（见 `02-execution.md` E-001 记录的实测偏差）。
- **`rebuildOperationLog` fail-closed 断言（`D-019` §5 改动 1）**：新增 `sideTablesAbsent`（SQLite 查 `sqlite_master`、PG 查 `pg_catalog.pg_class`，`'?'` 由各适配器重绑），在 `ALTER TABLE operation_log RENAME` **之前**拒绝仍在的 `operation_log_correlation` / `operation_log_session`；改动 2 遵守（舞步未内嵌进 `rebuildOperationLog`）；v75 descriptor 使用与 F-5 同一通用执行器按**字面 14 步**执行同一编排（改动 3 的「通用 F-5 helper」路径）。所有既有调用点按方言传入对应探测查询。

## 证据

- `go build ./...` exit 0；`gofmt` 对 15 个生成文件 clean。
- `go test ./internal/store/ -run TestCompiledMigrationCatalogOwnership|TestCompleteFingerprintTracksCatalogHead` → ok。
- `go test ./internal/store/ -run TestFullCatalogPostgresBootstrapIntegration` → **ok**（真实 PostgreSQL 15.4 上 v1–v87 全量应用 + 类型/精度断言）。
- `go test ./internal/w040contracttest/` → ok（真实迁移边界矩阵，见 `GOAL-004` E-002）。
- 目录内 15 个 descriptor 的 SQLite `Apply` 在 `TestMigrateFreshDB`（v1–v87 全量）与各模块 provider 测试中实际执行。

## 诚实边界

- 生成器随代码提交（`apps/api/internal/store/vp040_generate_test.go`，仅在 `VP040_GENERATE=1` 时运行，默认 skip）。**可复现性实测**：对同一份 v1–v72 历史重跑生成器，15 个 `vp040_temporal.go` 与语句清单附件 **字节级不变**（16/16 SHA-256 相同）——即「无转录漂移」可被复核。逐列类型编辑的机械性、m0/m4 进入 checksum 输入是否合规，**待 independent 复审**（见 `GOAL-004/03-audit.md` 待复审事项 3）。
- v86 descriptor 的 `ModuleID` 使用模块实际的 `channel.telegram`（台账 §1 写作 `admin.channel.telegram`）——属更正，待复审（待复审事项 2）。
- 本文件**不**声称 `D-021` residual 已闭合；residual 的复审触发 = 本批首个 v73+ 哈希已记录（= 已触发），复审由 checkpoint D 发起。
