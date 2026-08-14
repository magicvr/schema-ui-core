---
id: GOAL-013-nav-order-config
title: 导航顺序：默认清单 + 配置文件覆盖（方案 A）
status: active
parent: GOAL-001-admin-functional-modules
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
progress: 0/5
---

# GOAL-013 · 导航顺序：默认清单 + 配置文件覆盖

## 概述

用户 2026-08-14 裁决（方案 A）：导航菜单顺序由**产品级默认清单**决定（维护者在制作模块时排好），并提供**配置文件覆盖**路径（使用者改配置文件，重启生效）。当前各模块各自填 `Order` 且严重冲突（Users/Settings/Account 同 Order 1、Roles/Activity/Notifications 同 Order 2），实际顺序由 NodeID 字母兜底 → 不符合直觉。

## 当前边界

- 范围：全局默认导航顺序清单（冻结 + 测试锁定）、`sortNavigation` 改为清单优先 + Order 兜底、配置文件覆盖（依赖 workspace-10 W7 YAML 配置体系，`NAVIGATION_ORDER_FILE` 或 YAML 配置段）。
- **不**改变 Profile 默认集 / 模块矩阵 / 协议 pin；顺序属于 manifest 导航内容（内容扩展，go 不暂挂判定在 S4）。

## 成功标准与路线图（P-001）

- [ ] **S1 · 方案冻结**：默认清单顺位（Dashboard → Users → Roles → Settings → Activity → Account → Notifications → File library → Data dictionary → System monitoring → Scheduled tasks → Recycle bin）、覆盖语义（全量清单，缺项追加末尾）、非法配置行为（回退默认 + 告警）；方案级 self 审视
- [ ] **S2 · 实现**：默认清单常量 + 排序改造 + 快照测试；覆盖加载（W7 YAML 就绪后接入）
- [ ] **S3 · 验证**：排序快照 + 覆盖路径实测 + 全量回归
- [ ] **S4 · go 影响判定 + 自审**（manifest 导航内容变化 → go 判定）
- [ ] **S5 · 关门**：独立审计（grok，data 门禁）+ required 闭合 + goal-tree 同步

progress: 0/5 由五个等权检查点派生。

## 审计策略

独立审计沿用 grok build（用户书面偏好）；排序为 UI 结构非安全门禁，但含配置覆盖路径（data 门禁），S5 独立审计。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 |
|----|------|-----------------|----------|--------------|-----------------|------|
| I-001 | required | 默认清单顺位与分组需求（是否引入分组） | S1 方案 | 用户确认 + 业界对照 | open |
| I-002 | required | 覆盖载体（W7 YAML 配置段 vs 独立文件）与非法配置行为 | S1 方案 | W7 依赖对照 | open |
| I-003 | required | 新模块（S-05~S-14）加入时的清单维护规则 | S1 方案 | 波次流程对照 | open |
| I-004 | required | go 影响判定（导航内容扩展） | S4 | VP-008 接口对照 | open |

## 依赖

- **workspace-10 W7（GOAL-008-w7-yaml-config）**：配置文件覆盖载体（YAML 配置段 / env 指向文件）。默认清单部分不依赖，可先行。

## 父目标

- [GOAL-001-admin-functional-modules](../GOAL-001-admin-functional-modules/00-meta.md)

## 台账布局

本目标从首条记录起使用 `01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；索引与目录条目共同构成正式记录。
