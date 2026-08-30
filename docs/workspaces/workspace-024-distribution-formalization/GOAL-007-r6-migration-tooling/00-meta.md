---
id: GOAL-007-r6-migration-tooling
title: R6 · fork→包迁移工具化（schema-ui migrate-fork · 非破坏性辅助）
status: done
parent: GOAL-001-distribution-formalization
created: 2026-08-29
updated: 2026-08-29
version: 0.4.0
progress: 4/4
---

# GOAL-007 · R6 · fork→包迁移工具化

## 概述

承接 Root R6 与 VP-024 判据 #7：把迁移指南（workspace-023 fork-to-package-migration-guide）配套为**可执行工具**——`schema-ui migrate-fork` 子命令（非破坏性辅助：检查/引导/低侵入改写），面向 A/B 型 fork（C 型深度定制 = 建议保持 fork）。核销 go 后清单「fork→包迁移工具化」。

## 成功标准（可验证检查点）

- [x] C1：`schema-ui migrate-fork [--dir <path>] [--dry-run]` 子命令存在：dry-run 输出类型判定（A/B/C）+ 步骤清单（哈希不变 · 零写入）
- [x] C2：实跑 A/B 型旧态：go.mod require bump 至 registry @latest（无条件执行）· .npmrc GH→npmjs（备份 `.npmrc.migrate.bak`）· 旧组合根引导（不覆盖）· 报告含验证建议
- [x] C3：9510023（v0.3.0 旧态）：dry-run（B 型）→ 实跑（v0.3.0→v0.4.0 · npmjs 钉死）→ `go mod tidy` + `go build ./cmd/server` exit 0；块形态/无依赖夹具通过（F-002）
- [x] C4：独立审计（grok）→ A-002 `pass`（0 required · F-001~F-003 fixed）→ **关门（Root 6/7）**

## 方案与路线（P-001）

| 阶段 | 内容 | 状态 |
|------|------|------|
| S1 | 工具设计定档（D-001）：检查面/改写规则/非破坏原则 | **已关门**（2026-08-29 · D-001 · I-001 verified） |
| S2 | 实现 `schema-ui migrate-fork`（CLI 子命令 + dry-run） | **已关门**（2026-08-29 · E-002 · F-002 块解析补丁） |
| S3 | golden-field 旧态实测（dry-run B 型 + 实跑 + build 绿） | **已关门**（2026-08-29 · E-002） |
| S4 | 证据 + 独立审计（grok）→ 关门 | **已关门**（2026-08-29 · A-002 `pass`（0 required）· F-001~F-003 fixed · Root 6/7） |

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 迁移判定特征（A/B/C 的程序化识别面） | C1 | S2 | 读迁移指南 §1 + 代码特征扫描 | **verified**（A=薄封装 · B=组合根无覆盖 · C=手搓 kernel 无组合路径；9510023=B 实测校准 · A/C 夹具复现） | — | D-001 §7 · A-002 |

## 父目标

- `GOAL-001-distribution-formalization`

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger。

## 备注

- 审计模式：S4 关门 = independent（grok build 先例）。
- 非破坏原则：除 `.npmrc`（带备份）与 go.mod 依赖行外不做任何修改；用户代码一律引导不覆盖。