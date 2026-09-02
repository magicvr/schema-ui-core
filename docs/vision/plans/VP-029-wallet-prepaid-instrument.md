---
doc_type: vision-plan
id: VP-029-wallet-prepaid-instrument
title: 钱包预付资金凭证与外部主体接缝
status: active
vision_ref: schema-ui-core-admin-foundation@0.4.0
lead_workspace: workspace-029-wallet-prepaid-instrument
created: 2026-09-02
updated: 2026-09-02
version: 0.4.0
parent: null
---

# VP-029 · 钱包预付资金凭证与外部主体接缝

## 状态与门闩

| 项 | 值 |
|----|-----|
| status | **`active`**（2026-09-02 · v0.4.0 · 用户确认结构选型 **A**：reopen；不新开 VP / 不新开区） |
| lead_workspace | `workspace-029-wallet-prepaid-instrument`（唯一 delivery；Root 同步 `done → active`，加 R5） |
| Vision required | 计划阶段 self = [VRev-065](../reviews/VRev-065-c-end-paid-services-planned-self.md)；激活独立审视 = [VRev-066](../reviews/VRev-066-vp029-wallet-prepaid-instrument-independent.md) `pass`；首波关门就绪 = [VRev-067](../reviews/VRev-067-vp029-wallet-prepaid-instrument-close-out.md) self `pass`（0 required；**原文与 verdict 不改写**）；reopen 就绪 = [VRev-068](../reviews/VRev-068-vp029-reopen-my-wallet-self-redeem.md) self `pass`（0 required） |
| 组合位置 | **Admin 功能分支** · 扩展已交付的 `admin.wallet`（VP-011 S-14），**不是**支付/结算业务域 |

## 重开（2026-09-02 · 结构选型 A）

用户书面确认（P-004 / V5）：在已有「我的钱包」页加入预付凭证充值入口。判定 = **reopen 本 VP**，lead 仍为 workspace-029；**不**放 [workspace-010](../../workspaces/workspace-010-design-implementation-conformance/)（符合性程序不承载钱包资金路径），**不**重开 VP-011 / [GOAL-022](../../workspaces/workspace-011-admin-functional-modules/GOAL-022-my-wallet-self-service/00-meta.md)（当时冻结为只读自服务）。

本重开只撤销「组合层已成功交付预付凭证**全部**表面」的效力中 **HTTP 自助核销 + 我的钱包入口** 这一块。下列**保持原样**：

- R1～R4 子目标 `done`、代码、测试、Goal 审计 A 条目的原文与 verdict
- 2026-09-02 历史关门记录（下节）——作为当时分母（判据 #1～#7、无 HTTP Redeem）下的实施事实，**不再**单独构成现行 `closed`
- I-029-001～006 的历史闭合结论；I-029-005「R1 波次仅模块 API」仍为当时分母下的真实裁决，R5 用新信息项承接 HTTP 增量

## 历史关门记录（2026-09-02 · 组合效力已被 R5 重开取代）

> 下表是当时有界关门的原文摘要，保留为实施史。用户已确认 reopen，其作为本 VP 现行 `closed` 的效力被取代；判据 #1～#7 相对**当时分母**的核验结论不改写。

| date | 当时 outcome | summary | 现行效力 |
|------|----------------|---------|----------|
| 2026-09-02 | 当时 `closed` v0.3.0 | 主体接缝 + 模块 `Redeem` + Admin 批次面；VRev-067 `pass`；I-029-005 = 仅 Go 内部 API | **已被 R5 重开取代**。不得再把「无 HTTP 核销」写成现行成功交付 |

## 意图

为同进程 C 端（首个具名消费者 = Telegram 付费服务下游）补齐两条**可复用、通道无关**的基座能力：

1. **外部主体接缝**：`(issuer, external_id) → subject_id`。C 端用户**不得**被做成可登录的 `admin.users` 行。`subject_id` 作为钱包 `owner_id` 与后续权益主体的稳定主键。
2. **预付资金凭证（密钥 / 卡密）**：管理员批次生成、导出给卡网售卖、作废；核销后把面额记入目标账本。用于在接入在线支付之前完成充值。

**R5 增量（现行）**：在已登录 Admin 会话下，把预付凭证核销暴露为**身份作用域 HTTP**，并在已有「我的钱包」页提供可发现入口。入账目标 = 当前用户自己的 `owner_type=user` 钱包。这不是匿名公开核销，也不是把 Admin 用户做成 `subject`。

本 VP **不重开 VP-011**。现有账本原语保持：`adjust` / `freeze` / `unfreeze` / `deduct_frozen`、三余额恒等、不可变流水、幂等键、对账。凭证核销是**新的入金通道**，不是平行账本。GOAL-022 的「只读自服务」在本增量内被**有界放开**：仅增加本人预付凭证充值，不开放调账/冻结/提现/支付网关。

> **定位**：已有钱包是 Admin 能力面，**不等于**路线图业务域「支付/结算」已成立。在线支付网关、退款编排、卡组织对接仍归未来支付业务域，不进本 VP。

## 首波冻结（R1～R4 历史分母 · 已核销）

| 项 | 本 VP 已交付（R1～R4） | 不进本 VP |
|----|-----------|-----------|
| 外部主体 | 通道无关登记/查找：`(issuer, external_id) → subject_id`；幂等 get-or-create；钱包 `ownerExists` 承认已登记主体（不要求 `users` 行） | 完整 C 端 IAM；Admin 登录；OIDC/SSO；把 Bot 用户写入 `admin.users` |
| 凭证模型 | 批次、面额（最小货币单位）、状态（未用/已核销/作废）、过期（可选）；**只存密钥哈希**，明文仅在生成时一次性出示 | 卡网开放 API、渠道分润、面额 SKU 目录（那是 Offer 域） |
| 核销（模块） | 模块级 `Redeem(subjectID, code)`：原子「标记已用 + 账本入金」；幂等；`ref_type/ref_id` 指向凭证 | 公开无主体的 HTTP 核销；支付回调入金 |
| 账本 | 入金走现有账本（`adjust` + `ref_type=voucher`）；不破坏三余额恒等与对账 | 改写 freeze/deduct 语义；第二套余额 |
| Admin 批次面 | 批次生成 / 导出 / 作废 / 查询；权限键 `wallet.voucher.issue` | 卡网对账报表、营销发放规则 |
| Profile | **不**把新模块塞进 `mvp`/`admin` 默认集；本波是 `admin.wallet` 增量 | 改变 Manifest 装配红线 |

## R5 增量分母（现行 · 2026-09-02）

| 项 | 本增量交付 | 不进本增量 |
|----|-----------|-----------|
| HTTP | 已登录、身份作用域的自助核销（会话推导入账目标；默认候选路径 `POST /api/wallet/me/redeem`，S1 冻结） | 匿名/未登录公开核销；代他人核销；C 端 subject 会话 HTTP（仍归通道模块 / VP-030） |
| 入账 | 调用方自己的 `owner_type=user` 账本；CAS + `adjust`/`ref_type=voucher` 与模块核销同合同 | 把 Admin 自助核销走 `Redeem(subjectID)` 记入 subject 账；平行账本 |
| 页面 | 「我的钱包」可发现的预付凭证充值入口；成功后刷新余额与流水 | 重做 `/wallet` 管理端；支付/提现 UI |
| 权限 | identity-only（与 `GET/POST /api/wallet/me` 同款）；**不**复用 `wallet.voucher.issue` | 新权限键；把发卡权限授给全体登录用户 |
| 限流 | 对本 authenticated HTTP 完成 RT-Q05 精神评估并落盘（I-029-008） | 把限流义务推到未激活的 VP-030；消耗 RT-Q05 Redis trigger |

## 非目标

- 类目树、商品/SKU/税、库存、物流订单
- 在线支付网关、Telegram Stars、退款状态机
- 重开 VP-011；把本 VP 写成「支付业务域」
- 改 Charter；把 C 端用户做成 Admin 账号
- Redis / 多实例（不消耗 RT-Q03/Q05 trigger）；不解除 typed domain event gated
- 匿名公开核销 HTTP；把「我的钱包」做成管理端调账面

## 与相邻 VP 的边界

| VP / 分支 | 关系 |
|-----------|------|
| **VP-011 `admin.wallet`** | 只读基线。本 VP 增量扩展，不重开、不改已关门判据。GOAL-022 只读边界由本 VP R5 **有界放开**（仅本人凭证充值） |
| **VP-010** | **不**承接本增量。符合性程序不承载钱包资金路径（W12 D-005 先例） |
| **VP-003 / VP-004** | 主体接缝与凭证表走模块 Persistence + 全局迁移台账；公共面不泄漏 `*sql.Tx` |
| **VP-008 `go`** | Admin 类能力；R5 若只做 additive 路由/页面动作、不改 Profile 默认集 / 模块矩阵 / Manifest 装配语义，则不暂挂 `go`；实施若触碰共同门禁须 freshness |
| **VP-012** | 核销与生成必须有幂等与审计；不重开 Job/Token VP |
| **VP-026 / VP-027** | 本 VP 不是业务域激活，**不单独消耗** RT-Q03/Q05。R5 限流评估用已交付的内存限流器端口，不实现 Redis |
| **VP-030 `channel.telegram`** | 首个 `issuer=telegram` 消费者。本 VP **不得** import Telegram。硬前置方向：030 激活需要本 VP 的主体接缝已可用（R2 已交付，不因 R5 撤回） |
| **VP-031 数字 Offer** | 购买扣款消费钱包 `freeze`/`deduct_frozen`；权益主体 = 本 VP 的 `subject_id`。硬前置方向：031 激活需要本 VP 资金原语已可用（R2 已交付） |
| **业务域「支付」** | 仍未成立。凭证 ≠ 支付网关 |

## 方向级退出判据

**R1～R4 历史判据（#1～#7）**：相对当时分母已于 2026-09-02 verified（VRev-067）。重开**不改写**其核验结论。现行关门还须满足 R5 判据。

1. **主体接缝可用**：（历史 verified）`(issuer, external_id) → subject_id` 幂等登记/查找有测试；未登记主体不能开户；不创建 `admin.users` 行。查询与 get-or-create 不依赖 `admin.wallet` 已启。
2. **凭证生命周期**：（历史 verified）生成（哈希存储 + 一次性明文）/ 导出 / 作废 / 过期拒绝 有测试；明文不落库、不进审计原文。
3. **核销原子且幂等**：（历史 verified，模块 `Redeem`）成功则账本入金与凭证状态一致；重复核销不双记；并发双花 fail-closed。
4. **账本不变式保持**：（历史 verified）三余额恒等、流水快照链、对账 Job 仍通过既有钱包测试；新入金类型纳入 apply 表。
5. **Admin 可操作**：（历史 verified）批次生成/导出/作废有协议驱动页面 + 权限键 + 操作审计。
6. **边界保持**：未改 Charter；未改 `mvp`/`admin` 默认模块集的装配语义；未引入支付网关或 Telegram 依赖；未重开 VP-011。**R5 仍适用。**
7. **审计闭合**：开放 required finding = 0（或已合法闭合）。**R5 仍适用（含本增量 scope）。**
8. **Admin 已登录自助核销 HTTP**：（R5）会话身份推导入账目标；禁止匿名公开核销；入账为调用方自己的 `owner_type=user` 钱包；原子且幂等（与模块核销同 CAS+adjust 合同）；重复核销不双记；不得把 Admin 自助核销记入 `owner_type=subject` 账。
9. **我的钱包入口**：（R5）`/my-wallet` 有可发现的预付凭证充值入口；成功后余额与流水刷新。
10. **限流评估落盘**：（R5）对本 authenticated HTTP 完成 RT-Q05 精神评估并留痕（内存桶足够，或用户书面 residual）；不得推到未激活的 VP-030；不消耗 RT-Q05 Redis trigger。

详细纲领阶段由 lead Root（P-001）书写。R1～R4 已 done；现行 = **R5 Admin 自助核销 HTTP + 我的钱包入口**。本 VP 不写 Goal 五件套。

## 信息需求（P-005）

| id | 要回答的问题 | 级别 | 影响门禁 | 最晚阶段 | 状态 |
|----|--------------|------|----------|----------|------|
| I-029-001 | 主体落点：薄模块 vs `authsession` vs `admin.wallet` 表。公共契约必须通道无关。须同时冻结：`owner_type` 是否新增取值（现行 CHECK 仅 `user/business/system`）；`OwnerExistsFunc` 改为「已登记主体」、**禁止**回退 `UserByID`；查询/get-or-create 不依赖 `admin.wallet` 已启。W13 F-012 孤儿账本相对主体登记表。 | required | 方案冻结 + 判据 1 | R1 | closed（D-002 · 独立主体表 subjects） |
| I-029-002 | 核销入金的 `entry_type`：新类型 vs `adjust` + `ref_type=voucher`。 | required | 判据 3/4 | R1 | closed（D-002 · 复用 adjust） |
| I-029-003 | 生成权限键：复用 `wallet.adjust` vs 新 `wallet.voucher.issue`。 | required | 判据 5 | R1 | closed（D-002 · 新增 wallet.voucher.issue） |
| I-029-004 | 导出格式（CSV/TXT、是否含明文、一次性下载）。 | non-blocking | 判据 5 | R3 | closed（D-001 · API 一次性返回明文数组） |
| I-029-005 | C 端自助核销 HTTP 是否**首波**交付，或仅模块 API（Telegram 进程内调用）。若选 HTTP，本 VP 必须完成 RT-Q05 精神的限流评估，不得推到未激活的 VP-030。 | non-blocking | 当时判据面 | R1 | closed（D-002 · 仅交付 Go 模块内部 API；**不**表示 R5 禁止 HTTP） |
| I-029-006 | 凭证哈希与双花合同：哈希算法（默认候选 = 高熵码 SHA-256 或 HMAC-SHA256+pepper；禁止 6 位恢复码或 bcrypt 当卡密默认）；码字母表与长度（熵下限）；核销常时比较；`UNIQUE(code_hash)`（或等价）+ 同事务「未用→已核销 AND 账本入金」，并发失败者 fail-closed，重复 Redeem 不双记。 | required | 判据 2/3 | R1 | closed（D-002 · 高熵码 + SHA-256 + 单事务 CAS 原子核销入金） |
| I-029-007 | R5 HTTP 路径与服务函数形状（默认候选 = `POST /api/wallet/me/redeem`，body `{code}`，identity-scoped）。入账 `owner_type=user` 已由 reopen D-003 冻结，本项只收路径/函数名。 | required | 判据 8 · GOAL-005 S1 | R5 S1 | collecting |
| I-029-008 | 已登录核销 HTTP 的 RT-Q05 精神限流：专用内存桶（按 user id）vs 复用现有桶 vs 书面 residual。禁止匿名面；不消耗 Redis trigger。 | required | 判据 10 · GOAL-005 S1 | R5 S1 | collecting |
| I-029-009 | 自助核销权限模型：identity-only vs 新权限键 vs 复用 `wallet.voucher.issue`。 | required | 判据 8 · GOAL-005 S1 | R5 重开 | closed（Root D-003 · identity-only，与 `/api/wallet/me` 同款；不复用发卡键） |

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| workspace-029-wallet-prepaid-instrument | GOAL-001-wallet-prepaid-instrument | lead | 2026-09-02 | 唯一 delivery；2026-09-02 reopen 承接 R5，不新开区 |

## 关门记录

2026-09-02 · 当时 **`active → closed` v0.3.0**（用户书面确认 · VRev-067）。相对当时分母（判据 #1～#7、无 HTTP Redeem）的证据链保留。**现行 status = `active` v0.4.0**（R5 重开）；再次关门须满足判据 #6～#10 且用户确认。

## 规划修订短史

| date | change |
|------|--------|
| 2026-09-02 | 初创 `planned`：用户确认 Telegram 付费服务下游为真实触发；结构选型 C（基座一方可复用）+ 切分 1（钱包密钥 = Admin 功能，不是支付域）+ 外部主体接缝（不建 Admin 登录账号）。与 VP-030/031 同批落盘。 |
| 2026-09-02 | **激活** `planned → active` v0.2.0（用户指令：独立审视通过则激活并开区）。[VRev-066](../reviews/VRev-066-vp029-wallet-prepaid-instrument-independent.md) independent `pass`（0 required）。Admin 类 freshness **PASS**（`29727510` → `b5c39dfb`）。lead `workspace-029-wallet-prepaid-instrument`。 |
| 2026-09-02 | **首波关门** `active → closed` v0.3.0（用户书面确认）。[VRev-067](../reviews/VRev-067-vp029-wallet-prepaid-instrument-close-out.md) self `pass`（0 required）。Root 当时 done 4/4。 |
| 2026-09-02 | **reopen** `closed → active` v0.4.0（用户确认结构选型 A）。R5 = Admin 已登录自助核销 HTTP + 「我的钱包」入口。VRev-068 self `pass`。不新开区、不新开 VP、不放 workspace-010。 |
