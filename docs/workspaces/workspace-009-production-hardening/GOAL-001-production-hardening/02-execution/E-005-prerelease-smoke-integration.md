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
- `scripts/smoke.sh` `SM-004` 内置 W16-F01 首登强制改密：检测到 `mustChangePassword` 时走真实 `/api/account/password` 改密并切换新密码（`SMOKE_PASSWORD_NEW` 可选，默认 `<SMOKE_PASSWORD>-changed`），不作为测试预处理；后续步骤以新密码继续，smoke 因此覆盖真实首登旅程。
- 新增 `scripts/pre-release-smoke.sh` 发版前一键冒烟：独立 Compose project 构建/启动（仅本机 loopback 临时发布 API 端口供 smoke 使用）→ `smoke.sh --disposable + SMOKE_CSP=1`（改密交给 smoke 内完成）→ `down -v` 清理。
- 实测端到端运行：`SM-001~005 + SM-007 + SM-006 + SM-008` 全 PASS（`SM-004=PASS（W16-F01 首登强制改密已执行）`），`PRE-RELEASE SMOKE RESULT: PASS`。
- 文档：`QUICKSTART.md` / `README.md` 更新「发版前完整冒烟」用法与退出码 7。
- CI：`.github/workflows/r6-basic-matrix.yml` container-smoke 增加 `SMOKE_CSP=1` 与 `npm ci` + Playwright Chromium 安装（若该 job 在 W16-F01 fresh seed 上受强制改密影响，需按 wrapper 同样做改密 bootstrap 后再启用——已提示用户/维护者复核）。

## 证据

- `scripts/pre-release-smoke.sh`（隔离栈 + 临时 API loopback override；W16-F01 改密由 SM-004 内真实完成）
- `scripts/smoke.sh` SM-008 段
- `.github/workflows/r6-basic-matrix.yml`（container-smoke 职位）
- 运行输出：`SMOKE RESULT: PASS`