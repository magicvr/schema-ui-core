---
id: E-003-postgres-scan-fix
title: 关门后跟进——修复 postgres 方言下 users 列表 500（EXISTS bool 扫描）
status: recorded
created: 2026-08-21
updated: 2026-08-21
parent: GOAL-006-dual-path-evidence
version: 0.1.0
---

# E-003 · 关门后跟进：postgres 方言 users 列表 500 修复

## 背景

用户报告仪表盘 statCard / 用户列表页 500："could not list users"。排查链：

1. 复现：副本库 + dev-session 起真实服务，HTTP /api/users 稳定 500（sqlite 直查同文件却正常）。
2. 可观测性补强：resource 工厂 INTERNAL 分支原本完全吞掉底层错误——新增 slog 错误与堆栈日志，拿到根因：`scan user list row: sql: Scan error on column index 13, name "mfa_enabled": converting driver.Value type bool ("false") to a int`。
3. 根因：用户 dev 环境已切 **postgres 方言**（configs/.env DB_DIALECT=postgres）。pgx 对 `EXISTS(...)` 返回原生 bool，而 authsession 多处把 EXISTS 扫进 `int`——sqlite 返回 int64 故测试全绿、postgres 必炸。属 VP-013 双方言缺陷（非本工作区对象存储范围引入），由本工作区使用中发现。

## 修复

`*int` 目标改 `*bool`（database/sql driver.Bool 同时接受 pgx bool 与 sqlite int64 0/1；accounts.go FeaturesForUser 与 migrate.go 为既有正确先例）：

- authsession/users_repository.go ×5：ListUsers mfa_enabled、DeleteUser exists/isAdmin、batch exists/isAdmin。
- authsession/roles_repository.go ×4：create role exists、ValidatePermissionKeys、ValidateMenuItemIDs、replaceRoleMenuItems。
- 附带：resources.go INTERNAL 分支补错误日志（含临时堆栈定位用 debug.Stack，保留错误日志本体）。

## 验证

1. **postgres 实测**：r2-pg-probe 容器建一次性库 schema_ui_r6dbg → driver=postgres 启动真实服务 → GET /api/users = **200**（修复前同路径 500）；验证后库已删除。
2. sqlite 回归：全量 go test ./... exit 0。
3. throwaway cmd/dbgdump 诊断程序已删除。

## 归属备注

authsession 双方言扫描属 workspace-013（VP-013 Store 双方言）领域；本条为交叉发现与修复记录，建议 013 侧知悉（同类模式已全量清点：internal 非测试代码 EXISTS 共 14 处，本次修复 9 处 int 目标，其余为 bool 目标或 SQL 层正确）。
