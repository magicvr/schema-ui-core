---
doc_type: goal-audit
id: A-021-r3-c3-a020-response
parent: GOAL-004-r3-session-operator-console
date: 2026-09-04
source: self
auditor: Codex
audit_type: finding-response
scope: 响应 A-020 C3 合同修复独立复审；确认 A-018 required 已闭合并放行 C3 实现
verdict: pass
open_required: 0
version: 0.1.0
---

# A-021 · R3 C3 A-020 独立复审响应（2026-09-04）

## 响应结论

A-020 为 Grok `grok-4.6 · reasoning high` 的合同修复 independent `pass`，确认
A-018 F-001/F-002/F-003 在响应侧以 `fixed` 合法闭合，`open_required: 0`，并允许
进入 C3 生产代码。A-018 原件及 A-017～A-020 原文均保留；本条不把合同通过写成
C3 实现完成。

## C3 实施放行

- F-001：实现必须在 composition 以 `a.Middleware` 包 operator handler；
  `Public: false` 不能替代认证，401 → 403 `FORBIDDEN` → 409 顺序和 service-credential scope
  需用测试核对。
- F-002：实现既有 lease 的 `settings.read OR telegram.operator.read` 授权，
  保留未绑定 polling 按心跳进入 `running` 的门禁；settings API 不变。
- F-003：实现 request/root pending 的 `ON CONFLICT DO NOTHING` + `RowsAffected`
  事务分支、双方言 partial unique，并跑 gated PostgreSQL 冲突路径。
- A-020 recommended F-001～F-003 作为 C3 实现验收项处理：未知 chat/request
  采用稳定 404 catalog 码；partial unique 冲突不丢 `WHERE status = 'pending'`；
  `request_id` 只接受 mux-safe `[A-Za-z0-9._-]{1,128}`。

本条不修改 R3/C3 status 或 progress，不接受 residual，不 overrule；C3 进入生产
实现阶段，C3 检查点仍未完成。实现后必须另做 self + Grok independent 审计。
