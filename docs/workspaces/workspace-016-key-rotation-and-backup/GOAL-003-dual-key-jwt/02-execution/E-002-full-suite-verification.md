---
id: E-002
doc: execution-entry
goal: GOAL-003-dual-key-jwt
status: recorded
created: 2026-08-22
updated: 2026-08-22
version: 1.1.0
---

# E-002 · R2 全仓验证

## 事实（2026-08-22 · v1.1 按 A-002 F-003 收窄措辞）

- `go vet ./...`（apps/api 全模块）：0 finding。
- **JWT 相关包**（auth / config / composition / handler）：`go test -count=1` 全部 `ok`——本切片直接相关面零回归，且 independent 审计员独立复跑一致。
- 整包 `go test ./...`：编排器两次运行 exit 0（其中第二次 `internal/store` 为缓存结果）；independent 审计员以 `-count=1` 复跑时 `internal/store` 两条 **Postgres 集成测试**（`TestOpenPostgresProbeIntegration` / `TestPostgresMigrateRunnerIntegration`）因共享 probe DB 残留表报 `WasFresh() = false`；编排器随后 `-count=1` 复跑同样复现该两条失败。R2 diff 不含 `internal/store`（`git diff --name-only` 相对 `c96e963` = auth.go / auth_test.go / composition.go + 台账），失败机制为 `postgresWasFresh` 对 probe schema 内**任意**用户表敏感、而测试只清理自己的三张表（`store/postgres_test.go:599-615`、`postgres.go:230-254`）——属既有测试基建对共享 DB 状态的脆弱性，**非本切片回归**。
- 结论措辞：本切片证据 = vet 0 + JWT 相关包 ok（双方独立一致）；「整包 exit 0」不再作为本目标关门证据引用。

## 遗留注记（供 R3 参考）

R3 轮换后恢复证据要走 PG 方言：I-004 剧本须使用一次性/专用数据库或先做全 schema 清理，避免共享 probe DB 残留污染 WasFresh 类判定与恢复断言。
