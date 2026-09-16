---
id: D-001-r3-dirty-state-contract
doc: decision
status: accepted
goal_id: GOAL-004-r3-unsaved-change-protection
created: 2026-09-17
updated: 2026-09-17
parent: GOAL-004-r3-unsaved-change-protection
version: 0.1.0
---

# D-001 · R3 未保存变更保护合同

承接 Root R1 决策 D-004，R3 采用以下实现合同：

- `FormInner` 仅为 default-mode 表单注册只读 dirty predicate；表单值不写入全局 registry。初始化后的结构快照为 baseline，值回到 baseline 即清 dirty。
- search form 的查询变化不计为业务 dirty；Saved View 的选择、筛选、排序和分页也不阻断页面离开。
- App 的菜单、面包屑、Schema navigate 使用同一个确认入口。取消不执行 `pushState`；确认后才切换页面并丢弃当前草稿。
- 浏览器 `popstate` 已改变地址时，取消恢复最近 committed location；确认则解析目标并更新 committed location。`beforeunload` 只设置原生事件合同，不依赖自定义文案。
- 成功提交把当前值设为新 baseline；客户端校验、字段错误、服务端失败和 transport throw 保留值与 dirty。modal close/cancel 在 dirty 时先确认，reset 回到 baseline 后清 dirty。

本条不改变 R4 的反馈分类，也不把 R2/R4/R5 的阶段状态提前标为完成。
