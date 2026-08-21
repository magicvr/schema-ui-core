---
id: E-002-r4-records-surface-cleanup
doc: execution-entry
goal: GOAL-007-r4-records-retirement-closure
source: orchestrator
date: 2026-08-05
status: recorded
---

# E-002 · Records 运行面扫描与命名清理

扫描确认当前 `apps/api`/`apps/web` 没有 Records 产品 API、handler、store、seed、
manifest 页面、专属 hook 或运行 fixture。保留项为已应用 `0003`/`0006` 迁移、历史
`records.*` operation-log 兼容测试、通用 `recordView`/`recordSource`/`RecordID` 和
负向防复活测试。

本轮只修改通用测试/注释中的误导性命名：`RECORDS` → `SAMPLE_ROWS`、
`recordsFetcher` → `rowsFetcher`、`records-table` → `schema-table`，以及 legacy demo
说明文案。Web 定向测试 62/62 通过；API store/handler/auth 定向测试通过。
