---
doc_type: architecture-playbook
title: 一方模块贡献 Playbook
status: active
created: 2026-08-06
updated: 2026-08-06
parent: null
version: 1.0.0
vision_ref: schema-ui-core-admin-foundation@0.2.0
serves: VP-004-module-contribution-readiness
authority_of: product-module-contribution
architecture_boundary: module-architecture.md
---

# 一方模块贡献 Playbook

本文是 **产品模块贡献方法论 / 操作 playbook** 的权威入口：一方**标准 Admin 功能模块**从零接入时必须完成什么、明确不必/禁止做什么，以及演进应归入薄内核 / 组合根 / 独立模块的判定方法。

| 权威边界 | 路径 |
|----------|------|
| 架构终态（契约与边界） | [module-architecture.md](module-architecture.md) |
| 本 playbook（操作化清单） | 本文件 |
| 发现入口 | [overview.md](overview.md) · 根 [QUICKSTART.md](../../QUICKSTART.md) §5 |
| 过程台账 | `docs/workspaces/workspace-004-module-contribution-readiness/`（VP-004） |

**不是**本文件职责：Goal Governance 元规则（[principles.md](principles.md) P-001～P-006）、workspace-protocol、默认代码脚手架、`AGENTS.md`/Skills 改写。

参考正例（一方标准 Admin 功能模块）：`admin.users` → `apps/api/internal/modules/users/`（`provider.go`、`schema/`、`manifest/`）。  
横切基础设施对照：`core.operationlog` → `apps/api/internal/modules/operationlog/`。

---

## 1. 新增一方标准 Admin 功能模块 · 必须完成（MUST）

下列清单与 [module-architecture.md](module-architecture.md) §2.1 核心六项及仓库真实路径对齐。完成前不得宣称「模块已接入」。

| # | 必须完成项 | 落点 / 证据 |
|---|------------|-------------|
| M1 | **模块 id / 版本 / 内核 API 范围 / 依赖** | 实现 `kernel.Provider`：`Descriptor()` 返回稳定不可复用 `ID`、`Version`、`KernelAPIRange`、显式 `DependsOn`（可为空列表但必须声明）。契约类型：`apps/api/internal/kernel/module.go`。 |
| M2 | **核心六项贡献** | 标准 Admin 功能模块必须实现 HTTP、Schema、Authorization、Navigation、Manifest、Persistence（语义见 architecture §2.1）。通过 `Register(ctx, Registrar)` 注册；不得以「按需」永久缺省六项。 |
| M3 | **组合根静态候选注册** | 在 `apps/api/internal/composition/composition.go` 将 Provider 纳入已编译候选集（静态 import + 按 `plan.HasModule(...)` 装配）。不得依赖运行时下载/插件加载。 |
| M4 | **Profile / `modules.enabled` 成员关系** | 新模块若需进入默认 Profile，更新 `apps/api/internal/kernel/profile.go` 中 `profileDefaults`（`mvp` / `admin` 等）；custom 须显式 `modules.enabled`。解析优先级见 architecture §3 与 config。 |
| M5 | **全局迁移台账参与** | 若模块拥有 schema 迁移：经 `apps/api/internal/modules/compiled/persistence.go` 的 `PersistenceProviders()` 参与 **全局** 不可变 checksum 台账；**不以是否启用以过滤迁移**（architecture §4.1）。无迁移时 `CompiledPersistence()` 可返回空，但仍须声明 Persistence 能力语义。 |
| M6 | **验证 / 回归最小集** | 至少：模块级 Provider/契约测试；依赖/冲突 fail-closed；相关 Profile 启动路径；页面/权限/导航可观察。仓库范例：`apps/api/internal/modules/users/provider_test.go`、`apps/api/internal/kernel/provider_test.go`、`apps/api/internal/composition/composition_test.go`。 |

### 1.1 推荐目录骨架（与现网一致）

```text
apps/api/internal/modules/<short-name>/
  provider.go          # kernel.Provider：Descriptor + Register + CompiledPersistence
  provider_test.go
  schema/              # 页面 Schema 文档与贡献
  manifest/            # Manifest fragment
  # 可选：migration/、systemdata/、repository/
```

组合根与内核契约：

```text
apps/api/internal/composition/composition.go   # Fx 组合根、静态 import Provider
apps/api/internal/kernel/                      # Module / Provider / Registrar / Profile
apps/api/internal/modules/compiled/            # 全局迁移收集（全候选）
```

### 1.2 最小 Register 覆盖（对照 `admin.users`）

正例 `apps/api/internal/modules/users/provider.go` 展示了标准 Admin 模块在 `Register` 中应覆盖的六面：

1. **HTTP** — `reg.HTTP(...)` 路由（如 `/api/users`）
2. **Schema** — `reg.Schema(PageContribution{... Document: ...})`
3. **Authorization** — `reg.Authorization(PermissionContribution{...})`
4. **Navigation** — `reg.Navigation(NavigationContribution{...})`
5. **Manifest** — `reg.Manifest(FragmentContribution{...})`
6. **Persistence** — 迁移归全局台账；`admin.users` 将账户/RBAC 迁移归属 `core.auth-session`（`CompiledPersistence` 空返回 + 依赖声明），业务表迁移则应进入 `compiled.PersistenceProviders`

---

## 2. 明确不必做 / 禁止做（DO NOT）

与 architecture §1 / §5 / §6 及 VP-004 exit #2 对齐。

| # | 禁止 | 原因 / 正确做法 |
|---|------|-----------------|
| D1 | **不要为接模块改 Renderer / Shell 中央业务注册** | 增减标准模块不得要求改前端中央路由/导航/页面注册；由模块贡献 Schema/Navigation/Manifest，后端聚合后经 `/.well-known/schema-ui/app-manifest.json` 发布。 |
| D2 | **不要在生产路径静默使用静态 Manifest 兜底** | 禁止把生产 Manifest 放进 `apps/web/public/` 静默兜底；开发期兼容须有期限与删除门禁（architecture §5）。 |
| D3 | **不要私建平行认证 / 授权 / DB** | 认证、授权、数据库、错误协议与日志属内核能力；模块只消费稳定接口（architecture §6）。 |
| D4 | **不要把「按需能力」当成核心六项可永久缺省** | §2.2 按需（Configuration / Lifecycle / Observability）**不得**覆盖 §2.1 核心六项。横切模块若豁免 UI 相关项须有显式架构说明。 |
| D5 | **不要做运行时插件 / 热插拔 / `.so` / 远程下载幻想** | 仅在已编译候选集中由 Profile/`modules.enabled` 选择启用；不支持运行中启停（architecture §3）。 |

### 2.1 额外高风险反模式（推荐遵守）

| 反模式 | 说明 |
|--------|------|
| 改 `handler` 中央业务 Register 塞新资源 | 应走 Provider 贡献；组合根只装配 |
| 启用集过滤迁移 / 删禁用模块表 | 违反全局迁移与数据生命周期 |
| 依赖静默自动启用 | 依赖闭包 fail closed；不静默补全 |
| 在 `docs/schemas/` 放业务页面文档 | 该目录是上游协议 JSON Schema，不是业务页面库（见 QUICKSTART §5） |

---

## 3. Core vs 模块责任 · 归属判定

与 [module-architecture.md](module-architecture.md) §1 / §6 一致。按下列顺序判定；**不得**用本表推翻 architecture 边界。

### 3.1 判定树

```text
1. 是否为全站稳定基础契约（配置/日志/DB/HTTP 生命周期/认证授权接口/模块协议）？
   → 是：薄内核（apps/api/internal/kernel、config、store、auth 等）
2. 是否仅为「静态 import、解析 Profile、装配 Provider、生命周期编排」？
   → 是：组合根（apps/api/internal/composition）
3. 是否为可装配的一方业务/产品能力（有稳定 module id，可被 Profile 启停暴露）？
   → 是：独立（或既有）模块 apps/api/internal/modules/<name>
4. 是否仅为某模块内部实现细节（私有 repo、helper）？
   → 是：模块内 util；不提升为内核 API
5. 是否横切且始终启用、无标准管理 UI 要求？
   → 可能是横切基础设施模块（如 operationlog）；可豁免部分 UI 贡献，须显式说明
```

### 3.2 正反例

| 场景 | 归属 | 依据 |
|------|------|------|
| 新增「用户」CRUD API + 页面 + 权限 + 菜单 | **模块** `admin.users` | 标准 Admin 功能；正例已在 `modules/users` |
| 新增角色管理 | **模块** `admin.roles` | 同左；`modules/roles` |
| 写操作审计记录（始终开） | **横切模块** `operationlog` | architecture §6；非标准 Admin 六项全套 UI 要求 |
| Activity 只读查询 UI | **模块** `admin.activity`（可选启用） | 可关 UI 不得关 operationlog |
| JWT/会话校验接口 | **内核 / auth-session** | 模块不得私建平行认证 |
| 把 users 路由写进 `handler` 中央列表而不经 Provider | **反例** | 违反「模块不改中央业务注册」 |
| 为新模块改 Web Renderer 主路径 | **反例** | D1；应 Schema 驱动 |
| 新库表迁移只在启用模块时执行 | **反例** | 全局迁移台账；compiled 全候选收集 |
| 订单/钱包/类目业务 | **未来业务模块（另立 VP）** | Charter 非目标；不在本 playbook 交付范围 |

### 3.3 横切 vs 标准 Admin 功能模块

| 类型 | 示例 id | 核心六项 | Profile |
|------|---------|----------|---------|
| 标准 Admin 功能 | `admin.users`, `admin.roles`, `admin.settings` | **必须**满配 | 常进 `mvp`/`admin` 默认集 |
| 横切基础设施 | `core.operationlog` | 可豁免 Schema/Navigation/Manifest 中不适用项 | 通常始终启用 |
| 内核能力 | （非 module id 或 core.* 契约） | 不适用「接模块」清单 | 不由业务 Profile 装卸语义替代 |

---

## 4. 可发现性与 AI 充分条件

默认充分路径（不要求改 `AGENTS.md` / Skills）：

1. [overview.md](overview.md) 仓库布局 / 模块扩展指针 → 本文件  
2. 根 [QUICKSTART.md](../../QUICKSTART.md) §5「接业务」→ 本文件  
3. [module-architecture.md](module-architecture.md) 链出「操作 playbook」

读者**无需**阅读 closed VP-003 工作区过程树即可到达本文。

---

## 5. 最小验证清单（接模块后）

- [ ] `go test` 覆盖新 Provider（及组合根若改 Profile 默认集）
- [ ] 目标 Profile 下进程启动；`/readyz` 就绪
- [ ] Manifest 经 `/.well-known/schema-ui/app-manifest.json` 可见新页面/导航（非静态 public 兜底）
- [ ] 权限键与导航可见性符合 Authorization/Navigation 贡献
- [ ] 若有迁移：全局台账连续、checksum 稳定；禁用模块不删表

---

## 6. 修订

| 版本 | 日期 | 说明 |
|------|------|------|
| 1.0.0 | 2026-08-06 | VP-004 / workspace-004 Root 首版：MUST / DO NOT / 归属法；路径对齐现网 modules + composition + kernel |
