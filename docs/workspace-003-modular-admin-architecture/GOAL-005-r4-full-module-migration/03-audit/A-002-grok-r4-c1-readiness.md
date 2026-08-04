---
id: A-002-grok-r4-c1-readiness
doc: audit-entry
goal: GOAL-005-r4-full-module-migration
source: independent
auditor: Grok Build / grok-4.5
date: 2026-08-05
scope: R4 C1 capability inventory, provider contract, operationlog boundary, and Records scope conflict
verdict: conditional
---

# A-002 · Grok R4 C1 独立审计

## 结论

Grok Build（`grok-4.5`，reasoning `high`）独立复核了 R4 C1 的 canonical 目标
记录和当前代码，结论为 `conditional`。C1 资料诚实地保持 collecting，未错误
推进状态；但全量一方 Admin 迁移的冻结条件尚未满足。该意见不改变目标状态、
progress 或 goal-tree。

## 已核实成果

- Users/Roles 仍由中心 `handler.Register` 挂载，尚未宣称迁移完成。
- Settings/Activity 的 module 包入口仍委托中心 handler，尚非完整 module-owned
  provider。
- `0006 records_retire` 已删除 Records 表、权限和菜单，当前没有 Records 产品
  handler、`/api/records` 或产品 Schema fixture；历史 `records.*` operation-log
  事件仍保留。
- operationlog 当前是业务成功后的 append-only best-effort 写入，失败不回滚业务
  写入；已有 append、排序、事件 CHECK、迁移保留测试，但没有 retention/archival
  合同或失败注入证据。

## Required findings

### F-GROK-R4-001 · 一方能力盘点尚不足以冻结 R4-I001

- level: `required`
- status: `open`
- impact: C1 freeze and C3/C4 completeness claims
- finding: 现有盘点覆盖主要产品模块和 core 面，但尚未逐项处置所有已发布
  first-party 页面/Schema、中心 RBAC seed/menu ownership、完整 BuiltinModules 集合、
  六项能力矩阵以及 Web 中除 navigation projection 外的相关 surface。
- evidence: `attachments/r4-c1-capability-inventory.md`；
  `apps/api/internal/manifest/app-manifest.json:10-40,67-99`；
  `apps/api/internal/handler/schema.go:17-18,63-76`；
  `apps/api/internal/store/seed.go:39-78`；
  `apps/api/internal/kernel/profile.go:24-47,91-103`；
  `00-meta.md:37-39`。
- closure: 补充每个页面、模块、seed、迁移、测试 surface 的 owner、stage、include/
  exclude，并为每个标准 Admin 模块建立 HTTP/Schema/Authorization/Navigation/
  Manifest/Persistence 六项矩阵。

### F-GROK-R4-002 · Module Plan 与真实注册之间缺少框架无关 provider 契约

- level: `required`
- status: `open`
- impact: C2 contract freeze and all module migration
- finding: Kernel 只有 capability IDs、字符串 contribution keys 和生命周期 hooks，
  没有结构化 HTTP/Schema/Authorization/Navigation/Manifest/Persistence provider；
  Composition、Schema owner map 和 Manifest `adminModules` 仍由中心硬编码消费。
- evidence: `apps/api/internal/kernel/module.go:16-56,322-391`；
  `apps/api/internal/composition/composition.go:89-97`；
  `apps/api/internal/handler/schema.go:63-76`；
  `apps/api/internal/manifest/manifest.go:58-63`；
  `docs/architecture/module-architecture.md:41-55`。
- closure: 记录框架无关 provider 形状、Plan 消费顺序、冲突规则、Fx/DI 边界和
  fail-closed 生命周期/冲突测试，然后再进入 C2。

### F-GROK-R4-003 · operationlog 一致性和 retention 尚未决策

- level: `required`
- status: `open`
- impact: C3 behavior matrix and C5 data gate
- finding: 当前契约是业务成功后的 best-effort append，日志失败不逆转业务成功；
  未找到 retention duration、归档、purge、restore contract 或日志失败注入测试，
  R4 尚未决定保留或强化该语义。
- evidence: `apps/api/internal/store/operations.go:1-5,44-62`；
  `apps/api/internal/handler/users.go:273-303`；
  `apps/api/internal/handler/roles.go:228-258`；
  `apps/api/internal/handler/settings.go:143-157`；
  `apps/api/internal/handler/auth.go:120-151`；
  `apps/api/internal/handler/resources.go:465-530`；
  `apps/api/internal/store/operations_test.go:11-87`。
- closure: 明确保留当前 best-effort 或改为原子一致性，明确 retention 边界，并以
  对应的失败/读取/迁移恢复测试证明选择。

### F-GROK-R4-004 · VP-003 Records/Schema CRUD 与 `records_retire` 冲突

- level: `required`
- status: `open`
- impact: C1 scope freeze, C4, and R4 acceptance
- finding: VP-003 仍把 `records/Schema CRUD` 和记录/Schema CRUD 作为 R4/exit 范围，
  但 `0006 records_retire` 已删除产品表、权限和菜单，当前没有 Records CRUD surface。
  这是信息冲突，不是代码可以自动解决的架构决定。
- evidence: `docs/vision/plans/VP-003-modular-admin-architecture.md:55-56,70`；
  `apps/api/internal/store/migrate.go:291-313`；
  `apps/api/internal/kernel/profile.go:91-103`；
  `apps/api/README.md:109-110`；
  `00-meta.md` 的 R4-I003 与 `D-001`。
- closure: 由用户书面选择恢复 Records 产品 CRUD，或将 R4 明确收敛为仍存在的
  Schema-driven Admin 能力并保留历史事件；随后一致更新当前 canonical 范围记录。

## 推荐但不阻断

- `F-GROK-R4-005`: R4-I004 冻结时将 auth/settings writer 也纳入行为矩阵，因为它们
  共享 best-effort operationlog 语义。
- `F-GROK-R4-006`: provider 冻结后把 `schema.go` owner map 和 `manifest.go`
  `adminModules` 明确登记为待移除的中心特例/tombstone。

## 放行结论

四项 required findings 均为 `open`。R4-I001～I004 仍为 `collecting`，不得关闭
C1、进入 C2、推进 Root progress 或将本目标标为 `done`。该独立意见与 self
A-001 同向，没有需要 P-004 冲突仲裁的相反 verdict。
