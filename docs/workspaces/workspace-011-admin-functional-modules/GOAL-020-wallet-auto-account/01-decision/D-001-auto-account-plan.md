---
id: D-001
goal: GOAL-020-wallet-auto-account
title: 方案冻结：自动开户 get-or-create + 手动边界
date: 2026-08-16
status: accepted
parent: GOAL-020-wallet-auto-account
created: 2026-08-16
updated: 2026-08-16
version: 1.1.0
---

# D-001 · 方案冻结（自动开户与用户绑定）

> 依据：GOAL-020 00-meta 边界与 I-001~I-003；用户裁决（2026-08-16：惰性 get-or-create / user 手动禁止 / 新目标承载）。
> 证据：GOAL-019 D-002 §1/§3（账户 UNIQUE 约束、apply 表、端点形态）、store/repository.go（Mutate 乐观锁、CreateAccount UNIQUE 冲突）、handler/wallet.go（端点与审计）、errorcatalog/error_contract（错误码契约）。

## 1. 自动开户触发面（I-001 闭合）

- **读端点** `GET /api/wallet/by-owner/{ownerId}`（wallet.read；A-003 F-004 勘误：`/api/wallet/accounts/by-owner/...` 与既有 `GET /api/wallet/accounts/{id}/entries` 在 Go 1.22 ServeMux 模式重叠，改 `by-owner` 前缀）：get-or-create——账户不存在则自动创建零余额 user 账户（owner_type=user、owner_id=ownerId、currency=CNY）并返回；已存在直接返回。
- **写端点** `POST /api/wallet/by-owner/{ownerId}/adjust`（wallet.adjust；同勘误）：自动开户 + 调账（get-or-create 后走既有 Mutate，apply 表语义不变）；返回 account + entry。
- **幂等与并发**：创建走 UNIQUE(owner_type, owner_id, currency)；并发 INSERT 冲突 → 重读已有账户（get-or-create 语义，事务内实现）。调账幂等沿用 idempotencyKey。
- **审计**：自动开户记录 `wallet.account-create`（detail 含 `"auto":true` 与 ownerId）；调账沿用 `wallet.adjust`。
- 冻结/解冻不提供 by-owner 变体（账户经读端点开户后，既有 accountId 端点即可操作；控制面保持最小）。

## 2. 手动创建边界（I-002 闭合）

- `POST /api/wallet/accounts` 拒绝 owner_type=user → 409 **WALLET_USER_AUTO_ONLY**（"user wallet accounts are created automatically"）；business/system 手动创建保留。
- 新错误码进 errorcatalog（en/zh + messageKey）与 error_contract 冻结集。

## 3. 前端调整（I-003 闭合）

- wallet.json「新建账户」表单 ownerType select 移除 user 选项（保留 business/system）；schema.wallet.owner.user 键保留于 catalog（无害，移除会破坏双语键集一致性）。
- 不新增主动开户 UI：惰性语义下调账/读即开户；管理员可经 API 或未来页面扩展按用户开户。

## 4. 测试与验证

- store：GetOrCreateUserAccount（新建 / 已存在 / 并发冲突重读）。
- handler：by-owner 读端点（自动开户 + 审计 auto 标记）、by-owner 调账（自动开户 + apply）、user 手动创建 409、既有端点回归。
- 错误码：WALLET_USER_AUTO_ONLY 双语 + 冻结集。
- 组合根/迁移：无迁移变更（复用 0031 表）；快照不变（权限/导航不变——不新增权限键）。
- web：schema-keys 分母回归（user 选项移除不影响键集）+ 全量回归。

## 5. 未选方案（留痕）

- 用户注册钩子（users.create 联动开户）：跨模块耦合 + 存量回填迁移，不采纳（用户裁决惰性机制）。
- 完全移除手动创建：business/system 无自动机制，保留人工入口（用户裁决）。
- by-owner freeze/unfreeze 变体：控制面最小化，既有 accountId 端点可覆盖，不采纳。
- 显式「开户」操作端点：get-or-create 已覆盖语义，避免多余写面。

## 6. S2 实现清单

1. store：`GetOrCreateUserAccount(ownerID string, now) (*Account, error)`（WithTx 内 SELECT → INSERT → 冲突重读）。
2. handler：新增 by-owner 读/调账端点（复用 walletMutate 骨架 + auto 审计 detail）；POST /api/wallet/accounts 拒绝 user（409 WALLET_USER_AUTO_ONLY）。
3. errorcatalog + error_contract：WALLET_USER_AUTO_ONLY。
4. wallet.json：ownerType options 移除 user。
5. 测试：store/handler/错误码/web 回归。