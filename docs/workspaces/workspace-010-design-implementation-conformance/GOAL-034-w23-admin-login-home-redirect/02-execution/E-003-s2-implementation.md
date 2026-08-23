---
title: E-003 · S2 实施事实：挂具钉方言 + signInZh 去竞态
status: recorded
created: 2026-08-23
updated: 2026-08-23
parent: GOAL-034-w23-admin-login-home-redirect
version: 0.1.0
---

# E-003 · S2 实施事实（2026-08-23）

## 改动清单（仅 2 文件，无产品代码）

1. `apps/web/playwright.config.ts` —— API webServer env 增 `DB_DIALECT: "sqlite"`（显式钉死，含 W23 事故注释）。效果：无论本机 `configs/.env`/进程 env 如何，e2e API 永远使用挂具临时 SQLite + `DB_PATH` + `ADMIN_INITIAL_PASSWORD=admin` 隔离语义。
2. `apps/web/e2e/localization.spec.ts` —— `signInZh` 由一次性 `isVisible()` 分支改为等待式三段流程（对齐 `e2e/sign-in.ts`）：
   - 先 `waitForURL(/\/dashboard$/, 8s)` 直接命中「种子已被替换」快路径；
   - 未命中则 `forced.waitFor(8s)` 等强制改密屏（消除异步往返竞态），完成改密后 `waitForURL(15s)`；
   - 仍未命中（初始密码已失效）才 fallback 共享密码 `admin-e2e-pass` 重登。

## 退出条件

- `npx playwright test`（admin / mvp 两 profile 全量）恢复全绿；
- `vitest` / `tsc -b` + `vite build` / `go test ./...` 不受影响（改动面 = e2e 挂具与 e2e spec，产物代码零 diff）。