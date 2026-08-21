---
id: E-006-u05-menu-portal-fix
doc: execution-entry
goal: GOAL-012-w11-mfa-ux-review
date: 2026-08-15
status: done
parent: GOAL-001-design-implementation-conformance
created: 2026-08-15
updated: 2026-08-15
version: 0.1.0
---

# E-006 · U-05 行操作菜单层级修复（post-closeout follow-up）

## 用户报告

「列表点行 action 中的 ⋯ 按键时候弹窗弹出的层级不对，好像嵌套在表单控件里面了。」

## 根因

RowActionsMenu 原实现为 `absolute` 定位在行按钮的 `relative` 容器内；而桌面表格容器带 `overflow-x-auto`——按 CSS 规则，一个轴为 auto 时另一个轴会被强制计算为 auto，因此超出表格边界的菜单被容器**裁剪**，视觉上如同嵌套在表格/控件内部；且菜单 `z-30` 处于表格局部层叠上下文，无法盖住后续内容。

## 修复（apps/web/src/renderer/schema-table.tsx）

1. 菜单改用 **createPortal 渲染到 document.body**，彻底脱离表格容器的 overflow 裁剪与层叠上下文；
2. **position: fixed + zIndex 60**（高于 ModalHost 的 z-50），基于触发按钮 `getBoundingClientRect()` 视口坐标定位（右下对齐，空间不足时向上翻转，左右夹紧视口）；
3. 打开时聚焦首个可用菜单项（a11y）；外点/Escape 关闭；**scroll（capture，含表格内部滚动）与 resize 关闭**——fixed 定位下锚点失效即关闭，避免出现悬浮错位菜单；
4. aria-haspopup 补齐。

## 测试

- 新增 2 个回归用例（schema-table.test.tsx）：菜单渲染在 document.body 且不在表格容器内（fixed + zIndex≥50 + Escape 关闭）；外点 pointerdown 与 scroll 关闭。
- 既有测试更新：菜单项查找范围从 container 改为 document.body（schema-crud/representative-pages/error-localization/ui-bilingual/representative-pages.integration）；rowActionButton 支持菜单已打开场景（避免二次点击触发 toggle 关闭）。
- 验证：Web 全量 1004/1004；tsc 0。

## 结论

修复已落地并全量回归通过；提交见 E-005 后续 checkpoint（git log 最近一次）。