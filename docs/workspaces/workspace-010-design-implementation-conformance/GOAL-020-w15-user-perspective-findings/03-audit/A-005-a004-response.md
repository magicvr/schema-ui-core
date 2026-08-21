---
id: GOAL-020-w15-user-perspective-findings
doc: audit-entry
record_id: A-005
source: self
scope: S 关门响应 A-004
verdict: pass
status: recorded
auditor: grok-build /govern
parent: GOAL-001-design-implementation-conformance
created: 2026-08-17
updated: 2026-08-17
version: 0.1.0
---

# A-005 · 响应 A-004

- **source**：self
- **verdict**：pass

| finding | 路径 | 证据 |
|---------|------|------|
| A-004 F-001 required | **fixed** | `my-wallet.json` 增加 `POST /api/wallet/me` 动作 + toolbar「开通钱包」 |
| A-004 F-002 required | **fixed** | `account.json` 会话表增加 current / userAgent / ip 列 + i18n |
| A-004 F-003 recommended | **fixed** | CORS Allow-Headers 含 `Accept-Language, X-Refresh-Token` |
| A-004 F-004 recommended | **fixed** | create 字段 reason 已是规则码（PATCH 非本波 W15-F08 证据范围） |
| A-004 F-005 recommended | **fixed** | 用户菜单与通知铃铛 ArrowUp/Down 循环焦点 |
| A-004 F-006 recommended | **fixed** | 会话列表默认 `DefaultPageSize` |

无开放 required。
