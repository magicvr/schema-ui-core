---
id: GOAL-002-subject-module-and-wallet-integration
title: 外部主体接缝与钱包预付凭证入金集成
status: done
parent: GOAL-001-wallet-prepaid-instrument
created: 2026-09-02
updated: 2026-09-02
version: 0.1.0
progress: 4/4
---

# GOAL-002 · 外部主体接缝与钱包预付凭证入金集成

## 概述

承接 Root `GOAL-001` 纲领阶段 R2。依据 D-002 裁决结果，交付通道无关外部主体接缝与预付凭证（卡密）核销入金集成：
1. 通道无关外部主体表 `subjects` 与 `GetOrCreateSubject`/`SubjectExists` 服务；
2. 钱包 `owner_type` 扩充 `subject`，`OwnerExistsFunc` 适配主体校验；
3. 预付凭证哈希存储、高熵生成与单事务 CAS 原子核销入金（`Redeem` API）；
4. 保持三余额恒等与对账不变式，并发双花严格 fail-closed。

## 成功标准

- [x] 1. 迁移 0064 建立 `subjects` 表与 `vouchers` 表，双方言（SQLite / PG）兼容；
- [x] 2. 外部主体接缝可用：`(issuer, external_id) -> subject_id` 幂等 get-or-create，未登记主体不能开户，不创建 `admin.users`；
- [x] 3. 预付凭证核心服务：高熵码生成、SHA-256 哈希存储、一次性出示明文；
- [x] 4. 核销原子且防并发双花：单事务 CAS 保证核销与账本入金强一致，并发双花 fail-closed，三余额恒等式与流水不变式通过回归测试。

## 派生进度展示

`progress: 4/4`（4 个成功检查点已全部完成并通过交叉审计 A-003 闭合）。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-002-001 | non-blocking | 迁移 0064 在 SQLite 中对 `wallet_accounts` owner_type CHECK 的扩展方式（重建 vs 独立迁移） | 实施 | 方案冻结 | 代码设计核验 | closed | — | 统一通过标准迁移 DDL 管理 |

## 父目标

- `GOAL-001-wallet-prepaid-instrument`

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger。
