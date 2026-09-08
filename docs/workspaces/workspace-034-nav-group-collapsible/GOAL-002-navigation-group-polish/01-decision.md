---
id: GOAL-002-navigation-group-polish
doc: decision
status: active
parent: GOAL-001-nav-group-collapsible
created: 2026-09-07
updated: 2026-09-07
version: 0.1.0
---

# 决策记录 · GOAL-002-navigation-group-polish

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 证据 / 决策 |
|---|---|---|---|---|---|---|---|
| I-002-001 | non-blocking | 分组样式复用现有语义化视觉体系，不引入新主题变量 | P1 | P1 | 对照既有 link/accent/muted/border/focus 样式并补回归 | verified（用户要求 + 实现/测试） | D-001；E-001 |
| I-002-002 | required | 非 active 分组默认关闭；访问分组内直接页面或内页时自动展开 | P2 | P2 | 用户确认并通过 App/navigation 测试 | verified（用户决策 + 测试） | D-001；E-001 |
| I-002-003 | non-blocking | sessionStorage 偏好与 active 深链优先级保持 D-005 语义 | P2 | P2 | 补默认关闭/手动保持/active 优先级测试 | verified（继承 + 测试） | D-005；D-001；E-001 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|---|---|---|---|---|
| D-001 | 2026-09-07 | 导航分组语义样式与默认折叠行为 | accepted | `01-decision/D-001-navigation-group-polish-and-default-closed.md` |

## 当前方案边界

- 使用现有 `muted`、`accent`、`border`、focus ring、间距与圆角语义，增加分组 header 的层级容器、active/focus 状态、折叠指示与子项层级导轨。
- 非 active 分组默认关闭；active 直接页面或 D-005 登记的内页/动态深链自动展开。
- 手动折叠状态继续写入 `sessionStorage`；active 自动展开不覆盖用户离开 active 路由后的手动偏好。
- 不改 VP-034 既有 group key/成员/协议契约，不重开父目标。
