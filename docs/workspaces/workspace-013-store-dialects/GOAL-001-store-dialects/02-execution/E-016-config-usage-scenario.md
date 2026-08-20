---
id: E-016
doc: execution-entry
goal: GOAL-001-store-dialects
status: recorded
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
---

# E-016 · 配置面补强：config.yaml 切方言 + postgres 连接参数 + 密码走 env（真实库验证）

## 2026-08-20 · 使用场景确认（关门后维护补强）

### 已发生事实

- **需求**：`config.yaml` 可选 `sqlite|postgres`；可配置数据库连接参数；数据库连接密码只用 `.env`/环境变量。
- **补强**（`apps/api/internal/config`）：
  - `db.dialect`: `sqlite | postgres`（既有，空=sqlite，未知 fail-closed）。
  - 新增**展开式 postgres 连接参数**：`db.host / db.port / db.name / db.user / db.sslmode`（YAML 可配；env `DB_HOST/DB_PORT/DB_NAME/DB_USER/DB_SSLMODE` 覆盖；host/port/sslmode 有默认）。
  - **密码**：`db.password: ${DB_PASSWORD:-}` —— 一律从 `DB_PASSWORD` env / `configs/.env` 解析，**禁止硬编码到配置**；dialect=postgres 且未设 dsn 时，密码为空 → **fail-closed**（LoadError：请经 DB_PASSWORD 提供）。
  - `db.dsn`（`DB_DSN`）保留为**整串覆盖**：设置后优先于展开参数。
  - Load 末尾由展开参数 **`url.UserPassword` 安全拼装 DSN**（特殊字符自动转义）。
- **验证**：
  - config 单测：`TestDBPostgresExplodedParams`（env 密码拼 DSN ✓、缺密码 fail-closed ✓、dsn 覆盖 ✓、sqlite 忽略 ✓、env 覆盖 ✓）。
  - **config.yaml 驱动 + env 密码的真实 postgres 启动**：`TestCompositionPostgresConfigDriven`（live PG 192.168.31.213）——临时 `config.yaml` 声明 `dialect: postgres` + host/port/name/user/sslmode + `${DB_PASSWORD:-}`，密码经 `DB_PASSWORD` env 传入；`config.Load()` 拼出 DSN、`NewApp.Start` 全绿（ready 门禁）。
- `apps/api/internal/config/config.default.yaml`、`apps/api/configs/config.yaml` 同步更新示例（默认仍 sqlite；密码占位 `${DB_PASSWORD:-}`）。

### 证据

| 主张 | 路径 / 命令 |
|------|-------------|
| 连接面代码 | `apps/api/internal/config/config.go`（展开参数 + DSN 拼装 + fail-closed） |
| 单测 | `internal/config/config_test.go` `TestDBPostgresExplodedParams`（pass） |
| config 驱动实测 | `internal/composition/config_driven_postgres_test.go`（live PG 192.168.31.213，PASS） |
| 示例配置 | `internal/config/config.default.yaml`、`configs/config.yaml` |
