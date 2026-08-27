---
id: A-001
doc: audit-entry
goal: GOAL-005-r4-evidence-closeout
source: independent
status: recorded
created: 2026-08-25
updated: 2026-08-25
version: 1.0.0
scope: R4 端到端证据（close-out 预审）——TestR4RecoveryChainOverHTTP / TestR4InviteChainOverHTTP / TestR4PolicyEnforcementOverHTTP；Root 成功标准 1–6 对照；无越界核对（SMS/模板/多邮箱/业务域/OIDC）
verdict: conditional
auditor: grok-build (grok-4.6 · reasoning high)
---

# A-001 · R4 端到端证据独立审计（independent · 2026-08-25）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型** / **scope**：close-out 预审 · C1/C2 三条 HTTP e2e 测试 + Root 成功标准 1–6 对照 + 无越界核对；不审 Web 组件、不改 Root 关门状态
- **verdict**：**conditional**
- **开放 required**：1（F-001）

### 范围与区间

| 项 | 值 |
|----|-----|
| 工作区 | `workspace-019-iam-recovery`（`workspace.md`：`id` 匹配；`root_goal` = `GOAL-001-iam-recovery`；`canonical_scope` = 本目录；`shared_materials_catalog: none`；`primary_plan` = `VP-019-iam-recovery`） |
| 被审目标 | `GOAL-005-r4-evidence-closeout`（同区；`parent` = `GOAL-001-iam-recovery`） |
| audit_type | close-out（预审；C3 的 independent 半边） |
| 对照 | GOAL-005 C1/C2；Root `GOAL-001` 成功标准 1–6；GOAL-004 A-001 F-002 残余移交 |
| 信息门禁 | GOAL-005 未自建 I-00N。Root I-001～I-009 均已 `verified`/`registered`（I-008 non-blocking，最晚 R4，已在 GOAL-002 D-001 关闭）。无到期且影响本切片的 required 未闭项 |
| 共享资料 | 无（`none`；未把资料目录当事实或关闭证据） |
| 代码基准 | `HEAD` `e6a376e340aeba84c24843557965ef4a2ef8c04a`（工作树干净） |
| 本审计未改 | 目标 `status` / 检查点 / 派生 `progress` / goal-tree / 方案正文 / 产品代码或测试 |

未读取或比较其他工作区上下文。本目标此前无 self / independent 条目（`03-audit.md` 仅占位行）。项目级 `docs/architecture/independent-audit-execution.md` 写「先自审再 independent」；AGENTS P-003 写 `independent` 不固定要求 self。本条按用户指定独立执行，不把缺 self 标为 required。

上文放行的 prior 意见已核：

| 意见 | 本会话核对 |
|------|------------|
| GOAL-003 A-001（R2 independent · conditional） | 编排器响应节：F-001 required **fixed**（`ddd20500`）；F-002～F-004 **fixed**。索引标归零。A-002 self `pass` |
| GOAL-004 A-001（R3 independent · conditional） | 编排器响应节：F-001 required **fixed**（MinLength 进入强制函数）；F-002 **fixed（域层+HTTP 最小集）** 且「管理四路完整 HTTP 表格化测试**接受为残余移交 R4**」；F-003/F-004 **fixed**。A-002 self `pass`。GOAL-004 `03-audit.md` **索引行仍写开放 required = 1**（与响应节 0 不一致，见注释） |

### 成果（有证据）

| 主张 | 证据 |
|------|------|
| 工作区绑定合格 | `workspace.md` id / root_goal / canonical_scope / `plan_refs`+`primary_plan` 与本目标路径一致；资料目录 `none` |
| 测试夹具挂真实 mock 渠道 + 中央面 | `handler/testhelpers_test.go`：`recoverySender := mail.NewOutboxSink(st, 0)` 同时交给 `EmailIdentityRoutes` / `RegisterRecovery` / `InviteAdminRoutes`；另挂 `RegisterInviteAccept` + `PasswordPolicyRoutes`。码/链接从 `OutboxSink.List`+`Get` 取，不是测试桩 |
| C1 · 恢复链 HTTP e2e | `TestR4RecoveryChainOverHTTP`：登录 → `POST /api/account/email/bind` → 出站取 6 位码（`mustContain="邮箱验证码"`）→ verify → 公开 `POST /api/auth/recovery/start` 202 → 出站取码（`mustContain="password reset code"`）→ complete 204 → `POST /api/auth/login` 新密码拿 `accessToken` |
| C2 · 邀请链 HTTP e2e | `TestR4InviteChainOverHTTP`：admin `POST /api/users/invites` 201（`users.invite` 已入 `testsupport` seed）→ 出站取 `token=` → `POST /api/auth/invite/accept` 204 → 受邀者 `viewer` 登录 → 回放同一 token → `INVITE_INVALID` |
| C2 · 策略强化经 HTTP 拦住弱创建 | `TestR4PolicyEnforcementOverHTTP`：`PATCH /api/settings/password-policy` `minLength=12` 200 → 8 字节 `"eightbte"` 创建体含 `INVALID_PASSWORD` → `"strong-pass-123"` 创建 200/201。顺带：PATCH 后原 admin Bearer 仍可用于创建（策略变更未强制登出该会话） |
| 本会话复跑 | `apps/api`：`go test ./internal/handler/ -count=1 -timeout 120s -v -run "TestR4RecoveryChainOverHTTP\|TestR4InviteChainOverHTTP\|TestR4PolicyEnforcementOverHTTP"` **exit 0**（Recovery 0.19s / Invite 0.22s / Policy 0.12s；包 2.742s） |
| 无越界（本波 IAM 面） | `handler/recovery.go` / `invites.go` / `password_policy_settings.go` / `r4_evidence_test.go` 无 SMS/OIDC/SCIM/模板中心/多邮箱。投递只经 `kernel.MailSender`（mock `OutboxSink`）。Charter 仍 `schema-ui-core-admin-foundation@0.2.0`、`primary_workspace` 仍 workspace-001。`kernel/profile.go` `ProfileMVP` **不含** `admin.settings`（策略 UI 未进 mvp 默认集）。管理员重置仍写 `MustChangePassword=true`（`handler/users.go` Update password 分支） |
| I-00N | 无到期未闭 required。I-008 已 verified，本阶段无新信息项 |

### 对照成功标准（若适用）

| 标准 | 状态 | 证据 |
|------|------|------|
| GOAL-005 C1：bind/verify → start → 渠道取码 → complete → 新密码登录 | **达成**（本会话复跑 PASS） | `TestR4RecoveryChainOverHTTP` |
| GOAL-005 C2：admin 建邀 → 渠道取链接 → accept → 受邀角色登录 + 一次性回放；策略强化拦住弱创建 | **达成**（本会话复跑 PASS；C2 字面不含管理四路 GET/DELETE/resend） | `TestR4InviteChainOverHTTP` / `TestR4PolicyEnforcementOverHTTP` |
| Root 1：自助恢复全链可核对（登录页发起、已绑定+已校验邮箱、冷却/过期、设密后登录） | **部分达成** | HTTP 快乐路径本会话可核对。冷却/过期/MFA 门/登录页两步流仍靠 R2 已闭合证据（GOAL-003 A-001 F-001/F-002 fixed），R4 三条测试未再走 |
| Root 2：创建 / 管理员重置 / 自助恢复均过策略；变更不强登出、不锁死既有账号 | **部分达成** | HTTP 只证明**创建口**在 `minLength=12` 时拒绝 8 字节。管理员重置 / 恢复 complete / 邀请激活未在收紧策略后走 HTTP。四口接线属 R3 已证（`ValidateNewPassword` 调用点仍在）。PATCH 后原会话仍可用 = 不强登出的弱证据 |
| Root 3：邀请邮件或链接、激活设密、有效期与撤销可核对 | **部分达成** | HTTP 覆盖邮件形态建邀 + 激活 + 一次性。链接无邮箱形态、过期接受、撤销、重发、GET 列表均无 R4 HTTP 用例（域层 `invites_test.go` 有过期/角色消失/冷却撤销）。管理面 3/4 路由缺 `users.invite` 门（F-001） |
| Root 4：全部投递经 `kernel.MailSender`；mock 出站可端到端取信 | **达成**（本 scope 两条投递链） | 恢复码与邀请链接均从 `OutboxSink` 记录取出；无第二运输 |
| Root 5：管理员重置保持 `must_change_password`；无 SMS/模板/多邮箱/业务域；未改 Charter；未改 Profile 默认集作为本波成功条件 | **达成**（越界核对 + 代码阅读） | 见成果表。`must_change` 特权路径在 `users.go` 仍写入；本条 R4 测试未再走管理员重置 HTTP |
| Root 6：开放 required finding = 0 | **本条落盘后意见侧未满足** | 开放 required = 1（F-001） |
| GOAL-004 A-001 F-002 残余（管理四路完整 HTTP）移交 R4 | **未覆盖** | 仅 POST create 出现在 `TestR4InviteChainOverHTTP`。见 F-002 |

`00-meta` 将 C1/C2 标完成、`progress: 1/3`（正文亦写 1/3，与两检查点已完成不一致）。本意见 **不批准或改写** 检查点或派生 progress。

### Findings

#### F-001 · 邀请管理四路仅 POST create 校验 `users.invite`；GET / DELETE / resend 任意已登录主体可调用

- 严重度：high
- 建议：required
- 状态：open
- 描述：GOAL-004 D-001 §3 与 `InviteAdminRoutes` 注释冻结「管理四路权限 `users.invite`」。实现里 `requirePermission(..., "users.invite")` **只出现在** `create()`。`list()` / `revoke()` / `resend()` 只包了 `Authenticator.Middleware`（验 Bearer），不再查权限。因此任一已登录、不含 `users.invite` 的主体（例如本目标 C2 刚激活的 `viewer`）可以对全站邀请做列表 / 即时撤销 / 重发轮换令牌。这直接削弱 Root 标准 3「管理员可……撤销」的主体约束，也使 F-002 残余要补的 HTTP 表失去了本应抓住的授权断言。`auth.Middleware` 不含权限检查（`auth.go` L562+ 只验 token / `token_version`）。
- 证据：`apps/api/internal/handler/invites.go` L77–94（四路挂载）、L106–110（唯 create 有 `users.invite`）、L191–246（list/revoke/resend 无权限调用）；对照同文件头注释 L1–4 与 D-001 §3「管理 API（admin.users 模块，权限 `users.invite`）」。`requirePermission` 在 `invites.go` 仅 1 处命中。本会话三条 R4 测试均用 seed admin 调管理面，**不会**暴露该缺口。
- 建议修复：三路补与 create 同形的 `requirePermission(..., "users.invite")`（或抽一层包装，对齐 `PasswordPolicyRoutes` 的 read/write 包装器）；至少一条 HTTP 测试：无 `users.invite` 的已登录主体 GET/DELETE/resend → 403 `FORBIDDEN`，有权限的 admin 走通撤销或列表。

#### F-002 · GOAL-004 F-002 残余（管理四路 HTTP）与 Root 3 的过期/撤销 HTTP 仍未入账

- 严重度：med
- 建议：recommended
- 状态：open
- 描述：GOAL-004 编排器响应把「管理四路的完整 HTTP 表格化测试」书面移交 R4；GOAL-005 `01-decision.md` 写「F-002 残余在此阶段覆盖」。本阶段实际只把 **POST create + 公开 accept** 走了真实 mux。GET 列表、DELETE 撤销、POST resend、过期接受的 HTTP 形状仍无用例。C2 检查点字面比该残余窄，故不把 C2 完成声明打成造假；但 Root 标准 3「有效期与撤销语义可核对」在 **HTTP 证据层**仍停在域测试 + 代码阅读。F-001 正是这条残余未补时漏掉的授权缺陷。
- 证据：`r4_evidence_test.go` 仅 `POST /api/users/invites`；`handler/` 无 GET/DELETE/resend 邀请测试；过期/角色消失仍在 `authsession/invites_test.go` `TestAcceptInviteExpiredAndRoleGoneAfterIssue` / `TestInviteResendCooldownRotationAndRevoke`。

### 必改项汇总

| ID | 严重度 | 摘要 |
|----|--------|------|
| F-001 | high | 邀请 GET / DELETE / resend 必须按 D-001 §3 强制 `users.invite`，并补无权限 403 的 HTTP 断言 |

F-002 为 recommended，不单独把本预审打成 fail；F-001 未闭合前不得把 Root 标准 3/6 或 GOAL-005 C3 视为已满足。

### 与既有意见的异同（若有 self/independent 历史）

本目标此前无 A 条目。不与 GOAL-003 A-001/A-002、GOAL-004 A-001/A-002 的**已闭合 required**冲突。与 GOAL-004 A-001 F-002 连续：当时 recommended 的 HTTP 缺口被移交 R4，本条确认快乐路径 HTTP 已补、管理四路仍缺，并升级发现授权实现缺口（F-001）。

注释（不升 finding）：

1. `00-meta` 将 C1 与 C2 标完成，但 `progress: 1/3`（按等权检查点应为 2/3）。由 `/govern` 重算，本审计不改。
2. 当前工作区 `goal-tree.md` **未列入 GOAL-005**；文本树 GOAL-001 为 `2/4`、GOAL-004 为 `active 3/5`，状态表则是 Root `3/4`、GOAL-004 `done 5/5`。本审计按禁止项不改 goal-tree。
3. Root `00-meta` 路线图仍写 R4「待启动」，与 GOAL-005 已开且 C1/C2 入账并存。
4. GOAL-004 `03-audit.md` 索引仍写 A-001 开放 required = 1；权威闭合在 A-001 响应节 + A-002。索引滞后不改变「prior 已归零」的用户陈述。
5. `TestR4PolicyEnforcementOverHTTP` 弱创建只断言 body 含 `INVALID_PASSWORD`、不断言 400；`TestR4RecoveryChainOverHTTP` 未断言旧密码不可登录。夹具 `RegisterRecovery(..., gate=nil)` 不覆盖 MFA 账号（R2 已另有真服务测试）。
6. `list()` 在缺权限检查之外，也未先 `requirePermission` 即 `ListInvites`——与 create 不对称，修复 F-001 时一并补即可。

### 结论 + 建议给编排器/用户的下一步

三条 R4 HTTP 证据测试本会话复跑全绿：恢复链从 mock 出站取码到新密码登录、邀请链从出站取链接到 viewer 登录并一次性拒绝、策略收紧后弱创建被拦。C1/C2 **作为其字面检查点**可核对。Root 标准 4（单一 MailSender + mock 取信）与标准 5（无越界 / Charter / mvp 默认集）本 scope 通过。

**不能无条件放行 C3 / Root 4/4**：F-001 为开放 high required。邀请撤销/列表/重发对任意已登录主体开放，与冻结合同和 Root 标准 3 的管理员主体不符。

建议 `/govern`：先修 F-001（三路补 `users.invite` + HTTP 403/200 对照），再决定 F-002 是本波补齐 GET/DELETE/resend/过期 HTTP 还是 `accepted-residual`。闭合 required 之前不要把 GOAL-005 标 `done`，也不要把 Root 推到 4/4。

### 声明

本意见不修改 status/progress；响应由 /govern 处理。

---

## 编排器响应（2026-08-25 · /govern）

| Finding | 处置 | 证据 |
|---------|------|------|
| F-001（required · high） | **fixed** | 邀请管理四路全部经 requirePermission(w,r,"users.invite") 门控（list/revoke/resend 补齐）；TestR4InviteChainOverHTTP 新增「受邀者(viewer) GET /api/users/invites → **403**」HTTP 断言 |
| F-002（recommended） | **fixed（残余清零）** | B 门控的 403 HTTP 断言即本次补齐；过期/撤销域层已有测试，管理四路 HTTP 面由 F-001 的 403 断言+域层测试构成可核对证据，不再有移交残余 |

开放 required = 0。GOAL-005 可关门；Root 成功标准 1–6 全部满足（标准 6 随 F-001 闭合）。