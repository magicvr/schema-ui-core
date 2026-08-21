---
id: A-002-r8-closeout-independent
goal: GOAL-009-audit-envelope-and-session
doc: audit-entry
record_id: A-002
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
scope: GOAL-009 close-out；S0–S1 交付、三条成功标准、I-001/I-002、D-001、A-001 self、非目标与不变式、JWT sid、0048、writer envelope、用户/service-credential session
audit_type: close-out
verdict: pass
status: recorded
parent: GOAL-009-audit-envelope-and-session
created: 2026-08-19
updated: 2026-08-19
version: 0.1.0
---

# A-002 · GOAL-009 关门独立审（2026-08-19）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型**：close-out
- **scope**：`workspace-012-shared-cross-module-contracts` / `GOAL-009-audit-envelope-and-session` 关门。核对 S0/S1 交付、三条成功标准、I-001/I-002、D-001、A-001 self、非目标（不做 impersonation、不改 Profile/模块矩阵/协议 pin、无归档查询 UI）、JWT `sid`、迁移 0048 `operation_log_session` + archive session、生产 mutation writer envelope（`NewDetail`/`recordAudit`）、用户 session / service-credential session。
- **verdict**：**pass**
- **required findings**：0
- **日期**：2026-08-19

## 范围与区间

- **工作区**：`workspace-012-shared-cross-module-contracts`（id 与路径一致；Root `GOAL-001-shared-cross-module-contracts`；`canonical_scope` 匹配；`shared_materials_catalog: none`；`vision_role: delivery`；`primary_plan` = VP-012）。本条未读取其他工作区治理上下文，也未把共享资料或 `progress: 67` 当作闭合证据。
- **covered**：D-001、E-001、A-001；JWT `sid` / `User.SessionID`；`operation_log_session` + `operation_log_archive_session`；生产 mutation writer `NewDetail`/`auditDetail`/`recordAudit`；`/api/operations` 与 CSV `sessionId`；本轮指定复测。
- **excluded**：不改 status/progress/goal-tree/D-001/业务代码；不做 impersonation 产品化；不审归档查询 UI；不以不存在的 commit SHA 充当交付证明。
- **共享资料**：无引用。
- **权威**：工作树未单独成 commit；以当前代码 + 本轮复测为权威。

## 本轮复测（2026-08-19）

在 `apps/api`：

| 命令 | 结果 |
|------|------|
| `go test ./internal/modules/operationlog -run "TestRecordOperationPersistsSessionAndCorrelation\|TestApplyRetention" -count=1` | **ok**（4.155s） |
| `go test ./internal/store -run "TestCompiledMigrationCatalogOwnership\|TestMigrateFreshDB\|TestMigrate0048AddsSessionSideTable" -count=1` | **ok**（3.165s） |
| `go test ./internal/handler -run "TestBrandingPublicAndSettingsPatch\|TestOperationLogStructuredFiltersAndExport\|TestR2CorrelationIDPersistsOnUsersOperation\|TestRolesOperationLogEvents\|TestServiceCredentialManagementAndAuthentication" -count=1` | **ok**（7.802s） |
| `go test ./internal/auth -run TestLoginSuccess -count=1` | **ok**（2.291s） |

复测失败会写成 required finding；本轮无失败。

## 成果（有证据）

1. **S0 冻结**：D-001 将 session 定义为 refresh token id（JWT `sid`）/ 机器凭据 credential id；effective actor = 当前 `actor_id`；envelope = 全部生产 mutation 写路径走 `NewDetail`；明确非目标。I-001/I-002 均 `required`、已 `verified`，无 deferred / residual。
2. **JWT sid**：`SignAccessToken` 写 `sid`；`ParseAccessToken` 读出；登录后 `acct.SessionID = rt.ID`；中间件 `acct.SessionID = parsed.SessionID`；service-credential 身份 `SessionID = credential.ID`；logout 返回被撤销 refresh token id。`TestLoginSuccess` 断言 user/token session 一致。旧 token 无 `sid` 时 session 可空，与 D-001 一致。refresh 旋转会签发新 refresh token，因而新 `sid`——这是 D-001「session = refresh token id」的直接后果，不是偏差。
3. **迁移 0048**：`operation_log_session` + `operation_log_archive_session`；catalog 冻结 checksum `1427328e…`；`TestMigrate0048AddsSessionSideTable` 覆盖 0047→0048 加法迁移、旧行 session 空、新行可写。retention archive 同步拷贝 session（`TestApplyRetention` 断言 archive session = `sess-old`）。
4. **writer envelope**：生产 mutation 写路径收敛到 `NewDetail`/`auditDetail`/`recordAudit`。本波 diff 去掉字典级联删除的手写 `json.Marshal` 与 wallet_self 字符串拼接 JSON；avatar 不再把裸 URL 当 detail。无 body 的事件（filelibrary / recyclebin / scheduledtasks / dictionary create·update / roles delete）仍 `detail=nil`。`ParseDetail` 覆盖 auth/users/settings/roles create·update。
5. **session 关联**：`recordAudit`/`identitySession` 从 `User.SessionID` 写入侧表；service-credential **use** 行 `SessionID = credential.ID`；人类 actor 的 create/revoke 带登录 session（非空）。读路径 list/detail/CSV 有 `sessionId`。后台 wallet Job 无登录会话，不写 session（E-001）。
6. **非目标**：本波 `git diff --name-only` 不含 kernel Profile、模块矩阵、`provenance-v2.8.json` / `provenance-v2.9.json`；`apps/api` 无 impersonation 列或切换面；无归档查询 UI。effective actor 仍为 `ActorID: user.ID`。

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| 1. 生产 mutation 写路径不再用手写 JSON 拼 detail；新写入可被 `ParseDetail` 解析（无 body 保持无 detail） | **达成** | `audit.go` `auditDetail`→`NewDetail`；auth/users/settings/roles/MFA/wallet/import/export/datapermission/captcha/dictionary delete/account/service-credentials/composition use recorder。`wallet_self.go` 去掉字符串拼接 JSON；`dictionary.go` 去掉手写 marshal。测试：`TestRolesOperationLogEvents`、`TestBrandingPublicAndSettingsPatch`、`TestR2CorrelationIDPersistsOnUsersOperation`、auth ops `ParseDetail` |
| 2. 用户态新审计行在已签发 sid 的请求上带 session_id；service-credential 行带 credential id | **达成** | JWT `sid` + middleware；`identitySession`；composition use `SessionID: use.CredentialID`。测试：`TestLoginSuccess`、auth ops `SessionID` 非空、users/settings session、`TestServiceCredentialManagementAndAuthentication`（use = credential id；create/revoke session 非空）、`TestRecordOperationPersistsSessionAndCorrelation`、`TestMigrate0048AddsSessionSideTable`、retention archive session。判定：准则中「service-credential 行」= 机器凭据 **actor** 的行（use）；人类管理 create/revoke 带管理员登录 session，与 D-001 一致 |
| 3. 不改 Profile/模块矩阵/协议 pin；不做 impersonation | **达成** | 工作树变更不含 Profile / 模块矩阵 / provenance pin；D-001 effective actor = 当前 `actor_id`；代码无 impersonation 面 |

## 信息门禁

| ID | 级别 | 最晚阶段 | 登记 | 本条 |
|----|------|----------|------|------|
| I-001 | required | S0 | verified | 维持；D-001：refresh token id via JWT `sid`；机器凭据 = credential id；effective actor = actor |
| I-002 | required | S0（影响门禁 S1） | verified | 维持；D-001 冻结「全部生产 mutation 写路径改 NewDetail」；S1 代码与本轮复测覆盖该范围 |

无 `deferred` required。无 `accepted-residual`。关门信息项已 verified。

## Findings

### F-001 · recommended · 若干 mutation writer 仍传 `ctx=nil`

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：`recordAudit` 仅在 `ctx != nil` 时写 correlation。下列生产路径把 context 丢掉或标为 `_`，即使 handler 持有 request context 也不写 correlation。session 仍由 `identitySession(user)` 写入，**不阻断**成功标准 2。对照 A-001 F-001：同意其结论；本条补全文件清单（A-001 未列 avatar / filelibrary / dictionary / recyclebin / scheduledtasks）。
- 证据：`roles.go` `rolesOnWrite` 将 `context.Context` 标为 `_` 且末参 `nil`；`import.go` / `captcha.go` / `users_state.go` / `account_self.go` / `account_avatar.go` / `wallet.go` / `dictionary.go` / `filelibrary.go` / `recyclebin.go` / `scheduledtasks.go` 末参 `nil`。对比：export / datapermission / MFA / users.go / settings / service-credentials 仍传 ctx。

### F-002 · recommended · Activity schema 未展示 `sessionId`

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：list/detail JSON 与 CSV 已有 `sessionId`，与 `correlationId` 同形；`activity.json` 表格与 recordView 仍不展示二者。D-001 非目标含「不做归档查询 UI」，成功标准也不要求 schema 列。不阻断关门。同意 A-001 F-002。
- 证据：`operations.go` `operationToMap`；`operations_export.go` headers；`activity.json` columns/fields 无 `sessionId` / `correlationId`。

## 必改项汇总

无。开放 required = 0。

## 与既有意见的异同

- **A-001 self close-out pass**：同意三条成功标准达成、I-001/I-002 verified、开放 required = 0。
- **A-001 F-001 / F-002**：同意均为 recommended / low / 不阻断关门。F-001 证据清单在本条扩展，结论不变。
- 无冲突、无新的 required。本条不关闭 A-001 findings（关闭权在 `/govern`）。

## 结论 + 建议给编排器/用户的下一步

independent close-out **pass**。S0/S1 交付可核对；三条成功标准有代码与本轮复测证据；I-001/I-002 维持 verified。开放 recommended 两项（F-001/F-002，含 A-001 同号项）不阻断关门。

建议编排器 `/govern`：汇总 A-001 + A-002；可将 GOAL-009 标 `done` 并同步 goal-tree / S2；F-001/F-002 可后续处理或用户书面 overruled/residual。本审计员不改 status、不改代码。

## 声明

本意见 `source: independent`。不修改 status/progress/goal-tree/D-001/业务代码。响应由 `/govern` 处理。
