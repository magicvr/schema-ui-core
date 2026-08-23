---
title: D-001 · S1 设计冻结：收尾层双数据库方言矩阵
status: frozen
created: 2026-08-23
updated: 2026-08-23
parent: GOAL-035-w24-e2e-dual-dialect-matrix
version: 0.1.0
---

# D-001 · S1 设计冻结

## 背景与结论

- GOAL-034（W23）初始修法「钉死 `DB_DIALECT=sqlite`」被用户复审否决：收尾层 e2e 不该靠「防改道」自保，而应**两个方言各测一次**——目前显然只测了一个（sqlite 默认；历史上唯一一次「pg」还是被 `.env` 意外改道到共享开发库，无有效覆盖）。
- 实验先证：专用 scratch Postgres（全新 + `ADMIN_INITIAL_PASSWORD` 种子）跑全量浏览器 e2e **9/9 绿** → 产品 pg 方言无缺陷，缺口纯在挂具层（无方言感知、无 provisioning、无校验）。
- 修复方向：把方言变成挂具**显式契约**，postgres 成为一等公民（复用 pgtest / CI api-postgres 的 scratch-DB 先例），并加启动后契约校验 fail-fast。

## 决策表

| # | 决策 | 依据 |
|---|------|------|
| D1 | `playwright.config.ts` 声明 `DB_DIALECT` 契约：默认 `sqlite`；`postgres` 需显式设置。webServer env **始终显式写入**该值 → process env 优先于 `configs/.env`，`.env` 在结构上无法再静默改道 | W23 N-001 根因；config.Load 优先级（process env > env-file > YAML > default） |
| D2 | postgres 模式 provisioning：`apps/api/cmd/e2e-pgset`（pgx，`create`/`drop`/`verify`）每轮 create `schema_ui_e2e_<nonce>` → 跑 → globalTeardown drop；凭据来源 = process env 或 `configs/.env` 的 `DB_*`（与 API 进程一致；pgtest 同源先例） | CI api-postgres / pgtest 已验证的 scratch-DB 模式；本地「任意方言可跑」的落点 |
| D3 | `e2e/global-setup.ts` 启动后校验契约：sqlite = `DB_PATH` 文件出现；postgres = scratch 库 `schema_migrations` 就位（`verify`）。违反 → 立即失败并给出 W23 指向性诊断 | 杜绝「3 个 spec 后神秘 401」类失效；校验是契约的最后防线 |
| D4 | CI `browser-e2e` 扩为 `profile × dialect [sqlite, postgres]` 矩阵（postgres 用 service 容器 + `setup-go` 已有；`DB_*` 经 GITHUB_ENV 注入 pg 腿） | 用户的「双方言都测」要求在 CI 层永久生效 |
| D5 | npm `test:e2e:postgres`（跨平台小 runner 设 env 后 spawn playwright）+ README 双方言说明 | 本地一命令可达 |
| D6 | 关门审计 `self` | 无产品语义/协议变化；CI 基础设施面变更由回归 + 静态校验覆盖 |

## 未选方案

- **保留「对抗式钉值」**：仅在 sa 层防御 `.env`，掩盖收尾层从未覆盖 pg 方言的缺口（用户已否决）。
- **自动探测本机 pg 并自动切换**：隐式行为，违背「显式契约」原则；探测不可靠时仍要 fail-fast，不如直接显式。
- **去掉 sqlite 默认值强制双方言都显式**：破坏 CI/本地默认路径的可预测性，无收益；sqlite 默认 + pg opt-in 已满足矩阵语义。