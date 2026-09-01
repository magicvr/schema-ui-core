---
id: GOAL-017-w16-api-web-security-audit
title: W16 api/web 安全审计发现修复
status: active
parent: GOAL-001-production-hardening
created: 2026-08-30
updated: 2026-08-30
version: 0.2.0
plan_refs:
  - VP-009-production-hardening
primary_plan: VP-009-production-hardening
serves_summary: W16 波次——归档根目录遗留的安全审计报告并修复其中指出的安全问题
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
- [ ] 高危和中危发现已修复并通过验证
- [ ] 低危和信息级发现已处置（修复或合理接受）
- [ ] 所有开放 required finding 已按 P-003 三路径合法闭合
- [ ] 回归测试通过（go test、vitest、tsc、build）
- [ ] 独立审计确认修复有效（按 P-002 门禁要求）
- [ ] VP-008 go 宣称状态已更新（若受影响）

## 高层路线图（P-001）

按照 W7～W15 成熟模式：

1. **S1 · 审计报告归档与落盘**（✅ 已完成）
   - ✅ 移动 `SECURITY_AUDIT_REPORT.md` → `attachments/SECURITY_AUDIT_REPORT.md`
   - ✅ 建立 A-001 独立审计意见索引，链接到附件报告
   - ✅ 分类 findings 为 required (2) / recommended (3) / informational (7)
   - **Verdict**: conditional — 2 项高危 required (F-001, F-002)

2. **S2 · 范围冻结与方案决策**
   - 用户裁决：哪些必修、哪些可接受残余、哪些降级
   - 决定是否需要暂挂 VP-008 go 宣称
   - 冻结修复范围（D-001 或 D-002）

3. **S3～S5 · 实施与回归**
   - 按优先级修复 required findings
   - 处置 recommended findings（修复 / overruled / deferred）
   - 回归测试：go vet/test、vitest、tsc、vite build
   - Git checkpoint（按长流程策略）

4. **S6 · 审计关门**
   - Self 审计（A-002）
   - Independent 审计（A-003，grok build · grok-4.6 · high）
   - 闭合记录：所有 required finding 已 fixed / accepted-residual / user-overruled
   - 恢复 VP-008 go 宣称（若曾暂挂）
   - 用户书面关门确认

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 报告中各 finding 的准确分类与范围 | S2 方案冻结 | S2 开始前 | 移入 attachments 后分析报告全文 | open | — | 待报告归档后确认 |
| I-002 | required | 是否需要暂挂 VP-008 go 宣称 | S2 决策 | S2 | 根据高危严重性与影响面用户裁决 | open | — | 待 S2 用户确认 |
| I-003 | required | Independent audit provider 可用性 | S6 关门审计 | S6 | 验证 grok build 可用 | open | — | 默认沿用 workspace-008 D-002 配置 |

## 父目标

- [GOAL-001-production-hardening](../GOAL-001-production-hardening/00-meta.md)（持续安全程序容器）

## 台账布局

本目标使用三个平铺 ledger 目录：`01-decision/`、`02-execution/`、`03-audit/`。索引文件（`01-decision.md`、`02-execution.md`、`03-audit.md`）保留 frontmatter 与条目索引；独立记录使用 `D-NNN-*`、`E-NNN-*`、`A-NNN-*` 平铺文件。

## 备注

本波次承接根目录遗留审计报告；报告原始日期标注为"2025年度"，实际归档日期为 2026-08-30（本目标创建日）。
