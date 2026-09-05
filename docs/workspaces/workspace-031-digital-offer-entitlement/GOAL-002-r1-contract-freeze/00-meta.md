---
id: GOAL-002-r1-contract-freeze
title: R1 合同冻结（Offer 字段 / 购买状态机 / 权益形态 / 命令清单 / 事务与限流边界）
status: done
parent: GOAL-001-digital-offer-entitlement
created: 2026-09-05
updated: 2026-09-05
version: 0.1.0
progress: 3/3
plan_refs:
  - VP-031-digital-offer-entitlement
primary_plan: VP-031-digital-offer-entitlement
serves_summary: 承载 VP-031 R1：冻结数字 Offer 字段与状态机、购买同步一拍 fulfilled 状态机、权益时长/次数二者并存形态、Telegram 命令清单、单事务购买边界（freeze→deduct_frozen→凭证+权益）与业务限流桶语义。合同正文 = 本目标 D-002；R2～R4 实施与验收以之为分母。
---

# GOAL-002 · R1 合同冻结

## 概述

执行 Root 纲领 **R1**：在钱包原语（VP-029 `freeze / deduct_frozen / unfreeze`）、subject 身份（VP-029）与 Telegram 分发缝（VP-030 `kernel.TelegramDispatcher`）之上，冻结数字 Offer 业务域的**合同分母**——Offer 模型与 CRUD 面、购买单事务边界、权益形态与核验/消耗 API、可选 Telegram 命令、业务限流桶与错误码。**合同正文 = 本目标 D-002**。本目标只冻结合同与裁决，不写业务实现代码（实现归 R2/R3）。

对齐递归：GOAL-002 → Root GOAL-001（R1）→ VP-031（判据 1/2/3/5/6 与首波冻结表）→ Charter @0.4.0。不做类目树/SKU/税/库存/物流订单；不进默认 Profile；不解禁通用 Entitlement/Approval 接缝；不改 Charter。

## 纲领检查点（P-001）

| 检查点 | 内容 | 状态 |
|--------|------|------|
| C1 | **信息裁决**：I-031-001～003 required（P-004 用户书面裁决）+ I-031-004/005 non-blocking 默认冻结 | **已关门**（2026-09-05 用户书面全部采纳建议项——D-001） |
| C2 | **合同正文冻结**：D-002 冻结 Offer 字段/状态机、购买单事务边界、权益核验/消耗、命令清单、限流桶、错误码、迁移与装配面 | **已关门**（D-002 v1.0.0 `accepted`，2026-09-05 经 A-006 independent closure 复审通过） |
| C3 | **审视与关门**：self 审计（A-001）+ codex 独立审计（A-002，gpt 5.6 sol · medium）；意见响应整改；Root 信息台账回写 | **已关门**（A-001→A-007 循环；A-006 `pass` open required 0；A-007 登记闭合与关门判定） |

`progress` = 已关门检查点数 / 3。当前 **3/3**。

## 成功标准（方向级）

1. 信息裁决完成：I-031-001～003 required 全部 verified（用户书面），I-031-004/005 non-blocking 冻结（D-001）。
2. 合同正文冻结：D-002 覆盖 VP-031「首波冻结」表全部行（Offer/购买/权益/资金/通道/Profile/事件），且逐条可追溯到判据。
3. 事务边界可执行：购买 freeze→deduct_frozen→凭证+权益的单事务（或等价 fail-closed）设计、幂等键与失败路径在 D-002 中可照做（判据 2 前置）。
4. 限流语义冻结：业务桶 key/阈值/拒绝语义显式，且请求计数永不 key-wide `Clear`（V-F119）。
5. 审计闭合：self + independent 意见落盘，开放 required = 0 后关门。

## 信息就绪与未知项

与 Root / VP-031 同号镜像。本目标关闭 I-031-001～005（最晚阶段均为 R1）。

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-031-001 | required | 首波权益形态：仅时长、仅次数、或二者并存（一 Offer 一种） | 判据 1/3；R2 放行 | C1 | 用户裁决 | **verified** | — | 2026-09-05 用户裁决：二者并存，一 Offer 固定一种（D-001） |
| I-031-002 | required | 购买状态最小子集（pending 与否） | 判据 2；R2 放行 | C1 | 用户裁决 | **verified** | — | 2026-09-05 用户裁决：同步一拍 fulfilled，无 pending（D-001） |
| I-031-003 | required | 是否允许 Admin 人工发放/撤销权益 | 判据 3；R3 放行 | C1 | 用户裁决 | **verified** | — | 2026-09-05 用户裁决：只读 + 作废，不开放人工发放（D-001） |
| I-031-004 | non-blocking | Telegram 命令清单（channel.telegram 启用时） | 判据 5 | C1 | lead 建议默认 + 用户可否决 | **verified** | — | 2026-09-05 冻结默认：price / buy / entitlements 三命令（D-001；合同审查期用户可否决） |
| I-031-005 | non-blocking | 模块 id | 装配 | C1 | lead 建议默认 + 用户可否决 | **verified** | — | 2026-09-05 冻结默认：`biz.digital-offer`（D-001；合同审查期用户可否决） |

## 父目标

- `GOAL-001-digital-offer-entitlement`（Root · 纲领 R1）

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺记账；索引文件在本目标 `01-decision.md` / `02-execution.md` / `03-audit.md`。

## 备注

- 审计模式：R1 合同为资金路径事务设计的基础，C3 采用 **self（A-001）+ independent（A-002，本地 codex · gpt 5.6 sol · medium）**；R2/R3 阶段关门按 AGENTS P-003 风险表执行（资金/数据面 → independent）。
- D-002 冻结后，R2/R3/R4 实施与验收以本合同为分母；偏离合同需先修合同（02 决策）再实施。
