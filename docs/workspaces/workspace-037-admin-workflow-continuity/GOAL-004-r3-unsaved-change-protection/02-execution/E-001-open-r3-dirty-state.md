---
id: E-001-open-r3-dirty-state
doc: execution
goal_id: GOAL-004-r3-unsaved-change-protection
status: recorded
created: 2026-09-17
updated: 2026-09-17
parent: GOAL-004-r3-unsaved-change-protection
version: 0.1.0
---

# E-001 · 开设 R3 并承接 dirty-state 基础切片

2026-09-17，在 R2 `done · 4/4` 投影后开设 `GOAL-004-r3-unsaved-change-protection`。R3 不新增用户待裁决的方案选择，直接执行 R1 D-004 已冻结的 dirty-state 合同。

承接范围包括：`dirty-state.ts` registry 与结构比较、App 的 `onNavigate`/`popstate`/`beforeunload` 入口、`FormInner` default-mode baseline 注册，以及已有 3 项 registry 单元测试。后续 C1～C4 仍需补齐页面生命周期、导航事件、提交/reset/cancel 与独立审计证据。
