---
id: E-025-pagination-fix-and-closeout-authorization
doc: execution-entry
status: recorded
goal_id: GOAL-001-admin-workflow-continuity
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-001-admin-workflow-continuity
version: 1.0.0
---

# E-025 · 开设 GOAL-011 并完成分页与跳转文案修正

2026-09-18，用户报告两个缺陷并授权「修改这两个问题后，授权走根目标关闭流程」（`D-017`），据此开设整改子目标 `GOAL-011-pagination-page-size-contract`（非纲领，不计入 Root 分母）并于同日完成 C1～C4：

- C1 复现与根因：改动前 2 条测试失败（`expected '10' to be '20'`、`expected 'Search' to be 'Go'`）；根因是前端 `DEFAULT_PAGE_SIZE = 10` 与服务端 `handler.DefaultPageSize = 20` 不一致，叠加 `buildResourceQuery` 在等于默认值时省略 `pageSize` 的规则；跳转按钮复用了 `feedback.search` 文案（`E-001`）。
- C2 修正：默认值统一为 20、选 10 时显式发送 `pageSize=10`、按钮改用新增键 `feedback.jumpToPage`（跳转 / Go）、清理 `render.tsx`/`activity-export.tsx`/`schema-table.tsx` 中把 10 当默认的字面量（`E-002`）。
- C3 回归与防复发：新增 `src/pagination-size-contract.guard.test.ts`（解析 `apps/api/internal/handler/resources.go` 的 `DefaultPageSize` 并与前端常量绑定；锁定非默认尺寸必须上线路；锁定两个 catalog 的文案键与取值），修正两处固化错误契约的旧断言，并在 `e2e/list-visual-surface.spec.ts` 增加真实浏览器用例（`E-003`）。
- C4：self 审计 `A-001` `pass`、开放 required = 0，`GOAL-011` 以 `done · 4/4` 关门（`E-004`）。

验证：Vitest **113 文件 / 1434 测试**、`npm run typecheck` exit 0、`list-visual-surface` e2e **3 passed**（admin 50.8s / mvp 1.2m）、`git diff --check` 通过；未改 `apps/api`。

用户授权的前置条件（「修改这两个问题后」）据此解除；Root/VP 关门流程见 `E-026`。
