---
id: GOAL-001-modular-admin-architecture
doc: audit-entry
record_id: A-010
source: independent
scope: VP-003 终态意图 vs apps/api·web 代码内聚（R4 已关 / R5 进行中；重点 api store/handler/persistence）
verdict: conditional
status: recorded
parent: null
created: 2026-08-05
updated: 2026-08-05
version: 0.1.0
auditor: Grok Build / grok-4.5
audit_type: ad-hoc
---

# A-010 · VP-003 意图 vs apps 代码内聚独立审计（2026-08-05）

- **source**：independent
- **auditor**：Grok Build / grok-4.5（`/audit` 立场；用户请求将前序代码审视结论落盘）
- **类型**：ad-hoc · execution-facts（相对 VP-003 退出判据与 `module-architecture.md` 的代码对齐）
- **scope**：`apps/api`（主）与 `apps/web`（辅）；对照 [VP-003](../../../../vision/plans/VP-003-modular-admin-architecture.md) 方向级退出判据 #2/#3/#4/#6 与 [module-architecture.md](../../../../architecture/module-architecture.md) §1–§4；阶段语境为 R1–R4 已关、`GOAL-012` R5 active
- **verdict**：conditional
- **工作区上下文**：`workspace-003-modular-admin-architecture` · Root `GOAL-001-modular-admin-architecture` · `shared_materials_catalog: none`

## 范围与区间

### 覆盖

| 项 | 路径 / 说明 |
|----|-------------|
| 意图权威 | VP-003 退出判据 #2 薄内核/六项、#3 数据生命周期、#4 Manifest 聚合、#6 旧路径退出 |
| 架构权威 | `docs/architecture/module-architecture.md` §1–§4、§6–§8 |
| 冻结历史 | R4 C1 冻结包 §2.3 / §4（`CompiledPersistence`、store `core.persistence` 仅为历史） |
| API | `apps/api/internal/{store,handler,kernel,composition,modules/*}` |
| Web | `apps/web/src/{app,protocol,renderer,host}` 抽查 |
| 阶段台账 | Root progress `4/6`；GOAL-012 R5 residual 清单（E-002） |

### 排除

- 不复审 R4 关门合法性本身（GOAL-011 A-003 已 conditional 放行并登记 residual）。
- 不审业务领域模块（订单/钱包等，VP Non-goals）。
- 不改 status / progress / goal-tree / 方案正文。
- 不写入 `docs/vision/reviews.md`（非 Vision Review）。
- 不读取其他工作区过程状态。

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| 薄内核契约与框架无关 Provider/Registrar 存在 | `apps/api/internal/kernel/{module,provider,contribution,persistence}.go` |
| Fx 仅组合根 | `apps/api/internal/composition/composition.go`；modules 无 `go.uber.org/fx` |
| 四方 Admin 表面 Provider 化，composition finalize 挂载路由 | `modules/{users,roles,settings,activity}/provider.go`；`composition.newMux` |
| Schema **内容**已迁入模块 embed | `modules/*/schema/`；handler 仍静态合并（见 F-003） |
| Manifest 片段模块拥有并聚合 | `modules/*/manifest`；`manifest.ForModulesWithFragments` |
| Profile `mvp`/`admin` 默认集 | `kernel/profile.go` `profileDefaults` |
| 前端主路径 Manifest + 通用 Renderer | `apps/web/src/app/App.tsx`；`protocol/`；`renderer/` |
| R4 residual 已部分诚实登记 | GOAL-011 E-003；GOAL-012 E-002 / 00-meta C5.1 |

## 对照 VP-003 退出判据（代码快照）

| 退出 # | 主题 | 状态 | 说明 |
|--------|------|------|------|
| 2 | 薄内核、组合根、六项契约、无中央业务注册表 | **条件未齐** | 表面六项中 Persistence 实现为空；领域实现仍在中心 `store`/`handler`（F-001/F-002/F-005） |
| 3 | 全局迁移 catalog、与 enablement 解耦、bootstrap/reconcile | **条件未齐** | 生产路径仍 `store.compiledMigrations`；`CollectPersistence` 未接线（F-002） |
| 4 | 后端聚合 Manifest/Schema/导航为生产路径 | **部分** | Manifest 聚合成立；Schema 发布非 ContributionSet 驱动（F-003） |
| 5 | 横切 operationlog / auth 边界 | **部分** | operationlog Option A 行为可保留；writer/seed 所有权仍中心（F-001/F-005） |
| 6 | 旧装配路径退出 | **未完成** | 中心适配器、死 `module.Register`、test 双轨仍在（F-004/F-007）；属 R6 主战场但须可追踪 |
| Web | 同一 build、无 Shell 硬编码业务注册 | **大体对齐** | host 空层；无新的阻断级偏差（F-009 recommended） |

**总评**：相对 R4「表面 Provider 迁移」可接受；相对 VP-003 **终态**仍有系统级中心化持久化与中心 handler 实现。**不得**将 R4 done / progress `4/6` 解读为退出判据 #2/#3/#6 已取证。

## Findings

### F-001 · `internal/store` 仍为跨模块上帝对象（领域 SQL/迁移/seed 未内聚）

- **严重度**：high
- **建议**：required
- **状态**：open
- **影响门禁**：VP-003 退出 #2 / #3 / #6；Root 关门；R6 旧路径删除前须可证明所有权迁移或书面 residual
- **建议承接阶段**：R5 至少完成**所有权模型与 residual/信息项登记**；实现迁移可 R5 切片 + R6 收尾，但不得静默遗漏
- **描述**：`apps/api/internal/store` 同时拥有：DB 打开/ping、全局 `compiledMigrations`（0001–0008）、admin/RBAC seed、users/roles CRUD、site_settings、operation_log 读写。标准 Admin 模块 Provider 仅注入 `*store.Store` 并委托 handler，**无模块内仓储**。这与「Persistence 为模块核心贡献、内核只提供平台契约」终态冲突。合理拆分应为：内核/平台保留 runner+ledger；`core.auth-session` / `core.operationlog` 持横切表；`admin.*` 持领域仓储与后续迁移。
- **证据**：
  - `apps/api/internal/store/{store,migrate,seed,users,roles,settings,operations}.go`
  - 各模块 `CompiledPersistence() (nil, nil)` 注释「existing ledger is store-owned」：`modules/users/provider.go` 等
  - 冻结包：`GOAL-005/.../r4-c1-freeze-package-draft.md` §4.1「Store … core.persistence 只能作为当前历史事实，不能作为 R4 终态契约」
- **建议修复**：
  1. 决策记录：平台 `DB/Tx/runner` vs 模块 repository 边界；
  2. 将领域方法迁入 `modules/*/…` 或模块私有 store 包；
  3. 缩减中心 `store` 为平台 + 历史 ledger 宿主（过渡）直至迁移 descriptor 模块化；
  4. 在 Root/R5 信息门禁或 residual 清单显式登记本项（见 F-008）。

### F-002 · `CompiledPersistence` / `CollectPersistence` 未进入生产迁移路径

- **严重度**：high
- **建议**：required
- **状态**：open
- **影响门禁**：VP-003 退出 #3；数据生命周期与 module-owned 迁移追加规则
- **建议承接阶段**：R5 C5.2 相关深化 **或** 独立 R5/R6 工作包；在接线完成前禁止宣称「compiled-global Persistence 已生产化」
- **描述**：内核已实现 `Provider.CompiledPersistence` 与 `kernel.CollectPersistence`（含单测），但 composition/`store.Open` 仍执行包内 `compiledMigrations`，`ModuleID` 硬编码 `core.persistence`。四方模块返回空 migration 列表。契约与生产路径分裂。
- **证据**：
  - `kernel/persistence.go` `CollectPersistence`
  - `store/migrate.go` `compiledMigrations` + `ModuleID: "core.persistence"`
  - 模块 provider `CompiledPersistence` 均 `return nil, nil`
  - composition `openStore` → `store.Open`（无 CollectPersistence 调用）
- **建议修复**：历史 0001–0008 以 tombstone/descriptor 挂正确 module id 进入 catalog（可不重编号/改 checksum）；生产 Open 仅消费 Collect 结果；新迁移只允许 Provider 追加。

### F-003 · Schema 发布仍非 ContributionSet 驱动（中心静态合并）

> **状态更新（A-014/F-014-002）**：本条闭合口径已拆分——**F-003a** 门禁/owner 贡献
> 驱动 `fixed`（`RegisterSchemas` 接受 `set.Pages`，提交 `d1c372e`）；**F-003b**
> document 字节 ContributionSet 发布 `accepted-residual`→R6 C6.3（见 A-011/A-013）。

- **严重度**：med
- **建议**：required
- **状态**：open（历史快照；权威闭合见 A-011 拆分与 A-013）
- **影响门禁**：VP-003 退出 #4；R5 C5.1 residual
- **描述**：页面 JSON 虽在模块包内，但 `handler/schema.go` 通过静态 import 各 `modules/*/schema` 合并 map，再按 plan 过滤。不是 `RegisterContributions` 产生的 Page 贡献驱动发布。禁用/冲突语义依赖中心装配而非 finalize 后的 ContributionSet。
- **证据**：`apps/api/internal/handler/schema.go` `staticSchemaDocuments` / `schemaDocumentsForPlan`；GOAL-011 E-002 residual 行；GOAL-012 E-002
- **建议修复**：由 PageContribution（或等价已校验集合）发布 schema 字节；去掉 handler 对业务模块 schema 包的编译期枚举。

### F-004 · 中心 handler 仍承载业务协议实现与遗留 Register 适配器

- **严重度**：med
- **建议**：recommended（R6 删除主路径；R5 须保持可追踪）
- **状态**：open
- **描述**：通用 resource factory、Settings/Activity 路由适配、`RegisterSettings`/`RegisterActivity`、`modules/*/module.go` 旧 `Register` 仍在。生产 composition 已走 Provider，但「实现内聚」未完成；测试路径仍可走中心挂载。
- **证据**：`handler/{resources,settings,operations,health}.go`；`modules/{settings,activity}/module.go`；GOAL-012 E-002「中心适配器终态删除 pending」
- **建议修复**：协议工厂可上移为共享/kernel-http；业务适配进模块；删除死适配器与 test-only 双轨说明后的 `MountProviderRoutes` 中心业务挂载。

### F-005 · Seed / RBAC reconcile 未以 Authorization contribution 为唯一源

- **严重度**：med
- **建议**：required
- **状态**：open
- **影响门禁**：VP-003 退出 #2/#5；冻结包 §5 owner matrix
- **描述**：冻结包要求 permission key 在 seed、Manifest、handler、Schema action 间一致，且中心 reconcile 只消费已验证 module contributions。现状 `store.seedRBAC` / seedAdmin 仍为中心硬编码路径，不读 Provider Authorization 贡献。
- **证据**：`store/seed.go`；`store/store.go` Open 内 seed；冻结包 §5；modules 仅声明 Permissions 键未驱动 seed
- **建议修复**：reconcile 输入改为已 finalize 的 Permission/Navigation 贡献（或显式 system-data contribution）；保留 fresh bootstrap 最小 admin，与 versioned reconcile 分离（架构 §4.2）。

### F-006 · 组合根 if-chain 构造 Provider + `BuiltinModules` 与 Descriptor 双源

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **描述**：`composition.newMux` 手写 `plan.HasModule("admin.users")` 等分支；`kernel.BuiltinModules()` 与各 `Provider.Descriptor()` 字段平行维护，易漂移。
- **证据**：`composition/composition.go`；`kernel/profile.go` `BuiltinModules`；`modules/*/provider.go` Descriptor
- **建议修复**：编译期 provider catalog（id → factory）；Plan 元数据以 Descriptor 为单一来源或生成。

### F-007 · 旧装配双轨（module.Register / 中心测试挂载）仍未删除

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **影响门禁**：VP-003 退出 #6 / R6
- **描述**：与 F-004 相关但单独标记删除门禁：`RegisterSettings`/`RegisterActivity`、module 级 `Register(mux,…)`、`MountProviderRoutes` 注释为 test-only 仍留在 handler。R6 必须可列出删除清单与回归。
- **证据**：同上；GOAL-011 E-002 residual「中心 RegisterSettings/RegisterActivity 终态删除」

### F-008 · R5 residual / 信息门禁未覆盖 store·Persistence 内聚债

- **严重度**：high
- **建议**：required
- **状态**：open
- **影响门禁**：R5 C5.1（R5-I001）；Root 阶段门禁可见性；防止「R5 关完仍无所有权模型」
- **描述**：GOAL-012 已登记 Schema 贡献驱动、中心适配器、PolicyID/Visibility、双 Profile 矩阵、Configuration 等 residual，**未**将 F-001/F-002/F-005 级「store 上帝对象 + 生产迁移未走 CollectPersistence + seed 非贡献驱动」写入 residual 或 required 信息项。前序 R4 冻结已声明 store 非终态，但治理台账未形成可阻断的后续门禁。
- **证据**：GOAL-012 `00-meta.md` C5.1 / R5-I001；`02-execution/E-002-r5-readyz-and-residuals.md` residual 列表；对照本意见 F-001/F-002/F-005
- **建议修复**：`/govern` 响应时至少其一：  
  (a) 扩展 R5-I001 / C5.1 residual 清单含 Persistence/Repository 所有权与 CollectPersistence 接线；或  
  (b) 新建 Root/R5 子目标与 required 信息项，明确最晚阶段（建议不晚于 R6 方案冻结前完成模型，R6 完成删除/迁出）；  
  并写明在未 closed/residual 接受前，**不得**将 VP-003 退出 #2/#3 标为已取证。

### F-009 · Web 侧无阻断级偏离（记录）

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **描述**：Shell/Renderer/Manifest 消费主路径与 VP 一致；`host/` 仍为空预留；图标 registry 与 branding 事件属宿主能力，不构成「中央业务页面注册表」复发。无需为本条扩大 R5 范围。
- **证据**：`apps/web/src/app/App.tsx`；`host/README.md`；protocol/renderer 布局

## 必改项汇总（开放 required）

| ID | 摘要 | 建议阻断 |
|----|------|----------|
| F-001 | store 上帝对象 / 领域未内聚 | Root done / VP 退出 #2/#3/#6 取证前 |
| F-002 | CollectPersistence 未生产接线 | VP 退出 #3 取证前 |
| F-003 | Schema 非 ContributionSet 驱动 | R5 C5.1 residual 闭合或 VP 退出 #4 取证前 |
| F-005 | seed/RBAC 非贡献驱动 | VP 退出 #2/#5 取证前 |
| F-008 | residual/信息门禁未登记上述债 | **R5 C5.1 / R5-I001 响应前优先**（治理可见性） |

推荐项 F-004 / F-006 / F-007 / F-009 不单独阻断 R5 开工，但 F-004/F-007 应进入 R6 删除清单。

## 与既有意见的异同

| 既有 | 关系 |
|------|------|
| GOAL-011 A-003（R4 C5 independent） | 已记 Schema/适配器/readyz 等 **recommended** residual；本意见确认其方向，并将 **store/Persistence/seed** 与 **台账缺口** 升为 Root 级 **required**（R4 当时允许表面迁移，不免除终态） |
| GOAL-012 A-001 / E-002 | R5 residual 列表不完整（F-008）；readyz 已推进的部分不在本 scope 重审 |
| 冻结包 §4 | 与 F-001/F-002 一致：历史 store 非终态契约 |

## 结论与给编排器的下一步

**verdict: conditional**

- 代码相对 VP-003 **终态未对齐**的主因是 **持久化与领域仓储仍中心化**，叠加 Schema 发布与 seed 非贡献驱动；这不是「再写几个 Provider 方法」可掩盖的表面问题。
- R4 关门与 residual 诚实性 **不**被本意见推翻；但 **Root / VP 退出取证**在 F-001/F-002/F-003/F-005/F-008 开放期间 **不得**宣称完成。
- **建议 `/govern` 优先响应 F-008**（把债写入 R5 信息门禁或子目标），再排 F-001/F-002 模型与接线，F-003/F-005 并入同一所有权迁移包；F-004/F-007 链到 R6 删除清单。

### 声明

本意见 `source: independent`，**只写审计 ledger**（本文件 + `03-audit.md` 索引）。  
**未**修改任何目标 `status` / 检查点 / 派生 `progress` / goal-tree 状态列 / 方案或代码。  
响应、finding 闭合（fixed / accepted-residual / user-overruled）与推进由 **`/govern`** 执行。
