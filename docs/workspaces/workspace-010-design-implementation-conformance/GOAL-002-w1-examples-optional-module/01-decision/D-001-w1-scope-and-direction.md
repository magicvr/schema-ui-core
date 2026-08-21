---
id: D-001-w1-scope-and-direction
doc: decision-entry
goal: GOAL-002-w1-examples-optional-module
status: accepted
created: 2026-08-11
updated: 2026-08-11
version: 1.0.0
---

# D-001 · W1 范围与整改方向

## 已接受（范围）

1. **问题定性**：conformance gap（产品面卫生 + 模块装配对称性），不是安全 CVE 主叙事。  
2. **目标态**：范例/演示页与 Examples 导航成为**可选模块**；生产向 Profile 可且默认注销；Admin 标准模块不依赖演示包。  
3. **Web**：保持 schema 驱动 Shell，不为注销改中央业务路由表。  
4. **不在本波**：业务域模块；完整 WCAG/性能；物理从二进制 strip 未启用代码（仍服从 architecture「Profile 不承诺物理裁剪」）。

## 建议默认（待用户钉死 I-001～I-003）

| 项 | 建议 | 状态 |
|----|------|------|
| 模块 id | `dev.examples`（compiled 候选；默认 Profile 不含） | **待确认** |
| homePageRef（演示关） | `users`（mvp/admin 均有 admin.users） | **待确认** |
| mvp/admin 默认 | **不含**演示模块 | **待确认** |
| 启用方式 | `APP_MODULES_ENABLED` 显式列表或未来 `dogfood` profile | 方案可细化 |

## 未选方案

| 方案 | 原因 |
|------|------|
| 仅文档标注「生产请忽略 Examples」 | 不改变 as-built 强制启用与依赖 |
| 删除全部范例页与测试 | 损害协议 dogfood/回归；应可选而非消灭 |
| 重开 VP-003 全量迁移 | 过重；本波有界装配修正即可 |
