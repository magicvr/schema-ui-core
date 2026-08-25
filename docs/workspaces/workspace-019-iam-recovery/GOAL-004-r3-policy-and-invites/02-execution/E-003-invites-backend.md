---
id: E-003
doc: execution-entry
goal: GOAL-004-r3-policy-and-invites
status: recorded
created: 2026-08-25
updated: 2026-08-25
version: 1.0.0
---

# E-003 · C3 满：邀请域 + 管理/激活 API + settings 策略配置 API

2026-08-25 完成：

- 邀请域 authsession/invites.go：CreateInvite（角色≥1 且存在校验、令牌仅存 sha256）、ListInvites、RevokeInvite（即时失效）、ResendInvite（撤旧发新 + 60s 冷却 + 过期可重发）、AcceptInvite（单事务：live 校验→角色 revalidate INVITE_ROLE_GONE→用户名冲突 fail-closed→建号带角色+消费邀请；不签发会话）。
- handler/invites.go：管理四路（admin.users 模块，users.invite 权限经 PermissionsForUser 校验）+ 公开 POST /api/auth/invite/accept（中央注册）；带邮箱时经 MailSender 发邀请信（请求 Host 推导链接基址），发送失败补偿撤销。
- settings 配置面：GET/PATCH /api/settings/password-policy（settings.write 权限；范围校验 length 8–72/categories 0–4/depth 0–10）。
- kernel/profile.go 内置模块声明同步（admin.users +4 路由+users.invite；admin.settings +2 路由）；composition 装配更新；错误目录新增 INVALID_INVITE_BODY / INVITE_INVALID / INVITE_ROLE_GONE 并入契约冻结集；composition reconcile 黄金计数 mvp 10→11、admin 32→33。
- 测试：invites/policy 域 5 组全绿；全 internal 包 `-p 1` 全量 sweep 无 FAIL。

后续：C4 Web 面 → C5 independent + self 关门。
