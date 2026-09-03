---
doc_type: goal-execution
id: E-002-r5-frontend-tab
parent: GOAL-006-r5-telegram-settings-ui
date: 2026-09-03
status: recorded
version: 1.0.0
---

# E-002 · C2 前端 telegram-admin-tab 组件 + i18n

## 事实

1. **组件**：新增 `apps/web/src/components/telegram-admin-tab.tsx`——`registerCustomComponent("telegram-admin-tab", ...)`；GET `/api/channel/telegram/settings` 载入状态（configured/token_set/secret_set/captured），PATCH 提交 token/secret（write-only，空值保持当前，对标 mail）；`main.tsx` 注册 import。
2. **i18n**：`en-US.json` + `zh-CN.json` 新增 `schema.telegram.*` 与 `manifest.title.telegramSettings`/`manifest.nav.telegram` keys。
3. **测试**：新增 `apps/web/src/components/telegram-admin-tab.test.tsx`（jsdom）：载入状态 + write-only 输入 + PATCH 只提交非空值 + 保存后清空输入。
4. **结构测试**：`schema-keys.structural.test.ts` manifestFiles 增加 telegram fragment。

## 验证

- `npx vitest run src/components/telegram-admin-tab.test.tsx`：2/2 ok。
- `npx vitest run src/i18n/schema-keys.structural.test.ts`：4/4 ok。
- `npx tsc -b`：我的改动文件无 TS 错误（`form-controls.tsx` 的 `min/max` 类型错误为既有问题，非本次改动引入——该文件未被修改）。

## 评估

C2 完成：前端 tab 组件 + i18n + 测试齐备。判据 #5 的 Admin UI 补做交付面完成。
