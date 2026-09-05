---
doc_type: goal-audit
id: A-012-independent-runtime-closure-reaudit
parent: GOAL-005-r4-evidence-closeout
date: 2026-09-05
status: open
source: independent
auditor: codex (local independent audit)
audit_type: re-audit
scope: workspace-031 Root/GOAL-005 runtime closure after A-006 remediation; code-first re-audit of migration, Manifest, module registry and composition semantics, including prerequisites for a successor VP
verdict: conditional
open_required: 2
version: 1.0.0
---

# A-012 · workspace-031 运行时集成独立复审

## 结论摘要

- **source**：independent
- **verdict**：**conditional**
- **open required**：**2**（均为 medium）
- 本次把治理文档仅用于定位被声称已关闭的范围；完成性判断来自当前源码、真实组合根、迁移 runner、失败路径与本轮实际执行的测试。
- A-006 的核心实现缺陷已经被实质修正：`biz.digital-offer` 已进入编译模块注册表，可由 custom `app.modules` 解析；真实 `newAppWithOptions` 已能挂载其 Manifest/schema/Admin/public HTTP surface；真实 PostgreSQL 购买与完整迁移目录也在本轮实际执行通过。
- **当前方案方向正确，不建议更换为 profile-gated migration、另建 public `assembly` API 或实现 Offer Delete。** 注册表接入、Provider/Manifest contribution 和 compiled-global 0070 都是现有架构下的必要实现；更换这些方案反而会扩大核心语义变更面。
- 但现有关闭声明仍把“计划可解析”当成了 Telegram-enabled 的真实装配证据，并把迁移 runner 的事务实现当成了“0070 Apply 中途失败后可安全 reopen”的可执行证据。两项均尚未被对应失败/组合切片直接证明，因此不支持无条件复用为后继 VP 的已验证前置。

## 独立设计基线

本次不以既有治理裁决为设计前提，按当前代码约束独立推导出以下最小方案：

1. **模块启用**：`kernel.BuiltinModules()` 是 `ResolvePlan` 的编译候选真相源；`biz.digital-offer` 必须注册于此，同时继续不进入 mvp/admin 默认 profile，仅由 custom module list 显式启用。
2. **Manifest/schema**：继续由 `digitaloffer.Provider.Register` 贡献 route、schema、permission、navigation 和 fragment，再由现有 Registrar/聚合器装配；不应为本模块增加旁路 Manifest 或新的公开 assembly 工厂。
3. **持久化**：在本仓既有 compiled-global catalog 下，0070 必须随编译目录全局应用。改为 profile-gated/deferred migration 会制造按 profile 分叉的 schema/ledger，并改变“后启用模块”的启动语义，除非另立核心迁移目标并做 cross audit，否则不是本次整改的较优方案。
4. **Offer 生命周期**：购买凭证与权益引用要求 Offer 可下架但不可物理删除；Create/Read/Update/Status 是与现行数据模型一致的接口，不应为了字面 CRUD 增加 Delete。

## 已确认关闭的 A-006 范围

### F-003 · 模块注册与配置可达性：fixed

- `apps/api/kernel/profile.go:219-227` 已加入 `biz.digital-offer` descriptor；`apps/api/internal/composition/composition.go:87-100` 仍以 `BuiltinModules()` 构造 registry，因而 custom module list 现在可以真实解析该模块。
- registry descriptor 与 `apps/api/modules/digitaloffer/provider.go:40-61` 的 Provider descriptor 在路由、页面、导航、权限、fragment、依赖和 capability 上一致；真实组合根启动还会经过 Provider descriptor 一致性校验。
- `TestDigitalOfferPlanResolution` 对含/不含 `channel.telegram` 的两种 plan 解析均通过；本轮独立重放通过。

### F-004 · Offer Delete / 合同分母：closed（contract-conformant）

- 当前实现没有 DELETE route；这与数据模型的不可删除生命周期一致，而不是实现遗漏。Offer 状态机通过 `draft/on_sale/off_sale` 表达可售性，购买与权益保留对 Offer 的历史引用。
- 当前证据矩阵已经将字面 “Offer CRUD” 收窄为 Create/Read/Update/Status 且明确无删除。无需增加 Delete；建议用负向路由测试固化该边界，但它不是本轮 required。

### F-005 · 迁移策略：方案成立，但验证闭环不完整

- `apps/api/modules/compiled/persistence.go:26-52` 将 digital-offer migration 纳入全局目录；`apps/api/modules/digitaloffer/migration/migration.go:20-138` 为 SQLite/PostgreSQL 定义同一三表约束模型；这是当前 store/service 可运行所必需。
- `apps/api/internal/store/migrate.go:81-131` 确实把单个 migration 的 Apply 与 ledger insert 放在同一事务中，Apply 返回错误时 rollback。fresh、reopen/idempotency、SQLite 完整目录和真实 PostgreSQL 完整目录/购买验收均已通过。
- 缺口转为下述 F-007：目前仍没有直接覆盖多语句 Apply 中途失败、无部分 schema/ledger 残留、随后 corrected catalog reopen 成功的可执行切片。

### F-006 · 组合根：主体成立，但 enabled-channel 分支未被真实装配证明

- `apps/api/internal/composition/composition.go:619-646` 已按 plan 条件构造 digital-offer service/provider，并在 Telegram 启用时选择 live dispatcher/sender，否则使用 disabled 实现。
- `TestDigitalOfferCompositionRoot` 经真实 `newAppWithOptions`/Fx/mux 验证了 Telegram-disabled 启动、Manifest 页面、schema 认证/未知页 fail-closed、Admin route 与 public catalog；本轮独立重放通过。
- 缺口转为下述 F-008：该测试明确排除 `channel.telegram`；enabled plan 当前仅验证 `ResolvePlan`，尚未执行真实 Fx/runtime 装配。

## Findings

### F-007 · required · medium · 缺少 0070 多语句 Apply 失败后无残留并可 reopen 的直接证据

**证据**

- 0070 的每个方言都包含多条 `CREATE TABLE`/`CREATE INDEX` 语句（`apps/api/modules/digitaloffer/migration/migration.go:20-138`）。
- runner 的事务边界在代码上成立（`apps/api/internal/store/migrate.go:105-131`），但当前 migration 测试覆盖的是 fresh、成功应用、reopen/no-op、目录/校验和与若干身份拒绝路径；本次检索未发现“执行首条 DDL 后主动失败 → 断言 schema/ledger 均无残留 → corrected catalog reopen 成功”的迁移测试。
- 本轮真实 PostgreSQL `TestFullCatalogPostgresBootstrapIntegration`、`TestPostgresMigrateRunnerIntegration` 与 `TestPurchasePostgresAcceptance` 均通过，证明成功路径真实可用，但不能替代 Apply 中途失败的恢复证明。

**影响**

- 这不是要求增加 Down migration，也不否定 compiled-global；它只说明关门证据把实现推断当成了失败路径事实。
- 若未来修改 dialect adapter、DDL 顺序或 transaction wrapper，缺少该切片会使“失败后安全重启”成为未受保护的部署假设。

**必改要求**

1. 在 `internal/store` 增加迁移 runner 级失败/reopen 测试：使用多语句 migration，首条创建对象后返回注入错误；断言本次 migration 的 schema 对象和 ledger 行均不存在。
2. 使用修正后的同版本/同目录重新打开，断言迁移成功且仅有一条 ledger 记录。
3. SQLite 必测；既然本仓当前真实 PG 环境可用，PostgreSQL 同语义也应进入集成测试。无需给 0070 增加 Down，也无需改 profile-gated。

### F-008 · required · medium · Telegram-enabled 与完整 Manifest contribution 仍缺真实组合根验收

**证据**

- `TestDigitalOfferPlanResolution` 只证明含 `channel.telegram` 的模块图可解析（`apps/api/internal/composition/composition_digitaloffer_test.go:35-64`）。
- 唯一真实 app/mux 测试明确在配置中排除 Telegram，并说明只证明 disabled dispatcher（同文件 `66-82`）；因此没有执行 `composition.go:629-636` 的 live dispatcher/sender 选择，也没有证明 business commands 注册到 webhook 所使用的同一进程 runtime。
- 当前 Manifest 断言只搜索两个 page id（同文件 `108-118`）；route、schemaUrl、navigation refs 以及权限贡献分别存在于静态 Provider/fragment 中，但未在这条真实聚合结果上做结构化断言。

**影响**

- 代码结构强烈表明 enabled 分支应当工作，但“应当”不足以支持后继 VP 将 Telegram 命令集成和 Admin 页面装配视为稳定前置。
- `plan.HasModule("channel.telegram")` 成立而 runtime 状态不完整时，当前组合代码会退回 disabled dispatcher；只有真实 enabled 装配测试才能防止该分支被未来改动静默降级。

**必改要求**

1. 增加 `biz.digital-offer + channel.telegram` 的真实 `newAppWithOptions`/Fx 组合测试，使用现有 fake/runtime seam，禁止真实外网 Bot API 依赖。
2. 证明 digital-offer 的 `price`/`buy`/`entitlements` handlers 注册在 webhook/Telegram surface 使用的同一 dispatcher；至少验证一个命令经该组合根可达。
3. 将 Manifest 响应解析为结构，断言两个 page 的 route、schemaUrl、navigation refs，并保留当前 schema/HTTP fail-closed 断言。disabled 分支现有测试可保留，不需重写装配架构。

## 非 required 建议

1. 为 Offer Admin surface 增加 DELETE 请求的 404/405 负向测试，防止未来误把“CRUD”重新解释为物理删除。
2. 为非全权 Admin 增加 `digitaloffer.offer.manage` 与 `digitaloffer.entitlement.void` 的权限分离测试；当前 admin 200/anonymous 401 不能单独证明两把权限键的边界。
3. 可增加 registry descriptor 与 Provider descriptor 的专门等价测试；当前真实 app 启动已经能捕获漂移，因此优先级低于 F-007/F-008。

## 本轮实际验证

在 `apps/api/`：

- `go test -count=1 ./internal/composition -run 'TestDigitalOffer(PlanResolution|CompositionRoot)$' -v`：PASS。
- `go test -count=1 ./modules/digitaloffer/service -run '^TestPurchasePostgresAcceptance$' -v`：PASS，真实 PostgreSQL 子场景均实际执行，未 skip。
- `go test -count=1 ./internal/store -run '^(TestFullCatalogPostgresBootstrapIntegration|TestPostgresMigrateRunnerIntegration)$' -v`：PASS，真实 PostgreSQL，未 skip。
- `go test -count=1 ./...`：PASS。
- `go build ./...`：PASS。

在 `apps/web/`：

- `npm test -- --run src/protocol/app-manifest.test.ts src/protocol/load-page.test.ts`：PASS，2 files / 27 tests。

## 对后继 VP 的门禁建议

### AI 推荐

先在 workspace-031 内完成 F-007/F-008，并做一次 **focused independent finding-closure 复审**；无需再做一次全量 close-out，也无需更换迁移/Manifest/装配方案。复审只需核对：

1. SQLite/PG migration Apply 中途失败、零残留、corrected reopen；
2. Telegram-enabled 的真实组合根与同一 runtime 命令注册；
3. 结构化 Manifest page/route/schemaUrl/navigation 断言。

复审通过后，workspace-031 才适合作为依赖这些能力的后继 VP 的已验证前置。如果后继 VP 完全不依赖生产迁移、Telegram 命令或该 Admin surface，是否以书面 `accepted-residual` 限定范围后先立项，应由 `/govern` 提交用户裁决；本独立审计不替用户接受残余。

### 何时需要更换方案或升级审计模式

- 只有在决定放弃 dormant schema、改成 profile-gated/deferred migration，或新增 public `assembly` API 时，才需要单独立项并做 **cross audit**；这些会改变平台级迁移/装配语义。
- 若仅按 F-007/F-008 补测试与验证，不改变生产语义，一次 focused independent 复审足够。

## 声明

本独立意见只追加到 GOAL-005 审计台账，不修改目标 `status` / `progress`、`goal-tree.md`、VP 状态、决策正文或业务代码。A-011 的关门记录是否需要响应、重开或接受限定残余，由 `/govern` 汇总本意见后处理。
