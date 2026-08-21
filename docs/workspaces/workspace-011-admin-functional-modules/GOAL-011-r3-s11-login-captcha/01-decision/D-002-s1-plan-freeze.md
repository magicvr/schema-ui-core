---
id: D-002
goal: GOAL-011-r3-s11-login-captcha
title: 方案冻结：登录验证码设计（S1）
date: 2026-08-14
status: accepted
parent: GOAL-011-r3-s11-login-captcha
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# D-002 · 方案冻结（S-11 登录验证码）

## 1. 挑战模型（I-001）

- 生成：随机两数加减（操作数 1-50，结果非负），挑战 id = 16 随机字节 hex；答案哈希 = sha256(id + ":" + answer) 存入挑战表。
- 生命周期：5 分钟过期；单次有效（校验成功即删除）；创建时惰性清理过期行（best-effort）。
- 响应：GET /api/auth/captcha（public）→ {enabled: bool, challenge?: {id, question, expiresInSeconds}}——关闭时 enabled:false 无 challenge（web 预检用）。

## 2. 登录集成（I-002）

- 请求体扩展：credentials + {captchaId, captchaAnswer}（optional 字段；未启用时忽略）。
- 顺序：decode/字段校验 → 限流（既有 D-001 P1）→ **captcha 校验（启用时）** → Login。captcha 校验：id 存在、未过期、哈希匹配、一次性删除；失败 → 400 INVALID_CAPTCHA。
- 集成点：handler.CaptchaVerifier 接口（Required() bool; Verify(id, answer string, now time.Time) error）；authsHandler 增可选 verifier 参数；composition 在模块启用时注入（core 登录 handler 语义不变——nil verifier 即原行为）。
- 与 423 锁定/限流叠加：captcha 失败不计入锁定（锁定针对凭据失败）；限流先于 captcha（防挑战耗尽）。

## 3. 配置与页面

- 表 captcha_config（单行 id=1：enabled）；端点 GET/PATCH /api/captcha/settings（captcha.read/write，PolicyAdmin）。
- 页面 `captcha`（设置表单：启用开关；manifest + menu_captcha）。
- 审计事件：`captcha.settings-update`（0024 CHECK 扩展）。

## 4. 迁移

- **0023**（admin.login-captcha）：captcha_challenges + captcha_config。
- **0024**（core.operationlog）：CHECK + captcha.settings-update。

## 5. Web 集成

- LoginPage：挂载时 GET /api/auth/captcha 预检；enabled 时渲染问题输入 + 提交 captchaId/captchaAnswer；预检失败 fail-open（登录将按服务端要求报错）。
- 默认关闭 → e2e/冒烟登录链路零影响（既有断言不变）。

## 6. 测试与验证

- API：挑战生命周期（生成/校验/过期/一次性/惰性清理）、登录门禁（关闭=原行为 / 启用+错误答案 400 / 正确 200）、settings 门禁（401/403）、审计事件。
- Web：LoginPage 单测（预检渲染/提交字段）；e2e 默认关闭回归不变。
- 组合根：admin 权限 20→22、导航 11→12；迁移 22→24。
- go 判定（S4）：认证语义在未启用时字节级不变；启用时新门禁为模块自身语义（留痕）。

## 7. 未选方案

- 图形/滑动验证码（第三方或图像生成）：依赖/复杂度高，v1 否决（算术挑战可访问性更好）。
- 默认启用：破坏既有登录/冒烟/e2e 链路且无迁移窗口，v1 默认关闭 + 管理员启用（D-001 §5）。
- 挑战计入失败锁定：锁定语义保持针对凭据（文档化）。
