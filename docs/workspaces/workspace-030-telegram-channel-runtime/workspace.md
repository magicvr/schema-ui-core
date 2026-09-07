---
id: workspace-030-telegram-channel-runtime
title: Telegram Bot 通道运行时工作区（架构分支 · C 端通道）
status: done
root_goal: GOAL-001-telegram-channel-runtime
canonical_scope: docs/workspaces/workspace-030-telegram-channel-runtime/
shared_materials_catalog: none
vision_role: delivery
plan_refs:
  - VP-030-telegram-channel-runtime
primary_plan: VP-030-telegram-channel-runtime
created: 2026-09-03
updated: 2026-09-05
version: 1.1.0
parent: null
---

# 工作区上下文 · Telegram Bot 通道运行时

本工作区是 [VP-030-telegram-channel-runtime](../../vision/plans/VP-030-telegram-channel-runtime.md)（**`closed`** v0.3.0 · 2026-09-05 用户指令授权关门 · [VRev-076](../../vision/reviews/VRev-076-vp030-telegram-channel-runtime-close-out.md) self `pass`）的唯一 lead delivery workspace。**架构分支 · C 端通道**（对标 VP-017：内核端口 + 一方模块 + Admin 设置）：交付 Telegram Bot 通道运行时——HTTPS webhook（secret fail-closed）+ Update 分发端口（命令/callback Register）+ `SendMessage` 文本端口 + `issuer=telegram` 主体映射 + Admin bot 设置。**不是**业务域，**不是**付费命令实现。**工作区已顺利结项关门。**

- **Root** `GOAL-001-telegram-channel-runtime`：**`done`** · **4/4 + R5**（R1 合同冻结 GOAL-002 3/3 → R2 webhook+分发+身份 GOAL-003 3/3 → R3 出站+设置+限流核账 GOAL-004 3/3 → R4 证据与关门 GOAL-005 3/3 全量达成关门；A-006/A-008 审计闭环后，判据 #5 补做 Admin UI tab 由 GOAL-006 3/3 完成，Root 回归 done），纲领见 Root `00-meta.md`。
- 激活门禁已满足（2026-09-03）：[VRev-070](../../vision/reviews/VRev-070-vp030-telegram-channel-runtime-activation.md) self `pass`（0 required；V-F114/115 → 开区事务内 fixed）；**架构类轻量 freshness PASS**（`b5c39dfb` → `42036a3c`：协议 pin / 依赖锁 / Profile 默认集 / provenance 零变更；区间代码 = VP-029 已审结目）不暂挂 `go`；**限流评估落盘**：进程内够用、不需要 Redis，不消耗 RT-Q05 trigger。
- 不改变 Charter `primary_workspace`（仍为 workspace-001）。
- **消费基线**：VP-017 通道形态（端口 + 设置 + mock）· VP-027 RateLimiter 端口（已 closed）· VP-029 `GetOrCreateSubject`（已 closed；不得要求 `admin.wallet` HTTP 已启）· VP-021 停机 drain · VP-003/004 模块契约（横切 + 设置面，豁免业务导航）。
- **红线（激活即生效）**：不进 `mvp`/`admin` 默认集；不做 Mini App / Stars / 对话 FSM / 付费命令；不引入独立 Bot 进程 / 长轮询生产 / 多 bot；不把 Bot 用户写入 `admin.users`；密钥 fail-closed、不进配置包明文；不消耗 RT-Q05 Redis trigger；不重开 VP-017/026/027/028/029；内核不得 import `channel.telegram` 实现细节。

## 绑定

| 字段 | 当前值 | 说明 |
|------|--------|------|
| 工作区 ID | `workspace-030-telegram-channel-runtime` | 与本区目标及资料引用的 `workspace_id` 一致 |
| Root Goal | `GOAL-001-telegram-channel-runtime` | `parent: null`；**done** · 4/4（R1～R4 全部关门） |
| canonical 范围 | `docs/workspaces/workspace-030-telegram-channel-runtime/` | 本区唯一目标状态范围 |
| 共享资料目录 | `none` | 暂无固定共享资料 |
| 愿景角色 | `delivery` | VP-030 lead（closed v0.3.0）；不改变 Charter primary workspace |
| 规划对齐 | `primary_plan` = `VP-030-telegram-channel-runtime`（`closed` v0.3.0） | 2026-09-05 关门（VRev-076 self `pass`；八条判据 verified；Root done；R-009 按 A-009 保留 bounded accepted-residual） |

## 愿景对齐

Charter：`schema-ui-core-admin-foundation@0.4.0`（2026-08-31 strategic：同进程基座 · 成功边界 #6 · H-002）。
VP-030：Telegram Bot 通道运行时（vision_ref @0.4.0）——八条方向级退出判据 = webhook secret fail-closed / 分发 Register / SendMessage mock / 身份映射 / 设置与密钥 / 限流评估（本条已由 VRev-070 核销）/ 边界保持 / required=0 闭合；红线 = 不进默认集、不做 Mini App/Stars/FSM/付费命令、不消耗 Redis trigger。

## 纲领阶段（Root 路线图指针）

| 阶段 | 内容 | 状态 |
|------|------|------|
| R1 | **合同冻结**（判据 1/2/3/6 + I-030-001/002/003/006 裁决）：无 token 启动策略 · HTTP vs SDK · 桶分母 · 请求计数 vs 失败预算映射 · 分发 API / mock | **已关门**（GOAL-002 done 3/3；D-002 v0.1.0；kernel/telegram.go 落地并通过测试） |
| R2 | **webhook + 分发 + 身份**（判据 1/2/4 + I-030-007）：secret fail-closed 路由 · Register 分发 · `GetOrCreateSubject("telegram", id)`（不要求钱包 HTTP 已启）；入站三桶随 webhook 落地 | **已关门**（GOAL-003 done 3/3；grok 独立审 F-001 闭合；装配全绿） |
| R3 | **出站 + 设置 + 限流核账**（判据 3/5 + I-030-005）：SendMessage mock/生产供应商隔离 · Admin bot tab · 入站限流核账 | **已关门**（GOAL-004 done 3/3；HTTPSender + RuntimeManager 热切换 + 设置端点 + 限流核账） |
| R4 | **证据与关门**（判据 7/8）：证据矩阵 / 越界核账 / 审计闭合 | **已关门**（GOAL-005 done 3/3 · 证据矩阵全项 PASS + grok 独立复审 A-003 pass） |

## 固定共享资料引用

> `shared-materials/index.json` 只能提供候选路径与摘要。缺完整引用字段的行无效。

| reference_id | workspace_id | material_id | source | version | sha256 | purpose | local_record | status |
|--------------|--------------|-------------|--------|---------|--------|---------|--------------|--------|
| — | — | — | — | — | — | — | none | — |
