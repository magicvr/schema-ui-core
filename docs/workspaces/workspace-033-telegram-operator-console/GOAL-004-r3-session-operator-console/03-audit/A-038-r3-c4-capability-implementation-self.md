---
doc_type: goal-audit
id: A-038-r3-c4-capability-implementation-self
parent: GOAL-004-r3-session-operator-console
date: 2026-09-05
source: self
auditor: Codex
audit_type: implementation-response
scope: R3 C4 capability implementation；响应 A-037 F-037-1～F-037-4；不关闭 C4，等待 independent implementation audit
verdict: pass
open_required: 0
version: 0.1.0
---

# A-038 · R3 C4 capability 实现 self response（2026-09-05）

## 响应结论

A-037 是 `subagent (gpt-5.6-sol · reasoning medium)` 的 independent 合同审计，原文
判定 `fail` 并提出 F-037-1～F-037-4 四项 required。当前 checkpoint `cae40b3a` 已按
`fixed` 路径形成实现、测试与接线证据；本条不改写 A-037，不接受 residual、不作
overrule，也不把 self response 当作 independent 成功依据。

## Finding 响应

### F-037-1 · capability route/service/handler/composition 注入

- **状态：fixed**。
- `TelegramOperatorHandler` 增加 capability 接口、路由分类、`GET .../{chat_id}/capability`
  分支、read permission、runtime gate 和 fail-closed 502；Provider、BuiltinModules 与
  composition 均注册同一 canonical route，并注入 channel.telegram capability service。
- 证据：`apps/api/internal/handler/telegram_operator.go:50-175,251-290`；
  `apps/api/internal/composition/composition.go:624-662`；
  `apps/api/modules/channel/telegram/provider.go:60-151`；
  `apps/api/kernel/profile.go:217`；
  `apps/api/internal/handler/telegram_operator_test.go:273-321`；
  `apps/api/internal/composition/composition_telegram_test.go`；
  `apps/api/modules/channel/telegram/provider_test.go`。

### F-037-2 · GetChatMember、成员状态映射和结构化错误

- **状态：fixed**。
- `BotAPIClient.GetChatMember` 固定以当前 bot id 作为 `user_id`；`TelegramAPIError` 保留
  安全的 HTTP/API error metadata 并可识别两类 403，不将原始响应体放入结构化错误。
  capability service 对 creator/member/administrator/restricted/left/kicked/未知状态
  进行既定 fail-closed 映射；非 403 探测错误统一到稳定 service sentinel 和 cataloged
  `TELEGRAM_CAPABILITY_UNAVAILABLE`。
- 证据：`apps/api/internal/channel/telegram/bot_api.go:52-76,160-177,211-291`；
  `apps/api/internal/channel/telegram/capability.go:90-145`；
  `apps/api/internal/channel/telegram/capability_test.go:31-73,196-230`；
  `apps/api/internal/channel/telegram/bot_api_test.go:84-110`；
  `apps/api/internal/channel/telegram/http_sender_test.go:155-174`；
  `apps/api/internal/errorcatalog/errorcatalog.go:219`。

### F-037-3 · cache namespace/key/TTL/single-flight/refresh/403 精确失效

- **状态：fixed**。
- capability service 消费 composition 注入的 `kernel.Cache`，固定
  `telegram-operator-capability` namespace、`bot:{bot_id}/chat:{chat_id}` key、60 秒
  absolute expiry；同键在途请求合并，强制刷新绕过缓存但加入同键 flight；generation
  guard 防止 403 失效与在途探测竞态重新写回旧 allow。真实发送的结构化 403 经
  `CapabilityInvalidatingSender` 删除精确 bot/chat key，非 403 保留缓存。
- 证据：`apps/api/internal/channel/telegram/capability.go:12-24,50-183,226-263`；
  `apps/api/internal/channel/telegram/capability_test.go:72-195`；
  `apps/api/internal/channel/telegram/http_sender.go:140-160`；
  `apps/api/internal/composition/composition.go:629-643`。

### F-037-4 · Web capability 请求和生命周期接缝

- **状态：fixed**。
- Web 按 chat 维护 capability flight 与 identity guard；选中/重新进入/手动刷新使用
  `refresh=1`，普通 sessions/timeline refresh 不触发探测；unknown/denied/error/in-flight
  禁用 composer/retry；发送和 retry 使用新 request id，失败立即 fail-closed 并刷新成绩单，
  不自动重探。
- 证据：`apps/web/src/components/telegram-admin-tab.tsx:324-449,464-546,790-918`；
  `apps/web/src/components/telegram-admin-tab.test.tsx:157-340`。

## 验证与非阻断项

- API 全量 `go test ./... -count=1`、Telegram/handler race test 和 Web 全量 92/1213
  已通过；Telegram 管理台新增 capability/发送/重试/竞态覆盖为 15/15。
- 补齐了跨 bot/chat 缓存隔离、真实 HTTP 403 结构化识别、旧 chat capability 响应丢弃、
  Provider/composition route 数量维护等非阻断覆盖项。
- 用户随后明确要求修复该错误；`da9d955e` 已将 `form-controls.tsx:946-947` 的数字/日期
  边界按控件类型收窄，并让 `render.types.ts` 保留 datePicker 的 ISO 边界，新增解析与 DOM
  回归测试。修复后 `npm run build` 通过；仅有既有 chunk size warning，不再存在该构建错误。

## 门禁结论

- 本条 self response：`verdict: pass`，`open_required: 0`；A-037 F-037-1～F-037-4
  响应侧均为 `fixed`。
- C4 仍保持 active，不能仅凭 A-038 关闭；下一步是对当前 checkpoints `cae40b3a` +
  `da9d955e` 进行
  GPT-5.6-sol（reasoning medium）independent implementation audit。独立审计若产生
  required finding，须再次响应/修复/复审。
