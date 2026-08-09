---
id: GOAL-006-s5-evidence-and-closeout
title: S5 · 双 Profile 验证矩阵与关门
status: active
parent: GOAL-001-localization-and-system-settings
created: 2026-08-09
updated: 2026-08-09
version: 0.1.0
progress: 0/4
---

# GOAL-006 · S5 · 双 Profile 验证矩阵与关门

## 概述

承接 Root [GOAL-001-localization-and-system-settings](../GOAL-001-localization-and-system-settings/00-meta.md) 的 **S5 阶段**：完成两语种 × 两 Profile × 匿名/已认证证据矩阵（复用 F-V029 同一分母）、真实入口启动验证（API 二进制 + Web 构建/浏览器）、Root 关门独立审计（grok CLI，`-m grok-4.5 --effort high`）、用户书面关门确认后置 Root `done`、填 VP-007 关门记录、最终 goal-tree 同步与 checkpoint commit。

**方案依据**：Root D-002（审计模式：S5 关门 = `independent`）+ 计划验证步骤 4（真实入口两次启动一致）+ P-004（关门需用户书面确认）。本目标只验证与收口，不重新决策。

**范围纪律**：不新增业务功能；不改协议语义；不重开已关闭阶段。

## 成功标准（可验收 · 等权检查点 · 共 4 项）

- [ ] **C1**：证据矩阵完成：两语种 × 两 Profile × 匿名/认证分母，列覆盖固定 UI / 冻结 pageId-schema 并集 / M1～M4 / 权限正反例 / 缺失翻译 / 配置刷新 / 错误回退；非 N/A 单元格有证据路径，N/A 仅限 Profile 不可达并注明模块边界（落 Root `attachments/`）。
- [ ] **C2**：真实入口启动验证：API 二进制 `go build ./cmd/server` 后启动，请求公开启动配置（`/api/branding`）断言响应体内容（locale/品牌字段）正确，同一启动重复 ≥2 次一致成功；Web `npm run build` 成功；playwright 可用则 serve + 加载断言零页面错误、`lang` 随切换生效、一次设置保存产生可见变化并截图；不可用则捕获失败输出并以静态/结构性回退为验收线。输出捕获 `{SCRATCH}`。
- [ ] **C3**：关门独立审计：S5 关门审计由 grok CLI（`-m grok-4.5 --effort high`，`/audit` 提示词）执行并落盘 `03-audit/A-NNN-*`（source: independent）；required findings 全闭合后方可放行关门。
- [ ] **C4**：用户书面关门确认（P-004 留痕，含日期与范围）→ Root `status: done`、`progress: 6/6`；VP-007 关门记录填写（outcome/summary/evidence_links/residuals）；goal-tree 最终同步；checkpoint commit。

## 派生进度展示

`progress: X/4` 由上方 4 个等权检查点派生；仅为展示，不放行阶段、不关闭 finding。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 关门所需输入（矩阵分母、真实入口验证、审计、用户确认）是否齐备 | C1–C4 | 关门 | 逐项执行本目标检查点 | **closed** | — | S0–S4 证据链齐备（2026-08-09） |

## 父目标

- [GOAL-001-localization-and-system-settings](../GOAL-001-localization-and-system-settings/00-meta.md)（Root；本目标为 S5 阶段子目标）

## 台账布局

本目标使用 ledger 目录：`01-decision/`、`02-execution/`、`03-audit/`。索引文件保留 frontmatter 与条目表；独立记录使用 `D-NNN-*`、`E-NNN-*`、`A-NNN-*`。
