---
id: A-001-r8-closeout-self
goal: GOAL-009-audit-envelope-and-session
doc: audit-entry
record_id: A-001
source: self
auditor: /govern · 会话编排
scope: GOAL-009 close-out；S0–S1 交付、三条成功标准、I-001/I-002、D-001、非目标与不变式、JWT sid、0048、writer envelope
audit_type: close-out
verdict: pass
status: recorded
parent: GOAL-009-audit-envelope-and-session
created: 2026-08-19
updated: 2026-08-19
version: 0.1.0
---

# A-001 · GOAL-009 关门自审（2026-08-19）

- **source**：self
- **auditor**：`/govern` 会话编排
- **类型**：close-out
- **scope**：`workspace-012-shared-cross-module-contracts` / `GOAL-009-audit-envelope-and-session` 关门。核对 S0/S1、三条成功标准、I-001/I-002、D-001、非目标（不做 impersonation、不改 Profile/协议 pin、无归档查询 UI）、JWT `sid`、迁移 0048、writer envelope。
- **verdict**：**pass**
- **required findings**：0
- **日期**：2026-08-19

## 范围与区间

- **工作区**：`workspace-012-shared-cross-module-contracts`（id 与路径一致；Root `GOAL-001-shared-cross-module-contracts`；`canonical_scope` 匹配；`shared_materials_catalog: none`；`vision_role: delivery`；`primary_plan` = VP-012，`vision_ref` = `schema-ui-core-admin-foundation@0.2.0`）。
- **covered**：D-001、E-001；JWT `sid` / `User.SessionID`；`operation_log_session` + archive 拷贝；生产 mutation writer `NewDetail`/`recordAudit`；operations 读路径 `sessionId`；本会话已跑通的测试。
- **excluded**：不改 status/progress；不把本条当成 independent；不审 impersonation 产品化；不把 `progress: 67` 当闭合证据。
- **共享资料**：无引用。
- **审计模式**：JWT 声明 + 审计侧表 + 迁移 0048 属 security/data/migration 高影响门禁 → `independent`；项目路径要求 self 前置后再 grok-build independent。本条为 self 前置。关门仍待 independent。

## 本轮复验（2026-08-19）

| 命令 | 结果 |
|------|------|
| `apps/api`：`go test ./internal/modules/operationlog ./internal/auth ./internal/store -count=1` | **ok**（13.074s / 27.511s / 44.500s） |
| `apps/api`：`go test ./internal/handler ./internal/composition ./internal/modules/wallet -count=1` | **ok**（250.544s / 20.111s / 2.470s） |

未把派生 progress 当作完成证据。工作树尚未单独成 commit；本条以当前代码与上述复测为权威。

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| 1. 生产 mutation 写路径不再用手写 JSON 拼 detail；新写入可被 `ParseDetail` 解析（无 body 保持无 detail） | **达成** | `auditDetail`/`NewDetail`：auth/users/users_state/settings/roles/MFA/wallet/import/export/datapermission/captcha/dictionary/account/service-credentials/composition use recorder。filelibrary/recyclebin/scheduledtasks/dictionary create·update 仍 `detail=nil`。测试：`TestRolesOperationLogEvents`、`TestBrandingPublicAndSettingsPatch`、`TestR2CorrelationIDPersistsOnUsersOperation`、handler auth ops `ParseDetail` |
| 2. 用户态新审计行在已签发 sid 的请求上带 session_id；service-credential 行带 credential id | **达成** | `SignAccessToken` 写 `sid`；middleware 灌 `User.SessionID`；`recordAudit`/`identitySession`；service-credential `SessionID=credential.ID`。测试：`TestLoginSuccess`、auth ops `SessionID` 非空、users/settings session、`TestServiceCredentialManagementAndAuthentication` use=credential id、`TestRecordOperationPersistsSessionAndCorrelation`、`TestMigrate0048AddsSessionSideTable`、retention archive session |
| 3. 不改 Profile/模块矩阵/协议 pin；不做 impersonation | **达成** | 本波 git 变更不含 Profile / 模块矩阵 / `provenance-v2.8.json`；D-001 effective actor = 当前 `actor_id`；无 impersonation 列或切换面 |

## 信息门禁

| ID | 级别 | 最晚阶段 | 登记 | 本条 |
|----|------|----------|------|------|
| I-001 | required | S0 | verified | 维持；D-001：refresh token id via JWT `sid`；机器凭据 = credential id；effective actor = actor |
| I-002 | required | S0 | verified | 维持；D-001：全部生产 mutation 写路径改 `NewDetail` |

无 `deferred` required。无 `accepted-residual`。

## Findings

### F-001 · recommended · 若干 mutation writer 仍传 `ctx=nil`

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：`rolesOnWrite` / import / captcha / users_state / account_self / wallet 走 `recordAudit(..., nil)`，即使 handler 持有 request context 也不写 correlation。对照 HEAD，这些路径转换前也未写 correlation，不是本波回归。不阻断 session 成功标准。
- 证据：`roles.go` `rolesOnWrite` 把 `context.Context` 标为 `_`；`import.go` / `captcha.go` / `users_state.go` / `account_self.go` / `wallet.go` 末参 `nil`。export / datapermission / 部分 MFA / users.go 仍传 ctx。

### F-002 · recommended · Activity schema 未展示 `sessionId`

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：list/detail JSON 与 CSV 已有 `sessionId`，与 `correlationId` 同形；`activity.json` 表格与 recordView 仍不展示二者。D-001 非目标含「不做归档查询 UI」，不把 schema 列作为本波交付。不阻断关门。
- 证据：`operations.go` `operationToMap`；`operations_export.go` headers；`activity.json` columns/fields 无 `sessionId`。

## 必改项汇总

无。开放 required = 0。

## 结论 + 建议下一步

self close-out **pass**。三条成功标准可核对；I-001/I-002 verified。按项目路径仍须 grok-build independent 关门审；**不得**仅凭本条把 GOAL-009 标 `done`。

## 声明

本意见 `source: self`。不修改 status/progress。independent 由项目路径执行。
