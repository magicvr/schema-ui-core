---
doc_type: goal-execution
id: E-028-r4-polling-warning-remediation
parent: GOAL-001-telegram-operator-console
date: 2026-09-05
source: self
status: done
version: 0.1.0
---

# E-028 · Root R4 F-001 修复（2026-09-05）

## 已发生事实

- 针对 Root self audit A-001 的 required F-001，在 Telegram Admin 连接状态区域增加
  polling 单实例警示：明确多副本运行可能丢失 Update，且 polling 不是高可用生产路径。
- 警示使用英文/中文 i18n key `schema.telegram.status.pollingSingleInstanceWarning`，
  只在 `status.mode === polling` 时显示，并带 `data-telegram-polling-warning` 选择器。
- 增加 Web 回归断言：polling 显示警示，webhook 不显示；未改变 polling/webhook 路由、
  生命周期、缓存或 capability 方案。
- 同步 VP-033 计划中的 workspace 当前状态投影；VP-033 仍保持 `active`，Root 仍保持
  `active · 3/4`，等待 Root independent close-out。

## 验证事实

- `npx vitest run src/components/telegram-admin-tab.test.tsx src/i18n/catalog.test.ts src/i18n/ui-bilingual.test.tsx`：3 files、39 tests 通过。
- `npm test -- --run`：92 files、1213 tests 通过。
- `npm run build`：通过；仅有既有 chunk size warning，无 TypeScript/build error。
- 构建生成的 conformance projection 已恢复到 Git checkpoint；`git diff --check` 通过。
