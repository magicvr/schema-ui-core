---
id: E-001-user-decisions-recorded
doc: execution-entry
status: recorded
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-001 · R1 用户方案裁决已落盘

## 已发生事实

- 用户选择 PostgreSQL 字面 `timestamptz`，SQLite 使用固定 6 位 UTC RFC3339 `TEXT`。
- 用户选择全部绝对时刻列进入分母，排除 ID、duration、TOTP step、version、计数器、金额与 flag。
- 用户选择 SQLite 与 PostgreSQL 各自原地转换；不提供 SQLite→PostgreSQL 产品级搬运器。
- 用户选择将语义为“未发生/无期限”的 sentinel `0` 转为 `NULL`，非空绝对时刻禁止零值。

## 证据

- Root D-002：`../01-decision/D-002-r1-contract-freeze-user-decisions.md`
- 子目标 D-001：`../01-decision/D-001-r1-contract-freeze.md`

## 未发生的事项

尚未完成逐列 inventory、转换 DDL、回滚策略、代码编解码或审计；本条不放行 R2。
