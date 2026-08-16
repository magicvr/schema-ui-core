---
id: E-002
goal: GOAL-020-wallet-auto-account
title: S2 实现 + S3 验证完成（get-or-create + 手动边界）
date: 2026-08-16
status: recorded
parent: GOAL-020-wallet-auto-account
created: 2026-08-16
updated: 2026-08-16
version: 1.0.0
---

# E-002 · S2 实现 + S3 验证（2026-08-16）

## 事实

- **store**：`GetOrCreateUserAccount(ownerID, now) (*Account, bool, error)`——WithTx 内 SELECT → 不存在则 INSERT（UNIQUE 兜底，并发冲突重读）；bool 报告本次是否新建（审计 auto 标记的依据）。
- **handler**：`GET /api/wallet/by-owner/{ownerId}`（wallet.read；自动开户，新建时记 wallet.account-create + auto:true）；`POST /api/wallet/by-owner/{ownerId}/adjust`（wallet.adjust；自动开户 + 调账，account-create + adjust 双审计）；`POST /api/wallet/accounts` 拒绝 owner_type=user → 409 WALLET_USER_AUTO_ONLY（business/system 保留）。
- **路由形态修正（实现期）**：初版 `/api/wallet/accounts/by-owner/{ownerId}` 与既有 `GET /api/wallet/accounts/{id}/entries` 在 Go 1.22 ServeMux 上模式重叠冲突（交集 .../by-owner/entries）→ 改为无冲突的 `/api/wallet/by-owner/...` 前缀（kernel RouteKey/BuiltinModules 同步）。
- **错误码**：WALLET_USER_AUTO_ONLY 进 errorcatalog（en/zh + messageKey）与 error_contract 冻结集 + web 双语 catalog。
- **前端**：wallet.json 新建账户表单 ownerType 选项移除 user（business/system 保留）。
- **测试**：store（新建 created=true / 二次 created=false / 自动账户可正常 Mutate）；handler（by-owner 读自动开户 + 单次 auto 审计断言、by-owner 调账自动开户 + apply、user 手动 409 + business 手动 201、门禁 401/403 覆盖新路由）；GOAL-019 旧测试改用自动开户。
- **验证**：go 全量 + web 全量回归全绿（pwsh-402/403）。
