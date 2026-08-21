---
id: E-001
doc: execution-entry
goal: GOAL-004-r3-dual-dialect-ledger
status: recorded
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
---

# E-001 · R3 立项、方案与迁移方言债快照

## 2026-08-20 · R3 立项

### 已发生事实

- Root R2（GOAL-003）已 done（independent A-002 pass + A-003 fixed F-001～F-005；live PG 探测通过）。创建 GOAL-004 五件套；R3 方案写入 D-001（Apply→`kernel.Tx`、按方言分列 catalog、checksum 绑 sqlite 历史、T1–T4 阶段、R2「非空 catalog fail-closed」在 T2 解闸、运行时 SQL 债归 R4）。
- 迁移方言债快照（扫描 `apps/api/internal/modules`，ripgrep）：
  - `authsession/migration/migration.go`：`COLLATE NOCASE` DDL（service_credentials.name）；`sqlite_master` 探测与索引探测；`PRAGMA foreign_key_list(refresh_tokens)`；`PRAGMA table_info(<table>)`。
  - 模块运行时（**R4 范围**，仅登记）：`operationlog/retention.go` `INSERT OR IGNORE`×3；`wallet/store/repository.go` / `recyclebin/store/repository.go` `LIKE`；`users_repository.go` / `roles_repository.go` `ORDER BY … COLLATE NOCASE`。
  - 时间/宽度/布尔逐列完整证据待 T3 每迁移核对（I-002）。

### 证据

| 主张 | 路径 / 命令 |
|------|-------------|
| GOAL-004 建立 + 方案 | `docs/workspaces/workspace-013-store-dialects/GOAL-004-r3-dual-dialect-ledger/`（D-001） |
| 方言债快照 | 扫描输出（本会话 ripgrep）：authsession migration；operationlog/wallet/recyclebin/users/roles 运行时 |
| R2 已 done | `GOAL-003/00-meta.md` status=done；commit `1305754`/`a090227`/`079653a` |
