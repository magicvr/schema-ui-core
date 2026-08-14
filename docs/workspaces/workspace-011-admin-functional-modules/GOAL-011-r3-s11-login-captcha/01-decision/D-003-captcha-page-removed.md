---
id: D-003
goal: GOAL-011-r3-s11-login-captcha
title: 用户裁决：删除 captcha 独立页面，开关并入系统设置页
date: 2026-08-14
status: accepted
parent: GOAL-011-r3-s11-login-captcha
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# D-003 · 用户裁决（S-11 页面变更）

## 裁决

用户（直接书面指示，2026-08-14）：「不用排查这个页面了，删掉它。把验证码开关合并到系统设置页。」

- **删除**：`menu_captcha` 导航项、`captcha` 页面（schema 文档）、manifest fragment。
- **保留**：登录门禁、公开预检 `GET /api/auth/captcha`、设置端点 `GET/PATCH /api/captcha/settings`、权限键 `captcha.read/write`、审计事件 `captcha.settings-update`、迁移 0023/0024。
- **合并**：开关（enabled select）移入 `admin.settings` 设置页的新「Security/安全」区（settings.json 新增 `settings-security` 表单，recordSource `GET /api/captcha/settings`，submitAction `updateCaptcha` → `PATCH /api/captcha/settings`）。

## 理由（用户视角）

- 独立页面仅承载一个开关，功能单薄且当前渲染报错；开关属于站点级安全配置，归入系统设置页更自然。
- 权限说明：设置页本身 admin-only（settings.read）；表单 PATCH 实际走 `captcha.write`，admin 默认集同时持有两键，无权限断档。

## 影响

- admin 权限计数不变（24），导航 13→12（移除 menu_captcha）。
- 模块 `admin.login-captcha` 不再贡献 Pages/Navigation/Fragments；仅 Routes + Permissions。
- Web fixture/i18n/smoke/e2e 断言同步移除 captcha 页面引用。
