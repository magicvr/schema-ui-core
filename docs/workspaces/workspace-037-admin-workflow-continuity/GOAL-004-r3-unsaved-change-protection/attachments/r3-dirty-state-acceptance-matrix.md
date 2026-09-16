---
id: r3-dirty-state-acceptance-matrix
doc: attachment
status: active
created: 2026-09-17
updated: 2026-09-17
parent: GOAL-004-r3-unsaved-change-protection
version: 0.1.0
---

# R3 · dirty-state 验收矩阵

| 场景 | 预期 | 证据 | 状态 |
|------|------|------|------|
| default form 初始值 | clean，注册 baseline predicate | FormInner + Renderer test | collecting |
| default form 改值/还原 | dirty → clean | Renderer test | collecting |
| search q/筛选/排序/分页 | 不触发离开确认 | App/Renderer test | collecting |
| 内部导航取消/确认 | 取消留在当前页；确认后切换 | App integration test | collecting |
| popstate 取消/确认 | 取消恢复 committed URL；确认提交目标 | App integration test | collecting |
| beforeunload clean/dirty | clean 不拦截；dirty 设置原生 returnValue | event test | collecting |
| submit success/failure | success 清 dirty；failure 保留 dirty | Renderer test | collecting |
| modal close/cancel/reset | dirty 先确认；取消保留；确认/恢复 baseline 后清理 | Renderer test | collecting |
