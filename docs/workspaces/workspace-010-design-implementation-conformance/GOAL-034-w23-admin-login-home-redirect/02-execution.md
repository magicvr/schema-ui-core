---
title: 执行索引 · GOAL-034-w23-admin-login-home-redirect
status: active
created: 2026-08-23
updated: 2026-08-23
parent: GOAL-034-w23-admin-login-home-redirect
version: 0.1.0
---

# 执行时间线 · GOAL-034

| E-ID | 日期 | 事件 | 状态 | 文件 |
|------|------|------|------|------|
| E-001 | 2026-08-23 | S1 立项：承接 [workspace-010 GOAL-033 A-001 N-001](../../GOAL-033-w22-residual-closeout/03-audit/A-001-w22-closeout-self.md)（admin 登录后停留 `/` 未跳 `/dashboard`；基线实验证实先于 W22 存在）；五件套建立，goal-tree/workspace.md 同步 | recorded | 本文件 |
| E-002 | 2026-08-23 | S1 根因定位：全量 admin e2e 6/10 失败且全部为登录 401；`%TEMP%\schema-ui-e2e-*` 无 SQLite 产物；本地 `configs/.env`（2026-08-21 建立）`DB_DIALECT=postgres` 劫持挂具 DB_PATH；凭证级复现 A（未钉方言）401 vs B（`DB_DIALECT=sqlite`）200 全链 | recorded | `02-execution/E-002-s1-root-cause-evidence.md` |
| E-003 | 2026-08-23 | S2 实施：`playwright.config.ts` 钉死 `DB_DIALECT=sqlite`；`localization.spec.ts` signInZh 竞态修复（等待式三段流程） | recorded | `02-execution/E-003-s2-implementation.md` |
| E-004 | 2026-08-23 | S3 回归：go / vitest 1088 / tsc+build 全绿；e2e admin 连续 5 轮全绿；连带整改 F-1（菜单 scroll-close 竞态，产品面）与 F-2（sign-in fallback 等待竞态，测试面）均 fixed 并验证 | recorded | `02-execution/E-004-s3-regression-evidence.md` |
| E-005 | 2026-08-23 | 关门后用户复审：强制 sqlite 属绕过，收尾层应双方言各测一次；专用 scratch pg 全量实验 9/9 绿（产品方言无缺陷）→ 正确修法移交 [GOAL-035-w24-e2e-dual-dialect-matrix](../../GOAL-035-w24-e2e-dual-dialect-matrix/00-meta.md) | recorded | 本文件 |
