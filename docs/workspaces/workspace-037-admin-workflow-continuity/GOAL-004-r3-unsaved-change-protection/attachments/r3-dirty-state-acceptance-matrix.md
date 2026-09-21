---
id: r3-dirty-state-acceptance-matrix
doc: attachment
status: done
created: 2026-09-17
updated: 2026-09-17
parent: GOAL-004-r3-unsaved-change-protection
version: 0.2.0
---

# R3 · dirty-state 验收矩阵

| 场景 | 预期 | 证据 | 状态 |
|------|------|------|------|
| default form 初始值 | clean，注册 baseline predicate | E-002 · R3 UI test | verified |
| default form 改值/还原 | dirty → clean | E-002 · R3 UI test | verified |
| search form q | 不注册业务 dirty | E-002 · R3 UI test + FormInner `isSearch` guard | verified |
| 表级 query / Saved View 选择 | 不创建 dirty source；allowlist 状态不属于业务草稿 | R2 E-003 + provider 结构；R3 E-003 收窄证据范围 | verified |
| 内部导航取消/确认 | 取消留在当前页；确认后切换 | E-002 · App integration | verified |
| popstate 取消/确认 | 取消恢复 committed URL；确认提交目标 | E-002 · App integration | verified |
| beforeunload clean/dirty | clean 不拦截；dirty 设置原生 returnValue | E-002 · App integration | verified |
| submit success/failure | success 清 dirty；failure 保留 dirty | E-002 · R3 UI test | verified |
| modal close/cancel/reset | dirty 先确认；取消保留；确认/恢复 baseline 后清理 | E-002 · R3 UI test | verified |
