---
title: E-003 · R-001 处置——e2e 网络层 Bearer 冒烟并入本波
status: active
created: 2026-08-26
updated: 2026-08-26
parent: null
version: 1.0.0
---

# E-003 · R-001 处置——e2e 网络层 Bearer 冒烟并入本波

用户指令（2026-08-26）：「把 r-001 并入本波。处理完再关门。」按 [A-001](../../03-audit/A-001-w14-self-closeout.md) §2 的 recommended R-001 实施。

## 变更

| 文件 | 变更 |
|------|------|
| `apps/web/e2e/schema-auth-transport.spec.ts` | 新增 Playwright 冒烟：真实 Chromium + 真实网络路径 |

## 用例设计

1. 复用 `sign-in.ts` 真实登录流（含首登强制改密分支）；
2. 登录完成后挂请求/响应采集器，`goto("/users")` 触发全新文档启动（restoreSession → AuthGate → App → loadPageDocument）；
3. `waitForResponse` **确定性等待**首个 `/api/schema/*` 响应：断言 HTTP 200 且请求头 `Authorization: Bearer …`；
4. `expect.poll` 收敛后对**全部**捕获的 schema 请求逐条断言 Bearer 头 + 200（防「零捕获即空过」的正向控制内置）；
5. 断言失败面标题（`can't be displayed | 无法显示此页面` 双语）count = 0——页面渲染真实内容而非 PageSchemaErrorSurface。

## 迭代事实

- 首版在导航 `load` 后立即断言，与 SPA mount 后的异步取数竞态（requests 已见、responses 为 0）→ 失败；
- 改为 `waitForResponse` + `expect.poll` 收敛后复跑 → 通过。

## 验证记录

| 验证 | 结果 |
|------|------|
| 单 spec：`npx playwright test e2e/schema-auth-transport.spec.ts` | **1 passed（9.6s）** |
| 全量 e2e 套件（chromium · workers=1 · 自起 API+vite） | **10 passed / 1 skipped（admin 档用例于 mvp 档按设跳过）· exit 0**（含本 spec `ok`；首轮暴露的 shell.spec 陈旧契约另见 [E-004](E-004-w14-shell-spec-contract-fix.md)） |
| 单元回归锁不受影响 | vitest 侧无改动，1130/1130 维持 |

## 备注

- 本机 dev 栈（25080/25173）在运行前已被用户侧停止，端口让渡给 playwright webServer（其自起全新 SQLite API，互不污染）；e2e 结束后未自动拉起 dev 栈。
