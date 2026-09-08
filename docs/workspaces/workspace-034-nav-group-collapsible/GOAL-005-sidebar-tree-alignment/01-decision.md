---
id: GOAL-005-sidebar-tree-alignment
doc: decision
status: active
parent: GOAL-001-nav-group-collapsible
created: 2026-09-08
updated: 2026-09-08
version: 0.1.0
---

# 决策记录 · GOAL-005-sidebar-tree-alignment

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|---|---|---|---|---|---|---|---|---|
| I-005-001 | required | 组内常态竖线与左翼对齐 | P1 | P1 | 删除 group content `border-l`、`ml-2`、`pl-3` | verified | — | D-001 |
| I-005-002 | required | active secondary/pulse 互斥 | P1 | P1 | 有 secondary 渲染文本，无 secondary 渲染 `animate-ping` pulse | verified | — | D-001 |
| I-005-003 | non-blocking | desktop/mobile 作用面 | P2 | P2 | non-horizontal NavigationLink 共用；horizontal 不改 | verified | — | D-001 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|---|---|---|---|---|
| D-001 | 2026-09-08 | 组内导航左翼与 active secondary/pulse 互斥 | accepted | `01-decision/D-001-tree-alignment-and-pulse.md` |

## 当前方案边界

- 展开组的内容容器不再绘制常态竖向 `border-l`，同时移除额外 `ml-2` / `pl-3`，让叶子 link 的左翼回到组 header 的基础列。
- `NavigationLink` 的非 horizontal active 页面左侧保留静态高亮竖线。
- secondary 与 pulse 互斥：`item.secondary` 存在时只显示文字；缺失且 `item.active` 为 true 时显示小型 `animate-ping` 光点。
- inactive 无 secondary 页面不显示 pulse；带 secondary 的 inactive 页面仍显示 secondary，延续 GOAL-003 语义。

## 未选方案

1. **保留组内竖线并只改变颜色**：不采用，用户明确要求取消常态竖线。
2. **active 页面同时显示 secondary 和 pulse**：不采用，会让右侧两种状态语义重复。
3. **无 secondary 的页面完全不显示右侧 active 提示**：不采用，本轮明确要求恢复闪烁光点作为缺省 active 效果。

