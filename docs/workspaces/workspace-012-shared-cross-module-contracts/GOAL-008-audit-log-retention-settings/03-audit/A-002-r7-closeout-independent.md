---
id: A-002-r7-closeout-independent
goal: GOAL-008-audit-log-retention-settings
doc: audit-entry
record_id: A-002
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
scope: GOAL-008 关门。核对 S0/S1 交付、三条成功标准、I-001/I-002、D-001/D-002、A-001 self、非目标（无归档查询 UI、不改 Profile/协议 pin、不做 session/effective actor）、ApplyRetention archive/delete、sweeper 是否只读设置、迁移 0046/0047
audit_type: close-out
verdict: pass
status: recorded
parent: GOAL-008-audit-log-retention-settings
created: 2026-08-19
updated: 2026-08-19
version: 0.1.0
reviews: A-001
---

# A-002 · GOAL-008 关门独立审计（2026-08-19）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high；项目级路径见 `docs/architecture/independent-audit-execution.md`）
- **类型**：close-out
- **scope**：`workspace-012-shared-cross-module-contracts` / `GOAL-008-audit-log-retention-settings` 关门。核对 S0/S1 交付、三条成功标准、I-001/I-002、D-001/D-002、A-001 self、非目标（无归档查询 UI、不改 Profile/协议 pin、不做 session/effective actor）、ApplyRetention archive/delete、sweeper 是否只读设置、迁移 0046/0047。
- **verdict**：**pass**
- **required findings**：0
- **日期**：2026-08-19

## 范围与区间

- **工作区**：`workspace-012-shared-cross-module-contracts`（`workspace.md`：`id` 与路径一致；`root_goal` = `GOAL-001-shared-cross-module-contracts`；`canonical_scope` 覆盖本目标；`shared_materials_catalog: none`；`vision_role: delivery`；`primary_plan` = `VP-012-shared-cross-module-contracts`）。
- **covered**：GOAL-008 `00-meta` / D-001 / D-002 / E-001 / E-002 / A-001；实现与测试：settings GET/PATCH/reset 投影与校验、`settings.json` Audit log 页签、`ApplyRetention`、`StartRetentionSweep`、composition `loadPolicy`、0046/0047、本轮定向复测。
- **excluded**：不改 `status` / `progress` / goal-tree / 方案正文 / 业务代码；不读取或比较其他工作区治理上下文；不把 `progress: 67` 当闭合证据；不把本条写成 `/govern` 或建议本审计员改代码关门。
- **共享资料**：无固定引用；不得当作事实或 finding 关闭依据。
- **审计模式**：D-002 冻结 `independent` + grok-build；A-001 为项目路径 self 前置（pass / required=0）。本条为独立关门审。

## 本轮复验（2026-08-19）

| 命令 | 结果 |
|------|------|
| `apps/api`：`go test ./internal/modules/settings/repository ./internal/modules/operationlog ./internal/handler -run TestRepositoryOperationLogRetentionPatch\|TestApplyRetention\|TestSettingsValidationAndReset -count=1` | **ok**（2.453s / 3.725s / 3.169s） |
| `apps/api`：`go test ./internal/store -run TestCompiledMigrationCatalogOwnership\|TestMigrateFreshDB\|TestRestartPersistence -count=1 -v` | **PASS**（catalog 含 0046 `site_operation_log_retention` / 0047 `operation_log_archive`；fresh/restart 均为 47 条、末条 archive） |
| `apps/web`：`npm test -- --run src/i18n/schema-keys.structural.test.ts src/renderer/representative-pages.test.tsx` | **2 files / 12 tests passed** |

未把派生 `progress` 当作完成证据。工作树中 GOAL-008 实现尚未单独成 commit；本条以当前代码与上述复测为权威，不以未存在的 SHA 充当交付证明。

## 工作区与对齐（只读）

| 检查项 | 结论 | 证据 |
|--------|------|------|
| 工作区绑定 | 通过 | `workspace.md` Root / canonical / `plan_refs`+`primary_plan` 与 GOAL-008 `parent`、`primary_plan` 一致；`goal-tree.md` 含本目标且 `status: active` |
| 共享资料引用 | 无引用，不构成关闭证据 | `shared_materials_catalog: none` |
| 对齐链 | 未发现与 Root R7 / VP-012 的明显冲突 | Root 纲领 R7 指针指向本目标；交付为设置可改保留 + 过期归档/删除，无 Tier D |
| Vision Review required | 本 scope 未见开放 required | 本意见不写 `docs/vision/reviews.md` |
| 既有 Goal 审计 | A-001 self = pass；开放 required = 0；F-001 recommended 仍 open | `03-audit.md` |
| P-004 冲突 | 无互否必改项 | A-001 与本条均为 pass；recommended 可叠加 |

## 成果（有证据）

| 主张 | 证据 | 核验 |
|------|------|------|
| S0 冻结：默认 90 / archive、1–3650、`archive`\|`delete`、每小时读当前设置 | D-001；`settings/migration/migration.go` `DefaultOperationLogRetentionDays=90`、`DefaultOperationLogExpirationAction=archive`、Min/Max | 通过。未选硬编码 90、默认 delete、0=永久 |
| 0046 给 `site_settings` 加两列，DEFAULT 90 / `'archive'` | `settings/migration/migration.go` L131–154；catalog checksum `5038cd0d…` | 通过。本轮 `TestCompiledMigrationCatalogOwnership` + `TestMigrateFreshDB` |
| 0047 建 `operation_log_archive` + `operation_log_archive_correlation`；不重建热表 | `operationlog/migration/migration.go` L373–409；catalog checksum `6228b4e8…` | 通过。加法迁移，不触 0045 式 correlation rebuild |
| GET/PATCH/reset `/api/settings/default` 投影含两字段 | `handler/settings.go` `settingsRow` L123–124；`settingsDetail` / `settingsPatch` / `settingsReset` | 通过。公共 branding 投影**不含**这两项（`brandingRow` L95–107） |
| 默认 90 / archive；非法天数/动作拒绝 | repository 默认与 `ErrInvalidRetentionDays` / `ErrInvalidExpirationAction`；handler `writeSettingsError` → 400 `INVALID_RETENTION_DAYS` / `INVALID_EXPIRATION_ACTION`；`errorcatalog` + zh-CN/en-US | 通过。本轮 `TestRepositoryOperationLogRetentionPatch`、`TestSettingsValidationAndReset` |
| 设置页 Audit log 页签可改；页面级恢复默认回到 90 / archive | `settings.json` `tab-audit` + `updateAudit` + `resetSettings`；i18n 字段/选项键成对 | 通过。schema-keys 含 `settings.json`；representative-pages 含 `settings` |
| sweeper 每轮 `loadPolicy()`，天数/动作来自 `GetSiteSettings()` | `composition.go` L580–589；`retention.go` `StartRetentionSweep` L71–80 | 通过。未见硬编码 90 或 `archive` 作为 sweep 策略 |
| archive：先冷存行+correlation，再删热表；delete：不写归档 | `retention.go` L34–61；热表 correlation `ON DELETE CASCADE`（0041 DDL）；store `PRAGMA foreign_keys=ON` | 通过。本轮 `TestApplyRetentionArchivesThenRemovesHotRows`、`TestApplyRetentionDeleteDoesNotArchive` |
| 非目标：无归档查询 UI / 恢复 API | handler 无 archive 路由；activity 仍为 `GET /api/operations` 热表 | 通过 |
| 非目标：未改 Profile 默认集 / 模块矩阵 / 协议 pin | 本目标工作树 diff 不含 `kernel/profile.go`、`manifest/`、`apps/web/public/protocol`、`docs/schemas`；settings 路由/pages/nav 未扩 | 通过。改动限于 settings schema/i18n 与 API 字段 |
| 非目标：未做 session / effective actor / writer envelope | `retention.go` 无 session/actor 投影；设置 PATCH 仍走既有 `settings.update` + `requirePermission` | 通过 |

## 对照成功标准

| 标准 | 本轮 | 证据 |
|------|------|------|
| 1. GET/PATCH/reset `/api/settings/default` 暴露 `operationLogRetentionDays` 与 `operationLogExpirationAction`；默认 90 / archive；非法值 400 | **达成** | 三端点共用 `settingsRow`；DDL/repository/reset 默认 90/archive；天数 0 → HTTP 400 `INVALID_RETENTION_DAYS`；动作 `compress` → repository `ErrInvalidExpirationAction`，handler 映射 400。HTTP 层未单测非法 action，见 F-002 |
| 2. 设置页 Audit log 页签可改这两项；恢复默认回到 90 / archive | **达成** | `tab-audit` inputNumber 1–3650 + select archive/delete；`updateAudit` PATCH 两字段；页级 `resetSettings` → POST reset，handler 断言回到 90/archive |
| 3. sweeper 只读设置，不硬编码天数或动作；archive 与 delete 有仓库测试 | **达成** | composition 闭包每轮读 `GetSiteSettings()`；`StartRetentionSweep` 把 `policy.Days`/`policy.Action` 传入 `ApplyRetention`；两仓库测试本轮复测通过。间隔 `time.Hour` 是 D-001 冻结节奏，不是天数/动作硬编码 |

## 信息门禁

| ID | 级别 | 最晚阶段 | 登记 | 本条 |
|----|------|----------|------|------|
| I-001 | required | S0 | verified | 维持。用户书面策略 + D-001（90 / archive / 可改 / 不硬编码） |
| I-002 | required | S2（S1 实施前定模式） | verified | 维持。D-002：`independent` + grok-build；先 self。本条即该模式的独立关门意见，不回写信息项状态 |

无 `deferred` required。无 `accepted-residual`。无到期且影响关门的未关闭 required 信息项。

## 非目标与不变式

| 项 | 本轮 |
|----|------|
| 不提供归档查询 UI / 恢复 API | 成立。冷表仅 sweeper 写入 |
| 不改 Profile 默认集、模块矩阵、协议 pin | 成立。见上表 |
| 不做 session / effective actor 或其余 writer envelope | 成立 |

## Findings

### F-001 · recommended · `StartRetentionSweep` 无单独单测

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：与 A-001 F-001 同范围。间隔兜底 `time.Hour`、启动即跑一次、`loadPolicy` 失败只打日志跳过、`stop()` 关停，均无单测。策略加载与 `ApplyRetention` 已测；composition 已接线。不阻断关门。
- 证据：`retention.go` L65–104；`retention_test.go` 仅覆盖 `ApplyRetention`；`*_test.go` 无 `StartRetentionSweep`。

### F-002 · recommended · HTTP 未覆盖非法 `operationLogExpirationAction`

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：成功标准「非法值 400」在实现上对动作成立（repository 拒绝 + `writeSettingsError` → `INVALID_EXPIRATION_ACTION`）。`TestSettingsValidationAndReset` 只打了天数 0；非法动作仅在 `TestRepositoryOperationLogRetentionPatch`。不构成名不副实，但不把 HTTP 400 动作路径当作已黑盒锁死。
- 证据：`settings_test.go` L293–298；`repository_test.go` L205–208；`settings.go` L304–307。

## 必改项汇总

无。开放 required = 0。

## 与既有意见的异同

| 项 | A-001 self | 本条 independent |
|----|------------|------------------|
| verdict | pass | pass |
| 三条成功标准 | 达成 | 达成；独立复测 + 代码核对 |
| I-001 / I-002 | verified | 维持 verified |
| F-001 sweeper 无单测 | recommended / open | 同意，编号为本条 F-001 |
| HTTP 非法 action | 未单列 | 本条 F-002 recommended |
| 非目标 / 0046/0047 | 覆盖 | 独立确认；另核 FK ON + CASCADE |

无必改互否。A-001 不得单独把 GOAL-008 标 `done`；本条通过后仍须 `/govern` 响应并改状态。

## 结论 + 建议给编排器/用户的下一步

independent close-out **pass**。S0/S1 交付、三条成功标准、I-001/I-002、非目标、ApplyRetention archive/delete、sweeper 只读设置、0046/0047 均可重复核对。开放 required = 0。

建议用 **`/govern`** 响应本条（及 A-001 F-001）：闭合 recommended 或明确不纳入本波，再决定是否把 GOAL-008 标 `done` 并更新 goal-tree。本意见不改 status/progress，也不建议本审计员改代码。

## 声明

本意见 `source: independent`。不修改 status/progress/goal-tree/方案正文/业务代码。响应由 `/govern` 处理。
