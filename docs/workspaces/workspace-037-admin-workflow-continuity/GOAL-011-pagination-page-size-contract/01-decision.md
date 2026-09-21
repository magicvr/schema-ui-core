---
id: GOAL-011-pagination-page-size-contract-decisions
doc: decision
status: active
parent: GOAL-001-admin-workflow-continuity
created: 2026-09-18
updated: 2026-09-18
version: 1.0.0
---

# 决策台账 · GOAL-011

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-09-18 | 分页默认值契约与跳转文案修正方案 | accepted | [D-001-pagination-contract-fix.md](01-decision/D-001-pagination-contract-fix.md) |

## 当前事实

- 2026-09-18，用户报告两个缺陷并授权「修改这两个问题后，授权走根目标关闭流程」。
- 根因已确认：前端 `DEFAULT_PAGE_SIZE = 10` 与服务端 `DefaultPageSize = 20` 不一致，而 `buildResourceQuery()` 在该值上省略 `pageSize` 参数，导致「显示 10 / 实际 20」与「10 不生效」；跳转按钮复用了 `feedback.search` 文案。
