---
id: D-003
goal: GOAL-022-my-wallet-self-service
title: S4 go 判定（只读加法面，无门禁语义变化）
date: 2026-08-16
status: accepted
parent: GOAL-001-admin-functional-modules
created: 2026-08-16
updated: 2026-08-16
version: 1.0.0
---

# D-003 · S4 go 判定

## 决定

1. **放行 S5 关门**：S2/S3 证据充分（E-003/E-004），无开放必改 finding、无到期 required 信息项（I-001/I-002 verified）。
2. **go 影响判定 = 不加不减**：
   - 权限三键（wallet.read/write/adjust）未增未改；无新 capability；协议 pin（schema-ui-docs@v2.8.0）与 Profile 默认集 / Manifest 装配语义未变；
   - 新增路由为 **加法**：`GET /api/wallet/me`、`GET /api/wallet/me/entries`（identity-only，不参与权限矩阵）；
   - 无迁移（复用 0031 表）；系统数据仅新增 menu_items 1 行（menu_wallet_self），不影响既有授权；
   - 前端：manifest 增加一页一项，renderer/内核零改动。
3. 审计模式：只读（无资金操作面）→ S4 自审（A-001）+ S5 grok independent（用户偏好 grok-4.6 · high）。

## 未选方案

- 不安排独立审计直接关门：wallet 数据暴露边界（身份隔离）值得一次独立核验，且用户已提供 grok build 渠道。
- 把 /me 挂到管理端权限键：会破坏「任意已认证用户自服务」语义，与 D-002 §2 冲突。