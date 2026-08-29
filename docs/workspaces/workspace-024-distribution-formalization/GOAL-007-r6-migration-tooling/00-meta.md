---
id: GOAL-007-r6-migration-tooling
title: R6 · fork→包迁移工具化（schema-ui migrate-fork · 非破坏性辅助）
status: active
parent: GOAL-001-distribution-formalization
created: 2026-08-29
updated: 2026-08-29
version: 0.1.0
progress: 0/4
---

# GOAL-007 · R6 · fork→包迁移工具化

## 概述

承接 Root R6 与 VP-024 判据 #7：把迁移指南（workspace-023 fork-to-package-migration-guide）配套为**可执行工具**——`schema-ui migrate-fork` 子命令（非破坏性辅助：检查/引导/低侵入改写），面向 A/B 型 fork（C 型深度定制 = 建议保持 fork）。核销 go 后清单「fork→包迁移工具化」。

## 成功标准（可验证检查点）

- [ ] C1：`schema-ui migrate-fork [--dir <path>] [--dry-run]` 子命令存在：dry-run 输出类型判定（A/B/C）+ 迁移步骤清单（不写文件）
- [ ] C2：实跑对 A/B 型旧态：`go.mod` require bump 至 @latest（registry 语义）· `web/.npmrc` GH 映射 → npmjs 钉死（备份 `.npmrc.migrate.bak`）· 旧组合根 main 检测 → 引导（**不覆盖**用户代码）· 报告输出
- [ ] C3：golden-field `9510023`（v0.3.0 旧态）实测：migrate-fork dry-run → 实跑 → `go build ./cmd/server` 绿 + 报告含验证建议
- [ ] C4：独立审计（grok）→ 关门（Root 6/7）

## 方案与路线（P-001）

| 阶段 | 内容 | 状态 |
|------|------|------|
| S1 | 工具设计定档（D-001）：检查面/改写规则/非破坏原则 | 未开 |
| S2 | 实现 `schema-ui migrate-fork`（CLI 子命令 + dry-run） | 依赖 S1 |
| S3 | golden-field 旧态实测（dry-run + 实跑 + build 验证） | 依赖 S2 |
| S4 | 证据 + 独立审计（grok）→ 关门 | 依赖 S2/S3 |

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 迁移判定特征（A/B/C 的程序化识别面） | C1 | S2 | 读迁移指南 §1 + 代码特征扫描 | open | S2 前闭合 | 待确认（设计期闭合） |

## 父目标

- `GOAL-001-distribution-formalization`

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger。

## 备注

- 审计模式：S4 关门 = independent（grok build 先例）。
- 非破坏原则：除 `.npmrc`（带备份）与 go.mod 依赖行外不做任何修改；用户代码一律引导不覆盖。