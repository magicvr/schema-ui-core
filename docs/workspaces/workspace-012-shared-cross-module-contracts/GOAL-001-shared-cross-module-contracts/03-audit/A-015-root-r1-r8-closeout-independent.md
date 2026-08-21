---
id: A-015-root-r1-r8-closeout-independent
goal: GOAL-001-shared-cross-module-contracts
doc: audit-entry
record_id: A-015
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
scope: workspace-012 Root close-out；R1～R8 最终闭合链、四条方向成功标准、workspace/VP-012/Charter 对齐、I-001/I-002、A-014 self、开放门禁
audit_type: close-out
verdict: pass
status: recorded
parent: GOAL-001-shared-cross-module-contracts
created: 2026-08-19
updated: 2026-08-19
version: 0.1.0
reviews:
  - A-014
---

# A-015 · Workspace-012 Root independent close-out（R1～R8）（2026-08-19）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high；项目级路径见 `docs/architecture/independent-audit-execution.md`）
- **类型**：close-out
- **scope**：`workspace-012-shared-cross-module-contracts` Root `GOAL-001-shared-cross-module-contracts`。核对现行纲领 R1～R8、四条方向成功标准、工作区/VP-012/Charter 对齐、I-001/I-002、A-014 self、历史 A-001～A-013 开放门禁，以及 R7/R8 增量后的 Root 关门就绪。
- **verdict**：**pass**
- **required findings**：0
- **日期**：2026-08-19

## 范围与区间

- **工作区**：`workspace-012-shared-cross-module-contracts`（`workspace.md`：`id` 与路径一致；`root_goal` = `GOAL-001-shared-cross-module-contracts`；`canonical_scope` = `docs/workspaces/workspace-012-shared-cross-module-contracts/`；`shared_materials_catalog: none`；`vision_role: delivery`；`plan_refs` / `primary_plan` = `VP-012-shared-cross-module-contracts`，VP 已 `closed`）。
- **covered**：Root `00-meta` / A-014 / I-001/I-002；八个子目标 `00-meta` 与最终 A 条目；R6 fail-closed 现码；R7 `ApplyRetention` / sweeper `loadPolicy`；R8 JWT `sid` / `NewDetail` / session 侧表；本轮定向复测。
- **excluded**：不改 `status` / `progress` / `goal-tree` / `workspace` / 方案正文 / 业务代码 / `docs/vision`；不重开 VP-012；不读取或比较其他工作区；不把派生 `progress=100` 当作关门证据；不把 handler 全量包 VACUUM 超时重判为新的 required（见 A-008）。
- **共享资料**：目录为 `none`；无固定引用。
- **auditor 立场**：只出意见。响应与是否改 `done` 归 `/govern`。

## 本轮独立复验（2026-08-19）

在 `apps/api`：

| 命令 | 结果 |
|------|------|
| `go test ./internal/modules/operationlog ./internal/modules/settings/repository ./internal/auth ./internal/store ./internal/composition ./internal/modules/wallet ./internal/modules/authsession -count=1 -timeout 240s` | **ok**（operationlog 14.104s；settings/repository 14.095s；auth 33.055s；store 51.126s；composition 29.084s；wallet 1.349s；authsession 45.731s） |
| `go test ./internal/handler -count=1 -timeout 240s` | **超时**（4m0s）：卡在 `TestUsersPasswordPolicyPreservesBytesAndRevokesRefresh` 的 `store.snapshotBeforePending` / SQLite `VACUUM`。与 A-008「full API handler SQLite VACUUM 超时非阻断」同形；**不是** R7/R8 契约失败 |
| `go test ./internal/handler -run "TestServiceCredential\|TestSettingsValidationAndReset\|TestBrandingPublicAndSettingsPatch\|TestOperationLogStructuredFiltersAndExport\|TestR2CorrelationIDPersistsOnUsersOperation\|TestRolesOperationLogEvents" -count=1 -timeout 180s` | **ok** 11.352s |

在 `apps/web`：

| 命令 | 结果 |
|------|------|
| `npm test -- --run src/i18n/schema-keys.structural.test.ts src/renderer/representative-pages.test.tsx` | **2 files / 12 tests passed**（1.85s） |

## 工作区与对齐（只读）

| 检查项 | 结论 | 证据 |
|--------|------|------|
| 工作区绑定 | 通过 | `workspace.md` Root / canonical / `plan_refs`+`primary_plan` 与本 Root `parent: null`、路径一致；八个子目标 `parent` 均为完整 Root id |
| VP → Charter | 通过 | VP-012 `vision_ref` = `schema-ui-core-admin-foundation@0.2.0`；Charter `status: active`、版本 `0.2.0` |
| VP 状态 | 通过 | VP-012 已 `closed`（首波退出分母 R1～R6）。R7/R8 是移交项在本区的增量，不构成重开 VP，也不把方向表宽项写成已交付 |
| 共享资料 | 无引用 | `shared_materials_catalog: none` |
| Vision Review required | 本 scope 未见开放 required | `reviews.md` open required = 0；本意见不写 Vision Review |
| P-004 冲突 | 无 | A-014 self = pass / required=0；本条独立同意，无一要一否 |

现行 `00-meta` / `goal-tree` / `workspace.md` 均为 Root `active`、8/8 检查点完成。首波 A-003 的 `done/100` 投影已被 R7/R8 增量取代；本条审的是现行 8/8，不把 6/6 历史投影当作本轮已关门。

## 子目标闭合矩阵

| 阶段 | 子目标 | 最终闭合链 | 本轮独立核验 | 当前 required |
|------|--------|------------|--------------|---------------|
| R1 | GOAL-002 | A-001 self pass；`done` | frontmatter `done`；request-id / 错误包络仍在；非本轮新增缺口 | 0 |
| R2 | GOAL-003 | A-006 independent pass → A-007；`done` | 索引一致；结构化 detail / 脱敏 / correlation 仍在 | 0 |
| R3 | GOAL-004 | A-004 independent pass → A-005；`done` | 索引一致；wallet 包本轮 ok | 0 |
| R4 | GOAL-005 | A-012 independent pass → A-013；`done` | 索引一致；composition 本轮 ok | 0 |
| R5 | GOAL-006 | A-008 independent pass → A-009；`done` | 索引一致；composition 本轮 ok | 0 |
| R6 | GOAL-007 | A-009 F-001～F-005 fixed → A-010 close；Root A-010 F-010 `fixed` + A-012 pass / A-013 接收；`done` | 现码 `authenticateServiceCredential`：生产事务 `MarkServiceCredentialUsedWithAudit` 失败返回 503，无 `_ =` 丢弃；auth / authsession / `TestServiceCredential*` 本轮 ok | 0 |
| R7 | GOAL-008 | A-001 self / A-002 independent / A-003 close；`done` 3/3 | `ApplyRetention` archive 先拷行+correlation+session 再删热表；sweeper 每轮 `GetSiteSettings()`，不硬编码天数/动作；settings repository + operationlog + settings 定向 handler + Web schema-keys 本轮 ok | 0 |
| R8 | GOAL-009 | A-001 self / A-002 independent / A-003 close；`done` 3/3 | JWT `sid` / `User.SessionID`；0048 `operation_log_session`（store 末条 version 48）；`NewDetail` / `audit.go`；handler operations / auth / roles / users 定向测试本轮 ok | 0 |

R7/R8 子目标 recommended residual（sweeper 启停单测、HTTP 非法 action 映射、部分 writer `ctx=nil`、Activity schema 展示 `sessionId`）已在各自 A-003 点名非阻断。本条不把它们升为 Root required。

历史 Root required（A-005 F-008、A-010 F-010）保持 `fixed`；本轮未重开。

## 对照 Root 四条方向成功标准

| 标准 | 本轮 | 证据 |
|------|------|------|
| 每个契约有可验证实现路径 | **达成** | R1～R8 均有 `done` 子目标 + 最终审计链；R7/R8 现码与本轮定向复测可核对 |
| 至少一个真实模块或验证路径消费首波契约 | **达成** | operationlog / auth / settings / wallet / service-credential 仍为真实消费面；R7 由 settings + sweeper 消费；R8 由生产 mutation writer 与 `/api/operations` 消费 |
| Profile / 模块矩阵 / Manifest / protocol / 共同门禁语义不被意外改变 | **达成** | 各子目标非目标均排除这些变化；本轮复测未触及 Profile/pin 回归失败；未见相反证据 |
| Tier D 业务域不进入 Root | **达成** | R7 是审计日志生命周期；R8 是 envelope/session 横切；未新增业务域模块、导航或 fragment |

## 信息门禁

| ID | 级别 | 最晚阶段 | 登记 | 本条 |
|----|------|----------|------|------|
| I-001 | non-blocking | R1 开始前 | verified | 维持 |
| I-002 | required | Root 关门 | verified（分母已由 A-014 更新为 R1～R8） | 维持。八个子目标最终审计链与本轮复测足以支持「全部合法闭合且满足四条成功标准」。不把 A-014 的分母修订本身当成缺口 |

无到期未验证的 required。无 Root 级 `accepted-residual` / `user-overruled`。无信息冲突。

## Findings

无新的 required 或 recommended finding。

handler 全量包 VACUUM 超时沿用 A-008 的非阻断结论；本轮定向切片已通过，不另开 finding。

## 必改项汇总

无。

## 与既有意见的异同

- 与 A-014 self：同向 `pass`，required=0。本条补了独立复测，并单独核验 R6 fail-closed 现码与 R7/R8 实现。
- 与首波 A-002：同意当时 R1～R6 可关；不同意把 A-003 的 6/6 `done` 投影读成现行 Root 已关门。现行分母是 8/8。
- 与 A-008 / A-012：F-008 / F-010 闭合链仍可核对，不重开。

## 结论 + 建议给编排器/用户的下一步

现行 R1～R8 与四条方向成功标准可独立核对，independent verdict=`pass`，开放 required=0。可用 `/govern` 响应本条并决定是否将 Root 标 `done`。本意见不修改 status/progress。

## 声明

本意见不修改 status/progress；响应由 `/govern` 处理。
