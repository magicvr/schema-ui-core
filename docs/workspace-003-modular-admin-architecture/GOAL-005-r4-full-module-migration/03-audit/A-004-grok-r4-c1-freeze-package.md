---
id: A-004-grok-r4-c1-freeze-package
doc: audit-entry
goal: GOAL-005-r4-full-module-migration
source: independent
auditor: Grok Build / grok-4.5
date: 2026-08-05
scope: Candidate R4 C1 freeze package response to A-003 F-IND-R4-OPT-001 through 006
verdict: conditional
---

# A-004 · Grok R4 C1 冻结包候选独立审计

## 结论

Grok Build（`grok-4.5`，reasoning `high`）复核了当时工作树中的
`r4-c1-freeze-package-draft.md`、A-003、架构文档和当前 Go 实现，结论为
`conditional`。草案显著改善了 A-003 的响应，保持 Provider、Records 和
operationlog 取舍为 `pending_user`，但审计时仍有 4 项 required residual；该意见
不改变目标状态、progress、D-003 或任何 finding 状态。

## 已核实成果

- 草案明确 `status: draft`、`decision_state: pending_user`，没有把推荐写成用户决定。
- Provider/Registrar 公共形状不引入 Fx，composition-owned construction 和 compiled-global
  migration 原则与架构方向一致。
- 草案列出六类 contribution 的最小字段、Plan/Registrar 双检、失败闭锁、seed/security
  和 operationlog residual 表单。
- Records historical-only、operationlog A/B/C 和 Provider 采纳均仍是待裁决轴。

## Required findings

### F-IND-R4-FP-001 · Persistence collection path 在审计时仍有矛盾

- level: `required`
- status: `open`
- impact: R4-I002、C1 freeze、架构 §4.1
- finding: 当时草案一方面把 Persistence 放在 Registrar，另一方面规定只有 enabled
  Provider 调用 Register，同时又要求全部 compiled provider 的 migration 进入全局
  catalog。若不定义独立 compiled collector，C2 可能重新引入 Plan-gated migration。
- evidence: `attachments/r4-c1-freeze-package-draft.md` 的 Provider/Registrar、注册
  顺序和 Persistence 章节；`docs/architecture/module-architecture.md:79-84`；
  `apps/api/internal/store/migrate.go:60-119,459-472`。
- closure: 只能保留一条明确路径，例如由每个 compiled Provider 独立返回
  `CompiledPersistence`，且 Registrar 不拥有 Persistence 写入口；新增 disabled profile
  migration ledger 测试。

### F-IND-R4-FP-002 · Contribution contracts 在审计时仍未达到规范类型冻结

- level: `required`
- status: `open`
- impact: R4-I002、C2 public contract
- finding: 当时草案的六类字段仍是 prose 表格，没有规范字段名/类型；C2 实现者仍可能
  发明影响后续兼容性的 Route、Schema、Policy、Visibility、Manifest 和 Migration 结构。
- evidence: `attachments/r4-c1-freeze-package-draft.md` 的 contribution table；
  `apps/api/internal/kernel/module.go:28-56`。
- closure: 冻结 package-level struct 字段/类型，或明确可接受的 C2 invention bound，
  并为 identity、冲突键、授权策略和迁移 callback 增加 contract tests。

### F-IND-R4-FP-003 · Option A residual 仍只是模板

- level: `required`
- status: `open`
- impact: R4-I004、C3/C5 data gate
- finding: residual 的 owner、review date 仍为 `pending_user`。这是正确的非静默姿态，
  但不能据此关闭 Option A 的 required 信息或宣称 retention 已接受。
- evidence: `attachments/r4-c1-freeze-package-draft.md` 的 operationlog residual table；
  `apps/api/internal/store/operations.go:1-5,44-62`。
- closure: 用户书面选择 A 后填写范围、责任人和复核触发/日期，或选择 B/C 并完成其
  扩大后的行为和数据生命周期证据。

### F-IND-R4-FP-004 · P-004 三项决策尚未形成 D-003

- level: `required`
- status: `open`
- impact: R4-I002/I003/I004、C1 close
- finding: Provider 精确形状、Records 分叉和 operationlog 选项都仍是候选材料；
  R4-I002/I003/I004 不能由草案自动变为 verified。
- evidence: freeze package frontmatter and §1/§8；`GOAL-005/00-meta.md` 的信息门禁；
  `GOAL-005/03-audit.md` 当前结论。
- closure: 取得用户书面裁决，写入 D-003，逐项关联 evidence 和 finding response，
  再进行最终 self + independent freeze review。

## 推荐但不阻断

- `F-IND-R4-FP-005`: 将 `mvp`/`admin` 双 Profile 明列为 lifecycle/test gate。
- `F-IND-R4-FP-006`: 将 auth-session、RBAC、operationlog 的 cross-cutting owner
  matrix 落盘。
- `F-IND-R4-FP-007`: 枚举 HTTP status/payload、Schema IDs、operationlog event CHECK、
  migration names/checksums 的兼容清单。
- `F-IND-R4-FP-008`: 明确 Registrar 与 `Module.Hooks` 的 Lifecycle/Observability 归属。
- `F-IND-R4-FP-009`: 不得以当前 store-ping `/readyz` 的 200 响应冒充终态 module readiness。

## 放行结论

在本审计对应的草案版本上，FP-001 至 FP-004 均为 open required；C1 不能冻结，
C2 不能开始，Root progress 不变。审计之后已对草案做修订；修订本身需要新的
independent re-review，不能回写本意见或提前关闭本意见中的 finding。
