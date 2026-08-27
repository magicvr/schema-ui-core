---
id: E-001
doc: execution-entry
goal: GOAL-003-dual-dialect-email-schema
status: recorded
parent: GOAL-001-account-email-identity
created: 2026-08-24
updated: 2026-08-24
version: 1.0.0
---

# E-001 · 子目标开题 + R2 schema 设计冻结（2026-08-24）

## 已发生事实

- R1 关门后按纲领路线图启动 R2；五件套一次建齐。
- 勘察事实：全仓 catalog 头 = v53；新增迁移须同步五处黄金断言；users 全部 INSERT 显式列名、无非测试 `SELECT *`。
- schema 设计冻结（D-001）：0054 三条可移植 DDL（email 列 / email_status CHECK 列 / lower(email) 唯一表达式索引），ApplyPostgres nil。
- 未改代码。

## 证据

| 主张 | 路径 |
|------|------|
| 设计条款 | 本目标 D-001 |
| catalog 头 v53 | `apps/api/internal/store/identity.go` completeFingerprintCatalogHead |
