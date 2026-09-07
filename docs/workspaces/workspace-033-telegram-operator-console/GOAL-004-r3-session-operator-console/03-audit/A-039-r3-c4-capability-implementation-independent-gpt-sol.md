---
doc_type: goal-audit
id: A-039-r3-c4-capability-implementation-independent-gpt-sol
parent: GOAL-004-r3-session-operator-console
date: 2026-09-05
source: independent
auditor: subagent (gpt-5.6-sol · reasoning medium)
audit_type: implementation-closeout
scope: R3 C4 capability implementation at HEAD 964aa174；A-037 fixed-response verification；da9d955e Web form build-error remediation
verdict: pass
open_required: 0
version: 0.1.0
---

# A-039 · R3 C4 capability 实现 independent close-out（2026-09-05）

## 独立结论

本条由一次性 `subagent (gpt-5.6-sol · reasoning medium)` 独立执行，基于当前
`HEAD 964aa174` 重新读取源码、测试和 A-037/A-038 记录；不采信 self audit 作为成功依据。
结论为 `verdict: pass`、`open_required: 0`，未新增 required 或 recommended finding。
A-037 原文及其四项 required 保留；其 F-037-1～F-037-4 在当前代码、测试和接线中具备
`fixed` 证据。未调用 Grok。

## 独立核验范围与证据

### F-037-1 · capability route/service/handler/composition 注入

- handler 已定义 capability 接口，并在 capability GET 路由执行权限检查、runtime 检查、
  `refresh` 参数校验和 `TELEGRAM_CAPABILITY_UNAVAILABLE` fail-closed。
- composition 将同一 capability service 注入 handler，并以
  `CapabilityInvalidatingSender` 包装 sender。
- 证据：`apps/api/internal/handler/telegram_operator.go:50-55,101-129,165-171,251-290`；
  `apps/api/internal/composition/composition.go:624-662`。

### F-037-2 · GetChatMember、成员状态映射和结构化错误

- `GetChatMember` 使用 `{chat_id, user_id}` payload；结构化 Telegram 错误仅保留
  method/status/error metadata，不带原始响应体，并同时识别 HTTP/API 403。
- `creator`、`member`、`administrator` 允许；`restricted` 仅显式
  `CanSendMessages == true` 允许；其他状态 fail-closed。
- 证据：`apps/api/internal/channel/telegram/bot_api.go:52-76,96-113,160-177,211-291`；
  `apps/api/internal/channel/telegram/capability.go:185-209`。

### F-037-3 · cache、single-flight、refresh 和精确失效

- capability service 使用独立 namespace `telegram-operator-capability`、
  `bot:{bot_id}/chat:{chat_id}` key、60 秒 absolute TTL；同键 flight、force refresh、
  generation guard 和精确 Delete 均已接通。
- 结构化发送 403 触发精确 capability invalidation，非 403 不触发。
- 证据：`apps/api/internal/channel/telegram/capability.go:12-24,50-183,226-263`；
  `apps/api/internal/channel/telegram/http_sender.go:140-160`。

### F-037-4 · Web capability 生命周期

- Web 按 chat 维护 capability flight，并以当前 chat 与 operator-ready identity guard
  丢弃过期响应；选中/重入/手动 refresh 使用 `refresh=1`，普通周期刷新不隐式探测。
- composer/retry 仅在 allowed 且 ready 时可用；send/retry 使用新 request ID。
- 证据：`apps/web/src/components/telegram-admin-tab.tsx:324-362,396-449,464-537,865-912`；
  `apps/web/src/components/telegram-admin-tab.test.tsx:120-340`。

### Web 构建错误修复

- 日期修复后，`datePicker` 仅向 `DateField` 传入字符串 `min/max`，解析层保留
  datePicker 的 ISO 边界；数字控件仍只接收数字边界。
- 证据：`apps/web/src/renderer/form-controls.tsx:915-950`；
  `apps/web/src/renderer/render.types.ts:702-717`；
  `apps/web/src/renderer/render.test.ts:341-378`；
  `apps/web/src/renderer/visual-fidelity.test.tsx:85-97`。

## 验证事实

- `apps/api`：`go test ./... -count=1` 通过。
- `apps/web`：`npm test -- --run` 为 92 个测试文件、1213 个测试通过。
- `apps/web`：`npm run build` 通过，仅有 chunk size warning，无 TypeScript 构建错误。
- capability 定向测试为 Telegram 管理台 15/15；日期边界同时有解析和 DOM 属性回归测试。

## 边界与未覆盖项

- 未执行真实 Telegram 外部服务联调、浏览器 E2E、生产 token/权限配置或运行时外部依赖
  验证；这些不是本次源码实现 close-out 的已执行证据。
- 未逐一重读 A-001～A-036 全部历史意见；当前实现 scope 以 A-037、A-038 和本条列出的
  C4 capability / Web build 接缝为准。
- 构建会改写 3 个 conformance 生成文件；本次已按仓库约定定向恢复，未将其作为本次
  capability checkpoint 的源码变更。

## 独立门禁结论

- 无新增 required/recommended finding，`open_required: 0`。
- A-037 F-037-1～F-037-4 可由 `/govern` 按 `fixed` 响应；不接受 residual，不作
  overrule。
- 本条只提供 independent implementation close-out 意见，不直接修改 C4 status/progress；
  由编排器记录响应后再决定 C4 是否关门。
