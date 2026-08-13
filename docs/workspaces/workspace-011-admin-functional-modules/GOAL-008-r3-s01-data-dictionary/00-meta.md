---
id: GOAL-008-r3-s01-data-dictionary
title: R3-S01 · 数据字典（枚举/字典管理）
status: active
parent: GOAL-001-admin-functional-modules
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
progress: 0/5
---

# GOAL-008-r3-s01-data-dictionary · 数据字典

## 概述

常用档 S-01（I-011-001 §4）：新建标准 Admin 模块（`admin.data-dictionary` 候选名），提供字典/枚举管理——字典类型 + 字典条目两级 CRUD（键值/标签/启用/排序/备注），供业务表单下拉/选择复用。企业后台高频能力（中文生态尤甚），基架未覆盖。

## 当前边界

- 字典类型 + 字典条目两级 CRUD（schema 驱动代表页）
- 条目字段：键/值/标签（多语呈现按方案）/启用/排序/备注
- 权限键与操作审计（复用 operationlog）
- Profile 归属：进入 admin 默认集（S1 方案确认；Profile 内容扩展，不改装配语义）

## 成功标准与路线图（P-001）

- [ ] **S1 · 方案冻结**：两级模型（类型/条目）、字段与校验、权限键、审计设计、协议对照、Profile 归属；方案级 self 审视
- [ ] **S2 · 实现**：模块 provider + 迁移 + schema 代表页 + 字典端点 + 测试
- [ ] **S3 · 验证**：单元/集成 + 代表场景实测 + 全量回归（go test / web suite / 冒烟）
- [ ] **S4 · go 影响判定 + 自审**：go 影响判定（Profile 默认集变化触发失效检查）+ self 审计
- [ ] **S5 · 关门**：关门审计（按 P-004 确认独立 provider 或 self）+ required 闭合 + goal-tree 同步

progress: 0/5 由五个等权检查点派生（S1 完成后更新）。

## 审计策略

数据字典为常规 CRUD + 枚举数据，低风险 → 以 self 审计为主；关门门禁按 P-004 与用户确认是否引入独立审计。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 |
|----|------|-----------------|----------|--------------|-----------------|------|
| I-001 | required | 字典两级模型的字段/校验与协议对照（枚举呈现自由边界） | S1 方案 | 对照 protocol-inventory + node.schema.json 呈现自由 | open |
| I-002 | required | Profile 归属：S-01 进入 admin 默认集？mvp 保持精简？ | S1 方案 | F-01 先例（Profile 内容扩展 + adminFunctionalOrder） | open |

## 父目标

- [GOAL-001-admin-functional-modules](../GOAL-001-admin-functional-modules/00-meta.md)

## 台账布局

本目标从首条记录起使用 `01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；索引与目录条目共同构成正式记录。
