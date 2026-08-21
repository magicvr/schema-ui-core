---
title: R7 Records surface scan evidence
date: 2026-08-05
status: recorded
parent: GOAL-007-r4-records-retirement-closure
version: 0.1.0
---

# Records 退场扫描证据

## 当前产品面未发现

- `apps/api/internal/handler/records.go`
- `apps/api/internal/store/records.go`
- `apps/api/internal/store/seed_records.go`
- `/api/records` 产品路由、Records manifest page、Records 专属 frontend hook
- Records 专属 API/store/seed/manifest/fixture 运行实现

## 明确保留

- `apps/api/internal/store/migrate.go:80-101,199-217,297-311`：不可改写的 0003/0006
  migration ledger 和退场 DDL。
- `apps/api/internal/store/migrate.go:229,257,355`：历史 `records.*` operation-log
  CHECK 值。
- `apps/api/internal/store/operations_test.go:9,101-151`：历史事件、升级快照和
  0006 删除表测试。
- `apps/web/src/renderer/resource.test.ts:53-58`、`render.test.ts:78`：负向防复活
  测试。
- `apps/web/src/protocol/conformance/request-construction.ts` 及 upstream cases：
  通用 `recordSource`/`result.records` 协议，不是 Records 业务实体。

## 本轮变更与验证

`schema-table.test.tsx` 的演示命名已泛化为 `SAMPLE_ROWS`/`rowsFetcher`/`schema-table`；
DataTable、导航、Schema CRUD、E2E 和 API 注释中的 legacy demo wording 已同步。
Web 定向测试 4 files / 62 tests 通过；API `internal/store`, `internal/handler`,
`internal/auth` 测试通过。`git diff --check` 通过。
