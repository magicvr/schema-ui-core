---
doc_type: goal-execution
id: E-032-post-close-operator-page-scroll-containment
parent: GOAL-001-telegram-operator-console
date: 2026-09-05
source: self
status: done
version: 0.1.0
---

# E-032 · 人工会话页面级滚动隔离（2026-09-05）

## 已发生事实

- 根据用户反馈修正“打开页面后消息较多时页面自身出现滚动条”的问题：Telegram operator
  route 现在由应用 shell 的固定视口高度链承载，页面本身不再因消息或会话数量增长而溢出。
- `App.tsx` 对 `telegram-operator` 使用独立 capability route 的 contained shell：根 shell、shell
  body、`main`、page region 和 `PageSurface` 均补齐 `h-full`/`min-h-0`/`flex-1`/`overflow-hidden`
  接缝；普通页面继续保留 page-level `overflow-y-auto`。
- Telegram operator schema 将 operator surface 直接挂到 `custom` body，去除会破坏高度链的中间
  `SectionView` wrapper；`telegram-admin-tab` 内的 sessions 列表、消息列表继续分别承担内部纵向
  滚动，composer 保持在面板底部可见。
- Playwright 增加 `APP_PROFILE=custom` 测试配置与真实 Chromium E2E：用 80 个 sessions、120 条
  长消息验证 document/body/root/main 无垂直溢出、sessions/message list 可滚动且滚动不改变页面
  scroll；另用 100 个 Users 验证普通长页面仍保持 page-level 滚动。
- 代码变更已在 checkpoint `9e9102cb166d4b2acb7b2ad8b844ce11ba55c188` 提交，提交信息为
  `fix(telegram): contain operator page scrolling`。

## 审计与验证事实

- A-011 `subagent (gpt-5.6-sol · reasoning medium)` 对当前代码 checkpoint 独立复审为 `pass`、
  `open_required: 0`、`open_recommended: 0`；确认 operator route 高度/滚动链、刷新期间消息保留、
  schema 直接 custom body、普通页面回归和真实 Chromium 布局证据均可核对。
- Web 定向测试：`App.integration.test.tsx` 与 `telegram-admin-tab.test.tsx` 共 27 个测试通过；
  API Telegram schema 测试通过；Web 全量 `npm test` 为 92 个文件、1215 个测试通过。
- Web `npm run build`、E2E TypeScript 编译和 `APP_PROFILE=custom npm run test:e2e --
  telegram-operator-layout.spec.ts` 均通过；Chromium E2E 为 2 个测试通过。API `go test ./...`
  通过；构建仅保留既有 large chunk warning，自动生成的 conformance projection 未纳入提交。
- `git diff --check` 通过；实现 checkpoint 与本次治理记录均使用显式路径提交。

## 状态边界

本条是 Root 关门后的局部 UI 可用性修正，不改变 Root 六项成功标准、progress 或 status，因此不
重新打开 `GOAL-001-telegram-operator-console`/workspace-033；VP-033 仍保持 `active`。
