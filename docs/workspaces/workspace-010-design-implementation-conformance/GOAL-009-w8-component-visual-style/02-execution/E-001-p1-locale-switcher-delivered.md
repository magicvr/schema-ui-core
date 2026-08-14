---
id: E-001
goal: GOAL-009-w8-component-visual-style
date: 2026-08-14
status: recorded
parent: GOAL-009-w8-component-visual-style
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# E-001 · P1 语种切换器下拉重构（已交付事实登记）

## 事实

- 提交 46292e5（2026-08-14，本波次立项前已交付）：LocaleSwitcher 原生 select → 图标下拉；测试重写 6 例；web 960/960；活栈暗色 DOM 实测（面板 oklch(0.205)≈neutral-900、边框 white/10、选中 ✓、header 无 select）。
- 治理定位：用户裁决归入 W8 组件视觉样式优化目标（本目标）承载治理上下文；细节见提交信息与组件注释。

## 审计清单（P3 输入，I-001）

- 全仓 `<select`：仅 form-controls.tsx SelectField 1 处（其余无）。
NaN
NaN
NaN
