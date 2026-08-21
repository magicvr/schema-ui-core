---
id: E-018
doc: execution-entry
goal: GOAL-001-store-dialects
status: recorded
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
---

# E-018 · 本地 .env 配置（测试 PG）+ CI Docker postgres 服务模拟验证

## 2026-08-20 · 本地凭据落盘 + CI 拉取容器路径确认

### 已发生事实

- **本地 `apps/api/configs/.env`**（gitignored，`apps/api/.gitignore:.env`，`git ls-files` 未跟踪）：按用户提供的本地测试 PG `192.168.31.213:5432`（用户 `sa`/`Ss.110110`）写入 `DB_*` 与 `PG_TEST_*`（DB_DIALECT=postgres、DB_NAME=postgres、PG_TEST_DB=postgres）。
- **纯 `.env` 驱动验证**（不设任何 shell 变量，pgtest 自动读 `configs/.env`）：
  - `TestFullCatalogPostgresBootstrapIntegration` `-v` 确认 **PASS（非 SKIP）**；
  - `go test ./...` **0 FAIL**（sqlite 全量 + PG 门控全跑，目标 = 192.168.31.213）。
- **CI Docker 服务路径模拟**：`docker pull postgres:15-alpine`（GitHub Actions 会做同样的拉取）→ 起 `ci-pg-sim` 容器（POSTGRES_USER/PASSWORD/DATABASE=postgres）→ 设 `PG_TEST_*=127.0.0.1:55432/postgres` 跑 PG 门控套件，**全部 PASS**（store + composition）。
- `.github/workflows/r6-basic-matrix.yml` 的 `api-postgres` job 与本次模拟**完全一致**：`services: postgres:15-alpine` + `PG_TEST_*` → `go test/vet/build`；本仓库 CI 的 postgres 依赖来自 **Docker 拉取的服务容器**，与本地的 192.168.31.213 无关（CI 访问不到它，也不需要）。

### 证据

| 主张 | 证据 |
|------|------|
| .env 已写且未跟踪 | `apps/api/configs/.env`；`git check-ignore` / `git ls-files` 空 |
| .env 驱动全量 0 FAIL | `go test ./...`（无 shell env）|
| PG 门控 PASS（.env） | `TestFullCatalogPostgresBootstrapIntegration -v` PASS |
| CI 同款 docker 服务 PASS | `postgres:15-alpine` 容器 + `PG_TEST_*`，store/composition 全 PASS |
| CI job 已具备 | `r6-basic-matrix.yml` `api-postgres`（postgres:15-alpine service）|
