---
doc_type: goal-execution
id: E-021-r3-c4-capability-decision
parent: GOAL-004-r3-session-operator-console
date: 2026-09-05
source: self
status: done
version: 0.1.0
---

# E-021 · R3 C4 capability 路由裁决与合同门禁（2026-09-05）

## 已发生事实

- 用户通过裁决工具选择 **独立 capability 路由（推荐）**，关闭 I-033-023 的 API
  形状未知项；决策已记录为 D-011。
- D-011 固定独立 `GET /api/channel/telegram/operator/sessions/{chat_id}/capability`
  路由、`telegram.operator.read`、`{chatId,canSend}` 响应、channel.telegram
  capability service 缓存 owner、60 秒 bot/chat absolute TTL、single-flight、
  Telegram 403 精确失效和 `refresh=1` 显式重探。
- D-011 同时固定 Telegram member 状态映射、非 403 capability 错误及现有 D-008/D-009
  发送/retry 状态机和 C4 UI 接缝；未扩展 R3 首波边界。
- A-036 self contract gate 判定 `pass`、`open_required: 0`；本条不把合同门禁写成
  生产实现完成，等待本地 Grok independent 合同审计。

## 产物

- `01-decision/D-011-r3-c4-independent-capability-route.md`
- `03-audit/A-036-r3-c4-capability-contract-self.md`

## 状态边界

本条只记录用户裁决和合同门禁事实，不修改 GOAL-004/R3 status 或 progress，不关闭
C4；生产实现须在 independent 合同审计完成后进行。
