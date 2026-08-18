---
id: E-002
goal: GOAL-030-w19-my-wallet-lazy-open-empty-state
title: S2 实施惰性开通与空态
status: completed
created: 2026-08-18
updated: 2026-08-18
version: 0.1.0
parent: GOAL-001-design-implementation-conformance
---

# E-002 · S2 实施（2026-08-18）

## 已发生事实

1. 新增 `wallet-ensure`：进页一次 `POST /api/wallet/me`，成功后 `reloadList`；失败显示重试。
2. `my-wallet.json` 挂该组件；去掉表格常驻「开通钱包」。
3. `WALLET_NOT_FOUND` 在 statCard / SchemaTable 按空列表处理，不当硬错误。
4. GET `/api/wallet/me` 仍只读（W15-F11 未回退）。

## 下一步（计划）

S3 记录测试。
