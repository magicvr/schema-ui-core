---
id: GOAL-001-digital-offer-entitlement
title: 数字 Offer 与权益
status: active
parent: null
created: 2026-09-05
updated: 2026-09-05
version: 0.1.0
progress: 0/4
plan_refs:
  - VP-031-digital-offer-entitlement
primary_plan: VP-031-digital-offer-entitlement
serves_summary: 业务域分支 · 数字 Offer + 薄购买凭证 + 本域权益；复用 VP-029 subject 与钱包 freeze/deduct/unfreeze；Telegram 为可选 Register 接缝；不进入默认 Profile，不解禁通用 Entitlement/Approval Gate。
---

# GOAL-001 · 数字 Offer 与权益

## 概述

在同进程 Go 基座内交付可装配的数字 Offer、薄购买凭证与本域 entitlement。Root 只承接 VP-031 的实现与验证，不扩展为电商、支付网关或通用权益框架。

## 成功标准

1. Offer CRUD 具备 Admin 协议页面、权限键、审计与 C 端上架列表。
2. 购买路径满足余额不足拒绝、`freeze → deduct_frozen`、失败 `unfreeze`，并保证购买凭证与权益发放同事务或等价 fail-closed，有并发测试。
3. 权益有效、过期、次数耗尽均可测试，服务提供前统一核验。
4. 购买与权益只挂 VP-029 `subject_id`，不创建 `admin.users`。
5. Telegram 启用时注册 R1 冻结的演示命令；未启用时测试不依赖 Bot API。
6. 保留 VRev-080 freshness 与 RT-Q03/Q05 评估证据；触发条件变化时复审。
7. 不做类目树/SKU/税/库存/物流订单，不进默认 Profile，不解禁通用 Entitlement 接缝，不改 Charter。
8. 关门前开放 required finding = 0，或均已合法闭合。

## 纲领路线图（P-001）

`progress = 已完成纲领阶段 / 4`。

| 阶段 | 内容 | 检查点 / 状态 |
|------|------|---------------|
| R1 | 合同：Offer 字段、购买状态机、权益形态、命令清单、事务与限流边界 | 待开始；I-031-001～005 与 V-F119 待处理 |
| R2 | Offer + 购买 + 钱包扣款 | 待开始；依赖 R1 |
| R3 | 权益校验 + 可选 Telegram 注册 | 待开始；依赖 R2 |
| R4 | 证据与关门 | 待开始；依赖 R1～R3 |

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|----------|-----------------|------|-------------|-------------|
| I-031-001 | required | 首波权益形态：仅时长、仅次数、或二者并存（一 Offer 一种）。 | 判据 1/3；R2 放行 | R1 | R1 合同冻结并记录用户裁决/方案证据 | open | — | — |
| I-031-002 | required | 购买状态最小子集（是否要 `pending`，还是同步一拍 fulfilled）。 | 判据 2；R2 放行 | R1 | R1 状态机与事务边界设计 | open | — | — |
| I-031-003 | required | 是否允许 Admin 人工发放/撤销权益。 | 判据 3；R3 放行 | R1 | R1 权限与审计边界裁决 | open | — | — |
| I-031-004 | non-blocking | Telegram 命令清单（若 channel.telegram 启用）。 | 判据 5 | R1 | R1 冻结价目/购买/我的权益命令子集 | open | — | — |
| I-031-005 | non-blocking | 模块 id（建议 `biz.digital-offer`）。 | 装配 | R1 | 核验模块命名与 Manifest 约定 | open | — | — |

## freshness 与架构评估

- `consumer_vp`: `VP-031-digital-offer-entitlement`（`vision_ref = schema-ui-core-admin-foundation@0.4.0`）
- `go_issued_at`: 2026-08-10（VP-008 候选 `ed99e88`；2026-08-19 消费有效性恢复）
- `last_freshness_review_at`: 2026-09-05（VRev-080；写入前 clean HEAD `bd9ed5e062cd965ec4f2221ec5d00351023e76f2`；H-002 同进程；PASS）
- `next_freshness_review_trigger`: H-002、Profile、协议 provenance、依赖锁、迁移/Manifest、生产部署/密钥边界或多实例变化。
- RT-Q03 / RT-Q05：本波不需要 Redis；见 VRev-080。
