---
id: GOAL-001-wallet-prepaid-instrument
title: 钱包预付资金凭证与外部主体接缝
status: active
parent: null
created: 2026-09-02
updated: 2026-09-02
version: 0.3.0
progress: 5/5
plan_refs:
  - VP-029-wallet-prepaid-instrument
primary_plan: VP-029-wallet-prepaid-instrument
serves_summary: 钱包预付资金凭证 + 通道无关外部主体接缝（Admin 功能 · 扩展 admin.wallet · 不重开 VP-011 · 不是支付域）：(issuer, external_id) → subject_id（不创建 admin.users）+ 卡密批次生成/导出/作废/核销入账 + R5 Admin 已登录自助核销 HTTP / 我的钱包入口
---

# GOAL-001 · 钱包预付资金凭证与外部主体接缝

## 概述

承接 [VP-029-wallet-prepaid-instrument](../../../vision/plans/VP-029-wallet-prepaid-instrument.md)（**active v0.4.0** · 2026-09-02 用户确认结构选型 A reopen · [VRev-068](../../../vision/reviews/VRev-068-vp029-reopen-my-wallet-self-redeem.md) self `pass`）。R1～R4 首波已交付（当时 v0.3.0 关门 · [VRev-067](../../../vision/reviews/VRev-067-vp029-wallet-prepaid-instrument-close-out.md) 原文不改写）。**现行对象面**：外部主体接缝 + 预付资金凭证入金通道 + **R5 Admin 已登录自助核销 HTTP / 「我的钱包」入口**。**红线（激活即生效）**：不重开 VP-011；不把 C 端用户做成 `admin.users`；不引入支付网关或 Telegram 依赖；不改 Profile 默认集 / 模块矩阵 / Manifest 装配（VP-008 `go` 消费有效性）；不消耗 RT-Q03/Q05 trigger；Admin 自助核销不得记入 `owner_type=subject` 账。

## 成功标准（对应 VP-029 方向级退出判据；#1～#7 历史核销，#8～#10 为 R5）

- [x] 判据 #1（主体接缝可用）：`(issuer, external_id) → subject_id` 幂等登记/查找有测试；未登记主体不能开户；不创建 `admin.users`；查询/get-or-create 不依赖 `admin.wallet` 已启——R1/R2
- [x] 判据 #2（凭证生命周期）：生成（哈希存储 + 一次性明文）/ 导出 / 作废 / 过期拒绝有测试；明文不落库、不进审计原文——R1/R3
- [x] 判据 #3（核销原子且幂等）：Redeem 成功则账本入金与凭证状态一致；重复核销不双记；并发双花 fail-closed——R1/R2
- [x] 判据 #4（账本不变式保持）：三余额恒等、流水快照链、对账 Job 仍通过既有钱包测试；新入金类型纳入 apply 表——R2
- [x] 判据 #5（Admin 可操作）：批次生成/导出/作废有协议驱动页面 + 权限键 + 操作审计——R3
- [x] 判据 #6（边界保持）：未改 Charter；未改 `mvp`/`admin` 默认模块集装配语义；未引入支付网关或 Telegram 依赖；未重开 VP-011——全程
- [x] 判据 #7（审计闭合）：开放 required finding = 0（或已合法闭合）——R4 当时分母已核销；R5 仍适用
- [x] 判据 #8（Admin 已登录自助核销 HTTP）：会话身份入账 `owner_type=user`；禁止匿名；原子且幂等；不记入 subject 账——R5（GOAL-005 done）
- [x] 判据 #9（我的钱包入口）：`/my-wallet` 可发现预付凭证充值入口；成功刷新余额/流水——R5（GOAL-005 done）
- [x] 判据 #10（限流评估落盘）：对本 authenticated HTTP 完成 RT-Q05 精神评估——R5（GOAL-005 D-002 / I-029-008 closed）

## 纲领路线图（P-001）

阶段串行；同一阶段内可并行子目标。

| 阶段 | 内容 | 检查点/状态 |
|------|------|-------------|
| R1 | 合同冻结（判据 1/2/3/5 边界）：主体落点 + `owner_type` + `OwnerExists`（I-029-001）· 哈希/熵/常时比较/UNIQUE+同事务（I-029-006）· `entry_type`（I-029-002）· 权限键（I-029-003）· HTTP 核销是否本波（I-029-005） | **done**（D-002 裁决冻结 · 2026-09-02） |
| R2 | 主体接缝 + 账本入金（判据 1/3/4）：幂等 get-or-create · Redeem 原子入金 · 三余额/对账回归 | **done**（GOAL-002 完成 · 2026-09-02 交叉审计 pass 关门） |
| R3 | Admin 批次面 + 导出（判据 2/5 + I-029-004）：生成/导出/作废/查询 · 权限键 · 操作审计 | **done**（GOAL-003 完成 · 2026-09-02 交叉审计 pass 关门） |
| R4 | 证据与关门（判据 6/7；依赖 R1–R3） | **done**（GOAL-004 完成 · 2026-09-02 独立关门审计闭合；当时分母） |
| R5 | Admin 已登录自助核销 HTTP + 「我的钱包」入口（判据 8/9/10 + I-029-007/008/009） | **done**（GOAL-005 done 4/4 · A-001 independent + A-003 self；Root 仍 active） |

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-029-001 | required | 主体落点：薄模块 vs `authsession` vs `admin.wallet` 表。须同时冻结 `owner_type` 是否新增取值；`OwnerExistsFunc` 改为已登记主体、禁止回退 `UserByID`；查询/get-or-create 不依赖钱包已启。W13 F-012 孤儿账本相对主体登记表。 | 方案冻结 + 判据 1 | R1 | 用户裁决（R1 合同冻结前置） | closed | — | D-002：独立主体表 subjects，owner_type 扩充 subject，OwnerExists 校验主体（V-F112 closed） |
| I-029-002 | required | 核销入金的 `entry_type`：新类型 vs `adjust` + `ref_type=voucher`。 | 判据 3/4 | R1 | 用户裁决（R1 合同冻结前置） | closed | — | D-002：复用 adjust + ref_type='voucher'，不改动已有 CHECK 与 Apply |
| I-029-003 | required | 生成权限键：复用 `wallet.adjust` vs 新 `wallet.voucher.issue`。 | 判据 5 | R1 | 用户裁决（R1 合同冻结前置） | closed | — | D-002：新增细粒度 wallet.voucher.issue 权限键 |
| I-029-004 | non-blocking | 导出格式（CSV/TXT、是否含明文、一次性下载）。 | 判据 5 | R3 | lead 建议 + 用户确认 | closed | — | D-001：API 一次性返回明文数组作为导出数据源（明文不入库） |
| I-029-005 | non-blocking | C 端自助核销 HTTP 是否**首波**交付，或仅模块 API。若选 HTTP，本 VP 必须完成 RT-Q05 精神的限流评估，不得推到未激活的 VP-030。 | 当时判据面 | R1 | 用户裁决 | closed | — | D-002：首波仅模块 API。**不**表示 R5 禁止 HTTP（见 I-029-007/008） |
| I-029-006 | required | 凭证哈希与双花合同：哈希算法（默认候选 = 高熵码 SHA-256 或 HMAC-SHA256+pepper；禁止 6 位恢复码或 bcrypt 当卡密默认）；码字母表与长度（熵下限）；核销常时比较；`UNIQUE(code_hash)`（或等价）+ 同事务「未用→已核销 AND 账本入金」，并发失败者 fail-closed，重复 Redeem 不双记。 | 判据 2/3 | R1 | 用户裁决（R1 合同冻结前置） | closed | — | D-002：高熵码 + SHA-256 + 单事务 CAS 原子核销入金，并发失败 fail-closed |
| I-029-007 | required | R5 HTTP 路径与服务函数形状（默认候选 `POST /api/wallet/me/redeem`） | 判据 8 · GOAL-005 S1 | R5 S1 | GOAL-005 D-002 | closed | — | `POST /api/wallet/me/redeem` + `RedeemForUser` |
| I-029-008 | required | 已登录核销 HTTP 的 RT-Q05 精神限流（内存专用桶 vs 复用 vs residual） | 判据 10 · GOAL-005 S1 | R5 S1 | GOAL-005 D-002 | closed | — | 内存专用桶 15min/10/user id；不消耗 Redis trigger |
| I-029-009 | required | 自助核销权限：identity-only vs 新键 vs `wallet.voucher.issue` | 判据 8 | R5 重开 | 用户确认结构选型 A + D-003 | closed | — | D-003：identity-only，与 `/api/wallet/me` 同款 |

## 父目标

- `null`（Root）

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；D-001/E-001/A-001 已首条落盘。

## 备注

- 审计模式（D-001 已定）：阶段关门 default self；实证门禁（R2 核销并发 / R4 证据 / 关门）已通过 grok build 独立审计（A-001 independent pass）。**R5**：S4 关门 independent（资金路径）；不得静默降级。
- 2026-09-02 reopen（D-003 / E-010）：不回退 R1～R4 子目标 `done` 与 A-001～A-008 原文。
