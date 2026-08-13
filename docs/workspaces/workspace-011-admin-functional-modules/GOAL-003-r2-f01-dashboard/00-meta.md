---
id: GOAL-003-r2-f01-dashboard
title: R2-F01 · 仪表盘/控制台（生产 Profile home）
status: done
parent: GOAL-001-admin-functional-modules
created: 2026-08-14
updated: 2026-08-14
version: 0.4.0
progress: 5/5
---

# GOAL 003-r2-f01-dashboard · 仪表盘/控制台（生产 Profile home）

## 概述

一等公民 F-01（I-011-001 `3）：新建 `admin.dashboard` 模块，为生产 Profile（mvp/admin）提供 home 仪表盘面（指标/概览卡片）。现状：overview 仅存在于 demo（dev.examples，section+text 欢迎文案），生产面无 home dashboard。进入 mvp/admin 默认启用集 = Profile **内容**扩展（模块贡献机制 + adminFunctionalOrder 更新），**不改装配语义**。

## 当前边界

- 新增模块与 schema（dashboard 页 + 指标/卡片渲染）
- 进入 mvp/admin 默认启用集（内容扩展，非装配语义变更；R2 方案写清）
- 不改变 Profile 装配语义 / 模块矩阵 / Manifest 聚合规则 / 协议 pin / 共同门禁语义
- 协议对照按 `8 必办-1 独立做（grid-dashboard 信息性场景 → 呈现自由 + fail-open 留痕）

## 成功标准与路线图（P-001）

- [x] **S1 · 方案冻结**：方案冻结：必办-1（协议对照 grid-dashboard / 上游样例）+ 必办-3（home 装配内容扩展声明：贡献机制 + adminFunctionalOrder）+ 指标数据源设计；方案级 self 审视（D-002 / A-001，2026-08-14）
- [x] **S2 · 实现**：实现：模块 + schema + 渲染（指标卡片/概览）+ 范例页与测试（E-003 · a9642f4）
- [x] **S3 · 验证**：验证：单元/集成 + 代表页 + 全量回归（V-001～V-005）+ e2e 按可行性（E-003 · go 全绿 + web 892/892 + 冒烟）
- [x] **S4 · go 影响判定 + 自审**：go 影响判定（无影响/不暂挂 或 暂挂留痕）+ self 审计（D-003 无影响不暂挂；A-002 pass）
- [x] **S5 · 关门**：关门：全部 required 闭合 + 关门审计（A-003 conditional → F-001~F-009 处置）+ goal-tree 同步（E-004 · 2026-08-14）

progress: 0/5 由五个等权检查点派生（S1 完成后更新）。

## 审计策略

self（呈现层 + Profile 内容扩展；协议对照在 S1 留痕）；若涉及装配语义争议按 P-004 升级。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 |
|----|------|-----------------|----------|--------------|-----------------|------|
| I-001 | required | 协议面（v2.8.0）对 dashboard 是否有语义定义？ | S1 方案 | 核对 protocol-inventory `2.5 grid-dashboard、上游样例 | **closed**（D-002 §2） |
| I-002 | required | home 装配方案：admin.dashboard 进 mvp/admin 默认集 + adminFunctionalOrder 更新 | S1 方案 | 对照模块贡献机制（playbook M1–M6） | **closed**（D-002 §3） |
| I-003 | non-blocking | 指标数据源（API 聚合 vs 静态配置 vs 渲染器契约） | S2 实现 | 方案冻结时定默认 | **closed**（D-002 §1：既有端点 envelope） |

## 父目标

- [GOAL-001-admin-functional-modules](../GOAL-001-admin-functional-modules/00-meta.md)

## 台账布局

本目标从首条记录起使用 `01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；索引与目录条目共同构成正式记录。