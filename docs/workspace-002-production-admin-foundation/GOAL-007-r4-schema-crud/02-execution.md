---
title: 执行记录 · R4 · Schema 驱动 CRUD 与 SQLite 持久化闭环
status: active
created: 2026-08-02
updated: 2026-08-02
parent: GOAL-001-production-admin-foundation
version: 0.12.0
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

## 2026-08-02 · 实施 S4/S5 Schema 驱动读写主路径与权限负向闭环

- 用户通过 `/govern` 明确要求实施 S4/S5。信息门禁 `I-007-003` 已 verified（v0.2.2），首个 Schema 写交互代码变更放行维持；本轮无开放 required，无到期 required 信息门禁（`I-007-004` 仅约束 S6）。P-004 无新增裁决点（self 审计按既有用户裁决留待 S4/S5 放行或关门前）。
- **代表页 fixture 演进（v0.2.2 结构）**：`apps/api/internal/handler/fixtures/schema/list-edit-lifecycle.json` 改写为 I-007-003 §2.1/§9 唯一结构——`section + table`（`records-table`：`permissionCascade.keys [edit,delete]` + `permissions.edit/delete = $context.user.permissions contains "records.write"`、`props.toolbar` `create`→`openCreate`、`props.actions` `edit`→`openEdit` / `delete`→`deleteRecord`（`requestMapping.path.id: "$row.id"` + `confirm: "Delete this record?"`））+ `recordView`（`record-detail`，选中行拷贝）+ 顶层 **5 actions**（`createRecord`/`updateRecord`/`deleteRecord` RequestAction + `openCreate`/`openEdit` ModalAction，`onSuccess.behavior: "reload"`，两 modal 各含 default form：`create-form`/`edit-form` 各自 `submitAction` + `permissionCascade.keys [edit]` + `permissions.edit`）；`meta.requiredCapabilities` 更新为 §9.3 最小集（`permissions.inheritance`、`actions.row.request`、`actions.page.trigger`、`table.sort` 等），去掉 `record.view.load`/`form.controls.extended`。页面结构从「tabs+静态 recordView+无绑定 form」改为「table + toolbar/rowActions + 双 modal + 行删除确认」。
- **search 绑定 fixture**：`search-form-table.json` 的 search form 加 `mode: "search"` + `targetTable: "search-results"`（T-UI-03 form-to-query 归属该页，§9.4）。
- **records 客户端**：`renderer/records.ts` 新增 `createRecord`（POST，解析 201）；统一非 OK 响应解析为 `RecordApiError`（携带冻结 envelope `{error,message}` 的 `code`/`status`），`fetchRecords`/`updateRecord`/`deleteRecord` 共享 `readRecordApiError`；消息保留 `HTTP {status}` 子串，既有客户端测试保持绿。
- **renderer 一次性补齐（§9.5 白名单内）**：
  - `render.ts`：`RenderFormNode.props` 增加 `submitAction`/`mode`/`targetTable` 并透传解析；`RenderTableNode.props` 类型化 `actions`/`toolbar`。
  - `render.tsx`：新增 `SchemaCrudProvider` + `useSchemaCrud` 上下文——选中行（喂 recordView + edit-form 预填）、按 table id 的查询态（search form-to-query + post-write reload 共用）、`reloadToken`、`activeModal`、`pendingConfirm`、`feedback`、fetcher 注册、通用动作执行器（经冻结 `executeAction` 门禁 + `constructRequest` 构造，§9.1a `{id}` 槽从打开 modal 时捕获的行上下文解析——**有界扩展落在 `render.tsx` 内、不改 conformance 包**）。`FormView` 增加 default-mode submit（create POST / edit PATCH，§3 提交前 trim）与 search-mode 提交；`RecordView` 无静态 record 时渲染选中行拷贝；surface 渲染反馈区 + modal + confirm。`render.test.tsx` 等既有表单/反应测试保持绿（无 submitAction 的表单不渲染提交按钮）。
  - `components/data-table.tsx`：行点击选中 + `selectedKey` 高亮（`onRowClick`/`aria-selected`）。
  - `renderer/schema-table.tsx`：渲染 `props.toolbar`（create）与 `props.actions`（edit/delete）按钮，`effectivePermission` 驱动禁用（viewer/editor 只读）；行选中传递；查询挂 provider（无 provider 时回落本地）；把注入 fetcher 注册给 provider（保持首个注入值，避免测试内联 fetcher 身份变化导致重渲染循环）。
  - **新增文件**（§9.5 允许的一次性补齐）：`renderer/modal.tsx`（modal 宿主）、`renderer/confirm.tsx`（确认对话框）。
- **测试**：
  - 新增 `renderer/schema-crud.test.tsx`：T-UI-01（列表加载/空态）、T-UI-02（列排序 sort/order 刷新）、T-UI-03（search form-to-query 过滤）、T-UI-04（create POST 201 + 新行出现 + 空字段 INVALID_CREATE_FIELD）、T-UI-05（edit PATCH 预填 + `updatedAt` 刷新）、T-UI-06（delete 确认 → DELETE 204 / 取消无请求）、T-UI-07（统一 envelope 呈现 FORBIDDEN / RECORD_NOT_FOUND）、T-UI-08（admin 全启用 / viewer 只读禁用）、T-UI-09（后端权威不被前端隐藏替代）、T-UI-10（页级变更仅 fixture——record 动作 id 只在 fixture、不在 renderer 源码）。测试用 I-007-001 契约的 in-memory records API 仿真（GET/POST/PATCH/DELETE + envelope）。
  - 更新 `renderer/representative-pages.test.tsx` 与 `app/representative-pages.integration.test.tsx` 的 list-edit-lifecycle 断言为新的 table+toolbar+rowActions+recordView 结构（admin 上下文）。
  - 结果：`tsc -b` 干净；web `vitest run` **458/458** 全绿（23 文件）；`vite build` 成功；`go test ./...`（apps/api）全绿（fixture embed 正常）。
- **成功标准 S4、S5 勾选**：Schema 页面完成列表/搜索/详情/新建/编辑/删除，新增代表页仅改 fixture 不碰 Renderer 主路径（T-UI-10）；字段校验、加载/空态、成功反馈、删除确认、统一错误与权限负向（viewer/editor 禁用、后端 403 权威）闭环（T-UI-01～09）。派生进度 `3/6` → **`5/6`**。
- **未做（本轮）**：未实施 S6（`I-007-004` 仍 open，需 create/update/delete→重启→list/detail 机器可重复证据）；未勾选 Root R4（Root 保持 `3/5`）；未跑 playwright E2E（需真实服务）；未写正式审计（S4/S5 阶段审计留待放行或关门前选择 self 或 `/audit`）。

## 2026-08-02 · 响应 A-008 + self A-009 + 收集 I-007-004 + 实施 S6

- 用户通过 `/govern` 明确要求：`workspace-002` · `GOAL-007` 响应 A-008，并收集 I-007-004 后实施 S6。
- **P-004 §3.1（用户裁决）**：A-008 为 `source: independent` 且 S4/S5 scope 尚无 self 审计；用户裁决「**先补 S4/S5 self 审计**」。已写 **A-009（self · pass）**（复跑 schema-crud + representative-pages 29 项全绿、`go test ./...` 全绿、渲染源码 grep 无 action id 硬编码），随后**统一响应 A-008**（采纳 pass，A-009 作为 S4/S5 scope 的 self 覆盖）。
- **收集并冻结 I-007-004（S6 验收协议）**：落盘 [I-007-004-restart-e2e-protocol.md](attachments/I-007-004-restart-e2e-protocol.md) 并记录决策 **D-007**——「服务重启」分 **L1 HTTP 层**（同文件 store 关闭→重开，全 HTTP CRUD→list/detail）与 **L2 进程级**（真实 `cmd/server` 子进程终止→同 `DB_PATH` 重启）两层；每轮全新临时 DB（`t.TempDir()`）+ 空闲端口隔离；固定操作序列（admin 登录→POST create→PATCH rec-1→DELETE rec-2→重启→login→list/detail）；断言 create 新行存在、rec-1 已更新（含 `updatedAt` 毫秒精确一致）、rec-2 不复活、`total=8`（非空表 seedRecords 不复活）；迁移/seed 重跑与失败路径沿用既有 store/handler/browser E2E 覆盖；`go test ./...` / `npm test` / `npm run test:e2e` 为回归命令。`I-007-004` → **verified**。
- **实施 S6（L1 + L2）**：
  - **L1 · HTTP 层重启**：`apps/api/internal/handler/records_restart_test.go` `TestRecordsSurviveRestart`——临时 SQLite 文件上完整 handler/auth 栈，admin 登录→POST create（201，记录 id）→PATCH rec-1（200，记录 `updatedAt`）→DELETE rec-2（204）→`store.Close()`→以**同一路径**重开（不同 seed hash，不覆盖 admin 密码）→重新登录→list 断言新行存在 / rec-1 名称已更新 / rec-2 不复活 / `total=8`；detail 断言字段与 `updatedAt` 毫秒持久化一致。
  - **L2 · 进程级重启**：`apps/api/cmd/server/server_restart_test.go` `TestServerProcessRestartPersistsRecords`——`go build ./cmd/server` 生成真实二进制，以显式 env（`DB_PATH` 临时路径、`HTTP_ADDR` 空闲端口、`ADMIN_INITIAL_PASSWORD=admin`、`AUTH_JWT_SECRET=test-secret`、`AUTH_DEV_SESSION_ENABLED=false`、`APP_ENV=development`）启动 Phase 1 进程→登录→POST/PATCH/DELETE→`Process.Kill()` 终止→以**同一 `DB_PATH`** 启动 Phase 2 进程→登录→list/detail 断言同 L1。测试结束 Kill+Wait 不留进程，临时库随 `t.TempDir()` 清理。
  - **验证结果**：`go vet ./...` 干净；`go test ./... -count=1`（apps/api）**全绿**（含 L1 0.13s、L2 4.38s）；web `vitest run` **458/458** 全绿（本轮无 web 变更）。
- **成功标准 S6 勾选**：create/update/delete→重启→list/detail 的机器可重复证据已产出（L1+L2）；迁移/seed 重跑（ledger `{1,2,3}` 不重跑、空表才 seed、非空不复活）与关键失败路径（checksum 漂移、快照恢复、401/403）由既有 store/handler/browser 测试覆盖；API 回归 `go test ./...`、Web 回归 vitest 458/458。派生进度 `5/6` → **`6/6`**。
- **未做（本轮）**：本目标仍为 `active`，**未置 `done`**（关门需先做关门审计 self 或 `/audit` + 用户裁决 + Root R4 勾选）；Root R4 未勾选（Root 保持 `3/5`）；`npm run test:e2e`（playwright，需真实服务）本轮未跑，属可选回归。

## 2026-08-02 · 响应 A-010 · 修正 L2 updatedAt 跨进程 detail 断言（F-008 → fixed）

- 用户通过 `/govern` 明确要求：响应 A-010（independent · conditional · close-out）——修复 F-008 的 L2 `updatedAt` 跨进程 detail 断言并复跑；同步修正执行事实，随后请求 finding-closure 复核。
- **修正 `apps/api/cmd/server/server_restart_test.go`（L2，I-007-004 §3.6/§4）**：
  - `httpPatch` 改为返回 200 响应体；`httpCreate` 返回完整响应（含 `updatedAt`）；文件顶部注释补充 A-010 F-008 说明。
  - Phase 1：记录 create 响应 `createdAt` 与 PATCH `rec-1` 响应 `patchedAt`（毫秒 RFC3339）。
  - Phase 2：新增 `GET /api/records/{createdID}` 断言字段 + `updatedAt == createdAt`；新增 `GET /api/records/rec-1` 断言 `name == "Acme Rebrand"` 且 `updatedAt == patchedAt`（毫秒精确一致，跨进程持久化往返）。
- **执行事实修正**：此前「L2 … list/detail 断言同 L1」的表述偏满（A-010 指出）——L2 原仅以 list 检查 rec-1 名称、只 GET 新建记录 detail，未核对 rec-1 PATCH 的 `updatedAt` 跨进程往返。本轮补齐该断言后，该表述成立且与 L1 一致。
- **复跑**：`go vet ./cmd/server/` 干净；focused `go test ./cmd/server -run '^TestServerProcessRestartPersistsRecords$' -count=1` **PASS（4.32s）**；`go test ./... -count=1`（apps/api）全绿（cmd/server、handler、store、auth、account）。
- **F-008 → `fixed`**（03-audit A-010 响应节，2026-08-02）：L2 现按 I-007-004 §3.6/§4 对新建记录与 `rec-1` 均做 detail 断言且 `updatedAt` 毫秒精确一致；关闭证据已请求 finding-closure 复核（`/audit`）。
- **仍开放（非本意见 required）**：R-003（API README 端点表阶段标注，recommended）、R-004（真实浏览器 CRUD E2E，recommended，A-010 判定非阻断）。本目标仍 `active / 6/6`，**未置 `done`**；Root R4 未勾选（Root 保持 `3/5`）。后续意见从 A-011 起。

## 2026-08-02 · 响应 A-011/A-012 + close-out self 审计 A-013 + GOAL-007 关门

- 用户通过 `/govern` 明确要求：响应 A-012（及 A-011）；补 close-out self 审计；随后处理 GOAL-007 关门与 Root R4 勾选。
- **P-004 §3.1（用户裁决）**：A-011/A-012 均为 independent；用户裁决「**补 close-out self 审计**」→ 已写 **A-013（self · close-out · pass）**，为 GOAL-007 整体与 S6/L2 scope 补齐 `source: self` 覆盖（既有 A-009 self 仅覆盖 S4/S5）。
- **响应 A-011 / A-012（independent · finding-closure · pass）**：两轮独立复核均确认 A-010 F-008 的关闭证据充分可核对；**F-008 维持 `fixed`**；R-003/R-004 保持 recommended 非阻断。
- **A-013 close-out self 审计**：逐项核对成功标准 S1～S6（全 `6/6`）、四项 required 信息门禁 `I-007-001/002/003/004` verified、F-001～F-008 全部 fixed、无开放 required；复跑 `go test ./...`（apps/api）全绿 + web `npm test`（vitest）**458/458** 全绿；**verdict = pass**，确认具备关门条件。
- **关门（GOAL-007 → `done`）**：置 00-meta `status: done`；Root R4 检查点勾选（Root `3/5 → 4/5`）；同步 goal-tree（GOAL-007 状态列 + 台账）与 Root 00-meta 纲领 R4。R5 与 VP-002 保持 open，不受本轮关门影响。
- **未做（本轮）**：R-003（API README 端点表阶段标注）与 R-004（真实浏览器 CRUD E2E）作为 recommended 非阻断留待后续；Root R5 未立项。

## 下一步计划（非事实）

1. **R4 已关门**：GOAL-007 `done`，Root R4 已勾选（Root `4/5`）。Root 下一主路径为 **R5 · 工程化、fork 体验与集成关门**（需先收集 `I-005`（部署基线/15 分钟 fork 计时口径）并复核 `I-006`（最小操作日志取舍），再立项 R5 子目标）。
2. 可选补充：R-003（`apps/api/README.md` 端点表阶段标注统一为 R4）；R-004（真实浏览器 `list-edit-lifecycle` CRUD E2E）。
3. R5（容器/生产运维）与 fork 关门属后续目标，不在本目标范围。
