---
title: E-001 · 报障诊断与 hotfix 时间线
status: active
created: 2026-08-26
updated: 2026-08-26
parent: null
version: 1.0.0
---

# E-001 · 报障诊断与 hotfix 时间线

按时间线记事实（2026-08-26，schema-ui-core 工作树）。

## 事实

1. **报障**：用户报告「所有页面都显示无法显示此页面」。
2. **服务面排查（均健康）**：vite dev `http://localhost:25173/` → HTTP 200；API `/readyz`（:25080）→ HTTP 200；`/.well-known/schema-ui/app-manifest.json` → HTTP 200（21 个页面条目正常）；SPA 路由 fallback ×4 → 200。浏览器错误页假设排除。
3. **关键线索**：错误文案是应用自身渲染的——`apps/web/src/i18n/messages/zh-CN.json:94` `shell.pageSchemaError.title = "无法显示此页面"`（`App.tsx` `PageSchemaErrorSurface`）。渲染链路 = manifest → `GET /api/schema/{pageId}` → D-VAL 校验 → RenderPage，任一步失败即整页报错。
4. **根因实测**：匿名请求 `/api/schema/{dashboard|users|account|settings|activity}` **全部 HTTP 401**（经 vite 代理，等价浏览器路径）。
5. **语义溯源**：提交 `b7954235`（GOAL-013 S3/S5）含「F-010 /api/schema 挂认证」——该端点自此时起要求 Bearer。
6. **装配缺口定位**：
   - `main.tsx` `AuthGate` 渲染 `<App>` 时只传 `resourceFetcher={resourceFetcher}`（authFetch 包装），**未传 `schemaFetcher`**；
   - `App.tsx` 将页面文档加载交给 `loadPageDocument`，其缺省 fetcher 为裸 `withTimeout()`（`lib/fetch-timeout.ts` 仅加超时，不附加任何头）→ 页面文档请求恒为匿名；
   - 全仓 34 处测试均显式注入 `schemaFetcher=`，生产入口无测试覆盖——测试全绿与生产断裂并存。
7. **hotfix（措施 A）**：经用户确认（处置问答选「现在就修」），`main.tsx` 装配补 `schemaFetcher={authFetch}`。验证：`tsc -b` exit 0；vite transform `/src/main.tsx` HTTP 200 且产物含新接线。端到端登录验证未执行（dev 口令非默认 admin/admin，且 GOAL-014 刚上线 IP 来源锁，避免触发）——由用户刷新浏览器验证，未再收到异常反馈。
8. **追加指令**：用户指示补防回归测试并在工作区9立项承载治理上下文 → 开本目标（GOAL-015 · W14），S3 实施见 [E-002](E-002-w14-regression-lock-implementation.md)。

## 备注

- hotfix 先于本目标存在，属事实先行；本目标自 S1 起追认承载其治理上下文，时间线如实分列，不回填虚构顺序。
- 改动文件（截至 E-002）：`apps/web/src/main.tsx`、新增 `apps/web/src/app/AuthGate.tsx`、新增 `apps/web/src/app/auth-gate.wiring.test.tsx`。均未提交 git。
