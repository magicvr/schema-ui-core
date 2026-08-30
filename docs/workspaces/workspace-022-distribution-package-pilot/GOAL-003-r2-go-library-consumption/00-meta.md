---
id: GOAL-003-r2-go-library-consumption
title: R2 · Go 库包闭环
status: done
parent: GOAL-001-distribution-package-pilot
created: 2026-08-29
updated: 2026-08-29
version: 0.3.0
progress: 4/4
---

# GOAL-003 · R2 · Go 库包闭环

## 概述

承接 Root R2 检查点与 VP-022 退出判据 #1：验证「空下游仓仅 `go get` + 自建组合根装配 kernel + ≥1 标准模块，功能基线等价」的可行路径。**首个实验结论（E-001）**：Go `internal` 命名空间规则阻止外部模块导入 `apps/api/internal/*`——包化前必须将 A/B 层移出 `internal/`（外移方案 = 关键决策，D-001 提案待用户裁决）。本目标以「internal 外移重构 + 黄金下游仓装配闭环」完成判据 #1。

## 成功标准（阶段检查点）

- [x] **S1 · 阻断验证与外移方案**：internal 规则实验证据（E-001）；外移方案评估与**用户裁决**（D-001/D-002：方案 A 目录提升）
- [x] **S2 · 外移重构**：A 层（kernel）+ B 层（modules）移出 `internal/`；223 文件 import 改写；build 0 + 全量回归（A-001 pass；freeze-face v1.0.1 勘误）
- [x] **S3 · 黄金下游仓装配闭环**：`golden-consumer` 自建组合根装配 kernel + users（方案 β：`assembly` 公开装配工厂）——迁移 apply（fresh=true）+ 贡献注册 = Descriptor 声明 + 零 internal 命名；**F-001/F-002 fixed**（A-002）
- [x] **S4 · 关门**：self 审计 A-002 `pass`（0 required；PG external 消费 = 有界 residual F-005，R4/R5 复审）；**判据 #1 满足声明成立**；GOAL-003 `done 4/4`

## 信息就绪

| ID | 级别 | 所需信息 | 影响门禁 | 最晚阶段 | 状态 | 证据 |
|----|------|----------|----------|----------|------|------|
| I-001 | required | kernel 公共 API 冻结面清单 | R1 冻结 / R2/R3 实施 | R1 | **collecting**（S3 收尾：B 层符号回填 + 外移后路径同步） | `freeze-face` + GOAL-003 E 条目 |
| I-005 | required | internal 外移方案（目录提升 / 独立模块 / 保持+发布拷贝）及影响面 | S2 重构 | S1 | **collecting**（E-001 实验证据；D-001 提案待用户裁决） | GOAL-003 E-001 / D-001 |

## 父目标

- `GOAL-001-distribution-package-pilot`

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；条目见各索引。