---
id: GOAL-007-list-page-visual-alignment-decisions
doc: decision
status: done
parent: GOAL-001-admin-workflow-continuity
created: 2026-09-18
updated: 2026-09-18
version: 1.0.0
---

# 决策台账 · GOAL-007-list-page-visual-alignment

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-007-001 | required | 范例页的目标布局和视觉语义 | C2 | C2 | 读取 raw 范例与 DESIGN.md | verified | 2026-09-18；冲突时回到 C2 | `00-meta.md`、D-001 |
| I-007-002 | required | 通用列表调用链与 shell 不变边界 | C2/C3 | C2 | 盘点 renderer/component/app/test | verified | 2026-09-18；不得修改 topbar/sidenav | `00-meta.md`、D-001 |
| I-007-003 | required | 查询/重置/分页现有合同 | C2/C3 | C2 | 对照 handlers、query bridge、分页测试 | verified | 2026-09-18；保持行为，仅调整展示 | `00-meta.md`、D-001 |
| I-007-004 | required | “全部{对象}”的集中语义对象来源 | C3 | C3 前 | 集中解析并补测试 | verified | 2026-09-18 已验证集中解析器、可靠标题回退与对象标签测试 | `list-surface.test.ts`、E-002 |
| I-007-005 | non-blocking | 多选批量能力是否已存在 | C3/C4 | C3 | 扫描并复用；未实现则忽略 | deferred | 用户明确要求时另立范围 | 用户指令 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-09-18 | R6 列表页视觉与筛选实施合同 | accepted | [D-001-list-page-visual-contract.md](01-decision/D-001-list-page-visual-contract.md) |
| D-002 | 2026-09-18 | 回开 R6 并修订范例页布局合同 | accepted | [D-002-r6-layout-revision.md](01-decision/D-002-r6-layout-revision.md) |
| D-003 | 2026-09-18 | R6 追加 C7 控制位与搜索配对纠偏 | accepted | [D-003-r6-c7-control-placement-correction.md](01-decision/D-003-r6-c7-control-placement-correction.md) |
| D-004 | 2026-09-18 | R6 追加 C8：控制高度、折叠开关语义 token 与空展开抑制 | accepted | [D-004-r6-c8-control-height-and-toggle-token.md](01-decision/D-004-r6-c8-control-height-and-toggle-token.md) |

## 当前投影

- R6 沿用 VP-037、当前 workspace 和 Root；用户选择回开 GOAL-007 追加 C5/C6，不新建 VP/workspace 或 GOAL-008。
- 视觉层采用现有 `card`、`muted`、`border`、`input`、`foreground`、`primary`、`accent` 等 token，不新增或重定义全局 token。
- 查询、重置、表级 select 即时筛选、Saved View 持久化和字段归属保持原实现；新增层只负责展示、位置与折叠可访问性。
- 审计模式按常规、边界清楚且可逆的 UI 实施记录为 `self`；若发现跨表合并筛选、对象名无法本地化或多个列表争抢页级插槽，则先回到方案复查。
- D-002 局部修订 D-001：视图选择/保存视图仍使用标题区局部插槽，但页面 actions（包括列配置）不得进入该插槽；筛选操作必须是筛选网格的末尾 cell；分页 footer 仅改 DOM/样式，不改 query、Saved View 或分页状态逻辑。
- D-003 追加 C7 并修订 D-002 第 3 项：筛选网格末尾操作单元只保留重置与展开/收起；搜索提交按钮回到其关键词输入所属单元并贴合为同一控件（恢复 A-003 / W13 T-03）。D-002 第 1 项细化为：视图标签组所在 surface 同时承载视图管理操作与新建视图表单，页面 actions 行不承载视图表单。
- D-004 追加 C8 并局部修订 D-002 §5 / D-001 的“不新增全局语义 token”边界：允许新增 `--control`/`--control-foreground` 表达“可交互控件填充面”语义（须同时进入 `index.css` 双层声明、`@theme inline` 映射与 `theme.test.ts` 结构守卫），仍禁止重命名或重定义既有 token 值。同时统一页面 actions 行高度，并要求折叠首行无隐藏项时不渲染展开/收起按键。
