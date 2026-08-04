---
id: r4-c1-capability-inventory
doc: evidence
goal: GOAL-005-r4-full-module-migration
date: 2026-08-05
status: recorded
---

# R4 C1 能力与边界事实盘点

## 盘点结论

本附件记录 C1 阶段截至 2026-08-05 的代码事实和迁移边界。它是
`R4-I001` 的核验材料，不把尚未形成决策的范围写成已冻结方案；R4-I002、
R4-I003、R4-I004 仍需决策或契约证据后才能关闭。

## 能力映射

| 能力/资料面 | 当前所有权或注册事实 | C1 处理结论 | 关键证据 |
|---|---|---|---|
| HTTP server、health、readyz、auth/session、accounts | Composition Root 调用中心 `handler.Register`; core 路由在中心 handler | 保留为 core composition boundary，需在 C2 明确其与模块 provider 的边界 | `apps/api/internal/composition/composition.go:82-111`; `apps/api/internal/handler/health.go:20-35` |
| Manifest route、navigation、page projection | Manifest 聚合器从中心 `adminModules` 和静态 manifest 投影；Web 从 manifest 投影导航 | 作为统一 provider 聚合的现有基线，C2 需移除一方模块的硬编码 owner 映射 | `apps/api/internal/manifest/manifest.go:38-63,67-159`; `apps/web/src/app/navigation.ts:137-152` |
| Schema render | 中心 handler embed schema，并用 `schemaDocumentsForPlan` 的 owner map 过滤 | C2/C3 迁移到模块拥有的 Schema provider；当前 owner map 是待删除中心特例 | `apps/api/internal/handler/schema.go:13-18,36-58` |
| Users | `admin.users` 声明 page/navigation/permission；HTTP Resource 和 Store 仍由中心 handler/store 提供 | C3 迁移，行为与协议保持兼容 | `apps/api/internal/kernel/profile.go:99-100`; `apps/api/internal/handler/users.go:34-51`; `apps/api/internal/store/users.go:19-52` |
| Roles | `admin.roles` 声明 page/navigation/permission；HTTP Resource 和 Store 仍由中心 handler/store 提供 | C3 迁移，行为与协议保持兼容 | `apps/api/internal/kernel/profile.go:99-100`; `apps/api/internal/handler/roles.go:21-39`; `apps/api/internal/store/roles.go:16-60` |
| Settings | 已有 `settings` module 包入口，但入口仍委托中心 `handler.RegisterSettings` | C2/C4 迁移为真实 module-owned provider | `apps/api/internal/modules/settings/module.go:9-14`; `apps/api/internal/handler/settings.go` |
| Activity/operation query | 已有 `activity` module 包入口，但入口仍委托中心 `handler.RegisterActivity`; operation query 只读 | C2/C4 迁移为真实 module-owned provider；operationlog 写入保持横切基础设施 | `apps/api/internal/modules/activity/module.go:9-14`; `apps/api/internal/handler/operations.go:14-26,54-85`; `docs/architecture/module-architecture.md:105` |
| Operationlog persistence | `RecordOperation` append-only；业务写入后 best-effort 记录，失败不回滚业务写入；未发现 retention/archival contract | R4-I004 必须明确保留当前语义及 retention 边界，或形成变更决策和新测试 | `apps/api/internal/store/operations.go:44-62,75-119`; `apps/api/internal/handler/users.go:273-303`; `apps/api/internal/handler/roles.go:228-256` |
| Records product CRUD | migration `0006 records_retire` 删除表、权限和菜单；当前无 `admin.records`、`/api/records`、handler 或 fixture；历史 operation-log 事件保留 | R4-I003 信息冲突仍开放；不可把退役事实解释成 VP-003 的最终范围决定 | `apps/api/internal/store/migrate.go:291-310`; `apps/api/README.md:109-110`; `attachments/r4-initial-boundary-scan.md:37-47` |

## Provider 缺口

`kernel.Module` 当前提供 capability 依赖、贡献 key 和生命周期 hooks，但没有
结构化 HTTP、Schema、Authorization、Navigation、Manifest、Persistence provider
字段。现有冲突校验只覆盖字符串贡献 key。C2 进入条件是冻结框架无关 provider
形状、Plan 消费顺序、冲突规则、依赖注入和生命周期失败清理测试。

证据：`apps/api/internal/kernel/module.go:28-56,322-386`、
`apps/api/internal/kernel/lifecycle.go:17-65`。

## Operationlog 事实

当前业务写入和 operationlog INSERT 不在同一可核对事务契约中；日志写入失败只
记录服务日志。已有 append、排序、事件 CHECK、迁移保留和 reopen 测试，但未发现
日志失败注入、原子回滚、retention duration、归档表或归档恢复测试。该事实不
替代 R4-I004 的保留/变更决策。

证据：`apps/api/internal/store/operations_test.go:11-246`、
`apps/api/internal/store/migrate.go:222-238,249-286,348-382`。

## 未关闭门禁

- `R4-I001`：能力映射已形成，但“全部现有 Admin 能力”的 C1 核验仍需审计确认。
- `R4-I002`：provider gap 已核实，contract 尚未冻结。
- `R4-I003`：VP-003 与 `records_retire` 的范围冲突必须由用户或 canonical 决策解决。
- `R4-I004`：operationlog 失败语义和 retention 边界必须形成明确决策与证据。

## Freeze-grade surface disposition

下表将当前已发布的 page/Schema、API 入口、所有权和阶段处置逐项列出。
`keep-core` 与 `consumer-only` 是当前 C1 的范围处置，不代表它们已经完成 R5
或 R6；`pending-gate` 明确表示尚不能作 include/exclude 决定。

| page/Schema | API/route | 当前 owner | 阶段处置 | 关键证据 |
|---|---|---|---|---|
| `overview` | `/api/schema/overview`, `/overview` | core fixture / `core.schema-render` + `core.manifest-route` | `keep-core`; R5 验证 | `apps/api/internal/manifest/app-manifest.json:2-20`; `apps/api/internal/handler/fixtures/schema/overview.json:3-4` |
| `data-table` | `/api/schema/data-table`, `/data-table` | core fixture | `keep-core`; R5 验证 | `apps/api/internal/manifest/app-manifest.json:21-27,74-83`; `apps/api/internal/handler/fixtures/schema/data-table.json:3-4` |
| `search-form-table` | `/api/schema/search-form-table`, `/search-form-table` | core fixture | `keep-core`; R4/R5 protocol validation | `apps/api/internal/manifest/app-manifest.json:28-34,84-91`; `apps/api/internal/handler/fixtures/schema/search-form-table.json:3-4` |
| `form-controls` | `/api/schema/form-controls`, `/form-controls` | core fixture | `keep-core`; R5 验证 | `apps/api/internal/manifest/app-manifest.json:35-41,92-99`; `apps/api/internal/handler/fixtures/schema/form-controls.json:3-4` |
| `form-with-reactions` | `/api/schema/form-with-reactions`, `/form-with-reactions` | core fixture | `keep-core`; R5 验证 | `apps/api/internal/manifest/app-manifest.json:42-48,92-100`; `apps/api/internal/handler/fixtures/schema/form-with-reactions.json:3-4` |
| `users` | `/api/schema/users`, `/api/users`, `/users` | candidate `admin.users`; current HTTP/Schema central | `R4 C3 migrate` | `apps/api/internal/kernel/profile.go:99`; `apps/api/internal/handler/users.go:34-51`; `apps/api/internal/store/users.go:19-52` |
| `roles` | `/api/schema/roles`, `/api/roles`, `/roles` | candidate `admin.roles`; current HTTP/Schema central | `R4 C3 migrate` | `apps/api/internal/kernel/profile.go:100`; `apps/api/internal/handler/roles.go:21-39`; `apps/api/internal/store/roles.go:16-60` |
| `settings` | `/api/schema/settings`, `/api/settings*`, `/settings` | `admin.settings` package/schema with central handler/store adapter | `R4 C4 migrate after C1 gates` | `apps/api/internal/modules/settings/module.go:11-16`; `apps/api/internal/modules/settings/schema/schema.go:40-47`; `apps/api/internal/store/migrate.go:316-329` |
| `activity` | `/api/schema/activity`, `/api/operations`, `/activity` | `admin.activity` package/schema; operationlog is core | `R4 C4 migrate after C1 gates` | `apps/api/internal/modules/activity/module.go:11-16`; `apps/api/internal/handler/operations.go:17-24,94-97`; `apps/api/internal/store/operations.go:44-62` |
| Records | no current product route or Schema | historical event rows only after `0006 records_retire` | `pending-gate` under R4-I003; no silent exclude | `apps/api/internal/store/migrate.go:291-313`; `apps/api/README.md:105-110`; `r4-initial-boundary-scan.md:37-47` |

## Admin module six-capability matrix

The matrix describes current ownership and the migration disposition. It is not a
claim that any provider has already been implemented.

| module | HTTP | Schema | Authorization | Navigation | Manifest | Persistence | R4 disposition |
|---|---|---|---|---|---|---|---|
| `admin.users` | central generic Resource, `/api/users` | central fixture + owner map | central `users.read/write` | metadata + central RBAC seed `menu_users` | central `adminModules` projection | central users/user_roles store and `0002` schema | C3 migrate |
| `admin.roles` | central generic Resource, `/api/roles` | central fixture + owner map | central `roles.read/write/assign` | metadata + central RBAC seed `menu_roles` | central `adminModules` projection | central roles/permission/menu joins and `0002` schema | C3 migrate |
| `admin.settings` | module wrapper delegating to central handler | module-owned embedded schema | central `settings.read/write` | metadata + central RBAC seed `menu_settings` | central `adminModules` projection | central `site_settings` and `0007` schema | C4 migrate |
| `admin.activity` | module wrapper delegating to central read-only handler | module-owned embedded schema | central `operations.read` | metadata + central RBAC seed `menu_activity` | central `adminModules` projection | core operationlog; `activity` UI is optional | C4 migrate; keep log core |

Evidence for the matrix: `apps/api/internal/kernel/profile.go:91-107`,
`apps/api/internal/handler/schema.go:13-79`,
`apps/api/internal/manifest/manifest.go:27-126`,
`apps/api/internal/store/seed.go:19-113`, and the module/handler/store paths in the
surface table above.

## Core and central ownership disposition

`BuiltinModules()` has six core candidates and four admin candidates. The MVP profile
selects the six core modules plus Users/Roles; the Admin profile adds Settings/Activity;
an explicit `modules.enabled` list replaces profile defaults.

| core/cross-cutting surface | Current owner | R4 disposition |
|---|---|---|
| `core.server-registration` | central HTTP capability and composition root | `keep-core`; provider host boundary |
| `core.auth-session` | auth/session and central Store persistence capability | `keep-core`; provider dependency |
| `core.schema-render` | central `GET /api/schema/{pageId}` dispatcher | `keep-core`; remove module owner special cases in C2/C3 |
| `core.manifest-route` | API manifest route, embedded protocol baseline, aggregator | `keep-core`; remove `adminModules` hard-coded ownership in migration slices |
| `core.navigation-capability` | manifest navigation/expression projection | `keep-core`; consume module fragments |
| `core.operationlog` | append-only operation_log and migrations 0004/0005/0008 | `keep-core`; Activity UI does not own persistence |
| RBAC seed/menu | central `store.seedRBAC`, permissions/menu joins | R4 migration input; provider ownership and seed compatibility must be specified in C2 |
| migration ledger | central ordered `schema_migrations`, versions 1..8 | `keep-core` global order; module ownership mapping required before implementation |
| `0006 records_retire` | central tombstone migration | preserve as historical fact; final R4 treatment pending R4-I003 |

Evidence: `apps/api/internal/kernel/profile.go:25-52,68-104`,
`apps/api/internal/store/seed.go:19-113`, and
`apps/api/internal/store/migrate.go:31-38,60-119,291-313,387-456`.

## Web and delivery surfaces

These are shared consumers and delivery boundaries, not additional business modules:

| surface | owner | disposition | evidence |
|---|---|---|---|
| auth/bootstrap and Shell | `apps/web/src/main.tsx`, `apps/web/src/app/App.tsx` | `keep-core` | `apps/web/src/main.tsx:38-85`; `apps/web/src/app/App.tsx:61-76,313-433` |
| manifest validation, route matching, page loader | Web protocol | `keep-core` | `apps/web/src/protocol/app-manifest.ts:415-620,680-809`; `apps/web/src/protocol/load-page.ts:2-14,66-110` |
| generic Renderer, form controls, SchemaTable and resource actions | Web renderer | `keep-core` | `apps/web/src/renderer/render.tsx:44-64,986-1002`; `apps/web/src/renderer/schema-table.tsx:15-41` |
| protocol conformance fixtures and representative page tests | Web protocol/renderer tests | `keep-core` verification surface | `apps/web/src/protocol/conformance/stage3-fixtures.test.ts:132-398`; `apps/web/src/renderer/representative-pages.test.tsx:41-48,144-196` |
| `apps/web/src/host` | host boundary placeholder | `consumer-only`; no runtime implementation in scope | `apps/web/src/host/README.md:1-6` |
| Docker/nginx build and `/api` proxy | Web delivery runtime | `keep-core` delivery surface | `apps/web/Dockerfile:4-24`; `apps/web/package.json:7-12` |

## C1 disposition

All currently published page, module, RBAC, migration, test, Web and delivery
categories now have an explicit owner and stage disposition. Records is the only
`pending-gate` row and is linked to R4-I003 rather than silently excluded. This
completes the factual response for the inventory finding; it does not close the
separate Records scope conflict, provider contract, or operationlog decisions.
