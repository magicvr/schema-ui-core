---
title: 决策 · Shell 导航与 Schema fixture 洁净度
status: done
created: 2026-08-04
updated: 2026-08-04
parent: GOAL-001-production-admin-foundation
version: 0.2.0
---

# 决策 · GOAL-012

## D-001 · 承接 Root A-005 required finding 并立项

- **日期**：2026-08-04
- **状态**：accepted
- **用户裁决**：独立审计若发现阻断问题 → 在 workspace-002 新设子目标修正，并回退 Root 关门状态（本会话指令）。
- **决定**：以本目标承载 A-005 F-001（default Shell 死链）的修复与回归门禁；Root 由 `done` 回退为 `active`（派生进度保持 `5/5`），直至 F-001 合法闭合。
- **理由**：VP-002 要求可直接 fork 使用的生产级 Admin 基架；默认导航中存在必然失败的 pageId 属于产品面阻断，不能以「核心 CRUD 已通」掩盖。
- **未选方案**：
  - **仅聊天记录不立项**：违反 P-003 落盘与闭环。
  - **立刻补完整业务页（Activity/Settings 真功能）**：超出 A-005 范围与本 VP 非目标。
- **后续**：S1 消除死链 → S2 一致性测试 → S3 回归 → S4 关门；S5 文档可选。

## D-002 · S1 采用「移除占位」而非补假 Schema

- **日期**：2026-08-04
- **状态**：accepted
- **用户指令**：推进 GOAL-012（采纳立项时推荐路径）。
- **决定**：从 checked-in `app-manifest.json` 删除 `activity` / `settings` 的 `pages[]` 项，删除 sidebar「Workspace」组与 `user` 区 Settings 链接；**不**新增 minimal fixture 假页面。
- **理由**：占位入口无产品语义；补空页会制造第二套「能开但无能力」噪声。移除后 manifest 仅保留有 embed 文档的页面。
- **影响**：`I-012-001` → closed；S1 实施放行。
- **未选方案**：补 minimal text 页（保留死路由占位，产品价值为零）。
