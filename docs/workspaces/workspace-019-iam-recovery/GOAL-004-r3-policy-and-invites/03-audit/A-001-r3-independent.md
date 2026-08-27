---
id: A-001
doc: audit-entry
goal: GOAL-004-r3-policy-and-invites
source: independent
status: recorded
created: 2026-08-25
updated: 2026-08-25
version: 1.0.0
scope: R3 实施切片（C2 迁移 0057/0057b/0058 至 catalog v59 · C3 策略四口+邀请域/HTTP · C4 Web 四面）；execution-facts 对照本目标 D-001
verdict: conditional
auditor: grok-build (grok-4.6 · reasoning high)
---

# A-001 · R3 密码策略 + 邀请入职实施切片独立审计（independent · 2026-08-25）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型** / **scope**：execution-facts（对照本目标 D-001）· C2 迁移 0057 `password_policy` / 0057b=`user_password_history`（catalog v58）/ 0058=`user_invites`（catalog v59）与黄金断言；C3 `authsession/password_policy.go` 四口强制 + 历史捕获、`authsession/invites.go` 邀请域、`handler/invites.go` 管理四路 + 公开 accept、`handler/password_policy_settings.go`；C4 `password-policy-tab` / `user-invites-panel` / `invite-accept` / 新建用户表单角色字段
- **verdict**：**conditional**
- **开放 required**：1（F-001）

### 范围与区间

| 项 | 值 |
|----|-----|
| 工作区 | `workspace-019-iam-recovery`（`workspace.md`：`id` 匹配；`root_goal` = `GOAL-001-iam-recovery`；`canonical_scope` = 本目录；`shared_materials_catalog: none`；`primary_plan` = `VP-019-iam-recovery`） |
| 被审目标 | `GOAL-004-r3-policy-and-invites`（同区；`parent` = `GOAL-001-iam-recovery`） |
| audit_type | execution-facts（对照 D-001） |
| 对照 | GOAL-004 `00-meta` 成功标准 1～3 与 C2–C4；本目标 D-001 §1～§5；GOAL-002 D-001 §2/§3；用户裁决「受邀角色以发布邀请时指定为准」 |
| 信息门禁 | GOAL-004 未自建 I-00N。合同输入 I-003 / I-004 / I-005 / I-007 在 GOAL-002 已 **verified**（2026-08-25）。无到期且影响本切片的 required 未闭项 |
| 共享资料 | 无（`none`；未把资料目录当事实或关闭证据） |
| 代码基准 | `HEAD` `ff581731b62c9410626a753817feb81cb150cb64`（工作树干净） |
| 本审计未改 | 目标 `status` / 检查点 / 派生 `progress` / goal-tree / 方案正文 / 产品代码或测试 |

未读取或比较其他工作区上下文。本目标此前无 self / independent 条目（`03-audit.md` 仅占位行）。项目级 `docs/architecture/independent-audit-execution.md` 写「先自审再 independent」；AGENTS P-003 写 `independent` 不固定要求 self。本条按用户指定独立执行，不把缺 self 标为 required。

### 成果（有证据）

| 主张 | 证据 |
|------|------|
| 工作区绑定合格 | `workspace.md` id / root_goal / canonical_scope / `plan_refs`+`primary_plan` 与本目标路径一致；资料目录 `none` |
| C2 · 0057 单例策略行 + seed（可移植） | `migration.go` `passwordPolicyDDL`：`id INTEGER PK CHECK(id=1)`、`min_length`/`min_categories`/`history_depth` DEFAULT 8/0/0、`INSERT (id) VALUES (1)`；Descriptors v57、`ApplyPostgres` nil |
| C2 · 0057b 历史表双方言 | `userPasswordHistoryDDL` / `userPasswordHistoryPGDDL`：SQLite `created_at INTEGER` / PG `BIGINT`；catalog v58；checksum salt `0058:user-password-history:v1` |
| C2 · 0058 邀请表双方言 | `userInvitesDDL` / `userInvitesPGDDL`：`token_hash UNIQUE`、`roles` NOT NULL、`email` 可空、`consumed_at`/`revoked_at` 可空；catalog v59；salt `0059:user-invites:v1` |
| C2 · 冻结 checksum 与 DDL 实算一致 | 台账 want = `bfc7c4f2…` / `04f77fdc…` / `a35bbb21…`（`store/migrate_test.go` L710–712）。本会话按 `kernel.MigrationChecksum`（normalizeSQL + 上列 transform id）独立 sha256 复算 **三处 MATCH** |
| C2 · 黄金断言 catalog head = 59 | `identity.go` `completeFingerprintCatalogHead = 59`；`completeLostLedgerTables` / `postV1CatalogTables` 含三表；`identity_test.go` `lockedHeadExtraTables[57..59]`；`migrate_test.go` / `operations_test.go` / `restart_test.go` `len==59` 尾断言 `user_invites` |
| C3 · 四口接线存在 | `handler/users.go` Create / Update(password)；`handler/account_self.go` changePassword；`handler/recovery.go` complete；均在设密前 `ValidateNewPassword`，错误 `INVALID_PASSWORD`。邀请激活另走 `ValidateNewPassword("", …)`（D-001 §3） |
| C3 · 类别与历史强制 + 捕获 | `countCategories` 四类；`passwordInHistory` bcrypt 比较最近 N 条；`UpdateUser` 同事务读 `history_depth` 后 `capturePasswordHistory`。域测试 `TestValidateNewPasswordBaselineCategoriesHistory` 覆盖默认基线、`min_categories=3`、旋转后复用旧码 |
| C3 · I-003 默认 = 现行 8/关/关 | seed 列 DEFAULT；`GetPasswordPolicy` 缺行退回 `{MinLength:8}`。I-007：无存量扫描、无登录阻断路径 |
| C3 · 邀请域生命周期 + 角色裁决 | `CreateInvite` 角色 ≥1 且 `roles` 表存在；`AcceptInvite` 单事务 live 校验 → 角色 re-validate → 用户名冲突 fail-closed → 建号 `Roles=inv.Roles`、`MustChangePassword` 走列 DEFAULT 0 → 消费邀请、不签发会话。`TestCreateAndAcceptInviteAssignsIssuedRoles` 断言角色=邀请指定且 replay → `ErrInviteInvalid`。重发 60s 冷却 + 旧 token 失效 + 撤销即时失效：`TestInviteResendCooldownRotationAndRevoke` |
| C3 · HTTP 管理四路 + 公开激活 | `InviteAdminRoutes`：POST/GET `/api/users/invites`、DELETE `{id}`、POST `{id}/resend`；`users.invite` 经 `PermissionsForUser`（PolicyAdmin → 仅 admin）。`RegisterInviteAccept` 中央 `POST /api/auth/invite/accept` → 204 无 token。带邮箱创建失败补偿 `RevokeInvite` + `EMAIL_SEND_FAILED` |
| C3 · 配置 API | GET/PATCH `/api/settings/password-policy`；PATCH 范围 8–72 / 0–4 / 0–10。`kernel/profile.go` admin.users +4 路由 + `users.invite`；admin.settings +2 路由。错误目录 + `error_contract_test.go` 冻结集含 `INVALID_INVITE_BODY` / `INVITE_INVALID` / `INVITE_ROLE_GONE`。composition mvp 权限 11 / admin 33 |
| C4 · Web 四面存在 | 设置 tab `password-policy-tab`；用户页 `user-invites-panel`（角色随邀请、邮箱可选、一次性披露 link）；`main.tsx` 未认证分支 `/invite/accept` → `InviteAcceptPage`（成功回登录）；`users.json` `create-user-form` 补 `roles` checkboxGroup + `optionsSource /api/roles`。zh/en 各 31 个 `invite.*` 键；W25 `custom-components.schema.test.ts` side-effect import 两组件 |
| 本会话复跑 | `apps/api`：`go test ./internal/modules/authsession/ ./internal/store/ ./internal/handler/ ./internal/composition/ ./internal/modules/users/ ./internal/modules/settings/ -count=1` **exit 0**（authsession 7.0s / store 38.9s / handler 39.8s / composition 17.4s / users 1.3s / settings 1.2s）。`TestFullCatalogPostgresBootstrapIntegration` 本会话 **未跑**。web vitest / tsc **未复跑** |

### 对照成功标准（若适用）

| 标准 | 状态 | 证据 |
|------|------|------|
| 成功标准 1：策略可配置且四口统一强制；minLength/categories/history 生效可核对；默认保持现行行为 | **部分达成** | 配置面可 PATCH/UI 读写三字段。categories / history 域测试可核对。**minLength 只被写入、不被 `ValidateNewPassword` 读取**（F-001）。默认 8/0/0 与 I-003 一致 |
| 成功标准 2：邀请全链可核对（指定角色建邀 → 激活带角色 → 一次性/撤销/过期/重发冷却） | **部分达成** | 域层：角色随邀请、一次性、冷却、撤销、创建期未知角色。HTTP 管理/激活与邮件投递 **无 handler 测试**。过期接受与「签发后删角色再激活」无用例（实现有 `live()` / accept 期 re-validate）。见 F-002 |
| 成功标准 3：independent 开放 required = 0；Root progress → 3/4 | **本条落盘后意见侧未满足** | 开放 required = 1（F-001）。Root / 检查点不由本审计改写 |
| D-001 §1 迁移 + 黄金断言至 v59 | **达成** | 三表、双方言（策略可移植）、checksum 复算一致、head=59 |
| D-001 §2 四口 + 配置 API | **部分达成** | 四口已接线；类别/历史生效；配置范围校验在。minLength 强制缺（F-001）。GET 走 `settings.write` 而非与邮件 tab 同形的 `settings.read`（F-004） |
| D-001 §3 邀请 + 用户裁决角色 | **部分达成** | 域与接线符合角色-at-issuance。API `link` 为绝对 URL（方案写相对路径）；邮件正文字面 `` `n `` 而非换行；resend 忽略发信失败（F-003） |
| D-001 §4 Web 四面 | **达成**（存在性；无组件测试） | 四块文件 + schema 挂载 + i18n + AuthGate。本会话未复跑 vitest |
| GOAL-002 D-001 §2/§3（I-003/004/005/007） | **部分达成** | 起步宽松默认、渐进四口、双形态、7 天默认、一次性、撤销、60s 重发均有实现。§2「可配置最小长度过当时生效策略」被 F-001 否定其字面 |
| 用户裁决：受邀角色以发布时指定为准 | **达成**（域） | `AcceptInvite` 使用 `inv.Roles`；测试断言 `[admin]`；未选固定 viewer |

`00-meta` 已将 C2–C4 标完成、frontmatter `progress: 4/5`（正文仍写「当前 0/5」；goal-tree 为 3/5）。本意见 **不批准或改写** 这些检查点；C3「minLength 强制」完成声明被 F-001 打折。

### Findings

#### F-001 · `ValidateNewPassword` 不读取 `policy.MinLength`，收紧最小长度在四口与邀请激活上是空操作

- 严重度：high
- 建议：required
- 状态：open
- 描述：成功标准 1 与 GOAL-002 D-001 §2 要求可配置 **最小长度** 在四个设密时刻（及 D-001 §3 邀请激活）过当时生效策略。`UpdatePasswordPolicy` / PATCH / 设置 tab 均可把 `min_length` 写到 8–72，但 `ValidateNewPassword` 只硬编码 `policyMinLengthFloor/Ceiling`（8–72），读出策略行后只判断 `MinCategories` 与 `HistoryDepth`，**从不比较 `len(plain)` 与 `policy.MinLength`**。因此管理员把最小长度从 8 收到 12 后，8～11 字节密码在 users Create / users Update(password) / account_self changePassword / recovery complete / invite accept 上仍然通过（只要过基线 8 与其它开着的旋钮）。域测试同样只把 `MinLength` 维持在 8，没有「提高到 12 后拒绝 8 字节」的用例。D-001 §2 算法子弹写成「基线 8–72 → categories → history」，与同一条款的 PATCH `minLength∈[8,72]` 和成功标准 1 互相矛盾——实现跟了算法字面、没跟可配置项。
- 证据：`apps/api/internal/modules/authsession/password_policy.go` L79–100（长度门在 L80–82，随后 `GetPasswordPolicy` 不用 `MinLength`）；`handler/password_policy_settings.go` L109–129 持久化 `minLength`；`invites_test.go` `TestValidateNewPasswordBaselineCategoriesHistory` 无抬高 minLength 用例；四口调用点 `handler/users.go` L167/L203、`account_self.go` L260、`recovery.go` L208、`invites.go` L296。
- 建议修复：在基线 8–72 之后增加 `length < policy.MinLength` → `ErrPasswordPolicyViolation`；补域测试（minLength=12 拒绝 8 字节、接受 12 字节）；至少一条 handler 或域路径证明四口之一吃到该拒绝。

#### F-002 · 邀请 HTTP 面、过期接受、签发后角色消失缺少可核对测试

- 严重度：med
- 建议：recommended
- 状态：open
- 描述：C3/C4 的文件与域测试能证明「邀请域状态机 + Web 组件已挂上」。成功标准 2 的「全链可核对」还包含管理四路权限/错误码、公开 accept 204 不签发、邮件双形态、过期统一 `INVITE_INVALID`、激活前角色被删 → `INVITE_ROLE_GONE`。现状：`handler/` 无 `invites_test.go` / `password_policy_settings_test.go`；Web 无 `invite-accept` / `user-invites-panel` / `password-policy-tab` 测试。`TestAcceptInviteRejectsUnknownTokenUsernameConflictAndRoleGone` 注释写「Role deletion between issue and acceptance」，实际只断言 **创建时** 未知角色失败，从未走到 accept 期 re-validate。过期接受只靠 `Invite.live()` 代码阅读。不构成 C3/C4 造假，但关门前成功标准 2 不能称为 HTTP+过期+角色消失全链可核对。
- 证据：`apps/api/internal/handler/` 无邀请/策略测试文件；`invites_test.go` L54–81；`invites.go` `live()` L69–71 与 Accept L297–310；本会话 handler 包绿是既有套件，不覆盖新路由。

#### F-003 · 邀请链接形状与邮件正文偏离 D-001；resend 发信失败静默

- 严重度：med
- 建议：recommended
- 状态：open
- 描述：D-001 背景与 §3 冻结「返回 **相对** 链接；邮件正文以请求 Host 推导绝对基址」。实现 `inviteLink` 对 API 响应和邮件使用同一绝对 URL（默认 scheme `http`，仅当存在 `X-Forwarded-Proto` 才改）。功能上复制链接仍可用，但 TLS 反代未转发该头时邮件会指向 `http://`。邀请信 `TextBody` 在 **双引号** 字符串里写了字面 `` `n ``，Go 不会把它当换行，收件人看到的是挤在一行的 `` `nhttps://…`n`n ``。创建路径发信失败会补偿撤销并返回 `EMAIL_SEND_FAILED`；**resend 在已轮换 token 之后 `_ = sendInviteMail(...)` 忽略错误**——旧链接已死，新信可能没发出，HTTP 仍 200 并披露新 link（管理员若只依赖邮箱会丢链）。
- 证据：`handler/invites.go` `inviteLink` L176–182、`create` L152–160、`resend` L253–256、`sendInviteMail` L190–193。

#### F-004 · 策略 GET 权限与历史语义的次要偏差

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：（1）D-001 §2「复用既有 settings 管理权限」。同页邮件 tab：GET `settings.read` / 写 `settings.write`（`mail_admin.go`）。密码策略 GET **与 PATCH 都**要求 `settings.write`（`password_policy_settings.go` L38–39）。仅 `settings.read` 的主体能开设置页但策略 tab 加载失败。（2）历史比较只打 `user_password_history`，不含 `users.password_hash`；域测试显式断言「旋转前当前密码必须通过」。自愿改密 handler 另有 `new != current`；管理员重置与自助恢复可以把新密码设成 **当前** 密码，history_depth 挡不住。D-001 算法字面确实只说「最近 N 条历史」，故不升 required；若产品意图是「含当前的最近 N 个密码」，应补 live-hash 比较并回写方案。（3）`capturePasswordHistory` 插入/裁剪错误直接 `return` / `_, _ =`，与「同事务推入历史」的 fail-closed 字面不符（本会话 sqlite 同表 `DELETE … LIMIT` 裁剪探测成功，属吞错而非 trim SQL 不可用）。
- 证据：`password_policy_settings.go` L38–68；`mail_admin.go` L34–45；`invites_test.go` L154–156；`password_policy.go` L135–147；`users_repository.go` L240–248。

### 必改项汇总

| ID | 严重度 | 摘要 |
|----|--------|------|
| F-001 | high | `ValidateNewPassword` 必须强制 `policy.MinLength`，并补可核对测试（抬高最小长度后四口之一拒绝过短密码） |

F-002～F-004 为 recommended，不单独把本切片打成 fail，但 F-001 未闭合前不得把成功标准 1 或 C3「策略强制」视为已满足。F-002 建议在 C5 关门前至少补：accept 过期、accept 期角色消失、一条公开 accept HTTP 204/不签发断言。

### 与既有意见的异同（若有 self/independent 历史）

本目标此前无 A 条目。不与 GOAL-002 A-001（R1 合同 self `pass`）或 GOAL-003 A-001/A-002 冲突。E-002/E-003/E-004 将 C2–C4 记为完成：C2 与本审计一致；C3 四口接线与邀请域属实，**minLength 强制完成声明被 F-001 否定**；C4 文件与挂载可核对，成功标准 2 的 HTTP/过期/角色消失证据弱于检查点措辞。

注释（不升 finding）：D-001 §1 把历史表写在「0057」、邀请表写「0058」；catalog 实际为 v57 策略 / v58 历史 / v59 邀请（用户本轮 scope 已按 0057/0057b/0058→v59 陈述，checksum salt 与 version 一致）。`00-meta` frontmatter `4/5`、正文「当前 0/5」、goal-tree `3/5` 三处不一致，由 `/govern` 处理。邀请面板角色是逗号分隔文本而非 checkboxGroup（D-001 只强制新建用户表单用 checkboxGroup，不升 finding）。`/api/auth/invite/accept` 未进 operational allowlist（维护模式会 503）；恢复端点有明确合同要求，邀请没有。

### 结论 + 建议给编排器/用户的下一步

C2 迁移与黄金断言本会话复算/复跑通过。C3 邀请域与四口接线是真的；**可配置最小长度未进入强制函数**，成功标准 1 不能签字。建议 `/govern`：先修 F-001 并留测试，再决定 F-002/F-003 是本波补还是 residual；闭合 required 之前不要把 C3/C5 当完成，也不要把 Root 推到 3/4。

### 声明

本意见不修改 status/progress；响应由 /govern 处理。

---

## 编排器响应（2026-08-25 · /govern · git 2f088d55）

> 本节为编排器响应留痕，不改动上方独立意见正文。

| Finding | 处置 | 证据 |
|---------|------|------|
| F-001（required · high） | **fixed** | ValidateNewPassword 先读策略行，采用配置 MinLength（下限钳制 8）再校验；新增 TestValidateNewPasswordConfiguredMinLengthBites（minLength=12 时 8 字节拒绝/12+ 通过）；authsession 包全绿 |
| F-002（recommended） | **fixed（域层+HTTP 最小集补齐）**：①域层补 accept 过期与「签发后角色删除」两组测试；②公开 accept 面新增 TestInviteAcceptPublicSurface（204 不签发 + INVALID_PASSWORD + INVITE_INVALID 统一 + USERNAME_TAKEN 分明）。管理四路的完整 HTTP 表格化测试**接受为残余移交 R4 端到端证据阶段**（域逻辑已全测，HTTP 仅薄映射层），R4 e2e 将经真实 HTTP 覆盖 |
| F-003（recommended） | **fixed**：邀请信正文改真实换行（原 raw-string 字面 `n 缺陷确认并修复）；resend 发信失败显式 502 EMAIL_SEND_FAILED（邀请保持 live，冷却后可重试）；绝对链接口径回写 D-001 §3 |
| F-004（recommended） | **fixed**：策略 GET 分离 settings.read / PATCH settings.write（与邮件 tab 同形）；「历史不含轮换前当前密码」「捕获 best-effort」两处口径回写 D-001 §2 |

注释处置：00-meta 正文进度句已同步 4/5；`/api/auth/invite/accept` 已入 operational allowlist（含测试）；goal-tree 以状态表为准。

**开放 required = 0**（F-001 fixed）。