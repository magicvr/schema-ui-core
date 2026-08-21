---
id: A-003-grok-r4-c2-contract-audit
doc: audit-entry
goal: GOAL-008-r4-c2-module-contract-extension
source: independent
auditor: Grok Build / grok-4.5
date: 2026-08-05
scope: R4-C2 contract implementation slice, freeze-package conformance, C2.1-C2.4 readiness
audit_type: close-out
verdict: conditional
---

# A-003 · Grok R4-C2 契约实施独立交叉审计

## 声明

本意见 `source: independent`，只读审计。不修改任何文件、`status` / `progress` /
检查点 / goal-tree。响应与关门由 `/govern` 处理。

## 范围与区间

- 目标：`GOAL-008-r4-c2-module-contract-extension`；工作区 `workspace-003`
- 契约权威：冻结包 `r4-c1-freeze-package-draft.md`（`accepted`）§2/§3/§4/§7
- 代码：`apps/api/internal/kernel/{contribution,provider,persistence,module}.go` + 测试
- 非范围：C3 业务迁移、composition 生产接线、readyz 真实 readiness、Root/VP 关门

## 成果（有证据）

1. **六类 Contribution 与 Key 规范语义（C2.1 主体成立）**：类型与冻结 §2.2 一致，
   Handler/Apply 仅 `net/http`/`database/sql`，无 Fx；Key 规范语义有实现与校验
   （HTTP "METHOD pattern"、Schema PageID、Auth Permission、Navigation NodeID、
   Manifest FragmentID、Persistence Name）；owner 校验、tombstone/Apply 互斥。
2. **Registrar 无 Persistence；CompiledPersistence 为唯一收集入口（C2.2 收集侧）**：
   `CollectPersistence` 对每个 compiled provider 调用一次，做唯一性/缺口/tombstone/
   reconcile 元数据校验 + 按 version 排序；失败 fail closed。
3. **RegisterContributions 双检与 fail-closed 丢弃（C2.3 主体）**：仅启用 provider
   注册；无 provider 模块跳过（§7 过渡态）；Register 仅写声明 Kind+Key；finalize
   全局冲突/引用/capability/确定性排序；失败丢弃整个集合。
4. **Framework-agnostic（C2.4 静态检查）**：kernel/modules 无 Fx import；业务模块仍
   中心注册，未误称 C3；readyz 表述诚实。
5. **测试与编译**：`go test ./internal/kernel/`、`go test ./...` 全 ok、`go vet`
   干净。

## Findings

### F-IND-C2-001 · Descriptor「完全匹配」仅校验 ID + Version（required / med）

- evidence: `provider.go` 仅 `desc.ID != module.ID || desc.Version != module.Version`
- impact: Dependencies/KernelAPIRange/Contributions/Provides·Requires 可与 Plan 不一致
  仍注册成功
- closure: 按冻结语义比对全规范字段并补 mismatch 测试，或书面 residual

### F-IND-C2-002 · ContributionKeys.Fragments 未纳入 Plan 级 validateContributions（required / med）

- evidence: `module.go` `validateContributions` 无 Fragments
- closure: 加入 Fragments 全局唯一校验 + 测试，或 residual 说明延后

### F-IND-C2-003 · C2.2 成功标准写了 ledger drift/unknown，实现仅为 catalog 静态校验（required / med）

- closure: 收窄条文为 catalog 静态校验、drift/unknown 挂 C3/C5；或实现并举证

### recommended（不阻断 C2 主体）

- F-IND-C2-004 · Manifest secrecy 未在 finalize 实现 → C3 前实现或 residual 到 C3
- F-IND-C2-005 · Ready 失败不反向 Stop → C3/C5 运行时矩阵
- F-IND-C2-006 · PolicyID/Visibility/JSON 缺校验器 → C3 接线时补最小校验
- F-IND-C2-007 · navigation parent 引用依赖遍历顺序 → 两遍扫描 + 乱序测试
- F-IND-C2-008 · 冲突错误未统一稳定 error code → 包装 `*kernel.Error`
  `CodeModuleContributionConflict`

## 独立结论

**verdict: conditional**。主体工程交付成立（冻结 §2 类型形状、Key 语义、Registrar
无 Persistence、CollectPersistence 唯一入口、Register 双检主干、§7 无 provider
跳过、framework-agnostic、测试/vet/fx、不宣称 C3/readyz）。三条 required
（F-IND-C2-001/002/003）指向契约保真与成功标准诚实性。GOAL-008 **不能无条件关门**；
GOAL-005 **不能仅凭当前 C2 状态放行 C3**，需先闭合 required（或书面 residual）。

**明确声明：本独立审计员未修改任何 status / progress / 检查点 / goal-tree / 文件内容。**
