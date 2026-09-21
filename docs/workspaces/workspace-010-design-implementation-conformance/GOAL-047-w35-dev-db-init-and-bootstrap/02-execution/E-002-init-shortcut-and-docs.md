---
id: E-002-init-shortcut-and-docs
doc: execution-entry
status: active
parent: GOAL-047-w35-dev-db-init-and-bootstrap
created: 2026-09-21
updated: 2026-09-21
version: 0.1.0
---

# E-002 · 初始化快捷方式、自举修复、启动提示与文档

## 已完成事实

- 新增 `apps/api/internal/pgsetup`：按角色解析 `DB_*` / `PG_TEST_*`，加载 `.env`（进程环境优先）；维护连接按「目标库 → `postgres` → `template1`」回退；`Ensure` 幂等创建；`Drop` 幂等清理；`VerifyMigrated` 与 e2e 数据库列表。
- 新增 `apps/api/cmd/dbsetup`：默认一次确保 dev (`DB_NAME`) + test (`PG_TEST_DB`)；支持 `--dev-only`、`--test-only`、额外库名；不删除、不迁移。
- `apps/api/cmd/e2e-pgset` 重用 `pgsetup`，空实例自举不再依赖 `DB_NAME` 已存在；`create` 幂等、`drop` 对缺失库安全返回 `absent`。
- `dev.cmd init-db`：仓库根一键调用 `go run ./cmd/dbsetup`，help 文案同步。
- PG 3D000 启动错误在 `internal/store/postgres.go` 增加 actionable hint，指向 `dev.cmd init-db` / `go run ./cmd/dbsetup`；保留原错误链与分类，不泄露 DSN 秘密。
- 文档同步：根 `README.md`、`QUICKSTART.md`、`apps/api/README.md` 新增 PostgreSQL 初始化章节，说明 SQLite/PG 差异、权限、命令、e2e 专用库边界与 3D000 排查。
- 回归测试：`internal/pgsetup/pgsetup_test.go` 覆盖名称安全、维护库回退、DB/PG_TEST namespace 分离、DSN、真实 server 上 Ensure 幂等 + Drop；`recovery_wiring_test.go` 锁绝对 artifact dir；`postgres_open_test.go` 锁 3D000 actionable hint。

## 最终验证事实

- `go build ./...` exit 0。
- `go test -count=1 ./...`：**65/65 包 ok，0 fail**（新增 `internal/pgsetup` 测试包）。
- 实例 reset 场景：删除 `schema_ui_dev`/`schema_ui_test` → 无环境覆盖执行 `dev.cmd init-db` → 两库均 `created`；再次执行均 `exists`；`dev.cmd start` API `/readyz 200` + Web 200；`dev.cmd stop` 无残留。
- `e2e-pgset list` 在 `DB_NAME` 不存在时现在能回退维护库而不是 3D000；`cmd/dbsetup` 按 DB/PG_TEST 两个 namespace 正确创建对应库。

## 实测证据

- 模拟 reset：删除 `schema_ui_dev` 与 `schema_ui_test` 后，**无 DB_NAME 覆盖**执行 `dev.cmd init-db`：
  - maintenance connection 回退到 `postgres`；
  - `created schema_ui_dev (dev)`；
  - `created schema_ui_test (test)`；
  - 第二次执行均报告 `exists`（幂等）。
- `dev.cmd start --no-browser`（默认 `.env`，端口临时覆盖）：API `/readyz 200`、Web 监听成功；`dev.cmd stop` 无残留。
- `schema_ui_dev` / `schema_ui_test` 为一次性恢复后的专用数据库；PG 测试路径（`TestPGRestoreToNewDB`）通过。
- 直接启动缺库时，错误包含：`database "schema_ui_dev" does not exist` + `provision it first: dev.cmd init-db ...`。

## 边界事实

- 本次 `.env` 改动是用户授权的 gitignored 环境变更：`DB_NAME=postgres` → `schema_ui_dev`；`PG_TEST_DB=schema_ui_test` 保持不变；秘密未输出、未入库。
- 共享 `postgres` 数据库仍存在既有 `schema_migrations`（来自测试环境配置 `PG_TEST_DB=postgres` 的历史行为）；本目标不擅自重置或迁移它，另行记录为环境观察。

## 未完成

- 功能实现与文档已完成；本目标尚未做 self / grok independent 审计，`progress` 仍按 A/B/C 维护。
