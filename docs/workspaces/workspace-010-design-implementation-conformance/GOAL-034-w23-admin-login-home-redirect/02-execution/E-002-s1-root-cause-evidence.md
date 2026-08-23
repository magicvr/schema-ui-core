---
title: E-002 · S1 根因定位事实链
status: recorded
created: 2026-08-23
updated: 2026-08-23
parent: GOAL-034-w23-admin-login-home-redirect
version: 0.1.0
---

# E-002 · S1 根因定位事实链（2026-08-23）

## 复现（全量 admin e2e，2.6 分钟）

`APP_PROFILE=admin npx playwright test`（apps/web）：

- 10 用例：6 失败 / 3 通过 / 1 跳过；失败的 6 个（`force-password-change` → `host-failure:71` → `localization:34` → `schema-crud` → `shell:55` → `w4-long-content`）**全部**卡在登录门禁。
- 首个用例在**全新种子**上 `admin/admin` 提交即 401：快照 `alert: invalid username or password`（`error.unauthorized`），强制改密屏从未出现（`force-password-change.spec.ts:21` 超时）。
- 其余用例经 `sign-in.ts` fallback 以 `admin-e2e-pass` 提交同样 401（`sign-in.ts:50 waitForURL(/\/dashboard$/)` 超时，页面停留在登录表单）。

## 隔离失效证据

1. `%TEMP%\schema-ui-e2e-*`（两次运行 10:21 / 10:22 建立）与 `%TEMP%\schema-ui-e2e-cfg-*` 目录内**均无 `e2e.db`**：`DB_PATH` 未被使用。
2. `apps/api/configs/.env`（**gitignored**，CreationTime/LastWriteTime = 2026-08-21 08:05:50）内容：`DB_DIALECT=postgres`、`DB_HOST=192.168.31.213`、`DB_NAME=postgres` 等。
3. `config.Load`（config.go:314-318、835-850）：`CONFIG_ENV_FILE`（默认 `configs/.env`，相对 API 进程 CWD）「never overrides an already-set process env」——挂具未设 `DB_DIALECT` → `.env` 生效 → 方言 postgres；`DB_PATH` 仅 sqlite 语义，被忽略。
4. 仓库根 `e2e-baseline.log`（W22 期间遗留）：stash 后 API 以 postgres 方言启动直接报「unknown applied migration version 49」而 **webServer 启动失败**（exit 1）——W22「基线实验：stash 后 HEAD 同败（:62）」与自身日志矛盾；该实验无法移除 gitignored `.env`，结论「先于 W22 的路由回归」不可成立。

## 凭证级复现（API 直连）

- 环境 A（挂具同语义，未钉方言）：临时 config + `ADMIN_INITIAL_PASSWORD=admin` → `POST /api/auth/login {admin,admin}` → `401 UNAUTHORIZED`。
- 环境 B（仅加 `DB_DIALECT=sqlite`）：同参数 → `200`，返回 `accessToken`/`refreshToken`，`user.mustChangePassword=true`；强制改密→业务 API 链路完整示意（go 单测同源覆盖）。

## 结论

- 根因 = e2e 挂具未隔离 store 方言；N-001 观察到的「停留 `/`」是 401 后停留在登录页的必然结果，与 home 推导/路由代码无关（`resolveInitialRoute`/`StampHomePageRef`/manifest home 字段静态核查均无恙）。
- 修复面：仅 `apps/web/playwright.config.ts`（钉 `DB_DIALECT=sqlite`）+ `apps/web/e2e/localization.spec.ts`（signInZh 竞态，消除种子正确时仍可能误入 fallback 的分支抖动）。
- 影响判定：无产品代码 / 协议 / manifest 契约改动 → 按 00-meta 边界声明，关门审计模式 `self`。