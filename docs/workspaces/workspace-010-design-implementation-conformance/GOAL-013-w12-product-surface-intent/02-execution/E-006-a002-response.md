---
id: E-006-a002-response
doc: execution-entry
goal: GOAL-013-w12-product-surface-intent
date: 2026-08-16
status: recorded
parent: GOAL-001-design-implementation-conformance
created: 2026-08-16
updated: 2026-08-16
version: 0.1.0
---

# E-006 · A-002（grok independent）响应与 F-001～F-005 闭合

## 事实

### F-001（required）· 关门台账先于独立意见 —— **fixed**

- **纠正**：`00-meta.md` 撤回预写 `status: done` → `active`（progress 3/4，S4 检查点未勾）；`goal-tree.md` / `workspace.md` W12 行撤回 done → active 3/4；`03-audit.md` 索引按 A-002 真实 verdict（**conditional**，开放 F-001）登记（grok 已落盘）；`E-005` 审计段改为「conditional，响应见 E-006」。
- **过程教训**：独立意见落盘前不得预写 pass / done；先落 A 条目再动状态。
- **闭合路径**：fixed（台账已按本意见真实 verdict 纠正；本 E-006 即响应留痕）。

### F-002（recommended）· Web 全量一次未复现 1027 —— **fixed（复跑复现）**

- `npx vitest run` 全量复跑（2026-08-16，F-003/F-004 修复后）：**63 文件 / 1027/1027 通过**，含 `s5-denominator-render.test.tsx` 5/5。grok 观察到的 1025/1027 为并行负载 flake，非本波回归。

### F-003（recommended）· users 关键词标签含 ID 但 SQL 不搜 id —— **fixed**

- `authsession/users_repository.go` `usersWhere` 关键词子句增加 `instr(lower(u.id), ?) > 0`（与冻结标签「用户名 / 显示名 / ID」一致）；`ListUsers` COUNT 查询表别名 `FROM users u` 同步修正。
- 验证：`go test ./internal/modules/authsession/ ./internal/handler/ -run "Users|Import"` 通过；Go 全量 0 FAIL。

### F-004（recommended）· 窄屏隐藏显示名 —— **fixed**

- `App.tsx` `UserMenu` 触发器显示名去掉 `hidden sm:inline`（`max-w-[8rem] truncate` 始终显示），严格对齐 D-002 §1「头像圆标 + 显示名 + 小箭头；全断点同一控件」。`user-menu.test.tsx` 4/4 通过。

### F-005（recommended）· wallet-entries entryType 降级 —— **accepted-residual（D-003 既定条款）**

- D-003 冻结矩阵明确「若 list 未按类型可滤则仅关键词，不发明类型」；`wallet-entries.json` 仅 q 符合方案。账本 list 增加 entryType 谓词后可在后续波次再挂（与 A-001 R-001 同一观察）。范围：本波不实现；复审触发：账本 list 出现 entryType 参数时。

## 复核证据

- `go test ./... -count=1`（apps/api）：**0 FAIL**（F-003 修复后全量复跑）。
- `npx vitest run`（apps/web）：**1027/1027**（F-004 修复后全量复跑）。
- `tsc -b`：0 错误。
- 状态台账：E-005 + goal-tree / workspace.md 同步；A-002 conditional 的 required F-001 已按 fixed 闭合，S4 检查点可勾选。
