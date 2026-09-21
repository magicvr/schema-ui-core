---
id: E-002-r3-dirty-state-regression
doc: execution
goal_id: GOAL-004-r3-unsaved-change-protection
status: recorded
created: 2026-09-17
updated: 2026-09-17
parent: GOAL-004-r3-unsaved-change-protection
version: 0.1.0
---

# E-002 · dirty-state 生命周期与离开保护回归

已发生的实现与验证事实：

- `FormInner` 为 default-mode form 注册 predicate；结构性值比较、恢复 baseline 与成功提交后的 baseline 更新会同步 `data-form-dirty` 和全局 registry。
- default form 处理原生 `reset` 事件，回到 baseline 并清理字段/表单错误；search form 不注册业务 dirty。
- App 的内部菜单导航、`popstate` 取消/确认和 committed URL 恢复，以及 clean/dirty `beforeunload` 均由 `App.integration.test.tsx` 覆盖。
- modal close 在 dirty 时复用同一确认入口；取消保持 modal，确认关闭并卸载 dirty source。
- `r3-dirty-state.ui.test.tsx` 4 项、`dirty-state.test.ts` 3 项、App integration 14 项，加上既有 R2/Renderer/CRUD 回归共 8 个文件、136 项通过；`npx tsc -p tsconfig.app.json --noEmit` 通过。

本条不把 R4 的统一 FeedbackRegion retry 或 R5 组合验收写成完成事实。
