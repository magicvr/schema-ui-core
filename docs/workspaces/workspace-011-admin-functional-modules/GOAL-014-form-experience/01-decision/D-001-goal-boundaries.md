---
id: D-001
goal: GOAL-014-form-experience
title: 立项边界：表单体验（R4）
date: 2026-08-14
status: accepted
parent: GOAL-014-form-experience
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# D-001 · 立项边界（表单体验 R4）

## 1. 背景

- 用户 2026-08-14 反馈（数据字典页）：编辑时空缺项只显示 `INVALID_PATCH_FIELD: 更新字段无效`（无字段级提示）；弹窗两列布局疑似非业界主流。
- 用户裁决：新增完整三部分子目标；归属 workspace-011 → GOAL-001，编号 GOAL-014。

## 2. 三部分范围（S1 细化）

1. 服务端字段级错误明细（错误响应携带 field + reason；与既有 {error,message} 信封兼容策略 S1 定）。
2. 前端表单校验与内联报错（schema 字段约束声明最小集；提交前校验 + 错误内联）。
3. 弹窗布局参考业界主流（AntD/Element/shadcn 表单惯例；单列默认或列数可配）。

## 3. 排除项

- 不切换前端框架；不改既有错误码的 wire 形状破坏性（兼容策略 S1 定）。
- 不做运行时动态表单引擎重构（如 JSON Schema 全量生成器）。

## 4. 信息就绪

I-001~I-004（S1/S4）见 00-meta；均 open。
