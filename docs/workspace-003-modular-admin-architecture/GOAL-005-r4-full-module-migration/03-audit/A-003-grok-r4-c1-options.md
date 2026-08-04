---
id: A-003-grok-r4-c1-options
doc: audit-entry
goal: GOAL-005-r4-full-module-migration
source: independent
auditor: Grok Build / grok-4.5
date: 2026-08-05
scope: R4 C1 provider and operationlog decision options, framework boundary, migration, security, lifecycle, compatibility, and test gates
verdict: conditional
---

# A-003 · Grok R4 C1 方案材料独立审计

## 结论

Grok Build（`grok-4.5`，reasoning `high`）独立复核了 Provider/Registrar 与
operationlog 方案材料、模块架构和当前 Go 实现，结论为 `conditional`。方案是
诚实且方向一致的待裁决材料，但还不是可冻结的 D-003；本意见不改变目标状态、
progress 或任何用户决策。

## 已核实成果

- 方案明确标记为 `status: draft`、`decision_state: pending_user`，没有把
  Provider、Records 或 operationlog 选项伪装成已接受决定。
- Provider/Registrar 的概念接口不引入 Fx 类型，Plan -> 临时注册集 -> 全局校验
  -> finalize -> lifecycle 的顺序与架构的先注册、后校验、失败闭锁原则一致。
- 六类贡献与现有中心特例有对应关系：HTTP、Schema、Authorization、Navigation、
  Manifest、Persistence；operationlog 当前业务成功后 best-effort append 的事实也
  与实现和既有意见一致。
- Option A/B/C 覆盖了保持当前语义、同事务原子化和完整归档生命周期三条决策路径，
  但仍需用户选择，且 Option A 的 retention residual 尚未达到正式接受条件。

## Required findings

### F-IND-R4-OPT-001 · Provider 冻结材料仍缺少可实现的合约细节

- level: `required`
- status: `open`
- impact: R4-I002、C2 contract freeze
- finding: 六类 `*Contribution` 只有名称，没有字段、稳定 identity、冲突键，且未
  说明与现有 `ContributionKeys` 是并行双检、替换还是仅预声明。模块如何在不引入
  Fx 的情况下取得 Store/auth 依赖未定义；`Configuration` 及
  `ConfigNamespaces` 的处理也未冻结。
- evidence: `attachments/r4-c1-provider-operationlog-options.md:25-63`；
  `apps/api/internal/kernel/module.go:28-56,349-390`；
  `apps/api/internal/composition/composition.go:89-97`；
  `apps/api/internal/modules/settings/module.go:14-16`；
  `docs/architecture/module-architecture.md:41-67`。
- closure: D-003 或冻结附件需给出贡献字段和稳定键、Plan 与 Registrar 的关系、由
  composition 构造依赖且模块包不依赖 Fx 的规则，以及 Configuration 是 R4 实施项
  还是明确留空。

### F-IND-R4-OPT-002 · Persistence 收集路径可能违反全局迁移账本规则

- level: `required`
- status: `open`
- impact: R4 migration、R5 readiness inheritance
- finding: 架构要求所有已编译的一方迁移无论模块是否启用都进入唯一全局、不可变、
  有 checksum 的账本；当前方案按启用 Plan 描述 Provider 注册，未冻结 disabled-but-
  compiled migration 的收集、tombstone、reconcile 和 seed 规则。
- evidence: `docs/architecture/module-architecture.md:79-92`；
  `attachments/r4-c1-provider-operationlog-options.md:40-70`；
  `apps/api/internal/store/migrate.go:60-119`；
  `apps/api/internal/composition/composition_test.go:233-255`。
- closure: 明确迁移描述来自 compiled candidate/global catalog，启用过滤只影响
  HTTP/Schema/Nav/Manifest 等 surface，并补充 disabled module 仍应用迁移、数据保留
  和不覆盖用户数据的测试。

### F-IND-R4-OPT-003 · Authorization、seed、Manifest 与敏感信息边界未冻结

- level: `required`
- status: `open`
- impact: C2 security、C3 Users/Roles compatibility
- finding: 未冻结 Authorization contribution 如何归属 permission/menu seed 和
  system-data reconcile，也未把 permission key 一致性、公开 Manifest 不含秘密、
  operationlog detail 不含 token/secret 作为 provider/Option A 的验收门禁。
- evidence: `attachments/r4-c1-provider-operationlog-options.md:60-68`；
  `apps/api/internal/store/seed.go:19-114`；
  `docs/architecture/module-architecture.md:52,97-107`；
  `apps/api/internal/store/operations.go:31-33`。
- closure: 冻结 Authorization -> seed/reconcile 的 owner 和 fail-closed 一致性检查，
  并将 core.auth-session、RBAC、Manifest 与 operationlog 敏感字段约束加入 C2/C3 门禁。

### F-IND-R4-OPT-004 · Option A retention residual 尚未达到合法接受条件

- level: `required`
- status: `open`
- impact: R4-I004、C3 behavior matrix、C5 data gate
- finding: 当前 best-effort 失败语义已经足以作为 Option A 的候选事实，但 retention
  只写了 R4 不自动 purge/archive，未明确运行期保留事实、后续 duration 的 residual
  范围、责任人和复核触发条件，也未明确 Records 历史事件和 Activity UI 关闭时仍
  必须写 operationlog。
- evidence: `attachments/r4-c1-provider-operationlog-options.md:85-118`；
  `apps/api/internal/store/operations.go:1-5,44-62`；
  `docs/architecture/module-architecture.md:103-106,114`；
  `GOAL-005-r4-full-module-migration/00-meta.md:60`。
- closure: 若用户选择 A，必须书面冻结“R4 不自动删除 operation_log”，另附含范围、
  责任人和复核触发条件的 residual，并要求 Users/Roles/Auth/Settings 失败注入及
  Activity disabled 仍写入的测试；不得把“不删除”表述成永久 retention policy。

### F-IND-R4-OPT-005 · 生命周期与发布失败闭锁门禁不完整

- level: `required`
- status: `open`
- impact: C2 lifecycle verification
- finding: 方案提到临时 ContributionSet 和失败清理，但未冻结 Provider 注册失败、
  贡献冲突、Start 失败的分类与清理，HTTP/Manifest 在 listen 前后的发布边界，
  disabled surface 缺失且 core.operationlog 保持提供的验证，也未声明 `/readyz` 仍只
  是 store ping、不能冒充模块图 readiness。
- evidence: `attachments/r4-c1-provider-operationlog-options.md:48-52`；
  `apps/api/internal/composition/composition.go:61-75,89-160`；
  `apps/api/internal/handler/health.go:54-66`；
  `docs/architecture/module-architecture.md:109-114`。
- closure: C2 计划需列出 conflict 无部分发布、register 失败不 listen、Start 失败
  清理已启动模块、Activity disabled 仍记录 operationlog，以及双 Profile 验证。

### F-IND-R4-OPT-006 · 中心特例移除与兼容性切换顺序未冻结

- level: `required`
- status: `open`
- impact: C3/C4 cutover、non-regression
- finding: 方案知道 settings/activity wrapper 不是终态，但没有冻结 metadata-only、
  dual-run、移除 Schema owner map/Manifest `adminModules`/中心 Register 分支的顺序，
  也没有列出既有 HTTP 成功语义、operationlog event CHECK 和新迁移版本之间的兼容约束。
- evidence: `attachments/r4-c1-provider-operationlog-options.md:71-105`；
  `apps/api/internal/manifest/manifest.go:58-63`；
  `apps/api/internal/handler/schema.go:63-76`；
  `apps/api/internal/store/migrate.go:222-234,249-262,348-360`。
- closure: 冻结切换顺序、中心特例 tombstone、事件/HTTP 兼容列表，并明确永久双路径
  不属于终态；若选 B，必须另行接受其 breaking behavioral change。

## 推荐但不阻断

- `F-IND-R4-OPT-007`: 增加静态检查，禁止 `modules/**` 与 `kernel` import Fx。
- `F-IND-R4-OPT-008`: 明确 Registrar 与现有 `Module.Hooks` 的 Observability/Lifecycle
  归属。
- `F-IND-R4-OPT-009`: 若选择 Option B，要求所有 writer 共用同一 `sql.Tx` 并覆盖
  lock/retry 矩阵。
- `F-IND-R4-OPT-010`: Records historical-only 仍只是工程建议，不得绕过 R4-I003 的
  P-004 用户裁决。

## 放行结论

六项 required findings 均为 `open`。该意见与 A-001/A-002 同向，没有相反 verdict
需要 P-004 冲突仲裁；但 R4-I002/R4-I004 仍不能关闭，C1 不能冻结或进入 C2，
也不能推进 Root progress。
