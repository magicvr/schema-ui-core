---
id: GOAL-002-r1-contract-freeze
title: R1 · 契约冻结面落盘
status: done
parent: GOAL-001-distribution-package-pilot
created: 2026-08-29
updated: 2026-08-29
version: 0.2.0
progress: 4/4
---

# GOAL-002 · R1 · 契约冻结面落盘

## 概述

承接 Root R1 检查点：把「契约冻结面」从概念落成可核对清单——内核公共 API（A 层）、模块契约装配面（B 层）、内部实现面（C 层）三分层 + semver/breaking 流程 + changelog 模板。「冻结面 vs 内部自由演进面」分界成文。产出经用户确认后成为后续 R2～R5 的契约基座。

## 成功标准（阶段检查点）

- [x] **S1 · 扫描与草案**：`apps/api` Go 面扫描——kernel 导出面全量、模块契约模式、组合根装配面（E-001；冻结面清单 v0.1.0 草案）
- [x] **S2 · 成文**：semver/breaking 流程 + changelog 模板 + 分界规则（attachments ×2 + 清单 §0）
- [x] **S3 · 对账验证**：清单 vs 实际代码核对（A-001 self 审计：文件计数 / 导出锚点抽样 / 模块模式抽样）
- [x] **S4 · 关门**：**用户确认**冻结面清单 v0.1.0 → **v1.0.0 生效**（2026-08-29 · D-002；关键决策，P-004）

## 信息就绪

| ID | 级别 | 所需信息 | 影响门禁 | 最晚阶段 | 状态 | 证据 |
|----|------|----------|----------|----------|------|------|
| I-001 | required | kernel 公共 API 冻结面清单 | R1 冻结 / R2/R3 实施 | R1 | **collecting**（S1 已扫描，草案落盘；B 层逐包符号清单随 R2 收尾） | `attachments/freeze-face-v1.2.0.md` |

## 父目标

- `GOAL-001-distribution-package-pilot`

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；条目见各索引。