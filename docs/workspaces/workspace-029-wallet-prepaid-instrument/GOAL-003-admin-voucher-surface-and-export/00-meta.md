---
id: GOAL-003-admin-voucher-surface-and-export
title: 管理后台预付凭证批次管理、导出与操作审计
status: done
parent: GOAL-001-wallet-prepaid-instrument
created: 2026-09-02
updated: 2026-09-02
version: 0.1.0
progress: 4/4
---

# GOAL-003 · 管理后台预付凭证批次管理、导出与操作审计

## 概述

承接 Root `GOAL-001` 纲领阶段 R3 与 VP-029 判据 #2、#5。交付预付凭证的 Admin 操作面与完整生命周期管理：
1. 预付凭证批次生成 API（`POST /api/wallet/vouchers/batches`，权限 `wallet.voucher.issue`，明文一次性返回）；
2. 预付凭证列表与查询 API（`GET /api/wallet/vouchers`，按 batchId / status 过滤，仅出示 prefix，不出示明文）；
3. 预付凭证作废 API（`POST /api/wallet/vouchers/{id}/void`，权限 `wallet.voucher.issue`）；
4. 操作审计日志留痕（明文卡密绝不进审计日志详情）；
5. 路由与权限字典贡献对齐。

## 成功标准

- [x] 1. Admin 批次生成、列表查询、作废 HTTP API 交付，权限键 `wallet.voucher.issue` 严格鉴权；
- [x] 2. 导出与一次性出示机制：明文仅在生成批次时一次性出示/导出，数据库与审计日志中绝对不含明文；
- [x] 3. 凭证作废与生命周期操作审计：作废成功状态流转，记录操作审计；
- [x] 4. 路由与权限贡献注册到 `admin.wallet` 模块，全量测试与回归全绿。

## 派生进度展示

`progress: 4/4`（4 个成功检查点已全部完成并通过交叉审计 A-003 闭合）。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-003-001 | non-blocking | 凭证生成时导出文件结构（CSV / JSON） | 实施 | 方案冻结 | 设计核验 | closed | — | API 返回明文数组，作为一次性导出数据源（D-001） |

## 父目标

- `GOAL-001-wallet-prepaid-instrument`

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger。
