---
title: I-007-001 · records 精确 API 与错误契约
status: active
doc_type: info-collection
created: 2026-08-02
updated: 2026-08-02
parent: GOAL-007-r4-schema-crud
version: 0.2.0
related_info: I-007-001
related_decision: D-002
---

# I-007-001 · records 精确 API 与错误契约

> **结论**：本附件与 D-002 关闭 `I-007-001`，并完成成功标准 **S1**（契约冻结）。以下字段、端点、HTTP status、稳定 `code` 与正反矩阵是 R4 实施输入，**不是**已交付的 create/SQLite 产品事实。
> **扫描日期**：2026-08-02。工作区 `shared_materials_catalog: none`；事实来自本仓库代码、测试与 Root [I-004-schema-crud-collection.md](../../GOAL-001-production-admin-foundation/attachments/I-004-schema-crud-collection.md)。
> **修订（v0.2.0 · 响应 A-001 F-001）**：`updatedAt` 精度由 RFC3339 秒级统一为**含毫秒**；「严格晚于」断言保留并以单调钳制保证确定性。修订决策见 D-004。

## 1. 当前可继承基线

| 事实 | 证据 | R4 继承 |
|------|------|---------|
| 实体五字段：`id`、`name`、`status`、`owner`、`updatedAt` | `apps/api/internal/handler/records.go` `record` | 保持对外 JSON 形状；不新增 `createdAt` |
| 列表：`GET /api/records`，查询 `q`/`sort`/`order`/`page`/`pageSize`；响应 `{ items, total, page, pageSize }` | `list` / `recordList` | 继承；默认 `sort=name`、`order=asc`、`page=1`、`pageSize=10`；`pageSize` 上限 100 |
| 可排序字段：`name`、`status`、`owner`、`updatedAt` | `recordSortFields` | 继承 |
| 搜索：`q` 对 `name`/`status`/`owner` 做大小写不敏感子串匹配 | `matches` | 继承；**不**搜 `id` |
| 详情：`GET /api/records/{id}` → 200 + record | `detail` | 继承 |
| 更新：`PATCH /api/records/{id}`，body 指针字段 `name`/`status`/`owner`；成功 200 + 完整 record；刷新 `updatedAt` | `update` / `recordPatch` / `validatePatch` | 继承；可编辑字段仅此三 |
| 删除：`DELETE /api/records/{id}` → 204 空体 | `delete` | 继承 |
| 权限：读 `records.read`，写 `records.write`；匿名 401，缺权限 403 | `requirePermission` | 继承；**POST 归写权限** |
| 错误 envelope：`{"error":"<CODE>","message":"<text>"}` | `writeError` in `health.go` | 继承；不改字段名 |
| body 上限：4 KiB（`maxRecordBodyBytes`） | `records.go` const | create/update 共用 |
| 当前无 `POST`；数据为进程内 8 条静态切片 | `routes` / `staticRecords` | R4 必须补 create，并迁出进程切片（见 I-007-002） |

## 2. 冻结字段与 ID / 时间戳语义

| 字段 | 类型（JSON） | 来源 | 约束 |
|------|--------------|------|------|
| `id` | string | **仅服务端**生成 | 非空；主键；客户端 create/patch **不得**指定或改写 |
| `name` | string | 客户端 create 必填；patch 可选 | trim 后非空 |
| `status` | string | 客户端 create 必填；patch 可选 | trim 后非空；**不做枚举白名单**（与现 PATCH 一致；种子值可为 `active`/`archived`/`pending`） |
| `owner` | string | 客户端 create 必填；patch 可选 | trim 后非空 |
| `updatedAt` | string（RFC3339 **含毫秒**） | **仅服务端** | UTC；create 时写入；每次成功 update 刷新为 `time.Now().UTC()`；客户端不得写入。格式固定 `2006-01-02T15:04:05.000Z07:00`（如 `2026-08-02T03:04:05.123Z`）。**精度**：Unix 毫秒级（DB 存储），同一毫秒内连续更新由单调钳制保证严格递增（见 §3.1） |

### ID 生成（create）

- 格式：`rec-` + 16 位小写 hex（8 字节 `crypto/rand`）。
- 碰撞：若唯一约束冲突，重试有限次数后以 `INTERNAL` 失败（不暴露内部细节）；不回落为可预测序号。
- 既有种子 id（`rec-1`…`rec-8`）可继续作为演示数据；**不**要求新 create 使用递增整数。

### 未纳入本契约

- 乐观锁 / `If-Match` / `version` 字段 → R4 **不做**；并发语义为 **last-write-wins**（见 I-007-002）。
- `name` 全局唯一 → **不做**。
- soft-delete → **不做**；DELETE 为物理删除。

## 3. 端点契约

### 3.1 继承端点（行为不变，数据源改 SQLite）

| Method | Path | Permission | 成功 | 失败（稳定 code） |
|--------|------|------------|------|-------------------|
| GET | `/api/records` | `records.read` | 200 `recordList` | `UNAUTHENTICATED` 401；`FORBIDDEN` 403；`INVALID_SORT_FIELD` / `INVALID_SORT_ORDER` / `INVALID_PAGE` / `INVALID_PAGE_SIZE` 400 |
| GET | `/api/records/{id}` | `records.read` | 200 `record` | 401/403 同上；`RECORD_NOT_FOUND` 404 |
| PATCH | `/api/records/{id}` | `records.write` | 200 `record` | 401/403；`INVALID_PATCH_BODY` 400；`INVALID_PATCH_FIELD` 400；`RECORD_NOT_FOUND` 404；body 超限 → 400（与现 MaxBytesReader 行为一致，message 可实现侧） |
| DELETE | `/api/records/{id}` | `records.write` | **204** 无 JSON 体 | 401/403；`RECORD_NOT_FOUND` 404 |

PATCH 细节（冻结）：

1. 缺省键 = 不修改该字段；显式 `null` 按 JSON decode 为 Go `nil` 指针 → 不修改（与现 `*string` 语义一致）。
2. 提供的字符串经 trim 后为空 → `INVALID_PATCH_FIELD`，message 形如 `name must not be empty`。
3. 未知 JSON 键：当前 decoder 默认忽略；R4 **保持忽略**，不引入 `DisallowUnknownFields`（避免破坏既有客户端）。
4. 成功后 `updatedAt` 必须**严格晚于**更新前值，且随后 GET 一致。**确定性保证**（毫秒精度）：写入时若该行新时间戳 ≤ 前一 `updated_at`，则钳制为 `prev + 1ms`（单调钳制，仅同一毫秒内快速连续更新触发）；禁止人为跳秒/跳毫秒，也不退回 Unix 秒。测试无需 sleep 即可稳定断言严格递增。

### 3.2 新增 create

| Method | Path | Permission | 成功 | 失败 |
|--------|------|------------|------|------|
| POST | `/api/records` | `records.write` | **201** + 完整 `record` | 见下表 |

**请求 body**（全部必填，无 pointer 可选语义）：

```json
{
  "name": "string",
  "status": "string",
  "owner": "string"
}
```

| 场景 | HTTP | code | message 口径（稳定意图，文案可微调） |
|------|------|------|--------------------------------------|
| 无身份 | 401 | `UNAUTHENTICATED` | 与现 gate 一致 |
| 缺 `records.write` | 403 | `FORBIDDEN` | `permission required: records.write` |
| body 非 JSON / 截断 / 超 4KiB | 400 | `INVALID_CREATE_BODY` | body must be JSON（超限同属 body 失败） |
| 缺字段、非字符串、或 trim 后空 | 400 | `INVALID_CREATE_FIELD` | 指明字段，如 `name must not be empty` |
| 服务端 ID 生成最终失败 | 500 | `INTERNAL` | 不暴露内部细节 |

**成功响应**：201，`Content-Type: application/json`，body 为完整 record（含服务端 `id` 与 `updatedAt`）。

**Location 头**：R4 **不要求**；客户端以 body.id 为准。

## 4. 统一错误 envelope 与 code 全表

Envelope（全 API 一致）：

```json
{ "error": "<STABLE_CODE>", "message": "<human text>" }
```

| code | HTTP | 适用操作 | 来源 |
|------|------|----------|------|
| `UNAUTHENTICATED` | 401 | 全部 | 已有 |
| `FORBIDDEN` | 403 | 全部 | 已有 |
| `INVALID_SORT_FIELD` | 400 | list | 已有 |
| `INVALID_SORT_ORDER` | 400 | list | 已有 |
| `INVALID_PAGE` | 400 | list | 已有 |
| `INVALID_PAGE_SIZE` | 400 | list | 已有 |
| `RECORD_NOT_FOUND` | 404 | detail / patch / delete | 已有 |
| `INVALID_PATCH_BODY` | 400 | patch | 已有 |
| `INVALID_PATCH_FIELD` | 400 | patch | 已有 |
| `INVALID_CREATE_BODY` | 400 | **create（新）** | 本契约冻结 |
| `INVALID_CREATE_FIELD` | 400 | **create（新）** | 本契约冻结 |
| `INTERNAL` | 500 | create（ID 穷尽等稀有失败） | 对齐 account 等既有 INTERNAL 用法 |

**明确不引入（R4）**：

- `RECORD_CONFLICT` / `409`：无业务唯一约束，不需要。
- `INVALID_STATUS` 枚举错误：status 非枚举。
- 将 create 错误复用 `INVALID_PATCH_*`：禁止；create/patch code 分离，便于测试与前端映射。

## 5. 正反测试矩阵（S1 / 承接 S3）

| ID | 断言 | 基线 | R4 要求 |
|----|------|------|---------|
| T-API-01 | 默认 list：8 种子（或当前库全量）、`pageSize=10`、name 升序首条可预期 | 已有 `TestRecordsListDefault` | 数据改 SQLite 后语义不变 |
| T-API-02 | `q` / `sort` / `order` / `page` 行为 | 已有 | 保持 |
| T-API-03 | 非法 sort/order/page/pageSize → 400 + 对应 code | 已有 | 保持 |
| T-API-04 | detail 命中 / `RECORD_NOT_FOUND` | 已有 | 保持 |
| T-API-05 | PATCH 成功、`updatedAt` **严格晚于**前值（毫秒精度、单调钳制）、后续 GET 一致 | 已有 | 持久化后跨进程仍一致（S6） |
| T-API-06 | PATCH 空字段 / 非法 JSON / 过大 body → 400 | 已有 | 保持 |
| T-API-07 | DELETE 204，随后 list 减少且 detail 404 | 已有 | 持久化后保持 |
| T-API-08 | 匿名 PATCH/DELETE/POST → 401 `UNAUTHENTICATED` | 写已有；POST 待加 | 必测 |
| T-API-09 | viewer/editor 无 write → 403 `FORBIDDEN`（含 POST） | PATCH/DELETE 已有 | POST 必测 |
| T-API-10 | POST 合法 body → 201，返回 id/name/status/owner/updatedAt；list/detail 可见 | **新** | S3 必测 |
| T-API-11 | POST 缺字段/空字符串 → 400 `INVALID_CREATE_FIELD` | **新** | S3 必测 |
| T-API-12 | POST 非法 JSON / 过大 → 400 `INVALID_CREATE_BODY` | **新** | S3 必测 |
| T-API-13 | admin 持 `records.read`+`records.write` 可走全链路 | 种子 grants 已有 | 保持 |

对应 M-R4：`01`–`06`、`08`（错误码侧）、`09`（API 层）。Schema 交互与重启协议分别由 `I-007-003` / `I-007-004` 承接。

## 6. 与 I-007-002 的接口

- Handler 不再持有进程内 `[]record` 作为生产路径；读写经 store repository（见 I-007-002）。
- JSON `updatedAt` 为 RFC3339 **含毫秒**；DB 列 Unix 毫秒，映射层负责转换（见 I-007-002 v0.2.0）。
- 生产默认路径必须在进程重启后仍能通过 T-API-05/07/10；证明方式由 I-007-004 冻结。

## 7. 证据索引

- `apps/api/internal/handler/records.go`
- `apps/api/internal/handler/records_test.go`
- `apps/api/internal/handler/health.go`（`writeError`）
- `apps/api/internal/store/seed.go`（`records.read` / `records.write` grants）
- Root [I-004-schema-crud-collection.md](../../GOAL-001-production-admin-foundation/attachments/I-004-schema-crud-collection.md) M-R4-01～06/08/09
