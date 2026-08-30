---
title: E-002 · W15 S4 Web 修正实施与回归（F-005～F-006）
status: active
created: 2026-08-30
updated: 2026-08-30
parent: null
version: 0.1.0
---

# E-002 · W15 S4 Web 修正实施与回归（F-005～F-006）

日期：2026-08-30 · checkpoint：`609cd6d6`

## 实施事实

1. **F-005 · 邀请 token URL 清理**（`apps/web/src/components/invite-accept.tsx`）：
   - 挂载时读取 `token` 后立即 `history.replaceState` 删除 query 中的 `token`（地址栏/浏览历史/截图不再携带一次性 bearer）；成功跳转路径不变；读取失败仍回退空 token。
2. **F-006 · 测试 fixture 根统一**（13 个失败 suite + README）：
   - 全部 `apps/api/internal/modules` 引用切换为 canonical 根 `apps/api/modules`：`all-module-schemas-dval`、`load-page`（注释）、`custom-components`、`error-localization`、`row-action-bindings`（深度路径）、`ui-bilingual`、`schema-keys`（fragment 清单 10 条）、`s5-denominator-render`、`representative-pages`、`schema-dictionary-entries`、`schema-crud`、`wallet-navigate`、`startup-config`、`representative-pages.integration`。
   - `apps/web/README.md` 陈旧路径（`internal/modules/schemarender/schema`，目录已不存在）更新为 canonical 根 + dev examples。
   - 新增 `src/protocol/fixture-root.guard.test.ts`：断言 canonical 根存在 + `src` 下无文件再引用退役路径（A-001「CI 路径存在性检查」落地为 vitest 锁；guard 自身豁免）。

## 回归验证

- 基线（修复前）：`13 failed | 75 passed (88)` / `76 failed | 1081 passed (1157)`（与本波 A-001 F-006 记录一致）。
- 修复后：**89 passed (89) / 1183 passed (1183)**；`tsc -b` exit 0；`vite build` exit 0（仅既存 bundle 体积警告）。
- 中途一次并发 vite build 与 vitest 采集竞争导致 2 例假失败；串行复跑即全绿（非代码问题）。

## 状态

F-005～F-006 实现 + 回归完成；闭合证据待独立审计后在 03-audit 响应节正式标记 fixed。