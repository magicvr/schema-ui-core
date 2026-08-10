---
id: GOAL-004-r3-bounded-pilot
title: R3 · 有界试点
status: done
parent: GOAL-001-modular-admin-architecture
created: 2026-08-05
updated: 2026-08-05
version: 0.3.0
progress: 4/4
plan_refs:
  - VP-003-modular-admin-architecture
primary_plan: VP-003-modular-admin-architecture
serves_summary: 在 R2 的薄内核和组合根基础上，以 operationlog、Activity、Settings 为有界切口，验证模块启停、Schema/Manifest 贡献、Host 边界和数据/运维闭环；不得把试点通过当作 VP 关闭。
---

# GOAL-004 · R3 有界试点

## 概述

本子目标承接 Root 的 R3 阶段，范围限定为 operationlog、Activity、Settings
及其模块化注册、Schema/Manifest/导航投影和有界旧路径清理。它承接而不重写
R1/R2 的已冻结契约；R3 通过只允许进入 R4，不关闭 VP-003。

## 愿景对齐

| 字段 | 值 |
|------|----|
| `parent` | `GOAL-001-modular-admin-architecture` |
| `plan_refs` | `VP-003-modular-admin-architecture` |
| `primary_plan` | `VP-003-modular-admin-architecture` |
| Charter | `schema-ui-core-admin-foundation@0.2.0`（经 VP-003 间接对齐） |

## 成功标准

- [x] **C1 / I-006**：盘点并核验静态 Manifest、Shell、中心注册、模块
  禁用、兼容/告警、移除和回滚边界，然后才冻结 R3 实施方案。
- [x] **C2 / A+B 实施**：按统一模块契约实现 operationlog、Activity、Settings
  的有界切片，移除试点范围内四类旧架构病灶。
- [x] **C3 / V-1～V-4 验证**：在同一 Web 构建下验证双 Profile 的启停、
  Manifest/Schema/通用渲染、冲突 fail-closed、Settings 事件和数据保留。
- [x] **C4 / D 门**：完成 A/B/C 全部证据、自审、Grok 独立审计及 required
  finding 合法闭合，并记录是否允许进入 R4。

四个检查点等权；当前为 `progress: 4/4`。完成本子目标只表示 R3 子目标
关闭，不代表 VP-003 退出或 Root 关闭。

## 信息门禁：Root I-006

| 编号 | 级别 | 必须回答的问题 | 影响 | 最晚阶段 | 收集动作 | 状态 | 证据/延期 |
|------|------|----------------|------|----------|----------|------|-----------|
| R3-I006-01 | required | 哪些静态 Manifest、Shell、中心注册和 Schema fixture 必须移除，哪些需暂保留 | R3 方案冻结、B 门 | C1 关闭前 | 源码扫描、构建/代理核对、路由清单 | verified | `attachments/r3-c2-c3-v1-v4-evidence.md`；A-004 fixed response |
| R3-I006-02 | required | 开发兼容窗口、告警行为和移除完成触发条件是什么 | R3 方案冻结、R6 移除 | C1 关闭前 | 决策记录、告警测试和同构建证据 | verified | `attachments/r3-c2-c3-v1-v4-evidence.md`；A-004 fixed response |
| R3-I006-03 | required | 模块禁用或静态/Shell 清理失败时何时回滚、如何保留数据并复核 | R3 方案冻结、R6 回滚 | C1 关闭前 | 失败演练、数据保留检查、恢复后验证 | verified | `attachments/r3-c2-c3-v1-v4-evidence.md`；A-004 fixed response |

上述三项细化 Root I-006，但不替换 Root 的 canonical 状态。C1 审计通过前不得
把本子目标标为完成，也不得冻结 R4 迁移方案；A-004 记录了 C1/C2/C3 证据完成
后的 C4 close-out。R3 关闭只允许 Root 进入 R4 阶段评估，不代表 VP-003 或
Root 关闭。

## 阶段计划

1. 收集并核验 I-006 与试点入口清单（C1）。
2. 实现模块控制的路由、Schema、Manifest 切片并清理四类试点病灶（C2）。
3. 在同一构建上运行 V-1～V-4 及 fail-closed/数据保留验证（C3）。
4. 执行自审和 Grok 独立审计，响应 required finding，记录 R3 D 门结论（C4）。

## 范围与非目标

范围包括 Kernel 组合边界、operationlog always-on、可选 Activity/Settings、
其 Schema/Manifest/导航投影、现有 Web 通用渲染路径和有界旧路径清理。

非目标包括 users/roles 全量迁移、R4 批量迁移、插件系统、VP 关闭，以及改变
冻结的 `I-PROTO-001 v0.1.3` 边界。
