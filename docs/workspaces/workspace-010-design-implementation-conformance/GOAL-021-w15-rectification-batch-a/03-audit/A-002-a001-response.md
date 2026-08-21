---
id: GOAL-021-w15-rectification-batch-a
doc: audit-entry
record_id: A-002
source: self
scope: S4 响应 A-001
verdict: pass
status: recorded
auditor: grok-build /govern
parent: GOAL-020-w15-user-perspective-findings
created: 2026-08-17
updated: 2026-08-17
version: 0.1.0
---

# A-002 · 响应 A-001

- **source**：self
- **类型** / **scope**：response · A-001
- **verdict**：pass

## 关闭证据

| finding | 路径 | 证据 |
|---------|------|------|
| A-001 F-001 required | **fixed** | `authFetch` 仅在 refresh 已清空 refresh token（401/403）时 `clearTokens`+`onAuthLost`；网络/5xx 保留 Token。测试 `authFetch keeps tokens when refresh returns 500 after a 401` |
| A-001 F-002 required | **fixed** | `WithJSONRouteErrors` 对已注册 GET 的 HEAD 交给 mux（Go 默认 200），不再 405。测试 `HEAD on GET route is not JSON 405` |
| A-001 F-003 recommended | **fixed** | YAML 键保持仓内 `cors_origins`（snake_case），与 D-001 驼峰字面差异记录为本仓惯例 |
| A-001 F-004 recommended | **fixed** | 日志路径：`C:\Users\magicvr\AppData\Local\Temp\grok-goal-025f997f6e16\implementer\go-test-1.log` 等四份（E-002） |
| A-001 F-005 recommended | **fixed** | D-001 范围是不清 Token；`restoreSession` 返回 `reauth` 但 Token 仍在，符合冻结 |
| A-001 F-006 recommended | **fixed** | 已补 `authFetch` 500 与 HEAD 测试 |

## 仍开放

无 required。
