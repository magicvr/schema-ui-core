---
id: E-004-w8-real-browser-csp
doc: execution-entry
goal: GOAL-008-w8-api-web-security-audit
date: 2026-08-20
status: recorded
parent: GOAL-001-production-hardening
created: 2026-08-20
updated: 2026-08-20
version: 0.1.0
---

# E-004 · 真实浏览器 + 生产 CSP 响应头回归检查

## 触发

用户 `/govern` 指令：跑一次真实浏览器 + 生产 CSP 响应头的回归检查；没有问题则进行 Root 层汇总。

## 已发生事实

- 用 `docker compose up --build -d` 构建并启动生产 Compose 栈（API `APP_ENV=production` + web 生产 nginx），服务地址 `http://127.0.0.1:25081`。
- 新增并运行 `apps/web/scripts/check-prod-csp.mjs`：headless Chromium 打开生产首页，核对响应头与真实浏览器主题 bootstrap 行为。
- 结果：**全部通过**。
  - `Content-Security-Policy` 响应头存在（长度 200），包含 `default-src 'self'` 与 `script-src 'self'`。
  - 生产 HTML 引用 `<script src="/theme-init.js"></script>`，无 inline theme bootstrap。
  - `/theme-init.js` 返回 `200 application/javascript`，由同源 `'self'` 放行。
  - 首屏脚本执行：`colorScheme` 被设置为 `light`/`dark`（本次为 `light`）。
  - 设置 `localStorage.theme=dark` 并 reload 后，`html.dark` 类被添加且 `color-scheme: dark`。
  - 浏览器控制台**无** `Content Security Policy` / `Refused to execute inline script` 错误。

## 结论

A-003 原 recommended F-002（真实浏览器 CSP 回归）已完成闭合；F-001/F-002 相关生产浏览器面无回归。

## 证据

- 脚本：`apps/web/scripts/check-prod-csp.mjs`
- 运行输出：`{"ok":true,...}`（7 项检查全绿）
- Compose 生产栈：`compose.yaml` 构建并启动成功。