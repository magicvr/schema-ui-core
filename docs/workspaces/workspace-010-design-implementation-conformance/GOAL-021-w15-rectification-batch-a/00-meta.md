---
id: GOAL-021-w15-rectification-batch-a
title: W15 整改批 A · 核心安全与高危体验（F01 会话容灾 / F02 表格重试 / F04 JSON 404-405 / F05 CORS 安全头 / F07 refresh 错误码）
status: done
parent: GOAL-020-w15-user-perspective-findings
created: 2026-08-17
updated: 2026-08-17
version: 0.2.0
progress: 4/4
---

# GOAL-021 · W15 整改批 A

[GOAL-020](../GOAL-020-w15-user-perspective-findings/00-meta.md) 下级整改子目标（批 A）。承接 D-002：F01～F14 全部 in-scope。

## 当前边界

- **本波实施**：W15-F01、F02、F04、F05、F07。
- **非范围**：批 B/C；不改 Profile 默认集 / 模块矩阵 / Manifest。

## 成功标准与路线图

- [x] **S1 · 方案冻结**：D-001 冻结五条 as-built 改法
- [x] **S2 · 实施**：代码与测试
- [x] **S3 · 回归**：Go 全量 + Web 1046/1046 各两遍
- [x] **S4 · 自审与关门**：A-001 响应 + A-003 self pass

progress: **4/4**（2026-08-17 关门）。

## 审计策略

| 阶段 | 模式 | 说明 |
|------|------|------|
| S1 | self | 冻结对照 D-002 |
| S2 F01/F04/F05/F07 | independent | 会话/信封/CORS/错误码触及 security + compatibility |
| S4 | self | 关门自审 |

## 信息就绪

| ID | 级别 | 状态 | 结论 |
|----|------|------|------|
| I-001 | non-blocking | **closed** | 范围来自 GOAL-020 D-002 |

## 父目标

- [GOAL-020-w15-user-perspective-findings](../GOAL-020-w15-user-perspective-findings/00-meta.md)
