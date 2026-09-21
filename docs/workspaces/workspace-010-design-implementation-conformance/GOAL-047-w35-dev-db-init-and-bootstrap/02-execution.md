---
id: GOAL-047-w35-dev-db-init-and-bootstrap
doc: execution
status: active
parent: GOAL-001-design-implementation-conformance
created: 2026-09-21
updated: 2026-09-21
version: 0.1.0
---

# 执行记录 · GOAL-047-w35-dev-db-init-and-bootstrap（W35）

## 执行索引

| E-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|--------|------|
| E-001 | 2026-09-21 | W35 立项（用户报告 + 结构/范围/slug 裁决） | recorded | `02-execution/E-001-goal-created.md` |
| E-002 | 2026-09-21 | 初始化快捷方式、自举修复、启动提示与文档 | recorded | `02-execution/E-002-init-shortcut-and-docs.md` |

## 当前事实

> E-002 的实现、文档与实测已完成：`dev.cmd init-db` 可在目标库缺失时回退 maintenance DB 并一次创建 dev/test；再次执行幂等；`dev.cmd start`/`stop` 已实测；`e2e-pgset` 自举与 PG namespace 分离由 `internal/pgsetup` 测试锁定；全仓 Go 回归 **65/65 包 ok**。检查点 A/B 完成，`progress: 2/3`；尚未完成：self + grok independent 审计与 GOAL-047 关门。

## 事实边界

> 本目标承接用户 2026-09-21 报告：PG 实例重置后 `.\dev.cmd start` 因目标库不存在而失败；需给出明确的数据库初始化快捷方式。**尚未实施**：`e2e-pgset` 自举修复、初始化快捷方式、启动错误可操作化、文档章节、回归实测与审计均为待办；`progress: 0/3`。
>
> 立项时的临时恢复（不属交付，仅记录现场处置）：以 `DB_NAME=postgres` 覆盖后 `go run ./cmd/e2e-pgset create schema_ui_dev` / `create schema_ui_test` 建库，`dev.cmd start` 随即恢复通过（API `/readyz 200` + Web 200）。该覆盖正是 #1 要修掉的缺陷。
