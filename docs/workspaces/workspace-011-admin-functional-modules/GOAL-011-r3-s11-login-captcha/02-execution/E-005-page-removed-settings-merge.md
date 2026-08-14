---
id: E-005
goal: GOAL-011-r3-s11-login-captcha
date: 2026-08-14
status: recorded
parent: GOAL-011-r3-s11-login-captcha
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# E-005 · 页面移除 + 开关并入设置页

## 事实

- 2026-08-14（用户书面裁决 D-003）：
  - 删除 `modules/logincaptcha/schema/`（captcha.json + schema.go）与 `manifest/`（fragment.json + manifest.go）；provider.go 移除 Schema/Navigation/Manifest 注册与 Descriptor 的 Pages/Navigation/Fragments。
  - kernel/profile.go BuiltinModules 描述符同步；testsupport 移除 menu_captcha；composition_test admin 导航 13→12（权限 24 不变）；provider_test 页面/导航/fragment 断言 → 0。
  - settings.json 新增 `settings-security` 表单（enabled select → GET/PATCH /api/captcha/settings）；i18n +schema.settings.toolbar.security / field.requireCaptcha / option.captchaEnabled/Disabled。
  - Web：fixture 移除 captcha 页/导航并重钉 STATIC_MANIFEST_SHA256（b45a3e02…）；app-manifest.test / schema-keys / s5 列表移除；i18n 移除 manifest.title.captcha / manifest.nav.captcha / schema.captcha.*（保留 error.invalidCaptcha、login.captcha*）；smoke admin required_pages 移除 captcha；e2e shell.spec 移除 Login captcha 断言。
- 回归：go ./... 与 vitest 全量（随本记录后续提交确认）。
