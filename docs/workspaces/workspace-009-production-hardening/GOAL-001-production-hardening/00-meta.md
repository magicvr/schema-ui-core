---
id: GOAL-001-production-hardening
title: 生产加固（共享基架持续安全与健壮性）
status: active
parent: null
created: 2026-08-10
updated: 2026-08-21
version: 0.9.0
plan_refs:
  - VP-009-production-hardening
primary_plan: VP-009-production-hardening
serves_summary: 长期安全程序容器——周期扫描、波次修复、与 VP-008 go 消费有效性接口；波次=子目标，Root 不因单波完成而 done
---

# GOAL-001 · 生产加固（共享基架持续安全与健壮性）

## 概述

本 Root 是 `workspace-009-production-hardening` 的唯一总目标，承接 [VP-009-production-hardening](../../../vision/plans/VP-009-production-hardening.md) 的**长期实现层容器**。

- **程序语义**：持续安全扫描、finding 分流、有界波次修复、回归与（必要时）VP-008 `go` 重验证。  
- **波次语义**：每一次扫描/审查发现的修复 = **一个子目标**（可 `done`）；**Root 默认保持 `active`**，不因单波完成而关门。  
- **历史纠正**：2026-08-10 曾将「W1 修完 16 项」或「W2 修完上传 IDOR」当作 Root/`VP` 关门，属把有界波次误认为程序结束；用户同日书面纠正为长期意图。

本 Root **不**重开 VP-001～008 的历史 status，**不**修改 Charter 目的/边界/非目标，**不**实现订单/钱包/类目/通知等业务模块。

## 愿景对齐

| 字段 | 值 |
|------|-----|
| Charter | `schema-ui-core-admin-foundation@0.2.0` |
| `plan_refs` / `primary_plan` | `VP-009-production-hardening`（`active` 长期程序） |
| 工作区 | `workspace-009-production-hardening` (`vision_role: delivery`, VP lead) |
| independent provider | **`grok build` · 模型 `grok-4.5` · 思考强度 high · 执行 `audit` 命令**（沿用 workspace-008 D-002） |
| 审计模式 | 波次含 security 高影响时默认 `cross`（self + independent）；低风险波次可 `self`（P-004） |

## 成功标准（程序能力 · 非「修完即 done」）

下列检查点表示**程序已成立**；全部勾选后 Root **仍保持 `active`**，等待下一波扫描或 finding。

- [x] **P1 · 程序与波次模型**：Root = 长期容器；波次 = 子目标；单波完成 ≠ Root/VP 关门。（2026-08-10 用户纠正 + D-003）
- [x] **P2 · 与 go 的接口**：共享基架 Critical/High 可触发 VP-008 `go` 消费有效性暂挂/恢复的路径有台账约定。（见 VP-009；W1 曾完成一次恢复证据）
- [x] **W1 · 波次档案**：2026-08-10 审查 16 项（C1–C8 + D1–D8）修复 + cross 闭环 — [GOAL-002](../GOAL-002-audit-findings-remediation/00-meta.md) `done` 16/16
- [x] **W2 · 波次档案**：上传 owner/下载鉴权 + `ReadHeaderTimeout` — [GOAL-003](../GOAL-003-upload-ownership-hardening/00-meta.md) `done` 4/4
- [x] **P3 · 下一波就绪**：存在约定的触发（例行/发版前/变更后）时，可开新子目标承接扫描，无需重开 Root/VP。（W3–W10 已按此开波；最近：2026-08-21 W10）

> `progress`：不使用「n/n → Root done」推导。波次完成只更新子目标与下表档案；Root `status` 仅在用户明确废弃程序或迁移 `primary_plan` 时改为 `done`/`cancelled`。

## 波次台账（摘要）

| 波次 | 子目标 | status | 说明 |
|------|--------|--------|------|
| W1 | GOAL-002-audit-findings-remediation | done | 首批审查 16 项 |
| W2 | GOAL-003-upload-ownership-hardening | done | 上传 IDOR + ReadHeaderTimeout |
| W3 | GOAL-004-w3-security-audit-remediation | done | 2026-08-11 安全审计修复 |
| W4 | GOAL-005-w4-security-audit-remediation | done | 2026-08-11 限流/上传/token_version/前端异常 |
| W5 scan | — | 未开子目标 | 2026-08-14 全量审计 0 中高危；低危就地修补见 E-002 |
| W6 | GOAL-006-w6-scan-findings-remediation | done | 低危扫描修补；2026-08-17 关门 |
| W7 | GOAL-007-w7-api-web-security-audit | done | 2026-08-19 独立审计 A-001 fail → 整单采纳 → A-002/A-004 + A-005 pass；VP-008 go 宣称恢复（D-003） |
| W8 | GOAL-008-w8-api-web-security-audit | done | 2026-08-20 独立审计 A-001 fail → D-002 整单采纳 + go 暂挂 → E-002 修复 → A-002/A-003 pass → D-003 恢复 go；真实浏览器/CSP 回归 E-004 通过 |
| W9 | GOAL-009-w9-api-web-security-audit | done | 2026-08-21 独立审计 A-001 fail → A-002 conditional → D-003 整单 12 条 required → E-004 实施 12/12 → A-005 independent pass → A-006 闭合 → D-004 恢复 go |
| W10 | GOAL-010-w10-api-web-security-audit | active | 2026-08-21 独立审计 A-001 conditional（1 HIGH required）→ D-002 整单 7 条 + go 暂挂 → D-003 调和 4 误报作废 → E-002 修复 3 条 + 回归全绿 + A-002 self pass（开放 required = 0）；S4/关门/go 恢复待用户裁决 |

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 本程序是否为长期意图（非单波关门）？ | 程序定义 | 纠正当日 | 用户 2026-08-10 书面纠正 | verified | — | D-003；VP-009 v0.4.0 |
| I-002 | non-blocking | 例行扫描的具体日历/cron | 运营节奏 | 下一波前 | 用户或 CI 约定；可先事件触发 | open | deferred：事件触发足够启动 W3；责任人=维护者；复核=首次例行扫描前 | 待确认 |
| I-003 | required（波次级） | 每一波的 finding 清单与范围 | 该波实施 | 该波实施前 | 扫描报告落盘到子目标 | 按波次 | — | W1–W4、W6 在子目标 verified；W5 扫描 0 中高危见 E-002；W7 清单 = GOAL-007 A-001；W8 清单 = GOAL-008 A-001；W9 清单 = GOAL-009 A-001（均在波次方案前 verified）；W10 清单 = GOAL-010 A-001（2026-08-21 落盘 verified） |

## 台账布局

本 Root 使用 `01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；索引与目录条目共同构成正式记录。
