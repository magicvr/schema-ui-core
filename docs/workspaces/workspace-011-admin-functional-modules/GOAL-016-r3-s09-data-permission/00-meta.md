---
id: GOAL-016-r3-s09-data-permission
title: R3-S09 · 数据权限（行级/数据范围）
status: active
parent: GOAL-001-admin-functional-modules
created: 2026-08-15
updated: 2026-08-15
version: 0.3.0
progress: 3/5
---

# GOAL-016-r3-s09-data-permission · 数据权限（行级/数据范围）

## 概述

常用档 S-09（I-011-001 §4；R3 第三批次，2026-08-15 立项）：在既有资源级 RBAC（roles/permissions，I-011-001 C-02）之上扩展**行级/数据范围**权限——数据可见范围建模（全部 / 组织范围 / 本人等）、查询侧过滤下推与管理 UI。企业后台高频能力；基架 RBAC 覆盖资源/动作，无数据行级语义。

## 当前边界（立项；S1 方案冻结细化）

- 数据范围模型为 RBAC **扩展**（新增范围维度 + 合成语义），不替代既有角色/权限键与继承。
- 过滤下推：查询侧注入（服务/存储层），带审计；不改变 Manifest 装配语义。
- 管理 UI 与 Profile 归属（admin 默认集候选，内容扩展先例）。
- 不引入多租户/跨区语义（Charter 非目标）；领域数据（订单/钱包）范围规则留领域台账。

## 成功标准与路线图（P-001）

- [x] **S1 · 方案冻结**：数据范围模型（作用域集合 / 继承 / 与角色权限的合成）、过滤下推路径、权限键与端点、协议对照（独立口径，I-011-001 §7 必办）、Profile 归属；方案级 self 审视 + **grok build independent（data 门禁，grok-4.6 · high）**（D-002，2026-08-15）
- [x] **S2 · 实现**：模块 provider + 范围/过滤能力 + schema 页 + 测试（E-003，2026-08-15）
- [x] **S3 · 验证**：单元/集成 + 全量回归（go 全绿 / web 969/969；e2e 双 profile 归 S5 波次）（E-004，2026-08-15）
- [x] **S4 · go 影响判定 + 自审**（D-004 不暂挂 + A-006 pass，2026-08-15）
- [ ] **S5 · 关门**：独立审计（grok build）+ 关门 + goal-tree 同步

progress: 3/5 由五个等权检查点派生（S1~S4 完成后更新）。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 数据范围模型：作用域集合（全部/组织/本人…）、继承与角色权限合成语义 | S1 方案 | S1 | 对照既有 RBAC（roles/permissions-inheritance）+ 协议清单 | **verified** | — | D-002 §1（2026-08-15；all/self v1，org 按 I-004 排除） |
| I-002 | required | 过滤下推集成点：list/search 查询面注入位置与审计 | S1 方案 | S1 | 对照 renderer/schema-render 查询链路 | **verified** | — | D-002 §2（2026-08-15；resourceFilter→ResourceEntity.List 边界） |
| I-003 | non-blocking | Profile 归属与权限键（admin 默认集？） | S1 方案 | S1 | S-01/S-02 内容扩展先例 | **verified** | — | D-002 §3（2026-08-15；admin 默认集 + data-permission.read/write） |
| I-004 | non-blocking | 「组织范围」作用域与未立项 B-10（组织/部门/岗位）的依赖裁定：降级 / 桩 / 本波不纳入 | S1 方案 | S1 | A-002 016-F-003 登记；对照 B-10 触发条件（I-011-001 §5） | **verified** | — | D-002 §1（2026-08-15：本波不纳入，枚举留扩展位） |
| I-005 | non-blocking | owner_column 白名单校验（A-007 F-001）——首次登记生产资源前必办 | 登记首个生产资源 | 触发 | 登记时按白名单校验 owner 列并拒绝未知列 | open | 触发=首个生产资源登记（当前 enforceable=nil 无法落策略） | 待确认 |

## 审计策略

数据权限属 **data 门禁**（P-003 independent）：S1 方案冻结与 S5 关门必须 grok build independent（用户书面偏好沿用：grok-4.6 · reasoning high）。

## 父目标

- [GOAL-001-admin-functional-modules](../GOAL-001-admin-functional-modules/00-meta.md)

## 台账布局

本目标从首条记录起使用 01-decision/、02-execution/、03-audit/ 平铺 ledger；索引与目录条目共同构成正式记录。
