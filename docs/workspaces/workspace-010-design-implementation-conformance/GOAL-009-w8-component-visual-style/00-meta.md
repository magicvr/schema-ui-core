---
id: GOAL-009-w8-component-visual-style
title: W8 · 组件视觉样式优化（语种下拉 / 明暗切换按钮 / 下拉暗色审计）
status: done
parent: GOAL-001-design-implementation-conformance
created: 2026-08-14
updated: 2026-08-14
version: 0.2.0
progress: 5/5
---

# GOAL-009 · W8 · 组件视觉样式优化

## 概述

本子目标是 VP-010 / workspace-010 的**第八波**（用户 2026-08-14 裁决）：顶部导航与表单下拉类组件的视觉一致性整改，承载治理上下文，包含三部分：

1. **语种切换器下拉重构**（已交付，46292e5，2026-08-14）：原生 select → 图标触发（lucide Languages）+ 自制无 Portal 下拉，暗色 token 化样式 + 选中 ✓；**本目标登记该事实作为治理内容第一部分**。
2. **明暗色切换按钮统一**：与语种/通知铃图标按键样式不统一（outline Button vs size-9 ghost）；鼠标悬停文字提示为固定英文「Toggle color theme」，不随语种切换——修复为 i18n 语义提示。
3. **下拉控件暗色主题审计**：全部下拉/选择控件在暗色主题下选项弹层默认亮色、悬停才变暗的 bug 排查与修复（语种切换已通过替换控件消除；确认其余控件无同类问题，有则修正——重点：表单原生 select 的 color-scheme）。

## 当前边界

- 范围：apps/web 顶部导航图标按键（ThemeToggle）、语种下拉（LocaleSwitcher）、通知铃（复核）、表单原生 select（SelectField dark:scheme-dark）；i18n 键；相关测试。
NaN

## 成功标准与路线图（P-001）

- [x] **S1 · 方案冻结**：三部分范围 + 已交付事实登记（语种下拉 46292e5）+ 审计清单（D-001/E-001）
- [x] **S2 · 实施**：ThemeToggle 样式统一 + tooltip i18n；SelectField dark:scheme-dark；自定义面板复核（E-002）
- [x] **S3 · 验证**：单测 + 活栈暗色 DOM 复核（面板/选项/悬停）+ 全量回归 963/963（E-003）
- [x] **S4 · go 影响判定 + 自审**：go（不 held，不暂挂）+ A-001 pass
- [x] **S5 · 关门**：goal-tree 同步

progress: 5/5 由五个等权检查点派生（S1～S5 全勾）。

## 审计策略

样式一致性整改为低风险可逆维护：self 审计（A-001）；不触发独立审计（无 security/data/migration 门禁）。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 |
|----|------|-----------------|----------|--------------|-----------------|------|
| I-001 | required | 全量下拉/选择控件清单（原生 select、自制菜单、外部库） | S1 方案 | 全仓扫描（1 原生 select + 2 自制面板 + 无外部库） | **closed**（E-001） |
| I-002 | required | 暗色下原生 select 弹层渲染机制（root color-scheme vs 控件级声明） | S1 方案 | root 级联 + 控件级 scheme-light dark:scheme-dark（E-002/E-003 实测 color-scheme=dark） | **closed** |
| I-003 | required | go 影响判定（shell/表单样式不改 Profile/协议） | S4 | VP-010 接口对照 | **closed**（A-001：go 不 held） |

## 父目标

- [GOAL-001-design-implementation-conformance](../GOAL-001-design-implementation-conformance/00-meta.md)

## 台账布局

本目标从首条记录起使用 `01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；索引与目录条目共同构成正式记录。
