---
id: GOAL-047-w35-dev-db-init-and-bootstrap
doc: execution
status: done
parent: GOAL-001-design-implementation-conformance
created: 2026-09-21
updated: 2026-09-21
version: 1.0.0
---

# 执行记录 · GOAL-047-w35-dev-db-init-and-bootstrap（W35）

## 执行索引

| E-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|--------|------|
| E-001 | 2026-09-21 | W35 立项（用户报告 + 结构/范围/slug 裁决） | recorded | `02-execution/E-001-goal-created.md` |
| E-002 | 2026-09-21 | 初始化快捷方式、自举修复、启动提示与文档 | recorded | `02-execution/E-002-init-shortcut-and-docs.md` |
| E-003 | 2026-09-21 | 响应 independent A-002：F-001/F-002/F-003 修复、A-003 复审通过与关门 | recorded | `03-audit/A-003-independent-w35-findings-rereview.md`；`03-audit/A-004-response-to-a003-and-closeout.md` |

## 当前事实

> E-002 的实现、文档与实测已完成：`dev.cmd init-db` 可在目标库缺失时回退 maintenance DB 并一次创建 dev/test；再次执行幂等；`dev.cmd start`/`stop` 已实测；`e2e-pgset` 自举与 PG namespace 分离由 `internal/pgsetup` 测试锁定；全仓 Go 回归 **65/65 包 ok**。A-002 的 3 条 required 已 fixed，A-003 independent pass / open required = 0，A-004 关门记录已落盘；GOAL-047 `done · 3/3`。

## 事实边界

> 本目标承接用户 2026-09-21 报告：PG 实例重置后 `.\dev.cmd start` 因目标库不存在而失败；已交付显式数据库初始化、自举回退、actionable hint、文档与回归证据。SQLite 无需该初始化步骤；e2e 专用库继续由 e2e helper 独立管理。
>
> 现场恢复已转化为可复用路径：删除两个库 → 无环境覆盖执行 `dev.cmd init-db` → 两库均创建 → `dev.cmd start` API/Web 全绿；第二次 init-db 幂等。