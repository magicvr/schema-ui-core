---
id: A-004-w1-implementation-wave
doc: audit-entry
goal: GOAL-002-w1-examples-optional-module
source: self
status: closed
created: 2026-08-11
updated: 2026-08-11
version: 1.0.0
---

# A-004 · W1 实施波次审计（self · 关门准备）

> **闭合**：本意见无 required；recommended F-001/F-002/F-003 与 I-004 经 A-006 处置（fixed/保留），go 恢复经 A-006 留痕。

## 头字段

- **source**：self
- **auditor**：Claude Code（编排器自审）
- **类型 / scope**：execution-facts / W1 拆分与迁移实施产物 + 关门就绪（roadmap 阶段 5）
- **verdict**：**pass**（无未关闭 required；实施与 D-003 一致；回归证据完整）

## 范围与区间

审计 `E-004` 所记录的实施产物：kernel 模块拆分、`dev.examples` 包、manifest 装配、composition 条件装配 + home 推导、测试分母更新；对照 D-003 §1–§6 与成功标准 S1–S6。本审计不重审方案冻结（已在 A-001/A-002/A-003 闭环）。

## 成果（有证据）

- **S1 能力拆分**：`dev.examples` 持有 8 范例页 + Examples 导航（`modules/dev/examples/`）；`core.schema-render` 仅 CapabilitySchema/Validation 壳（`provider.go` 无页注册）。`schemarender/provider_test.go` 断言 0 页；`dev/examples/provider_test.go` 断言 8 页 + fragment。
- **S2 可选模块形态**：`dev.examples` 进 compiled 候选（`kernel/profile.go` BuiltinModules），`plan.HasModule("dev.examples")` 条件装配（`composition.go`），自有 schema + fragment。
- **S3 依赖剪枝**：`admin.*` / `core.manifest-route` 均不 DependsOn `dev.examples`；禁用时 admin 正常启动（`TestAppStartsAndStopsDualProfiles` 通过）。
- **S4 Profile 默认**：`profileDefaults` 不含 `dev.examples`；mvp/admin 默认无演示面。
- **S5 产品面与 home**：`TestManifestHomePageRefDerivation` 断言 mvp/admin home=users、无范例泄漏、**禁用时 `/api/schema/overview` 404**、dev.examples 启用时 home=overview + 范例恢复 + schema 200。
- **S6 回归与 go**：`go test ./...` 全绿；web 746 测试通过；mvp/admin e2e 3+3 passed；VP-008 `go` 消费**已暂挂**（E-004）。
- **装配层 home 推导**：`deriveHomePageRef` 按 D-003 §2 决策表实现；`StampHomePageRef` 以 map 级 app 块注入/删除 homePageRef，保留其余字段与 pages/navigation 原文。
- **fragment app 一致性**：baseline + 4 admin + dev.examples + probe fragment 的 `app` 块均无 homePageRef，`Aggregate` canonical 全等通过（composition/manifest 测试覆盖）。

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| S1 能力拆分 | 达成 | dev/examples 包；schemarender 0 页测试 |
| S2 可选模块形态 | 达成 | HasModule 装配 + dev/examples 包 |
| S3 依赖剪枝 | 达成 | BuiltinModules；双 Profile 启动测试 |
| S4 Profile 默认 | 达成 | profileDefaults 无 dev.examples |
| S5 产品面与 home | 达成 | TestManifestHomePageRefDerivation（含 schema 404） |
| S6 回归与 go 接口 | 达成 | 全测试 + 双 Profile e2e + E-004 go 暂挂 |

## Findings

### F-001 · web 渲染夹具仍模拟「范例启用」形态（low · recommended）
**描述**：`apps/web/src/test-fixtures/app-manifest.{admin,mvp}.json` 仍含 8 范例页 + `homePageRef: overview`，作为 Renderer 内部夹具（非 Profile 契约），内部自洽且测试通过。但 mvp 夹具与默认 mvp 现实不一致，可能误导后续维护者。
**证据**：`apps/web/src/test-fixtures/app-manifest.mvp.json`；`E-004` 已记录该取舍。
**建议**：非阻断。后续可增补「默认无范例」web 夹具在 Renderer 层覆盖卫生路径，或把现有夹具改名为 `dogfood` 语义。

### F-002 · StampHomePageRef 仅保留当前 envelope 顶层字段（low · recommended）
**描述**：`StampHomePageRef` 用 5 字段 struct 反序列化顶层 envelope 再重排，若未来 fragment 引入新顶层字段会被丢弃。当前协议 envelope 固定 5 字段，风险可控。
**证据**：`manifest.go` `StampHomePageRef`。
**建议**：非阻断。如协议增字段，回贴 `StampHomePageRef` 结构即可。

### F-003 · i18n 范例 key 不删（I-004 结论 · 非阻断）
**描述**：`manifest.title.overview` / `manifest.nav.examples` / `schema.overview.*` 等 key 在 `dev.examples` 启用（dogfood）路径仍被 fragment 引用，**并非死 key**。I-004「是否本波删除」结论 = **保留**（不需要清理）。
**证据**：`dev/examples/manifest/fragment.json` titleKeys；`apps/web/src/i18n/messages/*.json`。
**建议**：I-004 以「保留」闭合（non-blocking，无需动作）。

### F-004 · VP-008 `go` 恢复证据已齐备（建议恢复 · 待用户确认）
**描述**：E-004 §go 所列恢复证据：候选矩阵快照（本 commit 前后 BuiltinModules/profileDefaults）、digest（`4a2b8cd` + 本审计 commit）、双 Profile 烟测（mvp/admin e2e 通过）、新增断言（禁用无泄漏 + schema 404、homePageRef 正确、启用恢复）均已在 API 测试与 e2e 落盘。恢复证据充分。
**证据**：`TestManifestHomePageRefDerivation`；e2e 双 Profile；`E-004`。
**建议**：波次关门时由 `/govern` 留痕**恢复 VP-008 `go` 消费**（范围 = 本波变更后的模块矩阵），业务 VP 激活前仍须按 VP-008 完成消费前 freshness review。

## 必改项汇总（required）

- **无**。

## 结论 + 建议下一步

实施与 D-003 完全一致，S1–S6 达成，回归证据完整，**无开放 required**。建议：
1. 由 grok-build@grok-4.5 做 **independent 波次审计**（A-005）交叉核对实施产物。
2. 合并响应后进入 W1 关门：status `done`、goal-tree 同步、Root 波次台账归档、VP-008 `go` 留痕恢复。
3. 独立审计未触发前，本 self 意见不作为关门唯一依据（cross 补齐后闭环）。
