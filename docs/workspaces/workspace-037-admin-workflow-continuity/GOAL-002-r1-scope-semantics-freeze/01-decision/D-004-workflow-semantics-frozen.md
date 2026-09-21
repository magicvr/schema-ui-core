---
id: D-004-workflow-semantics-frozen
doc: decision
goal_id: GOAL-002-r1-scope-semantics-freeze
status: accepted
created: 2026-09-17
updated: 2026-09-17
parent: GOAL-002-r1-scope-semantics-freeze
version: 0.1.0
---

# D-004 · 工作流 dirty-state 语义冻结

本条把 VP-037 退出判据、当前 Renderer/App 合同和 R1 基线矩阵转成实现合同；它不扩大用户可见范围，也不新增业务提交语义。

- Saved View 保存 `q`、`filters`、`sort`、`order`、`pageSize`、列可见性和 UI active-view 指针；不保存当前页、选择集、record view 或 modal 草稿。
- 搜索、筛选、排序和分页只改变查询/视图状态，不构成业务 dirty；search form 不阻断离开。
- default-mode form 以挂载后的初始化快照作为 baseline；值回到 baseline 时清 dirty，值发生结构性变化时置 dirty。
- App 内部导航和面包屑/菜单导航在 dirty 时先确认；取消保持当前页面和历史状态，确认后丢弃草稿并继续 `pushState`。
- 浏览器刷新/关闭使用原生 `beforeunload` 合同；浏览器负责最终文案。`popstate` 在取消时恢复已提交 URL，在确认时提交目标历史状态。
- 成功提交清 dirty；客户端校验、字段错误、网络失败或服务端失败保留值与 dirty；写入路径不自动重复提交。
- dirty modal 的关闭/取消沿用同一确认语义；取消保持 modal，确认才关闭。

## 证据与边界

现有基线见 `attachments/r1-state-feedback-matrix.md`；dirty registry、表单接入与 App 守卫属于后续 R3 实现证据。本条冻结的是状态机语义，不宣称 R3 已完成。
