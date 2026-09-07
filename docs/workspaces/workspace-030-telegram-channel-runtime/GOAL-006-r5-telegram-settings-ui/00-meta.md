---
id: GOAL-006-r5-telegram-settings-ui
title: R5 补做 Telegram 设置 Admin UI tab（判据 #5）
status: done
parent: GOAL-001-telegram-channel-runtime
created: 2026-09-03
updated: 2026-09-03
version: 1.1.0
progress: 3/3
plan_refs:
  - VP-030-telegram-channel-runtime
primary_plan: VP-030-telegram-channel-runtime
serves_summary: 承载 VP-030 判据 #5 补做（用户 2026-09-03 裁决）：channel.telegram 设置 Schema/Nav/tab——后端 Schema/Page/Nav 贡献 + 前端 telegram-admin-tab custom component + i18n + 测试。
---

# GOAL-006 · R5 Telegram 设置 Admin UI tab

## 概述

A-006/A-008 将 VP-030 判据 #5 的 Schema/Nav/tab 降级为 recommended（API-only）。用户 2026-09-03 书面裁决（GOAL-001 A-008 处置）：**补做 Admin UI tab**，判据 #5 恢复为完整交付口径。

本目标在 `channel.telegram` 模块补齐 Admin 设置 UI：后端 Schema/Page/Navigation 贡献（telegram 设置页，embed custom component）+ 前端 `telegram-admin-tab` React 组件（token/secret 编辑，对标 mail-admin-tab）+ i18n keys + 测试。

对齐递归：GOAL-006 → Root GOAL-001（判据 #5）→ VP-030 → Charter @0.4.0。不引入业务域 / 付费命令 UI；密钥 write-only，不回显。

## 成功标准 / 纲领检查点（P-001）

| 检查点 | 内容 | 状态 |
|--------|------|------|
| C1 | **后端 Schema/Nav 贡献**：`channel.telegram` 设置页 schema（custom component 挂载）+ Navigation 菜单项 + 权限沿用 settings.read/write | **已关门**（E-001：Schema/Page/Nav/Manifest 贡献 + descriptor 对齐 + 测试绿） |
| C2 | **前端 tab 组件**：`telegram-admin-tab` custom component（GET/PATCH settings，token/secret 编辑 write-only）+ main.tsx 注册 + i18n keys | **已关门**（E-002：组件 + i18n + 测试绿） |
| C3 | **审视与关门**：self 审计 A-001；无开放 required；goal-tree / VP-030 判据 #5 回写 | **已关门**（A-001 self `pass` 0 required · goal-tree 同步 · 判据 #5 完整交付） |

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-006-001 | non-blocking | UI tab 权限门：沿用 `settings.read`/`settings.write`（mail 先例，无新权限键）vs 独立权限。 | C1/C2 | C1 | mail 先例（W26：复用 settings.read，红线无新权限键） | **verified** | — | 2026-09-03：沿用 settings.read/write（mail 同款红线） |
| I-006-002 | non-blocking | 导航挂点：独立菜单项（menu_telegram）vs settings 子项。 | C1 | C1 | mail 先例（menu_mail 独立项，settings.read 可见） | **verified** | — | 2026-09-03：独立菜单项 `menu_telegram`，权限 settings.read |

## 父目标

- `GOAL-001-telegram-channel-runtime`

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger。

## 备注

- 本目标为 VP-030 判据 #5 补做（A-006/A-008 recommended → 用户裁决补做 UI tab）。不重开 R1–R4 已关门事实。
- 后端端点 `GET/PATCH /api/channel/telegram/settings` 与 `RuntimeManager` 已存在（GOAL-004 R3）；本目标只补 Schema/Nav/UI 面。
