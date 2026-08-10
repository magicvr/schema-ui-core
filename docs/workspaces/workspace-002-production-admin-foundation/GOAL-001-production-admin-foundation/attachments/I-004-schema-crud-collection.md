---
title: I-004 · Schema CRUD 代表实体与 API/错误语义信息收集
status: active
doc_type: info-collection
created: 2026-08-02
updated: 2026-08-02
parent: GOAL-001-production-admin-foundation
version: 0.2.0
related_info: I-004
related_decision: D-010
---

# I-004 · R4 Schema CRUD 代表实体信息收集

> **性质**：回答“哪个代表性实体及 API/错误语义能够完整证明 Schema CRUD 闭环”所需的当前实现事实、候选比较、差量与验收矩阵。
> **裁决状态**：用户于 2026-08-02 确认采用 `records`，要求 SQLite 持久化与重启保持，并沿用统一错误 envelope；精确 `code` 在 R4 子目标方案中冻结。Root D-010 已将 `I-004` 置为 `verified`。
> **不是**：R4 详细实施方案冻结、R4 子目标立项、产品实现或验收通过。本轮没有修改产品代码。
> **扫描日期**：2026-08-02。工作区 `shared_materials_catalog: none`，全部事实来自本仓库代码、fixture 与本地测试。

## 0. 总览结论

| 维度 | 当前事实 | 对 R4 的含义 |
|------|----------|---------------|
| 代表实体 | `record`：`id`、`name`、`status`、`owner`、`updatedAt` | D-010 已采用；继承列表 API、Schema table、详情 API、PATCH、DELETE 与读写 permission key 基线 |
| 当前数据边界 | 八条静态开发数据；变更仅保存在 `recordHandler.records` 进程内切片 | 不能证明生产持久化或重启后 CRUD 结果，必须纳入 R4 差量 |
| 当前 API | `GET /api/records`、`GET /api/records/{id}`、`PATCH /api/records/{id}`、`DELETE /api/records/{id}` | 缺少新建端点；现有读、改、删语义可作为目标矩阵输入 |
| Schema 页面 | `data-table` / `search-form-table` 可请求真实列表；`list-edit-lifecycle` 含 `recordView` + `form` | 详情数据仍内联在 fixture，表单未绑定真实提交，尚无 Schema 驱动的创建/编辑/删除闭环 |
| 状态与错误 | table 已实现 loading / empty / fetch error / sort；后端已有参数、body、字段、not-found 错误 | 仍需统一创建/编辑/删除的成功反馈、服务端字段错误映射和动作失败呈现 |
| 权限 | 读取要求 `records.read`，写入要求 `records.write`；匿名 `401`，缺权限 `403` | 后端权限基线已存在；R4 仍需用真实 Schema 页面证明允许与拒绝路径 |
| 本轮验证 | API records 定向包测试通过；Web 相关 3 文件 `20/20` 通过 | 只证明现有候选基线可复核，不证明缺失能力已经实现 |

## 1. 当前 `records` 实体与 API 事实

| 事实 | 证据 |
|------|------|
| `record` 字段为 `id`、`name`、`status`、`owner`、`updatedAt` | `apps/api/internal/handler/records.go` `record` |
| 后端初始化八条稳定静态数据；注释明确无 records DB，mutation 仅存活于进程 | `apps/api/internal/handler/records.go` `staticRecords` / `recordHandler` |
| 列表支持 `q`、`sort`、`order`、`page`、`pageSize`，响应为 `{ items, total, page, pageSize }` | `apps/api/internal/handler/records.go` `list` / `recordList` |
| 详情、编辑、删除分别为 `GET`、`PATCH`、`DELETE /api/records/{id}`；当前没有 `POST /api/records` | `apps/api/internal/handler/records.go` `routes` |
| `PATCH` 只允许 `name`、`status`、`owner`，空白值被拒绝，成功后更新 `updatedAt` | `apps/api/internal/handler/records.go` `recordPatch` / `validatePatch` / `update` |
| 所有 records 路由经过请求身份中间件；读写分别要求 `records.read` / `records.write` | `apps/api/internal/handler/records.go` `requirePermission` 与各 handler |
| records 测试覆盖列表、搜索、排序、分页、详情、编辑、删除、401、403、404、非法参数、非法 JSON、空字段和过大 body | `apps/api/internal/handler/records_test.go` |

## 2. 当前错误与权限语义

| 场景 | HTTP | code | 当前状态 |
|------|------|------|----------|
| 无有效请求身份 | 401 | `UNAUTHENTICATED` | 已实现、已有测试 |
| 已认证但缺少目标 permission | 403 | `FORBIDDEN` | 已实现、已有测试 |
| 不支持的排序字段 | 400 | `INVALID_SORT_FIELD` | 已实现、已有测试 |
| 不支持的排序方向 | 400 | `INVALID_SORT_ORDER` | 已实现、已有测试 |
| 非法 page / pageSize | 400 | `INVALID_PAGE` / `INVALID_PAGE_SIZE` | 已实现、已有测试 |
| 详情、编辑或删除目标不存在 | 404 | `RECORD_NOT_FOUND` | 已实现、已有测试 |
| PATCH body 不是合法 JSON | 400 | `INVALID_PATCH_BODY` | 已实现、已有测试 |
| PATCH 可编辑字段为空 | 400 | `INVALID_PATCH_FIELD` | 已实现、已有测试 |
| 新建 body / 字段错误 | 待冻结 | 待冻结 | 当前无 create API，不得从 PATCH 语义静默推导 |
| 持久化/并发冲突 | 待冻结 | 待冻结 | 当前进程内切片没有相应语义 |

## 3. Schema 与 Web 现状

| 闭环环节 | 当前证据 | 判定 |
|----------|----------|------|
| 列表 | `SchemaTable` 按 `dataSource` 请求 `/api/records` 并渲染列、分页摘要 | 已有真实数据读取基线 |
| 搜索 | `search-form-table.json` 声明搜索表单与 records 数据源；后端 `q` 已实现 | 已有分离能力，仍需 R4 证明表单参数实际驱动列表 |
| 加载 / 空态 / 错误态 | `SchemaTable` / `DataTable` 传入 `loading`、`error`、`emptyMessage`；测试覆盖 fetch error 与无列 fail-closed | 已有列表状态基线 |
| 详情 | 后端 detail API 已有；`list-edit-lifecycle.json` 的 `recordView` 数据仍直接内联在 fixture | 缺真实 detail load 绑定 |
| 新建 | 未发现 records POST API 或 Schema 新建动作 | 缺失 |
| 编辑 | 后端 PATCH 已有；fixture form 只有字段与 `submitLabel` | 缺 form → PATCH → 成功/错误反馈绑定 |
| 删除 | 后端 DELETE 已有 | 缺 Schema action → DELETE → 确认/成功/错误反馈绑定 |
| 权限失败 | 后端 401/403 与持久化 permission key 已有 | 缺 Schema 页面上的允许/拒绝端到端证据 |
| 重启持久化 | R3 restart 测试覆盖 users/RBAC/menu/refresh token，不覆盖 records | 缺失 |

## 4. 候选比较

| 候选 | 现有优势 | 主要缺口 | 收集结论 |
|------|----------|----------|----------|
| `records`（已采用） | 唯一同时具备列表/搜索/详情/PATCH/DELETE、Schema fixtures、loading/empty/error、401/403 与集中测试的业务候选 | 缺 POST、SQLite 持久化、重启证据和 Schema 真实写动作 | D-010 已采用；缺口全部进入 R4 子目标方案与验收 |
| users / RBAC | 已有 SQLite、迁移、增量 seed、permission/menu 投影和重启证据 | 对外仅 `/api/accounts/me`；没有用户 CRUD API 或 Schema CRUD 页面 | 更适合作为身份/权限支撑域，不适合作为现成代表实体 |
| 新业务实体 | 可从零设计完整生产语义 | 仓库没有现成 API、模型、Schema 或测试；会重复建立已有 records 基线 | 除非用户明确要求真实业务域，否则不建议扩大 R4 |

## 5. R4 方案冻结前的最低验收矩阵

| ID | 必须固定或证明的行为 | 当前基线 | 冻结前状态 |
|----|----------------------|----------|------------|
| M-R4-01 | 代表实体、字段、可编辑字段、ID 与时间戳语义 | D-010 采用 `records` 与现有五字段基线 | 子目标冻结字段约束 |
| M-R4-02 | list/search/sort/page 请求与响应 envelope | 已实现并有 API 测试 | 可继承，仍需端到端联动 |
| M-R4-03 | create 请求、成功响应、默认值与字段校验 | 无 POST | 待冻结 |
| M-R4-04 | detail 请求与 not-found | API 已有 | 待接 Schema 真实 load |
| M-R4-05 | update 方法、并发/更新时间语义与字段错误 | PATCH 基线已有 | 持久化/并发语义待冻结 |
| M-R4-06 | delete 确认、成功响应、not-found 与后续列表一致性 | DELETE `204` 基线已有 | 待接 Schema action |
| M-R4-07 | SQLite 持久化、migration/seed 及重启保持 | D-010 设为 required；records 当前无 DB | 子目标冻结精确 DDL/migration/seed 与证据 |
| M-R4-08 | loading、empty、统一 error、成功反馈与字段错误呈现 | list 状态已有 | write lifecycle 待冻结 |
| M-R4-09 | `records.read` / `records.write` 的允许与 401/403 拒绝路径 | 后端已有 | 待补 Schema 端到端证据 |
| M-R4-10 | 只改 Schema 即可挂接代表页面，不改 Renderer 主路径 | R1 fixtures/Renderer 基线已有 | 待 R4 页面证明 |
| M-R4-11 | API 与 Web 定向回归、重启回归和失败路径证据 | 本轮基线测试通过 | R4 实施后重跑并扩展 |

## 6. 本轮验证证据

| 命令 | 工作目录 | 结果 |
|------|----------|------|
| `go test ./internal/handler -run Records -count=1` | `apps/api` | pass；package `github.com/magicvr/schema-ui-core/apps/api/internal/handler` |
| `npm test -- --run src/renderer/schema-table.test.tsx src/renderer/representative-pages.test.tsx src/app/representative-pages.integration.test.tsx` | `apps/web` | pass；3 files，20/20 tests |

这些结果验证当前候选基线仍可复现；它们不覆盖 create、records 持久化、Schema 写动作或 R4 端到端验收。

## 7. D-010 裁决与剩余门禁

1. **采用 `records`** 作为 R4 代表实体，不扩展 users/RBAC 或新建业务域。
2. **SQLite 持久化与重启保持为 required**：create/update/delete 后重启，list/detail 结果必须保持；精确 DDL/migration/seed 与失败证据在 R4 子目标方案中冻结。
3. **沿用统一错误 envelope**：精确 create、字段校验、持久化/并发冲突及 Schema action error code 在 R4 子目标方案中冻结，不在 Root 层臆造。

`I-004` 已由 D-010 置为 `verified`，只解除 Root 的 R4 方向冻结与子目标立项目门禁。未来 R4 子目标在创建时必须把精确 API/error code、SQLite schema/migration/seed、并发语义、Schema 端到端权限与重启证据登记为实施前 required；未关闭前不得实施受影响范围或验收 R4。R1～R3 结论不变。

## 8. 证据索引

- `apps/api/internal/handler/records.go`
- `apps/api/internal/handler/records_test.go`
- `apps/api/internal/handler/fixtures/schema/data-table.json`
- `apps/api/internal/handler/fixtures/schema/search-form-table.json`
- `apps/api/internal/handler/fixtures/schema/list-edit-lifecycle.json`
- `apps/api/internal/handler/account.go`
- `apps/api/internal/store/migrate.go`
- `apps/api/internal/store/restart_test.go`
- `apps/web/src/renderer/schema-table.tsx`
- `apps/web/src/renderer/schema-table.test.tsx`
- `apps/web/src/renderer/representative-pages.test.tsx`
- `apps/web/src/app/representative-pages.integration.test.tsx`
