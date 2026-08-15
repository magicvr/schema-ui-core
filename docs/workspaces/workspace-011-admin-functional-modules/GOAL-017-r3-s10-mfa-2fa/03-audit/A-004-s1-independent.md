---
id: A-004
goal: GOAL-017-r3-s10-mfa-2fa
title: S1 方案冻结独立审计（MFA/2FA · TOTP）
date: 2026-08-15
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
scope: S1 方案冻结（D-002 + I-001~I-004 闭合 + 登录集成/安全控制/协议/迁移数字）
audit_type: design-plan
verdict: conditional
status: recorded
parent: GOAL-017-r3-s10-mfa-2fa
created: 2026-08-15
updated: 2026-08-15
version: 1.0.0
---

# A-004 · 独立审计意见（S1 方案冻结 · S-10 MFA/2FA）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型** / **scope**：design-plan · S1 方案冻结（D-002 全文、I-001~I-004 证据、登录集成点、TOTP/proof/恢复码/分轨、协议对照、迁移/组合根、无越界）
- **verdict**：**conditional**

## 范围与区间

- **工作区**：`workspace-011-admin-functional-modules`（绑定与 `plan_refs`/`primary_plan` 已校验；`shared_materials_catalog: none`）。未读取或比较其他工作区目标状态。GOAL-011 仅作登录挑战先例只读对照，不审其状态。
- **已通读**：本目标 `00-meta`、`01-decision.md`、D-001、D-002、E-002、`03-audit.md`、A-001～A-003；I-011-001 §4 S-10/S-11；I-HOST-APP-001 AUTH-004/AUTH-006。
- **代码核对**：`handler/auth.go`（限流 → captcha → `Login`）、`internal/auth/auth.go`（`Login` / `issue` / 锁定 / disabled）、`handler/account_self.go`（自服务先例）、`modules/logincaptcha`（一次性挑战 + 5 分钟 TTL）、`config/config.go`（`AuthJWTSecret`）、`composition/composition_test.go`（24/12）、`operationlog` 迁移 max 26。
- **covered**：方案可实施性、security 控制完备性、信息项闭合、协议对照、迁移/组合根、无越界、未选方案与 S2 清单。
- **excluded**：S2 实现、S3～S5、不改 status/progress/goal-tree/方案正文。
- **保证等级**：L0。

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| 登录插入点与代码一致：密码成功后、`issue()` 前 | `auth.go` `Login` L123–145：验密 → 可选 `ResetLoginFailures` → `return a.issue(u, now)`（L145）。MFA 插在 L145 之前正确；`Refresh` L185 再 `issue` 不经 MFA，与「仅登录时第二因素」一致 |
| 链路顺序限流 → S-11 captcha → 密码 → TOTP 与现网一致 | `handler/auth.go` L85–88 限流、L89–98 captcha（失败不计锁定，L91–92）、L99 `Login`（锁/禁用先于验密，`auth.go` L114–122） |
| `MFAVerifier` nil = 契约不变，同 `CaptchaVerifier` 可选门 | `handler/auth.go` L16–23、L93；D-002 §3。`Required(userID)` 比 captcha 的 `Required()` 多 user 维，因 MFA 在验密后才知道用户——适配正确 |
| I-004 分轨：验证码失败不计锁定；MFA 失败不计入账号锁定 | GOAL-011 语义 + D-002 §3；锁定只在 `VerifyPassword` 失败走 `RecordLoginFailure`（`auth.go` L123–138） |
| TOTP 参数可实施：HMAC-SHA1 / 30s / 6 位 / ±1 窗 / 20 字节 secret；RFC 6238 附录 B 必过 | D-002 §1；stdlib `crypto/hmac`+`sha1` 路径成立 |
| 恢复码 bcrypt 一次性、整组轮换 | D-002 §1；仓库已有 bcrypt（`auth.go` import `golang.org/x/crypto/bcrypt`） |
| 解绑：二次验证 + `token_version+1` + 吊销 refresh，与 C-10/F-03 先例同构 | D-002 §1；`account_self.go` L5、L198 注释 |
| 自服务端点无权限键，先例 `AccountSelfRoutes` | `account_self.go` L2–4、L53–57 |
| 协议对照诚实：AUTH-004 `explicitly-out`、AUTH-006 `reserve-extension`；不声称协议覆盖 | I-HOST-APP-001 L64–66、L197–199；D-002 §5 本地身份层扩展、不改 pin v2.8.0 |
| 迁移预留 0029/0030 在 016 的 0027/0028 之后、相对 max 26 不碰撞 | `operationlog` `Version: 26`；两目标编号互斥 |
| 默认无人绑定 = 登录零变化；admin 内容扩展、mvp/demo 不启用 | D-002 §4；与 S-11 默认关闭同构 |
| 无越界：不改锁定/限流/改密语义，不改协议 pin / Manifest 装配 | D-001 L20；D-002 §3/§5 |
| 未选方案覆盖第三方库 / 短信邮件 / WebAuthn / step-up / MFA 计入锁定 | D-002 §7 |

## 对照成功标准（S1）

| 标准 | 状态 | 证据 |
|------|------|------|
| 因子与恢复流程（TOTP / 恢复码 / 解绑）可实施 | 部分 | 算法与存储方向成立；`user_mfa` 状态机与声明控制不一致（F-001） |
| 登录挑战集成点与会话/失效语义 | 部分 | 插入点正确；proof 传递与 admin reset 会话失效不完整（F-001/F-002/F-004） |
| 与 S-11 合成（I-004） | 满足 | 串行并存 + 分轨已写死，且与代码顺序一致 |
| 协议面不声称 MFA 已纳入页面协议 | 满足 | AUTH-004/006 → 本地扩展 |
| I-001/I-002 required 设计层闭合 | 满足（设计层） | 未把未实现 TOTP 写成已落地 |
| I-003 管理面含「强制启用候选」 | 部分 | 仅有 admin reset，强制启用未裁定（F-003） |

## Findings

### F-001 · 数据模型相对已声明安全控制不完整且自相矛盾

| 字段 | 值 |
|------|-----|
| level | **required**（med） |
| status | open |
| evidence | 三处不一致：① D-002 §2「行存在即启用」，§4 enroll「行状态 pending 至 confirm」，表无 `status`/`pending` 列。enroll 需已登录（§4），若写行即 `Required==true` 且会话过期未 confirm → 下一登录要 MFA、confirm 又要身份，**自锁**。② §3「同一 proof 连续失败 5 次失效」，`mfa_proofs` 列仅为 id/user_id/expires_at/created_at，**无 fail_count**。对比 logincaptcha：`ConsumeChallenge` 任意尝试即删（`store/repository.go` L50–74），与「5 次」不同，不能照搬。③ §6 测试要求 `lastUsedStep` 拒同窗重用；`user_mfa` 无该列。两段并行登录可发两个 proof，同一 TOTP 码可成功两次。 |
| closure | — |
| 影响门禁 | S1 方案冻结 / S2 迁移 0029 |

S2 写 0029 前必须在 D-002 补列并消解矛盾，至少：

- `user_mfa.status`（pending/active）或等价；`Required()` 仅 active；过期 pending 可被新 enroll 覆盖
- `mfa_proofs.fail_count`（或等价原子计数）；≥5 删除/作废
- `user_mfa.last_used_step`（成功校验后写入；同窗重放拒绝）
- proof `id` 熵（建议 ≥128-bit 随机），防枚举

### F-002 · 管理员重置 MFA 未规定会话失效，弱于自助解绑

| 字段 | 值 |
|------|-----|
| level | **required**（med） |
| status | open |
| evidence | 自助 disable：`删行 + token_version+1 + 吊销全部 refresh`（D-002 §1/§4）。admin reset：`POST /api/users/{id}/mfa/reset` 仅「解绑」+ 审计 `mfa.admin-reset`（§4 表）。F-03 停用/改密均 bump `token_version`（`auth.go` L415–426 注释；`account_self.go` L5）。 |
| closure | — |
| 影响门禁 | S1 / S2 安全语义 |

若重置原因是失窃或强制解绑，旧 access/refresh 在 TTL 内仍可用。S1 须规定 reset **至少**与 disable 同等：`token_version+1` + `RevokeAllRefreshTokensForUser`。

### F-003 · I-003「管理员强制启用」候选未裁定即标 verified

| 字段 | 值 |
|------|-----|
| level | recommended（med） |
| status | open |
| evidence | `00-meta` L24「管理员视角（强制启用候选，S1 冻结）」；I-003「自助入口 + 管理员强制启用候选」→ verified / D-002 §4。D-002 §4 只有 status/enroll/confirm/disable/rotate + admin **reset**。§7 未选方案未点名强制启用。 |
| closure | — |

I-003 为 non-blocking，不单独阻断 S1。S2 前应补一句：本波不纳入强制启用（策略/登录拦截），仅自助 + admin reset。

### F-004 · 第一段 `mfaProof` 如何从 `Login` 回到 HTTP 未写死

| 字段 | 值 |
|------|-----|
| level | recommended（med） |
| status | open |
| evidence | `Authenticator.Login` 现签名 `(access, refresh, user, err)`；错误时 user 为空（`auth.go` L105–145）。D-002 §3：`ErrMFARequired` 且 200 `{mfaRequired, mfaProof}`。未写：错误类型携带 proof、改签名、或 handler 在已知 userID 后签发 proof。auth 核心若直接依赖 `admin.mfa` store 会层破坏。 |
| closure | — |

建议：handler/mfa 服务在验密成功后签发 proof；`Login` 返回可识别哨兵 + userID（或带 proof 的 error）。S2 清单第 1 条补这一句即可。

### F-005 · 组合根数字：权限 22→24 与 +1 键不符；26→30 为两目标合计

| 字段 | 值 |
|------|-----|
| level | recommended（low） |
| status | open |
| evidence | 实测 admin **24** 权限 / **12** 导航（`composition_test.go` L465）。017 只增 `users.mfa-reset` → 应为 **24→25**（若 016 先落地则 26→27）。导航 12→12 正确。迁移 0029/0030 预留正确；「26→30」是 016+017 合计，017 单独是 26→28 或 28→30。jwt 签名密钥存在（`config.go` `AuthJWTSecret`）；HKDF context `"mfa/totp"` 可接受为已文档化残余，但 **jwt_secret 轮换会使全部 TOTP secret 不可解密**，S2 须写轮换程序或接受该残余。 |
| closure | — |

## 必改项汇总

1. **F-001（required · med）**：补 pending / fail_count / last_used_step（及 proof 熵），消除 §2 与 §4 矛盾。
2. **F-002（required · med）**：admin reset 必须 `token_version+1` + 吊销全部 refresh。

未闭合前**不可无条件放行 S2**（尤其不可先写 0029）。无 high required；I-001/I-002 设计层证据充分。

## 与既有意见的异同

- A-003（self · pass）认为可进独立审与 S2。本意见同意：插入点、S-11 合成、协议本地扩展、nil 门、恢复码方向、无越界成立。
- 不同意「可无条件进 S2」：self 未核对 schema 与声明控制是否同构（F-001）、reset 与 disable 是否对称（F-002）。
- A-002 立项 F-003（S-11 合成）已由 I-004 + D-002 §3 闭合。
- 不与 A-001/A-002 立项 pass 冲突。

## 结论 + 建议给编排器/用户的下一步

**verdict: conditional**。登录集成点与 S-11 分轨可核对；TOTP 自实现 + RFC 向量 + AES-GCM/HKDF 方向可接受（密钥面已留痕）。**开放 F-001、F-002 required**。

建议 `/govern`：

1. 修正 D-002 §2 表结构与 §4 reset 语义（F-001/F-002），不必重开 S1 信息收集。
2. F-003～F-005 一并勘误。
3. 修正落盘后再开 S2。勿用 `progress: 1/5` 放行。

## 声明

本意见不修改 `status` / `progress` / goal-tree / 方案正文。响应由 `/govern` 处理。保证等级 L0。
