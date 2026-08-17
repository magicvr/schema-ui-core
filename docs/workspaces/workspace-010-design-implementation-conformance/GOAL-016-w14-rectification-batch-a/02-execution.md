---
id: GOAL-016-w14-rectification-batch-a
doc: execution
status: active
parent: GOAL-015-w14-user-perspective-review
created: 2026-08-17
updated: 2026-08-17
version: 0.1.0
---

# 执行记录 · GOAL-016

## 执行索引

| E-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| E-001 | 2026-08-17 | 立项：由 GOAL-015 路线图 + W14 用户裁决（D-003）建立批 A 子目标（parent=GOAL-015，用户结构裁决后更名 w14） | recorded | `02-execution/E-001-w14-goal-established.md` |
| E-002 | 2026-08-17 | S1 方案冻结：F-01～F-04 设计落盘（D-001）；I-001/I-002 closed | recorded | `02-execution/E-002-s1-freeze.md` |
| E-003 | 2026-08-17 | S2/S3 实施与回归：F-01～F-04 代码/schema/UI 落盘；Go 全量 + Web 全量 + tsc + build 绿 | recorded | `02-execution/E-003-s2-s3-implementation.md` |

## 事实边界

> 只写已经发生且有证据的事实。计划、未知与建议留在决策。

- **2026-08-17**：用户经 GUI 对 W14 I-001 作出书面裁决（GOAL-015 D-003）：F-01～F-14 全部 in-scope、分批 A→C→D→B；F-01 新增端点 / F-04 存 messageKey / F-08 直接移除。
- **2026-08-17**：W14 用户结构裁决后，整改路线图归 GOAL-015（R1～R4），批 A 子目标 GOAL-016（parent=GOAL-015）建立（五件套落盘）；批 C/D/B 后续渐进添加。
- **2026-08-17**：S1 方案冻结（E-002）——D-001 落盘：F-01 handler 目录端点（GET /api/scheduled-tasks/handlers）、F-02 scopes UI 自定义组件、F-03 operations 结构化过滤 + CSV 导出、F-04 notification messageKey（title_key/body_key 迁移 0037）；I-001/I-002 关闭。
- **2026-08-17**：S2/S3 实施与回归（E-003）——F-01～F-04 代码/schema/UI 全部落盘；Go 全量、Web 全量、tsc、build 均通过。
