---
id: A-003
goal: GOAL-009-r3-s03-system-monitoring
source: independent
date: 2026-08-14
scope: S5 关门 · 安全/数据门禁（admin.system-monitoring vs D-002 冻结方案）
verdict: conditional
auditor: grok-build
audit_type: close-out
status: recorded
parent: GOAL-009-r3-s03-system-monitoring
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# A-003 · independent 安全/数据审计（S-03 实现）

## 范围与区间

- **auditor**：grok-build（independent cross-audit · grok-4.6）
- **type**：close-out / security-data gate
- **workspace**：`workspace-011-admin-functional-modules`（`root_goal` = `GOAL-001-admin-functional-modules`；`canonical_scope` 已核对；`shared_materials_catalog: none`）
- **covered**：
  - `apps/api/internal/handler/systemmonitoring.go`（GET status / GET errors 只读资源）
  - `apps/api/internal/modules/systemmonitoring/`（provider、schema、manifest fragment）
  - `apps/api/internal/kernel/profile.go`（admin 默认集）
  - 工厂门禁与列表契约：`apps/api/internal/handler/resources.go`、`operations.go`
  - Host 既有 statCard 数据契约：`apps/web/src/renderer/resource.ts`、`render.tsx`
  - 计划契约 `01-decision/D-002-s1-plan-freeze.md`（及 D-001/D-003 Profile 声明）
  - 页面文档对照 `docs/schemas/page.schema.json`（本轮 AJV draft-07 校验通过）
- **excluded**：端到端浏览器手测、生产部署、其它工作区上下文、非本模块的资源工厂/Host 实现变更
- **信息项**：I-001 / I-002 均已 closed（D-002 §1 / D-001 §2）；本 scope 无到期未关闭 required 信息项；无共享资料引用

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| status / errors 均 `monitoring.read` fail-closed | status 经 `a.Middleware` + `requirePermission(..., "monitoring.read")`（systemmonitoring.go:95–98）；errors 经 `ResourceRoutes` 的 list/detail，同样先 `requirePermission`（resources.go:207–210、286–288、514–518）。匿名无 Bearer → Middleware 401 `UNAUTHENTICATED`（auth.go:386–395），不进入探测或列表体；已认证无 key → 403 `FORBIDDEN`（resources.go:237–239）。测试：`TestSystemMonitoringStatus` 401/403；`TestSystemMonitoringErrors` viewer 403；`TestMonitoringProviderServesStatus` 匿名 401 |
| 401/403 体不泄漏 status/errors 载荷 | 失败路径只走 `writeLocalizedError`；`writeJSON(MonitoringStatus, …)` 仅在权限通过之后（systemmonitoring.go:96–125）。匿名响应为目录化 `UNAUTHENTICATED`，无 version/commit/modules/dbSize |
| DB 路径不回传，仅大小 | `MonitoringStatus` 无 path 字段（systemmonitoring.go:23–31）；`os.Stat(dbPath)` 只取 `info.Size()`，失败则 0（113–116）。`dbPath` 来自 composition `cfg.DBPath`（composition.go:234），非请求输入 |
| 只读：无写路由、无审计事件、无迁移 | `ReadOnly: true` 不挂 POST/PATCH/DELETE/batch-delete（resources.go:211–218）；entity 写方法一律 `errReadOnlyResource`（systemmonitoring.go:64–74）；无 `OnWrite`；`CompiledPersistence` 返回 `nil, nil`（provider.go:60–62）；本模块无 `RecordOperation` / 无 migration 贡献 |
| errors 列表契约镜像 operation-log | 同一 `ListOperationsFiltered` + `operationToMap`；`SortFields` = `createdAt,event,actorName`，`QSearch: true`（systemmonitoring.go:39–50、78–88 vs operations.go:13–25、53–65）。工厂白名单 sort/order、ReadOnly 默认 `desc`、`q` 仅在 QSearch 时传入（resources.go:291–332）。仓库层 `?` 绑定 + sort 开关（operationlog/repository.go:228–249） |
| readiness 门 nil-safe | `if ready != nil && !ready()`（systemmonitoring.go:109–112）；nil 视为就绪。测试环境显式传 `nil`（testhelpers_test.go:108）且 `TestSystemMonitoringStatus` 期望 `ready=true` |
| 权限/导航 admin-only | `PolicyID` / `Visibility` = `PolicyAdmin`（provider.go:81–99）；`rolesForPolicy` 只映射 `admin`（policy.go:15–16）。mvp/demo 默认集未加入该模块（profile.go:26–93） |
| Profile 仅为声明的内容扩展 | `profileDefaults[ProfileAdmin]` 追加 `"admin.system-monitoring"` 并注释 D-001（profile.go:71–73）；`ResolveProfile` / `ParseModuleList` 逻辑未改；`BuiltinModules` 仅登记描述符（162–163）。组合根 `plan.HasModule` 才挂 provider（composition.go:233–234）。admin 权限 18 / 导航 10、mvp 不变（composition_test.go:449–457） |
| 页面 AJV 结构合法；无新 renderer 类型 | 本轮 `page.schema.json` AJV 校验 `ok: true`。节点仅 `section` / `text` / `grid` / `statCard` / `table`（system-monitoring.json）；未改 `CUSTOM_HANDLER_URLS`、未新增 node type |
| I-001/I-002 不影响本门禁 | 均 closed 于 S1；共享资料目录 `none` |

## 对照成功标准（本 scope）

| 标准 | 结论 |
|------|------|
| 两端点 401/403 fail-closed；匿名不泄漏载荷 | **满足** |
| status 管理面门禁；不回传 DB 路径 | **满足** |
| 无写路径 / 无审计伪造 / 无迁移 | **满足** |
| errors 镜像 operation-log 列表契约（items/total、sort 白名单、q）；ready nil-safe | **满足**（列表路径）；detail 404 映射见 F-002 |
| 页面对照 `page.schema.json`；无新 renderer 扩展 | **结构满足**；**Host 可执行性不满足**（F-001） |
| Profile 仅 admin 默认集内容扩展 | **满足** |

## Findings

### F-001 · status 扁平对象无法被既有 statCard 列表信封消费，监控卡在 Host 上 fail-closed

| 字段 | 值 |
|------|-----|
| level | required |
| status | open |
| evidence | D-002 §2 冻结 GET `/api/system-monitoring/status` 为 `{status,ready,version,commit,uptimeSeconds,modules[],dbSizeBytes}`（systemmonitoring.go:23–31、117–125）。页面五张 statCard 的 `dataSource` 均指向该端点 + `valueField`（system-monitoring.json:30–84）。既有 Host：`StatCardView` → `useDisplayData` → `fetchResourceList`（render.tsx:1693–1747）；`parseResourceList` 在缺少 `items` 数组或非有限 `total/page/pageSize` 时抛错（resource.ts:149–163、167–186）。对照可用先例：dashboard statCard 使用 `/api/users` 列表信封 + `valueField: "total"`（dashboard.json:29–47） |
| severity | med（非授权绕过；D-002 §3 监控卡在当前 Host 上不可执行。已认证 200 响应也会被解析拒绝，五张卡走 `statCard data failed to load`） |

**说明**：AJV `page.schema.json` 只约束节点形状，不约束 dataSource 响应信封。本模块未新增 renderer（符合 D-002 §3「无新 renderer 扩展」字面），但把**非列表**汇总端点接到**只认 `{items,total,page,pageSize}`** 的既有 statCard 上。A-002 self「页面 schema 协议校验通过」不足以覆盖 Host 数据契约。errors 表走同一列表信封，该半页可工作。

次要：D-002 中文写「模块数」，schema `valueField` 为 `modules`（数组）。即便信封修好，`format: plain` 会 `String(array)` 打出模块 id 列表而非计数。

**建议修复**（不新增 renderer）：让 status 在保留 D-002 顶层字段的同时带上单行列表信封，例如：

```json
{
  "status": "ok",
  "ready": true,
  "version": "...",
  "commit": "...",
  "uptimeSeconds": 1,
  "modules": ["..."],
  "dbSizeBytes": 0,
  "items": [{ "status": "ok", "ready": true, "uptimeSeconds": 1, "modules": ["..."], "dbSizeBytes": 0, "moduleCount": 1 }],
  "total": 1,
  "page": 1,
  "pageSize": 1
}
```

`parseResourceList` 忽略额外顶层字段，现有 handler 字段断言仍可通过。若卡片要「模块数」，增加数值字段并改 `valueField`。补一条 Host 级测试：真实 status JSON 经 `fetchResourceList` 后五张卡能读到对应 `valueField`。修订 D-002 若顶层信封被视作契约变更。

### F-002 · errors detail 用 `==` 比较包装后的 `ErrNotFound`，缺 id 走 500 而非 404

| 字段 | 值 |
|------|-----|
| level | recommended |
| status | open |
| evidence | `monitoringEntity.Get`：`if err == operationlog.ErrNotFound`（systemmonitoring.go:53–61）。`GetOperation` 经 `withTx` 以 `fmt.Errorf("%s: %w", …)` 包装（operationlog/repository.go:178–201）；`Store.WithTx` 不剥包装（store.go:68–79）。activity 对照使用 `errors.Is`（operations.go:67–74）。结果：`GET /api/system-monitoring/errors/{missing}` 在已授权下走 `writeEntityError` 的 INTERNAL 分支（resources.go:255–269）→ 500，而不是 `OPERATION_NOT_FOUND` 404 |
| severity | low（授权仍先于 Get；不泄漏行内容；不是 IDOR。列表路径不受影响） |

**建议修复**：与 `operations.go` 对齐，改为 `errors.Is(err, operationlog.ErrNotFound)`；补一条缺 id → 404 `OPERATION_NOT_FOUND` 的测试。

## 必改项汇总

- **required / 必改**：F-001
- **recommended**：F-002（不阻断授权 fail-closed / 只读 / Profile 项）

## 与既有意见的异同

| 条目 | 关系 |
|------|------|
| A-002 self（S2–S4，verdict pass） | **同意**只读、401/403、nil ready、errors 复用 operationlog、无迁移、Profile 内容扩展。**不同意**「页面 schema / 无新 renderer」已足以放行 D-002 §3——A-002 以 AJV 结构校验代替 Host statCard 列表信封。另：A-002 写「ready 门翻转」有测试，`TestSystemMonitoringStatus` 仅断言默认 `ready=true`，未见翻转用例（代码路径仍存在，不单列 finding） |
| A-001 self（S1） | 方案级只读/权限/错误日志诚实化与 API 实现一致；S1 未核对 statCard 与冻结 status 形状是否被现有 Host 消费 |

## 结论 + 建议给编排器/用户的下一步

**verdict: conditional** — 安全核心（授权 fail-closed、匿名无载荷、不回传 DB 路径、只读无写/无审计伪造/无迁移、errors 列表契约、ready nil-safe、Profile 仅 admin 内容扩展）有代码与测试证据。不可无条件放行本门禁：D-002 §3 的 status statCard 在当前 Host 上不可执行（F-001）。

建议 `/govern`：响应本意见；先修 F-001（列表信封或改绑定，并补 Host 断言）；F-002 可同波或维护。闭合 required 前不要把本目标标为 `done`。

### 声明

本意见 `source: independent`，**不修改**目标 `status` / `progress` / goal-tree / 方案正文或实现代码；响应与状态变更由 `/govern` 与用户裁决处理。
