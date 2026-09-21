---
id: A-001-self-w35-db-init
doc: audit-entry
status: active
parent: GOAL-047-w35-dev-db-init-and-bootstrap
created: 2026-09-21
updated: 2026-09-21
version: 0.1.0
---

# A-001 · self 审计：W35 dev/test 数据库初始化快捷方式与自举修复

- **source**: `self`
- **日期**: 2026-09-21
- **scope**: GOAL-047 检查点 A/B；`e2e-pgset` 自举、`cmd/dbsetup`、`dev.cmd init-db`、3D000 actionable hint、文档、reset→init→start 实测；不含 checkpoint C 独立关门。
- **verdict**: `conditional`（开放 required = 0；recommended 留给 independent 复核）

## 成果

| # | 判据 | 结论 | 证据 |
|--:|------|------|------|
| 1 | PG 目标库缺失时可被明确初始化，且应用启动不静默建库 | **达成** | `openStore` 仍只 Ping + apply catalog；`cmd/dbsetup`/`dev.cmd init-db` 显式建库；3D000 hint 明确指向命令 |
| 2 | `e2e-pgset` 在 `DB_NAME` 不存在时可自举 | **达成** | `internal/pgsetup.Server.ConnectMaintenance` 候选顺序：目标库 → `postgres` → `template1`；reset 场景 `dev.cmd init-db` 无 DB_NAME 覆盖成功 |
| 3 | dev/test 一次性初始化快捷方式幂等 | **达成** | `cmd/dbsetup` 使用 dev `DB_*` 与 test `PG_TEST_*` 两个 namespace；第一次 `created schema_ui_dev/test`，第二次 `exists` |
| 4 | startup error actionable 且不泄漏秘密 | **达成** | `missingDatabaseHint` 仅处理 pgconn SQLSTATE 3D000；保留 error chain；提示 `dev.cmd init-db` / `go run ./cmd/dbsetup`；不输出 DSN |
| 5 | 文档与真实命令一致 | **达成** | 根 `README.md`、`QUICKSTART.md`、`apps/api/README.md` 均新增 PG 初始化章节；help 文案同步 |
| 6 | reset→init→start 实测 | **达成** | 见 E-002：reset 两库 → `dev.cmd init-db` → `dev.cmd start` API ready + Web 200 → stop；PG test path pass |
| 7 | 不越界 | **达成** | 未改迁移/codec/wire/checksum；只增 pgsetup/dbsetup、hint、dev launcher、docs/tests |

## 推荐事项（提交 independent 复核）

| ID | 级别 | 事项 | 复核点 |
|----|------|------|--------|
| `F-S-001` | recommended | `e2e-pgset create` 从严格创建变为幂等 `Ensure`，是否影响既有 e2e 生命周期与调用者对 exit 语义的假设 | 既有 e2e 使用随机 fresh 名且只需要成功；drop 仍显式独立；确认 docs/脚本不依赖「已存在必须失败」 |
| `F-S-002` | recommended | `pgsetup` 通过 `ServerForRole` 分离 DB_* 与 PG_TEST_*，但 `LoadEnvFile` 同时加载两个 namespace | 是否与 API/internal/pgtest 的 env 优先级一致；不同 server 的 dev/test 配置是否能分别工作 |
| `F-S-003` | recommended | 3D000 hint 是否被正确限定为 missing database，且错误链/秘密仍保护 | 认证错误/网络错误不得被改写；hint 不得回显密码 |
| `F-S-004` | recommended | `dev.cmd init-db` 的 Windows cmd 参数转发是否覆盖 `--dev-only`/`--test-only`/额外库名，并不会误删库 | 真实 cmd 调用 + dbsetup 参数解析 |

## 已核对的边界

- SQLite 路径无需此初始化：文件由首次启动创建，文档明确区分。
- e2e 专用 `schema_ui_e2e_*` 不并入 dev/test 默认集，仍由 `e2e-pgset` 独立管理。
- `.env` 为 gitignored；本轮用户授权的 `DB_NAME=schema_ui_dev` 环境变更不入库，`PG_TEST_*` 未被修改。
- 当前 root/workspace-040 已关门；本目标挂 workspace-010，是独立的 W35 运维/开发体验补强。

## 证据

- `go build ./...` exit 0。
- `go test -count=1 ./...`：**65/65 包 ok，0 fail**。
- 定向：`go test ./internal/pgsetup/`（含真实 PG Ensure/Drop）与 `go test ./internal/composition/` 均 green。
- reset 场景实测：3D000 → actionable hint；init-db 创建两库/幂等；dev.cmd start/stop 正常。

## 自审结论

A/B 判据均有可执行证据，开放 required = 0；`F-S-001`～`F-S-004` 交 grok independent 复核；checkpoint C 仅在 independent 审计通过后关门。
