---
id: workspace-029-wallet-prepaid-instrument
title: 钱包预付资金凭证与外部主体接缝工作区（Admin 功能）
status: active
root_goal: GOAL-001-wallet-prepaid-instrument
canonical_scope: docs/workspaces/workspace-029-wallet-prepaid-instrument/
shared_materials_catalog: none
vision_role: delivery
plan_refs:
  - VP-029-wallet-prepaid-instrument
primary_plan: VP-029-wallet-prepaid-instrument
created: 2026-09-02
updated: 2026-09-02
version: 0.1.0
parent: null
---

# 工作区上下文 · 钱包预付资金凭证与外部主体接缝

本工作区是 [VP-029-wallet-prepaid-instrument](../../vision/plans/VP-029-wallet-prepaid-instrument.md)（**`closed`** v0.3.0 · 2026-09-02 用户书面确认关门 · VRev-067 self `pass`）的唯一 lead delivery workspace。**Admin 功能分支**：扩展已交付的 `admin.wallet`——通道无关外部主体接缝 `(issuer, external_id) → subject_id`（不创建 `admin.users`）+ 预付资金凭证（批次生成/导出/作废/核销入账，哈希存储、幂等 Redeem）。**不是**支付/结算业务域。

- **Root** `GOAL-001-wallet-prepaid-instrument`：`done` · **4/4**（R1 合同冻结 [done] → R2 主体接缝+账本入金 [done] → R3 Admin 批次面+导出 [done] → R4 证据与关门 [done]），纲领见 Root `00-meta.md`。
- 激活门禁已满足（2026-09-02）：[VRev-066](../../vision/reviews/VRev-066-vp029-wallet-prepaid-instrument-independent.md) independent `pass`（0 required；V-F111/112/113 → 开区事务内 fixed；V-F110 核销）；**Admin 类轻量 freshness PASS**（`29727510` → `b5c39dfb`：协议 pin / 依赖锁 / 迁移台账 / Profile 装配 / provenance 五域零变更；区间代码 = VP-028 已审结目 + VP-009 W16/W17）不暂挂 `go`。
- 不改变 Charter `primary_workspace`（仍为 workspace-001）。
- **消费基线**：VP-011 `admin.wallet` 账本原语（`adjust` / `freeze` / `unfreeze` / `deduct_frozen`、三余额恒等、不可变流水、幂等键、对账）只读继承；VP-012 并发/幂等/审计；VP-003/004 模块 Persistence + 全局迁移台账。
- **红线（激活即生效）**：不重开 VP-011；不把 C 端用户做成 `admin.users`；不引入支付网关或 Telegram 依赖；不把新模块塞进 `mvp`/`admin` 默认集；不消耗 RT-Q03/Q05 trigger；不解除 typed domain event gated。

## 绑定

| 字段 | 当前值 | 说明 |
|------|--------|------|
| 工作区 ID | `workspace-029-wallet-prepaid-instrument` | 与本区目标及资料引用的 `workspace_id` 一致 |
| Root Goal | `GOAL-001-wallet-prepaid-instrument` | `parent: null`；done · 4/4 |
| canonical 范围 | `docs/workspaces/workspace-029-wallet-prepaid-instrument/` | 本区唯一目标状态范围 |
| 共享资料目录 | `none` | 暂无固定共享资料 |
| 愿景角色 | `delivery` | VP-029 lead（closed v0.3.0）；不改变 Charter primary workspace |
| 规划对齐 | `primary_plan` = `VP-029-wallet-prepaid-instrument`（`closed` v0.3.0） | 2026-09-02 关门（VRev-067 self `pass`；A-004 independent `pass`） |

## 愿景对齐

Charter：`schema-ui-core-admin-foundation@0.4.0`（2026-08-31 strategic：同进程基座 · 成功边界 #6 · H-002）。
VP-029：钱包预付资金凭证与外部主体接缝（vision_ref @0.4.0）——七条方向级退出判据 = 主体接缝 / 凭证生命周期 / 核销原子且幂等 / 账本不变式 / Admin 可操作 / 边界保持 / required=0 闭合；红线 = 不重开 VP-011、不改 Profile 默认集/Manifest、不引入支付网关或 Telegram、不把 Bot 用户写入 `admin.users`。

## 纲领阶段（Root 路线图指针）

| 阶段 | 内容 | 状态 |
|------|------|------|
| R1 | **合同冻结**（判据 1/2/3/5 + I-029-001/002/003/006 裁决）：主体落点与 `owner_type` / `OwnerExists` · 哈希/熵/常时比较/UNIQUE+同事务 · `entry_type` · 权限键 | **done**（D-002 裁决冻结） |
| R2 | **主体接缝 + 账本入金**（判据 1/3/4）：幂等 get-or-create · Redeem 原子入金 · 三余额/对账回归 | **done**（GOAL-002 完成关门） |
| R3 | **Admin 批次面 + 导出**（判据 2/5 + I-029-004）：生成/导出/作废/查询 · 权限键 · 操作审计 · 明文不进审计原文 | **done**（GOAL-003 完成关门） |
| R4 | **证据与关门**（判据 6/7）：证据矩阵 / 越界核账 / 审计闭合 | **done**（GOAL-004 关门 · Root GOAL-001 done 4/4） |

## 固定共享资料引用

> `shared-materials/index.json` 只能提供候选路径与摘要。缺完整引用字段的行无效。

| reference_id | workspace_id | material_id | source | version | sha256 | purpose | local_record | status |
|--------------|--------------|-------------|--------|---------|--------|---------|--------------|--------|
| — | — | — | — | — | — | — | none | — |
