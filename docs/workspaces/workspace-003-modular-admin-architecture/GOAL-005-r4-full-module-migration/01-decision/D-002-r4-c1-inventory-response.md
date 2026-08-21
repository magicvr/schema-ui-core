---
id: D-002-r4-c1-inventory-response
doc: decision-entry
goal: GOAL-005-r4-full-module-migration
date: 2026-08-05
status: accepted
finding_closure: fixed
---

# D-002 · R4 C1 能力盘点响应

## 决定

接受 `attachments/r4-c1-capability-inventory.md` 作为 C1 的 freeze-grade
事实盘点，并以 `fixed` 响应 self A-001 的 `F-R4-004` 和 Grok A-002 的
`F-GROK-R4-001`。当前已知 first-party page/Schema、10 个 BuiltinModules、
RBAC seed/menu、迁移 ledger、Web Shell/protocol/renderer/delivery surface 均有
owner、阶段处置和证据路径。

## 范围边界

- `admin.users`、`admin.roles` 进入 R4 C3 migrate。
- `admin.settings`、`admin.activity` 进入 R4 C4 migrate；operationlog persistence
  保持 core cross-cutting 候选，不随 Activity UI 关闭。
- overview、示例页、Web Shell、protocol、Renderer、Docker/nginx 属于 keep-core
  或 consumer-only/delivery surface，不冒充一方 Admin module migration 完成。
- Records 明确标为 `pending-gate`，仍由 R4-I003 管理；本决定不选择恢复 CRUD 或
  historical-only 解释，也不改写 VP-003。

## 放行影响

R4-I001 可标为 `verified`，但 C1 仍未完成。R4-I002 provider contract、R4-I003
Records 范围和 R4-I004 operationlog 语义/retention 仍为 collecting/open；本决定
不授权 C2、不推进 Root progress。
