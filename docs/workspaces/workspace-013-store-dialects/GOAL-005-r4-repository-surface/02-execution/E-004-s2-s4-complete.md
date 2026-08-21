---
id: E-004
doc: execution-entry
goal: GOAL-005-r4-repository-surface
status: recorded
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
---

# E-004 · S2 收尾（LIKE/COLLATE 改写）+ S4（postgres 完整启动）

## 2026-08-20 · R4 主体收口

### 已发生事实

- **S2 收尾（MOUNT A-002 F-001）**（commit `e8d2b67`）：
  - `usersSortSQL`/`rolesSortSQL`：`ORDER BY … COLLATE NOCASE` → **`LOWER(col)`**（sqlite/PG 双方言大小写不敏感等价，PG 无 NOCASE collation）。
  - wallet/recyclebin 检索：`col LIKE ?` → **`LOWER(col) LIKE LOWER(?)`**（sqlite LIKE 对 ASCII 不敏感、PG 敏感；LOWER 恢复两方言一致）。
  - 定向测试（authsession/wallet/recyclebin）全绿。
- **S4（postgres 完整启动）**（commit `e8d2b67`）：
  - `composition.openStore` 返回 **`kernel.Store`**（移除 `*store.Store` type-assert）；`newJobRuntime`/`newAuthSessionRepository`/`newOperationLogRepository`/`newSettingsRepository`/`newMux`/`registerLifecycle`/`withLifecycleHooks` 消费 `kernel.Store`。
  - `handler/health.go`（Register/RegisterWithReadiness/RegisterWithMFA/readyz）、`handler/systemmonitoring.go`、`systemmonitoring/provider.go` 改用 `kernel.Store`（其仅用 Ping/SystemDataReady）。
  - **`TestCompositionPostgresStartup`（live PG）**：postgres DSN 的 `NewApp.Start` 全绿——48 迁移 apply + admin 种子 + 系统数据 reconcile + 模块 Start/Ready 门禁（store Ping + SystemDataReady = readyz 等价）通过。
- 全量 `go test ./...` **0 FAIL**（含 live PG 三测试：full boot / migrated-repo / app startup）。

### 证据

| 主张 | 路径 / commit |
|------|---------------|
| LOWER 改写 | `authsession/{users,roles}_repository.go`、`wallet|recyclebin/store/repository.go`（`e8d2b67`） |
| 公共面 kernel.Store | `composition.go`、`handler/{health,systemmonitoring}.go`、`systemmonitoring/provider.go` |
| postgres 完整启动 | `composition/postgres_startup_test.go`（live PG，`e8d2b67`） |
| 全量回归 | `apps/api`: `go test ./...` 0 FAIL |
