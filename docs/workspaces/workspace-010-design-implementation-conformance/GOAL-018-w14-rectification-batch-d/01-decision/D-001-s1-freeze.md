---
id: GOAL-018-w14-rectification-batch-d
doc: decision
status: active
parent: GOAL-015-w14-user-perspective-review
created: 2026-08-17
updated: 2026-08-17
version: 0.1.0
---

# D-001 · GOAL-018 S1 方案冻结（F-11～F-14）

## 决策

### F-11 · 必填字段标记

- 表单字段 `required` 时：label 追加 `*`，主要控件（input/password/select/textarea/number/date/dateRange/upload）设置 `aria-required`。
- 不做额外视觉重设计，仅补可访问性/可见标记。

### F-12 · 确认对话框焦点

- `ConfirmDialog` 增加：初始焦点到 Cancel、ESC 取消、Tab 焦点圈（首尾循环）。
- 保持 `role="dialog"` + `aria-modal`。

### F-13 · 桌面表格键盘选中

- `DataTable` 桌面 `<tr>` 在可点击行上增加 `tabIndex=0` 与 Enter/Space 触发 `onRowClick`；交互单元格（按钮/链接等）不触发行选中。

### F-14 · 小缺口

- 移动卡片列表渲染全部剩余内容列（不再静默丢弃第 4+ 列）。
- 单选 select 空值显式显示占位项（`feedback.selectPlaceholder`），不再误导显示第一项。
- 列排序第三次点击（desc 后再点）清除排序。
- Tabs 增加方向键（ArrowLeft/Right/Home/End）、`aria-controls`、tabpanel `aria-labelledby`。
- 禁用按钮增加 `title` 原因（权限不足 / 先选行）。
- 字段错误增加 `aria-describedby` 关联。
- `<sm` 显示语言切换（LocaleSwitcher 不再 hidden）。
- 通知空收件箱使用 `schema.notifications.empty` 语义文案。

## 信息项更新

| ID | 状态 | 说明 |
|----|------|------|
| I-001 | **closed** | F-12 采用独立焦点圈实现（不强行复用 ModalHost，行为对齐） |
| I-002 | **closed** | F-14 排序清除交互 = 第三次点击清除 |
