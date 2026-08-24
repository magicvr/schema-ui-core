---
id: E-001
doc: execution-entry
goal: GOAL-004-r3-binding-flow
status: recorded
parent: GOAL-001-account-email-identity
created: 2026-08-24
updated: 2026-08-24
version: 1.0.0
---

# E-001 · 子目标开题 + R3 方案冻结（2026-08-24）

## 已发生事实

- 用户裁决（会话留痕 i005_ttl / i005_cooldown / i006_admin / scope_ui）：TTL 10 分钟；冷却 60 秒；允许代填待校验；分母含最小页面。
- 勘察事实：MailSender 端口同步单收件人；自助面路由先例 AccountSelfRoutes；错误面 errorcatalog.Catalog + writeError。
- R3 方案七条款冻结（D-001）：迁移 0055 成对 DDL、服务三操作语义、API 合同、最小页面、尝试上限。
- 未改代码。

## 证据

| 主张 | 路径 |
|------|------|
| 方案条款 | 本目标 D-001 |
| 端口契约 | `apps/api/internal/kernel/mail.go` |
