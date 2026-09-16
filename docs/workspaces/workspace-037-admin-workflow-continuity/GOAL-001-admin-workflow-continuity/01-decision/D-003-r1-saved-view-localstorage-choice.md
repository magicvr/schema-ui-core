---
id: D-003-r1-saved-view-localstorage-choice
doc: decision
goal_id: GOAL-001-admin-workflow-continuity
status: accepted
created: 2026-09-17
updated: 2026-09-17
parent: null
version: 0.1.0
---

# D-003 · R1 Saved View 持久化取舍

Root 记录用户在 R1 子目标中作出的关键选择：采用浏览器 `localStorage`，按 `user.id + pageId + tableId` 隔离 Saved View。完整决策、allowlist、fail-closed 与待验证边界见 [R1 D-003](../../GOAL-002-r1-scope-semantics-freeze/01-decision/D-003-saved-view-localstorage-accepted.md)。

该选择只关闭持久化方向；序列化校验、Schema/权限失效和读写异常的实现证据仍由 R2 收集，I-037-002 保持 `collecting`。
