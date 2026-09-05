---
id: VRev-080-vp031-digital-offer-entitlement-activation
doc_type: vision-review
title: VP-031 激活就绪 · 数字 Offer 与权益
source: self
date: 2026-09-05
scope: VP-031-digital-offer-entitlement 意图 / 边界 / 业务域 freshness（H-002 同进程）/ RT-Q03/Q05 评估 / 工作区对齐
verdict: pass
open_required: 0
status: active
created: 2026-09-05
updated: 2026-09-05
parent: null
version: 0.1.0
---

# VRev-080 · VP-031 激活就绪（数字 Offer 与权益）

## 背景与触发

用户先指令“`/vision` 走流程激活 VP-031，没有问题则交 `/govern` 开设工作区”，并于 2026-09-05 书面确认 H-002 主要部署形态为“同进程模块”。本审视承接 VRev-065 的计划阶段 `pass`，完成 VP-008 `go` 消费前 freshness、RT-Q03/Q05 评估与开区对齐核对。

## 1. 意图、边界与结构选型

**pass**。VP-031 是本仓库首个业务域 VP，以同进程可装配模块交付数字 Offer、薄购买凭证与本域权益；主体与钱包资金原语消费 VP-029。它不是电商 Catalog/SKU/税/库存/物流订单，不进入 `mvp`/`admin` 默认 Profile，不解禁 Admin 通用 Entitlement/Approval Gate，也不修改 Charter。

结构采用单一 lead delivery 工作区：`workspace-031-digital-offer-entitlement` / `GOAL-001-digital-offer-entitlement`。工作区建立只表示实现上下文与四阶段 Root 路线图就位，不构成 Offer、购买或权益能力已经交付。

## 2. P-005 信息需求

| 信息项 | 状态 | 激活判定 |
|--------|------|----------|
| I-031-001～003 | `required open`，最晚 R1 | 不阻断激活与 Root 建立；在 R1 合同冻结前必须由证据或用户裁决关闭，未关闭不得进入 R2 |
| I-031-004～005 | `non-blocking open`，最晚 R1 | 不阻断激活；随 R1 冻结 Telegram 命令清单与模块 id |

本次不把未决信息伪装为已冻结决策，也不创建 R1 子目标或实施业务代码。

## 3. 业务域 freshness（含 H-002）

候选身份：写入前 clean HEAD `bd9ed5e062cd965ec4f2221ec5d00351023e76f2`；`git status --short` 为空。消费基线为 VP-008 用户签发 `go` 的候选 `ed99e88`；其后变更均由 VP-009/010 持续程序或 VP-011～033 的规划、工作区、审计与关门记录承接。

| 域 | 当前事实 | 判定 |
|----|----------|------|
| H-002 部署形态 | 用户于 2026-09-05 书面确认“同进程模块”；C 端 API 与 Admin 复用同一 Go 基座 | ✅ 与 Charter H-002 一致 |
| VP-008 `go` | `go` 消费有效性已恢复；后续业务 VP 仍须逐次 freshness，本报告即本次记录 | ✅ 可消费，不替代实现验收 |
| 硬前置 VP-029 | VP-029 `closed` v0.5.0；Root done 5/5；主体接缝与钱包 freeze/deduct/unfreeze 原语已有审计证据 | ✅ 硬前置满足 |
| 基础设施端口 | VP-026 cache、VP-027/032 rate limiter 均已关闭并保留进程内默认实现；Redis 实现保持 trigger-gated | ✅ 可作为候选基础 |
| Profile / 装配边界 | 本 VP 仍不进入 `mvp`/`admin` 默认集；当前未发现 VP-031 业务实现 | ✅ 激活不冒充实现 |
| Goal / Vision 意见投影 | Vision Review 当前 open required = 0；未发现覆盖 VP-031 激活的开放 required finding | ✅ |

**freshness：PASS**。本次只确认方向、依赖与候选身份仍可用于开区；若 H-002 主要形态、Profile 默认集、协议 provenance、依赖锁、迁移/Manifest、生产部署或密钥边界变化，必须重新 review。

## 4. RT-Q03 缓存评估

| 项 | 评估 |
|----|------|
| 场景 | Offer 列表/标价/上架状态读取与权益核验；购买与扣款仍以持久化事务为权威，不以缓存作为正确性来源 |
| 同进程结论 | 首波单进程模块可直接读取权威存储；现有 `kernel.Cache`/内存实现仅在 R1/R2 有明确性能证据时选择性使用 |
| Redis | **本波不需要 Redis**；不实现跨实例共享缓存或分布式失效 |
| 复审触发 | 多实例部署、跨实例一致性/共享失效需求，或性能证据证明权威存储读取不能满足目标 |

结论：RT-Q03 评估已完成并登记；trigger 已被业务域激活触发，但结论为“不需要 Redis”，Redis 供应商仍 gated。

## 5. RT-Q05 限流评估

| 项 | 评估 |
|----|------|
| 场景 | Telegram ingress 已由 VP-030 的 IP/chat/user 桶覆盖；本域后续 HTTP/命令购买入口可按 subject/来源建立独立请求计数桶 |
| 端口 | 复用 VP-027/032 已交付的进程内 `kernel.RateLimiter` 与原子 `AllowRecord`/Reserve-Cancel 能力 |
| 语义边界 | 业务请求限流与钱包失败预算/购买事务分离；R1 必须冻结桶 key、阈值与拒绝语义，不得用 key-wide `Clear` 抹除既有历史 |
| Redis | **本波不需要 Redis**；同进程 limiter 足以覆盖当前单实例业务域流量 |
| 复审触发 | 多实例部署、需跨实例共享配额，或业务域流量无法共用/独立实例化进程内 limiter |

结论：VP-030 的 ingress 评估继续覆盖通道入口；VP-031 的业务端点使用独立进程内桶即可，RT-Q05 Redis 实现保持 gated。

## 6. 对齐与开区许可

**pass**。VP-031 `vision_ref = schema-ui-core-admin-foundation@0.4.0` 精确匹配唯一 active Charter；VP-029 硬前置已满足；VP-030 为已关闭的软前置。工作区采用 `vision_role: delivery`，`primary_plan = VP-031-digital-offer-entitlement`，Root 为 `GOAL-001-digital-offer-entitlement`，四阶段为 R1 合同 → R2 Offer/购买/钱包扣款 → R3 权益校验/可选 Telegram 注册 → R4 证据与关门。

## Verdict

**pass**（open required = 0）。H-002 同进程形态已由用户书面确认；业务域 freshness PASS；RT-Q03/RT-Q05 均完成评估且本波不需要 Redis。可将 VP-031 `planned → active`（v0.2.0），并交 `/govern` 建立 `workspace-031-digital-offer-entitlement` 与 Root。

## Findings

### 必改（required）

无。

### 建议（recommended）

- **V-F119**：R1 冻结业务域限流桶 key、阈值与拒绝语义，并明确请求计数不使用 key-wide `Clear`。该项由 I-031-001～005 所在 R1 合同阶段一并处理，不阻断本次激活或开区。

## 声明

- 本报告 `source: self`，不是独立 Vision Review。
- 本报告不创建 Goal 审计意见，不修改 Charter，也不宣称 VP-031 已实现。
- I-031-001～003 仍是 R1 required；未关闭前不得进入 R2。
