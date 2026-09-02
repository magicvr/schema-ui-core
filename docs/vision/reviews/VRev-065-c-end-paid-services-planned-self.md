---
id: VRev-065-c-end-paid-services-planned-self
doc_type: vision-review
title: VP-029/030/031 计划审视 · C 端付费服务触发（钱包卡密 / Telegram 通道 / 数字 Offer）
source: self
date: 2026-09-02
scope: VP-029-wallet-prepaid-instrument / VP-030-telegram-channel-runtime / VP-031-digital-offer-entitlement（planned · 意图与退出判据）
verdict: pass
open_required: 0
status: active
created: 2026-09-02
updated: 2026-09-02
parent: null
version: 0.1.0
---

# VRev-065 · VP-029/030/031 计划审视（C 端付费服务）

## 背景与触发

用户 2026-09-02 确认：消费本仓基座的下游仓将通过 Telegram Bot 向 C 端提供付费服务。结构选型：

- 基座做**可复用一方能力**，不是把该 Telegram 产品写进 Charter 成功条件
- **不要**电商类目 / 商品 / 物流订单三件套
- 切分：钱包卡密 = Admin 功能（扩展 `admin.wallet`）；Telegram = 架构通道运行时；数字 Offer+权益 = 业务域（本仓首个业务域 VP）
- H-002 **同进程**保持
- 外部主体接缝：`(issuer, external_id) → subject_id`，不创建 `admin.users`

三 VP 均 `planned`、0 区、`vision_ref` = `schema-ui-core-admin-foundation@0.4.0`。本审视覆盖意图、退出判据、非目标、P-005、相邻边界与组合索引同步。本审视**不是**激活门禁（激活前仍须 freshness + 激活 VRev）。

## 审视要点

### 1. 与 Charter 0.4.0 对齐

**pass**。非目标「不建设特定业务领域的终端产品」「不预制 C 端 API 的业务逻辑」仍然成立：三 VP 交付的是可装配模块与端口，不是某个 Bot 的运营终态。成功边界 #6 要求业务域模块能消费基础设施端口——VP-031 正是该边界的第一名消费者。H-002 同进程已由用户书面确认，发现机制（业务域激活前再确认）写进 VP-031 判据 6。

将「服务视为商品」收窄为 Offer 而非 Catalog/SKU，避免把 Charter 候选域 1/3 整包立项，符合「一域一 VP」与「登记 ≠ 立项」。

### 2. 三 VP 切分与关门独立性

**pass**。

| VP | 分支 | 关门能力 | 触发 |
|----|------|----------|------|
| 029 | Admin 功能 | 主体接缝 + 凭证入金 | 卡网售卖密钥、C 端身份不能进 Admin 用户表 |
| 030 | 架构 · C 端通道 | webhook/分发/出站/身份映射 | 同进程 Bot；即使无 Offer 也有 ingress |
| 031 | 业务域 | Offer + 薄购买 + 本域权益 | 真实付费服务；依赖 029 资金原语 |

合并三者会让通道运行时被业务域拖死，或把支付网关语义渗进钱包。切分与 VP-026/027/028「触发独立 × 关门独立」同构。建议激活序 029→030→031 与硬前置一致；三分支并行规则允许 029 closed 后 030 与 031 各一 active，031 仍硬前置 029。

### 3. 退出判据可判定性

**pass**。

| VP | 判据可判定性 |
|----|-------------|
| 029 | 主体幂等、凭证哈希、Redeem 原子/并发、账本不变式、Admin 页面均可测；明文不落库可核验 |
| 030 | secret fail-closed、Register 分发、SendMessage mock、主体映射、设置密钥、限流评估落盘均可核验 |
| 031 | Offer CRUD、freeze/deduct 同事务、权益三态、subject 对齐、命令可选、H-002/RT-Q 评估留痕均可核验 |

### 4. 非目标与反模式

**pass**。明文排除：类目树、SKU/税/库存、物流订单、Stars/支付网关、Mini App、对话 FSM、把 Bot 用户写入 `admin.users`、重开 VP-011、解禁通用 Entitlement 接缝。钱包卡密未冒充「支付域成立」。S-13 未被本波「顺便交付」。

### 5. 信息需求（P-005）

**pass**。主体落点、entry_type、权限键、无 token 启动策略、Bot HTTP vs SDK、限流桶、权益形态、购买状态机均 required 且最晚 R1，未伪装已决。

### 6. 组合编排同步

**pass**。roadmap 已加 29–31 行、RT-M03、RT-Q05 C 端评估前移、Admin 当前拍、业务域已立项表；revisions VR-059；workspaces 无需更新（0 区）。

### 7. 激活门禁（本审视不放行激活）

VP-029：Admin 类 freshness。VP-030：架构类 freshness + 限流评估。VP-031：业务域 freshness（含 H-002 再确认）+ RT-Q03/Q05 评估。三份均未激活、未开区。

## Verdict

**pass**

意图清晰、切分可关门、与 Charter 0.4.0 可调和、索引已同步。无 required finding。不构成激活许可。

## Findings

### 必改（required）

无。

### 建议（recommended）

| id | finding | 建议 | 状态 |
|----|---------|------|------|
| V-F108 | 「薄购买凭证」易被读成 VP-011 S-13 订单已立项。 | VP-031 边界表写明「不是 S-13 替代立项」。 | **fixed**（本轮写入 VP-031 相邻 VP 表） |
| V-F109 | 若主体接缝实现绑死在已启用的 `admin.wallet`，无钱包的权益/通道会断。 | 判据 1 要求查询/get-or-create 不依赖钱包已启。 | **fixed**（本轮写入 VP-029 判据 1） |
| V-F110 | 本波首次同时引入 C 端主体与资金入金，计划阶段仅 self Review。 | 激活 VP-029 前做一次 independent Vision Review（资金哈希、并发双花、主体与 Admin 用户隔离）。 | open（激活前执行，不阻断 `planned`） |

## 建议下一步

1. **不要**立刻 `/govern` 开区。先激活 VP-029：Admin 类 freshness +（建议）independent VRev + 用户确认 `planned → active`，再交 `/govern` scaffold `workspace-029-wallet-prepaid-instrument`。
2. VP-030 / VP-031 保持 `planned`，直至 029 主体接缝（030）/ 资金原语（031）可用。
3. V-F110 在 029 激活审视时核销。

## Finding 响应（2026-09-02 · `/vision` 计划落盘事务内）

| id | 路径 | 说明 |
|----|------|------|
| V-F108 | **fixed** | VP-031「与相邻 VP」表增加：本 VP 不是 S-13 替代立项 |
| V-F109 | **fixed** | VP-029 退出判据 1 增加：主体查询/get-or-create 不依赖 `admin.wallet` 已启 |
| V-F110 | 保持 open | 激活 VP-029 前独立 Vision Review；不阻断本批 `planned` |

## 声明

本意见不直接修改 Charter / VP / Goal status。required finding 的响应由 `/vision` 追加在本报告中；原 verdict 与 finding 原文不得改写。
