---
id: GOAL-003-r2-connection-settings
title: R2 · Telegram 连接与设置实现
status: active
parent: GOAL-001-telegram-operator-console
created: 2026-09-04
updated: 2026-09-04
version: 0.1.0
progress: 0/5
plan_refs:
  - VP-033-telegram-operator-console
primary_plan: VP-033-telegram-operator-console
serves_summary: 承载 Root R2：Telegram Bot API 管理调用、mode/显式 webhook URL 设置、connection manager、互斥切换、heartbeat/占用位和 Admin 设置页；实施源为 D-002 + D-003。
---

# GOAL-003 · R2 Telegram 连接与设置实现

## 概述

本目标是 `[workspace-033-telegram-operator-console]` Root 的 R2 阶段子目标，承接已完成的 R1 合同冻结。它将 D-002 + D-003 转化为 Telegram Bot API client、配置/迁移、connection manager、Fx 生命周期接缝、管理设置 API/UI 与可重复的 Fake Bot API 证据；不进入 R3 会话落盘和人工 IM。

## 已冻结的实现边界

- Telegram 入站 mode 仅为 `webhook` 或 `polling`；缺省为 `polling`，生产 webhook 必须显式配置。
- webhook 公网 base URL 使用 Telegram 专属字段 `webhook_public_base_url`；不复用全局 `runtime.mode` 或 `auth.public_base_url`。
- connection manager 是唯一 receiver owner；切换先 drain 当前 owner，再建立目标模式；`OnStop` 必须等待 manager drain。
- `webhook` 使用 `getMe → setWebhook`，且发送 `secret_token`；`polling` 使用 `getMe → deleteWebhook`，receiver 启停与模式建立分离。
- 首波保持 Telegram 模块 profile-gated，不修改默认 Profile；HTTP 路由与 Telegram 侧 webhook receiver 的语义分开。

## 成功标准

- [ ] mode、webhook URL、token/secret 的配置、持久化、迁移、重启回读与 settings API 边界可核对，且不泄漏密钥。
- [ ] Bot API client 支持 `getMe`、`setWebhook`、`deleteWebhook`、`getUpdates`，长轮询 client 与 `sendMessage` client 分离并遵守 D-003 timeout 语义。
- [ ] connection manager 实现 webhook/polling 互斥切换、idle/receiver 状态、业务占用位/heartbeat lease 与 shutdown drain，并接入 Fx composition。
- [ ] Admin Telegram settings 页面支持 mode、显式 URL 与 write-only secrets，保持现有权限、Profile 与 i18n 边界。
- [ ] Fake Bot API、API/Web/组合根和生命周期测试覆盖 R1-V-002～R1-V-009；R2 相关 required finding 归零后才进入 R3。

## 阶段路线图（P-001）

阶段串行；同一阶段内仅在写集不重叠时并行。`progress = 已完成检查点 / 5`。

| 检查点 | 内容 | 状态 |
|--------|------|------|
| C1 | R2 关键参数裁决、实施计划与 required 信息闭合 | 待开始；I-033-014～016 待用户裁决 |
| C2 | Telegram 配置 schema、迁移、runtime 回读与 settings API | 待开始；依赖 C1 |
| C3 | Bot API client、connection manager、互斥切换与 Fx 生命周期 | 待开始；依赖 C1，部分依赖 C2 |
| C4 | Admin settings UI、占用位/heartbeat 接缝与跨层集成 | 待开始；依赖 C2/C3 |
| C5 | Fake Bot API、退出/错误矩阵、self + independent 阶段审视 | 待开始；依赖 C2～C4 |

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-033-014 | required | mode 与 `webhook_public_base_url` 在 YAML/env、DB seed、DB authoritative、Admin PATCH 之间的优先级 | 方案 / C2 | C1 | 用户裁决后写 D-004；补 migration/runtime/settings 回读测试 | open | 未延期 | A-004 recommended F-001；当前代码仅 token/secret 为 DB authoritative |
| I-033-015 | required | 未绑定 polling 的 heartbeat 是引用计数还是单 lease，以及 TTL/失效 drain 语义 | 方案 / C3/C4 | C1 | 用户裁决后写 D-004；补 lease 并发/过期测试 | open | 未延期 | A-002 F-005 recommended；VP-033 使用活跃会话引用计数表述 |
| I-033-016 | required | `getUpdates` 长轮询请求 timeout 与独立 HTTP client 余量的默认值 | 方案 / C3 | C1 | 用户裁决后写 D-004；补正常等待/取消/错误/timeout 测试 | open | 未延期 | A-004 F-002 recommended；D-003 已冻结严格大于关系但未定数值 |
| I-033-017 | non-blocking | disabled profile 下 Telegram HTTP surface 是否继续按现有 module gating 处理 | 实施 / C4 | C3 | R2 计划核对 provider/composition 现状并记录 | open | 可沿用现有 profile 语义 | A-002 F-007 recommended；不重开默认 Profile 红线 |
| I-033-018 | non-blocking | `HasBusinessHandlers` 放在具体 dispatcher/adapter 还是扩展 kernel 端口 | 实施 / C3 | C3 | R2 实现决策与编译期/行为测试 | open | 可在 C3 记录 | A-002 F-006 recommended |

当前 C1 前有 3 项 required 信息未闭合；在用户裁决并写入 D-004 前，不开始依赖这些选择的生产代码实现。

## 父目标

`GOAL-001-telegram-operator-console`

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；本目标已建立索引与首条建立事实，后续按编号递增。
