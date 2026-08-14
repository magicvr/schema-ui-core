---
id: A-002
goal: GOAL-011-r3-s11-login-captcha
source: self
date: 2026-08-14
scope: S2 实现 + S3 验证 + S4 go 影响判定
verdict: pass
parent: GOAL-011-r3-s11-login-captcha
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# A-002 · self 审计（S2 实现 + S3 验证 + S4 go 判定）

## 结论

**verdict: pass**。实现与 D-002 逐条对应；验证全绿；go 判定为「未启用字节级不变，启用为模块自身语义」。

## 核对（S2 实现 vs D-002）

1. **挑战模型（§1）**：Generate 随机两数加减（操作数 1-50，结果非负），answer_hash = sha256(id + ":" + answer)（challenge.go），5 分钟 TTL、单次有效（Verify 任意尝试即删）、创建时惰性清理（repository.CreateChallenge）——与 D-002 §1 一致。
2. **登录集成（§2）**：login() 顺序 = decode → 限流 → captcha（h.captcha != nil && Required()）→ Login；INVALID_CAPTCHA 400；captcha 失败不计锁定（锁定仅在凭据失败路径 record）；credentials +captchaId/captchaAnswer omitempty——与 D-002 §2 一致。captcha == nil 时无任何行为分支（auth.go login 仅一处条件包裹）。
3. **配置与页面（§3）**：GET/PATCH /api/captcha/settings（captcha.read/write PolicyAdmin）、captcha 页面（settings 表单 + menu_captcha Order 7）、审计事件 captcha.settings-update（0024 CHECK）——一致。
4. **迁移（§4）**：0023 login_captcha（checksum 6bc0b556…）、0024 operation_log_captcha（51d6a3c1…）经 store ownership 测试冻结；compiled/persistence.go 注册。
5. **Web（§5）**：fixture/i18n/smoke/shell.spec 同步；默认关闭 → 登录链路断言零改动。

## 核对（S3 验证）

- go test ./...（apps/api）：全包通过（含 handler 111s、store 29s、composition、logincaptcha 17s）。
- vitest run（apps/web）：50 files / 900 tests 全绿（含 fixture sha 重钉、schema-keys、s5 分母、ui-bilingual）。
- 迁移计数 22→24 断言（migrate/restart/operations）全部更新并绿。
- composition admin 权限 22 / 导航 12 断言绿。

## go 影响判定（S4，核心认证路径）

- **接口**：POST /api/auth/login 请求体仅新增 omitempty 可选字段（captchaId/captchaAnswer）——字段校验不要求它们；响应不变。
- **语义**：verifier 为 nil（模块未启用，默认）或 Required()==false（启用但关闭）时，login() 与改造前逐字节相同（测试 TestCaptchaLoginGateDefaultOff 实证）。
- **启用路径**：新门禁属于 admin.login-captcha 模块自身语义（permission/审计/页面均由该模块贡献），非 core.auth-session 变更；core 登录 handler 无签名破坏（变参 verifier 向后兼容）。
- **结论**：无 breaking change；go 判定通过（D-002 §6 留痕）。

## Findings

- 无 required。self 视角无 open finding；S5 由 grok 独立审计兜底（security 门禁）。
