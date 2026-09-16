---
id: GOAL-004-r3-unsaved-change-protection
doc: execution
status: active
parent: GOAL-001-admin-workflow-continuity
created: 2026-09-17
updated: 2026-09-17
version: 0.2.0
---

# 执行台账 · GOAL-004 R3

## 执行索引

| E-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|--------|------|
| E-001 | 2026-09-17 | 开设 R3 并承接 dirty-state 基础切片 | recorded | [E-001-open-r3-dirty-state.md](02-execution/E-001-open-r3-dirty-state.md) |
| E-002 | 2026-09-17 | dirty-state 生命周期与离开保护回归 | recorded | [E-002-r3-dirty-state-regression.md](02-execution/E-002-r3-dirty-state-regression.md) |
| E-003 | 2026-09-17 | 响应 A-002 required/recommended 并补回归 | recorded | [E-003-r3-independent-response.md](02-execution/E-003-r3-independent-response.md) |

## 当前事实

- Root R1/R2 已分别以 `done · 3/3`、`done · 4/4` 完成；本目标按 VP-037 路线进入 R3。
- Git checkpoint `39c744ef` 已包含 `dirty-state.ts`、App 导航/离开守卫、FormInner 注册与定向单元测试；这些是 R3 的承接事实，不替代本目标 C1～C4 验收。
- 2026-09-17，补充 default form baseline/reset、提交成功/失败、search 非 dirty、modal close/cancel、内部导航、popstate 和 beforeunload 回归；R3 受影响测试 8 个文件、136 项通过，`npx tsc -p tsconfig.app.json --noEmit` 通过。
- E-002 关闭 C1～C3 的实现证据并将 R3-I-001～003 标为 `verified`；C4 的 independent audit 与 checkpoint 待完成。
- E-003 响应 A-002：F-001 按 `fixed` 修正内部导航确认测试；F-002～F-005 分别补真实 Schema form + App 菜单、客户端校验/transport throw/modal 成功、矩阵证据范围和同 href no-op。修正后 8 个受影响测试文件、140 项通过，TypeScript 检查通过；等待 independent recheck 后再关闭 C4。

## 事实边界

只写已发生的实现与验证事实；未完成的浏览器/导航/提交组合回归不得写成 R3 完成。
