---
id: GOAL-002-r1-contract-freeze
title: R1 合同冻结（Telegram 通道端口 / webhook / 分发 / 限流映射）
status: active
parent: GOAL-001-telegram-channel-runtime
created: 2026-09-03
updated: 2026-09-03
version: 0.1.0
progress: 1/3
plan_refs:
  - VP-030-telegram-channel-runtime
primary_plan: VP-030-telegram-channel-runtime
serves_summary: 承载 VP-030 R1：冻结无 token 启动策略、stdlib HTTP、三桶请求计数映射、模块 id、分发 API / mock / webhook 路径；合同正文 = D-002；端口本体落地归 C2。
---

# GOAL-002 · R1 合同冻结

## 概述

执行 Root 纲领 **R1**：在 VP-017 `kernel.MailSender` 与 VP-027 `kernel.RateLimiter` 先例之上，冻结 Telegram Bot 通道运行时的**合同分母**——入站 webhook（secret fail-closed）· 内核级命令/callback 分发 · `SendMessage` 文本端口 · 无 token mock · 限流桶与 Record/Clear 映射。**合同正文 = 本目标 D-002**。端口本体（`apps/api/kernel/telegram.go` + 合同级快测）归 C2；webhook 路由 / 身份映射归 R2；Admin 设置 / 出站生产适配器归 R3。

对齐递归：GOAL-002 → Root GOAL-001（R1）→ VP-030（判据 1/2/3/6）→ Charter @0.4.0。不进入 Mini App / Stars / 对话 FSM / 付费命令 / 独立 Bot 进程。

## 纲领检查点（P-001）

| 检查点 | 内容 | 状态 |
|--------|------|------|
| C1 | **信息裁决**：I-030-001/002/003/006 required + I-030-004 与分发/mock 建议包（P-004） | **已关门**（2026-09-03 用户书面全部采纳建议项——D-001） |
| C2 | **合同正文 + 端口落地**：D-002 冻结；`kernel/telegram.go`（Sender / Dispatcher / Update 薄类型 + Validate）+ 合同级快测绿 | **进行中**（D-002 正文已冻；端口代码未落地） |
| C3 | **审视与关门**：合同自审（self A-001）；无开放 required；Root 信息台账回写 | 待 C2 |

`progress` = 已关门检查点数 / 3。当前 **1/3**。

## 成功标准（方向级）

1. 入站合同冻结：无 token 进程可启动、webhook 503；secret 校验 fail-closed；路径与头字段钉死（判据 #1）。
2. 分发端口冻结：命令与 callback 的 Register/Unregister；未知命令确定回落、不 panic（判据 #2）。
3. 出站端口冻结：`SendMessage` 文本（可选 callback 按钮）；stdlib HTTP；无 token 走 mock 记录；公共面无 SDK 类型（判据 #3）。
4. 限流映射冻结：三桶（IP / chat_id / telegram_user_id）独立 limiter；每次入站 Record、永不 Clear（判据 #6 实施）。
5. 模块 id = `channel.telegram`；不进 `mvp`/`admin` 默认集。
6. 未越界：不改 Charter；不做 Mini App / Stars / FSM / 付费命令；不引入 Bot SDK / Redis / 独立 Bot 进程；不重开 VP-017/026/027/028/029。

## 信息就绪与未知项

与 Root / VP-030 同号镜像。I-030-005（token 热切换）最晚 R3、I-030-007（主体 Store 路径）最晚 R2，不在本目标关闭。

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-030-001 | required | 无 token 时：拒绝进程启动 vs 启动但 webhook 503 | 方案冻结 + 判据 1/3 | C1 | 用户裁决 | **verified** | — | 2026-09-03 用户裁决：进程可启动；webhook 拒绝入站（503）；出站走 mock（D-001） |
| I-030-002 | required | Bot API：标准库 HTTP vs 第三方 SDK | 判据 3 | C1 | 用户裁决 | **verified** | — | 2026-09-03 用户裁决：stdlib `net/http`，不引入 Telegram SDK（D-001） |
| I-030-003 | required | 限流桶分母：IP / chat_id / telegram_user_id 哪些本波必做 | 判据 6 实施 | C1 | 用户裁决 | **verified** | — | 2026-09-03 用户裁决：三桶全做（D-001） |
| I-030-004 | non-blocking | 模块 id 最终字符串 | 装配 | C1 | lead 建议 + 用户确认 | **verified** | — | 2026-09-03 用户确认：`channel.telegram`（D-001） |
| I-030-005 | non-blocking | 设置是否热切换 token | 判据 5 | R3 | 用户裁决或沿用 mail 先例 | open | — | 本目标不关闭 |
| I-030-006 | required | 入站 Update 映射 VP-027：请求计数 vs 失败预算 | 判据 6 实施 | C1 | 用户裁决 | **verified** | — | 2026-09-03 用户裁决：每次入站 Record，永不 Clear（D-001） |
| I-030-007 | non-blocking | 主体 Store 消费路径；不得要求 `admin.wallet` HTTP 已启 | 判据 4 | R2 | R2 方案记录 | open | — | 本目标不关闭 |

## 父目标

- `GOAL-001-telegram-channel-runtime`（Root · 纲领 R1）

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺记账；索引文件在本目标 `01-decision.md` / `02-execution.md` / `03-audit.md`。

## 备注

- 审计模式（Root D-001）：阶段关门 default **self**。R1 为合同 + 新内核端口面，C3 走 self；R2 webhook secret 与 R4 证据/关门按需 independent（grok build 先例）。
- D-002 冻结后，R2/R3/R4 实施与验收以本合同为分母。
