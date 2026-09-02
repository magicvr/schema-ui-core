---
doc_type: vision-plan
id: VP-029-wallet-prepaid-instrument
title: 钱包预付资金凭证与外部主体接缝
status: active
vision_ref: schema-ui-core-admin-foundation@0.4.0
lead_workspace: workspace-029-wallet-prepaid-instrument
created: 2026-09-02
updated: 2026-09-02
version: 0.2.0
parent: null
---

# VP-029 · 钱包预付资金凭证与外部主体接缝

## 状态与激活门禁

| 项 | 值 |
|----|-----|
| status | **`active`**（2026-09-02 · v0.2.0 · lead `workspace-029-wallet-prepaid-instrument`） |
| lead_workspace | `workspace-029-wallet-prepaid-instrument`（2026-09-02 `/govern` 开区） |
| Vision required | 计划阶段 self = [VRev-065](../reviews/VRev-065-c-end-paid-services-planned-self.md)；激活独立审视 = [VRev-066](../reviews/VRev-066-vp029-wallet-prepaid-instrument-independent.md) `pass`（0 required）+ Admin 类 freshness PASS（`29727510`→`b5c39dfb`） |
| 组合位置 | **Admin 功能分支** · 扩展已交付的 `admin.wallet`（VP-011 S-14），**不是**支付/结算业务域 |

## 意图

为同进程 C 端（首个具名消费者 = Telegram 付费服务下游）补齐两条**可复用、通道无关**的基座能力：

1. **外部主体接缝**：`(issuer, external_id) → subject_id`。C 端用户**不得**被做成可登录的 `admin.users` 行。`subject_id` 作为钱包 `owner_id` 与后续权益主体的稳定主键。
2. **预付资金凭证（密钥 / 卡密）**：管理员批次生成、导出给卡网售卖、作废；C 端或通道模块持主体身份核销后，把面额记入该主体的钱包账本。用于在接入在线支付之前完成充值。

本 VP **不重开 VP-011**。现有账本原语保持：`adjust` / `freeze` / `unfreeze` / `deduct_frozen`、三余额恒等、不可变流水、幂等键、对账。凭证核销是**新的入金通道**，不是平行账本。

> **定位**：已有钱包是 Admin 能力面，**不等于**路线图业务域「支付/结算」已成立。在线支付网关、退款编排、卡组织对接仍归未来支付业务域，不进本 VP。

## 首波冻结（退出分母）

| 项 | 本 VP 交付 | 不进本 VP |
|----|-----------|-----------|
| 外部主体 | 通道无关登记/查找：`(issuer, external_id) → subject_id`；幂等 get-or-create；钱包 `ownerExists` 承认已登记主体（不要求 `users` 行） | 完整 C 端 IAM；Admin 登录；OIDC/SSO；把 Bot 用户写入 `admin.users` |
| 凭证模型 | 批次、面额（最小货币单位）、状态（未用/已核销/作废）、过期（可选）；**只存密钥哈希**，明文仅在生成时一次性出示 | 卡网开放 API、渠道分润、面额 SKU 目录（那是 Offer 域） |
| 核销 | 模块级 `Redeem(subjectID, code)`：原子「标记已用 + 账本入金」；幂等（同一 code 重复核销不双记）；`ref_type/ref_id` 指向凭证 | 公开无主体的 HTTP 核销（伪造入账目标）；支付回调入金 |
| 账本 | 入金走现有账本（新 `entry_type` 或带 `ref_type=voucher` 的 `adjust`，R1 冻结）；不破坏三余额恒等与对账 | 改写 freeze/deduct 语义；第二套余额 |
| Admin 面 | 批次生成 / 导出 / 作废 / 查询；权限键从 `wallet.adjust` 拆或显式复用（R1 冻结） | 卡网对账报表、营销发放规则 |
| C 端 HTTP | 可选：在已有 C 端主体会话下的自助核销。无 C 端主体时，只暴露模块 API 供通道模块调用 | Telegram webhook（归 VP-030）；购买/权益（归 VP-031） |
| Profile | **不**把新模块塞进 `mvp`/`admin` 默认集；本波是 `admin.wallet` 增量（及主体接缝的持久化贡献） | 改变 Manifest 装配红线 |

## 非目标

- 类目树、商品/SKU/税、库存、物流订单
- 在线支付网关、Telegram Stars、退款状态机
- 重开 VP-011；把本 VP 写成「支付业务域」
- 改 Charter；把 C 端用户做成 Admin 账号
- Redis / 多实例（不消耗 RT-Q03/Q05）；不解除 typed domain event gated

## 与相邻 VP 的边界

| VP / 分支 | 关系 |
|-----------|------|
| **VP-011 `admin.wallet`** | 只读基线。本 VP 增量扩展，不重开、不改已关门判据 |
| **VP-003 / VP-004** | 主体接缝与凭证表走模块 Persistence + 全局迁移台账；公共面不泄漏 `*sql.Tx` |
| **VP-008 `go`** | Admin 类能力；激活前 freshness review |
| **VP-012** | 核销与生成必须有幂等与审计；不重开 Job/Token VP |
| **VP-026 / VP-027** | 本 VP 不是业务域激活，**不单独消耗** RT-Q03/Q05。C 端流量的限流评估挂 VP-030 激活门禁 |
| **VP-030 `channel.telegram`** | 首个 `issuer=telegram` 消费者。本 VP **不得** import Telegram。硬前置方向：030 激活需要本 VP 的主体接缝已可用 |
| **VP-031 数字 Offer** | 购买扣款消费钱包 `freeze`/`deduct_frozen`；权益主体 = 本 VP 的 `subject_id`。硬前置方向：031 激活需要本 VP 资金原语已可用 |
| **业务域「支付」** | 仍未成立。凭证 ≠ 支付网关 |

## 方向级退出判据

1. **主体接缝可用**：`(issuer, external_id) → subject_id` 幂等登记/查找有测试；未登记主体不能开户；不创建 `admin.users` 行。**查询与 get-or-create 不依赖 `admin.wallet` 已启**（权益/通道在未启用钱包时仍能挂主体；落点由 I-029-001 冻结）。
2. **凭证生命周期**：生成（哈希存储 + 一次性明文）/ 导出 / 作废 / 过期拒绝 有测试；明文不落库、不进审计原文。
3. **核销原子且幂等**：`Redeem` 成功则账本入金与凭证状态一致；重复核销不双记；并发双花 fail-closed。
4. **账本不变式保持**：三余额恒等、流水快照链、对账 Job 仍通过既有钱包测试；新入金类型纳入 apply 表。
5. **Admin 可操作**：批次生成/导出/作废有协议驱动页面 + 权限键 + 操作审计。
6. **边界保持**：未改 Charter；未改 `mvp`/`admin` 默认模块集的装配语义；未引入支付网关或 Telegram 依赖；未重开 VP-011。
7. **审计闭合**：开放 required finding = 0（或已合法闭合）。

详细纲领阶段由 lead Root（P-001）书写。建议：R1 合同（主体模型 / 凭证哈希 / entry_type / 权限键）→ R2 主体接缝 + 账本入金 → R3 Admin 批次面 + 导出 → R4 证据与关门。本 VP 不写 Goal 五件套。

## 信息需求（P-005）

| id | 要回答的问题 | 级别 | 影响门禁 | 最晚阶段 | 状态 |
|----|--------------|------|----------|----------|------|
| I-029-001 | 主体落点：薄模块 vs `authsession` vs `admin.wallet` 表。公共契约必须通道无关。须同时冻结：`owner_type` 是否新增取值（现行 CHECK 仅 `user/business/system`）；`OwnerExistsFunc` 改为「已登记主体」、**禁止**回退 `UserByID`；查询/get-or-create 不依赖 `admin.wallet` 已启。W13 F-012 孤儿账本相对主体登记表。 | required | 方案冻结 + 判据 1 | R1 | open |
| I-029-002 | 核销入金的 `entry_type`：新类型 vs `adjust` + `ref_type=voucher`。 | required | 判据 3/4 | R1 | open |
| I-029-003 | 生成权限键：复用 `wallet.adjust` vs 新 `wallet.voucher.issue`。 | required | 判据 5 | R1 | open |
| I-029-004 | 导出格式（CSV/TXT、是否含明文、一次性下载）。 | non-blocking | 判据 5 | R3 | open |
| I-029-005 | C 端自助核销 HTTP 是否本波交付，或仅模块 API（Telegram 进程内调用）。若选 HTTP，本 VP 必须完成 RT-Q05 精神的限流评估，不得推到未激活的 VP-030。 | non-blocking | 判据面 | R1 | open（默认倾向：模块 API 必做，HTTP 自助核销可选） |
| I-029-006 | 凭证哈希与双花合同：哈希算法（默认候选 = 高熵码 SHA-256 或 HMAC-SHA256+pepper；禁止 6 位恢复码或 bcrypt 当卡密默认）；码字母表与长度（熵下限）；核销常时比较；`UNIQUE(code_hash)`（或等价）+ 同事务「未用→已核销 AND 账本入金」，并发失败者 fail-closed，重复 Redeem 不双记。 | required | 判据 2/3 | R1 | open |

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| workspace-029-wallet-prepaid-instrument | GOAL-001-wallet-prepaid-instrument | lead | 2026-09-02 | 唯一 delivery；VRev-066 independent `pass` + Admin 类 freshness PASS |

## 关门记录

（仅 `closed` / `abandoned` 时填写。）

## 规划修订短史

| date | change |
|------|--------|
| 2026-09-02 | 初创 `planned`：用户确认 Telegram 付费服务下游为真实触发；结构选型 C（基座一方可复用）+ 切分 1（钱包密钥 = Admin 功能，不是支付域）+ 外部主体接缝（不建 Admin 登录账号）。与 VP-030/031 同批落盘。 |
| 2026-09-02 | **激活** `planned → active` v0.2.0（用户指令：独立审视通过则激活并开区）。[VRev-066](../reviews/VRev-066-vp029-wallet-prepaid-instrument-independent.md) independent `pass`（0 required）。Admin 类 freshness **PASS**（`29727510` → `b5c39dfb`：协议 pin / 依赖锁 / 迁移台账 / Profile 装配 / provenance 五域零变更；区间代码 = VP-028 已审结目 + VP-009 W16/W17；不暂挂 `go`）。V-F110 → fixed（本独立审视）；V-F111/112/113 → 激活+开区事务内 fixed（I-029-001 扩写 + I-029-006 + freshness 留痕 + VP-028 组合索引同步）。lead `workspace-029-wallet-prepaid-instrument`。 |
