---
id: GOAL-003-r2-cli
title: R2 · CLI 闭环
status: done
parent: GOAL-001-productionization-cli-package
created: 2026-08-29
updated: 2026-08-29
version: 0.2.0
progress: 4/4
---

# GOAL-003 · R2 · CLI 闭环

## 概述

承接 Root R2 与 VP-023 判据 #2：`schema-ui` CLI（用户定案：单命令 + Go 单二进制 + create/add/upgrade）——`create` 生成下游骨架、`add` 模块装配、`upgrade` registry 升级；与手工文档路径双轨对照；**CLI 产物 = golden-field 骨架等价**。R1 挂账（GOAL-002 A-001 F-001）：首次 CLI 发布 = registry 升级演练（bump → 安装 → 回归，含 breaking 预演）。

## 成功标准（阶段检查点）

- [x] **S1 · CLI 实现**：`apps/api/cmd/schema-ui`（标准库解析 · 零新依赖 · go:embed 模板）——`create`（骨架生成）+ `add`（可用模块清单/registry 装配）+ `upgrade`（registry 升级 + 探针回归）
- [x] **S2 · create 验证**：CLI 生成仓 == golden-field 骨架等价（双轨 diff）· 生成仓 go run/install/probes 全绿
- [x] **S3 · 发布与升级演练**：CLI 首版 registry 发布（Go tag apps/api/v0.2.0 或独立？——CLI 在 apps/api 模块内 → 随模块发布）+ **registry 升级演练**（F-001 核销：bump → 安装 → 回归，含 breaking 场景预演）
- [x] **S4 · 关门**：A-001 self ``pass``（0 required）→ 判据 #2 满足声明 · F-001 fixed（跨目标回填）

## 信息就绪

| ID | 级别 | 所需信息 | 影响门禁 | 最晚阶段 | 状态 | 证据 |
|----|------|----------|----------|----------|------|------|
| I-023-004 | required | CLI 分发形态（apps/api 模块内 cmd → `go install …/cmd/schema-ui@vX` 复用模块 tag 链） | S3 发布 | S3 | collecting（S1 实证后转 verified） | — |

## 父目标

- `GOAL-001-productionization-cli-package`

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；条目见各索引。