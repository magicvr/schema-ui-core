---
title: I-011-002 · records 版本化退场与迁移策略
status: active
doc_type: contract
created: 2026-08-03
updated: 2026-08-03
parent: GOAL-011-s4-semantic-admin-resources
version: 0.2.0
related_info: I-011-002
related_decision: D-002, D-003
supersedes: null
---

# I-011-002 · records 退场与迁移策略（冻结）

> **性质**：回答「records 如何从当前产品基线退场，同时保持既有迁移/checksum、已有数据库、数据处置、API/权限/菜单/operation log 与历史文档可追溯」。建立 fresh install 与 in-place upgrade 两条迁移矩阵，逐项盘点运行时代码、迁移、种子、文档与测试。冻结后 `I-011-002` 由 GOAL-011 **D-002** 置为 `verified`，解除 S1 方案冻结与 S3 退场实施门禁。
> **v0.2.0（2026-08-03 · A-002 响应，GOAL-011 D-003）**：修订 §2.1/§2.3 快照语义为「**每个**待应用数据变更迁移前快照」（或至少在 `0006` 前强制），保证 0005+0006 同批待应用时 `pre-v0006` 快照必然存在；§5 验收句与实现行为对齐。`I-011-002` 维持 `verified`（响应修订不改变冻结结论）。
> **依据**：`apps/api/internal/{handler,store}` 与 `apps/web/src` 全量静态核对（records 足迹）；[I-010-001 v0.2.1](../../GOAL-010-a002-schema-adapter/attachments/I-010-001-schema-resource-contract.md) §6（S1～S3 零对外变更历史；S4 终态由本目标承接）；GOAL-011 D-001（records 退场、users/roles 承接）；D-002（用户裁决：硬退场 DROP TABLE）。
> **硬约束**：不得改写 `0001`～`0004`（迁移账本 checksum 校验，fail closed）；不得改写 GOAL-004/007 历史决策/执行/审计文档（历史事实保留）。退场只能通过**新迁移** + **代码/种子/fixture 演进**。

## 1. records 当前足迹盘点（已核对）

| 面 | 位置 | 当前状态 | 退场动作（S3） |
|----|------|----------|----------------|
| API 注册 | `handler/health.go` `registerResource(recordsResource)` | `/api/records` 五路由 | 移除；改注册 users/roles |
| handler | `handler/records.go`（`recordsResource`+`recordsEntity`） | 工厂注册实例 | 删除 |
| store | `store/records.go`（List/Get/Create/Update/DeleteRecord） | SQLite 全 CRUD | 删除 |
| 种子数据 | `store/seed_records.go`（rec-1…rec-8） | 空表时 8 行 | 删除 |
| 种子 RBAC | `store/seed.go` `seedRBAC`：`perm-records-read/write`、`menu-list-edit-lifecycle` + grants | admin rw+menu、editor/viewer r | 移除 records 权限/菜单；新增 users/roles 权限/菜单 + grants（I-011-001 §4） |
| 迁移 | `migrate.go` `0003 records_persist`（records 表 + 3 索引） | 已应用（不可改写） | 由 `0006` 退场 |
| 操作日志 | `store/operations.go` + `0004` CHECK：`records.create/update/delete` | records 写事件 | CHECK 扩为含 users/roles 事件（`0005`）；records.* 事件**保留**为历史合法值 |
| fixture | `fixtures/schema/list-edit-lifecycle.json` | dataSource `/api/records` + 全 CRUD action | 删除（由 users/roles 页取代） |
| fixture | `fixtures/schema/{data-table,search-form-table,catalog}.json` | dataSource `/api/records` | 改指 users/roles 或删除（见 §3） |
| manifest | `public/.well-known/schema-ui/app-manifest.json` | `list-edit-lifecycle` 页 + `menu_list_edit_lifecycle` visibleWhen；`catalog` 页 | 移除 records/catalog 页与菜单；新增 users/roles 页（`menu_users`/`menu_roles` visibleWhen） |
| 前端专名 | `renderer/records.ts`（Record* 类型/函数名）、`use-records.ts`、`App.tsx` `recordsFetcher` prop | records 专名 | 泛化/更名（`fetchRecords` 保留为通用 transport；`use-records.ts` 删除；prop 更名 `resourceFetcher`） |
| 测试 | `handler/records_test.go`、`records_restart_test.go`、`operations_test.go`（records 事件）、`store/records_test.go`、`seed_test.go`（records 断言）、`resources_test.go`（records 权限引用） | records 回归 | 替换为 users/roles 回归 + 重启持久化 + 操作日志断言 |
| 测试 | `web` `records.test.ts`、`schema-crud.test.tsx`、`representative-pages.integration.test.tsx`（RECORDS emulator）、`app-examples.test.tsx`、e2e `schema-crud.spec.ts`/`shell.spec.ts` | records 前端/端到端 | 替换为 users/roles 端到端；emulator 形状复用 |
| 数据文件 | `apps/api/data/schema-ui.db`（dev，含 records）、`r2-browser-test.db` | 既有库 | 升级路径走迁移矩阵（§2）；`r2-browser-test.db` 为历史遗留，不迁移 |
| 文档 | `QUICKSTART.md`（「8 条演示 records」「Acme Console」） | 用户文档 | 更新为 users/roles |
| 历史治理 | GOAL-004/007 决策/执行/审计、I-007-001/002/003 契约、GOAL-010 S1～S3 执行记录 | 历史事实 | **保留**（不可改写）；以「历史证据」引用 |

## 2. 迁移矩阵

### 2.1 新迁移（退场通过新版本，不改既有）

| 版本 | 名称 | 内容 | 所属阶段 |
|------|------|------|----------|
| `0005` | `operation_log_expand` | 重建 `operation_log` 表：event CHECK 扩为含 `users.create/update/delete`、`roles.create/update/delete`（**保留** `records.*` 与 `auth.*`）；`SQLite` 改 CHECK 需 `CREATE 新表 → COPY → DROP → RENAME` 重建，事务内完成；checksum 以新 DDL + transformID `0005:operation-log-expand:v1` 冻结 | S2（users/roles 写事件需要） |
| `0006` | `records_retire` | **硬退场**：`DROP TABLE records`；清理 records 权限/菜单行（`permissions.perm-records-read/write`、`menu_items.menu-list-edit-lifecycle`、对应 `role_permissions`/`role_menu_items`）；`SQLite` 删行按 id 定位，不存在则跳过（幂等）；transformID `0006:records-retire:v1` | S3 |

### 2.2 fresh install（新库，fork 从零开始）

```
空库 → 0001 r2_baseline → 0002 rbac_expand → 0003 records_persist
     → 0004 operation_log → 0005 operation_log_expand → 0006 records_retire
```

终态：**无 `records` 表**；`operation_log` CHECK 含 users/roles/records/auth 事件；users/roles/permissions/menu 就绪。种子（seedAdmin→seedRBAC→seedUsersRoles）不创建任何 records 权限/菜单/数据。

### 2.3 in-place upgrade（既有库）

- 路径 A（已应用至 `0003`，无 0004）：`pre-v0004` 快照 → 0004 → 0005 → 0006。
- 路径 B（已应用至 `0004`）：`pre-v0005` 快照 → 0005 → `pre-v0006` 快照 → 0006。
- **快照机制（v0.2.0 · A-002 F-002）**：`snapshotBeforePending` 泛化为对**每个**待应用的数据变更迁移（版本 ≥2）在应用**前**各快照一次（`VACUUM INTO`，命名 `*.pre-v%04d-<UTC>.sqlite`）；至少须在 `0006` 前强制快照。因此 0005 与 0006 **同批待应用**时，`pre-v0005` 与 `pre-v0006` 均存在——`pre-v0006` 是 records 数据处置的可恢复兜底，`pre-v0005`（0004 态，仍含 records 表）为次级兜底。
- 0006 应用后，records 表及其数据**不可再访问**（表已删除）；如需恢复，回滚到 `pre-v0006` 快照。
- 校验：`validateApplied` 检查连续性、名称、checksum；任何漂移 fail closed。0001～0004 的 checksum 不受影响。

### 2.4 数据处置（用户裁决：硬退场）

- 既有库中的 records 行（种子 rec-1…rec-8 与用户自建行）随 0006 `DROP TABLE` 删除。
- **可恢复性**：升级前自动 `pre-v0006` 快照（既有机制）；用户如需 records 数据，从快照提取。
- 这是对「已有数据库数据」的显式处置决策（D-002，用户裁决采纳），不是静默丢数据。

## 3. 代码 / 种子 / fixture / 前端退场动作（S3）

1. **API**：`health.go` 移除 `recordsResource` 注册，改注册 `usersResource`、`rolesResource`；删除 `handler/records.go`。
2. **store**：删除 `store/records.go`、`store/seed_records.go`；`seed.go` 移除 records 权限/菜单，按 I-011-001 §4 新增 users/roles 权限/菜单 + grants；`migrate.go` 追加 `0005`/`0006`。
3. **fixture**：
   - 新增 `fixtures/schema/users.json`、`roles.json`（两语义资源列表/CRUD 页；S4 接入）。
   - 删除 `list-edit-lifecycle.json`（records 专属 CRUD 演示，由 users/roles 页取代）。
   - `data-table.json`、`search-form-table.json`：**改指** users/roles 的 `dataSource`（保留 conformance 示例页价值，去 records 耦合）。
   - `catalog.json`：**删除**（D-004 已降 catalog 为 genericity 历史示例；catalog 页无语义且 dataSource 指 records）。manifest 同步移除 `catalog` 页。
4. **manifest**：移除 `list-edit-lifecycle`、`catalog` 页与 `menu_list_edit_lifecycle`；新增 users/roles 页（sidebar，`visibleWhen` = `$context.features.menu_users == true` / `menu_roles == true`）。
5. **前端专名**：`records.ts` 的 `Record*` 别名与 `RecordsQuery` 等更名/泛化（`fetchRecords` 通用 transport 保留）；删除 `use-records.ts`；`App.tsx` 的 `recordsFetcher` prop 更名 `resourceFetcher`（渲染行为不变——S3 退场步骤，S4 不再动 Renderer）。
6. **测试**：records 后端测试替换为 users/roles 后端测试（含 401/403、restart 持久化、操作日志事件）；web `records.test.ts`/schema-crud/representative-pages/e2e 替换为 users/roles 形态；`resources_test.go` 的 genericity 测试改用 users/roles 权限键或内存测试资源（records 权限键退场）。
7. **文档**：`QUICKSTART.md` records 引用更新为 users/roles；历史治理文档保持，标注为「历史证据（已退场实体）」。
8. **保留**：I-007-001/002/003 契约作为 records 历史契约（不再演进）；GOAL-004/007/010 历史记录；`I-PROTO-001 v0.1.3` 覆盖不变。

## 4. 双轨边界

- **不做新旧双轨并行**（延续 I-010-002 §6）：records 退场是**一次性迁移 + 全量回归**；迁移期间不新增 records 业务 fixture 或测试。
- 迁移 `0005` 保留 `records.*` 事件为**历史合法值**（存量 operation_log 行 + 可能的历史写入路径不失效），但**不再有** records 写代码产生新事件。
- `operation_log` 表的 `record_id` 列保留（历史行语义不变）；users/roles 事件复用 `record_id` 承载 user/role id（列名不变量，不改列）。

## 5. 验收口径（S3 完成定义）

- fresh install 库中不存在 `records` 表；`sqlite_master` 查询验证。
- in-place upgrade 从 0004 基线到 0006 后，无 `records` 表、无 records 权限/菜单行；**每个待应用数据变更迁移前的快照均存在**（0005+0006 同批时 `pre-v0005` 与 `pre-v0006` 均存在，`pre-v0006` 为 records 数据可恢复兜底）。
- `go test ./...` 全绿（records 测试已替换）；web `vitest run` 全绿 + `tsc -b`/build 干净。
- `grep` 产品代码（`apps/api`、`apps/web/src` 非历史区）无 `api/records`/`records.read`/`records.write`/`menu_list_edit_lifecycle`/`list-edit-lifecycle` 残留（历史文档区除外）。
- 操作日志能记录 users/roles 写事件（`0005` 生效）。

## 6. 非目标

- 改写历史迁移/文档；records 数据在既有库中保留；双轨兼容层；扩大 `I-PROTO-001` 覆盖；`r2-browser-test.db` 等历史测试库的迁移。

## 7. 证据索引

- `apps/api/internal/handler/health.go`（注册面）、`handler/records.go`（handler）
- `apps/api/internal/store/{records.go,seed_records.go,seed.go,migrate.go,operations.go}`
- `apps/api/internal/handler/fixtures/schema/*.json`、`apps/web/public/.well-known/schema-ui/app-manifest.json`
- `apps/web/src/renderer/{records.ts,use-records.ts,schema-table.tsx}`、`apps/web/src/app/App.tsx`
- `apps/web/src/app/representative-pages.integration.test.tsx`、`apps/web/e2e/*.spec.ts`
- `apps/api/data/*.db`（既有库）、`QUICKSTART.md`
- [I-011-001-users-roles-contract.md](I-011-001-users-roles-contract.md)（users/roles 承接契约）

## 8. 修订记录

| 版本 | 日期 | 变更 |
|------|------|------|
| 0.1.0 | 2026-08-03 | 冻结（GOAL-011 D-002，用户裁决：硬退场 DROP TABLE + 操作日志纳入）；关闭 `I-011-002` |
| 0.2.0 | 2026-08-03 | A-002 响应（GOAL-011 D-003）：§2.1/§2.3 快照语义改为**每待应用数据变更迁移前快照**（0005+0006 同批时 `pre-v0006` 必然存在，F-002）；§5 验收句对齐。`I-011-002` 维持 `verified` |
