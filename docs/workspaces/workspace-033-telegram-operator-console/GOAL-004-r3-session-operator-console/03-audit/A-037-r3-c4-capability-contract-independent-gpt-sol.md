---
doc_type: goal-audit
id: A-037-r3-c4-capability-contract-independent-gpt-sol
parent: GOAL-004-r3-session-operator-console
date: 2026-09-05
source: independent
auditor: subagent (gpt-5.6-sol · reasoning medium)
audit_type: design-plan
scope: workspace-033 R3 C4 I-033-023 独立 capability 路由合同、当前 Go/Web 接缝与 C4 capability 实施前置条件；不修改代码，不修改目标状态，不关闭 C4
verdict: fail
open_required: 4
version: 0.1.0
---

# A-037 · R3 C4 capability 合同独立审计（2026-09-05）

- **source**：independent
- **auditor**：subagent (gpt-5.6-sol · reasoning medium)
- **类型** / **scope**：design-plan · `[workspace-033-telegram-operator-console]` `GOAL-004-r3-session-operator-console` 的 R3 C4 I-033-023 独立 capability 路由合同、当前 Go/Web 实现接缝、缓存/403 失效与 C4 实施前置条件。未调用 Grok；不修改生产代码、测试代码、方案正文、`status` / `progress` / `goal-tree`，不关闭 C4。
- **verdict**：fail
- **open_required**：4（F-037-1～F-037-4）
- **完整意见**：本文件（无附件）

本意见由主代理根据一次性子代理的独立审计结果代贴，保留实际 provider 与思考强度，不将其表述为 Grok 审计。A-001～A-036 原文保留、未改写；不接受 residual，不作 overrule。

## 覆盖与证据

- 工作区绑定为 `workspace-033-telegram-operator-console`，目标为 `GOAL-004-r3-session-operator-console`；`I-033-023` 已记录为 `verified (user decision)`，D-011 选择独立 capability 路由、60 秒缓存、single-flight、403 失效和 `refresh=1`。[`workspace.md:2-11`；`00-meta.md:48-58`]
- D-011 固定 `GET /api/channel/telegram/operator/sessions/{chat_id}/capability`、`telegram.operator.read`、`{chatId,canSend}` 载荷及 `(bot_id,chat_id)` 的 `telegram-operator-capability` namespace / 60 秒 absolute TTL。[`01-decision/D-011-r3-c4-independent-capability-route.md:19-39`]
- 当前 handler route kind 只有 `sessions`、`messages`、`send`、`retry`，classifier 没有 `/capability`；构造器也没有 capability service 接缝。[`apps/api/internal/handler/telegram_operator.go:50-67`；`70-85`；`119-159`]
- composition 当前只把 `tr.Sender` 接入 operator handler，尚无 capability service/cache 注入。[`apps/api/internal/composition/composition.go:636-650`]
- `BotAPIClient` 尚无 `GetChatMember`，其通用错误为字符串错误，未暴露可供精确识别的 Telegram 403 结构化类型。[`apps/api/internal/channel/telegram/bot_api.go:96-211`]
- `HTTPSender.Send` 的 HTTP/API 错误也是字符串路径，尚无 capability cache invalidation 接缝。[`apps/api/internal/channel/telegram/http_sender.go:76-154`]
- Web 组件已有 `operatorCapability` 状态，但当前生命周期只会重置为 `unknown`，未发起 capability 请求；composer/retry 依赖 `allowed`，而普通刷新没有 capability 或 `refresh=1` 请求。[`apps/web/src/components/telegram-admin-tab.tsx:266-384`；`376-384`；`717-755`]
- 目标审计序列此前止于 A-036；本条为下一条独立意见。审计未执行测试命令，未将截断工具输出作为证据。

## Findings

### F-037-1 · capability route/service/handler/composition 注入缺失

- 严重度：high
- 建议：required
- 状态：open
- 证据：`apps/api/internal/handler/telegram_operator.go:50-67`、`70-85`、`119-159`；`apps/api/internal/composition/composition.go:636-650`
- 风险：D-011 的独立 capability HTTP 合同无法执行，且权限、运行时、未绑定和当前 bot 门禁没有可验证的 handler 接缝。
- 修复/验证方向：增加 channel.telegram capability service、handler route kind/classifier/分支和 composition 注入；补权限、运行时、未绑定、当前 bot scope、成功载荷和非法 `refresh` 的测试。

### F-037-2 · GetChatMember、成员状态映射和结构化错误缺失

- 严重度：high
- 建议：required
- 状态：open
- 证据：`apps/api/internal/channel/telegram/bot_api.go:96-211`
- 风险：无法验证 creator/member/administrator/restricted/left/kicked/unknown 的 fail-closed 映射，也无法把 Telegram 403 作为可缓存拒绝、把其他探测错误稳定映射为 `502 TELEGRAM_CAPABILITY_UNAVAILABLE`。
- 修复/验证方向：增加 `GetChatMember` 请求与响应模型；结构化暴露 HTTP/API 403 和非 403；覆盖状态矩阵、畸形响应、403、非 403 和错误目录合同测试。

### F-037-3 · cache namespace/key/TTL/single-flight/refresh/403 精确失效未接通

- 严重度：high
- 建议：required
- 状态：open
- 证据：`apps/api/kernel/cache.go:138-164`；`apps/api/internal/composition/composition.go:636-650`；`apps/api/internal/channel/telegram/http_sender.go:76-154`
- 风险：可能发生跨 bot/chat 复用、并发重复探测、强制刷新覆盖错误、过期语义漂移，或真实发送 403 后继续使用过期 allow；这些都会破坏 D-011 的 fail-closed 与最终权威合同。
- 修复/验证方向：由 channel.telegram capability service 消费注入的 `kernel.Cache`，固定 namespace、bot/chat key、60 秒 absolute TTL；实现同键 single-flight、`refresh=1` 绕过并替换缓存，接入 sender 的结构化 403 精确删除；补同键并发、跨键隔离、TTL、强制重探和 403 失效测试。

### F-037-4 · Web capability 请求和生命周期接缝缺失

- 严重度：high
- 建议：required
- 状态：open
- 证据：`apps/web/src/components/telegram-admin-tab.tsx:266-384`、`376-384`、`717-755`
- 风险：composer/retry 不能取得 `allowed`，也无法验证选中 chat 切换、重新进入、手动显式重探、过期响应、403/非 403 错误和发送失败后的 fail-closed 行为。
- 修复/验证方向：实现独立 capability GET；选中/重新进入/手动刷新使用 `refresh=1`，普通 10 秒 sessions/timeline 刷新不触发探测；以 chat identity 丢弃过期响应，unknown/denied/error/in-flight 均禁用发送；补允许、拒绝、错误、切换竞态、显式重探、发送与 retry 测试。

## 门禁结论

当前 D-011 方案方向可识别，但在现有 Go/Web 接缝上尚不能形成可验证的实施闭环。四项 required finding 均为开放状态，故 **verdict = fail**，不允许以 A-036 self `pass` 直接放行 C4 capability 门禁或关闭 C4。上述 findings 应按 `fixed` 路径落实到代码、测试和错误目录，再由 self 响应并接受后续 independent 实现审计；未选择 residual 或 overrule。

本次审计未发现另列的 non-blocking finding；用户要求的非阻断项可在上述实现切片中一并补齐，但不能替代四项 required finding 的可核验证据。

## 覆盖缺口

未执行测试命令；未完整逐行读取全部历史 A-001～A-036；未完整读取 `docs/architecture/principles.md`、`workspace-protocol.md`、`independent-audit-execution.md`。这些限制不改变本条基于当前代码与 D-011 的 fail 判定；后续实现审计必须补充运行时和测试证据。

### 声明

本意见只写入 independent audit ledger，不修改目标状态；响应与 finding 闭合由 `/govern` 处理。
