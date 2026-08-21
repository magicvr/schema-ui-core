---
id: E-003-w8-recommended-disposition
doc: execution-entry
goal: GOAL-008-w8-api-web-security-audit
date: 2026-08-20
status: recorded
parent: GOAL-001-production-hardening
created: 2026-08-20
updated: 2026-08-20
version: 0.1.0
---

# E-003 · A-003 recommended 处置

## 事实

- A-003 独立审计在通过 closing 时提出 2 条 low recommended：F-001（`apps/web/src/main.tsx` 残留旧 inline 注释）与 F-002（未做真实浏览器 CSP 回归）。
- F-001：已顺手 fixed，`main.tsx` 注释统一为外部 `/theme-init.js` 事实。
- F-002：登记为后续维护项，不阻断关门；建议在后续 Web 浏览器回归/E2E 波次补真实 CSP 响应头检查。

## 证据

- `apps/web/src/main.tsx:37-40`：无 inline 字样，注释为 external script。
- `apps/web/nginx.conf:29`：`script-src 'self'` 不变；`apps/web/public/theme-init.js` 与 `dist/theme-init.js` 存在。