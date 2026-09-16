---
id: A-002-r3-dirty-state-independent
doc: audit-opinion
status: recorded
source: independent
verdict: conditional
scope: R3 C1-C3 implementation, C4 behavior and regression evidence
goal_id: GOAL-004-r3-unsaved-change-protection
auditor: grok-build (grok-4.6; reasoning high)
created: 2026-09-17
updated: 2026-09-17
parent: GOAL-004-r3-unsaved-change-protection
version: 0.1.0
---

# A-002 · R3 dirty-state 独立审计

## 审计结论

独立审计 verdict 为 `conditional`。C1～C3 的实现与 R1 D-004 / R3 D-001 合同大体一致，未发现把 R4/R5 越界写成完成；但 C2 内部导航确认证据存在一项 required 误报，须先以 `fixed` 路径处理才能关闭 R3 C4。

本意见由本地 `grok-build`（grok 4.6，reasoning high）于 2026-09-17 只读核对产生；未修改文件、status、progress 或目标台账。

## 证据判断

- `dirty-state.ts` 只保存 predicate，不保存表单草稿；throw fail-closed、注册/注销和结构比较有单元测试。
- `FormInner` 仅为非 search form 注册，baseline、reset、成功提交清理、失败保留 dirty 的实现与 R3 UI 测试基本一致；search query 不进入 registry。
- App 菜单/面包屑/Schema navigate 共用 `onNavigate`，popstate 取消会恢复 committed URL，beforeunload 使用原生事件合同；modal close 复用确认入口。
- 本地执行台账记录的 8 个测试文件、136 项回归和 TypeScript 检查通过，本独立审计未复跑测试，因此该数字作为执行层事实引用，不作为本意见的复跑结果。

## Findings

### F-001 · 内部导航确认测试误报（required / medium / open）

`App.integration.test.tsx` 的测试在第一次取消后把 `dirty` 设为 `false`，再设置 `confirm` 返回 `true` 并点击 Catalog。因此第二次导航由 `confirmDiscard` 在 clean 状态直接返回 `true`，并未调用 `window.confirm`；它证明的是 clean 导航，不是 dirty + 用户确认后再 `pushState`。

popstate 的 dirty 确认路径是真实覆盖的，实现本身没有据此判定损坏。建议保持 `dirty === true`，让 `confirm.mockReturnValue(true)` 后再次点击，并断言 `confirm` 被第二次调用且 pathname 变为 `/catalog`。

F-001 未按 `fixed`、`accepted-residual` 或 `user-overruled` 合法闭合前，不得关闭 C4、GOAL-004 或投影 Root R3。

### F-002 · 缺少真实表单 + App 导航组合证据（recommended / low）

App dirty 用例目前手工注册 global predicate，Renderer 用例是 hostless RenderPage；代码虽共用 `onNavigate`，但缺少真实 default form 改值后点击 App 菜单/面包屑的组合证据。建议至少补一条真实 Schema form + App 菜单确认回归。

### F-003 · 客户端校验、显式 throw、modal 成功提交证据不完整（recommended / low）

目前 UI 只测 HTTP 400；D-001 还要求客户端校验失败、transport throw 保留 dirty，且 modal 成功提交后卸载并清理 registry。建议补这些边界，避免把 400 扩写成全部失败子类。

### F-004 · search 非 dirty 矩阵证据写得过宽（recommended / low）

R3 UI 只直接改变了 search `q`，没有直接改变筛选/排序/分页/Saved View；应将矩阵改写为“search form 不注册、表级 query/Saved View 无 dirty source（结构性证据）”，或补对应断言。

### F-005 · 同 href 导航的丢弃语义未定义（recommended / low）

同一路径的侧栏/Command Palette 导航会 pushState + setPath，但 React 可能不卸载 FormInner，确认后草稿仍在。建议后续明确同 path 跳过确认或强制 remount/reset；换页主路径的卸载丢弃不受此 finding 影响。

## P-004 核对

本意见与 A-001 的 `pass` 不构成同范围 pass/fail 冲突；A-001 未声称 F-001 无需修正。无新的 required 信息项，R3-I-004 继续 deferred non-blocking。若不修 F-001 而要强行关门才需要用户裁决 residual/overrule；当前建议直接修测试并按 fixed 闭合。

## 是否可进入 C4

可以进入 C4 收口流程，但 F-001 合法闭合前不得把 C4、GOAL-004 或 Root R3 标为 done，也不得进行阶段放行 checkpoint。
