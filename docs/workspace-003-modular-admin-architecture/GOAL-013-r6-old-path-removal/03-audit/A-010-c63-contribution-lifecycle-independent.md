---
id: A-010-c63-contribution-lifecycle-independent
doc: audit-entry
goal: GOAL-013-r6-old-path-removal
source: independent
auditor: Grok Build / grok-4.5 / high
date: 2026-08-06
scope: >
  C6.3 Schema document bytes, Configuration runtime, PolicyID/Visibility v1,
  dual-profile lifecycle implementation facts; re-review Root A-010 F-003b fixed
  candidate and R6-I003 C6.3 gate
audit_type: execution-facts | finding-closure | stage
verdict: pass
status: recorded
parent: GOAL-001-modular-admin-architecture
created: 2026-08-06
updated: 2026-08-06
version: 0.1.0
---

# A-010 · C6.3 Contribution 与生命周期实施事实独立交叉审计

- **source**：independent
- **auditor**：Grok Build / grok-4.5 / high
- **类型 / scope**：execution-facts | finding-closure | stage；D-003 四个 C6.3
  实现切片、Root A-010 F-003b fixed candidate、R6-I003 C6.3 门禁
- **verdict**：**pass**
- **方法约束**：本审计会话只读（workspace / goal ledger / `apps/api` 源码与测试），
  未执行 `go test` / `go vet`、写盘或修改 status/progress/goal-tree。E-011/`8b76ab0`、
  E-012/`2548e42`、E-013/`9896a02` 的动态验证视为已绑定 checkpoint 的执行台账证据；
  静态核验未发现反证，故不因未重跑动态测试而降级为 required 缺口。

## 范围与区间

| 项 | 值 |
|----|-----|
| 工作区 | `workspace-003-modular-admin-architecture` |
| canonical | `docs/workspace-003-modular-admin-architecture/` |
| Root | `GOAL-001-modular-admin-architecture` |
| 被审目标 | `GOAL-013-r6-old-path-removal` |
| 契约 | D-003（accepted） |
| 执行证据 | E-011、E-012、E-013 |
| 既有 self | A-009（pass，required 0；cross 待本条） |
| 排除 | C6.4、Root close-out、VP-003 closed、exit #1-#7 终态取证 |

## 工作区与信息门禁

| 核对项 | 状态 | 备注 |
|--------|------|------|
| workspace Root / canonical / plan_refs | pass | `workspace.md`：Root=`GOAL-001-…`，`primary_plan=VP-003`，`shared_materials_catalog: none` |
| goal-tree 绑定 | pass | GOAL-013 `active 2/4` 挂于 Root；R6 进行中 |
| R6-I001 / R6-I002 | verified（meta） | 不在本 scope 重开 |
| R6-I003 | `collecting`（正确） | 四切片实现完成；本 independent 前不得 `verified` |
| R6-I004 | collecting | C6.4；本条不触及 |
| 共享资料 | n/a | catalog=none；未把外部资料当证据 |

## 成果（有证据）

### 1. Schema document bytes · 单一发布路径

| 核对项 | 状态 | 证据 |
|--------|------|------|
| `PageContribution.Document` 注册期校验 | pass | `kernel/contribution.go` `validatePage`：JSON object、`meta.pageId==PageID`、可确定性重编码 |
| Registrar defensive copy | pass | `validatingRegistrar.Schema`：Document/Resources/Actions 写入前 copy |
| 模块 provider 贡献字节 | pass | `modules/schemarender` 五个 core 页；users/roles/settings/activity 各提交模块 embed |
| composition → handler 单一路径 | pass | `composition.go`：`handler.RegisterSchemas(mux, set.Pages)` |
| handler 只发布 finalized pages | pass | `handler/schema.go`：由 `[]PageContribution` 建 map 并再 copy；无静态 merge |
| 旧中心静态路径移除 | pass | `apps/api` 对 `staticSchemaDocuments` / `schemaOwnerMap` / `schemaDocumentsForPlan` 源码零命中 |
| Profile 隔离 | pass | composition 测试：mvp 不发布 settings/activity schema；admin 发布；manifest 页逐一可取 schema |
| 负向校验 | pass | `TestSchemaDocumentValidationAndCopy`：非法 JSON / 非 object / 缺 meta / 错 pageId / 防别名 |

F-003b 要求由 PageContribution（或等价已校验集合）发布 schema 字节并去掉 handler 对业务
模块 schema 包的生产编译期枚举。当前生产路径满足；`handler/testhelpers_test.go` 在测试中
组装 PageContribution 属测试夹具，不是生产双轨。

### 2. Configuration runtime contribution

| 核对项 | 状态 | 证据 |
|--------|------|------|
| 类型 + Registrar + set 字段 | pass | `ConfigurationContribution`、`Registrar.Configuration`、`ContributionSet.Configurations` |
| identity / grammar / defaults / validator | pass | `validateConfiguration`：Key==Namespace、ASCII dotted grammar、JSON object、validator 必填且 defaults 通过 |
| defaults copy | pass | Registrar 写入前 copy；防别名测试 |
| 全局冲突 + namespace 排序 | pass | finalize `checkUnique(KindConfiguration)` + `sortConfigurations` |
| `admin.settings` owner `settings.branding` | pass | descriptor、`settings/configuration` defaults/validator、provider registration |
| handler 无私有 namespace 常量 | pass | `SettingsRoutes(..., configNamespace)`；PATCH header 使用入参；生产 handler 无该字面量 |
| 负向 | pass | kernel 非法 namespace/defaults/validator/conflict；settings 空标题/未知键/坏 logo 测试 |

### 3. PolicyID / Visibility v1

| 核对项 | 状态 | 证据 |
|--------|------|------|
| 单一 ASCII 小写点分 policy ref | pass | `validDottedIdentifier`；PolicyID 与 Visibility 共用 |
| 拒绝未版本化表达式 | pass | grammar tests：布尔表达式、大写、空段、`--`、`_`、非 ASCII 等 |
| 合法现有 policy | pass | `system.admin`、`system.admin-editor`、`system.admin-editor-viewer` |
| kernel / auth-session 分层 | pass | kernel 不导入 auth-session；well-formed unknown policy 在 `rolesForPolicy` allowlist 于 reconcile 前 fail closed |

### 4. Dual-profile lifecycle

| 核对项 | 状态 | 证据 |
|--------|------|------|
| Start 失败只逆序清理已 Start，并清空 started | pass | `Runtime.Start` + matrix |
| Ready 失败逆序清理全部已 Start，并清空 started | pass | `Runtime.Ready` + matrix |
| 保留 structured code/module，cleanup 追加 detail | pass | `lifecycleFailure`；matrix 断言 structured error + cleanup detail |
| Stop continuation、首错、清空、重复 no-op | pass | `Stop`/`stopModules` + matrix |
| composition 只清 listener/store | pass | Ready 失败关闭 listener/store，不重复 `runtime.Stop` |
| readiness 仅成功后 | pass | `gate.setReady()` 位于 Start+Ready 成功之后 |
| 真实 mvp/admin Runtime 矩阵 | pass | `TestDualProfileLifecycleMatrix` 解析真实 Plan |
| 两 Profile Fx Start/Stop + 端口占用 | pass | `TestAppStartsAndStopsDualProfiles` 与端口占用 stable failure |

### 5. 动态验证台账（本会话未重跑）

| 来源 | 记录 | 本审计态度 |
|------|------|------------|
| E-011 | `go test ./...`、`go vet ./...` exit 0；旧符号零命中 | 静态一致，不降级 |
| E-012 | 全量 + 定向测试/vet | 静态一致，不降级 |
| E-013 | kernel/composition 定向 + 全量 + vet | 静态一致，不降级 |
| A-009 self | 同 scope pass | 本条独立复审同意 |

## Finding 闭合复审

| finding / 信息项 | 独立结论 | 证据边界 |
|------------------|----------|----------|
| Root A-010 F-003b · Schema document bytes 未由 ContributionSet 发布 | **具备 fixed 闭合证据** | E-011/`8b76ab0`；生产 `set.Pages` → `RegisterSchemas`；旧源码路径移除；Profile schema 回归。闭合动作归 Root `/govern` 响应。 |
| R6-I003 · Schema 字节贡献发布 + 收尾 | **可由 /govern 在响应后改为 `verified`** | D-003 四切片代码+测试+self A-009+本 independent 均无开放 required。 |
| C6.3 检查点 / progress `2/4`→`3/4` | **可勾选候选** | 须编排器响应；不自动扩大为 R6/Root/VP 完成。 |

## Findings

### Required

无（required = 0）。本 scope 未发现阻断 C6.3 实施事实认定、F-003b fixed 候选或
R6-I003 门禁放行的 required 缺口。

### Recommended

无（recommended = 0）。说明（非 finding）：未重跑 test/vet 已诚实标明；台账动态证据
与静态路径一致。R6-I003 `collecting`、progress `2/4` 是正确程序门禁状态。空目录
`handler/fixtures/schema/` 无源文件，不构成生产路径残留。

## 必改项汇总

| 级别 | 数量 | 内容 |
|------|------|------|
| required | 0 | — |
| recommended | 0 | — |
| 意见冲突 | 0 | 与 A-009 self 结论一致 |

## 与 A-009 异同

| 维度 | A-009 (self) | A-010 (本条 independent) |
|------|--------------|---------------------------|
| scope | 四切片实施事实 | 同，并显式 F-003b / R6-I003 闭合复审 |
| verdict | pass | pass |
| required | 0 | 0 |
| 方法 | 执行侧 + 自审 | 只读静态交叉；动态依赖 E-011～E-013 台账 |
| F-003b | fixed candidate；independent gate pending | 候选确认：可 fixed |
| R6-I003 | collecting | 响应后可 verified |
| 边界 | independent 前不放行 C6.3 | 同意；不扩大至 C6.4/Root/VP |

无冲突。两意见可一并响应。

## 结论与给编排器的下一步

1. C6.3 实施事实通过独立交叉审计：`verdict=pass`，required=0，recommended=0。
2. 建议 `/govern` 响应 A-009 + A-010：R6-I003 → `verified`；勾选 C6.3 并重算
   `progress: 3/4`；在 Root 台账将 A-010 F-003b 标为 `fixed`，且不改写历史 A-010 原文。
3. 不得因此勾选 C6.4、将 GOAL-013/Root 标 `done`，或宣称 VP-003 closed / exit #1-#7
   已齐。下一焦点为 C6.4 完整回归与退出 #1-#7 逐条取证。

## 声明

本意见 `source: independent`，不修改目标 status/progress/检查点/方案正文/goal-tree。
Finding 闭合、信息项状态变更与阶段推进由 `/govern` 响应处理。C6.3 pass 不等于
R6、Root 或 VP 关门。
