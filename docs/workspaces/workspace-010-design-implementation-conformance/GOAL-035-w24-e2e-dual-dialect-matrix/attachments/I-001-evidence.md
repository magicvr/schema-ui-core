---
title: I-001 证据固化（A-002 F-002 响应）
status: recorded
created: 2026-08-23
updated: 2026-08-23
parent: GOAL-035-w24-e2e-dual-dialect-matrix
version: 0.1.0
---

# I-001 证据固化 · 「专用 Postgres 上浏览器 e2e 全绿」实验与回归记录

> 响应 [A-002 F-002](../../GOAL-035-w24-e2e-dual-dialect-matrix/03-audit/A-002-w24-fix-review-independent.md)：原始 `*.log` 被 `apps/web/.gitignore:16` 忽略、不随仓库传播；本文件是**仓库内持久化摘要**。原始日志保留在本机（见下），CI main 首跑后新增的矩阵运行证据回填于本文件「CI 回填」节（F-003 触发器）。

## 1. I-001 先证实验（W23 复审 → W24 立项前置）

- 日期：2026-08-23（上午）｜执行人：编排器（W23 用户复审调查）
- 目的：证明「浏览器 e2e 在产品 postgres 方言下能否全绿」——决定缺口在挂具层还是产品层。
- 方法：专用 scratch 库 **`schema_ui_e2e_w23_242e2651`**（全新、无数据、`ADMIN_INITIAL_PASSWORD=admin` 种子），`APP_PROFILE=admin` 全量浏览器 e2e。
- 结果：**10 tests → 1 skipped（mvp-only）→ 9 passed（1.8m），exit 0**。
- 结论：产品 pg 方言无缺陷；缺口在挂具层（无方言感知 / 无 provisioning / 无校验）。→ 支撑 GOAL-035 D-001 与 I-001 closed。
- 本机日志：`apps/web/e2e-w23-pg-experiment.log`（gitignored）。

## 2. W24 最终回归（实现后，最终代码态 018a4d7）

- sqlite 腿：`npm run test:e2e` → **9 passed（1.7m）**，exit 0。本机日志 `apps/web/e2e-w24-sqlite.log`。
- postgres 腿：`npm run test:e2e:postgres` → **9 passed（1.8m）**，exit 0；scratch 库 **`schema_ui_e2e_mt58t6po311j6l`** 于 teardown 可见 `dropped`；事后 `go run ./cmd/e2e-pgset list` = **0 遗留**（单建单删，F-1 配置双载守卫生效）。本机日志 `apps/web/e2e-w24-postgres2.log`。
- 密闭回归（同轮）：vitest **1088/1088**；`go test ./...` 全绿；tsc+build exit 0。日志 `e2e-w24-vitest.log` / `e2e-w24-build.log` / `e2e-w24-go.log`（apps/api 下）。

## 3. A-002 独立复审现场复跑（independent，2026-08-23）

- sqlite 腿：9 passed（2.1m），exit 0——本机 `.env` 陷阱（`DB_DIALECT=postgres`）仍存在时成立（契约防御实测）。
- postgres 腿：9 passed（1.8m），exit 0；scratch **`schema_ui_e2e_mt5a0ht1li16q8`** 可见 dropped；事后 `list` 0 遗留。
- 详见 A-002 成果表（现场复跑 job 输出）。

## 4. CI 回填（F-003 · 已触发并完成）

- **PR**：https://github.com/magicvr/schema-ui-core/pull/5（dev→main，71 commits）
- **工作流运行**：https://github.com/magicvr/schema-ui-core/actions/runs/32617287887（event=pull_request，headSha=916a280，conclusion=success）
- **运行时间**：2026-08-23（PR 合入 `cdb2308`，04:16Z）
- **矩阵结果（9/9 jobs SUCCESS）**：

| job | 结果 |
|-----|------|
| web (Linux, Node 22) | SUCCESS |
| api (Linux, Go 1.26) | SUCCESS |
| api + postgres (Linux, Go 1.26) | SUCCESS |
| browser E2E (mvp / sqlite) | SUCCESS |
| browser E2E (mvp / postgres) | **SUCCESS**（CI 首次真实跑 pg 腿） |
| browser E2E (admin / sqlite) | SUCCESS |
| browser E2E (admin / postgres) | **SUCCESS**（CI 首次真实跑 pg 腿） |
| container smoke (mvp, docker compose) | SUCCESS |
| container smoke (admin, docker compose) | SUCCESS |

- 状态：**F-003 closed（fixed）**——`browser-e2e` 矩阵首跑证据已入库（上方回填节）。若未来矩阵变更或失败，在此节追加记录。