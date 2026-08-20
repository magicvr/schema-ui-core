---
id: E-005
goal: GOAL-001-production-hardening
title: W8 CSP/真实浏览器冒烟纳入发版前流程
date: 2026-08-20
status: recorded
parent: null
created: 2026-08-20
updated: 2026-08-20
version: 0.1.0
---

# E-005 · W8 CSP/真实浏览器冒烟纳入发版前流程

## 已发生事实

- `scripts/smoke.sh` 新增可选 `SM-008`（`SMOKE_CSP=1` 时运行 `apps/web/scripts/check-prod-csp.mjs` 真实浏览器 + 生产 CSP 头检查），失败退出码 `7`。
- 新增 `scripts/pre-release-smoke.sh` 发版前一键冒烟：独立 Compose project 构建/启动（仅本机 loopback 临时发布 API 端口供 smoke 使用）→ admin 强制改密 bootstrap（W16-F01）→ `smoke.sh --disposable + SMOKE_CSP=1` → `down -v` 清理。
- 实测端到端运行：`SM-001~005 + SM-007 + SM-006 + SM-008` 全 PASS，`PRE-RELEASE SMOKE RESULT: PASS`。
- 文档：`QUICKSTART.md` / `README.md` 更新「发版前完整冒烟」用法与退出码 7。
- CI：`.github/workflows/r6-basic-matrix.yml` container-smoke 增加 `SMOKE_CSP=1` 与 `npm ci` + Playwright Chromium 安装（若该 job 在 W16-F01 fresh seed 上受强制改密影响，需按 wrapper 同样做改密 bootstrap 后再启用——已提示用户/维护者复核）。

## 证据

- `scripts/pre-release-smoke.sh`（含 forced-password bootstrap 与临时 API loopback override）
- `scripts/smoke.sh` SM-008 段
- `.github/workflows/r6-basic-matrix.yml`（container-smoke 职位）
- 运行输出：`SMOKE RESULT: PASS`