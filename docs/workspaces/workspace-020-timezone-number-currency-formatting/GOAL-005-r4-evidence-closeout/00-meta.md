---
id: GOAL-005-r4-evidence-closeout
title: R4 · 证据与关门（快测 + 双 locale 范例 · 无越界 · required = 0）
status: done
parent: GOAL-001-timezone-number-currency-formatting
created: 2026-08-26
updated: 2026-08-27
version: 0.3.0
progress: 4/4
---

# GOAL-005 · R4 证据与关门

## 概述

承载 Root 路线图 **R4（证据与关门）**：把 R1～R3 交付收成可核对的证据矩阵（快测 + `zh-CN`/`en-US` 双 locale 范例）、逐条核对 Root 成功标准与无越界声明、处置 R1～R3 移交的核账项（GOAL-002 F-001/F-002、GOAL-003 F-001/F-002、GOAL-004 F-002/F-005/F-006/F-007），随后执行 **Root 关门审计（self + 本地 grok build independent）**，经用户书面确认后 Root `GOAL-001` → `done 4/4`、工作区结项、VP-020 关门记录填写。

**审计模式**：`independent`（Root 关门 = 证据/无越界/required=0 门禁；按项目规则自审后调用本地 grok build（grok-4.6 · high）执行 `source: independent` 独立审）。

## 检查点（progress 派生源）

| # | 检查点 | 状态 |
|---|--------|------|
| C1 | 证据矩阵落盘：快测汇总（时区 15 + money 24 + runtime 接线 7 + switcher 4 + settings 保存/预填）映射到合同条款；双 locale 范例（zh-CN/en-US 各展示+输入场景） | **done**（2026-08-26 · `attachments/r4-evidence-matrix.md`） |
| C2 | 无越界核账：汇率/换算/计费、DB `timestamptz`（RT-T03）、Profile 默认集/模块矩阵/Manifest、`docs/contracts/`、热加载等逐项核对；API 机器合同不变量复证 | **done**（2026-08-26 · 矩阵 §5 逐项） |
| C3 | 核账项处置：GOAL-002 F-001/F-002（closed）；GOAL-003 F-001/F-002（epoch 输入控件·TIMEZONE_OPTIONS 留痕）；GOAL-004 F-002/F-005/F-006/F-007（grouping 位序评估 / FAIL 币种目录 / 安全整数评估→加拒） | **done**（2026-08-26 · F-007 **fixed**（safe-integer 守卫）；F-002/F-005/F-006 final residual；其余 closed/accepted） |
| C4 | Root 关门：自审 A-001（self）+ grok build independent（Root 03-audit）→ 用户书面确认 → Root done 4/4 + goal-tree/workspace/VP-020 关门记录同步 | **done**（Root A-001 self pass → A-002 grok independent pass（0 required）→ 用户 2026-08-27 书面确认关门；Root done 4/4，VP-020 收尾同步） |

`progress` = 4/4（2026-08-27 关门：用户书面确认；Root 关门审计双腿 pass）。

## 成功标准

1. 证据矩阵可核对：每个合同条款（§2/§3/§4）至少一个快测或双 locale 范例引用；R2/R3 已全量回归（Go 全量 + web 1180）。
2. 双 locale 范例：`zh-CN` / `en-US` 各至少一场景（时区展示 + 货币展示 + 输入解析），与快测证据同源。
3. 无越界声明逐项核对成立（汇率/计费/RT-T03/Profile 默认集/模块矩阵/Manifest/`docs/contracts/`/热加载）；API 机器合同不变量复证。
4. 核账项有明确处置或书面 residual（延续用户 2026-08-26 已接受的 F-005/F-006/F-007 范围；F-007 若评估为「加拒」则实施）。
5. Root 关门审计（self + independent）开放 required = 0；用户书面确认后 Root done、工作区结项、VP-020 关门记录落盘。

## 信息就绪与未知项

> 无新增 required 信息项。I-001/I-002/I-005 已 verified；I-003/I-004 registered（VP 冻结不进）。核账项为 recommended/residual（范围已书面明确）。

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| — | — | 无新增开放信息项 | — | — | — | — | — | 引用 Root D-002 / 各目标 audit 闭环 |

## 父目标

- `GOAL-001-timezone-number-currency-formatting`（Root；VP-020 `active` · primary_plan）

## 台账布局

新目标为三个可追加台账创建同名平铺目录：`01-decision/`、`02-execution/`、`03-audit/`。索引文件保留 frontmatter、摘要和条目索引；独立记录使用 `D-NNN-*`、`E-NNN-*`、`A-NNN-*` 文件。