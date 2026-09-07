---
id: D-005-r3-collapse-and-deep-link-scope
doc: decision-entry
parent_goal: GOAL-001-nav-group-collapsible
source: user-confirmed + /govern
status: accepted
date: 2026-09-07
created: 2026-09-07
updated: 2026-09-07
version: 0.1.0
---

# D-005 · R3 折叠状态与直接 URL 范围

## 决策

用户确认：

1. 分组折叠状态采用 `sessionStorage` 会话内持久化：同一浏览器会话内刷新/路由跳转保持；关闭会话后清理；不引入服务端存储。
2. 直接 URL 自动展开覆盖已分组模块的内页与动态路径，不仅是 sidebar NavLink 自身 route。现有父级关系纳入验收：
   - `users-invites` → `users`
   - `wallet-entries` → `wallet`
   - `dictionary-entries` → `data-dictionary`
   - `task-runs` → `scheduled-tasks`
   - `telegram-operator` → `telegram-settings`
3. 现行协议 `NavGroup` 不新增 `key`/`id`；Web 侧使用 Manifest `labelKey`（缺失时 literal label）作为会话状态 key，并以 child active/父级页面关系判断自动展开。
4. 用户手动折叠当前 active 组后不在同一次渲染中强制反复打开；当路由重新进入该组或进入其深链时，组再次自动展开。

## 方案边界

- 桌面 sidebar 与 mobile drawer 都使用同一 sessionStorage namespace；两个挂载实例各自响应本地状态，重新打开 drawer 时从会话状态读取。
- sessionStorage 不可用、内容损坏或结构不合法时，忽略存储并回退为默认展开；不得阻断导航。
- storage 只保存 group open/closed 布尔值，不保存用户、权限、路由参数或业务数据。
- R3 只负责 Shell 交互和导航激活投影；R4 负责当前所有 sidebar 节点、optional/custom/demo Profile 的全量矩阵回归。

## 信息项

- `I-034-003`：由本决策关闭为 `verified（用户决策）`；R3 仍需以交互测试证明 sessionStorage 读写与容错。
- `I-034-004`：范围已由本决策明确为 sidebar 直接 route + 已登记内页/动态路径；状态保持 `open`，直到 R3/R4 矩阵有证据。
