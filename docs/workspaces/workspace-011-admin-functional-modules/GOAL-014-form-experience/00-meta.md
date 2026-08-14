---
id: GOAL-014-form-experience
title: 表单体验：字段级校验与错误展示 + 弹窗布局（R4）
status: done
parent: GOAL-001-admin-functional-modules
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
progress: 5/5
---

# GOAL-014 · 表单体验：字段级校验与错误展示 + 弹窗布局

## 概述

用户 2026-08-14 裁决：新增 R4 波次子目标（完整三部分）：

1. **服务端字段级错误明细**：create/patch 校验失败时返回具体字段与原因（当前 INVALID_PATCH_FIELD 只有通用 message，前端无法定位字段）；
2. **前端表单校验与内联报错**：schema 字段可声明必填/格式约束，提交前校验 + 错误内联显示在对应输入项（当前空缺项只弹通用错误）；
3. **弹窗布局参考业界主流**：Ant Design / Element Plus / shadcn 表单弹窗惯例（单列或可配置列数、宽度/栅格响应式、内联错误样式）。

## 当前边界

- 范围：schema 协议扩展（字段约束声明 + 字段级错误契约）、API 错误响应结构、renderer 表单校验/展示、modal 布局样式。
- **不**改变既有错误码语义（INVALID_CREATE_FIELD / INVALID_PATCH_FIELD 保留为兼容别名或 supersede 需 S1 定）；不引入前端框架切换。

## 成功标准与路线图（P-001）

- [x] **S1 · 方案冻结**：字段级错误契约（可选 fieldErrors）+ 约束最小集 + 单列布局方案（D-002/A-001/E-002，2026-08-14）
- [x] **S2 · 实现**：服务端 fieldErrors + 前端 validateFieldValues/内联 + 单列布局 + schema 约束示范（E-003，2026-08-14）
- [x] **S3 · 验证**：validateFieldValues 单测 + 911/911 web + go 全绿 + HTTP fieldErrors 冒烟（E-004，2026-08-14）
- [x] **S4 · go 影响判定 + 自审**（E-005：go 不 held；A-002：pass，2026-08-14）
- [x] **S5 · 关门**：grok 审计（A-003）fail → F-001/F-002 fixed + F-003~F-007 处置；关门（E-006，2026-08-14）

progress: 5/5 由五个等权检查点派生。

## 审计策略

独立审计沿用 grok build（用户书面偏好）；错误契约变更属 compatibility 门禁，S5 独立审计。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 |
|----|------|-----------------|----------|--------------|-----------------|------|
| I-001 | required | 字段级错误响应结构（与既有 {error,message} 信封兼容策略） | S1 方案 | 现有错误契约对照（errorcatalog / 前端 readResourceApiError） | **closed**（D-002 §2：可选 fieldErrors 键） |
| I-002 | required | schema 字段约束最小集（required/pattern/min/max/长度）与协议版本影响 | S1 方案 | 业界对照（JSON Schema / AntD 规则） | **closed**（D-002 §3：required/pattern/minLength/maxLength + min/max；不 bump 版本） |
| I-003 | required | 弹窗布局方案（列数可配/单列默认）与既有两列 fixture 兼容 | S1 方案 | 现有 modal fixture 对照 | **closed**（D-002 §4：单列默认 + columns 可配 + modal width；渲染层变化） |
| I-004 | required | go 影响判定（错误契约/协议扩展） | S4 | VP-008 接口对照 | **closed**（E-005：go 不 held） |

## 依赖

- 无外部波次依赖；基于现有 schema 协议（protocolVersion 2.7）扩展需评估版本策略。

## 父目标

- [GOAL-001-admin-functional-modules](../GOAL-001-admin-functional-modules/00-meta.md)

## 台账布局

本目标从首条记录起使用 `01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；索引与目录条目共同构成正式记录。
