---
id: GOAL-011-w10-account-page-conformance
doc: execution
status: active
parent: GOAL-001-design-implementation-conformance
created: 2026-08-15
updated: 2026-08-15
version: 0.1.0
---

# E-004 · 参考样式裁决回退 + 通用表格组件样式刷新（关门记录）

## 参考样式对齐：裁决不采纳并回退

用户 2026-08-15 要求按 VP-005 参考样式（schema_ui_core_data_table）实现标题行与翻页控件行。实现后（一体化卡片页脚 + chevron-only + "Showing X to Y of Z entries" + 标题副标题，未提交），用户实测裁决：**「事实证明没有之前的好看，撤销修改要求」**，并手动回退全部相关改动（account.json descriptionKey、i18n showingEntries/session.description、DataTable footer 插槽、测试断言等）。

I-001 / I-003 按 P-003 **user-overruled** 闭合：不采纳参考样式，保持现有语义 token 与页码导航。

## 通用表格组件样式刷新（用户确认后提交 43ee84f）

在**保持原有页脚/导航形态**（页码按钮 + 摘要行）的前提下，按用户提出的现代表格规范优化 `DataTable` / `SchemaTable`：

1. **列宽与弹性布局**：`DataTableColumn` / `SchemaTableColumnSpec` 支持 `width` / `minWidth`（px 或 CSS 长度），th/td 应用；未配宽度的列保持内容驱动 auto 布局（不再平均拉扯）。
2. **通用单行截断 + Tooltip**：所有文本单元格默认 `truncate` + `title`（悬停显示完整内容）；`truncate` 列 16rem、普通列 20rem、显式 width 列由自身约束。
3. **空值兜底**：null / undefined / 空字符串 → `—`（muted）。
4. **表头层级**：`bg-muted/40→/30`（更轻背景）+ 底部分割线；`font-semibold→font-medium`；padding 与数据行对齐。
5. **操作列轻量化**：行内按钮实色带框 → ghost（无边框，hover 出背景）；toolbar 主按钮不动。
6. **行悬停**：所有行 `hover:bg-accent/40`（可点击行 /50），明暗自适应。
7. **内边距**：`px-3 py-2 → px-4 py-3`（水平 16 / 垂直 12），`align-middle` 居中。
8. **时间显示格式化**（新 `lib/datetime.ts` `formatDisplayTime`）：ISO-8601 时间戳 → 本地 `YYYY-MM-DD HH:mm`，应用于表格单元格（含移动卡片）与 recordView 详情值；表单编辑控件保持原始值。
9. **页脚文字偏移**：摘要行（X 条 · 第 X 页）按用户三轮微调定格 `pl-0.5`（右移 2px，与第一列视觉对齐）。

## 验证

- Web 全量 **991/991**（新增 datetime 4 测试 + 空值兜底/通用截断断言更新）、tsc 0；Go 无改动。
- 参考样式回退后工作树干净（git status 仅本轮 7 个文件）。

## 状态

- I-001 / I-002 / I-003 全部 **closed**；S1～S4 完成；本目标 2026-08-15 **关门（4/4）**。
