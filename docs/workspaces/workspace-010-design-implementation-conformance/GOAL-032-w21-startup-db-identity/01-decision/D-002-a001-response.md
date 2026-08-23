---
id: GOAL-032-w21-startup-db-identity
doc: decision-entry
record_id: D-002
status: accepted
parent: GOAL-001-design-implementation-conformance
created: 2026-08-22
updated: 2026-08-22
version: 0.1.0
---

## D-002 · 响应 A-001 F-001～F-003（收紧 restore / partial / sqlite 合同）

### 触发

A-001 independent close-out **conditional**，开放 required F-001（high）、F-002/F-003（med）。用户 `/govern 响应 GOAL-032 A-001`。

### 决定

1. **F-001 restore-ledger**：`lostLedgerLooksComplete` 不再只看四表。必须同时有 catalog 头对象 `service_credentials`（v44）与 `operation_log_session`（v48）。`completeFingerprintCatalogHead = 48`，catalog 再涨则 `TestCompleteFingerprintTracksCatalogHead` fail closed，逼更新指纹。缺头对象 → **不整表盖章**。
2. **F-002 partial**：无 ledger、我方 `users`、且已有任一 post-v1 catalog 表（`roles` / `operation_log` / `jobs` / …）但未达完整指纹 → 新身份 `lost-ledger-unsafe` → **refuse**（不 CREATE、不 stamp）。`adopt-then-pending` 仅保留「我方 users、且没有 post-v1 catalog 表」（users-only 或等价）。
3. **F-003 sqlite**：不废止 V-MIG-03。sqlite users-only 仍 fail closed、不留 ledger（v1 `fingerprintR2` 精确二表）。postgres users-only 仍可 v1 补 `refresh_tokens` 再 pending。枚举共享；**v1 Apply 方言差**写进合同，不是执行分叉隐瞒。

D-001 其余（ledger 权威、EF 用法、外库 refuse、精确 R2 adopt、健康 ledger noop/pending）仍有效。restore / partial 条款以本条为准。

### 为什么

四表停在 v42，盖章当前 1..48 会跳过 `service_credentials` 等。mid-catalog 丢 ledger 无法安全猜出版本。sqlite 精确 R2 是既有文件库合同，不在本波改掉。

### 未选方案

- 全部 DDL `IF NOT EXISTS` / 吞 42P07：D-001 已拒，且 PG 事务 25P02。
- 改 sqlite v1 放行 users-only：废 V-MIG-03，扩大损坏文件库的自动修复面。
- 无 ledger 一律 refuse：现场全量丢 ledger 的合法恢复会失败。
