---
doc_type: vision-plan
id: VP-019-iam-recovery
title: IAM：密码策略 / 邀请入职 / 自助恢复状态机
status: closed
vision_ref: schema-ui-core-admin-foundation@0.3.0
lead_workspace: workspace-019-iam-recovery
created: 2026-08-25
updated: 2026-08-26
version: 0.3.0
parent: null
---

# VP-019 · IAM：密码策略 / 邀请入职 / 自助恢复状态机

## 状态与门闩（2026-08-26 · **closed** · v0.3.0 交付后关门）

| 项 | 值 |
|----|-----|
| status | **`closed`**（v0.3.0 · 2026-08-26 用户书面确认关门；实现于 2026-08-25 同日全链交付，Root `GOAL-001-iam-recovery` `done 4/4`） |
| 组合位置 | Admin 功能分支下一拍（roadmap §Admin 功能「再下一截 = IAM VP」）；已交付收官 |
| 硬前置 | ✅ **VP-017 已按现行分母再关门**（v0.5.0 · 渠道模型 = mock 默认 + Resend 生产 + 设置/试发）；✅ **VP-018 已关门**（v1.0.0 · `users` 可空 email + 绑定/校验状态机 + 换绑）。两条身份/运输前置均满足 |
| 激活门禁 | ✅ 2026-08-25 全部满足：[VRev-043](../reviews/VRev-043-vp019-iam-recovery-intent-activation.md)（independent · grok build）`pass`（required = 0；V-F076/077/078 → fixed）；VP-008 `go` 消费有效性 **Admin 类 freshness** PASS（基线 `092bf37` → 现行 `66f5fd1f`，不暂挂 `go`）；VP-009/VP-010 无开放阻断 |
| 结构 | 新 VP + 单 delivery 工作区（`workspace-019-iam-recovery`，2026-08-25 开区）；不重开 VP-012～018；不改 Charter |

## 意图

在 VP-003 单主线、已交付的 `core.auth-session` 账号模型、VP-018 账号邮箱身份（`users.email` + 校验状态机）与 VP-017 内核 `MailSender`（可切换渠道）之上，把「忘记密码」从空转收成**可核对的自助恢复闭环**，并顺带把密码策略与邀请入职收成同一 IAM 波次的 Admin 能力：

1. **自助恢复状态机（忘记密码）**：登录页「忘记密码」发起；**证明依据 = 账号已绑定且已校验的邮箱**（产品事实 2026-08-22 用户确认）；投递形态（验证码 vs 魔术链接，I-019-001）经 `kernel.MailSender` 现行渠道投递；支持过期与重发冷却（I-019-002）；完成后设置新密码并恢复登录。无邮箱账号不自助恢复，走管理员重置。
2. **密码策略**：可配置密码策略产品面（长度 / 复杂性 / 历史等；默认参数与配置边界 R2 冻结，I-019-003）；创建、管理员重置、自助恢复设定新密码时统一强制校验；策略变更对既有账号**渐进生效**（I-019-007），不做强制登出。不解锁自助恢复的并行规划按 roadmap 仍成立，但本波三者同为分母。
3. **邀请入职**：管理员生成邀请并可经同一 `MailSender` 投递（邀请信），也可出示链接形态（roadmap 产品事实：邀请仍可用管理员出示链接）；受邀者在无既有账号情况下完成激活并设置密码（I-019-004）；邀请有效期 / 失效 / 撤销一次性语义可核对（I-019-005）。
4. **管理员重置**：继续走既有特权路径（`must_change_password`），不冒充自助恢复；两种恢复路径并存是产品事实（2026-08-22 用户确认）。

## 配置面与模块归属

- 邮箱身份、恢复、邀请、密码策略走**既有** `core.auth-session`（及消费它的 `admin.users` / `admin.account` / 设置面），**不是**新模块、不新增 Profile、不改 Profile 默认集 / 模块矩阵 / Manifest 装配语义。
- **缺省运输**：无生产渠道时校验/恢复/邀请信进 VP-017 现行默认渠道（mock 站内出站记录），可端到端取出核对；不把 Mailhog / 真实 SMTP / Resend 做成 mvp/dev 硬依赖。
- **生产**：显式配置后走 VP-017 当前生产渠道（Resend）；配置不完整 fail-closed（继承 017 语义）。
- **密码策略配置面**：系统设置既有 tab 的扩展（形状 R2 冻结），不是独立模块。

## 首波冻结（退出分母 = IAM 三件）

| 能力 | 本 VP 交付 | 不进本 VP |
|------|------------|-----------|
| 自助恢复 | 登录页发起 → 投递 → 校验 → 设新密码 → 恢复登录；过期 / 重发冷却可核对 | 恢复即改邮箱、恢复即解绑；SMS/推送恢复（RT-M02 gated） |
| 密码策略 | 可配置策略 + 新密码统一强制校验 + 渐进生效 | 强制登出、锁死既有账号、密码/凭证历史迁移工具、KMS |
| 邀请入职 | 管理员生成 + 邮件/链接投递 + 受邀激活（设置密码）+ 有效期/撤销 | 组织/部门/岗位邀请、SCIM/OIDC、批量导入 |
| 运输 | 只消费 `kernel.MailSender`（当前渠道） | 第二 SMTP、HTML/MIME 模板引擎、Notification Transport 产品、消息模板中心 |
| 管理员重置 | 保留 `must_change_password` 特权路径 | 管理员冒充自助恢复、重置证据代替投递校验 |
| 消费者 | 本波交付 IAM 三件的 API/页面（设置面 + 登录页 + 管理员邀请/重置） | 用户站内通知产品、SMS |

## 非目标

- **SMS / 推送 / IM**（RT-M02 仍 gated）
- 模板中心、消息模板管理、重做 F-04 站内通知
- 多邮箱 / 别名 / 邮箱即登录名；强制全站账号必须有邮箱才能启动
- 组织 / 部门 / 岗位、数据权限 `org` 扩展（Admin 分支后续候选，不在本波）
- OIDC / SSO / SCIM（扩展接缝，trigger-gated）
- 事务 outbox / 外部邮件队列（RT-Q06 仍 gated）
- 改 Profile 默认集 / 模块矩阵 / Manifest 装配语义
- 重开 VP-012～018；替代 VP-009 / VP-010
- 改变 Charter 边界；业务域页面

## 与相邻 VP 的边界

| VP / 分支 | 关系 |
|-----------|------|
| **VP-003** | 遵守薄内核。IAM 三件是 `core.auth-session` 账号能力扩展，不是平行认证模块 |
| **VP-017** | 运输前置（已 `closed` v0.5.0）。只消费 `kernel.MailSender` 现行渠道，不重开渠道模型 |
| **VP-018** | 身份前置（已 `closed` v1.0.0）。消费 `users.email` + 校验状态机；不重开邮箱身份波 |
| **VP-012** | 已 closed 的横切契约不重开。恢复/邀请/改密的审计用既有 envelope |
| **VP-008 `go`** | 本 VP 是 Admin 功能（IAM 面），不是 Tier D 业务域。激活前完成 Admin 类 freshness（基线 `092bf37`）；若实现改变 Profile 默认集 / 模块矩阵 / Manifest 装配 / 共同门禁，按消费有效性暂挂 |
| **VP-009 / VP-010** | 邮箱枚举、重放、open-relay、邀请滥用等安全 finding 与符合性 gap 仍归持续程序 |
| **业务域** | 不得在本 VP 加订单/营销邮件 |

## 方向级退出判据

1. 自助恢复全链落地且可核对：从登录页发起，经已绑定+已校验邮箱投递（现行渠道），重发冷却与过期语义可验证；完成后设新密码并正常登录。
2. 密码策略已配置并强制：创建 / 管理员重置 / 自助恢复的新密码均过策略校验；策略变更不强制登出、不锁死既有账号。
3. 邀请入职落地：管理员可生成邀请（邮件或链接形态），受邀者可激活账号并设置密码；有效期与撤销语义可核对。
4. 全部投递经 `kernel.MailSender`（当前渠道），无第二运输；无生产渠道时 mock 出站记录可端到端取信。
5. 管理员重置保持 `must_change_password` 特权路径；未引入 SMS / 模板中心 / 多邮箱 / 业务域；未改 Charter；未改 Profile 默认集作为本波成功条件。
6. 开放 required finding = 0（或已合法闭合）。

## 信息需求（P-005）

允许带未知立项。下列必须在对应阶段前关闭或经用户接受残余。

| id | 要回答的问题 | 级别 | 影响门禁 | 最晚阶段 | 状态 |
|----|--------------|------|----------|----------|------|
| I-019-001 | 自助恢复证明形态：验证码 vs 魔术链接（**默认候选 = VP-018 已冻结 6 位码** `email_identity.go:31-36`；R1 确认，不做空白选择题） | required | 方案冻结 | R1 合同冻结 | collecting |
| I-019-002 | 恢复令牌/验证码 TTL 与重发冷却 | required | 方案冻结 | **R2 方案冻结前** | collecting |
| I-019-003 | 密码策略默认参数（长度/复杂性/历史）与配置边界；**策略 UI 仅 `admin.settings`（与邮件 tab 同形），强制面在 `core.auth-session` 对全 Profile 生效，禁止为策略把 `admin.settings` 加入 mvp 默认集** | required | R2 方案冻结 | R2 | collecting |
| I-019-004 | 邀请形态：邮件投递 vs 管理员出示链接；目标是否须预置账号 or 邀请即建号 | required | R3 方案冻结 | R3 | collecting |
| I-019-005 | 邀请有效期 / 失效 / 撤销 / 一次性语义 | required | R3 方案 / 实施 | R3 接入前 | collecting |
| I-019-006 | 无邮箱账号的自助恢复边界：仅管理员重置（**2026-08-22 产品事实已确认**） | required | 方案冻结 | R1 | **registered**（2026-08-22 用户确认：无邮箱不自助、走管理员重置；R1 验收到位） |
| I-019-007 | 密码策略对既有账号生效边界（渐进 vs 登录强制改密） | required | R2 方案冻结 | R2 | collecting |
| I-019-008 | 恢复 / 邀请 / 改密后的会话语义（是否全端登出或撤销既有 session） | non-blocking | 方案冻结 | R4 | collecting |
| I-019-009 | **已登记 MFA 的账号如何走自助恢复**（绕过第二因子 / 要求 TOTP / 拒绝自助只走管理员重置）——`admin.mfa` 已在 admin 默认集、登录页已有 MFA 错误码；不登记则 R2 可能做成 MFA 旁路（V-F076） | required | 方案冻结 | R1 合同冻结 | collecting |

详细纲领阶段由 lead Root（P-001）书写：R1 合同冻结 → R2 自助恢复全链 → R3 密码策略 + 邀请入职 → R4 证据 / 关门。本 VP 不写 Goal 五件套。

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| workspace-019-iam-recovery | GOAL-001-iam-recovery | lead | 2026-08-25 | VRev-043 `pass`（V-F076/077/078 → fixed）后激活并开区；Root P-001 R1～R4 |

## 关门记录

（仅 `closed` / `abandoned` 时填写。）

| date | outcome | summary | evidence_links | residuals |
|------|---------|---------|----------------|-----------|
| 2026-08-26 | closed（v0.3.0 · 用户 ask_user_question 书面确认） | IAM 三件交付收官：①自助恢复全链（登录页发起 → 迁移 0056 挑战表 → 6 位码投递 → MFA 第二因子门 → 设新码 → 会话撤销语义；防枚举同形 202、挑战预算 ≤5、重发冷却 60s、IP\|identifier 限流）；②密码策略配置面 + 四口强制统一 `INVALID_PASSWORD`（MinLength 权威咬合）；③邀请入职双形态投递 + 即建号 + 有效期/撤销/重发 + Web 面（含 bodyMapping roles 缺陷修复）。全部投递经 `kernel.MailSender` 现行渠道（r4_evidence_test 三链真实 mux 取码）；无越界（composition 黄金计数吻合）。R1–R4 = GOAL-002～005 全 done（2026-08-25 同日）；关后维护 E-008/E-009 + 独立复审 A-001 `pass`（代码证据制全量复跑绿）/ 响应 A-002 | [workspace-019 goal-tree](../../workspaces/workspace-019-iam-recovery/goal-tree.md)；Root [E-007-root-closeout](../../workspaces/workspace-019-iam-recovery/GOAL-001-iam-recovery/02-execution/E-007-root-closeout.md)；关后 [03-audit/A-001-closeout-independent.md](../../workspaces/workspace-019-iam-recovery/GOAL-001-iam-recovery/03-audit/A-001-closeout-independent.md) + [A-002](../../workspaces/workspace-019-iam-recovery/GOAL-001-iam-recovery/03-audit/A-002-finding-response-self.md)；commits `299f8f52`/`9628ca8f`/`2f088d55`→`9ced003d`、`ce96df92` | 本波无新增未闭合残余；继承性平台边界照旧——SQLite lower() ASCII 折叠（VP-018 N-1 口径）、认证限流进程内单实例边界（[workspace-009] W12 D-002 2026-08-26 维持，复审触发 = 多实例形态出现） |

## 规划修订短史

| date | change |
|------|--------|
| 2026-08-25 | 初创 `planned`：按 roadmap Admin 功能分支「再下一截 = IAM VP」立项。硬前置 VP-017（v0.5.0 再关门）与 VP-018（v1.0.0 关门）均已满足；退出分母 = 密码策略 / 邀请入职 / 自助恢复状态机。VRev-043（independent · grok build）待办 |
| 2026-08-25 | **v0.2.0 `planned → active`**：VRev-043（independent · grok build）`pass`（0 required；V-F076/077/078 → **fixed**）；Admin 类 freshness **PASS**（`092bf37` → `66f5fd1f`，不暂挂 `go`；W11 已恢复 `go`、F-007 不升格）；lead `workspace-019-iam-recovery` 同日开区（Root `GOAL-001-iam-recovery`）。信息项按 V-F076 修正：I-019-002 最晚阶段 → R2 方案冻结前；增补 I-019-009（MFA，最晚 R1）；I-019-001 默认候选 = VP-018 6 位码；I-019-003 冻结策略面/Profile 边界；I-019-006 → registered（2026-08-22 产品事实） |
| 2026-08-26 | **v0.3.0 `active → closed`**：实现已于 2026-08-25 同日全链交付（Root done 4/4，六条方向级退出判据全达成）；关后独立复审 A-001 `pass`（代码证据制全量复跑绿）+ recommended ×2 于当日闭合（F-001 fixed / F-002 登记移交 [workspace-009] W12）。用户 ask_user_question 书面确认关门；历史绑定保留，不重开 |