---
id: E-002-r1-user-decisions-and-child-goal
doc: execution-entry
status: recorded
parent: null
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-002 · R1 用户裁决与子目标建立

## 已发生事实

- 用户选择 PostgreSQL 字面 `timestamptz(6)`，SQLite 使用固定 6 位 UTC RFC3339 `TEXT`。
- 用户选择全部表达绝对时刻的列进入分母；ID、duration、TOTP step、version、计数器、金额与 flag 排除。
- 用户选择 SQLite 与 PostgreSQL 各自原地转换；不提供 SQLite→PG 产品级搬运器。
- 用户选择将语义为“未发生/无期限”的 sentinel `0` 转为 `NULL`，非空绝对时刻禁止零值。
- 已建立 R1 子目标 `GOAL-002-r1-contract-and-denominator-freeze`，初始 `active · 0/4`，承接 C1～C4：inventory、合同冻结、转换/备份边界、self+grok independent 审计。

## 证据

- `01-decision/D-002-r1-contract-freeze-user-decisions.md`
- `../GOAL-002-r1-contract-and-denominator-freeze/00-meta.md`
- `../GOAL-002-r1-contract-and-denominator-freeze/01-decision/D-001-r1-contract-freeze.md`

## 未发生的事项

尚未修改 migration DDL、运行时编解码或业务代码；R2 未放行。
