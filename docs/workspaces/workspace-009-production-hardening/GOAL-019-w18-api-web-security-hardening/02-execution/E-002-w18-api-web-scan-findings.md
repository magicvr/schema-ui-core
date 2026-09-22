---
id: E-002-w18-api-web-scan-findings
doc: execution-entry
goal: GOAL-019-w18-api-web-security-hardening
date: 2026-09-22
status: recorded
parent: GOAL-001-production-hardening
version: 0.2.0
---

# E-002 · API/Web 扫描与 findings 分类

## 事实

- API 只读扫描与主线程定向核对确认：`apps/api/server/config.go` 的 `LoadConfig` 会调用 `cfg.validate()`，但 `apps/api/server/serve.go` 的 `Run` 对直接传入的 `Options.Config` 只检查 nil；`resolveSecret` 在空 secret 时返回公开固定字符串。现有 `serve_test.go` 只覆盖 nil Config 与弱 seed，不覆盖 Run 级非 development secret 门禁。
- Web 只读扫描与主线程定向核对确认：`apps/web/src/account/auth-client.ts` 的 `isSameOrigin` 与 `isSessionListRequest` 使用 `String(input)`；跨源 `Request` 的字符串化不保留真实 URL，现有测试只覆盖跨源字符串 URL，不覆盖跨源 Request 对象。
- 已将上述两项分类为 MAJOR、纳入 D-002 required 修复范围；API MFA 装配与 Web preview blob 观察按 D-002 保留为 NOTE。
- S1 扫描与 S2 范围冻结已完成；按 6 个检查点的当时状态应记录为 `2/6`（约 33%）。原 `1/6` 为算术勘误，未改变当时 required finding 尚未修复的事实。

## 阻塞 / 风险

I-003 的最终验证分母仍待修复后确定；required finding 尚未闭合，不能进入关门声明。

## 下一步（计划）

完成 API 与 Web 局部修复及针对性回归测试，随后补齐跨层验证矩阵并执行 self 审计。
