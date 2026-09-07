---
doc_type: goal-decision
id: D-002-r1-info-adjudication
parent: GOAL-001-telegram-channel-runtime
date: 2026-09-03
status: accepted
version: 0.1.0
---

# D-002 · R1 信息裁决（投影）

## 上下文

用户 2026-09-03 指令推进 R1 合同冻结，并对编排器提出的六项书面选项全部采纳建议项。权威正文在子目标：

- 裁决记录：[GOAL-002 D-001](../../GOAL-002-r1-contract-freeze/01-decision/D-001-info-adjudication.md)
- 合同正文：[GOAL-002 D-002](../../GOAL-002-r1-contract-freeze/01-decision/D-002-telegram-channel-contract.md) v0.1.0

## 决定（摘要）

| ID | 决定 |
|----|------|
| I-030-001 | 进程可启动；无 token 时 webhook **503**；出站 **mock** |
| I-030-002 | stdlib `net/http`；**不**引入 Telegram SDK |
| I-030-003 | 三桶全做：`tg:webhook:{ip}` / `tg:chat:{chat_id}` / `tg:user:{telegram_user_id}` |
| I-030-004 | 模块 id = `channel.telegram`（不进默认集） |
| I-030-006 | 请求计数：每次入站 `Record`，**永不** `Clear` |
| 分发/mock | Register/Unregister 命令+callback；未知命令确定回落；`POST /api/channel/telegram/webhook` + `X-Telegram-Bot-Api-Secret-Token`；mock 可检视出站记录 |

## 纲领微调

入站三桶使用点随 **R2 webhook** 落地（公开 ingress 不得无限流）。R3「限流接入」收窄为出站侧（若需）+ 设置面 + 核账。

## 未选方案

见 GOAL-002 D-001。本条不重复展开。
