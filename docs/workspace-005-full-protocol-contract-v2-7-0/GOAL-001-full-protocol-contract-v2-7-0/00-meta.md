---
id: GOAL-001-full-protocol-contract-v2-7-0
title: schema-ui-docs@v2.7.0 整份契约可验证兼容
status: active
parent: null
created: 2026-08-08
updated: 2026-08-08
version: 0.1.0
progress: 5/6
plan_refs:
  - VP-006-full-protocol-contract-v2-7-0
primary_plan: VP-006-full-protocol-contract-v2-7-0
serves_summary: 对 schema-ui-docs@v2.7.0 pin 形成整份契约可验证兼容（覆盖表 I-PROTO-FULL-001、Renderer/后端实现、范例与验证）；纠正长期停留在 I-PROTO-001 v0.1.3 MVP 子集的组合焦点；不启动 VP-005 视觉产品化。
---

# GOAL-001 · schema-ui-docs@v2.7.0 整份契约可验证兼容

## 概述

本 Root 承接 [VP-006 · 整份 v2.7.0 契约可验证兼容](../../vision/plans/VP-006-full-protocol-contract-v2-7-0.md)。在已关闭的 [VP-003](../../vision/plans/VP-003-modular-admin-architecture.md) 单主线架构与 [VP-004](../../vision/plans/VP-004-module-contribution-readiness.md) 贡献契约之上，交付：

1. 全量覆盖表 **`I-PROTO-FULL-001`**（新文件 + 新决策；默认 disposition = `include`）；
2. Renderer / 协议 runtime 与后端、模块贡献面对齐纳入面；
3. 每一纳入能力的范例路径 + 可复核验证；
4. 回归诚实：不回退 `I-PROTO-001 v0.1.3` 主路径，不在覆盖未闭合时宣称「已完整支持 v2.7.0」。

权威输入：[protocol-inventory-v2.7.0.md](../../vision/protocol-inventory-v2.7.0.md) + 上游 pin `ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b`。  
历史 MVP 基线（只读）：[I-PROTO-001 v0.1.3](../../workspace-001-mvp-admin-foundation/GOAL-001-mvp-admin-foundation/attachments/I-PROTO-001-coverage-draft.md)。

## 愿景对齐

| 字段 | 值 |
|------|----|
| `plan_refs` | `VP-006-full-protocol-contract-v2-7-0` |
| `primary_plan` | `VP-006-full-protocol-contract-v2-7-0` |
| Charter | `schema-ui-core-admin-foundation@0.2.0` |
| 工作区角色 | `delivery`（lead of VP-006） |
| 工作区 | `workspace-005-full-protocol-contract-v2-7-0` |

## 成功边界

### 阶段层（可验收 · 等权检查点）

- [x] **S0**：差距盘点 — 覆盖表 v0.1.3 vs inventory/registry/fixtures 差集；前端保真债 vs 未纳入 type 分列；产出可审计差集清单。
- [x] **S1**：覆盖表升版冻结 — 落盘 `I-PROTO-FULL-001` + Root 决策；默认 `include`；`include-partial` 仅保真/边角；范围收缩 → exclude 或用户书面 residual；相对 v0.1.3 差集「转为 include 计数 / 仍 residual 清单」。独立审计 A-001（source: independent，grok build / grok 4.5 / high）→ conditional；F-001 required 已 fixed，F-002/F-003 已勘误。
- [x] **S2**：核心缺口实现 — 未实现 registry type / 批量 selection / upload 等按表纳入批次交付（B1–B4；320/320 fixture case 全绿）。
- [x] **S3**：保真与 runtime — 钉死内降级控件提升到契约语义；表达式/权限边角 fail-closed（B5；白名单/门禁/空选/循环阻断测试）。
- [x] **S4**：范例 + conformance — 每纳入域可发现范例与验证入口；exclude 面有表可查（覆盖表验证入口列已登记真实路径）。
- [ ] **S5**：文档与关门 — 发现路径（overview/QUICKSTART）、兼容声明诚实、回归不回退、close-out 审计；开放 required = 0；可向 `/vision` 提出 VP-006 关门提案（须用户确认）。

### 阶段 ↔ VP 退出判据映射

| 阶段 | 主要服务的 VP 退出判据 | 证据（激活后填写） |
|------|------------------------|--------------------|
| S0 | 为 exit 1–4 提供差集输入 | 待 S0 产物 |
| S1 | exit 1 覆盖表决策 | 待 `I-PROTO-FULL-001` + D-* |
| S2 | exit 2–3 实现纳入面 | 待实现证据 |
| S3 | exit 2 保真 / fail-closed | 待实现证据 |
| S4 | exit 4 范例与验证 | 待 fixture/集成登记 |
| S5 | exit 5–6 文档诚实 + 过程可关门 | 待审计 + 关门提案 |

## 纲领路线图

| 阶段 | 名称 | 状态 | 说明 |
|------|------|------|------|
| S0 | 差距盘点 | **已完成** | 差集证据 `attachments/I-S0-001-*`；E-002；I-001 closed |
| S1 | 覆盖表升版冻结 | **已完成** | `I-PROTO-FULL-001` v1.0.0 + D-002；独立审计 A-001 conditional → F-001 fixed / F-002·F-003 勘误；I-PROTO-FULL-001 closed |
| S2 | 核心缺口实现 | **已完成** | B1–B4（E-003）；320/320 fixture case 全绿；vitest 569 / go test 全绿 |
| S3 | 保真与 runtime | **已完成** | B5 fail-closed 测试；表达式/权限/批量/上传边界（E-003） |
| S4 | 范例 + conformance | **已完成** | 覆盖表验证入口列登记真实路径（8 范例页 + 每域验证入口） |
| S5 | 文档与关门 | **进行中** | 依赖回归与关门审计；VP 关门提案须用户确认 |

阶段通常串行；同一纲领阶段内允许并行子目标。建区 **不**勾选任何检查点。

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-PROTO-FULL-001 | required | 整份 v2.7.0 契约覆盖 disposition 与冻结（新文件 + 新版本号） | S1 冻结；S2–S5 范围分母 | S1 方案冻结前 | S0 差集 → Root 决策落盘 attachments | **closed** | — | `attachments/I-PROTO-FULL-001-coverage-v2-7-0.md`（v1.0.0）+ D-002（2026-08-08）；独立审计 A-001 复核通过（F-001 fixed） |
| I-001 | required | v0.1.3 vs inventory/registry/fixtures 可审计差集 | S0 完成；S1 输入 | S0 结束前 | 盘点 inventory + 覆盖表 + 代码现状 | **closed** | — | 证据：`02-execution/E-002-s0-gap-analysis.md` + `attachments/I-S0-001-gap-analysis-v0-1-3-to-full.md`（2026-08-08） |
| I-002 | required | 范围收缩 / exclude 是否获用户书面 residual（若有） | S1 冻结 | S1 决策时 | P-004 用户裁决 + 决策留痕 | **N/A** | 差集全部可纳入，无收缩 | S0 结论：12/12 域、24/24 registry type、16/16 行为 fixture 套件默认 include；无 exclude / 范围收缩（E-002 §4） |
| I-003 | non-blocking | 上游 schemas/fixtures 是否 vendor 到本仓 | 验证策略 | S2 前可决 | 策略决策 | **closed** | 继承历史（I-PROTO-004 = vendor） | 6 schema + 17 fixture 套件全部 vendor + SHA pin（`provenance.json`，含 uploads `aaeb9683…` / permissions-inheritance `ac124fa1…`） |
| I-004 | non-blocking | 批次切分与并行子目标边界 | S2 立项 | S2 方案 | Root 决策 | **closed** | — | D-002 §4 批次 B1–B6 |

## 非目标（本 Root）

- **不**启动或吸收 VP-005 视觉/设计系统/Shell 产品化（VP-006 closed 前硬禁止）。
- **不**交付订单、钱包、类目、通知等业务领域模块。
- **不**重开 VP-003 架构迁移、不恢复长期双线、不引入运行时插件/热插拔。
- **不**在本项目内重新定义或替代上游协议语义。
- **不**修订 Goal Governance 核心方法论（principles P-001～P-006 等）。
- **不**把历史 `I-PROTO-001 v0.1.3` 或已关闭 VP 的「子集完成」改写成「全量协议已完成」。
- 不为 VP 在 `docs/vision/` 建立 Goal 五件套或 progress% 权威。

## 派生进度展示

`progress: 5/6` 由上方 S0～S5 六个等权检查点派生（S0–S4 已完成）。progress 仅为展示；不放行阶段、不关闭 finding、不覆盖信息门禁，也不自动推导 `status: done`。

## 台账布局

本目标使用 ledger 目录：`01-decision/`、`02-execution/`、`03-audit/`。索引文件保留 frontmatter 与条目表；独立记录使用 `D-NNN-*`、`E-NNN-*`、`A-NNN-*`。
