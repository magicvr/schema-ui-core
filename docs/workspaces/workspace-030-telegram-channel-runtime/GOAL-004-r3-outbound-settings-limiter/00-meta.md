---
id: GOAL-004-r3-outbound-settings-limiter
title: R3 出站生产适配器、Admin 设置与限流核账
status: done
parent: GOAL-001-telegram-channel-runtime
created: 2026-09-03
updated: 2026-09-03
version: 1.0.0
progress: 3/3
plan_refs:
  - VP-030-telegram-channel-runtime
primary_plan: VP-030-telegram-channel-runtime
serves_summary: 承载 VP-030 纲领 R3：落地 Telegram 出站 HTTP 适配器（stdlib net/http、10s 超时、SendMessage 文本与 InlineKeyboard 按钮、无 token 走 mock 记录）、Admin 设置（I-030-005 热切换、密钥 fail-closed、明文防泄漏）及限流核账。自审 A-001 pass 关门。
---

# GOAL-004 · R3 出站生产适配器、Admin 设置与限流核账

## 概述

承接 Root 纲领 **R3**（对应 VP-030 退出判据 #3/#5 及 R1 合同 [GOAL-002 D-002](../GOAL-002-r1-contract-freeze/01-decision/D-002-telegram-channel-contract.md)）：
1. **出站生产适配器（判据 #3）**：
   - 实现 `kernel.TelegramSender` 接口。
   - 生产环境基于标准库 `net/http` POST `https://api.telegram.org/bot<token>/sendMessage`。
   - 严格 10s 超时控制，无第三方 Telegram SDK 依赖。
   - 纯文本消息与可选 `inline_keyboard` 按钮（限制 `callback_data` ≤ 64 字节）。
   - 适配器在交付前先调用 `msg.Validate()`，违规立即 fail-closed。
   - 无 Token 或处于 Mock 模式时，降级记录到 `CaptureSender`，不发真实网络请求。
2. **Admin 设置与密钥管理（判据 #5）**：
   - 提供 Telegram 通道设置面（Token、Webhook Secret、当前模式、状态），密钥在只读查看与配置导出时不回显明文（脱敏/掩码展示），配置保存时校验强度与格式。
   - I-030-005：裁决设置生效方式（热切换 vs 重启生效）。
3. **入站限流核账**：
   - 核对入站三桶限流（IP 60/m, Chat 30/m, User 20/m）已随 R2 落地，确认出站端口无额外外部限流负担。

## 纲领检查点（P-001）

| 检查点 | 内容 | 状态 |
|--------|------|------|
| C1 | **方案与信息裁决**：I-030-005 设置生效方式与出站/设置技术方案裁决（D-001） | **已关门**（2026-09-03 用户书面裁决：热切换 + 基于 Token 自动降级 Mock，D-001） |
| C2 | **实现与回归**：出站 HTTP 适配器、Admin 设置路由与服务、限流核账与全量测试 | **已关门**（E-002 落地并通过全量测试） |
| C3 | **审计与关门**：自审与交叉审计（A-001 self），无开放必改项关门 | **已关门**（A-001 pass，0 required） |

`progress` = 已关门检查点数 / 3。当前 **3/3**。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-030-005 | non-blocking | 设置是否热切换 token（mail 有热切换先例）还是重启生效。 | 方案冻结 | C1 | 用户裁决 | **verified** | — | 2026-09-03 用户裁决：热切换（沿用 Mail 运行时通道先例），D-001 |

## 父目标

- `GOAL-001-telegram-channel-runtime`（Root · 纲领 R3）
