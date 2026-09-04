---
id: GOAL-004-r3-session-operator-console
title: R3 · 会话落盘与未绑定人工 IM
status: active
parent: GOAL-001-telegram-operator-console
created: 2026-09-04
updated: 2026-09-04
version: 0.8.0
progress: 1/4
plan_refs:
  - VP-033-telegram-operator-console
primary_plan: VP-033-telegram-operator-console
serves_summary: 承载 Root R3：Telegram 实际投递文本的会话落盘、私聊/群分栏、未绑定人工控制台、管理员代 bot 发言与发言权反馈；不进入历史回灌、FSM、群发、频道、多 bot、多实例 polling 或 SSE/WebSocket。
---

# GOAL-004 · R3 会话落盘与未绑定人工 IM

## 概述

本目标承接已关闭的 R2，交付 VP-033 的会话与人工台阶段：从 webhook/polling 的共同入站路径识别 Telegram 实际投递的文本更新，按 chat 形成私聊/群会话与文本成绩单，并为未绑定且连接成功的 Admin 提供人工代 bot 发言和发言权反馈。C1 已经 A-005 Grok independent `pass` 与 A-006 response 关闭；C2 已取得 D-004 用户裁决，D-005 经 A-008 F-001/F-002、D-006/A-009 fixed 响应和 A-010 Grok independent `pass` 修正，A-011 已响应，现进入代码实施，尚未完成 C2。

## 已冻结边界

- 只处理 Telegram 实际投递给 bot 的文本更新；不做历史回灌、媒体/文件/贴纸、FSM、群发、频道、多 bot、多实例 polling 或独立 Bot 进程。
- 复用 `Dispatcher.HasBusinessHandlers()` 作为业务占用位；已绑定时人工台入口必须隐藏，不能新增 Offer/租户绑定表。
- 人工发送必须经现有 `TelegramSender`/`HTTPSender` 以 bot 身份完成；不把人工消息写入业务领域事件总线。
- Admin 实时更新继续使用短轮询/现有页面接缝，不解除 SSE/WebSocket 门闩；具体刷新间隔仍登记为 `I-033-009`。

## 成功标准

- [ ] 入站文本在 webhook 与 polling 共同路径上按会话落盘，具备可核对的 chat/user 元数据、方向、Telegram message/update 标识和幂等边界；不回灌历史。
- [ ] 未绑定且连接成功时可列出有入站活动的私聊/群并展示该会话文本成绩单；已绑定时人工台入口隐藏，仍遵守 R2 占用位与 Profile 边界。
- [ ] 管理员可经鉴权 API 使用现有 Telegram sender 发送文本；发送成功与失败状态、重复提交和持久化顺序可核对，不泄漏 token/secret。
- [ ] 发言权策略已由用户裁决并落盘；composer 在无权时禁用，缓存/403 失效路径和非 403 错误均有测试事实。
- [ ] R3 API/Web/迁移/运行时/并发边界有 self + independent 审视，required finding 归零后才进入 R4。

## 阶段路线图（P-001）

阶段串行；同一阶段内仅在写集不重叠时并行。`progress = 已完成检查点 / 4`。在 C1 的用户裁决和 required 信息闭合前，不进入 C2～C4 实施。

| 检查点 | 内容 | 状态 |
|--------|------|------|
| C1 | R3 数据/权限/发言权合同、信息需求与用户裁决冻结 | **完成**：D-002+D-003/E-002～E-004；A-002 self `pass`；A-003 F-001 → A-004 `fixed`；A-005 Grok independent `pass`；A-006 response；开放 required = 0 |
| C2 | Telegram 文本入站、会话/消息持久化、迁移与幂等边界 | **实施中**：A-010 Grok independent `pass` 确认 A-008 F-001/F-002 合同 fixed，A-011 已响应；按 D-005 进入 v68、共同入站接线与测试，C2 尚未完成 |
| C3 | 会话列表/成绩单/人工发送 API、权限与运行时接线 | 待开始；依赖 C2 |
| C4 | Admin 人工台 UI、发言权反馈、端到端验证与 independent 审计 | 待开始；依赖 C3 |

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|----------|-----------------|------|-------------|-------------|
| I-033-009 | non-blocking | Admin 短轮询刷新间隔与请求并发/失焦行为 | C3/C4 | R3 | 用户裁决；补 UI timer/失焦测试 | **verified (decision)** | 未延期 | D-002：10 秒单飞、失焦暂停、恢复立即刷新 |
| I-033-010 | required（本 R3 gate） | `getChatMember` 预检 vs 发送 403 后灰掉，以及缓存 TTL/失效/重新探测策略 | C1/C4 | R3 | 用户裁决；补 Bot API、API、UI 状态机测试 | **verified (decision)** | 未延期 | D-002：混合策略；60 秒 bot/chat 缓存；403 失效；显式重探 |
| I-033-019 | required | 会话主键与分栏：`chat_id`、`subject_id` 或组合；私聊/群标题、排序、分页与未读语义 | C1/C2 | C1 | 用户裁决；写 schema/迁移与列表合同测试 | **verified (decision)** | 未延期 | D-002：`chat_id` 为会话边界 |
| I-033-020 | required | 入站文本准入、UpdateID/messageID 幂等、重复投递与 webhook 重试的落盘/分发顺序 | C1/C2 | C1 | 用户裁决；补 webhook/polling/retry/concurrency 测试 | **verified (decision + contract; independent pass)** | 未延期 | D-002 + D-003；A-005：持久化成功后才 webhook 2xx / polling offset；失败可重试；重复 bot-scoped `update_id` 不重复落盘/分发；C2 实现待验证 |
| I-033-021 | required | 人工台读取/发送 API 的权限边界：复用 `settings.read/write` 或专用 operator 权限 | C1/C3 | C1 | 用户裁决；更新 Provider/RBAC/auth 测试 | **verified (decision)** | 未延期 | D-002：`telegram.operator.read/write` |
| I-033-022 | required | 发送成功/失败的消息状态、落盘先后、失败重试与重复提交语义；是否记录失败消息 | C1/C3 | C1 | 用户裁决；补 sender/store/API 并发与失败测试 | **verified (decision)** | 未延期 | D-002：`pending`→`sent/failed`、`request_id` 幂等、失败显式重试 |

## 父目标

`GOAL-001-telegram-operator-console`

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；D-001～D-006、E-001～E-008、A-001～A-011 已记录 R3 建立、C1 审计闭合、C2 用户裁决、independent finding 响应与 C2 实施放行，后续按编号递增。
