---
id: GOAL-005-r4-fork-comparison
title: R4 · fork 对照计时实验（同一演进集 v0.3.0→v0.4.0：fork 同步 vs 包 bump 实测对比）
status: done
parent: GOAL-001-distribution-formalization
created: 2026-08-29
updated: 2026-08-29
version: 0.3.0
progress: 4/4
---

# GOAL-005 · R4 · fork 对照计时实验

## 概述

承接 Root R4 与 VP-024 判据 #4：以**用户定案样本 = v0.3.0→v0.4.0 真实演进集**（`apps/api/v0.3.0` → `apps/api/v0.4.0`：serve 面新增 + CLI/模板改造），在 **fork 同步模型 vs 包 bump 模型**双端实测对比——耗时 / 冲突计数 / 契约迁移成本，产出**定量结论**（核销 VP-022 判据 #6 遗留的「包 vs fork 实测对比」半项与 go 后清单「fork 对照计时」项）。

## 成功标准（可验证检查点）

- [x] C1：**fork 模型实测**：fork-sim（v0.3.0 基线 + 2 定制点）`git merge apps/api/v0.4.0` → **冲突 1**（main.go.tmpl 内容冲突 · QUICKSTART auto-merge）· 解冲突 2 改写点 · 构建 exit 0（12.9s；审计复跑 10.7s）
- [x] C2：**包模型实测**：golden-field 重演 v0.3.0→v0.4.0（go.mod 1 行 + 薄封装 + .npmrc）→ 冲突 0 · tidy+build 4.8s（审计复跑 4.4s）· serve/healthz 200
- [x] C3：**定量对比报告**：`attachments/fork-comparison-report.md`（耗时矩阵 ≈13.2s vs ≈4.8s · 冲突 1 vs 0 · 改写点 2 vs 0 · 暖口径注明）——v0.4.0 定量实证（VP-022 go/no-go 定性映射升级）
- [x] C4：独立审计（grok）→ A-002 `pass`（0 required）→ **关门（Root 4/7）**

## 方案与路线（P-001）

| 阶段 | 内容 | 状态 |
|------|------|------|
| S1 | 实验设计定档：样本集切分（tag 区间）· 双模型口径 · 计时规则（冷执行注明缓存） | **已关门**（2026-08-29 · D-001 · I-024-003 verified） |
| S2 | fork 模型实测（fork-sim 仓 + merge + 解冲突 + 构建验证） | **已关门**（2026-08-29 · E-002） |
| S3 | 包模型实测（golden-field worktree 重演 v0.3.0→v0.4.0） | **已关门**（2026-08-29 · E-002） |
| S4 | 对比报告 + 独立审计（grok）→ 关门 | **已关门**（2026-08-29 · A-002 `pass`（0 required）· F-001~F-004 fixed） |

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-024-003 | required | fork 对照的同一演进集样本（V 演进选择） | 判据 #4 | R4 | 样本设计与 fork 基线 | **verified**（2026-08-29 用户裁决：v0.3.0→v0.4.0 真实演进集） | — | 用户裁决（Round 3 末） |

## 父目标

- `GOAL-001-distribution-formalization`

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger。

## 备注

- 审计模式：S4 关门 = independent（grok build · R2/R3 先例）。
- 计时口径：本地 Windows 环境（docker 工具可用）；分别注明冷（清 Go 缓存）/暖口径；不主张跨机可复现绝对秒数（相对对比为主）。