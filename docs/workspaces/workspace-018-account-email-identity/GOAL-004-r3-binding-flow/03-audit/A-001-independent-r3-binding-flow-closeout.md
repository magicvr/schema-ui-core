---
id: A-001
doc: audit-entry
goal: GOAL-004-r3-binding-flow
source: independent
status: recorded
created: 2026-08-24
updated: 2026-08-24
version: 1.0.0
scope: R3 关门审计（绑定/校验流实现 · 迁移 0055 · 错误契约同步 · 最小页面；代码 checkpoint 0ae17f09）
verdict: conditional
auditor: grok-build (grok-4.6 · reasoning high)
---

# A-001 · R3 绑定/校验流独立关门审计（independent · 2026-08-24）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型** / **scope**：close-out · 绑定/校验流实现、迁移 0055 成对 DDL、错误契约同步、账号页最小绑定 UI（代码 checkpoint `0ae17f09`）
- **verdict**：**conditional**
- **开放 required**：1（F-001）

### 范围与区间

| 项 | 值 |
|----|-----|
| 工作区 | `workspace-018-account-email-identity`（`workspace.md`：`root_goal` = `GOAL-001-account-email-identity`；`canonical_scope` 匹配本目录；`shared_materials_catalog: none`；`primary_plan` = `VP-018-account-email-identity`） |
| 被审目标 | `GOAL-004-r3-binding-flow`（同区；`parent` = `GOAL-001-account-email-identity`） |
| audit_type | close-out |
| 对照 | GOAL-004 `00-meta` C1–C4；GOAL-004 D-001；GOAL-002 D-001 §2/§4/§5/§6；GOAL-003 A-001 F-002 / N-1 承接清单 |
| 信息门禁 | I-005 **required**（TTL 10 分钟 / 冷却 60 秒）、I-006 **non-blocking**（允许代填→pending）：用户书面裁决落在 GOAL-004 D-001（2026-08-24）。Root `01-decision.md` 已标 verified；Root `00-meta.md` 仍为 collecting（台账滞后，见 F-003） |
| 共享资料 | 无（`none`；未把资料目录当事实） |
| 代码基准 | `0ae17f092685d8c42a44f066755d50a325afc836`（与本会话 `HEAD` 一致） |
| 本审计未改 | 目标 `status` / 检查点 / 派生 `progress` / goal-tree / 方案正文 / 产品代码或测试 |

未读取或比较其他工作区上下文。本目标此前无 self / independent 条目（`03-audit.md` 仅占位行）。项目级 `docs/architecture/independent-audit-execution.md` 写「先自审再 independent」；AGENTS P-003 写 `independent` 不固定要求 self。本条按用户指定独立关门执行，不把缺 self 标为 required。

### 成果（有证据）

| 主张 | 证据 |
|------|------|
| 工作区绑定合格 | `workspace.md` id / root_goal / canonical_scope / `plan_refs`+`primary_plan` 与本目标路径一致；资料目录 `none` |
| 0055 成对 DDL（时间列 INTEGER\|BIGINT）+ 描述符 v55 + `ApplyPostgres` | `apps/api/internal/modules/authsession/migration/migration.go` `emailVerificationDDL` / `postgresEmailVerificationDDL` / Descriptors v55 / `migrateEmailVerification` + `migrateEmailVerificationPG` |
| 冻结 checksum 与 DDL 实算一致 | 目录行 `migrate_test.go` want 表 0055 = `1556bda28a7fb995807eea2b376a35ea79cf497fa76c73171b2973304ce5b754`；本会话按 `kernel.MigrationChecksum` 规则（normalizeSQL + `"0055:email-verification-challenges:v1"`）独立 sha256 复算 **一致** |
| 黄金断言随 catalog 头 = 55 | `identity.go` `completeFingerprintCatalogHead = 55`；`completeLostLedgerTables` / `postV1CatalogTables` 含 `email_verification_challenges`；`identity_test.go` `lockedHeadExtraTables[55]` + complete-no-ledger 夹具；`migrate_test.go` `len==55` 且 reopen `len(applied2)==55`；`operations_test.go` / `restart_test.go` 尾断言 |
| Bind 占槽 fail-closed + 原样存储 + 发信失败整单回滚 | `email_identity.go` `BindEmail`：`lower(email)` 他号 COUNT → `ErrEmailTaken`；trim 后原样写入；`Send` 在同一 `runner.Run` 内。`email_identity_test.go` `TestBindEmailReservesSlotAndSends` / `TestSendFailureRollsBindBack` |
| 校验流：6 位码、TTL 10 分钟、常量时间比较、失败计数独立事务、冷却 60 秒 | 常量 `emailCodeTTL` / `emailResendCooldown` / `emailMaxFailedAttempts`；`subtle.ConstantTimeCompare`；`registerFailedAttempt` 走独立 `Run`（避免 runner 对错误回滚吃掉计数）。测试：`TestVerifyEmailHappyPathAndIdempotentVerified` / `TestVerifyEmailExpiredDropsChallenge` / `TestFailedAttemptsVoidChallenge` / `TestResendCooldown` |
| 换绑 = 覆写释放旧址（合同 §5） | `TestRebindOverwriteReleasesOldSlot`：`old@example.com` 覆写后他号可绑 `OLD@example.com` |
| 无挑战直达 verified | 仓储 `TestAdminPrefillPendingCannotVerifyDirectly`：代填后无挑战，`VerifyEmail` → `ErrEmailNotPending`；清空 → `(nil, nil)`。`UserPatch` 仅 `Email *string`，无 status 字段 |
| HTTP 自助三端点 + Catalog 七码 + 冻结清单 | `handler/email_identity.go` `POST /api/account/email/bind\|verify\|resend`（`auth.Middleware`，无权限键）；`errorcatalog.Catalog` 七码；`error_contract_test.go` `frozenDomainCodes` 同步；`kernel/profile.go` 与 `account/provider.go` 路由清单同步；`composition.go` 向 `accountmodule.New` 注入 `mailSender` |
| 最小页面挂账号 profile tab；未复制设置页邮件 tab | `account.json` `email-identity-block` `component: "email-identity"`；`email-identity.tsx` 注册同名；`main.tsx` 导入。`settings.json` `tab-mail` / `mail-admin-tab` 未改（仍为出站渠道管理，非身份绑定复制品） |
| i18n 双语键齐 | `zh-CN.json` / `en-US.json` `schema.account.email.*` 各 15 键（非 E-002 所写「12」；见 N-5） |
| 本会话复跑 | `go test ./internal/store/ ./internal/modules/authsession/... ./internal/handler/ -count=1` **exit 0**；`TestFullCatalogPostgresBootstrapIntegration` **PASS**（1.91s，未 skip）；web `npx vitest run src/components/email-identity.test.tsx src/i18n/schema-keys.structural.test.ts` **7 passed** |

### 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| C1 迁移 0055 挑战表落地（双方言成对 DDL）+ 黄金断言全套同步 | **达成**（实现事实；检查点标记仍待编排器更新） | 成对 DDL + ApplyPostgres；冻结 checksum 与独立复算一致；head/lockedHead/completeLostLedger/postV1/夹具/migrate×2/operations/restart 均已改；PG 全 catalog bootstrap PASS |
| C2 绑定/校验/重发流实现且合同映射可核对（占槽冲突 fail-closed、覆写换绑、三态、常量时间比较、尝试上限） | **部分达成** | 自助 Bind/Verify/Resend 与 §2/§4/§5/§6 可核对（测试绿）。I-006 HTTP 面未接通（F-001）；超限作废码与 D-001 字面不一致（F-004）；Bind 再签发不受冷却（F-002） |
| C3 身份面 API + 账号页最小绑定 UI 可用；测试绿（sqlite 全量 + PG 集成） | **部分达成** | 自助三端点 + 最小卡 + sqlite/handler 本会话绿 + PG bootstrap PASS。I-006 声明的 users PATCH HTTP 面实际丢弃 `email`（F-001）。无 `migrate_0055` 语义专项（N-6） |
| C4 independent 审计落盘且开放 required = 0 | **本条落盘后意见侧条件未满足**（开放 required = 1） | 本文件；不代改 `00-meta` / goal-tree |
| GOAL-002 D-001 §2 绑定占槽 | **达成**（自助 Bind） | 他号同址（大小写折叠）`EMAIL_TAKEN`；多 NULL 不受影响（R2 物理槽仍在） |
| §4 6 位码 + MailSender + 不锁 CaptureSink | **达成** | `generateEmailCode` `%06d` + `crypto/rand`；只消费 `kernel.MailSender`；无 `Last()` |
| §5 换绑覆写释放旧址 | **达成**（服务层） | `TestRebindOverwriteReleasesOldSlot`。最小页 verified 态隐藏绑定表单，换绑 UI 未暴露（见 F-002 附带） |
| §6 三态；I-006 代填不得直达 verified | **负向达成 / 正向缺口** | 仓储无 `email_status=verified` 写入面（仅 Verify 成功路径）。HTTP 代填正向路径不可达（F-001） |
| GOAL-003 A-001 F-002 配对不变量进仓储 | **未关闭** | 写路径成对；无拒绝非法对（N-2） |
| GOAL-003 A-001 N-1 ASCII `lower()` 补偿 | **未关闭** | 唯一性仍只走 SQL `lower()`（N-1） |
| I-005 数值冻结并落地 | **达成（Resend 路径）** | TTL 10 分钟 / 冷却 60 秒有测试。Bind 再签发不受冷却（F-002，非 D-001 字面违约） |
| I-006 允许代填→pending | **未在 HTTP 合同面达成** | 仓储测试绿；users 资源工厂未列入 `email`（F-001） |

`00-meta` 成功标准第 1 条失败语义写 `email_already_in_use`，冻结方案与 Catalog 均为 `EMAIL_TAKEN`（N-4）。以 D-001 为准。

### Findings

#### F-001 · I-006 用户 PATCH HTTP 面未接入 `email`（代填合同面不可达）

- 严重度：high
- 建议：required
- 状态：open
- 描述：GOAL-004 D-001 §3 冻结「既有用户 PATCH 面新增可选 `email` 字段 → 置 pending；清空 = 回 unbound」。仓储 `UpdateUser` + `UserPatch.Email` 实现了该语义，且 `TestAdminPrefillPendingCannotVerifyDirectly` 覆盖。但 HTTP 资源工厂 `usersResource` 的 `PatchFields = ["name"]`、`JSONFields = ["roles"]`，`decodeResourcePatch` **忽略未知键**，因此 `PATCH /api/users/{id}` 的 `email` 从未进入 `usersEntity.Update`。`Update` 内 `body["email"]` + `v.(json.RawMessage)` 在工厂路径上是死代码；即便把 `email` 补进 `JSONFields`，工厂写入的是已反序列化的 `any`（string），`json.RawMessage` 断言会 panic。清空回 unbound 也不能走 `PatchFields`（空串会被 `patchFieldError` 拒绝），正确槽位是允许空串的 JSON 字段，并按 string 解码。无任何 handler 测试 PATCH `email`。E-002「users PATCH」主张在 HTTP 面上名不副实。负向门禁（无路径直达 verified）仍然成立。
- 证据：`apps/api/internal/handler/users.go:67-69, 215-223`；`apps/api/internal/handler/resources.go` `decodeResourcePatch`（未知键忽略；`JSONFields` 存 `any`）；`apps/api/internal/modules/authsession/repository.go` `UserPatch.Email`；`email_identity_test.go` `TestAdminPrefillPendingCannotVerifyDirectly`；handler 测试树无 `"email"` 用例。关联 I-006。

#### F-002 · pending 同址再次 Bind 不受 I-005 60s 冷却

- 严重度：med
- 建议：recommended
- 状态：open
- 描述：D-001 把冷却写在 **Resend**（`TestResendCooldown` 可核对）。`BindEmail` 对已 pending 的同址仍生成新码并发信（仅 verified 同址幂等不重发）。最小卡在 `status !== "verified"` 时仍渲染绑定表单；同会话 `address` 非空时「发送验证码」可再调 Bind，绕过 Resend 冷却。不把本条升为 required：冻结方案未要求 Bind 冷却。编排器若要把 I-005 控制目标理解为「一切再签发」，应改 Bind（pending 同址走冷却或拒绝）或 pending 态隐藏绑定提交。
- 证据：`email_identity.go` `BindEmail` vs `ResendEmailCode`；`email-identity.tsx:151-165`；D-001 §2。关联 I-005。

#### F-003 · 治理台账未同步 GOAL-004 / E-002

- 严重度：med
- 建议：recommended
- 状态：open
- 描述：GOAL-004 五件套与 E-002 已存在，且本审计针对该目标，但 `goal-tree.md` 树与状态表仍只有 GOAL-001 / GOAL-002 / GOAL-003；`workspace.md` 与 Root `00-meta` 纲领表 R3 仍为「待启动」。`02-execution.md` 索引只列 E-001，E-002 文件已在同一 checkpoint。Root `00-meta` I-005 / I-006 仍 collecting，与 Root `01-decision.md` / GOAL-004 D-001 的 verified 不一致（`01-decision.md` 自注「状态以 00-meta 为准」）。不否定代码事实，但 `/govern` 响应必须补登记，否则违反 AGENTS §7。本条独立审**禁止**改 goal-tree / 检查点。
- 证据：`docs/workspaces/workspace-018-account-email-identity/goal-tree.md`；`workspace.md` 纲领表；GOAL-001 `00-meta.md` vs `01-decision.md` 信息表；GOAL-004 `02-execution.md` vs `02-execution/E-002-r3-implementation.md`。

#### F-004 · 5 次失败作废返回 `EMAIL_CODE_EXPIRED`，D-001 写 `EMAIL_CODE_INVALID`

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：D-001 §2 Verify：「不匹配 → attempt_count++，≥5 作废挑战 `EMAIL_CODE_INVALID`」。实现与测试把第 5 次映射为 `ErrEmailCodeExpired`（挑战已删，须重发）。功能上作废成立；错误码分类与冻结方案字面不符。后续 Verify 无挑战走 `EMAIL_NOT_PENDING`。
- 证据：`email_identity.go:215-218`；`email_identity_test.go` `TestFailedAttemptsVoidChallenge`；D-001 §2。

#### N-1 · SQLite `lower()` ASCII 残留未做 locale-无关唯一性补偿

- 严重度：low（note）
- 建议：recommended
- 状态：open（GOAL-003 A-001 N-1 移交项，R3 未关）
- 描述：唯一性判定仍是 SQL `lower(email) = lower(?)`。应用层 `EqualFold` 只用于「同号已 verified」幂等，不参与他号占槽。非 ASCII 大小写对在 SQLite 上仍可能双占槽。D-001 N-1 写「比较与唯一性判定统一走 SQL lower()」——与「locale-无关折叠补偿」并不等同。本条**不**视为已关闭，也**不**升为 required（GOAL-003 已标 recommended 移交）。
- 证据：`email_identity.go:160-162`；`users_repository.go:214-216`；GOAL-003 A-001 N-1；`migration.go` 0054 注释。

#### N-2 · email / email_status 配对不是仓储拒绝不变量

- 严重度：low（note）
- 建议：recommended
- 状态：open（GOAL-003 A-001 F-002 移交项）
- 描述：Bind / 代填 / 清空写路径都成对（非空+pending，或双 NULL）。没有「拒绝 (email, NULL) / (NULL, pending)」的显式校验或测试。表级 CHECK 仍允许非法对（R2 已知）。
- 证据：`email_identity.go` / `users_repository.go` 写语句；无配对拒绝测试；GOAL-003 A-001 F-002。

#### N-3 · users 资源行不投影 email；管理 schema 无代填控件

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：`userToMap` 无 `email` / `emailStatus`；`users` schema 无 email 字段。即便 F-001 修复，PATCH 成功响应与列表/详情仍看不到代填结果。I-006 当前只是仓储能力，不是管理 UI。
- 证据：`handler/users.go` `userToMap`；`apps/api/internal/modules/users/schema` 无 email。

#### N-4 · 成功标准失败码文案与 Catalog 不一致

- 严重度：low（note）
- 建议：recommended
- 状态：open
- 描述：`00-meta` 成功标准 1 写 `email_already_in_use`；D-001 / Catalog / handler 均为 `EMAIL_TAKEN`。以冻结方案为准。
- 证据：GOAL-004 `00-meta.md`；`errorcatalog.go` `EMAIL_TAKEN`。

#### N-5 · E-002「12 键」与实际 15 键

- 严重度：low（note）
- 建议：recommended
- 状态：open
- 描述：zh/en `schema.account.email.*` 各 15 键且对齐；E-002 写 12。不影响正确性。
- 证据：`apps/web/src/i18n/messages/{zh-CN,en-US}.json`；E-002。

#### N-6 · 0055 无双方言语义专项；PG 靠全 catalog bootstrap

- 严重度：low（note）
- 建议：recommended（可选，GOAL-003 F-001 口径）
- 状态：open
- 描述：无 `migrate_0055_test.go` 对 INTEGER/BIGINT / PK / CASCADE 做专项。本会话 `TestFullCatalogPostgresBootstrapIntegration` PASS，说明 PG 应用 0055 成功，不等于列类型专项。GOAL-003 F-001 为 optional 移交。
- 证据：无 0055 专项测试文件；本会话 PG bootstrap verbose PASS。

### 必改项汇总

1. **F-001（required / high）**：把 `email` 接入 users PATCH HTTP 面（建议 `JSONFields`，允许 `""` 清空），按 string/`any` 解码而不是 `json.RawMessage`；补 HTTP 测试：非空 → pending、空串 → unbound、冲突 → `EMAIL_TAKEN`、响应或 GET 可核对。在此之前 **不得**把 I-006 / 成功标准 3 的正向路径写成已交付，也 **不得**将 GOAL-004 标 `done`。

开放 **required** = **1**。

### 与既有意见的异同

本目标此前无 self / independent 条目。

对照 GOAL-003 A-001 independent `pass`：F-001（PG 语义 harness）本条确认为仍可选未做（N-6）；F-002 配对与 N-1 ASCII `lower()` **未**在 R3 关闭，维持 recommended，不假装闭合。GOAL-003 F-003（goal-tree / 执行索引）同类问题在本目标再现（本条 F-003）。

对照 GOAL-002 A-001 self：§2/§4/§5/§6 自助流可核对；I-006 合同面未接通。

### 结论 + 建议给编排器/用户的下一步

**conditional** —— R3 的迁移 0055、自助 Bind/Verify/Resend、MailSender 消费、错误契约七码、最小账号卡与 sqlite/PG bootstrap 本会话可核对；**不能**无条件关门，因为 I-006 / D-001 §3 声明的管理员 PATCH 代填 HTTP 面实际不可达（F-001 required）。

建议 `/govern`：

1. **先响应 F-001**：修复 users 资源字段列表与解码，补 HTTP 测试；闭合路径 `fixed`。不要用仓储单测代替 HTTP 合同。
2. F-002 / F-003 / F-004 / N-1～N-6 按 recommended/note 处理。F-003 建议在响应事务中 **fixed**（登记 goal-tree、执行索引 E-002、同步 Root `00-meta` I-005/I-006 与 R3 指针）。N-1 / N-2 若本波不修，须书面留下残余范围与复审触发（`accepted-residual`），不得默认为已兑现承接清单。
3. 不要由本审计代改 `status` / 检查点 / `progress`。C1 实现侧可标完成；C2/C3 待 F-001；C4 待开放 required = 0。
4. 用户若书面驳回 I-006 HTTP 面（认为仓储即可，或把代填移出分母），走 `user-overruled` 并改 D-001 §3 / 成功标准 3 后再审。

### 声明

本意见 `source: independent`，**不修改** `status` / 检查点 / 派生 `progress` / goal-tree / 方案正文 / 产品代码。响应、finding 闭合与关门状态变更由 **`/govern`** 与用户书面裁决处理。
