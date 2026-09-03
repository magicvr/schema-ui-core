---
doc_type: vision-plan
id: VP-031-digital-offer-entitlement
title: 数字 Offer 与权益
status: planned
vision_ref: schema-ui-core-admin-foundation@0.4.0
lead_workspace:
created: 2026-09-02
updated: 2026-09-03
version: 0.1.1
parent: null
---

# VP-031 · 数字 Offer 与权益

## 状态与激活门禁

| 项 | 值 |
|----|-----|
| status | **`planned`**（2026-09-02 · v0.1.0 · 0 区） |
| lead_workspace | 未绑定（激活时按惯例 `workspace-031-digital-offer-entitlement`） |
| Vision required | 计划阶段 self = [VRev-065](../reviews/VRev-065-c-end-paid-services-planned-self.md)；**激活前必须**：① 业务域 freshness（含 H-002 同进程再确认）② RT-Q03/Q05 评估登记 ③ 激活审视 |
| 组合位置 | **业务域分支** · 本仓库**第一个**业务域 VP。卖的是数字服务/权益，**不是**电商 Catalog/SKU/税/库存/物流订单 |

## 意图

把「可售的数字服务」做成一方可复用业务域模块，供 Telegram（及未来其它通道）在同进程内购买与核验：

1. **Offer**：可上架的服务项（名称、标价、币种、上/下架、权益形态：时长和/或次数）。**不是**类目树，**不是**多规格 SKU，**不是**含税商品主数据。
2. **薄购买凭证**：一次成功扣款对应一条购买记录（主体、offer、金额、钱包流水引用、状态）。**不是**电商订单（无收货地址、无履约包裹、无售后工单）。
3. **权益（entitlement）**：购买成功后发放；通道/服务在提供能力前校验「该 subject 是否仍持有有效权益」。这是本波真正卖出的东西。

资金路径复用已有钱包原语：`freeze` → 成功则 `deduct_frozen`，失败则 `unfreeze`。入金不在本 VP（归 VP-029 凭证或未来支付域）。

若 `channel.telegram` 已启用，本模块 **Register** 命令/回调（例如余额、价目、购买）；命令文案与信息架构属本 VP，通道运行时属 VP-030。

> Charter 非目标仍然成立：本 VP 交付的是可装配业务域模块，不是某个 Telegram 产品的终态运营后台。下游仍可只启用本模块的一个 Offer 子集。

## 首波冻结（退出分母）

| 项 | 本 VP 交付 | 不进本 VP |
|----|-----------|-----------|
| Offer | 单层可售项：id、名称、标价（最小货币单位）、币种、状态、权益形态（R1 冻结：时长 / 次数 / 二者之一） | 类目树、多级类目、SKU/变体、税、库存、仓库 |
| 购买 | 薄凭证：subject、offer、金额、钱包 `ref`、状态（pending/paid/fulfilled/cancelled 的最小子集，R1 冻结） | 购物车、多行订单、收货、物流、拆单、发票 |
| 权益 | 发放与校验 API；过期或次数耗尽后无效；Admin 只读/作废（R1 冻结是否允许人工发放） | Admin 功能分支通用 Entitlement/Approval 框架（仍 trigger-gated；本 VP 只做本域权益） |
| 资金 | 只消费钱包 freeze/deduct/unfreeze；购买失败必须解冻 | 密钥生成、支付网关、退款编排（退款若本波需要，仅「取消未履约购买 + unfreeze」，不进渠道退款） |
| 通道 | 可选 Register Telegram 命令；无 Telegram 模块时 HTTP/模块 API 仍可测 | Bot 运行时、webhook、SendMessage 实现 |
| Profile | **不**进入 `mvp`/`admin` 默认集 | 改变装配红线 |
| 事件 | 可在模块内同步调用；**不**要求 typed domain event 接缝解禁 | 跨模块领域事件产品化（Admin 分支仍 gated） |

## 非目标

- 电商三件套：类目、商品主数据、物流订单
- 营销/优惠券/促销引擎
- 订阅计费/发票/用量（路线图候选 #6，另 VP）
- 支付网关、Telegram Stars
- 把 Telegram 通道实现打进本模块
- 解禁 Admin 通用 Entitlement / Approval Gate 接缝
- 改 Charter

## 与相邻 VP 的边界

| VP / 分支 | 关系 |
|-----------|------|
| **VP-011** | S-07/S-08/S-13 仍是未交付 residual，**不**在本 VP 兑现为电商模块。本 VP **不是** S-13 订单管理的替代立项 |
| **路线图业务域候选 1/3** | 本 VP **收窄**为数字 Offer + 薄购买 + 权益；不声称 Catalog 或「订单/支付/退款/退货」整域已成立 |
| **Admin 扩展接缝 · Entitlement** | 本 VP 实现**本域**权益表与校验。通用 Approval/Entitlement 框架仍 gated |
| **VP-008 `go`** | 业务域 freshness：候选身份 + 解锁 scope + **H-002 同进程再确认**（用户 2026-09-02 已口头确认同进程；激活时仍须走发现机制写进 freshness 记录） |
| **VP-026 / VP-027** | **业务域 VP 激活即触发**评估义务：缓存是否需要（Offer 读取可结论「不需要」）；限流是否已被 VP-030 覆盖。评估不可跳过 |
| **VP-029** | **硬前置**：主体 + 钱包资金原语（至少 freeze/deduct/unfreeze；购买扣款不走凭证核销） |
| **VP-030** | 软前置：无通道时本域仍可经 API 测通；有通道时本模块注册命令。不把 webhook 当本 VP 范围 |
| **VP-033** | 占用位的典型占用者：本模块 Register 之后，033 人工台入口必须隐藏。不把运营台/轮询模式当本 VP 范围 |
| **VP-009 / VP-010** | 购买/核销安全与符合性 gap 归持续程序 |

## 方向级退出判据

1. **Offer CRUD**：上/下架与标价变更有 Admin 协议页面 + 权限键 + 审计；C 端可列出上架项。
2. **购买扣款**：余额不足拒绝；成功路径 freeze→deduct 与购买凭证、权益发放同事务或等价 fail-closed；失败 unfreeze；有并发测试。
3. **权益校验**：有效/过期/耗尽三种可测；通道或服务在提供能力前走同一校验 API。
4. **主体对齐**：购买与权益只挂 VP-029 `subject_id`，不创建 `admin.users`。
5. **通道可选**：Telegram 启用时至少注册一套可演示命令（价目/购买/我的权益之一组，R1 冻结清单）；未启用时模块测试不依赖 Bot API。
6. **激活门禁留痕**：freshness 含 H-002 同进程再确认；RT-Q03/Q05 评估已写入路线图位置（允许「不需要 Redis」）。
7. **边界保持**：未做类目树/SKU/税/库存/物流订单；未进默认 Profile；未解禁通用 Entitlement 接缝；未改 Charter。
8. **审计闭合**：开放 required finding = 0（或已合法闭合）。

建议 Root 纲领：R1 合同（Offer 字段、购买状态机、权益形态、命令清单、事务边界）→ R2 Offer + 购买 + 钱包扣款 → R3 权益校验 + 可选 Telegram 注册 → R4 证据与关门。

## 信息需求（P-005）

| id | 要回答的问题 | 级别 | 影响门禁 | 最晚阶段 | 状态 |
|----|--------------|------|----------|----------|------|
| I-031-001 | 首波权益形态：仅时长、仅次数、或二者并存（一 Offer 一种）。 | required | 判据 1/3 | R1 | open |
| I-031-002 | 购买状态最小子集（是否要 `pending` 还是同步一拍 fulfilled）。 | required | 判据 2 | R1 | open |
| I-031-003 | 是否允许 Admin 人工发放/撤销权益（客服纠错）。 | required | 判据 3 | R1 | open |
| I-031-004 | Telegram 命令清单（若 030 已启用）。 | non-blocking | 判据 5 | R1 | open |
| I-031-005 | 模块 id（建议 `biz.digital-offer`）。 | non-blocking | 装配 | R1 | open |

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| — | — | lead | — | `planned` 0 区；硬前置 VP-029；软前置 VP-030 |

## 关门记录

（仅 `closed` / `abandoned` 时填写。）

## 规划修订短史

| date | change |
|------|--------|
| 2026-09-02 | 初创 `planned`：用户否决电商类目/商品/订单三件套；确认基座一方数字 Offer+权益；H-002 同进程。业务域分支首个 VP。 |
| 2026-09-03 | 边界指针：Register 后占用 VP-033 人工台入口；运营台不在本 VP。 |
