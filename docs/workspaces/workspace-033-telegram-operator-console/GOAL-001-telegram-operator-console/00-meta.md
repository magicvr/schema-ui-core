---
id: GOAL-001-telegram-operator-console
title: Telegram Bot 人工控制台
status: active
parent: null
created: 2026-09-04
updated: 2026-09-04
version: 0.7.2
progress: 0/4
plan_refs:
  - VP-033-telegram-operator-console
primary_plan: VP-033-telegram-operator-console
serves_summary: Admin 功能分支 · 消费 VP-030 Telegram runtime，交付连接状态、互斥 webhook/polling、业务占用位与未绑定人工文本控制台；不重开 VP-030，不进入默认 Profile。
---

# GOAL-001 · Telegram Bot 人工控制台

## 概述

承接 [VP-033-telegram-operator-console](../../../vision/plans/VP-033-telegram-operator-console.md)（`active` v0.2.0 · [VRev-075](../../../vision/reviews/VRev-075-vp033-telegram-operator-console-activation.md) self `pass`）：消费 VP-030 已交付 Telegram runtime，建立 Admin 运营台的连接状态、互斥入站模式、业务占用位和未绑定人工文本会话。

## 成功标准

- [ ] webhook 模式以 `getMe + setWebhook` 建立连接，polling 模式以 `getMe + deleteWebhook + getUpdates` 建立连接；公网 URL 显式配置且本地 Fake Bot API 可验收。
- [ ] webhook 与 polling 可互斥热切换；未绑定时只在控制台 heartbeat 存活期间轮询，已绑定且 polling 时常驻。
- [ ] Dispatcher 业务 Register 非空即视为已占用；已绑定会话隐藏人工台。
- [ ] 未绑定会话按用户与群分栏落盘，只收录 Telegram 实际投递给 bot 的文本更新；人工可代 bot 发言，无权限时灰掉。
- [ ] 保持首波边界：无历史回灌、FSM、群发、频道、多 bot、多实例 polling、SSE/WebSocket；不进入默认 Profile。
- [ ] 证据矩阵可核对，相关 required finding 归零后才允许关门。

## 纲领路线图（P-001）

阶段串行；同一阶段内可并行子目标。`progress = 已完成纲领阶段 / 4`。

| 阶段 | 内容 | 检查点/状态 |
|------|------|-------------|
| R1 | 合同冻结：模式切换、轮询生命周期、业务占用位、控制台 heartbeat、发言权探测、显式公网 base URL 与 Fake Bot API 验收 | **完成**：D-002+D-003；A-004 independent pass；GOAL-002 C3 done · 3/3 |
| R2 | 连接与设置：`getMe`、`setWebhook`/`deleteWebhook`、互斥热切换、占用位和 Admin 设置页 | **进行中**：`GOAL-003-r2-connection-settings` active · 3/5；C1/C2/C3 已完成（D-001；A-006/A-012 Grok independent pass；A-013 response），C4/C5 继续推进；实施源 D-002+D-003 |
| R3 | 会话与人工台：入站会话落盘、用户/群分栏、未绑定人工 IM、发言权反馈 | 待开始；依赖 R2；I-033-010 最晚本阶段冻结 |
| R4 | 证据与关门：退出判据矩阵、红线核账、审计 finding 闭合 | 待开始；依赖 R1～R3 |

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚阶段 | 状态 | 证据 / 结论 |
|----|------|-----------------|----------|----------|------|-------------|
| I-033-001 | required | 入站模式开关与互斥热切换 | 判据 2 | R1 | **verified** | 2026-09-03 用户书面：webhook \| polling |
| I-033-002 | required | 轮询启停生命周期 | 判据 3 | R1 | **verified** | 未绑定 + heartbeat 才跑；已绑定 + polling 常驻 |
| I-033-003 | required | 占用位与人工台可见性 | 判据 4 | R1 | **verified** | Dispatcher 业务 Register 非空即占用；已绑定隐藏人工台 |
| I-033-004 | required | 连接口径 | 判据 1 | R1 | **verified** | webhook = `getMe + setWebhook`；polling = `getMe + deleteWebhook` |
| I-033-005 | required | 首波产品边界 | 判据 5/6 | R1 | **verified** | 只文本、无历史/FSM/群发/频道/多 bot；同进程 |
| I-033-006 | non-blocking | 开发/生产默认 | 判据 2 默认值 | R1 | **verified** | 开发默认 polling；生产推荐 webhook |
| I-033-007 | required | Telegram Privacy Mode 与群消息可见性 | 判据 5 | R1 | **verified** | 2026-09-04：不要求关闭 Privacy Mode；仅收录 bot 实际可见更新 |
| I-033-008 | required | webhook 公网 URL 与本地验收 | 判据 1/2 | R1 | **verified** | 2026-09-04：显式公网 base URL；本地 Fake Bot API；不做运行时猜测 |
| I-033-009 | non-blocking | Admin 短轮询间隔 | 判据 5 | R1 | **open** | R1 给默认秒数；不解除 SSE 接缝 |
| I-033-010 | non-blocking | 发言权探测与缓存失效 | 判据 5 | R3 | **open** | `getChatMember` 预检或发送 403 后灰掉，待 R3 冻结 |
| I-033-011 | required | `mode` 与显式 webhook 公网 base URL 的持久化/配置边界 | R2 配置与重启语义 | R1 | **verified** | 用户书面裁决：Telegram 专属 `webhook_public_base_url` 配置/持久化面；D-002 |
| I-033-012 | required | 新安装/已有配置的 mode 默认及启动行为 | R2 连接建立 | R1 | **verified** | 用户书面裁决：缺省 `polling`、生产显式 `webhook`；D-002 |
| I-033-013 | required | polling/连接管理器的生命周期 owner 与 shutdown drain 接缝 | R2/R4 生命周期验证 | R1 | **verified** | 用户书面裁决：Telegram connection manager + composition `OnStop` drain；D-002 |

当前 R1 方案冻结的 required 信息 `I-033-011`～`I-033-013` 已由 D-002 记录为 `verified`；A-002 F-001～F-003 已按用户选择的 D-003 修正路径由 A-003 标记为 `fixed`，并经 A-004 Grok independent `pass` 复审；R1 C3 已由 A-005 完成。R2 子目标已完成 C1/C2/C3（D-001；I-033-014～016 verified；v67 migration、DB authoritative、settings PATCH；A-006/A-012 Grok independent pass；A-007/A-013 response；A-008 状态纠正），当前为 `active · 3/5`，C4/C5 尚未完成。`I-033-009/010` 仍为 non-blocking，必须在各自最晚阶段前处理。

## 父目标

- `null`（Root）

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；D-001/E-001 已首条落盘，后续按编号递增。
