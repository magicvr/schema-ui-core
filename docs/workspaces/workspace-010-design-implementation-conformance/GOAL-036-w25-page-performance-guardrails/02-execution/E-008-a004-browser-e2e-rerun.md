---
title: E-008 · A-004 建议 2 浏览器层可选背书——sqlite e2e 双 profile 复跑
status: recorded
date: 2026-08-24
created: 2026-08-24
updated: 2026-08-24
parent: GOAL-036-w25-page-performance-guardrails
version: 0.1.0
scope: GOAL-036 关门后可选复验（A-004 建议 2 / A-005 备注承接）；sqlite × admin/mvp 浏览器 e2e 全量回归事实记录
---

# E-008 · sqlite e2e 双 profile 复跑（浏览器层可选背书）

## 背景与定位

- 来源：`03-audit/A-004-correction-recheck-independent.md` 建议 2——「若需浏览器层最终背书，可选复跑一次 sqlite e2e 双 profile（非必需——机制层已被单测栅栏覆盖）」；A-005 备注记录「如用户需要可另约复跑」。**2026-08-24 用户发起本次复跑**，本条目记录执行事实。
- 性质：关门后可选背书（non-gate），**不重开门禁、不改 `done 6/6` 状态**；I-001 已 closed 的结论不受影响。

## 执行事实（2026-08-24 07:35–07:40）

命令（`apps/web`）：`APP_PROFILE={admin|mvp}` + `DB_DIALECT=sqlite`（方言契约默认，`configs/.env` 无法改道）→ `npm run test:e2e`（Playwright，每轮全新隔离临时 SQLite + `reuseExistingServer: false`）。

| Profile × Dialect | 结果 | 跳过 | 耗时 | Exit |
|---|---|---|---|---|
| `APP_PROFILE=admin` × sqlite | **9 passed / 0 failed** | 1（mvp 分支定位用例：`localization.spec.ts:120`） | 2.5m | **0** |
| `APP_PROFILE=mvp` × sqlite | **9 passed / 0 failed** | 1（admin 分支 settings projection 用例：`localization.spec.ts:56`） | 1.7m | **0** |

要点：

- 与基线（E-004 / I-001：各 **9 passed / 1 profile 专属跳过 / exit 0**）完全一致，互为镜像跳过符合预期，**无回归**。
- `schema-crud`（含删用户→删角色 CASCADE 断言，F-001 修复的直接证据面）双 profile 均完整通过。
- 收尾干净：25080/25173 端口释放，无残留服务。

## 原始日志

- `attachments/e2e-rerun-sqlite-admin.log`
- `attachments/e2e-rerun-sqlite-mvp.log`

## 结论

浏览器层可选背书成立：W25 连接面/渲染层改动在真实栈（Go API + Vite dev + Chromium）上双 profile 全绿。GOAL-036 维持 `done 6/6` 不变，无新增 finding。
