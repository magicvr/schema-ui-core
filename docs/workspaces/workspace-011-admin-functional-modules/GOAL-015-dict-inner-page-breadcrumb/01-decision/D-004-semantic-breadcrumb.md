---
id: D-004
goal: GOAL-015-dict-inner-page-breadcrumb
date: 2026-08-14
status: accepted
parent: GOAL-015-dict-inner-page-breadcrumb
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# D-004 · 面包屑改为语义层级 + 视觉重构（取代 D-002 §3.6 路由栈方案）

## 决策

- 用户 2026-08-14 裁决：面包屑逻辑从「页面访问历史」（D-002 §3.6 路由栈，用户此前确认）改为**有语义的上下级页面关系**；导航格式固定为 `首页 => 一级页 => ... => n级内页`，其中**首页 = manifest homePageRef（域名根直接指向的默认页，不一定是仪表盘）**。
- 语义来源（无协议增补，web-shell 层）：
-   1. manifest 导航树：页面在导航中的组标签（非点击纯标签段）作为祖先；
-   2. 消费侧声明父级 `BREADCRUMB_PAGE_PARENTS`（web-shell 常量，页面级声明，不改协议）：dictionary-entries → data-dictionary、task-runs → scheduled-tasks；
-   3. 首页恒为根（当前页即首页时不显示 trail）；未知声明父级 fail-safe 跳过；父链防环。
- 返回操作：紧凑圆形幽灵图标按钮（仅 ← 图标，hover 浅底色），位于面包屑最左；目标 = 最近的非首页路由祖先（语义父级），无语义父级时不显示。
- 视觉规范：12px（text-xs）常规字重；非当前项 text-muted-foreground（浅灰）+ hover 变亮/下划线；分隔符细斜杠 `/` 低透明度；当前项更亮（text-foreground/90）不可点击；面包屑与页面标题间距 10px（mt-2.5）；移除「ADMIN WORKSPACE」眉标与纯文本「← 返回上一页 |」拼接。

## 未选方案

- **新开子目标承载**：修订同属 GOAL-015 交付物（面包屑），且直接取代其内 D-002 决策——父目标已 done 而决策被外部子目标推翻会造成记录割裂；追加到原目标台账保证决策链连续。
- **schema/协议声明父级**（breadcrumbParent 等）：page.schema meta additionalProperties:false，需上游协议增补（P-005 门禁）；本决策用 web-shell 声明常量（消费侧），零协议成本。
