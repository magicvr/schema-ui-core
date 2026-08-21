---
id: D-002
goal: GOAL-017-r3-s10-mfa-2fa
title: 方案冻结：MFA / 2FA（TOTP）设计（S1）
date: 2026-08-15
status: accepted
parent: GOAL-017-r3-s10-mfa-2fa
created: 2026-08-15
updated: 2026-08-15
version: 1.0.0
---

# D-002 · 方案冻结（S-10 MFA/2FA · TOTP）

> 依据：I-011-001 §4 S-10；GOAL-017 00-meta 边界与 I-001~I-004；A-002 017-F-003（S-11 合成裁定）。
> 证据：handler/auth.go（登录链）、internal/auth/auth.go（Login 核心）、handler/account_self.go（自服务先例）、modules/logincaptcha（一次性挑战先例）、I-HOST-APP-001 AUTH-004/AUTH-006（协议对照）。

## 1. 因子模型（I-001 闭合）

- **TOTP（RFC 6238，自实现）**：HMAC-SHA1、30 秒周期、6 位码、时间步 counter = floor(now/30)；校验窗口 ±1 步（3 个候选）；secret = 20 字节随机（RFC 4226 建议 160-bit），Base32 编码展示。
- **自实现而非第三方库**：stdlib crypto/hmac + sha1 约 40 行；仓库纪律先例（GOAL-010 D-002 §7 否决 robfig/cron，D-001 §5）。RFC 6238 附录 B 官方测试向量必过。
- **恢复码**：10 个一次性随机码（8 位、字母表 23456789ABCDEFGHJKMNPQRSTUVWXYZ 去易混字符），**bcrypt 哈希存储**（复用既有 bcrypt 基建）；一次性消费（命中即删除该哈希）；轮换 = 重新生成整组覆盖。
- **解绑**：需当前 TOTP 码或恢复码二次验证；解绑后 token_version+1 + 吊销全部 refresh（强制重新登录，安全先例 C-10/F-03）。

## 2. 数据模型（migration 0029 admin.mfa）

- user_mfa：user_id TEXT PK REFERENCES users(id) ｜ status TEXT NOT NULL CHECK in (pending, active) ｜ totp_secret_ciphertext TEXT NOT NULL（AES-256-GCM，nonce 12B 前缀）｜ recovery_codes_hash TEXT NOT NULL（JSON 数组，bcrypt）｜ last_used_step INTEGER（最近成功消费的 TOTP 时间步，同窗重放拒绝）｜ created_at / updated_at INTEGER NOT NULL——**仅 status=active 触发登录 MFA**（enroll 后 pending，confirm 后 active；pending 不触发——防 enroll 后会话过期自锁，A-004 F-001 响应）。
- mfa_proofs：id TEXT PK ｜ user_id TEXT NOT NULL ｜ fail_count INTEGER NOT NULL DEFAULT 0 ｜ expires_at INTEGER NOT NULL ｜ created_at INTEGER NOT NULL——登录第二因素一次性证明（5 分钟过期，校验成功即删除；logincaptcha 挑战表先例）；fail_count 达 5 → proof 立即失效（防爆破，A-004 F-001 响应）。
- **secret 加密**：AES-256-GCM；密钥 = HKDF-SHA256(server secret（与 jwt 签名同源 config 密钥）, context "mfa/totp") 派生；无既有静态加密基建（仅 bcrypt/SHA-256 + jwt_secret），故复用 server secret 派生，留痕（独立 KMS/密钥面归后续安全波次）。
- operationlog CHECK 扩展：**0030**（mfa.enroll / mfa.confirm / mfa.disable / mfa.recovery-rotate / mfa.admin-reset / mfa.login）。

## 3. 登录集成（I-002 闭合 · 与 S-11 并存裁定 I-004）

- **MFAVerifier 接口**（handler/auth.go 同 CaptchaVerifier 模式）：Required(userID string) bool；nil = 未启用，登录契约逐字节不变（GOAL-011 先例）。
- **插入点**：Authenticator.Login 密码校验成功后、issue() 前（auth 核心增加可选 mfa 字段；nil 原行为不变）。Required==true → 返回哨兵 ErrMFARequired，**不签发凭据**。
- **两段登录**：第一段 POST /api/auth/login 密码段成功且需 MFA → 200 {mfaRequired: true, mfaProof: "<proof id>"}；第二段 POST /api/auth/mfa/verify {proof, code, recoveryCode?} → 校验 proof（存在/未过期/一次性删除）+ TOTP 或恢复码（二选一）→ 成功签发正式 access/refresh（提取 IssueTokensFor(userID) 复用签发路径）→ 审计 mfa.login。
- **失败计数分轨**（017-F-003 / I-004 闭合）：链路 = 限流 → S-11 验证码 → 密码 → **TOTP**；验证码失败不计锁定（GOAL-011 语义）；MFA 失败**不计入账号锁定**（锁定针对凭据，先例保留），但同一 proof 连续失败 5 次（fail_count 落库）→ 该 proof 立即失效（防爆破，需重新登录取新 proof）。恢复码失败同 TOTP 计数。
- **改密/停用联动**：改密（token_version+1 机制）后已签发 access 失效但 mfa 绑定保留；停用账号期间 mfa 校验不触发（enabled=0 早于 mfa 检查）。

## 4. 端点与权限（I-003 闭合）

| 端点 | 门禁 | 说明 |
|------|------|------|
| POST /api/auth/mfa/verify | public（持 proof） | 第二因素验证 → 签发凭据；MFA_INVALID 401 / MFA_PROOF_EXPIRED 401 |
| GET /api/mfa/status | 身份即可（自服务先例 account_self） | {enabled, enrolledAt} |
| POST /api/mfa/enroll | 身份即可 | 生成 secret + 恢复码 → {secretBase32, otpauthURL, recoveryCodes}（**仅此一次可见**）；行状态 pending 至 confirm |
| POST /api/mfa/confirm | 身份即可 | {code} 校验 TOTP → 激活 |
| POST /api/mfa/disable | 身份即可 | {code 或 recoveryCode} 二次验证 → 删行 + token_version+1 + 吊销全部 refresh；审计 mfa.disable |
| POST /api/mfa/recovery/rotate | 身份即可 | {code 或 recoveryCode} 验证 → 新恢复码整组覆盖 |
| POST /api/users/{id}/mfa/reset | users.mfa-reset（PolicyAdmin） | 管理员重置用户 MFA（解绑）+ token_version+1 + 吊销全部 refresh（与自助 disable 同级，A-004 F-002 响应）；审计 mfa.admin-reset |

- 权限键：users.mfa-reset（admin 命名空间，F-03 users.enable/disable 先例）；自服务端点无权限键（account_self.go 模式）。
- 页面：**个人中心（admin.account）区块**——MFA 状态卡 + 绑定流程（enroll → 展示 secret/恢复码 → confirm）+ 解绑 + 恢复码轮换；users 页行扩展 mfaEnabled 计算字段 + 行操作 Reset MFA（permissions.users.mfa-reset 表达式，fail-open 视觉，服务端 fail-closed）。无新导航项。
- Profile：admin.mfa 进入 **admin 默认集**（内容扩展先例 F-03/S-01/S-02）；mvp/demo 不启用；**默认无人绑定 = 登录行为零变化**（冒烟/e2e 回归零影响，GOAL-011 默认关闭先例）。

## 5. 协议对照（本地扩展，不声明协议覆盖）

- AUTH-004 explicitly-out：login/logout/revoke command 属身份服务，页面协议不携带凭据。
- AUTH-006 reserve-extension（上游 deferred）：MFA/step-up challenge 需独立安全 profile——协议未定义 MFA 语义。
- 处置：**本地身份层扩展**（auth 核心 + admin 模块），不新增协议 capability、不改协议 pin（v2.8.0）、不改 Manifest 装配语义；web 侧为本地扩展交互（LoginPage 两段登录 + auth-client 扩展），呈现自由不适用（非页面协议语义）；留痕于本决策与 03-audit。

## 6. 测试与验证

- TOTP：RFC 6238 附录 B 官方向量（SHA1/8 位与 6 位推导）+ 时钟漂移窗口 + 重放（同窗口二次校验拒绝——proof 单次 + 码滑动窗口记录 lastUsedStep 拒绝同窗重用）。
- auth：Login 无 mfa 时逐字节不变（nil 门测试）；绑定用户两段登录全链路（正确/错误码/恢复码/过期 proof/重放 proof/5 次失效）。
- 端点：enroll/confirm/disable/recovery-rotate 生命周期 + 权限（401/403）+ 审计事件。
- 组合根：admin 权限 **24→25**（+users.mfa-reset）、导航 **12→12**（composition_test.go L465 当前 24/12）；迁移 26→30。
- web：LoginPage 两段交互单测、个人中心 MFA 区块、users 行操作；fixture/schema-keys/s5-denominator/e2e 回归（默认无人绑定零影响）+ smoke.sh 页面集不变。

## 7. 未选方案（留痕）

- 第三方 OTP 库（pquerna/otp）：新依赖 + 仓库纪律不符（D-001 §5 先例），自实现 + RFC 向量测试。
- 短信/邮件 OTP：外部通道 + B-09 模板基建依赖，v1 否决（TOTP 无通道依赖）。
- WebAuthn/FIDO2：多因子形态超出 v1 单因子 TOTP；AUTH-006 独立安全 profile 语境后续立项。
- 会话内 step-up（敏感操作二次认证）：v1 仅登录时第二因素；step-up 属 AUTH-006 后续 profile，登记不实现。
- MFA 失败计入账号锁定：锁定语义保持针对凭据（先例 GOAL-011 文档化），改独立 proof 失效机制。

## 8. S2 实现清单

1. auth 核心：MFAVerifier 可选字段 + ErrMFARequired + IssueTokensFor 提取；mfaProof 签发/校验。
2. handler：/api/auth/mfa/verify + /api/mfa/* 端点（自服务 + admin reset）。
3. modules/mfa：provider（Descriptor + 路由 + 权限 + fragment）+ totp.go（RFC 6238 自实现）+ store（user_mfa / mfa_proofs）+ schema（个人中心区块动作）。
4. migration 0029（admin.mfa）+ 0030（operationlog CHECK）；compiled/persistence.go 注册。
5. web：LoginPage 两段登录 + auth-client 扩展 + 个人中心 MFA 区块 + users 行操作 + i18n。
6. 测试：RFC 向量 / auth 门 / 端点 / 组合根 / web fixture。
7. A-005 recommended：pending 状态下重复 enroll 覆盖旧 pending；管理员强制启用 v1 不裁定（留扩展位）；mfaProof 回传路径 = /api/auth/mfa/verify 响应体（与 login 同形状 accessToken/refreshToken/user）；迁移 0029/0030 为两目标合计（016 占 0027/0028）。
