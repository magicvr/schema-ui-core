---
id: E-003-w18-api-web-fixes-and-regression
doc: execution-entry
goal: GOAL-019-w18-api-web-security-hardening
date: 2026-09-22
status: recorded
parent: GOAL-001-production-hardening
version: 0.1.0
---

# E-003 · API/Web required 修复与回归验证

## 事实

- API 修复：`apps/api/server/serve.go` 的 `Run` 在 store、认证和 listener 装配前调用 `cfg.validate()`；`resolveSecret` 仅在 `development` 保留开发 fallback，非 development 缺失 secret 返回空值且不会绕过入口校验。
- API 回归：`apps/api/server/serve_test.go` 新增缺失/弱 JWT secret 的 Run 级测试，使用 store guard 与预占 listener 证明配置无效时不会启动 store 或 listener；另覆盖 development fallback 与非 development 空值行为。
- Web 修复：`apps/web/src/account/auth-client.ts` 新增统一目标 URL 解析，`Request` 读取 `.url`、`URL` 读取 `.href`、字符串按当前 origin 解析；same-origin、session-list、auth endpoint 与 password-change notice 均改用该解析。
- Web 回归：`apps/web/src/account/auth-client.test.ts` 新增同源 Request 双 token、跨源 Request 无认证头、Request 形式 auth endpoint 不刷新等测试。
- 当前工作树验证：API `go test ./...` exit 0；Web `npm run test` = 124 test files / 1510 tests passed；Web `npm run typecheck` exit 0；`git diff --check` exit 0。Web 全量测试输出含既有 React `act(...)` 环境 stderr，但无失败。
- S3、S4 已完成；I-003 验证分母已记录。S5 self 阶段审计另见 `03-audit/A-001-w18-self-stage.md`；S6 仍待独立审计响应。

## 阻塞 / 风险

两条原始 MAJOR 已有修复和回归证据，但仍需 self + clean-context independent Reviewer 审计；I-004 保持 open。

## 下一步（计划）

完成 self 阶段审计，等待并代贴 independent Reviewer 意见；若发现 required finding，进入整改环，不直接关门。
