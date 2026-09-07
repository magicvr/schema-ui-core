---
doc_type: goal-execution
id: E-012-r3-c3-contract-gate
parent: GOAL-004-r3-session-operator-console
date: 2026-09-04
source: self
status: done
version: 0.1.0
---

# E-012 · R3 C3 合同门禁与实现放行（2026-09-04）

## 已发生事实

- 用户通过裁决工具选择 D-010：`telegram.operator.read` 可与 `settings.read`
  并列获取/维持既有 polling lease；未绑定 polling 仍按心跳进入 `running`，
  operator API 不隐式自启。
- D-009 v0.2.0 补全了 A-018 F-001/F-002/F-003 所要求的 composition 认证包装、
  lease 授权和 PostgreSQL 安全冲突处理，并纳入 A-018 F-004～F-007 的非阻断项。
- A-020 为 Grok `grok-4.6 · reasoning high` independent `pass`，确认三项
  required 在合同侧 `fixed`；A-021 已响应。合同 checkpoint 可放行 C3 生产代码。

## 当前边界

R3 仍为 `active · 2/4`；C3 代码、测试、迁移和实现 independent 尚未完成。下一步
只实现 D-009/D-010 明确的 operator API、权限与 v69 outbound，不提前交付 C4 UI、
`getChatMember` 缓存、发言权反馈或全局 R3 关门。
