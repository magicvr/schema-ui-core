---
doc_type: goal-decision
id: D-001-r3-voucher-api-and-page-design
parent: GOAL-003-admin-voucher-surface-and-export
date: 2026-09-02
status: accepted
version: 0.1.0
---

# D-001 · R3 Admin 批次 API、导出与权限设计

## 触发

推进 Root `GOAL-001` 纲领阶段 R3，交付预付凭证的管理员操作端点与审计集成。

## 决定

1. **HTTP 端点**：
   - `POST /api/wallet/vouchers/batches`：批量生成预付凭证，权限 `wallet.voucher.issue`。接收 `{batchId, count, amount, currency, expiresAt}`，返回一次性明文列表供管理员下载/导出。
   - `GET /api/wallet/vouchers`：列表查询凭证，权限 `wallet.read`。支持 `batchId`、`status`、`page`、`pageSize`，仅返回 `codePrefix`，不返回明文与哈希。
   - `GET /api/wallet/vouchers/{id}`：详情查询，权限 `wallet.read`。
   - `POST /api/wallet/vouchers/{id}/void`：作废单张未核销凭证，权限 `wallet.voucher.issue`。
2. **权限与模块贡献**：
   - 在 `admin.wallet` Module Descriptor 的 Permissions 中注册 `"wallet.voucher.issue"`，并在 RBAC 默认策略中赋予 Admin 角色。
   - 在 Routes 中追加上述凭证路由贡献。
3. **审计策略**：
   - 生成批次与作废操作记录到 `operationlog`。
   - 审计 Payload 仅记录 `batchId`, `count`, `amount`, `currency`, `voucherId`，**绝对不记录**卡密明文。
