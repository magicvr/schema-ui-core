---
doc_type: vision-plan
id: VP-031-digital-offer-entitlement
title: 数字 Offer 与权益
status: closed
vision_ref: schema-ui-core-admin-foundation@0.4.0
lead_workspace: workspace-031-digital-offer-entitlement
created: 2026-09-02
updated: 2026-09-05
version: 0.3.4
parent: null
---

# VP-031 · 数字 Offer 与权益

## 状态与激活门禁

| 项 | 值 |
|----|-----|
| status | **`closed`**（2026-09-05 · v0.3.4 · 第 3 次关门：A-012 `conditional` 2 required → D-004 fixed ×2 + A-013 closed ×2 → A-014 independent `pass` 0 required → F-001 关门前置加固（A-015 fixed）→ 重新关门） |
| lead_workspace | `workspace-031-digital-offer-entitlement`（唯一 lead delivery；Root `GOAL-001-digital-offer-entitlement` `done · 4/4`，R1～R4 全部关门） |
| Vision required | 计划阶段 self = [VRev-065](../reviews/VRev-065-c-end-paid-services-planned-self.md)；激活审视 = [VRev-080](../reviews/VRev-080-vp031-digital-offer-entitlement-activation.md) self `pass`（0 required；业务域 freshness PASS；H-002 同进程书面确认；RT-Q03/Q05 = 本波不需要 Redis） |
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
| **VP-008 `go`** | 业务域 freshness：候选身份 + 解锁 scope + **H-002 同进程再确认**（用户 2026-09-05 已书面确认同进程；VRev-080 已完成 freshness 记录） |
| **VP-026 / VP-027** | **业务域 VP 激活即触发**评估义务：缓存是否需要（Offer 读取可结论「不需要」）；限流是否已被 VP-030 覆盖。评估不可跳过 |
| **VP-029** | **硬前置**：主体 + 钱包资金原语（至少 freeze/deduct/unfreeze；购买扣款不走凭证核销） |
| **VP-030** | 软前置：无通道时本域仍可经 API 测通；有通道时本模块注册命令。不把 webhook 当本 VP 范围 |
| **VP-033** | 占用位的典型占用者：本模块 Register 之后，033 人工台入口必须隐藏。不把运营台/轮询模式当本 VP 范围 |
| **VP-009 / VP-010** | 购买/核销安全与符合性 gap 归持续程序 |

## 激活记录（VRev-080）

- `consumer_vp`: `VP-031-digital-offer-entitlement`（`vision_ref = schema-ui-core-admin-foundation@0.4.0`）
- `go_issued_at`: 2026-08-10（VP-008 候选 `ed99e88`；消费有效性 2026-08-19 恢复）
- `last_freshness_review_at`: 2026-09-05（VRev-080；写入前 clean HEAD `bd9ed5e062cd965ec4f2221ec5d00351023e76f2`；H-002 同进程书面确认；业务域 freshness PASS）
- `next_freshness_review_trigger`: H-002 主要形态、Profile 默认集、协议 provenance、依赖锁、Offer/权益迁移与 Manifest/装配语义、生产部署/密钥边界变化，或多实例部署触发 RT-Q03/Q05 复核。
- `RT-Q03`: Offer/权益首波以权威存储为正确性来源；本波不需要 Redis。
- `RT-Q05`: Telegram ingress 继续由 VP-030 覆盖；业务端点使用独立进程内请求计数桶；本波不需要 Redis。R1 必须冻结桶 key/阈值/拒绝语义，且不得以 key-wide `Clear` 清除既有历史。

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
| I-031-001 | 首波权益形态：仅时长、仅次数、或二者并存（一 Offer 一种）。 | required | 判据 1/3 | R1 | **verified**（2026-09-05 用户裁决：二者并存，一 Offer 固定一种；workspace-031 GOAL-002 D-001） |
| I-031-002 | 购买状态最小子集（是否要 `pending` 还是同步一拍 fulfilled）。 | required | 判据 2 | R1 | **verified**（2026-09-05 用户裁决：同步一拍 fulfilled，无 pending；GOAL-002 D-001） |
| I-031-003 | 是否允许 Admin 人工发放/撤销权益（客服纠错）。 | required | 判据 3 | R1 | **verified**（2026-09-05 用户裁决：只读 + 作废，不开放人工发放；GOAL-002 D-001） |
| I-031-004 | Telegram 命令清单（若 030 已启用）。 | non-blocking | 判据 5 | R1 | **verified**（默认冻结 price / buy / entitlements；GOAL-002 D-001，用户可否决） |
| I-031-005 | 模块 id（建议 `biz.digital-offer`）。 | non-blocking | 装配 | R1 | **verified**（默认冻结 `biz.digital-offer`；GOAL-002 D-001，用户可否决） |

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| `workspace-031-digital-offer-entitlement` | `GOAL-001-digital-offer-entitlement` | lead | 2026-09-05 | 唯一 delivery；Root `active · 0/4`；R1～R4 尚未开始 |

## 关门记录

| 项 | 值 |
|----|-----|
| status | `closed`（第 3 次关门 2026-09-05 · A-012 `conditional` 2 required 经 D-004 fixed ×2 + A-013 closed ×2 → A-014 independent `pass` 0 required 确认 → F-001 关门前置加固（A-015）→ 重新关门） |
| lead_workspace | `workspace-031-digital-offer-entitlement`（Root `GOAL-001-digital-offer-entitlement` `done · 4/4`；纲领 R1～R4 全部关门） |
| 交付证据 | 判据 1～8 证据矩阵：workspace-031 GOAL-005 `02-execution/E-001-evidence-matrix.md`；代码 `apps/api/modules/digitaloffer/`；合同 `GOAL-002 D-002 v1.2.0` |
| 关门审计 | GOAL-005 `03-audit/A-002`（independent close-out · conditional → F-001/F-002 整改）+ `A-004`（independent closure · **pass · open required 0**）+ `A-005`（关门登记）；子目标审计链：GOAL-002 A-007 / GOAL-003 A-008 / GOAL-004 A-003 |
| freshness | 激活时 VRev-080（2026-09-05）留痕；本轮交付变更了迁移/Manifest/装配——后续任何新工作应先按 `next_freshness_review_trigger` 复审 |
| 边界保持 | 未做类目/SKU/税/库存/物流订单；未进默认 Profile；未解禁通用 Entitlement 接缝；Charter 未改（0.4.0） |

## 规划修订短史

| date | change |
|------|--------|
| 2026-09-02 | 初创 `planned`：用户否决电商类目/商品/订单三件套；确认基座一方数字 Offer+权益；H-002 同进程。业务域分支首个 VP。 |
| 2026-09-03 | 边界指针：Register 后占用 VP-033 人工台入口；运营台不在本 VP。 |
| 2026-09-05 | 用户书面确认 H-002 采用同进程模块；VRev-080 self `pass`（0 required），业务域 freshness PASS；RT-Q03/Q05 评估均结论“本波不需要 Redis”；VP-031 `planned → active` v0.2.0，绑定 `workspace-031-digital-offer-entitlement` 并交 `/govern` 建立 Root。 |
| 2026-09-05 | R1 信息裁决回写（v0.2.1，镜像同步）：I-031-001～003 用户书面裁决（二者并存 / 同步 fulfilled / 只读+作废），I-031-004/005 默认冻结；证据 = workspace-031 GOAL-002 D-001，合同正文 = GOAL-002 D-002。 |
| 2026-09-05 | **关门（v0.3.0）**：R1～R4 全部关门（GOAL-002/003/004/005 done）；Root `done · 4/4`；判据 1～8 证据矩阵见 GOAL-005 E-001；关门审计 GOAL-005 A-002→A-004（independent closure `pass` 0 required）。VP-031 `active → closed`。 |
| 2026-09-05 | **关门撤回（v0.3.1）**：A-006 independent runtime-integration 审计 `fail`（4 required：模块注册表缺失 F-003 / CRUD 语义 F-004 / 迁移策略 F-005 / 组合根验收 F-006）——按 P-003 撤回 Root/GOAL-005/VP-031 关门；裁决与修复见 GOAL-005 D-003，待 closure 复审重新关门。 |
| 2026-09-05 | **重新关门（v0.3.2）**：F-003（BuiltinModules 注册）/ F-005（compiled-global 裁决）/ F-006（组合根验收测试）修复 + F-004 的 E-001 证据分母收窄；经 A-008/A-010 两轮 independent closure 复审，A-010 `pass` 0 required。VP-031 `active → closed`（第 2 次关门）。 |
| 2026-09-05 | **关门撤回（v0.3.3）**：A-012 independent runtime closure re-audit `conditional`（2 required：F-007 迁移 Apply 中途失败/reopen 双方言证据、F-008 Telegram-enabled 真实组合根 + 结构化 Manifest 验收）。主方案确认无需更换。GOAL-005 D-004 fixed ×2（迁移失败/reopen 测试 SQLite+PG；Telegram-enabled 组合根测试 + seam 透传 + DELETE 负向），A-013 self 响应 closed ×2。Root/GOAL-005/VP-031 关门撤回，待 focused independent closure 复审 `pass` 后重新关门。 |
| 2026-09-05 | **重新关门（v0.3.4）**：A-014 independent focused finding-closure `pass` 0 required（F-007/F-008 关闭证据逐条对证成立，真实 PG 未 skip）；用户裁决关门前先落实 F-001 recommended（disabled 分支 Manifest 结构化断言 + dispatcher 指针同一断言，A-015 fixed）；GOAL-005 `done 2/2`、Root `done 4/4`、VP-031 `active → closed`（第 3 次关门）。workspace-031 可作为依赖生产迁移、Telegram 命令与 Admin surface 的后继 VP 的已验证前置（A-012/A-014 口径）。 |
