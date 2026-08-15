---
id: E-003
goal: GOAL-017-r3-s10-mfa-2fa
title: S2 实现完成（MFA/2FA）
date: 2026-08-15
status: recorded
parent: GOAL-017-r3-s10-mfa-2fa
created: 2026-08-15
updated: 2026-08-15
version: 1.0.0
---

# E-003 · S2 实现完成（2026-08-15）

## 事实

- **auth 核心**：Authenticator.mfa 可选门（MFAEnforcer，nil=字节不变）+ MFARequiredError{UserID}（登录返回、不签发凭据）+ IssueTokensFor（第二因素通过后签发）+ SetMFAEnforcer。
- **迁移 0029**（admin.mfa）：user_mfa（status pending/active、secret 密文、恢复码哈希、last_used_step）+ mfa_proofs（fail_count）；checksum daa2592f…；**0030**（core.operationlog）CHECK + 6 个 mfa.* 事件；checksum 70464abf…；compiled 注册。
- **TOTP 自实现**（totp.go）：RFC 6238 HMAC-SHA1/30s/6 位/±1 窗，官方附录 B 向量锁定；AES-256-GCM 加密 secret（HKDF(server secret, "mfa/totp")）。
- **service**：Required（仅 active 触发，无自锁）/ BeginChallenge（proof 5min 一次性）/ Verify（proof 存在性/过期/耗尽、TOTP 同窗重放拒绝、恢复码 bcrypt 一次性消费、5 次失败销毁 proof）/ Enroll/Confirm/Disable/RotateRecovery/AdminReset。
- **handler**：handler/mfa.go MFARoutes（POST /api/auth/mfa/verify + GET /api/mfa/status + enroll/confirm/disable/recovery-rotate + POST /api/users/{id}/mfa/reset）；auth.go 登录链 MFA 分支（密码成功后 proof，无 token）；SessionRevoker（authsession.BumpTokenVersionAndRevokeAll——disable/reset 同强度吊销，A-004 F-002）。
- **装配**：composition（mfaService + **接口类型 nil** 防 typed-nil 陷阱 + RegisterWithMFA）、kernel（ProfileAdmin + BuiltinModules）、testsupport 镜像（users.mfa-reset）；组合根 26→27 权限 / 13 导航（无新导航）。
- **web**：LoginPage 两段登录（code/recovery 输入级，onLogin 增 resolveMFA）、auth-client（LoginMFARequired + mfaVerify）、AuthContext 两段流程；users.json 增 mfaEnabled 列 + Reset MFA 行操作（users.mfa-reset 权限表达式 + disabledWhen）+ users API mfaEnabled 计算字段（跨模块 EXISTS 读 user_mfa）；i18n zh/en 30 键。
- **错误码**：INVALID_MFA_BODY / MFA_INVALID / MFA_PROOF_EXPIRED / MFA_PROOF_EXHAUSTED / MFA_NOT_ENROLLED / MFA_PENDING_ONLY / MFA_ALREADY_ACTIVE（errorcatalog + 契约钉住）。
- **登记残余（用户裁决点）**：个人中心（admin.account）MFA 管理区块 UI 未实现——renderer 仅支持 reload 行为、无载荷展示，enroll 的 secret/恢复码一次性展示需自定义组件；v1 web 交付 = 两段登录 + users Reset MFA + i18n，API 面完整且有测试覆盖。
