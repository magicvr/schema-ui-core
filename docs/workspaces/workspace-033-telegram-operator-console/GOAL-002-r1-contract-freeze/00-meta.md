---
id: GOAL-002-r1-contract-freeze
title: R1 · Telegram 连接与人工台合同冻结
status: done
parent: GOAL-001-telegram-operator-console
created: 2026-09-04
updated: 2026-09-04
version: 0.4.0
progress: 3/3
plan_refs:
  - VP-033-telegram-operator-console
primary_plan: VP-033-telegram-operator-console
serves_summary: 承载 Root R1 合同冻结：连接模式、轮询生命周期、业务占用位、控制台 heartbeat、发言权反馈、显式公网 URL 与 Fake Bot API 验收；不在本目标实现 R2 业务代码。
---

# GOAL-002 · R1 Telegram 连接与人工台合同冻结

## 概述

本目标是 `[workspace-033-telegram-operator-console]` Root 的 R1 阶段子目标，承接 VP-033 与 Root 已冻结的产品边界。它只负责把 R2 实施所需的行为合同、状态边界、失败语义和验证矩阵冻结为可执行记录；R2 的连接实现、设置页和轮询代码须在本目标完成后另行承载。

## 已继承的冻结边界

- 入站模式为 `webhook` 或 `polling`，二者互斥；开发默认 polling，生产推荐 webhook。
- webhook 使用 `getMe + setWebhook`；polling 使用 `getMe + deleteWebhook`，并遵守未绑定 heartbeat 懒启动、已绑定 polling 常驻。
- Dispatcher 至少存在一条业务 `RegisterCommand` / `RegisterCallback` 即视为已绑定；已绑定隐藏人工台。
- 首波只收录 Telegram 实际投递给 bot 的文本更新；不引入历史回灌、FSM、群发、频道、多 bot、多实例 polling 或 SSE/WebSocket。
- webhook 公网 URL 必须来自显式配置；本地以可注入 Fake Bot API 核对 `setWebhook` 与 fail-closed 语义。

## 成功标准

- [ ] 关键实现选择已由用户确认并记录，且 `I-033-011`～`I-033-013` 关闭或保留合规的用户裁决留痕。
- [ ] R1 行为合同、状态转换、失败语义、shutdown 接缝和 Fake Bot API 验证矩阵已形成可指向 R2 的决策/附件证据。
- [ ] R1 阶段 self 审视完成，相关 required finding 为 0，且 R2 入口条件可核对。

## 阶段检查点（P-001）

`progress = 已完成检查点 / 3`。

| 检查点 | 内容 | 状态 |
|--------|------|------|
| C1 | 用户方案裁决与 R1 合同冻结 | **完成**：D-002 accepted |
| C2 | 状态/接口/失败语义/验证矩阵落盘 | **完成**：D-002 + D-003 R1 合同与 R1-V-001～009 |
| C3 | R1 阶段审视与 R2 放行建议 | **完成**：A-004 independent pass；A-005 response；required `0` |

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚阶段 | 状态 | 证据 / 结论 |
|----|------|-----------------|----------|----------|------|-------------|
| I-033-001～008 | inherited | VP-033 入站、轮询、占用位、产品边界、Privacy Mode、公网 URL | R1 合同 | R1 | verified | Root `00-meta.md` 与 VP-033；已由用户书面冻结 |
| I-033-009 | non-blocking | Admin 短轮询间隔的默认秒数 | R1 UI 合同 | R1 | open | 可在 R1 合同中给出默认值，不解除 SSE 接缝 |
| I-033-010 | non-blocking | 发言权探测与缓存失效策略 | R3 | R3 | open | R3 冻结；不阻断本 R1 |
| I-033-011 | required | `mode` 与显式 webhook 公网 base URL 的持久化/配置边界 | R2 配置与重启语义 | R1 | **verified** | 用户 2026-09-04 书面选择 Telegram 专属 `webhook_public_base_url` 配置/持久化面；D-002 |
| I-033-012 | required | 新安装/已有配置的 mode 默认及启动行为 | R2 连接建立 | R1 | **verified** | 用户 2026-09-04 书面选择缺省 `polling`、生产显式 `webhook`；D-002 |
| I-033-013 | required | polling/连接管理器的生命周期 owner 与 shutdown drain 接缝 | R2/R4 生命周期验证 | R1 | **verified** | 用户 2026-09-04 书面选择 Telegram connection manager + composition `OnStop` drain；D-002 |

当前 R1 方案冻结所需的 required 信息 `I-033-011`～`I-033-013` 已由用户决定并以 D-002 记录为 `verified`；A-002 F-001～F-003 已由用户选择的 D-003 修正路径补足合同，并经 A-004 Grok independent `pass` 复审为 `closed/fixed`。A-002 F-004～F-009 与 A-004 recommended 转入 R2 计划，不阻断本 R1 子目标关门。其余开放项均为 non-blocking，不能被写成已验证。

## 父目标

`GOAL-001-telegram-operator-console`

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；本目标首条 D/E 已建立，后续按编号递增。
