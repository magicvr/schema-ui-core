---
id: D-002-w18-findings-and-fix-scope
doc: decision-entry
goal: GOAL-019-w18-api-web-security-hardening
date: 2026-09-22
status: accepted
parent: GOAL-001-production-hardening
version: 0.1.0
---

# D-002 · 扫描 findings 分类与修复范围冻结

## 决定

将以下两条扫描结果纳入 W18 required 修复范围：

1. **API MAJOR**：`apps/api/server/serve.go` 的 `Run` 未对直接传入的 `Options.Config` 执行 `Config.validate()`；非 development 且 `AuthJWTSecret` 为空时，`resolveSecret` 返回源码公开固定密钥，形成可伪造 JWT 的启动安全缺陷。修复要求：Run 入口 fail-closed、移除非 development 的静默 fallback 语义、增加 Run 级回归测试。
2. **Web MAJOR**：`apps/web/src/account/auth-client.ts` 的 `isSameOrigin` / `isSessionListRequest` 使用 `String(input)` 解析 `RequestInfo | URL`；跨源 `Request` 可能被误认为同源并注入当前 Bearer / refresh token。修复要求：按输入类型读取真实 URL，并增加跨源/同源 Request 回归测试。

## 不纳入当前 required 的观察

- API 下游 `serve` 的 MFA 未装配目前只有在该入口承载启用 MFA 的生产数据时才构成问题；先保留为部署边界 NOTE，若后续证据证明该组合可连接此类数据，再单独立项或追加 finding。
- Web 预览 blob iframe 已有空 sandbox 与 `opener = null`，当前扫描未证明可控内容进入拼接 HTML；只作为测试增强建议，不升级为本波 required。

## 理由与约束

两条 MAJOR 均有明确代码路径和测试缺口，属于 VP-009 共享基架安全边界，且可通过局部修复与回归测试闭环。保持修复集最小，避免将未被证据触发的架构重构或产品能力带入本波。

## 关联信息项

- I-001、I-002：扫描结论已 verified。
- I-003：修复后由执行记录冻结最终 API/Web/跨层验证分母。
- I-004：保留至独立 Reviewer 复核与关门判断。
