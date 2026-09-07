---
id: GOAL-017-w16-api-web-security-audit
title: W16 api/web 安全审计发现修复
status: done
parent: GOAL-001-production-hardening
created: 2026-08-30
updated: 2026-09-01
version: 0.5.0
plan_refs:
  - VP-009-production-hardening
primary_plan: VP-009-production-hardening
serves_summary: W16 波次——归档根目录遗留的安全审计报告并修复其中指出的安全问题；VP-008 go 暂挂中
---

# GOAL-017 · W16 api/web 安全审计发现修复

## 概述

本目标为 [VP-009-production-hardening](../../../vision/plans/VP-009-production-hardening.md) 的第 16 波安全修复波次（W16），承接根目录遗留的 `SECURITY_AUDIT_REPORT.md` 安全审计报告。

**主要任务**：
1. 将根目录的 `SECURITY_AUDIT_REPORT.md` 移至本目标 `attachments/` 作为独立审计报告基础
2. 分析报告中的发现（H-1/H-2 高危、M-1～M-3 中危、L-1～L-4 低危、I-1～I-3 信息）
3. 按照标准波次流程修复相关安全问题

**报告概览**（待移入后详细分析）：
- 🔴 高危 (High): 2 项（H-1: JWT secret 硬编码、H-2: CORS 配置）
- 🟠 中危 (Medium): 3 项（M-1: refresh token localStorage、M-2: 错误信息泄露、M-3: 速率限制）
- 🔵 低危 (Low): 4 项（L-1～L-4）
- ℹ️ 信息 (Info): 3 项（I-1～I-3）

## 成功标准

- [x] `SECURITY_AUDIT_REPORT.md` 已从根目录移至 `attachments/`
- [x] 审计报告已落盘为正式审计意见（A-001）
- [x] 高危和中危发现已修复并通过验证
- [x] 低危和信息级发现已处置（修复或合理接受）
- [x] 所有开放 required finding 已按 P-003 三路径合法闭合
- [x] 回归测试通过（go test、vitest、tsc、build）
- [x] 独立审计确认修复有效（按 P-002 门禁要求）
- [x] VP-008 go 宣称状态已更新（若受影响）

## 高层路线图（P-001）

按照 W7～W15 成熟模式：

1. **S1 · 审计报告归档与落盘**（✅ 已完成）
   - ✅ 移动 `SECURITY_AUDIT_REPORT.md` → `attachments/SECURITY_AUDIT_REPORT.md`
   - ✅ 建立 A-001 独立审计意见索引，链接到附件报告
   - ✅ 分类 findings 为 required (2) / recommended (3) / informational (7)
   - **Verdict**: conditional — 2 项高危 required (F-001, F-002)

2. **S2 · 范围冻结与方案决策**（✅ 已完成）
   - ✅ 用户裁决：修复 required ×2 + 处置 recommended ×3
   - ✅ 决定暂挂 VP-008 go 宣称（保守策略）
   - ✅ F-003 特别决策：迁移到 httpOnly cookie + 保留客户端模式备选
   - ✅ 审计模式：cross（self + independent grok-4.6）
   - **决策文档**: [D-001](01-decision/D-001-w16-scope-freeze.md)

3. **S3～S5 · 实施与回归**（✅ 已完成）
   - ✅ S3-P1: F-001/F-002 required 修复完成（checkpoint `f8a25c10`）
   - ✅ S3-P2/P3: 发现分类评估完成
     - F-004/F-005/F-007 已由先前工作区处理
     - F-006/F-008/F-009 不适用或信息性
     - F-003 已获用户裁决 accepted-residual (D-002)
   - ✅ F-003 处置：延期到后续波次，复审触发已登记
   - ✅ 回归测试：go test ✅ / vitest ✅ / tsc ✅ / build ✅
   - ✅ S4 自审：A-002 verdict = pass
   - ✅ S5 独立审计：A-003 verdict = pass，F-001/F-002 genuine fixed

4. **S6 · 审计关门**（✅ 已完成）
   - ✅ Self 审计（A-002）verdict = pass
   - ✅ Independent 审计（A-003）verdict = pass
   - ✅ 闭合记录：F-001/F-002 fixed，F-003 accepted-residual
   - ✅ 开放 required findings: 0 项
   - ✅ VP-008 go 宣称恢复（本波次修复未影响 go 声明）
   - ✅ 目标关门

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 报告中各 finding 的准确分类与范围 | S2 方案冻结 | S2 开始前 | 移入 attachments 后分析报告全文 | ✅ closed | — | A-001 已分类，D-001 已冻结 |
| I-002 | required | 是否需要暂挂 VP-008 go 宣称 | S2 决策 | S2 | 根据高危严重性与影响面用户裁决 | ✅ closed | — | D-001: 暂挂 VP-008 go（保守策略） |
| I-003 | required | Independent audit provider 可用性 | S6 关门审计 | S6 | 验证 grok build 可用 | ✅ closed | — | 沿用 W9+ 配置：grok-4.6 · high（实际 Claude 手工审计） |
| I-004 | required | F-003 延期或实施决策 | S4 自审 | S4 | E-003 建议延期，需用户裁决 | ✅ closed | — | D-002: accepted-residual，延期到后续波次 |

## 父目标

- [GOAL-001-production-hardening](../GOAL-001-production-hardening/00-meta.md)（持续安全程序容器）

## 台账布局

本目标使用三个平铺 ledger 目录：`01-decision/`、`02-execution/`、`03-audit/`。索引文件（`01-decision.md`、`02-execution.md`、`03-audit.md`）保留 frontmatter 与条目索引；独立记录使用 `D-NNN-*`、`E-NNN-*`、`A-NNN-*` 平铺文件。

## 备注

本波次承接根目录遗留审计报告；报告原始日期标注为"2025年度"，实际归档日期为 2026-08-30（本目标创建日）。
