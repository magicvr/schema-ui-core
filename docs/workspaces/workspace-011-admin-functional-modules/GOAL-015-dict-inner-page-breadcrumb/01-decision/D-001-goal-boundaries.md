---
id: D-001
goal: GOAL-015-dict-inner-page-breadcrumb
title: 立项边界：数据字典内页 + 面包屑（R4）
date: 2026-08-14
status: accepted
parent: GOAL-015-dict-inner-page-breadcrumb
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# D-001 · 立项边界（内页 + 面包屑 R4）

## 1. 背景

- 用户 2026-08-14 反馈：
  1. 条目页应只显示当前类型的条目；新增/编辑默认类型键指向该类型，表单显示类型名称（只读），提交传类型键。
  2. 2+ 级内页（首页→数据字典→条目）缺面包屑/返回按钮；需系统性修正。

## 2. 范围（S1 细化）

1. API：条目 List 新增 `dictKey` 精确过滤参数（与 q/sort/order/page/pageSize 组合）。
2. Schema：openEntries 导航带 `?dictKey=<key>`；条目页 dataSource 带过滤；create/edit 表单 dictKey 只读显示类型名称（值从 route query / row 注入，提交仍传 dictKey）。
3. Web：面包屑组件（路由栈驱动：首页 > 上级页 > 当前页）+ 返回按钮；所有 2+ 级页面通用。

## 3. 排除项

- 不改路由协议形状（path/query 保持）；面包屑为 shell 渲染层。
- 不做多级动态路由（如 /dictionary/:key/entries 路径形态变化）——query 参数足够且向后兼容。

## 4. 信息就绪

I-001~I-004（S1/S4）见 00-meta；均 open。
