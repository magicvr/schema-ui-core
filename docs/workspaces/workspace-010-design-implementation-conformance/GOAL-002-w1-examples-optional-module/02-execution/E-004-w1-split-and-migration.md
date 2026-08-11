---
id: E-004-w1-split-and-migration
doc: execution-entry
goal: GOAL-002-w1-examples-optional-module
status: recorded
created: 2026-08-11
updated: 2026-08-11
version: 1.0.0
---

# E-004 · W1 拆分与迁移实施（roadmap 阶段 2/3）+ VP-008 `go` 暂挂触发

## 事实（2026-08-11）

按 D-003 §1–§6 完成拆分与迁移，**首个改变模块矩阵/装配语义的代码合入**：

### Kernel（`apps/api/internal/kernel/`）
- `core.schema-render` 移除 8 范例页贡献，仅保留 CapabilitySchema/Validation 能力壳（`profile.go` BuiltinModules）。
- 新增 compiled 候选模块 **`dev.examples`**（DependsOn `core.schema-render` + `core.navigation-capability`；Pages 8 + Fragment examples；**无** Provides/Routes/Permissions/system-data nav —— 横切演示模块，六面豁免，D-003 §3）。
- `profileDefaults`（mvp/admin）**未含** `dev.examples`（默认关闭，S4）；`core.schema-render` 保留在默认集（能力闭包，F-005 落实）。
- `core.manifest-route` 依赖边**核验结论**：其不消费 schema 能力，但为最小化矩阵变更，**保留** DependsOn `core.schema-render`；不属演示依赖（演示模块现为 `dev.examples`），S3 满足。

### 新模块 `apps/api/internal/modules/dev/examples/`
- `schema/`：8 个范例 schema 文档自 `schemarender/schema/` **git mv** 迁入（保留历史），新 `schema.go`（ModuleID `dev.examples`）。
- `provider.go`：注册 8 页 + `examples` fragment；无 HTTP/权限/导航系统数据贡献。
- `manifest/fragment.json` + `manifest.go`：范例页 + Examples 导航 + top overview；`app` 块**无 homePageRef**。
- `provider_test.go`：8 页 + fragment + 无六面断言。

### 装配与 Manifest
- `composition.go`：`core.schema-render` 与 `dev.examples` 均按 `plan.HasModule` **条件装配**（F-007）；`deriveHomePageRef` 实现 D-003 §2 决策表；发布前 `manifest.StampHomePageRef` 注入。
- `manifest.go`：新增 `StampHomePageRef`（app 块 map 级 stamp/删除，保留其余字段）；baseline `app-manifest.json` 瘦身为 core-only（空 pages/nav，app 无 home）；4 个 admin fragment 与 S2 probe fragment 的 `app` 块去除 `homePageRef`（保持 Aggregate canonical 全等通过）。
- `schemarender/provider.go`：能力壳，无页注册；`schemarender/schema/` 包删除。

### 测试分母（D-003 §6 勾选清单）
- `schemarender/provider_test.go` → 0 页断言；新增 `dev/examples/provider_test.go`（8 页）。
- `handler/testhelpers_test.go`：范例 schema 引用切到 `dev/examples/schema`；`handler/schema_test.go` 分母来源随之更新。
- `composition/s2_access_drill_test.go`：probe fragment app 去 homePageRef。
- `manifest/manifest_test.go`：sidebar 计数 3 → 2（baseline 去 Examples 组）。
- `composition/composition_test.go`：新增 `TestManifestHomePageRefDerivation`（mvp/admin home=users、无范例泄漏；mvp+dev.examples home=overview + 范例恢复）。
- web：`schemarender/schema` → `dev/examples/schema` 路径回贴（5 个测试文件）；e2e `shell/localization/schema-crud` 的 `/overview` 与 manifest 必含 overview 断言改 `users`，并新增「无 Examples 导航」S5 卫生断言。

### 回归证据
- `go test ./...`（apps/api）：**全部通过**。
- `npm test`（apps/web）：**44 文件 / 746 测试通过**。
- Playwright e2e：**mvp** 3 passed / 1 skipped；**admin** 3 passed / 1 skipped（双 Profile 烟测通过）。

## VP-008 `go` 消费有效性 —— 暂挂触发（D-003 §5）

本合入为**首个改变 Profile 默认集（无）/模块矩阵（新增 `dev.examples`、`core.schema-render` 页贡献迁移）/Manifest 装配语义（homePageRef 装配层 stamp、fragment app 去 home）**的 commit → 按 VP-010 与 VP-008 §`go` 消费有效性规则，**业务对旧 `go` 的消费暂挂**，直至下列恢复证据落盘并由 `/govern` 留痕：

1. 候选身份快照：profile/module 矩阵（mvp/admin 默认集；`dev.examples` 显式启用路径）。
2. digest：本 commit hash。
3. 双 Profile 烟测：mvp/admin e2e 已通过（本记录）。
4. 新增断言：禁用无 Examples 组/8 pageId/schema 404；homePageRef 推导正确；启用恢复。
5. 台账落点：本 `02-execution/E-004` + Root 台账 + VP-010 波次档案指针。

## 尚未发生

- W1 波次审计（roadmap 阶段 5：self + 按风险 independent/cross；D-003 已定 cross）。
- 业务 VP 激活前的消费前 freshness review（VP-008 §go 消费有效性）。
