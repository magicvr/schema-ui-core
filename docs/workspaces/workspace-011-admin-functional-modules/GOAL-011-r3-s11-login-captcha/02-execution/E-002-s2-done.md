---
id: E-002
goal: GOAL-011-r3-s11-login-captcha
date: 2026-08-14
status: recorded
parent: GOAL-011-r3-s11-login-captcha
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# E-002 · S2 实现完成

## 事实

- handler/auth.go：登录链路新增验证码步骤（限流之后、凭据校验之前；h.captcha 为 nil 或 Required()==false 时行为与改造前逐字节一致）；credentials 增加 captchaId/captchaAnswer；authsHandler 增加 captcha CaptchaVerifier 参数。
- handler/health.go：RegisterWithReadiness 增加变参 captcha ...CaptchaVerifier 并传递给 authsHandler（composition 在模块启用时传入服务，禁用时传入 nil → 登录契约不变）。
- handler/captcha.go（新增）：GET /api/auth/captcha（公开预检，{enabled, challenge?}）、GET/PATCH /api/captcha/settings（captcha.read/captcha.write 权限，PATCH 写审计事件 captcha.settings-update）。
- modules/logincaptcha：challenge.go 重构（Generate 返回 (id, question, expiresInSeconds, err) 以结构满足 handler.CaptchaService；新增 SetEnabled）；provider.go（Descriptor + Register：3 路由 / 1 页 / 2 权限 / menu_captcha Order 7 / fragment）；schema/captcha.json（settings 表单页）；manifest/fragment.json。
- migration 0023 login_captcha（checksum 6bc0b556…）；operationlog 0024 operation_log_captcha（CHECK + captcha.settings-update，checksum 51d6a3c1…）；compiled/persistence.go 注册 logincaptcha migration provider。
- 接线：kernel/profile.go（ProfileAdmin 默认集 + BuiltinModules 描述符）；composition.go（单例 captchaService 同时供 verifier 与 provider）；testsupport/store.go（captcha.read/write + menu_captcha 镜像）。
- errorcatalog + error_contract_test + i18n（error.invalidCaptcha）；composition_test admin 22 权限 / 12 导航；store 迁移计数 22→24。
- 测试：store/repository_test.go、challenge_test.go、provider_test.go、handler/captcha_test.go（预检、设置、登录门默认关/开-错/开-对、审计）。
- Web：fixture app-manifest.admin.json +captcha 页/导航，STATIC_MANIFEST_SHA256 重钉 95e571e5…，i18n en/zh +8 键，schema-keys/s5 列表 + smoke.sh admin required_pages + shell.spec admin/mvp 断言。
