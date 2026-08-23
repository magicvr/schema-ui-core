---
title: E-002 · S2 实施事实：方言契约 + provisioning + 校验 + CI 矩阵
status: recorded
created: 2026-08-23
updated: 2026-08-23
parent: GOAL-035-w24-e2e-dual-dialect-matrix
version: 0.1.0
---

# E-002 · S2 实施事实（2026-08-23）

## 改动清单

1. **`apps/api/cmd/e2e-pgset/main.go`**（新工具）：`create` / `drop` / `verify` 三个动作；连接凭据 = process env 优先，其次 `apps/api/configs/.env` 的 `DB_*`（与 API 自身加载同源，pgtest 同款模式）；`verify` 检查 scratch 库 `public.schema_migrations` 是否存在（API 完成迁移的信号）。
2. **`apps/web/playwright.config.ts`**：`DB_DIALECT` 显式契约（默认 sqlite，非法值启动即抛）；postgres 模式在配置加载时经 `go run ./cmd/e2e-pgset create` 生成 `schema_ui_e2e_<nonce>` 并写入 webServer env `DB_NAME`；`E2E_DB_DIALECT` / `E2E_DB_PATH` / `E2E_PG_NAME` 传 `globalSetup/Teardown`；`globalSetup` + `globalTeardown` 接入。
3. **`apps/web/e2e/global-setup.ts`**：API 就绪后校验契约（60s 超时）：sqlite = `DB_PATH` 文件出现；postgres = `e2e-pgset verify` 通过；失败给出 W23 指向性诊断（不再出现「3 个 spec 后神秘 401」）。
4. **`apps/web/e2e/global-teardown.ts`**：postgres 模式 drop scratch 库（best-effort；遗留可用 `go run ./cmd/e2e-pgset drop` 清理）。
5. **`apps/web/scripts/run-e2e-postgres.mjs`** + package.json `test:e2e:postgres`：跨平台设 `DB_DIALECT=postgres` 后跑 playwright。
6. **`apps/web/README.md`**：双方言说明（默认 sqlite / pg opt-in / 校验语义 / CI 矩阵）。
7. **`.github/workflows/r6-basic-matrix.yml`**：`browser-e2e` 扩为 `profile × dialect [sqlite, postgres]`；postgres service 容器 + `setup-go`（已有）+ pg 腿经 `$GITHUB_ENV` 注入 `DB_HOST/PORT/USER/PASSWORD/SSLMODE`。

## 契约不变量（防回退核心）

- API webServer env **始终**含显式 `DB_DIALECT`（process env 优先于 `configs/.env`）→ `.env` 无法改道；
- `globalSetup` 校验不通过即整体失败 → 挂具契约断裂时第一时间暴露；
- CI 矩阵使 sqlite/postgres 两条腿每次 PR 都执行。