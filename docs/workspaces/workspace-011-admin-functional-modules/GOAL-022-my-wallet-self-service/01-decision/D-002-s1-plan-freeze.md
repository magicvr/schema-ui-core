---
id: D-002
goal: GOAL-022-my-wallet-self-service
title: 方案冻结：只读自服务 + /my-wallet 惰性开户
date: 2026-08-16
status: accepted
parent: GOAL-001-admin-functional-modules
created: 2026-08-16
updated: 2026-08-16
version: 1.0.0
---

# D-002 · 方案冻结（S1）

> 依据：00-meta I-001/I-002；用户裁决（2026-08-16，P-004）：**只读自服务** + **/my-wallet 惰性开户**；GOAL-020 D-001（by-owner get-or-create 语义）；GOAL-013 D-002（navigation.user 槽位）。

## 1. 范围（I-001 闭合 · verified）

- 当前登录用户**只读**查看自己的钱包：余额（总/可用/冻结）+ 自己的流水。**无任何资金操作面**（不提供调账/冻结/解冻/扣减的自服务入口）。
- 禁止调他人账：新端点不接受 ownerId 入参，账户恒由会话身份推导；返回账户 ownerId 必须等于会话用户 id。
- 复用 GOAL-019 账本、GOAL-020 get-or-create；不另起账本、不新增迁移（账户表 0031 复用）。

## 2. 路由与身份（I-002 闭合 · verified）

| 项 | 值 |
|----|----|
| pageId / route | `my-wallet` / `/my-wallet`（独立页面，非个人中心子页） |
| 新端点 | `GET /api/wallet/me`（身份作用域 get-or-create + 账户摘要，resourceList 信封：items=[account]）；`GET /api/wallet/me/entries`（身份作用域流水分页） |
| 权限模型 | **identity-only**（无权限键，与 /api/account/profile 同款）：任意已认证用户仅可读自己；管理端 /wallet（wallet.read）不动 |
| 开户时机 | 两个 /me 端点均惰性 get-or-create（GOAL-020 D-001 §1 语义）；首次创建审计 `wallet.account-create`（detail 含 `"auto":true`） |
| 导航 | `navigation.user` 新节点 `menu_wallet_self`（icon wallet），排序在个人中心（menu_account）之后、设置（menu_settings）之前；可见性 PolicyAdminEditorViewer |

## 3. 页面（schema 驱动）

- `my-wallet.json`：intro text + 3 张 statCard（总余额/可用/冻结，valueField=balanceTotal/balanceAvailable/balanceFrozen，dataSource=/api/wallet/me）+ 流水 table（data.url=/api/wallet/me/entries；列复用 walletEntries 既有 labelKey）。
- 金额沿用系统约定整数 min-units（F-011 前端金额格式化仍 deferred），label 注明单位。

## 4. 与管理端页分工 / 装配

- 管理端 `/wallet`（wallet.read）全量账户管理不变；本页只读自己的，两页并存互不干扰。
- Manifest fragment：pages 增 my-wallet；navigation.user 增 menu_wallet_self 项；DefaultNavigationOrder 插入 `menu_wallet_self`（menu_account 之后、menu_activity 之前）。
- Descriptor/BuiltinModules 路由声明同步；权限三键不增；无新 capability；协议 pin/Profile 默认集/装配语义不变。

## 5. 测试与验证

- handler：/me 401；开户（auto 审计标记 + 信封）；幂等（不重复开户/审计）；entries 分页与本人数据；身份隔离（返回账户 ownerId == 会话用户）。
- kernel：DefaultNavigationOrder 快照更新。
- web：admin fixture 增页/导航 + 哈希重钉；app-manifest 页面清单；schema-keys 双语键；D-VAL 全模块。
- 全量回归：go test ./... + vitest。

## 6. 审计策略

- 只读（无资金操作面）→ S2/S3 以 **self** 审计；S5 关门按用户偏好安排 **grok build independent**（grok-4.6 · high）核验身份隔离与数据暴露边界。

## 7. 未选方案（留痕）

- 有界自助操作（冻结/解冻自己的金额）：用户未选；如未来启用将升级 data 门禁并独立审计。
- `/account/wallet` 子页：用户未选；独立页利于顶栏入口与后续扩展。
- 无账户空态（不 get-or-create）：用户未选；惰性开户与 GOAL-020 语义一致，页面零特殊态。
