---
id: A-002-w1-plan-freeze-independent-grok
doc: audit-entry
goal: GOAL-002-w1-examples-optional-module
source: independent
status: closed
created: 2026-08-11
updated: 2026-08-11
version: 1.0.0
---

# A-002 · W1 方案冻结独立审计（grok-build · 实施前）

> **闭合**：required F-001～F-004 经 A-003 响应 → fixed（D-003）；recommended F-005/F-006/F-007 → fixed，F-008 保留 deferred（I-004）。

## 头字段

- **source**：independent
- **auditor**：grok-build@grok-4.5（reasoning-effort high；headless 只读，未改任何文件）
- **类型 / scope**：design-plan / W1 方案冻结 D-002 的完备性、可行性、设计意图对齐、测试分母、VP-008 `go` 消费有效性（实施前）
- **verdict**：**conditional**（方向正确、I-001～I-003 门禁已关；但 homePageRef「fragment 覆写」与现行 `Aggregate` 冲突语义不相容，且 `dev.examples` 模块形态与 home 推导算法未钉到可实施粒度）

## 范围与区间

对照 `00-meta` S1～S6、路线图阶段 1→2、D-001/D-002、E-001/E-002、VP-010 / Charter 边界；代码只读证据（composition / kernel / manifest / config / web protocol / 相关测试）。未做代码修改、status/progress 变更、落盘（本意见由编排器代贴，保留 `source: independent`）。

## 成果（有证据）

| 成果 | 证据 |
|------|------|
| 缺口 G1–G5 已盘点并与代码一致 | `E-001`；`composition.go:179`；`profile.go:24-47`；`module.go` BuiltinModules；`app-manifest.json` |
| 方案参数 I-001～I-003 用户确认并冻结 | `D-002`；`01-decision.md` 信息表 verified |
| 目标态方向正确（可选模块、默认关演示、能力与 demo 面解耦） | `D-002` §1–4；VP-010 产品面卫生；Charter 单主线/Profile |
| VP-008 `go` 接口需触发暂挂/重验证（方向级） | `D-002` §6；VP-010 接口表 |
| 测试分母调整有原则性要求 | `D-002` §5；`00-meta` S6 |

## 对照成功标准（方案冻结充分性）

| 标准 | 充分性 | 评注 |
|------|--------|------|
| S1 能力拆分 | 基本充分 | 8 pageId + Examples 拆到 dev.examples；core.schema-render 保留 Schema/Validation |
| S2 可选模块形态 | 部分充分 | 缺 BuiltinModules 条目形态、包路径、M1–M6 六面 vs 横切演示模块 |
| S3 依赖剪枝 | 基本充分（前提：schema-render 仍为能力模块） | 方案未**显式**写死「mvp/admin 默认仍含 core.schema-render」，需钉死防误拆 |
| S4 Profile 默认 | 充分 | 默认不含 dev.examples；显式启用路径与 ResolveProfile 模型一致 |
| S5 产品面与 home | **不充分（阻断）** | F-001/F-002 |
| S6 回归与 go | 部分充分 | 缺可枚举测试分母清单、go 暂挂触发时机/恢复字段 |

## 特别核查（决定性证据）

- **a) homePageRef 与 fragment 覆写**：`manifest.go:116-125` `Aggregate` 对 `app` 做 canonical 全等，首个非空 app 胜出、后续必须全等否则 `app identity conflict`（`manifest_test.go` 明确拒绝 app 冲突）。因此 `dev.examples` **无法**经 fragment 覆写 `homePageRef`。现行 baseline + 各 admin fragment + S2 probe **全部硬编码** `homePageRef: overview`（`app-manifest.json:10`；`users/manifest/fragment.json:10` 等同构；`s2_access_drill_test.go:100`）——任何派生 home 都必须全体同步。Web consumer 要求 pages 非空必须声明且指向已启用页（`app-manifest.ts:479-496`）。建议：(A) 组合根统一计算并写入 baseline + 全部 fragment 的同一 app 块；或 (B) 扩展 `Aggregate` 允许 home 合并（更大装配语义变更）。
- **b) schema-render / manifest-route / admin DependsOn**：`core.schema-render` 保留能力则依赖图可自洽（默认仍启用）；`admin.activity` 无 schema-render DependsOn 但 Requires StandardAdminCapabilities（含 Schema），能力闭包仍依赖 schema-render 在启用集。风险：误把「演示可选」理解成「整个 schema-render 移除」→ admin 解析失败。
- **c) 测试分母（已核实断言面）**：`schemarender/provider_test.go:23`（8 范例页）；`handler/schema_test.go`（overview 端点分母）；`composition/s2_access_drill_test.go:97-101,140`（probe app 块 + DependsOn）；`kernel/kernel_test.go:21-28`（mvp 默认列表）；`manifest/manifest_test.go`（ForModules 页投影）；web e2e `shell.spec.ts:20,44-45`、`localization.spec.ts:43,91,99`、`schema-crud.spec.ts:16`（默认 Manifest 必含 overview、home→`/overview`）——默认关演示将破坏，需双路径。
- **d) VP-008 `go`**：规则权威见 VP-008；缺口 = 暂挂触发时机（冻结时 vs 首个矩阵代码合入时）、恢复证据最低字段、台账落点未钉死。
- **e) dev.examples 自身形态**：缺 DependsOn / Provides / Contributions 键 / schema 包路径 / 是否豁免标准 Admin 六面。

## Findings

### F-001 · high · required
**描述**：D-002 声称「`dev.examples` 启用时可经 fragment **覆写** `homePageRef` 回 `overview`」与现行 Manifest 聚合模型**不相容**。`Aggregate` 要求所有 fragment 的 `app` canonical 全等、冲突 fail closed，无法字段级覆写；且 Web 要求非空 pages 时 home 必须指向已启用页。需改为「组合根统一推导并写入全部 fragment/baseline 的同一 app 块」或「扩展 Aggregate 合并语义」，否则 S5 在默认关/启演示两条路径均不可安全实现。
**证据**：`D-002` §2；`manifest.go:116-125`；`manifest_test.go:53`；`app-manifest.json:10`；`users/manifest/fragment.json:10`；`app-manifest.ts:479-496`。

### F-002 · high · required
**描述**：homePageRef「首个启用的 admin 功能页」的**确定性算法**未冻结：何谓「admin 功能页」（仅 `admin.*` 的 Pages？是否含 activity？排序键？）、与 `dev.examples` 同时启用时 overview 优先的单一决策表、「无 admin 功能页」fallback（fail closed vs 固定页 vs residual）。D-002 显式 defer 到拆分步 → S5 仍开放 required 设计门禁。
**证据**：`D-002` §2；`00-meta` S5；`module.go` BuiltinModules pages。

### F-003 · med · required
**描述**：`dev.examples` 模块契约不完整（DependsOn / Provides / Contributions 键 / schema 包路径 / 是否豁免标准 Admin 六面），不足以按 playbook M1–M3 对称落地，实施期易出现贡献冲突或错误 DependsOn。
**证据**：`D-002` §4；对照 `module-contribution-playbook.md` M1–M3 与 `admin.users` 正例；现行范例页在 `schemarender/schema/`。

### F-004 · med · required
**描述**：对 VP-008 `go` 的影响已定性，但**暂挂触发时机、恢复证据最低集、台账落点**未冻结。按 VP-010/VP-008，改变 Profile 默认集 / 模块矩阵 / Manifest 装配语义前应可判定 `go` 消费暂挂状态。
**证据**：`D-002` §6；`E-002`「尚未发生」；`VP-010` 与 VP-008 接口。

### F-005 · med · recommended
**描述**：应**显式**冻结：`mvp`/`admin` 默认集**继续包含** `core.schema-render`（能力壳），仅移除 `dev.examples`；避免实施误解导致 DependsOn/Requires 闭包失败。
**证据**：`profile.go:24-47`；`module.go` DependsOn/`StandardAdminCapabilities`；`D-002` 未写「默认仍含 schema-render」。

### F-006 · med · recommended
**描述**：测试分母仅原则性描述；硬分母至少包括 provider_test 8 页、handler schema 分母、S2 probe app 块、kernel 默认列表、manifest ForModules、Web e2e home→`/overview` 与 manifest 必含 overview。进入回归前应有勾选清单。
**证据**：见上文「测试分母」表；`D-002` §5。

### F-007 · low · recommended
**描述**：组合根 `providers := []kernel.Provider{schemarendermodule.New()}` 无条件装配（`composition.go:179`）——方案已要求改 `plan.HasModule`，实施时需同步测「未启用则不注册 / 无范例 schema」。
**证据**：`composition.go:179-195`；`D-002` §4。

### F-008 · low · recommended
**描述**：I-004（i18n 死 key）deferred non-blocking 合理，不阻断方案冻结。
**证据**：`00-meta` I-004。

## 必改项汇总（required）

1. **F-001**：修订 D-002——删除/改写「fragment 覆写 homePageRef」；冻结与 `Aggregate` 兼容的单一机制（推荐：组合根统一计算并写入 baseline + 全部 module fragment 的同一 app 身份；或显式变更 Aggregate 语义并评估 go 影响）。
2. **F-002**：冻结 home 推导表：admin 功能页集合、排序键、`dev.examples` 启用时是否固定 `overview`、无可用页时 fail closed/fallback。
3. **F-003**：冻结 `dev.examples` 最小模块契约（DependsOn/Provides/Contributions/schema 路径/非标准 Admin 六面豁免说明）。
4. **F-004**：冻结本波对 VP-008 `go` 的暂挂触发点与恢复证据/落盘位置（可与首个矩阵代码落地绑定，但必须书面）。

## 与既有意见的异同

`03-audit` 无历史 A 条目可对照；本意见为 scope 内首份 independent design-plan 意见（A-001 self 落盘于同日，编排器在合并响应中对照）。

## 结论 + 建议给用户 / 编排器的下一步

**结论**：D-002 在产品意图与 Profile 默认卫生上正确且 I-001～I-003 门禁已关，G1–G4 整改方向可验证；但 **S5 的 home/Manifest 装配路径存在与 as-built 聚合语义的硬冲突（F-001）**，模块形态与 go 暂挂执行细节不足（F-002～F-004）。Verdict = **conditional**：不可无条件进入完整「拆分与迁移」；可并行起草实现草稿，但合入/宣称阶段 2 前须先闭合 required findings。

**建议下一步（`/govern` 响应）**：1) 编排器汇总本意见；F-001～F-004 走 fixed（修订方案决策），勿用口头 residual 绕过 F-001；2) 用户确认 home 机制二选一（(A) 组合根统一 stamp，推荐 / (B) 扩展 Aggregate 合并）；3) 补「实施冻结附录」：dev.examples 契约 + home 算法表 + go 暂挂/恢复字段 + 测试分母勾选；4) 闭合后再开阶段 2。

## 声明

本意见不修改 status/progress/方案正文/goal-tree；响应、finding 闭合与阶段推进归 `/govern`；落盘由编排器代贴并保留 `source: independent`。
