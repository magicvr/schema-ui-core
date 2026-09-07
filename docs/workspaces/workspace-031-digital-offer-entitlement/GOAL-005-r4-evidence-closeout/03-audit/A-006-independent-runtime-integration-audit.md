---
doc_type: goal-audit
id: A-006-independent-runtime-integration-audit
parent: GOAL-005-r4-evidence-closeout
date: 2026-09-05
status: open
source: independent
auditor: codex (local independent audit)
audit_type: close-out
scope: workspace-031 Root GOAL-001 / GOAL-005 close-out; independent code audit of migration, Manifest, schema loading, module registry and composition/assembly semantics; assess whether the claimed completion is reachable from the configured application and what must precede a successor VP
verdict: fail
open_required: 4
version: 1.0.0
---

# A-006 · workspace-031 运行时集成独立审计

## 结论摘要

- **source**：independent
- **verdict**：**fail**
- **open required**：**4**
- 本意见把 workspace/goal/审计文档仅作为审计范围和声明定位；完成性判断以当前源码、可执行测试和失败路径为依据。
- R2/R3 的若干领域服务行为在当前 `apps/api` Go module 内有真实实现与通过的测试；这不能推出 workspace-031 已从应用配置、模块解析、组合根到 HTTP/Manifest 页面端到端可达。
- 最关键的代码事实是：`biz.digital-offer` 有 Provider 和组合根分支，但没有被加入 `kernel.BuiltinModules()`；`composition.ResolvePlan` 又只用该列表建立 registry。因此当前配置体系没有一条能解析并启用该模块的生产路径。

## 范围与方法

本次只审当前仓库和当前 workspace-031 所对应的代码范围：

- 持久化：`apps/api/modules/digitaloffer/migration/`、`apps/api/modules/compiled/persistence.go`、`apps/api/internal/store/migrate.go`、数字 Offer store/service。
- Manifest/协议：`apps/api/modules/digitaloffer/manifest/`、`apps/api/modules/digitaloffer/schema/`、`apps/api/internal/manifest/`、`apps/api/internal/handler/manifest.go`、`schema.go`、`apps/web/src/protocol/`。
- 装配：`apps/api/kernel/profile.go`、`apps/api/internal/composition/composition.go`、`apps/api/modules/digitaloffer/provider.go`、wallet/Telegram 依赖路径。
- 业务实现与测试：Offer handler、purchase/wallet/check/consume/Telegram 测试。

治理文档只用于确定“workspace-031/R4/Root close-out”这个范围，不作为实现完成、运行时可达或安全语义成立的证据。

## 成果（当前代码中可核对）

1. **领域持久化确实被实现依赖。** `digitaloffer` service/store 读写 `digital_offers`、`digital_purchases`、`digital_entitlements`；`migration.go:21-75` 定义三表及约束，`migration.go:153-163` 注册全局迁移 0070。购买、钱包冻结/扣除、权益 Check/Consume 的核心代码和失败/并发测试存在。
2. **Manifest/schema 组件本身存在。** `provider.go:84-103` 注册两个 PageContribution；`provider.go:139-145` 注册 Manifest fragment；`manifest/fragment.json:12-48` 绑定 pageId、schemaUrl、route 和导航；`schema.go:22-44` 以认证 handler 提供已注册 schema。聚合器和 Web loader 对冲突、未知字段、HTTP/JSON/结构错误采用 fail-closed 语义。
3. **组合根有条件装配意图。** `composition.go:619-646` 仅当 `plan.HasModule("biz.digital-offer")` 时构造共享 runner 上的 Offer store、wallet store、subject store、service 和 Provider；Telegram 也按计划注入真实或 disabled dispatcher。这一条件式装配方向本身是合理的。
4. **局部验证可重放。** 在 `apps/api` 目录执行 `go test ./...` 退出码 0，`go build ./...` 退出码 0；在 `apps/web` 执行 `npm test -- --run src/protocol/app-manifest.test.ts src/protocol/load-page.test.ts`，27 tests 全部通过。以上只能证明现有测试切片和编译，不证明 digital-offer 生产计划已经可解析或挂载。

## Findings

### F-003 · required · high · `biz.digital-offer` 不在编译模块注册表，配置启用路径不可达

**证据**

- `apps/api/modules/digitaloffer/provider.go:40-62` 定义了 `biz.digital-offer` 的 Provider Descriptor（版本 `1.0.0`、路由、页面、导航、权限和 fragment）。
- `apps/api/kernel/profile.go:158-218` 是当前 `BuiltinModules()` 返回的完整模块描述列表；其中有 wallet、channel.telegram 等候选，但没有任何 `ID: "biz.digital-offer"` 条目。
- `apps/api/internal/composition/composition.go:87-100` 的 `ResolvePlan` 仅以 `kernel.BuiltinModules()` 建立 registry，并对 `cfg.ModulesEnabled`/profile 解析结果调用 `registry.Resolve`。因此自定义 `app.modules.list` 也只能引用 registry 中已有模块；当前 digital-offer 会在解析阶段被判定为未编译/未知模块。
- `apps/api/internal/composition/composition.go:619-646` 虽然有 `plan.HasModule("biz.digital-offer")` 分支，但没有 descriptor 进入 plan 时该分支永远不会被生产 `NewApp`/配置路径命中。
- `apps/api/internal/composition/composition.go:121-158` 先解析 plan，再创建 Fx 图；后续 provider 注入不能绕过 plan 解析。

**影响**

当前代码可以编译数字 Offer 包，也可以单独测试 service/handler，但无法通过仓库现有 profile/custom-module 配置启用 `biz.digital-offer`。因此 HTTP routes、页面 schema、navigation、Manifest fragment 和可选 Telegram command 不构成已装配的应用能力。现有 `go test ./...` 通过不能消除这一缺口。

**必改要求**

1. 将 digital-offer 的运行时 descriptor 纳入编译模块注册表，或形成同等可核对的、不会绕过 `ResolvePlan` 的正式注册机制；版本、依赖、capability、ContributionKeys 必须与 `Provider.Descriptor()` 精确匹配。
2. 增加至少一条显式配置/测试路径（建议 custom preset 或专用 dogfood profile），包含其依赖模块，并以 `NewApp`/真实 mux 验证模块可解析、可注册和可访问。
3. 修复后必须重新做独立复审；不能以当前 service-level 绿色测试直接关闭本 finding。

### F-004 · required · high · “Offer CRUD”声明与实际 API/服务面不一致，Delete 缺失

**证据**

- `apps/api/modules/digitaloffer/provider.go:48-55` 的 Provider 路由声明只有 Offer `GET`、`POST`、`PATCH`，没有 `DELETE /api/digitaloffer/offers/{id}` 或等价删除/归档接口。
- `apps/api/internal/handler/digitaloffer.go:34-267` 的实际 Offer/Admin 路由同样只有列表、创建、更新；未发现 Offer 删除路由。
- `apps/api/modules/digitaloffer/service/service.go:156-265` 提供 `CreateOffer`、`UpdateOffer`、`ListOffers`、`ListOnSale`、查询购买和权益等方法，未提供 Offer 删除方法。
- `apps/api/internal/handler/digitaloffer_test.go` 的现有场景覆盖创建、更新/乐观锁、审计失败回滚和查询，但没有 Delete 语义测试。

**影响**

如果 “CRUD” 按通常字面含义包含 Delete，则成功标准/证据矩阵不能成立；如果产品设计有意禁止物理删除，只允许 `off_sale` 生命周期，则当前代码可以支持一种“Create/Read/Update/Status”子集，但必须把合同/成功标准/证据分母明确收窄，而不能继续以 CRUD 表述。

**必改要求**

在 `/govern` 中对该 finding 选择并落盘一种合法路径：

- `fixed`：实现受约束的删除/归档语义，并补 Provider、handler、service、权限、审计、迁移影响和失败/并发测试；或
- 书面收窄为不可删除的生命周期语义：明确 `off_sale`/作废是否替代 Delete、删除约束及其复审触发，并同步合同与证据分母。

本独立审计不替用户作该业务裁决。

### F-005 · required · high · 0070 迁移被全局应用，但当前运行时模块不可启用；迁移必要性与部署语义未闭合

**证据**

- `apps/api/modules/compiled/persistence.go:26-52` 将 `digitaloffermigration.Provider{}` 放入全局 `PersistenceCatalog` 的静态列表。
- `apps/api/internal/composition/composition.go:199-224` 的 `openStore` 在 profile 解析前后均按 compiled-global catalog 打开并应用迁移；这套路径不按 `plan.HasModule("biz.digital-offer")` 过滤持久化。
- `apps/api/modules/digitaloffer/migration/migration.go:21-75,140-163` 的 0070 仅创建三张新表和索引，且只提供正向 Apply；没有数据回填、旧 schema 转换或 Down migration。
- `apps/api/internal/store/migrate.go:81-102,105-131` 在升级前可创建 SQLite snapshot，并将 Apply 与 ledger 放在事务中，但当前路径没有“迁移失败后自动恢复 snapshot”逻辑；已提交迁移没有应用层逆向操作。
- 当前 `biz.digital-offer` 不能进入 plan（F-003），所以这项全局 schema side effect 会发生在没有对应业务 surface 的实例上。
- `SCHEMA_UI_R2_PG_DSN` 当前未设置；代码中 `apps/api/internal/store/postgres_test.go` 与 `apps/api/modules/digitaloffer/service/purchase_test.go:650-656` 对真实 PostgreSQL 场景采用环境门控跳过。本轮 `go test ./...` 的绿色结果不等于真实 PG 迁移/购买矩阵已执行。

**判断：迁移是否必要**

- 对“未来要运行的 digital-offer service”而言，三张表不是多余的：当前 store/service 的 SQL 直接依赖它们，因此持久化 schema 设计具有必要性。
- 对“当前声称已经完成且可部署的 workspace-031”而言，0070 的全局应用尚未证明是必要的，因为对应 runtime provider 还不可由配置启用；它目前是**提前产生的全局 dormant schema side effect**。
- 因此本 finding 不是要求立即删除 0070；它要求先闭合“模块可启用性 + 全局持久化策略”的设计契约。若继续采用本仓“compiled-global persistence”策略，应明确接受未启用模块也会创建 schema，并证明升级、失败、旧二进制/新 schema、SQLite/PG 的运维边界；若不接受，应另行设计并审视 profile-gated persistence，这将是核心方法/迁移 runner 级别的高影响变更，不应在 workspace-031 关门后静默改动。

**必改要求**

1. 先由 `/govern` 记录迁移策略裁决：保留 compiled-global dormant migration，或改为 profile-gated/deferred migration；不得以“以后再处理”代替裁决。
2. 在选定策略下补迁移证据：fresh、reopen/idempotency、Apply 中途失败与 reopen、SQLite/真实 PostgreSQL 双方言、部署/快照恢复处置；如果明确不支持回滚，必须将其作为有范围、有复审触发的残余风险书面接受，而不是把 snapshot 当成自动回滚证据。
3. 在 F-003 闭合前，不得把 0070 作为 workspace-031 已交付运行时能力的完成证据。

### F-006 · required · medium · 缺少从配置到 Manifest/schema/HTTP 的真实组合根证据

**证据**

- `apps/api/modules/digitaloffer/provider.go` 没有 provider-level test 文件；现有 `apps/api/modules/digitaloffer/` 测试主要集中在 service 层。
- `apps/api/server/serve_test.go:72-113` 的默认 server 集成测试覆盖 healthz、readyz、Manifest 和登录，但没有 digital-offer route、page schema、navigation 或 digital-offer Manifest entry 的断言。
- 当前 Web 协议测试（27 tests）验证通用 Manifest/page loader 的错误和缓存语义，但未证明 API 实际聚合出的 digital-offer fragment 能从启用 profile 到达页面。
- F-003 使得这种真实组合路径目前无法成立；因此“Manifest 组件存在”和“装配分支存在”仍只是静态组件证据。

**必改要求**

补一条最小可重复的组合根验收切片：通过显式 `app.modules.list`/preset 启用 `biz.digital-offer` 及依赖，构建实际 `NewApp` 或等价生产 mux，验证：

1. Manifest 含两个 digital-offer page、route、schemaUrl 和 navigation refs；
2. `/api/schema/digitaloffer-offers` 与 `/api/schema/digitaloffer-entitlements` 在正确认证下返回对应文档，未知 pageId 仍 fail-closed；
3. Offer Admin/public routes 的认证、权限和错误语义从真实组合根可达；
4. Telegram disabled 与 enabled 两种计划分别满足冻结的 no-op/注册语义；
5. SQLite 必测，真实 PostgreSQL 若作为 R2/R4 双库分母则必须在可用 PG 环境实际执行，而不是仅保留 skip-capable test。

## 对照成功标准的独立结论

| 范围 | 独立结论 | 依据 |
|------|------|------|
| Offer/购买/钱包/权益领域行为 | **部分成立** | service/handler/store 代码及 `apps/api` 可执行测试存在；不等于组合根可达 |
| Admin 协议页面、Manifest、schema 装载 | **组件成立，端到端不成立** | Provider/fragment/handler/Web loader 存在；F-003/F-006 阻断真实启用证明 |
| subject-only、冻结/扣除/回滚、权益 Check/Consume | **局部成立** | 当前实现与测试覆盖了主要路径；PG 环境门控仍是边界 |
| “Offer CRUD” | **不成立/定义未决** | F-004：Create/Read/Update 存在，Delete 缺失 |
| 不进默认 Profile | **静态上符合意图** | `profile.go:25-114` 默认集合不含该模块；但没有可用 custom/profile 注册路径证明“可选启用” |
| 迁移 0070 | **必要性 conditional，部署语义未闭合** | service 依赖表结构；compiled-global 应用和不可逆/未真实 PG 证据造成 F-005 |
| 关门完成 | **不支持** | 4 个 required findings 未闭合 |

## 独立设计建议（给后继 VP 前）

### AI 推荐：先在 workspace-031 内做“运行时集成补完 + 复审”，不立刻更换业务方案

1. **注册表/配置路径补完**：把 `biz.digital-offer` descriptor 纳入 `BuiltinModules()`（或同等正式 registry），保持 Provider Descriptor 与 Plan 精确匹配；建立专用 custom preset/dogfood 配置，明确该模块仍不进入 mvp/admin 默认 profile。
2. **组合根验收**：从配置解析开始跑真实 NewApp/mux smoke，验证 Manifest、schema、navigation、HTTP、权限和 Telegram disabled/enabled 两种计划；不增加新的 public `apps/api/assembly` 工厂，当前证据不足以证明公开 assembly API 扩展有必要。真正需要修正的是 plan/registry 与 composition integration。
3. **迁移策略裁决**：在保持 compiled-global persistence 的前提下，继续使用 0070 也可以，但必须把“未启用模块仍建表”视为明确产品/运维语义并补齐双库、失败处置和旧版本边界；若该 side effect 不可接受，再单独立项改 profile-gated persistence，不应在 workspace-031 里偷偷改变核心迁移语义。
4. **合同修正或 CRUD 补齐**：在实现 Delete 与“不可删除、仅 off_sale”之间作用户裁决并落盘；同步 Provider/权限/审计/测试或同步成功标准/证据分母。
5. **重新独立复审**：以上 required findings 按 fixed 或用户书面 accepted-residual/user-overruled 合法处理后，至少再做一次针对 migration + Manifest + composition 的 independent finding-closure/close-out audit；未完成该复审前，不应把 Root/GOAL-005 作为可复用的“已完成运行时能力”前置。

### 备选方案 B：保持当前 dormant migration，推迟运行时启用

- 保留 0070 和 Manifest/provider 代码，仅把 workspace-031 标记为未完成或 accepted-residual；补一份明确的运维残余接受（范围：所有 profile 会产生数字 Offer 空表；触发：首次启用、schema 变更、生产部署或多实例）并等待后继集成工作。
- 取舍：短期不改核心迁移 runner；但当前 workspace 不能 honest 地保持 done，且后继 VP 不能把它当作已验证的业务运行时基础。

### 备选方案 C：把持久化迁移改为 profile-gated/deferred

- 取消/延迟 compiled-global 应用 0070，改造迁移收集、版本台账、部署升级和跨 profile 兼容策略。
- 取舍：可消除 dormant schema side effect，但会改变全局迁移历史和运维模型，影响面明显超出 workspace-031，应先改核心方法/架构并进行独立 cross audit；本审计不推荐作为当前修复捷径。

## 在后继 VP 立项/激活前的最低工作

若后继 VP 会依赖 digital-offer 的实际运行时能力，最低前置不是再写一份治理说明，而是：

1. `/govern` 响应本 A-006，逐条对 F-003～F-006 走 `fixed`、`accepted-residual` 或 `user-overruled` 合法路径；不得以重写 done/progress 或保留旧 pass 代替。
2. 完成可启用 descriptor + 显式非默认 profile/preset + NewApp/mux 真实 smoke；证明 Manifest、schema、HTTP、权限和 Telegram 接缝。
3. 完成迁移策略裁决和对应 SQLite/真实 PostgreSQL/failure-path/部署处置证据；不要把“有 snapshot”表述成“自动回滚”。
4. 对 Offer 删除语义作用户裁决，并同步实现或收窄成功标准/证据分母。
5. 在上述工作后追加一次 independent re-audit（建议 scope：finding-closure + close-out；至少覆盖 F-003～F-006），确认 open required = 0，才可把 workspace-031 作为后继 VP 的已完成前置；若用户选择 residual，则后继 VP 必须显式继承残余范围和复审触发，不能隐式继承“done”。

## 与既有独立意见的异同

- 既有 A-002/A-004 主要审计关门台账投影、命令边界和既有业务证据；A-004 明确只复审 A-002 的两个治理/命令边界 finding，并不重新审计业务实现、合同或组合根。
- 本 A-006 不改写或否定 A-002 的历史文字，也不把 A-004 的限定 `pass` 扩大解释为运行时集成通过。
- 本意见新增的关键范围是：当前源码中的模块注册/可启用性、compiled-global migration side effect，以及从配置到 Manifest/schema/HTTP 的组合根可达性；因此对“Root 已完成”的总体结论更严格。

## 声明

本意见仅追加到审计 ledger，`source: independent`；不修改任何目标的 `status`、`progress`、goal-tree、决策/执行台账、合同正文或业务代码。后续 finding 响应、用户裁决、状态推进和关门由 `/govern` 处理。
