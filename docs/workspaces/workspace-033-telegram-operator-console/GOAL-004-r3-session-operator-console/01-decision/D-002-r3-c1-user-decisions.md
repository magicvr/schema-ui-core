---
doc_type: goal-decision
id: D-002-r3-c1-user-decisions
parent: GOAL-004-r3-session-operator-console
date: 2026-09-04
source: user
status: done
version: 0.1.0
---

# D-002 · R3 C1 用户方案裁决

## 用户已裁决

| 信息项 | 裁决 | 影响 |
|--------|------|------|
| I-033-010 发言权 | **混合策略**：进入/刷新会话做有限 TTL `getChatMember` 预检；发送时以真实发送结果为最终权威；Telegram 403 立即否决并使该 chat 缓存失效 | 需要预检状态、发送 403 降级和缓存失效测试 |
| I-033-010 缓存恢复 | **60 秒显式重探**：缓存按 bot/chat 维度保存 60 秒；403 后禁用 composer，只有重新进入或手动刷新会话才重探 | 不做后台自动重试，限制 Bot API 请求 |
| I-033-019 会话主键 | **`chat_id`** | 一个 Telegram chat 对应一个私聊/群分栏，参与者作为消息元数据保存 |
| I-033-021 权限 | **专用 `telegram.operator.read` / `telegram.operator.write`** | transcript 读取与人工发送独立于 `settings.read/write`；需 Provider/RBAC/角色和服务凭据 scope 验证 |
| I-033-009 Admin 刷新 | **10 秒单飞、失焦暂停**：每 10 秒最多一个请求，页面隐藏时暂停，恢复时立即刷新 | 保持短轮询，不解除 SSE/WebSocket 接缝 |
| I-033-020 入站幂等 | **`update_id` 主键**：按 bot 维度唯一记录；重复 update 不重复落盘或分发，同时保存 `message_id` 元数据 | webhook 重试和 polling 重复投递共享同一幂等边界 |
| I-033-022 人工发送状态 | **状态机幂等**：先记录 `pending` 与客户端 `request_id`，结果转为 `sent`/`failed`；同一 request 不重复外发，失败可显式重试但不自动重试 | 需要明确发送与持久化顺序、冲突响应和状态查询 |

## 裁决边界

- 以上是用户对 R3 方案的书面选择，替代候选分析中的待定项；不得把未选择的 A2/A3、B2/B3 或 C2 权限方案混入实现。
- `getChatMember` 查询目标是 bot 在目标 chat 的发言能力；具体 Telegram member 状态到 `can_send` 的映射、API 错误编码和数据库字段在实施合同中按该裁决细化，不改变已选策略。
- `request_id` 的格式、消息正文长度上限、列表页/成绩单分页大小和数据库保留上限属于实施参数；如改变上述数据一致性、权限或产品边界，必须重新进入 `/govern` 裁决。
- 既有 VP-033 首波边界保持不变：只文本、无历史回灌/FSM/群发/频道/多 bot/多实例 polling/独立进程/SSE/WebSocket。

## 依据

候选与取舍见 `attachments/r3-c1-option-analysis.md`；当前代码接缝见 `apps/api/internal/channel/telegram/webhook.go`、`apps/api/internal/channel/telegram/types.go`、`apps/api/modules/channel/telegram/provider.go` 与 `apps/web/src/components/telegram-admin-tab.tsx`。
