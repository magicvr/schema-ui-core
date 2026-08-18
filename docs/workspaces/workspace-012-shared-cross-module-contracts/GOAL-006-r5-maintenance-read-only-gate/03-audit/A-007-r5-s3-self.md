---
id: A-007-r5-s3-self
goal: GOAL-006-r5-maintenance-read-only-gate
source: self
verdict: pass
status: recorded
created: 2026-08-18
updated: 2026-08-18
parent: GOAL-006-r5-maintenance-read-only-gate
version: 0.1.0
---

# A-007 · R5 S3 消费与回归自审

## 审计头

| 项 | 值 |
|----|----|
| source | self |
| scope | S3 Host/前端消费、全量 API/Web 回归、Profile/Manifest/protocol/readiness 不变式 |
| verdict | pass |
| required findings | 0 |

## 核对结论

1. maintenance 是唯一 Host terminal；degraded/read-only 进入应用并由 `availabilityMode` + API `SERVICE_*` error 表达，符合 A-002 F-001/F-002 closure，不误触发 Host unavailable 或 account-lock terminal。
2. Web Host/bootstrap、resource error localization 定向测试与 production build 通过；API 全量测试通过；conformance claim 生成证据已 checkpoint。
3. 读路径、health/readiness、认证 allowlist 与正常态兼容均由 S2/S3 tests 覆盖；Profile/module/Manifest/protocol 约束无新增差异。

## Findings

| ID | 等级 | finding | disposition |
|----|------|---------|-------------|
| F-001 | recommended | `SERVICE_*` 文案若未来进入 HostFailureScreen，应新增专门 host/UI 文案；当前 API ResourceApiError 与 status 消费已足够，UI 横幅属后续范围。 | deferred non-blocking; outside R5 S3 scope |

## 结论

S3 self 通过，开放 required = 0；F-001 为明确的非阻塞后续建议。请求 A-008 independent 进行最终关门审计。
