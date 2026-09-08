---
id: GOAL-004-sidebar-active-indicator
doc: decision
status: active
parent: GOAL-001-nav-group-collapsible
created: 2026-09-08
updated: 2026-09-08
version: 0.1.0
---

# 决策记录 · GOAL-004-sidebar-active-indicator

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|---|---|---|---|---|---|---|---|---|
| I-004-001 | required | 组间距参考节奏 | P1 | P1 | Sidebar Engine 范例使用 `space-y-sm`；映射为 `space-y-2` | verified | — | D-001 |
| I-004-002 | required | active 页面指示关系 | P1 | P1 | 左侧竖线保留；secondary 替代右侧光点；缺省不渲染 | verified | — | D-001 |
| I-004-003 | non-blocking | desktop/mobile 对齐 | P2 | P2 | 共用 NavigationLink；DOM/class 回归 | collecting | P2 | D-001 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|---|---|---|---|---|
| D-001 | 2026-09-08 | 紧凑组间距与 active 指示语义 | accepted | `01-decision/D-001-spacing-and-active-marker.md` |

## 当前方案边界

- 将非 horizontal `NavigationItems` 的 `space-y-6` 改为参考页相近的 `space-y-2`；组内叶子间距保持原有紧凑值。
- `NavigationLink` 在 sidebar/mobile active 页面最左侧渲染静态 `data-navigation-active-marker="active"` 高亮竖线；不使用动画。
- active 页面右侧继续显示已注册的 `secondary` 副文本；没有 secondary 时不渲染 secondary 元素，也不渲染闪烁光点或空占位。
- active/non-active 的路由、权限、折叠状态和 deep-link 自动展开逻辑保持不变。

## 未选方案

1. **保留所有组的 24px 间距**：不采用，当前视觉密度明显高于参考页。
2. **无 secondary 时用静态圆点填补右侧**：不采用，用户明确要求不设置选中闪烁光点，且无注册副文本时不应制造占位。
3. **通过 pageRef 在 Shell 中推断 secondary**：不采用，副文本继续以注册 projection 为唯一来源。

