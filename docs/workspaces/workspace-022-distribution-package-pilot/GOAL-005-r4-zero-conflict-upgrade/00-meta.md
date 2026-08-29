---
id: GOAL-005-r4-zero-conflict-upgrade
title: R4 · 零冲突升级演练
status: active
parent: GOAL-001-distribution-package-pilot
created: 2026-08-29
updated: 2026-08-29
version: 0.1.0
progress: 0/4
---

# GOAL-005 · R4 · 零冲突升级演练

## 概述

承接 Root R4 检查点与 VP-022 退出判据 #3：上游施加一次**真实演进**（A 层 additive + 协议面 additive + 新增迁移），下游（golden-consumer / golden-web）**仅 bump 版本 + 执行 changelog 迁移说明** → 功能基线回归全绿、冲突计数 = 0、全程无 git merge。

## 成功标准（阶段检查点）

- [ ] **S1 · 样本设计与基线冻结**：演进样本集（≥3 类变更：kernel additive / protocol additive / 新迁移）；V1 基线证据（golden-consumer `go run` + golden-web 三探针，E-001）
- [ ] **S2 · 上游 V2 演进**：真实代码变更 + changelog 迁移说明（V1→V2）+ 主仓回归（build + 全量测试）
- [ ] **S3 · 下游升级演练**：golden 仓仅 bump（require/依赖版本号）→ 回归全绿 · 新迁移 apply · **冲突计数 = 0 · 无 git merge**
- [ ] **S4 · 关门**：self 审计 + 判据 #3 满足声明 + 挂账复核（F-005 PG external / I-007 pin / F-003 drain 时序 → R5 前状态声明）

## 信息就绪

| ID | 级别 | 所需信息 | 影响门禁 | 最晚阶段 | 状态 | 证据 |
|----|------|----------|----------|----------|------|------|
| I-004 | non-blocking | 演练样本选择（哪类上游变更最能代表冲突压力） | R4 演练质量 | R4 | **collecting**（S1 定案：kernel/protocol additive + 新迁移） | E-001 |
| I-007 | required | 协议 pin 漂移（2.9 vs 2.8）处置（`/vision` 裁决） | R5 发布 | R5 | registered（R4 后转 /vision 触发） | GOAL-004 E-001 |

## 父目标

- `GOAL-001-distribution-package-pilot`

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；条目见各索引。