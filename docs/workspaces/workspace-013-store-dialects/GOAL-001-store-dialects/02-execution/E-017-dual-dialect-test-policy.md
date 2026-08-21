---
id: E-017
doc: execution-entry
goal: GOAL-001-store-dialects
status: recorded
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
---

# E-017 · 双方言测试策略：凭据 env 化 + CI 双跑（sqlite & postgres）

## 2026-08-20 · 测试面与 CI 增强（关门后维护）

### 已发生事实

- **测试凭据 env 化**（新增 `apps/api/internal/pgtest`）：
  - `pgtest.DSN()` 从环境/`configs/.env` 读取后拼接 PG 测试 DSN；优先级 `SCHEMA_UI_R2_PG_DSN`（legacy 别名）> `PG_TEST_DSN` > `PG_TEST_HOST/PORT/USER/PASSWORD/DB/SSLMODE`；未配置返回空 → 门控测试 SKIP。
  - 会读取 gitignored 的 `apps/api/configs/.env` 中 `PG_TEST_*`（cwd 无关，按仓库根定位），本地开发者无需 shell 包装。
  - **全部 20 处**门控测试引用切换到 `pgtest.DSN()`（store + composition）。
- **实际双跑验证（真实 PG 192.168.31.213，仅走 `PG_TEST_*`，不带 legacy 变量）**：`go test ./...` **0 FAIL**——每个 DB 相关套件在 sqlite（默认，全量）与 postgres（门控，boot/启动/共事务/迁移原型/备份）下都过。
- **config.yaml 稳定/机密分离**：`configs/config.yaml` 与 `config.default.yaml` 的 `db:` 块只保留长期稳定、非机密项（`dialect`/`path`/`name`/`sslmode`）；**host/port/user/password 一律走 env**（`DB_HOST/DB_PORT/DB_USER/DB_PASSWORD`，host/port 有默认），`db.dsn`/`DB_DSN` 可整串覆盖；新增 `configs/env.example` 模板。
- **CI 双跑**（`.github/workflows/r6-basic-matrix.yml` 新增 `api-postgres` job）：
  - 既有 `api` job：`go test ./...`（sqlite 默认，PG 门控 SKIP）。
  - 新增 `api-postgres` job：`postgres:15-alpine` 服务容器 + `PG_TEST_*` 环境变量 → 同套 `go test ./...` + `go vet` + `go build`（superuser 建/删 scratch 库）。
- **未来 DB 测试约定**：需要 PG 侧的测试用 `pgtest.DSN()` 门控；sqlite 侧由默认全量回归覆盖；CI 两个 job 保证「sqlite & postgres 都要验证」自动成立。

### 证据

| 主张 | 路径 / 命令 / commit |
|------|----------------------|
| pgtest helper | `apps/api/internal/pgtest/pgtest.go` |
| 门控引用更新 | `internal/store/postgres_test.go`、`internal/composition/{postgres_startup,config_driven_postgres}_test.go` |
| 真实双跑 | `PG_TEST_*` 下 `go test ./...` 0 FAIL（192.168.31.213） |
| config 拆分 | `apps/api/internal/config/config.default.yaml`、`apps/api/configs/config.yaml`、`apps/api/configs/env.example` |
| CI 双 job | `.github/workflows/r6-basic-matrix.yml`（`api` + `api-postgres`） |
