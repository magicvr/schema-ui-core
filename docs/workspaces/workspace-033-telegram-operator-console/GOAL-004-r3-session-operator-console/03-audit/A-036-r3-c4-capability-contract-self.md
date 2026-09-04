---
doc_type: goal-audit
id: A-036-r3-c4-capability-contract-self
parent: GOAL-004-r3-session-operator-console
date: 2026-09-05
source: self
auditor: Codex
audit_type: stage-gate
scope: R3 C4 I-033-023 独立 capability 路由、getChatMember 映射、60 秒 bot/chat 缓存、single-flight、403 失效、显式重探及发送/UI 接缝合同；不审生产实现，不关闭 C4
verdict: pass
open_required: 0
version: 0.1.0
---

# A-036 · R3 C4 capability 合同自审（2026-09-05）

## 审视结论

用户已通过裁决工具选择独立 capability 路由；D-011 将该选择及其可执行的
HTTP、缓存、错误、发送和 UI 行为写成实施合同，可以放行 C4 capability 生产实现。
本条不把合同通过写成实现完成，不替代下一条本地 Grok independent 审计，也不关闭
C4。

## 对照证据

- D-002 已冻结混合策略：进入/刷新会话做有限 TTL `getChatMember` 预检，真实发送
  结果为最终权威，Telegram 403 立即否决并失效对应 chat 缓存；D-011 保留了该
  策略并明确 `refresh=1` 的显式重探边界。
- D-011 固定了独立路由
  `GET /api/channel/telegram/operator/sessions/{chat_id}/capability`、同一
  `telegram.operator.read` 权限和 `{chatId,canSend}` 成功载荷；没有把 capability
  扇出到 10 秒成绩单刷新。
- D-011 将缓存 owner 固定为 channel.telegram capability service，消费 composition
  注入的 `kernel.Cache`，按 `(bot_id,chat_id)` 使用 60 秒 absolute TTL，并要求同键
  single-flight、允许/拒绝均缓存、403 精确失效；与现有 cache port 的 namespace、
  key、absolute expiry 和 delete 接缝一致。
- D-011 明确了 member 状态的 fail-closed 映射、探测非 403 错误的 cataloged
  `TELEGRAM_CAPABILITY_UNAVAILABLE`、结构化 Telegram 403 识别，以及 D-008/D-009
  已冻结的 outbound pending/sent/failed、request id 和 retry_of 不变。
- 当前 C3 handler 已有 sessions/messages/send/retry、专用 read/write 权限、运行时
  占用门禁和现有发送状态机；当前 composition 已有 Bot API client、`kernel.Cache`
  和 `TelegramRuntime` 注入接缝，后续实施只扩展这些边界，不新增产品范围。

## 门禁判定

- **用户裁决**：I-033-023 已从 `collecting` 进入 `verified (user decision)`；未选的
  成绩单附带和会话列表附带方案不进入实现。
- **可实施性**：合同覆盖 API 形状、权限、bot/chat scope、缓存 TTL/并发/失效、403
  与非 403 错误、UI 显式重探以及发送/retry 状态机接缝；实现阶段仍须以代码和
  测试形成事实证据。
- **范围**：不扩展历史回灌、FSM、群发、频道、多 bot、多实例 polling、独立进程或
  SSE/WebSocket；不把 capability 合同误报为 C4/R3 完成。
- **required findings**：本条无 required finding，`open_required: 0`。下一门禁须
  由本地 `grok-4.6 · reasoning high` 独立审视；independent 失败不得由本条降级。

本条不修改 GOAL-004/R3 status 或 progress，不接受 residual，不 overrule；C4
capability 生产实现可在 independent 合同审计完成后开始。
