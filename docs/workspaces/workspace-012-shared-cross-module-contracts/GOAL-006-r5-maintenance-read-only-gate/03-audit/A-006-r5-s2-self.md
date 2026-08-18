---
id: A-006-r5-s2-self
goal: GOAL-006-r5-maintenance-read-only-gate
source: self
verdict: pass
status: recorded
created: 2026-08-18
updated: 2026-08-18
parent: GOAL-006-r5-maintenance-read-only-gate
version: 0.1.0
---

# A-006 · R5 S2 统一写门禁自审

## 审计头

| 项 | 值 |
|----|----|
| source | self |
| scope | S2：operational gate、SERVICE_* catalog、core/Provider composition 黑盒矩阵 |
| verdict | pass |
| required findings | 0 |

## 核对结论

1. Gate 的 mux method probe 保留未知路径/方法不匹配的 `NOT_FOUND` / `METHOD_NOT_ALLOWED`，不会把路由发现错误伪装成运行态拒绝。
2. 三种受控态均 fail closed 业务写；login/refresh/logout/MFA verify/强制改密 allowlist 保留认证生命周期；GET/HEAD 与 health/readiness 不受影响。
3. 新增 `SERVICE_MAINTENANCE`、`SERVICE_DEGRADED`、`SERVICE_READ_ONLY` 均进入 error catalog，错误包络携带 request correlation；read-only 使用 503，避免与现有 423 `ACCOUNT_LOCKED` 客户端分流。
4. 真实 composition Provider route、normal auth、unknown/mismatch/health 黑盒测试通过；Profile/Manifest/readiness 装配语义未改变。

## Findings

| ID | 等级 | finding | disposition |
|----|------|---------|-------------|
| F-001 | recommended | S3 需核对 Web Host 对 degraded/status 和 SERVICE_* API error 的消费边界，不将应用内拒绝误判为 Host unavailable。 | implementation gate：S3 |

## 结论

S2 self 通过，开放 required = 0；S3 进入最终消费、回归与 independent 关门审计。
