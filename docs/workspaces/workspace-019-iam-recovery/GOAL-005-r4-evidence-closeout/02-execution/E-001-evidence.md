---
id: E-001
doc: execution-entry
goal: GOAL-005-r4-evidence-closeout
status: recorded
created: 2026-08-25
updated: 2026-08-25
version: 1.0.0
---

# E-001 · R4 证据入账（C1/C2 满）

2026-08-25 完成：

- handler 测试环境补挂中央面（RegisterRecovery/RegisterInviteAccept/InviteAdminRoutes/PasswordPolicyRoutes/EmailIdentityRoutes），全部走真实 mock 渠道（OutboxSink）。
- 新增 r4_evidence_test.go：恢复链（bind→verify→start→渠道取码→complete→新密码登录 200）；邀请链（admin 建邀（users.invite 权限已入 testsupport seed）→ 渠道取链接 → accept 204 → 受邀者以 viewer 角色登录 → 一次性回放 INVITE_INVALID）；策略链（PATCH minLength=12 → 8 字节创建 400 INVALID_PASSWORD → 强密码创建成功）。
- Go 全包测试绿（handler 35s 全过）。
