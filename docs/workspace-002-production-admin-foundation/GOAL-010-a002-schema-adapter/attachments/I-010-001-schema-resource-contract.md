---
title: I-010-001 · Schema 驱动通用资源契约
status: active
doc_type: contract
created: 2026-08-03
updated: 2026-08-03
parent: GOAL-010-a002-schema-adapter
version: 0.2.0
related_info: I-010-001, I-010-002
related_decision: D-002, D-003
---

# I-010-001 · Schema 驱动通用资源契约（冻结）

> **性质**：回答「通用资源契约的精确形状是什么」——把 A-002 F-002-001 的关闭路径（表格/表单 transport、字段模型与 response mapping 提升为 Schema 驱动的通用适配层，records 降为注册实例）固化为可实施、可验收的版本化契约。冻结后 `I-010-001` 由 GOAL-010 **D-002** 置为 `verified`，解除 S1 方案冻结与 S2 实施门禁；§6 迁移策略随本契约冻结 `I-010-002`（最晚需要阶段 S3 首个前端变更前，提前关闭）。
> **v0.2.0（2026-08-03 · A-001 F-001/F-002 响应，GOAL-010 D-003）**：冻结 `dataSource` 单斜杠同源执行规则（§2，认证 fetch 前校验 + 反例）与 `rowKey` 行键不变量（§3，非空唯一标量、无效响应停止渲染并禁行 action）；S3 前端适配层 + 正反测试关闭 F-001/F-002。`I-010-001` 维持 `verified`（修订不改变冻结结论）。
> **不是**：S2～S5 的实施成品（handler/前端代码、新 fixture、回归证据属实施）；也不扩大 `I-PROTO-001 v0.1.3` 协议覆盖。
> **依据**：Root A-002 F-002-001 原文与建议关闭路径；GOAL-010 D-001 用户裁决（通用适配层改造，不降级 VP-002 主张）；[I-007-001 记录 API 契约](../../GOAL-007-r4-schema-crud/attachments/I-007-001-api-error-contract.md) v0.2.0（records 权威契约）；[I-007-003 v0.2.2](../../GOAL-007-r4-schema-crud/attachments/I-007-003-schema-crud-interaction.md) §9（action/form 写路径）；当前 `schema-table.tsx` / `records.ts` / `records.go` / `fixtures/schema/*.json` 静态核对。

## 1. 目的与范围

- 目标：**新增业务页面只修改 Schema（fixture JSON），不修改前端 Renderer 主路径**（VP-002 成功标准 1/4/6）。
- 范围：表格 list transport、表单/action 写 transport 的 response mapping、后端资源 CRUD 注册形态、错误 envelope 扩展边界。
- 保持冻结（不重开）：I-007-001 的 records 五字段形状、`{items,total,page,pageSize}` envelope、`{error,message}` envelope、`records.read`/`records.write` 权限键、毫秒 `updatedAt`；I-007-003 §9 的 action 声明（submitAction/rowAction/`{id}` 槽绑定/confirm）。

## 2. 资源标识（dataSource）

- **契约**：`table.props.dataSource` 为**协议相对 URL 字符串**（如 `/api/records`、`/api/catalog`），即该资源的 list 端点；禁止改写为资源名/映射键。
- **理由**：现状四个 fixture 均为 URL 形态；前端 transport 以 URL 直接 GET list；写端点由页面 action（`submitAction` / rowAction 的 `url`）显式声明（I-007-003 §9 已冻结），不引入第二套资源名解析层（避免双真相）。
- 默认：`dataSource` 缺省时表节点**不渲染数据**（fail-closed，`No data source` 提示），不再静默回落 `/api/records`。
- **执行规则（v0.2.0 · A-001 F-001）**：`dataSource` 必须匹配 `^/(?!\/)[^\s\\?#]*$` —— 仅允许**单斜杠同源绝对路径**：以单个 `/` 开头、非 `//`（禁止 protocol-relative host）、无 scheme（`http:`/`https:` 等）、无空白、无反斜杠、无 `?`、无 `#`。query 一律由 transport 的 `buildRecordsQuery` 追加（`dataSource` 内不得写 `?`）；fragment 禁止。`schema-table` 与 `fetchRecords` 必须在调用认证 transport（`authFetch`）**之前**校验该规则：不合法的 `dataSource` → **fail-closed**（不发起请求、不渲染数据、渲染可观察错误）。规则与 `DataRef.url` 的 `^/(?!\/)[^\s\\]*$` 一致且更严格（追加禁 `?`/`#`）。写端点由页面 action 显式声明，经 `request-construction` 的 URL 规则校验（I-007-003 §9），本规则约束 list transport。

## 3. 列表 envelope 与字段模型

- **统一列表 envelope（跨资源冻结）**：`{ "items": [...], "total": number, "page": number, "pageSize": number }`。
  - `items`：任意 JSON 对象数组；**不再**要求固定 `id/name/status/owner/updatedAt` 五字段白名单（解除 `RecordItem` 强制解析）。
  - `total`/`page`/`pageSize`：数值，语义同 I-007-001（`pageSize` 上限 100 属后端资源契约，见 §4）。
- **行键**：`table.props.rowKey`（string，默认 `"id"`）声明行唯一键，供表格 `rowKey` 与 recordView/预填选中；行对象整体透传给 action 构造器（`$row.*` 解析已有）。
- **行键执行规则（v0.2.0 · A-001 F-002）**：`rowKey` 为**直接字段名**（非空 string，默认 `"id"`），不是路径表达式，也不定义默认字段以外的推导。每行在该字段上必须为**非空且唯一的 JSON 标量**——允许 **string** 或 **finite number**；禁止 boolean/null/object/array/空串。**无效响应（缺失、空、非标量或重复键）必须 fail-closed：停止渲染数据行、禁止行 action（编辑/删除等）与选中态、渲染可观察错误**。表格的 React key、选中态与行 action 关联键统一使用该校验后的键值。
- **列模型**：沿用 `columns[].field/label/sortable`（现状）；`field` 仅声明展示与排序目标，不构成服务端字段白名单。`recordView` 直接用行对象条目渲染（现状）。
- **搜索**：搜索表单 `q` 绑定 table query（现状）；后端资源声明是否支持 `q`（§4），前端不假设字段。

## 4. 后端通用资源 CRUD 注册形态

- **资源定义（注册表条目）**，Go 侧 `Resource` 结构：
  | 字段 | 说明 |
  |------|------|
  | `id` | 小写资源 id（`records`、`catalog`…） |
  | `path` | 挂载路径（records = `/api/records`） |
  | `listable` | 是否暴露 list；`sortFields` 白名单（空 = 不可排序）、`qSearch` 布尔（是否支持 `q`） |
  | `entity` | 通用 store 接口：`List(filter) / Get(id) / Create(body) / Update(id, body) / Delete(id)`（或等价泛型仓库签名） |
  | `createFields` / `patchFields` | 必填 / 可编辑字段声明（trim 非空校验规则由实现复用） |
  | `permissionRead` / `permissionWrite` | 权限键：默认派生 `{id}.read` / `{id}.write`；records 显式保持 `records.read` / `records.write` |
- **通用 handler 工厂**：由资源定义生成 list/create/detail/update/delete 五路由，挂 `{path}` 与 `{path}/{id}`；统一 `requirePermission`、body 上限（4 KiB 保持）、`{error,message}` 写错误、`INTERNAL` 兜底。
- **records 注册实例**：`records` 注册到 `/api/records`，`sortFields = [name,status,owner,updatedAt]`、`qSearch = true`、create/patch 字段 `name/status/owner`、权限键 `records.read/write`——**对外 HTTP 契约与 I-007-001 逐项一致（零 API 变更）**；实现从手写 handler 收敛为注册条目 + 通用工厂（内部行为不变）。
- **新资源（S4 验证实体，示例 `catalog`）**：注册新条目 + 迁移/种子 + 权限键 `catalog.read`/`catalog.write` 注入种子 grants；fixture 的 `dataSource` 指向其路径。

## 5. 错误 envelope 与错误码

- **Envelope 冻结（全资源）**：`{ "error": "<STABLE_CODE>", "message": "<text>" }`；**不新增字段**。
- **通用错误码（全资源共享）**：`UNAUTHENTICATED`(401)、`FORBIDDEN`(403)、`INVALID_SORT_FIELD`/`INVALID_SORT_ORDER`/`INVALID_PAGE`/`INVALID_PAGE_SIZE`(400)、`INVALID_CREATE_BODY`/`INVALID_CREATE_FIELD`(400)、`INVALID_PATCH_BODY`/`INVALID_PATCH_FIELD`(400)、`INTERNAL`(500)。
- **资源特定**：NOT_FOUND 码 = `{ID}_NOT_FOUND`（records 保持 `RECORD_NOT_FOUND` 兼容；新资源如 `CATALOG_NOT_FOUND`）。前端 `readRecordApiError` 已按 envelope 泛读，无需改动。
- **不引入**：`409`/业务唯一冲突、枚举校验码（与 I-007-001 一致）；resources 元数据发现端点（如 `/api/resources` 列表）为非目标（注册表是代码级，不暴露）。

## 6. 迁移与兼容策略（I-010-002）

- **后端**：records handler 收敛为注册实例；路由、权限键、envelope、错误码、毫秒 `updatedAt`、body 上限全部保持——**对外零 API 变更**；回归由现有 records_test / T-API-01～13 承担。
- **前端**：`records.ts` 的 `RecordItem`/`RecordList` 固定解析迁移为通用 `ResourceItem = JsonRecord` + 统一 list 解析（envelope 不变）；`schema-table` 不再调用 `fetchRecords` 固定形状；`DEFAULT_RECORDS_URL` 回落逻辑删除（缺 `dataSource` fail-closed）。`readRecordApiError`/envelope 泛读保留。v0.2.0 起 `records.ts` 新增 `isValidDataSource`（§2 F-001，认证 fetch 前校验），`schema-table.tsx` 实施 `rowKey` 不变量（§3 F-002）。
- **fixture / 测试**：现有 fixture `dataSource: "/api/records"` 全部保持有效；`schema-crud.test.tsx` emulator 形状（`{items,total,page,pageSize}` + `{error,message}`）不变；`records.test.ts` 迁为通用解析用例。
- **双轨边界**：不提供「旧固定解析」并行层——一次性迁移 + 全量回归（S3）；迁移期间不新增业务 fixture。
- **权限**：新资源必须在其门禁范围内完成种子 grants 注入与匿名 401 / 缺权限 403 测试（S4 必测）。

## 7. 路线图范围映射

| 检查点 | 本契约对应 |
|--------|-----------|
| S2 后端通用资源 CRUD | §4 注册表 + 通用工厂；records 实例化 |
| S3 前端通用适配层 | §2/§3/§5/§6 前端泛化 |
| S4 新实体验证 | §4 新资源条目 + fixture 接入 |
| S5 回归、审计与关闭 | §6 兼容口径 + Root A-002 F-002-001 关闭证据 |

## 8. 非目标

- 协议层动态资源发现（`/api/resources` 元数据端点）、通用「任意表任意字段」泛 CRUD（无 schema 声明的自由写）——**不做**；资源必须显式注册。
- 扩大 `I-PROTO-001 v0.1.3` 覆盖；批量 action/上传；乐观锁/软删除；多租户。
- 不改变 I-007-003 §9 action 语义（`submitAction`/rowAction/`{id}` 槽绑定/confirm）。

## 9. 证据索引

- `apps/web/src/renderer/schema-table.tsx`（`schemaTableDataSource` F-001 校验、`schemaTableRowKey`/`checkRowKeys` F-002——S3 改造对象）
- `apps/web/src/renderer/records.ts`（`isValidDataSource`/`DATASOURCE_URL_PATTERN`、`parseRecordList`/`ResourceItem` 泛化——S3 改造对象）
- `apps/api/internal/handler/records.go`（硬编码路由/实体——S2 收敛对象）
- `apps/api/internal/handler/fixtures/schema/{list-edit-lifecycle,search-form-table,catalog,data-table}.json`（`dataSource` 现状）
- [I-007-001 v0.2.0](../../GOAL-007-r4-schema-crud/attachments/I-007-001-api-error-contract.md)（records 权威契约；本契约 §4/§5 的兼容基线）
- [I-007-003 v0.2.2](../../GOAL-007-r4-schema-crud/attachments/I-007-003-schema-crud-interaction.md) §9（action/写路径冻结）

## 10. 修订记录

| 版本 | 日期 | 变更 |
|------|------|------|
| 0.1.0 | 2026-08-03 | 冻结（GOAL-010 D-002；关闭 `I-010-001`；§6 迁移策略一并冻结 `I-010-002`） |
| 0.2.0 | 2026-08-03 | A-001 F-001/F-002 响应（GOAL-010 D-003）：§2 冻结 `dataSource` 单斜杠同源执行规则（认证 fetch 前校验 + 反例）；§3 冻结 `rowKey` 行键不变量（非空唯一标量、无效响应停止渲染并禁行 action）；S3 前端适配层实施 + 正反测试按 `fixed` 关闭 F-001/F-002 |
