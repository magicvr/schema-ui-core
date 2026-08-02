---
title: 执行记录 · R4 · Schema 驱动 CRUD 与 SQLite 持久化闭环
status: active
created: 2026-08-02
updated: 2026-08-02
parent: GOAL-001-production-admin-foundation
version: 0.8.0
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

## 2026-08-02 · 响应 A-003 并冻结 I-007-003

- 用户通过 `/govern` 明确要求：响应 A-003（pass）、同步 Root 台账 R-001、收集并冻结 I-007-003。
- **响应 A-003（independent · pass）**：在 03-audit 写入响应节——采纳 pass（scope 内无开放 required）；**A-002 R-001 → 已落实**（`TestUpdateRecordMonotonicClamp`/`TestRecordMillisecondRoundTrip`/`TestRecordsUpdateStrictlyIncreasesAcrossRapidPatches` 已在 S3 落地）；**A-003 R-001 → fixed**（Root `GOAL-001` 00-meta 纲领 R4 同步为 `active / 3/6`、S3 已实施；仍不勾选 Root R4）；**A-003 R-002 → fixed**（`store.UpdateRecord` 对 PATCH 提供字段做 `strings.TrimSpace` 入库，与 create 一致；新增 store/handler 回归 `TestUpdateRecordTrimsPatchValues` / `TestRecordsUpdateTrimsValues`，`go test ./...` 全绿）。P-004 §3.1 处置：用户明确指示「响应 A-003（pass）」并推进 I-007-003，本轮不补 S3 self 审计，留待 S4/S5 放行或关门前选择（self 或 `/audit`）。
- **收集并冻结 I-007-003（S4/S5 交互契约）**：只读扫描 `apps/web/src/renderer/*`（node whitelist、SchemaTable/FormControls、reactions、permissions/executeAction、row-action、records client）、`components/data-table.tsx`、`app/{App,navigation}.ts`、`protocol/app-manifest.ts` 与 `fixtures/schema/*.json`，以及 S3 后的 records API 与 `me` 会话形状。落盘 [I-007-003-schema-crud-interaction.md](attachments/I-007-003-schema-crud-interaction.md)：代表页 `list-edit-lifecycle`（table actions/toolbar + form `submitAction` + 搜索绑定）+ 字段映射（name/status/owner 可编辑，id/updatedAt 只读）+ 交互状态（成功/加载/空态/错误/确认）+ 权限矩阵（admin/viewer/anonymous，`$context.user.permissions contains "records.write"` 表达式门控，后端 403 权威）+ T-UI-01～10。记录决策 **D-005**。
- 信息台账：`I-007-003` → `verified`；**首个 Schema 写交互代码变更已放行**（S4/S5 实施输入就绪）。
- **未做（本轮）**：未写 Schema 写交互代码（S4/S5 尚未实施）；`I-007-004` 仍 open（S6）；Root R4 未勾选。

## 2026-08-02 · 响应 A-004 与 A-005，修订 I-007-003 至 v0.2.0

- 用户通过 `/govern` 明确要求：响应 A-004 和 A-005。
- **响应 A-004（independent · pass · S3 复核 + A-003 R-001/R-002 关闭证据）**：采纳 pass；独立复核证实——Root 台账同步为 `3/6`/S3 已实施且 Root R4 仍 `[ ]`；`UpdateRecord` 三字段 trim + store/handler 回归均在。scope 内无开放 required/recommended；R-001/R-002 继续作为已闭合 recommended 保留。
- **响应 A-005（independent · conditional · I-007-003 契约合理性）**：三处 required 均按 **fixed** 闭合——修订 [I-007-003-schema-crud-interaction.md](attachments/I-007-003-schema-crud-interaction.md) 至 **v0.2.0** 并在 01-decision 加 **D-005 补记**：
  - **F-002**：冻结唯一页面结构「table + toolbar `create` → modal `create-form`（`submitAction: createRecord` → POST）+ row `edit` → modal `edit-form`（`submitAction: updateRecord` → PATCH，选中行预填）+ row `delete` → `deleteRecord`（DELETE + 确认）」；禁止单 form `submitAction` 双 HTTP 语义。
  - **F-003**：冻结 `records.write` → table 祖先 `permissionCascade.keys: [edit,delete]` + `permissions.edit/delete: "$context.user.permissions contains \"records.write\""`；modal form 各自声明 `permissions.edit`（modal content 为新 permission 根）；禁止仅 `permissionIntent` 无表达式。
  - **F-004**：新增 §9 最小冻结实现规格——顶层 `actions`（createRecord/updateRecord/deleteRecord + `$row.id` 经 `requestMapping.path`）、`meta.requiredCapabilities` 最小集、`table.props.actions/toolbar` 形状、search 归属 `search-form-table`。
  - **R-001**（Renderer 文件白名单 → §9.5 冻结）与 **R-002**（recordView/预填用选中行拷贝 → §2.1/§9.6 冻结）一并 handled。
- `I-007-003` 保持 `verified`（v0.2.0）；**S4/S5 实施放行维持，可无歧义开工**。D-005 已加补记；03-audit 索引/边界更新；后续意见从 A-006 起。
- **未做（本轮）**：未写 Schema 写交互代码（S4/S5 尚未实施）；`I-007-004` 仍 open（S6）；Root R4 未勾选。

## 2026-08-02 · 响应 A-006，修订 I-007-003 §9 至 v0.2.1

- 用户通过 `/govern` 明确要求：响应 A-006——修订 I-007-003 §9（v0.2.1）闭合 F-005/F-006。
- **响应 A-006（independent · conditional · v0.2.0 修订复核）**：采纳 conditional。F-002/F-003 维持 fixed（A-006 复核关闭成立）；F-004 残余按 A-006 收窄由 F-005/F-006 承接。修订 [I-007-003-schema-crud-interaction.md](attachments/I-007-003-schema-crud-interaction.md) 至 **v0.2.1** 并在 01-decision 加 **D-005 v0.2.1 补记**：
  - **F-005**：`updateRecord`（form submit）PATCH `{id}` 槽不适用 row 专属 `requestMapping`/`$row`（核对 `request-construction.ts`：`buildFormAction` 仅组 body、不做 path 槽绑定；`buildRowAction` 才处理 `requestMapping`）。新增 **§9.1a** 冻结「form submit 行上下文槽绑定」——default form submit 且 `action.url` 含 `{id}` 槽时，从打开 modal 时捕获的选中行解析；为 `formAction` 的有界扩展，落入 §9.5 白名单并补测试（T-UI-05）。
  - **F-006**：§9.1–§9.2 改写对齐 `action.schema`/registry（核对 `OutcomeBehavior.behavior`、RequestAction 无 `confirm`、registry `actionRef`）——顶层 **5 个 action**（`createRecord`/`updateRecord`/`deleteRecord` RequestAction + `openCreate`/`openEdit` ModalAction）；`onSuccess` 用 **`behavior`**（enum，无 `type`）；挂载字段用 **`actionRef`**（无 `action` 键、无 `modal:` 前缀）；`confirm` 文案移到 **rowAction** 项；delete 的 `requestMapping.path.id: "$row.id"` 留在 rowAction。
  - **R-001**：§9.5 白名单扩展允许一次性新增 `apps/web/src/renderer/modal*.tsx` / `confirm*.tsx`（T-UI-10 说明）。
  - **R-002**：D-005 v0.2.1 补记显式取代主列表点 1–2 旧表述，声明以 I-007-003 v0.2.1 §2.1/§9 为权威。
- `I-007-003` 保持 `verified`（v0.2.1）；**S4 fixture/actions 接线可无歧义开工**。03-audit 索引/边界/版本更新；后续意见从 A-007 起。
- **未做（本轮）**：未写 Schema 写交互代码（S4/S5 尚未实施）；`I-007-004` 仍 open（S6）；Root R4 未勾选。

## 2026-08-02 · 响应 A-007，修订 I-007-003 §9.2 至 v0.2.2

- 用户通过 `/govern` 明确要求：响应 A-007——I-007-003 §9.2 `confirm` 改为 string（v0.2.2）闭合 F-007。
- **响应 A-007（independent · conditional · v0.2.1 修订复核）**：采纳 conditional。F-005 与 F-006 主体维持 fixed（A-007 复核关闭成立）；修订 [I-007-003-schema-crud-interaction.md](attachments/I-007-003-schema-crud-interaction.md) 至 **v0.2.2** 并在 01-decision 加 **D-005 v0.2.2 补记**：
  - **F-007**：§9.2 delete 的 `confirm` 由 `{ text: "Delete this record?" }` 改为 **`confirm: "Delete this record?"`**（string，与 registry `table.props.actions[].confirm: "type":"string"` 一致；挂载点仍在 rowAction 不变）。一行修补。
  - **R-001**：§9.5 白名单补入 `apps/web/src/protocol/conformance/request-construction.ts`（仅当槽绑定/rowAction 构造需在构造层实现时，改则补测试）与 `renderer/row-action.ts`。
  - **R-002**：§9.1 注明 create/edit/delete 的 `behavior: "reload"` 隐含关闭 modal 并清空选中态。
- `I-007-003` 保持 `verified`（v0.2.2）；**§9 字面形状已对齐 `action.schema`/registry，S4 代表页 fixture 可按 v0.2.2 编写**。03-audit 索引/边界/版本更新；后续意见从 A-008 起。
- **未做（本轮）**：未写 Schema 写交互代码（S4/S5 尚未实施）；`I-007-004` 仍 open（S6）；Root R4 未勾选。

## 下一步计划（非事实）

1. 实施 S4/S5：按 D-005 / I-007-003 **v0.2.2** 演进 `list-edit-lifecycle` fixture（table + `actionRef`→`openCreate`/`openEdit` modal + 行 delete `requestMapping`/`confirm` string + `records.write` 权限表达式 + `onSuccess.behavior: "reload"`（隐含关 modal）），渲染层一次性补齐 actions/toolbar/modal/form-submit（含 §9.1a `{id}` 槽绑定）/反馈；records client 新增 `createRecord`（POST）；`search-form-table` 做 form-to-query 绑定；T-UI-01～10 与权限负向闭环。
2. 在 S6 验收前关闭 `I-007-004`（重启保持与端到端验收协议），再补 create/update/delete→重启→list/detail 的机器可重复证据与 API/Web 回归。
3. 可选：S3/S4 阶段审计（self 或 `/audit`），为 S6 关门审计积累证据。
