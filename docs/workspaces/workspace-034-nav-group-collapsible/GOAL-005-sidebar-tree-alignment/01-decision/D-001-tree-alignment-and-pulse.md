---
id: D-001-tree-alignment-and-pulse
goal_id: GOAL-005-sidebar-tree-alignment
doc: decision-entry
source: user-request
status: accepted
date: 2026-09-08
parent: GOAL-001-nav-group-collapsible
version: 0.1.0
---

# D-001 · 组内导航左翼与 active secondary/pulse 互斥

## 决策

1. 展开组的内容容器取消常态左侧竖线，并移除额外左侧缩进，让组内菜单整体左移；组名与 active link 的最左侧视觉列基本平齐。
2. active 页面左侧保留静态高亮竖线。
3. active 页面有 secondary 时，右侧显示 secondary 文本；没有 secondary 时，右侧显示闪烁光点。secondary 与 pulse 互斥。
4. 仅在非 horizontal 的 sidebar/mobile NavigationLink 中启用上述 active 指示；不改变 top/user slot。

## 理由

- 竖线属于树形导轨，但当前组层级不需要长期显示该装饰；移除后可释放左侧空间并改善对齐。
- active 竖线负责稳定选中定位，secondary/pulse 负责右侧状态/注册元数据，两个位置语义清楚且不重复。
- 无 secondary 时仍保留 pulse，避免 active 页面失去右侧状态反馈；有 secondary 时以注册文本优先。

## 验收指向

- `nav-groups.test.tsx`：无 group content border、菜单左移、secondary/pulse 互斥。
- `navigation.test.ts` / `nav-groups-r4.test.ts`：active/deep-link projection 保持。
- `tsc -b` / `npm run build`：类型与生产构建通过。

