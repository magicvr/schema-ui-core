---
id: E-004-r5-s2-gate-validation
goal: GOAL-006-r5-maintenance-read-only-gate
status: recorded
created: 2026-08-18
updated: 2026-08-18
parent: GOAL-006-r5-maintenance-read-only-gate
version: 0.1.0
---

# E-004 · R5 S2 统一写门禁与 composition 黑盒验证

## 已核对事实

- `c4856f2` 中的 `WithOperationalGate` 位于 request-id/CORS 之后、JSON 404/405 handler 之前；仅对 mux 已注册的当前 mutation 方法执行 gate，allowlist 精确匹配认证生命周期与强制改密。
- `r5_operational_gate_test.go` 以真实 `admin` profile composition 验证 Provider `POST /api/data-dictionary/types` 在 maintenance/degraded/read-only 下分别返回 `503 SERVICE_*` + correlation；normal 保持既有 `401 UNAUTHENTICATED`。
- 同一黑盒验证确认 login allowlist 保留自身 400、`GET /healthz` 200、未知 POST 404、`POST /healthz` 405。
- handler operational tests、composition targeted test、system-monitoring/config tests 与此前 handler 全量均通过。

## 边界

S3 仍需核对 Host/前端消费、全量回归、独立关门审计与 required/recommended findings；本条不宣称 R5 已完成。
