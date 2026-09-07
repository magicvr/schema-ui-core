---
doc_type: goal-decision
id: D-001-workspace-root-establishment
parent: GOAL-001-digital-offer-entitlement
date: 2026-09-05
status: active
version: 0.1.0
---

# D-001 · 工作区与 Root 建立

## 已确认决策

- 用户授权顺序：先由 `/vision` 激活 VP-031；无 required 阻断后交 `/govern` 开设工作区。
- H-002：用户于 2026-09-05 书面确认采用**同进程模块**。
- 工作区：`workspace-031-digital-offer-entitlement`；Root：`GOAL-001-digital-offer-entitlement`；`vision_role: delivery`；`primary_plan: VP-031-digital-offer-entitlement`。
- 纲领路线图：R1 合同 → R2 Offer/购买/钱包扣款 → R3 权益校验/可选 Telegram 注册 → R4 证据与关门。
- 硬前置 VP-029 已满足；VP-030 为已关闭软前置。
- RT-Q03/RT-Q05 本波均不需要 Redis；多实例或跨实例共享状态出现时复审。
- 本轮只建立 Root 与治理台账，不创建 R1 子目标，不实施业务代码，不生成 Goal 审计结论。

## 信息与门禁

- I-031-001～003 保持 `required open`，最晚 R1；未关闭不得进入对应 R2/R3 门禁。
- I-031-004～005 保持 `non-blocking open`，随 R1 冻结。
- VRev-080 V-F119 为 recommended：R1 冻结业务域限流桶 key、阈值、拒绝语义；请求计数不得使用 key-wide `Clear`。

## freshness 三字段

| 字段 | 值 |
|------|----|
| consumer_vp | `VP-031-digital-offer-entitlement`（`vision_ref schema-ui-core-admin-foundation@0.4.0`） |
| last_freshness_review_at | 2026-09-05 · VRev-080 · 写入前 clean HEAD `bd9ed5e062cd965ec4f2221ec5d00351023e76f2` · H-002 同进程 · PASS |
| next_freshness_review_trigger | H-002 主要形态、Profile 默认集、协议 provenance、依赖锁、Offer/权益迁移与 Manifest/装配、生产部署/密钥边界变化，或多实例触发 RT-Q03/Q05 复核 |

## 未选方案

- 不拆分独立服务：会改变事务一致性、缓存与限流边界，需要先修订 VP。
- 不把本意图塞入其它工作区：VP-031 需要独立 goal-tree 与关门证据。
- 不预制 Redis：当前同进程首波没有跨实例共享状态或配额需求。
