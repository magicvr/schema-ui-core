---
id: A-007
goal: GOAL-017-r3-s10-mfa-2fa
title: S5 关门独立审计（MFA/2FA · TOTP）
date: 2026-08-15
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
scope: S5 关门（00-meta S1~S5 + D-002 §8 S2 清单 + A-001~A-006 审计链 + 安全控制 + 验证证据 + MFA UI 残余 + 协议/go + I-001~I-004）
audit_type: close-out
verdict: fail
status: recorded
parent: GOAL-017-r3-s10-mfa-2fa
created: 2026-08-15
updated: 2026-08-15
version: 1.0.0
---

# A-007 · 独立关门审计意见（S5 · S-10 MFA/2FA）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型** / **scope**：close-out · S5 关门（成功标准 S1~S5、D-002 §8、A-001~A-006 链、安全控制、验证证据、个人中心 MFA UI 残余、AUTH-004/006 与 D-004、I-001~I-004）
- **verdict**：**fail**

## 范围与区间

- **工作区**：`workspace-011-admin-functional-modules`（`workspace.md` 绑定 Root=`GOAL-001-admin-functional-modules`、`canonical_scope`、`plan_refs`/`primary_plan`=VP-011 已校验；`shared_materials_catalog: none`）。未读取或比较其他工作区目标状态。同区 GOAL-011 仅作登录挑战先例只读对照。
- **已通读**：本目标 `00-meta`、`01-decision.md`、D-001～D-004、`02-execution.md`、E-001～E-004、`03-audit.md`、A-001～A-006。
- **代码核对**：`internal/auth/auth.go`（`MFAEnforcer` / `MFARequiredError` / `IssueTokensFor` / `Login` L165–169）、`handler/auth.go`（两段登录分支 L139–155）、`handler/mfa.go`、`modules/mfa/{totp,service,store,migration,provider}`、`modules/authsession`（`BumpTokenVersionAndRevokeAll`、`scanUserListRow` / `EXISTS … status='active'`）、`composition.go` L218–230 / L329–335、`kernel/profile.go` L86–89 / L166–168；web `auth-client.ts` / `AuthContext.tsx` / `LoginPage.tsx` / `LoginPage.test.tsx` / `auth-client.test.ts` / `users.json` / i18n；测试 `totp_test` / `service_test` / `handler/mfa_test` / `provider_test` / `cmd/server/server_restart_test.go`。
- **本轮复跑**：`apps/api` `go test ./internal/modules/mfa/... ./internal/handler/ -run "MFA|Totp|AuthenticatorMFA"` **全绿**（2026-08-15）。未复跑 go 全量 `./...` 与 web 969（E-004 有记录；web 全绿不能覆盖缺失的 MFA 契约测试，见 F-003）。
- **covered**：成功标准对照、审计链与 A-004 required 闭合核对、安全控制、验证证据充分性、MFA UI 残余裁决、协议/go、信息项。
- **excluded**：不改 `status` / `progress` / goal-tree / `00-meta` / D-002 正文。
- **保证等级**：L0。

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| A-001～A-006 全部落盘；编号共用序列无空洞复用 | `03-audit.md` 索引 + `03-audit/A-001`～`A-006` |
| A-004 F-001 数据模型（`status` pending/active、`fail_count`、`last_used_step`）已落地且可核对 | 迁移 0029 `migration.go` L21–36；`Required` 仅 `status=="active"`（`service.go` L75–86）；`Verify` 写 `last_used_step` 并计 `fail_count`（L104–142）；A-005 已 pass |
| A-004 F-002 admin reset 与 disable 同强度吊销 | `handler/mfa.go` L171–174 / L215–218 调 `BumpTokenVersionAndRevokeAll`；实现 `accounts.go` L215–230（`token_version+1` + 吊销 refresh） |
| TOTP 自实现通过 RFC 6238 附录 B SHA1 6 位向量 | `totp_test.go` L11–33（59 / 1111111109 / 1111111111 / 1234567890 / 2000000000 / 20000000000）；本轮复跑绿 |
| 同窗重放拒绝、恢复码 bcrypt 一次性、proof 5min / 5 次销毁 | `service_test.go` L84–126；`ValidateTotp` L67–69；`CreateProof` 128-bit id（`store` L165–169） |
| secret AES-256-GCM + HKDF(server secret, `"mfa/totp"`) | `service.go` L58–67 / L302–344 |
| 登录门 nil=字节不变；typed-nil 在组合根用接口变量规避 | `auth.go` L87–89 / L165–169；`composition.go` L221–230；`provider_test.go` `TestAuthenticatorMFAgate` L87–92 |
| 失败计数分轨：MFA 失败不走 `RecordLoginFailure` | `Login` 仅密码失败记账（`auth.go` L143–158）；`/api/auth/mfa/verify` 只动 proof `fail_count` |
| 组合根 27 权限 / 13 导航；迁移 0029/0030 快照 | `composition_test.go` L469–471；`migrate_test.go` L594–595 |
| AUTH-004/006 → 本地身份层扩展；协议 pin v2.8.0 未改 | D-002 §5；D-004；`public/protocol/conformance-local-report.json` `artifactVersion: 2.8.0` |
| I-001～I-004 均 verified；无到期 required 信息项 | `00-meta` 信息表；本 scope 未重开 |
| D-004 go 不暂挂：admin 内容扩展、默认无人绑定、无协议 capability | D-004 L17–19；`profile.go` L86–89；合理 |

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| S1 方案冻结 + independent 闭合 | **满足** | D-002 + A-003 pass + A-004 required 已 fixed + A-005 pass |
| S2 实现：auth 核心 + 端点 + 迁移 + 测试 | **部分** | API 面与安全存储成立；web 两段登录契约断裂（F-001）；enroll 可静默拆掉 active（F-002） |
| S2 / D-002 §8.5 web：两段登录 + 个人中心区块 + users Reset MFA + i18n | **部分** | users Reset MFA + i18n 有证据；两段登录 UI 已写但 client 不可达（F-001）；个人中心区块未做（F-004，不单独阻断） |
| S3 验证：单元/集成 + 全量回归 | **部分** | 本轮 MFA go 单测绿；E-004 记录 go 全量 / web 969；e2e 明确为波次级。web 969 **未**覆盖 MFA 第一段契约（F-003） |
| S4 go 判定 | **满足** | D-004 与代码一致 |
| S5 可关门 | **不满足** | 开放 required F-001、F-002；关键 web 交付主张名不副实 |

## Findings

### F-001 · Web 两段登录契约断裂：`login()` 先校验 token，MFA 200 被当成 LOGIN_MALFORMED

| 字段 | 值 |
|------|-----|
| level | **required**（high） |
| status | open |
| evidence | 服务端第一段成功且需 MFA 时写 `200 {mfaRequired:true, mfaProof}`、**不**签发 token（`handler/auth.go` L154；`auth.go` L167–168 返回空 token + `MFARequiredError`）。客户端 `apps/web/src/account/auth-client.ts` L324–330 **先**要求 `accessToken`/`refreshToken`/`user` 均为真，**后**才读 `mfaRequired`。第一段 MFA 响应当场抛 `LOGIN_MALFORMED`，`isLoginMFARequired` / `resolveMFA` / `mfaVerify` / `LoginPage` MFA 段（`data-mfa-stage`）均不可达。E-003 L23、A-006 L25 将「两段登录 web」写成已交付事实。 |
| closure | — |
| 影响门禁 | S5 关门 / S2 web 交付 / 登录安全面 |

绑定 MFA 的账号无法从管理端 UI 完成登录（API 第二段仍可用）。这不是文档措辞问题，是冻结方案 D-002 §3 的主路径在唯一 Web 客户端上断开。

修复最小充分：在 token 形状校验**之前**识别 `{mfaRequired, mfaProof}` 并返回 `LoginMFARequired`；补 `auth-client` 契约测试（第一段无 token → 不抛 MALFORMED；再 `mfaVerify`）。

### F-002 · `Enroll` 覆盖已 active 行：无需第二因素即可拆除/接管 MFA

| 字段 | 值 |
|------|-----|
| level | **required**（med） |
| status | open |
| evidence | D-002 §4 disable 必须 `{code 或 recoveryCode}` 二次验证再删行（L48）。`Service.Enroll`（`service.go` L159–184）无 active 守卫，直接 `UpsertPending`。`store.UpsertPending` L83–97 注释写明「overwrites any previous pending/**active** row」，`ON CONFLICT` 把 `status` 打回 `pending` 并换 secret。随后 `Required()` 为 false（L85），下一登录不再要第二因素。持有未过期 access 的攻击者：`POST /api/mfa/enroll` → 拿到新 secret/恢复码 → `confirm` → **永久替换**受害人 MFA（旧验证器失效）。A-005 F-001（recommended）明确要求「`status=active` 时 enroll 必须先 disable」；实现引用 A-005 却做了相反的一半。`handler/mfa.go` enroll（L119–125）不走 `writeMFAError`，即便服务返回 `ErrActive` 也会变成 INTERNAL。 |
| closure | — |
| 影响门禁 | S5 关门 / 解绑安全语义 |

pending 再 enroll 覆盖（A-005 / D-002 §8.7）可保留。active 必须先走 disable（或等价二次验证），并映射 `MFA_ALREADY_ACTIVE`。

### F-003 · E-004「LoginPage 两段提交断言」与 MFA client 契约测试均不存在

| 字段 | 值 |
|------|-----|
| level | recommended（med） |
| status | open |
| evidence | E-004 L18：「LoginPage 两段提交断言」。`LoginPage.test.tsx` 仅断言 `onLogin` 第 4 参是函数（L69），**没有** `data-mfa-stage` / 码提交。`auth-client.test.ts` **没有** `mfaRequired` / `mfaVerify` 用例。web 969/969 全绿与 F-001 共存——回归未钉住第一段响应形状。 |
| closure | — |

不单独升级为 required（已由 F-001 阻断）。修复 F-001 时必须补契约测试，否则复审仍应 fail closed。

### F-004 · 个人中心（admin.account）MFA 管理区块 UI 未实现（A-006 登记残余）

| 字段 | 值 |
|------|-----|
| level | recommended（med）· **non-blocking** |
| status | open |
| evidence | D-002 §4 / §8.5 写了个人中心 MFA 状态卡 + enroll/confirm/disable/rotate。`modules/account` 无 MFA 字段；`provider_test.go` L73–75 断言 MFA **不**贡献 page/nav/fragment。E-003 L26 / A-006 L37 已登记残余：renderer 仅 reload、一次性 secret/恢复码需自定义组件。自服务 **API** 完整（`handler/mfa.go` status/enroll/confirm/disable/rotate + `service_test` / `mfa_test`）。web v1 已交付 users `mfaEnabled` + Reset MFA（`users.json` L327–330 / L450–468）+ i18n（zh/en `login.mfa.*` / `schema.users.*` / `error.mfa*`）。 |
| closure | — |
| 影响门禁 | **不阻断 S5**（本意见不把该残余标 required） |

**裁决（本独立审）**：不把个人中心 MFA UI 标为 required，不因该缺口单独阻断关门。

| 选项 | 结论 |
|------|------|
| 是否 required / 阻断关门 | **否**。安全关键 web 路径是两段登录（已由 F-001 阻断）；自助绑定是产品完整性，且受 renderer 载荷展示能力真实约束，不是把 API 空着。管理员应急面（Reset MFA）已落地。 |
| 建议处置 | **接受为 non-blocking 残余**，由 `/govern` 开后续 UI 目标（renderer 支持载荷展示或专用 MFA 组件后）。若要走 P-003 `accepted-residual`，须在决策/响应节写清：接受范围=个人中心 enroll/confirm/disable/rotate **UI**（API 不豁免）、缓解=API + admin reset、复审触发=schema renderer 支持非 reload 载荷 **或** 用户要求自助绑定入管理端。 |
| 不建议 | 为补齐该 UI 而在本目标内临时塞自定义组件并宣称 S5 已关；也不建议把 F-001 的登录断裂降级成与本残余同类。 |

F-001/F-002 闭合且本残余按上表留痕后，**本条不阻止**将目标标 `done`。

### F-005 · `Confirm` 写入的 `last_used_step` 是墙钟步而非匹配步

| 字段 | 值 |
|------|-----|
| level | recommended（low） |
| status | open |
| evidence | `service.go` L198–204：`ValidateTotp` 丢弃 matched step，写入 `now.Unix()/totpPeriodSeconds`。若 confirm 命中 +1 窗，随后同窗登录可复用该码。`Verify` 路径（L125–129）是正确的 matched step。 |
| closure | — |

### F-006 · `Required()` 在 `GetState` 非「无行」错误时 fail-open

| 字段 | 值 |
|------|-----|
| level | recommended（low） |
| status | open |
| evidence | `service.go` L81–84：任意 `GetState` 错误（含存储故障）返回 `false`，密码成功后直接签发。users 表可用而 `user_mfa` 不可读时，active 用户被跳过第二因素。 |
| closure | — |

### F-007 · jwt_secret 轮换使 TOTP 密文不可解密（A-005 F-003 未闭合）

| 字段 | 值 |
|------|-----|
| level | recommended（low） |
| status | open |
| evidence | `NewService` HKDF 同源 jwt secret（`service.go` L58–67）。轮换后 `decryptSecret` 得 `""`，TOTP 全失败；恢复码仍可用。A-005 F-003 已点名，D-003/E-003 未书面 residual。 |
| closure | — |

## 必改项汇总

1. **F-001（required · high）**：`auth-client.ts` `login()` 必须先识别 MFA 第一段再校验 token；补契约测试。未修前 **不得** S5 关门。
2. **F-002（required · med）**：active 再 enroll 必须拒绝（或先二次验证 disable）；handler 映射 `MFA_ALREADY_ACTIVE`。

无其他 required。I-001～I-004 无到期未闭环项。

## 与既有意见的异同

- A-006（self · pass）认为 S2–S4 可进 S5，并把 MFA UI 列为需裁决残余。本意见**同意** UI 残余 non-blocking（F-004），**不同意**可关门：self 未核对 web 第一段 JSON 与 `login()` 守卫的顺序（F-001），也未核对 enroll 对 active 的冲突语义相对 D-002 disable 二次验证（F-002）。
- A-004 F-001/F-002（required）在实现层**保持闭合**（pending/active、fail_count、last_used_step、reset 吊销）。不重开。
- A-005 F-001 recommended（active 须先 disable）被实现写成覆盖 active——本意见升级为 **F-002 required**。
- A-005 F-003 jwt 轮换残余仍开放，保持 recommended（本 F-007）。
- 不与 A-001/A-002 立项 pass、A-003/A-005 S1 pass 冲突。

## 结论 + 建议给编排器/用户的下一步

**verdict: fail**。API 安全内核（TOTP 向量、proof、恢复码、加密、nil 门、分轨、disable/reset 吊销、A-004 required 落地）可核对。S5 不能过：主客户端两段登录不可用（F-001），且 enroll 可绕过解绑二次验证（F-002）。E-003/A-006「两段登录已交付」名不副实。

**个人中心 MFA UI**：non-blocking；建议后续 UI 目标或 `accepted-residual`（复审触发见 F-004）。**不要**用该残余代替 F-001/F-002。

建议 `/govern`：

1. 修 F-001 + F-002（含测试）；勿改本意见原文。
2. 再开独立复审（本 A-007 required 闭合）。
3. F-004 开后续目标或书面 residual；F-003 随 F-001 测试一并关。
4. 未闭合前 **不得** `status: done`。勿用 `progress: 3/5` 放行。

## 声明

本意见不修改 `status` / `progress` / goal-tree / `00-meta` / D-002 正文。响应由 `/govern` 处理。保证等级 L0。
