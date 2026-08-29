---
id: GOAL-008-r7-topline-and-closeout
title: R7 · 方法 B 置顶宣告 + 收口报告 + 残余复核（Root 关门）
status: active
parent: GOAL-001-distribution-formalization
created: 2026-08-29
updated: 2026-08-29
version: 0.1.0
progress: 0/4
---

# GOAL-008 · R7 · 置顶宣告与收口报告

## 概述

承接 Root R7 与 VP-024 判据 #8：QUICKSTART **方法 B 置顶**（cli+包 为默认主路径首段 · fork 第二 · Charter 措辞不动——执行层动作，用户在 VP-024 判据 #8 已裁决）；**收口报告**（8 判据核销表 · 公开消费往返实证综述 · fork 对照结论回引 · 残余复核清单全套核销或登记）；**残余复核**（hosted CI 触发注记 · shell 类型面（消费端 tsc 未验证）· GH 私有包退役评述定稿 · C 类 fork 包化面登记）。Root 关门 = independent 审计（grok）→ VP-024 closed → Root done 7/7。

## 成功标准（可验证检查点）

- [ ] C1：QUICKSTART 方法 B 置顶（首段 cli+包 起步 · fork 为第二路径）；通读无断链（create/serve/upgrade 命令与六包终值一致）
- [ ] C2：收口报告（`attachments/closure-report.md`）：判据 #1–#8 核销表 + 往返实证 + fork 对照回引 + 残余复核清单（hosted CI 触发 / shell 类型面 / GH 包退役评述 / C 类面）逐项结论
- [ ] C3：残余复核定稿：hosted CI 触发 = 注记（等价证据已证 · 实触发随 R7 后用户指令）· shell 类型面 = 登记（JS 运行时自包含）· GH 私有包 = 保留不删（历史消费面；新消费 npmjs 公开）· C 类 fork = 保持 fork（assembly 扩展面 = R7 后候选）
- [ ] C4：Root 关门审计（grok independent）→ 用户确认 → Root done 7/7 · VP-024 closed · 收官提交

## 方案与路线（P-001）

| 阶段 | 内容 | 状态 |
|------|------|------|
| S1 | 方法 B 置顶（QUICKSTART 首段调整 + 一致性通读） | 未开 |
| S2 | 收口报告成文（核销表 + 残余复核清单 + 宣告） | 依赖 S1 |
| S3 | 残余复核定稿（四项）（C3） | 依赖 S2 |
| S4 | Root 关门独立审计（grok）→ 用户确认 → 收官 | 依赖 S2/S3 |

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