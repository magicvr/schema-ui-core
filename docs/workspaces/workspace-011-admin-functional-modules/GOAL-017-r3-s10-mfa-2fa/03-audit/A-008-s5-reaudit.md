---
id: A-008
goal: GOAL-017-r3-s10-mfa-2fa
title: S5 关门 A-007 required 闭合复审
date: 2026-08-15
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
scope: A-007 F-001/F-002 required 闭合验证（login() 第一段顺序 + 契约测试；Enroll active 守卫 + store WHERE + 回归；S5 关门门禁）
audit_type: finding-closure
verdict: pass
status: recorded
parent: GOAL-017-r3-s10-mfa-2fa
created: 2026-08-15
updated: 2026-08-15
version: 1.0.0
---

# A-008 · 独立复审意见（S5 · A-007 required 闭合）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型** / **scope**：finding-closure · A-007 F-001（web 两段登录契约）、F-002（enroll 覆盖 active）闭合验证；顺带核对 A-007 F-003 契约测试、F-004 MFA UI 残余
- **verdict**：**pass**

## 范围与区间

- **工作区**：`workspace-011-admin-functional-modules`（`workspace.md` 绑定 Root=`GOAL-001-admin-functional-modules`、`canonical_scope`、`plan_refs`/`primary_plan`=VP-011 已校验；`shared_materials_catalog: none`）。未读取或比较其他工作区目标状态。
- **已通读**：本目标 `00-meta` 信息表、D-002 §3/§4/§8、`03-audit.md`、A-005、A-007 全文。
- **代码核对**：`apps/web/src/account/auth-client.ts` `login()` L275–345；`auth-client.test.ts` L57–68；`LoginPage.tsx` L84–90 / L147–161 / L261–301；`LoginPage.test.tsx` L62–93；`AuthContext.tsx` L100–112；`modules/mfa/service.go` `Enroll` L159–194 / `Required` L75–86 / `Disable` L217–224；`store/repository.go` `UpsertPending` L84–110；`handler/mfa.go` enroll L113–126 / `writeMFAError` L237–255；`service_test.go` `TestServiceEnrollCannotOverwriteActive` L185–208；`handler/mfa_test.go` fake `Enroll` L55–57。
- **本轮复跑**（2026-08-15）：`apps/api` `go test ./internal/modules/mfa/... ./internal/handler/ -count=1` **全绿**（mfa ok 12.6s；handler ok 154.9s）。`apps/web` `vitest run src/account/auth-client.test.ts src/app/LoginPage.test.tsx` **29/29 绿**。未复跑 go `./...` 全量与 web 全量（E-004 有记录；本 scope 以关键包为准）。
- **covered**：A-007 F-001/F-002 关闭证据是否真实、充分、可重复核对；契约测试是否钉住无 token 第一段与 LoginPage 两段 UI；Enroll 三态（active 拒绝 / pending 覆盖语义 / disable 后再 enroll）；handler 映射；F-004 残余仍为 non-blocking。
- **excluded**：不改 `status` / `progress` / goal-tree / `00-meta` / D-002 正文。不重开 A-004 已闭合 required。A-007 F-005～F-007（recommended）不在本复审必改范围。
- **保证等级**：L0。

## 成果（有证据）

| A-007 主张 / 本复审核对项 | 闭合证据 |
|---------------------------|----------|
| F-001：`login()` 先识别 `{mfaRequired, mfaProof}` 再校验 token | `auth-client.ts` L324–328 在 L330–332 token 形状校验**之前**返回 `LoginMFARequired`；无 token 时不再抛 `LOGIN_MALFORMED` |
| F-001：第一段无 token 契约测试 | `auth-client.test.ts` L60–68：响应仅 `{mfaRequired:true, mfaProof:"proof-1"}` → 等值返回；`getAccessToken()` / `getRefreshToken()` 均为 `null` |
| F-001 / F-003：LoginPage 两段 UI | `LoginPage.test.tsx` L62–93：提交后 `[data-mfa-stage]` 出现 → `#mfaCode` 输入 `123456` → 点验证 → `resolveMFA` 解析为 `{code:"123456"}` → stage 消失 |
| AuthContext 第二段接线 | `AuthContext.tsx` L103–109：`isLoginMFARequired(first)` → `resolveMFA(first.mfaProof)` → `mfaVerifyRequest(...)` |
| F-002：active 再 enroll 拒绝 | `service.go` L164–167：`GetState` 且 `status=="active"` → `ErrActive`（= `handler.ErrMFAActive`） |
| F-002：store 层 WHERE 守卫 | `repository.go` L88–108：`ON CONFLICT … WHERE user_mfa.status='pending'`；`RowsAffected()==0` → `ErrActiveConflict` |
| F-002：disable 后再 enroll | `service_test.go` L200–207：`Disable` 后 `Enroll` 成功 |
| F-002：pending 覆盖语义仍保留 | `service.go` L164–170：仅 `active` 拒绝；pending / 无行走 `UpsertPending`（A-005 / D-002 §8.7） |
| F-003：契约测试缺口 | 上列两条 web 测试已补；不再作为复审 fail-closed 条件 |
| I-001～I-004 | `00-meta` 仍 verified；无到期 required 信息项；本 scope 未重开 |

## 对照 A-007 required 闭合标准

| 标准 | 状态 | 证据 |
|------|------|------|
| `login()` 在任何 token 形状校验之前识别 MFA 第一段 | **闭合** | L327 先于 L330；第一段 200 无 token 不再 `LOGIN_MALFORMED` |
| 契约测试覆盖「无 token 的第一段响应」 | **闭合** | `auth-client.test.ts` L60–68 |
| LoginPage 驱动完整两段（stage → 码 → 解析 → 完成） | **闭合**（页面契约） | `LoginPage.test.tsx` L83–92。未走 `AuthContext`→`mfaVerify` 网络链，见 F-002 recommended |
| Enroll 对 active 拒绝（不再覆盖 / 拆除 MFA） | **闭合** | `service.go` L164–167 + `store` WHERE + `TestServiceEnrollCannotOverwriteActive` L197–198 |
| pending 覆盖保留 | **闭合**（实现） | service 预检只拦 active；测试未断言 pending 覆盖，见 F-003 |
| disable 后重新 enroll | **闭合** | `service_test.go` L200–207 |
| store WHERE 与 service 一致 | **闭合**（安全不变量） | 两层都拒绝 active 覆盖；store 注释仍写 “pending/active”，见 F-003 |
| handler 映射 `400 MFA_ALREADY_ACTIVE` | **未做**（不恢复 required） | enroll 仍走 `INTERNAL` 500，见 F-001 recommended |
| 个人中心 MFA UI（A-007 F-004） | **保持 non-blocking** | `admin.account` schema 无 MFA 字段；`provider_test.go` L73–75 MFA 不贡献 page/nav/fragment |

## Findings

### F-001 · enroll handler 仍把 `ErrActive` 写成 500 INTERNAL，未映射 `MFA_ALREADY_ACTIVE`

| 字段 | 值 |
|------|-----|
| level | recommended（med） |
| status | open |
| evidence | A-007 F-002 补救清单含「handler 映射 `MFA_ALREADY_ACTIVE`」。`writeMFAError` 已有 `ErrMFAActive` → 400 `MFA_ALREADY_ACTIVE`（`handler/mfa.go` L250–251；errorcatalog / i18n / `error_contract_test.go` L51 均钉住该码）。但 `POST /api/mfa/enroll`（L119–122）仍 `writeLocalizedError(..., 500, "INTERNAL", ...)`，**不**走 `writeMFAError`。service 现已对 active 返回 `ErrActive`，该分支从死代码变成主路径：已绑定用户再 enroll 得到 500 而非 400。`service.go` L190–192 亦原样上抛 `store.ErrActiveConflict`（TOCTOU），handler 同样不认识。`mfa_test.go` fake `Enroll`（L55–57）从不返回 `ErrActive`，无 HTTP 断言。 |
| closure | — |
| 影响门禁 | **不阻断 S5** |

安全不变量已由 service 预检 + store WHERE 守住：active 行不会被覆盖，500 是 fail-closed 而非绕过。该缺口是冻结错误码契约，不是拆除 MFA。最小修复：enroll 改调 `writeMFAError`；`UpsertPending` 的 `ErrActiveConflict` 在 service 包一层为 `ErrActive`。

### F-002 · `mfaVerify` / AuthContext 两段链无客户端契约测试

| 字段 | 值 |
|------|-----|
| level | recommended（low） |
| status | open |
| evidence | A-007 F-001 曾写「第一段无 token → 不抛 MALFORMED；**再** `mfaVerify`」。现有 `auth-client.test.ts` 只钉第一段形状；`LoginPage.test.tsx` 用 stub `onLogin` 驱动 `resolveMFA`，不断言 `/api/auth/mfa/verify` 或 token 落盘。`AuthContext.tsx` L103–109 接线存在且与 D-002 §3 一致，本轮未找到反证。 |
| closure | — |

不升级。第一段断裂是 A-007 的 fail-closed 条件，现已钉住。

### F-003 · pending 覆盖无回归断言；store 注释与守卫不一致

| 字段 | 值 |
|------|-----|
| level | recommended（low） |
| status | open |
| evidence | 用户要求的三态：active 拒绝、pending 覆盖、disable 后再 enroll。`TestServiceEnrollCannotOverwriteActive` 覆盖 active 拒绝与 disable 后再 enroll，**没有** pending 再 `Enroll` 成功的断言。实现允许 pending 覆盖（`service.go` L164–170）。`repository.go` L84–85 注释仍写「overwrites any previous pending/**active** row」，与 L97 `WHERE status='pending'` 相反。`store` 包无单测。 |
| closure | — |

实现与 A-005 / D-002 §8.7 一致；缺口在测试与注释。

### F-004 · 个人中心 MFA 管理区块 UI 仍未实现（A-007 F-004）

| 字段 | 值 |
|------|-----|
| level | recommended（med）· **non-blocking** |
| status | open |
| evidence | 与 A-007 F-004 相同：`modules/account` schema 无 MFA 字段；MFA provider 不贡献 page/nav/fragment。自服务 API 仍完整。users `mfaEnabled` + Reset MFA + i18n 仍在。 |
| closure | — |
| 影响门禁 | **不阻断 S5** |

**裁决（本独立审，维持 A-007）**：不把个人中心 MFA UI 标为 required，不因该缺口单独阻断关门。

建议 `/govern` 走 P-003 `accepted-residual`（或后续 UI 目标），范围与复审触发与 A-007 F-004 相同：

| 项 | 值 |
|----|-----|
| 接受范围 | 个人中心 enroll/confirm/disable/rotate **UI**（API 不豁免） |
| 缓解 | 自服务 API + admin Reset MFA |
| 复审触发 | schema renderer 支持非 reload 载荷 **或** 用户要求自助绑定入管理端 |

## 必改项汇总

无 required。A-007 F-001、F-002（required）已合法闭合（fixed）。

A-007 F-002 原文「handler 映射 `MFA_ALREADY_ACTIVE`」未落地，本意见**不**将其恢复为 required：覆盖/接管路径已在两层拒绝，HTTP 500 为 fail-closed 契约缺口（本 F-001 recommended）。

## 与既有意见的异同

- A-007 independent fail：开放 F-001（high）+ F-002（med）。本意见核对修复后：**两条 required 均闭合**。
- A-007 F-003（缺契约测试）已被本轮测试覆盖，不再作为 fail-closed 条件；`mfaVerify` 链仍无单测（本 F-002 recommended）。
- A-007 F-004 MFA UI：维持 non-blocking / 建议 `accepted-residual`。
- A-007 F-005～F-007（confirm `last_used_step`、`Required` fail-open、jwt 轮换）仍开放 recommended，本复审不升级、不阻断关门。
- A-005 F-001 recommended（active 须先 disable）已由 service + store 落地；pending 覆盖实现保留。
- A-004 F-001/F-002（required）保持闭合。不重开。
- 不与 A-001/A-002 立项 pass、A-003/A-005 S1 pass、A-006 self pass 冲突。

## 结论 + 建议给编排器/用户的下一步

**verdict: pass**。A-007 F-001 / F-002 关闭证据充分、可重复核对：`login()` 先读 MFA 第一段；契约测试钉住无 token 响应与 LoginPage 两段 UI；active enroll 在 service 与 store 两层拒绝；disable 后再 enroll 可行。无 high/med required；无到期 required 信息项。

**A-007 required 已合法闭合，可关门**（`status: done` 由 `/govern` 执行）。本意见不改 status。

建议 `/govern`：

1. 记录 A-007 F-001 / F-002 → **fixed**（证据 = 本 A-008）。
2. F-004 书面 `accepted-residual`（范围/缓解/复审触发见上表）或开后续 UI 目标。
3. 本 F-001 recommended（enroll → `writeMFAError`）建议随手修，不阻断关门。
4. 执行 S5 关门：`00-meta` / goal-tree `status: done`、progress 按检查点重算。勿用既有 `progress: 3/5` 作为放行依据。

## 声明

本意见不修改 `status` / `progress` / goal-tree / `00-meta` / D-002 正文。响应由 `/govern` 处理。保证等级 L0。
