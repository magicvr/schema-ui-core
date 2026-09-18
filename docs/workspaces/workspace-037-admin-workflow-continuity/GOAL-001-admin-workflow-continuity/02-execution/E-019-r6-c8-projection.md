---
id: E-019-r6-c8-projection
doc: execution-entry
status: recorded
goal_id: GOAL-001-admin-workflow-continuity
created: 2026-09-18
updated: 2026-09-18
parent: null
version: 1.0.0
---

# E-019 · R6 C8 控制高度与折叠开关语义完成并保持 Root 开放

2026-09-18，R6 `GOAL-007-list-page-visual-alignment` 按 D-004 完成 C8（页面 actions 高度统一、新增 `--control` 折叠开关语义 token、折叠首行无隐藏项时不渲染展开/收起按键）与定向/全量回归，R6 由 `active · 6/7` 更新为 `active · 7/8`（检查点因追加 C8 由 7 项扩为 8 项）。

C8 新增了一个全局语义 token `--control`/`--control-foreground`，是对 GOAL-007 D-002 §5 / D-001“不新增全局语义 token”边界的**局部修订**；修订已在 D-004 显式记录，并附结构守卫（`theme.test.ts`）与深浅色双态声明。既有 token 未被重命名或重定义。

Root 的六阶段纲领分母不变，R6 检查点仍为未完成，因此 Root 保持 `active · 4/6`。R5 `GOAL-006` 仍为 `active · 3/4`，R5-I-004 用户书面关门确认仍开放。

本条只记录 R6 C8 的实现投影，不关闭 R6 C6、R5、Root 或 VP-037；C6 修订审计完成后再决定 R6 完成投影。
