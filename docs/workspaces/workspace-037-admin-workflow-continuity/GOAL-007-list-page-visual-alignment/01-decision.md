---
id: GOAL-007-list-page-visual-alignment-decisions
doc: decision
status: active
parent: GOAL-001-admin-workflow-continuity
created: 2026-09-18
updated: 2026-09-18
version: 0.3.0
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

## 当前投影

- R6 沿用 VP-037、当前 workspace 和 Root；用户选择回开 GOAL-007 追加 C5/C6，不新建 VP/workspace 或 GOAL-008。
- 视觉层采用现有 `card`、`muted`、`border`、`input`、`foreground`、`primary`、`accent` 等 token，不新增或重定义全局 token。
- 查询、重置、表级 select 即时筛选、Saved View 持久化和字段归属保持原实现；新增层只负责展示、位置与折叠可访问性。
- 审计模式按常规、边界清楚且可逆的 UI 实施记录为 `self`；若发现跨表合并筛选、对象名无法本地化或多个列表争抢页级插槽，则先回到方案复查。
- D-002 局部修订 D-001：视图选择/保存视图仍使用标题区局部插槽，但页面 actions（包括列配置）不得进入该插槽；筛选操作必须是筛选网格的末尾 cell；分页 footer 仅改 DOM/样式，不改 query、Saved View 或分页状态逻辑。
