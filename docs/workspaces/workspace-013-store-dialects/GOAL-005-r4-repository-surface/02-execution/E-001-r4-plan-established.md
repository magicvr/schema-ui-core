---
id: E-001
doc: execution-entry
goal: GOAL-005-r4-repository-surface
status: recorded
created: 2026-08-20
updated: 2026-08-20
version: 1.0.0
---

# E-001 · R4 立项与方案

## 2026-08-20 · R4 建立

### 已发生事实

- Root R3（GOAL-004）done（Root 3/5）。创建 GOAL-005 五件套；R4 方案写入 D-001（收口形状、运行时 SQL 债规则、postgres 启动 S4、self+independent 关门）。
- 泄漏面已知锚点（S0 将补全清单）：`withTx`/`runner.WithTx(ctx, func(*sql.Tx))` 遍布 authsession/jobs/operationlog/settings/datapermission/datadictionary/logincaptcha/mfa/scheduledtasks/recyclebin/wallet；`operationlog/retention.go` `INSERT OR IGNORE`×3；`wallet/store`、`recyclebin/store` `LIKE`；`users_repository`/`roles_repository` `ORDER BY … COLLATE NOCASE`。

### 证据

| 主张 | 路径 / 命令 |
|------|-------------|
| GOAL-005 建立 + 方案 | `docs/workspaces/workspace-013-store-dialects/GOAL-005-r4-repository-surface/`（D-001） |
| 锚点 | 本会话 grep（前置轮次）：`WithTx(*sql.Tx)` 140 处；`INSERT OR IGNORE`/`LIKE`/`COLLATE NOCASE` 运行时命中 |
| 前提 R3 done | `GOAL-004/00-meta.md` status=done；Root `00-meta.md` progress 3/5 |
