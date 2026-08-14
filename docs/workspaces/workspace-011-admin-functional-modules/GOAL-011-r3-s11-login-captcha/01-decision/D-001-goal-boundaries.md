---
id: D-001
goal: GOAL-011-r3-s11-login-captcha
title: 立项边界：模块身份、Profile 归属与审计策略
date: 2026-08-14
status: accepted
parent: GOAL-011-r3-s11-login-captcha
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# D-001 · 立项边界（S-11 登录验证码）

## 决定

1. **模块身份**：`admin.login-captcha`（标准 Admin 功能模块）；Descriptor 依赖 core.auth-session / core.schema-render / core.navigation-capability / core.operationlog。
2. **Profile 归属（I-003 闭合）**：进入 **admin 默认集**（内容扩展，先例一致）；mvp/demo 不启用。
3. **审计策略**：登录链路 + 挑战生命周期为 **security 门禁 → 独立审计必须**（grok，用户书面偏好）。
4. **验证码模型（I-001 闭合）**：**自研算术挑战**（加减小整数，如「7 + 5 = ?」）——无第三方依赖、可访问性好、无需图片生成。答案以 sha256(challengeID || answer) 哈希存储；挑战单次有效 + 5 分钟过期；惰性清理。
5. **默认关闭（I-002 部分）**：captcha_config.enabled 默认 0——默认不影响既有登录/冒烟/e2e 链路；管理员经设置页启用。启用后登录必须携带有效挑战（400 INVALID_CAPTCHA）。
