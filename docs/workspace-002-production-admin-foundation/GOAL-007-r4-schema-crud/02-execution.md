---
title: 执行记录 · R4 · Schema 驱动 CRUD 与 SQLite 持久化闭环
status: active
created: 2026-08-02
updated: 2026-08-02
parent: GOAL-001-production-admin-foundation
version: 0.4.0
---

# 执行记录 · GOAL-007

## 2026-08-02 · 目标立项

- 用户通过 `/govern` 明确要求按 Root D-010 创建 R4 子目标并登记实施前 required 信息项。
- 在工作区 canonical 根平铺建立本目标五件套与 `attachments/`，设定 `parent: GOAL-001-production-admin-foundation`、`status: active` 与六个顺序成功检查点；同步更新工作区 `goal-tree.md`。
- 记录 D-001，采用一个端到端目标承载 records SQLite 持久化与 Schema CRUD 闭环。
- 登记 `I-007-001`～`I-007-004` 四项 required 信息，分别约束精确 API/错误契约、SQLite 迁移/seed/并发契约、Schema 写交互绑定，以及重启/端到端验收协议；立项时均为 `open`。
- **未做（立项当拍）**：未修改产品代码、数据库、API、Schema fixtures 或 Web 行为。

## 2026-08-02 · 收集并冻结 I-007-001 / I-007-002（S1/S2）

- 用户通过 `/govern` 明确要求：`workspace-002` · `GOAL-007` 先收集 `I-007-001` 与 `I-007-002`。
- 只读对照 `apps/api/internal/handler/records.go`、`records_test.go`、`health.go`（`writeError`）、`apps/api/internal/store/migrate.go` / `store.go` / `seed.go` / `restart_test.go`，以及 Root `I-004` 附件 M-R4-01～07/08/09。
- 落盘：
  - [I-007-001-api-error-contract.md](attachments/I-007-001-api-error-contract.md)：字段/ID/时间戳、继承端点、POST create、稳定 error code 全表、T-API-01～13。
  - [I-007-002-sqlite-migration-plan.md](attachments/I-007-002-sqlite-migration-plan.md)：`records` DDL、`0003 records_persist`、空表 seed、repository/并发、T-DB-01～09 与静态切片退出路径。
- 记录决策 **D-002**（API/错误契约）与 **D-003**（SQLite 迁移/seed/repository）。
- 信息台账：`I-007-001` → `verified`；`I-007-002` → `verified`。
- 成功标准：**S1**、**S2** 勾选；派生进度 `0/6` → **`2/6`**。
- **未做**：未修改产品代码、未执行迁移、未新增 POST 实现、未跑 R4 产品测试；`I-007-003`/`I-007-004` 仍为 open；Root R4 未勾选。

## 2026-08-02 · 响应 A-001 F-001 并放行 S3

- 用户裁决（P-004）：F-001 修正方向选 **毫秒精度 + 保留严格递增**；按 A-001 为 independent 且无 self 审计，先补 **self 自审**再统一响应。
- 落盘 **D-004**：`updated_at` 存储精度由 Unix 秒统一为 Unix **毫秒**；API `updatedAt` 序列化为 RFC3339 **含毫秒**（`2006-01-02T15:04:05.000Z07:00`）；保留「严格晚于」，同一毫秒内以单调钳制（`prev + 1ms`）保证确定性，禁止人为跳秒。D-002/D-003 加修订注记；I-007-001/002 更新至 **v0.2.0**（精度、映射、seed、断言同步）。
- 写 **A-002**（self · `pass`）：复核 S1/S2 冻结与 F-001 修正证据，无新 required；R-001（recommended）要求 S3 覆盖「同一毫秒钳制」与毫秒往返测试。
- **F-001 → `fixed`**（03-audit 响应节）：证据 = D-004 + 附件 v0.2.0 + A-002 pass。
- **S3 实施放行**：`I-007-001`/`I-007-002` verified + F-001 fixed；`I-007-003`（S4/S5）、`I-007-004`（S6）仍为 open required。
- 派生进度保持 **`2/6`**（S1/S2 已勾选；S3 尚未实施，未改 status/progress）。
- **未做**：未修改产品代码、未执行迁移、未新增 POST/SQLite repository 实现。

## 2026-08-02 · 实施 S3 持久化 CRUD API（SQLite 路径 + T-API/T-DB）

- 用户通过 `/govern` 明确要求实施 S3：落地 0003 迁移（updated_at 毫秒）、repository、seedRecords、POST/list/detail/PATCH/DELETE 的 SQLite 路径与 T-API/T-DB 单测（含毫秒钳制/往返用例）。
- **迁移（0003 records_persist）**：`apps/api/internal/store/migrate.go` 追加 `version:3` / `records_persist` / `0003:records-persist:v1`；DDL 只建 `records` 表（`updated_at INTEGER` Unix 毫秒）+ `name`/`updated_at`/`owner` 索引，不在迁移内插入种子行。快照逻辑通用化为 `snapshotBeforePending(firstPendingVersion)`：已有 v2 库升级到 0003 前产生 `<db>.pre-v0003-<UTC>.sqlite`；`dbHasRows` 判定是否值得快照（排除 `schema_migrations` ledger 行，避免空库误快照）。
- **repository**：`apps/api/internal/store/records.go` 新增 `Record`/`RecordFilter`/`RecordPatch` 与 `ListRecords`（SQL 侧 q/sort/order/page，name 用 NOCASE、`id ASC` 决胜序保证分页确定性）、`GetRecord`、`CreateRecord`（事务内 `SELECT EXISTS` 判 PK 碰撞 → 哨兵 `ErrRecordExists` 供 handler 重试）、`UpdateRecord`（读改写 + 单调钳制 `prev+1ms`，last-write-wins）、`DeleteRecord`。
- **seedRecords**：`apps/api/internal/store/seed_records.go`；`Open(seedAdmin=true)` 在 `seedRBAC` 后执行，**仅当 `records` 表行数为 0** 时插入与旧 `staticRecords()` 对齐的 8 行（`rec-1…rec-8`，`2026-07-31T00:00:00Z` 起 +11h，Unix 毫秒）；非空整段跳过，删除/更新/新建在重启后保持。
- **handler**：`apps/api/internal/handler/records.go` 重写为 `recordHandler{st *store.Store}`；`Register` / `recordsHandler` 增加 `*store.Store` 注入（`health.go`、`cmd/server/main.go`、`testhelpers_test.go` 同步）；list/detail/PATCH/DELETE 全走 repository；新增 `POST /api/records`（`records.write`，成功 **201** + 完整 record；`crypto/rand` 生成 `rec-` + 16 位小写 hex，PK 碰撞有限重试后 `INTERNAL`；`INVALID_CREATE_BODY` 与 `INVALID_CREATE_FIELD` 400 分离）。**`staticRecords()` 与进程内切片作为生产路径已删除**（T-DB-09）。
- **毫秒契约**：handler `updatedAt` 自定义类型固定 RFC3339 3 位毫秒序列化（`2006-01-02T15:04:05.000Z07:00`）；DB 列 Unix 毫秒；同毫秒连续 update 由单调钳制保证「严格晚于」。
- **测试**：
  - store 新增 T-DB-01～09 与 R-001 毫秒钳制/往返用例：`TestRecordsTableEmptyBeforeSeed`、`TestSeedRecordsEmptyTable`、`TestSeedRecordsSkipsNonEmpty`、`TestRecordsPersistAcrossRestart`、`TestUpdateRecordMonotonicClamp`、`TestRecordMillisecondRoundTrip`、`TestMigrateExistingV2ToV3`、`TestRecordsSeedIdempotentAcrossOpens`、`TestMigrateFailClosedRecordsChecksumDrift`、`TestRecordsRepositoryNotFoundAndCollision`、`TestListRecordsFilterSortPagination`；更新 `TestMigrateFreshDB`/`TestRestartPersistence` 为 ledger {1,2,3}。
  - handler 新增 T-API-08～13 与毫秒格式/往返：`TestRecordsCreate`、`TestRecordsCreateInvalidField`、`TestRecordsCreateInvalidBody`、`TestRecordsAdminLifecycle`、`TestRecordsUpdateStrictlyIncreasesAcrossRapidPatches`、`TestRecordsHandlerReadsFromStore`；扩展匿名 401 与 editor/viewer 403 覆盖 POST。
  - 结果：`go vet ./...` 干净；`go test ./...`（apps/api）全绿；web `vitest run` 443/443 通过。E2E（playwright）需真实服务，未在本轮内执行。
- **API README 同步（R4）**：端点表新增 POST、记录数据源与毫秒语义、鉴权边界与测试覆盖更新至当前 R4 边界。
- **成功标准 S3 勾选**：POST/list/search/detail/PATCH/DELETE 的 SQLite repository 路径已实现；生产默认不再依赖进程内 records；认证、`records.read` / `records.write` 与统一错误 envelope 保持一致。派生进度 `2/6` → **`3/6`**。
- **未做（本轮）**：未修改 Schema 写交互代码（S4/S5 待 `I-007-003` 关闭）；未跑进程重启 E2E 与 S6 验收（待 `I-007-004`）；未写正式审计（S3 阶段审计可下一拍选择 self 或 `/audit`）。

## 下一步计划（非事实）

1. 在首个 Schema 写交互变更前关闭 `I-007-003`（Schema CRUD 页面/Node/action 绑定、字段映射、成功/加载/空态/错误/确认交互与 admin/viewer/匿名权限矩阵），再实施 S4/S5 页面读写与权限负向闭环。
2. 在 S6 验收前关闭 `I-007-004`（重启保持与端到端验收协议），再补 create/update/delete→重启→list/detail 的机器可重复证据与 API/Web 回归。
3. 可选：对 S3 做一次阶段审计（self 或 `/audit` 独立），为 S6 关门审计积累证据。
