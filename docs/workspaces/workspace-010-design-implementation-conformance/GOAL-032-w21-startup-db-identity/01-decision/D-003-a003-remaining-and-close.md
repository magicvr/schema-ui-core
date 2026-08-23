---
id: GOAL-032-w21-startup-db-identity
doc: decision-entry
record_id: D-003
status: accepted
parent: GOAL-001-design-implementation-conformance
created: 2026-08-22
updated: 2026-08-22
version: 0.1.0
---

## D-003 · 响应 A-003：确认 required 闭合、处理剩余 recommended、关门

### 触发

用户 `/govern 响应 A-003。确认 F-001～F-003 fixed。适当处理一下其他仍开放的审计意见。然后关门。`

A-003 independent **pass**：A-001 F-001～F-003 **fixed**。本条书面确认该复核。

### 决定

1. **确认**：A-001 F-001 / F-002 / F-003 以 A-003 为关闭权威，状态 **fixed**。不得重开。
2. **A-003 F-001**（recommended）：修测试——catalog 头版本必须与锁相等，且 `lockedHeadExtraTables[head]` 中的表名必须出现在 `completeLostLedgerTables`。
3. **A-003 F-002**（recommended）：扩大 `postV1CatalogTables`（含 wallet/dict/recycle/captcha/tasks/mfa 等），降低「抽样漏表 → 误走 partial」面。
4. **A-001 F-004b**（recommended）：补 sqlite 完整库丢 ledger 的 Open 级 restore 测试。
5. **A-001 F-007**（recommended）：删除未使用的 `kernel.ExecIdempotentDDL`，保留 `IsDuplicateObject`。避免以后接到 Apply 上吞 42P07。
6. **A-001 F-006**（recommended）：**accepted-residual**。Identify 与 v1 Apply 双探针保留为防御深度。范围：仅 `users` 身份列探测重复。复审触发：下一次改 `users` 基线列或再出现 Identify/v1 结论不一致。用户本轮书面要求「适当处理其余意见后关门」，不强制合并两处实现。
7. **关门**：S5 记完成。GOAL-032 `done` 5/5。Root/VP 仍 active。go 不暂挂。

### 为什么

A-003 已独立核对三条 required 的关闭证据。剩余均为 recommended。能用小改降低复发的就改；F-006 合并会把 store 与 authsession 迁模块耦在一起，收益低于本波关门成本。

### 未选方案

- 把 recommended 全部拖到下一波：用户已要求本波关门。
- 合并 `usersLooksLike*` 到单一包：循环依赖风险，本波不做。
