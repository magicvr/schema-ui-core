---
id: GOAL-008-r7-topline-and-closeout
title: R7 · 方法 B 置顶宣告 + 收口报告 + 残余复核（Root 关门）
status: done
parent: GOAL-001-distribution-formalization
created: 2026-08-29
updated: 2026-08-29
version: 0.3.0
progress: 4/4
---

# GOAL-008 · R7 · 置顶宣告与收口报告

## 概述

承接 Root R7 与 VP-024 判据 #8：QUICKSTART **方法 B 置顶**（cli+包 为默认主路径首段 · fork 第二 · Charter 措辞不动——执行层动作，用户在 VP-024 判据 #8 已裁决）；**收口报告**（8 判据核销表 · 公开消费往返实证综述 · fork 对照结论回引 · 残余复核清单全套核销或登记）；**残余复核**（hosted CI 触发注记 · shell 类型面（消费端 tsc 未验证）· GH 私有包退役评述定稿 · C 类 fork 包化面登记）。Root 关门 = independent 审计（grok）→ VP-024 closed → Root done 7/7。

## 成功标准（可验证检查点）

- [x] C1：QUICKSTART 方法 B 置顶（§0 决策块 + §1 cli+包 首节 + §2 migrate-fork + §3 fork 复现顺延）；create 钉终值 + upgrade 拉 latest 注记；命令族/终值一致性
- [x] C2：收口报告 `attachments/closure-report.md`：判据 #1–#8 核销表 + 公开消费往返实证 + fork 对照回引 + 残余复核四项逐项结论
- [x] C3：残余复核四项定稿（D-001）：hosted CI 触发 = 登记（不主张 acceptance）· shell 类型面 = 登记 · GH 私有包 = 保留 + npmjs 新消费面 · C 类 fork 面 = 未来候选
- [x] C4：Root 关门审计（grok independent）→ **pass（0 required · F-001~F-006 fixed）** → 用户书面确认 → **Root done 7/7 · VP-024 closed（VRev-053）**

## 方案与路线（P-001）

| 阶段 | 内容 | 状态 |
|------|------|------|
| S1 | 方法 B 置顶（QUICKSTART 首段调整 + 一致性通读） | **已关门**（2026-08-29 · E-001） |
| S2 | 收口报告成文（核销表 + 残余复核清单 + 宣告） | **已关门**（2026-08-29 · E-002 · closure-report） |
| S3 | 残余复核定稿（四项）（C3） | **已关门**（2026-08-29 · D-001） |
| S4 | Root 关门独立审计（grok）→ 用户确认 → 收官 | **已关门**（2026-08-29 · A-002 `pass`（0 required · F-001~F-006 fixed）· 用户书面确认 · VRev-053 · VP-024 closed） |

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| — | — | 无新增 required（R1–R6 信息项已全闭合） | — | — | — | — | — | — |

## 父目标

- `GOAL-001-distribution-formalization`

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger。

## 备注

- Root 关门审计 = independent（grok build · 项目 D-001 既例）；关门需用户书面确认（P-004 · 最终裁决点）。
- 遗留登记政策：仅接受「已书面 residual」或「未来候选登记」（hosted 实触发 / C 类 fork 面）。