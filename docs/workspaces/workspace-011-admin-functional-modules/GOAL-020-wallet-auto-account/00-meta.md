---
id: GOAL-020-wallet-auto-account
title: 钱包账户自动开户与用户绑定（get-or-create）
status: done
parent: GOAL-001-admin-functional-modules
created: 2026-08-16
updated: 2026-08-16
version: 0.2.0
progress: 5/5
---

# GOAL-020-wallet-auto-account · 钱包账户自动开户与用户绑定

## 概述

GOAL-019（S-14 钱包/账务）关门后用户反馈（2026-08-16）的产品演进：钱包账户应由**系统按用户自动创建并绑定**，而非管理员手动开户。经用户确认：**惰性 get-or-create**（按 ownerId 查询或调账时自动创建零余额账户，UNIQUE 约束幂等兜底）；**user 类型禁止手动创建**（保留 business/system 手动入口）。

## 当前边界

- 触发面：新增 by-owner 读端点（get-or-create）与 by-owner 调账端点（自动开户 + 调账）；现有 accountId 端点与数据兼容不动。
- 自动开户：owner_type=user 固定、owner_id=传入 ownerId、currency 默认 CNY；并发 INSERT 冲突 → 重读（UNIQUE(owner_type, owner_id, currency) 兜底）。
- 审计：自动开户同样记录 wallet.account-create（detail 标 auto）；调账沿用 wallet.adjust。
- 手动边界：POST /api/wallet/accounts 拒绝 owner_type=user（409 WALLET_USER_AUTO_ONLY）；business/system 手动入口保留。
- 前端：钱包页"新建账户"表单移除 user 选项（business/system 保留）；不新增主动开户 UI（调账即开户，惰性语义）。
- 无越界：不改变装配语义/协议 pin；跨模块零耦合（不依赖 users 模块事件）。

## 成功标准与路线图（P-001）

- [x] **S1 · 方案冻结**：by-owner 端点形态、get-or-create 幂等与并发、审计形态、手动边界错误码（D-001，2026-08-16）
- [x] **S2 · 实现**：store GetOrCreate + handler by-owner 端点 + 手动拒绝 + 前端选项（E-002，2026-08-16）
- [x] **S3 · 验证**：单元/集成 + 全量回归（go 全绿 / web 1004/1004）（E-002，2026-08-16）
- [x] **S4 · go 影响判定 + 自审**（D-002 不暂挂 + A-002 pass）（2026-08-16）
- [x] **S5 · 关门**：独立审计（A-003 conditional → F-001/F-002 fixed → A-004/A-005 pass）+ 关门 + goal-tree 同步（E-003，2026-08-16）

progress: 5/5 由五个等权检查点派生（S5 关门后更新）。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 自动开户触发面与幂等：by-owner 读/写端点形态、并发 UNIQUE 兜底、审计事件 detail 形态 | S1 方案 | S1 | 对照 GOAL-019 D-002 §3 端点 + store Mutate/乐观锁先例 | **verified** | — | D-001 §1（2026-08-16） |
| I-002 | required | 手动创建边界：user 类型禁止的语义与错误码（409 码名） | S1 方案 | S1 | errorcatalog/error_contract 先例 | **verified** | — | D-001 §2（2026-08-16；WALLET_USER_AUTO_ONLY） |
| I-003 | non-blocking | 前端调整：ownerType 选项移除 user 的副作用（既有 user 账户展示不受影响） | S1 方案 | S1 | wallet.json options + schema-keys 分母 | **verified** | — | D-001 §3（2026-08-16） |

## 审计策略

data 门禁（自动开户涉及资金账户创建）：S5 关门必须 grok build independent（用户书面偏好 grok-4.6 · high，沿用 GOAL-019）；本目标为小型演进，方案 + 实现合并为一次关门独立审（P-002 小目标合并审视）。

## 父目标

- [GOAL-001-admin-functional-modules](../GOAL-001-admin-functional-modules/00-meta.md)

## 台账布局

本目标从首条记录起使用 01-decision/、02-execution/、03-audit/ 平铺 ledger；索引与目录条目共同构成正式记录。