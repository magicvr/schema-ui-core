---
id: E-005
goal: GOAL-019-r3-s14-wallet-ledger
title: S2 实现完成（migration 0031/0032 + modules/wallet + handler + 组合根 + web）
date: 2026-08-16
status: recorded
parent: GOAL-019-r3-s14-wallet-ledger
created: 2026-08-16
updated: 2026-08-16
version: 1.0.0
---

# E-005 · S2 实现（2026-08-16）

## 事实

- **迁移**：0031（admin.wallet 三表：wallet_accounts / wallet_ledger_entries / wallet_reconciliation_runs；三余额恒等式 CHECK + 复合幂等 UNIQUE + 快照恒等式 CHECK + 索引）+ 0032（operationlog 事件 CHECK 超集重建）；compiled/persistence.go 注册。
- **store**（modules/wallet/store）：ListAccounts/GetAccount/CreateAccount（唯一持有方）/UpdateStatus（乐观锁）/Mutate（apply 表 + 幂等复合键 + 快照链）/ListEntries/ReconcileRun（恒等式 + 链重放）/ListReconcileRuns；NULL 安全扫描。
- **链序修正（实现期发现）**：provider 生成**毫秒时间序 id**（UnixMilli 前缀 + 随机后缀），保证同秒流水按 D-002 §1 (created_at ASC, id ASC) 重放即创建序——handler 测试暴露的同秒乱序问题已修复。
- **handler**（WalletRoutes）：10 路由（accounts CRUD/状态、entries 路径+查询双变体、adjust/freeze/unfreeze、reconcile、runs）；错误码 WALLET_NOT_FOUND / WALLET_OWNER_TAKEN / WALLET_DISABLED / INSUFFICIENT_BALANCE / LEDGER_VERSION_CONFLICT / LEDGER_IDEMPOTENCY_CONFLICT / INVALID_LEDGER_ENTRY / INVALID_WALLET_BODY 全部进 errorcatalog 双语 + error_contract 冻结集。
- **组合根**：composition.go 接线 + BuiltinModules + ProfileAdmin += admin.wallet + DefaultNavigationOrder += menu_wallet + 快照测试（权限 27→30、导航 13→14）。
- **schema/fragment**：wallet 页（账户表 + 创建/调账/冻结/解冻/状态弹窗 + 对账 + 流水导航）+ wallet-entries 页（data.route-binding 流水表）+ fragment.json（页面 + menu_wallet 导航）。
- **web**：en-US/zh-CN 双语键（manifest/schema/error 全量）+ schema-keys 分母清单 + admin fixture（wallet/wallet-entries 页面 + 导航 + SHA 重钉）+ app-manifest 断言。
- **测试**：store（apply 表/原子性/幂等/乐观锁/对账链）+ handler（门禁 401/403、生命周期流、幂等、状态、错误码）+ 迁移冻结（30→32）。
- **验证**：go test ./... **全量全绿**；web vitest **1004/1004**。
- 检查点：S2 完成 → progress 2/5；git checkpoint 提交（hash 见 commit 消息）。
