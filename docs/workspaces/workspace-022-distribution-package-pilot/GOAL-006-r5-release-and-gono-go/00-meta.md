---
id: GOAL-006-r5-release-and-gono-go
title: R5 · 发布可复现与 go/no-go 报告
status: done
parent: GOAL-001-distribution-package-pilot
created: 2026-08-29
updated: 2026-08-29
version: 0.3.0
progress: 4/4
---

# GOAL-006 · R5 · 发布可复现与 go/no-go 报告

## 概述

承接 Root R5 检查点与 VP-022 退出判据 #4/#5/#6：契约冻结面（含 B 层盘点）定稿 → 发布流水线（Go tag 语义 + npm tarball 一键产出）→ golden 消费回归（tarball 安装语义）→ **go/no-go 报告**（实测对比 + Charter strategic 修订建议，按 VP-022 触发框架判向）。**R5 与 Root 关门 = `independent` 审计（项目默认执行路径：grok build · grok-4.6 · 思考 high）**。

## 成功标准（阶段检查点）

- [x] **S1 · 发布流水线**：`scripts/pack-npm-packages.mjs`（双 tgz 一键产出）+ Go tag `v0.0.2` + proxy 文档语义（E-001）
- [x] **S2 · golden 消费回归**：golden-web **tarball 安装**（registry 语义）→ 三探针全绿 + V2 能力可用（E-001）
- [x] **S3 · go/no-go 报告**：判据表 + 触发框架判向 + Charter 草案 + pin 建议（gono-go-report-v1）→ **用户裁决 = GO**（VR-050 执行 · Charter 0.3.0 · pin 2.9.0 · 22 VP re-align · VRev-050 pass）
- [x] **S4 · 关门**：冻结面 v1.2.0 定稿 + **independent 审计（grok build · A-002 conditional → 用户 P-004 全部闭合）** → 判据 #4/#5/#6 满足声明 → **GOAL-006 done 4/4 · Root done 5/5**

## 信息就绪

| ID | 级别 | 所需信息 | 影响门禁 | 最晚阶段 | 状态 | 证据 |
|----|------|----------|----------|----------|------|------|
| I-003 | required | 发布通道（npm 私有 registry vs GitHub Packages；Go tag 策略） | R5 发布可复现 | R5 | **verified**（D-001 定案 + residual：registry/proxy 上传 = go 后，复审触发 = 首次对外发布） | S1 |
| I-007 | required | 协议 pin 漂移处置（`/vision` 裁决 2.8 → 2.9？） | R5 发布 / Root 关门 | R5 | **verified**（VR-050 已执行 pin 2.9.0） | GOAL-004 E-001 |

## 父目标

- `GOAL-001-distribution-package-pilot`

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；条目见各索引。