---
title: 视觉方向摘要 · Stitch 定稿指针
status: active
doc_type: design-input
created: 2026-08-09
updated: 2026-08-09
parent: GOAL-001-design-system-and-ui-experience
version: 0.1.0
related_decision: D-004
related_execution: E-004
---

# 视觉方向摘要（仓库内指针）

> **权威决策**：`01-decision/D-004-visual-direction-freeze.md`  
> **事实台账**：`02-execution/E-004-stitch-visual-refs.md`  
> 完整截图在本地 `raw/`（**gitignore**），本文件不复制 PNG。

## 本地路径（开发机）

```text
raw/stitch-vp005-visual-refs/
  00–07 … 提示词与清单
  exports/
    notes.md
    stitch_schema_ui_core_admin_console/
      schema_ui_core_sign_in/
      schema_ui_core_sign_in_mobile/
      schema_ui_core_overview/
      schema_ui_core_overview_dark_mode/
      schema_ui_core_overview_mobile/
      schema_ui_core_data_table/
      schema_ui_core_data_table_mobile_list/
      schema_ui_core_users_management/
      schema_ui_core_users_management_mobile/
      schema_ui_core_search_table/
      schema_ui_core_data_display/
      schema_ui_core_form_controls/
      schema_ui_core_form_with_reactions/
```

每屏至少 `screen.png`；多数含 `code.html`（**非生产源**）。

## 冻结约束（操作摘要）

| 面 | 约束 |
|----|------|
| 气质 | Linear / Vercel Dashboard；中性近黑主按钮；shadcn new-york |
| 桌面壳 | 顶栏 + 左栏 ~256px |
| 移动壳 | 汉堡 + 导航抽屉；无底部 Tab |
| 列表 | 桌面密表；移动卡片列表 |
| 详情 | 右栏/Drawer 或移动 Sheet；Modal 仅短编辑/Confirm |
| Token 代码 | 仍以 D-002/D-003 + `apps/web/src/index.css` 为准 |

## 实施参考优先级

1. overview（浅 / 暗 / 移动）  
2. data_table + mobile list  
3. users_management + mobile  
4. sign_in  
5. search_table / data_display / form_*  

## 非结论

- 本摘要 **不**勾选 S1–S5。  
- 本摘要 **不**关闭 F-002。  
- 字段以协议 schema 为准，不以 Stitch 示例字段为准。  
