---
id: E-006-r3-shell-interaction
doc: execution-entry
parent_goal: GOAL-001-nav-group-collapsible
status: recorded
date: 2026-09-07
created: 2026-09-07
updated: 2026-09-07
version: 0.1.0
---

# E-006 · R3 Shell 分组交互实施

## 已发生事实

1. 用户确认 D-005：折叠状态使用 `sessionStorage`；直接 URL 自动展开覆盖已登记内页/动态路径；现行 NavGroup 不增加 key/id。
2. Web 导航 projection 增加稳定 group key 与 group active 状态；统一页面父级映射驱动 sidebar 父项在 `users-invites`、`wallet-entries` 等深链上 active。
3. Admin Shell 分组标题改为原生 button，支持 `aria-expanded`、`aria-controls`、Enter/Space 切换；折叠状态按 `schema-ui:nav-groups:v1` 保存，存储损坏/不可用时回退内存默认展开。
4. desktop sidebar 与 mobile drawer 均消费同一可折叠组组件；top/user slot 仍不进入该组件的业务归一化路径。
5. 新增 `nav-groups.test.tsx`：覆盖 Enter/Space、sessionStorage 写入、损坏存储容错、内页深链自动展开；`navigation.test.ts` 新增父链接在内页/动态路径上的 active 断言。

## 验证事实

- R3 专项：`nav-groups.test.tsx` 3/3、`navigation.test.ts` 7/7 通过。
- Web 全量：98 个测试文件、1337 个测试通过；直接 `tsc -b` 通过；直接 `vite build` 通过（仅既有 chunk size warning）。
- R3 未引入协议 NavGroup key/id；现有 App 集成与移动抽屉焦点测试继续通过。

## 当前门禁

- `I-034-003`：用户决策与实现测试均已 verified。
- `I-034-004`：范围已冻结为 sidebar 直接 route + 已登记内页/动态路径；当前有 projection/UI 证据，完整 default/optional/custom/demo route matrix/e2e 留到 R4。
- R3 代码已实施并通过 Web 验证；独立审计、R3 检查点同步与 checkpoint 尚待完成，不提前宣称 Root R3 已关闭。
