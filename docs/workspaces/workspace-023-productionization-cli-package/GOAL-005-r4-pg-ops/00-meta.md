---
id: GOAL-005-r4-pg-ops
title: R4 · 覆盖与运维
status: done
parent: GOAL-001-productionization-cli-package
created: 2026-08-29
updated: 2026-08-29
version: 0.2.0
progress: 4/4
---

# GOAL-005 · R4 · 覆盖与运维

## 概述

承接 Root R4 与 VP-023 判据 #4：PG external 消费实测（workspace-022 F-005 核销）；包形态下游运维路径成文（启动/升级/迁移/备份/优雅停机 + compose 样例）；golden 仓团队化（CI 消费回归槽位）。

## 成功标准（阶段检查点）

- [x] **S1 · PG external 实测**：docker postgres:16 → golden-field 组合根 PG 方言装配（迁移台账全部 apply）+ 与 SQLite 行为一致性核对；**F-005 核销**
- [x] **S2 · 运维路径文档**：`ops-playbook`（启动/升级/迁移/备份/停机，对照主仓契约）+ compose 样例（postgres + 下游应用容器）
- [x] **S3 · golden 仓团队化**：golden-field CI 消费回归 workflow（workflow_dispatch + repository_dispatch(published) → 最新包安装 + 探针全绿）
- [x] **S4 · 关门**：A-001 self `pass`（0 required · F-001 推荐 compose 实跑留 CI）→ 判据 #4 满足声明 · F-005 核销

## 信息就绪

| ID | 级别 | 所需信息 | 影响门禁 | 最晚阶段 | 状态 | 证据 |
|----|------|----------|----------|----------|------|------|
| I-023-003 | required | PG 实例（docker postgres:16 已起 · 127.0.0.1:15432） | S1 实测 | R4 | **collecting**（S1 执行中） | — |

## 父目标

- `GOAL-001-productionization-cli-package`

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；条目见各索引。