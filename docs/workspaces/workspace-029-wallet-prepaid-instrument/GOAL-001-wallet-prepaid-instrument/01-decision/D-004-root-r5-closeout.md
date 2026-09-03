---
id: GOAL-001-wallet-prepaid-instrument
doc: decision-entry
record_id: D-004
status: accepted
parent: GOAL-001-wallet-prepaid-instrument
created: 2026-09-02
updated: 2026-09-02
version: 0.1.0
---

# D-004 · 根目标 GOAL-001 与工作区 029 结项关门裁决（2026-09-02）

- **status**：accepted
- **deciders**：用户书面指令 + govern 编排器
- **scope**：`workspace-029-wallet-prepaid-instrument` Root `GOAL-001-wallet-prepaid-instrument` 关门（`active → done`）与工作区结项

### 决定

1. **Root 目标正式关门**：`GOAL-001-wallet-prepaid-instrument` 状态由 `active` 变更为 `done`，`progress: 5/5`。
2. **工作区结项**：`workspace-029-wallet-prepaid-instrument` 的 `workspace.md` 状态变更为 `done`。
3. **交付成果确认**：
   - 外部主体接缝：`(issuer, external_id) → subject_id` 幂等登记与查找，不创建 `admin.users`。
   - 钱包预付凭证：批次生成、导出、作废、过期拒绝，哈希存储与一次性出示。
   - 账本入金：单事务 CAS + 原子 Redeem 入金，并发防双花，三余额恒等保持。
   - 管理后台操作面：协议驱动批次管理与声明化导出，`wallet.voucher.issue` 细粒度权限，操作审计脱敏。
   - R5 自助核销增量：`POST /api/wallet/me/redeem` 身份作用域核销，不串 subject 账，限流专用桶，`/my-wallet` 具名卡片入口与即时刷新。
4. **红线保持确认**：未重开 VP-011，未修改 Profile 默认集与 Manifest 装配，未引入外部支付网关或 Telegram SDK，不消耗 RT-Q03/Q05 trigger。
