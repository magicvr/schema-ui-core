---
doc_type: goal-decision
id: D-001-r5-ui-shape
parent: GOAL-006-r5-telegram-settings-ui
date: 2026-09-03
status: accepted
version: 1.0.0
---

# D-001 · R5 UI tab 形态（沿用 mail 先例）

## 决定

`channel.telegram` 设置 Admin UI tab 采用 **mail-admin-tab 同款形态**：

1. **后端 Schema/Nav**：`channel.telegram` 模块新增 Schema 页面贡献（PageID `telegram-settings`，embed 一个 `telegram-settings.json`，body = `custom` 节点挂 `telegram-admin-tab` 组件）+ Navigation 菜单项 `menu_telegram`。
2. **权限**：沿用 `settings.read` / `settings.write`（mail W26 先例：红线 = 无新权限键）。
3. **前端**：`apps/web/src/components/telegram-admin-tab.tsx`（GET/PATCH `/api/channel/telegram/settings`，token/secret 编辑 write-only，status 展示）+ main.tsx 注册 + i18n keys（zh-CN / en-US）。
4. **不进默认集**：`channel.telegram` 不在 `mvp`/`admin` 默认 Profile，UI 随模块启用才出现（VP-030 红线保持）。

## 未选方案

- 独立权限键（`telegram.settings.*`）：违反 mail W26 无新权限键红线。
- Settings 内嵌子 tab（非独立菜单）：mail 已先例为独立菜单项，保持一致性。

## 依据

- I-006-001 / I-006-002 verified（mail 先例）。
- VP-030 判据 #5：Admin 可配置 token/secret；密钥 fail-closed；不进配置包明文。
