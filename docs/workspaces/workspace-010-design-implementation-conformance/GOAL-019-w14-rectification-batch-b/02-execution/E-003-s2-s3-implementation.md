---
id: GOAL-019-w14-rectification-batch-b
doc: execution
status: active
parent: GOAL-015-w14-user-perspective-review
created: 2026-08-17
updated: 2026-08-17
version: 0.1.0
---

# E-003 · S2/S3 实施与回归

## 事实

- **2026-08-17**：F-05 实施——recycle-bin/wallet/per-task runs 分页校验；data-permission policies 内存分页；recycle sort/order。
- **2026-08-17**：F-06 实施——`OPERATION_NOT_FOUND`、`INVALID_SCOPE_BODY`、`INVALID_WALLET_OWNER`、`INVALID_WALLET_ACCOUNT`、`INVALID_WALLET_STATUS` 入目录/契约/i18n。
- **2026-08-17**：F-07 实施——通知 q 大小写不敏感；wallet 账户搜索扩展；wallet ledger `q` + `entryType` 过滤；wallet-entries schema 增加 entryType 筛选；recycle sort/order。
- **2026-08-17**：S3 回归——Go 全量 `go test ./...` 通过；Web 全量 1041/1041、tsc、build 通过。
