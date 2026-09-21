---
id: E-003-inventory-complete-v0.2
doc: execution-entry
status: recorded
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-003 · 完整时间列 inventory v0.2

## 已发生事实

- 完成 compiled migration DDL + runtime repository 的只读逐模块盘点。
- 当前 live schema 共识别 **90 个绝对时刻列**；另有历史 retired `records.updated_at`，不计 live 分母但保留来源记录。
- 确认 compiled DDL 时间列当前均为 SQLite `INTEGER` / PostgreSQL `BIGINT`，未发现现存时间列 TEXT。
- 完成 seconds/milliseconds 分组：jobs、operationlog、mail outbox/config 为 milliseconds；其余主要模块为 seconds；nullable 与 sentinel 列已标记。
- inventory v0.3/v0.3.1 已落盘，C1 inventory 逐列 90 条、历史 retired source、compiled catalog v1–v72 与 v67–v72 无额外时间列均已登记；具备再次 self/independent 复审的证据基础。

## 证据

- `attachments/r1-time-column-inventory-v0.2.md`
- `attachments/r1-time-column-inventory-v0.1.md`（初版历史）
- VP-013 R1/R3 既有合同与 migration DDL

## 未发生事项

C2 目标 codec/DDL、C3 双方言原地转换、失败回滚与备份验证尚未完成；不放行 R2。
