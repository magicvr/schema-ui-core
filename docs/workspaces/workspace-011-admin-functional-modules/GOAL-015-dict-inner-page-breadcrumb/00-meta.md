---
id: GOAL-015-dict-inner-page-breadcrumb
title: 数据字典内页（按类型过滤）+ 面包屑层级导航（R4）
status: active
parent: GOAL-001-admin-functional-modules
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
progress: 2/5
---

# GOAL-015 · 数据字典内页（按类型过滤）+ 面包屑层级导航

## 概述

用户 2026-08-14 反馈两个系统性问题：

1. **数据字典内页**：从类型行「条目」进入的应该是该类型的内页——只显示属于该类型的条目；新增/编辑默认类型键指向该类型（表单显示类型名称只读，提交传类型键，不需要用户填）。
2. **面包屑层级导航**：超过 2 级的页面（如 首页→数据字典→条目内页）缺乏返回上层的方式；业界通用实践是面包屑 + 返回按钮，需系统性修正（所有 2+ 级内页通用）。

## 当前边界

- 范围：API 条目 List 增加 dictKey 精确过滤；openEntries 导航带 dictKey query；条目页表单 dictKey 只读（显示类型名/传类型键）；web 面包屑组件（路由栈驱动）+ 返回按钮。
- **不**改变现有路由协议形状；面包屑为 shell 层 UI（渲染层）。

## 成功标准与路线图（P-001）

- [x] **S1 · 方案冻结**：dictKey 过滤契约 + 内页导航 + 面包屑路由栈；协议增补门禁 P-1/P-2 登记 I-005（D-002/A-001/E-002，2026-08-14）
- [x] **S2（部分）**：服务端 dictKey 过滤 + 面包屑组件（门禁期先行，E-003）；schema 改造待 I-005 解除
- [ ] **S3 · 验证**：过滤回归 + 内页实测 + 面包屑测试 + 全量回归
- [ ] **S4 · go 影响判定 + 自审**
- [ ] **S5 · 关门**：独立审计（grok）+ required 闭合 + goal-tree 同步

progress: 2/5 由五个等权检查点派生。

## 审计策略

独立审计沿用 grok build（用户书面偏好）；过滤参数为 API 契约扩展（compatibility 门禁），S5 独立审计。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 |
|----|------|-----------------|----------|--------------|-----------------|------|
| I-001 | required | dictKey 过滤参数形状（?dictKey= 与现有 q/sort/page 组合） | S1 方案 | 现有 resourceFilter 对照 | open |
| I-002 | required | 内页表单 dictKey 只读绑定（显示类型名/传类型键；create 默认值来自 route query） | S1 方案 | 现有 recordSource/context.route.query 机制对照 | open |
| I-003 | required | 面包屑路由栈方案（history 驱动；返回按钮语义；与 HOST_OWNED_PATHS/route-not-found 交互） | S1 方案 | 现有 App.tsx 路由对照 | open |
| I-004 | required | go 影响判定（List 过滤参数/契约扩展） | S4 | VP-008 接口对照 | open |
| I-005 | required | **协议增补门禁**：table dataSource query 注入（queryMapping，P-2）——用户 2026-08-14 裁决先去上游 schema-ui-docs 增补，本目标暂停实施依赖项；P-1（toolbar navigateMapping）已撤销（「条目」是行 action，协议已有） | S1 方案 → 实施 | 上游协议仓库变更 + vendor 重 pin | **open（用户增补中）** |

## 依赖

- 无外部波次依赖；基于现有 schema 协议（recordSource path 绑定、route.query）扩展。

## 父目标

- [GOAL-001-admin-functional-modules](../GOAL-001-admin-functional-modules/00-meta.md)

## 台账布局

本目标从首条记录起使用 `01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；索引与目录条目共同构成正式记录。
